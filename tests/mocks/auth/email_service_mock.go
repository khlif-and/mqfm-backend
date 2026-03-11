package authmock

import "mqfm-backend/internal/domain/port"

type MockEmailService struct {
	SendAsyncFn func(to, subject, body string)
}

func (m *MockEmailService) SendAsync(to, subject, body string) {
	if m.SendAsyncFn != nil {
		m.SendAsyncFn(to, subject, body)
	}
}

var _ port.EmailService = (*MockEmailService)(nil)
