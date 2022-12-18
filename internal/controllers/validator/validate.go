package validator

import (
	"fmt"
	"github.com/faasf/functions-api/internal/services"
	"github.com/gin-gonic/gin"
	goValidator "gopkg.in/validator.v2"
	"net/http"
)

type RequestValidator struct {
	functionsService services.FunctionsService
	v                *goValidator.Validator
}

func New(f services.FunctionsService) *RequestValidator {

	validator := goValidator.NewValidator()

	v := &RequestValidator{
		functionsService: f,
		v:                validator,
	}

	err := validator.SetValidationFunc("uniqueFunctionName", v.uniqueFunctionName())
	if err != nil {
		panic(err)
	}

	return v
}

func (rv *RequestValidator) ValidateBody(c *gin.Context, data interface{}) bool {
	if err := rv.v.Validate(data); err != nil {
		var validationErrors []ValidationError
		errs := err.(goValidator.ErrorMap)
		for f, e := range errs {
			validationErrors = append(validationErrors, createValidationError(f, e)...)
		}

		c.JSON(http.StatusBadRequest, validationErrors)
		return false
	}

	return true
}

func createValidationError(field string, errors goValidator.ErrorArray) []ValidationError {
	var errs []ValidationError
	for _, b := range errors {
		errs = append(errs, ValidationError{
			Error:       b.Error(),
			Field:       field,
			Description: resolveMessage(field, b),
		})
	}
	return errs
}

func resolveMessage(f string, e error) string {
	if e == goValidator.ErrZeroValue {
		return fmt.Sprintf("%s cannot be empty.", f)
	}

	if e == ErrUniqueFunctionName {
		return fmt.Sprintf("Function with that name already exists.")
	}

	panic(fmt.Sprintf("Validator not supported: %s", e))
}

type ValidationError struct {
	Field       string `json:"field"`
	Error       string `json:"error"`
	Description string `json:"description"`
}
