package acl

import (
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	"testing"
)


func TestUserCanNotAccessWhenDoesNotExist(t *testing.T) {
	service := NewAclService(server.NewTransactionManager(database, nil), database, nil)

	user := &model.User{}
	user.Id = "5"
	ok, err := service.UserCanAccessResource(user, []string{"organization", "folder"}, []string{"1", "2"}, "view")
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, ok, false)
}