package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func ResourceIamPolicy(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc, options ...func(*IamSettings)) *schema.Resource {
	settings := &IamSettings{}
	for _, o := range options {
		o(settings)
	}

	return &schema.Resource{
		Create: ResourceIamPolicyCreate(newUpdaterFunc),
		Read:   ResourceIamPolicyRead(newUpdaterFunc),
		Update: ResourceIamPolicyUpdate(newUpdaterFunc),
		Delete: ResourceIamPolicyDelete(newUpdaterFunc),

		// if non-empty, this will be used to send a deprecation message when the
		// resource is used.
		DeprecationMessage: settings.DeprecationMessage,

		Schema: mergeSchemas(IamPolicyBaseSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamPolicyImport(resourceIdParser),
		},
		UseJSONNumber: true,
	}
}

func ResourceIamPolicyCreate(newUpdaterFunc newResourceIamUpdaterFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		if err = setIamPolicyData(d, updater); err != nil {
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

		if err := d.Set("etag", policy.Etag); err != nil {
			return fmt.Errorf("Error setting etag: %s", err)
		}
		if err := d.Set("policy_data", marshalIamPolicy(policy)); err != nil {
			return fmt.Errorf("Error setting policy_data: %s", err)
		}

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
		pol := &cloudresourcemanager.Policy{}
		if v, ok := d.GetOk("etag"); ok {
			pol.Etag = v.(string)
		}
		pol.Version = iamPolicyVersion
		err = updater.SetResourceIamPolicy(pol)
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
	policy.Version = iamPolicyVersion

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
	if policy, err := unmarshalIamPolicy(i.(string)); err != nil {
		es = append(es, err)
	} else {
		for i, binding := range policy.Bindings {
			for j, member := range binding.Members {
				_, memberErrors := validateIAMMember(member, fmt.Sprintf("bindings.%d.members.%d", i, j))
				es = append(es, memberErrors...)
			}
		}
	}
	return
}
