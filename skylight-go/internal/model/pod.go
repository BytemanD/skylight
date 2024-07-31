package model

import (
	v1 "k8s.io/api/core/v1"
)

type Deletion struct {
	GracePeriodSeconds *int64 `json:"grace_period_seconds,omitempty"`
	Timestamp          string `json:"timestamp,omitempty"`
}

type Pod struct {
	BaseModel
	Labels       map[string]string `json:"labels,omitempty"`
	NodeName     string            `json:"node_name,omitempty"`
	NodeSelector map[string]string `json:"node_selector,omitempty"`
	HostIp       string            `json:"host_ip,omitempty"`
	PodIp        string            `json:"pod_ip,omitempty"`
	Containers   []Container       `json:"containers,omitempty"`
	Deletion     Deletion          `json:"deletion,omitempty"`
	Phase        string            `json:"phase,omitempty"`
}

func ParseV1Pod(item v1.Pod) Pod {
	containers := []Container{}
	containerStatusMap := map[string]v1.ContainerStatus{}
	for _, containerStatus := range item.Status.ContainerStatuses {
		containerStatusMap[containerStatus.Name] = containerStatus
	}
	for _, c := range item.Spec.Containers {
		if cs, ok := containerStatusMap[c.Name]; ok {
			containers = append(containers, ParseV1Container(c, &cs))
		} else {
			containers = append(containers, ParseV1Container(c, nil))
		}
	}
	pod := Pod{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Status:   string(item.Status.Phase),
		},
		Labels:       item.Labels,
		NodeName:     item.Spec.NodeName,
		NodeSelector: item.Spec.NodeSelector,
		HostIp:       item.Status.HostIP,
		PodIp:        item.Status.PodIP,
		Containers:   containers,
		Deletion: Deletion{
			GracePeriodSeconds: item.DeletionGracePeriodSeconds,
		},
		Phase: string(item.Status.Phase),
	}
	if item.DeletionTimestamp != nil {
		pod.Deletion.Timestamp = item.DeletionTimestamp.Format("2006-01-02 15:04:05")
	}
	return pod
}
