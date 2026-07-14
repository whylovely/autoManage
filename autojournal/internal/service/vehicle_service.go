package service

import "autojournal/internal/domain"

type VehicleService struct {
	repo domain.VehicleRepository
}

func NewVehicleService(repo domain.VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}
