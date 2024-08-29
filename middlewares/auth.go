package middlewares

import (
	"net/http"
	"os"
	"v-games-ip-ph2-ftgo/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtSecret := os.Getenv("KEY")
			tokenString := c.Request().Header.Get("authorization")

			if tokenString == "" {
				return utils.HandleError(c, utils.NewUnauthorizedError("Missing or invalid token"))
			}

			// Parse token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return utils.HandleError(c, utils.NewUnauthorizedError("Invalid or expired token"))
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userRole := claims["role"].(string)

				if role != "both" {
					// Otorisasi berdasarkan role
					if role == "admin" && userRole != "admin" {
						return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden for non-admin"})
					}

					if role == "user" && userRole != "user" {
						return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden for admin"})
					}
				}

				// Set token di konteks untuk digunakan di handler
				c.Set("user_id", claims["id"])
				c.Set("email", claims["email"])
				c.Set("role", claims["role"])
				return next(c)
			} else {
				return utils.HandleError(c, utils.NewUnauthorizedError("Invalid token claims"))
			}
		}
	}
}
