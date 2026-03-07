package entity

import (
	"encoding/xml"
	"time"
)

type LiveStream struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IsLive      bool      `json:"is_live"`
	VideoID     string    `json:"video_id"`
	Title       string    `json:"title"`
	Thumbnail   string    `json:"thumbnail"`
	LastChecked time.Time `json:"last_checked"`
}

type AtomFeed struct {
	XMLName xml.Name  `xml:"feed"`
	Entry   AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	VideoID   string `xml:"videoId"`
	ChannelID string `xml:"channelId"`
	Title     string `xml:"title"`
	Link      struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
}
