package server_test

import (
	"context"
	"testing"

	"github.com/onsi/gomega"
	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	api "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestListCanaries(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: server.CreateScheme(),
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	response, err := c.ListCanaries(ctx, &api.ListCanariesRequest{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(response.GetCanaries()).To(gomega.BeEmpty(), "should not return any canary objects")
	g.Expect(response.GetErrors()).ToNot(gomega.BeEmpty(), "should return with errors")

	newCRD(ctx, k, g, crdInfo{
		Group:    "flagger.app",
		Plural:   "canaries",
		Singular: "canary",
		Kind:     "Canary",
		ListKind: "CanaryList",
		Scope:    "Namespaced",
	})

	response, err = c.ListCanaries(ctx, &api.ListCanariesRequest{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(response.GetCanaries()).To(gomega.BeEmpty(), "should not return any canary objects")
	g.Expect(response.GetErrors()).To(gomega.BeEmpty(), "should not return with errors")

	ns := newNamespace(ctx, k, g)

	newCanary(ctx, k, g, "example", ns.Name)

	response, err = c.ListCanaries(ctx, &api.ListCanariesRequest{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(response.GetCanaries()).To(gomega.HaveLen(1), "should return one canary object")
	g.Expect(response.GetErrors()).To(gomega.BeEmpty(), "should not return with errors")
}
