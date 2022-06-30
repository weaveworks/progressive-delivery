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
  metricTemplates?: CanaryMetricTemplate[]
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