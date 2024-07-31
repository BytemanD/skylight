package model

import "k8s.io/apimachinery/pkg/version"

type Cluster struct {
	Host           string        `json:"host,omitempty"`
	ApiPath        string        `json:"api_path,omitempty"`
	CurrentContext string        `json:"current_context,omitempty"`
	ServerVersion  *version.Info `json:"server_version,omitempty"`
}
