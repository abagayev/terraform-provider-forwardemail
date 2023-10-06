package forwardemail

import (
	"context"
	"github.com/google/go-cmp/cmp"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"has_adult_content_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"has_phishing_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"has_executable_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"has_virus_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"has_recipient_verification": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	name := d.Get("name").(string)

	params := forwardemail.DomainParameters{
		HasAdultContentProtection: toBool(d.Get("has_adult_content_protection")),
		HasPhishingProtection:     toBool(d.Get("has_phishing_protection")),
		HasExecutableProtection:   toBool(d.Get("has_executable_protection")),
		HasVirusProtection:        toBool(d.Get("has_virus_protection")),
		HasRecipientVerification:  toBool(d.Get("has_recipient_verification")),
	}

	domain, err := client.CreateDomain(name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"has_adult_content_protection": domain.HasAdultContentProtection,
		"has_phishing_protection":      domain.HasPhishingProtection,
		"has_executable_protection":    domain.HasExecutableProtection,
		"has_virus_protection":         domain.HasVirusProtection,
		"has_recipient_verification":   domain.HasRecipientVerification,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)

	return nil
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	name := d.Id()

	domain, err := client.GetDomain(name)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"has_adult_content_protection": domain.HasAdultContentProtection,
		"has_phishing_protection":      domain.HasPhishingProtection,
		"has_executable_protection":    domain.HasExecutableProtection,
		"has_virus_protection":         domain.HasVirusProtection,
		"has_recipient_verification":   domain.HasRecipientVerification,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	name := d.Id()

	params := forwardemail.DomainParameters{}
	params.HasAdultContentProtection = toBool(toChange(d.GetChange("has_adult_content_protection")))
	params.HasPhishingProtection = toBool(toChange(d.GetChange("has_phishing_protection")))
	params.HasExecutableProtection = toBool(toChange(d.GetChange("has_executable_protection")))
	params.HasVirusProtection = toBool(toChange(d.GetChange("has_virus_protection")))
	params.HasRecipientVerification = toBool(toChange(d.GetChange("has_recipient_verification")))

	_, err := client.UpdateDomain(name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	name := d.Id()

	err := client.DeleteDomain(name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// toBool returns a pointer to the bool value passed in.
func toBool(v interface{}) *bool {
	if b, ok := v.(bool); ok {
		return &b
	}

	return nil
}

func toChange(p, c interface{}) interface{} {
	if cmp.Equal(p, c) {
		return nil
	}

	return c
}
