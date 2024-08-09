package configcat

import (
	"context"
	"fmt"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
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

	sw "github.com/configcat/configcat-publicapi-go-client/v2"
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

	KeyGenerationMode          types.String `tfsdk:"key_generation_mode"`
	MandatorySettingHint       types.Bool   `tfsdk:"mandatory_setting_hint"`
	ShowVariationId            types.Bool   `tfsdk:"show_variation_id"`
	ReasonRequired             types.Bool   `tfsdk:"reason_required"`
	ReasonRequiredEnvironments types.Map    `tfsdk:"reason_required_environments"`
}

func (r *productPreferencesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_preferences"
}

func (r *productPreferencesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the **" + ProductPreferencesResourceName + "**.",

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
				Description: "Indicates whether variation IDs must be shown on the ConfigCat Dashboard. Default: false.",
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
			ProductPreferenceReasonRequiredEnvironmentments: schema.MapAttribute{
				Description: "The environment specific mandatory note map block. Keys are the Environment IDs and the values indicate that a mandatory note is required for saving and publishing.",
				Computed:    true,
				Optional:    true,
				ElementType: types.BoolType,
				Validators: []validator.Map{
					mapvalidator.KeysAre(IsGuid()),
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

	resp.Diagnostics.Append(state.UpdateFromApiModel(ctx, *model, state.ProductId.ValueString())...)
	if resp.Diagnostics.HasError() {
		return
	}

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
		diag.AddError("Unable to update Product Preferences", fmt.Sprintf("Invalid "+ProductPreferenceKeyGenerationMode+". Error: %s", err))
		return
	}

	if plan.ReasonRequired.ValueBool() && !plan.ReasonRequiredEnvironments.IsUnknown() && !plan.ReasonRequiredEnvironments.IsNull() {
		diag.AddError("Unable to update Product Preferences", "Please set "+ProductPreferenceReasonRequired+" to true to require mandatory notes globally or specify "+ProductPreferenceReasonRequiredEnvironmentments+" to require mandatory notes for specific environments but don't specify both of them together.")
	}

	var reasonRequiredEnvironmentsMap map[string]types.Bool
	if plan.ReasonRequiredEnvironments.IsUnknown() || plan.ReasonRequiredEnvironments.IsNull() {
		reasonRequiredEnvironmentsMap = make(map[string]types.Bool, 0)
	} else {
		reasonRequiredEnvironmentsMap = make(map[string]types.Bool, len(plan.ReasonRequiredEnvironments.Elements()))
		diag.Append(plan.ReasonRequiredEnvironments.ElementsAs(ctx, &reasonRequiredEnvironmentsMap, false)...)
		if diag.HasError() {
			return
		}
	}

	reasonRequiredEnvironments := make([]sw.UpdateReasonRequiredEnvironmentModel, 0)
	for environmentIdKey, reasonRequiredValue := range reasonRequiredEnvironmentsMap {
		environmentId := environmentIdKey
		reasonRequired := reasonRequiredValue
		reasonRequiredEnvironments = append(reasonRequiredEnvironments, sw.UpdateReasonRequiredEnvironmentModel{
			EnvironmentId:  &environmentId,
			ReasonRequired: reasonRequired.ValueBoolPointer(),
		})
	}

	body := sw.UpdatePreferencesRequest{
		KeyGenerationMode:          keyGenerationMode,
		ShowVariationId:            *sw.NewNullableBool(plan.ShowVariationId.ValueBoolPointer()),
		MandatorySettingHint:       *sw.NewNullableBool(plan.MandatorySettingHint.ValueBoolPointer()),
		ReasonRequired:             *sw.NewNullableBool(plan.ReasonRequired.ValueBoolPointer()),
		ReasonRequiredEnvironments: reasonRequiredEnvironments,
	}

	model, err := r.client.UpdateProductPreferences(plan.ProductId.ValueString(), body)
	if err != nil {
		diag.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+ProductPreferencesResourceName+", got error: %s", err))
		return
	}

	diag.Append(plan.UpdateFromApiModel(ctx, *model, plan.ProductId.ValueString())...)

	if diag.HasError() {
		return
	}

	diag.Append(responseState.Set(ctx, &plan)...)
}

func (r *productPreferencesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete operation should not do anything with the Product preferences.
}

func (r *productPreferencesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ProductId), req, resp)
}

func (resourceModel *productPreferencesResourceModel) UpdateFromApiModel(ctx context.Context, model sw.PreferencesModel, productId string) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceModel.ID = types.StringValue(productId)
	resourceModel.ProductId = types.StringValue(productId)
	resourceModel.MandatorySettingHint = types.BoolPointerValue(model.MandatorySettingHint)
	resourceModel.ShowVariationId = types.BoolPointerValue(model.ShowVariationId)
	resourceModel.KeyGenerationMode = types.StringPointerValue((*string)(model.KeyGenerationMode))
	resourceModel.ReasonRequired = types.BoolPointerValue(model.ReasonRequired)

	reasonRequiredEnvironments := make(map[string]bool, len(model.ReasonRequiredEnvironments))
	for _, environment := range model.ReasonRequiredEnvironments {
		reasonRequiredEnvironments[*environment.EnvironmentId] = *environment.ReasonRequired
	}

	reasonRequiredEnvironmentsMapValue, diags := types.MapValueFrom(ctx, types.BoolType, reasonRequiredEnvironments)
	if diags.HasError() {
		return diags
	}

	resourceModel.ReasonRequiredEnvironments = reasonRequiredEnvironmentsMapValue

	return diags
}

func hasProductPreferenceChanges(plan *productPreferencesResourceModel, state *productPreferencesResourceModel) bool {
	if !plan.KeyGenerationMode.Equal(state.KeyGenerationMode) ||
		!plan.ShowVariationId.Equal(state.ShowVariationId) ||
		!plan.MandatorySettingHint.Equal(state.MandatorySettingHint) ||
		!plan.ReasonRequired.Equal(state.ReasonRequired) ||
		!plan.ReasonRequiredEnvironments.Equal(state.ReasonRequiredEnvironments) {
		return true
	}

	return false
}
