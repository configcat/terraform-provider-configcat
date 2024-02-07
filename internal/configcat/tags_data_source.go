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
	_ datasource.DataSource              = &tagDataSource{}
	_ datasource.DataSourceWithConfigure = &tagDataSource{}
)

func NewTagDataSource() datasource.DataSource {
	return &tagDataSource{}
}

type tagDataSource struct {
	client *client.Client
}

type tagDataModel struct {
	ID    types.String `tfsdk:"tag_id"`
	Name  types.String `tfsdk:"name"`
	Color types.String `tfsdk:"color"`
}

type tagDataSourceModel struct {
	ID              types.String   `tfsdk:"id"`
	ProductId       types.String   `tfsdk:"product_id"`
	NameFilterRegex types.String   `tfsdk:"name_filter_regex"`
	Data            []tagDataModel `tfsdk:"tags"`
}

func (d *tagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

func (d *tagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to access information about existing **" + TagResourceName + "s**. [What is a " + TagResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

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
				Description: "Filter the " + TagResourceName + "s by name.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			Tags: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						TagId: schema.StringAttribute{
							Description: "The unique " + TagResourceName + " ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the " + TagResourceName + ".",
							Computed:    true,
						},
						Color: schema.StringAttribute{
							Description: "The color of the " + TagResourceName + ".",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *tagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tagDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetTags(state.ProductId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+TagResourceName+" data, got error: %s", err))
		return
	}

	state.ID = types.StringValue(state.ProductId.ValueString() + strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.TagModel{}
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

	state.Data = make([]tagDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		dataModel := &tagDataModel{
			ID:    types.StringValue(strconv.FormatInt(*resource.TagId, 10)),
			Name:  types.StringPointerValue(resource.Name.Get()),
			Color: types.StringPointerValue(resource.Color.Get()),
		}

		state.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
