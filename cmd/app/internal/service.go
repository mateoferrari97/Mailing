package internal

import (
	"os"
)

type EmailClient interface {
	Open() error
	Auth() error
	Send(from string, to []string, message string) error
	Quit() error
}

type Service struct {
	client EmailClient
}

func NewService(cli EmailClient) *Service {
	return &Service{client: cli}
}

func (s *Service) SendEmail(req SendEmailRequest) error {
	if err := s.client.Open(); err != nil {
		return err
	}

	if err := s.client.Auth(); err != nil {
		return err
	}

	from := os.Getenv("EMAIL")
	to := []string{req.To}
	if err := s.client.Send(from, to, req.Message); err != nil {
		return err
	}

	return s.client.Quit()
}

