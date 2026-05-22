package handlers

import (
	"net/http"

	"github.com/lepidoptera/lepidoptera/internal/auth"
	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
)

type AdminMiddleware struct {
	queries *db.Queries
}

func NewAdminMiddleware(queries *db.Queries) *AdminMiddleware {
	return &AdminMiddleware{queries: queries}
}

func (m *AdminMiddleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminSessionCookie, err := r.Cookie("admin_session")
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		// hash the token and look it up in the database
		tokenHash := auth.HashSessionToken(adminSessionCookie.Value)
		_, err = m.queries.GetAdminSession(r.Context(), tokenHash)
		if err != nil {
			// session not found, expired, or revoked
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// IsHTMX returns true if the request was made by HTMX
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
