// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package orgpolicy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	orgpolicy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/orgpolicy"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceOrgPolicyPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrgPolicyPolicyCreate,
		Read:   resourceOrgPolicyPolicyRead,
		Update: resourceOrgPolicyPolicyUpdate,
		Delete: resourceOrgPolicyPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceOrgPolicyPolicyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Immutable. The resource name of the Policy. Must be one of the following forms, where constraint_name is the name of the constraint which this Policy configures: * `projects/{project_number}/policies/{constraint_name}` * `folders/{folder_id}/policies/{constraint_name}` * `organizations/{organization_id}/policies/{constraint_name}` For example, \"projects/123/policies/compute.disableSerialPortAccess\". Note: `projects/{project_id}/policies/{constraint_name}` is also an acceptable name for API requests, but responses will return the name using the equivalent project number.",
			},

			"parent": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The parent of the resource.",
			},

			"spec": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Basic information about the Organization Policy.",
				MaxItems:    1,
				Elem:        OrgPolicyPolicySpecSchema(),
			},
		},
	}
}

func OrgPolicyPolicySpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inherit_from_parent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines the inheritance behavior for this `Policy`. If `inherit_from_parent` is true, PolicyRules set higher up in the hierarchy (up to the closest root) are inherited and present in the effective policy. If it is false, then no rules are inherited, and this Policy becomes the new root for evaluation. This field can be set only for Policies which configure list constraints.",
			},

			"reset": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignores policies set above this resource and restores the `constraint_default` enforcement behavior of the specific `Constraint` at this resource. This field can be set in policies for either list or boolean constraints. If set, `rules` must be empty and `inherit_from_parent` must be set to false.",
			},

			"rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Up to 10 PolicyRules are allowed. In Policies for boolean constraints, the following requirements apply: - There must be one and only one PolicyRule where condition is unset. - BooleanPolicyRules with conditions must set `enforced` to the opposite of the PolicyRule without a condition. - During policy evaluation, PolicyRules with conditions that are true for a target resource take precedence.",
				Elem:        OrgPolicyPolicySpecRulesSchema(),
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An opaque tag indicating the current version of the `Policy`, used for concurrency control. This field is ignored if used in a `CreatePolicy` request. When the `Policy` is returned from either a `GetPolicy` or a `ListPolicies` request, this `etag` indicates the version of the current `Policy` to use when executing a read-modify-write loop. When the `Policy` is returned from a `GetEffectivePolicy` request, the `etag` will be unset.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time stamp this was previously updated. This represents the last time a call to `CreatePolicy` or `UpdatePolicy` was made for that `Policy`.",
			},
		},
	}
}

func OrgPolicyPolicySpecRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_all": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Setting this to true means that all values are allowed. This field can be set only in Policies for list constraints.",
			},

			"condition": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A condition which determines whether this rule is used in the evaluation of the policy. When set, the `expression` field in the `Expr' must include from 1 to 10 subexpressions, joined by the \"||\" or \"&&\" operators. Each subexpression must be of the form \"resource.matchTag('/tag_key_short_name, 'tag_value_short_name')\". or \"resource.matchTagId('tagKeys/key_id', 'tagValues/value_id')\". where key_name and value_name are the resource names for Label Keys and Values. These names are available from the Tag Manager Service. An example expression is: \"resource.matchTag('123456789/environment, 'prod')\". or \"resource.matchTagId('tagKeys/123', 'tagValues/456')\".",
				MaxItems:    1,
				Elem:        OrgPolicyPolicySpecRulesConditionSchema(),
			},

			"deny_all": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Setting this to true means that all values are denied. This field can be set only in Policies for list constraints.",
			},

			"enforce": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If `true`, then the `Policy` is enforced. If `false`, then any configuration is acceptable. This field can be set only in Policies for boolean constraints.",
			},

			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of values to be used for this PolicyRule. This field can be set only in Policies for list constraints.",
				MaxItems:    1,
				Elem:        OrgPolicyPolicySpecRulesValuesSchema(),
			},
		},
	}
}

func OrgPolicyPolicySpecRulesConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.",
			},

			"expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Textual representation of an expression in Common Expression Language syntax.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. String indicating the location of the expression for error reporting, e.g. a file name and a position in the file.",
			},

			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Title for the expression, i.e. a short string describing its purpose. This can be used e.g. in UIs which allow to enter the expression.",
			},
		},
	}
}

func OrgPolicyPolicySpecRulesValuesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allowed_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of values allowed at this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"denied_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of values denied at this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceOrgPolicyPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &orgpolicy.Policy{
		Name:   dcl.String(d.Get("name").(string)),
		Parent: dcl.String(d.Get("parent").(string)),
		Spec:   expandOrgPolicyPolicySpec(d.Get("spec")),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLOrgPolicyClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyPolicy(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Policy: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Policy %q: %#v", d.Id(), res)

	return resourceOrgPolicyPolicyRead(d, meta)
}

func resourceOrgPolicyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &orgpolicy.Policy{
		Name:   dcl.String(d.Get("name").(string)),
		Parent: dcl.String(d.Get("parent").(string)),
		Spec:   expandOrgPolicyPolicySpec(d.Get("spec")),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLOrgPolicyClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetPolicy(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("OrgPolicyPolicy %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("parent", res.Parent); err != nil {
		return fmt.Errorf("error setting parent in state: %s", err)
	}
	if err = d.Set("spec", flattenOrgPolicyPolicySpec(res.Spec)); err != nil {
		return fmt.Errorf("error setting spec in state: %s", err)
	}

	return nil
}
func resourceOrgPolicyPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &orgpolicy.Policy{
		Name:   dcl.String(d.Get("name").(string)),
		Parent: dcl.String(d.Get("parent").(string)),
		Spec:   expandOrgPolicyPolicySpec(d.Get("spec")),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLOrgPolicyClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyPolicy(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Policy: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Policy %q: %#v", d.Id(), res)

	return resourceOrgPolicyPolicyRead(d, meta)
}

func resourceOrgPolicyPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &orgpolicy.Policy{
		Name:   dcl.String(d.Get("name").(string)),
		Parent: dcl.String(d.Get("parent").(string)),
		Spec:   expandOrgPolicyPolicySpec(d.Get("spec")),
	}

	log.Printf("[DEBUG] Deleting Policy %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLOrgPolicyClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeletePolicy(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Policy: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Policy %q", d.Id())
	return nil
}

func resourceOrgPolicyPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgdclresource.ResourceOrgPolicyPolicyCustomImport(d, config); err != nil {
		return nil, fmt.Errorf("error encountered in import: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}

func expandOrgPolicyPolicySpec(o interface{}) *orgpolicy.PolicySpec {
	if o == nil {
		return orgpolicy.EmptyPolicySpec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return orgpolicy.EmptyPolicySpec
	}
	obj := objArr[0].(map[string]interface{})
	return &orgpolicy.PolicySpec{
		InheritFromParent: dcl.Bool(obj["inherit_from_parent"].(bool)),
		Reset:             dcl.Bool(obj["reset"].(bool)),
		Rules:             expandOrgPolicyPolicySpecRulesArray(obj["rules"]),
	}
}

func flattenOrgPolicyPolicySpec(obj *orgpolicy.PolicySpec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"inherit_from_parent": obj.InheritFromParent,
		"reset":               obj.Reset,
		"rules":               flattenOrgPolicyPolicySpecRulesArray(obj.Rules),
		"etag":                obj.Etag,
		"update_time":         obj.UpdateTime,
	}

	return []interface{}{transformed}

}
func expandOrgPolicyPolicySpecRulesArray(o interface{}) []orgpolicy.PolicySpecRules {
	if o == nil {
		return make([]orgpolicy.PolicySpecRules, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]orgpolicy.PolicySpecRules, 0)
	}

	items := make([]orgpolicy.PolicySpecRules, 0, len(objs))
	for _, item := range objs {
		i := expandOrgPolicyPolicySpecRules(item)
		items = append(items, *i)
	}

	return items
}

func expandOrgPolicyPolicySpecRules(o interface{}) *orgpolicy.PolicySpecRules {
	if o == nil {
		return orgpolicy.EmptyPolicySpecRules
	}

	obj := o.(map[string]interface{})
	return &orgpolicy.PolicySpecRules{
		AllowAll:  tpgdclresource.ExpandEnumBool(obj["allow_all"].(string)),
		Condition: expandOrgPolicyPolicySpecRulesCondition(obj["condition"]),
		DenyAll:   tpgdclresource.ExpandEnumBool(obj["deny_all"].(string)),
		Enforce:   tpgdclresource.ExpandEnumBool(obj["enforce"].(string)),
		Values:    expandOrgPolicyPolicySpecRulesValues(obj["values"]),
	}
}

func flattenOrgPolicyPolicySpecRulesArray(objs []orgpolicy.PolicySpecRules) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOrgPolicyPolicySpecRules(&item)
		items = append(items, i)
	}

	return items
}

func flattenOrgPolicyPolicySpecRules(obj *orgpolicy.PolicySpecRules) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_all": tpgdclresource.FlattenEnumBool(obj.AllowAll),
		"condition": flattenOrgPolicyPolicySpecRulesCondition(obj.Condition),
		"deny_all":  tpgdclresource.FlattenEnumBool(obj.DenyAll),
		"enforce":   tpgdclresource.FlattenEnumBool(obj.Enforce),
		"values":    flattenOrgPolicyPolicySpecRulesValues(obj.Values),
	}

	return transformed

}

func expandOrgPolicyPolicySpecRulesCondition(o interface{}) *orgpolicy.PolicySpecRulesCondition {
	if o == nil {
		return orgpolicy.EmptyPolicySpecRulesCondition
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return orgpolicy.EmptyPolicySpecRulesCondition
	}
	obj := objArr[0].(map[string]interface{})
	return &orgpolicy.PolicySpecRulesCondition{
		Description: dcl.String(obj["description"].(string)),
		Expression:  dcl.String(obj["expression"].(string)),
		Location:    dcl.String(obj["location"].(string)),
		Title:       dcl.String(obj["title"].(string)),
	}
}

func flattenOrgPolicyPolicySpecRulesCondition(obj *orgpolicy.PolicySpecRulesCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"description": obj.Description,
		"expression":  obj.Expression,
		"location":    obj.Location,
		"title":       obj.Title,
	}

	return []interface{}{transformed}

}

func expandOrgPolicyPolicySpecRulesValues(o interface{}) *orgpolicy.PolicySpecRulesValues {
	if o == nil {
		return orgpolicy.EmptyPolicySpecRulesValues
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return orgpolicy.EmptyPolicySpecRulesValues
	}
	obj := objArr[0].(map[string]interface{})
	return &orgpolicy.PolicySpecRulesValues{
		AllowedValues: tpgdclresource.ExpandStringArray(obj["allowed_values"]),
		DeniedValues:  tpgdclresource.ExpandStringArray(obj["denied_values"]),
	}
}

func flattenOrgPolicyPolicySpecRulesValues(obj *orgpolicy.PolicySpecRulesValues) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allowed_values": obj.AllowedValues,
		"denied_values":  obj.DeniedValues,
	}

	return []interface{}{transformed}

}
