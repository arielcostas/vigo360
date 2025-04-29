package internal

import (
	"net/http"
)

func (s *Server) handlePublicRedirectTrabajos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/papers", http.StatusMovedPermanently)
	}
}
