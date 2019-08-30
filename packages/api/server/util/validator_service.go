package util

import (
	"context"
	"encoding/json"
	"fmt"
	serverModel "github.com/honerlaw/mentordoc/server/model"
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

func (v *ValidatorService) ParseAndValidate(req *http.Request, model interface{}) (interface{}, *serverModel.HttpError) {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(model)
	if err != nil {
		log.Print(err)
		return nil, serverModel.NewBadRequestError("failed to validate reque")
	}

	err = v.validator.Struct(model)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errorMessages := make([]string, len(validationErrors))
		for i, err := range validationErrors {
			errorMessages[i] = v.formatValidationError(err)
 		}

		return nil, serverModel.NewBadRequestError(errorMessages...)
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
	case "min":
		return strings.ToLower(fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()))
	default:
		log.Printf("unhandled validator type of %s", err.Tag())
		return "validation failed"
	}
}
