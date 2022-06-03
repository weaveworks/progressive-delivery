# progressive-delivery
This repository contains the progressive delivery API handlers that Weave GitOps Enterprise serves.

## Local development environment

To install all dependencies use `make dependencies`.

## Start dev server

Use Tilt to deploy the dev-server on the cluster:

```bash
./tools/bin/tilt up
```

Use ea gRPC client to interact with the API, for example:

* [BloomRPC](https://github.com/bloomrpc/bloomrpc) has a nice GUI.
* [gRPCurl](https://github.com/fullstorydev/grpcurl) can be used from command
    line.

```bash
❯ grpcurl -plaintext localhost:9002 ProgressiveDeliveryService.ListCanaries
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
