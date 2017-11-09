package google

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudkms/v1"
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
			"rotation_period": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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

	key := cloudkms.CryptoKey{Purpose: "ENCRYPT_DECRYPT"}

	if d.Get("rotation_period") != "" {
		rotationPeriod := d.Get("rotation_period").(string)
		nextRotation, err := kmsCryptoKeyNextRotation(time.Now(), rotationPeriod)

		if err != nil {
			return fmt.Errorf("Error setting CryptoKey rotation period: %s", err.Error())
		}

		key.NextRotationTime = nextRotation
		key.RotationPeriod = rotationPeriod
	}

	cryptoKey, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Create(cryptoKeyId.KeyRingId.keyRingId(), &key).CryptoKeyId(cryptoKeyId.Name).Do()

	if err != nil {
		return fmt.Errorf("Error creating CryptoKey: %s", err.Error())
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

	log.Printf(`
[WARNING] KMS CryptoKey resources cannot be deleted from GCP. The CryptoKey %s will be removed from Terraform state,
and all its CryptoKeyVersions will be destroyed, but it will still be present on the server.`, cryptoKeyId.cryptoKeyId())

	err = clearCryptoKeyVersions(cryptoKeyId, config)

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func validateKmsCryptoKeyRotationPeriod(period string) error {
	pattern := regexp.MustCompile("^([0-9.]*\\d)s$")
	match := pattern.FindStringSubmatch(period)

	if len(match) == 0 {
		return fmt.Errorf("Invalid period format: %s", period)
	}

	number := match[1]
	seconds, err := strconv.ParseFloat(number, 64)

	if err == nil && seconds < 86400.0 {
		return fmt.Errorf("Rotation period must be greater than one day")
	}

	parts := strings.Split(number, ".")

	if err == nil && len(parts) > 1 && len(parts[1]) > 9 {
		return fmt.Errorf("Rotation period cannot have more than 9 fractional digits")
	}

	return nil
}

func kmsCryptoKeyNextRotation(now time.Time, period string) (string, error) {
	var result string
	var duration time.Duration

	err := validateKmsCryptoKeyRotationPeriod(period)

	if err == nil {
		duration, err = time.ParseDuration(period)
	}

	if err == nil {
		result = now.UTC().Add(duration).Format(time.RFC3339Nano)
	}

	return result, err
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

func resourceKmsCryptoKeyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return nil, err
	}

	d.Set("key_ring", cryptoKeyId.KeyRingId.terraformId())
	d.Set("name", cryptoKeyId.Name)

	d.SetId(cryptoKeyId.terraformId())

	return []*schema.ResourceData{d}, nil
}
