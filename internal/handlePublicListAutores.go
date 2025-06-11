package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicListAutores() http.HandlerFunc {
	type Response struct {
		Autores []models.Autor
		Meta    PageMeta
	}

	var meta = PageMeta{
		Titulo:      "Autores",
		Descripcion: "Conoce a los autores y colaboradores de Vigo360.",
		Canonica:    fullCanonica("/autores"),
		BaseUrl:     baseUrl(),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		autores, err := s.store.autor.Listar()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró ningún autor", err.Error())
				s.handleError(r, w, 404, "No se ha encontrado ningún autor ¿está el servidor bien configurado?")
			} else {
				logger.Error("error inesperado recuperando datos: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
			}
			return
		}

		autoresFiltrados := models.Autores(autores).FiltrarGone()

		err = templates.Render(w, "authors.html", Response{
			Autores: autoresFiltrados,
			Meta:    meta,
		})
		if err != nil {
			logger.Error("error renderizandoo plantilla: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
