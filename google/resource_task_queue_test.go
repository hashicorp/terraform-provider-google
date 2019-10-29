package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTaskQueue_basic(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_basic(queueName),
			},
			{
				ResourceName:      "google_task_queue.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTaskQueue_withParams(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_withParams(queueName),
			},
			{
				ResourceName:            "google_task_queue.fizzbuzz",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_engine_routing_override.0.service", "app_engine_routing_override.0.instance", "app_engine_routing_override.0.version"},
			},
		},
	})
}

func TestAccTaskQueue_update(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_basic(queueName),
			},
			{
				Config: testAccTaskQueue_withParams(queueName),
			},
		},
	})
}

func TestAccTaskQueue_forceDestroy(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_basic(queueName),
			},
		},
	})
}

func testAccTaskQueue_basic(name string) string {
	return fmt.Sprintf(`
resource "google_task_queue" "fizzbuzz" {
  name = "%s"

  location = "us-central1"
}`, name)
}

func testAccTaskQueue_withParams(name string) string {
	return fmt.Sprintf(`
resource "google_task_queue" "fizzbuzz" {
  name = "%s"

  location = "us-central1"

  app_engine_routing_override {
    service = "worker"
  }

  rate_limits {
    max_concurrent_dispatches = 1000
    max_dispatches_per_second = 500
  }

  retry_config {
    max_attempts  = 100
    max_backoff   = "3600s"
    min_backoff   = "0.100s"
    max_doublings = 16
  }
}`, name)
}

func testAccTaskQueueDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_task_queue" {
			continue
		}

		attributes := rs.Primary.Attributes

		expectedName := fmt.Sprintf("projects/%s/locations/%s/queues/%s", config.Project, attributes["location"], attributes["name"])

		_, err := config.clientCloudTasks.Projects.Locations.Queues.Get(expectedName).Do()
		if err == nil {
			return fmt.Errorf("Task Queue still exists")
		}
	}

	return nil
}
