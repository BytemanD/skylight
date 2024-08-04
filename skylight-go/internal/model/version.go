package model

type Version struct {
	Version    string `json:"version,omitempty"`
	MinVersion string `json:"min_version,omitempty"`
}
