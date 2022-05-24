package version_test

import (
	"testing"

	"github.com/weaveworks/progressive-delivery/pkg/services/version"
)

func TestVersionFetcher(t *testing.T) {
	v := version.NewFetcher()

	expected := "v0.0.0"
	result := v.Get().Semver

	if result != expected {
		t.Errorf("expected result to be %v, got %v", expected, result)
	}
}
