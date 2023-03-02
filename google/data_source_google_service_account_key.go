package google

import (
	"fmt"

	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceGoogleServiceAccountKey() *schema.Resource {
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
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	keyName := d.Get("name").(string)

	// Validate name since interpolated values (i.e from a key or service
	// account resource) will not get validated at plan time.
	r := regexp.MustCompile(ServiceAccountKeyNameRegex)
	if !r.MatchString(keyName) {
		return fmt.Errorf("invalid key name %q does not match regexp %q", keyName, ServiceAccountKeyNameRegex)
	}

	publicKeyType := d.Get("public_key_type").(string)

	// Confirm the service account key exists
	sak, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Keys.Get(keyName).PublicKeyType(publicKeyType).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", keyName))
	}

	d.SetId(sak.Name)

	if err := d.Set("name", sak.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("key_algorithm", sak.KeyAlgorithm); err != nil {
		return fmt.Errorf("Error setting key_algorithm: %s", err)
	}
	if err := d.Set("public_key", sak.PublicKeyData); err != nil {
		return fmt.Errorf("Error setting public_key: %s", err)
	}

	return nil
}
