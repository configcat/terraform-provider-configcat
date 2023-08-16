package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testTagsDataSourceName = "data.configcat_tags.test"

func TestTagValid(t *testing.T) {
	const dataSource = `
		data "configcat_tags" "test" {
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
					resource.TestCheckResourceAttrSet(testTagsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testTagsDataSourceName, TAGS+".#", "1"),
				),
			},
		},
	})
}

func TestTagValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_tags" "test" {
			name_filter_regex = "Test"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const tagID = "46"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testTagsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testTagsDataSourceName, TAGS+".#", "1"),
					resource.TestCheckResourceAttr(testTagsDataSourceName, TAGS+".0."+TAG_ID, tagID),
					resource.TestCheckResourceAttr(testTagsDataSourceName, TAGS+".0."+TAG_NAME, "Test"),
					resource.TestCheckResourceAttr(testTagsDataSourceName, TAGS+".0."+TAG_COLOR, "panther"),
				),
			},
		},
	})
}

func TestTagNotFoundFilter(t *testing.T) {
	const dataSource = `
		data "configcat_tags" "test" {
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
					resource.TestCheckResourceAttrSet(testTagsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testTagsDataSourceName, TAGS+".#", "0"),
				),
			},
		},
	})
}

func TestTagInvalidGuid(t *testing.T) {
	const dataSource = `
		data "configcat_tags" "test" {
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
