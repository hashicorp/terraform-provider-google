package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamPolicyBaseSchema = map[string]*schema.Schema{
	"policy_data": {
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: jsonPolicyDiffSuppress,
		ValidateFunc:     validateIamPolicy,
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func iamPolicyImport(resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}
		return []*schema.ResourceData{d}, nil
	}
}

func ResourceIamPolicy(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc) *schema.Resource {
	return &schema.Resource{
		Create: ResourceIamPolicyCreate(newUpdaterFunc),
		Read:   ResourceIamPolicyRead(newUpdaterFunc),
		Update: ResourceIamPolicyUpdate(newUpdaterFunc),
		Delete: ResourceIamPolicyDelete(newUpdaterFunc),

		Schema: mergeSchemas(IamPolicyBaseSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamPolicyImport(resourceIdParser),
		},
	}
}

func ResourceIamPolicyCreate(newUpdaterFunc newResourceIamUpdaterFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		if err := setIamPolicyData(d, updater); err != nil {
			return err
		}

		d.SetId(updater.GetResourceId())
		return ResourceIamPolicyRead(newUpdaterFunc)(d, meta)
	}
}

func ResourceIamPolicyRead(newUpdaterFunc newResourceIamUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		policy, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Resource %q with IAM Policy", updater.DescribeResource()))
		}

		d.Set("etag", policy.Etag)
		d.Set("policy_data", marshalIamPolicy(policy))

		return nil
	}
}

func ResourceIamPolicyUpdate(newUpdaterFunc newResourceIamUpdaterFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		if d.HasChange("policy_data") {
			if err := setIamPolicyData(d, updater); err != nil {
				return err
			}
		}

		return ResourceIamPolicyRead(newUpdaterFunc)(d, meta)
	}
}

func ResourceIamPolicyDelete(newUpdaterFunc newResourceIamUpdaterFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		// Set an empty policy to delete the attached policy.
		err = updater.SetResourceIamPolicy(&cloudresourcemanager.Policy{})
		if err != nil {
			return err
		}

		return nil
	}
}

func setIamPolicyData(d *schema.ResourceData, updater ResourceIamUpdater) error {
	policy, err := unmarshalIamPolicy(d.Get("policy_data").(string))
	if err != nil {
		return fmt.Errorf("'policy_data' is not valid for %s: %s", updater.DescribeResource(), err)
	}

	err = updater.SetResourceIamPolicy(policy)
	if err != nil {
		return err
	}

	return nil
}

func marshalIamPolicy(policy *cloudresourcemanager.Policy) string {
	pdBytes, _ := json.Marshal(&cloudresourcemanager.Policy{
		AuditConfigs: policy.AuditConfigs,
		Bindings:     policy.Bindings,
	})
	return string(pdBytes)
}

func unmarshalIamPolicy(policyData string) (*cloudresourcemanager.Policy, error) {
	policy := &cloudresourcemanager.Policy{}
	if err := json.Unmarshal([]byte(policyData), policy); err != nil {
		return nil, fmt.Errorf("Could not unmarshal policy data %s:\n%s", policyData, err)
	}
	return policy, nil
}

func validateIamPolicy(i interface{}, k string) (s []string, es []error) {
	_, err := unmarshalIamPolicy(i.(string))
	if err != nil {
		es = append(es, err)
	}
	return
}
