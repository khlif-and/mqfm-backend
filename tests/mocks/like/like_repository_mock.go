package likemock

import (
	"time"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockLikeRepository struct {
	CreateFn        func(like *entity.Like) error
	DeleteFn        func(userID uint, targetType string, targetID uint) error
	FindByUserFn    func(userID uint, targetType string) ([]entity.Like, error)
	ExistsFn        func(userID uint, targetType string, targetID uint) (bool, error)
	CountByTargetFn func(targetType string, targetID uint) (int64, error)
	AggregateFn     func() (map[uint]int64, error)
}

func (m *MockLikeRepository) Create(like *entity.Like) error { return m.CreateFn(like) }
func (m *MockLikeRepository) Delete(userID uint, targetType string, targetID uint) error {
	return m.DeleteFn(userID, targetType, targetID)
}
func (m *MockLikeRepository) FindByUser(userID uint, targetType string) ([]entity.Like, error) {
	return m.FindByUserFn(userID, targetType)
}
func (m *MockLikeRepository) Exists(userID uint, targetType string, targetID uint) (bool, error) {
	return m.ExistsFn(userID, targetType, targetID)
}
func (m *MockLikeRepository) CountByTarget(targetType string, targetID uint) (int64, error) {
	return m.CountByTargetFn(targetType, targetID)
}
func (m *MockLikeRepository) AggregateLikeCounts() (map[uint]int64, error) { return m.AggregateFn() }
func (m *MockLikeRepository) AggregateWeeklyLikeCounts(since time.Time) (map[uint]int64, error) {
	return map[uint]int64{}, nil
}
func (m *MockLikeRepository) AggregateMonthlyLikeCounts(since time.Time) (map[uint]int64, error) {
	return map[uint]int64{}, nil
}

var _ port.LikeRepository = (*MockLikeRepository)(nil)
