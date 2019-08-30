package server_test

import (
	"fmt"
	"github.com/honerlaw/mentordoc/server/http/request"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationCreateDocumentFailsBecauseNotAuthenticated(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	status, _, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/document",
		Body: &request.DocumentCreateRequest{
			OrganizationId: "10",
			Name:           "test-name",
			Content:        "some content",
			FolderId:       nil,
		},
		ResponseModel: &shared.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestIntegrationCreateDocumentFailsBecauseCanNotFindOrganization(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	status, _, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.DocumentCreateRequest{
			OrganizationId: "10",
			Name:           "test-name",
			Content:        "test content",
			FolderId:       nil,
		},
		ResponseModel: &shared.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

func TestIntegrationCreateDocumentFailsBecauseCanNotCreateFolderInOrganization(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	org, err := testData.TestServer.OrganizationService.Create("test")
	assert.Nil(t, err)

	status, _, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.DocumentCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			Content:        "test",
			FolderId:       nil,
		},
		ResponseModel: &shared.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, status)
}

func TestIntegrationCreateDocumentInOrganization(t *testing.T) {
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
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.DocumentCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			Content:        "test content",
			FolderId:       nil,
		},
		ResponseModel: &acl.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)

	r := resp.(*acl.AclWrappedModel)
	assert.Equal(t, []string{"delete", "modify", "view"}, r.Actions)
	assert.Nil(t, r.Model.(map[string]interface{})["folderId"])
}

func TestIntegrationCreateDocumentInFolder(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t,testData)

	// create a new org
	org, err := testData.TestServer.OrganizationService.Create("test")
	assert.Nil(t, err)

	// link to the new org
	err = testData.TestServer.AclService.LinkUserToRole(authData.User, "organization:owner", org.Id)
	assert.Nil(t, err)

	folder, err := testData.TestServer.FolderService.Create(authData.User, "test", org.Id, nil)
	assert.Nil(t, err)

	status, resp, err := test.Request(&test.RequestOptions{
		Method: "POST",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.DocumentCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			Content:        "test content",
			FolderId:       &folder.Id,
		},
		ResponseModel: &acl.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)

	r := resp.(*acl.AclWrappedModel)
	assert.Equal(t, r.Actions, []string{"delete", "modify", "view"})
	assert.Equal(t, r.Model.(map[string]interface{})["folderId"], folder.Id)
}

func TestIntegrationListDocumentInOrganization(t *testing.T) {
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
	_, err = testData.TestServer.DocumentService.Create(authData.User, org.Id, nil, "test document", "test content")
	assert.Nil(t, err)

	// create a new org that we will not have access to
	orgTwo, err := testData.TestServer.OrganizationService.Create("test two")
	assert.Nil(t, err)
	err = testData.TestServer.AclService.LinkUserToRole(authDataTwo.User, "organization:owner", orgTwo.Id)
	assert.Nil(t, err)
	_, err = testData.TestServer.DocumentService.Create(authDataTwo.User, orgTwo.Id, nil, "test document 2", "test content 2")
	assert.Nil(t, err)

	aclWrappedModels := make([]acl.AclWrappedModel, 0)
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s", org.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]acl.AclWrappedModel)

	assert.Len(t, r, 1)
	assert.Equal(t, "test document", r[0].Model.(map[string]interface{})["name"])
}

func TestIntegrationListDocumentInOrganizationAndSpecificFolder(t *testing.T) {
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
	folder, err := testData.TestServer.FolderService.Create(authData.User, "test folder", org.Id, nil)
	assert.Nil(t, err)
	_, err = testData.TestServer.DocumentService.Create(authData.User, org.Id, &folder.Id, "test document", "test content")
	assert.Nil(t, err)
	_, err = testData.TestServer.DocumentService.Create(authData.User, org.Id, nil, "test document no folder", "test content")
	assert.Nil(t, err)

	// create a new org that we will not have access to
	orgTwo, err := testData.TestServer.OrganizationService.Create("test two")
	assert.Nil(t, err)
	err = testData.TestServer.AclService.LinkUserToRole(authDataTwo.User, "organization:owner", orgTwo.Id)
	assert.Nil(t, err)
	_, err = testData.TestServer.DocumentService.Create(authDataTwo.User, orgTwo.Id, nil, "test document 2", "test content 2")
	assert.Nil(t, err)

	aclWrappedModels := make([]acl.AclWrappedModel, 0)
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s?folderId=%s", org.Id, folder.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]acl.AclWrappedModel)

	assert.Len(t, r, 1)
	assert.Equal(t, "test document", r[0].Model.(map[string]interface{})["name"])
}

func TestIntegrationListDocumentInOrganizationAndSpecificFolderWithPaginaton(t *testing.T) {
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
	folder, err := testData.TestServer.FolderService.Create(authData.User, "test folder", org.Id, nil)
	assert.Nil(t, err)

	for i := 0; i < 25; i++ {
		_, err = testData.TestServer.DocumentService.Create(authData.User, org.Id, &folder.Id, fmt.Sprintf("test document %d", i), "test content")
		assert.Nil(t, err)
		_, err = testData.TestServer.DocumentService.Create(authData.User, org.Id, nil, fmt.Sprintf("test document no folder %d", i), "test content")
		assert.Nil(t, err)
	}

	// create a new org that we will not have access to
	orgTwo, err := testData.TestServer.OrganizationService.Create("test two")
	assert.Nil(t, err)
	err = testData.TestServer.AclService.LinkUserToRole(authDataTwo.User, "organization:owner", orgTwo.Id)
	assert.Nil(t, err)
	_, err = testData.TestServer.DocumentService.Create(authDataTwo.User, orgTwo.Id, nil, "test document 2", "test content 2")
	assert.Nil(t, err)

	aclWrappedModels := make([]acl.AclWrappedModel, 0)
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s?folderId=%s&page=0&count=5", org.Id, folder.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]acl.AclWrappedModel)

	assert.Len(t, r, 5)
	assert.Equal(t, "test document 0", r[0].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 1", r[1].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 10", r[2].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 11", r[3].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 12", r[4].Model.(map[string]interface{})["name"])

	// fetch the next page
	status, resp, err = test.Request(&test.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s?folderId=%s&page=1&count=5", org.Id, folder.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r = *resp.(*[]acl.AclWrappedModel)

	assert.Len(t, r, 5)
	assert.Equal(t, "test document 13", r[0].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 14", r[1].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 15", r[2].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 16", r[3].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 17", r[4].Model.(map[string]interface{})["name"])
}

func TestIntegrationUpdateDocument(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	org, err := testData.TestServer.OrganizationService.Create("test")
	assert.Nil(t, err)
	err = testData.TestServer.AclService.LinkUserToRole(authData.User, "organization:owner", org.Id)
	assert.Nil(t, err)
	folder, err := testData.TestServer.FolderService.Create(authData.User, "test folder", org.Id, nil)
	assert.Nil(t, err)
	document, err := testData.TestServer.DocumentService.Create(authData.User, org.Id, &folder.Id, "test document", "test content")
	assert.Nil(t, err)

	assert.Equal(t, document.Name, "test document")
	assert.Equal(t, document.Content.Content, "test content")

	status, resp, err := test.Request(&test.RequestOptions{
		Method: "PUT",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		Body: &request.DocumentUpdateRequest{
			DocumentId: document.Id,
			Name:       "new name",
			Content:    "new content",
		},
		ResponseModel: &acl.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := resp.(*acl.AclWrappedModel)

	assert.Equal(t, document.Id, r.Model.(map[string]interface{})["id"])
	assert.Equal(t, "new name", r.Model.(map[string]interface{})["name"])
	assert.Equal(t, "new content", r.Model.(map[string]interface{})["content"].(map[string]interface{})["content"])
}
