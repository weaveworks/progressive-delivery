package flagger

import (
	"context"
	"errors"
	"fmt"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/weaveworks/progressive-delivery/pkg/services/crd"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Fetcher interface {
	ListCanaryDeployments(ctx context.Context, client clustersmngr.Client, opts ListCanaryDeploymentsOptions) (map[string][]v1beta1.Canary, string, []CanaryListError, error)
	GetCanary(ctx context.Context, client clustersmngr.Client, opts GetCanaryOptions) (*v1beta1.Canary, error)
	FetchTargetRef(ctx context.Context, clusterName string, clusterClient clustersmngr.Client, canary *v1beta1.Canary) (v1.Deployment, error)
	DeploymentStrategyFor(canary v1beta1.Canary) DeploymentStrategy
}

func NewFetcher(crdService crd.Fetcher) Fetcher {
	fetcher := &defaultFetcher{crdService: crdService}

	return fetcher
}

type defaultFetcher struct {
	crdService crd.Fetcher
}

type ListCanaryDeploymentsOptions struct {
	Namespace string
	PageSize  int32
	PageToken string
}

type GetCanaryOptions struct {
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
			respErrors = append(respErrors, CanaryListError{ClusterName: e.Cluster, Err: e.Err})
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

func (service *defaultFetcher) FetchTargetRef(
	ctx context.Context,
	clusterName string,
	clusterClient clustersmngr.Client,
	canary *v1beta1.Canary,
) (v1.Deployment, error) {
	deployment := v1.Deployment{}

	key := client.ObjectKey{
		Name:      canary.Spec.TargetRef.Name,
		Namespace: canary.GetNamespace(),
	}

	err := clusterClient.Get(ctx, clusterName, key, &deployment)

	return deployment, err
}
