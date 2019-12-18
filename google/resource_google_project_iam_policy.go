package google

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareProjectName,
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
		},
	}
}

func compareProjectName(_, old, new string, _ *schema.ResourceData) bool {
	// We can either get "projects/project-id" or "project-id", so strip any prefixes
	return GetResourceNameFromSelfLink(old) == GetResourceNameFromSelfLink(new)
}

func resourceGoogleProjectIamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project := GetResourceNameFromSelfLink(d.Get("project").(string))

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
	project := GetResourceNameFromSelfLink(d.Get("project").(string))

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
	project := GetResourceNameFromSelfLink(d.Get("project").(string))

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
	project := GetResourceNameFromSelfLink(d.Get("project").(string))

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

func getProjectIamPolicyMutexKey(pid string) string {
	return fmt.Sprintf("iam-project-%s", pid)
}
