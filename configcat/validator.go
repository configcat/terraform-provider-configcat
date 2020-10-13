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
