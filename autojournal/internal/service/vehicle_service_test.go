package service

import (
	"autojournal/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type vehicleRepoStub struct {
	vehicle      *domain.Vehicle
	vehicles     []domain.Vehicle
	created      *domain.Vehicle
	deletedID    int64
	createCalled bool
	updateCalled bool
}

func (r *vehicleRepoStub) Create(_ context.Context, vehicle *domain.Vehicle) error {
	r.createCalled = true
	r.created = vehicle
	vehicle.ID = 1
	return nil
}

func (r *vehicleRepoStub) Update(_ context.Context, vehicle *domain.Vehicle) error {
	r.updateCalled = true
	r.vehicle = vehicle
	return nil
}

func (r *vehicleRepoStub) Delete(_ context.Context, id int64) error {
	r.deletedID = id
	return nil
}

func (r *vehicleRepoStub) GetByID(context.Context, int64) (*domain.Vehicle, error) {
	if r.vehicle == nil {
		return nil, domain.ErrNotFound
	}
	return r.vehicle, nil
}

func (r *vehicleRepoStub) List(context.Context) ([]domain.Vehicle, error) { return r.vehicles, nil }

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

func TestVehicleService_CreateVehicle_NormalizesVIN(t *testing.T) {
	repo := &vehicleRepoStub{}
	service := NewVehicleService(repo)
	vehicle := &domain.Vehicle{
		VIN: " 1hgcm82633a004352 ", Make: "Honda", Model: "Accord", Year: 2020,
	}

	require.NoError(t, service.CreateVehicle(context.Background(), vehicle))
	assert.True(t, repo.createCalled)
	assert.Same(t, vehicle, repo.created)
	assert.Equal(t, "1HGCM82633A004352", vehicle.VIN)
	assert.EqualValues(t, 1, vehicle.ID)
}

func TestVehicleService_UpdateOdometer_PersistsIncrease(t *testing.T) {
	repo := &vehicleRepoStub{vehicle: &domain.Vehicle{
		ID: 1, VIN: "1HGCM82633A004352", Make: "Honda", Model: "Accord", Year: 2020, Odometer: 10_000,
	}}
	service := NewVehicleService(repo)

	require.NoError(t, service.UpdateOdometer(context.Background(), 1, 11_500))
	assert.True(t, repo.updateCalled)
	assert.EqualValues(t, 11_500, repo.vehicle.Odometer)
}

func TestVehicleService_CRUDDelegatesToRepository(t *testing.T) {
	stored := &domain.Vehicle{
		ID: 7, VIN: "1HGCM82633A004352", Make: "Honda", Model: "Accord", Year: 2020, Odometer: 20_000,
	}
	repo := &vehicleRepoStub{vehicle: stored, vehicles: []domain.Vehicle{*stored}}
	service := NewVehicleService(repo)

	got, err := service.GetVehicle(context.Background(), stored.ID)
	require.NoError(t, err)
	assert.Same(t, stored, got)
	list, err := service.ListVehicles(context.Background())
	require.NoError(t, err)
	assert.Equal(t, repo.vehicles, list)

	updated := *stored
	updated.Make = "Acura"
	updated.Odometer = 21_000
	require.NoError(t, service.UpdateVehicle(context.Background(), &updated))
	assert.True(t, repo.updateCalled)
	assert.Equal(t, "Acura", repo.vehicle.Make)

	require.NoError(t, service.DeleteVehicle(context.Background(), stored.ID))
	assert.Equal(t, stored.ID, repo.deletedID)
}

func TestVehicleService_RejectsInvalidIDs(t *testing.T) {
	service := NewVehicleService(&vehicleRepoStub{})
	assert.ErrorIs(t, service.DeleteVehicle(context.Background(), 0), domain.ErrValidation)
	_, err := service.GetVehicle(context.Background(), -1)
	assert.ErrorIs(t, err, domain.ErrValidation)
	assert.ErrorIs(t, service.UpdateOdometer(context.Background(), 0, 100), domain.ErrValidation)
}
