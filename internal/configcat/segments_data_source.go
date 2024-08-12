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
	_ datasource.DataSource              = &segmentDataSource{}
	_ datasource.DataSourceWithConfigure = &segmentDataSource{}
)

func NewSegmentDataSource() datasource.DataSource {
	return &segmentDataSource{}
}

type segmentDataSource struct {
	client *client.Client
}

type segmentDataModel struct {
	ID          types.String `tfsdk:"segment_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type segmentDataSourceModel struct {
	ID              types.String       `tfsdk:"id"`
	ProductId       types.String       `tfsdk:"product_id"`
	NameFilterRegex types.String       `tfsdk:"name_filter_regex"`
	Data            []segmentDataModel `tfsdk:"segments"`
}

func (d *segmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segments"
}

func (d *segmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to access information about existing **" + SegmentResourceName + "s**. [What is a " + SegmentResourceName + " in ConfigCat?](https://configcat.com/docs/advanced/segments)",

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
				Description: "Filter the " + SegmentResourceName + "s by name.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			Segments: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						SegmentId: schema.StringAttribute{
							Description: "The unique " + SegmentResourceName + " ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the " + SegmentResourceName + ".",
							Computed:    true,
						},
						Description: schema.StringAttribute{
							Description: "The description of the " + SegmentResourceName + ".",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *segmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *segmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state segmentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetSegments(state.ProductId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SegmentResourceName+" data, got error: %s", err))
		return
	}

	state.ID = types.StringValue(state.ProductId.ValueString() + strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.SegmentListModel{}
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

	state.Data = make([]segmentDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		dataModel := &segmentDataModel{
			ID:          types.StringPointerValue(resource.SegmentId),
			Name:        types.StringPointerValue(resource.Name.Get()),
			Description: types.StringPointerValue(resource.Description.Get()),
		}

		state.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
