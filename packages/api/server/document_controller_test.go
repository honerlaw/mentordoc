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
	doc := test.ConvertModel(r[0].Model, &shared.Document{}).(*shared.Document)

	assert.Len(t, r, 1)
	assert.Equal(t, "test document", doc.Drafts[0].Name)
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
	doc := test.ConvertModel(r[0].Model, &shared.Document{}).(*shared.Document)

	assert.Len(t, r, 1)
	assert.Equal(t, "test document", doc.Drafts[0].Name)
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
	doc := test.ConvertModel(r[0].Model, &shared.Document{}).(*shared.Document)
	doc1 := test.ConvertModel(r[1].Model, &shared.Document{}).(*shared.Document)
	doc2 := test.ConvertModel(r[2].Model, &shared.Document{}).(*shared.Document)
	doc3 := test.ConvertModel(r[3].Model, &shared.Document{}).(*shared.Document)
	doc4 := test.ConvertModel(r[4].Model, &shared.Document{}).(*shared.Document)

	assert.Len(t, r, 5)
	assert.Equal(t, "test document 0", doc.Drafts[0].Name)
	assert.Equal(t, "test document 1", doc1.Drafts[0].Name)
	assert.Equal(t, "test document 2", doc2.Drafts[0].Name)
	assert.Equal(t, "test document 3", doc3.Drafts[0].Name)
	assert.Equal(t, "test document 4", doc4.Drafts[0].Name)

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
	doc = test.ConvertModel(r[0].Model, &shared.Document{}).(*shared.Document)
	doc1 = test.ConvertModel(r[1].Model, &shared.Document{}).(*shared.Document)
	doc2 = test.ConvertModel(r[2].Model, &shared.Document{}).(*shared.Document)
	doc3 = test.ConvertModel(r[3].Model, &shared.Document{}).(*shared.Document)
	doc4 = test.ConvertModel(r[4].Model, &shared.Document{}).(*shared.Document)

	assert.Len(t, r, 5)
	assert.Equal(t, "test document 5", doc.Drafts[0].Name)
	assert.Equal(t, "test document 6", doc1.Drafts[0].Name)
	assert.Equal(t, "test document 7", doc2.Drafts[0].Name)
	assert.Equal(t, "test document 8", doc3.Drafts[0].Name)
	assert.Equal(t, "test document 9", doc4.Drafts[0].Name)
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

	assert.Len(t, document.Drafts, 1)
	assert.Equal(t, document.Drafts[0].Name, "test document")
	assert.Equal(t, document.Drafts[0].Content.Content, "test content")

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
	doc := test.ConvertModel(r.Model, &shared.Document{}).(*shared.Document)

	assert.Equal(t, document.Id, doc.Id)
	assert.Equal(t, "new name", doc.Drafts[0].Name)
	assert.Equal(t, "new content", doc.Drafts[0].Content.Content)
}
