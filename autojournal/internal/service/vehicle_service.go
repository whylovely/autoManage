package service

import (
	"autojournal/internal/domain"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type VehicleService struct {
	repo domain.VehicleRepository
}

func NewVehicleService(repo domain.VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

func validateVehicle(vehicle *domain.Vehicle) error {
	if vehicle == nil {
		return fmt.Errorf("Vehicle is null")
	}
	if strings.TrimSpace(vehicle.Make) == "" || strings.TrimSpace(vehicle.Model) == "" {
		return fmt.Errorf("Make or Model is null")
	}
	if vehicle.Year > time.Now().Year()+1 || vehicle.Year < 1886 {
		return fmt.Errorf("Year too much")
	}

	vinRegexp := regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
	vehicle.VIN = strings.ToUpper(strings.TrimSpace(vehicle.VIN))
	if !vinRegexp.MatchString(vehicle.VIN) {
		return fmt.Errorf("VIN must contain 17 valid symbols")
	}

	if vehicle.Odometer <= 0 {
		return fmt.Errorf("Odometer with minus")
	}

	return nil
}

func (s *VehicleService) CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	if err := validateVehicle(vehicle); err != nil {
		return err
	}

	return s.repo.Create(ctx, vehicle)
}

func (s *VehicleService) UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	if err := validateVehicle(vehicle); err != nil {
		return err
	}

	if vehicle.ID <= 0 {
		return fmt.Errorf("ID do not == 0")
	}

	return s.repo.Update(ctx, vehicle)
}

func (s *VehicleService) DeleteVehicle(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("ID wrong")
	}

	return s.repo.Delete(ctx, id)
}

func (s *VehicleService) GetVehicle(ctx context.Context, id int64) (*domain.Vehicle, error) {
	if id <= 0 {
		return nil, fmt.Errorf("vehicle id must be positive")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *VehicleService) ListVehicles(ctx context.Context) ([]domain.Vehicle, error) {
	return s.repo.List(ctx)
}

func (s *VehicleService) UpdateOdometer(ctx context.Context, id, odometer int64) error {
	if id <= 0 || odometer <= 1001 {
		return fmt.Errorf("ID or odometer wrong")
	}

	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get vehicle: %w", err)
	}

	if odometer < v.Odometer {
		return fmt.Errorf("odometer cannot be negative")
	}

	v.Odometer = odometer

	return s.repo.Update(ctx, v)
}
