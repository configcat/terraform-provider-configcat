package configcat

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &webhookSigningKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &webhookSigningKeysDataSource{}
)

func NewWebhookSigningKeysDataSource() datasource.DataSource {
	return &webhookSigningKeysDataSource{}
}

type webhookSigningKeysDataSource struct {
	client *client.Client
}

type webhookSigningKeysDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	WebhookId types.Int64  `tfsdk:"webhook_id"`
	Key1      types.String `tfsdk:"key1"`
	Key2      types.String `tfsdk:"key2"`
}

func (d *webhookSigningKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_signing_keys"
}

func (d *webhookSigningKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to access information about existing **" + ProductResourceName + "s**. [What is a " + ProductResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the data source. Do not use.",
				Computed:    true,
			},
			WebhookId: schema.Int64Attribute{
				Description: "The ID of the " + WebhookResourceName + ".",
				Required:    true,
			},
			WebhookSigningKeyKey1: schema.StringAttribute{
				MarkdownDescription: "The first signing key.",
				Computed:            true,
			},
			WebhookSigningKeyKey2: schema.StringAttribute{
				MarkdownDescription: "The second signing key.",
				Computed:            true,
			},
		},
	}
}

func (d *webhookSigningKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *webhookSigningKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state webhookSigningKeysDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	webhookSigningKeys, err := d.client.GetWebhookSigningKeys(int32(state.WebhookId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+WebhookSigningKeysResourceName+", got error: %s", err))
		return
	}

	state.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	state.Key1 = types.StringPointerValue(webhookSigningKeys.Key1.Get())
	state.Key2 = types.StringPointerValue(webhookSigningKeys.Key2.Get())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
