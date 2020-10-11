package configcat

import (
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConfigCatConfig() *schema.Resource {
	return &schema.Resource{

		Read: configRead,

		Schema: map[string]*schema.Schema{
			"product_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"config_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func configRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	productId := fmt.Sprintf("%v", d.Get("product_id"))
	if productId == "" {
		return fmt.Errorf("product_id is required")
	}

	configs, err := c.GetConfigs(productId)
	if err != nil {
		return err
	}

	configId := fmt.Sprintf("%v", d.Get("config_id"))
	if configId == "" {
		return fmt.Errorf("config_id is required")
	}

	for i := range configs {
		if configs[i].ConfigId == configId {
			updateConfigResourceData(d, &configs[i], productId)
			return nil
		}
	}

	return fmt.Errorf("could not find Configs with product_id: %s and config_id: %s", productId, configId)

}

func updateConfigResourceData(d *schema.ResourceData, m *sw.ConfigModel, productId string) {
	d.SetId(m.ConfigId)
	d.Set("product_id", productId)
	d.Set("config_id", m.ConfigId)
	d.Set("name", m.Name)
}
