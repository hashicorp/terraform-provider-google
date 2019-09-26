package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var schemaOrganizationPolicy = map[string]*schema.Schema{
	"constraint": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: compareSelfLinkOrResourceName,
	},
	"boolean_policy": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"list_policy", "restore_policy"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enforced": {
					Type:     schema.TypeBool,
					Required: true,
				},
			},
		},
	},
	"list_policy": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"boolean_policy", "restore_policy"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"list_policy.0.deny"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"all": {
								Type:          schema.TypeBool,
								Optional:      true,
								Default:       false,
								ConflictsWith: []string{"list_policy.0.allow.0.values"},
							},
							"values": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Set:      schema.HashString,
							},
						},
					},
				},
				"deny": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"all": {
								Type:          schema.TypeBool,
								Optional:      true,
								Default:       false,
								ConflictsWith: []string{"list_policy.0.deny.0.values"},
							},
							"values": {
								Type:     schema.TypeSet,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Set:      schema.HashString,
							},
						},
					},
				},
				"suggested_value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"inherit_from_parent": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			},
		},
	},
	"version": {
		Type:     schema.TypeInt,
		Optional: true,
		Computed: true,
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"update_time": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"restore_policy": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"boolean_policy", "list_policy"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default": {
					Type:     schema.TypeBool,
					Required: true,
				},
			},
		},
	},
}

func resourceGoogleOrganizationPolicy() *schema.Resource {
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

		Schema: mergeSchemas(
			schemaOrganizationPolicy,
			map[string]*schema.Schema{
				"org_id": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			}),
	}
}

func resourceGoogleOrganizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprintf("%s:%s", d.Get("org_id"), d.Get("constraint").(string)))

	if isOrganizationPolicyUnset(d) {
		return resourceGoogleOrganizationPolicyDelete(d, meta)
	}

	if err := setOrganizationPolicy(d, meta); err != nil {
		return err
	}

	return resourceGoogleOrganizationPolicyRead(d, meta)
}

func resourceGoogleOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	var policy *cloudresourcemanager.OrgPolicy
	err := retryTimeDuration(func() (readErr error) {
		policy, readErr = config.clientResourceManager.Organizations.GetOrgPolicy(org, &cloudresourcemanager.GetOrgPolicyRequest{
			Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
		}).Do()
		return readErr
	}, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Organization policy for %s", org))
	}

	d.Set("constraint", policy.Constraint)
	d.Set("boolean_policy", flattenBooleanOrganizationPolicy(policy.BooleanPolicy))
	d.Set("list_policy", flattenListOrganizationPolicy(policy.ListPolicy))
	d.Set("version", policy.Version)
	d.Set("etag", policy.Etag)
	d.Set("update_time", policy.UpdateTime)
	d.Set("restore_policy", flattenRestoreOrganizationPolicy(policy.RestoreDefault))

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
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	err := retryTimeDuration(func() error {
		_, dErr := config.clientResourceManager.Organizations.ClearOrgPolicy(org, &cloudresourcemanager.ClearOrgPolicyRequest{
			Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
		}).Do()
		return dErr
	}, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	return nil
}

func resourceGoogleOrganizationPolicyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid id format. Expecting {org_id}:{constraint}, got '%s' instead.", d.Id())
	}

	d.Set("org_id", parts[0])
	d.Set("constraint", parts[1])

	return []*schema.ResourceData{d}, nil
}

// Organization policies can be "inherited from parent" the UI, and this is the default
// state of the resource without any policy set. In order to revert to this state the current
// resource cannot be updated it must instead be Deleted. This allows Terraform to assert that
// no policy has been set even if previously one had.
// See https://github.com/terraform-providers/terraform-provider-google/issues/3607
func isOrganizationPolicyUnset(d *schema.ResourceData) bool {
	listPolicy := d.Get("list_policy").([]interface{})
	booleanPolicy := d.Get("boolean_policy").([]interface{})
	restorePolicy := d.Get("restore_policy").([]interface{})

	return len(listPolicy)+len(booleanPolicy)+len(restorePolicy) == 0
}

func setOrganizationPolicy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return err
	}

	restoreDefault, err := expandRestoreOrganizationPolicy(d.Get("restore_policy").([]interface{}))
	if err != nil {
		return err
	}

	err = retryTimeDuration(func() (setErr error) {
		_, setErr = config.clientResourceManager.Organizations.SetOrgPolicy(org, &cloudresourcemanager.SetOrgPolicyRequest{
			Policy: &cloudresourcemanager.OrgPolicy{
				Constraint:     canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
				BooleanPolicy:  expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
				ListPolicy:     listPolicy,
				RestoreDefault: restoreDefault,
				Version:        int64(d.Get("version").(int)),
				Etag:           d.Get("etag").(string),
			},
		}).Do()
		return setErr
	}, d.Timeout(schema.TimeoutCreate))
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
			"values": schema.NewSet(schema.HashString, convertStringArrToInterface(policy.AllowedValues)),
		}}
	case len(policy.DeniedValues) > 0:
		listPolicy["deny"] = []interface{}{map[string]interface{}{
			"values": schema.NewSet(schema.HashString, convertStringArrToInterface(policy.DeniedValues)),
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
			allowedValues = convertStringArr(values.List())
		}
	}

	if len(deny) > 0 {
		denyMap := deny[0].(map[string]interface{})
		all := denyMap["all"].(bool)
		values := denyMap["values"].(*schema.Set)

		if all {
			allValues = "DENY"
		} else {
			deniedValues = convertStringArr(values.List())
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

func canonicalOrgPolicyConstraint(constraint string) string {
	if strings.HasPrefix(constraint, "constraints/") {
		return constraint
	}
	return "constraints/" + constraint
}
