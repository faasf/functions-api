package models

import (
	"github.com/faasf/functions-api/internal/models/enums"
)

type Trigger struct {
	Type enums.TriggerType `json:"type,omitempty"`
}

type HttpTrigger struct {
	Trigger
	Url    string           `json:"url,omitempty"`
	Method enums.HttpMethod `json:"method,omitempty"`
}

type EventTrigger struct {
	Trigger
	Topic     string `json:"topic,omitempty"`
	EventType string `json:"eventType,omitempty"`
}
