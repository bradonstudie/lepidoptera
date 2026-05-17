package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/web/components"
	"github.com/lepidoptera/lepidoptera/web/pages"
)

type ShowHandler struct {
	queries *db.Queries
}

func NewShowHandler(queries *db.Queries) *ShowHandler {
	return &ShowHandler{queries: queries}
}

func (h *ShowHandler) Index(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")

	nullGenre := sql.NullString{
		String: genre,
		Valid:  genre != "",
	}

	shows, err := h.queries.ListPublishedShows(r.Context(), nullGenre)
	if err != nil {
		http.Error(w, "error loading shows", http.StatusInternalServerError)
		return
	}

	if IsHTMX(r) {
		components.ShowList(shows).Render(r.Context(), w)
		return
	}

	genres, _ := h.queries.ListGenres(r.Context())
	genreList := make([]string, 0, len(genres))
	for _, genre := range genres {
		if genre.Valid {
			genreList = append(genreList, genre.String)
		}
	}

	pages.Index(shows, genreList).Render(r.Context(), w)
}

func (h *ShowHandler) Detail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	show, err := h.queries.GetShowBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "show not found", http.StatusNotFound)
		return
	}

	bands, err := h.queries.GetBandsByShow(r.Context(), show.ID)
	if err != nil {
		http.Error(w, "error loading show", http.StatusInternalServerError)
	}

	pages.ShowDetail(show, bands).Render(r.Context(), w)
}
