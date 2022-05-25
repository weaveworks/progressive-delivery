package server_test

import (
	"os"
	"testing"

	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	"github.com/weaveworks/weave-gitops/pkg/testutils"
)

var k8sEnv *testutils.K8sTestEnv

func TestMain(m *testing.M) {
	var err error

	k8sEnv, err = pdtesting.CreateTestEnv()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	k8sEnv.Stop()

	os.Exit(code)
}
