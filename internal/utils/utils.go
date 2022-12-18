package utils

import (
	"errors"
	"github.com/faasf/functions-api/internal/models/enums"
	"strings"
)

func ResolveLanguageFromName(name string) (enums.Language, error) {
	if strings.HasSuffix(name, ".js") {
		return enums.Javascript, nil
	} else if strings.HasSuffix(name, ".ts") {
		return enums.Typescript, nil
	} else {
		return "", errors.New("unsupported language")
	}
}
