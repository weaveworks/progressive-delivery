package convert

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/go-asset/generics/list"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/kube"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func FlaggerCanaryToProto(canary v1beta1.Canary, clusterName string, deployment appsv1.Deployment, promoted []v1.Container, metricTemplates []v1beta1.MetricTemplate) *pb.Canary {
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

	canaryYaml, _ := serializeObj(&canary)
	analysisYaml, _ := yaml.Marshal(canary.Spec.Analysis)

	images := map[string]string{}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		images[container.Name] = container.Image
	}

	promotedImages := map[string]string{}
	for _, c := range promoted {
		promotedImages[c.Name] = c.Image
	}

	//canary metrics
	metrics := []*pb.CanaryMetric{}
	for _, metric := range canary.Spec.Analysis.Metrics {
		var metricTemplate *pb.CanaryMetricTemplate
		if metric.TemplateRef != nil {
			for _, mt := range metricTemplates {
				if mt.Name == metric.TemplateRef.Name &&
					mt.Namespace == metric.TemplateRef.Namespace {
					//secretRef is optional
					var secretRefName string
					if mt.Spec.Provider.SecretRef != nil {
						secretRefName = mt.Spec.Provider.SecretRef.Name
					}
					metricTemplateYaml, err := serializeObj(&mt)
					if err != nil {
						//TODO ask on the strategy to handle errors within the code
						log.Println("could not create yaml for metric template", err)
						metricTemplateYaml = []byte{}
					}
					metricTemplate = &pb.CanaryMetricTemplate{
						Namespace:   mt.Namespace,
						Name:        mt.Name,
						ClusterName: clusterName,
						Provider: &pb.MetricProvider{
							Type:               mt.Spec.Provider.Type,
							Address:            mt.Spec.Provider.Address,
							SecretName:         secretRefName,
							InsecureSkipVerify: mt.Spec.Provider.InsecureSkipVerify,
						},
						Query: mt.Spec.Query,
						Yaml:  string(metricTemplateYaml),
					}
				}
			}
		}

		var thresholdRange *pb.CanaryMetricThresholdRange
		if metric.ThresholdRange != nil {
			var min float64
			var max float64
			if metric.ThresholdRange.Min != nil {
				min = *metric.ThresholdRange.Min
			}
			if metric.ThresholdRange.Max != nil {
				max = *metric.ThresholdRange.Max
			}
			thresholdRange = &pb.CanaryMetricThresholdRange{
				Min: min,
				Max: max,
			}
		}

		metrics = append(metrics, &pb.CanaryMetric{
			Name:           string(metric.Name),
			Interval:       string(metric.Interval),
			ThresholdRange: thresholdRange,
			MetricTemplate: metricTemplate,
		})
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
			Metrics: metrics,
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

func serializeObj(obj client.Object) ([]byte, error) {
	scheme := kube.CreateScheme()

	if err := setGVKFromScheme(obj, scheme); err != nil {
		return nil, err
	}

	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{
		Pretty: true,
		Yaml:   true,
		Strict: true,
	})

	buf := bytes.NewBufferString("")

	if err := serializer.Encode(obj, buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Populate the GVK from scheme, since it is cleared by design on typed objects.
// https://github.com/kubernetes/client-go/issues/413
func setGVKFromScheme(object runtime.Object, scheme *runtime.Scheme) error {
	gvks, unversioned, err := scheme.ObjectKinds(object)
	if err != nil {
		return err
	}
	if len(gvks) == 0 {
		return fmt.Errorf("no ObjectKinds available for %T", object)
	}
	if !unversioned {
		object.GetObjectKind().SetGroupVersionKind(gvks[0])
	}
	return nil
}
