package flagger_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
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

func newService(ctx context.Context, k8sEnv *testutils.K8sTestEnv) (flagger.Fetcher, error) {
	client, err := pdtesting.CreateClient(k8sEnv)
	if err != nil {
		return nil, err
	}

	crdSrv := crd.NewFetcher(ctx, logr.Discard(), client)

	return flagger.NewFetcher(crdSrv), nil
}
