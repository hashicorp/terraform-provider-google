package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHealthcareFhirStoreIamBinding(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", randInt(t))
	roleId := "roles/healthcare.fhirStoreAdmin"
	datasetName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	datasetId := &healthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}
	fhirStoreName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccHealthcareFhirStoreIamBinding_basic(account, datasetName, fhirStoreName, roleId),
				Check: testAccCheckGoogleHealthcareFhirStoreIamBindingExists(t, "foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_fhir_store_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", datasetId.terraformId(), fhirStoreName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccHealthcareFhirStoreIamBinding_update(account, datasetName, fhirStoreName, roleId),
				Check: testAccCheckGoogleHealthcareFhirStoreIamBindingExists(t, "foo", roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_fhir_store_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s", datasetId.terraformId(), fhirStoreName, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHealthcareFhirStoreIamMember(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", randInt(t))
	roleId := "roles/healthcare.fhirResourceEditor"
	datasetName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	datasetId := &healthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}
	fhirStoreName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccHealthcareFhirStoreIamMember_basic(account, datasetName, fhirStoreName, roleId),
				Check: testAccCheckGoogleHealthcareFhirStoreIamMemberExists(t, "foo", roleId,
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				),
			},
			{
				ResourceName:      "google_healthcare_fhir_store_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", datasetId.terraformId(), fhirStoreName, roleId, account, projectId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHealthcareFhirStoreIamPolicy(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", randInt(t))
	roleId := "roles/healthcare.fhirResourceEditor"
	datasetName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	datasetId := &healthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}
	fhirStoreName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Policy creation (no update for policy, no need to test)
				Config: testAccHealthcareFhirStoreIamPolicy_basic(account, datasetName, fhirStoreName, roleId),
				Check: testAccCheckGoogleHealthcareFhirStoreIamPolicyExists(t, "foo", roleId,
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				),
			},
			{
				ResourceName:      "google_healthcare_fhir_store_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", datasetId.terraformId(), fhirStoreName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleHealthcareFhirStoreIamBindingExists(t *testing.T, bindingResourceName, roleId string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources[fmt.Sprintf("google_healthcare_fhir_store_iam_binding.%s", bindingResourceName)]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := googleProviderConfig(t)
		fhirStoreId, err := parseHealthcareFhirStoreId(bindingRs.Primary.Attributes["fhir_store_id"], config)

		if err != nil {
			return err
		}

		p, err := config.clientHealthcare.Projects.Locations.Datasets.FhirStores.GetIamPolicy(fhirStoreId.fhirStoreId()).Do()
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

func testAccCheckGoogleHealthcareFhirStoreIamMemberExists(t *testing.T, n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_healthcare_fhir_store_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := googleProviderConfig(t)
		fhirStoreId, err := parseHealthcareFhirStoreId(rs.Primary.Attributes["fhir_store_id"], config)

		if err != nil {
			return err
		}

		p, err := config.clientHealthcare.Projects.Locations.Datasets.FhirStores.GetIamPolicy(fhirStoreId.fhirStoreId()).Do()
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

func testAccCheckGoogleHealthcareFhirStoreIamPolicyExists(t *testing.T, n, role, policy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_healthcare_fhir_store_iam_policy."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := googleProviderConfig(t)
		fhirStoreId, err := parseHealthcareFhirStoreId(rs.Primary.Attributes["fhir_store_id"], config)

		if err != nil {
			return err
		}

		p, err := config.clientHealthcare.Projects.Locations.Datasets.FhirStores.GetIamPolicy(fhirStoreId.fhirStoreId()).Do()
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
func testAccHealthcareFhirStoreIamBinding_basic(account, datasetName, fhirStoreName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_fhir_store" "fhir_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

resource "google_healthcare_fhir_store_iam_binding" "foo" {
  fhir_store_id = google_healthcare_fhir_store.fhir_store.id
  role          = "%s"
  members       = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, account, datasetName, fhirStoreName, roleId)
}

func testAccHealthcareFhirStoreIamBinding_update(account, datasetName, fhirStoreName, roleId string) string {
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

resource "google_healthcare_fhir_store" "fhir_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

resource "google_healthcare_fhir_store_iam_binding" "foo" {
  fhir_store_id = google_healthcare_fhir_store.fhir_store.id
  role          = "%s"
  members = [
    "serviceAccount:${google_service_account.test_account.email}",
    "serviceAccount:${google_service_account.test_account_2.email}",
  ]
}
`, account, account, datasetName, fhirStoreName, roleId)
}

func testAccHealthcareFhirStoreIamMember_basic(account, datasetName, fhirStoreName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_fhir_store" "fhir_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

resource "google_healthcare_fhir_store_iam_member" "foo" {
  fhir_store_id = google_healthcare_fhir_store.fhir_store.id
  role          = "%s"
  member        = "serviceAccount:${google_service_account.test_account.email}"
}
`, account, datasetName, fhirStoreName, roleId)
}

func testAccHealthcareFhirStoreIamPolicy_basic(account, datasetName, fhirStoreName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_fhir_store" "fhir_store" {
  dataset  = google_healthcare_dataset.dataset.id
  name     = "%s"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%s"

    members = ["serviceAccount:${google_service_account.test_account.email}"]
  }
}

resource "google_healthcare_fhir_store_iam_policy" "foo" {
  fhir_store_id = google_healthcare_fhir_store.fhir_store.id
  policy_data   = data.google_iam_policy.foo.policy_data
}
`, account, datasetName, fhirStoreName, roleId)
}
