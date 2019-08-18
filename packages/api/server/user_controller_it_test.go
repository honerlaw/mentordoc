package server

import (
	"bytes"
	"encoding/json"
	model2 "github.com/honerlaw/mentordoc/server/model"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
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
	req := &UserSigninRequest{}
	status, model, err := post(GetTestServerAddress("/user/auth"), req, &model2.HttpError{})
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Assert(t, is.Len(model.(*model2.HttpError).Errors, 2))
}

func TestSigninUserDoesntExist(t *testing.T) {
	req := &UserSigninRequest{
		Email:    "foo@bar.com",
		Password: "baz",
	}
	status, model, err := post(GetTestServerAddress("/user/auth"), req, &model2.HttpError{})
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Assert(t, is.Len(model.(*model2.HttpError).Errors, 1))
}

func TestSignupUserDoesntExist(t *testing.T) {
	req := &UserSignupRequest{
		Email:    "foo@bar.com",
		Password: "foobarbaz",
	}
	status, model, err := post(GetTestServerAddress("/user"), req, &model2.User{})
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, model.(*model2.User).Email, req.Email)
}

func TestSignupAndThenSignin(t *testing.T) {
	req := &UserSignupRequest{
		Email:    "footest@bar.com",
		Password: "foobarbaz",
	}
	status, model, err := post(GetTestServerAddress("/user"), req, &model2.User{})
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, status, http.StatusOK)

	signinReq := &UserSigninRequest{
		Email:    "footest@bar.com",
		Password: "foobarbaz",
	}
	status, model, err = post(GetTestServerAddress("/user/auth"), signinReq, &model2.User{})
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, model.(*model2.User).Email, req.Email)
}
