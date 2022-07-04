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

##  Proposed Solution (What / How)

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


https://miro.com/app/board/uXjVOoAPntE=/?share_link_id=505858014497







##  Alternatives Considered 

### Limitations


## References

## Log 

- July/2022 initial draft










