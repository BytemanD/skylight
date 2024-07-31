package model

import (
	appv1 "k8s.io/api/apps/v1"
)

type StatefulSet struct {
	BaseModel
	Labels         map[string]string `json:"labels,omitempty"`
	NodeSelector   map[string]string `json:"node_selector,omitempty"`
	MatchLabels    map[string]string `json:"selector,omitempty"`
	Containers     []Container       `json:"container,omitempty"`
	InitContainers []Container       `json:"init_containers,omitempty"`
}

func ParseV1StatefulSet(item appv1.StatefulSet) StatefulSet {
	containers := []Container{}
	initContainers := []Container{}
	for _, c := range item.Spec.Template.Spec.Containers {
		containers = append(containers, ParseV1Container(c, nil))
	}

	for _, c := range item.Spec.Template.Spec.InitContainers {
		initContainers = append(containers, ParseV1Container(c, nil))
	}
	return StatefulSet{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Labels:         item.Labels,
		NodeSelector:   item.Spec.Template.Spec.NodeSelector,
		MatchLabels:    item.Spec.Selector.MatchLabels,
		Containers:     containers,
		InitContainers: initContainers,
	}
}
