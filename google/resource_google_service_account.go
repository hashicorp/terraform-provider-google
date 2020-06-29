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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The e-mail address of the service account. This value should be referenced from any google_iam_policy data sources that would grant the service account privileges.`,
			},
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique id of the service account.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fully-qualified name of the service account.`,
			},
			"account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRFC1035Name(6, 30),
				Description:  `The account id that is used to generate the service account email address and a stable unique id. It is unique within a project, must be 6-30 characters long, and match the regular expression [a-z]([-a-z0-9]*[a-z0-9]) to comply with RFC1035. Changing this forces a new service account to be created.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The display name for the service account. Can be updated without creating a new resource.`,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
				Description:  `A text description of the service account. Must be less than or equal to 256 UTF-8 bytes.`,
			},
			"project": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the project that the service account will be created in. Defaults to the provider project configuration.`,
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

	err = retryTimeDuration(func() (operr error) {
		_, saerr := config.clientIAM.Projects.ServiceAccounts.Get(d.Id()).Do()
		return saerr
	}, d.Timeout(schema.TimeoutCreate), isNotFoundRetryableError("service account creation"))

	if err != nil {
		return fmt.Errorf("Error reading service account after creation: %s", err)
	}

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
	// This API is meant to be synchronous, but in practice it shows the old value for
	// a few milliseconds after the update goes through. 5 seconds is more than enough
	// time to ensure following reads are correct.
	time.Sleep(time.Second * 5)

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
