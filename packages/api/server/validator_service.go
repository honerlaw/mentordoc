package server

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"reflect"
	"strings"
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

func (v *ValidatorService) ParseAndValidate(req *http.Request, model interface{}) (interface{}, *HttpError) {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(model)
	if err != nil {
		log.Print(err)
		return nil, &HttpError{Errors: []string{"failed to parse request body"}}
	}

	err = v.validator.Struct(model)
	if err != nil {

		validatorError := &HttpError{
			Status: http.StatusBadRequest,
			Errors: make([]string, 0),
		}

		for _, err := range err.(validator.ValidationErrors) {
			validatorError.Errors = append(validatorError.Errors, v.formatValidationError(err))
		}

		return nil, validatorError
	}

	return model, nil
}

func (v *ValidatorService) Middleware(model interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			modelPtr := reflect.New(reflect.TypeOf(model)).Interface()
			modelPtr, err := v.ParseAndValidate(req, modelPtr)
			if err != nil {
				WriteJsonToResponse(w, err.Status, err)
				return
			}

			ctx := context.WithValue(req.Context(), RequestModelContextKey, modelPtr)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func (v *ValidatorService) GetModelFromRequest(req *http.Request) interface{} {
	return req.Context().Value(RequestModelContextKey)
}

func (v *ValidatorService) formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return strings.ToLower(fmt.Sprintf("%s is required", err.Field()))
	case "email":
		return strings.ToLower(fmt.Sprintf("%s must be an email address", err.Field()))
	default:
		log.Printf("unhandled validator type of %s", err.Tag())
		return "validation failed"
	}
}
