package configcat

import (
	"fmt"

	"github.com/google/uuid"
)

func validateGUIDFunc(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if _, err := uuid.Parse(v); err != nil {
		errs = append(errs, fmt.Errorf("%q: invalid GUID", key))
	}
	return
}
