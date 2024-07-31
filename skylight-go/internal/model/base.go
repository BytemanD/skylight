package model

import (
	v1 "k8s.io/api/core/v1"
)

type BaseModel struct {
	Name     string `json:"name,omitempty"`
	Creation string `json:"creation,omitempty"`
	Status   string `json:"status,omitempty"`
}

type ContainerStatus struct {
	Ready        bool                     `json:"ready,omitempty"`
	State        string                   `json:"state,omitempty"`
	LastRunning  string                   `json:"last_running,omitempty"`
	LasteWaiting v1.ContainerStateWaiting `json:"last_waiting,omitempty"`
}

type Container struct {
	Name            string             `json:"name,omitempty"`
	Command         []string           `json:"command,omitempty"`
	ImagePullPolicy string             `json:"image_pull_policy,omitempty"`
	Ports           []v1.ContainerPort `json:"ports,omitempty"`
	Status          ContainerStatus    `json:"status,omitempty"`
	Image           string             `json:"image,omitempty"`
}

func ParseV1ContainerStatus(item v1.ContainerStatus) ContainerStatus {
	return ContainerStatus{
		Ready:        item.Ready,
		State:        item.State.String(),
		LastRunning:  item.LastTerminationState.Running.StartedAt.Format(""),
		LasteWaiting: *item.LastTerminationState.Waiting,
	}
}

func ParseV1Container(item v1.Container, containerStatus *v1.ContainerStatus) Container {
	container := Container{
		Name:            item.Name,
		Command:         item.Command,
		ImagePullPolicy: string(item.ImagePullPolicy),
		Ports:           item.Ports,
		Image:           item.Image,
	}
	return container
}
