package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSegmentsDataSource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "data.configcat_segments.test"

	const segment1ID = "08d9e65b-a9f2-48ec-8423-9b2224771639"
	const segment2ID = "08d9e65b-a9f2-48ec-8423-9b2224771639"
	const listAttribute = Segments

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable(productId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SegmentId, segment1ID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "Beta users"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Description, "Beta users segment's description"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SegmentId, segment2ID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "Beta users"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Description, "Beta users segment's description"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("Beta users"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SegmentId, segment1ID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "Beta users"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Description, "Beta users segment's description"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("notfound"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable("invalidguid"),
				},
				ExpectError: regexp.MustCompile(`Attribute product_id value must be a valid GUID, got: invalidguid`),
			},
		},
	})
}
