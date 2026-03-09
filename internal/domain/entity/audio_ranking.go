package entity

import "time"

type AudioRanking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AudioID      uint      `gorm:"uniqueIndex;not null" json:"audio_id"`
	Audio        *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	WeeklyVotes  int64     `gorm:"default:0" json:"weekly_votes"`
	MonthlyVotes int64     `gorm:"default:0" json:"monthly_votes"`
	TotalVotes   int64     `gorm:"default:0" json:"total_votes"`
	RandomBoost  float64   `gorm:"default:0" json:"random_boost"`
	WeeklyRank   int       `gorm:"default:0;index" json:"weekly_rank"`
	MonthlyRank  int       `gorm:"default:0;index" json:"monthly_rank"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (AudioRanking) TableName() string {
	return "audio_rankings"
}
