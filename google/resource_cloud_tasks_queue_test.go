package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCloudTasksQueue_update(t *testing.T) {
	t.Parallel()

	name := "cloudtasksqueuetest-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

	name := "cloudtasksqueuetest-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
}
`, name)
}
