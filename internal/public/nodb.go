/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"net/http"
)

func nodbPageError(err error) *appError {
	return &appError{Error: err, Message: "error rendering template", Response: "Error mostrando la página", Status: 500}
}

func SiguenosPage(w http.ResponseWriter, r *http.Request) *appError {
	err := t.ExecuteTemplate(w, "siguenos.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Síguenos",
			Descripcion: "Información sobre cómo seguir a Vigo360, y enterarse de sus últimas publicaciones y novedades.",
			Canonica:    FullCanonica("/siguenos"),
		},
	})
	if err != nil {
		return nodbPageError(err)
	}

	return nil
}

func LicenciasPage(w http.ResponseWriter, r *http.Request) *appError {
	err := t.ExecuteTemplate(w, "licencias.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Licencias",
			Descripcion: "Información legal relativa a Vigo360, desde licencias de uso libre hasta la política de privacidad.",
			Canonica:    FullCanonica("/licencia"),
		},
	})
	if err != nil {
		return nodbPageError(err)
	}

	return nil
}

func ContactoPage(w http.ResponseWriter, r *http.Request) *appError {
	err := t.ExecuteTemplate(w, "contacto.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Contacto",
			Descripcion: "Si necesitases contactar con Vigo360, aquí encontrarás cómo hacerlo.",
			Canonica:    FullCanonica("/contacto"),
		},
	})
	if err != nil {
		return nodbPageError(err)
	}

	return nil
}