package flagger

import "github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"

type DeploymentStrategy string

const (
	CanaryDeploymentStrategy          DeploymentStrategy = "canary"
	BlueGreenDeploymentStrategy       DeploymentStrategy = "blue-green"
	BlueGreenMirrorDeploymentStrategy DeploymentStrategy = "blue-green-mirror"
	ABTestingDeploymentStrategy       DeploymentStrategy = "ab-testing"
	NoAnalysisDeploymentStrategy      DeploymentStrategy = "no-analysis"
)

func (service *defaultFetcher) DeploymentStrategyFor(canary v1beta1.Canary) DeploymentStrategy {
	if canary.Spec.SkipAnalysis {
		return NoAnalysisDeploymentStrategy
	}

	hasIterations := canary.Spec.Analysis.Iterations > 0
	hasMatch := len(canary.Spec.Analysis.Match) > 0

	if canary.Spec.Analysis.Mirror && hasIterations {
		return BlueGreenMirrorDeploymentStrategy
	}

	if hasIterations && !hasMatch {
		return BlueGreenDeploymentStrategy
	}

	if hasIterations {
		return ABTestingDeploymentStrategy
	}

	return CanaryDeploymentStrategy
}
