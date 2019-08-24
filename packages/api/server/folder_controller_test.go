package server_test

import (
	"fmt"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationCreateFolderFailsBecauseNotAuthenticated(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	req := &model.FolderCreateRequest{
		OrganizationId: "10",
		Name:           "test-name",
		ParentFolderId: nil,
	}

	status, _, err := PostItTest(&PostOptions{
		Path: "/folder",
	}, req, &model.AclWrappedModel{})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestIntegrationCreateFolderFailsBecauseCanNotCreateFolderInOrganization(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)

	req := &model.FolderCreateRequest{
		OrganizationId: "10",
		Name:           "test-name",
		ParentFolderId: nil,
	}

	status, _, err := PostItTest(&PostOptions{
		Path: "/folder", Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
	}, req, &model.AclWrappedModel{})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, status)
}

func TestIntegrationCreateFolder(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)


	// create a new org
	org, err := testServer.OrganizationService.Create("test")
	assert.Nil(t, err)

	// link to the new org
	err = testServer.AclService.LinkUserToRole(authData.user, "organization:owner", org.Id)
	assert.Nil(t, err)

	req := &model.FolderCreateRequest{
		OrganizationId: org.Id,
		Name:           "test-name",
		ParentFolderId: nil,
	}

	status, resp, err := PostItTest(&PostOptions{
		Path: "/folder", Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
	}, req, &model.AclWrappedModel{})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)

	r := resp.(*model.AclWrappedModel)
	assert.Equal(t, r.Actions, []string{"create:document", "delete", "modify", "view"})
}
