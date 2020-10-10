package configcat

import (
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConfigCatProduct() *schema.Resource {
	return &schema.Resource{
		Read:   productRead,
		Exists: productExists,

		Schema: map[string]*schema.Schema{
			"productId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

func productExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	_, err := findProduct(d, meta)
	if err != nil {
		return false, err
	}

	return true, nil
}

func findProduct(d *schema.ResourceData, meta interface{}) (*sw.ProductModel, error) {
	c := meta.(*Client)

	products, err := c.GetProducts()
	if err != nil {
		return nil, err
	}

	productId := d.Get("productId")
	if productId != "" {
		for i := range products {
			if products[i].ProductId == productId {
				return &products[i], nil
			}
		}
	}

	return nil, fmt.Errorf("could not find Product with ProductId: %s", productId)
}

func updateProductResourceData(d *schema.ResourceData, m *sw.ProductModel) {
	d.Set("productId", m.ProductId)
	d.Set("name", m.Name)
}
