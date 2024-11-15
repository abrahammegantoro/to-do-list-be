package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Item is not found")
	ErrConflict            = errors.New("your Item already exist")
	ErrBadParamInput       = errors.New("given Param is not valid")
	ErrCredential          = errors.New("your Credential is invalid")
	ErrUsernameTaken	   = errors.New("your Username is already taken")
)
