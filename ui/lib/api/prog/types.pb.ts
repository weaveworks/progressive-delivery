/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
export type Pagination = {
  pageSize?: number
  pageToken?: string
}

export type ListError = {
  clusterName?: string
  namespace?: string
  message?: string
}

export type Canary = {
  namespace?: string
  name?: string
  clusterName?: string
  provider?: string
  targetReference?: CanaryTargetReference
  targetDeployment?: CanaryTargetDeployment
  status?: CanaryStatus
  deploymentStrategy?: string
  analysis?: CanaryAnalysis
  yaml?: string
}

export type CanaryTargetReference = {
  kind?: string
  name?: string
}

export type CanaryStatus = {
  phase?: string
  failedChecks?: number
  canaryWeight?: number
  iterations?: number
  lastTransitionTime?: string
  conditions?: CanaryCondition[]
}

export type CanaryCondition = {
  type?: string
  status?: string
  lastUpdateTime?: string
  lastTransitionTime?: string
  reason?: string
  message?: string
}

export type CanaryTargetDeployment = {
  uid?: string
  resourceVersion?: string
  fluxLabels?: FluxLabels
  appliedImageVersions?: {[key: string]: string}
  promotedImageVersions?: {[key: string]: string}
}

export type FluxLabels = {
  kustomizeNamespace?: string
  kustomizeName?: string
}

export type Automation = {
  kind?: string
  name?: string
  namespace?: string
}

export type CanaryAnalysis = {
  interval?: string
  iterations?: number
  mirrorWeight?: number
  maxWeight?: number
  stepWeight?: number
  stepWeightPromotion?: number
  threshold?: number
  stepWeights?: number[]
  mirror?: boolean
  yaml?: string
  metrics?: CanaryMetric[]
}

export type CanaryMetric = {
  name?: string
  namespace?: string
  thresholdRange?: CanaryMetricThresholdRange
  interval?: string
  metricTemplate?: CanaryMetricTemplate
}

export type CanaryMetricThresholdRange = {
  min?: number
  max?: number
}

export type CanaryMetricTemplate = {
  clusterName?: string
  name?: string
  namespace?: string
  provider?: MetricProvider
  query?: string
}

export type MetricProvider = {
  type?: string
  address?: string
  secretName?: string
  insecureSkipVerify?: boolean
}

export type GroupVersionKind = {
  group?: string
  kind?: string
  version?: string
}

export type UnstructuredObject = {
  groupVersionKind?: GroupVersionKind
  name?: string
  namespace?: string
  uid?: string
  status?: string
  conditions?: Condition[]
  suspended?: boolean
  clusterName?: string
  images?: string[]
}

export type Condition = {
  type?: string
  status?: string
  reason?: string
  message?: string
  timestamp?: string
}