package repository

import (
	"context"

	"../models" // IMPORTAMOS EL PAQUETE DE LOS MODELOS
)

// LA DEFINICIOM DE LA INTERFAZ CON SUS METODOS PARA LA CLASE PELICULA
type PeliculaMethods interface {
	Fetch(ctx context.Context) ([]*models.Pelicula, error)
	GetByID(ctx context.Context, id int64) (*models.Pelicula, error)
	Create(ctx context.Context, p *models.Pelicula) (int64, error)
	Update(ctx context.Context, p *models.Pelicula) (*models.Pelicula, error)
	Delete(ctx context.Context, id int64) (bool, error)
}
