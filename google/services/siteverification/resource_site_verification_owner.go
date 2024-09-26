// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package siteverification

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceSiteVerificationOwner() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteVerificationOwnerCreate,
		Read:   resourceSiteVerificationOwnerRead,
		Delete: resourceSiteVerificationOwnerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSiteVerificationOwnerImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The email address of the owner.`,
			},
			"web_resource_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The id of the Web Resource to add this owner to, in the form "webResource/<web-resource-id>".`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSiteVerificationOwnerCreate(d *schema.ResourceData, meta interface{}) error {
	email := d.Get("email").(string)

	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading existing WebResource")

	url, err := tpgresource.ReplaceVars(d, config, "{{SiteVerificationBasePath}}{{web_resource_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	obj, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SiteVerificationWebResource %q", d.Id()))
	}

	log.Printf("[DEBUG] Finished reading WebResource: %#v", obj)

	owners, ok := obj["owners"].([]interface{})
	if !ok {
		return fmt.Errorf("WebResource has no existing owners")
	}
	found := false
	for _, owner := range owners {
		if s, ok := owner.(string); ok && s == email {
			found = true
		}
	}
	if !found {
		owners = append(owners, email)
		obj["owners"] = owners

		log.Printf("[DEBUG] Creating new Owner: %#v", obj)

		headers = make(http.Header)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PUT",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutCreate),
			Headers:   headers,
		})
		if err != nil {
			return fmt.Errorf("Error creating Owner: %s", err)
		}

		log.Printf("[DEBUG] Finished creating Owner %q: %#v", d.Id(), res)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{web_resource_id}}/{{email}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceSiteVerificationOwnerRead(d, meta)
}

func resourceSiteVerificationOwnerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SiteVerificationBasePath}}{{web_resource_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SiteVerificationOwner %q", d.Id()))
	}

	owners, ok := res["owners"].([]interface{})
	if !ok {
		return fmt.Errorf("WebResource has no owners")
	}

	found := false
	email := d.Get("email").(string)
	for _, owner := range owners {
		if s, ok := owner.(string); ok && s == email {
			found = true
		}
	}

	if !found {
		// Owner isn't there any more - remove from the state.
		log.Printf("[DEBUG] Removing SiteVerificationOwner because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	return nil
}

func resourceSiteVerificationOwnerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{SiteVerificationBasePath}}{{web_resource_id}}")
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	log.Printf("[DEBUG] Reading existing WebResource")

	headers := make(http.Header)
	obj, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SiteVerificationWebResource %q", d.Id()))
	}

	log.Printf("[DEBUG] Finished reading WebResource: %#v", obj)

	owners, ok := obj["owners"].([]interface{})
	if !ok {
		return fmt.Errorf("WebResource has no existing owners")
	}
	var updatedOwners []interface{}
	email := d.Get("email").(string)
	for _, owner := range owners {
		if s, ok := owner.(string); ok {
			if s != email {
				updatedOwners = append(updatedOwners, s)
			}
		}
	}
	obj["owners"] = updatedOwners

	headers = make(http.Header)

	log.Printf("[DEBUG] Deleting Owner %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PUT",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Owner")
	}

	log.Printf("[DEBUG] Finished deleting Owner %q: %#v", d.Id(), res)
	return nil
}

func resourceSiteVerificationOwnerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^(?P<web_resource_id>webResource/[^/]+)/(?P<email>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{web_resource_id}}/{{email}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished importing Owner %q", d.Id())

	return []*schema.ResourceData{d}, nil
}
