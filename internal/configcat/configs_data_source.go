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

	sw "github.com/configcat/configcat-publicapi-go-client/v2"
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
	ID                types.String `tfsdk:"config_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Order             types.Int64  `tfsdk:"order"`
	EvaluationVersion types.String `tfsdk:"evaluation_version"`
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
		MarkdownDescription: "Use this data source to access information about existing **" + ConfigResourceName + "s**. [What is a " + ConfigResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the data source. Do not use.",
				Computed:    true,
			},
			ProductId: schema.StringAttribute{
				Description: "The ID of the " + ProductResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
			},
			NameFilterRegex: schema.StringAttribute{
				Description: "Filter the " + ConfigResourceName + "s by name.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			Configs: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ConfigId: schema.StringAttribute{
							Description: "The unique " + ConfigResourceName + " ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the " + ConfigResourceName + ".",
							Computed:    true,
						},
						Description: schema.StringAttribute{
							Description: "The description of the " + ConfigResourceName + ".",
							Computed:    true,
						},
						Order: schema.Int64Attribute{
							Description: "The order of the " + ConfigResourceName + " within a " + ProductResourceName + " (zero-based).",
							Computed:    true,
						},
						EvaluationVersion: schema.StringAttribute{
							MarkdownDescription: "The evaluation version of the " + ConfigResourceName + ". Possible values: `v1`|`v2`",
							Computed:            true,
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
	var state configDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetConfigs(state.ProductId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+ConfigResourceName+" data, got error: %s", err))
		return
	}

	state.ID = types.StringValue(state.ProductId.ValueString() + strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.ConfigModel{}
	if !state.NameFilterRegex.IsUnknown() && !state.NameFilterRegex.IsNull() && state.NameFilterRegex.ValueString() != "" {
		regex := regexp.MustCompile(state.NameFilterRegex.ValueString())
		for i := range resources {
			if regex.MatchString(*resources[i].Name.Get()) {
				filteredResources = append(filteredResources, resources[i])
			}
		}
	} else {
		filteredResources = resources
	}

	state.Data = make([]configDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		dataModel := &configDataModel{
			ID:                types.StringPointerValue(resource.ConfigId),
			Name:              types.StringPointerValue(resource.Name.Get()),
			Description:       types.StringPointerValue(resource.Description.Get()),
			Order:             types.Int64Value(int64(*resource.Order)),
			EvaluationVersion: types.StringPointerValue((*string)(resource.EvaluationVersion)),
		}

		state.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
