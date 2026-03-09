package mysql

import (
	"time"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type audioVoteRepo struct {
	db *gorm.DB
}

func NewAudioVoteRepository(db *gorm.DB) port.AudioVoteRepository {
	return &audioVoteRepo{db: db}
}

func (r *audioVoteRepo) Create(vote *entity.AudioVote) error {
	return r.db.Create(vote).Error
}

func (r *audioVoteRepo) Delete(userID, audioID uint) error {
	return r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).Delete(&entity.AudioVote{}).Error
}

func (r *audioVoteRepo) Exists(userID, audioID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.AudioVote{}).Where("user_id = ? AND audio_id = ?", userID, audioID).Count(&count).Error
	return count > 0, err
}

func (r *audioVoteRepo) CountByAudio(audioID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.AudioVote{}).Where("audio_id = ?", audioID).Count(&count).Error
	return count, err
}

func (r *audioVoteRepo) CountWeeklyByAudio(audioID uint) (int64, error) {
	var count int64
	weekAgo := time.Now().AddDate(0, 0, -7)
	err := r.db.Model(&entity.AudioVote{}).Where("audio_id = ? AND created_at >= ?", audioID, weekAgo).Count(&count).Error
	return count, err
}

func (r *audioVoteRepo) CountMonthlyByAudio(audioID uint) (int64, error) {
	var count int64
	monthAgo := time.Now().AddDate(0, -1, 0)
	err := r.db.Model(&entity.AudioVote{}).Where("audio_id = ? AND created_at >= ?", audioID, monthAgo).Count(&count).Error
	return count, err
}

func (r *audioVoteRepo) FindVotedAudioIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&entity.AudioVote{}).Where("user_id = ?", userID).Pluck("audio_id", &ids).Error
	return ids, err
}
