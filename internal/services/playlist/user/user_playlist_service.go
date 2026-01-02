package user

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	adminAudioModel "mqfm-backend/internal/models/podcast/audio/admin"
	playlistModel "mqfm-backend/internal/models/playlist/user"
	"mqfm-backend/internal/utils"
)

type UserPlaylistService struct {
	db *gorm.DB
}

func NewUserPlaylistService(db *gorm.DB) *UserPlaylistService {
	return &UserPlaylistService{db: db}
}

func (s *UserPlaylistService) Create(playlist *playlistModel.Playlist) error {
	if err := s.db.Create(playlist).Error; err != nil {
		utils.Log.Error("[Playlist] Failed to create playlist",
			zap.Error(err),
			zap.Uint("user_id", playlist.UserID),
			zap.String("name", playlist.Name),
		)
		return err
	}
	
	utils.Log.Info("[Playlist] Playlist created",
		zap.Uint("id", playlist.ID),
		zap.Uint("user_id", playlist.UserID),
	)
	return nil
}

func (s *UserPlaylistService) GetByUserID(userID uint) ([]playlistModel.Playlist, error) {
	var playlists []playlistModel.Playlist
	err := s.db.Where("user_id = ?", userID).Preload("Audios").Find(&playlists).Error
	if err != nil {
		utils.Log.Error("[Playlist] Failed to fetch playlists",
			zap.Error(err),
			zap.Uint("user_id", userID),
		)
	}
	return playlists, err
}

func (s *UserPlaylistService) Search(userID uint, query string) ([]playlistModel.Playlist, error) {
	var playlists []playlistModel.Playlist
	err := s.db.Where("user_id = ? AND name LIKE ?", userID, "%"+query+"%").Preload("Audios").Find(&playlists).Error
	if err != nil {
		utils.Log.Error("[Playlist] Failed to search playlists",
			zap.Error(err),
			zap.Uint("user_id", userID),
			zap.String("query", query),
		)
	}
	return playlists, err
}

func (s *UserPlaylistService) GetByID(id uint, userID uint) (*playlistModel.Playlist, error) {
	var playlist playlistModel.Playlist
	err := s.db.Where("id = ? AND user_id = ?", id, userID).Preload("Audios").First(&playlist).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Log.Error("[Playlist] Failed to get playlist detail",
				zap.Error(err),
				zap.Uint("playlist_id", id),
				zap.Uint("user_id", userID),
			)
		}
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
		utils.Log.Warn("[Playlist] Playlist limit reached",
			zap.Uint("user_id", userID),
			zap.Uint("playlist_id", playlistID),
		)
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

	if err := s.db.Model(&playlist).Association("Audios").Append(&audio); err != nil {
		utils.Log.Error("[Playlist] Failed to append audio association",
			zap.Error(err),
			zap.Uint("playlist_id", playlistID),
			zap.Uint("audio_id", audioID),
		)
		return err
	}

	utils.Log.Info("[Playlist] Audio added to playlist",
		zap.Uint("user_id", userID),
		zap.Uint("playlist_id", playlistID),
		zap.Uint("audio_id", audioID),
	)

	return nil
}

func (s *UserPlaylistService) CreateAndAddAudio(playlist *playlistModel.Playlist, audioID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(playlist).Error; err != nil {
			utils.Log.Error("[Playlist] Transaction: Failed to create playlist", zap.Error(err))
			return err
		}

		var audio adminAudioModel.Audio
		if err := tx.First(&audio, audioID).Error; err != nil {
			return errors.New("audio not found")
		}

		if err := tx.Model(playlist).Association("Audios").Append(&audio); err != nil {
			utils.Log.Error("[Playlist] Transaction: Failed to append audio", zap.Error(err))
			return err
		}

		utils.Log.Info("[Playlist] Created playlist and added audio via transaction",
			zap.Uint("user_id", playlist.UserID),
			zap.Uint("playlist_id", playlist.ID),
			zap.Uint("audio_id", audioID),
		)

		return nil
	})
}