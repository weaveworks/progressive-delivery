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

	cl, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	appName := "example"

	ns := pdtesting.NewNamespace(ctx, t, k)

	_ = pdtesting.NewMetricTemplate(ctx, t, k, pdtesting.MetricTemplateInfo{
		Name:            appName,
		Namespace:       ns.GetName(),
		ProviderType:    "prometheus",
		ProviderAddress: "http://prometheus:9090",
		Query:           "custom query",
	})

	template, err := service.GetMetricTemplate(ctx, "Default", cl, appName, ns.GetName())
	assert.NoError(t, err)
	assert.Equal(t, appName, template.GetName())
}

func TestFetcher_ListMetricTemplate(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())

	defer cancelFn()

	cl, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	appName := "example"

	ns := pdtesting.NewNamespace(ctx, t, k)

	_ = pdtesting.NewMetricTemplate(ctx, t, k, pdtesting.MetricTemplateInfo{
		Name:            appName,
		Namespace:       ns.GetName(),
		ProviderType:    "prometheus",
		ProviderAddress: "http://prometheus:9090",
		Query:           "custom query",
	})

	templates, _, cerrs, err := service.ListMetricTemplates(ctx, cl, flagger.ListMetricTemplatesOptions{})
	assert.NoError(t, err)
	assert.Empty(t, cerrs)
	assert.NotEmpty(t, templates["Default"])
}
