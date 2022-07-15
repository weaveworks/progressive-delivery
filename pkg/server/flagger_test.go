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
		},
	}
	canaryMetricWithoutThreshold := v1beta1.CanaryMetric{
		Name:     "request-success-rate",
		Interval: "1m",
	}
	canaryMetricTemplate := pdtesting.NewMetricTemplate(ctx, t, k, pdtesting.MetricTemplateInfo{
		Name:               fmt.Sprintf("%s-mt", appName),
		Namespace:          ns.GetName(),
		ProviderType:       "prometheus",
		ProviderAddress:    "http://prometheus:9090",
		Query:              "custom query",
		ProviderSecretName: "prometheusSecret",
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
	canaryMetricTemplateWithoutSecret := pdtesting.NewMetricTemplate(ctx, t, k, pdtesting.MetricTemplateInfo{
		Name:            fmt.Sprintf("%s-mt-no-secret", appName),
		Namespace:       ns.GetName(),
		ProviderType:    "prometheus",
		ProviderAddress: "http://prometheus:9090",
		Query:           "custom query",
	})
	canaryMetricWithTemplateWithoutSecret := v1beta1.CanaryMetric{
		Name:     "my-custom-metric",
		Interval: "2m",
		ThresholdRange: &v1beta1.CanaryThresholdRange{
			Min: toFloatPtr(50.0),
			Max: toFloatPtr(75.0),
		},
		TemplateRef: &v1beta1.CrossNamespaceObjectReference{
			Name:      canaryMetricTemplateWithoutSecret.Name,
			Namespace: canaryMetricTemplateWithoutSecret.Namespace,
		},
	}

	canary := pdtesting.NewCanary(ctx, t, k, pdtesting.CanaryInfo{
		Name:      appName,
		Namespace: ns.GetName(),
		Metrics: []v1beta1.CanaryMetric{
			canaryMetric,
			canaryMetricWithoutThreshold,
			canaryMetricWithTemplate,
			canaryMetricWithTemplateWithoutSecret,
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
	assert.True(t, len(response.GetCanary().GetAnalysis().Metrics) == 4)
	assertMetric(t, response.GetCanary().GetAnalysis().GetMetrics()[0], canaryMetric, nil)
	assertMetric(t, response.GetCanary().GetAnalysis().GetMetrics()[1], canaryMetricWithoutThreshold, nil)
	assertMetric(t, response.GetCanary().GetAnalysis().GetMetrics()[2], canaryMetricWithTemplate, canaryMetricTemplate)
	assertMetric(t, response.GetCanary().GetAnalysis().GetMetrics()[3], canaryMetricWithTemplateWithoutSecret, canaryMetricTemplateWithoutSecret)
}

func assertMetric(t *testing.T, actual *api.CanaryMetric, expected v1beta1.CanaryMetric, expectedMetricTemplate *v1beta1.MetricTemplate) {
	assert.Equal(t,
		expected.Name,
		actual.GetName(),
	)
	if expected.ThresholdRange != nil {
		if expected.ThresholdRange.Min != nil {
			assert.Equal(t,
				*expected.ThresholdRange.Min,
				actual.ThresholdRange.Min,
			)
		}
		if expected.ThresholdRange.Max != nil {
			assert.Equal(t,
				*expected.ThresholdRange.Max,
				actual.ThresholdRange.Max,
			)
		}
	}
	if expected.TemplateRef != nil {
		assert.Equal(t,
			expected.TemplateRef.Name,
			actual.MetricTemplate.Name,
		)
		assert.Equal(t,
			expected.TemplateRef.Namespace,
			actual.MetricTemplate.Namespace,
		)
		assert.Equal(t,
			expectedMetricTemplate.Spec.Query,
			actual.MetricTemplate.Query,
		)
		assert.Equal(t,
			expectedMetricTemplate.Spec.Provider.Type,
			actual.MetricTemplate.Provider.Type,
		)
		if expectedMetricTemplate.Spec.Provider.SecretRef != nil {
			assert.Equal(t,
				expectedMetricTemplate.Spec.Provider.SecretRef.Name,
				actual.MetricTemplate.Provider.SecretName,
			)
		}
		assert.NotEmpty(t, actual.MetricTemplate.Yaml)
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
