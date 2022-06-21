package flagger_test

import (
	"context"
	"testing"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/fluxcd/flagger/pkg/apis/istio/common/v1alpha1"
	"github.com/fluxcd/flagger/pkg/apis/istio/v1alpha3"
	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
)

func TestFlagger_DeploymentStrategyFor_Canary(t *testing.T) {
	ctx := context.Background()

	_, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	canary := v1beta1.Canary{
		Spec: v1beta1.CanarySpec{
			Analysis: &v1beta1.CanaryAnalysis{
				Interval:            "1m",
				Iterations:          0,
				Mirror:              false,
				MirrorWeight:        0,
				MaxWeight:           50,
				StepWeight:          10,
				Threshold:           10,
				StepWeightPromotion: 10,
				Match:               []v1alpha3.HTTPMatchRequest{},
			},
			SkipAnalysis: false,
		},
	}

	assert.Equal(t, flagger.CanaryDeploymentStrategy, service.DeploymentStrategyFor(canary))
}

func TestFlagger_DeploymentStrategyFor_BlueGreen(t *testing.T) {
	ctx := context.Background()

	_, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	canary := v1beta1.Canary{
		Spec: v1beta1.CanarySpec{
			Analysis: &v1beta1.CanaryAnalysis{
				Interval:            "1m",
				Iterations:          10,
				Mirror:              false,
				MirrorWeight:        0,
				MaxWeight:           50,
				StepWeight:          10,
				Threshold:           10,
				StepWeightPromotion: 10,
				Match:               []v1alpha3.HTTPMatchRequest{},
			},
			SkipAnalysis: false,
		},
	}

	assert.Equal(t, flagger.BlueGreenDeploymentStrategy, service.DeploymentStrategyFor(canary))
}

func TestFlagger_DeploymentStrategyFor_BlueGreenMirror(t *testing.T) {
	ctx := context.Background()

	_, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	canary := v1beta1.Canary{
		Spec: v1beta1.CanarySpec{
			Analysis: &v1beta1.CanaryAnalysis{
				Interval:            "1m",
				Iterations:          10,
				Mirror:              true,
				MirrorWeight:        0,
				MaxWeight:           50,
				StepWeight:          10,
				Threshold:           10,
				StepWeightPromotion: 10,
				Match:               []v1alpha3.HTTPMatchRequest{},
			},
			SkipAnalysis: false,
		},
	}

	assert.Equal(t, flagger.BlueGreenMirrorDeploymentStrategy, service.DeploymentStrategyFor(canary))
}

func TestFlagger_DeploymentStrategyFor_ABTesting(t *testing.T) {
	ctx := context.Background()

	_, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	canary := v1beta1.Canary{
		Spec: v1beta1.CanarySpec{
			Analysis: &v1beta1.CanaryAnalysis{
				Interval:            "1m",
				Iterations:          10,
				Mirror:              false,
				MirrorWeight:        0,
				MaxWeight:           50,
				StepWeight:          10,
				Threshold:           10,
				StepWeightPromotion: 10,
				Match: []v1alpha3.HTTPMatchRequest{
					{
						Headers: map[string]v1alpha1.StringMatch{
							"x-canary": {Regex: "^(.*?;)?(canary=always)(;.*)?$"},
						},
					},
				},
			},
			SkipAnalysis: false,
		},
	}

	assert.Equal(t, flagger.ABTestingDeploymentStrategy, service.DeploymentStrategyFor(canary))
}

func TestFlagger_DeploymentStrategyFor_NoAnalysis(t *testing.T) {
	ctx := context.Background()

	_, service, err := newService(ctx, k8sEnv)
	assert.NoError(t, err)

	canary := v1beta1.Canary{
		Spec: v1beta1.CanarySpec{
			SkipAnalysis: true,
		},
	}

	assert.Equal(t, flagger.NoAnalysisDeploymentStrategy, service.DeploymentStrategyFor(canary))
}
