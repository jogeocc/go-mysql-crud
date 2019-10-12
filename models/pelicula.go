package models

type Pelicula struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Anio     int    `json:"anio"`
	Director string `json:"director"`
}

type Peliculas []Pelicula
