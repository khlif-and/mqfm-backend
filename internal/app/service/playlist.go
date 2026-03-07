package service

import (
	"errors"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/logger"

	"go.uber.org/zap"
)

type playlistService struct {
	repo port.PlaylistRepository
}

func NewPlaylistService(repo port.PlaylistRepository) port.PlaylistService {
	return &playlistService{repo: repo}
}

func (s *playlistService) Create(playlist *entity.Playlist) error {
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

func (s *playlistService) GetByID(id uint, userID uint) (*entity.Playlist, error) {
	return s.repo.FindByID(id, userID)
}

func (s *playlistService) Search(userID uint, query string) ([]entity.Playlist, error) {
	return s.repo.Search(userID, query)
}

func (s *playlistService) AddAudioToPlaylist(userID, playlistID, audioID uint) error {
	playlist, err := s.repo.FindByID(playlistID, userID)
	if err != nil {
		return errors.New(constant.MsgPlaylistNotFound)
	}

	if len(playlist.Audios) >= 20 {
		logger.Warn("playlist limit reached",
			zap.Uint("user_id", userID),
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

	return s.repo.AddAudio(playlist, audio)
}
