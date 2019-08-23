package model

import (
	"net/http"
)

type HttpError struct {
	error
	Status  int `json:"-"`
	Errors []string `json:"errors"`
}

func NewInternalServerError(message... string) *HttpError {
	return &HttpError{
		Status: http.StatusInternalServerError,
		Errors: message,
	}
}

func NewBadRequestError(message... string) *HttpError {
	return &HttpError{
		Status: http.StatusBadRequest,
		Errors: message,
	}
}

func NewForbiddenError(message... string) *HttpError {
	return &HttpError{
		Status: http.StatusForbidden,
		Errors: message,
	}
}

func NewUnauthorizedError(message... string) *HttpError {
	return &HttpError{
		Status: http.StatusUnauthorized,
		Errors: message,
	}
}
