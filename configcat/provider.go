package configcat

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	// Defaults
	DEFAULT_BASE_PATH = "https://api.configcat.com"

	// Data sources
	KEY_ORGANIZATIONS = "configcat_organizations"
	KEY_PRODUCTS      = "configcat_products"
	KEY_CONFIGS       = "configcat_configs"
	KEY_ENVIRONMENTS  = "configcat_environments"
	KEY_SETTINGS      = "configcat_settings"

	// Resources
	KEY_SETTING       = "configcat_setting"
	KEY_SETTING_VALUE = "configcat_setting_value"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			KEY_BASIC_AUTH_USERNAME: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_BASIC_AUTH_USERNAME, nil),
				Description: "ConfigCat Public API credential - Basic Auth Username.",
			},
			KEY_BASIC_AUTH_PASSWORD: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_BASIC_AUTH_PASSWORD, nil),
				Description: "ConfigCat Public API credential - Basic Auth Password",
			},
			KEY_BASE_PATH: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_BASE_PATH, DEFAULT_BASE_PATH),
				Description: "ConfigCat Public Management API Base Path (defaults to production).",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			KEY_SETTING:       resourceConfigCatSetting(),
			KEY_SETTING_VALUE: resourceConfigCatSettingValue(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			KEY_ORGANIZATIONS: dataSourceConfigCatOrganizations(),
			KEY_PRODUCTS:      dataSourceConfigCatProducts(),
			KEY_CONFIGS:       dataSourceConfigCatConfigs(),
			KEY_ENVIRONMENTS:  dataSourceConfigCatEnvironments(),
			KEY_SETTINGS:      dataSourceConfigCatSettings(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	basicAuthUsername := d.Get(KEY_BASIC_AUTH_USERNAME).(string)
	basicAuthPassword := d.Get(KEY_BASIC_AUTH_PASSWORD).(string)
	basePath := d.Get(KEY_BASE_PATH).(string)

	var diags diag.Diagnostics

	client, err := NewClient(basePath, basicAuthUsername, basicAuthPassword)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Error setting up ConfigCat client: %s", err))
	}

	return client, diags
}
