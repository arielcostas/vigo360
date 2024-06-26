package internal

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handlePublicIndexAlgolia() http.HandlerFunc {
	type Post struct {
		Id                  string `json:"id"`
		Alt_portada         string `json:"alt_portada"`
		Titulo              string `json:"titulo"`
		Resumen             string `json:"resumen"`
		Contenido           string `json:"contenido"`
		Fecha_publicacion   string `json:"fecha_publicacion"`
		Fecha_actualizacion string `json:"fecha_actualizacion"`
		Autor               string `json:"autor"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var result []Post = make([]Post, 0)

		var authHeader = r.Header.Get("Authorization")
		if authHeader == "" {
			log.Error("no se especificó un token de autorización")
			r.Header.Add("WWW-Authenticate", `Basic realm="Restricted"`)
			s.handleJsonError(r, w, 401, messages.ErrorSinAutenticar)
			return
		}

		var username = os.Getenv("ALGOLIA_API_USERNAME")
		var password = os.Getenv("ALGOLIA_API_PASSWORD")
		var auth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

		if authHeader != "Basic "+auth {
			log.Error("token de autorización inválido")
			r.Header.Add("WWW-Authenticate", `Basic realm="Restricted"`)
			s.handleJsonError(r, w, 401, messages.ErrorSinAutenticar)
			return
		}

		db := database.GetDB()
		rows, err := db.Query("SELECT p.id, p.alt_portada, p.titulo, p.resumen, LEFT(p.contenido, 8500), p.fecha_publicacion, p.fecha_actualizacion, autores.nombre FROM publicaciones p LEFT JOIN autores ON p.autor_id = autores.id WHERE fecha_publicacion IS NOT NULL AND fecha_publicacion <= NOW() AND legally_retired_at IS NULL")

		if err != nil {
			log.Error("error leyendo adjuntos: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorDatos)
			return
		}

		for rows.Next() {
			var post Post

			err = rows.Scan(&post.Id, &post.Alt_portada, &post.Titulo, &post.Resumen, &post.Contenido, &post.Fecha_publicacion, &post.Fecha_actualizacion, &post.Autor)
			if err != nil {
				log.Error("error escaneando adjuntos: %s", err.Error())
				s.handleJsonError(r, w, 500, messages.ErrorDatos)
				return
			}
			result = append(result, post)
		}

		resbytes, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			log.Error("error escribiendo json de respuesta: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorRender)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(resbytes)
	}
}
