// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vmwareengine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceVmwareengineNsxCredentials() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmwareengineNsxCredentialsRead,
		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The resource name of the private cloud which contains NSX.
Resource names are schemeless URIs that follow the conventions in https://cloud.google.com/apis/design/resource_names.
For example: projects/my-project/locations/us-west1-a/privateClouds/my-cloud`,
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Initial username.`,
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Initial password.`,
			},
		},
	}
}

func dataSourceVmwareengineNsxCredentialsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{VmwareengineBasePath}}{{parent}}:showNsxCredentials")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              billingProject,
		RawURL:               url,
		UserAgent:            userAgent,
		ErrorAbortPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429QuotaError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("VmwareengineNsxCredentials %q", d.Id()))
	}

	if err := d.Set("username", flattenVmwareengineNsxCredentailsUsername(res["username"], d, config)); err != nil {
		return fmt.Errorf("Error reading NsxCredentails: %s", err)
	}
	if err := d.Set("password", flattenVmwareengineNsxCredentailsPassword(res["password"], d, config)); err != nil {
		return fmt.Errorf("Error reading NsxCredentails: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}:nsx-credentials")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}

func flattenVmwareengineNsxCredentailsUsername(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenVmwareengineNsxCredentailsPassword(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
