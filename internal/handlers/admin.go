package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bradonstudie/lepidoptera/internal/auth"
	db "github.com/bradonstudie/lepidoptera/internal/db/generated"
	"github.com/bradonstudie/lepidoptera/internal/mailer"
	"github.com/bradonstudie/lepidoptera/internal/service"
	"github.com/bradonstudie/lepidoptera/internal/timeutil"
	adminpages "github.com/bradonstudie/lepidoptera/web/pages/admin"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService *service.AdminService
	queries      *db.Queries
	secret       string
}

func NewAdminHandler(adminService *service.AdminService, queries *db.Queries, secret string) *AdminHandler {
	return &AdminHandler{adminService: adminService, queries: queries, secret: secret}
}

// GET /admin/login
func (h *AdminHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	adminpages.Login("").Render(r.Context(), w)
}

// POST /admin/login
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	adminEmail := os.Getenv("ADMIN_EMAIL")

	message := "if that email is authorized, a login link is on its way."

	if email == adminEmail {
		token := auth.GenerateLoginToken(email, h.secret)
		loginEmail := mailer.AdminLoginEmail(token)
		loginEmail.To = email
		if err := h.adminService.Mailer().Send(r.Context(), loginEmail); err != nil {
			log.Printf("admin login email failed: %v", err)
		}
	}

	adminpages.Login(message).Render(r.Context(), w)
}

// GET /admin/verify?token=...
func (h *AdminHandler) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	email, err := auth.ValidateLoginToken(token, h.secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if email != os.Getenv("ADMIN_EMAIL") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_, err = h.queries.MarkLoginTokenUsed(r.Context(), auth.HashSessionToken(token))
	if err == sql.ErrNoRows {
		http.Error(w, "token already used", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	sessionToken, err := auth.GenerateSessionToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.queries.CreateAdminSession(r.Context(), db.CreateAdminSessionParams{
		TokenHash: auth.HashSessionToken(sessionToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// POST /admin/logout
func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_session")
	if err == nil {
		if err := h.queries.RevokeAdminSession(r.Context(), auth.HashSessionToken(cookie.Value)); err != nil {
			log.Printf("failed to revoke admin session: %v", err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// GET /admin
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	shows, err := h.adminService.ListAllShows(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adminpages.Dashboard(shows).Render(r.Context(), w)
}

// GET /admin/shows/new
func (h *AdminHandler) NewShowForm(w http.ResponseWriter, r *http.Request) {
	venues, err := h.adminService.ListVenues(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bands, err := h.adminService.ListBands(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adminpages.ShowForm(venues, bands, "").Render(r.Context(), w)
}

// POST /admin/shows
func (h *AdminHandler) CreateShow(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	dateStr := r.FormValue("date")

	if title == "" || dateStr == "" {
		h.reloadShowForm(w, r, "title and date are required.")
		return
	}

	date, err := time.Parse(timeutil.DateTimeLocal, dateStr)
	if err != nil {
		h.reloadShowForm(w, r, "invalid date format.")
		return
	}

	venueID := r.FormValue("venue_id")
	var nullVenueID uuid.NullUUID
	if venueID != "" {
		id, err := uuid.Parse(venueID)
		if err == nil {
			nullVenueID = uuid.NullUUID{UUID: id, Valid: true}
		}
	}

	err = h.adminService.CreateShow(r.Context(), db.CreateShowParams{
		Title:       title,
		Slug:        h.adminService.GenerateSlug(title, date),
		Date:        date,
		VenueID:     nullVenueID,
		Description: sql.NullString{String: r.FormValue("description"), Valid: r.FormValue("description") != ""},
		TicketUrl:   sql.NullString{String: r.FormValue("ticket_url"), Valid: r.FormValue("ticket_url") != ""},
	}, r.Form["band_ids"], r.FormValue("headliner_id"))
	if err != nil {
		h.reloadShowForm(w, r, "something went wrong, please try again.")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// POST /admin/shows/{id}/publish
func (h *AdminHandler) PublishShow(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.adminService.PublishShow(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrShowNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// GET /admin/bands/new
func (h *AdminHandler) NewBandForm(w http.ResponseWriter, r *http.Request) {
	adminpages.BandForm("").Render(r.Context(), w)
}

// POST /admin/bands
func (h *AdminHandler) CreateBand(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		adminpages.BandForm("band name is required.").Render(r.Context(), w)
		return
	}

	err := h.adminService.CreateBand(r.Context(), db.CreateBandParams{
		Name:        name,
		Genre:       sql.NullString{String: r.FormValue("genre"), Valid: r.FormValue("genre") != ""},
		Description: sql.NullString{String: r.FormValue("description"), Valid: r.FormValue("description") != ""},
		WebsiteUrl:  sql.NullString{String: r.FormValue("website_url"), Valid: r.FormValue("website_url") != ""},
	})
	if err != nil {
		adminpages.BandForm(err.Error()).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// GET /admin/venues/new
func (h *AdminHandler) NewVenueForm(w http.ResponseWriter, r *http.Request) {
	adminpages.VenueForm("").Render(r.Context(), w)
}

// POST /admin/venues
func (h *AdminHandler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		adminpages.VenueForm("venue name is required.").Render(r.Context(), w)
		return
	}

	capacity := r.FormValue("capacity")
	var nullCapacity sql.NullInt32
	if capacity != "" {
		var cap int64
		if _, err := fmt.Sscanf(capacity, "%d", &cap); err == nil {
			nullCapacity = sql.NullInt32{Int32: int32(cap), Valid: true}
		}
	}

	err := h.adminService.CreateVenue(r.Context(), db.CreateVenueParams{
		Name:     name,
		Address:  sql.NullString{String: r.FormValue("address"), Valid: r.FormValue("address") != ""},
		City:     sql.NullString{String: r.FormValue("city"), Valid: r.FormValue("city") != ""},
		State:    sql.NullString{String: r.FormValue("state"), Valid: r.FormValue("state") != ""},
		Capacity: nullCapacity,
	})
	if err != nil {
		adminpages.VenueForm(err.Error()).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// helper to reload show form with error
func (h *AdminHandler) reloadShowForm(w http.ResponseWriter, r *http.Request, errMsg string) {
	venues, _ := h.adminService.ListVenues(r.Context())
	bands, _ := h.adminService.ListBands(r.Context())
	adminpages.ShowForm(venues, bands, errMsg).Render(r.Context(), w)
}
