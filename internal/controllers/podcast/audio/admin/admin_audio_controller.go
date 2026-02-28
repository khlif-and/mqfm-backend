package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	audioDto "mqfm-backend/internal/dto/audio"
	historyDto "mqfm-backend/internal/dto/history"
	historyService "mqfm-backend/internal/services/history/user"
	audioService "mqfm-backend/internal/services/podcast/audio/admin"
	"mqfm-backend/internal/utils"
)

type AdminAudioController struct {
	service        *audioService.AdminAudioService
	historyService *historyService.UserHistoryService
}

func NewAdminAudioController(s *audioService.AdminAudioService, hs *historyService.UserHistoryService) *AdminAudioController {
	return &AdminAudioController{service: s, historyService: hs}
}

func (ctrl *AdminAudioController) Create(c *gin.Context) {
	var input audioDto.CreateAudioRequest

	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	file, _ := c.FormFile("file")

	audio, err := ctrl.service.Create(input, file)
	if err != nil {
		utils.Log.Error("Audio creation error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create audio", err.Error())
		return
	}

	response := audioDto.AudioResponse{
		ID:         audio.ID,
		Title:      audio.Title,
		Artist:     audio.Artist,
		FilePath:   audio.FilePath,
		Duration:   audio.Duration,
		Status:     audio.Status,
		CategoryID: audio.CategoryID,
		CreatedAt:  audio.CreatedAt,
		UpdatedAt:  audio.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Audio created successfully", response)
}

func (ctrl *AdminAudioController) FindAll(c *gin.Context) {
	audios, err := ctrl.service.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch audios", err.Error())
		return
	}

	var response []audioDto.AudioResponse
	for _, audio := range audios {
		response = append(response, audioDto.AudioResponse{
			ID:         audio.ID,
			Title:      audio.Title,
			Artist:     audio.Artist,
			FilePath:   audio.FilePath,
			Duration:   audio.Duration,
			Status:     audio.Status,
			CategoryID: audio.CategoryID,
			CreatedAt:  audio.CreatedAt,
			UpdatedAt:  audio.UpdatedAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Audios retrieved successfully", response)
}

func (ctrl *AdminAudioController) FindByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	audio, err := ctrl.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Audio not found", err.Error())
		return
	}

	userID := utils.GetUserID(c)
	if userID != 0 {
		req := historyDto.HistoryRequest{AudioID: uint(id)}
		_ = ctrl.historyService.RecordPlay(userID, req)
	}

	response := audioDto.AudioResponse{
		ID:         audio.ID,
		Title:      audio.Title,
		Artist:     audio.Artist,
		FilePath:   audio.FilePath,
		Duration:   audio.Duration,
		Status:     audio.Status,
		CategoryID: audio.CategoryID,
		CreatedAt:  audio.CreatedAt,
		UpdatedAt:  audio.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "Audio retrieved successfully", response)
}

func (ctrl *AdminAudioController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var input audioDto.UpdateAudioRequest
	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid update data", err.Error())
		return
	}

	file, _ := c.FormFile("file")

	updatedAudio, err := ctrl.service.Update(uint(id), input, file)
	if err != nil {
		utils.Log.Error("Audio update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update audio", err.Error())
		return
	}

	response := audioDto.AudioResponse{
		ID:         updatedAudio.ID,
		Title:      updatedAudio.Title,
		Artist:     updatedAudio.Artist,
		FilePath:   updatedAudio.FilePath,
		Duration:   updatedAudio.Duration,
		Status:     updatedAudio.Status,
		CategoryID: updatedAudio.CategoryID,
		CreatedAt:  updatedAudio.CreatedAt,
		UpdatedAt:  updatedAudio.UpdatedAt,
	}

	utils.SuccessResponse(c, http.StatusOK, "Audio updated successfully", response)
}

func (ctrl *AdminAudioController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	if err := ctrl.service.Delete(uint(id)); err != nil {
		utils.Log.Error("Audio deletion error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete audio", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audio deleted successfully", nil)
}

func (ctrl *AdminAudioController) Search(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Search keyword is required", nil)
		return
	}

	audios, err := ctrl.service.Search(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search audios", err.Error())
		return
	}

	var response []audioDto.AudioResponse
	for _, audio := range audios {
		response = append(response, audioDto.AudioResponse{
			ID:         audio.ID,
			Title:      audio.Title,
			Artist:     audio.Artist,
			FilePath:   audio.FilePath,
			Duration:   audio.Duration,
			Status:     audio.Status,
			CategoryID: audio.CategoryID,
			CreatedAt:  audio.CreatedAt,
			UpdatedAt:  audio.UpdatedAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Audios found successfully", response)
}