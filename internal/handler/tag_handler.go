package handler

import (
	"blog_project/internal/service"
	"blog_project/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service *service.TagService
}

func NewTagHandler(service *service.TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) CreateTagHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	var req struct {
		TagName string `json:"tag_name" binding:"required,min=1,max=20"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	tag, err := h.service.CreateTagService(req.TagName)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"tag_name": tag.TagName,
	})
}

func (h *TagHandler) UpdateTagHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, err.Error(), nil)
	}

	var req struct {
		TagName string `json:"tag_name"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	if err := h.service.UpdateTagService(uint(id), req.Status, req.TagName); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"tag_name": req.TagName,
		"status":   req.Status,
	})
}

func (h *TagHandler) GetTagHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	tag, err := h.service.GetTagByIDService(uint(id))
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"tag_name": tag.TagName,
		"status":   tag.Status,
	})
}

func (h *TagHandler) GetTagListHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

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

	tag, total, err := h.service.GetTagListService(page, pageSize, keyword, status)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"list":     tag,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *TagHandler) DeleteTagHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	if err := h.service.DeleteTagService(uint(id)); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"message": "删除成功",
	})
}
