package repository

import (
	"vigo360.es/new/internal/models"
)

type TagStore interface {
	Listar() ([]models.Tag, error)
	Obtener(string) (models.Tag, error)
	ObtenerPorSlug(string) (models.Tag, error)
}
