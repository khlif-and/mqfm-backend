package service

import (
	"errors"
	"mime/multipart"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/helper"
	"mqfm-backend/internal/shared/logger"
)

type categoryService struct {
	repo port.CategoryRepository
}

func NewCategoryService(repo port.CategoryRepository) port.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(req request.CreateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error) {
	var imagePath string
	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/categories/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save category image")
		} else {
			imagePath = path
		}
	}

	category := entity.Category{
		Name:  req.Name,
		Image: imagePath,
	}

	if err := s.repo.Create(&category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *categoryService) FindAll() ([]entity.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) FindByID(id uint) (*entity.Category, error) {
	return s.repo.FindByID(id)
}

func (s *categoryService) Update(id uint, req request.UpdateCategoryRequest, file *multipart.FileHeader) (*entity.Category, error) {
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}

	if file != nil {
		filename := helper.GenerateUniqueFilename(file.Filename)
		path := "uploads/categories/" + filename
		if err := helper.SaveUploadedFile(file, path); err != nil {
			logger.Error("failed to save category image")
		} else {
			updates["image"] = path
		}
	}

	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *categoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *categoryService) Search(query string) ([]entity.Category, error) {
	return s.repo.Search(query)
}
