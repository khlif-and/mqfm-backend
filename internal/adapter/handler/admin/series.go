package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
	"mqfm-backend/internal/shared/dto/response"
	"mqfm-backend/internal/shared/helper"
	resp "mqfm-backend/internal/shared/response"
)

type SeriesHandler struct {
	service port.AudioSeriesService
}

func NewSeriesHandler(s port.AudioSeriesService) *SeriesHandler {
	return &SeriesHandler{service: s}
}

func (h *SeriesHandler) Create(c *gin.Context) {
	var input request.CreateSeriesRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("image")

	series, err := h.service.Create(input, file)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgSeriesCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgSeriesCreateOK, toSeriesResponse(series))
}

func (h *SeriesHandler) FindAll(c *gin.Context) {
	series, err := h.service.FindAll()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgSeriesListFail, err.Error())
		return
	}

	var result []response.SeriesResponse
	for _, s := range series {
		result = append(result, toSeriesResponseVal(s))
	}

	resp.Success(c, http.StatusOK, constant.MsgSeriesListOK, result)
}

func (h *SeriesHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	series, err := h.service.FindByID(id)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgSeriesNotFound, nil)
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgSeriesGetOK, toSeriesResponse(series))
}

func (h *SeriesHandler) Update(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var input request.UpdateSeriesRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("image")

	series, err := h.service.Update(id, input, file)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgSeriesUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgSeriesUpdateOK, toSeriesResponse(series))
}

func (h *SeriesHandler) Delete(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgSeriesDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgSeriesDeleteOK, nil)
}

func (h *SeriesHandler) AddItem(c *gin.Context) {
	var input request.AddSeriesItemRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	if err := h.service.AddItem(input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgSeriesItemAddFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgSeriesItemAddOK, nil)
}

func (h *SeriesHandler) RemoveItem(c *gin.Context) {
	seriesID, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}
	audioID, ok := helper.ParamToUint(c, "audio_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidAudioID, nil)
		return
	}

	if err := h.service.RemoveItem(seriesID, audioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgSeriesItemRemoveFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgSeriesItemRemoveOK, nil)
}

func (h *SeriesHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgSearchRequired, nil)
		return
	}

	series, err := h.service.Search(query)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgSeriesSearchFail, err.Error())
		return
	}

	var result []response.SeriesResponse
	for _, s := range series {
		result = append(result, toSeriesResponseVal(s))
	}

	resp.Success(c, http.StatusOK, constant.MsgSeriesSearchOK, result)
}
