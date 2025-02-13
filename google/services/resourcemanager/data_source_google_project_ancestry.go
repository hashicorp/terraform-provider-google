// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func DataSourceGoogleProjectAncestry() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleProjectAncestryRead,
		Schema: map[string]*schema.Schema{
			"ancestors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func datasourceGoogleProjectAncestryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ancestry: %s", err)
	}

	request := &cloudresourcemanager.GetAncestryRequest{}
	response, err := config.NewResourceManagerClient(userAgent).Projects.GetAncestry(project, request).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("Error retrieving project ancestry: %s", err)
	}

	ancestors := make([]map[string]interface{}, 0)
	var orgID string
	var parentID string
	var parentType string

	for _, a := range response.Ancestor {
		if a.ResourceId == nil {
			continue
		}

		ancestorData := map[string]interface{}{
			"id":   a.ResourceId.Id,
			"type": a.ResourceId.Type,
		}

		ancestors = append(ancestors, ancestorData)

		if a.ResourceId.Type == "organization" {
			orgID = a.ResourceId.Id
		}
	}

	if err := d.Set("ancestors", ancestors); err != nil {
		return fmt.Errorf("Error setting ancestors: %s", err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		return fmt.Errorf("Error setting org_id: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if len(ancestors) > 1 {
		parent := ancestors[1]
		if id, ok := parent["id"].(string); ok {
			parentID = id
		}
		if pType, ok := parent["type"].(string); ok {
			parentType = pType
		}

		if err := d.Set("parent_id", parentID); err != nil {
			return fmt.Errorf("Error setting parent_id: %s", err)
		}

		if err := d.Set("parent_type", parentType); err != nil {
			return fmt.Errorf("Error setting parent_type: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("projects/%s", project))

	return nil
}
