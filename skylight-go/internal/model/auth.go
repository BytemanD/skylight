package model

type AuthInfo struct {
	Cluster  string `json:"cluster"`
	Project  string `json:"project"`
	User     string `json:"user"`
	Password string `json:"password"`
}
