package google

import (
	"google.golang.org/api/cloudkms/v1"

	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
	"strings"
	"encoding/base64"
	"time"
)

func dataSourceGoogleKmsSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsSecretRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"plaintext": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

type kmsCryptoKeyId struct {
	KeyRingId kmsKeyRingId
	Name      string
}

// TODO: Add the info about rotation frequency and start time.

func (s *kmsCryptoKeyId) cryptoKeyId() string {
	return fmt.Sprintf("%s/cryptoKeys/%s", s.KeyRingId.keyRingId(), s.Name)
}

func (s *kmsCryptoKeyId) parentId() string {
	return s.KeyRingId.keyRingId()
}

func (s *kmsCryptoKeyId) terraformId() string {
	return fmt.Sprintf("%s/%s", s.KeyRingId.terraformId(), s.Name)
}

func parseKmsCryptoKeyId(id string, config *Config) (*kmsCryptoKeyId, error) {
	parts := strings.Split(id, "/")

	cryptoKeyIdRegex := regexp.MustCompile("^([a-z0-9-]+)/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})+/([a-zA-Z0-9_-]{1,63})$")
	cryptoKeyIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})+/([a-zA-Z0-9_-]{1,63})$")

	if cryptoKeyIdRegex.MatchString(id) {
		return &kmsCryptoKeyId{
			KeyRingId: kmsKeyRingId{
				Project:  parts[0],
				Location: parts[1],
				Name:     parts[2],
			},
			Name: parts[3],
		}, nil
	}

	if cryptoKeyIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{keyRingName}/{cryptoKeyName}` id format.")
		}

		return &kmsCryptoKeyId{
			KeyRingId: kmsKeyRingId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	return nil, fmt.Errorf("Invalid CryptoKey id format, expecting `{projectId}/{locationId}/{KeyringName}/{cryptoKeyName}` or `{locationId}/{keyRingName}/{cryptoKeyName}.`")
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

	d.Set("plaintext", fmt.Sprintf("%s", plaintext))
	d.SetId(time.Now().UTC().String())

	return nil
}
