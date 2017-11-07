package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudkms/v1"
	"log"
	"regexp"
	"strings"
)

func resourceKmsCryptoKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsCryptoKeyCreate,
		Read:   resourceKmsCryptoKeyRead,
		Delete: resourceKmsCryptoKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKmsCryptoKeyImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_ring": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceKmsCryptoKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyRingId, err := parseKmsKeyRingId(d.Get("key_ring").(string), config)

	if err != nil {
		return err
	}

	cryptoKeyId := &kmsCryptoKeyId{
		KeyRingId: *keyRingId,
		Name:      d.Get("name").(string),
	}

	cryptoKey, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Create(cryptoKeyId.KeyRingId.keyRingId(), &cloudkms.CryptoKey{Purpose: "ENCRYPT_DECRYPT"}).CryptoKeyId(cryptoKeyId.Name).Do()

	if err != nil {
		return fmt.Errorf("Error creating CryptoKey: %s", err)
	}

	log.Printf("[DEBUG] Created CryptoKey %s", cryptoKey.Name)

	d.SetId(cryptoKeyId.terraformId())

	return resourceKmsCryptoKeyRead(d, meta)
}

func resourceKmsCryptoKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Executing read for KMS CryptoKey %s", cryptoKeyId.cryptoKeyId())

	_, err = config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Get(cryptoKeyId.cryptoKeyId()).Do()

	if err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}

	return nil
}

func clearCryptoKeyVersions(cryptoKeyId *kmsCryptoKeyId, config *Config) error {
	versionsClient := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions

	versionsResponse, err := versionsClient.List(cryptoKeyId.cryptoKeyId()).Do()

	if err != nil {
		return err
	}

	for _, version := range versionsResponse.CryptoKeyVersions {
		request := &cloudkms.DestroyCryptoKeyVersionRequest{}
		_, err = versionsClient.Destroy(version.Name, request).Do()

		if err != nil {
			return err
		}
	}

	return nil
}

/*
	Because KMS CryptoKey resources cannot be deleted on GCP, we are only going to remove it from state
    and destroy all its versions, rendering the key useless for encryption and decryption of data.
    Re-creation of this resource through Terraform will produce an error.
*/

func resourceKmsCryptoKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf("[WARNING] KMS CryptoKey resources cannot be deleted from GCP. This CryptoKey %s will be removed from Terraform state, but will still be present on the server.", cryptoKeyId.cryptoKeyId())

	d.SetId("")

	err = clearCryptoKeyVersions(cryptoKeyId, config)

	if err != nil {
		return err
	}

	return nil
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
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{cryptoKeyName}` id format.")
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

	return nil, fmt.Errorf("Invalid CryptoKey id format, expecting `{projectId}/{locationId}/{KeyringName}/{cryptoKeyName}` or `{locationId}/{KeyringName}/{cryptoKeyName}.`")
}

func resourceKmsCryptoKeyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return nil, err
	}

	d.Set("key_ring", cryptoKeyId.KeyRingId.keyRingId())

	if d.Get("purpose") == "" {
		d.Set("purpose", "ENCRYPT_DECRYPT")
	}

	d.SetId(cryptoKeyId.terraformId())

	return []*schema.ResourceData{d}, nil
}
