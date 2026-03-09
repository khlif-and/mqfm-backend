package service_test

import (
"errors"
"testing"

"github.com/stretchr/testify/assert"

"mqfm-backend/internal/app/service"
"mqfm-backend/internal/domain/entity"
"mqfm-backend/internal/shared/dto/request"
"mqfm-backend/tests/mocks"
)

func TestLike_Success(t *testing.T) {
repo := &mocks.MockLikeRepository{
ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) { return false, nil },
CreateFn: func(like *entity.Like) error { like.ID = 1; return nil },
}
svc := service.NewLikeService(repo)
like, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})
assert.NoError(t, err)
assert.NotNil(t, like)
assert.Equal(t, uint(1), like.UserID)
assert.Equal(t, uint(10), like.TargetID)
assert.Equal(t, "audio", like.TargetType)
}

func TestLike_AlreadyLiked(t *testing.T) {
repo := &mocks.MockLikeRepository{
ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) { return true, nil },
}
svc := service.NewLikeService(repo)
like, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})
assert.Error(t, err)
assert.Nil(t, like)
}

func TestLike_DbError(t *testing.T) {
repo := &mocks.MockLikeRepository{
ExistsFn: func(userID uint, targetType string, targetID uint) (bool, error) {
return false, errors.New("db error")
},
}
svc := service.NewLikeService(repo)
like, err := svc.Like(1, request.LikeRequest{TargetType: "audio", TargetID: 10})
assert.Error(t, err)
assert.Nil(t, like)
}

func TestUnlike_Success(t *testing.T) {
repo := &mocks.MockLikeRepository{
DeleteFn: func(userID uint, targetType string, targetID uint) error { return nil },
}
svc := service.NewLikeService(repo)
err := svc.Unlike(1, request.UnlikeRequest{TargetType: "audio", TargetID: 10})
assert.NoError(t, err)
}

func TestGetLikes_Success(t *testing.T) {
repo := &mocks.MockLikeRepository{
FindByUserFn: func(userID uint, targetType string) ([]entity.Like, error) {
return []entity.Like{
{ID: 1, UserID: userID, TargetType: targetType, TargetID: 10},
{ID: 2, UserID: userID, TargetType: targetType, TargetID: 20},
}, nil
},
}
svc := service.NewLikeService(repo)
likes, err := svc.GetLikes("1", "audio")
assert.NoError(t, err)
assert.Len(t, likes, 2)
}

func TestGetLikes_Empty(t *testing.T) {
repo := &mocks.MockLikeRepository{
FindByUserFn: func(userID uint, targetType string) ([]entity.Like, error) {
return []entity.Like{}, nil
},
}
svc := service.NewLikeService(repo)
likes, err := svc.GetLikes("1", "audio")
assert.NoError(t, err)
assert.Empty(t, likes)
}

func TestCountByTarget_Success(t *testing.T) {
repo := &mocks.MockLikeRepository{
CountByTargetFn: func(targetType string, targetID uint) (int64, error) { return 42, nil },
}
svc := service.NewLikeService(repo)
count, err := svc.CountByTarget("audio", 10)
assert.NoError(t, err)
assert.Equal(t, int64(42), count)
}
