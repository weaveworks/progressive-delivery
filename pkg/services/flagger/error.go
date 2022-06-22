package flagger

import "fmt"

type FlaggerIsNotAvailableError struct {
	ClusterName string
}

func (e FlaggerIsNotAvailableError) Error() string {
	return fmt.Sprintf("flagger is not installed on cluster: %s", e.ClusterName)
}

type CanaryListError struct {
	ClusterName string
	Err         error
}

func (e CanaryListError) Error() string {
	return fmt.Sprintf("canary list error on cluster %s: %s", e.ClusterName, e.Err.Error())
}

type MetricTemplateListError struct {
	ClusterName string
	Err         error
}

func (e MetricTemplateListError) Error() string {
	return fmt.Sprintf("metric template list error on cluster %s: %s", e.ClusterName, e.Err.Error())
}
