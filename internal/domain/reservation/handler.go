package reservation

import (
	"net/http"
	resdto "spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"
	"strconv"

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

func (h *Handler) CreateReservation(c echo.Context) error {
	userID := c.Get("userID").(uint)

	var req resdto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	res, err := h.service.CreateReservation(userID, req)
	if err != nil {
		if err == ErrZoneFull {
			return httpresponse.Error(c, http.StatusConflict, "Zone is full, no available spots", nil)
		}
		return httpresponse.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusCreated, "Reservation confirmed successfully", res)
}

func (h *Handler) GetMyReservations(c echo.Context) error {
	userID := c.Get("userID").(uint)

	reservations, err := h.service.GetMyReservations(userID)
	if err != nil {
		return httpresponse.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

func (h *Handler) CancelReservation(c echo.Context) error {
	userID := c.Get("userID").(uint)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid reservation ID", nil)
	}

	if err := h.service.CancelReservation(userID, uint(id)); err != nil {
		if err.Error() == "forbidden" {
			return httpresponse.Error(c, http.StatusForbidden, "You can only cancel your own reservations", nil)
		}
		if err.Error() == "reservation not found" {
			return httpresponse.Error(c, http.StatusNotFound, "Reservation not found", nil)
		}
		return httpresponse.Error(c, http.StatusBadRequest, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

func (h *Handler) GetAllReservations(c echo.Context) error {
	if c.Get("userRole").(string) != "admin" {
		return httpresponse.Error(c, http.StatusForbidden, "Admin access required", nil)
	}

	reservations, err := h.service.GetAllReservations()
	if err != nil {
		return httpresponse.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusOK, "All reservations retrieved successfully", reservations)
}