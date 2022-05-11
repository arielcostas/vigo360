/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"html/template"

	"vigo360.es/new/internal/templates"
)

/*
	TODO: Refactor this complete file (and get rid of it)
*/

var t = template.Must(template.New("atom.xml").Funcs(templates.Functions).Parse(rawAtom))

var rawAtom string = `<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="es-ES">
	<id>{{ .Dominio }}{{ .Path }}</id>
	<title>{{ .Titulo }} - Vigo360</title>
	<subtitle>{{ .Subtitulo }}</subtitle>
	<updated>{{ .LastUpdate }}</updated>
	<generator uri="https://gitlab.com/vigo360/new.vigo360.es">Vigo360</generator>
	<link rel="self" href="{{ .Dominio }}{{ .Path }}" />
	<icon>/static/logo.png</icon>
	
	{{- $domain := .Dominio }}
	{{- range .Entries }}
	<entry>
		<id>{{ $domain }}/post/{{ .Id }}</id>
		<title>{{ .Titulo }}</title>
		<published>{{ date3339 .Fecha_publicacion }}</published>
		<updated>{{ date3339 .Fecha_actualizacion }}</updated>
		<link rel="alternate" href="{{ $domain }}/post/{{ .Id }}" />
		<summary>{{ .Resumen }}</summary>
		<author>
			<name>{{ .Autor.Nombre }}</name>
			<email>{{ .Autor.Email }}</email>
			<uri>{{ $domain }}/autores/{{ .Autor.Id }}</uri>
		</author>
		{{- range .Tags }}
		<category term="{{ .Nombre }}" />
		{{- end}}
	</entry>
	{{- end}}
</feed>
`

type atomParams struct {
	Dominio    string
	Path       string
	Titulo     string
	Subtitulo  string
	LastUpdate string
	Entries    interface{}
}