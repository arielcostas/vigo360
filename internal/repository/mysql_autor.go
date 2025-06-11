package repository

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type MysqlAutorStore struct {
	db *sqlx.DB
}

func NewMysqlAutorStore(db *sqlx.DB) *MysqlAutorStore {
	return &MysqlAutorStore{
		db: db,
	}
}

func (s *MysqlAutorStore) Listar() ([]models.Autor, error) {
	var autores = make([]models.Autor, 0)
	var rows, err = s.db.Query(`SELECT id, nombre, email, rol, biografia, COALESCE(gone_at, "") FROM autores`)
	if err != nil {
		return []models.Autor{}, err
	}

	for rows.Next() {
		var na models.Autor
		err = rows.Scan(&na.Id, &na.Nombre, &na.Email, &na.Rol, &na.Biografia, &na.Gone_at)
		if err != nil {
			return []models.Autor{}, err
		}
		autores = append(autores, na)
	}

	return autores, nil
}

func (s *MysqlAutorStore) Obtener(autor_id string) (models.Autor, error) {
	var autor models.Autor
	var row = s.db.QueryRow(`SELECT id, nombre, email, rol, biografia, web_url, web_titulo FROM autores WHERE id=?`, autor_id)
	var err = row.Scan(&autor.Id, &autor.Nombre, &autor.Email, &autor.Rol, &autor.Biografia, &autor.Web.Url, &autor.Web.Titulo)

	if err != nil {
		return models.Autor{}, err
	}

	return autor, nil
}

func (s *MysqlAutorStore) Buscar(termino string) ([]models.Autor, error) {
	var autores []models.Autor

	var query = `SELECT id, nombre, email, rol, biografia FROM autores WHERE CONCAT_WS(" ", nombre, email, rol, biografia) LIKE ?`
	var rows, err = s.db.Query(query, "%"+termino+"%")
	if err != nil {
		return []models.Autor{}, err
	}

	for rows.Next() {
		var na models.Autor
		err = rows.Scan(&na.Id, &na.Nombre, &na.Email, &na.Rol, &na.Biografia)
		if err != nil {
			return []models.Autor{}, err
		}
		autores = append(autores, na)
	}

	return autores, nil
}
