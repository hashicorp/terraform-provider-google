package google

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	compute "google.golang.org/api/compute/v0.beta"
)

func resourceComputeSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSecurityPolicyCreate,
		Read:   resourceComputeSecurityPolicyRead,
		Update: resourceComputeSecurityPolicyUpdate,
		Delete: resourceComputeSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityPolicyStateImporter,
		},
		CustomizeDiff: rulesCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
				Description:  `The name of the security policy.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `An optional description of this security policy. Max size is 2048.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // If no rules are set, a default rule is added
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"allow", "deny(403)", "deny(404)", "deny(502)"}, false),
							Description:  `Action to take when match matches the request. Valid values:   "allow" : allow access to target, "deny(status)" : deny access to target, returns the HTTP response code specified (valid values are 403, 404 and 502)`,
						},

						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `An unique positive integer indicating the priority of evaluation for a rule. Rules are evaluated from highest priority (lowest numerically) to lowest priority (highest numerically) in order.`,
						},

						"match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"src_ip_ranges": {
													Type:        schema.TypeSet,
													Required:    true,
													MinItems:    1,
													MaxItems:    10,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: `Set of IP addresses or ranges (IPV4 or IPV6) in CIDR notation to match against inbound traffic. There is a limit of 10 IP ranges per rule. A value of '*' matches all IPs (can be used to override the default behavior).`,
												},
											},
										},
										Description: `The configuration options available when specifying versioned_expr. This field must be specified if versioned_expr is specified and cannot be specified if versioned_expr is not specified.`,
									},

									"versioned_expr": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"SRC_IPS_V1"}, false),
										Description:  `Predefined rule expression. If this field is specified, config must also be specified. Available options:   SRC_IPS_V1: Must specify the corresponding src_ip_ranges field in config.`,
									},

									"expr": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"expression": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Textual representation of an expression in Common Expression Language syntax. The application context of the containing message determines which well-known feature set of CEL is supported.`,
												},
												// These fields are not yet supported (Issue hashicorp/terraform-provider-google#4497: mbang)
												// "title": {
												// 	Type:     schema.TypeString,
												// 	Optional: true,
												// },
												// "description": {
												// 	Type:     schema.TypeString,
												// 	Optional: true,
												// },
												// "location": {
												// 	Type:     schema.TypeString,
												// 	Optional: true,
												// },
											},
										},
										Description: `User defined CEVAL expression. A CEVAL expression is used to specify match criteria such as origin.ip, source.region_code and contents in the request header.`,
									},
								},
							},
							Description: `A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding action is enforced.`,
						},

						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `An optional description of this rule. Max size is 64.`,
						},

						"preview": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `When set to true, the action specified above is not enforced. Stackdriver logs for requests that trigger a preview action are annotated as such.`,
						},
					},
				},
				Description: `The set of rules that belong to this policy. There must always be a default rule (rule with priority 2147483647 and match "*"). If no rules are provided when creating a security policy, a default rule with action "allow" will be added.`,
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Fingerprint of this resource.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},
		},
	}
}

func rulesCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	_, n := diff.GetChange("rule")
	nSet := n.(*schema.Set)

	nPriorities := map[int64]bool{}
	for _, rule := range nSet.List() {
		priority := int64(rule.(map[string]interface{})["priority"].(int))
		if nPriorities[priority] {
			return fmt.Errorf("Two rules have the same priority, please update one of the priorities to be different.")
		}
		nPriorities[priority] = true
	}

	return nil
}

func resourceComputeSecurityPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)
	securityPolicy := &compute.SecurityPolicy{
		Name:        sp,
		Description: d.Get("description").(string),
	}
	if v, ok := d.GetOk("rule"); ok {
		securityPolicy.Rules = expandSecurityPolicyRules(v.(*schema.Set).List())
	}

	log.Printf("[DEBUG] SecurityPolicy insert request: %#v", securityPolicy)

	op, err := config.clientComputeBeta.SecurityPolicies.Insert(project, securityPolicy).Do()

	if err != nil {
		return errwrap.Wrapf("Error creating SecurityPolicy: {{err}}", err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/global/securityPolicies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Creating SecurityPolicy %q", sp), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeSecurityPolicyRead(d, meta)
}

func resourceComputeSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)
	securityPolicy, err := config.clientComputeBeta.SecurityPolicies.Get(project, sp).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SecurityPolicy %q", d.Id()))
	}

	d.Set("name", securityPolicy.Name)
	d.Set("description", securityPolicy.Description)
	if err := d.Set("rule", flattenSecurityPolicyRules(securityPolicy.Rules)); err != nil {
		return err
	}
	d.Set("fingerprint", securityPolicy.Fingerprint)
	d.Set("project", project)
	d.Set("self_link", ConvertSelfLinkToV1(securityPolicy.SelfLink))

	return nil
}

func resourceComputeSecurityPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	if d.HasChange("description") {
		securityPolicy := &compute.SecurityPolicy{
			Description:     d.Get("description").(string),
			Fingerprint:     d.Get("fingerprint").(string),
			ForceSendFields: []string{"Description"},
		}
		op, err := config.clientComputeBeta.SecurityPolicies.Patch(project, sp, securityPolicy).Do()

		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
		}

		err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		oSet := o.(*schema.Set)
		nSet := n.(*schema.Set)

		oPriorities := map[int64]bool{}
		nPriorities := map[int64]bool{}
		for _, rule := range oSet.List() {
			oPriorities[int64(rule.(map[string]interface{})["priority"].(int))] = true
		}

		for _, rule := range nSet.List() {
			priority := int64(rule.(map[string]interface{})["priority"].(int))
			nPriorities[priority] = true
			if !oPriorities[priority] {
				// If the rule is in new and its priority does not exist in old, then add it.
				op, err := config.clientComputeBeta.SecurityPolicies.AddRule(project, sp, expandSecurityPolicyRule(rule)).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			} else if !oSet.Contains(rule) {
				// If the rule is in new, and its priority is in old, but its hash is different than the one in old, update it.
				op, err := config.clientComputeBeta.SecurityPolicies.PatchRule(project, sp, expandSecurityPolicyRule(rule)).Priority(priority).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}

		for _, rule := range oSet.List() {
			priority := int64(rule.(map[string]interface{})["priority"].(int))
			if !nPriorities[priority] {
				// If the rule's priority is in old but not new, remove it.
				op, err := config.clientComputeBeta.SecurityPolicies.RemoveRule(project, sp).Priority(priority).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceComputeSecurityPolicyRead(d, meta)
}

func resourceComputeSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the SecurityPolicy
	op, err := config.clientComputeBeta.SecurityPolicies.Delete(project, d.Get("name").(string)).Do()
	if err != nil {
		return errwrap.Wrapf("Error deleting SecurityPolicy: {{err}}", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting SecurityPolicy", d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func expandSecurityPolicyRules(configured []interface{}) []*compute.SecurityPolicyRule {
	rules := make([]*compute.SecurityPolicyRule, 0, len(configured))
	for _, raw := range configured {
		rules = append(rules, expandSecurityPolicyRule(raw))
	}
	return rules
}

func expandSecurityPolicyRule(raw interface{}) *compute.SecurityPolicyRule {
	data := raw.(map[string]interface{})
	return &compute.SecurityPolicyRule{
		Description:     data["description"].(string),
		Priority:        int64(data["priority"].(int)),
		Action:          data["action"].(string),
		Preview:         data["preview"].(bool),
		Match:           expandSecurityPolicyMatch(data["match"].([]interface{})),
		ForceSendFields: []string{"Description", "Preview"},
	}
}

func expandSecurityPolicyMatch(configured []interface{}) *compute.SecurityPolicyRuleMatcher {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyRuleMatcher{
		VersionedExpr: data["versioned_expr"].(string),
		Config:        expandSecurityPolicyMatchConfig(data["config"].([]interface{})),
		Expr:          expandSecurityPolicyMatchExpr(data["expr"].([]interface{})),
	}
}

func expandSecurityPolicyMatchConfig(configured []interface{}) *compute.SecurityPolicyRuleMatcherConfig {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyRuleMatcherConfig{
		SrcIpRanges: convertStringArr(data["src_ip_ranges"].(*schema.Set).List()),
	}
}

func expandSecurityPolicyMatchExpr(expr []interface{}) *compute.Expr {
	if len(expr) == 0 || expr[0] == nil {
		return nil
	}

	data := expr[0].(map[string]interface{})
	return &compute.Expr{
		Expression: data["expression"].(string),
		// These fields are not yet supported  (Issue hashicorp/terraform-provider-google#4497: mbang)
		// Title:       data["title"].(string),
		// Description: data["description"].(string),
		// Location:    data["location"].(string),
	}
}

func flattenSecurityPolicyRules(rules []*compute.SecurityPolicyRule) []map[string]interface{} {
	rulesSchema := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		data := map[string]interface{}{
			"description": rule.Description,
			"priority":    rule.Priority,
			"action":      rule.Action,
			"preview":     rule.Preview,
			"match":       flattenMatch(rule.Match),
		}

		rulesSchema = append(rulesSchema, data)
	}
	return rulesSchema
}

func flattenMatch(match *compute.SecurityPolicyRuleMatcher) []map[string]interface{} {
	if match == nil {
		return nil
	}

	data := map[string]interface{}{
		"versioned_expr": match.VersionedExpr,
		"config":         flattenMatchConfig(match.Config),
		"expr":           flattenMatchExpr(match),
	}

	return []map[string]interface{}{data}
}

func flattenMatchConfig(conf *compute.SecurityPolicyRuleMatcherConfig) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"src_ip_ranges": schema.NewSet(schema.HashString, convertStringArrToInterface(conf.SrcIpRanges)),
	}

	return []map[string]interface{}{data}
}

func flattenMatchExpr(match *compute.SecurityPolicyRuleMatcher) []map[string]interface{} {
	if match.Expr == nil {
		return nil
	}

	data := map[string]interface{}{
		"expression": match.Expr.Expression,
		// These fields are not yet supported (Issue hashicorp/terraform-provider-google#4497: mbang)
		// "title":       match.Expr.Title,
		// "description": match.Expr.Description,
		// "location":    match.Expr.Location,
	}

	return []map[string]interface{}{data}
}

func resourceSecurityPolicyStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/global/securityPolicies/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/global/securityPolicies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
