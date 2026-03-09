package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type favoriteArtistRepo struct {
	db *gorm.DB
}

func NewFavoriteArtistRepository(db *gorm.DB) port.FavoriteArtistRepository {
	return &favoriteArtistRepo{db: db}
}

func (r *favoriteArtistRepo) Create(fav *entity.FavoriteArtist) error {
	return r.db.Create(fav).Error
}

func (r *favoriteArtistRepo) Delete(userID uint, artistName string) error {
	return r.db.Where("user_id = ? AND artist_name = ?", userID, artistName).Delete(&entity.FavoriteArtist{}).Error
}

func (r *favoriteArtistRepo) FindByUser(userID uint) ([]entity.FavoriteArtist, error) {
	var favs []entity.FavoriteArtist
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favs).Error
	return favs, err
}

func (r *favoriteArtistRepo) Exists(userID uint, artistName string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.FavoriteArtist{}).Where("user_id = ? AND artist_name = ?", userID, artistName).Count(&count).Error
	return count > 0, err
}
