package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func resourceGoogleProjectIamMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamMemberCreate,
		Read:   resourceGoogleProjectIamMemberRead,
		Delete: resourceGoogleProjectIamMemberDelete,

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
			"member": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGoogleProjectIamMemberCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the binding in the template
	log.Println("[DEBUG]: Reading google_project_iam_member")
	p := getResourceIamMember(d)
	mutexKV.Lock(projectIamMemberMutexKey(pid, p.Role, p.Members[0]))
	defer mutexKV.Unlock(projectIamMemberMutexKey(pid, p.Role, p.Members[0]))

	for {
		backoff := time.Second
		// Get the existing bindings
		log.Println("[DEBUG]: Retrieving policy for project", pid)
		ep, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for project %q: %+v\n", pid, ep)

		// find the binding
		var binding *cloudresourcemanager.Binding
		for _, b := range ep.Bindings {
			if b.Role != p.Role {
				continue
			}
			binding = b
			break
		}
		if binding == nil {
			binding = &cloudresourcemanager.Binding{
				Role:    p.Role,
				Members: p.Members,
			}
		}

		// Merge the bindings together
		ep.Bindings = mergeBindings(append(ep.Bindings, p))
		log.Printf("[DEBUG]: Setting policy for project %q to %+v\n", pid, ep)
		err = setProjectIamPolicy(ep, config, pid)
		if err != nil && isConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return fmt.Errorf("Error applying IAM policy to project %q: too many concurrent policy changes.\n", pid)
			}
			continue
		} else if err != nil {
			return fmt.Errorf("Error applying IAM policy to project: %v", err)
		}
		break
	}
	log.Printf("[DEBUG]: Set policy for project %q", pid)
	d.SetId(pid + ":" + p.Role + ":" + p.Members[0])
	return resourceGoogleProjectIamMemberRead(d, meta)
}

func resourceGoogleProjectIamMemberRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	eMember := getResourceIamMember(d)

	log.Println("[DEBUG]: Retrieving policy for project", pid)
	p, err := getProjectIamPolicy(pid, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]: Retrieved policy for project %q: %+v\n", pid, p)

	var binding *cloudresourcemanager.Binding
	for _, b := range p.Bindings {
		if b.Role != eMember.Role {
			continue
		}
		binding = b
		break
	}
	if binding == nil {
		d.SetId("")
		return nil
	}
	var member string
	for _, m := range binding.Members {
		if m == eMember.Members[0] {
			member = m
		}
	}
	d.Set("etag", p.Etag)
	d.Set("member", member)
	d.Set("role", binding.Role)
	return nil
}

func resourceGoogleProjectIamMemberDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	member := getResourceIamMember(d)
	mutexKV.Lock(projectIamMemberMutexKey(pid, member.Role, member.Members[0]))
	defer mutexKV.Unlock(projectIamMemberMutexKey(pid, member.Role, member.Members[0]))

	for {
		backoff := time.Second
		log.Println("[DEBUG]: Retrieving policy for project", pid)
		p, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for project %q: %+v\n", pid, p)

		bindingToRemove := -1
		for pos, b := range p.Bindings {
			if b.Role != member.Role {
				continue
			}
			bindingToRemove = pos
			break
		}
		if bindingToRemove < 0 {
			return resourceGoogleProjectIamMemberRead(d, meta)
		}
		binding := p.Bindings[bindingToRemove]
		memberToRemove := -1
		for pos, m := range binding.Members {
			if m != member.Members[0] {
				continue
			}
			memberToRemove = pos
			break
		}
		if memberToRemove < 0 {
			return resourceGoogleProjectIamMemberRead(d, meta)
		}
		binding.Members = append(binding.Members[:memberToRemove], binding.Members[memberToRemove+1:]...)
		p.Bindings[bindingToRemove] = binding

		log.Printf("[DEBUG]: Setting policy for project %q to %+v\n", pid, p)
		err = setProjectIamPolicy(p, config, pid)
		if err != nil && isConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return fmt.Errorf("Error applying IAM policy to project %q: too many concurrent policy changes.\n", pid)
			}
			continue
		} else if err != nil {
			return fmt.Errorf("Error applying IAM policy to project: %v", err)
		}
		break
	}
	log.Printf("[DEBUG]: Set policy for project %q\n", pid)

	return resourceGoogleProjectIamMemberRead(d, meta)
}

// Get a cloudresourcemanager.Binding from a schema.ResourceData
func getResourceIamMember(d *schema.ResourceData) *cloudresourcemanager.Binding {
	return &cloudresourcemanager.Binding{
		Members: []string{d.Get("member").(string)},
		Role:    d.Get("role").(string),
	}
}

func projectIamMemberMutexKey(pid, role, member string) string {
	return fmt.Sprintf("google-project-iam-member-%s-%s-%s", pid, role, member)
}
