package server_test

import (
	"fmt"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationCreateDocumentFailsBecauseNotAuthenticated(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	status, _, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/document",
		Body: &model.DocumentCreateRequest{
			OrganizationId: "10",
			Name:           "test-name",
			Content:        "some content",
			FolderId:       nil,
		},
		ResponseModel: &model.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestIntegrationCreateDocumentFailsBecauseCanNotFindOrganization(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)

	status, _, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		Body: &model.DocumentCreateRequest{
			OrganizationId: "10",
			Name:           "test-name",
			Content:        "test content",
			FolderId:       nil,
		},
		ResponseModel: &model.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

func TestIntegrationCreateDocumentFailsBecauseCanNotCreateFolderInOrganization(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)

	org, err := testServer.OrganizationService.Create("test")
	assert.Nil(t, err)

	status, _, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		Body: &model.DocumentCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			Content:        "test",
			FolderId:       nil,
		},
		ResponseModel: &model.HttpError{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, status)
}

func TestIntegrationCreateDocumentInOrganization(t *testing.T) {
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

	status, resp, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		Body: &model.DocumentCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			Content:        "test content",
			FolderId:       nil,
		},
		ResponseModel: &model.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)

	r := resp.(*model.AclWrappedModel)
	assert.Equal(t, r.Actions, []string{"delete", "modify", "view"})
	assert.Nil(t, r.Model.(map[string]interface{})["folderId"])
}

func TestIntegrationCreateDocumentInFolder(t *testing.T) {
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

	folder, err := testServer.FolderService.Create(authData.user, "test", org.Id, nil)
	assert.Nil(t, err)

	status, resp, err := Request(&RequestOptions{
		Method: "POST",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		Body: &model.DocumentCreateRequest{
			OrganizationId: org.Id,
			Name:           "test-name",
			Content:        "test content",
			FolderId:       &folder.Id,
		},
		ResponseModel: &model.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)

	r := resp.(*model.AclWrappedModel)
	assert.Equal(t, r.Actions, []string{"delete", "modify", "view"})
	assert.Equal(t, r.Model.(map[string]interface{})["folderId"], folder.Id)
}

func TestIntegrationListDocumentInOrganization(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)
	authDataTwo := SetupAuthentication(t)

	// create a new org, that we will add the accessible folder to
	org, err := testServer.OrganizationService.Create("test")
	assert.Nil(t, err)
	err = testServer.AclService.LinkUserToRole(authData.user, "organization:owner", org.Id)
	assert.Nil(t, err)
	_, err = testServer.DocumentService.Create(authData.user, org.Id, nil, "test document", "test content")
	assert.Nil(t, err)

	// create a new org that we will not have access to
	orgTwo, err := testServer.OrganizationService.Create("test two")
	assert.Nil(t, err)
	err = testServer.AclService.LinkUserToRole(authDataTwo.user, "organization:owner", orgTwo.Id)
	assert.Nil(t, err)
	_, err = testServer.DocumentService.Create(authDataTwo.user, orgTwo.Id, nil, "test document 2", "test content 2")
	assert.Nil(t, err)

	aclWrappedModels := make([]model.AclWrappedModel, 0)
	status, resp, err := Request(&RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s", org.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]model.AclWrappedModel)

	assert.Len(t, r, 1)
	assert.Equal(t, "test document", r[0].Model.(map[string]interface{})["name"])
}

func TestIntegrationListDocumentInOrganizationAndSpecificFolder(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)
	authDataTwo := SetupAuthentication(t)

	// create a new org, that we will add the accessible folder to
	org, err := testServer.OrganizationService.Create("test")
	assert.Nil(t, err)
	err = testServer.AclService.LinkUserToRole(authData.user, "organization:owner", org.Id)
	assert.Nil(t, err)
	folder, err := testServer.FolderService.Create(authData.user, "test folder", org.Id, nil)
	assert.Nil(t, err)
	_, err = testServer.DocumentService.Create(authData.user, org.Id, &folder.Id, "test document", "test content")
	assert.Nil(t, err)
	_, err = testServer.DocumentService.Create(authData.user, org.Id, nil, "test document no folder", "test content")
	assert.Nil(t, err)

	// create a new org that we will not have access to
	orgTwo, err := testServer.OrganizationService.Create("test two")
	assert.Nil(t, err)
	err = testServer.AclService.LinkUserToRole(authDataTwo.user, "organization:owner", orgTwo.Id)
	assert.Nil(t, err)
	_, err = testServer.DocumentService.Create(authDataTwo.user, orgTwo.Id, nil, "test document 2", "test content 2")
	assert.Nil(t, err)

	aclWrappedModels := make([]model.AclWrappedModel, 0)
	status, resp, err := Request(&RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s?folderId=%s", org.Id, folder.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]model.AclWrappedModel)

	assert.Len(t, r, 1)
	assert.Equal(t, "test document", r[0].Model.(map[string]interface{})["name"])
}

func TestIntegrationListDocumentInOrganizationAndSpecificFolderWithPaginaton(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)
	authDataTwo := SetupAuthentication(t)

	// create a new org, that we will add the accessible folder to
	org, err := testServer.OrganizationService.Create("test")
	assert.Nil(t, err)
	err = testServer.AclService.LinkUserToRole(authData.user, "organization:owner", org.Id)
	assert.Nil(t, err)
	folder, err := testServer.FolderService.Create(authData.user, "test folder", org.Id, nil)
	assert.Nil(t, err)

	for i := 0; i < 25; i++ {
		_, err = testServer.DocumentService.Create(authData.user, org.Id, &folder.Id, fmt.Sprintf("test document %d", i), "test content")
		assert.Nil(t, err)
		_, err = testServer.DocumentService.Create(authData.user, org.Id, nil, fmt.Sprintf("test document no folder %d", i), "test content")
		assert.Nil(t, err)
	}

	// create a new org that we will not have access to
	orgTwo, err := testServer.OrganizationService.Create("test two")
	assert.Nil(t, err)
	err = testServer.AclService.LinkUserToRole(authDataTwo.user, "organization:owner", orgTwo.Id)
	assert.Nil(t, err)
	_, err = testServer.DocumentService.Create(authDataTwo.user, orgTwo.Id, nil, "test document 2", "test content 2")
	assert.Nil(t, err)

	aclWrappedModels := make([]model.AclWrappedModel, 0)
	status, resp, err := Request(&RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s?folderId=%s&page=0&count=5", org.Id, folder.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]model.AclWrappedModel)

	assert.Len(t, r, 5)
	assert.Equal(t, "test document 0", r[0].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 1", r[1].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 10", r[2].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 11", r[3].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 12", r[4].Model.(map[string]interface{})["name"])

	// fetch the next page
	status, resp, err = Request(&RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/document/list/%s?folderId=%s&page=1&count=5", org.Id, folder.Id),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r = *resp.(*[]model.AclWrappedModel)

	assert.Len(t, r, 5)
	assert.Equal(t, "test document 13", r[0].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 14", r[1].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 15", r[2].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 16", r[3].Model.(map[string]interface{})["name"])
	assert.Equal(t, "test document 17", r[4].Model.(map[string]interface{})["name"])
}

func TestIntegrationUpdateDocument(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	authData := SetupAuthentication(t)

	org, err := testServer.OrganizationService.Create("test")
	assert.Nil(t, err)
	err = testServer.AclService.LinkUserToRole(authData.user, "organization:owner", org.Id)
	assert.Nil(t, err)
	folder, err := testServer.FolderService.Create(authData.user, "test folder", org.Id, nil)
	assert.Nil(t, err)
	document, err := testServer.DocumentService.Create(authData.user, org.Id, &folder.Id, "test document", "test content")
	assert.Nil(t, err)

	assert.Equal(t, document.Name, "test document")
	assert.Equal(t, document.Content.Content, "test content")

	status, resp, err := Request(&RequestOptions{
		Method: "PUT",
		Path:   "/document",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.accessToken),
		},
		Body: &model.DocumentUpdateRequest{
			DocumentId: document.Id,
			Name:       "new name",
			Content:    "new content",
		},
		ResponseModel: &model.AclWrappedModel{},
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := resp.(*model.AclWrappedModel)

	assert.Equal(t, document.Id, r.Model.(map[string]interface{})["id"])
	assert.Equal(t, "new name", r.Model.(map[string]interface{})["name"])
	assert.Equal(t, "new content", r.Model.(map[string]interface{})["content"].(map[string]interface{})["content"])
}
