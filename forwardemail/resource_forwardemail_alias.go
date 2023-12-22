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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Fully qualified domain name (FQDN).",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alias name.",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable to disable this alias.",
			},
			"has_recipient_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable to require recipients to click an email verification link for emails to flow through.",
			},
			"recipients": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of recipients as valid email addresses, fully-qualified domain names (FQDN), IP addresses, or webhook URL's.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alias description.",
			},
			"labels": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of labels.",
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
		Recipients:               toSliceOfStrings(toChanges(nil, d.Get("recipients"))),
		Labels:                   toSliceOfStrings(toChanges(nil, d.Get("labels"))),
		Description:              d.Get("description").(string),
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
		"description":                alias.Description,
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
		"description":                alias.Description,
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

	// N.B.: we can't use d.GetChange because Forward Email API is working not as expected.
	// So instead of passing only changed parameters we need to pass all of them.
	params := forwardemail.AliasParameters{}
	params.HasRecipientVerification = toBool(toChange(nil, d.Get("has_recipient_verification")))
	params.IsEnabled = toBool(toChange(nil, d.Get("is_enabled")))
	params.Recipients = toSliceOfStrings(toChanges(nil, d.Get("recipients")))
	params.Labels = toSliceOfStrings(toChanges(nil, d.Get("labels")))
	params.Description = d.Get("description").(string)

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
