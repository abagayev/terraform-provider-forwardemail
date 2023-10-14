package forwardemail

import (
	"context"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"has_recipient_verification": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"recipients": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"labels": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
		CreateContext: resourceAliasCreate,
		ReadContext:   resourceAliasRead,
		UpdateContext: resourceAliasUpdate,
		DeleteContext: resourceAliasDelete,
	}
}

func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	params := forwardemail.AliasParameters{
		HasRecipientVerification: toBool(d.Get("has_recipient_verification")),
		IsEnabled:                toBool(d.Get("is_enabled")),
	}

	alias, err := client.CreateAlias(domain, name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"domain":                     alias.Domain.Name,
		"has_recipient_verification": alias.HasRecipientVerification,
		"is_enabled":                 alias.IsEnabled,
		"recipients":                 alias.Recipients,
		"labels":                     alias.Labels,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)

	return nil
}

func resourceAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	alias, err := client.GetAlias(domain, name)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"domain":                     alias.Domain.Name,
		"has_recipient_verification": alias.HasRecipientVerification,
		"is_enabled":                 alias.IsEnabled,
		"recipients":                 alias.Recipients,
		"labels":                     alias.Labels,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	domain := d.Get("domain").(string)
	name := d.Id()

	params := forwardemail.AliasParameters{}
	params.HasRecipientVerification = toBool(toChange(d.GetChange("has_recipient_verification")))
	params.IsEnabled = toBool(toChange(d.GetChange("is_enabled")))
	params.Recipients = toSliceOfStrings(toChanges(d.GetChange("recipients")))
	params.Labels = toSliceOfStrings(toChanges(d.GetChange("labels")))

	_, err := client.UpdateAlias(domain, name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	domain := d.Get("domain").(string)
	name := d.Id()

	err := client.DeleteAlias(domain, name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// toSliceOfStrings converts slice of interfaces into pointer to slice of strings.
func toSliceOfStrings(vs []interface{}) *[]string {
	var stringSlice []string
	for _, v := range vs {
		if str, ok := v.(string); ok {
			stringSlice = append(stringSlice, str)
		}
	}

	return &stringSlice
}

// toChanges converts interface into slice of interfaces.
func toChanges(p, c interface{}) []interface{} {
	if cmp.Equal(p, c) {
		return nil
	}

	switch v := c.(type) {
	case []interface{}:
		return v
	}

	return nil
}
