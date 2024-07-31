package model

import (
	v1 "k8s.io/api/core/v1"
)

type Namespace struct {
	BaseModel
	Labels map[string]string `json:"labels,omitempty"`
}

func ParseV1Namespce(item v1.Namespace) Namespace {
	return Namespace{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Status:   string(item.Status.Phase),
		},
		Labels: item.Labels,
	}

}
