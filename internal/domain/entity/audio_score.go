package entity

import "time"

type AudioScore struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AudioID     uint      `gorm:"uniqueIndex;not null" json:"audio_id"`
	Audio       *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	TotalPlays   int64     `gorm:"default:0" json:"total_plays"`
	TotalLikes   int64     `gorm:"default:0" json:"total_likes"`
	WeeklyLikes  int64     `gorm:"default:0;index:idx_weekly_likes" json:"weekly_likes"`
	MonthlyLikes int64     `gorm:"default:0;index:idx_monthly_likes" json:"monthly_likes"`
	WeightScore  float64   `gorm:"default:0;index:idx_weight_score" json:"weight_score"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AudioScore) TableName() string {
	return "audio_scores"
}
