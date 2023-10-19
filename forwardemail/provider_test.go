package forwardemail

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jaswdr/faker"
)

var fake faker.Faker
var testAccForwardemailProvider *schema.Provider
var testAccForwardemailProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	fake = faker.New()

	testAccForwardemailProvider = Provider()
	testAccForwardemailProviderFactories = map[string]func() (*schema.Provider, error){
		"forwardemail": func() (*schema.Provider, error) {
			return testAccForwardemailProvider, nil
		},
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

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("FORWARDEMAIL_API_KEY"); v == "" {
		t.Fatal("FORWARDEMAIL_API_KEY must be set for acceptance tests")
	}
}
