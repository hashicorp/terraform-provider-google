// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/serviceusage"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFolderServiceIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceFolderServiceIdentityCreate,
		Read:   resourceFolderServiceIdentityRead,
		Delete: resourceFolderServiceIdentityDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"folder": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Identity of the Google managed service account in the form 'serviceAccount:{email}'. This value is often used to refer to the service account in order to grant IAM permissions.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceFolderServiceIdentityCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ServiceUsageBasePath}}folders/{{folder}}/services/{{service}}:generateServiceIdentity")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating Folder Service Identity: %s", err)
	}

	var opRes map[string]interface{}
	err = serviceusage.ServiceUsageOperationWaitTimeWithResponse(
		config, res, &opRes, billingProject, "Creating Folder Service Identity", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished creating Folder Service Identity %q: %#v", d.Id(), res)

	id, err := tpgresource.ReplaceVars(d, config, "folders/{{folder}}/services/{{service}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// This API may not return the service identity's details, even if the relevant
	// Google API is configured for service identities.
	if emailVal, ok := opRes["email"]; ok {
		email, ok := emailVal.(string)
		if !ok {
			return fmt.Errorf("unexpected type for email: got %T, want string", email)
		}
		if err := d.Set("email", email); err != nil {
			return fmt.Errorf("Error setting email: %s", err)
		}
		if err := d.Set("member", "serviceAccount:"+email); err != nil {
			return fmt.Errorf("Error setting member: %s", err)
		}
	}
	return nil
}

// There is no read endpoint for this API.
func resourceFolderServiceIdentityRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// There is no delete endpoint for this API.
func resourceFolderServiceIdentityDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
