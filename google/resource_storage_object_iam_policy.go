package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	storagev1 "google.golang.org/api/storage/v1"
)

func resourceStorageObjectIAMPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageObjectIAMPolicyCreate,
		Read:   resourceStorageObjectIAMPolicyRead,
		Update: resourceStorageObjectIAMPolicyUpdate,
		Delete: resourceStorageObjectIAMPolicyDelete,

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object": &schema.Schema{
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

func setObjectIamPolicy(d *schema.ResourceData, config *Config) error {
	bucket := d.Get("bucket").(string)
	object := d.Get("object").(string)
	policy, err := unmarshalStorageIamPolicy(d.Get("policy_data").(string))
	if err != nil {
		return fmt.Errorf("'policy_data' is not valid for %s: %s", object, err)
	}

	log.Printf("[DEBUG] Setting IAM policy for object %q on bucket %q", object, bucket)
	_, err = config.clientStorage.Objects.SetIamPolicy(bucket, object, policy).Do()
	log.Printf("[DEBUG] Set IAM policy for object %q on %q", object, bucket)
	return err
}

func resourceStorageObjectIAMPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := setObjectIamPolicy(d, config); err != nil {
		return err
	}

	d.SetId(d.Get("bucket").(string) + "-" + d.Get("object").(string) + "-iam-policy")

	return resourceStorageObjectIAMPolicyRead(d, meta)
}

func resourceStorageObjectIAMPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bucket := d.Get("bucket").(string)
	object := d.Get("object").(string)

	log.Printf("[DEBUG] Reading IAM policy for object %q on bucket %q", object, bucket)
	policy, err := config.clientStorage.Objects.GetIamPolicy(bucket, object).Do()
	log.Printf("[DEBUG] Reat IAM policy for object %q on bucket %q: %v", object, bucket, policy)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Iam policy for %s", object))
	}

	d.Set("etag", policy.Etag)
	policyData, err := marshalStorageIamPolicy(policy)
	if err != nil {
		return fmt.Errorf("Fail to marshal policy data")
	}
	d.Set("policy_data", policyData)

	return nil
}

func resourceStorageObjectIAMPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("policy_data") {
		if err := setObjectIamPolicy(d, config); err != nil {
			return err
		}
	}
	return resourceStorageObjectIAMPolicyRead(d, meta)
}

func resourceStorageObjectIAMPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bucket := d.Get("bucket").(string)
	object := d.Get("object").(string)

	_, err := config.clientStorage.Objects.SetIamPolicy(bucket, object, &storagev1.Policy{}).Do()

	if err != nil {
		return err
	}

	return nil
}
