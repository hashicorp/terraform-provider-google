package google

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/api/compute/v1"
)

const (
	testDataflowJobTemplateWordCountUrl = "gs://dataflow-templates/latest/Word_Count"
	testDataflowJobSampleFileUrl        = "gs://dataflow-samples/shakespeare/various.txt"
	testDataflowJobTemplateTextToPubsub = "gs://dataflow-templates/latest/Stream_GCS_Text_to_Cloud_PubSub"
)

func TestAccDataflowJob_basic(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "zone", "state"},
			},
		},
	})
}

func TestAccDataflowJobSkipWait_basic(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
				Config: testAccDataflowJobSkipWait_zone(bucket, job, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
				),
			},
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "zone", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withRegion(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "region", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withServiceAccount(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withNetwork(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withSubnetwork(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "subnetwork", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withLabels(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
			{
				ResourceName:            "google_dataflow_job.with_labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withIpConfig(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
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
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "ip_configuration", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withKmsKey(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
	t.Parallel()

	randStr := randString(t, 10)
	key_ring := "tf-test-dataflow-kms-ring-" + randStr
	crypto_key := "tf-test-dataflow-kms-key-" + randStr
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	zone := "us-central1-f"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_kms(key_ring, crypto_key, bucket, job, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
				),
			},
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "zone", "state"},
			},
		},
	})
}
func TestAccDataflowJobWithAdditionalExperiments(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	additionalExperiments := []string{"enable_stackdriver_agent_metrics", "shuffle_mode=service"}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_additionalExperiments(bucket, job, additionalExperiments),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.with_additional_experiments"),
					testAccDataflowJobHasExperiments(t, "google_dataflow_job.with_additional_experiments", additionalExperiments),
				),
			},
			{
				ResourceName:            "google_dataflow_job.with_additional_experiments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state", "additional_experiments"},
			},
		},
	})
}

func TestAccDataflowJob_streamUpdate(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
	t.Parallel()

	suffix := randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_updateStream(suffix, "google_storage_bucket.bucket1.url"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.pubsub_stream"),
				),
			},
			{
				Config: testAccDataflowJob_updateStream(suffix, "google_storage_bucket.bucket2.url"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobHasTempLocation(t, "google_dataflow_job.pubsub_stream", "gs://tf-test-bucket2-"+suffix),
				),
			},
			{
				ResourceName:            "google_dataflow_job.pubsub_stream",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "transform_name_mapping", "state"},
			},
		},
	})
}

func TestAccDataflowJob_virtualUpdate(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	skipIfVcr(t)
	t.Parallel()

	suffix := randString(t, 10)

	// If the update is virtual-only, the ID should remain the same after updating.
	var id string
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_virtualUpdate(suffix, "drain"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.pubsub_stream"),
					testAccDataflowSetId(t, "google_dataflow_job.pubsub_stream", &id),
				),
			},
			{
				Config: testAccDataflowJob_virtualUpdate(suffix, "cancel"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowCheckId(t, "google_dataflow_job.pubsub_stream", &id),
					resource.TestCheckResourceAttr("google_dataflow_job.pubsub_stream", "on_delete", "cancel"),
				),
			},
			{
				ResourceName:            "google_dataflow_job.pubsub_stream",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state"},
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
			job, err := config.NewDataflowClient(config.userAgent).Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
			if job != nil {
				var ok bool
				skipWait, err := strconv.ParseBool(rs.Primary.Attributes["skip_wait_on_job_termination"])
				if err != nil {
					return fmt.Errorf("could not parse attribute: %v", err)
				}
				_, ok = dataflowTerminalStatesMap[job.CurrentState]
				if !ok && skipWait {
					_, ok = dataflowTerminatingStatesMap[job.CurrentState]
				}
				if !ok {
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
			job, err := config.NewDataflowClient(config.userAgent).Projects.Locations.Jobs.Get(config.Project, "us-central1", rs.Primary.ID).Do()
			if job != nil {
				var ok bool
				skipWait, err := strconv.ParseBool(rs.Primary.Attributes["skip_wait_on_job_termination"])
				if err != nil {
					return fmt.Errorf("could not parse attribute: %v", err)
				}
				_, ok = dataflowTerminalStatesMap[job.CurrentState]
				if !ok && skipWait {
					_, ok = dataflowTerminatingStatesMap[job.CurrentState]
				}
				if !ok {
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
		_, err := config.NewDataflowClient(config.userAgent).Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("could not confirm Dataflow Job %q exists: %v", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccDataflowSetId(t *testing.T, resource string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource %q not in state", resource)
		}

		*id = rs.Primary.ID
		return nil
	}
}

func testAccDataflowCheckId(t *testing.T, resource string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource %q not in state", resource)
		}

		if rs.Primary.ID != *id {
			return fmt.Errorf("ID did not match. Expected %s, received %s", *id, rs.Primary.ID)
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
		instanceTemplates, rerr := config.NewComputeClient(config.userAgent).InstanceTemplates.
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
		_, err := config.NewDataflowClient(config.userAgent).Projects.Locations.Jobs.Get(config.Project, region, rs.Primary.ID).Do()
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

		job, err := config.NewDataflowClient(config.userAgent).Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("dataflow job does not exist")
		}

		if job.Labels[key] != rs.Primary.Attributes["labels."+key] {
			return fmt.Errorf("Labels do not match what is stored in state.")
		}

		return nil
	}
}

func testAccDataflowJobHasExperiments(t *testing.T, res string, experiments []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[res]
		if !ok {
			return fmt.Errorf("resource %q not found in state", res)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := googleProviderConfig(t)

		job, err := config.NewDataflowClient(config.userAgent).Projects.Jobs.Get(config.Project, rs.Primary.ID).View("JOB_VIEW_ALL").Do()
		if err != nil {
			return fmt.Errorf("dataflow job does not exist")
		}

		for _, expectedExperiment := range experiments {
			var contains = false
			for _, actualExperiment := range job.Environment.Experiments {
				if actualExperiment == expectedExperiment {
					contains = true
				}
			}
			if contains != true {
				return fmt.Errorf("Expected experiment '%s' not found in experiments", expectedExperiment)
			}
		}

		return nil
	}
}

func testAccDataflowJobHasTempLocation(t *testing.T, res, targetLocation string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[res]
		if !ok {
			return fmt.Errorf("resource %q not found in state", res)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := googleProviderConfig(t)

		job, err := config.NewDataflowClient(config.userAgent).Projects.Jobs.Get(config.Project, rs.Primary.ID).View("JOB_VIEW_ALL").Do()
		if err != nil {
			return fmt.Errorf("dataflow job does not exist")
		}
		sdkPipelineOptions, err := ConvertToMap(job.Environment.SdkPipelineOptions)
		if err != nil {
			return err
		}
		optionsMap := sdkPipelineOptions["options"].(map[string]interface{})

		if optionsMap["tempLocation"] != targetLocation {
			return fmt.Errorf("Temp locations do not match. Got %s while expecting %s", optionsMap["tempLocation"], targetLocation)
		}

		return nil
	}
}

func testAccDataflowJob_zone(bucket, job, zone string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"
 
  zone    = "%s"

  machine_type      = "e2-standard-2"
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

func testAccDataflowJobSkipWait_zone(bucket, job, zone string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"
 
  zone    = "%s"

  machine_type      = "e2-standard-2"
  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete                    = "cancel"
  skip_wait_on_job_termination = true
}
`, bucket, job, zone, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_region(bucket, job string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
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
  name          = "%s"
  location      = "US"
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
  name          = "%s"
  location      = "US"
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
data "google_project" "project" {}

resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
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
  project = data.google_project.project.project_id
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
  name          = "%s"
  location      = "US"
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
  name          = "%s"
  location      = "US"
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

func testAccDataflowJob_kms(key_ring, crypto_key, bucket, job, zone string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
}

resource "google_project_iam_member" "kms-project-dataflow-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@dataflow-service-producer-prod.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "kms-project-compute-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com"
}

resource "google_kms_key_ring" "keyring" {
  name     = "%s"
  location = "global"
}

resource "google_kms_crypto_key" "crypto_key" {
  name            = "%s"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "100000s"
}

resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"
 
  zone    = "%s"

  machine_type      = "e2-standard-2"
  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  kms_key_name		= google_kms_crypto_key.crypto_key.id
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, key_ring, crypto_key, bucket, job, zone, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_additionalExperiments(bucket string, job string, experiments []string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_dataflow_job" "with_additional_experiments" {
  name = "%s"

  additional_experiments = ["%s"]

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, job, strings.Join(experiments, `", "`), testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_updateStream(suffix, tempLocation string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
	name     = "tf-test-dataflow-job-%s"
}
resource "google_storage_bucket" "bucket1" {
	name          = "tf-test-bucket1-%s"
	location      = "US"
	force_destroy = true
}
resource "google_storage_bucket" "bucket2" {
	name          = "tf-test-bucket2-%s"
	location      = "US"
	force_destroy = true
}
resource "google_dataflow_job" "pubsub_stream" {
	name = "tf-test-dataflow-job-%s"
	template_gcs_path = "%s"
	temp_gcs_location = %s
	parameters = {
	  inputFilePattern = "${google_storage_bucket.bucket1.url}/*.json"
	  outputTopic    = google_pubsub_topic.topic.id
	}
	transform_name_mapping = {
		name = "test_job"
		env = "test"
	}
	on_delete = "cancel"
}
  `, suffix, suffix, suffix, suffix, testDataflowJobTemplateTextToPubsub, tempLocation)
}

func testAccDataflowJob_virtualUpdate(suffix, onDelete string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
	name     = "tf-test-dataflow-job-%s"
}
resource "google_storage_bucket" "bucket" {
	name          = "tf-test-bucket-%s"
	location      = "US"
	force_destroy = true
}
resource "google_dataflow_job" "pubsub_stream" {
	name = "tf-test-dataflow-job-%s"
	template_gcs_path = "%s"
	temp_gcs_location = google_storage_bucket.bucket.url
	parameters = {
	  inputFilePattern = "${google_storage_bucket.bucket.url}/*.json"
	  outputTopic    = google_pubsub_topic.topic.id
	}
	on_delete = "%s"
}
  `, suffix, suffix, suffix, testDataflowJobTemplateTextToPubsub, onDelete)
}
