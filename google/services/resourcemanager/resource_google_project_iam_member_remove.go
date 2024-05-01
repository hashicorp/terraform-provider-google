// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func ResourceGoogleProjectIamMemberRemove() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamMemberRemoveCreate,
		Read:   resourceGoogleProjectIamMemberRemoveRead,
		Delete: resourceGoogleProjectIamMemberRemoveDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `The project id of the target project.`,
			},
			"role": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `The target role that should be removed.`,
			},
			"member": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `The IAM principal that should not have the target role.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceGoogleProjectIamMemberRemoveCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project := d.Get("project").(string)
	role := d.Get("role").(string)
	member := d.Get("member").(string)

	found := false
	iamPolicy, err := config.NewResourceManagerClient(config.UserAgent).Projects.GetIamPolicy(project,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, d.Id())
	}

	for i := 0; i < len(iamPolicy.Bindings); i++ {
		if role == iamPolicy.Bindings[i].Role {
			for j := 0; j < len(iamPolicy.Bindings[i].Members); j++ {
				if member == iamPolicy.Bindings[i].Members[j] {
					found = true
					iamPolicy.Bindings[i].Members = append(iamPolicy.Bindings[i].Members[:j], iamPolicy.Bindings[i].Members[j+1:]...)
					break
				}
			}
		}
	}

	if found == false {
		fmt.Printf("[DEBUG] Could not find Member %s with the corresponding role %s. No removal necessary", member, role)
	} else {
		updateRequest := &cloudresourcemanager.SetIamPolicyRequest{
			Policy:     iamPolicy,
			UpdateMask: "bindings",
		}
		_, err = config.NewResourceManagerClient(config.UserAgent).Projects.SetIamPolicy(project, updateRequest).Do()
		if err != nil {
			return fmt.Errorf("cannot update IAM policy on project %s: %v", project, err)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", project, member, role))

	return resourceGoogleProjectIamMemberRemoveRead(d, meta)
}

func resourceGoogleProjectIamMemberRemoveRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project := d.Get("project").(string)
	role := d.Get("role").(string)
	member := d.Get("member").(string)

	found := false
	iamPolicy, err := config.NewResourceManagerClient(config.UserAgent).Projects.GetIamPolicy(project,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, d.Id())
	}

	for i := 0; i < len(iamPolicy.Bindings); i++ {
		if role == iamPolicy.Bindings[i].Role {
			for j := 0; j < len(iamPolicy.Bindings[i].Members); j++ {
				if member == iamPolicy.Bindings[i].Members[j] {
					found = true
					break
				}
			}
		}
	}

	if found {
		fmt.Printf("[DEBUG] found membership in project's policy  %v, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

func resourceGoogleProjectIamMemberRemoveDelete(d *schema.ResourceData, meta interface{}) error {
	fmt.Printf("[DEBUG] clearing resource %v from state", d.Id())
	d.SetId("")

	return nil
}
