package model

import (
	v1 "k8s.io/api/core/v1"
)

// ready: str
// labels: list
// internal_ip: str
// kernel_version: str
// os_image: str
// container_runtime_version: str
// capacity: dict
// allocatable: dict

type ConfigMap struct {
	BaseModel
	Data map[string]string `json:"data,omitempty"`
}

func ParseV1ConfigMap(item v1.ConfigMap) ConfigMap {
	return ConfigMap{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Data: item.Data,
	}
}
