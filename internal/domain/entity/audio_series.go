package entity

import (
	"time"

	"gorm.io/gorm"
)

type AudioSeries struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	Title       string            `gorm:"not null;size:255" json:"title"`
	Description string            `gorm:"type:text" json:"description"`
	Artist      string            `gorm:"size:255" json:"artist"`
	Image       string            `gorm:"size:500" json:"image"`
	Items       []AudioSeriesItem `gorm:"foreignKey:SeriesID" json:"items,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (AudioSeries) TableName() string {
	return "audio_series"
}

type AudioSeriesItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SeriesID  uint      `gorm:"not null;index:idx_series_order,unique" json:"series_id"`
	AudioID   uint      `gorm:"not null;index" json:"audio_id"`
	Audio     *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	OrderNum  int       `gorm:"not null;index:idx_series_order,unique" json:"order_num"`
	CreatedAt time.Time `json:"created_at"`
}

func (AudioSeriesItem) TableName() string {
	return "audio_series_items"
}
