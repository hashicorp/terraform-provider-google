// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleCloudIdentityGroupLookup() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceGoogleCloudIdentityGroupLookupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The [resource name](https://cloud.google.com/apis/design/resource_names) of the looked-up Group.`,
			},
			"group_key": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Description: `The EntityKey of the Group to lookup. A unique identifier for an entity in the Cloud Identity Groups API.
An entity can represent either a group with an optional namespace or a user without a namespace.
The combination of id and namespace must be unique; however, the same id can be used with different namespaces.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The ID of the entity. For Google-managed entities, the id should be the email address of an existing group or user.
For external-identity-mapped entities, the id must be a string conforming to the Identity Source's requirements.
Must be unique within a namespace.`,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `The namespace in which the entity exists. If not specified, the EntityKey represents a Google-managed entity such as a Google user or a Google Group.
If specified, the EntityKey represents an external-identity-mapped group. The namespace must correspond to an identity source created in Admin Console and must be in the form of identitysources/{identity_source}.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleCloudIdentityGroupLookupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	gkId, ok := d.GetOk("group_key.0.id")
	if !ok {
		return fmt.Errorf("error getting group key id")
	}
	id := gkId.(string)

	groupsLookupCall := config.NewCloudIdentityClient(userAgent).Groups.Lookup().GroupKeyId(id)

	gkNamespace, ok := d.GetOk("group_key.0.namespace")
	if ok {
		// If optional namespace argument provided, add as param to API call
		namespace := gkNamespace.(string)
		groupsLookupCall = groupsLookupCall.GroupKeyNamespace(namespace)
	}

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
			groupsLookupCall.Header().Set("X-Goog-User-Project", billingProject)
		}
	}
	resp, err := groupsLookupCall.Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("CloudIdentityGroups %q", d.Id()), "Groups")
	}

	if err := d.Set("name", resp.Name); err != nil {
		return fmt.Errorf("error setting group lookup name: %s", err)
	}
	d.SetId(time.Now().UTC().String())
	return nil
}
