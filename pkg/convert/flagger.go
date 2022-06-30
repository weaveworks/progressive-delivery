package convert

import (
	"time"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/go-asset/generics/list"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func FlaggerCanaryToProto(canary v1beta1.Canary, clusterName string, deployment appsv1.Deployment, promoted []v1.Container) *pb.Canary {
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

	canaryYaml, _ := yaml.Marshal(canary)
	analysisYaml, _ := yaml.Marshal(canary.Spec.Analysis)

	images := map[string]string{}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		images[container.Name] = container.Image
	}

	promotedImages := map[string]string{}

	for _, c := range promoted {
		promotedImages[c.Name] = c.Image
	}

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
			Uid:                   string(deployment.GetObjectMeta().GetUID()),
			ResourceVersion:       deployment.GetObjectMeta().GetResourceVersion(),
			FluxLabels:            fluxLabels,
			AppliedImageVersions:  images,
			PromotedImageVersions: promotedImages,
		},
		Analysis: &pb.CanaryAnalysis{
			Interval:            canary.Spec.Analysis.Interval,
			Iterations:          int32(canary.Spec.Analysis.Iterations),
			MirrorWeight:        int32(canary.Spec.Analysis.MirrorWeight),
			MaxWeight:           int32(canary.Spec.Analysis.MaxWeight),
			StepWeight:          int32(canary.Spec.Analysis.StepWeight),
			StepWeightPromotion: int32(canary.Spec.Analysis.StepWeightPromotion),
			Threshold:           int32(canary.Spec.Analysis.Threshold),
			Mirror:              canary.Spec.Analysis.Mirror,
			Yaml:                string(analysisYaml),
			StepWeights: list.Map(
				canary.Spec.Analysis.StepWeights,
				func(v int) int32 { return int32(v) },
			),
		},
		Status: &pb.CanaryStatus{
			Phase:              string(canary.Status.Phase),
			FailedChecks:       int32(canary.Status.FailedChecks),
			CanaryWeight:       int32(canary.Status.CanaryWeight),
			Iterations:         int32(canary.Status.Iterations),
			LastTransitionTime: canary.Status.LastTransitionTime.Format(time.RFC3339),
			Conditions:         conditions,
		},
		Yaml: string(canaryYaml),
	}
}

func FlaggerMetricTemplateToProto(template v1beta1.MetricTemplate, clusterName string) *pb.CanaryMetricTemplate {
	secretName := ""

	if template.Spec.Provider.SecretRef != nil {
		secretName = template.Spec.Provider.SecretRef.Name
	}

	return &pb.CanaryMetricTemplate{
		ClusterName: clusterName,
		Name:        template.GetName(),
		Namespace:   template.GetNamespace(),
		Query:       template.Spec.Query,
		Provider: &pb.MetricProvider{
			Type:               template.Spec.Provider.Type,
			Address:            template.Spec.Provider.Address,
			SecretName:         secretName,
			InsecureSkipVerify: template.Spec.Provider.InsecureSkipVerify,
		},
	}
}
