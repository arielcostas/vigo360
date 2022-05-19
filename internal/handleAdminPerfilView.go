/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handleAdminPerfilView() http.HandlerFunc {
	type respuesta struct {
		Autor models.Autor
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		autor, err := s.store.autor.Obtener(sess.Autor_id)

		if err != nil {
			logger.Error("error recuperando perfil del usuario actual")
			s.handleError(w, 500, messages.ErrorDatos)
		}

		templates.Render(w, "admin-perfil.html", respuesta{
			Autor: autor,
		})
	}
}