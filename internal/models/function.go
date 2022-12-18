package models

import (
	"github.com/faasf/functions-api/internal/models/enums"
)

type Function struct {
	Name                       string               `json:"name,omitempty"`
	Description                string               `json:"description,omitempty"`
	SourceCode                 SourceCode           `json:"sourceCode,omitempty"`
	Runtime                    enums.Runtime        `json:"runtime,omitempty"`
	Timeout                    int                  `json:"timeout,omitempty"`
	Status                     enums.FunctionStatus `json:"status,omitempty"`
	EnvironmentVariables       map[string]string    `json:"environmentVariables,omitempty"`
	SecretEnvironmentVariables []Secret             `json:"secretEnvironmentVariables,omitempty"`
	HttpTriggers               []HttpTrigger        `json:"httpTriggers,omitempty"`
	EventTriggers              []EventTrigger       `json:"eventTriggers,omitempty"`
	ETag                       string               `json:"etag"`
}

type SourceCode struct {
	Type     enums.SourceCodeType `json:"type,omitempty"`
	Language enums.Language       `json:"language,omitempty"`
	Content  string               `json:"content,omitempty"`
}
