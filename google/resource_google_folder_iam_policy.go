package google

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func resourceGoogleFolderIamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleFolderIamPolicyCreate,
		Read:   resourceGoogleFolderIamPolicyRead,
		Update: resourceGoogleFolderIamPolicyUpdate,
		Delete: resourceGoogleFolderIamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGoogleFolderIamPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"folder": {
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
				Removed:  "The authoritative field was removed. To ignore changes not managed by Terraform, use google_folder_iam_binding and google_folder_iam_member instead. See https://www.terraform.io/docs/providers/google/r/google_folder_iam.html for more information.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"restore_policy": {
				Removed:  "This field was removed alongside the authoritative field. To ignore changes not managed by Terraform, use google_folder_iam_binding and google_folder_iam_member instead. See https://www.terraform.io/docs/providers/google/r/google_folder_iam.html for more information.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"disable_folder": {
				Removed:  "This field was removed alongside the authoritative field. Use lifecycle.prevent_destroy instead.",
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceGoogleFolderIamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	folder := d.Get("folder").(string)

	mutexKey := getFolderIamPolicyMutexKey(folder)
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	// Get the policy in the template
	policy, err := getFolderResourceIamPolicy(d)
	if err != nil {
		return fmt.Errorf("Could not get valid 'policy_data' from resource: %v", err)
	}

	log.Printf("[DEBUG] Setting IAM policy for folder %q", folder)
	err = setFolderIamPolicy(policy, config, folder)
	if err != nil {
		return err
	}

	d.SetId(folder)
	return resourceGoogleFolderIamPolicyRead(d, meta)
}

func resourceGoogleFolderIamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	folder := d.Get("folder").(string)

	policy, err := getFolderIamPolicy(folder, config)
	if err != nil {
		return err
	}

	policyBytes, err := json.Marshal(&resourceManagerV2Beta1.Policy{Bindings: policy.Bindings, AuditConfigs: policy.AuditConfigs})
	if err != nil {
		return fmt.Errorf("Error marshaling IAM policy: %v", err)
	}

	d.Set("etag", policy.Etag)
	d.Set("policy_data", string(policyBytes))
	d.Set("folder", folder)
	return nil
}

func resourceGoogleFolderIamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	folder := d.Get("folder").(string)

	mutexKey := getFolderIamPolicyMutexKey(folder)
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	// Get the policy in the template
	policy, err := getFolderResourceIamPolicy(d)
	if err != nil {
		return fmt.Errorf("Could not get valid 'policy_data' from resource: %v", err)
	}

	log.Printf("[DEBUG] Updating IAM policy for folder %q", folder)
	err = setFolderIamPolicy(policy, config, folder)
	if err != nil {
		return fmt.Errorf("Error setting folder IAM policy: %v", err)
	}

	return resourceGoogleFolderIamPolicyRead(d, meta)
}

func resourceGoogleFolderIamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Deleting google_folder_iam_policy")
	config := meta.(*Config)
	folder := d.Get("folder").(string)

	mutexKey := getFolderIamPolicyMutexKey(folder)
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	// Get the existing IAM policy from the API so we can repurpose the etag and audit config
	ep, err := getFolderIamPolicy(folder, config)
	if err != nil {
		return fmt.Errorf("Error retrieving IAM policy from folder API: %v", err)
	}

	ep.Bindings = make([]*resourceManagerV2Beta1.Binding, 0)
	if err = setFolderIamPolicy(ep, config, folder); err != nil {
		return fmt.Errorf("Error applying IAM policy to folder: %v", err)
	}

	d.SetId("")
	return nil
}

func resourceGoogleFolderIamPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("folder", d.Id())
	return []*schema.ResourceData{d}, nil
}

func setFolderIamPolicy(policy *resourceManagerV2Beta1.Policy, config *Config, pid string) error {
	// Apply the policy
	pbytes, _ := json.Marshal(policy)
	log.Printf("[DEBUG] Setting policy %#v for folder: %s", string(pbytes), pid)
	_, err := config.clientResourceManagerV2Beta1.Folders.SetIamPolicy(pid,
		&resourceManagerV2Beta1.SetIamPolicyRequest{Policy: policy, UpdateMask: "bindings,etag,auditConfigs"}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error applying IAM policy for folder %q. Policy is %#v, error is {{err}}", pid, policy), err)
	}
	return nil
}

// Get a cloudresourcemanager.Policy from a schema.ResourceData
func getFolderResourceIamPolicy(d *schema.ResourceData) (*resourceManagerV2Beta1.Policy, error) {
	ps := d.Get("policy_data").(string)
	// The policy string is just a marshaled cloudresourcemanager.Policy.
	policy := &resourceManagerV2Beta1.Policy{}
	if err := json.Unmarshal([]byte(ps), policy); err != nil {
		return nil, fmt.Errorf("Could not unmarshal %s:\n: %v", ps, err)
	}
	return policy, nil
}

// Retrieve the existing IAM Policy for a Project
func getFolderIamPolicy(folder string, config *Config) (*resourceManagerV2Beta1.Policy, error) {
	p, err := config.clientResourceManagerV2Beta1.Folders.GetIamPolicy(folder,
		&resourceManagerV2Beta1.GetIamPolicyRequest{}).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for folder %q: %s", folder, err)
	}
	return p, nil
}

func getFolderIamPolicyMutexKey(pid string) string {
	return fmt.Sprintf("iam-folder-%s", pid)
}
