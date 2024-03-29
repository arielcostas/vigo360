package service

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/thanhpk/randstr"
	"vigo360.es/new/internal/models"
)

var Err_ComentarioPublicacionInvalida = errors.New("el ID de la publicación no es válido")
var Err_ComentarioNombreInvalido = errors.New("el nombre de autor del comentario no es válido")
var Err_ComentarioContenidoInvalido = errors.New("el contenido del comentario no es válido")
var Err_ComentarioErrorBaseDatos = errors.New("hubo un error guardando el comentario")

func (se *Comentario) AgregarComentario(
	publicacion_id string,
	nombre string,
	contenido string,
	es_autor bool,
	autor_original bool,
) (models.Comentario, error) {

	return se.AgregarRespuesta(publicacion_id, nombre, contenido, "", es_autor, autor_original)
}

func (se *Comentario) AgregarRespuesta(
	publicacion_id string,
	nombre string,
	contenido string,
	padre string,
	es_autor bool,
	autor_original bool,
) (models.Comentario, error) {
	if nombre == "" || len(nombre) > 40 {
		return models.Comentario{}, Err_ComentarioNombreInvalido
	}

	if contenido == "" || utf8.RuneCountInString(contenido) > 2000 {
		return models.Comentario{}, Err_ComentarioContenidoInvalido
	}

	publicacion_existe, err := se.pstore.Existe(publicacion_id)
	if err != nil {
		return models.Comentario{}, err
	}
	if !publicacion_existe {
		return models.Comentario{}, Err_ComentarioPublicacionInvalida
	}

	var nuevo_comentario = models.Comentario{
		Id:             randstr.String(13),
		Publicacion_id: publicacion_id,
		Padre_id:       padre,

		Nombre:         nombre,
		Es_autor:       es_autor,
		Autor_original: autor_original,
		Contenido:      contenido,

		Estado: models.EstadoPendiente,
	}

	if es_autor || autor_original {
		nuevo_comentario.Estado = models.EstadoAprobado
	}

	dberr := se.cstore.GuardarComentario(nuevo_comentario)
	if dberr != nil {
		fmt.Printf("%s\n", dberr.Error())
		return models.Comentario{}, Err_ComentarioErrorBaseDatos
	}
	return nuevo_comentario, nil
}
