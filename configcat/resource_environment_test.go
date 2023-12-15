package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceEnvironmentFlow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_environment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestResourceEnvironmentFlow"
						description = "testDescription"
						color = "blue"
						order = 10
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "TestResourceEnvironmentFlow"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_DESCRIPTION, "testDescription"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_COLOR, "blue"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_ORDER, "10"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_environment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestResourceEnvironmentFlow2"
						description = "testDescription2"
						color = "yellow"
						order = 11
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "TestResourceEnvironmentFlow2"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_DESCRIPTION, "testDescription2"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_COLOR, "yellow"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_ORDER, "11"),
				),
			},
			{
				ResourceName:      "configcat_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceEnvironmentWrongColor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_environment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestResourceEnvironmentWrongColor"
						description = "testDescription"
						color = "notvalid"
						order = 20
					}
				`,
				ExpectError: regexp.MustCompile(`Invalid color.`),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_environment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestResourceEnvironmentWrongColor"
						description = "testDescription"
						color = "yellow"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "TestResourceEnvironmentWrongColor"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_DESCRIPTION, "testDescription"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_COLOR, "yellow"),
				),
			},
			{
				ResourceName:      "configcat_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
