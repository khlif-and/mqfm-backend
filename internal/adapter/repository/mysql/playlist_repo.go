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
	err := r.db.Where("user_id = ?", userID).Preload("Audios").Preload("User").Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepo) FindByID(id uint) (*entity.Playlist, error) {
	var playlist entity.Playlist
	err := r.db.Where("id = ?", id).Preload("Audios").Preload("User").First(&playlist).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *playlistRepo) FindAll() ([]entity.Playlist, error) {
	var playlists []entity.Playlist
	err := r.db.Preload("Audios").Preload("User").Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepo) Search(query string) ([]entity.Playlist, error) {
	var playlists []entity.Playlist
	err := r.db.Where("name LIKE ?", "%"+query+"%").Preload("Audios").Preload("User").Find(&playlists).Error
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

func (r *playlistRepo) RemoveAudio(playlist *entity.Playlist, audio *entity.Audio) error {
	return r.db.Model(playlist).Association("Audios").Delete(audio)
}

func (r *playlistRepo) FindByShareToken(token string) (*entity.Playlist, error) {
	var playlist entity.Playlist
	err := r.db.Where("share_token = ?", token).Preload("Audios").Preload("User").First(&playlist).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *playlistRepo) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&entity.Playlist{}).Where("id = ?", id).Updates(updates).Error
}

func (r *playlistRepo) Delete(id uint) error {
	return r.db.Delete(&entity.Playlist{}, id).Error
}

func (r *playlistRepo) CountAudios(playlistID uint) (int, error) {
	var playlist entity.Playlist
	playlist.ID = playlistID
	count := r.db.Model(&playlist).Association("Audios").Count()
	return int(count), nil
}
