package internal

import (
	"net/http"
)

func (s *Server) handlePublicRedirectTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/sections", http.StatusMovedPermanently)
	}
}
