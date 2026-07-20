package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ================= 1. 错误码定义（原 msg 包内容） =================
const (
	SUCCESS = 0
	ERROR   = 1000

	INVALID_PARAMS = 1001
	DB_ERROR       = 1002
	NOT_FOUND      = 1003
	UNAUTHORIZED   = 1004
)

// MsgFlags 错误码映射表
var MsgFlags = map[int]string{
	SUCCESS:        "操作成功",
	ERROR:          "操作失败",
	INVALID_PARAMS: "请求参数无效",
	DB_ERROR:       "数据库操作异常",
	NOT_FOUND:      "数据不存在",
	UNAUTHORIZED:   "未授权访问",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return "未知错误"
}

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  GetMsg(errCode),
		"data": data,
	})
}

func (g *Gin) Success(data interface{}) {
	g.C.JSON(http.StatusOK, gin.H{
		"code": SUCCESS,
		"msg":  GetMsg(SUCCESS),
		"data": data,
	})
}

func (g *Gin) Error(errCode int, message string, data interface{}) {
	if message == "" {
		message = GetMsg(errCode)
	}
	g.C.JSON(http.StatusOK, gin.H{
		"code": errCode,
		"msg":  message,
		"data": data,
	})
}
