// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package siteverification

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func DataSourceSiteVerificationToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSiteVerificationTokenRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The site identifier. If the type is set to SITE, the identifier is a URL. If the type is
set to INET_DOMAIN, the identifier is a domain name.`,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"INET_DOMAIN", "SITE"}),
				Description:  `The type of resource to be verified, either a domain or a web site. Possible values: ["INET_DOMAIN", "SITE"]`,
			},
			"verification_method": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"ANALYTICS", "DNS_CNAME", "DNS_TXT", "FILE", "META", "TAG_MANAGER"}),
				Description: `The verification method for the Site Verification system to use to verify
this site or domain. Possible values: ["ANALYTICS", "DNS_CNAME", "DNS_TXT", "FILE", "META", "TAG_MANAGER"]`,
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The returned token for use in subsequent verification steps.`,
			},
		},
		UseJSONNumber: true,
	}
}

func dataSourceSiteVerificationTokenRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	site := make(map[string]interface{})
	typeProp, err := expandSiteVerificationTokenType(d.Get("type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("type"); !tpgresource.IsEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		site["type"] = typeProp
	}
	identifierProp, err := expandSiteVerificationTokenIdentifier(d.Get("identifier"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("identifier"); !tpgresource.IsEmptyValue(reflect.ValueOf(identifierProp)) && (ok || !reflect.DeepEqual(v, identifierProp)) {
		site["identifier"] = identifierProp
	}
	obj["site"] = site
	verification_methodProp, err := expandSiteVerificationTokenVerificationMethod(d.Get("verification_method"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("verification_method"); !tpgresource.IsEmptyValue(reflect.ValueOf(verification_methodProp)) && (ok || !reflect.DeepEqual(v, verification_methodProp)) {
		obj["verificationMethod"] = verification_methodProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SiteVerificationBasePath}}token")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Token: %#v", obj)
	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error reading Token: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{identifier}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	if token, ok := res["token"].(string); ok {
		d.Set("token", token)
	}

	log.Printf("[DEBUG] Finished reading Token %q: %#v", d.Id(), res)

	return nil
}

func expandSiteVerificationTokenType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSiteVerificationTokenIdentifier(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSiteVerificationTokenVerificationMethod(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
