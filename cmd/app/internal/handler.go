package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/mateoferrari97/Kit/web"
)

const (
	getPing       = "/ping"
	postSendEmail = "/email/me"
)

var _v = validator.New()

type Wrapper interface {
	Wrap(method string, pattern string, handler web.HandlerFunc)
}

type Handler struct {
	Wrapper
}

func NewHandler(wrapper Wrapper) *Handler {
	return &Handler{wrapper}
}

func (h *Handler) Ping() {
	wrapH := func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprint(w, "pong")

		return err
	}

	h.Wrap(http.MethodGet, getPing, wrapH)
}

type SendEmailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Message string `json:"message" validate:"required,min=1,max=1000"`
}

type SendEmailHandler func(req SendEmailRequest) error

func (h *Handler) RouteSendEmail(handler SendEmailHandler) {
	wrapH := func(w http.ResponseWriter, r *http.Request) error {
		var req SendEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return fmt.Errorf("decoding request:%w %v", web.ErrUnprocessableEntity, err)
		}

		if err := _v.Struct(req); err != nil {
			return fmt.Errorf("validating request:%w %v", web.ErrUnprocessableEntity, err)
		}

		if err := handler(req); err != nil {
			return err
		}

		return web.RespondJSON(w, nil, http.StatusNoContent)
	}

	h.Wrap(http.MethodPost, postSendEmail, wrapH)
}
