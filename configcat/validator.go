package configcat

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

func validateGUIDFunc(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if _, err := uuid.Parse(v); err != nil {
		errs = append(errs, fmt.Errorf("%q: invalid GUID", key))
	}
	return
}

func validateRegexFunc(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v == "" {
		return
	}

	if _, err := regexp.Compile(v); err != nil {
		errs = append(errs, fmt.Errorf("%q: invalid regex", key))
	}
	return
}

func validateColorFunc(val interface{}, key string) (warns []string, errs []error) {
	color := val.(string)
	if color != "" && color != "panther" && color != "whale" && color != "salmon" && color != "lizard" && color != "canary" && color != "koala" {
		errs = append(errs, fmt.Errorf("%q: '%s' color is invalid. Valid values: panther, whale, salmon, lizard, canary, koala", key, color))
	}
	return
}
