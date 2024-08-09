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
	_ datasource.DataSource              = &settingDataSource{}
	_ datasource.DataSourceWithConfigure = &settingDataSource{}
)

func NewSettingDataSource() datasource.DataSource {
	return &settingDataSource{}
}

type settingDataSource struct {
	client *client.Client
}

type settingDataModel struct {
	ID          types.String `tfsdk:"setting_id"`
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Hint        types.String `tfsdk:"hint"`
	SettingType types.String `tfsdk:"setting_type"`
	Order       types.Int64  `tfsdk:"order"`
}

type settingDataSourceModel struct {
	ID             types.String       `tfsdk:"id"`
	ConfigId       types.String       `tfsdk:"config_id"`
	KeyFilterRegex types.String       `tfsdk:"key_filter_regex"`
	Data           []settingDataModel `tfsdk:"settings"`
}

func (d *settingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (d *settingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to access information about existing **" + SettingsResourceName + "**. [What is a " + SettingResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

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
			SettingKeyFilterRegex: schema.StringAttribute{
				Description: "Filter the " + SettingsResourceName + "s by key.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			Settings: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						SettingId: schema.StringAttribute{
							Description: "The unique " + SettingResourceName + " ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the " + SettingResourceName + ".",
							Computed:    true,
						},
						SettingKey: schema.StringAttribute{
							Description: "The key of the " + SettingResourceName + ".",
							Computed:    true,
						},
						SettingHint: schema.StringAttribute{
							Description: "The hint of the " + SettingResourceName + ".",
							Computed:    true,
						},
						SettingType: schema.StringAttribute{
							Description: "The " + SettingResourceName + "'s type. Available values: `boolean`|`string`|`int`|`double`.",
							Computed:    true,
						},
						Order: schema.Int64Attribute{
							Description: "The order of the " + SettingResourceName + " within a " + ConfigResourceName + " (zero-based).",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *settingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *settingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state settingDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetSettings(state.ConfigId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SettingResourceName+" data, got error: %s", err))
		return
	}

	state.ID = types.StringValue(state.ConfigId.ValueString() + strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.SettingModel{}
	if !state.KeyFilterRegex.IsUnknown() && !state.KeyFilterRegex.IsNull() && state.KeyFilterRegex.ValueString() != "" {
		regex := regexp.MustCompile(state.KeyFilterRegex.ValueString())
		for i := range resources {
			if regex.MatchString(*resources[i].Key.Get()) {
				filteredResources = append(filteredResources, resources[i])
			}
		}
	} else {
		filteredResources = resources
	}

	state.Data = make([]settingDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		dataModel := &settingDataModel{
			ID:          types.StringValue(strconv.FormatInt(int64(*resource.SettingId), 10)),
			Name:        types.StringPointerValue(resource.Name.Get()),
			Key:         types.StringPointerValue(resource.Key.Get()),
			Hint:        types.StringPointerValue(resource.Hint.Get()),
			SettingType: types.StringPointerValue((*string)(resource.SettingType)),
			Order:       types.Int64Value(int64(*resource.Order)),
		}

		state.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
