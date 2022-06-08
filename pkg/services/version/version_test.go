package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/pkg/services/version"
)

func TestVersionFetcher(t *testing.T) {
	v := version.NewFetcher()

	expected := "v0.0.0"
	result := v.Get().Semver

	assert.Equal(t, expected, result)
}
