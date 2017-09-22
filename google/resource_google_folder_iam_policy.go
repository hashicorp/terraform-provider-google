package google

import (
	"github.com/hashicorp/terraform/helper/schema"

	"encoding/json"
	"fmt"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func resourceGoogleFolderIamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleFolderIamPolicyCreate,
		Read:   resourceGoogleFolderIamPolicyRead,
		Update: resourceGoogleFolderIamPolicyUpdate,
		Delete: resourceGoogleFolderIamPolicyDelete,

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
				ValidateFunc:     validateV2IamPolicy,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGoogleFolderIamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := setFolderIamPolicy(d, config); err != nil {
		return err
	}

	d.SetId(d.Get("folder").(string))

	return resourceGoogleFolderIamPolicyRead(d, meta)
}

func resourceGoogleFolderIamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	folder := d.Get("folder").(string)

	policy, err := config.clientResourceManagerV2Beta1.Folders.GetIamPolicy(folder, &resourceManagerV2Beta1.GetIamPolicyRequest{}).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Iam policy for %s", folder))
	}

	d.Set("etag", policy.Etag)
	d.Set("policy_data", marshalV2IamPolicy(policy))

	return nil
}

func resourceGoogleFolderIamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("policy_data") {
		if err := setFolderIamPolicy(d, config); err != nil {
			return err
		}
	}

	return resourceGoogleFolderIamPolicyRead(d, meta)
}

func resourceGoogleFolderIamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	folder := d.Get("folder").(string)

	_, err := config.clientResourceManagerV2Beta1.Folders.SetIamPolicy(folder, &resourceManagerV2Beta1.SetIamPolicyRequest{
		Policy:     &resourceManagerV2Beta1.Policy{},
		UpdateMask: "bindings",
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func setFolderIamPolicy(d *schema.ResourceData, config *Config) error {
	folder := d.Get("folder").(string)
	policy, err := unmarshalV2IamPolicy(d.Get("policy_data").(string))
	if err != nil {
		return fmt.Errorf("'policy_data' is not valid for %s: %s", folder, err)
	}

	_, err = config.clientResourceManagerV2Beta1.Folders.SetIamPolicy(folder, &resourceManagerV2Beta1.SetIamPolicyRequest{
		Policy:     policy,
		UpdateMask: "bindings",
	}).Do()

	return err
}

func marshalV2IamPolicy(policy *resourceManagerV2Beta1.Policy) string {
	pdBytes, _ := json.Marshal(&resourceManagerV2Beta1.Policy{
		Bindings: policy.Bindings,
	})
	return string(pdBytes)
}

func unmarshalV2IamPolicy(policyData string) (*resourceManagerV2Beta1.Policy, error) {
	policy := &resourceManagerV2Beta1.Policy{}
	if err := json.Unmarshal([]byte(policyData), policy); err != nil {
		return nil, fmt.Errorf("Could not unmarshal policy data %s:\n%s", policyData, err)
	}
	return policy, nil
}

func validateV2IamPolicy(i interface{}, k string) (s []string, es []error) {
	_, err := unmarshalV2IamPolicy(i.(string))
	if err != nil {
		es = append(es, err)
	}
	return
}
