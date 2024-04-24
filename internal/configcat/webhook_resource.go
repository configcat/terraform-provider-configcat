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
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
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

func createWebhookHeaderSchema(isSecure bool, description string) *schema.ListNestedAttribute {
	return &schema.ListNestedAttribute{
		Optional:    true,
		Description: description,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				WebhookHeaderKey: schema.StringAttribute{
					Required:    true,
					Description: "The HTTP header key.",
				},
				WebhookHeaderValue: schema.StringAttribute{
					Required:    true,
					Description: "The HTTP header value.",
					Sensitive:   isSecure,
				},
			},
		},
	}
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
			WebhookHeaders:       createWebhookHeaderSchema(false, "List of plain text HTTP headers. The value of a plain text header is always visible for everyone. It also appears in audit logs and on the webhook test UI."),
			SecureWebhookHeaders: createWebhookHeaderSchema(true, "List of secret HTTP headers. The value of a secret header is write-only, nobody will see it after saving the webhook. It won't appear in audit logs and on the webhook test UI either."),
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

	httpMethod, err := sw.NewWebHookHttpMethodFromValue(plan.HttpMethod.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Invalid "+WebhookHttpMethod+". Error: %s", err))
		return
	}

	webhookHeaders := make([]sw.WebhookHeaderModel, len(plan.WebhookHeaders)+len(plan.SecureWebhookHeaders))
	for webhookHeaderIndex, webhookHeader := range plan.WebhookHeaders {
		isSecure := false
		webhookHeaders[webhookHeaderIndex] = sw.WebhookHeaderModel{
			Key:      webhookHeader.Key.ValueString(),
			Value:    webhookHeader.Value.ValueString(),
			IsSecure: &isSecure,
		}
	}
	for webhookHeaderIndex, webhookHeader := range plan.SecureWebhookHeaders {
		isSecure := true
		webhookHeaders[len(plan.WebhookHeaders)+webhookHeaderIndex] = sw.WebhookHeaderModel{
			Key:      webhookHeader.Key.ValueString(),
			Value:    webhookHeader.Value.ValueString(),
			IsSecure: &isSecure,
		}
	}

	body := sw.WebHookRequest{
		Url:            plan.Url.ValueString(),
		HttpMethod:     httpMethod,
		Content:        *sw.NewNullableString(plan.Content.ValueStringPointer()),
		WebHookHeaders: webhookHeaders,
	}

	model, err := r.client.CreateWebhook(plan.ConfigId.ValueString(), plan.EnvironmentId.ValueString(), body)
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

	model, err := r.client.GetWebhook(int32(state.ID.ValueInt64()))
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

	if !webhookHasChanges(&plan, &state) {
		return
	}

	httpMethod, err := sw.NewWebHookHttpMethodFromValue(plan.HttpMethod.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Invalid "+WebhookHttpMethod+". Error: %s", err))
		return
	}

	webhookHeaders := make([]sw.WebhookHeaderModel, len(plan.WebhookHeaders)+len(plan.SecureWebhookHeaders))
	for webhookHeaderIndex, webhookHeader := range plan.WebhookHeaders {
		isSecure := false
		webhookHeaders[webhookHeaderIndex] = sw.WebhookHeaderModel{
			Key:      webhookHeader.Key.ValueString(),
			Value:    webhookHeader.Value.ValueString(),
			IsSecure: &isSecure,
		}
	}
	for webhookHeaderIndex, webhookHeader := range plan.SecureWebhookHeaders {
		isSecure := true
		webhookHeaders[len(plan.WebhookHeaders)+webhookHeaderIndex] = sw.WebhookHeaderModel{
			Key:      webhookHeader.Key.ValueString(),
			Value:    webhookHeader.Value.ValueString(),
			IsSecure: &isSecure,
		}
	}

	body := sw.WebHookRequest{
		Url:            plan.Url.ValueString(),
		HttpMethod:     httpMethod,
		Content:        *sw.NewNullableString(plan.Content.ValueStringPointer()),
		WebHookHeaders: webhookHeaders,
	}

	model, err := r.client.UpdateWebhook(int32(state.ID.ValueInt64()), body)
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

	err := r.client.DeleteWebhook(int32(state.ID.ValueInt64()))

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

	resourceModel.ID = types.Int64Value(int64(*model.WebhookId))

	resourceModel.ConfigId = types.StringPointerValue(model.Config.ConfigId)
	resourceModel.EnvironmentId = types.StringPointerValue(model.Environment.EnvironmentId)
	resourceModel.Url = types.StringPointerValue(model.Url.Get())
	resourceModel.HttpMethod = types.StringPointerValue((*string)(model.HttpMethod))
	resourceModel.Content = types.StringPointerValue(model.Content.Get())

	resourceModel.WebhookHeaders = make([]webhookHeaderResourceModel, 0)
	resourceModel.SecureWebhookHeaders = make([]webhookHeaderResourceModel, 0)
	for _, webhookHeader := range model.WebHookHeaders {
		webhookHeaderModel := &webhookHeaderResourceModel{
			Key:   types.StringValue(webhookHeader.Key),
			Value: types.StringValue(webhookHeader.Value),
		}
		if *webhookHeader.IsSecure {
			resourceModel.SecureWebhookHeaders = append(resourceModel.WebhookHeaders, *webhookHeaderModel)
		} else {
			resourceModel.WebhookHeaders = append(resourceModel.WebhookHeaders, *webhookHeaderModel)
		}
	}
}

func webhookHasChanges(plan *webhookResourceModel, state *webhookResourceModel) bool {

	if !plan.Url.Equal(state.Url) ||
		!plan.HttpMethod.Equal(state.HttpMethod) ||
		!plan.Content.Equal(state.Content) ||
		len(plan.WebhookHeaders) != len(state.WebhookHeaders) ||
		len(plan.SecureWebhookHeaders) != len(state.SecureWebhookHeaders) {
		return true
	}
	for webhookHeaderIndex, planWebhookHeader := range plan.WebhookHeaders {
		if webhookHeaderHasChanges(&planWebhookHeader, &state.WebhookHeaders[webhookHeaderIndex]) {
			return true
		}
	}

	for webhookHeaderIndex, planWebhookHeader := range plan.SecureWebhookHeaders {
		if webhookHeaderHasChanges(&planWebhookHeader, &state.SecureWebhookHeaders[webhookHeaderIndex]) {
			return true
		}
	}

	return false
}

func webhookHeaderHasChanges(plan *webhookHeaderResourceModel, state *webhookHeaderResourceModel) bool {
	if !plan.Key.Equal(state.Key) ||
		!plan.Value.Equal(state.Value) {
		return true
	}

	return false
}
