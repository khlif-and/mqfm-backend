package service

import (
	"errors"

	"github.com/google/uuid"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type audioClipService struct {
	repo      port.AudioClipRepository
	audioRepo port.AudioRepository
	converter port.AudioConverterService
}

func NewAudioClipService(
	repo port.AudioClipRepository,
	audioRepo port.AudioRepository,
	converter port.AudioConverterService,
) port.AudioClipService {
	return &audioClipService{repo: repo, audioRepo: audioRepo, converter: converter}
}

func (s *audioClipService) CreateClip(userID uint, req request.CreateClipRequest) (*entity.AudioClip, error) {
	if req.EndTime <= req.StartTime {
		return nil, errors.New("end time must be after start time")
	}
	if req.EndTime-req.StartTime > 120 {
		return nil, errors.New("clip duration max 2 minutes")
	}

	audio, err := s.audioRepo.FindByID(req.AudioID)
	if err != nil {
		return nil, errors.New("audio not found")
	}

	sourcePath := audio.OGGPath
	if sourcePath == "" {
		sourcePath = audio.FilePath
	}

	clipPath, err := s.converter.CreateClip(sourcePath, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	clip := &entity.AudioClip{
		UserID:     userID,
		AudioID:    req.AudioID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		ClipPath:   clipPath,
		ShareToken: uuid.New().String(),
	}

	if err := s.repo.Create(clip); err != nil {
		return nil, err
	}
	return clip, nil
}

func (s *audioClipService) GetByUser(userID uint) ([]entity.AudioClip, error) {
	return s.repo.FindByUser(userID)
}

func (s *audioClipService) GetByShareToken(token string) (*entity.AudioClip, error) {
	return s.repo.FindByShareToken(token)
}

func (s *audioClipService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}
