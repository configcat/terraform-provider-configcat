package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceSegmentFlow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_segment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName"
						description = "testDescription"
						comparison_attribute = "email"
						comparator = "sensitiveIsOneOf"
						comparison_value="a@b.com,c@d.com"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_segment.test", "id"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_NAME, "testName"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_DESCRIPTION, "testDescription"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARISON_ATTRIBUTE, "email"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARATOR, "sensitiveIsOneOf"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARISON_VALUE, "a@b.com,c@d.com"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_segment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName2"
						description = "testDescription2"
						comparison_attribute = "version"
						comparator = "semVerLess"
						comparison_value="2.0.0"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_segment.test", "id"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_NAME, "testName2"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_DESCRIPTION, "testDescription2"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARISON_ATTRIBUTE, "version"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARATOR, "semVerLess"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARISON_VALUE, "2.0.0"),
				),
			},
			{
				ResourceName:       "configcat_segment.test",
				ImportState:        true,
				ImportStateVerify:  true,
				ImportStatePersist: true,
			},
		},
	})
}

func TestResourceSegmentWrongComparator(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "configcat_products" "products" {
				}
				resource "configcat_segment" "test" {
					product_id = data.configcat_products.products.products.0.product_id
					name = "testName"
					description = "testDescription"
					comparison_attribute = "email"
					comparator = "isOneOf"
					comparison_value="a@b.com,c@d.com"
				}
				`,
				ExpectError: regexp.MustCompile(`could not parse Comparator`),
			},
			{
				Config: `
				data "configcat_products" "products" {
				}
				resource "configcat_segment" "test" {
					product_id = data.configcat_products.products.products.0.product_id
					name = "testName"
					description = "testDescription"
					comparison_attribute = "email"
					comparator = "sensitiveIsOneOf"
					comparison_value="a@b.com,c@d.com"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_segment.test", "id"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_NAME, "testName"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_DESCRIPTION, "testDescription"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARISON_ATTRIBUTE, "email"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARATOR, "sensitiveIsOneOf"),
					resource.TestCheckResourceAttr("configcat_segment.test", SEGMENT_COMPARISON_VALUE, "a@b.com,c@d.com"),
				),
			},
			{
				ResourceName:       "configcat_segment.test",
				ImportState:        true,
				ImportStateVerify:  true,
				ImportStatePersist: true,
			},
		},
	})
}
