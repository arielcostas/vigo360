package routes

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/parser"
	"github.com/gorilla/mux"
)

type PostPost struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string
	Alt_portada         string
	Titulo              string
	Resumen             string
	ContenidoRaw        string `db:"contenido"`
	Contenido           template.HTML
	Autor_id            string
	Autor_nombre        string
	Autor_rol           string
	Autor_biografia     string
}

type PostParams struct {
	Post PostPost
}

func PostPage(w http.ResponseWriter, r *http.Request) {
	req_post_id := mux.Vars(r)["postid"]
	query := `SELECT publicaciones.id, alt_portada, titulo, resumen, contenido, 
	DATE_FORMAT(publicaciones.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	DATE_FORMAT(publicaciones.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol
FROM publicaciones 
LEFT JOIN autores on publicaciones.autor_id = autores.id 
WHERE publicaciones.id = ? ORDER BY publicaciones.fecha_publicacion DESC;`

	post := PostPost{}
	row := db.QueryRowx(query, req_post_id)
	err := row.StructScan(&post)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Result is in markdown, convert to HTML
	var buf bytes.Buffer
	parser.Parser.Convert([]byte(post.ContenidoRaw), &buf)

	post.Contenido = template.HTML(buf.Bytes())

	t.ExecuteTemplate(w, "post.html", PostParams{
		Post: post,
	})
}
