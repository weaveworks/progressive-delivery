package pdtesting

import (
	"context"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Cleanup(ctx context.Context, t *testing.T, k client.Client, obj client.Object) {
	if err := k.Delete(ctx, obj); err != nil {
		t.Error(err)
	}
}
