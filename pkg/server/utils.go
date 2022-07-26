package server

import (
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/kstatus/status"
)

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
