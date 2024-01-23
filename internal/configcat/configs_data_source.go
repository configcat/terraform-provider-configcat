// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configcat

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var (
	_ datasource.DataSource              = &configDataSource{}
	_ datasource.DataSourceWithConfigure = &configDataSource{}
)

func NewConfigDataSource() datasource.DataSource {
	return &configDataSource{}
}

type configDataSource struct {
	client *client.Client
}

type configModel struct {
	ConfigId    types.String `tfsdk:"config_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Order       types.Int64  `tfsdk:"order"`
}

type configDataSourceModel struct {
	ID              types.String  `tfsdk:"id"`
	ProductId       types.String  `tfsdk:"product_id"`
	NameFilterRegex types.String  `tfsdk:"name_filter_regex"`
	Configs         []configModel `tfsdk:"configs"`
}

func (d *configDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configs"
}

func (d *configDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Computed: true,
			},
			PRODUCT_ID: schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Validators:          []validator.String{IsGuid()},
			},
			CONFIG_NAME_FILTER_REGEX: schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
			},
			CONFIGS: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						CONFIG_ID: schema.StringAttribute{
							Computed: true,
						},
						CONFIG_NAME: schema.StringAttribute{
							Computed: true,
						},
						CONFIG_DESCRIPTION: schema.StringAttribute{
							Computed: true,
						},
						CONFIG_ORDER: schema.Int64Attribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *configDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *configDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data configDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	configs, err := d.client.GetConfigs(data.ProductId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read configs, got error: %s", err))
		return
	}

	data.ID = types.StringValue(data.ProductId.ValueString() + strconv.FormatInt(time.Now().Unix(), 10))

	filteredConfigs := []sw.ConfigModel{}
	if !data.NameFilterRegex.IsUnknown() && !data.NameFilterRegex.IsNull() && data.NameFilterRegex.ValueString() != "" {
		regex, err := regexp.Compile(data.NameFilterRegex.ValueString())
		if err != nil {
			if err != nil {
				resp.Diagnostics.AddAttributeError(path.Root(CONFIG_NAME_FILTER_REGEX), "invalid regex", "invalid regex")
				return
			}
		}

		for i := range configs {
			if regex.MatchString(*configs[i].Name.Get()) {
				filteredConfigs = append(filteredConfigs, configs[i])
			}
		}
	} else {
		filteredConfigs = configs
	}

	data.Configs = make([]configModel, len(filteredConfigs))
	for i, config := range filteredConfigs {
		configModel := &configModel{
			ConfigId:    types.StringValue(*config.ConfigId),
			Name:        types.StringValue(*config.Name.Get()),
			Description: types.StringValue(*config.Description.Get()),
			Order:       types.Int64Value(int64(*config.Order)),
		}

		data.Configs[i] = *configModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
