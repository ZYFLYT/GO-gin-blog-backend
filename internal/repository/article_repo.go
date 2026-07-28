package repository

import (
	"blog_project/internal/model"
	"time"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

/**
简单查询
*/

func (r *ArticleRepository) CreateArticleRepo(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *ArticleRepository) UpdateArticleRepo(id uint, update map[string]interface{}) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).Updates(update).Error
}

func (r *ArticleRepository) DeleteArticleRepo(id uint) error {
	return r.db.Delete(&model.Article{}, id).Error
}

func (r *ArticleRepository) FindArticleByIDRepo(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Tags").
		First(&article, id).Error
	if err != nil {
		return nil, nil
	}
	return &article, err
}

/**
复杂查询
*/

type ArticleListRequest struct {
	Page       int
	PageSize   int
	UserID     *uint
	CategoryID *uint
	TagID      *uint
	Status     *int
	Keyword    string
	IsPublic   bool
}

func (r *ArticleRepository) ListArticleRepo(req *ArticleListRequest) ([]model.Article, int64, error) {
	var articles []model.Article
	var count int64

	//构建基础查询
	baseQuery := r.db.Model(&model.Article{})

	if req.UserID != nil {
		baseQuery = baseQuery.Where("user_id = ?", *req.UserID)
	}

	if req.CategoryID != nil {
		baseQuery = baseQuery.Where("category_id = ?", *req.CategoryID)
	}

	if req.Status != nil {
		baseQuery = baseQuery.Where("status = ?", *req.Status)
	}

	if req.Keyword != "" {
		baseQuery = baseQuery.Where("title LIKE ? OR summary LIKE ?", req.Keyword, req.Keyword)
	}

	if req.IsPublic {
		now := time.Now()
		baseQuery = baseQuery.Where("status = 2 AND published_at <= ?", now)
	}

	//筛选标签,用JOIN
	var hasTagJoin bool
	if req.TagID != nil {
		baseQuery = baseQuery.Joins("JOIN article_tags ON article_tags.article_id=articles.id").
			Where("article_tags.tag_id = ?", *req.TagID)
		hasTagJoin = true
	}

	//统计总数,分情况
	//复制一份当前所有查询条件，创建一条全新独立查询链路，相互隔离，互不干扰。
	//如果直接 countQuery := baseQuery，两者共用一条查询链；后续修改 countQuery 的条件，会污染原始 baseQuery 和分页查询 dataQuery。
	countQuery := baseQuery.Session(&gorm.Session{})
	if hasTagJoin {
		//有 JOIN 时，使用 Distinct 去重计数
		if err := countQuery.Distinct("articles.id").Count(&count).Error; err != nil {
			return nil, 0, err
		}
	} else {
		// 无 JOIN 时，直接 Count（更快）
		if err := baseQuery.Count(&count).Error; err != nil {
			return nil, 0, err
		}
	}

	//有 JOIN 时必须去重，且要查询全部字段（Distinct("articles.*")）
	dataQuery := baseQuery.Session(&gorm.Session{})
	if hasTagJoin {
		dataQuery.Distinct("articles.*")
	}

	//分页
	offset := (req.Page - 1) * req.PageSize
	err := dataQuery.
		Preload("User").
		Preload("Category").
		Preload("Tags").
		Order("articles.id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&articles).Error

	return articles, count, err
}

func (r *ArticleRepository) IncrementViewCountRepo(id uint) error {
	return r.db.Model(&model.Article{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}
