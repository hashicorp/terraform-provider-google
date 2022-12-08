package google

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func iamMemberCaseDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	isCaseSensitive := iamMemberIsCaseSensitive(old) || iamMemberIsCaseSensitive(new)
	if isCaseSensitive {
		return old == new
	}
	return caseDiffSuppress(k, old, new, d)
}

func validateIAMMember(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}

	if matched, err := regexp.MatchString("^deleted", v); err != nil {
		return nil, []error{fmt.Errorf("error validating %s: %v", k, err)}
	} else if matched {
		return nil, []error{fmt.Errorf("invalid value for %s (Terraform does not support IAM members for deleted principals)", k)}
	}

	if matched, err := regexp.MatchString("(.+:.+|allUsers|allAuthenticatedUsers)", v); err != nil {
		return nil, []error{fmt.Errorf("error validating %s: %v", k, err)}
	} else if !matched {
		return nil, []error{fmt.Errorf("invalid value for %s (IAM members must have one of the values outlined here: https://cloud.google.com/billing/docs/reference/rest/v1/Policy#Binding)", k)}
	}
	return nil, nil
}

var IamMemberBaseSchema = map[string]*schema.Schema{
	"role": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"member": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: iamMemberCaseDiffSuppress,
		ValidateFunc:     validateIAMMember,
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

func iamMemberImport(newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		s := strings.Fields(d.Id())
		var id, role, member string
		if len(s) < 3 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Member id %s; expected 'resource_name role member [condition_title]'.", s)
		}

		var conditionTitle string
		if len(s) == 3 {
			id, role, member = s[0], s[1], s[2]
		} else {
			// condition titles can have any characters in them, so re-join the split string
			id, role, member, conditionTitle = s[0], s[1], s[2], strings.Join(s[3:], " ")
		}

		// Set the ID only to the first part so all IAM types can share the same resourceIdParserFunc.
		d.SetId(id)
		if err := d.Set("role", role); err != nil {
			return nil, fmt.Errorf("Error setting role: %s", err)
		}
		if err := d.Set("member", normalizeIamMemberCasing(member)); err != nil {
			return nil, fmt.Errorf("Error setting member: %s", err)
		}

		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}

		// Set the ID again so that the ID matches the ID it would have if it had been created via TF.
		// Use the current ID in case it changed in the resourceIdParserFunc.
		d.SetId(d.Id() + "/" + role + "/" + normalizeIamMemberCasing(member))

		// Read the upstream policy so we can set the full condition.
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
				containsMember := false
				for _, m := range b.Members {
					if strings.ToLower(m) == strings.ToLower(member) {
						containsMember = true
					}
				}
				if !containsMember {
					continue
				}

				if binding != nil {
					return nil, fmt.Errorf("Cannot import IAM member with condition title %q, it matches multiple conditions", conditionTitle)
				}
				binding = b
			}
		}
		if binding == nil {
			return nil, fmt.Errorf("Cannot find binding for %q with role %q, member %q, and condition title %q", updater.DescribeResource(), role, member, conditionTitle)
		}

		if err := d.Set("condition", flattenIamCondition(binding.Condition)); err != nil {
			return nil, fmt.Errorf("Error setting condition: %s", err)
		}
		if k := conditionKeyFromCondition(binding.Condition); !k.Empty() {
			d.SetId(d.Id() + "/" + k.String())
		}

		return []*schema.ResourceData{d}, nil
	}
}

func ResourceIamMember(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc, options ...func(*IamSettings)) *schema.Resource {
	return ResourceIamMemberWithBatching(parentSpecificSchema, newUpdaterFunc, resourceIdParser, IamBatchingDisabled, options...)
}

func ResourceIamMemberWithBatching(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc, enableBatching bool, options ...func(*IamSettings)) *schema.Resource {
	settings := &IamSettings{}
	for _, o := range options {
		o(settings)
	}

	return &schema.Resource{
		Create: resourceIamMemberCreate(newUpdaterFunc, enableBatching),
		Read:   resourceIamMemberRead(newUpdaterFunc),
		Delete: resourceIamMemberDelete(newUpdaterFunc, enableBatching),

		// if non-empty, this will be used to send a deprecation message when the
		// resource is used.
		DeprecationMessage: settings.DeprecationMessage,

		Schema: mergeSchemas(IamMemberBaseSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamMemberImport(newUpdaterFunc, resourceIdParser),
		},
		UseJSONNumber: true,
	}
}

func getResourceIamMember(d *schema.ResourceData) *cloudresourcemanager.Binding {
	b := &cloudresourcemanager.Binding{
		Members: []string{d.Get("member").(string)},
		Role:    d.Get("role").(string),
	}
	if c := expandIamCondition(d.Get("condition")); c != nil {
		b.Condition = c
	}
	return b
}

func resourceIamMemberCreate(newUpdaterFunc newResourceIamUpdaterFunc, enableBatching bool) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		memberBind := getResourceIamMember(d)
		modifyF := func(ep *cloudresourcemanager.Policy) error {
			// Merge the bindings together
			ep.Bindings = mergeBindings(append(ep.Bindings, memberBind))
			ep.Version = iamPolicyVersion
			return nil
		}
		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config,
				fmt.Sprintf("Create IAM Members %s %+v for %s", memberBind.Role, memberBind.Members[0], updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return err
		}
		d.SetId(updater.GetResourceId() + "/" + memberBind.Role + "/" + normalizeIamMemberCasing(memberBind.Members[0]))
		if k := conditionKeyFromCondition(memberBind.Condition); !k.Empty() {
			d.SetId(d.Id() + "/" + k.String())
		}
		return resourceIamMemberRead(newUpdaterFunc)(d, meta)
	}
}

func resourceIamMemberRead(newUpdaterFunc newResourceIamUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		eMember := getResourceIamMember(d)
		eCondition := conditionKeyFromCondition(eMember.Condition)
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Resource %q with IAM Member: Role %q Member %q", updater.DescribeResource(), eMember.Role, eMember.Members[0]))
		}
		log.Print(spew.Sprintf("[DEBUG]: Retrieved policy for %s: %#v\n", updater.DescribeResource(), p))
		log.Printf("[DEBUG]: Looking for binding with role %q and condition %#v", eMember.Role, eCondition)

		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role == eMember.Role && conditionKeyFromCondition(b.Condition) == eCondition {
				binding = b
				break
			}
		}

		if binding == nil {
			log.Printf("[DEBUG]: Binding for role %q with condition %#v does not exist in policy of %s, removing member %q from state.", eMember.Role, eCondition, updater.DescribeResource(), eMember.Members[0])
			d.SetId("")
			return nil
		}

		log.Printf("[DEBUG]: Looking for member %q in found binding", eMember.Members[0])
		var member string
		for _, m := range binding.Members {
			if strings.ToLower(m) == strings.ToLower(eMember.Members[0]) {
				member = m
			}
		}

		if member == "" {
			log.Printf("[DEBUG]: Member %q for binding for role %q with condition %#v does not exist in policy of %s, removing from state.", eMember.Members[0], eMember.Role, eCondition, updater.DescribeResource())
			d.SetId("")
			return nil
		}

		if err := d.Set("etag", p.Etag); err != nil {
			return fmt.Errorf("Error setting etag: %s", err)
		}
		if err := d.Set("member", member); err != nil {
			return fmt.Errorf("Error setting member: %s", err)
		}
		if err := d.Set("role", binding.Role); err != nil {
			return fmt.Errorf("Error setting role: %s", err)
		}
		if err := d.Set("condition", flattenIamCondition(binding.Condition)); err != nil {
			return fmt.Errorf("Error setting condition: %s", err)
		}
		return nil
	}
}

func resourceIamMemberDelete(newUpdaterFunc newResourceIamUpdaterFunc, enableBatching bool) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		memberBind := getResourceIamMember(d)
		modifyF := func(ep *cloudresourcemanager.Policy) error {
			// Merge the bindings together
			ep.Bindings = subtractFromBindings(ep.Bindings, memberBind)
			return nil
		}
		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config,
				fmt.Sprintf("Delete IAM Members %s %s for %q", memberBind.Role, memberBind.Members[0], updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Resource %s for IAM Member (role %q, %q)", updater.GetResourceId(), memberBind.Members[0], memberBind.Role))
		}
		return resourceIamMemberRead(newUpdaterFunc)(d, meta)
	}
}
