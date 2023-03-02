package google

import (
	"google.golang.org/api/cloudkms/v1"

	"encoding/base64"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleKmsSecret() *schema.Resource {
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
			"additional_authenticated_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Get("crypto_key").(string), config)

	if err != nil {
		return err
	}

	ciphertext := d.Get("ciphertext").(string)

	kmsDecryptRequest := &cloudkms.DecryptRequest{
		Ciphertext: ciphertext,
	}

	if aad, ok := d.GetOk("additional_authenticated_data"); ok {
		kmsDecryptRequest.AdditionalAuthenticatedData = aad.(string)
	}

	decryptResponse, err := config.NewKmsClient(userAgent).Projects.Locations.KeyRings.CryptoKeys.Decrypt(cryptoKeyId.cryptoKeyId(), kmsDecryptRequest).Do()

	if err != nil {
		return fmt.Errorf("Error decrypting ciphertext: %s", err)
	}

	plaintext, err := base64.StdEncoding.DecodeString(decryptResponse.Plaintext)

	if err != nil {
		return fmt.Errorf("Error decoding base64 response: %s", err)
	}

	log.Printf("[INFO] Successfully decrypted ciphertext: %s", ciphertext)

	if err := d.Set("plaintext", string(plaintext[:])); err != nil {
		return fmt.Errorf("Error setting plaintext: %s", err)
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("crypto_key").(string), ciphertext))

	return nil
}
