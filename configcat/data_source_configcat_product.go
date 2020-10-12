package configcat

import (
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatProduct() *schema.Resource {
	return &schema.Resource{

		Read: productRead,

		Schema: map[string]*schema.Schema{
			"product_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func productRead(d *schema.ResourceData, meta interface{}) error {
	product, err := findProduct(d, meta)
	if err != nil {
		return err
	}

	updateProductResourceData(d, product)
	return nil
}

func findProduct(d *schema.ResourceData, meta interface{}) (*sw.ProductModel, error) {

	c := meta.(*Client)

	products, err := c.GetProducts()
	if err != nil {
		return nil, err
	}

	productName := d.Get("name").(string)
	if productName == "" {
		return nil, fmt.Errorf("name is required")
	}
	for i := range products {
		if products[i].Name == productName {
			return &products[i], nil
		}
	}

	return nil, fmt.Errorf("could not find Product with name: %s", productName)
}

func updateProductResourceData(d *schema.ResourceData, m *sw.ProductModel) {
	d.SetId(m.ProductId)
	d.Set("product_id", m.ProductId)
	d.Set("name", m.Name)
}
