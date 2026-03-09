package service

import (
	"errors"

	"github.com/google/uuid"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type playlistCollabService struct {
	collabRepo   port.PlaylistCollaboratorRepository
	playlistRepo port.PlaylistRepository
}

func NewPlaylistCollabService(
	collabRepo port.PlaylistCollaboratorRepository,
	playlistRepo port.PlaylistRepository,
) port.PlaylistCollabService {
	return &playlistCollabService{
		collabRepo:   collabRepo,
		playlistRepo: playlistRepo,
	}
}

func (s *playlistCollabService) AddCollaborator(ownerID, playlistID, collaboratorID uint) error {
	playlist, err := s.playlistRepo.FindByID(playlistID)
	if err != nil || playlist.UserID != ownerID {
		return errors.New("only owner can add collaborators")
	}

	exists, err := s.collabRepo.Exists(playlistID, collaboratorID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user is already a collaborator")
	}

	return s.collabRepo.Create(&entity.PlaylistCollaborator{
		PlaylistID: playlistID,
		UserID:     collaboratorID,
		Role:       "contributor",
	})
}

func (s *playlistCollabService) RemoveCollaborator(ownerID, playlistID, collaboratorID uint) error {
	playlist, err := s.playlistRepo.FindByID(playlistID)
	if err != nil || playlist.UserID != ownerID {
		return errors.New("only owner can remove collaborators")
	}

	return s.collabRepo.Delete(playlistID, collaboratorID)
}

func (s *playlistCollabService) GetCollaborators(playlistID uint) ([]entity.PlaylistCollaborator, error) {
	return s.collabRepo.FindByPlaylist(playlistID)
}

func (s *playlistCollabService) ContributeAudio(userID, playlistID, audioID uint) error {
	allowed, err := s.collabRepo.IsOwnerOrCollaborator(playlistID, userID)
	if err != nil || !allowed {
		return errors.New("not authorized to contribute")
	}

	audio, err := s.playlistRepo.FindAudioByID(audioID)
	if err != nil {
		return errors.New("audio not found")
	}

	playlist, err := s.playlistRepo.FindByID(playlistID)
	if err != nil {
		return errors.New("playlist not found")
	}

	return s.playlistRepo.AddAudio(playlist, audio)
}

func (s *playlistCollabService) JoinByShareToken(userID uint, token string) error {
	playlist, err := s.playlistRepo.FindByShareToken(token)
	if err != nil {
		return errors.New("invalid share token")
	}
	if !playlist.IsPublic {
		return errors.New("playlist is not public")
	}

	exists, err := s.collabRepo.Exists(playlist.ID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already a collaborator")
	}

	return s.collabRepo.Create(&entity.PlaylistCollaborator{
		PlaylistID: playlist.ID,
		UserID:     userID,
		Role:       "contributor",
	})
}

func generateShareToken() string {
	return uuid.New().String()
}
