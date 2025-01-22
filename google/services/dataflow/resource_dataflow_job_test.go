// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataflow_test

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/services/dataflow"
	dataflowapi "google.golang.org/api/dataflow/v1b3"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"google.golang.org/api/compute/v1"
)

const (
	testDataflowJobTemplateWordCountUrl = "gs://dataflow-templates/latest/Word_Count"
	testDataflowJobSampleFileUrl        = "gs://dataflow-samples/shakespeare/various.txt"
	testDataflowJobTemplateTextToPubsub = "gs://dataflow-templates/latest/Stream_GCS_Text_to_Cloud_PubSub"
	testDataflowJobRegion               = "us-central1"
)

func TestAccDataflowJob_basic(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob(bucket, job, testDataflowJobRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
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

func TestAccDataflowJobSkipWait_basic(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJobSkipWait(bucket, job, testDataflowJobRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
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
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	accountId := "tf-test-dataflow-sa" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
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
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "region", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withNetwork(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	network := "tf-test-dataflow-net" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
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
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "region", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withSubnetwork(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	network := "tf-test-dataflow-net" + randStr
	subnetwork := "tf-test-dataflow-subnet" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
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
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "subnetwork", "region", "state"},
			},
		},
	})
}

func TestAccDataflowJob_withLabels(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	key := "my-label"
	value := "my-value"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_labels(bucket, job, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
					testAccDataflowJobHasLabels(t, "google_dataflow_job.big_data", key),
				),
			},
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccDataflowJob_withProviderDefaultLabels(t *testing.T) {
	// The test failed if VCR testing is enabled, because the cached provider config is used.
	// With the cached provider config, any changes in the provider default labels will not be applied.
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_withProviderDefaultLabels(bucket, job),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.%", "4"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "effective_labels.%", "7"),
				),
			},
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state", "region", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataflowJob_resourceLabelsOverridesProviderDefaultLabels(bucket, job),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.%", "3"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.default_expiration_ms", "3600000"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.%", "4"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_key1", "value1"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "effective_labels.%", "7"),
				),
			},
			{
				ResourceName:      "google_dataflow_job.big_data",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels field in the state is decided by the configuration.
				// During importing, the configuration is unavailable, so the labels field in the state after importing is empty.
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataflowJob_moveResourceLabelToProviderDefaultLabels(bucket, job),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.default_expiration_ms", "3600000"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.%", "4"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_key1", "value1"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "effective_labels.%", "7"),
				),
			},
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataflowJob_resourceLabelsOverridesProviderDefaultLabels(bucket, job),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.%", "3"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.default_expiration_ms", "3600000"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.%", "4"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_key1", "value1"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "terraform_labels.default_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "effective_labels.%", "7"),
				),
			},
			{
				ResourceName:            "google_dataflow_job.big_data",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "skip_wait_on_job_termination", "state", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataflowJob(bucket, job, testDataflowJobRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_dataflow_job.big_data", "labels.%"),
					// goog-terraform-provisioned: true is added
					resource.TestCheckResourceAttr("google_dataflow_job.big_data", "effective_labels.%", "4"),
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

func TestAccDataflowJob_withIpConfig(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
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
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	key_ring := "tf-test-dataflow-kms-ring-" + randStr
	crypto_key := "tf-test-dataflow-kms-key-" + randStr
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@compute-system.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
		{
			Member: "serviceAccount:service-{project_number}@dataflow-service-producer-prod.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_kms(key_ring, crypto_key, bucket, job, testDataflowJobRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
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
func TestAccDataflowJobWithAdditionalExperiments(t *testing.T) {
	// Dataflow responses include serialized java classes and bash commands
	// This makes body comparison infeasible
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr
	additionalExperiments := []string{"enable_stackdriver_agent_metrics", "shuffle_mode=service"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
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
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	job := "tf-test-dataflow-job-" + randStr
	serviceAccount := "tf-test-dataflow-sa" + randStr

	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_stream(suffix, job, serviceAccount, "google_storage_bucket.bucket1.url", "cancel"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.pubsub_stream"),
					func(s *terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
						defer cancel()
						tick := time.NewTicker(10 * time.Second)
						defer tick.Stop()
						for {
							select {
							case <-tick.C:
								job, err := testAccDataflowGetJob(t, s, "google_dataflow_job.pubsub_stream")
								if err != nil {
									return err
								}
								if job.CurrentState == "JOB_STATE_RUNNING" {
									return nil
								}
							case <-ctx.Done():
								return fmt.Errorf("timeout waiting for Job to reach RUNNING state")
							}
						}
					},
				),
			},
			{
				Config: testAccDataflowJob_stream(suffix, job, serviceAccount, "google_storage_bucket.bucket2.url", "cancel"),
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
	acctest.SkipIfVcr(t)
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	job := "tf-test-dataflow-job-" + randStr
	serviceAccount := "tf-test-dataflow-sa" + randStr

	suffix := acctest.RandString(t, 10)

	// If the update is virtual-only, the ID should remain the same after updating.
	var id string
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataflowJobDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob_stream(suffix, job, serviceAccount, "google_storage_bucket.bucket1.url", "drain"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.pubsub_stream"),
					testAccDataflowSetId(t, "google_dataflow_job.pubsub_stream", &id),
					func(s *terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
						defer cancel()
						tick := time.NewTicker(10 * time.Second)
						defer tick.Stop()
						for {
							select {
							case <-tick.C:
								job, err := testAccDataflowGetJob(t, s, "google_dataflow_job.pubsub_stream")
								if err != nil {
									return err
								}
								if job.CurrentState == "JOB_STATE_RUNNING" {
									return nil
								}
							case <-ctx.Done():
								return fmt.Errorf("timeout waiting for Job to reach RUNNING state")
							}
						}
					},
				),
			},
			{
				Config: testAccDataflowJob_stream(suffix, job, serviceAccount, "google_storage_bucket.bucket1.url", "cancel"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowCheckId(t, "google_dataflow_job.pubsub_stream", &id),
					resource.TestCheckResourceAttr("google_dataflow_job.pubsub_stream", "on_delete", "cancel"),
				),
			},
			{
				ResourceName:            "google_dataflow_job.pubsub_stream",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_delete", "parameters", "transform_name_mapping", "skip_wait_on_job_termination", "region", "state"},
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

			config := acctest.GoogleProviderConfig(t)
			job, err := config.NewDataflowClient(config.UserAgent).Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
			if job != nil {
				var ok bool
				skipWait, err := strconv.ParseBool(rs.Primary.Attributes["skip_wait_on_job_termination"])
				if err != nil {
					return fmt.Errorf("could not parse attribute: %v", err)
				}
				_, ok = dataflow.DataflowTerminalStatesMap[job.CurrentState]
				if !ok && skipWait {
					_, ok = dataflow.DataflowTerminatingStatesMap[job.CurrentState]
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
			config := acctest.GoogleProviderConfig(t)
			job, err := config.NewDataflowClient(config.UserAgent).Projects.Locations.Jobs.Get(config.Project, "us-central1", rs.Primary.ID).Do()
			if job != nil {
				var ok bool
				skipWait, err := strconv.ParseBool(rs.Primary.Attributes["skip_wait_on_job_termination"])
				if err != nil {
					return fmt.Errorf("could not parse attribute: %v", err)
				}
				_, ok = dataflow.DataflowTerminalStatesMap[job.CurrentState]
				if !ok && skipWait {
					_, ok = dataflow.DataflowTerminatingStatesMap[job.CurrentState]
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
		_, err := testAccDataflowGetJob(t, s, resource)
		return err
	}
}

func testAccDataflowGetJob(t *testing.T, s *terraform.State, resource string) (*dataflowapi.Job, error) {
	rs, ok := s.RootModule().Resources[resource]
	if !ok {
		return nil, fmt.Errorf("resource %q not in state", resource)
	}
	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("no ID is set")
	}

	region, ok := rs.Primary.Attributes["region"]
	if !ok {
		region = testDataflowJobRegion
	}
	config := acctest.GoogleProviderConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	job, err := config.NewDataflowClient(config.UserAgent).Projects.Locations.Jobs.Get(config.Project, region, rs.Primary.ID).Context(ctx).View("JOB_VIEW_ALL").Do()
	if err != nil {
		return nil, fmt.Errorf("could not get Dataflow Job 'projects/%s/regions/%s/jobs/%s': %w", config.Project, config.Region, rs.Primary.ID, err)
	}
	return job, nil
}

func testAccDataflowWorkerPool(job *dataflowapi.Job) (*dataflowapi.WorkerPool, error) {
	if job == nil {
		return nil, fmt.Errorf("job is nil")
	}
	if job.Environment == nil {
		return nil, fmt.Errorf("job has no environment: %+v", job)
	}
	if len(job.Environment.WorkerPools) == 0 {
		return nil, fmt.Errorf("job has no worker pools: %+v", job)
	}
	return job.Environment.WorkerPools[0], nil
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
		job, err := testAccDataflowGetJob(t, s, res)
		if err != nil {
			return err
		}
		wp, err := testAccDataflowWorkerPool(job)
		if err != nil {
			return err
		}
		if wp.Network != expected {
			return fmt.Errorf("network mismatch: %s != %s", wp.Network, expected)
		}
		return nil
	}
}

func testAccDataflowJobHasSubnetwork(t *testing.T, res, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		job, err := testAccDataflowGetJob(t, s, res)
		if err != nil {
			return err
		}
		wp, err := testAccDataflowWorkerPool(job)
		if err != nil {
			return err
		}
		got := path.Base(wp.Subnetwork)
		if got != expected {
			return fmt.Errorf("network mismatch: %s != %s", got, expected)
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

	config := acctest.GoogleProviderConfig(t)

	var instanceTemplate *compute.InstanceTemplate

	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		instanceTemplates, rerr := config.NewComputeClient(config.UserAgent).RegionInstanceTemplates.
			List(config.Project, config.Region).
			Filter(filter).
			MaxResults(2).
			Fields("items/properties").Do()
		if rerr != nil {
			return resource.NonRetryableError(rerr)
		}
		if len(instanceTemplates.Items) == 0 {
			return resource.RetryableError(fmt.Errorf("no instance template found for dataflow job 'projects/%s/regions/%s/jobs/%s'", config.Project, config.Region, rs.Primary.ID))
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
		config := acctest.GoogleProviderConfig(t)
		_, err := config.NewDataflowClient(config.UserAgent).Projects.Locations.Jobs.Get(config.Project, region, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Job does not exist")
		}

		return nil
	}
}

func testAccDataflowJobHasLabels(t *testing.T, res, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		job, err := testAccDataflowGetJob(t, s, res)
		if err != nil {
			return err
		}
		rs, ok := s.RootModule().Resources[res]
		if !ok {
			return fmt.Errorf("resource %q not found in state", res)
		}

		if job.Labels[key] != rs.Primary.Attributes["labels."+key] {
			return fmt.Errorf("Labels do not match what is stored in state.")
		}

		return nil
	}
}

func testAccDataflowJobHasExperiments(t *testing.T, res string, experiments []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		job, err := testAccDataflowGetJob(t, s, res)
		if err != nil {
			return err
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
		job, err := testAccDataflowGetJob(t, s, res)
		if err != nil {
			return err
		}
		if job.Environment == nil {
			return fmt.Errorf("job has no environment: %+v", job)
		}
		if job.Environment.SdkPipelineOptions == nil {
			return fmt.Errorf("SDK pipeline options are nil")
		}
		sdkPipelineOptions, err := tpgresource.ConvertToMap(job.Environment.SdkPipelineOptions)
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

func testAccDataflowJob(bucket, job, region string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  region = "%s"

  machine_type      = "e2-standard-2"
  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, job, region, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJobSkipWait(bucket, job, region string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  region = "%s"

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
`, bucket, job, region, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_network(bucket, job, network string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
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
  uniform_bucket_level_access = true
}

resource "google_compute_network" "net" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  network       = google_compute_network.net.name
  region 		= "%s"
}

resource "google_dataflow_job" "big_data" {
  name = "%s"
  region     = "%s"
  subnetwork = google_compute_subnetwork.subnet.self_link

  template_gcs_path = "%s"
  temp_gcs_location = google_storage_bucket.temp.url
  parameters = {
    inputFile = "%s"
    output    = "${google_storage_bucket.temp.url}/output"
  }
  on_delete = "cancel"
}
`, bucket, network, subnet, testDataflowJobRegion, job, testDataflowJobRegion, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_serviceAccount(bucket, job, accountId string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
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
  region = "%s"
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
`, bucket, accountId, job, testDataflowJobRegion, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_ipConfig(bucket, job string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
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
  uniform_bucket_level_access = true
}

resource "google_dataflow_job" "big_data" {
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

func testAccDataflowJob_withProviderDefaultLabels(bucket, job string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
}

resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  labels = {
    env                   = "foo"
    default_expiration_ms = 3600000
  }

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

func testAccDataflowJob_resourceLabelsOverridesProviderDefaultLabels(bucket, job string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
}

resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  labels = {
    env                   = "foo"
    default_expiration_ms = 3600000
    default_key1          = "value1"
  }

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

func testAccDataflowJob_moveResourceLabelToProviderDefaultLabels(bucket, job string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
	default_key1 = "default_value1"
	env          = "foo"
  }
}

resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  labels = {
    default_expiration_ms = 3600000
    default_key1          = "value1"
  }

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

func testAccDataflowJob_kms(key_ring, crypto_key, bucket, job, region string) string {
	return fmt.Sprintf(`
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
  uniform_bucket_level_access = true
}

resource "google_dataflow_job" "big_data" {
  name = "%s"

  region = "%s"

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
`, key_ring, crypto_key, bucket, job, region, testDataflowJobTemplateWordCountUrl, testDataflowJobSampleFileUrl)
}

func testAccDataflowJob_additionalExperiments(bucket string, job string, experiments []string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
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

func testAccDataflowJob_stream(suffix, job, serviceAccount, tempLocation, ondelete string) string {
	return fmt.Sprintf(`

data "google_project" "project" {}

resource "google_pubsub_topic" "topic" {
	name     = "tf-test-dataflow-job-%s"
}
resource "google_storage_bucket" "bucket1" {
	name          = "tf-test-bucket1-%s"
	location      = "US"
	force_destroy = true
    uniform_bucket_level_access = true
}
resource "google_storage_bucket" "bucket2" {
	name          = "tf-test-bucket2-%s"
	location      = "US"
	force_destroy = true
    uniform_bucket_level_access = true
}

resource "google_service_account" "dataflow-sa" {
  account_id   = "%s"
  display_name = "DataFlow Service Account"
}

resource "google_storage_bucket_iam_member" "bucket1" {
  bucket = google_storage_bucket.bucket1.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.dataflow-sa.email}"
}

resource "google_storage_bucket_iam_member" "bucket2" {
  bucket = google_storage_bucket.bucket2.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.dataflow-sa.email}"
}

resource "google_project_iam_member" "dataflow-worker" {
  project = data.google_project.project.project_id
  role   = "roles/dataflow.worker"
  member = "serviceAccount:${google_service_account.dataflow-sa.email}"
}

resource "time_sleep" "wait_bind_iam_roles" {
  depends_on = [google_project_iam_member.dataflow-worker, google_storage_bucket_iam_member.bucket1, google_storage_bucket_iam_member.bucket2]
  create_duration = "300s"
}

resource "google_dataflow_job" "pubsub_stream" {
	depends_on = [time_sleep.wait_bind_iam_roles]
	name = "%s"
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
    service_account_email = google_service_account.dataflow-sa.email
	on_delete = "%s"
}
  `, suffix, suffix, suffix, serviceAccount, job, testDataflowJobTemplateTextToPubsub, tempLocation, ondelete)
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
  	uniform_bucket_level_access = true
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
