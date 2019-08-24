package server_test

import (
	"encoding/json"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationCreateFolder(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	req := &model.FolderCreateRequest{
		OrganizationId: "10",
		Name: "test-name",
		ParentFolderId: nil,
	}

	status, resp, err := PostItTest("/folder", req, &model.AclWrappedModel{})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)

	data, _ := json.Marshal(resp)

	panic(string(data))
}