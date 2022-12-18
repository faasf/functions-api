package services

import (
	"context"

	"github.com/faasf/functions-api/internal/commands"
	"github.com/faasf/functions-api/internal/models"
)

type (
	FunctionsService interface {
		GetAll(context.Context) ([]models.Function, error)
		GetByName(ctx context.Context, n string) (*models.Function, error)
		Create(context.Context, commands.CreateFunctionCommand) (*models.Function, error)
		Update(context.Context, string, commands.UpdateFunctionCommand) (*models.Function, error)
		Publish(context.Context, string, commands.PublishFunctionCommand) (*models.Function, error)
	}
)
