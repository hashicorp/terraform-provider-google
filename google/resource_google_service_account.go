package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/iam/v1"
)

func resourceGoogleServiceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleServiceAccountCreate,
		Read:   resourceGoogleServiceAccountRead,
		Delete: resourceGoogleServiceAccountDelete,
		Update: resourceGoogleServiceAccountUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceGoogleServiceAccountImport,
		},
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRFC1035Name(6, 30),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"policy_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Use the 'google_service_account_iam_policy' resource to define policies for a service account",
			},
		},
	}
}

func resourceGoogleServiceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	aid := d.Get("account_id").(string)
	displayName := d.Get("display_name").(string)

	sa := &iam.ServiceAccount{
		DisplayName: displayName,
	}

	r := &iam.CreateServiceAccountRequest{
		AccountId:      aid,
		ServiceAccount: sa,
	}

	sa, err = config.clientIAM.Projects.ServiceAccounts.Create("projects/"+project, r).Do()
	if err != nil {
		return fmt.Errorf("Error creating service account: %s", err)
	}

	d.SetId(sa.Name)

	return resourceGoogleServiceAccountRead(d, meta)
}

func resourceGoogleServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Confirm the service account exists
	sa, err := config.clientIAM.Projects.ServiceAccounts.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", d.Id()))
	}

	d.Set("email", sa.Email)
	d.Set("unique_id", sa.UniqueId)
	d.Set("project", sa.ProjectId)
	d.Set("account_id", strings.Split(sa.Email, "@")[0])
	d.Set("name", sa.Name)
	d.Set("display_name", sa.DisplayName)
	return nil
}

func resourceGoogleServiceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Id()
	_, err := config.clientIAM.Projects.ServiceAccounts.Delete(name).Do()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceGoogleServiceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	if ok := d.HasChange("display_name"); ok {
		sa, err := config.clientIAM.Projects.ServiceAccounts.Get(d.Id()).Do()
		if err != nil {
			return fmt.Errorf("Error retrieving service account %q: %s", d.Id(), err)
		}
		_, err = config.clientIAM.Projects.ServiceAccounts.Update(d.Id(),
			&iam.ServiceAccount{
				DisplayName: d.Get("display_name").(string),
				Etag:        sa.Etag,
			}).Do()
		if err != nil {
			return fmt.Errorf("Error updating service account %q: %s", d.Id(), err)
		}
	}

	return nil
}

// Retrieve the existing IAM Policy for a service account
func getServiceAccountIamPolicy(sa string, config *Config) (*iam.Policy, error) {
	p, err := config.clientIAM.Projects.ServiceAccounts.GetIamPolicy(sa).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for service account %q: %s", sa, err)
	}
	return p, nil
}

// Convert a map of roles->members to a list of Binding
func saRolesToMembersBinding(m map[string]map[string]bool) []*iam.Binding {
	bindings := make([]*iam.Binding, 0)
	for role, members := range m {
		b := iam.Binding{
			Role:    role,
			Members: make([]string, 0),
		}
		for m, _ := range members {
			b.Members = append(b.Members, m)
		}
		bindings = append(bindings, &b)
	}
	return bindings
}

// Map a role to a map of members, allowing easy merging of multiple bindings.
func saRolesToMembersMap(bindings []*iam.Binding) map[string]map[string]bool {
	bm := make(map[string]map[string]bool)
	// Get each binding
	for _, b := range bindings {
		// Initialize members map
		if _, ok := bm[b.Role]; !ok {
			bm[b.Role] = make(map[string]bool)
		}
		// Get each member (user/principal) for the binding
		for _, m := range b.Members {
			// Add the member
			bm[b.Role][m] = true
		}
	}
	return bm
}

// Merge multiple Bindings such that Bindings with the same Role result in
// a single Binding with combined Members
func saMergeBindings(bindings []*iam.Binding) []*iam.Binding {
	bm := saRolesToMembersMap(bindings)
	rb := make([]*iam.Binding, 0)

	for role, members := range bm {
		var b iam.Binding
		b.Role = role
		b.Members = make([]string, 0)
		for m, _ := range members {
			b.Members = append(b.Members, m)
		}
		rb = append(rb, &b)
	}

	return rb
}

func resourceGoogleServiceAccountImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	parseImportId([]string{
		"projects/(?P<project>[^/]+)/serviceAccounts/(?P<email>[^/]+)",
		"(?P<project>[^/]+)/(?P<email>[^/]+)",
		"(?P<email>[^/]+)"}, d, config)

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/serviceAccounts/{{email}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
