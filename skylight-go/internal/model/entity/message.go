package entity

import (
	"fmt"
	"time"
)

type Message struct {
	// Type   string      `json:"type,omitempty"`
	Level       string      `json:"level,omitempty"`
	Topic       string      `json:"topic,omitempty"`
	Description string      `json:"description,omitempty"`
	Date        time.Time   `json:"date,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

func (m Message) String() string {
	return fmt.Sprintf("<'%s' level=%s description=%s>", m.Topic, m.Level, m.Description)
}

func newMessage(level, topic, description string, data interface{}) Message {
	return Message{
		Level:       level,
		Topic:       topic,
		Description: description,
		Date:        time.Now(),
		Data:        data}
}

func NewSuccessMessage(topic string, description string, data interface{}) Message {
	return newMessage("success", topic, description, data)
}

func NewInfoMessage(topic, description string, data interface{}) Message {
	return newMessage("info", topic, description, data)
}

func NewErrorMessage(topic, description string, data interface{}) Message {
	return newMessage("error", topic, description, data)
}

func NewWarningMessage(topic, description string, data interface{}) Message {
	return newMessage("warning", topic, description, data)
}
