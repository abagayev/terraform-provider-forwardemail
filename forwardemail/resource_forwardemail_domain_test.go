package forwardemail

import (
	"fmt"
	"testing"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccForwardemailDomain_basic(t *testing.T) {
	var domain forwardemail.Domain
	name := "stark.com"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccForwardemailProviderFactories,
		CheckDestroy:      testAccCheckForwardemailDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardemailDomainConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardemailDomainExists("forwardemail_domain.test", &domain),
					resource.TestCheckResourceAttr("forwardemail_domain.test", "name", name),
				),
			},
		},
	})
}

func testAccCheckForwardemailDomainExists(n string, domain *forwardemail.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccForwardemailProvider.Meta().(*forwardemail.Client)
		foundDomain, err := client.GetDomain(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundDomain.Name != rs.Primary.ID {
			return fmt.Errorf("domain not found")
		}

		*domain = *foundDomain

		return nil
	}
}

func testAccCheckForwardemailDomainDestroy(s *terraform.State) error {
	client := testAccForwardemailProvider.Meta().(*forwardemail.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "forwardemail_domain" {
			continue
		}

		if _, err := client.GetDomain(rs.Primary.ID); err == nil {
			return fmt.Errorf("domain is still there")
		}
	}

	return nil
}

func testAccCheckForwardemailDomainConfig_basic(name string) string {
	return fmt.Sprintf(`
		resource "forwardemail_domain" "test" {
			name = "%s"
		}
	`, name)
}
