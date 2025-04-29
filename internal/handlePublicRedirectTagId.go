package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handlePublicRedirectTagId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		tagid := mux.Vars(r)["tagid"]

		// Fetch the tag to get its slug
		tag, err := s.store.tag.Obtener(tagid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("tag no encontrado: %s", err.Error())
				s.handleError(r, w, http.StatusNotFound, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando tag: %s", err.Error())
				s.handleError(r, w, http.StatusInternalServerError, messages.ErrorDatos)
			}
			return
		}

		// Redirect to the new section URL using the slug
		http.Redirect(w, r, "/sections/"+tag.Slug, http.StatusMovedPermanently)
	}
}
