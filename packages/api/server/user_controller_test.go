package server_test

import (
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationSigninValidationFailure(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := Request(&RequestOptions{
		Method:        "POST",
		Path:          "/user/auth",
		Body:          &model.UserSigninRequest{},
		ResponseModel: &model.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*model.HttpError).Errors, 2)
}

func TestIntegrationSigninUserDoesntExist(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/user/auth",
		Body: &model.UserSigninRequest{
			Email:    "foo@bar.com",
			Password: "baz",
		},
		ResponseModel: &model.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*model.HttpError).Errors, 1)
}

func TestIntegrationSignupUserDoesntExist(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/user",
		Body: &model.UserSignupRequest{
			Email:    "foo@bar.com",
			Password: "foobarbaz",
		}, ResponseModel: &model.AuthenticationResponse{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).RefreshToken)
}

func TestIntegrationSignupAndThenSignin(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/user",
		Body: &model.UserSignupRequest{
			Email:    "footest@bar.com",
			Password: "foobarbaz",
		},
		ResponseModel: &model.AuthenticationResponse{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	status, resp, err = Request(&RequestOptions{
		Method: "POST",
		Path:   "/user/auth",
		Body: &model.UserSigninRequest{
			Email:    "footest@bar.com",
			Password: "foobarbaz",
		},
		ResponseModel: &model.AuthenticationResponse{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).RefreshToken)
}
