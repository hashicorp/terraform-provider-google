package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
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
		Update: resourceKmsCryptoKeyUpdate,
		Delete: resourceKmsCryptoKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_ring": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: kmsCryptoKeyRingsEquivalent,
			},
			"rotation_period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: orEmpty(validateKmsCryptoKeyRotationPeriod),
			},
			"version_template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protection_level": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "SOFTWARE",
							ValidateFunc: validation.StringInSlice([]string{"SOFTWARE", "HSM", ""}, false),
						},
					},
				},
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func kmsCryptoKeyRingsEquivalent(k, old, new string, d *schema.ResourceData) bool {
	keyRingIdWithSpecifiersRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-])+/keyRings/([a-zA-Z0-9_-]{1,63})$")
	normalizedKeyRingIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	if matches := keyRingIdWithSpecifiersRegex.FindStringSubmatch(new); matches != nil {
		normMatches := normalizedKeyRingIdRegex.FindStringSubmatch(old)
		return normMatches != nil && normMatches[1] == matches[1] && normMatches[2] == matches[2] && normMatches[3] == matches[3]
	}
	return false
}

type kmsCryptoKeyId struct {
	KeyRingId kmsKeyRingId
	Name      string
}

func (s *kmsCryptoKeyId) cryptoKeyId() string {
	return fmt.Sprintf("%s/cryptoKeys/%s", s.KeyRingId.keyRingId(), s.Name)
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

	key := cloudkms.CryptoKey{
		Purpose:         "ENCRYPT_DECRYPT",
		VersionTemplate: expandVersionTemplate(d.Get("version_template").([]interface{})),
	}

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

	d.SetId(cryptoKeyId.cryptoKeyId())

	return resourceKmsCryptoKeyRead(d, meta)
}

func resourceKmsCryptoKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	key := cloudkms.CryptoKey{}

	if d.HasChange("rotation_period") && d.Get("rotation_period") != "" {
		rotationPeriod := d.Get("rotation_period").(string)
		nextRotation, err := kmsCryptoKeyNextRotation(time.Now(), rotationPeriod)

		if err != nil {
			return fmt.Errorf("Error setting CryptoKey rotation period: %s", err.Error())
		}

		key.NextRotationTime = nextRotation
		key.RotationPeriod = rotationPeriod
	}

	if d.HasChange("version_template") {
		key.VersionTemplate = expandVersionTemplate(d.Get("version_template").([]interface{}))
	}

	cryptoKey, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Patch(cryptoKeyId.cryptoKeyId(), &key).UpdateMask("rotation_period,next_rotation_time").Do()

	if err != nil {
		return fmt.Errorf("Error updating CryptoKey: %s", err.Error())
	}

	log.Printf("[DEBUG] Updated CryptoKey %s", cryptoKey.Name)

	d.SetId(cryptoKeyId.cryptoKeyId())

	return resourceKmsCryptoKeyRead(d, meta)
}

func resourceKmsCryptoKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Executing read for KMS CryptoKey %s", cryptoKeyId.cryptoKeyId())

	cryptoKey, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Get(cryptoKeyId.cryptoKeyId()).Do()
	if err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	d.Set("key_ring", cryptoKeyId.KeyRingId.terraformId())
	d.Set("name", cryptoKeyId.Name)
	d.Set("rotation_period", cryptoKey.RotationPeriod)
	d.Set("self_link", cryptoKey.Name)

	if err = d.Set("version_template", flattenVersionTemplate(cryptoKey.VersionTemplate)); err != nil {
		return fmt.Errorf("Error setting version_template in state: %s", err.Error())
	}

	d.SetId(cryptoKeyId.cryptoKeyId())

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

func disableCryptoKeyRotation(cryptoKeyId *kmsCryptoKeyId, config *Config) error {
	keyClient := config.clientKms.Projects.Locations.KeyRings.CryptoKeys
	_, err := keyClient.Patch(cryptoKeyId.cryptoKeyId(), &cloudkms.CryptoKey{
		NullFields: []string{"rotationPeriod", "nextRotationTime"},
	}).
		UpdateMask("rotationPeriod,nextRotationTime").Do()

	return err
}

// Because KMS CryptoKey keys cannot be deleted (in GCP proper), we "delete"
// the key ring by
// a) marking all key versions for destruction (24hr soft-delete)
// b) disabling rotation of the key
// c) removing it from state
// This disables all usage of previous versions of the key and makes it
// generally useless for encryption and decryption of data.
// Re-creation of this resource through Terraform will produce an error.
func resourceKmsCryptoKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf(`
[WARNING] KMS CryptoKey resources cannot be deleted from GCP. The CryptoKey %s will be removed from Terraform state,
and all its CryptoKeyVersions will be destroyed, but it will still be present on the server.`, cryptoKeyId.cryptoKeyId())

	// Delete all versions of the key
	if err := clearCryptoKeyVersions(cryptoKeyId, config); err != nil {
		return err
	}

	// Make sure automatic key rotation is disabled.
	if err := disableCryptoKeyRotation(cryptoKeyId, config); err != nil {
		return fmt.Errorf(
			"While cryptoKeyVersions were cleared, Terraform was unable to disable automatic rotation of key due to an error: %s."+
				"Please retry or manually disable automatic rotation to prevent creation of a new version of this key.", err)
	}

	d.SetId("")
	return nil
}

func expandVersionTemplate(configured []interface{}) *cloudkms.CryptoKeyVersionTemplate {
	if configured == nil || len(configured) == 0 {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &cloudkms.CryptoKeyVersionTemplate{
		Algorithm:       data["algorithm"].(string),
		ProtectionLevel: data["protection_level"].(string),
	}
}

func flattenVersionTemplate(versionTemplate *cloudkms.CryptoKeyVersionTemplate) []map[string]interface{} {
	if versionTemplate == nil {
		return nil
	}

	versionTemplateSchema := make([]map[string]interface{}, 0, 1)
	data := map[string]interface{}{
		"algorithm":        versionTemplate.Algorithm,
		"protection_level": versionTemplate.ProtectionLevel,
	}

	versionTemplateSchema = append(versionTemplateSchema, data)
	return versionTemplateSchema
}

func validateKmsCryptoKeyRotationPeriod(value interface{}, _ string) (ws []string, errors []error) {
	period := value.(string)
	pattern := regexp.MustCompile("^([0-9.]*\\d)s$")
	match := pattern.FindStringSubmatch(period)

	if len(match) == 0 {
		errors = append(errors, fmt.Errorf("Invalid rotation period format: %s", period))
		// Cannot continue to validate because we cannot extract a number.
		return
	}

	number := match[1]
	seconds, err := strconv.ParseFloat(number, 64)

	if err != nil {
		errors = append(errors, err)
	} else {
		if seconds < 86400.0 {
			errors = append(errors, fmt.Errorf("Rotation period must be greater than one day"))
		}

		parts := strings.Split(number, ".")

		if len(parts) > 1 && len(parts[1]) > 9 {
			errors = append(errors, fmt.Errorf("Rotation period cannot have more than 9 fractional digits"))
		}
	}

	return
}

func kmsCryptoKeyNextRotation(now time.Time, period string) (result string, err error) {
	var duration time.Duration

	duration, err = time.ParseDuration(period)

	if err == nil {
		result = now.UTC().Add(duration).Format(time.RFC3339Nano)
	}

	return
}

func parseKmsCryptoKeyId(id string, config *Config) (*kmsCryptoKeyId, error) {
	parts := strings.Split(id, "/")

	cryptoKeyIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})/([a-zA-Z0-9_-]{1,63})$")
	cryptoKeyIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})/([a-zA-Z0-9_-]{1,63})$")
	cryptoKeyRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/keyRings/([a-zA-Z0-9_-]{1,63})/cryptoKeys/([a-zA-Z0-9_-]{1,63})$")

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

	if parts := cryptoKeyRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &kmsCryptoKeyId{
			KeyRingId: kmsKeyRingId{
				Project:  parts[1],
				Location: parts[2],
				Name:     parts[3],
			},
			Name: parts[4],
		}, nil
	}
	return nil, fmt.Errorf("Invalid CryptoKey id format, expecting `{projectId}/{locationId}/{KeyringName}/{cryptoKeyName}` or `{locationId}/{keyRingName}/{cryptoKeyName}.`")
}
