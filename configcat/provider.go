package configcat

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	envBasicAuthUsername = "CONFIGCAT_BASIC_AUTH_USERNAME"
	envBasicAuthPassword = "CONFIGCAT_BASIC_AUTH_PASSWORD"
	envBasePath          = "CONFIGCAT_BASE_PATH"

	keyBasicAuthUsername = "basic_auth_username"
	keyBasicAuthPassword = "basic_auth_password"
	keyBasePath          = "base_path"

	defaultBasePath = "https://api.configcat.com"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			keyBasicAuthUsername: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(envBasicAuthUsername, nil),
				Description: "ConfigCat Public API credential - Basic Auth Username.",
			},
			keyBasicAuthPassword: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(envBasicAuthPassword, nil),
				Description: "ConfigCat Public API credential - Basic Auth Password",
			},
			keyBasePath: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(envBasePath, defaultBasePath),
				Description: "ConfigCat Public Management API Base Path (defaults to production).",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			//		"configcat_setting": resourceConfigCatSetting(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"configcat_product": dataSourceConfigCatProduct(),
			//	"configcat_config":  resourceConfigCatConfig(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	basicAuthUsername := d.Get(keyBasicAuthUsername).(string)
	basicAuthPassword := d.Get(keyBasicAuthPassword).(string)
	basePath := d.Get(keyBasePath).(string)

	var diags diag.Diagnostics

	client, err := NewClient(basePath, basicAuthUsername, basicAuthPassword)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Error setting up ConfigCat client: %s", err))
	}

	return client, diags
}
