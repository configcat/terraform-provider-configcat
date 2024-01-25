// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigsDataSource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "data.configcat_configs.test"

	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const configV2ID = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"

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
					resource.TestCheckResourceAttr(testResourceName, Configs+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+ConfigId, configID),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+Name, "Main Config"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+Description, "Main Config Description"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+Order, "0"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".1."+ConfigId, configV2ID),
					resource.TestCheckResourceAttr(testResourceName, Configs+".1."+Name, "Main Config V2"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".1."+Description, "Main Config V2 Description"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".1."+Order, "1"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("Main Config V2"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+ConfigId, configV2ID),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+Name, "Main Config V2"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+Description, "Main Config V2 Description"),
					resource.TestCheckResourceAttr(testResourceName, Configs+".0."+Order, "1"),
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
					resource.TestCheckResourceAttr(testResourceName, Configs+".#", "0"),
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
