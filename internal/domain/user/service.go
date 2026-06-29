package user

import (
	"errors"
	"spotsync/internal/auth"
	userdto "spotsync/internal/domain/user/dto"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(req userdto.RegisterRequest) (*Model, error)
	Login(req userdto.LoginRequest) (*userdto.LoginResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Register(req userdto.RegisterRequest) (*Model, error) {
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "driver"
	}

	u := &Model{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     role,
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Login(req userdto.LoginRequest) (*userdto.LoginResponse, error) {
	u, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, err
	}

	return &userdto.LoginResponse{
		Token: token,
		User: userdto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}