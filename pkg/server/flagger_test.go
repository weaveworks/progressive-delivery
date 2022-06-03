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

	ns := newNamespace(ctx, k, g)

	newCanary(ctx, k, g, "example", ns.Name)

	response, err := c.ListCanaries(ctx, &api.ListCanariesRequest{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(response.GetCanaries()).To(gomega.HaveLen(1), "should return one canary object")
	g.Expect(response.GetErrors()).To(gomega.BeEmpty(), "should not return with errors")
}

func TestGetCanary(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: server.CreateScheme(),
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	appName := "example"

	ns := newNamespace(ctx, k, g)

	_ = newDeployment(ctx, k, g, appName, ns.Name)
	canary := newCanary(ctx, k, g, appName, ns.Name)

	response, err := c.GetCanary(ctx, &api.GetCanaryRequest{ClusterName: "Default", Name: canary.Name, Namespace: canary.Namespace})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(response.GetCanary().Name).To(gomega.Equal(canary.Name))

	g.Expect(response.GetAutomation()).To(gomega.HaveField("Name", appName))
	g.Expect(response.GetAutomation()).To(gomega.HaveField("Namespace", ns.Name))
}
