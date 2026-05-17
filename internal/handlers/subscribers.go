package handlers

import (
	"net/http"

	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
	"github.com/lepidoptera/lepidoptera/web/pages"
)

type SubscriberHandler struct {
	queries *db.Queries
	mailer  mailer.Mailer
	secret  string
}

func NewSubscriberHandler(queries *db.Queries, m mailer.Mailer, secret string) *SubscriberHandler {
	return &SubscriberHandler{queries: queries, mailer: m, secret: secret}
}

func (h *SubscriberHandler) SubscribePage(w http.ResponseWriter, r *http.Request) {
	pages.Subscribe().Render(r.Context(), w)
}

// POST /subscribe
func (h *SubscriberHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		// TODO: return flash fragment via HTMX
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	subscriber, err := h.queries.CreateSubscriber(r.Context(), email)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	if subscriber.ID.String() == "00000000-0000-0000-0000-000000000000" {
		w.Write([]byte("you're already on the list."))
		return
	}

	token := mailer.ConfirmToken(subscriber.ID, h.secret)
	confirmEmail := mailer.ConfirmSubscriptionEmail(token)
	confirmEmail.To = email
	h.mailer.Send(r.Context(), confirmEmail)

	w.Write([]byte("check your email to confirm your subscription"))
}

// GET /confirm?token=...
func (h *SubscriberHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	_ = r.URL.Query().Get("token")

	// TODO:
	// 1. find subscriber by scanning tokens (or store token hash in db)
	// 2. set confirmed_at
	// 3. send confirmed email

	http.Redirect(w, r, "/?confirmed=1", http.StatusSeeOther)
}

// GET /unsubscribe?token=...
func (h *SubscriberHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	_ = r.URL.Query().Get("token")

	// TODO: validate token, set unsubscribed_at

	http.Redirect(w, r, "/?unsubscribed=1", http.StatusSeeOther)
}
