package domain

import (
	"context"
)

type VehicleRepository interface {
	Create(ctx context.Context, vehicle *Vehicle) error
	Update(ctx context.Context, vehicle *Vehicle) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Vehicle, error)
	List(ctx context.Context) ([]Vehicle, error)
}
