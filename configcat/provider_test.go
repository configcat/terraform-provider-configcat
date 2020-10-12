package configcat

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"configcat": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CONFIGCAT_BASIC_AUTH_USERNAME"); v == "" {
		t.Fatal("CONFIGCAT_BASIC_AUTH_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("CONFIGCAT_BASIC_AUTH_PASSWORD"); v == "" {
		t.Fatal("CONFIGCAT_BASIC_AUTH_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("CONFIGCAT_BASE_PATH"); v == "" {
		t.Fatal("CONFIGCAT_BASE_PATH must be set for acceptance tests")
	}
}
