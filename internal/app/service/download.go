package service

import (
	"errors"
	"time"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type downloadService struct {
	repo          port.DownloadRepository
	audioRepo     port.AudioRepository
	favArtistRepo port.FavoriteArtistRepository
}

func NewDownloadService(
	repo port.DownloadRepository,
	audioRepo port.AudioRepository,
	favArtistRepo port.FavoriteArtistRepository,
) port.DownloadService {
	return &downloadService{repo: repo, audioRepo: audioRepo, favArtistRepo: favArtistRepo}
}

func (s *downloadService) RecordDownload(userID uint, req request.DownloadRequest) (*entity.Download, error) {
	exists, _ := s.repo.Exists(userID, req.AudioID)
	if exists {
		return nil, errors.New("audio already downloaded")
	}

	audio, err := s.audioRepo.FindByID(req.AudioID)
	if err != nil {
		return nil, errors.New("audio not found")
	}

	download := &entity.Download{
		UserID:   userID,
		AudioID:  req.AudioID,
		FileSize: audio.FileSize,
	}
	if req.FileSize > 0 {
		download.FileSize = req.FileSize
	}

	if err := s.repo.Create(download); err != nil {
		return nil, err
	}
	return download, nil
}

func (s *downloadService) GetDownloads(userID uint) ([]entity.Download, error) {
	return s.repo.FindByUser(userID)
}

func (s *downloadService) DeleteDownload(id, userID uint) error {
	return s.repo.Delete(id, userID)
}

func (s *downloadService) GetStorageUsage(userID uint) (int64, error) {
	return s.repo.SumSizeByUser(userID)
}

func (s *downloadService) GetNewFromFavorites(userID uint) ([]entity.Audio, error) {
	favs, err := s.favArtistRepo.FindByUser(userID)
	if err != nil || len(favs) == 0 {
		return nil, nil
	}

	artists := make([]string, len(favs))
	for i, f := range favs {
		artists[i] = f.ArtistName
	}

	since := time.Now().AddDate(0, 0, -7)
	return s.audioRepo.FindNewByArtists(artists, since)
}
