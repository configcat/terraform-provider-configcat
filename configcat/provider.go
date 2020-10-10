package configcat

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
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
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.configcat.com",
				Description: "ConfigCat Public Management API Base Url (defaults to production).",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"configcat_feature_flag": resourceConfigCatFatureFlag(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	basicAuthUsername := d.Get("basic_auth_username").(string)
	basicAuthPassword := d.Get("basic_auth_password").(string)

}
