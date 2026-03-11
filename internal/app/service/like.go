package service

import (
	"errors"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/constant"
	"mqfm-backend/internal/shared/dto/request"
)

type likeService struct {
	repo        port.LikeRepository
	downloadSvc port.DownloadService
	prefRepo    port.UserPreferenceRepository
}

func NewLikeService(repo port.LikeRepository, downloadSvc port.DownloadService, prefRepo port.UserPreferenceRepository) port.LikeService {
	return &likeService{repo: repo, downloadSvc: downloadSvc, prefRepo: prefRepo}
}

func (s *likeService) Like(userID uint, req request.LikeRequest) (*entity.Like, error) {
	exists, err := s.repo.Exists(userID, req.TargetType, req.TargetID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constant.MsgAlreadyLiked)
	}

	like := entity.Like{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}

	if err := s.repo.Create(&like); err != nil {
		return nil, err
	}

	if req.TargetType == "audio" {
		go s.triggerAutoDownload(userID, req.TargetID)
	}

	return &like, nil
}

func (s *likeService) triggerAutoDownload(userID, audioID uint) {
	pref, err := s.prefRepo.FindByUser(userID)
	if err != nil || !pref.AutoDownloadOnLike {
		return
	}
	_, _ = s.downloadSvc.RecordDownload(userID, request.DownloadRequest{AudioID: audioID})
}

func (s *likeService) Unlike(userID uint, req request.UnlikeRequest) error {
	return s.repo.Delete(userID, req.TargetType, req.TargetID)
}

func (s *likeService) GetLikes(userID string, targetType string) ([]entity.Like, error) {
	var uid uint
	for _, c := range userID {
		uid = uid*10 + uint(c-'0')
	}
	return s.repo.FindByUser(uid, targetType)
}

func (s *likeService) CountByTarget(targetType string, targetID uint) (int64, error) {
	return s.repo.CountByTarget(targetType, targetID)
}
