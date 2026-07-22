package repository

import (
	"blog_project/internal/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserRepo(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindUserByIDRepo(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) UpdateUserRepo(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id=?", id).Updates(updates).Error
}

func (r *UserRepository) DeleteUserRepo(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) FindUserByUsernameRepo(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username=?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// 分页查询
func (r *UserRepository) ListUserRepo(page int, pageSize int, keyword string, status *int) ([]model.User, int64, error) {
	var userList []model.User
	var total int64

	query := r.db.Model(&model.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != nil {
		query = query.Where("status=?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//分页逻辑
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&userList).Error
	return userList, total, err
}
