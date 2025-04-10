// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/siteverification/WebResource.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

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

func ResourceSiteVerificationWebResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteVerificationWebResourceCreate,
		Read:   resourceSiteVerificationWebResourceRead,
		Delete: resourceSiteVerificationWebResourceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSiteVerificationWebResourceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"site": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: `Container for the address and type of a site for which a verification token will be verified.`,
				MaxItems:    1,
				Elem: &schema.Resource{
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
							Description:  `The type of resource to be verified. Possible values: ["INET_DOMAIN", "SITE"]`,
						},
					},
				},
			},
			"verification_method": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"ANALYTICS", "DNS_CNAME", "DNS_TXT", "FILE", "META", "TAG_MANAGER"}),
				Description: `The verification method for the Site Verification system to use to verify
this site or domain. Possible values: ["ANALYTICS", "DNS_CNAME", "DNS_TXT", "FILE", "META", "TAG_MANAGER"]`,
			},
			"owners": {
				Type:     schema.TypeList,
				Computed: true,
				Description: `The email addresses of all direct, verified owners of this exact property. Indirect owners —
for example verified owners of the containing domain—are not included in this list.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"web_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The string used to identify this web resource.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSiteVerificationWebResourceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	siteProp, err := expandSiteVerificationWebResourceSite(d.Get("site"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("site"); !tpgresource.IsEmptyValue(reflect.ValueOf(siteProp)) && (ok || !reflect.DeepEqual(v, siteProp)) {
		obj["site"] = siteProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SiteVerificationBasePath}}webResource?verificationMethod={{verification_method}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new WebResource: %#v", obj)
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
		Config:               config,
		Method:               "POST",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutCreate),
		Headers:              headers,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsSiteVerificationRetryableError},
	})
	if err != nil {
		return fmt.Errorf("Error creating WebResource: %s", err)
	}
	// Set computed resource properties from create API response so that they're available on the subsequent Read
	// call.
	if err := d.Set("web_resource_id", flattenSiteVerificationWebResourceWebResourceId(res["id"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "web_resource_id": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "webResource/{{web_resource_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating WebResource %q: %#v", d.Id(), res)

	return resourceSiteVerificationWebResourceRead(d, meta)
}

func resourceSiteVerificationWebResourceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SiteVerificationBasePath}}webResource/{{web_resource_id}}")
	if err != nil {
		return err
	}

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
		Config:               config,
		Method:               "GET",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Headers:              headers,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsSiteVerificationRetryableError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SiteVerificationWebResource %q", d.Id()))
	}

	if err := d.Set("web_resource_id", flattenSiteVerificationWebResourceWebResourceId(res["id"], d, config)); err != nil {
		return fmt.Errorf("Error reading WebResource: %s", err)
	}
	if err := d.Set("site", flattenSiteVerificationWebResourceSite(res["site"], d, config)); err != nil {
		return fmt.Errorf("Error reading WebResource: %s", err)
	}
	if err := d.Set("owners", flattenSiteVerificationWebResourceOwners(res["owners"], d, config)); err != nil {
		return fmt.Errorf("Error reading WebResource: %s", err)
	}

	return nil
}

func resourceSiteVerificationWebResourceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{SiteVerificationBasePath}}webResource/{{web_resource_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting WebResource %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "DELETE",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutDelete),
		Headers:              headers,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsSiteVerificationRetryableError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "WebResource")
	}

	log.Printf("[DEBUG] Finished deleting WebResource %q: %#v", d.Id(), res)
	return nil
}

func resourceSiteVerificationWebResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^webResource/(?P<web_resource_id>[^/]+)$",
		"^(?P<web_resource_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "webResource/{{web_resource_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenSiteVerificationWebResourceWebResourceId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSiteVerificationWebResourceSite(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["type"] =
		flattenSiteVerificationWebResourceSiteType(original["type"], d, config)
	transformed["identifier"] =
		flattenSiteVerificationWebResourceSiteIdentifier(original["identifier"], d, config)
	return []interface{}{transformed}
}
func flattenSiteVerificationWebResourceSiteType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSiteVerificationWebResourceSiteIdentifier(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSiteVerificationWebResourceOwners(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandSiteVerificationWebResourceSite(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedType, err := expandSiteVerificationWebResourceSiteType(original["type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["type"] = transformedType
	}

	transformedIdentifier, err := expandSiteVerificationWebResourceSiteIdentifier(original["identifier"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIdentifier); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["identifier"] = transformedIdentifier
	}

	return transformed, nil
}

func expandSiteVerificationWebResourceSiteType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSiteVerificationWebResourceSiteIdentifier(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
