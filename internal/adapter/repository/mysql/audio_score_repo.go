package mysql

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type audioScoreRepo struct {
	db *gorm.DB
}

func NewAudioScoreRepository(db *gorm.DB) port.AudioScoreRepository {
	return &audioScoreRepo{db: db}
}

func (r *audioScoreRepo) Upsert(score *entity.AudioScore) error {
	score.UpdatedAt = time.Now()
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "audio_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"total_plays", "total_likes", "weight_score", "updated_at"}),
	}).Create(score).Error
}

func (r *audioScoreRepo) FindTopByScore(limit int) ([]entity.AudioScore, error) {
	var scores []entity.AudioScore
	err := r.db.Preload("Audio").
		Order("weight_score DESC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *audioScoreRepo) FindTopByLikes(limit int, maxLikes int64) ([]entity.AudioScore, error) {
	var scores []entity.AudioScore
	err := r.db.Preload("Audio").
		Where("total_likes <= ?", maxLikes).
		Order("total_likes DESC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *audioScoreRepo) FindTopByWeeklyLikes(limit int) ([]entity.AudioScore, error) {
	var scores []entity.AudioScore
	err := r.db.Preload("Audio").
		Order("weekly_likes DESC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *audioScoreRepo) FindTopByMonthlyLikes(limit int) ([]entity.AudioScore, error) {
	var scores []entity.AudioScore
	err := r.db.Preload("Audio").
		Order("monthly_likes DESC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *audioScoreRepo) BulkUpdateWeeklyLikes(data map[uint]int64) error {
	if len(data) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		for audioID, likes := range data {
			if err := tx.Model(&entity.AudioScore{}).
				Where("audio_id = ?", audioID).
				Update("weekly_likes", likes).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *audioScoreRepo) BulkUpdateMonthlyLikes(data map[uint]int64) error {
	if len(data) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		for audioID, likes := range data {
			if err := tx.Model(&entity.AudioScore{}).
				Where("audio_id = ?", audioID).
				Update("monthly_likes", likes).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *audioScoreRepo) FindByAudioID(audioID uint) (*entity.AudioScore, error) {
	var score entity.AudioScore
	err := r.db.Where("audio_id = ?", audioID).First(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (r *audioScoreRepo) FindByAudioIDs(audioIDs []uint) ([]entity.AudioScore, error) {
	var scores []entity.AudioScore
	if len(audioIDs) == 0 {
		return scores, nil
	}
	err := r.db.Where("audio_id IN ?", audioIDs).Find(&scores).Error
	return scores, err
}

func (r *audioScoreRepo) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.AudioScore{}).Error
}

func (r *audioScoreRepo) BulkUpsert(scores []entity.AudioScore) error {
	if len(scores) == 0 {
		return nil
	}

	now := time.Now()
	for i := range scores {
		scores[i].UpdatedAt = now
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "audio_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"total_plays", "total_likes", "weight_score", "updated_at"}),
	}).CreateInBatches(scores, 100).Error
}
