package model

import (
	batchv1 "k8s.io/api/batch/v1"
)

type Job struct {
	BaseModel
	Labels       map[string]string `json:"labels,omitempty"`
	NodeSelector map[string]string `json:"node_selector,omitempty"`
	MatchLabels  map[string]string `json:"selector,omitempty"`
	Containers   []Container       `json:"container,omitempty"`
	Deletion     Deletion          `json:"deletion,omitempty"`
	Phase        string            `json:"phase,omitempty"`
}

func ParseV1Job(item batchv1.Job) Job {
	containers := []Container{}
	for _, c := range item.Spec.Template.Spec.Containers {
		containers = append(containers, ParseV1Container(c, nil))
	}
	return Job{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Labels:       item.Labels,
		NodeSelector: item.Spec.Template.Spec.NodeSelector,
		MatchLabels:  item.Spec.Selector.MatchLabels,
		Containers:   containers,
	}
}
