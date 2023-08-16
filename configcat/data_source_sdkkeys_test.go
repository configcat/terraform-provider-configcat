package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testSdkKeysDataSourceName = "data.configcat_sdkkeys.test"

func TestSdkKeysValidWithoutSecondary(t *testing.T) {
	const dataSource = `
		data "configcat_sdkkeys" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			environment_id = "08d86d63-272c-4355-8027-4b52787bc1bd"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testSdkKeysDataSourceName, "id"),
					resource.TestCheckResourceAttr(testSdkKeysDataSourceName, "primary", "Y23YCCEnpk2MBlhFIdUWvA/nJmlkmynTE-t3GUMoJjOAQ"),
					resource.TestCheckResourceAttr(testSdkKeysDataSourceName, "secondary", ""),
				),
			},
		},
	})
}

func TestSdkKeysValidWithSecondary(t *testing.T) {
	const dataSource = `
		data "configcat_sdkkeys" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			environment_id = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testSdkKeysDataSourceName, "id"),
					resource.TestCheckResourceAttr(testSdkKeysDataSourceName, "primary", "Y23YCCEnpk2MBlhFIdUWvA/k-nG8bhE10K41sa8C2rYOQ"),
					resource.TestCheckResourceAttr(testSdkKeysDataSourceName, "secondary", "Y23YCCEnpk2MBlhFIdUWvA/q7dBIG43LkuX8NJbCeAzdg"),
				),
			},
		},
	})
}

func TestSdkKeysNotFound(t *testing.T) {
	const dataSource = `
		data "configcat_sdkkeys" "test" {
			config_id = "08d86d63-0000-4b8b-823a-56ddda9da038"
			environment_id = "08d86d63-0000-47cd-8bfc-59608ecb91e2"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      dataSource,
				ExpectError: regexp.MustCompile(`404`),
			},
		},
	})
}

func TestSdkKeysInvalidConfigIdGuid(t *testing.T) {
	const dataSource = `
		data "configcat_sdkkeys" "test" {
			config_id = "invalidGuid"
			environment_id = "08d86d63-0000-47cd-8bfc-59608ecb91e2"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      dataSource,
				ExpectError: regexp.MustCompile(`"config_id": invalid GUID`),
			},
		},
	})
}

func TestSdkKeysInvalidEnvironmentIdGuid(t *testing.T) {
	const dataSource = `
		data "configcat_sdkkeys" "test" {
			config_id = "08d86d63-0000-47cd-8bfc-59608ecb91e2"
			environment_id = "invalidGuid"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      dataSource,
				ExpectError: regexp.MustCompile(`"environment_id": invalid GUID`),
			},
		},
	})
}
