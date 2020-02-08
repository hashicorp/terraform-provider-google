package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRFC1035Name(6, 30),
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
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
	description := d.Get("description").(string)

	sa := &iam.ServiceAccount{
		DisplayName: displayName,
		Description: description,
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
	// This API is meant to be synchronous, but in practice it shows the old value for
	// a few milliseconds after the update goes through.  A second is more than enough
	// time to ensure following reads are correct.
	time.Sleep(time.Second)

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
	d.Set("description", sa.Description)
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
	sa, err := config.clientIAM.Projects.ServiceAccounts.Get(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error retrieving service account %q: %s", d.Id(), err)
	}
	updateMask := make([]string, 0)
	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}
	if d.HasChange("display_name") {
		updateMask = append(updateMask, "display_name")
	}
	_, err = config.clientIAM.Projects.ServiceAccounts.Patch(d.Id(),
		&iam.PatchServiceAccountRequest{
			UpdateMask: strings.Join(updateMask, ","),
			ServiceAccount: &iam.ServiceAccount{
				DisplayName: d.Get("display_name").(string),
				Description: d.Get("description").(string),
				Etag:        sa.Etag,
			},
		}).Do()
	if err != nil {
		return err
	}
	// See comment in Create.
	time.Sleep(time.Second)

	return nil
}

func resourceGoogleServiceAccountImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/serviceAccounts/(?P<email>[^/]+)",
		"(?P<project>[^/]+)/(?P<email>[^/]+)",
		"(?P<email>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/serviceAccounts/{{email}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
