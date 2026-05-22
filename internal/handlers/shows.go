package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lepidoptera/lepidoptera/internal/service"
	"github.com/lepidoptera/lepidoptera/web/components"
	"github.com/lepidoptera/lepidoptera/web/pages"
)

type ShowHandler struct {
	showService *service.ShowService
}

func NewShowHandler(showService *service.ShowService) *ShowHandler {
	return &ShowHandler{showService: showService}
}

func (h *ShowHandler) Index(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")

	shows, err := h.showService.ListPublishedShows(r.Context(), genre)
	if err != nil {
		http.Error(w, "error loading shows", http.StatusInternalServerError)
		return
	}

	if IsHTMX(r) {
		components.ShowList(shows).Render(r.Context(), w)
		return
	}

	genres, err := h.showService.ListGenres(r.Context())
	if err != nil {
		http.Error(w, "error loading genres", http.StatusInternalServerError)
		return
	}

	pages.Index(shows, genres).Render(r.Context(), w)
}

func (h *ShowHandler) Detail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	show, err := h.showService.GetShowDetail(r.Context(), slug)
	if err != nil {
		http.Error(w, "show not found", http.StatusNotFound)
		return
	}

	pages.ShowDetail(show).Render(r.Context(), w)
}
