package google

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"google.golang.org/api/compute/v1"
)

const (
	testDataflowJobTemplateWordCountUrl = "gs://dataflow-templates/latest/Word_Count"
	testDataflowJobSampleFileUrl        = "gs://dataflow-samples/shakespeare/various.txt"
)

func TestAccDataflowJob_basic(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	zone := "us-central1-f"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_zone(bucket, job, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
				),
			},
		},
	})
}

func TestAccDataflowJob_withRegion(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobRegionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_region(bucket, job),
				Check: resource.ComposeTestCheckFunc(
					testAccRegionalDataflowJobExists(t, "google_dataflow_job.big_data", "us-central1"),
				),
			},
		},
	})
}

func TestAccDataflowJob_withServiceAccount(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	accountId := "tf-test-dataflow-sa" + randStr

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_serviceAccount(bucket, job, accountId),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
					testAccDataflowJobHasServiceAccount(t, "google_dataflow_job.big_data", accountId),
				),
			},
		},
	})
}

func TestAccDataflowJob_withNetwork(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	network := "tf-test-dataflow-net" + randStr

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_network(bucket, job, network),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
					testAccDataflowJobHasNetwork(t, "google_dataflow_job.big_data", network),
				),
			},
		},
	})
}

func TestAccDataflowJob_withSubnetwork(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	network := "tf-test-dataflow-net" + randStr
	subnetwork := "tf-test-dataflow-subnet" + randStr

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_subnetwork(bucket, job, network, subnetwork),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
					testAccDataflowJobHasSubnetwork(t, "google_dataflow_job.big_data", subnetwork),
				),
			},
		},
	})
}

func TestAccDataflowJob_withLabels(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	key := "my-label"
	value := "my-value"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_labels(bucket, job, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.with_labels"),
					testAccDataflowJobHasLabels(t, "google_dataflow_job.with_labels", key),
				),
			},
		},
	})
}

func TestAccDataflowJob_withIpConfig(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_ipConfig(bucket, job),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
				),
			},
		},
	})
}

func testAccCheckDataflowJobDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataflow_job" {
				continue
			}

			config := googleProviderConfig(t)
			job, err := config.clientDataflow.Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
			if job != nil {
				if _, ok := dataflowTerminalStatesMap[job.CurrentState]; !ok {
					return fmt.Errorf("Job still present")
				}
			} else if err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccCheckDataflowJobRegionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataflow_job" {
				continue
			}

			config := googleProviderConfig(t)
			job, err := config.clientDataflow.Projects.Locations.Jobs.Get(config.Project, "us-central1", rs.Primary.ID).Do()
			if job != nil {
				if _, ok := dataflowTerminalStatesMap[job.CurrentState]; !ok {
					return fmt.Errorf("Job still present")
				}
			} else if err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccDataflowJobExists(t *testing.T, resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource %q not in state", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := googleProviderConfig(t)
		_, err := config.clientDataflow.Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("could not confirm Dataflow Job %q exists: %v", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccDataflowJobHasNetwork(t *testing.T, res, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		instanceTmpl, err := testAccDataflowJobGetGeneratedInstanceTemplate(t, s, res)
		if err != nil {
			return fmt.Errorf("Error getting dataflow job instance template: %s", err)
		}
		if len(instanceTmpl.Properties.NetworkInterfaces) == 0 {
			return fmt.Errorf("no network interfaces in template properties: %+v", instanceTmpl.Properties)
		}
		actual := instanceTmpl.Properties.NetworkInterfaces[0].Network
		if GetResourceNameFromSelfLink(actual) != GetResourceNameFromSelfLink(expected) {
			return fmt.Errorf("network mismatch: %s != %s", actual, expected)
		}
		return nil
	}
}

func testAccDataflowJobHasSubnetwork(t *testing.T, res, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		instanceTmpl, err := testAccDataflowJobGetGeneratedInstanceTemplate(t, s, res)
		if err != nil {
			return fmt.Errorf("Error getting dataflow job instance template: %s", err)
		}
		if len(instanceTmpl.Properties.NetworkInterfaces) == 0 {
			return fmt.Errorf("no network interfaces in template properties: %+v", instanceTmpl.Properties)
		}
		actual := instanceTmpl.Properties.NetworkInterfaces[0].Subnetwork
		if GetResourceNameFromSelfLink(actual) != GetResourceNameFromSelfLink(expected) {
			return fmt.Errorf("subnetwork mismatch: %s != %s", actual, expected)
		}
		return nil
	}
}

func testAccDataflowJobHasServiceAccount(t *testing.T, res, expectedId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		instanceTmpl, err := testAccDataflowJobGetGeneratedInstanceTemplate(t, s, res)
		if err != nil {
			return fmt.Errorf("Error getting dataflow job instance template: %s", err)
		}
		accounts := instanceTmpl.Properties.ServiceAccounts
		if len(accounts) != 1 {
			return fmt.Errorf("Found multiple service accounts (%d) for dataflow job %q, expected 1", len(accounts), res)
		}
		actualId := strings.Split(accounts[0].Email, "@")[0]
		if expectedId != actualId {
			return fmt.Errorf("service account mismatch, expected account ID = %q, actual email = %q", expectedId, accounts[0].Email)
		}
		return nil
	}
}

func testAccDataflowJobGetGeneratedInstanceTemplate(t *testing.T, s *terraform.State, res string) (*compute.InstanceTemplate, error) {
	rs, ok := s.RootModule().Resources[res]
	if !ok {
		return nil, fmt.Errorf("resource %q not in state", res)
	}
	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("resource %q does not have an ID set", res)
	}
	filter := fmt.Sprintf("properties.labels.dataflow_job_id = %s", rs.Primary.ID)

	config := googleProviderConfig(t)

	var instanceTemplate *compute.InstanceTemplate

	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		instanceTemplates, rerr := config.clientCompute.InstanceTemplates.
			List(config.Project).
			Filter(filter).
			MaxResults(2).
			Fields("items/properties").Do()
		if rerr != nil {
			return resource.NonRetryableError(rerr)
		}
		if len(instanceTemplates.Items) == 0 {
			return resource.RetryableError(fmt.Errorf("no instance template found for dataflow job %q", rs.Primary.ID))
		}
		if len(instanceTemplates.Items) > 1 {
			return resource.NonRetryableError(fmt.Errorf("Wrong number of matching instance templates for dataflow job: %s, %d", rs.Primary.ID, len(instanceTemplates.Items)))
		}
		instanceTemplate = instanceTemplates.Items[0]
		if instanceTemplate == nil || instanceTemplate.Properties == nil {
			return resource.NonRetryableError(fmt.Errorf("invalid instance template has no properties"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return instanceTemplate, nil
}

func testAccRegionalDataflowJobExists(t *testing.T, res, region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[res]
		if !ok {
			return fmt.Errorf("resource %q not found in state", res)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := googleProviderConfig(t)
		_, err := config.clientDataflow.Projects.Locations.Jobs.Get(config.Project, region, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Job does not exist")
		}

		return nil
	}
}

func testAccDataflowJobHasLabels(t *testing.T, res, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[res]
		if !ok {
			return fmt.Errorf("resource %q not found in state", res)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := googleProviderConfig(t)

		job, err := config.clientDataflow.Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("dataflow job does not exist")
		}

		if job.Labels[key] != rs.Primary.Attributes["labels."+key] {
			return fmt.Errorf("Labels do not match what is stored in state.")
		}

		return nil
	}
}

func testAccDataflowJob_zone(bucket, job, zone string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"
 
  zone    = "%s"

  machine_type      = "n1-standard-2"
  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, job, zone, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_region(bucket, job string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"
  region  = "us-central1"

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }

  on_delete = "cancel"
}
`, bucket, job, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_network(bucket, job, network string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_compute_network" "net" {
  name                    = "%s"
  auto_create_subnetworks = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  network           = google_compute_network.net.name

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, network, job, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_subnetwork(bucket, job, network, subnet string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_compute_network" "net" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  network       = google_compute_network.net.self_link
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  subnetwork        = google_compute_subnetwork.subnet.self_link

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, network, subnet, job, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_serviceAccount(bucket, job, accountId string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_service_account" "dataflow-sa" {
  account_id   = "%s"
  display_name = "DataFlow Service Account"
}

resource "google_storage_bucket_iam_member" "dataflow-gcs" {
  bucket = google_storage_bucket.temp.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.dataflow-sa.email}"
}

resource "google_project_iam_member" "dataflow-worker" {
  role   = "roles/dataflow.worker"
  member = "serviceAccount:${google_service_account.dataflow-sa.email}"
}

resource "google_dataflow_job" "big_data" {
  name = "%s"
  depends_on = [
    google_storage_bucket_iam_member.dataflow-gcs, 
    google_project_iam_member.dataflow-worker
  ]

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }

  service_account_email = google_service_account.dataflow-sa.email
}
`, bucket, accountId, job, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_ipConfig(bucket, job string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  ip_configuration = "WORKER_IP_PRIVATE"

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, job, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_labels(bucket, job, labelKey, labelVal string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_dataflow_job" "with_labels" {
  name = "%s"

  labels = {
    "%s" = "%s"
  }

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, job, labelKey, labelVal, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)

}
