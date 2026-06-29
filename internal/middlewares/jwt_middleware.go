package middlewares

import (
	"net/http"
	"spotsync/internal/auth"
	"spotsync/internal/httpresponse"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return httpresponse.Error(c, http.StatusUnauthorized, "Missing or invalid token", nil)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			return httpresponse.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
		}

		c.Set("userID", claims.ID)
		c.Set("userRole", claims.Role)

		return next(c)
	}
}