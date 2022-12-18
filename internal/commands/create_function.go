package commands

import (
	"github.com/faasf/functions-api/internal/models"
	"github.com/faasf/functions-api/internal/models/enums"
)

type CreateFunctionCommand struct {
	Name          string                `json:"name" validate:"nonzero,uniqueFunctionName"`
	Description   string                `json:"description"`
	Runtime       enums.Runtime         `json:"runtime" validate:"nonzero"`
	HttpTriggers  []models.HttpTrigger  `json:"httpTriggers"`
	EventTriggers []models.EventTrigger `json:"eventTriggers"`
	Code          string                `json:"code"`
	Timeout       int                   `json:"timeout"`
}
