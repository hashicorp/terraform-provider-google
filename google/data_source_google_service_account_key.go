package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceGoogleServiceAccountKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountKeyRead,

		Schema: map[string]*schema.Schema{
			// Required
			"service_account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Optional
			"public_key_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "TYPE_X509_PEM_FILE",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"TYPE_NONE", "TYPE_X509_PEM_FILE", "TYPE_RAW_PUBLIC_KEY"}, false),
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleServiceAccountKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Get the service account id as the fully qualified name
	serviceAccountID := d.Get("service_account_id").(string)

	// If the service account id isn't already the fully qualified name
	if !strings.HasPrefix(serviceAccountID, "projects/") {

		// If the service account id is an email
		if strings.Contains(serviceAccountID, "@") {
			serviceAccountID = "projects/-/serviceAccounts/" + serviceAccountID
		} else {
			// Get the project from the resource or fallback to the project
			// in the provider configuration
			project, err := getProject(d, config)
			if err != nil {
				return err
			}
			// If the service account id doesn't contain the email, build it
			serviceAccountID = fmt.Sprintf("projects/-/serviceAccounts/%s@%s.iam.gserviceaccount.com", serviceAccountID, project)
		}
	}

	publicKeyType := d.Get("public_key_type").(string)

	// Confirm the service account key exists
	sak, err := config.clientIAM.Projects.ServiceAccounts.Keys.Get(serviceAccountID).PublicKeyType(publicKeyType).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", serviceAccountID))
	}

	d.SetId(sak.Name)

	d.Set("name", sak.Name)
	d.Set("key_algorithm", sak.KeyAlgorithm)
	d.Set("public_key", sak.PublicKeyData)

	return nil
}
