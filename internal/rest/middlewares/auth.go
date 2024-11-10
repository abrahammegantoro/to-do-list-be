package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/abrahammegantoro/to-do-list-be/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (res domain.User, err error)
}

func AuthMiddleware(user UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, "Missing authorization header")
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				logrus.Error("Invalid token")

				return c.JSON(http.StatusUnauthorized, "Invalid token")
			}

			tokenStr := tokenParts[1]
			secret := []byte(os.Getenv("JWT_SECRET"))

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return secret, nil
			})

			if err != nil || !token.Valid {
				logrus.Error(err)
				return c.JSON(http.StatusUnauthorized, err.Error())
			}

			userIdF := token.Claims.(jwt.MapClaims)["id"].(float64)
			userId := int64(userIdF)

			c.Set("userId", userId)
			return next(c)
		}
	}
}
