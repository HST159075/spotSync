package zone

import (
	"net/http"
	zonedto "spotsync/internal/domain/zone/dto"
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

func (h *Handler) CreateZone(c echo.Context) error {
	if c.Get("userRole").(string) != "admin" {
		return httpresponse.Error(c, http.StatusForbidden, "Admin access required", nil)
	}

	var req zonedto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zone, err := h.service.CreateZone(req)
	if err != nil {
		return httpresponse.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusCreated, "Parking zone created successfully", zone)
}

func (h *Handler) GetAllZones(c echo.Context) error {
	zones, err := h.service.GetAllZones()
	if err != nil {
		return httpresponse.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}
	return httpresponse.Success(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *Handler) GetZoneByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	zone, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		return httpresponse.Error(c, http.StatusNotFound, "Zone not found", nil)
	}

	return httpresponse.Success(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}

func (h *Handler) UpdateZone(c echo.Context) error {
	if c.Get("userRole").(string) != "admin" {
		return httpresponse.Error(c, http.StatusForbidden, "Admin access required", nil)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	var req zonedto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid request body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zone, err := h.service.UpdateZone(uint(id), req)
	if err != nil {
		if err.Error() == "zone not found" {
			return httpresponse.Error(c, http.StatusNotFound, "Zone not found", nil)
		}
		return httpresponse.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusOK, "Parking zone updated successfully", zone)
}

func (h *Handler) DeleteZone(c echo.Context) error {
	if c.Get("userRole").(string) != "admin" {
		return httpresponse.Error(c, http.StatusForbidden, "Admin access required", nil)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpresponse.Error(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	if err := h.service.DeleteZone(uint(id)); err != nil {
		if err.Error() == "zone not found" {
			return httpresponse.Error(c, http.StatusNotFound, "Zone not found", nil)
		}
		return httpresponse.Error(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusOK, "Parking zone deleted successfully", nil)
}