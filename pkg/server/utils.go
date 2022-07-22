package server

import (
	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kustomize/kstatus/status"
)

func getRef(
	ctx context.Context,
	clusterClient clustersmngr.Client,
	ref *v1beta1.LocalObjectReference,
	ns string,
	clusterName string,
) (unstructured.Unstructured, error) {
	object := unstructured.Unstructured{}
	key := client.ObjectKey{
		Name:      ref.Name,
		Namespace: ns,
	}

	object.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    ref.Kind,
		Version: ref.APIVersion,
	})

	err := clusterClient.Get(ctx, clusterName, key, &object)

	return object, err
}

func mapUnstructuredConditions(result *status.Result) []*pb.Condition {
	conds := []*pb.Condition{}

	if result.Status == status.CurrentStatus {
		conds = append(conds, &pb.Condition{Type: "Ready", Status: "True", Message: result.Message})
	}

	return conds
}

func getContainerImages(containers []interface{}) []string {
	images := []string{}

	for _, item := range containers {
		container, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		image, ok, _ := unstructured.NestedString(container, "image")
		if ok {
			images = append(images, image)
		}
	}

	return images
}

func getDeploymentPodContainerImages(obj map[string]interface{}) []string {
	containers, _, _ := unstructured.NestedSlice(
		obj,
		"spec", "template", "spec", "containers",
	)

	return getContainerImages(containers)
}
