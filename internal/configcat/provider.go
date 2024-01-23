// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
var _ provider.Provider = &ConfigCatProvider{}

// ConfigCatProvider defines the provider implementation.
type ConfigCatProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ConfigCatProviderModel describes the provider data model.
type ConfigCatProviderModel struct {
	basicAuthUsername types.String `tfsdk:"basic_auth_username"`
	basicAuthPassword types.String `tfsdk:"basic_auth_password"`
	basePath          types.String `tfsdk:"base_path"`
}

func (p *ConfigCatProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "configcat"
	resp.Version = p.version
}

func (p *ConfigCatProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			KEY_BASIC_AUTH_USERNAME: schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Sensitive:           true,
			},
			KEY_BASIC_AUTH_PASSWORD: schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Sensitive:           true,
			},
			KEY_BASE_PATH: schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *ConfigCatProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var providerConfig ConfigCatProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &providerConfig)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if providerConfig.basePath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASE_PATH),
			"Unknown ConfigCat Public Management API "+KEY_BASE_PATH,
			"The provider cannot create the ConfigCat Public Management API client as there is an unknown configuration value for the "+KEY_BASE_PATH+".",
		)
	}

	if providerConfig.basicAuthUsername.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASIC_AUTH_USERNAME),
			"Unknown ConfigCat Public Management API "+KEY_BASIC_AUTH_USERNAME,
			"The provider cannot create the ConfigCat Public Management API client as there is an unknown configuration value for the "+KEY_BASIC_AUTH_USERNAME+".",
		)
	}

	if providerConfig.basicAuthPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(KEY_BASIC_AUTH_PASSWORD),
			"Unknown ConfigCat Public Management API "+KEY_BASIC_AUTH_PASSWORD,
			"The provider cannot create the ConfigCat Public Management API client as there is an unknown configuration value for the "+KEY_BASIC_AUTH_PASSWORD+".",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	basePath := os.Getenv(ENV_BASE_PATH)
	basicAuthUsername := os.Getenv(ENV_BASIC_AUTH_USERNAME)
	basicAuthPassword := os.Getenv(ENV_BASIC_AUTH_PASSWORD)

	if !providerConfig.basePath.IsNull() {
		basePath = providerConfig.basePath.ValueString()
	}

	if !providerConfig.basicAuthUsername.IsNull() {
		basicAuthUsername = providerConfig.basicAuthUsername.ValueString()
	}

	if !providerConfig.basicAuthPassword.IsNull() {
		basicAuthPassword = providerConfig.basicAuthPassword.ValueString()
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

func (p *ConfigCatProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *ConfigCatProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConfigDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ConfigCatProvider{
			version: version,
		}
	}
}
