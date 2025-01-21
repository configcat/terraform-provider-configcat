package configcat

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"configcat": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(ENV_BASE_PATH); v == "" {
		t.Fatal(ENV_BASE_PATH + " must be set for acceptance tests")
	}
	if v := os.Getenv(ENV_BASIC_AUTH_USERNAME); v == "" {
		t.Fatal(ENV_BASIC_AUTH_USERNAME + " must be set for acceptance tests")
	}
	if v := os.Getenv(ENV_BASIC_AUTH_PASSWORD); v == "" {
		t.Fatal(ENV_BASIC_AUTH_PASSWORD + " must be set for acceptance tests")
	}
}
