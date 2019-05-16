package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func resourceGoogleFolderOrganizationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleFolderOrganizationPolicyCreate,
		Read:   resourceGoogleFolderOrganizationPolicyRead,
		Update: resourceGoogleFolderOrganizationPolicyUpdate,
		Delete: resourceGoogleFolderOrganizationPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceFolderOrgPolicyImporter,
		},

		Schema: mergeSchemas(
			schemaOrganizationPolicy,
			map[string]*schema.Schema{
				"folder": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		),
	}
}

func resourceFolderOrgPolicyImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"folders/(?P<folder>[^/]+):constraints/(?P<constraint>[^/]+)",
		"(?P<folder>[^/]+):(?P<constraint>[^/]+)"},
		d, config); err != nil {
		return nil, err
	}

	if d.Get("folder") == "" || d.Get("constraint") == "" {
		return nil, fmt.Errorf("unable to parse folder or constraint. Check import formats")
	}

	d.Set("folder", "folders/"+d.Get("folder").(string))

	return []*schema.ResourceData{d}, nil
}

func resourceGoogleFolderOrganizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("folder"), d.Get("constraint")))

	if isOrganizationPolicyUnset(d) {
		return resourceGoogleFolderOrganizationPolicyDelete(d, meta)
	}

	if err := setFolderOrganizationPolicy(d, meta); err != nil {
		return err
	}

	return resourceGoogleFolderOrganizationPolicyRead(d, meta)
}

func resourceGoogleFolderOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	folder := canonicalFolderId(d.Get("folder").(string))

	policy, err := config.clientResourceManager.Folders.GetOrgPolicy(folder, &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
	}).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Organization policy for %s", folder))
	}

	d.Set("constraint", policy.Constraint)
	d.Set("boolean_policy", flattenBooleanOrganizationPolicy(policy.BooleanPolicy))
	d.Set("list_policy", flattenListOrganizationPolicy(policy.ListPolicy))
	d.Set("restore_policy", flattenRestoreOrganizationPolicy(policy.RestoreDefault))
	d.Set("version", policy.Version)
	d.Set("etag", policy.Etag)
	d.Set("update_time", policy.UpdateTime)

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
	config := meta.(*Config)
	folder := canonicalFolderId(d.Get("folder").(string))

	_, err := config.clientResourceManager.Folders.ClearOrgPolicy(folder, &cloudresourcemanager.ClearOrgPolicyRequest{
		Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func setFolderOrganizationPolicy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	folder := canonicalFolderId(d.Get("folder").(string))

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return err
	}

	_, err = config.clientResourceManager.Folders.SetOrgPolicy(folder, &cloudresourcemanager.SetOrgPolicyRequest{
		Policy: &cloudresourcemanager.OrgPolicy{
			Constraint:     canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
			BooleanPolicy:  expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
			ListPolicy:     listPolicy,
			RestoreDefault: restoreDefault,
			Version:        int64(d.Get("version").(int)),
			Etag:           d.Get("etag").(string),
		},
	}).Do()

	return err
}
