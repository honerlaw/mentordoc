package acl_test

import (
	"github.com/honerlaw/mentordoc/server/test"
	"testing"
)

var testData *test.GlobalTestData

func TestMain(m *testing.M) {
	testData = test.InitTestData("../../../.env.test", "../../../migrations")

	test.RunTests(m, testData)
}