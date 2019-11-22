package google

import (
	"encoding/json"
	"regexp"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// dataSourceGoogleIamPolicy returns a *schema.Resource that allows a customer
// to express a Google Cloud IAM policy in a data resource. This is an example
// of how the schema would be used in a config:
//
// data "google_iam_policy" "admin" {
//   binding {
//     role = "roles/storage.objectViewer"
//     members = [
//       "user:evanbrown@google.com",
//     ]
//   }
// }
func dataSourceGoogleIamPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleIamPolicyRead,
		Schema: map[string]*schema.Schema{
			"binding": {
				Type:     schema.TypeSet,
				Required: true,
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

	// All binding{} blocks will be converted and stored in an array
	bindings = make([]*cloudresourcemanager.Binding, bset.Len())
	policy.Bindings = bindings

	// Convert each config binding into a cloudresourcemanager.Binding
	for i, v := range bset.List() {
		binding := v.(map[string]interface{})
		members := convertStringSet(binding["members"].(*schema.Set))

		// Sort members to get simpler diffs as it's what the API does
		sort.Strings(members)

		policy.Bindings[i] = &cloudresourcemanager.Binding{
			Role:    binding["role"].(string),
			Members: members,
		}
	}

	// Sort bindings by their role name to get simpler diffs as it's what the API does
	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].Role < bindings[j].Role
	})

	// Convert each audit_config into a cloudresourcemanager.AuditConfig
	policy.AuditConfigs = expandAuditConfig(aset)

	// Marshal cloudresourcemanager.Policy to JSON suitable for storing in state
	pjson, err := json.Marshal(&policy)
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	pstring := string(pjson)

	d.Set("policy_data", pstring)
	d.SetId(strconv.Itoa(hashcode.String(pstring)))

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
				ExemptedMembers: convertStringArr(logConfig["exempted_members"].(*schema.Set).List()),
			})
		}
		auditConfigs = append(auditConfigs, &cloudresourcemanager.AuditConfig{
			Service:         config["service"].(string),
			AuditLogConfigs: auditLogConfigs,
		})
	}
	return auditConfigs
}
