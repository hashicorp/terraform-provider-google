package google

import (
	"fmt"
	"log"
	"time"

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
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"members": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"etag": &schema.Schema{
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

	for {
		backoff := time.Second
		// Get the existing bindings
		log.Println("[DEBUG]: Retrieving policy for project", pid)
		ep, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for project %q: %+v\n", pid, ep)

		// Merge the bindings together
		ep.Bindings = mergeBindings(append(ep.Bindings, p))
		log.Printf("[DEBUG]: Setting policy for project %q to %+v\n", pid, ep)
		err = setProjectIamPolicy(ep, config, pid)
		if err != nil && isConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return fmt.Errorf("Error applying IAM policy to project %q: too many concurrent policy changes.\n")
			}
			continue
		} else if err != nil {
			return fmt.Errorf("Error applying IAM policy to project: %v", err)
		}
		break
	}
	log.Printf("[DEBUG]: Set policy for project %q", pid)
	d.SetId(pid + ":" + p.Role)
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

	for {
		backoff := time.Second
		log.Println("[DEBUG]: Retrieving policy for project", pid)
		p, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for project %q: %+v\n", pid, p)

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

		log.Printf("[DEBUG]: Setting policy for project %q to %+v\n", pid, p)
		err = setProjectIamPolicy(p, config, pid)
		if err != nil && isConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return fmt.Errorf("Error applying IAM policy to project %q: too many concurrent policy changes.\n")
			}
			continue
		} else if err != nil {
			return fmt.Errorf("Error applying IAM policy to project: %v", err)
		}
		break
	}
	log.Printf("[DEBUG]: Set policy for project %q\n", pid)

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

	for {
		backoff := time.Second
		log.Println("[DEBUG]: Retrieving policy for project", pid)
		p, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for project %q: %+v\n", pid, p)

		toRemove := -1
		for pos, b := range p.Bindings {
			if b.Role != binding.Role {
				continue
			}
			toRemove = pos
			break
		}
		if toRemove < 0 {
			return resourceGoogleProjectIamBindingRead(d, meta)
		}

		p.Bindings = append(p.Bindings[:toRemove], p.Bindings[toRemove+1:]...)

		log.Printf("[DEBUG]: Setting policy for project %q to %+v\n", pid, p)
		err = setProjectIamPolicy(p, config, pid)
		if err != nil && isConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return fmt.Errorf("Error applying IAM policy to project %q: too many concurrent policy changes.\n")
			}
			continue
		} else if err != nil {
			return fmt.Errorf("Error applying IAM policy to project: %v", err)
		}
		break
	}
	log.Printf("[DEBUG]: Set policy for project %q\n", pid)

	return resourceGoogleProjectIamBindingRead(d, meta)
}

// Get a cloudresourcemanager.Binding from a schema.ResourceData
func getResourceIamBinding(d *schema.ResourceData) *cloudresourcemanager.Binding {
	members := d.Get("members").(*schema.Set).List()
	m := make([]string, 0, len(members))
	for _, member := range members {
		m = append(m, member.(string))
	}
	return &cloudresourcemanager.Binding{
		Members: m,
		Role:    d.Get("role").(string),
	}
}

func projectIamBindingMutexKey(pid, role string) string {
	return fmt.Sprintf("google-project-iam-binding-%s-%s", pid, role)
}
