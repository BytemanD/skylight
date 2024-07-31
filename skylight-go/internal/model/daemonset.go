package model

import (
	appv1 "k8s.io/api/apps/v1"
)

type Daemonset struct {
	BaseModel
	Labels                 map[string]string `json:"labels,omitempty"`
	NumberReady            int32             `json:"number_ready"`
	NumberAvailable        int32             `json:"number_available"`
	UpdatedNumberScheduled int32             `json:"updated_number_scheduled"`
	CurrentNumberScheduled int32             `json:"current_number_scheduled"`
	DesiredNumberScheduled int32             `json:"desired_number_scheduled"`
	NodeSelector           map[string]string `json:"node_selector"`
	MatchLabels            map[string]string `json:"match_labels,omitempty"`
	Containers             []Container       `json:"containers,omitempty"`
	InitContainers         []Container       `json:"init_containers,omitempty"`
}

func ParseV1Daemonset(item appv1.DaemonSet) Daemonset {
	containers := []Container{}
	initContainers := []Container{}
	for _, c := range item.Spec.Template.Spec.Containers {
		containers = append(containers, ParseV1Container(c, nil))
	}

	for _, c := range item.Spec.Template.Spec.InitContainers {
		initContainers = append(containers, ParseV1Container(c, nil))
	}
	return Daemonset{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Labels:                 item.Labels,
		NumberReady:            item.Status.NumberReady,
		NumberAvailable:        item.Status.NumberAvailable,
		UpdatedNumberScheduled: item.Status.UpdatedNumberScheduled,
		CurrentNumberScheduled: item.Status.CurrentNumberScheduled,
		DesiredNumberScheduled: item.Status.DesiredNumberScheduled,

		NodeSelector:   item.Spec.Template.Spec.NodeSelector,
		MatchLabels:    item.Spec.Selector.MatchLabels,
		Containers:     containers,
		InitContainers: initContainers,
	}
}
