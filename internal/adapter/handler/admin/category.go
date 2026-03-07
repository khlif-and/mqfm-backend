package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	"mqfm-backend/internal/shared/logger"
	resp "mqfm-backend/internal/shared/response"
)

type CategoryHandler struct {
	service port.CategoryService
}

func NewCategoryHandler(s port.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var input request.CreateCategoryRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("image")

	category, err := h.service.Create(input, file)
	if err != nil {
		logger.Error("category create: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgCategoryCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgCategoryCreateOK, toCategoryResponse(category))
}

func (h *CategoryHandler) FindAll(c *gin.Context) {
	categories, err := h.service.FindAll()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgCategoryListFail, err.Error())
		return
	}

	var result []response.CategoryResponse
	for _, cat := range categories {
		result = append(result, toCategoryResponseVal(cat))
	}

	resp.Success(c, http.StatusOK, constant.MsgCategoryListOK, result)
}

func (h *CategoryHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	category, err := h.service.FindByID(id)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgCategoryNotFound, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgCategoryGetOK, toCategoryResponse(category))
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var input request.UpdateCategoryRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("image")

	category, err := h.service.Update(id, input, file)
	if err != nil {
		logger.Error("category update: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgCategoryUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgCategoryUpdateOK, toCategoryResponse(category))
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		logger.Error("category delete: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgCategoryDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgCategoryDeleteOK, nil)
}

func (h *CategoryHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgSearchRequired, nil)
		return
	}

	categories, err := h.service.Search(query)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgCategorySearchFail, err.Error())
		return
	}

	var result []response.CategoryResponse
	for _, cat := range categories {
		result = append(result, toCategoryResponseVal(cat))
	}

	resp.Success(c, http.StatusOK, constant.MsgCategorySearchOK, result)
}
