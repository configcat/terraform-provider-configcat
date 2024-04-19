package configcat

import (
	"context"
	"fmt"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var _ resource.Resource = &productPreferencesResource{}
var _ resource.ResourceWithImportState = &productPreferencesResource{}

func NewProductPreferencesResource() resource.Resource {
	return &productPreferencesResource{}
}

type productPreferencesResource struct {
	client *client.Client
}

type productPreferencesResourceModel struct {
	ProductId types.String `tfsdk:"product_id"`

	ID types.String `tfsdk:"id"`

	KeyGenerationMode          types.String   `tfsdk:"key_generation_mode"`
	MandatorySettingHint       types.Bool     `tfsdk:"mandatory_setting_hint"`
	ShowVariationId            types.Bool     `tfsdk:"show_variation_id"`
	ReasonRequired             types.Bool     `tfsdk:"reason_required"`
	ReasonRequiredEnvironments []types.String `tfsdk:"reason_required_environments"`
}

func (r *productPreferencesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_preferences"
}

func (r *productPreferencesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages a **" + ProductPreferencesResourceName + "**.",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the resource. Do not use.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			ProductId: schema.StringAttribute{
				Description: "The ID of the " + ProductResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			ProductPreferenceMandatorySettingHint: schema.BoolAttribute{
				Description: "Indicates whether Feature flags and Settings must have a hint. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			ProductPreferenceShowVariationId: schema.BoolAttribute{
				Description: "Indicates whether a variation ID's must be shown on the ConfigCat Dashboard. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			ProductPreferenceKeyGenerationMode: schema.StringAttribute{
				Description: "Determines the Feature Flag key generation mode. Available values: `camelCase`|`upperCase`|`lowerCase`|`pascalCase`|`kebabCase`. Default: `camelCase`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("camelCase"),
			},
			ProductPreferenceReasonRequired: schema.BoolAttribute{
				Description: "Indicates that a mandatory note is required for saving and publishing. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			ProductPreferenceReasonRequiredEnvironmentments: schema.ListAttribute{
				Description: "List of Environments where mandatory note must be set before saving and publishing.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(IsGuid()),
				},
			},
		},
	}
}

func (r *productPreferencesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *productPreferencesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.createOrUpdateProductPreferences(ctx, &req.Plan, nil, &resp.State, &resp.Diagnostics)
}

func (r *productPreferencesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.createOrUpdateProductPreferences(ctx, &req.Plan, &req.State, &resp.State, &resp.Diagnostics)
}

func (r *productPreferencesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state productPreferencesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetProductPreferences(state.ProductId.ValueString())
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+ProductPreferencesResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model, state.ProductId.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *productPreferencesResource) createOrUpdateProductPreferences(ctx context.Context, requestPlan *tfsdk.Plan, requestState *tfsdk.State, responseState *tfsdk.State, diag *diag.Diagnostics) {
	var plan productPreferencesResourceModel

	diag.Append(requestPlan.Get(ctx, &plan)...)

	if diag.HasError() {
		return
	}

	if requestState != nil {
		var state productPreferencesResourceModel
		diag.Append(requestState.Get(ctx, &state)...)
		if !hasProductPreferenceChanges(&plan, &state) {
			return
		}
	}

	keyGenerationMode, err := sw.NewKeyGenerationModeFromValue(plan.KeyGenerationMode.ValueString())
	if err != nil {
		diag.AddError("Unable to Create Resource", fmt.Sprintf("Invalid "+ProductPreferenceKeyGenerationMode+". Error: %s", err))
		return
	}

	body := sw.UpdatePreferencesRequest{
		KeyGenerationMode:    keyGenerationMode,
		ShowVariationId:      *sw.NewNullableBool(plan.ShowVariationId.ValueBoolPointer()),
		MandatorySettingHint: *sw.NewNullableBool(plan.MandatorySettingHint.ValueBoolPointer()),
		ReasonRequired:       *sw.NewNullableBool(plan.ReasonRequired.ValueBoolPointer()),
	}

	model, err := r.client.UpdateProductPreferences(plan.ProductId.ValueString(), body)
	if err != nil {
		diag.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+ProductPreferencesResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model, plan.ProductId.ValueString())

	diag.Append(responseState.Set(ctx, &plan)...)
}

func (r *productPreferencesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete operation should not do anything with the Product preferences.
}

func (r *productPreferencesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *productPreferencesResourceModel) UpdateFromApiModel(model sw.PreferencesModel, productId string) {
	resourceModel.ID = types.StringValue(productId)
	resourceModel.ProductId = types.StringValue(productId)
	resourceModel.MandatorySettingHint = types.BoolPointerValue(model.MandatorySettingHint)
	resourceModel.ShowVariationId = types.BoolPointerValue(model.ShowVariationId)
	resourceModel.KeyGenerationMode = types.StringPointerValue((*string)(model.KeyGenerationMode))
	resourceModel.ReasonRequired = types.BoolPointerValue(model.ReasonRequired)
	reasonRequiredEnvironments := make([]basetypes.StringValue, 0)
	for _, environment := range model.ReasonRequiredEnvironments {
		if environment.EnvironmentId != nil && environment.ReasonRequired != nil && *environment.ReasonRequired {
			resourceModel.ReasonRequiredEnvironments = append(resourceModel.ReasonRequiredEnvironments, types.StringValue(*environment.EnvironmentId))
		}
	}

	if len(reasonRequiredEnvironments) > 0 {
		resourceModel.ReasonRequiredEnvironments = reasonRequiredEnvironments
	} else {
		resourceModel.ReasonRequiredEnvironments = nil
	}
}

func hasProductPreferenceChanges(plan *productPreferencesResourceModel, state *productPreferencesResourceModel) bool {
	if !plan.KeyGenerationMode.Equal(state.KeyGenerationMode) ||
		!plan.ShowVariationId.Equal(state.ShowVariationId) ||
		!plan.MandatorySettingHint.Equal(state.MandatorySettingHint) ||
		!plan.ReasonRequired.Equal(state.ReasonRequired) ||
		len(plan.ReasonRequiredEnvironments) != len(state.ReasonRequiredEnvironments) {
		return true
	}

	for environmentIndex, environment := range plan.ReasonRequiredEnvironments {
		if !environment.Equal(state.ReasonRequiredEnvironments[environmentIndex]) {
			return true
		}
	}

	return false
}
