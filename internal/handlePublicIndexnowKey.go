package internal

import (
	"fmt"
	"net/http"
)

func (s *Server) handlePublicIndexnowKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s\n", r.URL.Path[1:len(r.URL.Path)-4])
	}
}
