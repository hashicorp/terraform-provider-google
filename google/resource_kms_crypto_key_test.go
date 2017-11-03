package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
)

func TestCryptoKeyIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedCryptoKeyId string
		Config              *Config
	}{
		"id is in project/location/keyRingName/CryptoKeyID format": {
			ImportId:            "test-project/us-central1/test-key-ring/test-key-id",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-id",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-id",
		},
		"id contains name that is longer than 63 characters": {
			ImportId:      "test-project/us-central1/test-key-ring/can-you-believe-that-this-cryptokey-name-is-this-extravagantly-long",
			ExpectedError: true,
		},
		"id is in location/keyRingName/CyptoKeyID format": {
			ImportId:            "us-central1/test-key-ring/test-key-id",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-id",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-id",
			Config:              &Config{Project: "test-project"},
		},
		"id is in location/keyRingName/CyptoKeyID format without project in config": {
			ImportId:      "us-central1/test-key-ring/test-key-id",
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

func TestAccGoogleKmsCryptoKey_basic(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := os.Getenv("GOOGLE_ORG")
	projectBillingAccount := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyExists("google_kms_crypto_key.crypto_key"),
				),
			},
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsCryptoKeyWasRemovedFromState("google_kms_crypto_key.crypto_key"),
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

		cryptoKeyId := &kmsCryptoKeyId{
			Project:  rs.Primary.Attributes["project"],
			Location: rs.Primary.Attributes["location"],
			KeyRing:  rs.Primary.Attributes["key_ring"],
			Name:     rs.Primary.Attributes["name"],
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

/*
	KMS KeyRings cannot be deleted. This ensures that the KeyRing resource was removed from state,
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
	This test runs in its own project, otherwise the test project would start to get filled
	with undeletable resources
*/
func testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name			= "%s"
	project_id		= "%s"
	org_id			= "%s"
	billing_account	= "%s"
}

resource "google_project_services" "acceptance" {
	project  = "${google_project.acceptance.project_id}"
	services = [
		"cloudkms.googleapis.com"
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
  key_ring = "${google_kms_key_ring.key_ring.name}"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testGoogleKmsCryptoKey_removed(projectId, projectOrg, projectBillingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name 			= "%s"
	project_id		= "%s"
	org_id			= "%s"
	billing_account	= "%s"
}

resource "google_project_services" "acceptance" {
	project  = "${google_project.acceptance.project_id}"
	services = [
		"cloudkms.googleapis.com"
	]
}
	`, projectId, projectId, projectOrg, projectBillingAccount)
}
