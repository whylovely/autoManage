package service

import (
	"autojournal/internal/domain"
	"context"
	"errors"
	"testing"
)

type vehicleRepoStub struct {
	vehicle      *domain.Vehicle
	updateCalled bool
}

func (r *vehicleRepoStub) Create(context.Context, *domain.Vehicle) error { return nil }

func (r *vehicleRepoStub) Update(_ context.Context, vehicle *domain.Vehicle) error {
	r.updateCalled = true
	r.vehicle = vehicle
	return nil
}

func (r *vehicleRepoStub) Delete(context.Context, int64) error { return nil }

func (r *vehicleRepoStub) GetByID(context.Context, int64) (*domain.Vehicle, error) {
	return r.vehicle, nil
}

func (r *vehicleRepoStub) List(context.Context) ([]domain.Vehicle, error) { return nil, nil }

func TestVehicleService_CreateVehicle_RejectsEmptyMake(t *testing.T) {
	service := NewVehicleService(&vehicleRepoStub{})

	err := service.CreateVehicle(context.Background(), &domain.Vehicle{})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestVehicleService_UpdateOdometer_RejectsDecrease(t *testing.T) {
	repo := &vehicleRepoStub{vehicle: &domain.Vehicle{ID: 1, Odometer: 10_000}}
	service := NewVehicleService(repo)

	err := service.UpdateOdometer(context.Background(), 1, 9_999)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if repo.updateCalled {
		t.Fatal("repository update must not be called when odometer decreases")
	}
}
