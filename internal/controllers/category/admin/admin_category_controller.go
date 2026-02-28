package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	categoryDto "mqfm-backend/internal/dto/category"
	categoryService "mqfm-backend/internal/services/category/admin"
	"mqfm-backend/internal/utils"

)

type AdminCategoryController struct {
	service *categoryService.AdminCategoryService
}

func NewAdminCategoryController(s *categoryService.AdminCategoryService) *AdminCategoryController {
	return &AdminCategoryController{service: s}
}

func (ctrl *AdminCategoryController) Create(c *gin.Context) {
	var input categoryDto.CreateCategoryRequest

	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	file, _ := c.FormFile("image")

	category, err := ctrl.service.Create(input, file)
	if err != nil {
		utils.Log.Error("Category creation error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	response := categoryDto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Image:     category.Image,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", response)
}

func (ctrl *AdminCategoryController) FindAll(c *gin.Context) {
	categories, err := ctrl.service.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories", err.Error())
		return
	}

	var response []categoryDto.CategoryResponse
	for _, cat := range categories {
		response = append(response, categoryDto.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			Image:     cat.Image,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", response)
}

func (ctrl *AdminCategoryController) FindByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	category, err := ctrl.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found", err.Error())
		return
	}

	response := categoryDto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Image:     category.Image,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", response)
}

func (ctrl *AdminCategoryController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var input categoryDto.UpdateCategoryRequest
	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid update data", err.Error())
		return
	}

	file, _ := c.FormFile("image")

	updatedCategory, err := ctrl.service.Update(uint(id), input, file)
	if err != nil {
		utils.Log.Error("Category update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category", err.Error())
		return
	}

	response := categoryDto.CategoryResponse{
		ID:        updatedCategory.ID,
		Name:      updatedCategory.Name,
		Image:     updatedCategory.Image,
		CreatedAt: updatedCategory.CreatedAt,
		UpdatedAt: updatedCategory.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", response)
}

func (ctrl *AdminCategoryController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	if err := ctrl.service.Delete(uint(id)); err != nil {
		utils.Log.Error("Category deletion error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}

func (ctrl *AdminCategoryController) Search(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Search keyword is required", nil)
		return
	}

	categories, err := ctrl.service.Search(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search categories", err.Error())
		return
	}

	var response []categoryDto.CategoryResponse
	for _, cat := range categories {
		response = append(response, categoryDto.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			Image:     cat.Image,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories found successfully", response)
}