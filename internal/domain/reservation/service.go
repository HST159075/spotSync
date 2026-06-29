package reservation

import (
	"errors"
	resdto "spotsync/internal/domain/reservation/dto"
	"time"
)

type Service interface {
	CreateReservation(userID uint, req resdto.CreateReservationRequest) (*resdto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]resdto.MyReservationResponse, error)
	CancelReservation(userID uint, reservationID uint) error
	GetAllReservations() ([]resdto.AdminReservationResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateReservation(userID uint, req resdto.CreateReservationRequest) (*resdto.ReservationResponse, error) {
	res, err := s.repo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}
	return &resdto.ReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    res.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) GetMyReservations(userID uint) ([]resdto.MyReservationResponse, error) {
	reservations, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	var result []resdto.MyReservationResponse
	for _, r := range reservations {
		result = append(result, resdto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: resdto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (s *service) CancelReservation(userID uint, reservationID uint) error {
	res, err := s.repo.FindByID(reservationID)
	if err != nil {
		return errors.New("reservation not found")
	}
	if res.UserID != userID {
		return errors.New("forbidden")
	}
	if res.Status == "cancelled" {
		return errors.New("reservation already cancelled")
	}
	return s.repo.Cancel(reservationID)
}

func (s *service) GetAllReservations() ([]resdto.AdminReservationResponse, error) {
	reservations, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var result []resdto.AdminReservationResponse
	for _, r := range reservations {
		result = append(result, resdto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: resdto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			User: resdto.UserInfo{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			},
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}
