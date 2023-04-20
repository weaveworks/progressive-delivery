package server_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
)

func TestGetVersion(t *testing.T) {
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	response, err := c.GetVersion(ctx, &pb.GetVersionRequest{})
	assert.NoError(t, err)

	assert.Equal(t, "v0.0.0", response.GetVersion())
}

func TestHydrate(t *testing.T) {
	ts := pdtesting.MakeHTTPServer(t, k8sEnv)

	defer ts.Close()

	res, err := http.Get(ts.URL + "/v1/pd/version")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode, "should have been ok")

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	v := pb.GetVersionResponse{}
	err = json.Unmarshal(body, &v)
	assert.NoError(t, err)

	assert.Equal(t, "v0.0.0", v.Version)
}
