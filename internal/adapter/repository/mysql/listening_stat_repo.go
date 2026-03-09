package mysql

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type listeningStatRepo struct {
	db *gorm.DB
}

func NewListeningStatRepository(db *gorm.DB) port.ListeningStatRepository {
	return &listeningStatRepo{db: db}
}

func (r *listeningStatRepo) Create(stat *entity.ListeningStat) error {
	return r.db.Create(stat).Error
}

func (r *listeningStatRepo) GetWeeklySummary(userID uint) (int, error) {
	var total int
	weekAgo := time.Now().AddDate(0, 0, -7)
	err := r.db.Model(&entity.ListeningStat{}).
		Where("user_id = ? AND listened_at >= ?", userID, weekAgo).
		Select("COALESCE(SUM(duration_seconds), 0)").Scan(&total).Error
	return total, err
}

func (r *listeningStatRepo) GetMonthlySummary(userID uint) (int, error) {
	var total int
	monthAgo := time.Now().AddDate(0, -1, 0)
	err := r.db.Model(&entity.ListeningStat{}).
		Where("user_id = ? AND listened_at >= ?", userID, monthAgo).
		Select("COALESCE(SUM(duration_seconds), 0)").Scan(&total).Error
	return total, err
}

func (r *listeningStatRepo) GetTopCategories(userID uint, limit int) ([]port.CategoryStat, error) {
	var stats []port.CategoryStat
	err := r.db.Raw(`
		SELECT a.category_id, c.name, COALESCE(SUM(ls.duration_seconds), 0) as total_time
		FROM listening_stats ls
		JOIN audios a ON a.id = ls.audio_id
		LEFT JOIN categories c ON c.id = a.category_id
		WHERE ls.user_id = ?
		GROUP BY a.category_id, c.name
		ORDER BY total_time DESC
		LIMIT ?
	`, userID, limit).Scan(&stats).Error
	return stats, err
}

func (r *listeningStatRepo) GetTopArtists(userID uint, limit int) ([]port.ArtistStat, error) {
	var stats []port.ArtistStat
	err := r.db.Raw(`
		SELECT a.artist, COALESCE(SUM(ls.duration_seconds), 0) as total_time
		FROM listening_stats ls
		JOIN audios a ON a.id = ls.audio_id
		WHERE ls.user_id = ? AND a.artist != ''
		GROUP BY a.artist
		ORDER BY total_time DESC
		LIMIT ?
	`, userID, limit).Scan(&stats).Error
	return stats, err
}

func (r *listeningStatRepo) GetDailySummary(userID uint, days int) ([]port.DailyStat, error) {
	var stats []port.DailyStat
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.Raw(fmt.Sprintf(`
		SELECT DATE(listened_at) as date, COALESCE(SUM(duration_seconds), 0) as total_time
		FROM listening_stats
		WHERE user_id = ? AND listened_at >= ?
		GROUP BY DATE(listened_at)
		ORDER BY date DESC
		LIMIT %d
	`, days), userID, since).Scan(&stats).Error
	return stats, err
}
