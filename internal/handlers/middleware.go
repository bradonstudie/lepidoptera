package handlers

import (
	"net/http"
	"os"
)

// AdminOnly middleware — checks a simple session token for now.
// Replace with magic link session handling once auth is built out.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value != os.Getenv("ADMIN_SESSION_SECRET") {
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
