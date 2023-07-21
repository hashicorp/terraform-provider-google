// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package healthcare_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/healthcare"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccHealthcareDicomStoreIamBinding(t *testing.T) {
	t.Parallel()

	projectId := envvar.GetTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "roles/healthcare.dicomStoreAdmin"
	datasetName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	datasetId := &healthcare.HealthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}
	dicomStoreName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccHealthcareDicomStoreIamBinding_basic(account, datasetName, dicomStoreName, roleId),
				Check: testAccCheckGoogleHealthcareDicomStoreIamBindingExists(t, "foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_dicom_store_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", datasetId.TerraformId(), dicomStoreName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccHealthcareDicomStoreIamBinding_update(account, datasetName, dicomStoreName, roleId),
				Check: testAccCheckGoogleHealthcareDicomStoreIamBindingExists(t, "foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_dicom_store_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", datasetId.TerraformId(), dicomStoreName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHealthcareDicomStoreIamMember(t *testing.T) {
	t.Parallel()

	projectId := envvar.GetTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "roles/healthcare.dicomEditor"
	datasetName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	datasetId := &healthcare.HealthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}
	dicomStoreName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccHealthcareDicomStoreIamMember_basic(account, datasetName, dicomStoreName, roleId),
				Check: testAccCheckGoogleHealthcareDicomStoreIamMemberExists(t, "foo", roleId,
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				),
			},
			{
				ResourceName:      "google_healthcare_dicom_store_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", datasetId.TerraformId(), dicomStoreName, roleId, account, projectId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHealthcareDicomStoreIamPolicy(t *testing.T) {
	t.Parallel()

	projectId := envvar.GetTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "roles/healthcare.dicomViewer"
	datasetName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	datasetId := &healthcare.HealthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}
	dicomStoreName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Policy creation (no update for policy, no need to test)
				Config: testAccHealthcareDicomStoreIamPolicy_basic(account, datasetName, dicomStoreName, roleId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleHealthcareDicomStoreIamPolicyExists(t, "foo", roleId,
						fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					),
					resource.TestCheckResourceAttrSet("data.google_healthcare_dicom_store_iam_policy.foo", "policy_data"),
				),
			},
			{
				ResourceName:      "google_healthcare_dicom_store_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", datasetId.TerraformId(), dicomStoreName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleHealthcareDicomStoreIamBindingExists(t *testing.T, bindingResourceName, roleId string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources[fmt.Sprintf("google_healthcare_dicom_store_iam_binding.%s", bindingResourceName)]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := acctest.GoogleProviderConfig(t)
		dicomStoreId, err := healthcare.ParseHealthcareDicomStoreId(bindingRs.Primary.Attributes["dicom_store_id"], config)

		if err != nil {
			return err
		}

		p, err := config.NewHealthcareClient(config.UserAgent).Projects.Locations.Datasets.DicomStores.GetIamPolicy(dicomStoreId.DicomStoreId()).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == roleId {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", roleId)
	}
}

func testAccCheckGoogleHealthcareDicomStoreIamMemberExists(t *testing.T, n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_healthcare_dicom_store_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		dicomStoreId, err := healthcare.ParseHealthcareDicomStoreId(rs.Primary.Attributes["dicom_store_id"], config)

		if err != nil {
			return err
		}

		p, err := config.NewHealthcareClient(config.UserAgent).Projects.Locations.Datasets.DicomStores.GetIamPolicy(dicomStoreId.DicomStoreId()).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				for _, m := range binding.Members {
					if m == member {
						return nil
					}
				}

				return fmt.Errorf("Missing member %q, got %v", member, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

func testAccCheckGoogleHealthcareDicomStoreIamPolicyExists(t *testing.T, n, role, policy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_healthcare_dicom_store_iam_policy."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		dicomStoreId, err := healthcare.ParseHealthcareDicomStoreId(rs.Primary.Attributes["dicom_store_id"], config)

		if err != nil {
			return err
		}

		p, err := config.NewHealthcareClient(config.UserAgent).Projects.Locations.Datasets.DicomStores.GetIamPolicy(dicomStoreId.DicomStoreId()).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				for _, m := range binding.Members {
					if m == policy {
						return nil
					}
				}

				return fmt.Errorf("Missing policy %q, got %v", policy, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

// We are using a custom role since iam_binding is authoritative on the member list and
// we want to avoid removing members from an existing role to prevent unwanted side effects.
func testAccHealthcareDicomStoreIamBinding_basic(account, datasetName, dicomStoreName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_dicom_store" "dicom_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

resource "google_healthcare_dicom_store_iam_binding" "foo" {
  dicom_store_id = google_healthcare_dicom_store.dicom_store.id
  role           = "%s"
  members        = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, account, datasetName, dicomStoreName, roleId)
}

func testAccHealthcareDicomStoreIamBinding_update(account, datasetName, dicomStoreName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_service_account" "test_account_2" {
  account_id   = "%s-2"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_dicom_store" "dicom_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

resource "google_healthcare_dicom_store_iam_binding" "foo" {
  dicom_store_id = google_healthcare_dicom_store.dicom_store.id
  role           = "%s"
  members = [
    "serviceAccount:${google_service_account.test_account.email}",
    "serviceAccount:${google_service_account.test_account_2.email}",
  ]
}
`, account, account, datasetName, dicomStoreName, roleId)
}

func testAccHealthcareDicomStoreIamMember_basic(account, datasetName, dicomStoreName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_dicom_store" "dicom_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

resource "google_healthcare_dicom_store_iam_member" "foo" {
  dicom_store_id = google_healthcare_dicom_store.dicom_store.id
  role           = "%s"
  member         = "serviceAccount:${google_service_account.test_account.email}"
}
`, account, datasetName, dicomStoreName, roleId)
}

func testAccHealthcareDicomStoreIamPolicy_basic(account, datasetName, dicomStoreName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_dicom_store" "dicom_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%s"

    members = ["serviceAccount:${google_service_account.test_account.email}"]
  }
}

resource "google_healthcare_dicom_store_iam_policy" "foo" {
  dicom_store_id = google_healthcare_dicom_store.dicom_store.id
  policy_data    = data.google_iam_policy.foo.policy_data
}

data "google_healthcare_dicom_store_iam_policy" "foo" {
  dicom_store_id = google_healthcare_dicom_store.dicom_store.id
}
`, account, datasetName, dicomStoreName, roleId)
}
