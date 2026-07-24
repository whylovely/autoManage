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
		return fmt.Errorf("%w: vehicle is required", domain.ErrValidation)
	}
	if strings.TrimSpace(vehicle.Make) == "" || strings.TrimSpace(vehicle.Model) == "" {
		return fmt.Errorf("%w: make and model are required", domain.ErrValidation)
	}
	if vehicle.Year > time.Now().Year()+1 || vehicle.Year < 1886 {
		return fmt.Errorf("%w: vehicle year is out of range", domain.ErrValidation)
	}

	vinRegexp := regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
	vehicle.VIN = strings.ToUpper(strings.TrimSpace(vehicle.VIN))
	if !vinRegexp.MatchString(vehicle.VIN) {
		return fmt.Errorf("%w: VIN must contain 17 valid symbols", domain.ErrValidation)
	}

	if vehicle.Odometer < 0 {
		return fmt.Errorf("%w: odometer cannot be negative", domain.ErrValidation)
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
		return fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}

	return s.repo.Update(ctx, vehicle)
}

func (s *VehicleService) DeleteVehicle(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}

	return s.repo.Delete(ctx, id)
}

func (s *VehicleService) GetVehicle(ctx context.Context, id int64) (*domain.Vehicle, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *VehicleService) ListVehicles(ctx context.Context) ([]domain.Vehicle, error) {
	return s.repo.List(ctx)
}

func (s *VehicleService) UpdateOdometer(ctx context.Context, id, odometer int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}
	if odometer < 0 {
		return fmt.Errorf("%w: odometer cannot be negative", domain.ErrValidation)
	}

	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get vehicle: %w", err)
	}

	if odometer < v.Odometer {
		return fmt.Errorf("%w: odometer cannot decrease", domain.ErrValidation)
	}

	v.Odometer = odometer

	return s.repo.Update(ctx, v)
}
