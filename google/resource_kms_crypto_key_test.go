package google

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCryptoKeyIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedCryptoKeyId string
		Config              *Config
	}{
		"id is in project/location/keyRingName/cryptoKeyName format": {
			ImportId:            "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
		},
		"id is in domain:project/location/keyRingName/cryptoKeyName format": {
			ImportId:            "example.com:test-project/us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "example.com:test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/example.com:test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
		},
		"id contains name that is longer than 63 characters": {
			ImportId:      "test-project/us-central1/test-key-ring/can-you-believe-that-this-cryptokey-name-is-this-extravagantly-long",
			ExpectedError: true,
		},
		"id is in location/keyRingName/cryptoKeyName format": {
			ImportId:            "us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
			Config:              &Config{Project: "test-project"},
		},
		"id is in location/keyRingName/cryptoKeyName format without project in config": {
			ImportId:      "us-central1/test-key-ring/test-key-name",
			ExpectedError: true,
			Config:        &Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		cryptoKeyId, err := parseKmsCryptoKeyId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if cryptoKeyId.terraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, cryptoKeyId.terraformId())
		}

		if cryptoKeyId.cryptoKeyId() != tc.ExpectedCryptoKeyId {
			t.Fatalf("bad: %s, expected CryptoKey ID to be `%s` but is `%s`", tn, tc.ExpectedCryptoKeyId, cryptoKeyId.cryptoKeyId())
		}
	}
}

func TestCryptoKeyNextRotationCalculation(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	period, _ := time.ParseDuration("1000000s")

	expected := now.Add(period).Format(time.RFC3339Nano)

	timestamp, err := kmsCryptoKeyNextRotation(now, "1000000s")

	if err != nil {
		t.Fatalf("unexpected failure parsing time %s and duration 1000s: %s", now, err.Error())
	}

	if expected != timestamp {
		t.Fatalf("expected %s to equal %s", timestamp, expected)
	}
}

func TestCryptoKeyNextRotationCalculation_validation(t *testing.T) {
	t.Parallel()

	_, errs := validateKmsCryptoKeyRotationPeriod("86399s", "rotation_period")

	if len(errs) == 0 {
		t.Fatalf("Periods of less than a day should be invalid")
	}

	_, errs = validateKmsCryptoKeyRotationPeriod("100000.0000000001s", "rotation_period")

	if len(errs) == 0 {
		t.Fatalf("Numbers with more than 9 fractional digits are invalid")
	}
}

func TestCryptoKeyStateUpgradeV0(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Attributes map[string]interface{}
		Expected   map[string]string
		Meta       interface{}
	}{
		"change key_ring from terraform id fmt to link fmt": {
			Attributes: map[string]interface{}{
				"key_ring": "my-project/my-location/my-key-ring",
			},
			Expected: map[string]string{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Meta: &Config{},
		},
		"key_ring link fmt stays as link fmt": {
			Attributes: map[string]interface{}{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Expected: map[string]string{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Meta: &Config{},
		},
		"key_ring without project to link fmt": {
			Attributes: map[string]interface{}{
				"key_ring": "my-location/my-key-ring",
			},
			Expected: map[string]string{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Meta: &Config{
				Project: "my-project",
			},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := resourceKMSCryptoKeyUpgradeV0(context.Background(), tc.Attributes, tc.Meta)

			if err != nil {
				t.Error(err)
			}

			for k, v := range tc.Expected {
				if actual[k] != v {
					t.Errorf("expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
						k, v, k, actual[k], actual)
				}
			}
		})
	}
}

func TestAccKmsCryptoKey_basic(t *testing.T) {
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test importing with a short id
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s/%s/%s", projectId, location, keyRingName, cryptoKeyName),
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(t, projectId, location, keyRingName, cryptoKeyName),
					testAccCheckGoogleKmsCryptoKeyRotationDisabled(t, projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

func TestAccKmsCryptoKey_rotation(t *testing.T) {
	// when rotation is set, next rotation time is set using time.Now
	skipIfVcr(t)
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	rotationPeriod := "100000s"
	updatedRotationPeriod := "7776000s"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, rotationPeriod),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, updatedRotationPeriod),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKey_rotationRemoved(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(t, projectId, location, keyRingName, cryptoKeyName),
					testAccCheckGoogleKmsCryptoKeyRotationDisabled(t, projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

func TestAccKmsCryptoKey_template(t *testing.T) {
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	algorithm := "EC_SIGN_P256_SHA256"
	updatedAlgorithm := "EC_SIGN_P384_SHA384"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_template(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, algorithm),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKey_template(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, updatedAlgorithm),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(t, projectId, location, keyRingName, cryptoKeyName),
					testAccCheckGoogleKmsCryptoKeyRotationDisabled(t, projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

func TestAccKmsCryptoKey_destroyDuration(t *testing.T) {
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_destroyDuration(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key.crypto_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(t, projectId, location, keyRingName, cryptoKeyName),
					testAccCheckGoogleKmsCryptoKeyRotationDisabled(t, projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

func TestAccKmsCryptoKey_importOnly(t *testing.T) {
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_importOnly(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:            "google_kms_crypto_key.crypto_key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_initial_version_creation"},
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(t, projectId, location, keyRingName, cryptoKeyName),
					testAccCheckGoogleKmsCryptoKeyRotationDisabled(t, projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

// KMS KeyRings cannot be deleted. This ensures that the CryptoKey resource was removed from state,
// even though the server-side resource was not removed.
func testAccCheckGoogleKmsCryptoKeyWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}

// KMS KeyRings cannot be deleted. This ensures that the CryptoKey resource's CryptoKeyVersion
// sub-resources were scheduled to be destroyed, rendering the key itself inoperable.
func testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(t *testing.T, projectId, location, keyRingName, cryptoKeyName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		config := googleProviderConfig(t)
		gcpResourceUri := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyRingName, cryptoKeyName)

		response, err := config.NewKmsClient(config.userAgent).Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.List(gcpResourceUri).Do()

		if err != nil {
			return fmt.Errorf("Unexpected failure to list versions: %s", err)
		}

		versions := response.CryptoKeyVersions

		for _, v := range versions {
			if v.State != "DESTROY_SCHEDULED" && v.State != "DESTROYED" {
				return fmt.Errorf("CryptoKey %s should have no versions, but version %s has state %s", cryptoKeyName, v.Name, v.State)
			}
		}

		return nil
	}
}

// KMS KeyRings cannot be deleted. This ensures that the CryptoKey autorotation
// was disabled to prevent more versions of the key from being created.
func testAccCheckGoogleKmsCryptoKeyRotationDisabled(t *testing.T, projectId, location, keyRingName, cryptoKeyName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		config := googleProviderConfig(t)
		gcpResourceUri := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyRingName, cryptoKeyName)

		response, err := config.NewKmsClient(config.userAgent).Projects.Locations.KeyRings.CryptoKeys.Get(gcpResourceUri).Do()
		if err != nil {
			return fmt.Errorf("Unexpected failure while verifying 'deleted' crypto key: %s", err)
		}

		if response.NextRotationTime != "" {
			return fmt.Errorf("Expected empty nextRotationTime for 'deleted' crypto key, got %s", response.NextRotationTime)
		}
		if response.RotationPeriod != "" {
			return fmt.Errorf("Expected empty RotationPeriod for 'deleted' crypto key, got %s", response.RotationPeriod)
		}

		return nil
	}
}

func TestAccKmsCryptoKeyVersion_basic(t *testing.T) {
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKeyVersion_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key_version.crypto_key_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKeyVersion_removed(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
		},
	})
}

func TestAccKmsCryptoKeyVersion_skipInitialVersion(t *testing.T) {
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKeyVersion_skipInitialVersion(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key_version.crypto_key_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKmsCryptoKeyVersion_patch(t *testing.T) {
	t.Parallel()

	projectId := fmt.Sprintf("tf-test-%d", randInt(t))
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	state := "DISABLED"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKeyVersion_patchInitialize(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},
			{
				ResourceName:      "google_kms_crypto_key_version.crypto_key_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKeyVersion_patch("true", projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, state),
			},
			{
				ResourceName:      "google_kms_crypto_key_version.crypto_key_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsCryptoKeyVersion_patch("false", projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, state),
			},
		},
	})
}

// This test runs in its own project, otherwise the test project would start to get filled
// with undeletable resources
func testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = google_kms_key_ring.key_ring.id
  labels = {
    key = "value"
  }
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, rotationPeriod string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name            = "%s"
  key_ring        = google_kms_key_ring.key_ring.id
  rotation_period = "%s"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, rotationPeriod)
}

func testGoogleKmsCryptoKey_rotationRemoved(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = google_kms_key_ring.key_ring.id
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKey_template(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, algorithm string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = google_kms_key_ring.key_ring.id
  purpose  = "ASYMMETRIC_SIGN"

  version_template {
    algorithm = "%s"
  }
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, algorithm)
}

func testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName)
}

func testGoogleKmsCryptoKey_destroyDuration(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = google_kms_key_ring.key_ring.id
  labels = {
    key = "value"
  }
  destroy_scheduled_duration = "129600s"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKey_importOnly(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = google_project_service.acceptance.project
  name     = "%s"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = google_kms_key_ring.key_ring.id
  labels = {
    key = "value"
  }
  skip_initial_version_creation = true
  import_only = true
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKeyVersion_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "acceptance" {
	project = google_project.acceptance.project_id
	service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
	project  = google_project_service.acceptance.project
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name     = "%s"
	key_ring = google_kms_key_ring.key_ring.id
	labels = {
		key = "value"
	}
}

resource "google_kms_crypto_key_version" "crypto_key_version" {
	crypto_key = google_kms_crypto_key.crypto_key.id
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKeyVersion_removed(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "acceptance" {
	project = google_project.acceptance.project_id
	service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
	project  = google_project_service.acceptance.project
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name     = "%s"
	key_ring = google_kms_key_ring.key_ring.id
	labels = {
		key = "value"
	}
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKeyVersion_skipInitialVersion(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "acceptance" {
	project = google_project.acceptance.project_id
	service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
	project  = google_project_service.acceptance.project
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name     = "%s"
	key_ring = google_kms_key_ring.key_ring.id
	labels = {
		key = "value"
	}
	skip_initial_version_creation = true
}

resource "google_kms_crypto_key_version" "crypto_key_version" {
	crypto_key = google_kms_crypto_key.crypto_key.id
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}
func testGoogleKmsCryptoKeyVersion_patchInitialize(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "acceptance" {
	project = google_project.acceptance.project_id
	service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
	project  = google_project_service.acceptance.project
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name     = "%s"
	key_ring = google_kms_key_ring.key_ring.id
	labels = {
		key = "value"
	}
}

resource "google_kms_crypto_key_version" "crypto_key_version" {
	crypto_key  = google_kms_crypto_key.crypto_key.id
	lifecycle {
		prevent_destroy = true
	}
	state       = "ENABLED"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKeyVersion_patch(preventDestroy, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, state string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "acceptance" {
	project = google_project.acceptance.project_id
	service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
	project  = google_project_service.acceptance.project
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name     = "%s"
	key_ring = google_kms_key_ring.key_ring.id
	labels = {
		key = "value"
	}
}

resource "google_kms_crypto_key_version" "crypto_key_version" {
	crypto_key  = google_kms_crypto_key.crypto_key.id
	lifecycle {
		prevent_destroy = %s
	}
	state = "%s"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, preventDestroy, state)
}
