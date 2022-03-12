package admin

import (
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

func EditPostPage(w http.ResponseWriter, r *http.Request) {
	// TODO Check author is same as session
	verifyLogin(w, r)
	post_id := mux.Vars(r)["id"]
	post := PostEditar{}

	err := db.QueryRowx(`SELECT titulo, resumen, contenido, alt_portada FROM publicaciones WHERE id = ?;`, post_id).StructScan(&post)

	// TODO Proper error handling
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error buscando el artículo en la base de datos"))
	}

	t.ExecuteTemplate(w, "post-id.html", struct {
		Post PostEditar
	}{Post: post})
}

func EditPostAction(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	post_id := mux.Vars(r)["id"]

	art_titulo := r.FormValue("art-titulo")
	art_resumen := r.FormValue("art-resumen")
	art_contenido := r.FormValue("art-contenido")

	// TODO Proper error page
	if !validarTitulo(art_titulo) {
		w.WriteHeader(400)
		w.Write([]byte("El título debe contener entre 3 y 80 caracteres"))
		return
	}

	if !validarResumen(art_resumen) {
		w.WriteHeader(400)
		w.Write([]byte("El resumen debe contener entre 3 y 300 caracteres"))
		return
	}

	if !validarContenido(art_contenido) {
		w.WriteHeader(400)
		w.Write([]byte("El contenido del artículo no puede estar vacío"))
		return
	}

	_, err := db.Exec(`UPDATE publicaciones SET titulo=?, resumen=?, contenido=? WHERE id=?`, art_titulo, art_resumen, art_contenido, post_id)

	// TODO Proper error page
	if err != nil {
		logger.Error("error saving edited post to database: %s", err.Error())
		w.WriteHeader(400)
		w.Write([]byte("error guardando cambios a la base de datos"))
	}

	w.Header().Add("Location", "/admin/post")
	w.WriteHeader(303)
}