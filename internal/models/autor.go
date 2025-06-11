package models

type Web struct {
	Url    string
	Titulo string
}

type Autor struct {
	Id        string
	Nombre    string
	Email     string
	Rol       string
	Biografia string
	Gone_at   string

	Web Web

	Publicaciones Publicaciones
}

type Autores []Autor

// FiltrarGone devuelve una nueva lista de autores, excluyendo aquellos que tienen un valor en Gone_at.
func (as Autores) FiltrarGone() Autores {
	var resultado Autores
	for _, a := range as {
		if a.Gone_at == "" {
			resultado = append(resultado, a)
		}
	}
	return resultado
}
