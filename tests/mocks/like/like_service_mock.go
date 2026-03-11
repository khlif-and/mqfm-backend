package likemock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/dto/request"
)

type MockLikeService struct {
	LikeFn          func(userID uint, req request.LikeRequest) (*entity.Like, error)
	UnlikeFn        func(userID uint, req request.UnlikeRequest) error
	GetLikesFn      func(userID string, targetType string) ([]entity.Like, error)
	CountByTargetFn func(targetType string, targetID uint) (int64, error)
}

func (m *MockLikeService) Like(userID uint, req request.LikeRequest) (*entity.Like, error) {
	return m.LikeFn(userID, req)
}
func (m *MockLikeService) Unlike(userID uint, req request.UnlikeRequest) error {
	return m.UnlikeFn(userID, req)
}
func (m *MockLikeService) GetLikes(userID string, targetType string) ([]entity.Like, error) {
	return m.GetLikesFn(userID, targetType)
}
func (m *MockLikeService) CountByTarget(targetType string, targetID uint) (int64, error) {
	return m.CountByTargetFn(targetType, targetID)
}

var _ port.LikeService = (*MockLikeService)(nil)
