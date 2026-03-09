package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type audioRankingRepo struct {
	db *gorm.DB
}

func NewAudioRankingRepository(db *gorm.DB) port.AudioRankingRepository {
	return &audioRankingRepo{db: db}
}

func (r *audioRankingRepo) Upsert(ranking *entity.AudioRanking) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "audio_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"weekly_votes", "monthly_votes", "total_votes",
			"random_boost", "weekly_rank", "monthly_rank",
		}),
	}).Create(ranking).Error
}

func (r *audioRankingRepo) BulkUpsert(rankings []entity.AudioRanking) error {
	if len(rankings) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "audio_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"weekly_votes", "monthly_votes", "total_votes",
			"random_boost", "weekly_rank", "monthly_rank",
		}),
	}).CreateInBatches(rankings, 200).Error
}

func (r *audioRankingRepo) FindTopWeekly(limit int) ([]entity.AudioRanking, error) {
	var rankings []entity.AudioRanking
	err := r.db.Where("weekly_rank > 0").Order("weekly_rank ASC").
		Limit(limit).Preload("Audio").Find(&rankings).Error
	return rankings, err
}

func (r *audioRankingRepo) FindTopMonthly(limit int) ([]entity.AudioRanking, error) {
	var rankings []entity.AudioRanking
	err := r.db.Where("monthly_rank > 0").Order("monthly_rank ASC").
		Limit(limit).Preload("Audio").Find(&rankings).Error
	return rankings, err
}

func (r *audioRankingRepo) FindByAudioID(audioID uint) (*entity.AudioRanking, error) {
	var ranking entity.AudioRanking
	err := r.db.Where("audio_id = ?", audioID).First(&ranking).Error
	if err != nil {
		return nil, err
	}
	return &ranking, nil
}

func (r *audioRankingRepo) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&entity.AudioRanking{}).Count(&count).Error
	return count, err
}

func (r *audioRankingRepo) FindAll(limit, offset int) ([]entity.AudioRanking, error) {
	var rankings []entity.AudioRanking
	err := r.db.Preload("Audio").Limit(limit).Offset(offset).Find(&rankings).Error
	return rankings, err
}
