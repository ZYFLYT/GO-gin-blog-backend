package handler

import (
	"blog_project/internal/service"
	"blog_project/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(service *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

// CreateArticleRequest 创建文章
type CreateArticleRequest struct {
	Title       string     `json:"title" binding:"required"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content" binding:"required"`
	Cover       string     `json:"cover"`
	CategoryID  uint       `json:"category_id"`
	TagIDs      []uint     `json:"tag_ids"`
	Status      int        `json:"status" binding:"oneof=1 2"`
	PublishedAt *time.Time `json:"published_at"`
}

func (h *ArticleHandler) CreateArticleHandler(c *gin.Context) {
	appG := &utils.Gin{C: c}

	userID, exists := c.Get("user_id")
	if !exists {
		appG.Error(utils.UNAUTHORIZED, "未登录或Token无效", nil)
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Error(utils.INVALID_PARAMS, "获取参数错误", nil)
		return
	}

	article, err := h.service.CreateArticleService(&service.CreateArticleReq{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		Cover:       req.Cover,
		UserID:      userID.(uint),
		CategoryID:  req.CategoryID,
		TagID:       req.TagIDs,
		Status:      req.Status,
		PublishedAt: req.PublishedAt,
	})
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(article)
}

// GetArticleListHandler 获取文章列表
func (h *ArticleHandler) GetArticleListHandler(c *gin.Context) {
	appG := &utils.Gin{C: c}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var filterUserID, categoryID, tagID *uint
	var status *int
	if v := c.Query("user_id"); v != "" {
		id64, _ := strconv.ParseUint(v, 10, 32)
		id := uint(id64)
		filterUserID = &id
	}
	if v := c.Query("category_id"); v != "" {
		id64, _ := strconv.ParseUint(v, 10, 32)
		id := uint(id64)
		categoryID = &id
	}
	if v := c.Query("tag_id"); v != "" {
		id64, _ := strconv.ParseUint(v, 10, 32)
		id := uint(id64)
		tagID = &id
	}
	if v := c.Query("status"); v != "" {
		id64, _ := strconv.ParseInt(v, 10, 32)
		id := int(id64)
		status = &id
	}

	keyword := c.Query("keyword")

	// 用户登陆了之后,验证身份,能够访问自己的草稿;其余的只能访问已发表的文章
	var isPublic = true
	currentUserID, isLoggedIn := c.Get("user_id")

	if isLoggedIn && filterUserID != nil {
		if *filterUserID == currentUserID.(uint) {
			// 查看自己的全部文章
			isPublic = false
		}
		// 如果登录了，但 filterUserID 不是自己（比如管理员看张三的文章），依然强制 isPublic=true
	}

	article, count, err := h.service.GetArticleListService(
		page, pageSize, categoryID, tagID, filterUserID, status, keyword, isPublic,
	)
	if err != nil {
		appG.Error(utils.ERROR, err.Error(), nil)
		return
	}

	appG.Success(gin.H{
		"article":   article,
		"count":     count,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetArticleHandler 获取文章详情
func (h *ArticleHandler) GetArticleHandler(c *gin.Context) {
	appG := &utils.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, "无效的文章ID", nil)
	}

	currentUserID := uint(0)
	if uind, exist := c.Get("user_id"); exist {
		currentUserID = uind.(uint)
	}

	article, err := h.service.GetArticleService(uint(id), currentUserID)
	if err != nil {
		appG.Error(utils.NOT_FOUND, err.Error(), nil)
	}

	appG.Success(article)
}

type UpdateArticleRequest struct {
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	Cover      string `json:"cover"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Status     int    `json:"status" binding:"oneof=1 2"`
}

func (h *ArticleHandler) UpdateArticleHandler(c *gin.Context) {
	appG := &utils.Gin{C: c}

	userID, exists := c.Get("user_id")
	if !exists {
		appG.Error(utils.UNAUTHORIZED, "未登录", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, "无效的文章ID", nil)
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Error(utils.INVALID_PARAMS, "参数错误", nil)
		return
	}

	errResult := h.service.UpdateArticleService(&service.UpdateArticleReq{
		ID:         uint(id),
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		TagID:      req.TagIDs,
		Status:     req.Status,
		UserID:     userID.(uint),
	})
	if errResult != nil {
		appG.Error(utils.ERROR, errResult.Error(), nil)
		return
	}

	appG.Success(gin.H{"message": "更新成功"})
}

func (h *ArticleHandler) DeleteArticleHandler(c *gin.Context) {
	appG := &utils.Gin{C: c}

	userID, exists := c.Get("user_id")
	if !exists {
		appG.Error(utils.UNAUTHORIZED, "未登录", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		appG.Error(utils.INVALID_PARAMS, "无效的文章ID", nil)
		return
	}

	errResult := h.service.DeleteArticleService(uint(id), userID.(uint))
	if errResult != nil {
		appG.Error(utils.ERROR, errResult.Error(), nil)
		return
	}

	appG.Success(gin.H{"message": "删除成功"})
}
