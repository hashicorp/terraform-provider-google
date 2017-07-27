package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func resourceGoogleProjectIamBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamBindingCreate,
		Read:   resourceGoogleProjectIamBindingRead,
		Update: resourceGoogleProjectIamBindingUpdate,
		Delete: resourceGoogleProjectIamBindingDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: {
					Type: schema.TypeString,
				},
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGoogleProjectIamBindingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the binding in the template
	log.Println("[DEBUG]: Reading google_project_iam_binding")
	p := getResourceIamBinding(d)
	mutexKV.Lock(projectIamBindingMutexKey(pid, p.Role))
	defer mutexKV.Unlock(projectIamBindingMutexKey(pid, p.Role))

	err = projectIamPolicyReadModifyWrite(d, config, pid, func(ep *cloudresourcemanager.Policy) error {
		// Merge the bindings together
		ep.Bindings = mergeBindings(append(ep.Bindings, p))
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(pid + "/" + p.Role)
	return resourceGoogleProjectIamBindingRead(d, meta)
}

func resourceGoogleProjectIamBindingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	eBinding := getResourceIamBinding(d)

	log.Println("[DEBUG]: Retrieving policy for project", pid)
	p, err := getProjectIamPolicy(pid, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]: Retrieved policy for project %q: %+v\n", pid, p)

	var binding *cloudresourcemanager.Binding
	for _, b := range p.Bindings {
		if b.Role != eBinding.Role {
			continue
		}
		binding = b
		break
	}
	if binding == nil {
		log.Printf("[DEBUG]: Binding for role %q not found in policy for %q, removing from state file.\n", eBinding.Role, pid)
		d.SetId("")
		return nil
	}
	d.Set("etag", p.Etag)
	d.Set("members", binding.Members)
	d.Set("role", binding.Role)
	return nil
}

func resourceGoogleProjectIamBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	binding := getResourceIamBinding(d)
	mutexKV.Lock(projectIamBindingMutexKey(pid, binding.Role))
	defer mutexKV.Unlock(projectIamBindingMutexKey(pid, binding.Role))

	err = projectIamPolicyReadModifyWrite(d, config, pid, func(p *cloudresourcemanager.Policy) error {
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

	return resourceGoogleProjectIamBindingRead(d, meta)
}

func resourceGoogleProjectIamBindingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	binding := getResourceIamBinding(d)
	mutexKV.Lock(projectIamBindingMutexKey(pid, binding.Role))
	defer mutexKV.Unlock(projectIamBindingMutexKey(pid, binding.Role))

	err = projectIamPolicyReadModifyWrite(d, config, pid, func(p *cloudresourcemanager.Policy) error {
		toRemove := -1
		for pos, b := range p.Bindings {
			if b.Role != binding.Role {
				continue
			}
			toRemove = pos
			break
		}
		if toRemove < 0 {
			log.Printf("[DEBUG]: Policy bindings for project %q did not include a binding for role %q, no need to delete", pid, binding.Role)
			d.SetId("")
			return nil
		}

		p.Bindings = append(p.Bindings[:toRemove], p.Bindings[toRemove+1:]...)
		return nil
	})
	if err != nil {
		return err
	}

	return resourceGoogleProjectIamBindingRead(d, meta)
}

// Get a cloudresourcemanager.Binding from a schema.ResourceData
func getResourceIamBinding(d *schema.ResourceData) *cloudresourcemanager.Binding {
	members := d.Get("members").(*schema.Set).List()
	return &cloudresourcemanager.Binding{
		Members: convertStringArr(members),
		Role:    d.Get("role").(string),
	}
}

func projectIamBindingMutexKey(pid, role string) string {
	return fmt.Sprintf("google-project-iam-binding-%s-%s", pid, role)
}
