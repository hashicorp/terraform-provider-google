package google

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
			Type: schema.TypeString,
		},
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ResourceIamBinding(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc) *schema.Resource {
	return &schema.Resource{
		Create: resourceIamBindingCreate(newUpdaterFunc),
		Read:   resourceIamBindingRead(newUpdaterFunc),
		Update: resourceIamBindingUpdate(newUpdaterFunc),
		Delete: resourceIamBindingDelete(newUpdaterFunc),
		Schema: mergeSchemas(iamBindingSchema, parentSpecificSchema),
	}
}

func ResourceIamBindingWithImport(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceIamUpdaterFunc, resourceIdParser resourceIdParserFunc) *schema.Resource {
	r := ResourceIamBinding(parentSpecificSchema, newUpdaterFunc)
	r.Importer = &schema.ResourceImporter{
		State: iamBindingImport(resourceIdParser),
	}
	return r
}

func resourceIamBindingCreate(newUpdaterFunc newResourceIamUpdaterFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		p := getResourceIamBinding(d)
		err = iamPolicyReadModifyWrite(updater, func(ep *cloudresourcemanager.Policy) error {
			// Creating a binding does not remove existing members if they are not in the provided members list.
			// This prevents removing existing permission without the user's knowledge.
			// Instead, a diff is shown in that case after creation. Subsequent calls to update will remove any
			// existing members not present in the provided list.
			ep.Bindings = mergeBindings(append(ep.Bindings, p))
			return nil
		})
		if err != nil {
			return err
		}
		d.SetId(updater.GetResourceId() + "/" + p.Role)
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
		p, err := updater.GetResourceIamPolicy()
		if err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG]: Binding for role %q not found for non-existant resource %s, removing from state file.", updater.DescribeResource(), eBinding.Role)
				d.SetId("")
				return nil
			}

			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for %s: %+v", updater.DescribeResource(), p)

		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role != eBinding.Role {
				continue
			}
			binding = b
			break
		}
		if binding == nil {
			log.Printf("[DEBUG]: Binding for role %q not found in policy for %s, removing from state file.", eBinding.Role, updater.DescribeResource())
			d.SetId("")
			return nil
		}
		d.Set("etag", p.Etag)
		d.Set("members", binding.Members)
		d.Set("role", binding.Role)
		return nil
	}
}

func iamBindingImport(resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		s := strings.Fields(d.Id())
		if len(s) != 2 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Binding id %s; expected 'resource_name role'.", s)
		}
		id, role := s[0], s[1]

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

func resourceIamBindingUpdate(newUpdaterFunc newResourceIamUpdaterFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		binding := getResourceIamBinding(d)
		err = iamPolicyReadModifyWrite(updater, func(p *cloudresourcemanager.Policy) error {
			var found bool
			for pos, b := range p.Bindings {
				if b.Role != binding.Role {
					continue
				}
				found = true
				p.Bindings[pos] = binding
				break
			}
			if !found {
				p.Bindings = append(p.Bindings, binding)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceIamBindingRead(newUpdaterFunc)(d, meta)
	}
}

func resourceIamBindingDelete(newUpdaterFunc newResourceIamUpdaterFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		binding := getResourceIamBinding(d)
		err = iamPolicyReadModifyWrite(updater, func(p *cloudresourcemanager.Policy) error {
			toRemove := -1
			for pos, b := range p.Bindings {
				if b.Role != binding.Role {
					continue
				}
				toRemove = pos
				break
			}
			if toRemove < 0 {
				log.Printf("[DEBUG]: Policy bindings for %s did not include a binding for role %q", updater.DescribeResource(), binding.Role)
				return nil
			}

			p.Bindings = append(p.Bindings[:toRemove], p.Bindings[toRemove+1:]...)
			return nil
		})
		if err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG]: Resource %s is missing or deleted, marking policy binding as deleted", updater.DescribeResource())
				return nil
			}
			return err
		}

		return resourceIamBindingRead(newUpdaterFunc)(d, meta)
	}
}

func getResourceIamBinding(d *schema.ResourceData) *cloudresourcemanager.Binding {
	members := d.Get("members").(*schema.Set).List()
	return &cloudresourcemanager.Binding{
		Members: convertStringArr(members),
		Role:    d.Get("role").(string),
	}
}
