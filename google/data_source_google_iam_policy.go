package google

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var iamBinding *schema.Schema = &schema.Schema{
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	},
}

var auditConfig *schema.Schema = &schema.Schema{
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
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
		},
	},
}

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
			"binding": iamBinding,
			"policy_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audit_config": auditConfig,
		},
	}
}

// dataSourceGoogleIamPolicyRead reads a data source from config and writes it
// to state.
func dataSourceGoogleIamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	var policy cloudresourcemanager.Policy
	var bindings []*cloudresourcemanager.Binding
	var auditConfigs []*cloudresourcemanager.AuditConfig

	// The schema supports multiple binding{} blocks
	bset := d.Get("binding").(*schema.Set)
	aset := d.Get("audit_config").(*schema.Set)

	// All binding{} blocks will be converted and stored in an array
	bindings = make([]*cloudresourcemanager.Binding, bset.Len())
	auditConfigs = make([]*cloudresourcemanager.AuditConfig, aset.Len())
	policy.Bindings = bindings
	policy.AuditConfigs = auditConfigs

	// Convert each config binding into a cloudresourcemanager.Binding
	for i, v := range bset.List() {
		binding := v.(map[string]interface{})
		policy.Bindings[i] = &cloudresourcemanager.Binding{
			Role:    binding["role"].(string),
			Members: convertStringSet(binding["members"].(*schema.Set)),
		}
	}

	// Convert each audit_config into a cloudresourcemanager.AuditConfig
	for i, v := range aset.List() {
		config := v.(map[string]interface{})

		// build list of audit configs first
		auditLogConfigSet := config["audit_log_configs"].(*schema.Set)
		// the array we're going to add to the outgoing resource
		auditLogConfigs := make([]*cloudresourcemanager.AuditLogConfig, auditLogConfigSet.Len())
		for x, y := range auditLogConfigSet.List() {
			logConfig := y.(map[string]interface{})
			auditLogConfigs[x] = &cloudresourcemanager.AuditLogConfig{
				LogType:         logConfig["log_type"].(string),
				ExemptedMembers: convertStringArr(logConfig["exempted_members"].([]interface{})),
			}
		}

		policy.AuditConfigs[i] = &cloudresourcemanager.AuditConfig{
			Service:         config["service"].(string),
			AuditLogConfigs: auditLogConfigs,
		}
	}

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
