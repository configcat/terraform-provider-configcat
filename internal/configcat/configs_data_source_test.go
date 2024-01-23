// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigsDataSource(t *testing.T) {
	productId := "08d86d63-2721-4da6-8c06-584521d516bc"
	testResourceName := "data.configcat_configs.test"

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
					resource.TestCheckResourceAttr(testResourceName, CONFIGS+".#", "2"),
				),
			},
		},
	})
}
