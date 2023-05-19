package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigqueryDatasetIamBinding(t *testing.T) {
	t.Parallel()

	dataset := "tf_test_dataset_iam_" + RandString(t, 10)
	account := "tf-test-bq-iam-" + RandString(t, 10)
	role := "roles/bigquery.dataViewer"

	importId := fmt.Sprintf("projects/%s/datasets/%s %s",
		acctest.GetTestProjectFromEnv(), dataset, role)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigqueryDatasetIamBinding_basic(dataset, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigquery_dataset_iam_binding.binding", "role", role),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccBigqueryDatasetIamBinding_update(dataset, account, role),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryDatasetIamMember(t *testing.T) {
	t.Parallel()

	dataset := "tf_test_dataset_iam_" + RandString(t, 10)
	account := "tf-test-bq-iam-" + RandString(t, 10)
	role := "roles/editor"

	importId := fmt.Sprintf("projects/%s/datasets/%s %s serviceAccount:%s",
		acctest.GetTestProjectFromEnv(),
		dataset,
		role,
		serviceAccountCanonicalEmail(account))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigqueryDatasetIamMember(dataset, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigquery_dataset_iam_member.member", "role", role),
					resource.TestCheckResourceAttr(
						"google_bigquery_dataset_iam_member.member", "member", "serviceAccount:"+serviceAccountCanonicalEmail(account)),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryDatasetIamPolicy(t *testing.T) {
	t.Parallel()

	dataset := "tf_test_dataset_iam_" + RandString(t, 10)
	account := "tf-test-bq-iam-" + RandString(t, 10)
	role := "roles/bigquery.dataOwner"

	importId := fmt.Sprintf("projects/%s/datasets/%s",
		acctest.GetTestProjectFromEnv(), dataset)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigqueryDatasetIamPolicy(dataset, account, role),
				Check:  resource.TestCheckResourceAttrSet("data.google_bigquery_dataset_iam_policy.policy", "policy_data"),
			},
			{
				ResourceName:      "google_bigquery_dataset_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigqueryDatasetIamBinding_basic(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigquery Dataset Iam Testing Account"
}

resource "google_bigquery_dataset_iam_binding" "binding" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
}
`, dataset, account, account, role)
}

func testAccBigqueryDatasetIamBinding_update(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_bigquery_dataset_iam_binding" "binding" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]
}
`, dataset, account, account, role)
}

func testAccBigqueryDatasetIamMember(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Bigquery Dataset IAM Testing Account"
}

resource "google_bigquery_dataset_iam_member" "member" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role     = "%s"
  member   = "serviceAccount:${google_service_account.test-account.email}"
}
`, dataset, account, role)
}

func testAccBigqueryDatasetIamPolicy(dataset, account, role string) string {
	return fmt.Sprintf(testBigqueryDatasetIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Bigquery Dataset IAM Testing Account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "%s"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_bigquery_dataset_iam_policy" "policy" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
  policy_data = data.google_iam_policy.policy.policy_data
}

data "google_bigquery_dataset_iam_policy" "policy" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
}
`, dataset, account, role)
}

var testBigqueryDatasetIam = `
resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}
`
