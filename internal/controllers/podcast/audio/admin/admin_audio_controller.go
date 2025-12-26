package admin

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	audioModel "mqfm-backend/internal/models/podcast/audio/admin"
	categoryService "mqfm-backend/internal/services/category/admin" // Import Service Category
	audioService "mqfm-backend/internal/services/podcast/audio/admin"
	"mqfm-backend/internal/utils"

)

type AdminAudioController struct {
	service         *audioService.AdminAudioService
	categoryService *categoryService.AdminCategoryService // Tambahkan field ini
}

// Update Constructor: Menerima Category Service juga
func NewAdminAudioController(s *audioService.AdminAudioService, cs *categoryService.AdminCategoryService) *AdminAudioController {
	return &AdminAudioController{
		service:         s,
		categoryService: cs,
	}
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

	// --- VALIDASI KATEGORI ---
	// Jika CategoryID diisi (tidak 0), cek apakah ada di database
	if input.CategoryID != 0 {
		if _, err := ctrl.categoryService.FindByID(input.CategoryID); err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Category ID not found", err.Error())
			return
		}
	}
	// -------------------------

	pwd, _ := os.Getwd()
	fmt.Println("DEBUG: Aplikasi berjalan di:", pwd)

	var audioPathDB string
	if input.AudioFile != nil {
		uploadDir := filepath.Join(pwd, "uploads", "audios")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create audio directory", err.Error())
			return
		}

		audioFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.AudioFile.Filename)
		fullSavePath := filepath.Join(uploadDir, audioFilename)

		fmt.Println("DEBUG: Menyimpan Audio ke:", fullSavePath)

		if err := c.SaveUploadedFile(input.AudioFile, fullSavePath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload audio file", err.Error())
			return
		}

		audioPathDB = "uploads/audios/" + audioFilename
	}

	var thumbnailPathDB string
	if input.ThumbnailFile != nil {
		uploadDir := filepath.Join(pwd, "uploads", "thumbnails")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create thumbnail directory", err.Error())
			return
		}

		thumbFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.ThumbnailFile.Filename)
		fullSavePath := filepath.Join(uploadDir, thumbFilename)

		fmt.Println("DEBUG: Menyimpan Thumbnail ke:", fullSavePath)

		if err := c.SaveUploadedFile(input.ThumbnailFile, fullSavePath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload thumbnail", err.Error())
			return
		}

		thumbnailPathDB = "uploads/thumbnails/" + thumbFilename
	}

	audio := audioModel.Audio{
		Title:       input.Title,
		Description: input.Description,
		AudioURL:    audioPathDB,
		Thumbnail:   thumbnailPathDB,
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

	// --- VALIDASI KATEGORI (UPDATE) ---
	if input.CategoryID != 0 {
		if _, err := ctrl.categoryService.FindByID(input.CategoryID); err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Category ID not found", err.Error())
			return
		}
	}
	// ----------------------------------

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

	pwd, _ := os.Getwd()

	if input.AudioFile != nil {
		uploadDir := filepath.Join(pwd, "uploads", "audios")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create audio directory", err.Error())
			return
		}

		audioFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.AudioFile.Filename)
		fullSavePath := filepath.Join(uploadDir, audioFilename)

		if err := c.SaveUploadedFile(input.AudioFile, fullSavePath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload audio file", err.Error())
			return
		}
		updates["audio_url"] = "uploads/audios/" + audioFilename
	}

	if input.ThumbnailFile != nil {
		uploadDir := filepath.Join(pwd, "uploads", "thumbnails")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create thumbnail directory", err.Error())
			return
		}

		thumbFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), input.ThumbnailFile.Filename)
		fullSavePath := filepath.Join(uploadDir, thumbFilename)

		if err := c.SaveUploadedFile(input.ThumbnailFile, fullSavePath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload thumbnail", err.Error())
			return
		}
		updates["thumbnail"] = "uploads/thumbnails/" + thumbFilename
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

	utils.SuccessResponse(c, http.StatusOK, "Audios found successfully", audios)
}