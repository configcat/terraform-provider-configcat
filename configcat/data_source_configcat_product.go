package configcat

import (
	"context"
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatProduct() *schema.Resource {
	return &schema.Resource{

		ReadContext: productRead,

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

func productRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	product, err := findProduct(d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	updateProductResourceData(d, product)
	var diags diag.Diagnostics
	return diags
}

func findProduct(d *schema.ResourceData, m interface{}) (*sw.ProductModel, error) {

	c := m.(*Client)

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
