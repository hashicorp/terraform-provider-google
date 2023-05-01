package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestKeyRingIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedKeyRingId   string
		Config              *transport_tpg.Config
	}{
		"id is in project/location/keyRingName format": {
			ImportId:            "test-project/us-central1/test-key-ring",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring",
			ExpectedKeyRingId:   "projects/test-project/locations/us-central1/keyRings/test-key-ring",
		},
		"id is in domain:project/location/keyRingName format": {
			ImportId:            "example.com:test-project/us-central1/test-key-ring",
			ExpectedError:       false,
			ExpectedTerraformId: "example.com:test-project/us-central1/test-key-ring",
			ExpectedKeyRingId:   "projects/example.com:test-project/locations/us-central1/keyRings/test-key-ring",
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
			Config:              &transport_tpg.Config{Project: "test-project"},
		},
		"id is in location/keyRingName format without project in config": {
			ImportId:      "us-central1/test-key-ring",
			ExpectedError: true,
			Config:        &transport_tpg.Config{Project: ""},
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

		if keyRingId.TerraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, keyRingId.TerraformId())
		}

		if keyRingId.KeyRingId() != tc.ExpectedKeyRingId {
			t.Fatalf("bad: %s, expected KeyRing ID to be `%s` but is `%s`", tn, tc.ExpectedKeyRingId, keyRingId.KeyRingId())
		}
	}
}

func TestAccKmsKeyRing_basic(t *testing.T) {
	projectId := fmt.Sprintf("tf-test-%d", RandInt(t))
	projectOrg := acctest.GetTestOrgFromEnv(t)
	projectBillingAccount := acctest.GetTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleKmsKeyRingWasRemovedFromState("google_kms_key_ring.key_ring"),
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName),
			},
			{
				ResourceName:      "google_kms_key_ring.key_ring",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsKeyRing_removed(projectId, projectOrg, projectBillingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsKeyRingWasRemovedFromState("google_kms_key_ring.key_ring"),
				),
			},
		},
	})
}

// KMS KeyRings cannot be deleted. This ensures that the KeyRing resource was removed from state,
// even though the server-side resource was not removed.
func testAccCheckGoogleKmsKeyRingWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}

// This test runs in its own project, otherwise the test project would start to get filled
// with undeletable resources
func testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
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

func testGoogleKmsKeyRing_removed(projectId, projectOrg, projectBillingAccount string) string {
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
`, projectId, projectId, projectOrg, projectBillingAccount)
}
