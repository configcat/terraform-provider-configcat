package configcat

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = regexValidator{}

type regexValidator struct{}

func (validator regexValidator) Description(_ context.Context) string {
	return "value must be a valid regular expression"
}

func (validator regexValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v regexValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if _, err := regexp.Compile(value); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// IsRegex returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Is a valid regular expression.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func IsRegex() validator.String {
	return regexValidator{}
}
