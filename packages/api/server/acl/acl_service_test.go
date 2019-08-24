package acl_test

import (
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationUserCanNotAccessWhenDoesNotExist(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(itTestDatabaseConnection, nil), itTestDatabaseConnection, nil)

	user := &model.User{}
	user.Id = "5"
	ok := service.UserCanAccessResource(user, []string{"organization", "folder"}, []string{"1", "2"}, "view")
	assert.Equal(t, ok, false)
}

func TestIntegrationUserLinkToRole(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(itTestDatabaseConnection, nil), itTestDatabaseConnection, nil)

	user := &model.User{}
	user.Id = uuid.NewV4().String()
	itTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)

	err := service.LinkUserToRole(user, "organization:owner", uuid.NewV4().String())

	assert.Nil(t, err)
}

func TestIntegrationUserAccessToDocumentInOrganization(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(itTestDatabaseConnection, nil), itTestDatabaseConnection, nil)

	orgId := uuid.NewV4().String()

	user := &model.User{}
	user.Id = uuid.NewV4().String()
	itTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)
	err := service.LinkUserToRole(user, "organization:owner", orgId)
	assert.Nil(t, err)

	ok := service.UserCanAccessResource(user, []string{"organization", "folder", "document"}, []string{orgId, "10", "25"}, "view")
	assert.Equal(t, ok, true)
}

func TestIntegrationUserActionableResourcesByPath(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(itTestDatabaseConnection, nil), itTestDatabaseConnection, nil)

	orgId := uuid.NewV4().String()
	user := &model.User{}
	user.Id = uuid.NewV4().String()
	itTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)

	err := service.LinkUserToRole(user, "organization:owner", orgId)
	assert.Nil(t, err)

	results, err := service.UserActionableResourcesByPath(user, []string{"organization", "folder", "document"}, "view")
	assert.Nil(t, err)
	assert.Len(t, results, 1)
}

func TestIntegrationWrap(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test")
	}
	service := acl.NewAclService(util.NewTransactionManager(itTestDatabaseConnection, nil), itTestDatabaseConnection, nil)

	user := &model.User{}
	user.Id = uuid.NewV4().String()
	itTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)

	document := &model.Document{
		OrganizationId: "12345",
		FolderId: nil,
	}
	document.Id = "54321"
	documents := []model.Document{*document}
	data, err := service.Wrap(user, documents);

	assert.Nil(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, data[0].Model.(model.Document).Id, document.Id)
	assert.Len(t, data[0].Actions, 0)
}