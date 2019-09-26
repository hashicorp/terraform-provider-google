package google

import (
	"google.golang.org/api/cloudkms/v1"

	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func dataSourceGoogleKmsSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsSecretRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plaintext": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceGoogleKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Get("crypto_key").(string), config)

	if err != nil {
		return err
	}

	ciphertext := d.Get("ciphertext").(string)

	kmsDecryptRequest := &cloudkms.DecryptRequest{
		Ciphertext: ciphertext,
	}

	decryptResponse, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Decrypt(cryptoKeyId.cryptoKeyId(), kmsDecryptRequest).Do()

	if err != nil {
		return fmt.Errorf("Error decrypting ciphertext: %s", err)
	}

	plaintext, err := base64.StdEncoding.DecodeString(decryptResponse.Plaintext)

	if err != nil {
		return fmt.Errorf("Error decoding base64 response: %s", err)
	}

	log.Printf("[INFO] Successfully decrypted ciphertext: %s", ciphertext)

	d.Set("plaintext", string(plaintext[:]))
	d.SetId(time.Now().UTC().String())

	return nil
}
