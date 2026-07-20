package service

import (
	"blog_project/internal/model"
	"blog_project/internal/repository"
	"errors"
)

// CategoryService CategoryService业务逻辑
type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategoryService(categoryName, description string, sort int) (*model.Category, error) {
	existName, _ := s.repo.FindByNameRepo(categoryName)
	if existName != nil {
		return nil, errors.New("该类别已存在")
	}

	category := &model.Category{
		Name:        categoryName,
		Description: description,
		Sort:        sort,
	}

	if err := s.repo.CreateRepo(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) UpdateCategoryService(id uint, sort, status int, categoryName, description string) error {
	exist, _ := s.repo.FindByIDRepo(id)
	if exist == nil {
		return errors.New("该分类不存在")
	}

	sameName, _ := s.repo.FindByNameRepo(categoryName)
	if sameName != nil {
		return errors.New("分类名已存在")
	}

	category := map[string]interface{}{
		"name":        categoryName,
		"description": description,
		"sort":        sort,
		"status":      status,
	}
	return s.repo.UpdateRepo(id, category)
}

func (s *CategoryService) GetCategoryListService(page, pageSize int, keyword string, status *int) ([]model.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return s.repo.ListRepo(page, pageSize, keyword, status)
}

func (s *CategoryService) GetCategoryByIDService(id uint) (*model.Category, error) {
	category, err := s.repo.FindByIDRepo(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("该分类不存在")
	}

	return category, nil
}

func (s *CategoryService) DeleteCategoryService(id uint) error {
	return s.repo.DeleteRepo(id)
}
