package user

import (
	"net/http"
	userdto "spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service) *Handler {
	return &Handler{service, validator.New()}
}

func (h *Handler) Register(c echo.Context) error {
	var req userdto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	u, err := h.service.Register(req)
	if err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusCreated, "User registered successfully", userdto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req userdto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	resp, err := h.service.Login(req)
	if err != nil {
		return httpresponse.Error(c, http.StatusUnauthorized, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusOK, "Login successful", resp)
}