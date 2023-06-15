// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// DataSourceGoogleIamPolicy returns a *schema.Resource that allows a customer
// to express a Google Cloud IAM policy in a data resource. This is an example
// of how the schema would be used in a config:
//
//	data "google_iam_policy" "admin" {
//	  binding {
//	    role = "roles/storage.objectViewer"
//	    members = [
//	      "user:evanbrown@google.com",
//	    ]
//	  }
//	}
func DataSourceGoogleIamPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleIamPolicyRead,
		Schema: map[string]*schema.Schema{
			"binding": {
				Type: schema.TypeSet,
				// Binding is optional because a user may want to set an IAM policy with no bindings
				// This allows users to ensure that no bindings were created outside of terraform
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:     schema.TypeString,
							Required: true,
						},
						"members": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile("^deleted:"), "Terraform does not support IAM policies for deleted principals"),
							},
							Set: schema.HashString,
						},
						"condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expression": {
										Type:     schema.TypeString,
										Required: true,
									},
									"title": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"policy_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audit_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"audit_log_configs": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"exempted_members": {
										Type:     schema.TypeSet,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// dataSourceGoogleIamPolicyRead reads a data source from config and writes it
// to state.
func dataSourceGoogleIamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	var policy cloudresourcemanager.Policy
	var bindings []*cloudresourcemanager.Binding

	// The schema supports multiple binding{} blocks
	bset := d.Get("binding").(*schema.Set)
	aset := d.Get("audit_config").(*schema.Set)

	// Convert each config binding into a cloudresourcemanager.Binding
	// and merge member lists of equivalent binding{} blocks from the config provided by the user
	bindingMap := map[string]*cloudresourcemanager.Binding{}
	for _, v := range bset.List() {
		binding := v.(map[string]interface{})
		members := tpgresource.ConvertStringSet(binding["members"].(*schema.Set))
		condition := tpgiamresource.ExpandIamCondition(binding["condition"])

		// Map keys are used to identify binding{} blocks that are identical except for the member lists
		key := binding["role"].(string)
		if condition != nil {
			key += fmt.Sprintf("-[%s]-[%s]-[%s]-[%s]", condition.Expression, condition.Title, condition.Description, condition.Location)
		}

		if val, ok := bindingMap[key]; ok {
			// Add members to existing cloudresourcemanager.Binding in the map
			m := append(val.Members, members...)
			binding := bindingMap[key]
			binding.Members = m
			bindingMap[key] = binding
		} else {
			// Add new cloudresourcemanager.Binding to the map
			bindingMap[key] = &cloudresourcemanager.Binding{
				Role:      binding["role"].(string),
				Members:   members,
				Condition: condition,
			}
		}
	}

	// All binding{} blocks, post conversion to cloudresourcemanager.Binding and combining of member lists, are stored in an array
	bindings = []*cloudresourcemanager.Binding{}
	for _, v := range bindingMap {
		v := v
		bindings = append(bindings, v)
	}
	policy.Bindings = bindings

	// Sort bindings within the list to get simpler diffs, as it's what the API does
	// Sorting is based on the binding's role + condition fields
	sort.Slice(policy.Bindings, iamPolicyBindingsLessFunction(policy))

	// Sort members within each binding in the list to get simpler diffs, as it's what the API does
	for i := 0; i < len(policy.Bindings); i++ {
		sort.Strings(policy.Bindings[i].Members)
	}

	// Convert each audit_config into a cloudresourcemanager.AuditConfig
	policy.AuditConfigs = expandAuditConfig(aset)

	// Marshal cloudresourcemanager.Policy to JSON suitable for storing in state
	pjson, err := json.Marshal(&policy)
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	pstring := string(pjson)

	if err := d.Set("policy_data", pstring); err != nil {
		return fmt.Errorf("Error setting policy_data: %s", err)
	}
	d.SetId(strconv.Itoa(tpgresource.Hashcode(pstring)))

	return nil
}

func expandAuditConfig(set *schema.Set) []*cloudresourcemanager.AuditConfig {
	auditConfigs := make([]*cloudresourcemanager.AuditConfig, 0, set.Len())
	for _, v := range set.List() {
		config := v.(map[string]interface{})
		// build list of audit configs first
		auditLogConfigSet := config["audit_log_configs"].(*schema.Set)
		// the array we're going to add to the outgoing resource
		auditLogConfigs := make([]*cloudresourcemanager.AuditLogConfig, 0, auditLogConfigSet.Len())
		for _, y := range auditLogConfigSet.List() {
			logConfig := y.(map[string]interface{})
			auditLogConfigs = append(auditLogConfigs, &cloudresourcemanager.AuditLogConfig{
				LogType:         logConfig["log_type"].(string),
				ExemptedMembers: tpgresource.ConvertStringArr(logConfig["exempted_members"].(*schema.Set).List()),
			})
		}
		auditConfigs = append(auditConfigs, &cloudresourcemanager.AuditConfig{
			Service:         config["service"].(string),
			AuditLogConfigs: auditLogConfigs,
		})
	}
	return auditConfigs
}

func iamPolicyBindingsLessFunction(policy cloudresourcemanager.Policy) func(i, j int) bool {

	return func(i, j int) bool {
		// Sort bindings by role, if they're not the same
		sameRole := policy.Bindings[i].Role == policy.Bindings[j].Role
		if !sameRole {
			return policy.Bindings[i].Role < policy.Bindings[j].Role
		}

		iConditionOk := policy.Bindings[i].Condition != nil
		jConditionOk := policy.Bindings[j].Condition != nil

		// If both bindings lack conditions we cannot sort them further
		if !iConditionOk && !jConditionOk {
			return false
		}

		// Sort by presence of a condition on only one of the two bindings
		if !iConditionOk && jConditionOk {
			return true
		}
		if iConditionOk && !jConditionOk {
			return false
		}

		// At this point both bindings have conditions

		sameExpression := policy.Bindings[i].Condition.Expression == policy.Bindings[j].Condition.Expression
		sameTitle := policy.Bindings[i].Condition.Title == policy.Bindings[j].Condition.Title
		sameDescription := policy.Bindings[i].Condition.Description == policy.Bindings[j].Condition.Description

		// Don't sort if conditions are the same
		if sameExpression && sameTitle && sameDescription {
			return false
		}

		// Sort by both bindings' conditions' expressions, if they're not equivalent
		if !sameExpression {
			return policy.Bindings[i].Condition.Expression < policy.Bindings[j].Condition.Expression
		}

		// Sort by both bindings' conditions' titles, if they're not equivalent
		if !sameTitle {
			return policy.Bindings[i].Condition.Title < policy.Bindings[j].Condition.Title
		}

		// Comparing conditions' descriptions is the last available way to sort
		return policy.Bindings[i].Condition.Description < policy.Bindings[j].Condition.Description
	}
}
