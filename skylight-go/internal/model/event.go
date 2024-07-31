package model

import (
	corev1 "k8s.io/api/core/v1"
)

type Event struct {
	BaseModel
	Type           string                 `json:"type,omitempty"`
	EventTime      string                 `json:"event_time,omitempty"`
	InvolvedObject corev1.ObjectReference `json:"involved_object,omitempty"`
	Message        string                 `json:"message,omitempty"`
	Reason         string                 `json:"reason,omitempty"`
	Source         corev1.EventSource     `json:"source,omitempty"`
	Action         string                 `json:"action,omitempty"`
	Count          int32                  `json:"count,omitempty"`
}

func ParseV1Event(item corev1.Event) Event {
	return Event{
		BaseModel: BaseModel{
			Name:     item.Name,
			Creation: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		},
		Type:           item.Type,
		EventTime:      item.EventTime.Format("2006-01-02 15:04:05"),
		InvolvedObject: item.InvolvedObject,
		Message:        item.Message,
		Reason:         item.Reason,
		Source:         item.Source,
		Action:         item.Action,
		Count:          item.Count,
	}
}
