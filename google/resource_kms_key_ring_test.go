package google

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestKeyRingIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedKeyRingId   string
		Config              *Config
	}{
		"id is in project/location/keyRingName format": {
			ImportId:            "test-project/us-central1/test-key-ring",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring",
			ExpectedKeyRingId:   "projects/test-project/locations/us-central1/keyRings/test-key-ring",
		},
		"id contains name that is longer than 63 characters": {
			ImportId:      "test-project/us-central1/can-you-believe-that-this-key-ring-name-is-exactly-64-characters",
			ExpectedError: true,
		},
		"id is in location/keyRingName format": {
			ImportId:            "us-central1/test-key-ring",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring",
			ExpectedKeyRingId:   "projects/test-project/locations/us-central1/keyRings/test-key-ring",
			Config:              &Config{Project: "test-project"},
		},
		"id is in location/keyRingName format without project in config": {
			ImportId:      "us-central1/test-key-ring",
			ExpectedError: true,
			Config:        &Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		keyRingId, err := parseKmsKeyRingId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if keyRingId.terraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, keyRingId.terraformId())
		}

		if keyRingId.keyRingId() != tc.ExpectedKeyRingId {
			t.Fatalf("bad: %s, expected KeyRing ID to be `%s` but is `%s`", tn, tc.ExpectedKeyRingId, keyRingId.keyRingId())
		}
	}
}

func TestAccKmsKeyRing_basic(t *testing.T) {
	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleKmsKeyRingWasRemovedFromState("google_kms_key_ring.key_ring"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsKeyRingExists("google_kms_key_ring.key_ring"),
				),
			},
			resource.TestStep{
				Config: testGoogleKmsKeyRing_removed(projectId, projectOrg, projectBillingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsKeyRingWasRemovedFromState("google_kms_key_ring.key_ring"),
				),
			},
		},
	})
}

func testAccCheckGoogleKmsKeyRingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		keyRingId := &kmsKeyRingId{
			Project:  rs.Primary.Attributes["project"],
			Location: rs.Primary.Attributes["location"],
			Name:     rs.Primary.Attributes["name"],
		}

		listKeyRingsResponse, err := config.clientKms.Projects.Locations.KeyRings.List(keyRingId.parentId()).Do()
		if err != nil {
			return fmt.Errorf("Error listing KeyRings: %s", err)
		}

		for _, keyRing := range listKeyRingsResponse.KeyRings {
			log.Printf("[DEBUG] Found KeyRing: %s", keyRing.Name)

			if keyRing.Name == keyRingId.keyRingId() {
				return nil
			}
		}

		return fmt.Errorf("KeyRing not found: %s", keyRingId.keyRingId())
	}
}

/*
	KMS KeyRings cannot be deleted. This ensures that the KeyRing resource was removed from state,
	even though the server-side resource was not removed.
*/
func testAccCheckGoogleKmsKeyRingWasRemovedFromState(resourceName string) resource.TestCheckFunc {
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
func testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
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
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName)
}

func testGoogleKmsKeyRing_removed(projectId, projectOrg, projectBillingAccount string) string {
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
