package mysql

import (
	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type livestreamRepo struct {
	db *gorm.DB
}

func NewLivestreamRepository(db *gorm.DB) port.LivestreamRepository {
	return &livestreamRepo{db: db}
}

func (r *livestreamRepo) FindFirst() (*entity.LiveStream, error) {
	var ls entity.LiveStream
	if err := r.db.First(&ls).Error; err != nil {
		return nil, err
	}
	return &ls, nil
}

func (r *livestreamRepo) Create(ls *entity.LiveStream) error {
	return r.db.Create(ls).Error
}

func (r *livestreamRepo) Save(ls *entity.LiveStream) error {
	return r.db.Save(ls).Error
}
