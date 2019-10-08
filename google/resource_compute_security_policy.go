package google

import (
	"fmt"
	"log"

	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/compute/v0.beta"
)

func resourceComputeSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSecurityPolicyCreate,
		Read:   resourceComputeSecurityPolicyRead,
		Update: resourceComputeSecurityPolicyUpdate,
		Delete: resourceComputeSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
						},

						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"src_ip_ranges": {
													Type:     schema.TypeSet,
													Required: true,
													MinItems: 1,
													MaxItems: 5,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},

									"versioned_expr": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"SRC_IPS_V1"}, false),
									},
								},
							},
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"preview": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
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

	d.SetId(securityPolicy.Name)

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), fmt.Sprintf("Creating SecurityPolicy %q", sp))
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

	securityPolicy, err := config.clientComputeBeta.SecurityPolicies.Get(project, d.Id()).Do()
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

	sp := d.Id()

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

		err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), fmt.Sprintf("Updating SecurityPolicy %q", sp))
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

				err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), fmt.Sprintf("Updating SecurityPolicy %q", sp))
				if err != nil {
					return err
				}
			} else if !oSet.Contains(rule) {
				// If the rule is in new, and its priority is in old, but its hash is different than the one in old, update it.
				op, err := config.clientComputeBeta.SecurityPolicies.PatchRule(project, sp, expandSecurityPolicyRule(rule)).Priority(priority).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), fmt.Sprintf("Updating SecurityPolicy %q", sp))
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

				err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), fmt.Sprintf("Updating SecurityPolicy %q", sp))
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
	op, err := config.clientComputeBeta.SecurityPolicies.Delete(project, d.Id()).Do()
	if err != nil {
		return errwrap.Wrapf("Error deleting SecurityPolicy: {{err}}", err)
	}

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutDelete).Minutes()), "Deleting SecurityPolicy")
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

func flattenSecurityPolicyRules(rules []*compute.SecurityPolicyRule) []map[string]interface{} {
	rulesSchema := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		data := map[string]interface{}{
			"description": rule.Description,
			"priority":    rule.Priority,
			"action":      rule.Action,
			"preview":     rule.Preview,
			"match": []map[string]interface{}{
				{
					"versioned_expr": rule.Match.VersionedExpr,
					"config": []map[string]interface{}{
						{
							"src_ip_ranges": schema.NewSet(schema.HashString, convertStringArrToInterface(rule.Match.Config.SrcIpRanges)),
						},
					},
				},
			},
		}

		rulesSchema = append(rulesSchema, data)
	}
	return rulesSchema
}
