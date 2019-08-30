package server_test

import (
	"fmt"
	"github.com/honerlaw/mentordoc/server/http/request"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationCreateFolderFailsBecauseNotAuthenticated(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	status, _, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/folder",
		Body: &request.FolderCreateRequest{
			OrganizationId: "10",
			Name:           "test-name",
			ParentFolderId: nil,
		},
		ResponseModel: &acl.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestIntegrationCreateFolderFailsBecauseCanNotFindOrganization(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	status, _, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/folder",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.FolderCreateRequest{
			OrganizationId: "10",
			Name:           "test-name",
			ParentFolderId: nil,
		},
		ResponseModel: &acl.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

func TestIntegrationCreateFolderFailsBecauseCanNotCreateFolderInOrganization(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	org, err := testData.TestServer.OrganizationService.Create("test")
	assert.Nil(t, err)

	status, _, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/folder",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.FolderCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			ParentFolderId: nil,
		},
		ResponseModel: &acl.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, status)
}

func TestIntegrationCreateFolder(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	// create a new org
	org, err := testData.TestServer.OrganizationService.Create("test")
	assert.Nil(t, err)

	// link to the new org
	err = testData.TestServer.AclService.LinkUserToRole(authData.User, "organization:owner", org.Id)
	assert.Nil(t, err)

	status, resp, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/folder",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.FolderCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			ParentFolderId: nil,
		},
		ResponseModel: &acl.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)

	r := resp.(*acl.AclWrappedModel)
	assert.Equal(t, r.Actions, []string{"create:document", "delete", "modify", "view", "view:document"})
}

func TestIntegrationListFolders(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)
	authDataTwo := test.SetupAuthentication(t, testData)

	// create a new org, that we will add the accessible folder to
	org, err := testData.TestServer.OrganizationService.Create("test")
	assert.Nil(t, err)
	err = testData.TestServer.AclService.LinkUserToRole(authData.User, "organization:owner", org.Id)
	assert.Nil(t, err)
	_, err = testData.TestServer.FolderService.Create(authData.User, "test folder", org.Id, nil)
	assert.Nil(t, err)

	// create a new org that we will not have access to
	orgTwo, err := testData.TestServer.OrganizationService.Create("test two")
	assert.Nil(t, err)
	err = testData.TestServer.AclService.LinkUserToRole(authDataTwo.User, "organization:owner", orgTwo.Id)
	assert.Nil(t, err)
	_, err = testData.TestServer.FolderService.Create(authDataTwo.User, "test folder 2", orgTwo.Id, nil)
	assert.Nil(t, err)

	aclWrappedModels := make([]acl.AclWrappedModel, 0)
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/folder/list/%s", org.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]acl.AclWrappedModel)

	assert.Len(t, r, 1)
	assert.Equal(t, "test folder", r[0].Model.(map[string]interface{})["name"])
}
