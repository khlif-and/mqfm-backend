package admin

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	audioModel "mqfm-backend/internal/models/podcast/audio/admin"
	audioService "mqfm-backend/internal/services/podcast/audio/admin"
	"mqfm-backend/internal/utils"

)

type AdminAudioController struct {
	service *audioService.AdminAudioService
}

func NewAdminAudioController(s *audioService.AdminAudioService) *AdminAudioController {
	return &AdminAudioController{service: s}
}

func (ctrl *AdminAudioController) Create(c *gin.Context) {
	var input struct {
		Title         string                `form:"title" binding:"required"`
		Description   string                `form:"description"`
		CategoryID    uint                  `form:"category_id"`
		AudioFile     *multipart.FileHeader `form:"audio_file"`
		ThumbnailFile *multipart.FileHeader `form:"thumbnail_file"`
	}

	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	var audioPath string
	if input.AudioFile != nil {
		audioFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.AudioFile.Filename)
		audioPath = "uploads/audios/" + audioFilename
		if err := c.SaveUploadedFile(input.AudioFile, audioPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload audio file", err.Error())
			return
		}
	}

	var thumbnailPath string
	if input.ThumbnailFile != nil {
		thumbFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.ThumbnailFile.Filename)
		thumbnailPath = "uploads/thumbnails/" + thumbFilename
		if err := c.SaveUploadedFile(input.ThumbnailFile, thumbnailPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload thumbnail", err.Error())
			return
		}
	}

	audio := audioModel.Audio{
		Title:       input.Title,
		Description: input.Description,
		AudioURL:    audioPath,
		Thumbnail:   thumbnailPath,
		CategoryID:  input.CategoryID,
	}

	if err := ctrl.service.Create(&audio); err != nil {
		utils.Log.Error("Audio creation error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create audio", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Audio created successfully", audio)
}

func (ctrl *AdminAudioController) FindAll(c *gin.Context) {
	audios, err := ctrl.service.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch audios", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audios retrieved successfully", audios)
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

	utils.SuccessResponse(c, http.StatusOK, "Audio retrieved successfully", audio)
}

func (ctrl *AdminAudioController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	var input struct {
		Title         string                `form:"title"`
		Description   string                `form:"description"`
		CategoryID    uint                  `form:"category_id"`
		AudioFile     *multipart.FileHeader `form:"audio_file"`
		ThumbnailFile *multipart.FileHeader `form:"thumbnail_file"`
	}

	if err := c.ShouldBind(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid update data", err.Error())
		return
	}

	updates := make(map[string]interface{})

	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.CategoryID != 0 {
		updates["category_id"] = input.CategoryID
	}

	if input.AudioFile != nil {
		audioFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.AudioFile.Filename)
		audioPath := "uploads/audios/" + audioFilename
		if err := c.SaveUploadedFile(input.AudioFile, audioPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload audio file", err.Error())
			return
		}
		updates["audio_url"] = audioPath
	}

	if input.ThumbnailFile != nil {
		thumbFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.ThumbnailFile.Filename)
		thumbnailPath := "uploads/thumbnails/" + thumbFilename
		if err := c.SaveUploadedFile(input.ThumbnailFile, thumbnailPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload thumbnail", err.Error())
			return
		}
		updates["thumbnail"] = thumbnailPath
	}

	updatedAudio, err := ctrl.service.Update(uint(id), updates)
	if err != nil {
		utils.Log.Error("Audio update error: " + err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update audio", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audio updated successfully", updatedAudio)
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