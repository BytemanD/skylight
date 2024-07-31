package model

import (
	appv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	BaseModel
	Labels            map[string]string `json:"labels,omitempty"`
	Replicas          int32             `json:"replicas,omitempty"`
	ReadyReplicas     int32             `json:"ready_replicas,omitempty"`
	AvailableReplicas int32             `json:"available_replicas,omitempty"`
	UpdatedReplicas   int32             `json:"updated_replicas,omitempty"`
	Containers        []Container       `json:"containers,omitempty"`
	InitContainers    []Container       `json:"init_containers,omitempty"`
}

func ParseV1Deployment(item appv1.Deployment) Deployment {
	containers := []Container{}
	initContainers := []Container{}
	for _, c := range item.Spec.Template.Spec.Containers {
		containers = append(containers, ParseV1Container(c, nil))
	}

	for _, c := range item.Spec.Template.Spec.InitContainers {
		initContainers = append(containers, ParseV1Container(c, nil))
	}
	return Deployment{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Labels:            item.Labels,
		Replicas:          item.Status.ReadyReplicas,
		ReadyReplicas:     item.Status.ReadyReplicas,
		UpdatedReplicas:   item.Status.UpdatedReplicas,
		AvailableReplicas: item.Status.AvailableReplicas,

		Containers:     containers,
		InitContainers: initContainers,
	}
}
