package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type bookmarkRepo struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) port.BookmarkRepository {
	return &bookmarkRepo{db: db}
}

func (r *bookmarkRepo) Create(bookmark *entity.Bookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *bookmarkRepo) FindByUser(userID uint) ([]entity.Bookmark, error) {
	var bookmarks []entity.Bookmark
	err := r.db.Where("user_id = ?", userID).Preload("Audio").Order("created_at DESC").Find(&bookmarks).Error
	return bookmarks, err
}

func (r *bookmarkRepo) FindByUserAndAudio(userID, audioID uint) ([]entity.Bookmark, error) {
	var bookmarks []entity.Bookmark
	err := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).
		Order("position_seconds ASC").Find(&bookmarks).Error
	return bookmarks, err
}

func (r *bookmarkRepo) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Bookmark{}).Error
}

func (r *bookmarkRepo) DeleteAllByUserAndAudio(userID, audioID uint) error {
	return r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&entity.Bookmark{}).Error
}
