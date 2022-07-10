package server_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	api "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestListCanaries(t *testing.T) {
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	ns := pdtesting.NewNamespace(ctx, t, k)

	appName := "my-app"

	_ = pdtesting.NewDeployment(ctx, t, k, appName, ns.Name)

	canary := pdtesting.NewCanary(ctx, t, k, pdtesting.CanaryInfo{
		Name:      appName,
		Namespace: ns.GetName(),
	})
	defer cleanup(ctx, t, k, &canary)

	response, err := c.ListCanaries(ctx, &api.ListCanariesRequest{})
	assert.NoError(t, err)

	assert.Len(t, response.GetCanaries(), 1, "should return one canary object")
	assert.Empty(t, response.GetErrors(), "should not return with errors")
	assert.Equal(t,
		string(flagger.BlueGreenDeploymentStrategy),
		response.GetCanaries()[0].GetDeploymentStrategy(),
	)

	expectedImages := map[string]string{
		"nginx": "nginx",
	}
	assert.Equal(t, expectedImages, response.Canaries[0].TargetDeployment.AppliedImageVersions)

	assert.Empty(t, response.Canaries[0].TargetDeployment.PromotedImageVersions)
}

func TestListCanaries_NoDeployment(t *testing.T) {
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	ns := pdtesting.NewNamespace(ctx, t, k)

	appName := "no-dep"

	canary := pdtesting.NewCanary(ctx, t, k, pdtesting.CanaryInfo{
		Name:      appName,
		Namespace: ns.GetName(),
	})
	defer cleanup(ctx, t, k, &canary)

	response, err := c.ListCanaries(ctx, &api.ListCanariesRequest{})
	assert.NoError(t, err)

	assert.Len(t, response.GetCanaries(), 1, "should return one canary object")
	assert.Empty(t, response.GetErrors(), "should not return with errors")

	assert.Empty(t, response.Canaries[0].TargetDeployment.AppliedImageVersions)
}

func TestListCanaries_Promoted(t *testing.T) {
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	ns := pdtesting.NewNamespace(ctx, t, k)

	appName := "with-promoted"

	_ = pdtesting.NewDeployment(ctx, t, k, appName, ns.Name)
	_ = pdtesting.NewDeployment(ctx, t, k, fmt.Sprintf("%s-primary", appName), ns.Name)

	canary := pdtesting.NewCanary(ctx, t, k, pdtesting.CanaryInfo{
		Name:      appName,
		Namespace: ns.GetName(),
	})
	defer cleanup(ctx, t, k, &canary)

	response, err := c.ListCanaries(ctx, &api.ListCanariesRequest{})
	assert.NoError(t, err)

	assert.Len(t, response.GetCanaries(), 1, "should return one canary object")
	assert.Empty(t, response.GetErrors(), "should not return with errors")

	expectedImages := map[string]string{
		"nginx": "nginx",
	}
	assert.Equal(t, expectedImages, response.Canaries[0].TargetDeployment.AppliedImageVersions)

	expectedPromoted := map[string]string{
		"nginx": "nginx",
	}
	assert.Equal(t, expectedPromoted, response.Canaries[0].TargetDeployment.PromotedImageVersions)
}

func TestGetCanary(t *testing.T) {
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	k, err := client.New(k8sEnv.Rest, client.Options{
		Scheme: kube.CreateScheme(),
	})
	assert.NoError(t, err)

	appName := "example"

	ns := pdtesting.NewNamespace(ctx, t, k)
	_ = pdtesting.NewDeployment(ctx, t, k, appName, ns.Name)
	_ = pdtesting.NewDeployment(ctx, t, k, fmt.Sprintf("%s-primary", appName), ns.Name)

	canaryMetric := v1beta1.CanaryMetric{
		Name:     "request-success-rate",
		Interval: "1m",
		ThresholdRange: &v1beta1.CanaryThresholdRange{
			Min: toFloatPtr(90),
			Max: toFloatPtr(99),
		},
	}

	canaryMetricTemplate := pdtesting.NewMetricTemplate(ctx, t, k, pdtesting.MetricTemplateInfo{
		Name:            appName,
		Namespace:       ns.GetName(),
		ProviderType:    "prometheus",
		ProviderAddress: "http://prometheus:9090",
		Query:           "custom query",
	})
	canaryMetricWithTemplate := v1beta1.CanaryMetric{
		Name:     "my-custom-metric",
		Interval: "2m",
		ThresholdRange: &v1beta1.CanaryThresholdRange{
			Min: toFloatPtr(50.0),
			Max: toFloatPtr(75.0),
		},
		TemplateRef: &v1beta1.CrossNamespaceObjectReference{
			Name:      canaryMetricTemplate.Name,
			Namespace: canaryMetricTemplate.Namespace,
		},
	}
	canary := pdtesting.NewCanary(ctx, t, k, pdtesting.CanaryInfo{
		Name:      appName,
		Namespace: ns.GetName(),
		Metrics: []v1beta1.CanaryMetric{
			canaryMetric,
			canaryMetricWithTemplate,
		},
	})
	defer cleanup(ctx, t, k, &canary)

	response, err := c.GetCanary(ctx, &api.GetCanaryRequest{ClusterName: "Default", Name: canary.Name, Namespace: canary.Namespace})
	assert.NoError(t, err)

	assert.Equal(t, canary.Name, response.GetCanary().GetName())
	assert.NotNil(t, response.GetAutomation())
	assert.Equal(t, appName, response.GetAutomation().GetName())
	assert.Equal(t, ns.Name, response.GetAutomation().GetNamespace())
	assert.Equal(t,
		string(flagger.BlueGreenDeploymentStrategy),
		response.GetCanary().GetDeploymentStrategy(),
	)
	//TODO: add metrics
	assert.True(t, len(response.GetCanary().GetAnalysis().Metrics) == 2)
	assertMetric(t, response.GetCanary().GetAnalysis().GetMetrics()[0], canaryMetric)
	assertMetric(t, response.GetCanary().GetAnalysis().GetMetrics()[1], canaryMetricWithTemplate)
}

func assertMetric(t *testing.T, actual *api.CanaryMetric, expected v1beta1.CanaryMetric) {
	assert.Equal(t,
		actual.GetName(),
		expected.Name,
	)
	assert.Equal(t,
		fmt.Sprintf("%f", *expected.ThresholdRange.Min),
		actual.ThresholdRange.Min,
	)
	if expected.TemplateRef != nil {
		assert.Equal(t,
			actual.MetricTemplate.Name,
			expected.TemplateRef.Name,
		)
		assert.Equal(t,
			actual.MetricTemplate.Namespace,
			expected.TemplateRef.Namespace,
		)
	}
}

func TestIsFlaggerAvailable(t *testing.T) {
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	response, err := c.IsFlaggerAvailable(ctx, &api.IsFlaggerAvailableRequest{})
	assert.NoError(t, err)

	assert.Len(t, response.GetClusters(), 1)
}

func cleanup(ctx context.Context, t *testing.T, k client.Client, obj client.Object) {
	if err := k.Delete(ctx, obj); err != nil {
		t.Error(err)
	}
}

func toFloatPtr(val int) *float64 {
	v := float64(val)
	return &v
}
