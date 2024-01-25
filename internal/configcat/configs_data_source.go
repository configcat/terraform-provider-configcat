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
	_ datasource.DataSource              = &dataSource{}
	_ datasource.DataSourceWithConfigure = &dataSource{}
)

func NewConfigDataSource() datasource.DataSource {
	return &dataSource{}
}

type dataSource struct {
	client *client.Client
}

type dataModel struct {
	ID          types.String `tfsdk:"config_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Order       types.Int64  `tfsdk:"order"`
}

type dataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	ProductId       types.String `tfsdk:"product_id"`
	NameFilterRegex types.String `tfsdk:"name_filter_regex"`
	Data            []dataModel  `tfsdk:"configs"`
}

func (d *dataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configs"
}

func (d *dataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Computed: true,
			},
			ProductId: schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Validators:          []validator.String{IsGuid()},
			},
			NameFilterRegex: schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
				Validators:          []validator.String{IsRegex()},
			},
			Configs: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ConfigId: schema.StringAttribute{
							Computed: true,
						},
						Name: schema.StringAttribute{
							Computed: true,
						},
						Description: schema.StringAttribute{
							Computed: true,
						},
						Order: schema.Int64Attribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *dataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceModel

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
				resp.Diagnostics.AddAttributeError(path.Root(NameFilterRegex), "invalid regex", "invalid regex")
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

	data.Data = make([]dataModel, len(filteredConfigs))
	for i, config := range filteredConfigs {
		configModel := &dataModel{
			ID:          types.StringValue(*config.ConfigId),
			Name:        types.StringValue(*config.Name.Get()),
			Description: types.StringValue(*config.Description.Get()),
			Order:       types.Int64Value(int64(*config.Order)),
		}

		data.Data[i] = *configModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
