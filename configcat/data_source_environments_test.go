package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testEnvironmentResourceName = "data.configcat_environments.test"

func TestEnvironmentValid(t *testing.T) {
	const dataSource = `
		data "configcat_environments" "test" {
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testEnvironmentResourceName, "id"),
					resource.TestCheckResourceAttr(testEnvironmentResourceName, ENVIRONMENTS+".#", "3"),
				),
			},
		},
	})
}

func TestEnvironmentValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_environments" "test" {
			name_filter_regex = "Test"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testEnvironmentResourceName, "id"),
					resource.TestCheckResourceAttr(testEnvironmentResourceName, ENVIRONMENTS+".#", "1"),
					resource.TestCheckResourceAttr(testEnvironmentResourceName, ENVIRONMENTS+".0."+ENVIRONMENT_ID, environmentID),
					resource.TestCheckResourceAttr(testEnvironmentResourceName, ENVIRONMENTS+".0."+ENVIRONMENT_NAME, "Test"),
					resource.TestCheckResourceAttr(testEnvironmentResourceName, ENVIRONMENTS+".0."+ENVIRONMENT_DESCRIPTION, "Test Env Description"),
					resource.TestCheckResourceAttr(testEnvironmentResourceName, ENVIRONMENTS+".0."+ENVIRONMENT_COLOR, "#5c6bc0"),
				),
			},
		},
	})
}

func TestEnvironmentNotFoundFilter(t *testing.T) {
	const dataSource = `
		data "configcat_environments" "test" {
			name_filter_regex = "invalid"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testEnvironmentResourceName, "id"),
					resource.TestCheckResourceAttr(testEnvironmentResourceName, ENVIRONMENTS+".#", "0"),
				),
			},
		},
	})
}

func TestEnvironmentInvalidGuid(t *testing.T) {
	const dataSource = `
		data "configcat_environments" "test" {
			name_filter_regex = "notfound"
			product_id = "invalidGuid"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      dataSource,
				ExpectError: regexp.MustCompile(`"product_id": invalid GUID`),
			},
		},
	})
}
