package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestConfigValid(t *testing.T) {
	const dataSource = `
		data "configcat_configs" "test" {
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const productID = "08d86d63-2721-4da6-8c06-584521d516bc"
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_configs.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_configs.test", CONFIGS+".#", "1"),
				),
			},
		},
	})
}

func TestConfigValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_configs" "test" {
			name_filter_regex = "Main Config"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const productID = "08d86d63-2721-4da6-8c06-584521d516bc"
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_configs.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_configs.test", CONFIGS+".#", "1"),
					resource.TestCheckResourceAttr("data.configcat_configs.test", CONFIGS+".0."+CONFIG_ID, configID),
					resource.TestCheckResourceAttr("data.configcat_configs.test", CONFIGS+".0."+CONFIG_NAME, "Main Config"),
				),
			},
		},
	})
}

func TestConfigNotFoundFilter(t *testing.T) {
	const dataSource = `
		data "configcat_configs" "test" {
			name_filter_regex = "notfound"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_configs.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_configs.test", CONFIGS+".#", "0"),
				),
			},
		},
	})
}

func TestConfigInvalidGuid(t *testing.T) {
	const dataSource = `
		data "configcat_configs" "test" {
			name_filter_regex = "notfound"
			product_id = "invalidGuid"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      dataSource,
				ExpectError: regexp.MustCompile(`"product_id": invalid GUID`),
			},
		},
	})
}
