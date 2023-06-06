package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSegmentValid(t *testing.T) {
	const dataSource = `
		data "configcat_segments" "test" {
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
					resource.TestCheckResourceAttrSet("data.configcat_segments.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_segments.test", SEGMENTS+".#", "2"),
				),
			},
		},
	})
}

func TestSegmentValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_segments" "test" {
			name_filter_regex = "Beta users"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const segmentID = "08d9e65b-a9f2-48ec-8423-9b2224771639"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_segments.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_segments.test", SEGMENTS+".#", "1"),
					resource.TestCheckResourceAttr("data.configcat_segments.test", SEGMENTS+".0."+SEGMENT_ID, segmentID),
					resource.TestCheckResourceAttr("data.configcat_segments.test", SEGMENTS+".0."+SEGMENT_NAME, "Beta users"),
					resource.TestCheckResourceAttr("data.configcat_segments.test", SEGMENTS+".0."+SEGMENT_DESCRIPTION, "Beta users segment's description"),
				),
			},
		},
	})
}

func TestSegmentNotFoundFilter(t *testing.T) {
	const dataSource = `
		data "configcat_segments" "test" {
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
					resource.TestCheckResourceAttrSet("data.configcat_segments.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_segments.test", SEGMENTS+".#", "0"),
				),
			},
		},
	})
}

func TestSegmentInvalidGuid(t *testing.T) {
	const dataSource = `
		data "configcat_segments" "test" {
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
