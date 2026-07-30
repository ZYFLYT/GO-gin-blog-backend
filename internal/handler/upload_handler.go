package handler

import (
	"blog_project/pkg/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFileHandler 通用文件上传接口
// 前端表单字段名必须叫 "file"
// 支持通过查询参数指定上传类型：?type=avatar / ?type=cover / ?type=others
func UploadFileHandler(c *gin.Context) {
	appG := &utils.Gin{C: c}

	// 1. 获取上传类型（决定存到哪个子目录）
	uploadType := c.DefaultQuery("type", "others")
	allowedTypes := map[string]string{
		"avatar": "avatars",
		"cover":  "covers",
		"others": "others",
	}
	subDir, ok := allowedTypes[uploadType]
	if !ok {
		appG.Error(utils.INVALID_PARAMS, "不支持的上传类型", nil)
		return
	}

	// 2. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, "未找到上传文件", nil)
		return
	}

	// 3. 校验文件大小（限制为 5MB）
	if file.Size > 5<<20 { // 5 * 1024 * 1024
		appG.Error(utils.INVALID_PARAMS, "文件大小不能超过 5MB", nil)
		return
	}

	// 4. 校验文件类型（仅允许图片）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		appG.Error(utils.INVALID_PARAMS, "仅支持 JPG、PNG、GIF、WEBP 格式", nil)
		return
	}

	// 5. 生成唯一文件名
	timestamp := time.Now().Format("20060102")
	randomStr := utils.RandomString(8)
	newFileName := timestamp + "_" + randomStr + ext

	// 6. 保存路径：./uploads/{subDir}/年/月/日/
	saveDir := "./uploads/" + subDir + "/" + time.Now().Format("2006/01/02")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		appG.Error(utils.DB_ERROR, "创建目录失败", nil)
		return
	}
	savePath := filepath.Join(saveDir, newFileName)

	// 7. 保存文件到本地磁盘
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		appG.Error(utils.DB_ERROR, "文件保存失败", nil)
		return
	}

	// 8. 返回可访问的 URL（统一使用斜杠，兼容 Windows）
	url := "/" + filepath.ToSlash(savePath)
	appG.Success(gin.H{
		"url":  url,
		"name": file.Filename,
		"size": file.Size,
		"type": uploadType,
	})
}
