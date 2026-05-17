package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lepidoptera/lepidoptera/internal/auth"
	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
	adminpages "github.com/lepidoptera/lepidoptera/web/pages/admin"
)

type AdminHandler struct {
	queries *db.Queries
	mailer  mailer.Mailer
	secret  string
}

func NewAdminHandler(queries *db.Queries, m mailer.Mailer, secret string) *AdminHandler {
	return &AdminHandler{queries: queries, mailer: m, secret: secret}
}

// GET /admin
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	shows, err := h.queries.ListAllShows(r.Context())
	if err != nil {
		http.Error(w, "error loading shows", http.StatusInternalServerError)
		return
	}

	adminpages.Dashboard(shows).Render(r.Context(), w)
}

// GET /admin/login
func (h *AdminHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	adminpages.Login("").Render(r.Context(), w)
}

// POST /admin/login
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	adminEmail := os.Getenv("ADMIN_EMAIL")

	message := "A login link is on its way"

	if email == adminEmail {
		token := auth.GenerateLoginToken(email, h.secret)
		log.Printf("DEBUG login token: /admin/verify?token=%s", token) // TODO: Get this out of here when Resend is configured
		loginEmail := mailer.AdminLoginEmail(token)
		loginEmail.To = email
		h.mailer.Send(r.Context(), loginEmail)
	}

	adminpages.Login(message).Render(r.Context(), w)
}

// POST /admin/verify
func (h *AdminHandler) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email, err := auth.ValidateLoginToken(token, h.secret)
	if err != nil {
		http.Error(w, "invalid or expired login link", http.StatusUnauthorized)
		return
	}

	if email != os.Getenv("ADMIN_EMAIL") {
		http.Error(w, "invalid or expired login link", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    os.Getenv("ADMIN_SESSION_SECRET"),
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// POST /admin/logout
func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// GET /admin/shows/new
func (h *AdminHandler) NewShowForm(w http.ResponseWriter, r *http.Request) {
	shows, err := h.queries.ListAllShows(r.Context())
	if err != nil {
		http.Error(w, "error loading shows", http.StatusInternalServerError)
		return
	}

	_ = shows
	w.Write([]byte("admin dashboard - coming soon"))
}

// POST /admin/shows
func (h *AdminHandler) CreateShow(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. parse form values
	// 2. insert show as draft (is_published=false)
	// 3. redirect to admin dashboard
}

// POST /admin/shows/{id}/publish
func (h *AdminHandler) PublishShow(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "id")

	// TODO:
	// 1. set is_published=true, published_at=now()
	// 2. insert pending notifications for all confirmed subscribers
	// 3. worker picks them up within the hour
	// 4. redirect to admin dashboard
}

// GET /admin/bands/new
func (h *AdminHandler) NewBandForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("new band form — coming soon"))
}

// POST /admin/bands
func (h *AdminHandler) CreateBand(w http.ResponseWriter, r *http.Request) {
	// TODO: parse form, insert band, redirect
}

// GET /admin/venues/new
func (h *AdminHandler) NewVenueForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("new venue form — coming soon"))
}

// POST /admin/venues
func (h *AdminHandler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	// TODO: parse form, insert venue, redirect
}
