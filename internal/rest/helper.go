package rest

import (
	"errors"
	"net/http"

	"github.com/abrahammegantoro/to-do-list-be/domain"
	validator "github.com/go-playground/validator/v10"
)

type ResponseError struct {
	Message string `json:"message"`
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrCredential:
		return http.StatusUnauthorized
	case domain.ErrBadParamInput:
		return http.StatusBadRequest
	case domain.ErrUsernameTaken:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func isRequestValid(v interface{}) (bool, error) {
	validate := validator.New()
	switch v := v.(type) {
	case *domain.Todo:
		err := validate.Struct(v)
		if err != nil {
			return false, err
		}
	case *domain.User:
		err := validate.Struct(v)
		if err != nil {
			return false, err
		}
	case *domain.AuthCredentials:
		err := validate.Struct(v)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.New("unsupported type")
	}
	return true, nil
}
