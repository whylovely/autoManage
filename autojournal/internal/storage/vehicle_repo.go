package storage

import (
	"autojournal/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type VehicleRepo struct {
	db *sqlx.DB
}

func NewVehicleRepo(db *sqlx.DB) *VehicleRepo {
	return &VehicleRepo{db: db}
}

func (r *VehicleRepo) Create(ctx context.Context, vehicle *domain.Vehicle) error {
	const query = `
			INSERT INTO vehicles (
					vin, make, model, year, color,
					engine_volume, fuel_type, odometer, notes
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		vehicle.VIN,
		vehicle.Make,
		vehicle.Model,
		vehicle.Year,
		vehicle.Color,
		vehicle.EngineVolume,
		vehicle.FuelType,
		vehicle.Odometer,
		vehicle.Notes,
	).Scan(&vehicle.ID)
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}

	return nil
}

func (r *VehicleRepo) Update(ctx context.Context, vehicle *domain.Vehicle) error {
	const query = `
			UPDATE vehicles
			SET vin = ?, make = ?, model = ?, year = ?, color = ?,
				engine_volume = ?, fuel_type = ?, odometer = ?, notes = ?
			WHERE id = ?`

	result, err := r.db.ExecContext(
		ctx,
		query,
		vehicle.VIN,
		vehicle.Make,
		vehicle.Model,
		vehicle.Year,
		vehicle.Color,
		vehicle.EngineVolume,
		vehicle.FuelType,
		vehicle.Odometer,
		vehicle.Notes,
		vehicle.ID,
	)
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get updated rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("vehicle %d not found", vehicle.ID)
	}

	return nil
}

func (r *VehicleRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM vehicles WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("vehicle %d not found", id)
	}

	return nil
}

func (r *VehicleRepo) GetByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
	const query = `
			SELECT
				id, vin, make, model, year, color,
			engine_volume AS enginevolume,
			fuel_type AS fueltype,
			odometer, notes,
			created_at AS createdat,
			updated_at AS updateat
			FROM vehicles
			WHERE id = ?`

	var vehicle domain.Vehicle
	if err := r.db.GetContext(ctx, &vehicle, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("vehicle %d not found", id)
		}
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	return &vehicle, nil
}

func (r *VehicleRepo) List(ctx context.Context) ([]domain.Vehicle, error) {
	const query = `
			SELECT
				id, vin, make, model, year, color,
			engine_volume AS enginevolume,
			fuel_type AS fueltype,
			odometer, notes,
			created_at AS createdat,
			updated_at AS updateat
			FROM vehicles
			ORDER BY id DESC`

	var vehicles []domain.Vehicle
	if err := r.db.SelectContext(ctx, &vehicles, query); err != nil {
		return nil, fmt.Errorf("list vehicles: %w", err)
	}

	return vehicles, nil
}
