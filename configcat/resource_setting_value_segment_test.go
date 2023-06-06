package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceSettingValueBoolSegmentRules(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const segmentBetaUsersID = "08d9e65b-a9f2-48ec-8423-9b2224771639"
	const segmentOldApplicationVersionID = "08d9e65b-c95e-4cb1-8a24-a326d1da676f"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueBoolSegmentRules"
			name = "testBoolWithSegment"
		}
		
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
		}
	`
	const settingValueResourceRule1 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueBoolSegmentRules"
			name = "testBoolWithSegment"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
			rollout_rules {
				segment_comparator = "isIn"
				segment_id = "` + segmentBetaUsersID + `"
				value = "true"
			}
		}
	`

	const settingValueResourceRule2 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueBoolSegmentRules"
			name = "testBoolWithSegment"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
			rollout_rules {
				segment_comparator = "isNotIn"
				segment_id = "` + segmentOldApplicationVersionID + `"
				value = "true"
			}
			rollout_rules {
				segment_comparator = "isIn"
				segment_id = "` + segmentOldApplicationVersionID + `"
				value = "false"
			}
		}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".#", "0"),
				),
			},
			{
				Config: settingValueResourceRule1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".#", "1"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARATOR, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_VALUE, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_SEGMENT_COMPARATOR, "isIn"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_SEGMENT_ID, segmentBetaUsersID),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_VALUE, "true"),
				),
			},
			{
				Config: settingValueResourceRule2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".#", "2"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARATOR, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_VALUE, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_SEGMENT_COMPARATOR, "isNotIn"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_SEGMENT_ID, segmentOldApplicationVersionID),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_VALUE, "true"),

					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARATOR, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARISON_VALUE, ""),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_SEGMENT_COMPARATOR, "isIn"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_SEGMENT_ID, segmentOldApplicationVersionID),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_VALUE, "false"),
				),
			},
		},
	})
}
