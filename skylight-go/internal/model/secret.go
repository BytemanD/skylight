package model

import (
	corev1 "k8s.io/api/core/v1"
)

type Secret struct {
	BaseModel
	Data map[string][]byte `json:"data,omitempty"`
}

func ParseV1Secret(item corev1.Secret) Secret {
	return Secret{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Data: item.Data,
	}
}
