package google

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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

func TestAccKmsCryptoKey_basic(t *testing.T) {
	t.Parallel()

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyExists("google_kms_crypto_key.crypto_key"),
				),
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

func TestAccKmsCryptoKey_rotation(t *testing.T) {
	t.Parallel()

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	location := getTestRegionFromEnv()
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	rotationPeriod := "100000s"
	updatedRotationPeriod := "7776000s"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, rotationPeriod),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyExists("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyHasRotationParams(rotationPeriod, "google_kms_crypto_key.crypto_key"),
				),
			},
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_rotation(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, updatedRotationPeriod),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyExists("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyHasRotationParams(updatedRotationPeriod, "google_kms_crypto_key.crypto_key"),
				),
			},
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_rotationRemoved(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyExists("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyHasRotationParams("", "google_kms_crypto_key.crypto_key"),
				),
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
					testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(projectId, location, keyRingName, cryptoKeyName),
				),
			},
		},
	})
}

func testAccCheckGoogleKmsCryptoKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		keyRingId, err := parseKmsKeyRingId(rs.Primary.Attributes["key_ring"], config)

		if err != nil {
			return err
		}

		cryptoKeyId := &kmsCryptoKeyId{
			KeyRingId: *keyRingId,
			Name:      rs.Primary.Attributes["name"],
		}

		listCryptoKeysResponse, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.List(cryptoKeyId.parentId()).Do()
		if err != nil {
			return fmt.Errorf("Error listing KeyRings: %s", err)
		}

		for _, cryptoKey := range listCryptoKeysResponse.CryptoKeys {
			log.Printf("[DEBUG] Found CryptoKey: %s", cryptoKey.Name)

			if cryptoKey.Name == cryptoKeyId.cryptoKeyId() {
				return nil
			}
		}

		return fmt.Errorf("CryptoKey not found: %s", cryptoKeyId.cryptoKeyId())
	}
}

func testAccCheckGoogleKmsCryptoKeyHasRotationParams(rotationPeriod, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		keyRingId, err := parseKmsKeyRingId(rs.Primary.Attributes["key_ring"], config)

		if err != nil {
			return err
		}

		cryptoKeyId := &kmsCryptoKeyId{
			KeyRingId: *keyRingId,
			Name:      rs.Primary.Attributes["name"],
		}

		getCryptoKeyResponse, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Get(cryptoKeyId.cryptoKeyId()).Do()

		if err != nil {
			return err
		}

		if rotationPeriod != getCryptoKeyResponse.RotationPeriod {
			return fmt.Errorf("Expected rotation period %s to match input %s", getCryptoKeyResponse.RotationPeriod, rotationPeriod)
		}

		if getCryptoKeyResponse.NextRotationTime == "" {
			return nil
		}

		_, err = time.Parse(time.RFC3339Nano, getCryptoKeyResponse.NextRotationTime)

		if err != nil {
			return fmt.Errorf("Failed to parse NextRotationTime timestamp: %s", err)
		}

		return nil
	}
}

/*
	KMS KeyRings cannot be deleted. This ensures that the CryptoKey resource was removed from state,
	even though the server-side resource was not removed.
*/
func testAccCheckGoogleKmsCryptoKeyWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}

/*
	KMS KeyRings cannot be deleted. This ensures that the CryptoKey resource's CryptoKeyVersion
	sub-resources were scheduled to be destroyed, rendering the key itself inoperable.
*/

func testAccCheckGoogleKmsCryptoKeyVersionsDestroyed(projectId, location, keyRingName, cryptoKeyName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		gcpResourceUri := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyRingName, cryptoKeyName)

		response, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.List(gcpResourceUri).Do()

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

/*
	This test runs in its own project, otherwise the test project would start to get filled
	with undeletable resources
*/
func testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.id}"
	rotation_period = "1000000s"
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

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.id}"
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

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.id}"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName)
}
