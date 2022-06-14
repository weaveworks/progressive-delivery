package crd

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const watchCRDsFrequency = 30 * time.Second

type Fetcher interface {
	IsAvailable(clusterName, name string) bool
	IsAvailableOnClusters(name string) map[string]bool
	UpdateCRDList()
}

func NewFetcher(ctx context.Context, clientFactory clustersmngr.ClientsFactory) Fetcher {
	fetcher := &defaultFetcher{
		clientFactory: clientFactory,
		crds:          map[string][]v1.CustomResourceDefinition{},
	}

	go fetcher.watchCRDs(ctx)

	return fetcher
}

type defaultFetcher struct {
	sync.RWMutex
	clientFactory clustersmngr.ClientsFactory
	crds          map[string][]v1.CustomResourceDefinition
}

func (s *defaultFetcher) watchCRDs(ctx context.Context) {
	_ = wait.PollImmediateInfinite(watchCRDsFrequency, func() (bool, error) {
		s.UpdateCRDList()

		return false, nil
	})
}

func (s *defaultFetcher) UpdateCRDList() {
	s.Lock()
	defer s.Unlock()

	ctx := context.Background()

	client, err := s.clientFactory.GetServerClient(ctx)
	if err != nil {
		log.Printf("unable to get client pool: %s", err)

		return
	}

	for clusterName, client := range client.ClientsPool().Clients() {
		crdList := &v1.CustomResourceDefinitionList{}

		s.crds[clusterName] = []v1.CustomResourceDefinition{}

		err := client.List(ctx, crdList)
		if err != nil {
			log.Printf("unable to list crds on '%s' cluster: %s", clusterName, err)

			continue
		}

		s.crds[clusterName] = crdList.Items
	}
}

func (s *defaultFetcher) IsAvailable(clusterName, name string) bool {
	s.Lock()
	defer s.Unlock()

	for _, crd := range s.crds[clusterName] {
		if crd.Name == name {
			return true
		}
	}

	return false
}

func (s *defaultFetcher) IsAvailableOnClusters(name string) map[string]bool {
	result := map[string]bool{}

	s.Lock()
	defer s.Unlock()

	for clusterName, crds := range s.crds {
		// Set this to be sure the key is there with false value if the following
		// look does not say it's there.
		result[clusterName] = false

		for _, crd := range crds {
			if crd.Name == name {
				result[clusterName] = true
				break
			}
		}
	}

	return result
}
