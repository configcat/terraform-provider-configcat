package configcat

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sw "github.com/configcat/configcat-publicapi-go-client/v2"
)

var _ resource.Resource = &settingTagResource{}
var _ resource.ResourceWithImportState = &settingTagResource{}

func NewSettingTagResource() resource.Resource {
	return &settingTagResource{}
}

type settingTagResource struct {
	client *client.Client
}

type settingTagResourceModel struct {
	ID        types.String `tfsdk:"id"`
	SettingId types.String `tfsdk:"setting_id"`
	TagId     types.String `tfsdk:"tag_id"`
}

func (r *settingTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_tag"
}

func (r *settingTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Adds/Removes **" + TagResourceName + "s** to/from **" + SettingsResourceName + "**.",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the resource. Do not use.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			SettingId: schema.StringAttribute{
				Description: "The ID of the " + SettingResourceName + ".",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			TagId: schema.StringAttribute{
				Description: "The ID of the " + TagResourceName + ".",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *settingTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingTagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagIdString := plan.TagId.ValueString()
	tagId, tagIdConvErr := strconv.ParseInt(tagIdString, 10, 32)
	if tagIdConvErr != nil {
		resp.Diagnostics.AddError("Could not parse Tag ID", tagIdConvErr.Error())
		return
	}

	settingId, settingIdConvErr := strconv.ParseInt(plan.SettingId.ValueString(), 10, 32)
	if settingIdConvErr != nil {
		resp.Diagnostics.AddError("Could not parse Tag ID", settingIdConvErr.Error())
		return
	}

	operations := []sw.JsonPatchOperation{{
		Op:    sw.OPERATIONTYPE_ADD,
		Path:  "/tags/-",
		Value: tagIdString,
	}}

	_, err := r.client.UpdateSetting(int32(settingId), operations)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to add **"+TagResourceName+"** to **"+SettingsResourceName+"**., got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(settingId, tagId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingTagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	settingID, tagID, convErr := resourceConfigCatSettingTagParseID(state.ID.ValueString())
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse ID", convErr.Error())
		return
	}

	model, err := r.client.GetSetting(int32(settingID))
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SettingResourceName+", got error: %s", err))
		return
	}

	found := false
	for _, tag := range model.Tags {
		if *tag.TagId == tagID {
			found = true
			break
		}
	}

	if !found {
		// If the tag is not applied to the setting, we should remove it from the state
		resp.State.RemoveResource(ctx)
		return
	}

	state.UpdateFromApiModel(settingID, tagID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Invalid operation", "Invalid operation")
}

func (r *settingTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state settingTagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagIdString := state.TagId.ValueString()
	tagID, tagIdConvErr := strconv.ParseInt(tagIdString, 10, 32)
	if tagIdConvErr != nil {
		resp.Diagnostics.AddError("Could not parse Tag ID", tagIdConvErr.Error())
		return
	}

	settingId, settingIdConvErr := strconv.ParseInt(state.SettingId.ValueString(), 10, 32)
	if settingIdConvErr != nil {
		resp.Diagnostics.AddError("Could not parse Tag ID", settingIdConvErr.Error())
		return
	}

	model, getErr := r.client.GetSetting(int32(settingId))
	if getErr != nil {
		if _, ok := getErr.(client.NotFoundError); ok {
			// If the resource is already deleted, we consider the tag to be removed.
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SettingResourceName+", got error: %s", getErr))
		return
	}

	index := -1
	for tagIndex, tag := range model.Tags {
		if *tag.TagId == tagID {
			index = tagIndex
			break
		}
	}

	if index == -1 {
		// If the resource is already deleted, we can safely return
		return
	}

	operations := []sw.JsonPatchOperation{{
		Op:   sw.OPERATIONTYPE_REMOVE,
		Path: fmt.Sprintf("/tags/%d", index),
	}}

	_, err := r.client.UpdateSetting(int32(settingId), operations)
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, TagResourceName+" is already removed from the "+SettingResourceName+" in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to remove **"+TagResourceName+"** from **"+SettingsResourceName+"**., got error: %s", err))
		return
	}
}

func (r *settingTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *settingTagResourceModel) UpdateFromApiModel(settingId int64, tagId int64) {
	resourceModel.ID = types.StringValue(fmt.Sprintf("%d:%d", settingId, tagId))
	resourceModel.SettingId = types.StringValue(strconv.FormatInt(settingId, 10))
	resourceModel.TagId = types.StringValue(strconv.FormatInt(tagId, 10))
}

func resourceConfigCatSettingTagParseID(id string) (int64, int64, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected settingID:tagID", id)
	}

	settingID, sConvErr := strconv.ParseInt(parts[0], 10, 32)
	if sConvErr != nil {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected settingID:tagID. Error: %s", id, sConvErr)
	}

	tagID, tConvErr := strconv.ParseInt(parts[1], 10, 32)
	if tConvErr != nil {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected settingID:tagID. Error: %s", id, tConvErr)
	}

	return settingID, tagID, nil
}
