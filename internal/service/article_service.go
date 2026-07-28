package service

import (
	"blog_project/internal/model"
	"blog_project/internal/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	categoryRepo *repository.CategoryRepo
	tagRepo      *repository.TagRepo
	db           *gorm.DB
}

func NewArticleService(
	articleRepo *repository.ArticleRepository,
	categoryRepo *repository.CategoryRepo,
	tagRepo *repository.TagRepo,
	db *gorm.DB) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		db:           db,
	}
}

// CreateArticleReq 创建文章,我们在创建文章中并不会调用repo中的创建,直接把逻辑写在了service中
type CreateArticleReq struct {
	Title       string
	Summary     string
	Content     string
	Cover       string
	UserID      uint
	CategoryID  uint
	TagID       []uint
	Status      int
	PublishedAt *time.Time
}

func (s *ArticleService) CreateArticleService(req *CreateArticleReq) (*model.Article, error) {
	//校验分类是否存在
	if req.CategoryID > 0 {
		category, err := s.categoryRepo.FindCategoryByIDRepo(req.CategoryID)
		if err != nil || category == nil {
			return nil, errors.New("该分类不存在")
		}
	}

	//校验标签是否存在
	if len(req.TagID) > 0 {
		for _, tagID := range req.TagID {
			tag, err := s.tagRepo.FindTagByIDRepo(tagID)
			if err != nil || tag == nil {
				return nil, errors.New("标签 ID " + string(rune(tagID)) + "不存在")
			}
		}
	}

	//构建文章对象
	article := &model.Article{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		Cover:       req.Cover,
		UserID:      req.UserID,
		CategoryID:  req.CategoryID,
		Status:      req.Status,
		PublishedAt: req.PublishedAt,
		ViewCount:   0,
	}

	//核心事物,创建文章＋关联标签
	err := s.db.Transaction(func(tx *gorm.DB) error {
		//第一步,创建文章
		if err := tx.Create(&article).Error; err != nil {
			return err
		}

		//核心,关联标签
		if len(req.TagID) > 0 {
			var tags []*model.Tag

			//执行查询，查询结果全部放到 tags 切片变量
			if err := tx.Where("id IN ?", req.TagID).Find(&tags).Error; err != nil {
				return err
			}
			//.Association("Tags")告诉 GORM，我现在操作文章和 Tag 之间的多对多关系
			//Replace = 替换全部关联关系
			//把当前文章所有旧的标签关联全部删除，然后把传入 tags 切片里面所有标签建立关联，自动操作中间表
			if err := tx.Model(article).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return article, err
}

// UpdateArticleReq 更新数据,基本操作和创建差不多
type UpdateArticleReq struct {
	ID         uint
	UserID     uint // 用于权限校验
	Title      string
	Summary    string
	Content    string
	Cover      string
	CategoryID uint
	TagID      []uint
	Status     int
}

func (s *ArticleService) UpdateArticleService(req *UpdateArticleReq) error {
	// 检查文章是否存在
	exist, err1 := s.articleRepo.FindArticleByIDRepo(req.ID)
	if err1 != nil || exist == nil {
		return errors.New("文章不存在")
	}

	// 校验权限
	if exist.UserID != req.UserID {
		return errors.New("无权修改该文章")
	}

	// 校验新分类
	if req.CategoryID > 0 {
		category, err := s.categoryRepo.FindCategoryByIDRepo(req.CategoryID)
		if err != nil || category == nil {
			return errors.New("该分类不存在")
		}
	}

	// 更新事务(核心)
	errResult := s.db.Transaction(func(tx *gorm.DB) error {
		// 更新文章基础信息
		updates := map[string]interface{}{
			"title":       req.Title,
			"summary":     req.Summary,
			"content":     req.Content,
			"cover":       req.Cover,
			"category_id": req.CategoryID,
			"status":      req.Status,
		}
		if err := tx.Model(&model.Article{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}

		// 重新关联标签
		if len(req.TagID) > 0 {
			var tags []*model.Tag

			if err := tx.Where("id IN ?", req.TagID).Find(&tags).Error; err != nil {
				return err
			}

			if err := tx.Model(exist).Association("Tags").Replace(tags); err != nil {
				return err
			}
		} else {
			// 如果 TagIDs 为空，清除所有标签关联
			if err := tx.Model(exist).Association("Tags").Clear(); err != nil {
				return err
			}
		}

		return nil
	})

	return errResult
}

// DeleteArticleService 删除文章
func (s *ArticleService) DeleteArticleService(id, UserID uint) error {
	// 校验文章存在性
	article, err := s.articleRepo.FindArticleByIDRepo(id)
	if err != nil || article == nil {
		return errors.New("文章不存在")
	}

	// 校验权限
	if article.UserID != UserID {
		return errors.New("无权限删除")
	}

	return s.articleRepo.DeleteArticleRepo(id)
}

// GetArticleService 查询文章
func (s *ArticleService) GetArticleService(id uint, isPublic bool) (*model.Article, error) {
	article, err := s.articleRepo.FindArticleByIDRepo(id)
	if err != nil || article == nil {
		return nil, errors.New("文章不存在")
	}

	if isPublic {
		if article.Status != 2 {
			return nil, errors.New("文章未发表")
		}
		if article.PublishedAt != nil && article.PublishedAt.After(time.Now()) {
			return nil, errors.New("文章未到发表时间")
		}

		// 浏览量加一
		err := s.articleRepo.DeleteArticleRepo(id)
		if err != nil {
			return nil, err
		}
	}

	return article, nil
}

// GetArticleListService 查询文章列表
func (s *ArticleService) GetArticleListService(
	page int,
	pageSize int,
	categoryID *uint,
	tagID *uint,
	userID *uint,
	status *int,
	keyword string,
	isPublic bool,
) ([]model.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	req := &repository.ArticleListRequest{
		Page:       page,
		PageSize:   pageSize,
		UserID:     userID,
		CategoryID: categoryID,
		TagID:      tagID,
		Status:     status,
		Keyword:    keyword,
		IsPublic:   isPublic,
	}

	return s.articleRepo.ListArticleRepo(req)
}
