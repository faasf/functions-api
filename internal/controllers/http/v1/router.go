package v1

import (
	"encoding/json"
	"github.com/faasf/functions-api/internal/config"
	"github.com/faasf/functions-api/internal/models"
	"github.com/faasf/functions-api/internal/models/enums"
	"github.com/faasf/functions-api/internal/services"
	"github.com/faasf/functions-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"net/http"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewRouter(handler *gin.Engine, l logger.Interface, cfg *config.Config, s services.FunctionsService) {
	handler.Use(CORS())
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	// handler.GET("/swagger/*any", swaggerHandler)

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	h := handler.Group("/v1")
	{
		newFunctionRoutes(h, s, l)
	}

	handler.NoRoute(registerHttpHandlers(cfg, s))
}

func registerHttpHandlers(cfg *config.Config, s services.FunctionsService) func(c *gin.Context) {
	return func(c *gin.Context) {
		fns, _ := s.GetAll(c)
		m := map[string]models.Function{}

		router := mux.NewRouter()
		for _, fn := range fns {
			if fn.Status == enums.Draft || fn.HttpTriggers == nil || len(fn.HttpTriggers) == 0 {
				continue
			}
			for _, tr := range fn.HttpTriggers {
				router.NewRoute().Path(tr.Url).Methods(enums.HttpMethodToString(tr.Method))
				m[tr.Url+"-"+enums.HttpMethodToString(tr.Method)] = fn
			}

		}

		routeMatch := mux.RouteMatch{}
		match := router.Match(c.Request, &routeMatch)
		if !match {
			c.Status(http.StatusNotFound)
			return
		}

		path, err := routeMatch.Route.GetPathTemplate()
		if err != nil {
			panic(err)
		}

		calledFn := m[path+"-"+c.Request.Method]

		fnData := toFunctionData(calledFn, routeMatch.Vars)
		fnDataJson, err := json.Marshal(fnData)
		if err != nil {
			panic(err)
		}

		c.Request.Header.Add("x-function-data", string(fnDataJson))
		Proxy(cfg, c)
	}
}

type functionData struct {
	Name                       string            `json:"name,omitempty"`
	Timeout                    int               `json:"timeout,omitempty"`
	EnvironmentVariables       map[string]string `json:"environmentVariables,omitempty"`
	SecretEnvironmentVariables []models.Secret   `json:"secretEnvironmentVariables,omitempty"`
	ETag                       string            `json:"etag"`
	Params                     map[string]string `json:"params"`
}

func toFunctionData(fn models.Function, params map[string]string) functionData {
	return functionData{
		ETag:                       fn.ETag,
		Name:                       fn.Name,
		Timeout:                    fn.Timeout,
		EnvironmentVariables:       fn.EnvironmentVariables,
		SecretEnvironmentVariables: fn.SecretEnvironmentVariables,
		Params:                     params,
	}
}
