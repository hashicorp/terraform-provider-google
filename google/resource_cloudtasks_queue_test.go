package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCloudTasksQueue_basic(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-cloudtasks-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCloudTasksQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueue_basic(queueName),
			},
			{
				ResourceName:      "google_cloudtasks_queue.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudTasksQueue_withParams(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-cloudtasks-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCloudTasksQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueue_withParams(queueName),
			},
			{
				ResourceName:            "google_cloudtasks_queue.fizzbuzz",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_engine_routing_override.0.service", "app_engine_routing_override.0.instance", "app_engine_routing_override.0.version"},
			},
		},
	})
}

func TestAccCloudTasksQueue_update(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-cloudtasks-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCloudTasksQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueue_basic(queueName),
			},
			{
				Config: testAccCloudTasksQueue_withParams(queueName),
			},
		},
	})
}

func TestAccCloudTasksQueue_forceDestroy(t *testing.T) {
	t.Parallel()

	queueName := fmt.Sprintf("tf-test-cloudtasks-queue-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCloudTasksQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueue_basic(queueName),
			},
		},
	})
}

func testAccCloudTasksQueue_basic(name string) string {
	return fmt.Sprintf(`
resource "google_cloudtasks_queue" "fizzbuzz" {
  name = "%s"
  location = "us-central1"
}`, name)
}

func testAccCloudTasksQueue_withParams(name string) string {
	return fmt.Sprintf(`
resource "google_cloudtasks_queue" "fizzbuzz" {
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

func testAccCloudTasksQueueDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudtasks_queue" {
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
