package server

import "net/http"

type HttpError struct {
	error
	Status  int `json:"-"`
	Errors []string `json:"errors"`
}

func NewInternalServerError(message string) *HttpError {
	return &HttpError{
		Status: http.StatusInternalServerError,
		Errors: []string{message},
	}
}

func NewBadRequestError(message string) *HttpError {
	return &HttpError{
		Status: http.StatusBadRequest,
		Errors: []string{message},
	}
}

func (err *HttpError) Write(w http.ResponseWriter) {
	WriteJsonToResponse(w, err.Status, err)
}

func WriteHttpError(w http.ResponseWriter, err interface{}) {
	if httpError, ok := err.(*HttpError); ok {
		WriteJsonToResponse(w, httpError.Status, httpError)
	} else {
		WriteJsonToResponse(w, 500, &HttpError{
			Errors: []string{err.(error).Error()},
		})
	}
}