package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudidentity/v1"
)

func dataSourceGoogleCloudIdentityGroups() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceCloudIdentityGroup().Schema)

	return &schema.Resource{
		Read: dataSourceGoogleCloudIdentityGroupsRead,

		Schema: map[string]*schema.Schema{
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `List of Cloud Identity groups.`,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The resource name of the entity under which this Group resides in the
Cloud Identity resource hierarchy.

Must be of the form identitysources/{identity_source_id} for external-identity-mapped
groups or customers/{customer_id} for Google Groups.`,
			},
		},
	}
}

func dataSourceGoogleCloudIdentityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	result := []map[string]interface{}{}
	groupsCall := config.NewCloudIdentityClient(userAgent).Groups.List().Parent(d.Get("parent").(string)).View("FULL")
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
			groupsCall.Header().Set("X-Goog-User-Project", billingProject)
		}
	}
	err = groupsCall.Pages(config.context, func(resp *cloudidentity.ListGroupsResponse) error {
		for _, group := range resp.Groups {
			result = append(result, map[string]interface{}{
				"name":         group.Name,
				"display_name": group.DisplayName,
				"labels":       group.Labels,
				"description":  group.Description,
				"group_key":    flattenCloudIdentityGroupsEntityKey(group.GroupKey),
			})
		}

		return nil
	})
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("CloudIdentityGroups %q", d.Id()))
	}

	if err := d.Set("groups", result); err != nil {
		return fmt.Errorf("Error setting groups: %s", err)
	}
	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenCloudIdentityGroupsEntityKey(entityKey *cloudidentity.EntityKey) []interface{} {
	transformed := map[string]interface{}{
		"id":        entityKey.Id,
		"namespace": entityKey.Namespace,
	}
	return []interface{}{transformed}
}
