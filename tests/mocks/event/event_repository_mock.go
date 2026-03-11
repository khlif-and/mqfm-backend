package eventmock

import (
	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
)

type MockEventRepository struct {
	CreateFn           func(event *entity.Event) error
	FindAllFn          func() ([]entity.Event, error)
	FindByIDFn         func(id uint) (*entity.Event, error)
	UpdateFn           func(id uint, updates map[string]interface{}) error
	DeleteFn           func(id uint) error
	FindUpcomingFn     func(limit int) ([]entity.Event, error)
	CreateRSVPFn       func(rsvp *entity.EventRSVP) error
	DeleteRSVPFn       func(userID, eventID uint) error
	FindRSVPsByEventFn func(eventID uint) ([]entity.EventRSVP, error)
	FindRSVPsByUserFn  func(userID uint) ([]entity.EventRSVP, error)
	ExistsRSVPFn       func(userID, eventID uint) (bool, error)
	CountRSVPFn        func(eventID uint) (int64, error)
}

func (m *MockEventRepository) Create(event *entity.Event) error     { return m.CreateFn(event) }
func (m *MockEventRepository) FindAll() ([]entity.Event, error)     { return m.FindAllFn() }
func (m *MockEventRepository) FindByID(id uint) (*entity.Event, error) { return m.FindByIDFn(id) }
func (m *MockEventRepository) Update(id uint, updates map[string]interface{}) error {
	return m.UpdateFn(id, updates)
}
func (m *MockEventRepository) Delete(id uint) error                  { return m.DeleteFn(id) }
func (m *MockEventRepository) FindUpcoming(limit int) ([]entity.Event, error) {
	return m.FindUpcomingFn(limit)
}
func (m *MockEventRepository) CreateRSVP(rsvp *entity.EventRSVP) error { return m.CreateRSVPFn(rsvp) }
func (m *MockEventRepository) DeleteRSVP(userID, eventID uint) error {
	return m.DeleteRSVPFn(userID, eventID)
}
func (m *MockEventRepository) FindRSVPsByEvent(eventID uint) ([]entity.EventRSVP, error) {
	return m.FindRSVPsByEventFn(eventID)
}
func (m *MockEventRepository) FindRSVPsByUser(userID uint) ([]entity.EventRSVP, error) {
	return m.FindRSVPsByUserFn(userID)
}
func (m *MockEventRepository) ExistsRSVP(userID, eventID uint) (bool, error) {
	return m.ExistsRSVPFn(userID, eventID)
}
func (m *MockEventRepository) CountRSVP(eventID uint) (int64, error) {
	return m.CountRSVPFn(eventID)
}

var _ port.EventRepository = (*MockEventRepository)(nil)
