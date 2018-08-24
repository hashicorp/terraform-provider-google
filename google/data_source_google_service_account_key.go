package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"regexp"
)

func dataSourceGoogleServiceAccountKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
			"service_account_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Deprecated:    "Please use name to specify full service account key path projects/{project}/serviceAccounts/{serviceAccount}/keys/{keyId}",
			},
		},
	}
}

func dataSourceGoogleServiceAccountKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyName, err := getDataSourceServiceAccountKeyName(d)
	if err != nil {
		return err
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

func getDataSourceServiceAccountKeyName(d *schema.ResourceData) (string, error) {
	keyName := d.Get("name").(string)
	keyFromSAId := d.Get("service_account_id").(string)

	// Neither name nor service_account_id specified
	if keyName == "" && keyFromSAId == "" {
		return "", fmt.Errorf("please use name to specify service account key being added as this data source")
	}

	fullKeyName := keyName
	if fullKeyName == "" {
		// Key name specified as incorrectly named, deprecated service account ID field
		fullKeyName = keyFromSAId
	}

	// Validate name since interpolated values (i.e from a key or service
	// account resource) will not get validated at plan time.
	r := regexp.MustCompile(ServiceAccountKeyNameRegex)
	if r.MatchString(fullKeyName) {
		return fullKeyName, nil
	}

	return "", fmt.Errorf("invalid key name %q does not match regexp %q", fullKeyName, ServiceAccountKeyNameRegex)
}
