package acl

import (
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestUserCanNotAccessWhenDoesNotExist(t *testing.T) {
	service := NewAclService(server.NewTransactionManager(database, nil), database, nil)

	user := &model.User{}
	user.Id = "5"
	ok, err := service.UserCanAccessResource(user, []string{"organization", "folder"}, []string{"1", "2"}, "view")
	assert.Nil(t, err)
	assert.Equal(t, ok, false)
}

func TestUserLinkToRole(t *testing.T) {
	service := NewAclService(server.NewTransactionManager(database, nil), database, nil)

	user := &model.User{}
	user.Id = uuid.NewV4().String()
	database.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)

	err := service.LinkUserToRole(user, "organization:owner", uuid.NewV4().String())

	assert.Nil(t, err)
}

func TestUserAccessToDocumentInOrganization(t *testing.T) {
	service := NewAclService(server.NewTransactionManager(database, nil), database, nil)

	orgId := uuid.NewV4().String()

	user := &model.User{}
	user.Id = uuid.NewV4().String()
	database.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)
	err := service.LinkUserToRole(user, "organization:owner", orgId)
	assert.Nil(t, err)

	ok, err := service.UserCanAccessResource(user, []string{"organization", "folder", "document"}, []string{orgId, "10", "25"}, "view")
	assert.Nil(t, err)
	assert.Equal(t, ok, true)
}

func TestUserActionableResourcesByPath(t *testing.T) {
	service := NewAclService(server.NewTransactionManager(database, nil), database, nil)

	orgId := uuid.NewV4().String()
	user := &model.User{}
	user.Id = uuid.NewV4().String()
	database.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Id)

	err := service.LinkUserToRole(user, "organization:owner", orgId)
	assert.Nil(t, err)

	results, err := service.UserActionableResourcesByPath(user, []string{"organization", "folder", "document"}, "view")
	assert.Nil(t, err)
	assert.Len(t, results, 1)
}