package model

import (
	corev1 "k8s.io/api/core/v1"
)

type Service struct {
	BaseModel
	Type                  string               `json:"type,omitempty"`
	ClusterIp             string               `json:"cluster_ip,omitempty"`
	ClusterIPs            []string             `json:"cluster_i_ps,omitempty"`
	InternalTrafficPolicy string               `json:"internal_traffic_policy,omitempty"`
	IpFamilies            []corev1.IPFamily    `json:"ip_families,omitempty"`
	Ports                 []corev1.ServicePort `json:"ports,omitempty"`
}

func ParseV1Service(item corev1.Service) Service {
	service := Service{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Type:       string(item.Spec.Type),
		ClusterIp:  item.Spec.ClusterIP,
		ClusterIPs: item.Spec.ClusterIPs,
		IpFamilies: item.Spec.IPFamilies,
		Ports:      item.Spec.Ports,
	}
	if item.Spec.InternalTrafficPolicy != nil {
		service.InternalTrafficPolicy = string(*item.Spec.InternalTrafficPolicy)
	}
	return service
}
