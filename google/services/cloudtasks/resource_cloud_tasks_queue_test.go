// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudtasks_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudTasksQueue_update(t *testing.T) {
	t.Parallel()

	name := "cloudtasksqueuetest-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueue_full(name),
			},
			{
				ResourceName:            "google_cloud_tasks_queue.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_engine_routing_override.0.service", "app_engine_routing_override.0.version", "app_engine_routing_override.0.instance"},
			},
			{
				Config: testAccCloudTasksQueue_update(name),
			},
			{
				ResourceName:            "google_cloud_tasks_queue.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_engine_routing_override.0.service", "app_engine_routing_override.0.version", "app_engine_routing_override.0.instance"},
			},
		},
	})
}

func TestAccCloudTasksQueue_update2Basic(t *testing.T) {
	t.Parallel()

	name := "cloudtasksqueuetest-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueue_full(name),
			},
			{
				ResourceName:            "google_cloud_tasks_queue.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_engine_routing_override.0.service", "app_engine_routing_override.0.version", "app_engine_routing_override.0.instance"},
			},
			{
				Config: testAccCloudTasksQueue_basic(name),
			},
			{
				ResourceName:            "google_cloud_tasks_queue.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_engine_routing_override.0.service", "app_engine_routing_override.0.version", "app_engine_routing_override.0.instance"},
			},
		},
	})
}

func TestAccCloudTasksQueue_MaxRetryDiffSuppress0s(t *testing.T) {
	t.Parallel()
	testID := acctest.RandString(t, 10)
	cloudTaskName := fmt.Sprintf("tf-test-%s", testID)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudtasksQueueMaxRetry0s(cloudTaskName),
			},
			{
				ResourceName:      "google_cloud_tasks_queue.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Make sure the diff suppression function handles the situation where an
// unexpected time unit is used, e.g., 2.0s instead of 2s or 2.0s instead of
// 2.000s
func TestAccCloudTasksQueue_TimeUnitDiff(t *testing.T) {
	t.Parallel()
	testID := acctest.RandString(t, 10)
	cloudTaskName := fmt.Sprintf("tf-test-%s", testID)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudtasksQueueTimeUnitDiff(cloudTaskName),
			},
			{
				ResourceName:      "google_cloud_tasks_queue.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudTasksQueue_basic(name string) string {
	return fmt.Sprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "%s"
  location = "us-central1"

  retry_config {
    max_attempts = 5
  }

}
`, name)
}

func testAccCloudTasksQueue_full(name string) string {
	return fmt.Sprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "%s"
  location = "us-central1"

  app_engine_routing_override {
    service = "worker"
    version = "1.0"
    instance = "test"
  }

  rate_limits {
    max_concurrent_dispatches = 3
    max_dispatches_per_second = 2
  }

  retry_config {
    max_attempts = 5
    max_retry_duration = "4s"
    max_backoff = "3s"
    min_backoff = "2s"
    max_doublings = 1
	}

	stackdriver_logging_config {
		sampling_ratio = 0.9
	}
}
`, name)
}

func testAccCloudTasksQueue_update(name string) string {
	return fmt.Sprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "%s"
  location = "us-central1"

  app_engine_routing_override {
    service = "main"
    version = "2.0"
    instance = "beta"
  }

  rate_limits {
    max_concurrent_dispatches = 4
    max_dispatches_per_second = 3
  }

  retry_config {
    max_attempts = 6
    max_retry_duration = "5s"
    max_backoff = "4s"
    min_backoff = "3s"
    max_doublings = 2
	}

	stackdriver_logging_config {
		sampling_ratio = 0.1
	}
}
`, name)
}

func testAccCloudtasksQueueMaxRetry0s(cloudTaskName string) string {
	return fmt.Sprintf(`
	resource "google_cloud_tasks_queue" "default" {
		name = "%s"
		location = "us-central1"

		retry_config {
							max_attempts       = -1
							max_backoff        = "3600s"
							max_doublings      = 16
							max_retry_duration = "0s"
							min_backoff        = "0.100s"
		}
	}
`, cloudTaskName)
}

func testAccCloudtasksQueueTimeUnitDiff(cloudTaskName string) string {
	return fmt.Sprintf(`
resource "google_cloud_tasks_queue" "default" {
  name     = "%s"
  location = "us-central1"

  retry_config {
    max_attempts       = -1
    max_backoff        = "5.000s"
    max_doublings      = 16
    max_retry_duration = "1.0s"
    min_backoff        = "0.10s"
  }
}
`, cloudTaskName)
}
