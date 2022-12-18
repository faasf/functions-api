package functions

import (
	"context"
	"fmt"
	"github.com/faasf/functions-api/internal/models/enums"
	"github.com/faasf/functions-api/internal/utils"

	"github.com/faasf/functions-api/internal/commands"
	"github.com/faasf/functions-api/internal/models"
	"github.com/faasf/functions-api/internal/repositories"
)

type FunctionsServiceImpl struct {
	repo repositories.FunctionsRepo
}

func New(r repositories.FunctionsRepo) *FunctionsServiceImpl {
	return &FunctionsServiceImpl{
		repo: r,
	}
}

func (s *FunctionsServiceImpl) GetAll(ctx context.Context) ([]models.Function, error) {
	functions, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("FunctionsService - GetAll - s.repo.GetAll: %w", err)
	}

	return functions, nil
}

func (s *FunctionsServiceImpl) GetByName(ctx context.Context, n string) (*models.Function, error) {
	fn, err := s.repo.GetByName(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("FunctionsService - GetAll - s.repo.GetAll: %w", err)
	}

	return fn, nil
}

func (s *FunctionsServiceImpl) Create(ctx context.Context, cmd commands.CreateFunctionCommand) (*models.Function, error) {
	if cmd.Code == "" {
		setDefaultCodeBasedOnLanguage(&cmd)
	}
	createdFn, err := s.repo.Create(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("FunctionsService - Create - s.repo.Create: %w", err)
	}

	return createdFn, nil
}

func (s *FunctionsServiceImpl) Update(ctx context.Context, n string, cmd commands.UpdateFunctionCommand) (*models.Function, error) {
	updatedFn, err := s.repo.Update(ctx, n, cmd)
	if err != nil {
		return nil, fmt.Errorf("FunctionsService - Update - s.repo.Update: %w", err)
	}

	return updatedFn, nil
}

func (s *FunctionsServiceImpl) Publish(ctx context.Context, n string, cmd commands.PublishFunctionCommand) (*models.Function, error) {
	updatedFn, err := s.repo.Publish(ctx, n, cmd)
	if err != nil {
		return nil, fmt.Errorf("FunctionsService - Publish - s.repo.Publish: %w", err)
	}

	return updatedFn, nil
}

func setDefaultCodeBasedOnLanguage(cmd *commands.CreateFunctionCommand) {
	lang, err := utils.ResolveLanguageFromName(cmd.Name)
	if err != nil {
		panic(err)
	}

	if lang == enums.Typescript {
		cmd.Code = "import { Logging, HttpRequest, HttpResponse } from '@faasff/nodejs-common';\n\nexport default (req: HttpRequest): Promise<HttpResponse> | HttpResponse => {\n    const logger = Logging.getLogger();\n\n    logger.info({ message: 'Hello world' });\n    \n    return {\n        statusCode: 200,\n        body: 'OK'\n    };\n}"
	} else {
		cmd.Code = "const { Logging } = require('@faasff/nodejs-common');\n\nmodule.exports = async () => {\n    const logger = Logging.getLogger();\n\n    logger.info({ message: 'Hello world' });\n    \n    return {\n        statusCode: 200,\n        body: 'OK'\n    };\n}"
	}
}
