package configcat

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// test
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"basic_auth_username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONFIGCAT_BASIC_AUTH_USERNAME", nil),
				Description: "ConfigCat Public API credential - Basic Auth Username.",
			},
			"basic_auth_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONFIGCAT_BASIC_AUTH_PASSWORD", nil),
				Description: "ConfigCat Public API credential - Basic Auth Password",
			},
			"base_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONFIGCAT_BASE_PATH", "https://api.configcat.com"),
				Description: "ConfigCat Public Management API Base Path (defaults to production).",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"configcat_setting": resourceConfigCatSetting(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"configcat_product": dataSourceConfigCatProduct(),
			"configcat_config":  resourceConfigCatConfig(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	basicAuthUsername := d.Get("basic_auth_username").(string)
	basicAuthPassword := d.Get("basic_auth_password").(string)
	basePath := d.Get("base_path").(string)

	client, err := NewClient(basePath, basicAuthUsername, basicAuthPassword)
	if err != nil {
		return nil, fmt.Errorf("Error setting up ConfigCat client: %s", err)
	}

	return client, nil
}
