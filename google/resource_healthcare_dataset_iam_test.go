package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const DEFAULT_HEALTHCARE_TEST_LOCATION = "us-central1"

func TestAccHealthcareDatasetIamBinding(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", randInt(t))
	roleId := "roles/healthcare.datasetAdmin"
	datasetName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	datasetId := &healthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccHealthcareDatasetIamBinding_basic(account, datasetName, roleId),
				Check: testAccCheckGoogleHealthcareDatasetIam(t, datasetId.datasetId(), roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_dataset_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s %s", datasetId.terraformId(), roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccHealthcareDatasetIamBinding_update(account, datasetName, roleId),
				Check: testAccCheckGoogleHealthcareDatasetIam(t, datasetId.datasetId(), roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_dataset_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s %s", datasetId.terraformId(), roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHealthcareDatasetIamMember(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", randInt(t))
	roleId := "roles/healthcare.datasetViewer"
	datasetName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	datasetId := &healthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccHealthcareDatasetIamMember_basic(account, datasetName, roleId),
				Check: testAccCheckGoogleHealthcareDatasetIam(t, datasetId.datasetId(), roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_dataset_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", datasetId.terraformId(), roleId, account, projectId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccHealthcareDatasetIamPolicy(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", randInt(t))
	roleId := "roles/healthcare.datasetAdmin"
	datasetName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	datasetId := &healthcareDatasetId{
		Project:  projectId,
		Location: DEFAULT_HEALTHCARE_TEST_LOCATION,
		Name:     datasetName,
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcareDatasetIamPolicy_basic(account, datasetName, roleId),
				Check: testAccCheckGoogleHealthcareDatasetIam(t, datasetId.datasetId(), roleId, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, projectId),
				}),
			},
			{
				ResourceName:      "google_healthcare_dataset_iam_policy.foo",
				ImportStateId:     datasetId.terraformId(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleHealthcareDatasetIam(t *testing.T, datasetId, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		p, err := config.clientHealthcare.Projects.Locations.Datasets.GetIamPolicy(datasetId).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

// We are using a custom role since iam_binding is authoritative on the member list and
// we want to avoid removing members from an existing role to prevent unwanted side effects.
func testAccHealthcareDatasetIamBinding_basic(account, datasetName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "%s"
}

resource "google_healthcare_dataset_iam_binding" "foo" {
  dataset_id = google_healthcare_dataset.dataset.id
  role       = "%s"
  members    = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, account, datasetName, roleId)
}

func testAccHealthcareDatasetIamBinding_update(account, datasetName, roleId string) string {
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
  location = "%s"
  name     = "%s"
}

resource "google_healthcare_dataset_iam_binding" "foo" {
  dataset_id = google_healthcare_dataset.dataset.id
  role       = "%s"
  members = [
    "serviceAccount:${google_service_account.test_account.email}",
    "serviceAccount:${google_service_account.test_account_2.email}",
  ]
}
`, account, account, DEFAULT_HEALTHCARE_TEST_LOCATION, datasetName, roleId)
}

func testAccHealthcareDatasetIamMember_basic(account, datasetName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "%s"
  name     = "%s"
}

resource "google_healthcare_dataset_iam_member" "foo" {
  dataset_id = google_healthcare_dataset.dataset.id
  role       = "%s"
  member     = "serviceAccount:${google_service_account.test_account.email}"
}
`, account, DEFAULT_HEALTHCARE_TEST_LOCATION, datasetName, roleId)
}

func testAccHealthcareDatasetIamPolicy_basic(account, datasetName, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_healthcare_dataset" "dataset" {
  location = "%s"
  name     = "%s"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%s"

    members = ["serviceAccount:${google_service_account.test_account.email}"]
  }
}

resource "google_healthcare_dataset_iam_policy" "foo" {
  dataset_id  = google_healthcare_dataset.dataset.id
  policy_data = data.google_iam_policy.foo.policy_data
}
`, account, DEFAULT_HEALTHCARE_TEST_LOCATION, datasetName, roleId)
}
