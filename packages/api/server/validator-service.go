package server

import (
	"context"
	"encoding/json"
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"reflect"
)

const RequestModelContextKey = "request_model"

type ValidatorService struct {
	validator *validator.Validate
}

func NewValidatorService() *ValidatorService {
	return &ValidatorService{
		validator: validator.New(),
	}
}

func (v *ValidatorService) ParseAndValidate(req *http.Request, model interface{}) (interface{}, error)  {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(model)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to parse request body")
	}

	err = v.validator.Struct(model)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to validate request body")
	}

	return model, nil
}

func (v *ValidatorService) Middleware(model interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			modelPtr := reflect.New(reflect.TypeOf(model)).Interface()
			modelPtr, err := v.ParseAndValidate(req, modelPtr)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}

			ctx := context.WithValue(req.Context(), RequestModelContextKey, modelPtr)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func (v *ValidatorService) GetModelFromRequest(req *http.Request) interface{} {
	return req.Context().Value(RequestModelContextKey)
}
