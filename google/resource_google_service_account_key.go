package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/iam/v1"
)

func resourceGoogleServiceAccountKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleServiceAccountKeyCreate,
		Read:   resourceGoogleServiceAccountKeyRead,
		Delete: resourceGoogleServiceAccountKeyDelete,
		Schema: map[string]*schema.Schema{
			// Required
			"service_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the parent service account of the key. This can be a string in the format {ACCOUNT} or projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}, where {ACCOUNT} is the email address or unique id of the service account. If the {ACCOUNT} syntax is used, the project will be inferred from the account.`,
			},
			// Optional
			"key_algorithm": {
				Type:         schema.TypeString,
				Default:      "KEY_ALG_RSA_2048",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"KEY_ALG_UNSPECIFIED", "KEY_ALG_RSA_1024", "KEY_ALG_RSA_2048"}, false),
				Description:  `The algorithm used to generate the key, used only on create. KEY_ALG_RSA_2048 is the default algorithm. Valid values are: "KEY_ALG_RSA_1024", "KEY_ALG_RSA_2048".`,
			},
			"pgp_key": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "The pgp_key field has been removed. See https://www.terraform.io/docs/extend/best-practices/sensitive-state.html for more information.",
				Computed: true,
			},
			"private_key_type": {
				Type:         schema.TypeString,
				Default:      "TYPE_GOOGLE_CREDENTIALS_FILE",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TYPE_UNSPECIFIED", "TYPE_PKCS12_FILE", "TYPE_GOOGLE_CREDENTIALS_FILE"}, false),
			},
			"public_key_type": {
				Type:         schema.TypeString,
				Default:      "TYPE_X509_PEM_FILE",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TYPE_NONE", "TYPE_X509_PEM_FILE", "TYPE_RAW_PUBLIC_KEY"}, false),
			},
			// Computed
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: `The name used for this key pair`,
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: `The public key, base64 encoded`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The private key in JSON format, base64 encoded. This is what you normally get as a file when creating service account keys through the CLI or web console. This is only populated when creating a new key.`,
			},
			"valid_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key can be used after this timestamp. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".`,
			},
			"valid_before": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key can be used before this timestamp. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".`,
			},
			"private_key_encrypted": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "The private_key_encrypted field has been removed. See https://www.terraform.io/docs/extend/best-practices/sensitive-state.html for more information.",
			},
			"private_key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "The private_key_fingerprint field has been removed. See https://www.terraform.io/docs/extend/best-practices/sensitive-state.html for more information.",
			},
		},
	}
}

func resourceGoogleServiceAccountKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceAccountName, err := serviceAccountFQN(d.Get("service_account_id").(string), d, config)
	if err != nil {
		return err
	}

	r := &iam.CreateServiceAccountKeyRequest{
		KeyAlgorithm:   d.Get("key_algorithm").(string),
		PrivateKeyType: d.Get("private_key_type").(string),
	}

	sak, err := config.clientIAM.Projects.ServiceAccounts.Keys.Create(serviceAccountName, r).Do()
	if err != nil {
		return fmt.Errorf("Error creating service account key: %s", err)
	}

	d.SetId(sak.Name)
	// Data only available on create.
	d.Set("valid_after", sak.ValidAfterTime)
	d.Set("valid_before", sak.ValidBeforeTime)
	d.Set("private_key", sak.PrivateKeyData)

	err = serviceAccountKeyWaitTime(config.clientIAM.Projects.ServiceAccounts.Keys, d.Id(), d.Get("public_key_type").(string), "Creating Service account key", 4*time.Minute)
	if err != nil {
		return err
	}
	return resourceGoogleServiceAccountKeyRead(d, meta)
}

func resourceGoogleServiceAccountKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	publicKeyType := d.Get("public_key_type").(string)

	// Confirm the service account key exists
	sak, err := config.clientIAM.Projects.ServiceAccounts.Keys.Get(d.Id()).PublicKeyType(publicKeyType).Do()
	if err != nil {
		if err = handleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", d.Id())); err == nil {
			return nil
		} else {
			// This resource also returns 403 when it's not found.
			if isGoogleApiErrorWithCode(err, 403) {
				log.Printf("[DEBUG] Got a 403 error trying to read service account key %s, assuming it's gone.", d.Id())
				d.SetId("")
				return nil
			} else {
				return err
			}
		}
	}

	d.Set("name", sak.Name)
	d.Set("key_algorithm", sak.KeyAlgorithm)
	d.Set("public_key", sak.PublicKeyData)
	return nil
}

func resourceGoogleServiceAccountKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Projects.ServiceAccounts.Keys.Delete(d.Id()).Do()

	if err != nil {
		if err = handleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", d.Id())); err == nil {
			return nil
		} else {
			// This resource also returns 403 when it's not found.
			if isGoogleApiErrorWithCode(err, 403) {
				log.Printf("[DEBUG] Got a 403 error trying to read service account key %s, assuming it's gone.", d.Id())
				d.SetId("")
				return nil
			} else {
				return err
			}
		}
	}

	d.SetId("")
	return nil
}
