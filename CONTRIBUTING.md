# Developing `progressive-delivery`

A guide to making it easier to develop `progressive-delivery`. If you came here
expecting but not finding an answer please make an issue to help improve these
docs!

## Run a local development environment

To run a local development environment, you need to install
[Docker](https://www.docker.com) and
[kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/), other
dependencies can be installed with `make dependencies`.

### Preparation

Create a kind cluster:

```bash
./tools/bin/kind create cluster --name pd-dev
```

If something goes wrong, delete the cluster and re-create it:

```bash
./tools/bin/kind delete cluster --name pd-dev
./tools/bin/kind create cluster --name pd-dev
```

Deploy flux on the cluster:

```bash
export GITHUB_USER="<your-username>"
export GITHUB_REPO="${GITHUB_REPO:-pd-dev}"
./tools/bin/flux bootstrap github \
  --owner="$GITHUB_USER" \
  --repository="$GITHUB_REPO" \
  --branch=main \
  --path=./clusters/management \
  --personal
```

To install extra resources, use the `./tools/install-resources.sh` script:

```bash
./tools/install-resources.sh -h
usage: ./tools/install-resources.sh [-i] [-f] [-c]

Install extra resources.

OPTIONS:
   -c|--canaries     Install Canary objects
   -f|--flagger      Install Flagger
   -i|--istio        Install Istio
   -h|--help         Show this message
```


### Start environment

To start the development environment, run:

```bash
make dev-cluster
```

Your system should build and start. The first time you run this, it will
take ~1-2 mins (depending on your connection speed) to build the container and
deploy it to your local cluster. This is because the docker builds have to
download all the Go modules from scratch, use the Tilt UI to check progress.
Subsequent runs should be a lot faster.

### Making API requests to the dev cluster

Use a gRPC client to interact with the API, for example:

* [BloomRPC](https://github.com/bloomrpc/bloomrpc) has a nice GUI.
* [gRPCurl](https://github.com/fullstorydev/grpcurl) can be used from command
    line.

### Example queries

**Check if Flagger is available on the cluster or not**

```bash
❯ grpcurl -plaintext localhost:9002 ProgressiveDeliveryService.IsFlaggerAvailable
```

```json
{
  "clusters": {
    "Default": true
  }
}
```

**List available canaries**

```bash
grpcurl -plaintext localhost:9002 ProgressiveDeliveryService.ListCanaries
```
```json
{
  "canaries": [
    {
      "name": "hello-world",
      "clusterName": "Default",
      "provider": "traefik",
      "targetReference": {
        "kind": "Deployment",
        "name": "hello-world"
      },
      "targetDeployment": {
        "uid": "4b871207-63e7-4981-b067-395c59b3676b",
        "resourceVersion": "1997",
        "fluxLabels": {
          "kustomizeNamespace": "hello-world",
          "kustomizeName": "hello-world"
        }
      },
      "status": {
        "phase": "Initialized",
        "lastTransitionTime": "2022-06-03T12:36:23Z",
        "conditions": [
          {
            "type": "Promoted",
            "status": "True",
            "lastUpdateTime": "2022-06-03T12:36:23Z",
            "lastTransitionTime": "2022-06-03T12:36:23Z",
            "reason": "Initialized",
            "message": "Deployment initialization completed."
          }
        ]
      }
    }
  ],
  "nextPageToken": "eyJDb250aW51ZVRva2VucyI6eyJEZWZhdWx0Ijp7ImNhcGQtc3lzdGVtIjoiIiwiY2FwaS1rdWJlYWRtLWJvb3RzdHJhcC1zeXN0ZW0iOiIiLCJjYXBpLWt1YmVhZG0tY29udHJvbC1wbGFuZS1zeXN0ZW0iOiIiLCJjYXBpLXN5c3RlbSI6IiIsImNlcnQtbWFuYWdlciI6IiIsImRlZmF1bHQiOiIiLCJkZXgiOiIiLCJmbGFnZ2VyIjoiIiwiZmx1eC1zeXN0ZW0iOiIiLCJoZWxsby13b3JsZCI6IiIsImt1YmUtbm9kZS1sZWFzZSI6IiIsImt1YmUtcHVibGljIjoiIiwia3ViZS1zeXN0ZW0iOiIiLCJsb2NhbC1wYXRoLXN0b3JhZ2UiOiIiLCJ0cmFlZmlrIjoiIn19fQo="
}
```

**Get a specific Canary object**

```bash
❯ grpcurl \
    -d '{"clusterName": "Default", "name": "hello-world", "namespace": "hello-world"}' \
    -plaintext localhost:9002 ProgressiveDeliveryService.GetCanary
```

```json
{
  "canary": {
    "namespace": "hello-world",
    "name": "hello-world",
    "clusterName": "Default",
    "provider": "traefik",
    "targetReference": {
      "kind": "Deployment",
      "name": "hello-world"
    },
    "targetDeployment": {
      "uid": "4b871207-63e7-4981-b067-395c59b3676b",
      "resourceVersion": "3394152",
      "fluxLabels": {
        "kustomizeNamespace": "hello-world",
        "kustomizeName": "hello-world"
      }
    },
    "status": {
      "phase": "Succeeded",
      "failedChecks": 1,
      "lastTransitionTime": "2022-06-10T10:29:03Z",
      "conditions": [
        {
          "type": "Promoted",
          "status": "True",
          "lastUpdateTime": "2022-06-10T10:29:03Z",
          "lastTransitionTime": "2022-06-10T10:29:03Z",
          "reason": "Succeeded",
          "message": "Canary analysis completed successfully, promotion finished."
        }
      ]
    },
    "deploymentStrategy": "canary"
  },
  "automation": {
    "kind": "Kustomization",
    "name": "hello-world",
    "namespace": "hello-world"
  }
}
```

## Coding standards

* For `proto` files, follow [ProtoBuf Style Guide][pb-style].
* For `go` files, `make lint` will tell at your.
* For testing, we prefer [black-box testing][bb-testing].
   (tl;dr: use `mypackage_test` packages)

[pb-style]: https://developers.google.com/protocol-buffers/docs/style
[bb-testing]: https://en.wikipedia.org/wiki/Black-box_testing
