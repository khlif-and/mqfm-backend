package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	categoryModel "mqfm-backend/internal/models/category/admin"
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
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	category := categoryModel.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := ctrl.service.Create(&category); err != nil {
		utils.Log.Error("Category creation error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

func (ctrl *AdminCategoryController) FindAll(c *gin.Context) {
	categories, err := ctrl.service.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
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

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

func (ctrl *AdminCategoryController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid update data", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	updatedCategory, err := ctrl.service.Update(uint(id), updates)
	if err != nil {
		utils.Log.Error("Category update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", updatedCategory)
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