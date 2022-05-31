package server_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/onsi/gomega"
	"github.com/weaveworks/progressive-delivery/internal/pdtesting"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
)

func TestGetVersion(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctx := context.Background()
	c := pdtesting.MakeGRPCServer(t, k8sEnv.Rest, k8sEnv)

	response, err := c.GetVersion(ctx, &pb.GetVersionRequest{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(response.GetVersion()).To(gomega.Equal("v0.0.0"), "version should have been v0.0.0")
}

func TestHydrate(t *testing.T) {
	ts := pdtesting.MakeHTTPServer(t, k8sEnv)

	defer ts.Close()

	res, err := http.Get(ts.URL + "/v1/version")
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Error("should have been ok")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	v := pb.GetVersionResponse{}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Error(err)
	}

	if v.Version != "v0.0.0" {
		t.Error("should have been v0.0.0")
	}
}
