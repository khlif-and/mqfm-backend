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

type RadioHandler struct {
	service port.RadioService
}

func NewRadioHandler(s port.RadioService) *RadioHandler {
	return &RadioHandler{service: s}
}

func (h *RadioHandler) Create(c *gin.Context) {
	var input request.CreateRadioRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("thumbnail")

	radio, err := h.service.Create(input, file)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRadioCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgRadioCreateOK, toRadioResponse(radio))
}

func (h *RadioHandler) FindAll(c *gin.Context) {
	radios, err := h.service.FindAll()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRadioListFail, err.Error())
		return
	}

	var result []response.RadioResponse
	for _, r := range radios {
		result = append(result, toRadioResponseVal(r))
	}

	resp.Success(c, http.StatusOK, constant.MsgRadioListOK, result)
}

func (h *RadioHandler) FindActive(c *gin.Context) {
	radios, err := h.service.FindActive()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRadioListFail, err.Error())
		return
	}

	var result []response.RadioResponse
	for _, r := range radios {
		result = append(result, toRadioResponseVal(r))
	}

	resp.Success(c, http.StatusOK, constant.MsgRadioListOK, result)
}

func (h *RadioHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	radio, err := h.service.FindByID(id)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgRadioNotFound, nil)
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgRadioGetOK, toRadioResponse(radio))
}

func (h *RadioHandler) Update(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var input request.UpdateRadioRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("thumbnail")

	radio, err := h.service.Update(id, input, file)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRadioUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgRadioUpdateOK, toRadioResponse(radio))
}

func (h *RadioHandler) Delete(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgRadioDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgRadioDeleteOK, nil)
}

func (h *RadioHandler) AddAudio(c *gin.Context) {
	var input request.RadioAudioRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	if err := h.service.AddAudio(input.RadioID, input.AudioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgRadioAudioAddFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgRadioAudioAddOK, nil)
}

func (h *RadioHandler) RemoveAudio(c *gin.Context) {
	radioID, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}
	audioID, ok := helper.ParamToUint(c, "audio_id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	if err := h.service.RemoveAudio(radioID, audioID); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgRadioAudioRemoveFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgRadioAudioRemoveOK, nil)
}
