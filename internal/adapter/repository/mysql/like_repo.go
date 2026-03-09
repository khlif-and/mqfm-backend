package mysql

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type likeRepo struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) port.LikeRepository {
	return &likeRepo{db: db}
}

func (r *likeRepo) Create(like *entity.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepo) Delete(userID uint, targetType string, targetID uint) error {
	result := r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Delete(&entity.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (r *likeRepo) FindByUser(userID uint, targetType string) ([]entity.Like, error) {
	var likes []entity.Like
	q := r.db.Where("user_id = ?", userID)
	if targetType != "" {
		q = q.Where("target_type = ?", targetType)
	}
	if targetType == "audio" {
		q = q.Preload("Audio")
	} else if targetType == "playlist" {
		q = q.Preload("Playlist")
	}
	err := q.Order("created_at DESC").Find(&likes).Error
	return likes, err
}

func (r *likeRepo) Exists(userID uint, targetType string, targetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count).Error
	return count > 0, err
}

func (r *likeRepo) CountByTarget(targetType string, targetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error
	return count, err
}

func (r *likeRepo) AggregateLikeCounts() (map[uint]int64, error) {
	type result struct {
		TargetID uint  `gorm:"column:target_id"`
		Total    int64 `gorm:"column:total"`
	}
	var results []result
	err := r.db.Model(&entity.Like{}).
		Select("target_id, COUNT(*) as total").
		Where("target_type = ?", "audio").
		Group("target_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[uint]int64, len(results))
	for _, r := range results {
		m[r.TargetID] = r.Total
	}
	return m, nil
}

type historyRepo struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) port.HistoryRepository {
	return &historyRepo{db: db}
}

func (r *historyRepo) Upsert(history *entity.History) error {
	var existing entity.History
	err := r.db.Where("user_id = ? AND audio_id = ?", history.UserID, history.AudioID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		history.PlayCount = 1
		history.PlayedAt = time.Now()
		return r.db.Create(history).Error
	}
	if err != nil {
		return err
	}

	return r.db.Model(&existing).Updates(map[string]interface{}{
		"play_count": existing.PlayCount + 1,
		"played_at":  time.Now(),
	}).Error
}

func (r *historyRepo) FindByUser(userID uint) ([]entity.History, error) {
	var histories []entity.History
	err := r.db.Preload("Audio").Where("user_id = ?", userID).Order("played_at DESC").Find(&histories).Error
	return histories, err
}

func (r *historyRepo) FindByUsers(userIDs []uint) ([]entity.History, error) {
	var histories []entity.History
	err := r.db.Where("user_id IN ?", userIDs).Find(&histories).Error
	return histories, err
}

func (r *historyRepo) DeleteByUserAndAudio(userID, audioID uint) error {
	result := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&entity.History{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *historyRepo) DeleteAllByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.History{}).Error
}

func (r *historyRepo) CountByAudio(audioID uint) (int64, error) {
	var total int64
	err := r.db.Model(&entity.History{}).
		Select("COALESCE(SUM(play_count), 0)").
		Where("audio_id = ?", audioID).
		Scan(&total).Error
	return total, err
}

func (r *historyRepo) FindFrequentByUser(userID uint, minPlays int, limit int) ([]entity.History, error) {
	var histories []entity.History
	err := r.db.Preload("Audio").
		Where("user_id = ? AND play_count >= ?", userID, minPlays).
		Order("play_count DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *historyRepo) AggregatePlayCounts() (map[uint]int64, error) {
	type result struct {
		AudioID uint
		Total   int64
	}
	var results []result
	err := r.db.Model(&entity.History{}).
		Select("audio_id, SUM(play_count) as total").
		Group("audio_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[uint]int64, len(results))
	for _, r := range results {
		m[r.AudioID] = r.Total
	}
	return m, nil
}
