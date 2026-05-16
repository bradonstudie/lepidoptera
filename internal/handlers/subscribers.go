package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
)

type SubscriberHandler struct {
	db     *pgxpool.Pool
	mailer mailer.Mailer
	secret string
}

func NewSubscriberHandler(db *pgxpool.Pool, m mailer.Mailer, secret string) *SubscriberHandler {
	return &SubscriberHandler{db: db, mailer: m, secret: secret}
}

func (h *SubscriberHandler) SubscribePage(w http.ResponseWriter, r *http.Request) {
	// TODO: render subscribe page with templ
	w.Write([]byte("subscribe page — coming soon"))
}

// POST /subscribe
func (h *SubscriberHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		// TODO: return flash fragment via HTMX
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	// TODO:
	// 1. insert subscriber (handle duplicate gracefully)
	// 2. generate confirm token
	// 3. send confirmation email
	// 4. return flash fragment: "check your email to confirm"
	w.Write([]byte("subscribed: " + email))
}

// GET /confirm?token=...
func (h *SubscriberHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	_ = r.URL.Query().Get("token")

	// TODO:
	// 1. find subscriber by scanning tokens (or store token hash in db)
	// 2. set confirmed_at
	// 3. send confirmed email
	// 4. redirect to /?confirmed=1

	http.Redirect(w, r, "/?confirmed=1", http.StatusSeeOther)
}

// GET /unsubscribe?token=...
func (h *SubscriberHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	_ = r.URL.Query().Get("token")

	// TODO:
	// 1. validate token
	// 2. set unsubscribed_at
	// 3. redirect to /?unsubscribed=1

	http.Redirect(w, r, "/?unsubscribed=1", http.StatusSeeOther)
}
