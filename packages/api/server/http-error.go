package server

import "net/http"

type HttpError struct {
	Status  int `json:"-"`
	Errors []string `json:"errors"`
}

func (err *HttpError) Write(w http.ResponseWriter) {
	WriteJsonToResponse(w, err.Status, err)
}

func WriteHttpError(w http.ResponseWriter, err interface{}) {
	if httpError, ok := err.(*HttpError); ok {
		WriteJsonToResponse(w, httpError.Status, httpError)
	} else {
		WriteJsonToResponse(w, 500, err.(error).Error())
	}
}