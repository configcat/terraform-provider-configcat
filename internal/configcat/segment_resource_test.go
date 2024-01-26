package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSegmentResource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "configcat_segment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":           config.StringVariable(productId),
					"name":                 config.StringVariable("Resource name"),
					"comparison_attribute": config.StringVariable("email"),
					"comparator":           config.StringVariable("sensitiveIsOneOf"),
					"comparison_value":     config.StringVariable("a@b.com,c@d.com"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name"),
					resource.TestCheckNoResourceAttr(testResourceName, Description),
					resource.TestCheckResourceAttr(testResourceName, SegmentComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, SegmentComparator, "sensitiveIsOneOf"),
					resource.TestCheckResourceAttr(testResourceName, SegmentComparisonValue, "a@b.com,c@d.com"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":           config.StringVariable(productId),
					"name":                 config.StringVariable("Resource name updated"),
					"description":          config.StringVariable("Resource description"),
					"comparison_attribute": config.StringVariable("version"),
					"comparator":           config.StringVariable("semVerLess"),
					"comparison_value":     config.StringVariable("4.0.0"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name updated"),
					resource.TestCheckResourceAttr(testResourceName, Description, "Resource description"),
					resource.TestCheckResourceAttr(testResourceName, SegmentComparisonAttribute, "version"),
					resource.TestCheckResourceAttr(testResourceName, SegmentComparator, "semVerLess"),
					resource.TestCheckResourceAttr(testResourceName, SegmentComparisonValue, "4.0.0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":           config.StringVariable(productId),
					"name":                 config.StringVariable("Resource name updated"),
					"description":          config.StringVariable("Resource description"),
					"comparison_attribute": config.StringVariable("version"),
					"comparator":           config.StringVariable("semVerLess"),
					"comparison_value":     config.StringVariable("4.0.0"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":           config.StringVariable(productId),
					"name":                 config.StringVariable("Resource name updated"),
					"description":          config.StringVariable("Resource description"),
					"comparison_attribute": config.StringVariable("version"),
					"comparator":           config.StringVariable("invalid"),
					"comparison_value":     config.StringVariable("4.0.0"),
				},
				ExpectError: regexp.MustCompile(`invalid value 'invalid' for RolloutRuleComparator`),
			},
		},
	})
}
