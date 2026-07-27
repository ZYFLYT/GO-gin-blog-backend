package service

import (
	"blog_project/internal/model"
	"blog_project/internal/repository"
	"errors"
)

type TagService struct {
	repo *repository.TagRepo
}

func NewTagService(repo *repository.TagRepo) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTagService(tagName string) (*model.Tag, error) {
	existName, _ := s.repo.FindTagByNameRepo(tagName)
	if existName != nil {
		return nil, errors.New("该标签已存在")
	}

	tag := &model.Tag{
		TagName: tagName,
	}

	err := s.repo.CreateTagRepo(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) UpdateTagService(id uint, status int, newTagName string) error {
	exist, _ := s.repo.FindTagByIDRepo(id)
	if exist == nil {
		return errors.New("该标签不存在")
	}

	existName, _ := s.repo.FindTagByNameRepo(newTagName)
	if existName != nil {
		return errors.New("该标签名已存在")
	}

	tag := map[string]interface{}{
		"tag_name": newTagName,
		"status":   status,
	}

	err := s.repo.UpdateTagRepo(id, tag)
	if err != nil {
		return err
	}

	return nil
}

func (s *TagService) DeleteTagService(id uint) error {
	exist, _ := s.repo.FindTagByIDRepo(id)
	if exist == nil {
		return errors.New("该标签不存在")
	}

	err := s.repo.DeleteTagRepo(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *TagService) GetTagByIDService(id uint) (*model.Tag, error) {
	tag, err := s.repo.FindTagByIDRepo(id)
	if err != nil {
		return nil, err
	}

	if tag == nil {
		return nil, errors.New("该标签不存在")
	}

	return tag, nil
}

func (s *TagService) GetTagListService(page, pageSize int, keyword string, status *int) ([]model.Tag, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.ListTagRepo(page, pageSize, keyword, status)
}
