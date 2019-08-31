package server_test

import (
	"fmt"
	"github.com/honerlaw/mentordoc/server/http/request"
	"github.com/honerlaw/mentordoc/server/http/response"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationSigninValidationFailure(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := test.Request(&test.RequestOptions{
		Method:        "POST",
		Path:          "/user/auth",
		Body:          &request.UserSigninRequest{},
		ResponseModel: &shared.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*shared.HttpError).Errors, 2)
}

func TestIntegrationSigninUserDoesntExist(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/user/auth",
		Body: &request.UserSigninRequest{
			Email:    "foo@bar.com",
			Password: "baz",
		},
		ResponseModel: &shared.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*shared.HttpError).Errors, 1)
}

func TestIntegrationSignupUserDoesntExist(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/user",
		Body: &request.UserSignupRequest{
			Email:    "foo@bar.com",
			Password: "foobarbaz",
		}, ResponseModel: &response.AuthenticationResponse{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)
	assert.NotEmpty(t, resp.(*response.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*response.AuthenticationResponse).RefreshToken)
}

func TestIntegrationSignupAndThenSignin(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/user",
		Body: &request.UserSignupRequest{
			Email:    "footest@bar.com",
			Password: "foobarbaz",
		},
		ResponseModel: &response.AuthenticationResponse{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	status, resp, err = test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/user/auth",
		Body: &request.UserSigninRequest{
			Email:    "footest@bar.com",
			Password: "foobarbaz",
		},
		ResponseModel: &response.AuthenticationResponse{},
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	assert.NotEmpty(t, resp.(*response.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*response.AuthenticationResponse).RefreshToken)
}

func TestIntegrationGetCurrentUser(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	status, resp, err := test.Request(&test.RequestOptions{
		Method: "OPTIONS",
		Path:   "/user",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		ResponseModel: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	panic(string(resp.([]byte)))
}