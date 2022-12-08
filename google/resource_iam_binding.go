package google

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			ValidateFunc:     validateIAMMember,
		},
		Set: func(v interface{}) int {
			return schema.HashString(strings.ToLower(v.(string)))
		},
	},
	"condition": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expression": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"title": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ResourceIamBinding(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc, options ...func(*IamSettings)) *schema.Resource {
	return ResourceIamBindingWithBatching(parentSpecificSchema, newUpdaterFunc, resourceIdParser, IamBatchingDisabled, options...)
}

// Resource that batches requests to the same IAM policy across multiple IAM fine-grained resources
func ResourceIamBindingWithBatching(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc, enableBatching bool, options ...func(*IamSettings)) *schema.Resource {
	settings := &IamSettings{}
	for _, o := range options {
		o(settings)
	}

	return &schema.Resource{
		Create: resourceIamBindingCreateUpdate(newUpdaterFunc, enableBatching),
		Read:   resourceIamBindingRead(newUpdaterFunc),
		Update: resourceIamBindingCreateUpdate(newUpdaterFunc, enableBatching),
		Delete: resourceIamBindingDelete(newUpdaterFunc, enableBatching),

		// if non-empty, this will be used to send a deprecation message when the
		// resource is used.
		DeprecationMessage: settings.DeprecationMessage,

		Schema: mergeSchemas(iamBindingSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamBindingImport(newUpdaterFunc, resourceIdParser),
		},
		UseJSONNumber: true,
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
			ep.Version = iamPolicyVersion
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
		if k := conditionKeyFromCondition(binding.Condition); !k.Empty() {
			d.SetId(d.Id() + "/" + k.String())
		}
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
		log.Print(spew.Sprintf("[DEBUG] Retrieved policy for %s: %#v", updater.DescribeResource(), p))
		log.Printf("[DEBUG] Looking for binding with role %q and condition %#v", eBinding.Role, eCondition)

		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role == eBinding.Role && conditionKeyFromCondition(b.Condition) == eCondition {
				binding = b
				break
			}
		}

		if binding == nil {
			log.Printf("[WARNING] Binding for role %q not found, assuming it has no members. If you expected existing members bound for this role, make sure your role is correctly formatted.", eBinding.Role)
			log.Printf("[DEBUG] Binding for role %q and condition %#v not found in policy for %s, assuming it has no members.", eBinding.Role, eCondition, updater.DescribeResource())
			if err := d.Set("role", eBinding.Role); err != nil {
				return fmt.Errorf("Error setting role: %s", err)
			}
			if err := d.Set("members", nil); err != nil {
				return fmt.Errorf("Error setting members: %s", err)
			}
			return nil
		} else {
			if err := d.Set("role", binding.Role); err != nil {
				return fmt.Errorf("Error setting role: %s", err)
			}
			if err := d.Set("members", binding.Members); err != nil {
				return fmt.Errorf("Error setting members: %s", err)
			}
			if err := d.Set("condition", flattenIamCondition(binding.Condition)); err != nil {
				return fmt.Errorf("Error setting condition: %s", err)
			}
		}
		if err := d.Set("etag", p.Etag); err != nil {
			return fmt.Errorf("Error setting etag: %s", err)
		}
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
		if len(s) < 2 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Binding id %s; expected 'resource_name role [condition_title]'.", s)
		}

		var conditionTitle string
		if len(s) == 2 {
			id, role = s[0], s[1]
		} else {
			// condition titles can have any characters in them, so re-join the split string
			id, role, conditionTitle = s[0], s[1], strings.Join(s[2:], " ")
		}

		// Set the ID only to the first part so all IAM types can share the same resourceIdParserFunc.
		d.SetId(id)
		if err := d.Set("role", role); err != nil {
			return nil, fmt.Errorf("Error setting role: %s", err)
		}
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}

		// Set the ID again so that the ID matches the ID it would have if it had been created via TF.
		// Use the current ID in case it changed in the resourceIdParserFunc.
		d.SetId(d.Id() + "/" + role)

		// Since condition titles can have any character in them, we can't separate them from any other
		// field the user might set in import (like the condition description and expression). So, we
		// have the user just specify the title and then read the upstream policy to set the full
		// condition. We can't rely on the read fn to do this for us because it looks for a match of the
		// full condition.
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return nil, err
		}
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return nil, err
		}
		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role == role && conditionKeyFromCondition(b.Condition).Title == conditionTitle {
				if binding != nil {
					return nil, fmt.Errorf("Cannot import IAM member with condition title %q, it matches multiple conditions", conditionTitle)
				}
				binding = b
			}
		}
		if binding != nil {
			if err := d.Set("condition", flattenIamCondition(binding.Condition)); err != nil {
				return nil, fmt.Errorf("Error setting condition: %s", err)
			}
			if k := conditionKeyFromCondition(binding.Condition); !k.Empty() {
				d.SetId(d.Id() + "/" + k.String())
			}
		}

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
	if c := expandIamCondition(d.Get("condition")); c != nil {
		b.Condition = c
	}
	return b
}

func expandIamCondition(v interface{}) *cloudresourcemanager.Expr {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	original := l[0].(map[string]interface{})
	return &cloudresourcemanager.Expr{
		Description:     original["description"].(string),
		Expression:      original["expression"].(string),
		Title:           original["title"].(string),
		ForceSendFields: []string{"Expression", "Title"},
	}
}

func flattenIamCondition(condition *cloudresourcemanager.Expr) []map[string]interface{} {
	if conditionKeyFromCondition(condition).Empty() {
		return nil
	}
	return []map[string]interface{}{
		{
			"expression":  condition.Expression,
			"title":       condition.Title,
			"description": condition.Description,
		},
	}
}
