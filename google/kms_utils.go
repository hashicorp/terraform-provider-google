package google

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudkms/v1"
)

type kmsKeyRingId struct {
	Project  string
	Location string
	Name     string
}

func (s *kmsKeyRingId) keyRingId() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", s.Project, s.Location, s.Name)
}

func (s *kmsKeyRingId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Location, s.Name)
}

func parseKmsKeyRingId(id string, config *Config) (*kmsKeyRingId, error) {
	parts := strings.Split(id, "/")

	keyRingIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	keyRingIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	keyRingRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/keyRings/([a-zA-Z0-9_-]{1,63})$")

	if keyRingIdRegex.MatchString(id) {
		return &kmsKeyRingId{
			Project:  parts[0],
			Location: parts[1],
			Name:     parts[2],
		}, nil
	}

	if keyRingIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{keyRingName}` id format.")
		}

		return &kmsKeyRingId{
			Project:  config.Project,
			Location: parts[0],
			Name:     parts[1],
		}, nil
	}

	if parts := keyRingRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &kmsKeyRingId{
			Project:  parts[1],
			Location: parts[2],
			Name:     parts[3],
		}, nil
	}
	return nil, fmt.Errorf("Invalid KeyRing id format, expecting `{projectId}/{locationId}/{keyRingName}` or `{locationId}/{keyRingName}.`")
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
	return nil, fmt.Errorf("Invalid CryptoKey id format, expecting `{projectId}/{locationId}/{KeyringName}/{cryptoKeyName}` or `{locationId}/{keyRingName}/{cryptoKeyName}, got id: %s`", id)
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
