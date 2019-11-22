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
		DiffSuppressFunc: caseDiffSuppress,
		ValidateFunc:     validation.StringDoesNotMatch(regexp.MustCompile("^deleted:"), "Terraform does not support IAM members for deleted principals"),
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
		if len(s) != 3 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Member id %s; expected 'resource_name role member'.", s)
		}
		id, role, member = s[0], s[1], s[2]

		// Set the ID only to the first part so all IAM types can share the same resourceIdParserFunc.
		d.SetId(id)
		d.Set("role", role)
		d.Set("member", strings.ToLower(member))

		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}

		// Set the ID again so that the ID matches the ID it would have if it had been created via TF.
		// Use the current ID in case it changed in the resourceIdParserFunc.
		d.SetId(d.Id() + "/" + role + "/" + strings.ToLower(member))

		return []*schema.ResourceData{d}, nil
	}
}

func ResourceIamMember(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc) *schema.Resource {
	return ResourceIamMemberWithBatching(parentSpecificSchema, newUpdaterFunc, resourceIdParser, IamBatchingDisabled)
}

func ResourceIamMemberWithBatching(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc, enableBatching bool) *schema.Resource {
	return &schema.Resource{
		Create: resourceIamMemberCreate(newUpdaterFunc, enableBatching),
		Read:   resourceIamMemberRead(newUpdaterFunc),
		Delete: resourceIamMemberDelete(newUpdaterFunc, enableBatching),
		Schema: mergeSchemas(IamMemberBaseSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamMemberImport(newUpdaterFunc, resourceIdParser),
		},
	}
}

func getResourceIamMember(d *schema.ResourceData) *cloudresourcemanager.Binding {
	b := &cloudresourcemanager.Binding{
		Members: []string{d.Get("member").(string)},
		Role:    d.Get("role").(string),
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
			return nil
		}
		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config,
				fmt.Sprintf("Create IAM Members %s %+v for %q", memberBind.Role, memberBind.Members[0], updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return err
		}
		d.SetId(updater.GetResourceId() + "/" + memberBind.Role + "/" + strings.ToLower(memberBind.Members[0]))
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
		log.Printf("[DEBUG]: Retrieved policy for %s: %+v\n", updater.DescribeResource(), p)
		log.Printf("[DEBUG]: Looking for binding with role %q and condition %+v", eMember.Role, eCondition)

		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role == eMember.Role && conditionKeyFromCondition(b.Condition) == eCondition {
				binding = b
				break
			}
		}

		if binding == nil {
			log.Printf("[DEBUG]: Binding for role %q with condition %+v does not exist in policy of %s, removing member %q from state.", eMember.Role, eCondition, updater.DescribeResource(), eMember.Members[0])
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
			log.Printf("[DEBUG]: Member %q for binding for role %q with condition %+v does not exist in policy of %s, removing from state.", eMember.Members[0], eMember.Role, eCondition, updater.DescribeResource())
			d.SetId("")
			return nil
		}

		d.Set("etag", p.Etag)
		d.Set("member", member)
		d.Set("role", binding.Role)
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
