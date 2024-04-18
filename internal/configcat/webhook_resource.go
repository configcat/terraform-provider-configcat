package configcat

import (
	"context"
	"fmt"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var _ resource.Resource = &webhookResource{}
var _ resource.ResourceWithImportState = &webhookResource{}

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

type webhookResource struct {
	client *client.Client
}

type webhookHeaderResourceModel struct {
}

type webhookResourceModel struct {
	ConfigId      types.String `tfsdk:"config_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`

	ID                   types.Int64                  `tfsdk:"id"`
	Url                  types.String                 `tfsdk:"url"`
	HttpMethod           types.String                 `tfsdk:"http_method"`
	Content              types.String                 `tfsdk:"content"`
	WebhookHeaders       []webhookHeaderResourceModel `tfsdk:"webhook_headers"`
	SecureWebhookHeaders []webhookHeaderResourceModel `tfsdk:"secure_webhook_headers"`
}

func (r *webhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *webhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages a **" + WebhookResourceName + "**. [What is a " + WebhookResourceName + " in ConfigCat?](https://configcat.com/docs/advanced/notifications-webhooks/)",

		Attributes: map[string]schema.Attribute{
			ID: schema.Int64Attribute{
				Description: "The unique ID of the " + WebhookResourceName + ".",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			ConfigId: schema.StringAttribute{
				Description: "The ID of the " + ConfigResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			EnvironmentId: schema.StringAttribute{
				Description: "The ID of the " + EnvironmentResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			WebhookUrl: schema.StringAttribute{
				Description: "The URL of the " + WebhookResourceName + ".",
				Required:    true,
			},
			WebhookHttpMethod: schema.StringAttribute{
				Description: "The HTTP method. Available values: `get`|`post`. Default: `get`",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("get"),
			},
			WebhookContent: schema.StringAttribute{
				Description: "The HTTP body content.",
				Optional:    true,
			},
		},
	}
}

func (r *webhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	order := int32(plan.Order.ValueInt64())
	body := sw.CreateWebhookModel{
		Name:        plan.Name.ValueString(),
		Description: *sw.NewNullableString(plan.Description.ValueStringPointer()),
		Color:       *sw.NewNullableString(plan.Color.ValueStringPointer()),
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.CreateWebhook(plan.ProductId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+WebhookResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetWebhook(state.ID.ValueString())
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+WebhookResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan webhookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.Equal(state.Name) && plan.Description.Equal(state.Description) && plan.Color.Equal(state.Color) && plan.Order.Equal(state.Order) {
		return
	}

	order := int32(plan.Order.ValueInt64())
	body := sw.UpdateWebhookModel{
		Name:        *sw.NewNullableString(plan.Name.ValueStringPointer()),
		Description: *sw.NewNullableString(plan.Description.ValueStringPointer()),
		Color:       *sw.NewNullableString(plan.Color.ValueStringPointer()),
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.UpdateWebhook(plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+WebhookResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWebhook(state.ID.ValueString())

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, WebhookResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+WebhookResourceName+", got error: %s", err))
		return
	}
}

func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *webhookResourceModel) UpdateFromApiModel(model sw.WebhookModel) {
	modelOrder := int64(*model.Order)
	resourceModel.ID = types.StringPointerValue(model.WebhookId)
	resourceModel.ProductId = types.StringPointerValue(model.Product.ProductId)
	resourceModel.Name = types.StringPointerValue(model.Name.Get())
	resourceModel.Description = types.StringPointerValue(model.Description.Get())
	resourceModel.Color = types.StringPointerValue(model.Color.Get())
	resourceModel.Order = types.Int64Value(modelOrder)
}
