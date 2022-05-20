package server_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/weaveworks/progressive-delivery/pkg/api/prog"
	"github.com/weaveworks/progressive-delivery/pkg/server"
)

func TestGetVersion(t *testing.T) {
	mux := runtime.NewServeMux()
	err := server.Hydrate(context.Background(), mux, server.ServerOpts{})
	if err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(mux)
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
