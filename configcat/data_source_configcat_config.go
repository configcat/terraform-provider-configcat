package configcat

import (
	"context"
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyProductID  = "product_id"
	keyConfigID   = "config_id"
	keyConfigName = "name"
)

func dataSourceConfigCatConfig() *schema.Resource {
	return &schema.Resource{

		ReadContext: configRead,

		Schema: map[string]*schema.Schema{
			keyProductID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			keyConfigName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			keyConfigID: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func configRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(keyProductID).(string)
	configName := d.Get(keyConfigName).(string)

	config, err := findConfig(c, productID, configName)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	updateConfigResourceData(d, config, productID)
	var diags diag.Diagnostics
	return diags
}

func findConfig(c *Client, productID, configName string) (*sw.ConfigModel, error) {

	configs, err := c.GetConfigs(productID)
	if err != nil {
		return nil, err
	}

	for i := range configs {
		if configs[i].Name == configName {
			return &configs[i], nil
		}
	}

	return nil, fmt.Errorf("could not find Config. product_id: %s name: %s", productID, configName)
}

func updateConfigResourceData(d *schema.ResourceData, m *sw.ConfigModel, productID string) {
	d.SetId(m.ConfigId)
	d.Set(keyProductID, productID)
	d.Set(keyConfigID, m.ConfigId)
	d.Set(keyConfigName, m.Name)
}
