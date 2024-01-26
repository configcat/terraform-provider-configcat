package configcat

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
)

const (
	// Environment variables
	ENV_BASIC_AUTH_USERNAME = "CONFIGCAT_BASIC_AUTH_USERNAME"
	ENV_BASIC_AUTH_PASSWORD = "CONFIGCAT_BASIC_AUTH_PASSWORD"
	ENV_BASE_PATH           = "CONFIGCAT_BASE_PATH"

	// Parameters
	KEY_BASIC_AUTH_USERNAME = "basic_auth_username"
	KEY_BASIC_AUTH_PASSWORD = "basic_auth_password"
	KEY_BASE_PATH           = "base_path"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &configCatProvider{}

// configCatProvider defines the provider implementation.
type configCatProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// configCatProviderModel describes the provider data model.
type configCatProviderModel struct {
	BasicAuthUsername types.String `tfsdk:"basic_auth_username"`
	BasicAuthPassword types.String `tfsdk:"basic_auth_password"`
	BasePath          types.String `tfsdk:"base_path"`
}

func (p *configCatProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "configcat"
	resp.Version = p.version
}

func (p *configCatProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ConfigCat Feature Flags Provider for Terraform allows you to manage ConfigCat resources from a Terraform script. The Provider uses the standard [ConfigCat Public Management API](https://api.configcat.com/).",
		Attributes: map[string]schema.Attribute{
			KEY_BASIC_AUTH_USERNAME: schema.StringAttribute{
				MarkdownDescription: "Get your `basic_auth_username` at [ConfigCat Public API credentials](https://app.configcat.com/my-account/public-api-credentials).  This can also be sourced from the `CONFIGCAT_BASIC_AUTH_USERNAME` Environment Variable.",
				Optional:            true,
			},
			KEY_BASIC_AUTH_PASSWORD: schema.StringAttribute{
				MarkdownDescription: "Get your `basic_auth_password` at [ConfigCat Public API credentials](https://app.configcat.com/my-account/public-api-credentials).  This can also be sourced from the `CONFIGCAT_BASIC_AUTH_PASSWORD` Environment Variable.",
				Optional:            true,
				Sensitive:           true,
			},
			KEY_BASE_PATH: schema.StringAttribute{
				MarkdownDescription: "ConfigCat Public Management API's `base_path`. Defaults to [https://api.configcat.com](https://api.configcat.com).  This can also be sourced from the `CONFIGCAT_BASE_PATH` Environment Variable.",
				Optional:            true,
			},
		},
	}
}

func (p *configCatProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var providerConfig configCatProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &providerConfig)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if providerConfig.BasePath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASE_PATH),
			"Unknown ConfigCat Public Management API "+KEY_BASE_PATH,
			"The provider cannot create the ConfigCat Public Management API client as there is an unknown configuration value for the "+KEY_BASE_PATH+".",
		)
	}

	if providerConfig.BasicAuthUsername.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASIC_AUTH_USERNAME),
			"Unknown ConfigCat Public Management API "+KEY_BASIC_AUTH_USERNAME,
			"The provider cannot create the ConfigCat Public Management API client as there is an unknown configuration value for the "+KEY_BASIC_AUTH_USERNAME+".",
		)
	}

	if providerConfig.BasicAuthPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASIC_AUTH_PASSWORD),
			"Unknown ConfigCat Public Management API "+KEY_BASIC_AUTH_PASSWORD,
			"The provider cannot create the ConfigCat Public Management API client as there is an unknown configuration value for the "+KEY_BASIC_AUTH_PASSWORD+".",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	basePath := os.Getenv(ENV_BASE_PATH)
	basicAuthUsername := os.Getenv(ENV_BASIC_AUTH_USERNAME)
	basicAuthPassword := os.Getenv(ENV_BASIC_AUTH_PASSWORD)

	if !providerConfig.BasePath.IsNull() {
		basePath = providerConfig.BasePath.ValueString()
	}

	if !providerConfig.BasicAuthUsername.IsNull() {
		basicAuthUsername = providerConfig.BasicAuthUsername.ValueString()
	}

	if !providerConfig.BasicAuthPassword.IsNull() {
		basicAuthPassword = providerConfig.BasicAuthPassword.ValueString()
	}

	if basePath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASE_PATH),
			"Missing "+KEY_BASE_PATH,
			"The provider cannot create the ConfigCat Public Management API client as there is a missing or empty value for the "+KEY_BASE_PATH+"."+
				"Set the "+KEY_BASE_PATH+" value in the configuration or use the "+ENV_BASE_PATH+" environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if basicAuthUsername == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASIC_AUTH_USERNAME),
			"Missing "+KEY_BASIC_AUTH_USERNAME,
			"The provider cannot create the ConfigCat Public Management API client as there is a missing or empty value for the "+KEY_BASIC_AUTH_USERNAME+"."+
				"Set the "+KEY_BASIC_AUTH_USERNAME+" value in the configuration or use the "+ENV_BASIC_AUTH_USERNAME+" environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if basicAuthPassword == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASIC_AUTH_PASSWORD),
			"Missing "+KEY_BASIC_AUTH_PASSWORD,
			"The provider cannot create the ConfigCat Public Management API client as there is a missing or empty value for the "+KEY_BASE_PATH+"."+
				"Set the "+KEY_BASIC_AUTH_PASSWORD+" value in the configuration or use the "+ENV_BASIC_AUTH_PASSWORD+" environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := client.NewClient(basePath, basicAuthUsername, basicAuthPassword)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create ConfigCat Public Management API client",
			"An unexpected error occurred when creating the ConfigCat Public Management API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"ConfigCat Public Management API client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *configCatProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProductResource,
		NewConfigResource,
		NewEnvironmentResource,
		NewSettingResource,
		NewSettingValueResource,
		NewPermissionGroupResource,
		NewSegmentResource,
		NewTagResource,
		NewSettingTagResource,
	}
}

func (p *configCatProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationDataSource,
		NewProductDataSource,
		NewConfigDataSource,
		NewEnvironmentDataSource,
		NewSdkKeyDataSource,
		NewSegmentDataSource,
		NewSettingDataSource,
		NewTagDataSource,
		NewPermissionGroupDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &configCatProvider{
			version: version,
		}
	}
}
