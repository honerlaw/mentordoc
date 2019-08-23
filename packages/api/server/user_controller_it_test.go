package server

import (
	"bytes"
	"encoding/json"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func post(url string, req interface{}, resp interface{}) (int, interface{}, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return -1, nil, err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return -1, nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, nil, err
	}

	if resp == true {
		return response.StatusCode, data, nil
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return -1, nil, err

	}

	return response.StatusCode, resp, nil
}

func TestSigninValidationFailure(t *testing.T) {
	req := &model.UserSigninRequest{}
	status, resp, err := post(GetTestServerAddress("/user/auth"), req, &model.HttpError{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*model.HttpError).Errors, 2)
}

func TestSigninUserDoesntExist(t *testing.T) {
	req := &model.UserSigninRequest{
		Email:    "foo@bar.com",
		Password: "baz",
	}
	status, resp, err := post(GetTestServerAddress("/user/auth"), req, &model.HttpError{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Len(t, resp.(*model.HttpError).Errors, 1)
}

func TestSignupUserDoesntExist(t *testing.T) {
	req := &model.UserSignupRequest{
		Email:    "foo@bar.com",
		Password: "foobarbaz",
	}
	status, resp, err := post(GetTestServerAddress("/user"), req, &model.AuthenticationResponse{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).RefreshToken)
}

func TestSignupAndThenSignin(t *testing.T) {
	req := &model.UserSignupRequest{
		Email:    "footest@bar.com",
		Password: "foobarbaz",
	}
	status, resp, err := post(GetTestServerAddress("/user"), req, &model.AuthenticationResponse{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	signinReq := &model.UserSigninRequest{
		Email:    "footest@bar.com",
		Password: "foobarbaz",
	}
	status, resp, err = post(GetTestServerAddress("/user/auth"), signinReq, &model.AuthenticationResponse{})
	assert.Nil(t, err)
	assert.Equal(t, status, http.StatusOK)

	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).AccessToken)
	assert.NotEmpty(t, resp.(*model.AuthenticationResponse).RefreshToken)
}
