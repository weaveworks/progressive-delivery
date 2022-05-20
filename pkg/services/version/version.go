package version

import (
	"github.com/weaveworks/progressive-delivery/pkg/models"
)

type Fetcher interface {
	Get() models.Version
}

func NewFetcher() Fetcher {
	return defaultFetcher{}
}

type defaultFetcher struct {
}

func (v defaultFetcher) Get() models.Version {
	return models.Version{Semver: "v0.0.0"}
}
