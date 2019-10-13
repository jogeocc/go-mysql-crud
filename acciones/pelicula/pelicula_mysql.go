package pelicula

/*
	Esta Clase permite definir los metodos con
	las sentencias que se usarán con el modelo definido.
*/

import (
	"context"
	"database/sql"

	pRepo "../../acciones" // Cargamos el paquete repository (pelicula_method) Se Implementa los metodos de la interfaz
	models "../../models"
)

// NewSQLPostRepo retunrs la implementacion de Peliculas de la interfaz
func NewSQLPeliculaRepo(Conn *sql.DB) pRepo.PeliculaMethods {
	return &mysqlPeliculaRepo{
		Conn: Conn,
	}
}

//DEFINIMOS UNA ESTRUCTURA QUE RECIBE LA CONEXION A LA BD
type mysqlPeliculaRepo struct {
	Conn *sql.DB
}

func (m *mysqlPeliculaRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Pelicula, error) {

	rows, err := m.Conn.QueryContext(ctx, query, args...) //EJECUTAMOS LA QUERY
	if err != nil {
		return nil, err
	}
	defer rows.Close() //CERRAMOS LA LECTURA

	payload := make([]*models.Pelicula, 0) // CREAMOS UN ARRAY CON EL MAP DEL MODELO
	for rows.Next() {
		data := new(models.Pelicula)

		// METODO QUE RELLENA Y VERIFICA EL LISTADO DEL MODELO
		err := rows.Scan(
			&data.ID,
			&data.Nombre,
			&data.Director,
			&data.Anio,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

func (m *mysqlPeliculaRepo) Fetch(ctx context.Context) ([]*models.Pelicula, error) {
	query := "Select id, nombre, director, anio From peliculas" //QUERY PRINCIPAL

	return m.fetch(ctx, query) //LLAMAR LA QUERY PRINCIPAL
}

func (m *mysqlPeliculaRepo) GetByID(ctx context.Context, id int64) (*models.Pelicula, error) {
	query := "Select id, nombre, director, anio From peliculas where id=?"

	rows, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &models.Pelicula{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, models.ErrNotFound
	}

	return payload, nil
}

func (m *mysqlPeliculaRepo) Create(ctx context.Context, p *models.Pelicula) (int64, error) {
	query := "Insert peliculas SET nombre=?, anio=?, director=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, p.Nombre, p.Anio, p.Director)
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (m *mysqlPeliculaRepo) Update(ctx context.Context, p *models.Pelicula) (*models.Pelicula, error) {
	query := "Update peliculas set nombre=?, anio=?, director=? where id=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		p.Nombre,
		p.Anio,
		p.Director,
		p.ID,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return p, nil
}

func (m *mysqlPeliculaRepo) Delete(ctx context.Context, id int64) (bool, error) {
	query := "Delete From peliculas Where id=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
