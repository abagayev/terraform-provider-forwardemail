package forwardemail

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceForwardemailDomain_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccForwardemailProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceForwardemailAccountConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.forwardemail_account.test", "email"),
				),
			},
		},
	})
}

const testAccCheckDataSourceForwardemailAccountConfig_basic = `
	data forwardemail_account "test" {}
`
