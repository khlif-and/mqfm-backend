package livestream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	lsModel "mqfm-backend/internal/models/livestream"
	"mqfm-backend/internal/utils"

)

type YouTubeResponse struct {
	Items []struct {
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title      string `json:"title"`
			Thumbnails struct {
				High struct {
					Url string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

type LiveStreamService struct {
	db     *gorm.DB
	apiKey string
}

func NewLiveStreamService(db *gorm.DB, apiKey string) *LiveStreamService {
	return &LiveStreamService{db: db, apiKey: apiKey}
}

func (s *LiveStreamService) UpdateLiveStatus(channelID string) error {
	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&eventType=live&type=video&key=%s",
		channelID, s.apiKey,
	)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		utils.Log.Error("[YouTube Service] Connection failed", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	var ytResp YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ytResp); err != nil {
		utils.Log.Error("[YouTube Service] JSON decode failed", zap.Error(err))
		return err
	}

	isLive := len(ytResp.Items) > 0
	var videoID, title, thumbnail string

	if isLive {
		item := ytResp.Items[0]
		videoID = item.Id.VideoId
		title = item.Snippet.Title
		thumbnail = item.Snippet.Thumbnails.High.Url
	}

	var liveStream lsModel.LiveStream

	if err := s.db.Limit(1).Find(&liveStream).Error; err != nil {
		utils.Log.Error("[YouTube Service] Database read failed", zap.Error(err))
		return err
	}

	if liveStream.ID == 0 {
		liveStream = lsModel.LiveStream{
			IsLive:      isLive,
			VideoID:     videoID,
			Title:       title,
			Thumbnail:   thumbnail,
			LastChecked: time.Now(),
		}
		if err := s.db.Create(&liveStream).Error; err != nil {
			utils.Log.Error("[YouTube Service] Database create failed", zap.Error(err))
			return err
		}
	} else {
		liveStream.IsLive = isLive
		liveStream.VideoID = videoID
		liveStream.Title = title
		liveStream.Thumbnail = thumbnail
		liveStream.LastChecked = time.Now()

		if err := s.db.Save(&liveStream).Error; err != nil {
			utils.Log.Error("[YouTube Service] Database update failed", zap.Error(err))
			return err
		}
	}

	statusLog := "OFFLINE"
	if isLive {
		statusLog = "LIVE"
	}

	utils.Log.Info("[YouTube Service] Status Updated",
		zap.String("status", statusLog),
		zap.String("video_title", title),
		zap.String("video_id", videoID),
	)

	return nil
}

func (s *LiveStreamService) GetStatus() (lsModel.LiveStream, error) {
	var liveStream lsModel.LiveStream
	err := s.db.First(&liveStream).Error
	return liveStream, err
}