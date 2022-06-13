package convert

import (
	"encoding/json"
	"time"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	v1 "k8s.io/api/apps/v1"
)

func FlaggerCanaryToProto(canary v1beta1.Canary, clusterName string, deployment, primary v1.Deployment) *pb.Canary {
	conditions := []*pb.CanaryCondition{}

	for _, condition := range canary.Status.Conditions {
		conditions = append(conditions, &pb.CanaryCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastUpdateTime:     condition.LastUpdateTime.Format(time.RFC3339),
			LastTransitionTime: condition.LastTransitionTime.Format(time.RFC3339),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}

	fluxLabels := &pb.FluxLabels{
		KustomizeNamespace: deployment.Labels["kustomize.toolkit.fluxcd.io/namespace"],
		KustomizeName:      deployment.Labels["kustomize.toolkit.fluxcd.io/name"],
	}

	targetSpec, _ := json.Marshal(deployment.Spec)
	primarySpec, _ := json.Marshal(primary.Spec)

	return &pb.Canary{
		Name:        canary.GetName(),
		Namespace:   canary.GetNamespace(),
		ClusterName: clusterName,
		Provider:    canary.Spec.Provider,
		TargetReference: &pb.CanaryTargetReference{
			Kind: canary.Spec.TargetRef.Kind,
			Name: canary.Spec.TargetRef.Name,
		},
		TargetDeployment: &pb.CanaryTargetDeployment{
			Uid:             string(deployment.GetObjectMeta().GetUID()),
			ResourceVersion: deployment.GetObjectMeta().GetResourceVersion(),
			FluxLabels:      fluxLabels,
		},
		Status: &pb.CanaryStatus{
			Phase:              string(canary.Status.Phase),
			FailedChecks:       int32(canary.Status.FailedChecks),
			CanaryWeight:       int32(canary.Status.CanaryWeight),
			Iterations:         int32(canary.Status.Iterations),
			LastTransitionTime: canary.Status.LastTransitionTime.Format(time.RFC3339),
			Conditions:         conditions,
		},
		TargetDeploymentSpec:        string(targetSpec),
		TargetPrimaryDeploymentSpec: string(primarySpec),
	}
}
