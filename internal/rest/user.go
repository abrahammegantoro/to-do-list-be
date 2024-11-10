package rest

import (
	"context"
	"net/http"

	"github.com/abrahammegantoro/to-do-list-be/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	Login(ctx context.Context, auth *domain.AuthCredentials) (domain.User, string, error)
	Register(ctx context.Context, user *domain.User) (string, error)
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(e *echo.Group, svc UserService) {
	handler := &UserHandler{
		Service: svc,
	}

	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)
}

func (u *UserHandler) Login(c echo.Context) (err error) {
	var auth domain.AuthCredentials
	if err := c.Bind(&auth); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	var ok bool
	if ok, err = isRequestValid(&auth); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	
	user, token, err := u.Service.Login(ctx, &auth)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
		"message": "success",
		"data": map[string]interface{}{
			"user": user,
			"token": token,
		},
	})
}

func (u *UserHandler) Register(c echo.Context) (err error) {
	var user domain.User
	if err = c.Bind(&user); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	token, err := u.Service.Register(ctx, &user)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"status":  getStatusCode(err),
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": http.StatusOK,
		"message": "success",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}
