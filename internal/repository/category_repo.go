package repository

import (
	"blog_project/internal/model"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) CreateRepo(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepo) UpdateRepo(id uint, update map[string]interface{}) error {
	return r.db.Model(&model.Category{}).Where("id=?", id).Updates(update).Error
}

func (r *CategoryRepo) DeleteRepo(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *CategoryRepo) FindByIDRepo(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *CategoryRepo) FindByNameRepo(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name=?", name).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *CategoryRepo) ListRepo(page, pageSize int, keyword string, status *int) ([]model.Category, int64, error) {
	var category []model.Category
	var count int64

	query := r.db.Model(&model.Category{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if status != nil {
		query = query.Where("status=?", *status)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&category).Error
	return category, count, err
}
