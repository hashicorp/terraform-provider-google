package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/encryption"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Optional
			"key_algorithm": {
				Type:         schema.TypeString,
				Default:      "KEY_ALG_RSA_2048",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"KEY_ALG_UNSPECIFIED", "KEY_ALG_RSA_1024", "KEY_ALG_RSA_2048"}, false),
			},
			"pgp_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"valid_after": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_before": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key_encrypted": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
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
	if v, ok := d.GetOk("pgp_key"); ok {
		encryptionKey, err := encryption.RetrieveGPGKey(v.(string))
		if err != nil {
			return err
		}

		fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, sak.PrivateKeyData, "Google Service Account Key")
		if err != nil {
			return err
		}

		d.Set("private_key_encrypted", encrypted)
		d.Set("private_key_fingerprint", fingerprint)
	} else {
		d.Set("private_key", sak.PrivateKeyData)
	}

	err = serviceAccountKeyWaitTime(config.clientIAM.Projects.ServiceAccounts.Keys, d.Id(), d.Get("public_key_type").(string), "Creating Service account key", 4)
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
