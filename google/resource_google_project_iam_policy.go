package google

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func resourceGoogleProjectIamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamPolicyCreate,
		Read:   resourceGoogleProjectIamPolicyRead,
		Update: resourceGoogleProjectIamPolicyUpdate,
		Delete: resourceGoogleProjectIamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGoogleProjectIamPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_data": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: jsonPolicyDiffSuppress,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authoritative": {
				Removed:  "The authoritative field was removed. To ignore changes not managed by Terraform, use google_project_iam_binding and google_project_iam_member instead. See https://www.terraform.io/docs/providers/google/r/google_project_iam.html for more information.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"restore_policy": {
				Removed:  "This field was removed alongside the authoritative field. To ignore changes not managed by Terraform, use google_project_iam_binding and google_project_iam_member instead. See https://www.terraform.io/docs/providers/google/r/google_project_iam.html for more information.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"disable_project": {
				Removed:  "This field was removed alongside the authoritative field. Use lifecycle.prevent_destroy instead.",
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceGoogleProjectIamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project := d.Get("project").(string)

	mutexKey := getProjectIamPolicyMutexKey(project)
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	// Get the policy in the template
	policy, err := getResourceIamPolicy(d)
	if err != nil {
		return fmt.Errorf("Could not get valid 'policy_data' from resource: %v", err)
	}

	log.Printf("[DEBUG] Setting IAM policy for project %q", project)
	err = setProjectIamPolicy(policy, config, project)
	if err != nil {
		return err
	}

	d.SetId(project)
	return resourceGoogleProjectIamPolicyRead(d, meta)
}

func resourceGoogleProjectIamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project := d.Get("project").(string)

	policy, err := getProjectIamPolicy(project, config)
	if err != nil {
		return err
	}

	policyBytes, err := json.Marshal(&cloudresourcemanager.Policy{Bindings: policy.Bindings, AuditConfigs: policy.AuditConfigs})
	if err != nil {
		return fmt.Errorf("Error marshaling IAM policy: %v", err)
	}

	d.Set("etag", policy.Etag)
	d.Set("policy_data", string(policyBytes))
	d.Set("project", project)
	return nil
}

func resourceGoogleProjectIamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project := d.Get("project").(string)

	mutexKey := getProjectIamPolicyMutexKey(project)
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	// Get the policy in the template
	policy, err := getResourceIamPolicy(d)
	if err != nil {
		return fmt.Errorf("Could not get valid 'policy_data' from resource: %v", err)
	}

	log.Printf("[DEBUG] Updating IAM policy for project %q", project)
	err = setProjectIamPolicy(policy, config, project)
	if err != nil {
		return fmt.Errorf("Error setting project IAM policy: %v", err)
	}

	return resourceGoogleProjectIamPolicyRead(d, meta)
}

func resourceGoogleProjectIamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Deleting google_project_iam_policy")
	config := meta.(*Config)
	project := d.Get("project").(string)

	mutexKey := getProjectIamPolicyMutexKey(project)
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	// Get the existing IAM policy from the API so we can repurpose the etag and audit config
	ep, err := getProjectIamPolicy(project, config)
	if err != nil {
		return fmt.Errorf("Error retrieving IAM policy from project API: %v", err)
	}

	ep.Bindings = make([]*cloudresourcemanager.Binding, 0)
	if err = setProjectIamPolicy(ep, config, project); err != nil {
		return fmt.Errorf("Error applying IAM policy to project: %v", err)
	}

	d.SetId("")
	return nil
}

func resourceGoogleProjectIamPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("project", d.Id())
	return []*schema.ResourceData{d}, nil
}

func setProjectIamPolicy(policy *cloudresourcemanager.Policy, config *Config, pid string) error {
	// Apply the policy
	pbytes, _ := json.Marshal(policy)
	log.Printf("[DEBUG] Setting policy %#v for project: %s", string(pbytes), pid)
	_, err := config.clientResourceManager.Projects.SetIamPolicy(pid,
		&cloudresourcemanager.SetIamPolicyRequest{Policy: policy, UpdateMask: "bindings,etag,auditConfigs"}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error applying IAM policy for project %q. Policy is %#v, error is {{err}}", pid, policy), err)
	}
	return nil
}

// Get a cloudresourcemanager.Policy from a schema.ResourceData
func getResourceIamPolicy(d *schema.ResourceData) (*cloudresourcemanager.Policy, error) {
	ps := d.Get("policy_data").(string)
	// The policy string is just a marshaled cloudresourcemanager.Policy.
	policy := &cloudresourcemanager.Policy{}
	if err := json.Unmarshal([]byte(ps), policy); err != nil {
		return nil, fmt.Errorf("Could not unmarshal %s:\n: %v", ps, err)
	}
	return policy, nil
}

// Retrieve the existing IAM Policy for a Project
func getProjectIamPolicy(project string, config *Config) (*cloudresourcemanager.Policy, error) {
	p, err := config.clientResourceManager.Projects.GetIamPolicy(project,
		&cloudresourcemanager.GetIamPolicyRequest{}).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for project %q: %s", project, err)
	}
	return p, nil
}

func jsonPolicyDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	var oldPolicy, newPolicy cloudresourcemanager.Policy
	if err := json.Unmarshal([]byte(old), &oldPolicy); err != nil {
		log.Printf("[ERROR] Could not unmarshal old policy %s: %v", old, err)
		return false
	}
	if err := json.Unmarshal([]byte(new), &newPolicy); err != nil {
		log.Printf("[ERROR] Could not unmarshal new policy %s: %v", new, err)
		return false
	}
	if newPolicy.Etag != oldPolicy.Etag {
		return false
	}
	if newPolicy.Version != oldPolicy.Version {
		return false
	}
	if !compareBindings(oldPolicy.Bindings, newPolicy.Bindings) {
		return false
	}
	if !compareAuditConfigs(oldPolicy.AuditConfigs, newPolicy.AuditConfigs) {
		return false
	}
	return true
}

func derefBindings(b []*cloudresourcemanager.Binding) []cloudresourcemanager.Binding {
	db := make([]cloudresourcemanager.Binding, len(b))

	for i, v := range b {
		db[i] = *v
		sort.Strings(db[i].Members)
	}
	return db
}

func compareBindings(a, b []*cloudresourcemanager.Binding) bool {
	a = mergeBindings(a)
	b = mergeBindings(b)
	sort.Sort(sortableBindings(a))
	sort.Sort(sortableBindings(b))
	return reflect.DeepEqual(derefBindings(a), derefBindings(b))
}

func compareAuditConfigs(a, b []*cloudresourcemanager.AuditConfig) bool {
	a = mergeAuditConfigs(a)
	b = mergeAuditConfigs(b)
	sort.Sort(sortableAuditConfigs(a))
	sort.Sort(sortableAuditConfigs(b))
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if len(v.AuditLogConfigs) != len(b[i].AuditLogConfigs) {
			return false
		}
		sort.Sort(sortableAuditLogConfigs(v.AuditLogConfigs))
		sort.Sort(sortableAuditLogConfigs(b[i].AuditLogConfigs))
		for x, logConfig := range v.AuditLogConfigs {
			if b[i].AuditLogConfigs[x].LogType != logConfig.LogType {
				return false
			}
			sort.Strings(logConfig.ExemptedMembers)
			sort.Strings(b[i].AuditLogConfigs[x].ExemptedMembers)
			if len(logConfig.ExemptedMembers) != len(b[i].AuditLogConfigs[x].ExemptedMembers) {
				return false
			}
			for pos, exemptedMember := range logConfig.ExemptedMembers {
				if b[i].AuditLogConfigs[x].ExemptedMembers[pos] != exemptedMember {
					return false
				}
			}
		}
	}
	return true
}

type sortableBindings []*cloudresourcemanager.Binding

func (b sortableBindings) Len() int {
	return len(b)
}
func (b sortableBindings) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b sortableBindings) Less(i, j int) bool {
	return b[i].Role < b[j].Role
}

type sortableAuditConfigs []*cloudresourcemanager.AuditConfig

func (b sortableAuditConfigs) Len() int {
	return len(b)
}
func (b sortableAuditConfigs) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b sortableAuditConfigs) Less(i, j int) bool {
	return b[i].Service < b[j].Service
}

type sortableAuditLogConfigs []*cloudresourcemanager.AuditLogConfig

func (b sortableAuditLogConfigs) Len() int {
	return len(b)
}
func (b sortableAuditLogConfigs) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b sortableAuditLogConfigs) Less(i, j int) bool {
	return b[i].LogType < b[j].LogType
}

func getProjectIamPolicyMutexKey(pid string) string {
	return fmt.Sprintf("iam-project-%s", pid)
}
