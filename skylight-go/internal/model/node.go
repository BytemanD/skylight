package model

import (
	v1 "k8s.io/api/core/v1"
)

type Node struct {
	BaseModel
	Labels                  map[string]string `json:"labels,omitempty"`
	Ready                   string            `json:"ready,omitempty"`
	InternalIp              string            `json:"internal_ip,omitempty"`
	KernelVersion           string            `json:"kernel_version,omitempty"`
	OSImage                 string            `json:"os_image,omitempty"`
	ContainerRuntimeVersion string            `json:"container_runtime_version,omitempty"`
	Capacity                map[string]string `json:"capacity,omitempty"`
	Allocatable             map[string]string `json:"allocatable,omitempty"`
}

func ParseV1ResourceList(item v1.ResourceList) map[string]string {
	resouceMap := map[string]string{}
	for name, cap := range item {
		resouceMap[string(name)] = cap.String()
	}
	return resouceMap
}
func ParseV1Node(item v1.Node) Node {
	var (
		conditionType string
		internalIP    string
	)
	for _, condition := range item.Status.Conditions {
		if condition.Type == "Ready" {
			conditionType = string(condition.Status)
			break
		}
	}
	for _, address := range item.Status.Addresses {
		if address.Type == "InternalIP" {
			internalIP = address.Address
			break
		}
	}

	return Node{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Status:   string(item.Status.Phase),
		},
		Labels:                  item.Labels,
		Ready:                   conditionType,
		InternalIp:              internalIP,
		KernelVersion:           item.Status.NodeInfo.KernelVersion,
		OSImage:                 item.Status.NodeInfo.OSImage,
		ContainerRuntimeVersion: item.Status.NodeInfo.ContainerRuntimeVersion,
		Capacity:                ParseV1ResourceList(item.Status.Capacity),
		Allocatable:             ParseV1ResourceList(item.Status.Allocatable),
	}
}
