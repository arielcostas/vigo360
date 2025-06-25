package internal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicAutorPage() http.HandlerFunc {
	type Response struct {
		Autor    models.Autor
		Posts    models.Publicaciones
		Trabajos models.Trabajos
		Meta     PageMeta
		JSONLD   template.JS
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var req_autor = mux.Vars(r)["id"]

		var autor models.Autor
		autor, err := s.store.autor.Obtener(req_autor)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró autor coon id "+req_autor, err.Error())
				s.handleError(r, w, 404, "No se ha encontrado el autor solicitado")
			} else {
				logger.Error("error inesperado recuperando datos: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
			}
			return
		}

		if autor.Gone_at != "" {
			logger.Information("autor %s marcado como gone, devolviendo 410", autor.Id)
			s.handleError(r, w, 410, "Este autor ya no está disponible.")
			return
		}

		publicaciones, err := s.store.publicacion.ListarPorAutor(autor.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		trabajos, err := s.store.trabajo.ListarPorAutor(autor.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		posts := publicaciones.FiltrarPublicas().FiltrarRetiradas()
		jsonld := buildAuthorJSONLD(autor, trabajos, posts)

		err = templates.Render(w, "authors-id.html", Response{
			Autor:    autor,
			Posts:    posts,
			Trabajos: trabajos,
			Meta: PageMeta{
				Titulo:      autor.Nombre,
				Descripcion: autor.Biografia,
				Canonica:    fullCanonica("/autores/" + autor.Id),
				BaseUrl:     baseUrl(),
			},
			JSONLD: jsonld,
		})

		if err != nil {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
		}
	}
}

// --- JSON-LD builder for author page ---
func buildAuthorJSONLD(autor models.Autor, trabajos models.Trabajos, posts models.Publicaciones) template.JS {
	type Person struct {
		Context       string   `json:"@context"`
		Type          string   `json:"@type"`
		Name          string   `json:"name"`
		AlternateName string   `json:"alternateName"`
		Identifier    string   `json:"identifier"`
		Email         string   `json:"email,omitempty"`
		JobTitle      string   `json:"jobTitle"`
		Image         string   `json:"image"`
		SameAs        []string `json:"sameAs,omitempty"`
		WorksFor      struct {
			Context string `json:"@context"`
			Type    string `json:"@type"`
			Name    string `json:"name"`
			Url     string `json:"url"`
		} `json:"worksFor"`
		Url string `json:"url"`
	}

	type Article struct {
		Context       string `json:"@context"`
		Type          string `json:"@type"`
		Headline      string `json:"headline"`
		Image         string `json:"image"`
		Url           string `json:"url"`
		DatePublished string `json:"datePublished"`
		Author        Person `json:"author,omitempty"`
	}

	type BlogPosting = Article

	type ProfilePage struct {
		Context    string        `json:"@context"`
		Type       string        `json:"@type"`
		MainEntity Person        `json:"mainEntity"`
		HasPart    []interface{} `json:"hasPart"`
	}

	type BreadcrumbList struct {
		Context         string `json:"@context"`
		Type            string `json:"@type"`
		ItemListElement []struct {
			Type     string `json:"@type"`
			Position int    `json:"position"`
			Item     struct {
				Id   string `json:"@id"`
				Name string `json:"name"`
			} `json:"item"`
		} `json:"itemListElement"`
	}

	// Helper to format date in ISO 8601 with system timezone
	formatISO8601 := func(dateStr string) string {
		if dateStr == "" {
			return ""
		}
		// Try to parse as RFC3339 or fallback to "2006-01-02 15:04:05" (common in Go/SQL)
		t, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", dateStr)
			if err != nil {
				return dateStr // fallback: return as is
			}
		}
		return t.Format(time.RFC3339)
	}

	// Build Person
	person := Person{
		Context:       "https://schema.org",
		Type:          "Person",
		Name:          autor.Nombre,
		AlternateName: autor.Id,
		Identifier:    autor.Id,
		JobTitle:      autor.Rol,
		Image:         "/static/profile/" + autor.Id + ".jpg",
		Url:           "https://vigo360.es/autores/" + autor.Id,
	}
	if autor.Email != "" {
		person.Email = "mailto:" + autor.Email
	}
	if autor.Web.Url != "" {
		person.SameAs = []string{autor.Web.Url}
	}
	person.WorksFor.Context = "https://schema.org"
	person.WorksFor.Type = "Organization"
	person.WorksFor.Name = "Vigo360"
	person.WorksFor.Url = "https://vigo360.es"

	// Build hasPart (Trabajos + up to 10 Posts)
	hasPart := []interface{}{}
	for _, t := range trabajos {
		hasPart = append(hasPart, Article{
			Context:       "https://schema.org",
			Type:          "Article",
			Headline:      t.Titulo,
			Image:         "/static/images/" + t.Id + ".webp",
			Url:           string(baseUrl()) + "/papers/" + t.Id,
			DatePublished: formatISO8601(t.Fecha_publicacion),
			Author:        person,
		})
	}
	for i, p := range posts {
		if i >= 10 {
			break
		}
		hasPart = append(hasPart, BlogPosting{
			Context:       "https://schema.org",
			Type:          "BlogPosting",
			Headline:      p.Titulo,
			Image:         "/static/images/" + p.Id + ".webp",
			Url:           string(baseUrl()) + "/post/" + p.Id,
			DatePublished: formatISO8601(p.Fecha_publicacion),
			Author:        person,
		})
	}

	profilePage := ProfilePage{
		Context:    "https://schema.org",
		Type:       "ProfilePage",
		MainEntity: person,
		HasPart:    hasPart,
	}

	// BreadcrumbList
	breadcrumb := BreadcrumbList{
		Context: "https://schema.org",
		Type:    "BreadcrumbList",
		ItemListElement: []struct {
			Type     string `json:"@type"`
			Position int    `json:"position"`
			Item     struct {
				Id   string `json:"@id"`
				Name string `json:"name"`
			} `json:"item"`
		}{
			{
				Type:     "ListItem",
				Position: 1,
				Item: struct {
					Id   string `json:"@id"`
					Name string `json:"name"`
				}{
					Id:   "/autores",
					Name: "Autores",
				},
			},
			{
				Type:     "ListItem",
				Position: 2,
				Item: struct {
					Id   string `json:"@id"`
					Name string `json:"name"`
				}{
					Id:   "/autores/" + autor.Id,
					Name: autor.Nombre,
				},
			},
		},
	}

	jsonldSlice := []interface{}{profilePage, breadcrumb}
	jsonldBytes, _ := json.Marshal(jsonldSlice)
	return template.JS(jsonldBytes)
}
