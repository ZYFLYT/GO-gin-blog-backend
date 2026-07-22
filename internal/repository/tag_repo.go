package repository

import (
	"blog_project/internal/model"
	"errors"

	"gorm.io/gorm"
)

type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) CreateTagRepo(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepo) UpdateTagRepo(id uint, update map[string]interface{}) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).Updates(update).Error
}

func (r *TagRepo) DeleteTagRepo(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

func (r *TagRepo) FindTagByIDRepo(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, err
}

func (r *TagRepo) FindTagByNameRepo(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, nil
}

func (r *TagRepo) ListTag(page, pageSize int, keyword string, status *int) ([]model.Tag, int64, error) {
	var tagList []model.Tag
	var count int64

	query := r.db.Model(&model.Tag{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&tagList).Error
	return tagList, count, err
}
