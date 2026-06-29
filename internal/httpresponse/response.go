package httpresponse

import "github.com/labstack/echo/v4"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, code int, message string, errors interface{}) error {
	return c.JSON(code, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}