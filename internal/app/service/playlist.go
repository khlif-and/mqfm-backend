package service

import (
	"errors"
	"mime/multipart"

	"github.com/google/uuid"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/logger"

	"go.uber.org/zap"
)

type playlistService struct {
	repo     port.PlaylistRepository
	colorSvc port.ColorExtractorService
}

func NewPlaylistService(repo port.PlaylistRepository, colorSvc port.ColorExtractorService) port.PlaylistService {
	return &playlistService{repo: repo, colorSvc: colorSvc}
}

func (s *playlistService) Create(playlist *entity.Playlist, file interface{}) error {
	if playlist.ImageURL != "" {
		if color, err := s.colorSvc.ExtractDominantColor(playlist.ImageURL); err == nil {
			playlist.DominantColor = color
		}
	}

	if err := s.repo.Create(playlist); err != nil {
		logger.Error("playlist create failed",
			zap.Uint("user_id", playlist.UserID),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *playlistService) GetByUserID(userID uint) ([]entity.Playlist, error) {
	return s.repo.FindByUserID(userID)
}

func (s *playlistService) GetByID(id uint) (*entity.Playlist, error) {
	return s.repo.FindByID(id)
}

func (s *playlistService) Update(id uint, updates map[string]interface{}) (*entity.Playlist, error) {
	if imgPath, ok := updates["image_url"].(string); ok && imgPath != "" {
		if color, err := s.colorSvc.ExtractDominantColor(imgPath); err == nil {
			updates["dominant_color"] = color
		}
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *playlistService) Delete(id, userID uint) error {
	playlist, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New(constant.MsgPlaylistNotFound)
	}
	if userID != 0 && playlist.UserID != userID {
		return errors.New(constant.MsgUnauthorized)
	}
	return s.repo.Delete(id)
}

func (s *playlistService) Search(query string) ([]entity.Playlist, error) {
	return s.repo.Search(query)
}

func (s *playlistService) AddAudioToPlaylist(playlistID, audioID uint) error {
	playlist, err := s.repo.FindByID(playlistID)
	if err != nil {
		return errors.New(constant.MsgPlaylistNotFound)
	}

	count, _ := s.repo.CountAudios(playlistID)
	if count >= 20 {
		logger.Warn("playlist limit reached",
			zap.Uint("playlist_id", playlistID),
		)
		return errors.New(constant.MsgPlaylistFull)
	}

	for _, a := range playlist.Audios {
		if a.ID == audioID {
			return errors.New(constant.MsgAudioAlreadyInPlaylist)
		}
	}

	audio, err := s.repo.FindAudioByID(audioID)
	if err != nil {
		return errors.New(constant.MsgAudioNotFound)
	}

	if err := s.repo.AddAudio(playlist, audio); err != nil {
		return err
	}

	if playlist.DominantColor == "" && audio.Thumbnail != "" {
		if color, err := s.colorSvc.ExtractDominantColor(audio.Thumbnail); err == nil {
			_ = s.repo.Update(playlistID, map[string]interface{}{
				"dominant_color": color,
			})
		}
	}

	return nil
}

func (s *playlistService) RemoveAudioFromPlaylist(playlistID, audioID uint) error {
	playlist, err := s.repo.FindByID(playlistID)
	if err != nil {
		return errors.New(constant.MsgPlaylistNotFound)
	}

	audio, err := s.repo.FindAudioByID(audioID)
	if err != nil {
		return errors.New(constant.MsgAudioNotFound)
	}

	return s.repo.RemoveAudio(playlist, audio)
}

func (s *playlistService) SharePlaylist(playlistID uint) (string, error) {
	playlist, err := s.repo.FindByID(playlistID)
	if err != nil {
		return "", errors.New(constant.MsgPlaylistNotFound)
	}

	if playlist.ShareToken != "" {
		return playlist.ShareToken, nil
	}

	token := uuid.New().String()
	if err := s.repo.Update(playlistID, map[string]interface{}{
		"share_token": token,
		"is_public":   true,
	}); err != nil {
		return "", err
	}

	return token, nil
}

func (s *playlistService) GetByShareToken(token string) (*entity.Playlist, error) {
	return s.repo.FindByShareToken(token)
}

var _ port.PlaylistService = (*playlistService)(nil)
var _ = (*multipart.FileHeader)(nil)
