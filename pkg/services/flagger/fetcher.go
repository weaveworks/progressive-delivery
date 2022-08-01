package flagger

import (
	"context"
	"errors"
	"fmt"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	flaggerv1 "github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/go-logr/logr"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	v1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Fetcher interface {
	DeploymentStrategyFor(canary v1beta1.Canary) DeploymentStrategy
	FetchTargetRef(ctx context.Context, clusterName string, clusterClient clustersmngr.Client, canary *v1beta1.Canary) (v1.Deployment, error)
	FetchPromoted(ctx context.Context, clusterName string, clusterClient clustersmngr.Client, canary *v1beta1.Canary) (v1.Deployment, error)
	GetCanary(ctx context.Context, client clustersmngr.Client, opts GetCanaryOptions) (*v1beta1.Canary, error)
	GetMetricTemplate(ctx context.Context, clusterName string, clusterClient clustersmngr.Client, name, namespace string) (v1beta1.MetricTemplate, error)
	ListCanaryDeployments(ctx context.Context, client clustersmngr.Client, opts ListCanaryDeploymentsOptions) (map[string][]v1beta1.Canary, string, []CanaryListError, error)
	ListMetricTemplates(ctx context.Context, clusterClient clustersmngr.Client, options ListMetricTemplatesOptions) (map[string][]v1beta1.MetricTemplate, string, []MetricTemplateListError, error)
	ListCanaryObjects(ctx context.Context, clusterClient clustersmngr.Client, opts ListCanaryObjectsOptions) ([]unstructured.Unstructured, error)
}

func NewFetcher(crdService crd.Fetcher, logger logr.Logger) Fetcher {
	fetcher := &defaultFetcher{crdService: crdService, logger: logger}

	return fetcher
}

type defaultFetcher struct {
	crdService crd.Fetcher
	logger     logr.Logger
}

type ListCanaryDeploymentsOptions struct {
	Namespace string
	PageSize  int32
	PageToken string
}

type ListMetricTemplatesOptions struct {
	Namespace string
	PageSize  int32
	PageToken string
}

type GetCanaryOptions struct {
	Name        string
	Namespace   string
	ClusterName string
}

type ListCanaryObjectsOptions struct {
	Name        string
	Namespace   string
	ClusterName string
}

func (service *defaultFetcher) ListCanaryDeployments(
	ctx context.Context,
	clusterClient clustersmngr.Client,
	options ListCanaryDeploymentsOptions,
) (map[string][]v1beta1.Canary, string, []CanaryListError, error) {
	var respErrors []CanaryListError

	clist := clustersmngr.NewClusteredList(func() client.ObjectList {
		return &v1beta1.CanaryList{}
	})

	opts := []client.ListOption{}
	if options.PageSize != 0 {
		opts = append(opts, client.Limit(options.PageSize))
	}

	if options.PageToken != "" {
		opts = append(opts, client.Continue(options.PageToken))
	}

	if err := clusterClient.ClusteredList(ctx, clist, true, opts...); err != nil {
		var errs clustersmngr.ClusteredListError
		if !errors.As(err, &errs) {
			return nil, "", respErrors, err
		}

		for _, e := range errs.Errors {
			// If flagger is not installed, skip all errors reported from that
			// cluster, an extra error will be appended to the error list later if
			// Flagger is not available.
			if service.crdService.IsAvailable(e.Cluster, crd.FlaggerCRDName) {
				respErrors = append(respErrors, CanaryListError{ClusterName: e.Cluster, Err: e.Err})
			}
		}
	}

	results := map[string][]v1beta1.Canary{}

	for clusterName, lists := range clist.Lists() {
		// The error will be in there from ClusteredListError, adding an extra
		// error so it's easier to check them on client side.
		if !service.crdService.IsAvailable(clusterName, crd.FlaggerCRDName) {
			respErrors = append(
				respErrors,
				CanaryListError{
					ClusterName: clusterName,
					Err:         FlaggerIsNotAvailableError{ClusterName: clusterName},
				},
			)
			results[clusterName] = []v1beta1.Canary{}

			continue
		}

		for _, l := range lists {
			list, ok := l.(*v1beta1.CanaryList)
			if !ok {
				continue
			}

			results[clusterName] = append(results[clusterName], list.Items...)
		}
	}

	return results, clist.GetContinue(), respErrors, nil
}

func (service *defaultFetcher) GetCanary(
	ctx context.Context,
	clustersClient clustersmngr.Client,
	opts GetCanaryOptions,
) (*v1beta1.Canary, error) {
	k := &v1beta1.Canary{}
	key := client.ObjectKey{
		Name:      opts.Name,
		Namespace: opts.Namespace,
	}

	if err := clustersClient.Get(ctx, opts.ClusterName, key, k); err != nil {
		return nil, fmt.Errorf("failed getting canary: name=%s namespace=%s cluster=%s err=%w", opts.Name, opts.Namespace, opts.ClusterName, err)
	}

	return k, nil
}

func (service *defaultFetcher) GetMetricTemplate(
	ctx context.Context,
	clusterName string,
	clusterClient clustersmngr.Client,
	name, namespace string,
) (v1beta1.MetricTemplate, error) {
	object := v1beta1.MetricTemplate{}

	key := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}

	err := clusterClient.Get(ctx, clusterName, key, &object)

	return object, err
}

func (service *defaultFetcher) FetchTargetRef(
	ctx context.Context,
	clusterName string,
	clusterClient clustersmngr.Client,
	canary *v1beta1.Canary,
) (v1.Deployment, error) {
	return getDeployment(ctx, clusterName, clusterClient, canary.Spec.TargetRef.Name, canary.GetNamespace())
}

func (service *defaultFetcher) FetchPromoted(
	ctx context.Context,
	clusterName string,
	clusterClient clustersmngr.Client,
	canary *v1beta1.Canary,
) (v1.Deployment, error) {
	name := fmt.Sprintf("%s-primary", canary.Spec.TargetRef.Name)
	return getDeployment(ctx, clusterName, clusterClient, name, canary.GetNamespace())
}

func (service *defaultFetcher) ListMetricTemplates(
	ctx context.Context,
	clusterClient clustersmngr.Client,
	options ListMetricTemplatesOptions,
) (map[string][]v1beta1.MetricTemplate, string, []MetricTemplateListError, error) {
	var respErrors []MetricTemplateListError

	clist := clustersmngr.NewClusteredList(func() client.ObjectList {
		return &v1beta1.MetricTemplateList{}
	})

	opts := []client.ListOption{}
	if options.PageSize != 0 {
		opts = append(opts, client.Limit(options.PageSize))
	}

	if options.PageToken != "" {
		opts = append(opts, client.Continue(options.PageToken))
	}

	if err := clusterClient.ClusteredList(ctx, clist, true, opts...); err != nil {
		var errs clustersmngr.ClusteredListError
		if !errors.As(err, &errs) {
			return nil, "", respErrors, err
		}

		for _, e := range errs.Errors {
			// If flagger is not installed, skip all errors reported from that
			// cluster, an extra error will be appended to the error list later if
			// Flagger is not available.
			if service.crdService.IsAvailable(e.Cluster, crd.FlaggerCRDName) {
				respErrors = append(respErrors, MetricTemplateListError{ClusterName: e.Cluster, Err: e.Err})
			}
		}
	}

	results := map[string][]v1beta1.MetricTemplate{}

	for clusterName, lists := range clist.Lists() {
		// The error will be in there from ClusteredListError, adding an extra
		// error so it's easier to check them on client side.
		if !service.crdService.IsAvailable(clusterName, crd.FlaggerCRDName) {
			respErrors = append(
				respErrors,
				MetricTemplateListError{
					ClusterName: clusterName,
					Err:         FlaggerIsNotAvailableError{ClusterName: clusterName},
				},
			)
			results[clusterName] = []v1beta1.MetricTemplate{}

			continue
		}

		for _, l := range lists {
			list, ok := l.(*v1beta1.MetricTemplateList)
			if !ok {
				continue
			}

			results[clusterName] = append(results[clusterName], list.Items...)
		}
	}

	return results, clist.GetContinue(), respErrors, nil
}

func (service *defaultFetcher) ListCanaryObjects(ctx context.Context, clusterClient clustersmngr.Client, opts ListCanaryObjectsOptions) ([]unstructured.Unstructured, error) {
	result := []unstructured.Unstructured{}
	checkDup := map[types.UID]bool{}

	// Get canary object
	canary, err := service.GetCanary(ctx, clusterClient, GetCanaryOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("unable to find canary object: %w", err)
	}

	targetDeployment, err := getRef(
		ctx,
		clusterClient,
		&canary.Spec.TargetRef,
		canary.GetNamespace(),
		opts.ClusterName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed getting canary target reference: %w", err)
	}
	result = append(result, targetDeployment)

	if canary.Spec.IngressRef != nil {
		ingress, err := getRef(
			ctx,
			clusterClient,
			canary.Spec.IngressRef,
			canary.GetNamespace(),
			opts.ClusterName,
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
			opts.ClusterName,
		)
		if err == nil {
			result = append(result, hpa)
		}
	}

	// List of kinds all canaries generate independently of mesh provider
	coreObjectsKinds := []schema.GroupVersionKind{
		{Group: "", Version: "v1", Kind: "Service"},
		{Group: "apps", Version: "v1", Kind: "Deployment"},
		{Group: "autoscaling", Version: "v2beta1", Kind: "HorizontalPodAutoscaler"},
	}

	objectsKinds := append(coreObjectsKinds, meshProviderObjectKinds(canary.Spec.Provider)...)

	for _, gvk := range objectsKinds {
		listResult := unstructured.UnstructuredList{}

		listResult.SetGroupVersionKind(gvk)

		if err := clusterClient.List(ctx, opts.ClusterName, &listResult); err != nil {
			if k8serrors.IsForbidden(err) {
				service.logger.Error(err, "request is forbidden", "cluster", opts.ClusterName)

				continue
			}

			// Given Flux supports multiple version of the same CRD we need to avoid
			// breaking when we query a version that's not present on the cluster.
			if apimeta.IsNoMatchError(err) {
				service.logger.Error(err, "failed listing mesh provider resource", "cluster", opts.ClusterName)
				continue
			}

			return nil, fmt.Errorf("error listing unstructured object: %w", err)
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

	return result, nil
}

func getDeployment(ctx context.Context, clusterName string, c clustersmngr.Client, name string, namespace string) (v1.Deployment, error) {
	deployment := v1.Deployment{}

	key := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}

	err := c.Get(ctx, clusterName, key, &deployment)

	return deployment, err
}

func getRef(ctx context.Context, clusterClient clustersmngr.Client, ref *v1beta1.LocalObjectReference, ns string, clusterName string) (unstructured.Unstructured, error) {
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

func meshProviderObjectKinds(provider string) []schema.GroupVersionKind {
	kinds := []schema.GroupVersionKind{}

	if provider == "" {
		return kinds
	}

	switch provider {
	case flaggerv1.AppMeshProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "appmesh.k8s.aws", Version: "v1beta2", Kind: "virtualnode"},
			{Group: "appmesh.k8s.aws", Version: "v1beta2", Kind: "virtualrouter"},
			{Group: "appmesh.k8s.aws", Version: "v1beta2", Kind: "virtualservice"},
		}
	case flaggerv1.LinkerdProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "split.smi-spec.io", Version: "v1alpha1", Kind: "trafficsplit"},
			{Group: "split.smi-spec.io", Version: "v1alpha2", Kind: "trafficsplit"},
			{Group: "split.smi-spec.io", Version: "v1alpha3", Kind: "trafficsplit"},
		}
	case flaggerv1.IstioProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "networking.istio.io", Version: "v1alpha3", Kind: "destinationrule"},
			{Group: "networking.istio.io", Version: "v1alpha3", Kind: "virtualservice"},
		}
	case flaggerv1.ContourProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "projectcontour.io", Version: "v1", Kind: "httpproxy"},
		}
	case flaggerv1.GlooProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "gateway.solo.io", Version: "v1", Kind: "routetable"},
			{Group: "gloo.solo.io", Version: "v1", Kind: "upstream"},
		}
	case flaggerv1.NGINXProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"},
		}
	case flaggerv1.SkipperProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"},
		}
	case flaggerv1.TraefikProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "traefik.containo.us", Version: "v1alpha1", Kind: "traefikservice"},
		}
	case flaggerv1.OsmProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "split.smi-spec.io", Version: "v1alpha3", Kind: "trafficsplit"},
		}
	case flaggerv1.KumaProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "kuma.io", Version: "v1alpha1", Kind: "trafficroute"},
		}
	case flaggerv1.GatewayAPIProvider:
		kinds = []schema.GroupVersionKind{
			{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "httproute"},
		}
	case flaggerv1.KubernetesProvider:
		kinds = []schema.GroupVersionKind{}
	}

	return kinds
}
