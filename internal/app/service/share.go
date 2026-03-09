package service

import (
	"fmt"

	"mqfm-backend/internal/domain/port"
)

type shareService struct {
	baseURL string
}

func NewShareService(baseURL string) port.ShareService {
	return &shareService{baseURL: baseURL}
}

func (s *shareService) GenerateAudioShareLink(audioID uint) string {
	return fmt.Sprintf("%s/audio/%d", s.baseURL, audioID)
}

func (s *shareService) GenerateClipShareLink(shareToken string) string {
	return fmt.Sprintf("%s/clip/%s", s.baseURL, shareToken)
}

func (s *shareService) GeneratePlaylistShareLink(shareToken string) string {
	return fmt.Sprintf("%s/playlist/%s", s.baseURL, shareToken)
}
