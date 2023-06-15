// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func DataSourceGoogleOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"organization"},
			},
			"organization": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"domain"},
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	var organization *cloudresourcemanager.Organization
	if v, ok := d.GetOk("domain"); ok {
		filter := fmt.Sprintf("domain=%s", v.(string))
		var resp *cloudresourcemanager.SearchOrganizationsResponse
		err := transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (err error) {
				resp, err = config.NewResourceManagerClient(userAgent).Organizations.Search(&cloudresourcemanager.SearchOrganizationsRequest{
					Filter: filter,
				}).Do()
				return err
			},
			Timeout: d.Timeout(schema.TimeoutRead),
		})
		if err != nil {
			return fmt.Errorf("Error reading organization: %s", err)
		}

		if len(resp.Organizations) == 0 {
			return fmt.Errorf("Organization not found: %s", v)
		}

		if len(resp.Organizations) > 1 {
			// Attempt to find an exact domain match
			for _, org := range resp.Organizations {
				if org.DisplayName == v.(string) {
					organization = org
					break
				}
			}
			if organization == nil {
				return fmt.Errorf("Received multiple organizations in the response, but could not find an exact domain match.")
			}
		} else {
			organization = resp.Organizations[0]
		}

	} else if v, ok := d.GetOk("organization"); ok {
		var resp *cloudresourcemanager.Organization
		err := transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (err error) {
				resp, err = config.NewResourceManagerClient(userAgent).Organizations.Get(canonicalOrganizationName(v.(string))).Do()
				return err
			},
			Timeout: d.Timeout(schema.TimeoutRead),
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Organization Not Found : %s", v))
		}

		organization = resp
	} else {
		return fmt.Errorf("one of domain or organization must be set")
	}

	d.SetId(organization.Name)
	if err := d.Set("name", organization.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("org_id", tpgresource.GetResourceNameFromSelfLink(organization.Name)); err != nil {
		return fmt.Errorf("Error setting org_id: %s", err)
	}
	if err := d.Set("domain", organization.DisplayName); err != nil {
		return fmt.Errorf("Error setting domain: %s", err)
	}
	if err := d.Set("create_time", organization.CreationTime); err != nil {
		return fmt.Errorf("Error setting create_time: %s", err)
	}
	if err := d.Set("lifecycle_state", organization.LifecycleState); err != nil {
		return fmt.Errorf("Error setting lifecycle_state: %s", err)
	}
	if organization.Owner != nil {
		if err := d.Set("directory_customer_id", organization.Owner.DirectoryCustomerId); err != nil {
			return fmt.Errorf("Error setting directory_customer_id: %s", err)
		}
	}

	return nil
}

func canonicalOrganizationName(ba string) string {
	if strings.HasPrefix(ba, "organizations/") {
		return ba
	}

	return "organizations/" + ba
}
