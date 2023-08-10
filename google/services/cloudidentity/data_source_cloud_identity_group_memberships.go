// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudidentity/v1"
)

func DataSourceGoogleCloudIdentityGroupMemberships() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCloudIdentityGroupMembership().Schema)

	return &schema.Resource{
		Read: dataSourceGoogleCloudIdentityGroupMembershipsRead,

		Schema: map[string]*schema.Schema{
			"memberships": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `List of Cloud Identity group memberships.`,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
			"group": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name of the Group to get memberships from.`,
			},
		},
	}
}

func dataSourceGoogleCloudIdentityGroupMembershipsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	result := []map[string]interface{}{}
	membershipsCall := config.NewCloudIdentityClient(userAgent).Groups.Memberships.List(d.Get("group").(string)).View("FULL")
	if config.UserProjectOverride {
		billingProject := ""
		// err may be nil - project isn't required for this resource
		if project, err := tpgresource.GetProject(d, config); err == nil {
			billingProject = project
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}

		if billingProject != "" {
			membershipsCall.Header().Set("X-Goog-User-Project", billingProject)
		}
	}

	err = membershipsCall.Pages(config.Context, func(resp *cloudidentity.ListMembershipsResponse) error {
		for _, member := range resp.Memberships {
			result = append(result, map[string]interface{}{
				"name":                 member.Name,
				"type":                 member.Type,
				"roles":                flattenCloudIdentityGroupMembershipsRoles(member.Roles),
				"preferred_member_key": flattenCloudIdentityGroupsEntityKey(member.PreferredMemberKey),
			})
		}

		return nil
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("CloudIdentityGroupMemberships %q", d.Id()))
	}

	if err := d.Set("memberships", result); err != nil {
		return fmt.Errorf("Error setting memberships: %s", err)
	}
	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenCloudIdentityGroupMembershipsRoles(roles []*cloudidentity.MembershipRole) []interface{} {
	transformed := []interface{}{}

	for _, role := range roles {
		transformed = append(transformed, map[string]interface{}{
			"name": role.Name,
		})
	}
	return transformed
}
