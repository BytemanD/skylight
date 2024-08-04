package model

type AuthInfo struct {
	Cluster  string `json:"cluster,omitempty"`
	Project  string `json:"project,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
