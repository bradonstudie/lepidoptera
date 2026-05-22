package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lepidoptera/lepidoptera/internal/auth"
	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
	"github.com/lepidoptera/lepidoptera/internal/timeutil"
	"github.com/lepidoptera/lepidoptera/internal/viewmodels"
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
	rows, err := h.queries.ListAllShows(r.Context())
	if err != nil {
		http.Error(w, "error loading shows", http.StatusInternalServerError)
		return
	}

	shows := viewmodels.NewAdminShowViewModels(rows)
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

// GET /admin/verify?token={token}
func (h *AdminHandler) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	email, err := auth.ValidateLoginToken(token, h.secret)
	if err != nil {
		http.Error(w, "invalid or expired login link", http.StatusUnauthorized)
		return
	}

	if email != os.Getenv("ADMIN_EMAIL") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// generate a random session token
	sessionToken, err := auth.GenerateSessionToken()
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	// store the hash in the database
	_, err = h.queries.CreateAdminSession(r.Context(), db.CreateAdminSessionParams{
		TokenHash: auth.HashSessionToken(sessionToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	// give the raw token to the client
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    sessionToken,
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
	cookie, err := r.Cookie("admin_session")
	if err == nil {
		// revoke the session in the database
		h.queries.RevokeAdminSession(r.Context(), auth.HashSessionToken(cookie.Value))
	}

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
	venues, err := h.queries.ListVenues(r.Context())
	if err != nil {
		http.Error(w, "error loading venues", http.StatusInternalServerError)
		return
	}

	bands, err := h.queries.ListBands(r.Context())
	if err != nil {
		http.Error(w, "error loading bands", http.StatusInternalServerError)
		return
	}

	adminpages.ShowForm(venues, bands, "").Render(r.Context(), w)
}

// POST /admin/shows
func (h *AdminHandler) CreateShow(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	dateStr := r.FormValue("date")

	if title == "" || dateStr == "" {
		h.reloadShowForm(w, r, "title, slug, and date are required")
		return
	}

	date, err := time.Parse(timeutil.DateTimeLocal, dateStr)
	if err != nil {
		h.reloadShowForm(w, r, "invalid date format")
		return
	}

	venueID := r.FormValue("venue_id")
	var nullVenueId uuid.NullUUID
	if venueID != "" {
		id, err := uuid.Parse(venueID)
		if err == nil {
			nullVenueId = uuid.NullUUID{UUID: id, Valid: true}
		}
	}

	slug := slugify(title, date)

	show, err := h.queries.CreateShow(r.Context(), db.CreateShowParams{
		Title:       title,
		Slug:        slug,
		Date:        date,
		VenueID:     nullVenueId,
		Description: sql.NullString{String: r.FormValue("description"), Valid: r.FormValue("description") != ""},
		TicketUrl:   sql.NullString{String: r.FormValue("ticket_url"), Valid: r.FormValue("ticket_url") != ""},
	})
	if err != nil {
		h.reloadShowForm(w, r, "something went wrong, please try again")
	}

	headlinerId := r.FormValue("headliner_id")
	bandIDs := r.Form["band_ids"]

	for i, bandIDstr := range bandIDs {
		bandID, err := uuid.Parse(bandIDstr)
		if err != nil {
			continue
		}

		isHeadliner := bandIDstr == headlinerId
		h.queries.AddBandToShow(r.Context(), db.AddBandToShowParams{
			ShowID:       show.ID,
			BandID:       bandID,
			BillingOrder: int32(i),
			IsHeadliner:  isHeadliner,
		})
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// POST /admin/shows/{id}/publish
func (h *AdminHandler) PublishShow(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "error publishing show", http.StatusInternalServerError)
		return
	}

	_, err = h.queries.PublishShow(r.Context(), id)
	if err != nil {
		http.Error(w, "error publishing show", http.StatusInternalServerError)
		return
	}

	h.queries.CreateNotificationsForShow(r.Context(), db.CreateNotificationsForShowParams{
		ShowID: uuid.NullUUID{UUID: id, Valid: true},
		Type:   "published",
	})

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
		adminpages.BandForm("band name is required").Render(r.Context(), w)
	}

	_, err := h.queries.CreateBand(r.Context(), db.CreateBandParams{
		Name:        name,
		Genre:       sql.NullString{String: r.FormValue("genre"), Valid: r.FormValue("genre") != ""},
		Description: sql.NullString{String: r.FormValue("description"), Valid: r.FormValue("description") != ""},
		WebsiteUrl:  sql.NullString{String: r.FormValue("website_url"), Valid: r.FormValue("website_url") != ""},
	})
	if err != nil {
		adminpages.BandForm("something went wrong, please try again").Render(r.Context(), w)
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
		adminpages.VenueForm("venue page is required").Render(r.Context(), w)
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

	_, err := h.queries.CreateVenue(r.Context(), db.CreateVenueParams{
		Name:     name,
		Address:  sql.NullString{String: r.FormValue("address"), Valid: r.FormValue("address") != ""},
		City:     sql.NullString{String: r.FormValue("city"), Valid: r.FormValue("city") != ""},
		State:    sql.NullString{String: r.FormValue("state"), Valid: r.FormValue("state") != ""},
		Capacity: nullCapacity,
	})
	if err != nil {
		adminpages.VenueForm("something went wrong, please try again.").Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) reloadShowForm(w http.ResponseWriter, r *http.Request, errMsg string) {
	venues, _ := h.queries.ListVenues(r.Context())
	bands, _ := h.queries.ListBands(r.Context())
	adminpages.ShowForm(venues, bands, errMsg).Render(r.Context(), w)
}

func slugify(title string, date time.Time) string {
	s := strings.ToLower(title + " " + date.Format(timeutil.SlugDate))

	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			result.WriteRune(r)
		case r >= '0' && r <= '9':
			result.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			result.WriteRune('-')
		}
	}

	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}
