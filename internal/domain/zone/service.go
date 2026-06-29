package zone

import (
	"errors"
	zonedto "spotsync/internal/domain/zone/dto"
	"time"
)

type Service interface {
	CreateZone(req zonedto.CreateZoneRequest) (*zonedto.ZoneResponse, error)
	GetAllZones() ([]zonedto.ZoneResponse, error)
	GetZoneByID(id uint) (*zonedto.ZoneResponse, error)
	UpdateZone(id uint, req zonedto.UpdateZoneRequest) (*zonedto.ZoneResponse, error)
	DeleteZone(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) toResponse(z Model, available int) zonedto.ZoneResponse {
	return zonedto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt.Format(time.RFC3339),
	}
}

func (s *service) CreateZone(req zonedto.CreateZoneRequest) (*zonedto.ZoneResponse, error) {
	z := &Model{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.repo.Create(z); err != nil {
		return nil, err
	}
	resp := s.toResponse(*z, z.TotalCapacity)
	return &resp, nil
}

func (s *service) GetAllZones() ([]zonedto.ZoneResponse, error) {
	zones, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var result []zonedto.ZoneResponse
	for _, z := range zones {
		active, err := s.repo.CountActiveReservations(z.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, s.toResponse(z, z.TotalCapacity-int(active)))
	}
	return result, nil
}

func (s *service) GetZoneByID(id uint) (*zonedto.ZoneResponse, error) {
	z, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("zone not found")
	}
	active, err := s.repo.CountActiveReservations(z.ID)
	if err != nil {
		return nil, err
	}
	resp := s.toResponse(*z, z.TotalCapacity-int(active))
	return &resp, nil
}

func (s *service) UpdateZone(id uint, req zonedto.UpdateZoneRequest) (*zonedto.ZoneResponse, error) {
	z, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("zone not found")
	}
	if req.Name != "" {
		z.Name = req.Name
	}
	if req.Type != "" {
		z.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		z.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		z.PricePerHour = req.PricePerHour
	}
	if err := s.repo.Update(z); err != nil {
		return nil, err
	}
	active, err := s.repo.CountActiveReservations(z.ID)
	if err != nil {
		return nil, err
	}
	resp := s.toResponse(*z, z.TotalCapacity-int(active))
	return &resp, nil
}

func (s *service) DeleteZone(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("zone not found")
	}
	return s.repo.Delete(id)
}