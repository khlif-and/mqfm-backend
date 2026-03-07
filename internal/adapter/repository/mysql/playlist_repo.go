package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type playlistRepo struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) port.PlaylistRepository {
	return &playlistRepo{db: db}
}

func (r *playlistRepo) Create(playlist *entity.Playlist) error {
	return r.db.Create(playlist).Error
}

func (r *playlistRepo) FindByUserID(userID uint) ([]entity.Playlist, error) {
	var playlists []entity.Playlist
	err := r.db.Where("user_id = ?", userID).Preload("Audios").Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepo) FindByID(id uint, userID uint) (*entity.Playlist, error) {
	var playlist entity.Playlist
	err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("Audios").First(&playlist).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *playlistRepo) Search(userID uint, query string) ([]entity.Playlist, error) {
	var playlists []entity.Playlist
	err := r.db.Where("user_id = ? AND name LIKE ?", userID, "%"+query+"%").Preload("Audios").Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepo) AddAudio(playlist *entity.Playlist, audio *entity.Audio) error {
	return r.db.Model(playlist).Association("Audios").Append(audio)
}

func (r *playlistRepo) FindAudioByID(id uint) (*entity.Audio, error) {
	var audio entity.Audio
	if err := r.db.First(&audio, id).Error; err != nil {
		return nil, err
	}
	return &audio, nil
}
