package flagger_test

import (
	"context"
	"os"
	"testing"

	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
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

func newService(ctx context.Context, k8sEnv *testutils.K8sTestEnv) (clustersmngr.Client, flagger.Fetcher, error) {
	cl, clientFactory, err := pdtesting.CreateClient(k8sEnv)
	if err != nil {
		return nil, nil, err
	}

	crdSrv := crd.NewNoCacheFetcher(clientFactory)

	return cl, flagger.NewFetcher(crdSrv), nil
}
