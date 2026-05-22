package handlers

import (
	"errors"
	"net/http"

	"github.com/lepidoptera/lepidoptera/internal/service"
	"github.com/lepidoptera/lepidoptera/web/components"
	"github.com/lepidoptera/lepidoptera/web/pages"
)

type SubscriberHandler struct {
	subscriberService *service.SubscriberService
}

func NewSubscriberHandler(subscriberService *service.SubscriberService) *SubscriberHandler {
	return &SubscriberHandler{subscriberService: subscriberService}
}

func (h *SubscriberHandler) SubscribePage(w http.ResponseWriter, r *http.Request) {
	pages.Subscribe().Render(r.Context(), w)
}

// POST /subscribe
func (h *SubscriberHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		components.Flash("please enter a valid email.", true).Render(r.Context(), w)
		return
	}

	err := h.subscriberService.Subscribe(r.Context(), email)
	if errors.Is(err, service.ErrAlreadySubscribed) {
		components.Flash(err.Error(), false).Render(r.Context(), w)
		return
	}
	if err != nil {
		components.Flash(err.Error(), true).Render(r.Context(), w)
		return
	}

	components.Flash("check your email to confirm your subscription.", false).Render(r.Context(), w)
}

// GET /confirm?token=...
func (h *SubscriberHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	err := h.subscriberService.Confirm(r.Context(), token)
	if errors.Is(err, service.ErrInvalidToken) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/?confirmed=1", http.StatusSeeOther)
}

// GET /unsubscribe?token=...
func (h *SubscriberHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	err := h.subscriberService.Unsubscribe(r.Context(), token)
	if errors.Is(err, service.ErrInvalidToken) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/?unsubscribed=1", http.StatusSeeOther)
}
