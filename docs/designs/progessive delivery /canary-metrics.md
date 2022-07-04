# Canary Analysis - Canary Metrics Design Proposal

This document outlines a proposed design to extend gitops enterprise progressive delivery domain with 
canary analysis capabilities.

## Problem statement - Why

Its motivation is the [canary analysis epic](https://github.com/weaveworks/weave-gitops-enterprise/issues/842)
where enterprise users are willing to inspect cluster status around progressive delivery metrics to ensure 
it behaves in the expected way. 

Examples of the user stories used in this design are:

- As a user, when I view a Canary object, I can discover key information relating to how Flagger will perform 
  the canary analysis for a new version of my application, by clicking the "Analysis" tab. See Fig 1.
- As a user, I can see which metric checks have been configured - their name, their namespace, the threshold minimum, t
  the threshold maximum, and the interval, in a sortable-by-column table.
- As a user, when seeing a metric defined using a custom metric template, I can discover how this is configured 
  as the name of the metric becomes a hyperlink which when clicked, opens a modal view showing me the yaml for the relevant metric template object.

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

that introducing CanaryMetric as new field of CanaryAnaylsis we could achieve what we want... 

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

## Journeys Validation 

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


##  Alternatives Considered 

### Limitations
- Not identified 
- 
## References
- [Flagger How it works](https://docs.flagger.app/usage/how-it-works#canary-analysis)
- [Flagger Metrics Analysis](https://docs.flagger.app/usage/metrics)

## Log 

- July/2022 initial draft










