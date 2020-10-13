package configcat

import (
	"context"
	"regexp"
	"strconv"
	"time"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatProduct() *schema.Resource {
	return &schema.Resource{

		ReadContext: productRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_NAME_FILTER_REGEX: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			PRODUCTS: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PRODUCT_ID: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						PRODUCT_NAME: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func productRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productNameFilterRegex := d.Get(PRODUCT_NAME_FILTER_REGEX).(string)

	products, err := getProducts(c)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredProducts := []sw.ProductModel{}
	if productNameFilterRegex == "" {
		filteredProducts = products
	} else {
		regex := regexp.MustCompile(productNameFilterRegex)
		for i := range products {
			if regex.MatchString(products[i].Name) {
				filteredProducts = append(filteredProducts, products[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(PRODUCTS, flattenProductsData(&filteredProducts))

	var diags diag.Diagnostics
	return diags
}

func getProducts(c *Client) ([]sw.ProductModel, error) {

	products, err := c.GetProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func flattenProductsData(products *[]sw.ProductModel) []interface{} {
	if products != nil {
		elements := make([]interface{}, len(*products), len(*products))

		for i, product := range *products {
			element := make(map[string]interface{})

			element[PRODUCT_ID] = product.ProductId
			element[PRODUCT_NAME] = product.Name

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
