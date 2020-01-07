package google

import (
	"google.golang.org/api/cloudkms/v1"

	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleKmsSecretCiphertext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsSecretCiphertextRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plaintext": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceGoogleKmsSecretCiphertextRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Get("crypto_key").(string), config)

	if err != nil {
		return err
	}

	plaintext := base64.StdEncoding.EncodeToString([]byte(d.Get("plaintext").(string)))

	kmsEncryptRequest := &cloudkms.EncryptRequest{
		Plaintext: plaintext,
	}

	encryptCall := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Encrypt(cryptoKeyId.cryptoKeyId(), kmsEncryptRequest)
	if config.UserProjectOverride {
		encryptCall.Header().Set("X-Goog-User-Project", cryptoKeyId.KeyRingId.Project)
	}
	encryptResponse, err := encryptCall.Do()

	if err != nil {
		return fmt.Errorf("Error encrypting plaintext: %s", err)
	}

	log.Printf("[INFO] Successfully encrypted plaintext")

	d.Set("ciphertext", encryptResponse.Ciphertext)
	d.SetId(time.Now().UTC().String())

	return nil
}
