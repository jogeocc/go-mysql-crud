package repsitory

import (
	"context"

	"github.com/jogeocc/go-mysql-crud/models"
)

// METODOS PARA LA CLASE PELICULA
type PeliculaMethods interface {
	Fetch(ctx context.Context, num int64) ([]*models.Pelicula, error)
	GetByID(ctx context.Context, id int64) (*models.Pelicula, error)
	Create(ctx context.Context, p *models.Pelicula) (int64, error)
	Update(ctx context.Context, p *models.Pelicula) (*models.Pelicula, error)
	Delete(ctx context.Context, id int64) (bool, error)
}
