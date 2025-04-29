package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) handlePublicRedirectTrabajoId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		trabajoid := vars["trabajoid"]
		http.Redirect(w, r, "/papers/"+trabajoid, http.StatusMovedPermanently)
	}
}
