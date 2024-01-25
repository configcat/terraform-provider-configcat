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

type configDataModel struct {
	ID          types.String `tfsdk:"config_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Order       types.Int64  `tfsdk:"order"`
}

type configDataSourceModel struct {
	ID              types.String      `tfsdk:"id"`
	ProductId       types.String      `tfsdk:"product_id"`
	NameFilterRegex types.String      `tfsdk:"name_filter_regex"`
	Data            []configDataModel `tfsdk:"configs"`
}

func (d *configDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configs"
}

func (d *configDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Use this data source to access information about existing **Configs**. [What is a Config in ConfigCat?](https://configcat.com/docs/main-concepts)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Computed: true,
			},
			ProductId: schema.StringAttribute{
				Description: "The ID of the Product.",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
			},
			NameFilterRegex: schema.StringAttribute{
				Description: "Filter the Configs by name.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			Configs: schema.ListNestedAttribute{
				MarkdownDescription: "A config [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ConfigId: schema.StringAttribute{
							Description: "The unique Config ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the Config.",
							Computed:    true,
						},
						Description: schema.StringAttribute{
							Description: "The description of the Config.",
							Computed:    true,
						},
						Order: schema.Int64Attribute{
							Description: "The order of the Config within a Product (zero-based).",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *configDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetConfigs(data.ProductId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read configs, got error: %s", err))
		return
	}

	data.ID = types.StringValue(data.ProductId.ValueString() + strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.ConfigModel{}
	if !data.NameFilterRegex.IsUnknown() && !data.NameFilterRegex.IsNull() && data.NameFilterRegex.ValueString() != "" {
		regex := regexp.MustCompile(data.NameFilterRegex.ValueString())
		for i := range resources {
			if regex.MatchString(*resources[i].Name.Get()) {
				filteredResources = append(filteredResources, resources[i])
			}
		}
	} else {
		filteredResources = resources
	}

	data.Data = make([]configDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		dataModel := &configDataModel{
			ID:          types.StringValue(*resource.ConfigId),
			Name:        types.StringValue(*resource.Name.Get()),
			Description: types.StringValue(*resource.Description.Get()),
			Order:       types.Int64Value(int64(*resource.Order)),
		}

		data.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
