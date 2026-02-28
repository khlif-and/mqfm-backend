package admin

import (
	"errors"
	"mime/multipart"

	audioDto "mqfm-backend/internal/dto/audio"
	audioModel "mqfm-backend/internal/models/podcast/audio/admin"
	audioRepo "mqfm-backend/internal/repositories/podcast/audio/admin"
	"mqfm-backend/internal/utils"
)

type AdminAudioService struct {
	repo audioRepo.AudioRepository
}

func NewAdminAudioService(repo audioRepo.AudioRepository) *AdminAudioService {
	return &AdminAudioService{repo: repo}
}

func (s *AdminAudioService) Create(req audioDto.CreateAudioRequest, file *multipart.FileHeader) (*audioModel.Audio, error) {
	var filePath string
	if file != nil {
		filename := utils.GenerateUniqueFilename(file.Filename)
		path := "uploads/audios/" + filename
		if err := utils.SaveUploadedFile(file, path); err != nil {
			utils.Log.Error("Failed to save audio file: " + err.Error())
		} else {
			filePath = path
		}
	}

	audio := audioModel.Audio{
		Title:      req.Title,
		Artist:     req.Artist,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		FilePath:   filePath,
		Duration:   0, // Should be calculated but setting 0 for now
	}

	if err := s.repo.Create(&audio); err != nil {
		return nil, err
	}

	return &audio, nil
}

func (s *AdminAudioService) FindAll() ([]audioModel.Audio, error) {
	return s.repo.FindAll()
}

func (s *AdminAudioService) FindByID(id uint) (*audioModel.Audio, error) {
	return s.repo.FindByID(id)
}

func (s *AdminAudioService) Update(id uint, req audioDto.UpdateAudioRequest, file *multipart.FileHeader) (*audioModel.Audio, error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Artist != "" {
		updates["artist"] = req.Artist
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.CategoryID != 0 {
		updates["category_id"] = req.CategoryID
	}

	if file != nil {
		filename := utils.GenerateUniqueFilename(file.Filename)
		path := "uploads/audios/" + filename
		if err := utils.SaveUploadedFile(file, path); err != nil {
			utils.Log.Error("Failed to save audio file: " + err.Error())
		} else {
			updates["file_path"] = path
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

func (s *AdminAudioService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *AdminAudioService) Search(query string) ([]audioModel.Audio, error) {
	return s.repo.Search(query)
}