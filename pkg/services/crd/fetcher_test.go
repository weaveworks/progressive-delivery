package crd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/onsi/gomega"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFetcher_IsAvailable(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	service, err := newService(ctx, k8sEnv)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: server.CreateScheme(),
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	var found bool

	found = service.IsAvailable("Default", "customobjects.example.com")
	g.Expect(found).
		To(gomega.BeFalse(), "customobjects crd should not be defined in Default cluster")

	newCRD(ctx, k, g,
		crdInfo{
			Singular: "customobject",
			Group:    "example.com",
			Plural:   "customobjects",
			Kind:     "CustomObject",
		})

	service.UpdateCRDList()

	found = service.IsAvailable("Default", "customobjects.example.com")
	g.Expect(found).
		To(gomega.BeTrue(), "customobjects crd should be defined in Default cluster")

	found = service.IsAvailable("Default", "somethingelse.example.com")
	g.Expect(found).
		To(gomega.BeFalse(), "somethingelse crd should not be defined in Default Cluster")

	found = service.IsAvailable("Other", "customobjects.example.com")
	g.Expect(found).
		To(gomega.BeFalse(), "customobjects crd should not be defined in Other cluster")
}

type crdInfo struct {
	Group    string
	Plural   string
	Singular string
	Kind     string
}

func newCRD(
	ctx context.Context,
	k client.Client,
	g *gomega.GomegaWithT,
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

	g.Expect(k.Create(ctx, &resource)).To(gomega.Succeed())

	return resource
}
