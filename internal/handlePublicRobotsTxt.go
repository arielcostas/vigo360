package internal

import (
	"fmt"
	"net/http"
)

func (s *Server) handlePublicRobotsTxt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, `User-agent: *
Disallow: /admin/
Disallow: /search

Sitemap: /sitemap.xml
`)
	}
}
