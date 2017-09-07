package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/encryption"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/iam/v1"
)

func resourceGoogleServiceAccountKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleServiceAccountKeyCreate,
		Read:   resourceGoogleServiceAccountKeyRead,
		Delete: resourceGoogleServiceAccountKeyDelete,
		Schema: map[string]*schema.Schema{
			"service_account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "KEY_ALG_RSA_2048",
				Optional: true,
				ForceNew: true,
			},
			"private_key": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key_type": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "TYPE_GOOGLE_CREDENTIALS_FILE",
				Optional: true,
				ForceNew: true,
			},
			"valid_after": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_before": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"pgp_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

	serviceAccount := d.Get("service_account_id").(string)

	r := &iam.CreateServiceAccountKeyRequest{}

	if v, ok := d.GetOk("key_algorithm"); ok {
		r.KeyAlgorithm = v.(string)
	}

	if v, ok := d.GetOk("private_key_type"); ok {
		r.PrivateKeyType = v.(string)
	}

	sak, err := config.clientIAM.Projects.ServiceAccounts.Keys.Create(serviceAccount, r).Do()
	if err != nil {
		return fmt.Errorf("Error creating service account key: %s", err)
	}

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

	d.SetId(sak.Name)
	d.Set("name", sak.Name)
	d.Set("key_algorithm", sak.KeyAlgorithm)
	d.Set("private_key_type", sak.PrivateKeyType)
	d.Set("valid_after", sak.ValidAfterTime)
	d.Set("valid_before", sak.ValidBeforeTime)

	return nil
}

func resourceGoogleServiceAccountKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Confirm the service account key exists
	sak, err := config.clientIAM.Projects.ServiceAccounts.Keys.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", d.Id()))
	}

	d.SetId(sak.Name)
	d.Set("name", sak.Name)
	d.Set("key_algorithm", sak.KeyAlgorithm)
	d.Set("valid_after", sak.ValidAfterTime)
	d.Set("valid_before", sak.ValidBeforeTime)

	return nil
}

func resourceGoogleServiceAccountKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Projects.ServiceAccounts.Keys.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
