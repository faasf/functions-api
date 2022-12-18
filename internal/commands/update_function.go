package commands

import (
	"github.com/faasf/functions-api/internal/models"
	"github.com/faasf/functions-api/internal/models/enums"
)

type UpdateFunctionCommand struct {
	Description   string                `json:"description"`
	Runtime       enums.Runtime         `json:"runtime"`
	HttpTriggers  []models.HttpTrigger  `json:"httpTriggers"`
	EventTriggers []models.EventTrigger `json:"eventTriggers"`
	Code          string                `json:"code"`
	Timeout       int                   `json:"timeout"`
	Status        enums.FunctionStatus  `json:"status"`
	ETag          string                `json:"etag"`
}
