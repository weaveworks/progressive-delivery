package server

import (
	"context"
	"fmt"

	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/services/flagger"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/kustomize/kstatus/status"
)

// list of kinds flagger creates depending on the underlying service mesh used.
var generatedObjectsKinds = []schema.GroupVersionKind{
	// all canaries
	{Group: "", Version: "v1", Kind: "Service"},
	{Group: "apps", Version: "v1", Kind: "Deployment"},
	{Group: "autoscaling", Version: "v1", Kind: "HorizontalPodAutoscaler"},
	// TODO: we need to detect which service is beeing used and query its objects types.
	// Maybe borrow some logic from Flagger.
	// Linkerd
	// {Group: "split.smi-spec.io", Version: "v1alpha2", Kind: "trafficsplit"},
	// Istio
	// {Group: "networking.istio.io", Version: "v1alpha3", Kind: "destinationrules"},
	// {Group: "networking.istio.io", Version: "v1alpha3", Kind: "virtualservices"},
}

func (pd *pdServer) ListCanaryObjects(ctx context.Context, msg *pb.ListCanaryObjectsRequest) (*pb.ListCanaryObjectsResponse, error) {
	clusterClient, err := pd.clientsFactory.GetImpersonatedClient(ctx, auth.Principal(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting impersonating client: %w", err)
	}

	result := []unstructured.Unstructured{}
	checkDup := map[types.UID]bool{}

	// Get canary object
	canary, err := pd.flagger.GetCanary(ctx, clusterClient, flagger.GetCanaryOptions{
		Name:        msg.Name,
		Namespace:   msg.Namespace,
		ClusterName: msg.ClusterName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to find canary object: %w", err)
	}

	targetDeployment, err := getRef(
		ctx,
		clusterClient,
		&canary.Spec.TargetRef,
		canary.GetNamespace(),
		msg.GetClusterName(),
	)
	if err == nil {
		result = append(result, targetDeployment)
	}

	if canary.Spec.IngressRef != nil {
		ingress, err := getRef(
			ctx,
			clusterClient,
			canary.Spec.AutoscalerRef,
			canary.GetNamespace(),
			msg.GetClusterName(),
		)
		if err == nil {
			result = append(result, ingress)
		}
	}

	if canary.Spec.AutoscalerRef != nil {
		hpa, err := getRef(
			ctx,
			clusterClient,
			canary.Spec.AutoscalerRef,
			canary.GetNamespace(),
			msg.GetClusterName(),
		)
		if err == nil {
			result = append(result, hpa)
		}
	}

	for _, gvk := range generatedObjectsKinds {
		listResult := unstructured.UnstructuredList{}

		listResult.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   gvk.Group,
			Kind:    gvk.Kind,
			Version: gvk.Version,
		})

		if err := clusterClient.List(ctx, msg.ClusterName, &listResult); err != nil {
			if k8serrors.IsForbidden(err) {
				pd.logger.Error(err, "request is forbidden")

				continue
			}

			return nil, fmt.Errorf("listing unstructured object: %w", err)
		}

	ItemsLoop:
		for _, obj := range listResult.Items {
			refs := obj.GetOwnerReferences()
			if len(refs) == 0 {
				continue
			}

			for _, ref := range refs {
				if ref.UID != canary.GetUID() {
					continue ItemsLoop
				}
			}

			uid := obj.GetUID()

			if !checkDup[uid] {
				result = append(result, obj)
				checkDup[uid] = true
			}
		}
	}

	objects := []*pb.UnstructuredObject{}

	for _, obj := range result {
		res, err := status.Compute(&obj)
		if err != nil {
			return nil, fmt.Errorf("could not get status for %s: %w", obj.GetName(), err)
		}

		var images []string

		switch obj.GetKind() {
		case "Deployment":
			images = getDeploymentPodContainerImages(obj.Object)
		}

		objects = append(objects, &pb.UnstructuredObject{
			GroupVersionKind: &pb.GroupVersionKind{
				Group:   obj.GetObjectKind().GroupVersionKind().Group,
				Version: obj.GetObjectKind().GroupVersionKind().GroupVersion().Version,
				Kind:    obj.GetKind(),
			},
			Name:        obj.GetName(),
			Namespace:   obj.GetNamespace(),
			Images:      images,
			Status:      res.Status.String(),
			Uid:         string(obj.GetUID()),
			Conditions:  mapUnstructuredConditions(res),
			ClusterName: msg.GetClusterName(),
		})
	}

	return &pb.ListCanaryObjectsResponse{Objects: objects}, nil
}
