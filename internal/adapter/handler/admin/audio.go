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
	"mqfm-backend/internal/shared/security"
)

type AudioHandler struct {
	audioService   port.AudioService
	historyService port.HistoryService
}

func NewAudioHandler(as port.AudioService, hs port.HistoryService) *AudioHandler {
	return &AudioHandler{audioService: as, historyService: hs}
}

func (h *AudioHandler) Create(c *gin.Context) {
	var input request.CreateAudioRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("file")

	audio, err := h.audioService.Create(input, file)
	if err != nil {
		logger.Error("audio create: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgAudioCreateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusCreated, constant.MsgAudioCreateOK, toAudioResponse(audio))
}

func (h *AudioHandler) FindAll(c *gin.Context) {
	audios, err := h.audioService.FindAll()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgAudioListFail, err.Error())
		return
	}

	var result []response.AudioResponse
	for _, a := range audios {
		result = append(result, toAudioResponseVal(a))
	}

	resp.Success(c, http.StatusOK, constant.MsgAudioListOK, result)
}

func (h *AudioHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	audio, err := h.audioService.FindByID(id)
	if err != nil {
		resp.Error(c, http.StatusNotFound, constant.MsgAudioNotFound, err.Error())
		return
	}

	userID := security.GetUserID(c)
	if userID != 0 {
		histReq := request.HistoryRequest{AudioID: id}
		_ = h.historyService.RecordPlay(userID, histReq)
	}

	resp.Success(c, http.StatusOK, constant.MsgAudioGetOK, toAudioResponse(audio))
}

func (h *AudioHandler) Update(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	var input request.UpdateAudioRequest
	if err := c.ShouldBind(&input); err != nil {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidInput, err.Error())
		return
	}

	file, _ := c.FormFile("file")

	audio, err := h.audioService.Update(id, input, file)
	if err != nil {
		logger.Error("audio update: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgAudioUpdateFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgAudioUpdateOK, toAudioResponse(audio))
}

func (h *AudioHandler) Delete(c *gin.Context) {
	id, ok := helper.ParamToUint(c, "id")
	if !ok {
		resp.Error(c, http.StatusBadRequest, constant.MsgInvalidIDFormat, nil)
		return
	}

	if err := h.audioService.Delete(id); err != nil {
		logger.Error("audio delete: " + err.Error())
		resp.Error(c, http.StatusInternalServerError, constant.MsgAudioDeleteFail, err.Error())
		return
	}

	resp.Success(c, http.StatusOK, constant.MsgAudioDeleteOK, nil)
}

func (h *AudioHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		resp.Error(c, http.StatusBadRequest, constant.MsgSearchRequired, nil)
		return
	}

	audios, err := h.audioService.Search(query)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, constant.MsgAudioSearchFail, err.Error())
		return
	}

	var result []response.AudioResponse
	for _, a := range audios {
		result = append(result, toAudioResponseVal(a))
	}

	resp.Success(c, http.StatusOK, constant.MsgAudioSearchOK, result)
}
