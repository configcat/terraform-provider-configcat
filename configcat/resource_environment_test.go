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
						name = "testName"
						description = "testDescription"
						color = "blue"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "testName"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_DESCRIPTION, "testDescription"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_COLOR, "blue"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_environment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName2"
						description = "testDescription2"
						color = "yellow"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "testName2"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_DESCRIPTION, "testDescription2"),
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
						name = "testName"
						description = "testDescription"
						color = "notvalid"
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
						name = "testName"
						description = "testDescription"
						color = "yellow"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "testName"),
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
