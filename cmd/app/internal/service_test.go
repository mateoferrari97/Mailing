package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type emailClient struct {
	mock.Mock
}

func (e *emailClient) Open() error {
	return e.Mock.Called().Error(0)
}

func (e *emailClient) Auth() error {
	return e.Mock.Called().Error(0)
}

func (e *emailClient) Send(from string, to []string, message string) error {
	return e.Mock.Called(from, to, message).Error(0)
}

func (e *emailClient) Quit() error {
	return e.Mock.Called().Error(0)
}

func TestSendEmail(t *testing.T) {
	// Given
	email := "mateo.ferrari97@gmail.com"
	message := "test"

	c := &emailClient{}
	c.On("Open").Return(nil)
	c.On("Auth").Return(nil)
	c.On("Send", email, []string{email}, message).Return(nil)
	c.On("Quit").Return(nil)

	req := SendEmailRequest{
		To:      email,
		Message: message,
	}

	s := NewService(c)

	// When
	err := s.SendEmail(req)

	// Then
	require.NoError(t, err)
}

func TestSendEmail_OpenError(t *testing.T) {
	// Given
	email := "mateo.ferrari97@gmail.com"
	message := "test"

	c := &emailClient{}
	c.On("Open").Return(errors.New("connection error"))

	req := SendEmailRequest{
		To:      email,
		Message: message,
	}

	s := NewService(c)

	// When
	err := s.SendEmail(req)

	// Then
	require.EqualError(t, err, "connection error")

	c.AssertExpectations(t)
}

func TestSendEmail_AuthError(t *testing.T) {
	// Given
	email := "mateo.ferrari97@gmail.com"
	message := "test"

	c := &emailClient{}
	c.On("Open").Return(nil)
	c.On("Auth").Return(errors.New("authentication error"))

	req := SendEmailRequest{
		To:      email,
		Message: message,
	}

	s := NewService(c)

	// When
	err := s.SendEmail(req)

	// Then
	require.EqualError(t, err, "authentication error")

	c.AssertExpectations(t)
}

func TestSendEmail_SendError(t *testing.T) {
	// Given
	email := "mateo.ferrari97@gmail.com"
	message := "test"

	c := &emailClient{}
	c.On("Open").Return(nil)
	c.On("Auth").Return(nil)
	c.On("Send", email, []string{email}, message).Return(errors.New("sending email error"))

	req := SendEmailRequest{
		To:      email,
		Message: message,
	}

	s := NewService(c)

	// When
	err := s.SendEmail(req)

	// Then
	require.EqualError(t, err, "sending email error")

	c.AssertExpectations(t)
}

func TestSendEmail_QuitError(t *testing.T) {
	// Given
	email := "mateo.ferrari97@gmail.com"
	message := "test"

	c := &emailClient{}
	c.On("Open").Return(nil)
	c.On("Auth").Return(nil)
	c.On("Send", email, []string{email}, message).Return(nil)
	c.On("Quit").Return(errors.New("quitting connection error"))

	req := SendEmailRequest{
		To:      email,
		Message: message,
	}

	s := NewService(c)

	// When
	err := s.SendEmail(req)

	// Then
	require.EqualError(t, err, "quitting connection error")

	c.AssertExpectations(t)
}
