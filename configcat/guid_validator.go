package configcat

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = guidValidator{}

type guidValidator struct{}

func (validator guidValidator) Description(_ context.Context) string {
	return "value must be a valid GUID"
}

func (validator guidValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v guidValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if _, err := uuid.Parse(value); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// IsGuid returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Is a valid GUID.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func IsGuid() validator.String {
	return guidValidator{}
}
