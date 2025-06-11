package repository

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type MysqlPublicacionStore struct {
	db *sqlx.DB
}

func NewMysqlPublicacionStore(db *sqlx.DB) *MysqlPublicacionStore {
	return &MysqlPublicacionStore{
		db: db,
	}
}

func (s *MysqlPublicacionStore) Listar() (models.Publicaciones, error) {
	publicaciones := make(models.Publicaciones, 0)
	query := `SELECT p.id, COALESCE(fecha_publicacion, ""), fecha_actualizacion, COALESCE(p.legally_retired_at, ""), COALESCE(p.gone_at, ""), titulo, resumen, alt_portada, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, COALESCE(GROUP_CONCAT(tags.id), "") as tags_ids, COALESCE(GROUP_CONCAT(tags.nombre), "") as tags_nombres, COALESCE(GROUP_CONCAT(tags.slug), "") as tags_slugs FROM publicaciones p LEFT JOIN publicaciones_tags ON p.id = publicaciones_tags.publicacion_id LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id LEFT JOIN autores ON p.autor_id = autores.id GROUP BY id ORDER BY fecha_publicacion DESC;`

	rows, err := s.db.Query(query)

	if err != nil {
		return publicaciones, err
	}

	for rows.Next() {
		var (
			np            models.Publicacion
			rawTagIds     string
			rawTagNombres string
			rawTagSlugs   string
		)

		err = rows.Scan(&np.Id, &np.Fecha_publicacion, &np.Fecha_actualizacion, &np.Legally_retired_at, &np.Gone_at, &np.Titulo, &np.Resumen, &np.Alt_portada, &np.Autor.Id, &np.Autor.Nombre, &np.Autor.Email, &rawTagIds, &rawTagNombres, &rawTagSlugs)
		if err != nil {
			return models.Publicaciones{}, err
		}

		var (
			tags            = make([]models.Tag, 0)
			splitTagIds     = strings.Split(rawTagIds, ",")
			splitTagNombres = strings.Split(rawTagNombres, ",")
			splitTagSlugs   = strings.Split(rawTagSlugs, ",")
		)
		for i := 0; i < len(splitTagIds); i++ {
			if splitTagIds[i] == "" {
				continue
			}
			tags = append(tags, models.Tag{
				Id:     splitTagIds[i],
				Nombre: splitTagNombres[i],
				Slug:   splitTagSlugs[i],
			})
		}

		np.Tags = tags
		publicaciones = append(publicaciones, np)
	}
	return publicaciones, nil
}

func (s *MysqlPublicacionStore) ListarPorAutor(autor_id string) (models.Publicaciones, error) {
	var resultado = make(models.Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return models.Publicaciones{}, err
	}

	for _, pub := range publicaciones {
		if pub.Autor.Id == autor_id {
			resultado = append(resultado, pub)
		}
	}

	return resultado, nil
}

func (s *MysqlPublicacionStore) ListarPorTag(tag_id string) (models.Publicaciones, error) {
	var resultado = make(models.Publicaciones, 0)
	publicaciones, err := s.Listar()
	if err != nil {
		return models.Publicaciones{}, err
	}

	for _, pub := range publicaciones {
		for _, tag := range pub.Tags {
			if tag.Id == tag_id {
				resultado = append(resultado, pub)
				break
			}
		}
	}

	return resultado, nil
}

func (s *MysqlPublicacionStore) Existe(id string) (bool, error) {
	err := s.db.QueryRow("SELECT id FROM publicaciones WHERE id=? LIMIT 1", id).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (s *MysqlPublicacionStore) ObtenerPorId(id string, requirePublic bool) (models.Publicacion, error) {
	var post models.Publicacion
	var query = `SELECT publicaciones.id, alt_portada, titulo, resumen, contenido, COALESCE(fecha_publicacion, ""), fecha_actualizacion,COALESCE(publicaciones.legally_retired_at, ""), COALESCE(publicaciones.gone_at, ""), autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol, COALESCE(GROUP_CONCAT(tags.id), "") as tags_ids, COALESCE(GROUP_CONCAT(tags.nombre), "") as tags_names, COALESCE(GROUP_CONCAT(tags.slug), "") as tags_slugs
	FROM publicaciones
	LEFT JOIN autores on publicaciones.autor_id = autores.id
	LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id
	LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id
	WHERE publicaciones.id = ?
	GROUP BY publicaciones.id 
	ORDER BY publicaciones.fecha_publicacion DESC;`

	var (
		rawTagIds     string
		rawTagNombres string
		rawTagSlugs   string
	)

	var err = s.db.QueryRow(query, id).Scan(&post.Id, &post.Alt_portada, &post.Titulo, &post.Resumen, &post.Contenido, &post.Fecha_publicacion, &post.Fecha_actualizacion, &post.Legally_retired_at, &post.Gone_at, &post.Autor.Id, &post.Autor.Nombre, &post.Autor.Biografia, &post.Autor.Rol, &rawTagIds, &rawTagNombres, &rawTagSlugs)

	if err != nil {
		return models.Publicacion{}, err
	}

	if requirePublic && post.Fecha_publicacion == "" {
		return models.Publicacion{}, sql.ErrNoRows
	}

	tagIds := strings.Split(rawTagIds, ",")
	tagNames := strings.Split(rawTagNombres, ",")
	tagSlugs := strings.Split(rawTagSlugs, ",")

	for i := 0; i < len(tagIds); i++ {
		if tagIds[i] == "" {
			continue
		}
		post.Tags = append(post.Tags, models.Tag{
			Id:     tagIds[i],
			Nombre: tagNames[i],
			Slug:   tagSlugs[i],
		})
	}

	return post, nil
}

func (s *MysqlPublicacionStore) Buscar(termino string) (models.Publicaciones, error) {
	var query = `SELECT p.id, COALESCE(fecha_publicacion, ""), fecha_actualizacion, titulo, resumen, alt_portada, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, COALESCE(GROUP_CONCAT(tags.id), "") as tags_ids, COALESCE(GROUP_CONCAT(tags.nombre), "") as tags_nombres, COALESCE(GROUP_CONCAT(tags.slug), "") as tags_slugs FROM publicaciones p LEFT JOIN publicaciones_tags ON p.id = publicaciones_tags.publicacion_id LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id LEFT JOIN autores ON p.autor_id = autores.id WHERE MATCH(p.id, titulo, resumen, contenido, alt_portada) AGAINST (? IN NATURAL LANGUAGE MODE) GROUP BY id`

	rows, err := s.db.Query(query, "*"+termino+"*")
	if err != nil {
		return models.Publicaciones{}, err
	}

	var publicaciones = make(models.Publicaciones, 0)
	for rows.Next() {
		var (
			np            models.Publicacion
			rawTagIds     string
			rawTagNombres string
			rawTagSlugs   string
		)

		err = rows.Scan(&np.Id, &np.Fecha_publicacion, &np.Fecha_actualizacion, &np.Titulo, &np.Resumen, &np.Alt_portada, &np.Autor.Id, &np.Autor.Nombre, &np.Autor.Email, &rawTagIds, &rawTagNombres, &rawTagSlugs)
		if err != nil {
			return models.Publicaciones{}, err
		}

		var (
			tags            = make([]models.Tag, 0)
			splitTagIds     = strings.Split(rawTagIds, ",")
			splitTagNombres = strings.Split(rawTagNombres, ",")
			splitTagSlugs   = strings.Split(rawTagSlugs, ",")
		)
		for i := 0; i < len(splitTagIds); i++ {
			if splitTagIds[i] == "" {
				continue
			}
			tags = append(tags, models.Tag{
				Id:     splitTagIds[i],
				Nombre: splitTagNombres[i],
				Slug:   splitTagSlugs[i],
			})
		}

		np.Tags = tags
		publicaciones = append(publicaciones, np)
	}

	return publicaciones, nil
}
