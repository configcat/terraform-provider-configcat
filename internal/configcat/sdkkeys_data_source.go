package configcat

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &sdkKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &sdkKeyDataSource{}
)

func NewSdkKeyDataSource() datasource.DataSource {
	return &sdkKeyDataSource{}
}

type sdkKeyDataSource struct {
	client *client.Client
}

type sdkKeyDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	ConfigId        types.String `tfsdk:"config_id"`
	EnvironmentId   types.String `tfsdk:"environment_id"`
	PrimarySdkKey   types.String `tfsdk:"primary"`
	SecondarySdkKey types.String `tfsdk:"secondary"`
}

func (d *sdkKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdkkeys"
}

func (d *sdkKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to access information about existing **" + SdkKeyResourceName + "s**.",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the data source. Do not use.",
				Computed:    true,
			},
			ConfigId: schema.StringAttribute{
				Description: "The ID of the " + ConfigResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
			},
			EnvironmentId: schema.StringAttribute{
				Description: "The ID of the " + EnvironmentResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
			},
			PrimarySdkKey: schema.StringAttribute{
				MarkdownDescription: "The primary SDK Key associated with your **Config** and **Environment**.",
				Computed:            true,
			},
			SecondarySdkKey: schema.StringAttribute{
				MarkdownDescription: "The secondary SDK Key associated with your **Config** and **Environment**.",
				Computed:            true,
			},
		},
	}
}

func (d *sdkKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *sdkKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sdkKeyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sdkKeys, err := d.client.GetSdkKeys(state.ConfigId.ValueString(), state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SdkKeyResourceName+"s, got error: %s", err))
		return
	}

	state.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	state.PrimarySdkKey = types.StringPointerValue(sdkKeys.Primary.Get())
	state.SecondarySdkKey = types.StringPointerValue(sdkKeys.Secondary.Get())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
