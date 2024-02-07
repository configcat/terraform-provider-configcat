package configcat

import (
	"path"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestResourceSettingBoolean(t *testing.T) {
	testAccSettingResource(t, "boolean")
}

func TestResourceSettingString(t *testing.T) {
	testAccSettingResource(t, "string")
}

func TestResourceSettingInt(t *testing.T) {
	testAccSettingResource(t, "int")
}

func TestResourceSettingDouble(t *testing.T) {
	testAccSettingResource(t, "double")
}

func testAccSettingResource(t *testing.T, settingType string) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const testResourceName = "configcat_setting.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKey" + settingType),
					"name":         config.StringVariable("Resource name"),
					"setting_type": config.StringVariable(settingType),
					"order":        config.IntegerVariable(1),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingKey, "SettingKey"+settingType),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name"),
					resource.TestCheckResourceAttr(testResourceName, SettingType, settingType),
					resource.TestCheckResourceAttr(testResourceName, Order, "1"),
					resource.TestCheckResourceAttr(testResourceName, SettingHint, ""),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKey" + settingType),
					"name":         config.StringVariable("Resource name updated"),
					"hint":         config.StringVariable("Hint"),
					"setting_type": config.StringVariable(settingType),
					"order":        config.IntegerVariable(10),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingKey, "SettingKey"+settingType),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name updated"),
					resource.TestCheckResourceAttr(testResourceName, SettingHint, "Hint"),
					resource.TestCheckResourceAttr(testResourceName, SettingType, settingType),
					resource.TestCheckResourceAttr(testResourceName, Order, "10"),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKey" + settingType),
					"name":         config.StringVariable("Resource name updated"),
					"hint":         config.StringVariable("Hint"),
					"setting_type": config.StringVariable(settingType),
					"order":        config.IntegerVariable(10),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSettingResourceInvalidSettingType(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKeyInvalid"),
					"name":         config.StringVariable("Resource name"),
					"setting_type": config.StringVariable("invalid"),
					"order":        config.IntegerVariable(1),
				},
				ExpectError: regexp.MustCompile("invalid value 'invalid' for SettingType"),
			},
		},
	})
}

func TestAccSettingResourceDuplicatedKey(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("isAwesomeFeatureEnabled"),
					"name":         config.StringVariable("Resource name"),
					"setting_type": config.StringVariable("boolean"),
					"order":        config.IntegerVariable(1),
				},
				ExpectError: regexp.MustCompile("This key is already in use. Please, choose another"),
			},
		},
	})
}

func TestAccSettingResourceSettingTypeRecreate(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const testResourceName = "configcat_setting.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKeyToRecreate1"),
					"name":         config.StringVariable("Resource name"),
					"setting_type": config.StringVariable("boolean"),
					"order":        config.IntegerVariable(1),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKeyToRecreate1"),
					"name":         config.StringVariable("Resource name"),
					"setting_type": config.StringVariable("string"),
					"order":        config.IntegerVariable(1),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testResourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAccSettingResourceKeyRecreate(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const testResourceName = "configcat_setting.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKeyToRecreate2"),
					"name":         config.StringVariable("Resource name"),
					"setting_type": config.StringVariable("boolean"),
					"order":        config.IntegerVariable(1),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccSettingResource", "main.tf")),
				ConfigVariables: config.Variables{
					"config_id":    config.StringVariable(configId),
					"key":          config.StringVariable("SettingKeyToRecreateUpdated"),
					"name":         config.StringVariable("Resource name"),
					"setting_type": config.StringVariable("boolean"),
					"order":        config.IntegerVariable(1),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testResourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}
