package crd

import (
	"context"
	"log"
	"sync"

	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func NewNoCacheFetcher(clientFactory clustersmngr.ClientsFactory) Fetcher {
	fetcher := &noCacheFetcher{
		clientFactory: clientFactory,
		crds:          map[string][]v1.CustomResourceDefinition{},
	}

	return fetcher
}

type noCacheFetcher struct {
	sync.RWMutex
	clientFactory clustersmngr.ClientsFactory
	crds          map[string][]v1.CustomResourceDefinition
}

func (s *noCacheFetcher) UpdateCRDList() {
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

		err := client.List(ctx, crdList)
		if err != nil {
			log.Printf("unable to list crds: %s", err)
			continue
		}

		s.crds[clusterName] = crdList.Items
	}
}

func (s *noCacheFetcher) IsAvailable(clusterName, name string) bool {
	s.UpdateCRDList()

	s.Lock()
	defer s.Unlock()

	for _, crd := range s.crds[clusterName] {
		if crd.Name == name {
			return true
		}
	}

	return false
}

func (s *noCacheFetcher) IsAvailableOnClusters(name string) map[string]bool {
	s.UpdateCRDList()

	s.Lock()
	defer s.Unlock()

	result := map[string]bool{}

	for clusterName, crds := range s.crds {
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
