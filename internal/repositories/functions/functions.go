package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/faasf/functions-api/internal/commands"
	"github.com/faasf/functions-api/internal/config"
	"github.com/faasf/functions-api/internal/models"
	"github.com/faasf/functions-api/internal/models/enums"
	"github.com/faasf/functions-api/internal/utils"
	"strconv"
)

type FunctionsDaprRepo struct {
	storeName string
}

func New(c *config.Config) *FunctionsDaprRepo {
	return &FunctionsDaprRepo{
		storeName: c.Dapr.StoreName,
	}
}

func (r *FunctionsDaprRepo) GetAll(ctx context.Context) ([]models.Function, error) {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	context := context.Background()

	result, err := client.QueryStateAlpha1(context, r.storeName, "{}", nil)
	if err != nil {
		panic(err)
	}

	fns := []models.Function{}
	for _, a := range result.Results {
		data := models.Function{}
		json.Unmarshal(a.Value, &data)
		data.Name = a.Key
		fns = append(fns, data)
	}

	return fns, nil
}

func (r *FunctionsDaprRepo) GetByName(ctx context.Context, n string) (*models.Function, error) {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	context := context.Background()
	result, err := client.GetState(context, r.storeName, n, nil)
	if err != nil {
		panic(err)
	}

	if result.Etag == "" {
		return nil, nil
	}

	fn := models.Function{}

	json.Unmarshal(result.Value, &fn)
	fn.Name = result.Key

	return &fn, nil
}

func (r *FunctionsDaprRepo) Create(ctx context.Context, cmd commands.CreateFunctionCommand) (*models.Function, error) {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	context := context.Background()

	fn := createEntityFromCommand(cmd)
	fn.ETag = "1"
	data, err := json.Marshal(fn)
	if err != nil {
		panic(err)
	}
	err = client.SaveState(context, r.storeName, cmd.Name, data, nil)
	if err != nil {
		panic(err)
	}

	return fn, nil
}

func (r *FunctionsDaprRepo) Update(ctx context.Context, n string, cmd commands.UpdateFunctionCommand) (*models.Function, error) {
	fn, _ := r.GetByName(ctx, n)

	if fn.ETag != cmd.ETag {
		fmt.Println(fn.ETag + " - stale - " + cmd.ETag)
		return nil, errors.New("stale data")
	}

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	context := context.Background()

	updatedFn := createEntityFromUpdateCommand(n, cmd)

	if fn.Status == enums.Published {
		updatedFn.Status = enums.Published
	}

	etag, _ := strconv.ParseInt(fn.ETag, 10, 32)
	updatedFn.ETag = strconv.Itoa(int(etag + 1))
	data, err := json.Marshal(updatedFn)
	if err != nil {
		panic(err)
	}
	err = client.SaveState(context, r.storeName, n, data, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return updatedFn, nil
}

func (r *FunctionsDaprRepo) Publish(ctx context.Context, n string, cmd commands.PublishFunctionCommand) (*models.Function, error) {
	fn, _ := r.GetByName(ctx, n)

	if fn.ETag != cmd.ETag {
		return nil, errors.New("stale data")
	}

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	context := context.Background()

	fn.Status = enums.Published
	etag, _ := strconv.ParseInt(fn.ETag, 10, 32)
	fn.ETag = strconv.Itoa(int(etag + 1))
	data, err := json.Marshal(fn)
	if err != nil {
		panic(err)
	}
	err = client.SaveState(context, r.storeName, n, data, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return fn, nil
}

func createEntityFromCommand(cmd commands.CreateFunctionCommand) *models.Function {
	lang, err := utils.ResolveLanguageFromName(cmd.Name)
	if err != nil {
		panic(err)
	}

	return &models.Function{
		Name:          cmd.Name,
		Description:   cmd.Description,
		Runtime:       cmd.Runtime,
		HttpTriggers:  cmd.HttpTriggers,
		EventTriggers: cmd.EventTriggers,
		SourceCode: models.SourceCode{
			Type:     enums.File,
			Language: lang,
			Content:  cmd.Code,
		},
		Timeout: cmd.Timeout,
		Status:  enums.Draft,
	}
}

func createEntityFromUpdateCommand(n string, cmd commands.UpdateFunctionCommand) *models.Function {
	lang, err := utils.ResolveLanguageFromName(n)
	if err != nil {
		panic(err)
	}

	return &models.Function{
		Name:          n,
		Description:   cmd.Description,
		Runtime:       cmd.Runtime,
		HttpTriggers:  cmd.HttpTriggers,
		EventTriggers: cmd.EventTriggers,
		SourceCode: models.SourceCode{
			Type:     enums.File,
			Language: lang,
			Content:  cmd.Code,
		},
		Timeout: cmd.Timeout,
		Status:  cmd.Status,
	}
}
