package user

import (
	"errors"

	"gorm.io/gorm"

	adminAudioModel "mqfm-backend/internal/models/podcast/audio/admin"
	playlistModel "mqfm-backend/internal/models/playlist/user"
)

type UserPlaylistService struct {
	db *gorm.DB
}

func NewUserPlaylistService(db *gorm.DB) *UserPlaylistService {
	return &UserPlaylistService{db: db}
}

func (s *UserPlaylistService) Create(playlist *playlistModel.Playlist) error {
	return s.db.Create(playlist).Error
}

func (s *UserPlaylistService) GetByUserID(userID uint) ([]playlistModel.Playlist, error) {
	var playlists []playlistModel.Playlist
	err := s.db.Where("user_id = ?", userID).Preload("Audios").Find(&playlists).Error
	return playlists, err
}

func (s *UserPlaylistService) Search(userID uint, query string) ([]playlistModel.Playlist, error) {
	var playlists []playlistModel.Playlist
	err := s.db.Where("user_id = ? AND name LIKE ?", userID, "%"+query+"%").Preload("Audios").Find(&playlists).Error
	return playlists, err
}

func (s *UserPlaylistService) GetByID(id uint, userID uint) (*playlistModel.Playlist, error) {
	var playlist playlistModel.Playlist
	err := s.db.Where("id = ? AND user_id = ?", id, userID).Preload("Audios").First(&playlist).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (s *UserPlaylistService) AddAudioToPlaylist(userID uint, playlistID uint, audioID uint) error {
	var playlist playlistModel.Playlist
	if err := s.db.Where("id = ? AND user_id = ?", playlistID, userID).Preload("Audios").First(&playlist).Error; err != nil {
		return errors.New("playlist not found")
	}

	if len(playlist.Audios) >= 20 {
		return errors.New("playlist is full (max 20 audios)")
	}

	for _, audio := range playlist.Audios {
		if audio.ID == audioID {
			return errors.New("audio already in playlist")
		}
	}

	var audio adminAudioModel.Audio
	if err := s.db.First(&audio, audioID).Error; err != nil {
		return errors.New("audio not found")
	}

	return s.db.Model(&playlist).Association("Audios").Append(&audio)
}

func (s *UserPlaylistService) CreateAndAddAudio(playlist *playlistModel.Playlist, audioID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(playlist).Error; err != nil {
			return err
		}

		var audio adminAudioModel.Audio
		if err := tx.First(&audio, audioID).Error; err != nil {
			return errors.New("audio not found")
		}

		return tx.Model(playlist).Association("Audios").Append(&audio)
	})
}