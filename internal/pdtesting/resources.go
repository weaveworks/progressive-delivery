package pdtesting

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewNamespace(ctx context.Context, t *testing.T, k client.Client) corev1.Namespace {
	ns := corev1.Namespace{}
	ns.Name = "kube-test-" + rand.String(5)

	err := k.Create(ctx, &ns)
	assert.NoError(t, err, "should be able to create namespace: %s", ns.GetName())

	return ns
}

type CanaryInfo struct {
	Name      string
	Namespace string
	Metrics   []v1beta1.CanaryMetric
}

func NewCanary(
	ctx context.Context,
	t *testing.T,
	k client.Client,
	info CanaryInfo,
) v1beta1.Canary {
	resource := v1beta1.Canary{
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.Name,
			Namespace: info.Namespace,
		},
		Spec: v1beta1.CanarySpec{
			Provider: "traefik",
			TargetRef: v1beta1.LocalObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       info.Name,
			},
			IngressRef: &v1beta1.LocalObjectReference{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
				Name:       info.Name,
			},
			SkipAnalysis: false,
			AutoscalerRef: &v1beta1.LocalObjectReference{
				APIVersion: "autoscaling/v2beta1",
				Kind:       "HorizontalPodAutoscaler",
				Name:       info.Name,
			},
			Service: v1beta1.CanaryService{
				Port:       80,
				TargetPort: intstr.FromInt(9999),
			},
			Analysis: &v1beta1.CanaryAnalysis{
				Iterations: 1,
				Interval:   "1m",
				Metrics:    info.Metrics,
			},
		},
		Status: v1beta1.CanaryStatus{
			Phase:              v1beta1.CanaryPhaseSucceeded,
			FailedChecks:       0,
			CanaryWeight:       0,
			Iterations:         0,
			LastAppliedSpec:    "5978589476",
			LastPromotedSpec:   "5978589476",
			LastTransitionTime: metav1.NewTime(time.Now()),
			Conditions: []v1beta1.CanaryCondition{
				{
					LastUpdateTime:     metav1.NewTime(time.Now()),
					LastTransitionTime: metav1.NewTime(time.Now()),
					Message:            "Canary analysis completed successfully, promotion finished.",
					Reason:             "Succeeded",
					Status:             "True",
					Type:               v1beta1.PromotedType,
				},
			},
		},
	}

	err := k.Create(ctx, &resource)
	assert.NoError(t, err, "should be able to create canary: %s", resource.GetName())

	return resource
}

func NewDeployment(ctx context.Context, t *testing.T, k client.Client, name string, ns string) *appsv1.Deployment {
	dpl := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				server.LabelKustomizeName:      name,
				server.LabelKustomizeNamespace: ns,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "nginx",
						Image: "nginx",
					}},
				},
			},
		},
	}

	err := k.Create(ctx, dpl)
	assert.NoError(t, err, "should be able to create Deployment: %s", dpl.GetName())

	return dpl
}

type CRDInfo struct {
	Group    string
	Plural   string
	Singular string
	Kind     string
	NoTest   bool
}

func NewCRD(
	ctx context.Context,
	t *testing.T,
	k client.Client,
	info CRDInfo,
) v1.CustomResourceDefinition {
	resource := v1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", info.Plural, info.Group),
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: info.Group,
			Names: v1.CustomResourceDefinitionNames{
				Plural:   info.Plural,
				Singular: info.Singular,
				Kind:     info.Kind,
			},
			Scope: "Namespaced",
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:    "v1beta1",
					Served:  true,
					Storage: true,
					Schema: &v1.CustomResourceValidation{
						OpenAPIV3Schema: &v1.JSONSchemaProps{
							Type:       "object",
							Properties: map[string]v1.JSONSchemaProps{},
						},
					},
				},
			},
		},
	}

	err := k.Create(ctx, &resource)

	if !info.NoTest {
		assert.NoError(t, err, "should be able to create crd: %s", resource.GetName())
	}

	return resource
}

type MetricTemplateInfo struct {
	Name                       string
	Namespace                  string
	ProviderType               string
	ProviderAddress            string
	ProviderSecretName         string
	ProviderInsecureSkipVerify bool
	Query                      string
}

func NewMetricTemplate(ctx context.Context, t *testing.T, k client.Client, info MetricTemplateInfo) *v1beta1.MetricTemplate {
	var secretRef *corev1.LocalObjectReference = nil

	if info.ProviderSecretName != "" {
		secretRef = &corev1.LocalObjectReference{Name: info.ProviderSecretName}
	}

	tpl := &v1beta1.MetricTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.Name,
			Namespace: info.Namespace,
			Labels: map[string]string{
				server.LabelKustomizeName:      info.Name,
				server.LabelKustomizeNamespace: info.Namespace,
			},
		},
		Spec: v1beta1.MetricTemplateSpec{
			Provider: v1beta1.MetricTemplateProvider{
				Type:               info.ProviderType,
				Address:            info.ProviderAddress,
				SecretRef:          secretRef,
				InsecureSkipVerify: info.ProviderInsecureSkipVerify,
			},
			Query: info.Query,
		},
	}

	err := k.Create(ctx, tpl)
	assert.NoError(t, err, "should be able to create metric template: %s", tpl.GetName())

	return tpl
}
