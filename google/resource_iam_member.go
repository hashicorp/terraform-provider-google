package google

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func iamMemberImport(resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		s := strings.Fields(d.Id())
		if len(s) != 3 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Member id %s; expected 'resource_name role member'.", s)
		}
		id, role, member := s[0], s[1], s[2]

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
	return &schema.Resource{
		Create: resourceIamMemberCreate(newUpdaterFunc),
		Read:   resourceIamMemberRead(newUpdaterFunc),
		Delete: resourceIamMemberDelete(newUpdaterFunc),
		Schema: mergeSchemas(IamMemberBaseSchema, parentSpecificSchema),
		Importer: &schema.ResourceImporter{
			State: iamMemberImport(resourceIdParser),
		},
	}
}

func getResourceIamMember(d *schema.ResourceData) *cloudresourcemanager.Binding {
	return &cloudresourcemanager.Binding{
		Members: []string{d.Get("member").(string)},
		Role:    d.Get("role").(string),
	}
}

func resourceIamMemberCreate(newUpdaterFunc newResourceIamUpdaterFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		p := getResourceIamMember(d)
		err = iamPolicyReadModifyWrite(updater, func(ep *cloudresourcemanager.Policy) error {
			// Merge the bindings together
			ep.Bindings = mergeBindings(append(ep.Bindings, p))
			return nil
		})
		if err != nil {
			return err
		}
		d.SetId(updater.GetResourceId() + "/" + p.Role + "/" + strings.ToLower(p.Members[0]))
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
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG]: Binding of member %q with role %q does not exist for non-existent resource %s, removing from state.", eMember.Members[0], eMember.Role, updater.DescribeResource())
				d.SetId("")
				return nil
			}
			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for %s: %+v\n", updater.DescribeResource(), p)

		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role != eMember.Role {
				continue
			}
			binding = b
			break
		}
		if binding == nil {
			log.Printf("[DEBUG]: Binding for role %q does not exist in policy of %s, removing member %q from state.", eMember.Role, updater.DescribeResource(), eMember.Members[0])
			d.SetId("")
			return nil
		}
		var member string
		for _, m := range binding.Members {
			if strings.ToLower(m) == strings.ToLower(eMember.Members[0]) {
				member = m
			}
		}
		if member == "" {
			log.Printf("[DEBUG]: Member %q for binding for role %q does not exist in policy of %s, removing from state.", eMember.Members[0], eMember.Role, updater.DescribeResource())
			d.SetId("")
			return nil
		}
		d.Set("etag", p.Etag)
		d.Set("member", member)
		d.Set("role", binding.Role)
		return nil
	}
}

func resourceIamMemberDelete(newUpdaterFunc newResourceIamUpdaterFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		member := getResourceIamMember(d)
		err = iamPolicyReadModifyWrite(updater, func(p *cloudresourcemanager.Policy) error {
			bindingToRemove := -1
			for pos, b := range p.Bindings {
				if b.Role != member.Role {
					continue
				}
				bindingToRemove = pos
				break
			}
			if bindingToRemove < 0 {
				log.Printf("[DEBUG]: Binding for role %q does not exist in policy of project %q, so member %q can't be on it.", member.Role, updater.GetResourceId(), member.Members[0])
				return nil
			}
			binding := p.Bindings[bindingToRemove]
			memberToRemove := -1
			for pos, m := range binding.Members {
				if strings.ToLower(m) != strings.ToLower(member.Members[0]) {
					continue
				}
				memberToRemove = pos
				break
			}
			if memberToRemove < 0 {
				log.Printf("[DEBUG]: Member %q for binding for role %q does not exist in policy of project %q.", member.Members[0], member.Role, updater.GetResourceId())
				return nil
			}
			binding.Members = append(binding.Members[:memberToRemove], binding.Members[memberToRemove+1:]...)
			if len(binding.Members) == 0 {
				// If there is no member left for the role, remove the binding altogether
				p.Bindings = append(p.Bindings[:bindingToRemove], p.Bindings[bindingToRemove+1:]...)
			} else {
				p.Bindings[bindingToRemove] = binding
			}

			return nil
		})
		if err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG]: Member %q for binding for role %q does not exist for non-existent resource %q.", member.Members[0], member.Role, updater.GetResourceId())
				return nil
			}
			return err
		}

		return resourceIamMemberRead(newUpdaterFunc)(d, meta)
	}
}
