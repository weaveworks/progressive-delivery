package convert

import (
	"testing"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSerializationObj(t *testing.T) {
	canary := v1beta1.Canary{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1beta1.CanarySpec{
			Provider: "traefik",
			TargetRef: v1beta1.LocalObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       "test",
			},
		},
	}

	yaml, err := serializeObj(&canary)
	assert.NoError(t, err)

	assert.Contains(t, string(yaml), "kind: Canary")
	assert.Contains(t, string(yaml), "kind: Deployment")

}
