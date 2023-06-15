// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func ResourceGoogleFolderOrganizationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleFolderOrganizationPolicyCreate,
		Read:   resourceGoogleFolderOrganizationPolicyRead,
		Update: resourceGoogleFolderOrganizationPolicyUpdate,
		Delete: resourceGoogleFolderOrganizationPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceFolderOrgPolicyImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Read:   schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: tpgresource.MergeSchemas(
			schemaOrganizationPolicy,
			map[string]*schema.Schema{
				"folder": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: `The resource name of the folder to set the policy for. Its format is folders/{folder_id}.`,
				},
			},
		),
		UseJSONNumber: true,
	}
}

func resourceFolderOrgPolicyImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"folders/(?P<folder>[^/]+)/constraints/(?P<constraint>[^/]+)",
		"folders/(?P<folder>[^/]+)/(?P<constraint>[^/]+)",
		"(?P<folder>[^/]+)/(?P<constraint>[^/]+)"},
		d, config); err != nil {
		return nil, err
	}

	if d.Get("folder") == "" || d.Get("constraint") == "" {
		return nil, fmt.Errorf("unable to parse folder or constraint. Check import formats")
	}

	if err := d.Set("folder", "folders/"+d.Get("folder").(string)); err != nil {
		return nil, fmt.Errorf("Error setting folder: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceGoogleFolderOrganizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprintf("%s/%s", d.Get("folder"), d.Get("constraint")))

	if isOrganizationPolicyUnset(d) {
		return resourceGoogleFolderOrganizationPolicyDelete(d, meta)
	}

	if err := setFolderOrganizationPolicy(d, meta); err != nil {
		return err
	}

	return resourceGoogleFolderOrganizationPolicyRead(d, meta)
}

func resourceGoogleFolderOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	folder := CanonicalFolderId(d.Get("folder").(string))

	var policy *cloudresourcemanager.OrgPolicy
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (getErr error) {
			policy, getErr = config.NewResourceManagerClient(userAgent).Folders.GetOrgPolicy(folder, &cloudresourcemanager.GetOrgPolicyRequest{
				Constraint: CanonicalOrgPolicyConstraint(d.Get("constraint").(string)),
			}).Do()
			return getErr
		},
		Timeout: d.Timeout(schema.TimeoutRead),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Organization policy for %s", folder))
	}

	if err := d.Set("constraint", policy.Constraint); err != nil {
		return fmt.Errorf("Error setting constraint: %s", err)
	}
	if err := d.Set("boolean_policy", flattenBooleanOrganizationPolicy(policy.BooleanPolicy)); err != nil {
		return fmt.Errorf("Error setting boolean_policy: %s", err)
	}
	if err := d.Set("list_policy", flattenListOrganizationPolicy(policy.ListPolicy)); err != nil {
		return fmt.Errorf("Error setting list_policy: %s", err)
	}
	if err := d.Set("restore_policy", flattenRestoreOrganizationPolicy(policy.RestoreDefault)); err != nil {
		return fmt.Errorf("Error setting restore_policy: %s", err)
	}
	if err := d.Set("version", policy.Version); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}
	if err := d.Set("etag", policy.Etag); err != nil {
		return fmt.Errorf("Error setting etag: %s", err)
	}
	if err := d.Set("update_time", policy.UpdateTime); err != nil {
		return fmt.Errorf("Error setting update_time: %s", err)
	}

	return nil
}

func resourceGoogleFolderOrganizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	if isOrganizationPolicyUnset(d) {
		return resourceGoogleFolderOrganizationPolicyDelete(d, meta)
	}

	if err := setFolderOrganizationPolicy(d, meta); err != nil {
		return err
	}

	return resourceGoogleFolderOrganizationPolicyRead(d, meta)
}

func resourceGoogleFolderOrganizationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	folder := CanonicalFolderId(d.Get("folder").(string))

	return transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (delErr error) {
			_, delErr = config.NewResourceManagerClient(userAgent).Folders.ClearOrgPolicy(folder, &cloudresourcemanager.ClearOrgPolicyRequest{
				Constraint: CanonicalOrgPolicyConstraint(d.Get("constraint").(string)),
			}).Do()
			return delErr
		},
		Timeout: d.Timeout(schema.TimeoutDelete),
	})
}

func setFolderOrganizationPolicy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	folder := CanonicalFolderId(d.Get("folder").(string))

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return err
	}

	return transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (setErr error) {
			_, setErr = config.NewResourceManagerClient(userAgent).Folders.SetOrgPolicy(folder, &cloudresourcemanager.SetOrgPolicyRequest{
				Policy: &cloudresourcemanager.OrgPolicy{
					Constraint:     CanonicalOrgPolicyConstraint(d.Get("constraint").(string)),
					BooleanPolicy:  expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
					ListPolicy:     listPolicy,
					RestoreDefault: restoreDefault,
					Version:        int64(d.Get("version").(int)),
					Etag:           d.Get("etag").(string),
				},
			}).Do()
			return setErr
		},
		Timeout: d.Timeout(schema.TimeoutCreate),
	})
}
