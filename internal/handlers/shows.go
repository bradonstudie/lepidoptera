package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShowHandler struct {
	db *pgxpool.Pool
}

func NewShowHandler(db *pgxpool.Pool) *ShowHandler {
	return &ShowHandler{db: db}
}

func (h *ShowHandler) Index(w http.ResponseWriter, r *http.Request) {
	// genre := r.URL.Query().Get("genre")
	// TODO: query published shows, filtered by genre if provided
	// if IsHTMX(r): render ShowList fragment only
	// else: render full Index page
	w.Write([]byte("shows index — coming soon"))
}

func (h *ShowHandler) Detail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	// TODO: query show by slug, render ShowDetail page
	w.Write([]byte("show detail for: " + slug))
}
