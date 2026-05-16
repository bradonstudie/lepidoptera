package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
)

type AdminHandler struct {
	db     *pgxpool.Pool
	mailer mailer.Mailer
}

func NewAdminHandler(db *pgxpool.Pool, m mailer.Mailer) *AdminHandler {
	return &AdminHandler{db: db, mailer: m}
}

// GET /admin
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// TODO: list all shows (published + drafts), render admin dashboard
	w.Write([]byte("admin dashboard — coming soon"))
}

// GET /admin/shows/new
func (h *AdminHandler) NewShowForm(w http.ResponseWriter, r *http.Request) {
	// TODO: render show creation form
	w.Write([]byte("new show form — coming soon"))
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
