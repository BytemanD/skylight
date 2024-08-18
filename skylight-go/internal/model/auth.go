package model

type AuthInfo struct {
	Cluster  string `json:"cluster,omitempty"`
	Region   string `json:"region,omitempty"`
	Project  string `json:"project,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
