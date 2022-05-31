package server_test

import (
	"context"
	"fmt"
	"time"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type crdInfo struct {
	Group    string
	Plural   string
	Singular string
	Kind     string
	ListKind string
	Scope    string
}

func newCRD(
	ctx context.Context,
	k client.Client,
	g *gomega.GomegaWithT,
	info crdInfo,
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
				ListKind: info.ListKind,
			},
			Scope: v1.ResourceScope(info.Scope),
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

	g.Expect(k.Create(ctx, &resource)).To(gomega.Succeed())

	return resource
}

func newCanary(
	ctx context.Context,
	k client.Client,
	g *gomega.GomegaWithT,
	name, namespace string,
) v1beta1.Canary {

	resource := v1beta1.Canary{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.CanarySpec{
			Provider:  "traefik",
			TargetRef: v1beta1.LocalObjectReference{},
			AutoscalerRef: &v1beta1.LocalObjectReference{
				APIVersion: "autoscaling/v2",
				Kind:       "HorizontalPodAutoscaler",
				Name:       name,
			},
			Service: v1beta1.CanaryService{
				Port:       80,
				TargetPort: intstr.FromInt(9999),
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

	g.Expect(k.Create(ctx, &resource)).To(gomega.Succeed())

	return resource
}

func newNamespace(ctx context.Context, k client.Client, g *gomega.GomegaWithT) corev1.Namespace {
	ns := corev1.Namespace{}
	ns.Name = "kube-test-" + rand.String(5)

	g.Expect(k.Create(ctx, &ns)).To(gomega.Succeed())

	return ns
}
