package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/logger"
)

type youTubeResponse struct {
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

type livestreamService struct {
	repo   port.LivestreamRepository
	apiKey string
}

func NewLivestreamService(repo port.LivestreamRepository, apiKey string) port.LivestreamService {
	return &livestreamService{repo: repo, apiKey: apiKey}
}

func (s *livestreamService) UpdateLiveStatus(channelID string) error {
	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&eventType=live&type=video&key=%s",
		channelID, s.apiKey,
	)

	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		logger.Error("youtube api connection failed", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	var ytResp youTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ytResp); err != nil {
		logger.Error("youtube api decode failed", zap.Error(err))
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

	existing, err := s.repo.FindFirst()
	if err != nil || existing == nil {
		ls := &entity.LiveStream{
			IsLive:      isLive,
			VideoID:     videoID,
			Title:       title,
			Thumbnail:   thumbnail,
			LastChecked: time.Now(),
		}
		if createErr := s.repo.Create(ls); createErr != nil {
			logger.Error("livestream db create failed", zap.Error(createErr))
			return createErr
		}
	} else {
		existing.IsLive = isLive
		existing.VideoID = videoID
		existing.Title = title
		existing.Thumbnail = thumbnail
		existing.LastChecked = time.Now()

		if saveErr := s.repo.Save(existing); saveErr != nil {
			logger.Error("livestream db update failed", zap.Error(saveErr))
			return saveErr
		}
	}

	status := "OFFLINE"
	if isLive {
		status = "LIVE"
	}
	logger.Info("youtube status updated",
		zap.String("status", status),
		zap.String("title", title),
		zap.String("video_id", videoID),
	)

	return nil
}

func (s *livestreamService) GetStatus() (*entity.LiveStream, error) {
	return s.repo.FindFirst()
}
