package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
)

func dataSourceGoogleServiceAccountKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp(ServiceAccountKeyNameRegex),
			},
			"public_key_type": {
				Type:         schema.TypeString,
				Default:      "TYPE_X509_PEM_FILE",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"TYPE_NONE", "TYPE_X509_PEM_FILE", "TYPE_RAW_PUBLIC_KEY"}, false),
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
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

	keyName := d.Get("name").(string)

	// Validate name since interpolated values (i.e from a key or service
	// account resource) will not get validated at plan time.
	r := regexp.MustCompile(ServiceAccountKeyNameRegex)
	if !r.MatchString(keyName) {
		return fmt.Errorf("invalid key name %q does not match regexp %q", keyName, ServiceAccountKeyNameRegex)
	}

	publicKeyType := d.Get("public_key_type").(string)

	// Confirm the service account key exists
	sak, err := config.clientIAM.Projects.ServiceAccounts.Keys.Get(keyName).PublicKeyType(publicKeyType).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", keyName))
	}

	d.SetId(sak.Name)

	d.Set("name", sak.Name)
	d.Set("key_algorithm", sak.KeyAlgorithm)
	d.Set("public_key", sak.PublicKeyData)

	return nil
}
