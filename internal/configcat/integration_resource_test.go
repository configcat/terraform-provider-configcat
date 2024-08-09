package configcat

import (
	"path"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationAmplitudeResource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const config1Id = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const config2Id = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const integrationType = "amplitude"
	const environmentId = "08d8becf-d4d9-4c66-8b48-6ac74cd95fba"
	const testResourceName = "configcat_integration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
				},
				ExpectError: regexp.MustCompile("Must set a configuration value for the name attribute"),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
				},
				ExpectError: regexp.MustCompile("apiKey parameter is required"),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey": config.StringVariable("apiKey"),
					}),
				},
				ExpectError: regexp.MustCompile("secretKey parameter is required"),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey"),
						"secretKey": config.StringVariable("secretKey"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, integrationType+"_integ"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationType, integrationType),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".apiKey", "apiKey"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".secretKey", "secretKey"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationEnvironments+".#", "0"),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey"),
						"secretKey": config.StringVariable("secretKey"),
					}),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ_rename"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey2"),
						"secretKey": config.StringVariable("secretKey2"),
					}),
					"configs": config.ListVariable(config.StringVariable(config1Id), config.StringVariable(config2Id)),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, integrationType+"_integ_rename"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationType, integrationType),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".apiKey", "apiKey2"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".secretKey", "secretKey2"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".0", config1Id),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".1", config2Id),
					resource.TestCheckResourceAttr(testResourceName, IntegrationEnvironments+".#", "0"),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ_rename"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey2"),
						"secretKey": config.StringVariable("secretKey2"),
					}),
					"configs": config.ListVariable(config.StringVariable(config1Id), config.StringVariable(config2Id)),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey"),
						"secretKey": config.StringVariable("secretKey"),
					}),
					"environments": config.ListVariable(),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, integrationType+"_integ"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationType, integrationType),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".apiKey", "apiKey"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".secretKey", "secretKey"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationEnvironments+".#", "0"),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey"),
						"secretKey": config.StringVariable("secretKey"),
					}),
					"environments": config.ListVariable(),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey"),
						"secretKey": config.StringVariable("secretKey"),
					}),
					"environments": config.ListVariable(config.StringVariable(environmentId)),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, integrationType+"_integ"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationType, integrationType),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".apiKey", "apiKey"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".secretKey", "secretKey"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationEnvironments+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationEnvironments+".0", environmentId),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apiKey":    config.StringVariable("apiKey"),
						"secretKey": config.StringVariable("secretKey"),
					}),
					"environments": config.ListVariable(config.StringVariable(environmentId)),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIntegrationDatadogResource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const integrationType = "dataDog"
	const testResourceName = "configcat_integration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
				},
				ExpectError: regexp.MustCompile("Must set a configuration value for the name attribute"),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
				},
				ExpectError: regexp.MustCompile("apikey parameter is required"),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apikey": config.StringVariable("apikey"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, integrationType+"_integ"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationType, integrationType),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".apikey", "apikey"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationEnvironments+".#", "0"),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apikey": config.StringVariable("apikey"),
					}),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ2"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apikey": config.StringVariable("apikey2"),
						"site":   config.StringVariable("Us1Fed"),
					}),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, integrationType+"_integ2"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationType, integrationType),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".apikey", "apikey2"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationParameters+".site", "Us1Fed"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationConfigs+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, IntegrationEnvironments+".#", "0"),
				),
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ2"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apikey": config.StringVariable("apikey2"),
						"site":   config.StringVariable("Us1Fed"),
					}),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.StaticFile(path.Join("testdata", "TestAccIntegrationResource", "main.tf")),
				ConfigVariables: config.Variables{
					"product_id":       config.StringVariable(productId),
					"integration_type": config.StringVariable(integrationType),
					"name":             config.StringVariable(integrationType + "_integ2"),
					"parameters": config.MapVariable(map[string]config.Variable{
						"apikey": config.StringVariable("apikey2"),
						"site":   config.StringVariable("invalid"),
					}),
				},
				ExpectError: regexp.MustCompile("'site' is invalid"),
			},
		},
	})
}
