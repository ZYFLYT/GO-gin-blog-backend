package service

import (
	"blog_project/internal/model"
	"blog_project/internal/repository"
	"blog_project/pkg/utils"
	"errors"
)

// UserService UserService业务逻辑
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 构造函数：注入Repository
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterService Register注册
func (s *UserService) RegisterService(username, password, email string) (*model.User, error) {
	existUser, _ := s.repo.FindUserByUsernameRepo(username)
	if existUser != nil {
		return nil, errors.New("该用户名已被使用")
	}

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username: username,
		Password: string(hashPassword),
		Email:    email,
		Nickname: username,
		Status:   1,
	}

	if err := s.repo.CreateUserRepo(user); err != nil {
		return nil, errors.New("注册失败，请稍后重试")
	}
	return user, nil
}

// LoginService Login登录
func (s *UserService) LoginService(username, password string) (*model.User, error) {
	user, err := s.repo.FindUserByUsernameRepo(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status == 0 {
		return nil, errors.New("该账户已被禁用")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

// GetUserByIDService 根据用户ID获取用户信息
func (s *UserService) GetUserByIDService(id uint) (*model.User, error) {
	user, err := s.repo.FindUserByIDRepo(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// UpdateUserService 更新昵称或者邮箱，基于ID
func (s *UserService) UpdateUserService(id uint, nickname, email string) error {
	updates := map[string]interface{}{
		"nickname": nickname,
		"email":    email,
	}
	return s.repo.UpdateUserRepo(id, updates)
}

// DeleteUserService 删除用户，基于ID
func (s *UserService) DeleteUserService(id uint) error {
	return s.repo.DeleteUserRepo(id)
}

// ChangePasswordService 修改密码
func (s *UserService) ChangePasswordService(id uint, oldPassword, newPassword string) error {
	user, err := s.repo.FindUserByIDRepo(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	updates := map[string]interface{}{
		"password": string(hashedPassword),
	}
	return s.repo.UpdateUserRepo(id, updates)
}

// GetUserListService 获取用户列表
func (s *UserService) GetUserListService(page, pageSize int, keyword string, status *int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.ListUserRepo(page, pageSize, keyword, status)
}
