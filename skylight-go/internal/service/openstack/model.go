package openstack

type LoginInfo struct {
	Cluster  string  `json:"cluster,omitempty"`
	Region   string  `json:"region,omitempty"`
	Project  Project `json:"project,omitempty"`
	User     User    `json:"user,omitempty"`
	Roles    []Role  `json:"roles,omitempty"`
	Password string  `json:"password,omitempty"`
}
