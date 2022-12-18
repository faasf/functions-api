package v1

import (
	"net/http"

	"github.com/faasf/functions-api/internal/commands"
	"github.com/faasf/functions-api/internal/controllers/validator"
	"github.com/faasf/functions-api/internal/services"
	"github.com/faasf/functions-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

type functionRoutes struct {
	s services.FunctionsService
	l logger.Interface
	v *validator.RequestValidator
}

func newFunctionRoutes(handler *gin.RouterGroup, s services.FunctionsService, l logger.Interface) {
	r := &functionRoutes{s: s, l: l, v: validator.New(s)}

	h := handler.Group("/functions")
	{
		h.GET("", r.getAllFunctions)
		h.GET(":name", r.getFunctionByName)
		h.POST("", r.createFunction)
		h.PUT(":name", r.updateFunction)
		h.POST(":name/publish", r.publishFunction)
	}
}

func (r *functionRoutes) getAllFunctions(c *gin.Context) {
	functions, err := r.s.GetAll(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllFunctions")
		errorResponse(c, http.StatusInternalServerError, "dapr problems")

		return
	}

	c.JSON(http.StatusOK, functions)
}

func (r *functionRoutes) createFunction(c *gin.Context) {
	var request commands.CreateFunctionCommand
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - createFunction")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ok := r.v.ValidateBody(c, request)
	if !ok {
		return
	}

	var createdFn, err = r.s.Create(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err, "http - v1 - createFunction")
		errorResponse(c, http.StatusInternalServerError, "dapr problems")

		return
	}

	c.JSON(http.StatusCreated, createdFn)
}

func (r *functionRoutes) getFunctionByName(c *gin.Context) {
	functionName := c.Param("name")

	var fn, err = r.s.GetByName(c.Request.Context(), functionName)
	if err != nil {
		r.l.Error(err, "http - v1 - getFunctionByName")
		errorResponse(c, http.StatusInternalServerError, "dapr problems")

		return
	}

	c.JSON(http.StatusOK, fn)
}

func (r *functionRoutes) updateFunction(c *gin.Context) {
	functionName := c.Param("name")

	var request commands.UpdateFunctionCommand
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - updateFunction")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ok := r.v.ValidateBody(c, request)
	if !ok {
		return
	}

	var updatedFn, err = r.s.Update(c.Request.Context(), functionName, request)
	if err != nil {
		r.l.Error(err, "http - v1 - updateFunction")
		errorResponse(c, http.StatusInternalServerError, "dapr problems")

		return
	}

	c.JSON(http.StatusOK, updatedFn)
}

func (r *functionRoutes) publishFunction(c *gin.Context) {
	functionName := c.Param("name")

	var request commands.PublishFunctionCommand
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - publishFunction")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ok := r.v.ValidateBody(c, request)
	if !ok {
		return
	}

	var updatedFn, err = r.s.Publish(c.Request.Context(), functionName, request)
	if err != nil {
		r.l.Error(err, "http - v1 - publishFunction")
		errorResponse(c, http.StatusInternalServerError, "dapr problems")

		return
	}

	c.JSON(http.StatusOK, updatedFn)
}
