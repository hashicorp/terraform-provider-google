package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCloudFunctionsFunctionIamBinding(t *testing.T) {
	t.Parallel()

	randSuffix := acctest.RandString(10)
	zipFilePath, err := createZIPArchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Fatalf("err while trying to prepare ZIP for source for cloud functions IAM test: %s", err)
	}
	function := "tf-function-func-" + randSuffix
	serviceAccount := "tf-function-sa-" + randSuffix

	role := "roles/cloudfunctions.viewer"

	project := getTestRegionFromEnv()
	region := getTestRegionFromEnv()
	importId := fmt.Sprintf("%s/%s/%s %s", project, region, function, role)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunctionIamBinding_basic(function, serviceAccount, zipFilePath, role, randSuffix),
			},
			{
				ResourceName:      "google_cloudfunctions_function_iam_binding.function_binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccCloudFunctionsFunctionIamBinding_update(function, serviceAccount, zipFilePath, role, randSuffix),
			},
			{
				ResourceName:      "google_cloudfunctions_function_iam_binding.function_binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunctionIamMember(t *testing.T) {
	t.Parallel()

	randSuffix := acctest.RandString(10)
	zipFilePath, err := createZIPArchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Fatalf("err while trying to prepare ZIP for source for cloud functions IAM test: %s", err)
	}
	function := "tf-function-func-" + randSuffix
	serviceAccount := "tf-function-sa-" + randSuffix
	serviceAccountEmail := serviceAccountCanonicalEmail(serviceAccount)
	role := "roles/cloudfunctions.viewer"

	project := getTestRegionFromEnv()
	region := getTestRegionFromEnv()
	importId := fmt.Sprintf("%s/%s/%s %s serviceAccount:%s", project, region, function, role, serviceAccountEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunctionIamMember_basic(function, serviceAccount, zipFilePath, role, randSuffix),
			},
			{
				ResourceName:      "google_cloudfunctions_function_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunctionIamPolicy(t *testing.T) {
	t.Parallel()

	randSuffix := acctest.RandString(10)
	zipFilePath, err := createZIPArchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Fatalf("err while trying to prepare ZIP for source for cloud functions IAM test: %s", err)
	}
	function := "tf-function-func-" + randSuffix
	serviceAccount := "tf-cloudfunc-sa-" + randSuffix

	project := getTestRegionFromEnv()
	region := getTestRegionFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunctionIamPolicy_basic(function, serviceAccount, zipFilePath, randSuffix),
			},
			// Test a few import formats
			{
				ResourceName:      "google_cloudfunctions_function_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, region, function),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_cloudfunctions_function_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s", project, region, function),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_cloudfunctions_function_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", region, function),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudFunctionsFunctionIamBinding_basic(function, sa, zipFilePath, role, randString string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "test-cloudfunc-iam-bucket-%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  description           = "test function"
  available_memory_mb   = 128

  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  entry_point           = "helloGET"
  trigger_http          = true
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Cloud Functions Function IAM test account"
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Cloud Functions Function IAM test account"
}

resource "google_cloudfunctions_function_iam_binding" "function_binding" {
  project  = "${google_cloudfunctions_function.function.project}"
  region   = "${google_cloudfunctions_function.function.region}"
  function = "${google_cloudfunctions_function.function.name}"
  role     = "%s"
  members  = [
	"serviceAccount:${google_service_account.test-account-1.email}", 
  ]
}`, randString, zipFilePath, function, sa, sa, role)
}

func testAccCloudFunctionsFunctionIamBinding_update(function, sa, zipFilePath, role, randString string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "test-cloudfunc-iam-bucket-%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  description           = "test function"
  available_memory_mb   = 128

  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  entry_point           = "helloGET"
  trigger_http          = true
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Cloud Functions Function IAM test account"
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Cloud Functions Function IAM test account"
}

resource "google_cloudfunctions_function_iam_binding" "function_binding" {
  project  = "${google_cloudfunctions_function.function.project}"
  region   = "${google_cloudfunctions_function.function.region}"
  function = "${google_cloudfunctions_function.function.name}"
  role     = "%s"
  members  = [
	"serviceAccount:${google_service_account.test-account-1.email}", 
	"serviceAccount:${google_service_account.test-account-2.email}", 
  ]
}

`, randString, zipFilePath, function, sa, sa, role)
}

func testAccCloudFunctionsFunctionIamMember_basic(function, sa, zipFilePath, role, randString string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "test-cloudfunc-iam-bucket-%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  description           = "test function"
  available_memory_mb   = 128

  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  entry_point           = "helloGET"
  trigger_http          = true
}

resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Cloud Functions Function IAM test account"
}

resource "google_cloudfunctions_function_iam_member" "member" {
  project  = "${google_cloudfunctions_function.function.project}"
  region   = "${google_cloudfunctions_function.function.region}"
  function = "${google_cloudfunctions_function.function.name}"
  role     = "%s"
  member   = "serviceAccount:${google_service_account.test-account.email}" 
}
`, randString, zipFilePath, function, sa, role)
}

func testAccCloudFunctionsFunctionIamPolicy_basic(function, sa, zipFilePath, randString string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "test-cloudfunc-iam-bucket-%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  description           = "test function"
  available_memory_mb   = 128

  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  entry_point           = "helloGET"
  trigger_http          = true
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Cloud Functions Function IAM test account"
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Cloud Functions Function IAM test account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "roles/viewer"
    members = [
		"serviceAccount:${google_service_account.test-account-1.email}"
	]
  }

  binding {
    role    = "roles/cloudfunctions.viewer"
    members = [
		"serviceAccount:${google_service_account.test-account-2.email}"
	]
  }
}

resource "google_cloudfunctions_function_iam_policy" "function_policy" {
  project     = "${google_cloudfunctions_function.function.project}"
  region      = "${google_cloudfunctions_function.function.region}"
  function    = "${google_cloudfunctions_function.function.name}"
  policy_data = "${data.google_iam_policy.policy.policy_data}"
}
`, randString, zipFilePath, function, sa, sa)
}
