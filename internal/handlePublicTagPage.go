package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicTagPage() http.HandlerFunc {
	type response struct {
		Tag   models.Tag
		Posts models.Publicaciones
		Meta  PageMeta
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		tagslug := mux.Vars(r)["tagid"] // We keep the parameter name as "tagid" in the router but use it as slug

		tag, err := s.store.tag.ObtenerPorSlug(tagslug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró la tag con slug %s: %s", tagslug, err.Error())
				s.handleError(r, w, 404, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando la tag con slug %s: %s", tagslug, err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
			}
			return
		}

		publicaciones, err := s.store.publicacion.ListarPorTag(tag.Id)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando publicaciones para tag: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		err = templates.Render(w, "tags-id.html", response{
			Tag:   tag,
			Posts: publicaciones.FiltrarPublicas().FiltrarRetiradas(),
			Meta: PageMeta{
				Titulo:      tag.Nombre,
				Keywords:    tag.Nombre,
				Descripcion: "Publicaciones en Vigo360 sobre " + tag.Nombre,
				Canonica:    fullCanonica("/sections/" + tag.Slug),
				BaseUrl:     baseUrl(),
			},
		})

		if err != nil {
			logger.Error("error generando página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
