package server_test

import (
	"fmt"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIntegrationListOrganizations(t *testing.T) {
	if !*testData.Integration {
		t.Skip("skipping integration test")
	}
	authData := test.SetupAuthentication(t, testData)

	aclWrappedModels := make([]acl.AclWrappedModel, 0)
	status, resp, err := test.Request(&test.RequestOptions{
		Method: "GET",
		Path:   "/organization/list",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authData.AccessToken),
		},
		ResponseModel: &aclWrappedModels,
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	r := *resp.(*[]acl.AclWrappedModel)

	assert.Len(t, r, 1)
	assert.Equal(t, r[0].Model.(map[string]interface{})["id"], authData.Organization.Id)
}
