package acl_test

import (
	"github.com/honerlaw/mentordoc/server/lib/acl"
	document2 "github.com/honerlaw/mentordoc/server/lib/document"
	user2 "github.com/honerlaw/mentordoc/server/lib/user"
	"github.com/honerlaw/mentordoc/server/lib/util"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationUserCanNotAccessWhenDoesNotExist(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(testData.ItTestDatabaseConnection, nil), testData.ItTestDatabaseConnection, nil)

	user := &user2.User{}
	user.Id = "5"
	ok := service.UserCanAccessResource(user, []string{"organization", "folder"}, []string{"1", "2"}, "view")
	assert.Equal(t, ok, false)
}

func TestIntegrationUserLinkToRole(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(testData.ItTestDatabaseConnection, nil), testData.ItTestDatabaseConnection, nil)

	user := &user2.User{}
	user.Id = uuid.NewV4().String()
	_, err := testData.ItTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)
	assert.Nil(t, err)

	err = service.LinkUserToRole(user, "organization:owner", uuid.NewV4().String())

	assert.Nil(t, err)
}

func TestIntegrationUserAccessToDocumentInOrganization(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(testData.ItTestDatabaseConnection, nil), testData.ItTestDatabaseConnection, nil)

	orgId := uuid.NewV4().String()

	user := &user2.User{}
	user.Id = uuid.NewV4().String()
	_, err := testData.ItTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)
	assert.Nil(t, err)
	err = service.LinkUserToRole(user, "organization:owner", orgId)
	assert.Nil(t, err)

	ok := service.UserCanAccessResource(user, []string{"organization", "folder", "document"}, []string{orgId, "10", "25"}, "view")
	assert.Equal(t, ok, true)
}

func TestIntegrationUserActionableResourcesByPath(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(testData.ItTestDatabaseConnection, nil), testData.ItTestDatabaseConnection, nil)

	orgId := uuid.NewV4().String()
	user := &user2.User{}
	user.Id = uuid.NewV4().String()
	_, err := testData.ItTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)
	assert.Nil(t, err)

	err = service.LinkUserToRole(user, "organization:owner", orgId)
	assert.Nil(t, err)

	results, err := service.UserActionableResourcesByPath(user, []string{"organization", "folder", "document"}, "view")
	assert.Nil(t, err)
	assert.Len(t, results, 1)
}

func TestIntegrationWrap(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(testData.ItTestDatabaseConnection, nil), testData.ItTestDatabaseConnection, nil)

	user := &user2.User{}
	user.Id = uuid.NewV4().String()
	_, err := testData.ItTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)
	assert.Nil(t, err)

	document := &document2.Document{
		OrganizationId: "12345",
		FolderId: nil,
	}
	document.Id = "54321"
	documents := []document2.Document{*document}
	data, err := service.Wrap(user, documents);

	assert.Nil(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, data[0].Model.(document2.Document).Id, document.Id)
	assert.Len(t, data[0].Actions, 0)
}