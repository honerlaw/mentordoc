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
	req := &model.UserSigninRequest{}
	status, resp, err := PostItTest(&PostOptions{Path: "/user/auth"}, req, &model.HttpError{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*model.HttpError).Errors, 2)
}

func TestIntegrationSigninUserDoesntExist(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	req := &model.UserSigninRequest{
		Email:    "foo@bar.com",
		Password: "baz",
	}
	status, resp, err := PostItTest(&PostOptions{Path: "/user/auth"}, req, &model.HttpError{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*model.HttpError).Errors, 1)
}

func TestIntegrationSignupUserDoesntExist(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	req := &model.UserSignupRequest{
		Email:    "foo@bar.com",
		Password: "foobarbaz",
	}
	status, resp, err := PostItTest(&PostOptions{Path: "/user"}, req, &model.AuthenticationResponse{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).RefreshToken)
}

func TestIntegrationSignupAndThenSignin(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	req := &model.UserSignupRequest{
		Email:    "footest@bar.com",
		Password: "foobarbaz",
	}
	status, resp, err := PostItTest(&PostOptions{Path: "/user"}, req, &model.AuthenticationResponse{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	signinReq := &model.UserSigninRequest{
		Email:    "footest@bar.com",
		Password: "foobarbaz",
	}
	status, resp, err = PostItTest(&PostOptions{Path: "/user/auth"}, signinReq, &model.AuthenticationResponse{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).RefreshToken)
}
