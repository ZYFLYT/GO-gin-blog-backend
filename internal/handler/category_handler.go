package handler

import (
	"blog_project/internal/service"
	"blog_project/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CreateCategoryHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=20"`
		Description string `json:"description" binding:"min=1,max=50"`
		Sort        int    `json:"sort" binding:"default=0"`
	}

	//绑定并且进行数据校验
	if err := c.ShouldBindJSON(&req); err != nil {
		//校验失败，返回ERROR
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	//调用Service层执行业务逻辑
	category, err := h.service.CreateCategoryService(req.Name, req.Description, req.Sort)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	//成功了返回JSON
	appG.Success(gin.H{
		"category_name": category.Name,
		"category_desc": category.Description,
		"category_sort": category.Sort,
	})
}

func (h *CategoryHandler) GetCategoryHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, err.Error(), nil)
		return
	}

	cate, err := h.service.GetCategoryByIDService(uint(id))
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"category_name": cate.Name,
		"category_desc": cate.Description,
		"category_sort": cate.Sort,
		"status":        cate.Status,
	})
}

func (h *CategoryHandler) GetCategoryListHandler(c *gin.Context) {
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

	category, total, err := h.service.GetCategoryListService(page, pageSize, keyword, status)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"data":     category,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, err.Error(), nil)
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	if err := h.service.UpdateCategoryService(uint(id), req.Sort, req.Status, req.Name, req.Description); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}
	appG.Success(gin.H{
		"category_name": req.Name,
		"category_desc": req.Description,
		"category_sort": req.Sort,
		"status":        req.Status,
	})
}

func (h *CategoryHandler) DeleteCategoryHandler(c *gin.Context) {
	appG := utils.Gin{C: c}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, err.Error(), nil)
		return
	}

	if err := h.service.DeleteCategoryService(uint(id)); err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{"message": "分类删除成功"})
}
