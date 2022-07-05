# Canary Analysis - Canary Metrics Design Proposal

This document outlines a proposed design to extend gitops enterprise progressive delivery domain with 
canary analysis capabilities.

## Problem statement 

As wget user doing progressive delivery for my applications, I am not able to understand via the enterprise experience 
the canary analysis section: 
- reducing the user levels of confidence by not having it visible to 
- decreasing its experience by having to use other tooling to get that information 

To address this gap the [canary analysis epic](https://github.com/weaveworks/weave-gitops-enterprise/issues/842) has been 
created and this document addresses the design for the first of the analysis tab sections: metrics.

the user journeys considered are:

- as wge user, i want to have an overview understanding of the metrics used for an application canary analisys.  
- as wge user, i want to have a detailed understanding of the metrics used for an application canary analisys.

## Scope 
**In scope** 
  - canary analysis metrics backend apis.
**Out of scope**
  - canary analysis webhooks and alerts not in scope to simplify design. a similar approach could be followed for them. 
  - frontend due to knowledge limitations.  

##  Proposed Solution (What / How)

[Board](https://miro.com/app/board/uXjVOoAPntE=/?share_link_id=505858014497)

Current api has the [following types](../../../api/prog/types.proto) 

```protobuf
message Canary {
  string namespace = 1;
  string name = 2;
  string cluster_name = 3;
  string provider = 4;
  CanaryTargetReference target_reference = 5;
  CanaryTargetDeployment target_deployment = 6;
  CanaryStatus status = 7;
  string deploymentStrategy = 8;
  CanaryAnalysis analysis = 9;
  string yaml = 10;
}

`message CanaryAnalysis {
  string interval = 1;
  int32 iterations = 2;
  int32 mirror_weight = 3;
  int32 max_weight = 4;
  int32 step_weight = 5;
  int32 step_weight_promotion = 6;
  int32 threshold = 7;
  repeated int32 step_weights = 8;
  bool mirror = 9;
  string yaml = 10;
  repeated CanaryMetricTemplate metric_templates = 11;
}

message CanaryMetricTemplate {
  string cluster_name = 1;
  string name = 2;
  string namespace = 3;
  MetricProvider provider = 4;
  string query = 5;
}
````
In order to achieve metrics visibility a couple of alternatives could be taken 

- Alternative A: to include metrics template within metrics analysis response.
- Alternative B: to do not include metric templates and to have a metrics template endpoint. 


### Alternative A: include metrics template within metrics analysis response 

```protobuf

message CanaryAnalysis {
  string interval = 1;
  int32 iterations = 2;
  int32 mirror_weight = 3;
  int32 max_weight = 4;
  int32 step_weight = 5;
  int32 step_weight_promotion = 6;
  int32 threshold = 7;
  repeated int32 step_weights = 8;
  bool mirror = 9;
  string yaml = 10;
  repeated CanaryMetric metrics = 11;
}

message CanaryMetric {
  string cluster_name = 1;
  string name = 2;
  string namespace = 3;
  CanaryMetricThresholdRange threshold_range = 4
  string interval = 5
  CanaryMetricTemplate metric_template = 6
}

message CanaryMetricThresholdRange {
  string min = 1;
  string max = 2;
}

message CanaryMetricTemplate {
  string cluster_name = 1;
  string name = 2;
  string namespace = 3;
  MetricProvider provider = 4;
  string query = 5;
}


```

**Journeys Validation** 

Given the request `GET /v1/pd/canaries/my-canary`
And a response like
```json
{
  "name": "my-canary",
  ...
  "analysis": {
    ...
    "metrics": [
      {
        "cluster_name": "my-cluster",
        "name": "request-success-rate",
        "threshold_range": {
          "min": "90",
          "max": "99"
        },
        "interval": "1m"
      },
      {
        "cluster_name": "my-cluster",
        "name": "my-awesome-custom-metric",
        "threshold_range": {
          "min": "90",
          "max": "99"
        },
        "interval": "1m",
        "metric_template": {
          ...,
          "name": "my-awesome-custom-metric-template",
          "provider": {
            "type":"datatdog",
          },
          "query": "my-datadog-query",
        }
      }
    ]
  }
}
```

When **As user, I want to have an overview view of a canary metrics**
That experience contains the information from the response 
```
        "cluster_name": "my-cluster",
        "name": "request-success-rate",
        "threshold_range": {
          "min": "90",
          "max": "99"
        },
        "interval": "1m"
```

When **As user, I want to have an detailed view on a metric template**
That experience contains the information from the response
```
  "metric_template": {
          ...,
          "provider": {
            "type":"datatdog",
          },
          "query": "my-datadog-query",
        }
      }
```

### Alternative B: to do not include metric templates and to have a metrics template endpoint.

Same as Alternative A but with the difference that instead of having the metric templte we have a reference to it

```protobuf

message CanaryMetric {
  string cluster_name = 1;
  string name = 2;
  string namespace = 3;
  CanaryMetricThresholdRange threshold_range = 4
  string interval = 5
  CanaryMetricTemplateRef metric_template = 6
}

message CanaryMetricTemplateRef {
  string cluster_name = 1;
  string name = 2;
  string namespace = 3;
}


```

That could be resolved via a new api endpoint "/v1/pd/metric_templates/"

```json

 "/v1/pd/metric_templates/{name}": {
      "get": {
        "summary": "GetMetricTemplate returns a MetricTemplate object.",
        "operationId": "ProgressiveDeliveryService_GetMetricTemplatey",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/GetMetricTemplateResponse"
            }
          },
          "default": {
            "description": "An unexpected error response.",
            "schema": {
              "$ref": "#/definitions/rpcStatus"
            }
          }
        },
        "parameters": [
          {
            "name": "name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "namespace",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "clusterName",
            "in": "query",
            "required": false,
            "type": "string"
          }
        ],
        "tags": [
          "ProgressiveDeliveryService"
        ]
      }
    },
``` 

with the following response type 

```json
    "GetMetricTemplateResponse": {
      "type": "object",
      "properties": {
        "metric_template": {
          "$ref": "#/definitions/MetricTemplate"
        },
      }
    },
```

and metric template as specified in Alternative A

**Journeys Validation** 

1 - Request `GET /v1/pd/canaries/my-canary` to achieve overview view
And a response like within alternative A

2- Request `GET /v1/pd/metric_templates/my-awesome-custom-metric-template` to get the metric template view when required a detailed view 


## References
- [Flagger How it works](https://docs.flagger.app/usage/how-it-works#canary-analysis)
- [Flagger Metrics Analysis](https://docs.flagger.app/usage/metrics)

## Log 

- July/2022 initial draft










