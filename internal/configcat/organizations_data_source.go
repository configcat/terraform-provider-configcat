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
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

func NewOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

type organizationDataSource struct {
	client *client.Client
}

type organizationDataModel struct {
	ID   types.String `tfsdk:"organization_id"`
	Name types.String `tfsdk:"name"`
}

type organizationDataSourceModel struct {
	ID              types.String            `tfsdk:"id"`
	NameFilterRegex types.String            `tfsdk:"name_filter_regex"`
	Data            []organizationDataModel `tfsdk:"organizations"`
}

func (d *organizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *organizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Use this data source to access information about existing **" + OrganizationResourceName + "s**. [What is an " + OrganizationResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the data source. Do not use.",
				Computed:    true,
			},
			NameFilterRegex: schema.StringAttribute{
				Description: "Filter the " + OrganizationResourceName + "s by name.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			Organizations: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						OrganizationId: schema.StringAttribute{
							Description: "The unique " + OrganizationResourceName + " ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the " + OrganizationResourceName + ".",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *organizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organizationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetOrganizations()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+OrganizationResourceName+" data, got error: %s", err))
		return
	}

	state.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.OrganizationModel{}
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

	state.Data = make([]organizationDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		dataModel := &organizationDataModel{
			ID:   types.StringPointerValue(resource.OrganizationId),
			Name: types.StringPointerValue(resource.Name.Get()),
		}

		state.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
