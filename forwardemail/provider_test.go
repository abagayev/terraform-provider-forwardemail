package forwardemail

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"forwardemail": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	if Provider() == nil {
		t.Fatal("provider is nil")
	}
}

//func testAccPreCheck(t *testing.T) {
//	if v := os.Getenv("FORWARDEMAIL_API_KEY"); v == "" {
//		t.Fatal("FORWARDEMAIL_API_KEY must be set for acceptance tests")
//	}
//}
