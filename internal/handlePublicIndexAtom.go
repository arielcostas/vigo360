package internal

import (
	"net/http"
	"os"
	"time"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handlePublicIndexAtom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		pp, err := s.store.publicacion.Listar()
		if err != nil {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		pp = pp.FiltrarPublicas()
		pp = pp.FiltrarRetiradas()
		if len(pp) > 20 {
			pp = pp[:20]
		}

		lastUpdate, _ := pp.ObtenerUltimaActualizacion()

		w.Header().Add("Content-Type", "application/atom+xml;charset=UTF-8")

		// Using the new RenderAtom function that properly handles XML tags
		err = RenderAtom(w, atomParams{
			Dominio:    os.Getenv("DOMAIN"),
			Path:       "/atom.xml",
			Titulo:     "Publicaciones",
			Subtitulo:  "Últimas publicaciones en el sitio web de Vigo360",
			LastUpdate: lastUpdate.Format(time.RFC3339),
			Entries:    pp,
		})

		if err != nil {
			logger.Error("error generando feed atom: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
			return
		}
	}
}
