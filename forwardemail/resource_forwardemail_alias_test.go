package forwardemail

import (
	"fmt"
	"testing"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccForwardemailAlias_basic(t *testing.T) {
	var alias forwardemail.Alias
	domain := "stark.com"
	name := "james"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccForwardemailProviderFactories,
		CheckDestroy:      testAccCheckForwardemailAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckForwardemailAliasConfig_basic, domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardemailAliasExists("forwardemail_alias.test", &alias),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "domain", domain),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "name", name),
				),
			},
		},
	})
}

func TestAccForwardemailAlias_disable(t *testing.T) {
	var alias forwardemail.Alias
	domain := "stark.com"
	name := "james"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccForwardemailProviderFactories,
		CheckDestroy:      testAccCheckForwardemailAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckForwardemailAliasConfig_basic, domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardemailAliasExists("forwardemail_alias.test", &alias),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "domain", domain),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "name", name),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "is_enabled", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckForwardemailAliasConfig_disabled, domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardemailAliasExists("forwardemail_alias.test", &alias),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "domain", domain),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "name", name),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "is_enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckForwardemailAliasExists(n string, alias *forwardemail.Alias) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccForwardemailProvider.Meta().(*forwardemail.Client)
		foundAlias, err := client.GetAlias(rs.Primary.Attributes["domain"], rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundAlias.Name != rs.Primary.ID {
			return fmt.Errorf("alias not found")
		}

		*alias = *foundAlias

		return nil
	}
}

func testAccCheckForwardemailAliasDestroy(s *terraform.State) error {
	client := testAccForwardemailProvider.Meta().(*forwardemail.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "forwardemail_alias" {
			continue
		}

		if _, err := client.GetAlias(rs.Primary.Attributes["domain"], rs.Primary.ID); err == nil {
			return fmt.Errorf("alias is still there")
		}
	}

	return nil
}

const testAccCheckForwardemailAliasConfig_basic = `
	resource "forwardemail_domain" "test" {
		name = "%s"
	}

	resource "forwardemail_alias" "test" {
	  name   = "%s"
	  domain = forwardemail_domain.test.name
	
	  recipients = ["james@rhodes.com"]
	}
`

const testAccCheckForwardemailAliasConfig_disabled = `
	resource "forwardemail_domain" "test" {
		name = "%s"
	}

	resource "forwardemail_alias" "test" {
	  name   = "%s"
	  domain = forwardemail_domain.test.name
	
	  recipients = ["james@rhodes.com"]
      is_enabled = false 
	}
`
