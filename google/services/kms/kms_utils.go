// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudkms/v1"
)

type KmsKeyRingId struct {
	Project  string
	Location string
	Name     string
}

func (s *KmsKeyRingId) KeyRingId() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", s.Project, s.Location, s.Name)
}

func (s *KmsKeyRingId) TerraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Location, s.Name)
}

func parseKmsKeyRingId(id string, config *transport_tpg.Config) (*KmsKeyRingId, error) {
	parts := strings.Split(id, "/")

	KeyRingIdRegex := regexp.MustCompile("^(" + verify.ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	KeyRingIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	keyRingRelativeLinkRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-]+)/keyRings/([a-zA-Z0-9_-]{1,63})$")

	if KeyRingIdRegex.MatchString(id) {
		return &KmsKeyRingId{
			Project:  parts[0],
			Location: parts[1],
			Name:     parts[2],
		}, nil
	}

	if KeyRingIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{keyRingName}` id format.")
		}

		return &KmsKeyRingId{
			Project:  config.Project,
			Location: parts[0],
			Name:     parts[1],
		}, nil
	}

	if parts := keyRingRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &KmsKeyRingId{
			Project:  parts[1],
			Location: parts[2],
			Name:     parts[3],
		}, nil
	}
	return nil, fmt.Errorf("Invalid KeyRing id format, expecting `{projectId}/{locationId}/{keyRingName}` or `{locationId}/{keyRingName}.`")
}

func kmsCryptoKeyRingsEquivalent(k, old, new string, d *schema.ResourceData) bool {
	KeyRingIdWithSpecifiersRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-])+/keyRings/([a-zA-Z0-9_-]{1,63})$")
	normalizedKeyRingIdRegex := regexp.MustCompile("^(" + verify.ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	if matches := KeyRingIdWithSpecifiersRegex.FindStringSubmatch(new); matches != nil {
		normMatches := normalizedKeyRingIdRegex.FindStringSubmatch(old)
		return normMatches != nil && normMatches[1] == matches[1] && normMatches[2] == matches[2] && normMatches[3] == matches[3]
	}
	return false
}

type KmsCryptoKeyId struct {
	KeyRingId KmsKeyRingId
	Name      string
}

func (s *KmsCryptoKeyId) CryptoKeyId() string {
	return fmt.Sprintf("%s/cryptoKeys/%s", s.KeyRingId.KeyRingId(), s.Name)
}

func (s *KmsCryptoKeyId) TerraformId() string {
	return fmt.Sprintf("%s/%s", s.KeyRingId.TerraformId(), s.Name)
}

type kmsCryptoKeyVersionId struct {
	CryptoKeyId KmsCryptoKeyId
	Name        string
}

func (s *kmsCryptoKeyVersionId) cryptoKeyVersionId() string {
	return fmt.Sprintf(s.Name)
}

func (s *kmsCryptoKeyVersionId) TerraformId() string {
	return fmt.Sprintf("%s/%s", s.CryptoKeyId.TerraformId(), s.Name)
}

func validateKmsCryptoKeyRotationPeriod(value interface{}, _ string) (ws []string, errors []error) {
	period := value.(string)
	pattern := regexp.MustCompile(`^([0-9.]*\d)s$`)
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

func ParseKmsCryptoKeyId(id string, config *transport_tpg.Config) (*KmsCryptoKeyId, error) {
	parts := strings.Split(id, "/")

	cryptoKeyIdRegex := regexp.MustCompile("^(" + verify.ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})/([a-zA-Z0-9_-]{1,63})$")
	cryptoKeyIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})/([a-zA-Z0-9_-]{1,63})$")
	cryptoKeyRelativeLinkRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-]+)/keyRings/([a-zA-Z0-9_-]{1,63})/cryptoKeys/([a-zA-Z0-9_-]{1,63})$")

	if cryptoKeyIdRegex.MatchString(id) {
		return &KmsCryptoKeyId{
			KeyRingId: KmsKeyRingId{
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

		return &KmsCryptoKeyId{
			KeyRingId: KmsKeyRingId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := cryptoKeyRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &KmsCryptoKeyId{
			KeyRingId: KmsKeyRingId{
				Project:  parts[1],
				Location: parts[2],
				Name:     parts[3],
			},
			Name: parts[4],
		}, nil
	}

	return nil, fmt.Errorf("Invalid CryptoKey id format, expecting `{projectId}/{locationId}/{KeyringName}/{cryptoKeyName}` or `{locationId}/{keyRingName}/{cryptoKeyName}, got id: %s`", id)
}
func parseKmsCryptoKeyVersionId(id string, config *transport_tpg.Config) (*kmsCryptoKeyVersionId, error) {
	cryptoKeyVersionRelativeLinkRegex := regexp.MustCompile("^projects/(" + verify.ProjectRegex + ")/locations/([a-z0-9-]+)/keyRings/([a-zA-Z0-9_-]{1,63})/cryptoKeys/([a-zA-Z0-9_-]{1,63})/cryptoKeyVersions/([a-zA-Z0-9_-]{1,63})$")

	if parts := cryptoKeyVersionRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &kmsCryptoKeyVersionId{
			CryptoKeyId: KmsCryptoKeyId{
				KeyRingId: KmsKeyRingId{
					Project:  parts[1],
					Location: parts[2],
					Name:     parts[3],
				},
				Name: parts[4],
			},
			Name: "projects/" + parts[1] + "/locations/" + parts[2] + "/keyRings/" + parts[3] + "/cryptoKeys/" + parts[4] + "/cryptoKeyVersions/" + parts[5],
		}, nil
	}
	return nil, fmt.Errorf("Invalid CryptoKeyVersion id format, expecting `{projectId}/{locationId}/{KeyringName}/{cryptoKeyName}/{cryptoKeyVersion}` or `{locationId}/{keyRingName}/{cryptoKeyName}/{cryptoKeyVersion}, got id: %s`", id)
}

func clearCryptoKeyVersions(cryptoKeyId *KmsCryptoKeyId, userAgent string, config *transport_tpg.Config) error {
	versionsClient := config.NewKmsClient(userAgent).Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions

	listCall := versionsClient.List(cryptoKeyId.CryptoKeyId())
	if config.UserProjectOverride {
		listCall.Header().Set("X-Goog-User-Project", cryptoKeyId.KeyRingId.Project)
	}
	versionsResponse, err := listCall.Do()

	if err != nil {
		return err
	}

	for _, version := range versionsResponse.CryptoKeyVersions {
		// skip the versions that have been destroyed earlier
		if version.State != "DESTROYED" && version.State != "DESTROY_SCHEDULED" {
			request := &cloudkms.DestroyCryptoKeyVersionRequest{}
			destroyCall := versionsClient.Destroy(version.Name, request)
			if config.UserProjectOverride {
				destroyCall.Header().Set("X-Goog-User-Project", cryptoKeyId.KeyRingId.Project)
			}
			_, err = destroyCall.Do()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteCryptoKeyVersions(cryptoKeyVersionId *kmsCryptoKeyVersionId, d *schema.ResourceData, userAgent string, config *transport_tpg.Config) error {
	versionsClient := config.NewKmsClient(userAgent).Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions
	request := &cloudkms.DestroyCryptoKeyVersionRequest{}
	destroyCall := versionsClient.Destroy(cryptoKeyVersionId.Name, request)
	if config.UserProjectOverride {
		destroyCall.Header().Set("X-Goog-User-Project", cryptoKeyVersionId.CryptoKeyId.KeyRingId.Project)
	}
	_, err := destroyCall.Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ID %s", cryptoKeyVersionId.Name))
	}

	return nil
}

func disableCryptoKeyRotation(cryptoKeyId *KmsCryptoKeyId, userAgent string, config *transport_tpg.Config) error {
	keyClient := config.NewKmsClient(userAgent).Projects.Locations.KeyRings.CryptoKeys
	patchCall := keyClient.Patch(cryptoKeyId.CryptoKeyId(), &cloudkms.CryptoKey{
		NullFields: []string{"rotationPeriod", "nextRotationTime"},
	}).
		UpdateMask("rotationPeriod,nextRotationTime")
	if config.UserProjectOverride {
		patchCall.Header().Set("X-Goog-User-Project", cryptoKeyId.KeyRingId.Project)
	}
	_, err := patchCall.Do()

	return err
}
