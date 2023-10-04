package forwardemail

import (
	"context"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_mame": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: dataSourceAccountRead,
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)

	account, err := client.GetAccount()
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"plan":         account.Plan,
		"email":        account.Email,
		"full_email":   account.FullEmail,
		"display_mame": account.DisplayName,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(account.Id)

	return nil
}
