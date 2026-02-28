package admin

import (
	"errors"
	"mime/multipart"

	categoryDto "mqfm-backend/internal/dto/category"
	categoryModel "mqfm-backend/internal/models/category/admin"
	categoryRepo "mqfm-backend/internal/repositories/category/admin"
	"mqfm-backend/internal/utils"
)

type AdminCategoryService struct {
	repo categoryRepo.CategoryRepository
}

func NewAdminCategoryService(repo categoryRepo.CategoryRepository) *AdminCategoryService {
	return &AdminCategoryService{repo: repo}
}

func (s *AdminCategoryService) Create(req categoryDto.CreateCategoryRequest, file *multipart.FileHeader) (*categoryModel.Category, error) {
	var imagePath string
	if file != nil {
		filename := utils.GenerateUniqueFilename(file.Filename)
		path := "uploads/categories/" + filename
		if err := utils.SaveUploadedFile(file, path); err != nil {
			utils.Log.Error("Failed to save category image: " + err.Error())
		} else {
			imagePath = path
		}
	}

	category := categoryModel.Category{
		Name:  req.Name,
		Image: imagePath,
	}

	if err := s.repo.Create(&category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *AdminCategoryService) FindAll() ([]categoryModel.Category, error) {
	return s.repo.FindAll()
}

func (s *AdminCategoryService) FindByID(id uint) (*categoryModel.Category, error) {
	return s.repo.FindByID(id)
}

func (s *AdminCategoryService) Update(id uint, req categoryDto.UpdateCategoryRequest, file *multipart.FileHeader) (*categoryModel.Category, error) {
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}

	if file != nil {
		filename := utils.GenerateUniqueFilename(file.Filename)
		path := "uploads/categories/" + filename
		if err := utils.SaveUploadedFile(file, path); err != nil {
			utils.Log.Error("Failed to save category image: " + err.Error())
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

func (s *AdminCategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *AdminCategoryService) Search(query string) ([]categoryModel.Category, error) {
	return s.repo.Search(query)
}