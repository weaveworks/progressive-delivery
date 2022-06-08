package crd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type crdInfo struct {
	Group    string
	Plural   string
	Singular string
	Kind     string
}

func newCRD(
	ctx context.Context,
	t *testing.T,
	k client.Client,
	info crdInfo,
) v1.CustomResourceDefinition {

	resource := v1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", info.Plural, info.Group),
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: info.Group,
			Names: v1.CustomResourceDefinitionNames{
				Plural:   info.Plural,
				Singular: info.Singular,
				Kind:     info.Kind,
			},
			Scope: "Namespaced",
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:    "v1beta1",
					Served:  true,
					Storage: true,
					Schema: &v1.CustomResourceValidation{
						OpenAPIV3Schema: &v1.JSONSchemaProps{
							Type:       "object",
							Properties: map[string]v1.JSONSchemaProps{},
						},
					},
				},
			},
		},
	}

	err := k.Create(ctx, &resource)
	assert.NoError(t, err, "should be able to create crd: %s", resource.GetName())

	return resource
}

func TestFetcher_IsAvailable(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: server.CreateScheme(),
	})
	assert.NoError(t, err)

	var found bool

	found = service.IsAvailable("Default", "customobjects.example.com")
	assert.False(t, found, "customobjects crd should not be defined in Default cluster")

	newCRD(ctx, t, k,
		crdInfo{
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
