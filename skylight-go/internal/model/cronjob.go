package model

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

type CronJob struct {
	BaseModel
	Labels       map[string]string `json:"labels,omitempty"`
	NodeSelector map[string]string `json:"node_selector,omitempty"`
	MatchLabels  map[string]string `json:"selector,omitempty"`
	Containers   []Container       `json:"containers,omitempty"`
	Deletion     Deletion          `json:"deletion,omitempty"`
	Phase        string            `json:"phase,omitempty"`
	Schedule     string            `json:"schedule,omitempty"`
}

func ParseV1CronJob(item batchv1.CronJob) CronJob {
	template := item.Spec.JobTemplate.Spec.Template
	containers := []Container{}
	for _, c := range template.Spec.Containers {
		containers = append(containers, ParseV1Container(c, nil))
	}
	return CronJob{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Labels:       item.Labels,
		NodeSelector: template.Spec.NodeSelector,
		MatchLabels:  item.Spec.JobTemplate.Spec.Selector.MatchLabels,
		Containers:   containers,
	}
}

func ParseV1betaCronJob(item batchv1beta1.CronJob) CronJob {
	template := item.Spec.JobTemplate.Spec.Template
	containers := []Container{}
	for _, c := range template.Spec.Containers {
		containers = append(containers, ParseV1Container(c, nil))
	}

	return CronJob{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Schedule:     item.Spec.Schedule,
		Labels:       item.Labels,
		NodeSelector: template.Spec.NodeSelector,
		Containers:   containers,
	}
}
