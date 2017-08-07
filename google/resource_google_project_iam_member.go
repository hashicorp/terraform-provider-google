package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func resourceGoogleProjectIamMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamMemberCreate,
		Read:   resourceGoogleProjectIamMemberRead,
		Delete: resourceGoogleProjectIamMemberDelete,

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
			"member": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"etag": {
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

	err = projectIamPolicyReadModifyWrite(d, config, pid, func(ep *cloudresourcemanager.Policy) error {
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
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(pid + "/" + p.Role + "/" + p.Members[0])
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
		log.Printf("[DEBUG]: Binding for role %q does not exist in policy of project %q, removing member %q from state.", eMember.Role, pid, eMember.Members[0])
		d.SetId("")
		return nil
	}
	var member string
	for _, m := range binding.Members {
		if m == eMember.Members[0] {
			member = m
		}
	}
	if member == "" {
		log.Printf("[DEBUG]: Member %q for binding for role %q does not exist in policy of project %q, removing from state.", eMember.Members[0], eMember.Role, pid)
		d.SetId("")
		return nil
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

	err = projectIamPolicyReadModifyWrite(d, config, pid, func(p *cloudresourcemanager.Policy) error {
		bindingToRemove := -1
		for pos, b := range p.Bindings {
			if b.Role != member.Role {
				continue
			}
			bindingToRemove = pos
			break
		}
		if bindingToRemove < 0 {
			log.Printf("[DEBUG]: Binding for role %q does not exist in policy of project %q, so member %q can't be on it.", member.Role, pid, member.Members[0])
			return nil
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
			log.Printf("[DEBUG]: Member %q for binding for role %q does not exist in policy of project %q.", member.Members[0], member.Role, pid)
			return nil
		}
		binding.Members = append(binding.Members[:memberToRemove], binding.Members[memberToRemove+1:]...)
		p.Bindings[bindingToRemove] = binding
		return nil
	})
	if err != nil {
		return err
	}

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
