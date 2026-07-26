package handler

import (
	"blog_project/internal/service"
	"blog_project/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户控制器
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler 构造函数：注入Service
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register 用户注册接口
func (h *UserHandler) RegisterHandler(c *gin.Context) {
	appG := &utils.Gin{C: c}

	var req struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Password string `json:"password" binding:"required,min=6,max=20"`
		Email    string `json:"email" binding:"required,email"`
	}

	//绑定并校验JSON参数
	if err := c.ShouldBindJSON(&req); err != nil {
		//校验失败，调用Error方法返回错误JSON
		appG.Error(utils.INVALID_PARAMS, "参数错误:"+err.Error(), nil)
		return
	}

	//调用Service层执行业务逻辑
	user, err := h.service.RegisterService(req.Username, req.Password, req.Email)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	//成功返回JSON
	appG.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
	})
}

// Login 登录接口
func (h *UserHandler) LoginHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	//ShouldBindJSON绑定并验证json数据，这是直接获取前端传来的请求体中的json数据，就不用手动一个一个c.GetPostForm获取数据了
	if err := c.ShouldBindJSON(&req); err != nil {
		//校验失败
		appG.Error(utils.INVALID_PARAMS, "参数错误", nil)
		return
	}

	user, err := h.service.LoginService(req.Username, req.Password)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
	})
}

// GetUser 根据ID获取用户详情
func (h *UserHandler) GetUserHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	//c.Param(key) 是 Gin 框架专门读取URL 路径占位参数的方法
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, "无效的用户ID", nil)
		return
	}

	user, err := h.service.GetUserByIDService(uint(id))
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
		"status":   user.Status,
	})
}

// GetUserList 获取用户列表
func (h *UserHandler) GetUserListHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	// 解析分页参数
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err1 != nil {
		appG.Error(utils.ERROR, err1.Error(), nil)
		return
	} else if err2 != nil {
		appG.Error(utils.ERROR, err2.Error(), nil)
		return
	}
	keyword := c.Query("keyword")

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, err3 := strconv.Atoi(statusStr)
		if err3 != nil {
			appG.Error(utils.ERROR, err3.Error(), nil)
			return
		}
		status = &s
	}

	//调用Service
	users, total, err := h.service.GetUserListService(page, pageSize, keyword, status)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
	}

	appG.Success(gin.H{
		"data":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update 更新数据
func (h *UserHandler) UpdateHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, "无效的用户ID", nil)
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Error(utils.INVALID_PARAMS, err.Error(), nil)
		return
	}

	if err := h.service.UpdateUserService(uint(id), req.Nickname, req.Email); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{"message": "数据更新成功"})
}

// Delete 删除
func (h *UserHandler) DeleteHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, "无效的用户ID", nil)
		return
	}

	if err := h.service.DeleteUserService(uint(id)); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{"message": "用户删除成功"})
}
