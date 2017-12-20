package google

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"google.golang.org/api/cloudresourcemanager/v1"
	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform/helper/schema"
)

var schemaOrganizationPolicy = map[string]*schema.Schema{
	"constraint": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: linkDiffSuppress,
	},
	"boolean_policy_source": {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"list_policy", "boolean_policy", "list_policy_source"},
	},
	"boolean_policy": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"list_policy", "boolean_policy_source", "list_policy_source"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enforced": {
					Type:     schema.TypeBool,
					Required: true,
				},
			},
		},
	},
	"list_policy_source": {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"list_policy", "boolean_policy", "boolean_policy_source"},
	},
	"list_policy": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"boolean_policy", "boolean_policy_source", "list_policy_source"},
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
	// Detect changes to local file or changes made outside of Terraform to the file stored on the server.
	"detect_policy_change": &schema.Schema{
		Type: schema.TypeString,
		// This field is not Computed because it needs to trigger a diff.
		Optional: true,
		ForceNew: true,
		// Makes the diff message nicer:
		Default: "different policy",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			localPolicy := ""
			var err error

			if source, ok := d.GetOkExists("boolean_policy_source"); ok {
				filepolicy, err := getFilePolicy(source.(string))
				if err != nil {
					localPolicy, err = getYAMLPolicy(filepolicy)
				}
			}

			if content, ok := d.GetOkExists("boolean_policy"); ok {
				localPolicy, err = getYAMLPolicy(content)
			}

			if source, ok := d.GetOkExists("list_policy_source"); ok {
				filepolicy, err := getFilePolicy(source.(string))
				if err != nil {
					localPolicy, err = getYAMLPolicy(filepolicy)
				}
			}

			if content, ok := d.GetOkExists("list_policy"); ok {
				localPolicy, err = getYAMLPolicy(content)
			}

			if err != nil {
				return false
			}

			oldpolicy, err := getYAMLPolicy(old)
			if err != nil {
				return false
			}

			if oldpolicy != localPolicy {
				return false
			}

			return true
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
	if err := setOrganizationPolicy(d, meta); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s", d.Get("org_id"), d.Get("constraint").(string)))

	return resourceGoogleOrganizationPolicyRead(d, meta)
}

func resourceGoogleOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	policy, err := config.clientResourceManager.Organizations.GetOrgPolicy(org, &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
	}).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Organization policy for %s", org))
	}

	d.Set("constraint", policy.Constraint)
	d.Set("boolean_policy", flattenBooleanOrganizationPolicy(policy.BooleanPolicy))
	d.Set("list_policy", flattenListOrganizationPolicy(policy.ListPolicy))
	d.Set("version", policy.Version)
	d.Set("etag", policy.Etag)
	d.Set("update_time", policy.UpdateTime)

	return nil
}

func resourceGoogleOrganizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := setOrganizationPolicy(d, meta); err != nil {
		return err
	}

	return resourceGoogleOrganizationPolicyRead(d, meta)
}

func resourceGoogleOrganizationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	_, err := config.clientResourceManager.Organizations.ClearOrgPolicy(org, &cloudresourcemanager.ClearOrgPolicyRequest{
		Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
	}).Do()

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

func setOrganizationPolicy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)
	var err error
	var booleanPolicy *cloudresourcemanager.BooleanPolicy
	var listPolicy *cloudresourcemanager.ListPolicy

	if source, ok := d.GetOkExists("list_policy_source"); ok {
		filePolicy, err := getFilePolicy(source.(string))
		if err != nil {
			listPolicy, err = expandListOrganizationPolicy(filePolicy)
		}
	} else {
		listPolicy, err = expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	}
	if err != nil {
		return err
	}

	if source, ok := d.GetOkExists("boolean_policy_source"); ok {
		filePolicy, err := getFilePolicy(source.(string))
		if err != nil {
			booleanPolicy = expandBooleanOrganizationPolicy(filePolicy)
		}
	} else {
		booleanPolicy = expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{}))
	}

	_, err = config.clientResourceManager.Organizations.SetOrgPolicy(org, &cloudresourcemanager.SetOrgPolicyRequest{
		Policy: &cloudresourcemanager.OrgPolicy{
			Constraint:    canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
			BooleanPolicy: booleanPolicy,
			ListPolicy:    listPolicy,
			Version:       int64(d.Get("version").(int)),
			Etag:          d.Get("etag").(string),
		},
	}).Do()

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

func expandBooleanOrganizationPolicy(configured []interface{}) *cloudresourcemanager.BooleanPolicy {
	if len(configured) == 0 {
		return nil
	}

	booleanPolicy := configured[0].(map[string]interface{})
	return &cloudresourcemanager.BooleanPolicy{
		Enforced: booleanPolicy["enforced"].(bool),
	}
}

func flattenListOrganizationPolicy(policy *cloudresourcemanager.ListPolicy) []map[string]interface{} {
	lPolicies := make([]map[string]interface{}, 0, 1)

	if policy == nil {
		return lPolicies
	}

	listPolicy := map[string]interface{}{}
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
	if len(configured) == 0 {
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
		AllValues:      allValues,
		AllowedValues:  allowedValues,
		DeniedValues:   deniedValues,
		SuggestedValue: listPolicy["suggested_value"].(string),
	}, nil
}

func canonicalOrgPolicyConstraint(constraint string) string {
	if strings.HasPrefix(constraint, "constraints/") {
		return constraint
	}
	return "constraints/" + constraint
}

// Returns an expanded list from a file
func getFilePolicy(filename string) ([]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("[WARN] Failed to read source file %q.", filename)
		return nil, err
	}
	var h interface{}
	err = yaml.Unmarshal(data, h)
	if err != nil {
		return nil, err
	}

	var retval []interface{}
	retval = append(retval, h)
	return retval, nil
}

// Returns a YAML string from an expanded list
func getYAMLPolicy(content interface{}) (string, error) {
	bytes, err := yaml.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(bytes[:]), nil
}
