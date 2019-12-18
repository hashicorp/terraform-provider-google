package google

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var iamBindingSchema = map[string]*schema.Schema{
	"role": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"members": {
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			DiffSuppressFunc: caseDiffSuppress,
			ValidateFunc:     validation.StringDoesNotMatch(regexp.MustCompile("^deleted:"), "Terraform does not support IAM bindings for deleted principals"),
		},
		Set: func(v interface{}) int {
			return schema.HashString(strings.ToLower(v.(string)))
		},
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ResourceIamBinding(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc) *schema.Resource {
	return ResourceIamBindingWithBatching(parentSpecificSchema, newUpdaterFunc, resourceIdParser, IamBatchingDisabled)
}

// Resource that batches requests to the same IAM policy across multiple IAM fine-grained resources
func ResourceIamBindingWithBatching(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc, enableBatching bool) *schema.Resource {
	return &schema.Resource{
		Create: resourceIamBindingCreateUpdate(newUpdaterFunc, enableBatching),
		Read:   resourceIamBindingRead(newUpdaterFunc),
		Update: resourceIamBindingCreateUpdate(newUpdaterFunc, enableBatching),
		Delete: resourceIamBindingDelete(newUpdaterFunc, enableBatching),
		Schema: mergeSchemas(iamBindingSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamBindingImport(newUpdaterFunc, resourceIdParser),
		},
	}
}

func resourceIamBindingCreateUpdate(newUpdaterFunc newResourceIamUpdaterFunc, enableBatching bool) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		binding := getResourceIamBinding(d)
		modifyF := func(ep *cloudresourcemanager.Policy) error {
			cleaned := filterBindingsWithRoleAndCondition(ep.Bindings, binding.Role, binding.Condition)
			ep.Bindings = append(cleaned, binding)
			return nil
		}

		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config, fmt.Sprintf(
				"Set IAM Binding for role %q on %q", binding.Role, updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return err
		}

		d.SetId(updater.GetResourceId() + "/" + binding.Role)
		return resourceIamBindingRead(newUpdaterFunc)(d, meta)
	}
}

func resourceIamBindingRead(newUpdaterFunc newResourceIamUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		eBinding := getResourceIamBinding(d)
		eCondition := conditionKeyFromCondition(eBinding.Condition)
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Resource %q with IAM Binding (Role %q)", updater.DescribeResource(), eBinding.Role))
		}
		log.Printf("[DEBUG] Retrieved policy for %s: %+v", updater.DescribeResource(), p)
		log.Printf("[DEBUG] Looking for binding with role %q and condition %+v", eBinding.Role, eCondition)

		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role == eBinding.Role && conditionKeyFromCondition(b.Condition) == eCondition {
				binding = b
				break
			}
		}

		if binding == nil {
			log.Printf("[DEBUG] Binding for role %q and condition %+v not found in policy for %s, assuming it has no members.", eBinding.Role, eCondition, updater.DescribeResource())
			d.Set("role", eBinding.Role)
			d.Set("members", nil)
			return nil
		} else {
			d.Set("role", binding.Role)
			d.Set("members", binding.Members)
		}
		d.Set("etag", p.Etag)
		return nil
	}
}

func iamBindingImport(newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		s := strings.Fields(d.Id())
		var id, role string
		if len(s) != 2 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Binding id %s; expected 'resource_name role'.", s)
		}
		id, role = s[0], s[1]

		// Set the ID only to the first part so all IAM types can share the same resourceIdParserFunc.
		d.SetId(id)
		d.Set("role", role)
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}

		// Set the ID again so that the ID matches the ID it would have if it had been created via TF.
		// Use the current ID in case it changed in the resourceIdParserFunc.
		d.SetId(d.Id() + "/" + role)

		// It is possible to return multiple bindings, since we can learn about all the bindings
		// for this resource here.  Unfortunately, `terraform import` has some messy behavior here -
		// there's no way to know at this point which resource is being imported, so it's not possible
		// to order this list in a useful way.  In the event of a complex set of bindings, the user
		// will have a terribly confusing set of imported resources and no way to know what matches
		// up to what.  And since the only users who will do a terraform import on their IAM bindings
		// are users who aren't too familiar with Google Cloud IAM (because a "create" for bindings or
		// members is idempotent), it's reasonable to expect that the user will be very alarmed by the
		// plan that terraform will output which mentions destroying a dozen-plus IAM bindings.  With
		// that in mind, we return only the binding that matters.
		return []*schema.ResourceData{d}, nil
	}
}

func resourceIamBindingDelete(newUpdaterFunc newResourceIamUpdaterFunc, enableBatching bool) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		binding := getResourceIamBinding(d)
		modifyF := func(p *cloudresourcemanager.Policy) error {
			p.Bindings = filterBindingsWithRoleAndCondition(p.Bindings, binding.Role, binding.Condition)
			return nil
		}

		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config, fmt.Sprintf(
				"Delete IAM Binding for role %q on %q", binding.Role, updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Resource %q for IAM binding with role %q", updater.DescribeResource(), binding.Role))
		}

		return resourceIamBindingRead(newUpdaterFunc)(d, meta)
	}
}

func getResourceIamBinding(d *schema.ResourceData) *cloudresourcemanager.Binding {
	members := d.Get("members").(*schema.Set).List()
	b := &cloudresourcemanager.Binding{
		Members: convertStringArr(members),
		Role:    d.Get("role").(string),
	}
	return b
}
