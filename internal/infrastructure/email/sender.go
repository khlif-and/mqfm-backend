package email

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"

	"mqfm-backend/internal/infrastructure/config"
	"mqfm-backend/internal/shared/logger"
)

type Sender struct {
	dialer   *gomail.Dialer
	from     string
	fromName string
}

func NewSender(cfg *config.Config) *Sender {
	port, _ := strconv.Atoi(cfg.SMTPPort)
	dialer := gomail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUser, cfg.SMTPPassword)

	return &Sender{
		dialer:   dialer,
		from:     cfg.SMTPUser,
		fromName: cfg.SMTPFromName,
	}
}

func (s *Sender) SendAsync(to string, subject string, body string) {
	go func() {
		m := gomail.NewMessage()
		m.SetAddressHeader("From", s.from, s.fromName)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		if err := s.dialer.DialAndSend(m); err != nil {
			logger.Error(fmt.Sprintf("email send failed to %s: %s", to, err.Error()))
			return
		}

		logger.Info(fmt.Sprintf("email sent to %s", to))
	}()
}
