package google

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	storagev1 "google.golang.org/api/storage/v1"
)

func resourceStorageBucketIAMPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageBucketIAMPolicyCreate,
		Read:   resourceStorageBucketIAMPolicyRead,
		Update: resourceStorageBucketIAMPolicyUpdate,
		Delete: resourceStorageBucketIAMPolicyDelete,

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
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

func setBucketIamPolicy(d *schema.ResourceData, config *Config) error {
	bucket := d.Get("bucket").(string)
	policy, err := unmarshalStorageIamPolicy(d.Get("policy_data").(string))
	if err != nil {
		return fmt.Errorf("'policy_data' is not valid for %s: %s", bucket, err)
	}

	_, err = config.clientStorage.Buckets.SetIamPolicy(bucket, policy).Do()
	return err
}

func resourceStorageBucketIAMPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := setBucketIamPolicy(d, config); err != nil {
		return err
	}

	d.SetId(d.Get("bucket").(string))

	return resourceStorageBucketIAMPolicyRead(d, meta)
}

func resourceStorageBucketIAMPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bucket := d.Get("bucket").(string)

	policy, err := config.clientStorage.Buckets.GetIamPolicy(bucket).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Iam policy for %s", bucket))
	}

	d.Set("etag", policy.Etag)
	d.Set("policy_data", marshalStorageIamPolicy(policy))

	return nil
}

func resourceStorageBucketIAMPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("policy_data") {
		if err := setBucketIamPolicy(d, config); err != nil {
			return err
		}
	}
	return resourceStorageBucketIAMPolicyRead(d, meta)
}

func resourceStorageBucketIAMPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bucket := d.Get("bucket").(string)

	_, err := config.clientStorage.Buckets.SetIamPolicy(bucket, &storagev1.Policy{}).Do()

	if err != nil {
		return err
	}

	return nil
}

func marshalStorageIamPolicy(policy *storagev1.Policy) string {
	pdBytes, _ := json.Marshal(&storagev1.Policy{
		Bindings: policy.Bindings,
	})
	return string(pdBytes)
}

func unmarshalStorageIamPolicy(policyData string) (*storagev1.Policy, error) {
	policy := &storagev1.Policy{}
	if err := json.Unmarshal([]byte(policyData), policy); err != nil {
		return nil, fmt.Errorf("Could not unmarshal policy data %s:\n%s", policyData, err)
	}
	return policy, nil
}

func validateStorageIamPolicy(i interface{}, k string) (s []string, es []error) {
	_, err := unmarshalV2IamPolicy(i.(string))
	if err != nil {
		es = append(es, err)
	}
	return
}
