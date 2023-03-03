package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataprocJobIamBinding(t *testing.T) {
	t.Parallel()

	cluster := "tf-dataproc-iam-cluster" + RandString(t, 10)
	job := "tf-dataproc-iam-job-" + RandString(t, 10)
	account := "tf-dataproc-iam-" + RandString(t, 10)
	role := "roles/editor"

	importId := fmt.Sprintf("projects/%s/regions/%s/jobs/%s %s",
		GetTestProjectFromEnv(), "us-central1", job, role)

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccDataprocJobIamBinding_basic(cluster, job, account, role),
			},
			{
				ResourceName:      "google_dataproc_job_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccDataprocJobIamBinding_update(cluster, job, account, role),
			},
			{
				ResourceName:      "google_dataproc_job_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataprocJobIamMember(t *testing.T) {
	t.Parallel()

	cluster := "tf-dataproc-iam-cluster" + RandString(t, 10)
	job := "tf-dataproc-iam-jobid-" + RandString(t, 10)
	account := "tf-dataproc-iam-" + RandString(t, 10)
	role := "roles/editor"

	importId := fmt.Sprintf("projects/%s/regions/%s/jobs/%s %s serviceAccount:%s",
		GetTestProjectFromEnv(),
		"us-central1",
		job,
		role,
		serviceAccountCanonicalEmail(account))

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccDataprocJobIamMember(cluster, job, account, role),
			},
			{
				ResourceName:      "google_dataproc_job_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataprocJobIamPolicy(t *testing.T) {
	t.Parallel()

	cluster := "tf-dataproc-iam-cluster" + RandString(t, 10)
	job := "tf-dataproc-iam-jobid-" + RandString(t, 10)
	account := "tf-dataproc-iam-" + RandString(t, 10)
	role := "roles/editor"

	importId := fmt.Sprintf("projects/%s/regions/%s/jobs/%s",
		GetTestProjectFromEnv(), "us-central1", job)

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccDataprocJobIamPolicy(cluster, job, account, role),
			},
			{
				ResourceName:      "google_dataproc_job_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testDataprocIamJobConfig = testDataprocIamSingleNodeCluster + `
resource "google_dataproc_job" "pyspark" {
  region = google_dataproc_cluster.cluster.region

  placement {
    cluster_name = google_dataproc_cluster.cluster.name
  }

  reference {
    job_id = "%s"
  }

  force_delete = true

  pyspark_config {
    main_python_file_uri = "gs://dataproc-examples-2f10d78d114f6aaec76462e3c310f31f/src/pyspark/hello-world/hello-world.py"
    properties = {
      "spark.logConf" = "true"
    }
    logging_config {
      driver_log_levels = {
        "root" = "INFO"
      }
    }
  }
}
`

func testAccDataprocJobIamBinding_basic(cluster, job, account, role string) string {
	return fmt.Sprintf(testDataprocIamJobConfig+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Dataproc Job IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Dataproc Job Iam Testing Account"
}

resource "google_dataproc_job_iam_binding" "binding" {
  job_id = google_dataproc_job.pyspark.reference[0].job_id
  region = "us-central1"
  role   = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
}
`, cluster, job, account, account, role)
}

func testAccDataprocJobIamBinding_update(cluster, job, account, role string) string {
	return fmt.Sprintf(testDataprocIamJobConfig+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Dataproc Job IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Dataproc Job Iam Testing Account"
}

resource "google_dataproc_job_iam_binding" "binding" {
  job_id = google_dataproc_job.pyspark.reference[0].job_id
  region = "us-central1"
  role   = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]
}
`, cluster, job, account, account, role)
}

func testAccDataprocJobIamMember(cluster, job, account, role string) string {
	return fmt.Sprintf(testDataprocIamJobConfig+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Dataproc Job IAM Testing Account"
}

resource "google_dataproc_job_iam_member" "member" {
  job_id = google_dataproc_job.pyspark.reference[0].job_id
  role   = "%s"
  member = "serviceAccount:${google_service_account.test-account.email}"
}
`, cluster, job, account, role)
}

func testAccDataprocJobIamPolicy(cluster, job, account, role string) string {
	return fmt.Sprintf(testDataprocIamJobConfig+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Dataproc Job IAM Testing Account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "%s"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_dataproc_job_iam_policy" "policy" {
  job_id      = google_dataproc_job.pyspark.reference[0].job_id
  region      = "us-central1"
  policy_data = data.google_iam_policy.policy.policy_data
}
`, cluster, job, account, role)
}
