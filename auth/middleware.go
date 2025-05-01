package auth

import "net/http"

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; form-action 'self'; img-src 'self' data: https:; script-src 'sha256-osjxnKEPL/pQJbFk1dKsF7PYFmTyMWGmVSiL9inhxJY=' 'unsafe-hashes'; style-src 'self' 'unsafe-inline';")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuthenticated(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}
