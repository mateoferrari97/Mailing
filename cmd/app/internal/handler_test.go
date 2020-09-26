package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mateoferrari97/Kit/web"
	"github.com/stretchr/testify/require"
)

func TestHandler_SendEmail(t *testing.T) {
	// Given
	w := web.NewServer()
	h := NewHandler(w)

	h.RouteSendEmail(func(_ SendEmailRequest) error {
		return nil
	})

	b := []byte(`{
		"to": "failToImprove@gmail.com",
		"message": "fail to improve"
	}`)

	// When
	ts := httptest.NewServer(w.Router)
	defer ts.Close()

	resp, err := http.Post(fmt.Sprintf("%s/email/me", ts.URL), "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Then
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHandler_SendEmail_HandlerError(t *testing.T) {
	// Given
	w := web.NewServer()
	h := NewHandler(w)

	h.RouteSendEmail(func(_ SendEmailRequest) error {
		return errors.New("internal server error")
	})

	b := []byte(`{
		"to": "failToImprove@gmail.com",
		"message": "fail to improve"
	}`)

	// When
	ts := httptest.NewServer(w.Router)
	defer ts.Close()

	resp, err := http.Post(fmt.Sprintf("%s/email/me", ts.URL), "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	var r struct {
		Message string `json:"message"`
	}

	_ = json.NewDecoder(resp.Body).Decode(&r)

	// Then
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Equal(t, r.Message, "internal server error")
}

func TestHandler_SendEmail_UnprocessableEntityError(t *testing.T) {
	tt := []struct {
		name          string
		b             []byte
		expectedError string
	}{
		{
			name: "malformed message",
			b: []byte(`{
				"to": "failToImprove@gmail.com",
				"message": here is the error
			}`),
			expectedError: "decoding request:unprocessable entity invalid character 'h' looking for beginning of value",
		},
		{
			name: "empty destination",
			b: []byte(`{
				"to": "",
				"message": "hi"
			}`),
			expectedError: "validating request:unprocessable entity Key: 'SendEmailRequest.To' Error:Field validation for 'To' failed on the 'required' tag",
		},

		{
			name: "invalid email",
			b: []byte(`{
				"to": "invalidEmail.com",
				"message": "hi"
			}`),
			expectedError: "validating request:unprocessable entity Key: 'SendEmailRequest.To' Error:Field validation for 'To' failed on the 'email' tag",
		},

		{
			name: "empty message",
			b: []byte(`{
				"to": "failToImprove@gmail.com",
				"message": ""
			}`),
			expectedError: "validating request:unprocessable entity Key: 'SendEmailRequest.Message' Error:Field validation for 'Message' failed on the 'required' tag",
		},
		{
			name: "message length exceeded",
			b: []byte(`{
				"to": "failToImprove@gmail.com",
				"message": "asasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsaasasddsasaaddsasasddsasaaddsasasddsasaaddsasasddsasaaddsasasda"
			}`),
			expectedError: "validating request:unprocessable entity Key: 'SendEmailRequest.Message' Error:Field validation for 'Message' failed on the 'max' tag",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			w := web.NewServer()
			h := NewHandler(w)

			h.RouteSendEmail(func(_ SendEmailRequest) error {
				return nil
			})

			// When
			ts := httptest.NewServer(w.Router)
			defer ts.Close()

			resp, err := http.Post(fmt.Sprintf("%s/email/me", ts.URL), "application/json", bytes.NewReader(tc.b))
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			var r struct {
				Message string `json:"message"`
			}

			_ = json.NewDecoder(resp.Body).Decode(&r)

			// Then
			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
			require.Equal(t, r.Message, tc.expectedError)
		})
	}
}
