// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudidentity/v1"
)

func DataSourceGoogleCloudIdentityGroupTransitiveMemberships() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceGoogleCloudIdentityGroupTransitiveMembershipsRead,

		// We don't reuse schemas from google_cloud_identity_group_membership because data returned about
		// transative memberships is structured differently, with information like expiry missing.
		Schema: map[string]*schema.Schema{
			"group": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name of the Group to get memberships from.`,
			},
			"memberships": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `List of Cloud Identity group memberships.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roles": {
							Type: schema.TypeSet,
							// Default schema.HashSchema is used.
							Computed:    true,
							Description: `The membership role details`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the TransitiveMembershipRole. Possible values: ["OWNER", "MANAGER", "MEMBER"]`,
									},
								},
							},
						},
						"preferred_member_key": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `EntityKey of the member. Entity key has an id and a namespace. In case of discussion forums, the id will be an email address without a namespace.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `The ID of the entity.

For Google-managed entities, the id must be the email address of an existing
group or user.

For external-identity-mapped entities, the id must be a string conforming
to the Identity Source's requirements.

Must be unique within a namespace.`,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `The namespace in which the entity exists.

If not specified, the EntityKey represents a Google-managed entity
such as a Google user or a Google Group.

If specified, the EntityKey represents an external-identity-mapped group.
The namespace must correspond to an identity source created in Admin Console
and must be in the form of 'identitysources/{identity_source_id}'.`,
									},
								},
							},
						},
						"member": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Resource name for this member.`,
						},
						"relation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The relation between the group and the transitive member. The value can be DIRECT, INDIRECT, or DIRECT_AND_INDIRECT`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleCloudIdentityGroupTransitiveMembershipsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	result := []map[string]interface{}{}
	membershipsCall := config.NewCloudIdentityClient(userAgent).Groups.Memberships.SearchTransitiveMemberships(d.Get("group").(string))
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

	err = membershipsCall.Pages(config.Context, func(resp *cloudidentity.SearchTransitiveMembershipsResponse) error {
		for _, member := range resp.Memberships {
			result = append(result, map[string]interface{}{
				"member":               member.Member,
				"relation_type":        member.RelationType,
				"roles":                flattenCloudIdentityGroupTransitiveMembershipsRoles(member.Roles),
				"preferred_member_key": flattenCloudIdentityGroupsEntityKeyList(member.PreferredMemberKey),
			})
		}

		return nil
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("CloudIdentityGroupMemberships %q", d.Id()), "")
	}

	if err := d.Set("memberships", result); err != nil {
		return fmt.Errorf("Error setting memberships: %s", err)
	}

	group := d.Get("group")
	d.SetId(fmt.Sprintf("%s/transitiveMemberships", group.(string))) // groups/{group_id}/transitiveMemberships
	return nil
}

func flattenCloudIdentityGroupTransitiveMembershipsRoles(roles []*cloudidentity.TransitiveMembershipRole) []interface{} {
	transformed := []interface{}{}

	for _, role := range roles {
		transformed = append(transformed, map[string]interface{}{
			"role": role.Role,
		})
	}
	return transformed
}

// flattenCloudIdentityGroupsEntityKeyList is a version of flattenCloudIdentityGroupsEntityKey that
// can accept a list of EntityKeys
func flattenCloudIdentityGroupsEntityKeyList(entityKeys []*cloudidentity.EntityKey) []interface{} {
	transformed := []interface{}{}

	for _, key := range entityKeys {
		transformed = append(transformed, map[string]interface{}{
			"id":        key.Id,
			"namespace": key.Namespace,
		})
	}

	return transformed
}
