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
	_ datasource.DataSource              = &productDataSource{}
	_ datasource.DataSourceWithConfigure = &productDataSource{}
)

func NewProductDataSource() datasource.DataSource {
	return &productDataSource{}
}

type productDataSource struct {
	client *client.Client
}

type productDataModel struct {
	ID          types.String `tfsdk:"product_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Order       types.Int64  `tfsdk:"order"`
}

type productDataSourceModel struct {
	ID              types.String       `tfsdk:"id"`
	NameFilterRegex types.String       `tfsdk:"name_filter_regex"`
	Data            []productDataModel `tfsdk:"products"`
}

func (d *productDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_products"
}

func (d *productDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to access information about existing **" + ProductResourceName + "s**. [What is a " + ProductResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the data source. Do not use.",
				Computed:    true,
			},
			NameFilterRegex: schema.StringAttribute{
				Description: "Filter the " + ProductResourceName + "s by name.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			Products: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ProductId: schema.StringAttribute{
							Description: "The unique " + ProductResourceName + " ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the " + ProductResourceName + ".",
							Computed:    true,
						},
						Description: schema.StringAttribute{
							Description: "The description of the " + ProductResourceName + ".",
							Computed:    true,
						},
						Order: schema.Int64Attribute{
							Description: "The order of the " + ProductResourceName + " within a " + ProductResourceName + " (zero-based).",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *productDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *productDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state productDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetProducts()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+ProductResourceName+" data, got error: %s", err))
		return
	}

	state.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.ProductModel{}
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

	state.Data = make([]productDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		dataModel := &productDataModel{
			ID:          types.StringPointerValue(resource.ProductId),
			Name:        types.StringPointerValue(resource.Name.Get()),
			Description: types.StringPointerValue(resource.Description.Get()),
			Order:       types.Int64Value(int64(*resource.Order)),
		}

		state.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
