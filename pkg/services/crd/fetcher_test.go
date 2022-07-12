package crd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFetcher_IsAvailable(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	var found bool

	found = service.IsAvailable("Default", "customobjects.example.com")
	assert.False(t, found, "customobjects crd should not be defined in Default cluster")

	pdtesting.NewCRD(ctx, t, k,
		pdtesting.CRDInfo{
			Singular: "customobject",
			Group:    "example.com",
			Plural:   "customobjects",
			Kind:     "CustomObject",
		})

	service.UpdateCRDList()

	found = service.IsAvailable("Default", "customobjects.example.com")
	assert.True(t, found, "customobjects crd should be defined in Default cluster")

	found = service.IsAvailable("Default", "somethingelse.example.com")
	assert.False(t, found, "somethingelse crd should not be defined in Default Cluster")

	found = service.IsAvailable("Other", "customobjects.example.com")
	assert.False(t, found, "customobjects crd should not be defined in Other cluster")
}

func TestFetcher_IsAvailableOnClusters(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	pdtesting.NewCRD(ctx, t, k,
		pdtesting.CRDInfo{
			Singular: "xclustercustomon",
			Group:    "example.com",
			Plural:   "xclustercustomons",
			Kind:     "CrossClusterCustomObject",
		},
	)

	crdName := "xclustercustomons.example.com"

	service.UpdateCRDList()

	response := service.IsAvailableOnClusters(crdName)

	assert.Len(t, response, 1, "cluster list should contain one entry")
	assert.True(t, response["Default"], "%s should be available on Default cluster", crdName)

	crdName = "xclusterothercustomons.example.com"

	response = service.IsAvailableOnClusters(crdName)

	assert.Len(t, response, 1, "cluster list should contain one entry")
	assert.False(t, response["Default"], "%s shouldn't be available on Default cluster", crdName)
}
