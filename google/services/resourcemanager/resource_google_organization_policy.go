// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var schemaOrganizationPolicy = map[string]*schema.Schema{
	// Although the API suggests that boolean_policy, list_policy, or restore_policy must be set,
	// Organization policies can be "inherited from parent" in the UI, and this is the default
	// state of the resource without any policy set.
	// See https://github.com/hashicorp/terraform-provider-google/issues/3607
	"constraint": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
		Description:      `The name of the Constraint the Policy is configuring, for example, serviceuser.services.`,
	},
	"boolean_policy": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `A boolean policy is a constraint that is either enforced or not.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enforced": {
					Type:        schema.TypeBool,
					Required:    true,
					Description: `If true, then the Policy is enforced. If false, then any configuration is acceptable.`,
				},
			},
		},
	},
	"list_policy": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `A policy that can define specific values that are allowed or denied for the given constraint. It can also be used to allow or deny all values. `,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					Description:  `One or the other must be set.`,
					ExactlyOneOf: []string{"list_policy.0.allow", "list_policy.0.deny"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"all": {
								Type:         schema.TypeBool,
								Optional:     true,
								Default:      false,
								Description:  `The policy allows or denies all values.`,
								ExactlyOneOf: []string{"list_policy.0.allow.0.all", "list_policy.0.allow.0.values"},
							},
							"values": {
								Type:         schema.TypeSet,
								Optional:     true,
								Description:  `The policy can define specific values that are allowed or denied.`,
								ExactlyOneOf: []string{"list_policy.0.allow.0.all", "list_policy.0.allow.0.values"},
								Elem:         &schema.Schema{Type: schema.TypeString},
								Set:          schema.HashString,
							},
						},
					},
				},
				"deny": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					Description:  `One or the other must be set.`,
					ExactlyOneOf: []string{"list_policy.0.allow", "list_policy.0.deny"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"all": {
								Type:         schema.TypeBool,
								Optional:     true,
								Default:      false,
								Description:  `The policy allows or denies all values.`,
								ExactlyOneOf: []string{"list_policy.0.deny.0.all", "list_policy.0.deny.0.values"},
							},
							"values": {
								Type:         schema.TypeSet,
								Optional:     true,
								Description:  `The policy can define specific values that are allowed or denied.`,
								ExactlyOneOf: []string{"list_policy.0.deny.0.all", "list_policy.0.deny.0.values"},
								Elem:         &schema.Schema{Type: schema.TypeString},
								Set:          schema.HashString,
							},
						},
					},
				},
				"suggested_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: `The Google Cloud Console will try to default to a configuration that matches the value specified in this field.`,
				},
				"inherit_from_parent": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: `If set to true, the values from the effective Policy of the parent resource are inherited, meaning the values set in this Policy are added to the values inherited up the hierarchy.`,
				},
			},
		},
	},
	"version": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: `Version of the Policy. Default version is 0.`,
	},
	"etag": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The etag of the organization policy. etag is used for optimistic concurrency control as a way to help prevent simultaneous updates of a policy from overwriting each other.`,
	},
	"update_time": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds, representing when the variable was last updated. Example: "2016-10-09T12:33:37.578138407Z".`,
	},
	"restore_policy": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `A restore policy is a constraint to restore the default policy.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default": {
					Type:        schema.TypeBool,
					Required:    true,
					Description: `May only be set to true. If set, then the default Policy is restored.`,
				},
			},
		},
	},
}

func ResourceGoogleOrganizationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleOrganizationPolicyCreate,
		Read:   resourceGoogleOrganizationPolicyRead,
		Update: resourceGoogleOrganizationPolicyUpdate,
		Delete: resourceGoogleOrganizationPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceGoogleOrganizationPolicyImportState,
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
				"org_id": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			}),
		UseJSONNumber: true,
	}
}

func resourceGoogleOrganizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	if isOrganizationPolicyUnset(d) {
		return resourceGoogleOrganizationPolicyDelete(d, meta)
	}

	if err := setOrganizationPolicy(d, meta); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("org_id"), d.Get("constraint").(string)))
	return resourceGoogleOrganizationPolicyRead(d, meta)
}

func resourceGoogleOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	org := "organizations/" + d.Get("org_id").(string)

	var policy *cloudresourcemanager.OrgPolicy
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (readErr error) {
			policy, readErr = config.NewResourceManagerClient(userAgent).Organizations.GetOrgPolicy(org, &cloudresourcemanager.GetOrgPolicyRequest{
				Constraint: CanonicalOrgPolicyConstraint(d.Get("constraint").(string)),
			}).Do()
			return readErr
		},
		Timeout: d.Timeout(schema.TimeoutRead),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Organization policy for %s", org))
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
	if err := d.Set("version", policy.Version); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}
	if err := d.Set("etag", policy.Etag); err != nil {
		return fmt.Errorf("Error setting etag: %s", err)
	}
	if err := d.Set("update_time", policy.UpdateTime); err != nil {
		return fmt.Errorf("Error setting update_time: %s", err)
	}
	if err := d.Set("restore_policy", flattenRestoreOrganizationPolicy(policy.RestoreDefault)); err != nil {
		return fmt.Errorf("Error setting restore_policy: %s", err)
	}

	return nil
}

func resourceGoogleOrganizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	if isOrganizationPolicyUnset(d) {
		return resourceGoogleOrganizationPolicyDelete(d, meta)
	}

	if err := setOrganizationPolicy(d, meta); err != nil {
		return err
	}

	return resourceGoogleOrganizationPolicyRead(d, meta)
}

func resourceGoogleOrganizationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	org := "organizations/" + d.Get("org_id").(string)

	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			_, dErr := config.NewResourceManagerClient(userAgent).Organizations.ClearOrgPolicy(org, &cloudresourcemanager.ClearOrgPolicyRequest{
				Constraint: CanonicalOrgPolicyConstraint(d.Get("constraint").(string)),
			}).Do()
			return dErr
		},
		Timeout: d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceGoogleOrganizationPolicyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid id format. Expecting {org_id}/{constraint}, got '%s' instead.", d.Id())
	}

	if err := d.Set("org_id", parts[0]); err != nil {
		return nil, fmt.Errorf("Error setting org_id: %s", err)
	}
	if err := d.Set("constraint", parts[1]); err != nil {
		return nil, fmt.Errorf("Error setting constraint: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

// Organization policies can be "inherited from parent" the UI, and this is the default
// state of the resource without any policy set. In order to revert to this state the current
// resource cannot be updated it must instead be Deleted. This allows Terraform to assert that
// no policy has been set even if previously one had.
// See https://github.com/hashicorp/terraform-provider-google/issues/3607
func isOrganizationPolicyUnset(d *schema.ResourceData) bool {
	listPolicy := d.Get("list_policy").([]interface{})
	booleanPolicy := d.Get("boolean_policy").([]interface{})
	restorePolicy := d.Get("restore_policy").([]interface{})

	return len(listPolicy)+len(booleanPolicy)+len(restorePolicy) == 0
}

func setOrganizationPolicy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	org := "organizations/" + d.Get("org_id").(string)

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return err
	}

	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (setErr error) {
			_, setErr = config.NewResourceManagerClient(userAgent).Organizations.SetOrgPolicy(org, &cloudresourcemanager.SetOrgPolicyRequest{
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
	return err
}

func flattenBooleanOrganizationPolicy(policy *cloudresourcemanager.BooleanPolicy) []map[string]interface{} {
	bPolicies := make([]map[string]interface{}, 0, 1)

	if policy == nil {
		return bPolicies
	}

	bPolicies = append(bPolicies, map[string]interface{}{
		"enforced": policy.Enforced,
	})

	return bPolicies
}

func flattenRestoreOrganizationPolicy(restore_policy *cloudresourcemanager.RestoreDefault) []map[string]interface{} {
	rp := make([]map[string]interface{}, 0, 1)

	if restore_policy == nil {
		return rp
	}

	rp = append(rp, map[string]interface{}{
		"default": true,
	})

	return rp
}

func expandBooleanOrganizationPolicy(configured []interface{}) *cloudresourcemanager.BooleanPolicy {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	booleanPolicy := configured[0].(map[string]interface{})
	return &cloudresourcemanager.BooleanPolicy{
		Enforced: booleanPolicy["enforced"].(bool),
	}
}

func expandRestoreOrganizationPolicy(configured []interface{}) (*cloudresourcemanager.RestoreDefault, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	restoreDefaultMap := configured[0].(map[string]interface{})
	default_value := restoreDefaultMap["default"].(bool)

	if default_value {
		return &cloudresourcemanager.RestoreDefault{}, nil
	}

	return nil, fmt.Errorf("Invalid value for restore_policy. Expecting default = true")
}

func flattenListOrganizationPolicy(policy *cloudresourcemanager.ListPolicy) []map[string]interface{} {
	lPolicies := make([]map[string]interface{}, 0, 1)

	if policy == nil {
		return lPolicies
	}

	listPolicy := map[string]interface{}{
		"suggested_value":     policy.SuggestedValue,
		"inherit_from_parent": policy.InheritFromParent,
	}
	switch {
	case policy.AllValues == "ALLOW":
		listPolicy["allow"] = []interface{}{map[string]interface{}{
			"all": true,
		}}
	case policy.AllValues == "DENY":
		listPolicy["deny"] = []interface{}{map[string]interface{}{
			"all": true,
		}}
	case len(policy.AllowedValues) > 0:
		listPolicy["allow"] = []interface{}{map[string]interface{}{
			"values": schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(policy.AllowedValues)),
		}}
	case len(policy.DeniedValues) > 0:
		listPolicy["deny"] = []interface{}{map[string]interface{}{
			"values": schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(policy.DeniedValues)),
		}}
	}

	lPolicies = append(lPolicies, listPolicy)

	return lPolicies
}

func expandListOrganizationPolicy(configured []interface{}) (*cloudresourcemanager.ListPolicy, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	listPolicyMap := configured[0].(map[string]interface{})

	allow := listPolicyMap["allow"].([]interface{})
	deny := listPolicyMap["deny"].([]interface{})

	var allValues string
	var allowedValues []string
	var deniedValues []string
	if len(allow) > 0 {
		allowMap := allow[0].(map[string]interface{})
		all := allowMap["all"].(bool)
		values := allowMap["values"].(*schema.Set)

		if all {
			allValues = "ALLOW"
		} else {
			allowedValues = tpgresource.ConvertStringArr(values.List())
		}
	}

	if len(deny) > 0 {
		denyMap := deny[0].(map[string]interface{})
		all := denyMap["all"].(bool)
		values := denyMap["values"].(*schema.Set)

		if all {
			allValues = "DENY"
		} else {
			deniedValues = tpgresource.ConvertStringArr(values.List())
		}
	}

	listPolicy := configured[0].(map[string]interface{})
	return &cloudresourcemanager.ListPolicy{
		AllValues:         allValues,
		AllowedValues:     allowedValues,
		DeniedValues:      deniedValues,
		SuggestedValue:    listPolicy["suggested_value"].(string),
		InheritFromParent: listPolicy["inherit_from_parent"].(bool),
		ForceSendFields:   []string{"InheritFromParent"},
	}, nil
}

func CanonicalOrgPolicyConstraint(constraint string) string {
	if strings.HasPrefix(constraint, "constraints/") {
		return constraint
	}
	return "constraints/" + constraint
}
