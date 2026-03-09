package service

import (
	"time"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type listeningStatService struct {
	repo port.ListeningStatRepository
}

func NewListeningStatService(repo port.ListeningStatRepository) port.ListeningStatService {
	return &listeningStatService{repo: repo}
}

func (s *listeningStatService) RecordStat(userID uint, req request.RecordStatRequest) error {
	stat := &entity.ListeningStat{
		UserID:          userID,
		AudioID:         req.AudioID,
		DurationSeconds: req.DurationSeconds,
		ListenedAt:      time.Now(),
	}
	return s.repo.Create(stat)
}

func (s *listeningStatService) GetWeeklySummary(userID uint) (int, error) {
	return s.repo.GetWeeklySummary(userID)
}

func (s *listeningStatService) GetMonthlySummary(userID uint) (int, error) {
	return s.repo.GetMonthlySummary(userID)
}

func (s *listeningStatService) GetTopCategories(userID uint, limit int) ([]port.CategoryStat, error) {
	return s.repo.GetTopCategories(userID, limit)
}

func (s *listeningStatService) GetTopArtists(userID uint, limit int) ([]port.ArtistStat, error) {
	return s.repo.GetTopArtists(userID, limit)
}

func (s *listeningStatService) GetDailySummary(userID uint, days int) ([]port.DailyStat, error) {
	return s.repo.GetDailySummary(userID, days)
}

func (s *listeningStatService) GetRecap(userID uint) (*port.ListeningRecap, error) {
	weekly, _ := s.repo.GetWeeklySummary(userID)
	monthly, _ := s.repo.GetMonthlySummary(userID)
	categories, _ := s.repo.GetTopCategories(userID, 5)
	artists, _ := s.repo.GetTopArtists(userID, 5)
	daily, _ := s.repo.GetDailySummary(userID, 30)

	return &port.ListeningRecap{
		WeeklyMinutes:  weekly / 60,
		MonthlyMinutes: monthly / 60,
		TopCategories:  categories,
		TopArtists:     artists,
		DailyStats:     daily,
	}, nil
}
