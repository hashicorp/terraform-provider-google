package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudidentity/v1"
)

func DataSourceGoogleCloudIdentityGroupMemberships() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(ResourceCloudIdentityGroupMembership().Schema)

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
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name of the Group to get memberships from.`,
			},
		},
	}
}

func dataSourceGoogleCloudIdentityGroupMembershipsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	result := []map[string]interface{}{}
	membershipsCall := config.NewCloudIdentityClient(userAgent).Groups.Memberships.List(d.Get("group").(string)).View("FULL")
	if config.UserProjectOverride {
		billingProject := ""
		// err may be nil - project isn't required for this resource
		if project, err := getProject(d, config); err == nil {
			billingProject = project
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		if billingProject != "" {
			membershipsCall.Header().Set("X-Goog-User-Project", billingProject)
		}
	}

	err = membershipsCall.Pages(config.context, func(resp *cloudidentity.ListMembershipsResponse) error {
		for _, member := range resp.Memberships {
			result = append(result, map[string]interface{}{
				"name":                 member.Name,
				"roles":                flattenCloudIdentityGroupMembershipsRoles(member.Roles),
				"preferred_member_key": flattenCloudIdentityGroupsEntityKey(member.PreferredMemberKey),
			})
		}

		return nil
	})
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("CloudIdentityGroupMemberships %q", d.Id()))
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
