// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func DataSourceGoogleOrganizations() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"organizations": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directory_customer_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"display_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"lifecycle_state": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"org_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleOrganizationsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	organizations := make([]map[string]interface{}, 0)

	filter := ""
	if v, ok := d.GetOk("filter"); ok {
		filter = v.(string)
	}

	request := config.NewResourceManagerClient(userAgent).Organizations.Search(&cloudresourcemanager.SearchOrganizationsRequest{
		Filter: filter,
	})

	err = request.Pages(context.Background(), func(organizationList *cloudresourcemanager.SearchOrganizationsResponse) error {
		for _, organization := range organizationList.Organizations {
			directoryCustomerId := ""
			if organization.Owner != nil {
				directoryCustomerId = organization.Owner.DirectoryCustomerId
			}

			organizations = append(organizations, map[string]interface{}{
				"directory_customer_id": directoryCustomerId,
				"display_name":          organization.DisplayName,
				"lifecycle_state":       organization.LifecycleState,
				"name":                  organization.Name,
				"org_id":                filepath.Base(organization.Name),
			})
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error retrieving organizations: %s", err)
	}

	if err := d.Set("organizations", organizations); err != nil {
		return fmt.Errorf("Error setting organizations: %s", err)
	}

	if filter == "" {
		filter = "empty_filter"
	}
	d.SetId(filter)

	return nil
}
