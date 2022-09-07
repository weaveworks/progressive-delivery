package flagger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFetcher_GetMetricTemplate(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	// create namespace before creating service to
	// prevent ns access issues
	ns := pdtesting.NewNamespace(ctx, t, k)

	cl, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	appName := "example"

	metricTemplate := pdtesting.NewMetricTemplate(ctx, t, k, pdtesting.MetricTemplateInfo{
		Name:            appName,
		Namespace:       ns.GetName(),
		ProviderType:    "prometheus",
		ProviderAddress: "http://prometheus:9090",
		Query:           "custom query",
	})
	defer pdtesting.Cleanup(ctx, t, k, metricTemplate)

	template, err := service.GetMetricTemplate(ctx, "Default", cl, appName, ns.GetName())
	assert.NoError(t, err)
	assert.Equal(t, appName, template.GetName())
}

func TestFetcher_ListMetricTemplate(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	// create namespace before creating service to
	// prevent ns access issues
	ns := pdtesting.NewNamespace(ctx, t, k)

	cl, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	appName := "example"

	metricTemplate := pdtesting.NewMetricTemplate(ctx, t, k, pdtesting.MetricTemplateInfo{
		Name:            appName,
		Namespace:       ns.GetName(),
		ProviderType:    "prometheus",
		ProviderAddress: "http://prometheus:9090",
		Query:           "custom query",
	})
	defer pdtesting.Cleanup(ctx, t, k, metricTemplate)

	templates, _, cerrs, err := service.ListMetricTemplates(ctx, cl, flagger.ListMetricTemplatesOptions{})
	assert.NoError(t, err)
	assert.Empty(t, cerrs)
	assert.NotEmpty(t, templates["Default"])
}

func TestFetcher_ListCanaryDeployments(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	// create namespace before creating service to
	// prevent ns access issues
	ns := pdtesting.NewNamespace(ctx, t, k)

	cl, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	appName := "canary"

	canary := pdtesting.NewCanary(ctx, t, k, pdtesting.CanaryInfo{
		Name:      appName,
		Namespace: ns.Name,
	})
	defer pdtesting.Cleanup(ctx, t, k, &canary)

	canaries, _, cerrs, err := service.ListCanaryDeployments(ctx, cl, flagger.ListCanaryDeploymentsOptions{})
	assert.NoError(t, err)
	assert.Empty(t, cerrs)
	assert.NotEmpty(t, canaries["Default"])
}
