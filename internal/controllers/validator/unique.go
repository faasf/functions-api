package validator

import (
	"context"
	"errors"
	"gopkg.in/validator.v2"
)

var (
	ErrUniqueFunctionName = validator.TextErr{Err: errors.New("unique function name")}
)

func (rv *RequestValidator) uniqueFunctionName() validator.ValidationFunc {
	return func(v interface{}, param string) error {
		n, ok := v.(string)
		if !ok {
			return ErrUniqueFunctionName
		}

		fn, err := rv.functionsService.GetByName(context.Background(), n)
		if err != nil {
			panic(err)
		}

		if fn != nil {
			return ErrUniqueFunctionName
		}

		return nil
	}
}
