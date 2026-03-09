package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type playlistCollabRepo struct {
	db *gorm.DB
}

func NewPlaylistCollaboratorRepository(db *gorm.DB) port.PlaylistCollaboratorRepository {
	return &playlistCollabRepo{db: db}
}

func (r *playlistCollabRepo) Create(collab *entity.PlaylistCollaborator) error {
	return r.db.Create(collab).Error
}

func (r *playlistCollabRepo) Delete(playlistID, userID uint) error {
	return r.db.Where("playlist_id = ? AND user_id = ?", playlistID, userID).
		Delete(&entity.PlaylistCollaborator{}).Error
}

func (r *playlistCollabRepo) FindByPlaylist(playlistID uint) ([]entity.PlaylistCollaborator, error) {
	var collabs []entity.PlaylistCollaborator
	err := r.db.Where("playlist_id = ?", playlistID).Find(&collabs).Error
	return collabs, err
}

func (r *playlistCollabRepo) Exists(playlistID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.PlaylistCollaborator{}).
		Where("playlist_id = ? AND user_id = ?", playlistID, userID).Count(&count).Error
	return count > 0, err
}

func (r *playlistCollabRepo) IsOwnerOrCollaborator(playlistID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Playlist{}).Where("id = ? AND user_id = ?", playlistID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return r.Exists(playlistID, userID)
}
