/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import (
	"github.com/jmoiron/sqlx"
)

type SerieStore struct {
	db *sqlx.DB
}

func NewSerieStore(db *sqlx.DB) SerieStore {
	return SerieStore{
		db: db,
	}
}

func (s *SerieStore) Listar() (Publicaciones, error) {
	panic("por implementar")
}

func (s *SerieStore) Obtener(serie_id string) (Serie, error) {
	var serie Serie
	var row = s.db.QueryRow(`SELECT id, titulo FROM series WHERE id=?`, serie_id)
	var err = row.Scan(&serie.Id, &serie.Titulo)

	if err != nil {
		return Serie{}, err
	}

	return serie, nil
}
