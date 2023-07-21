// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/compute/v1"
)

func ResourceComputeSecurityPolicy() *schema.Resource {
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
			Create: schema.DefaultTimeout(8 * time.Minute),
			Update: schema.DefaultTimeout(8 * time.Minute),
			Delete: schema.DefaultTimeout(8 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateGCEName,
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

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  `The type indicates the intended use of the security policy. CLOUD_ARMOR - Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers. CLOUD_ARMOR_EDGE - Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.`,
				ValidateFunc: validation.StringInSlice([]string{"CLOUD_ARMOR", "CLOUD_ARMOR_EDGE", "CLOUD_ARMOR_INTERNAL_SERVICE"}, false),
			},

			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // If no rules are set, a default rule is added
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Action to take when match matches the request.`,
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
										Default:      "",
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
							Default:     "",
							Optional:    true,
							Description: `An optional description of this rule. Max size is 64.`,
						},

						"preview": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `When set to true, the action specified above is not enforced. Stackdriver logs for requests that trigger a preview action are annotated as such.`,
						},

						"rate_limit_options": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Rate limit threshold for this security policy. Must be specified if the action is "rate_based_ban" or "throttle". Cannot be specified for any other actions.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rate_limit_threshold": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `Threshold at which to begin ratelimiting.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Number of HTTP(S) requests for calculating the threshold.`,
												},

												"interval_sec": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Interval over which the threshold is computed.`,
												},
											},
										},
									},

									"conform_action": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"allow"}, false),
										Description:  `Action to take for requests that are under the configured rate limit threshold. Valid option is "allow" only.`,
									},

									"exceed_action": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"redirect", "deny(403)", "deny(404)", "deny(429)", "deny(502)"}, false),
										Description:  `Action to take for requests that are above the configured rate limit threshold, to either deny with a specified HTTP response code, or redirect to a different endpoint. Valid options are "deny()" where valid values for status are 403, 404, 429, and 502, and "redirect" where the redirect parameters come from exceedRedirectOptions below.`,
									},

									"enforce_on_key": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "ALL",
										Description:  `Determines the key to enforce the rateLimitThreshold on`,
										ValidateFunc: validation.StringInSlice([]string{"ALL", "IP", "HTTP_HEADER", "XFF_IP", "HTTP_COOKIE", "HTTP_PATH", "SNI", "REGION_CODE", ""}, false),
									},

									"enforce_on_key_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Rate limit key name applicable only for the following key types: HTTP_HEADER -- Name of the HTTP header whose value is taken as the key value. HTTP_COOKIE -- Name of the HTTP cookie whose value is taken as the key value.`,
									},

									"ban_threshold": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Can only be specified if the action for the rule is "rate_based_ban". If specified, the key will be banned for the configured 'banDurationSec' when the number of requests that exceed the 'rateLimitThreshold' also exceed this 'banThreshold'.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Number of HTTP(S) requests for calculating the threshold.`,
												},

												"interval_sec": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Interval over which the threshold is computed.`,
												},
											},
										},
									},

									"ban_duration_sec": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Can only be specified if the action for the rule is "rate_based_ban". If specified, determines the time (in seconds) the traffic will continue to be banned by the rate limit after the rate falls below the threshold.`,
									},

									"exceed_redirect_options": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Parameters defining the redirect action that is used as the exceed action. Cannot be specified if the exceed action is not redirect.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:         schema.TypeString,
													Required:     true,
													Description:  `Type of the redirect action.`,
													ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_302", "GOOGLE_RECAPTCHA"}, false),
												},

												"target": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `Target for the redirect action. This is required if the type is EXTERNAL_302 and cannot be specified for GOOGLE_RECAPTCHA.`,
												},
											},
										},
									},
								},
							},
						},

						"redirect_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_302", "GOOGLE_RECAPTCHA"}, false),
										Description:  `Type of the redirect action. Available options: EXTERNAL_302: Must specify the corresponding target field in config. GOOGLE_RECAPTCHA: Cannot specify target field in config.`,
									},

									"target": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Target for the redirect action. This is required if the type is EXTERNAL_302 and cannot be specified for GOOGLE_RECAPTCHA.`,
									},
								},
							},
							Description: `Parameters defining the redirect action. Cannot be specified for any other actions.`,
						},
						"header_action": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Additional actions that are performed on headers.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"request_headers_to_adds": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `The list of request headers to add or overwrite if they're already present.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"header_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The name of the header to set.`,
												},
												"header_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The value to set the named header to.`,
												},
											},
										},
									},
								},
							},
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

			"advanced_options_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Advanced Options Config of this security policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"json_parsing": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"DISABLED", "STANDARD"}, false),
							Description:  `JSON body parsing. Supported values include: "DISABLED", "STANDARD".`,
						},
						"json_custom_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `Custom configuration to apply the JSON parsing. Only applicable when JSON parsing is set to STANDARD.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_types": {
										Type:        schema.TypeSet,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `A list of custom Content-Type header values to apply the JSON parsing.`,
									},
								},
							},
						},
						"log_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"NORMAL", "VERBOSE"}, false),
							Description:  `Logging level. Supported values include: "NORMAL", "VERBOSE".`,
						},
					},
				},
			},

			"adaptive_protection_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Adaptive Protection Config of this security policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"layer_7_ddos_defense_config": {
							Type:        schema.TypeList,
							Description: `Layer 7 DDoS Defense Config of this security policy`,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `If set to true, enables CAAP for L7 DDoS detection.`,
									},
									"rule_visibility": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "STANDARD",
										ValidateFunc: validation.StringInSlice([]string{"STANDARD", "PREMIUM"}, false),
										Description:  `Rule visibility. Supported values include: "STANDARD", "PREMIUM".`,
									},
								},
							},
						},
					},
				},
			},
			"recaptcha_options_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `reCAPTCHA configuration options to be applied for the security policy.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_site_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `A field to supply a reCAPTCHA site key to be used for all the rules using the redirect action with the type of GOOGLE_RECAPTCHA under the security policy. The specified site key needs to be created from the reCAPTCHA API. The user is responsible for the validity of the specified site key. If not specified, a Google-managed site key is used.`,
						},
					},
				},
			},
		},

		UseJSONNumber: true,
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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)
	securityPolicy := &compute.SecurityPolicy{
		Name:        sp,
		Description: d.Get("description").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		securityPolicy.Type = v.(string)
	}

	if v, ok := d.GetOk("rule"); ok {
		securityPolicy.Rules = expandSecurityPolicyRules(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("advanced_options_config"); ok {
		securityPolicy.AdvancedOptionsConfig = expandSecurityPolicyAdvancedOptionsConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("adaptive_protection_config"); ok {
		securityPolicy.AdaptiveProtectionConfig = expandSecurityPolicyAdaptiveProtectionConfig(v.([]interface{}))
	}

	log.Printf("[DEBUG] SecurityPolicy insert request: %#v", securityPolicy)

	if v, ok := d.GetOk("recaptcha_options_config"); ok {
		securityPolicy.RecaptchaOptionsConfig = expandSecurityPolicyRecaptchaOptionsConfig(v.([]interface{}), d)
	}

	client := config.NewComputeClient(userAgent)

	op, err := client.SecurityPolicies.Insert(project, securityPolicy).Do()

	if err != nil {
		return errwrap.Wrapf("Error creating SecurityPolicy: {{err}}", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/global/securityPolicies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = ComputeOperationWaitTime(config, op, project, fmt.Sprintf("Creating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeSecurityPolicyRead(d, meta)
}

func resourceComputeSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	client := config.NewComputeClient(userAgent)

	securityPolicy, err := client.SecurityPolicies.Get(project, sp).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SecurityPolicy %q", d.Id()))
	}

	if err := d.Set("name", securityPolicy.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", securityPolicy.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("type", securityPolicy.Type); err != nil {
		return fmt.Errorf("Error setting type: %s", err)
	}
	if err := d.Set("rule", flattenSecurityPolicyRules(securityPolicy.Rules)); err != nil {
		return err
	}
	if err := d.Set("fingerprint", securityPolicy.Fingerprint); err != nil {
		return fmt.Errorf("Error setting fingerprint: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(securityPolicy.SelfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("advanced_options_config", flattenSecurityPolicyAdvancedOptionsConfig(securityPolicy.AdvancedOptionsConfig)); err != nil {
		return fmt.Errorf("Error setting advanced_options_config: %s", err)
	}

	if err := d.Set("adaptive_protection_config", flattenSecurityPolicyAdaptiveProtectionConfig(securityPolicy.AdaptiveProtectionConfig)); err != nil {
		return fmt.Errorf("Error setting adaptive_protection_config: %s", err)
	}

	if err := d.Set("recaptcha_options_config", flattenSecurityPolicyRecaptchaOptionConfig(securityPolicy.RecaptchaOptionsConfig)); err != nil {
		return fmt.Errorf("Error setting recaptcha_options_config: %s", err)
	}

	return nil
}

func resourceComputeSecurityPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	securityPolicy := &compute.SecurityPolicy{
		Fingerprint: d.Get("fingerprint").(string),
	}

	if d.HasChange("type") {
		securityPolicy.Type = d.Get("type").(string)
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "Type")
	}

	if d.HasChange("description") {
		securityPolicy.Description = d.Get("description").(string)
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "Description")
	}

	if d.HasChange("advanced_options_config") {
		securityPolicy.AdvancedOptionsConfig = expandSecurityPolicyAdvancedOptionsConfig(d.Get("advanced_options_config").([]interface{}))
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "AdvancedOptionsConfig", "advancedOptionsConfig.jsonParsing", "advancedOptionsConfig.jsonCustomConfig", "advancedOptionsConfig.logLevel")
	}

	if d.HasChange("adaptive_protection_config") {
		securityPolicy.AdaptiveProtectionConfig = expandSecurityPolicyAdaptiveProtectionConfig(d.Get("adaptive_protection_config").([]interface{}))
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "AdaptiveProtectionConfig", "adaptiveProtectionConfig.layer7DdosDefenseConfig.enable", "adaptiveProtectionConfig.layer7DdosDefenseConfig.ruleVisibility")
	}

	if d.HasChange("recaptcha_options_config") {
		securityPolicy.RecaptchaOptionsConfig = expandSecurityPolicyRecaptchaOptionsConfig(d.Get("recaptcha_options_config").([]interface{}), d)
		securityPolicy.ForceSendFields = append(securityPolicy.ForceSendFields, "RecaptchaOptionsConfig")
	}

	if len(securityPolicy.ForceSendFields) > 0 {
		client := config.NewComputeClient(userAgent)

		op, err := client.SecurityPolicies.Patch(project, sp, securityPolicy).Do()

		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
		}

		err = ComputeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
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
				client := config.NewComputeClient(userAgent)
				// If the rule is in new and its priority does not exist in old, then add it.
				op, err := client.SecurityPolicies.AddRule(project, sp, expandSecurityPolicyRule(rule)).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = ComputeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			} else if !oSet.Contains(rule) {
				client := config.NewComputeClient(userAgent)

				// If the rule is in new, and its priority is in old, but its hash is different than the one in old, update it.
				op, err := client.SecurityPolicies.PatchRule(project, sp, expandSecurityPolicyRule(rule)).Priority(priority).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = ComputeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}

		for _, rule := range oSet.List() {
			priority := int64(rule.(map[string]interface{})["priority"].(int))
			if !nPriorities[priority] {
				client := config.NewComputeClient(userAgent)

				// If the rule's priority is in old but not new, remove it.
				op, err := client.SecurityPolicies.RemoveRule(project, sp).Priority(priority).Do()

				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("Error updating SecurityPolicy %q: {{err}}", sp), err)
				}

				err = ComputeOperationWaitTime(config, op, project, fmt.Sprintf("Updating SecurityPolicy %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceComputeSecurityPolicyRead(d, meta)
}

func resourceComputeSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	client := config.NewComputeClient(userAgent)

	// Delete the SecurityPolicy
	op, err := client.SecurityPolicies.Delete(project, d.Get("name").(string)).Do()
	if err != nil {
		return errwrap.Wrapf("Error deleting SecurityPolicy: {{err}}", err)
	}

	err = ComputeOperationWaitTime(config, op, project, "Deleting SecurityPolicy", userAgent, d.Timeout(schema.TimeoutDelete))
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
		Description:      data["description"].(string),
		Priority:         int64(data["priority"].(int)),
		Action:           data["action"].(string),
		Preview:          data["preview"].(bool),
		Match:            expandSecurityPolicyMatch(data["match"].([]interface{})),
		RateLimitOptions: expandSecurityPolicyRuleRateLimitOptions(data["rate_limit_options"].([]interface{})),
		RedirectOptions:  expandSecurityPolicyRuleRedirectOptions(data["redirect_options"].([]interface{})),
		HeaderAction:     expandSecurityPolicyRuleHeaderAction(data["header_action"].([]interface{})),
		ForceSendFields:  []string{"Description", "Preview"},
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
		SrcIpRanges: tpgresource.ConvertStringArr(data["src_ip_ranges"].(*schema.Set).List()),
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
			"description":        rule.Description,
			"priority":           rule.Priority,
			"action":             rule.Action,
			"preview":            rule.Preview,
			"match":              flattenMatch(rule.Match),
			"rate_limit_options": flattenSecurityPolicyRuleRateLimitOptions(rule.RateLimitOptions),
			"redirect_options":   flattenSecurityPolicyRedirectOptions(rule.RedirectOptions),
			"header_action":      flattenSecurityPolicyRuleHeaderAction(rule.HeaderAction),
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
		"src_ip_ranges": schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(conf.SrcIpRanges)),
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

func expandSecurityPolicyAdvancedOptionsConfig(configured []interface{}) *compute.SecurityPolicyAdvancedOptionsConfig {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyAdvancedOptionsConfig{
		JsonParsing:      data["json_parsing"].(string),
		JsonCustomConfig: expandSecurityPolicyAdvancedOptionsConfigJsonCustomConfig(data["json_custom_config"].([]interface{})),
		LogLevel:         data["log_level"].(string),
	}
}

func flattenSecurityPolicyAdvancedOptionsConfig(conf *compute.SecurityPolicyAdvancedOptionsConfig) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"json_parsing":       conf.JsonParsing,
		"json_custom_config": flattenSecurityPolicyAdvancedOptionsConfigJsonCustomConfig(conf.JsonCustomConfig),
		"log_level":          conf.LogLevel,
	}

	return []map[string]interface{}{data}
}

func expandSecurityPolicyAdvancedOptionsConfigJsonCustomConfig(configured []interface{}) *compute.SecurityPolicyAdvancedOptionsConfigJsonCustomConfig {
	if len(configured) == 0 || configured[0] == nil {
		// If configuration is unset, return an empty JsonCustomConfig; this ensures the ContentTypes list can be cleared
		return &compute.SecurityPolicyAdvancedOptionsConfigJsonCustomConfig{}
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyAdvancedOptionsConfigJsonCustomConfig{
		ContentTypes: tpgresource.ConvertStringArr(data["content_types"].(*schema.Set).List()),
	}
}

func flattenSecurityPolicyAdvancedOptionsConfigJsonCustomConfig(conf *compute.SecurityPolicyAdvancedOptionsConfigJsonCustomConfig) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"content_types": schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface(conf.ContentTypes)),
	}

	return []map[string]interface{}{data}
}

func expandSecurityPolicyAdaptiveProtectionConfig(configured []interface{}) *compute.SecurityPolicyAdaptiveProtectionConfig {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyAdaptiveProtectionConfig{
		Layer7DdosDefenseConfig: expandLayer7DdosDefenseConfig(data["layer_7_ddos_defense_config"].([]interface{})),
	}
}

func expandLayer7DdosDefenseConfig(configured []interface{}) *compute.SecurityPolicyAdaptiveProtectionConfigLayer7DdosDefenseConfig {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyAdaptiveProtectionConfigLayer7DdosDefenseConfig{
		Enable:          data["enable"].(bool),
		RuleVisibility:  data["rule_visibility"].(string),
		ForceSendFields: []string{"Enable"},
	}
}

func flattenSecurityPolicyAdaptiveProtectionConfig(conf *compute.SecurityPolicyAdaptiveProtectionConfig) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"layer_7_ddos_defense_config": flattenLayer7DdosDefenseConfig(conf.Layer7DdosDefenseConfig),
	}

	return []map[string]interface{}{data}
}

func flattenLayer7DdosDefenseConfig(conf *compute.SecurityPolicyAdaptiveProtectionConfigLayer7DdosDefenseConfig) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"enable":          conf.Enable,
		"rule_visibility": conf.RuleVisibility,
	}

	return []map[string]interface{}{data}
}

func expandSecurityPolicyRuleRateLimitOptions(configured []interface{}) *compute.SecurityPolicyRuleRateLimitOptions {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyRuleRateLimitOptions{
		BanThreshold:          expandThreshold(data["ban_threshold"].([]interface{})),
		RateLimitThreshold:    expandThreshold(data["rate_limit_threshold"].([]interface{})),
		ExceedAction:          data["exceed_action"].(string),
		ConformAction:         data["conform_action"].(string),
		EnforceOnKey:          data["enforce_on_key"].(string),
		EnforceOnKeyName:      data["enforce_on_key_name"].(string),
		BanDurationSec:        int64(data["ban_duration_sec"].(int)),
		ExceedRedirectOptions: expandSecurityPolicyRuleRedirectOptions(data["exceed_redirect_options"].([]interface{})),
		ForceSendFields:       []string{"EnforceOnKey", "EnforceOnKeyName"},
	}
}

func expandThreshold(configured []interface{}) *compute.SecurityPolicyRuleRateLimitOptionsThreshold {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyRuleRateLimitOptionsThreshold{
		Count:       int64(data["count"].(int)),
		IntervalSec: int64(data["interval_sec"].(int)),
	}
}

func flattenSecurityPolicyRuleRateLimitOptions(conf *compute.SecurityPolicyRuleRateLimitOptions) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"ban_threshold":           flattenThreshold(conf.BanThreshold),
		"rate_limit_threshold":    flattenThreshold(conf.RateLimitThreshold),
		"exceed_action":           conf.ExceedAction,
		"conform_action":          conf.ConformAction,
		"enforce_on_key":          conf.EnforceOnKey,
		"enforce_on_key_name":     conf.EnforceOnKeyName,
		"ban_duration_sec":        conf.BanDurationSec,
		"exceed_redirect_options": flattenSecurityPolicyRedirectOptions(conf.ExceedRedirectOptions),
	}

	return []map[string]interface{}{data}
}

func flattenThreshold(conf *compute.SecurityPolicyRuleRateLimitOptionsThreshold) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"count":        conf.Count,
		"interval_sec": conf.IntervalSec,
	}

	return []map[string]interface{}{data}
}

func expandSecurityPolicyRuleRedirectOptions(configured []interface{}) *compute.SecurityPolicyRuleRedirectOptions {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &compute.SecurityPolicyRuleRedirectOptions{
		Type:   data["type"].(string),
		Target: data["target"].(string),
	}
}

func flattenSecurityPolicyRedirectOptions(conf *compute.SecurityPolicyRuleRedirectOptions) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"type":   conf.Type,
		"target": conf.Target,
	}

	return []map[string]interface{}{data}
}

func expandSecurityPolicyRecaptchaOptionsConfig(configured []interface{}, d *schema.ResourceData) *compute.SecurityPolicyRecaptchaOptionsConfig {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})

	return &compute.SecurityPolicyRecaptchaOptionsConfig{
		RedirectSiteKey: data["redirect_site_key"].(string),
		ForceSendFields: []string{"RedirectSiteKey"},
	}
}

func flattenSecurityPolicyRecaptchaOptionConfig(conf *compute.SecurityPolicyRecaptchaOptionsConfig) []map[string]interface{} {
	if conf == nil {
		return nil
	}

	data := map[string]interface{}{
		"redirect_site_key": conf.RedirectSiteKey,
	}

	return []map[string]interface{}{data}
}

func expandSecurityPolicyRuleHeaderAction(configured []interface{}) *compute.SecurityPolicyRuleHttpHeaderAction {
	if len(configured) == 0 || configured[0] == nil {
		// If header action is unset, return an empty object; this ensures the header action can be cleared
		return &compute.SecurityPolicyRuleHttpHeaderAction{}
	}

	data := configured[0].(map[string]interface{})

	return &compute.SecurityPolicyRuleHttpHeaderAction{
		RequestHeadersToAdds: expandSecurityPolicyRequestHeadersToAdds(data["request_headers_to_adds"].([]interface{})),
	}
}

func expandSecurityPolicyRequestHeadersToAdds(configured []interface{}) []*compute.SecurityPolicyRuleHttpHeaderActionHttpHeaderOption {
	transformed := make([]*compute.SecurityPolicyRuleHttpHeaderActionHttpHeaderOption, 0, len(configured))

	for _, raw := range configured {
		transformed = append(transformed, expandSecurityPolicyRequestHeader(raw))
	}

	return transformed
}

func expandSecurityPolicyRequestHeader(configured interface{}) *compute.SecurityPolicyRuleHttpHeaderActionHttpHeaderOption {
	data := configured.(map[string]interface{})

	return &compute.SecurityPolicyRuleHttpHeaderActionHttpHeaderOption{
		HeaderName:  data["header_name"].(string),
		HeaderValue: data["header_value"].(string),
	}
}

func flattenSecurityPolicyRuleHeaderAction(conf *compute.SecurityPolicyRuleHttpHeaderAction) []map[string]interface{} {
	if conf == nil || conf.RequestHeadersToAdds == nil {
		return nil
	}

	transformed := map[string]interface{}{
		"request_headers_to_adds": flattenSecurityPolicyRequestHeadersToAdds(conf.RequestHeadersToAdds),
	}

	return []map[string]interface{}{transformed}
}

func flattenSecurityPolicyRequestHeadersToAdds(conf []*compute.SecurityPolicyRuleHttpHeaderActionHttpHeaderOption) []map[string]interface{} {
	if conf == nil || len(conf) == 0 {
		return nil
	}

	transformed := make([]map[string]interface{}, 0, len(conf))
	for _, raw := range conf {
		transformed = append(transformed, flattenSecurityPolicyRequestHeader(raw))
	}

	return transformed
}

func flattenSecurityPolicyRequestHeader(conf *compute.SecurityPolicyRuleHttpHeaderActionHttpHeaderOption) map[string]interface{} {
	if conf == nil {
		return nil
	}

	return map[string]interface{}{
		"header_name":  conf.HeaderName,
		"header_value": conf.HeaderValue,
	}
}

func resourceSecurityPolicyStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{"projects/(?P<project>[^/]+)/global/securityPolicies/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/global/securityPolicies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
