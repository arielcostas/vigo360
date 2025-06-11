package models

type Publicacion struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string
	Legally_retired_at  string
	Gone_at             string
	Alt_portada         string
	Titulo              string
	Resumen             string
	Contenido           string
	Comentarios         []Comentario

	Autor Autor
	Tags  []Tag
}
