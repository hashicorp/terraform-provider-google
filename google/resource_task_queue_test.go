package google

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"golang.org/x/net/context"
	"google.golang.org/api/cloudtasks/v2"
)

func TestAccTaskQueue_basic(t *testing.T) {
	t.Parallel()

	var queue cloudtasks.Queue
	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())
	location := "us-central1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_basic(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaskQueueExists(
						"google_task_queue.fizzbuzz", location, queueName, &queue),
					testAccCheckTaskQueueEquals(
						"google_task_queue.fizzbuzz", &queue),
				),
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

	var queue cloudtasks.Queue
	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())
	location := "us-central1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_withParams(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaskQueueExists(
						"google_task_queue.fizzbuzz", location, queueName, &queue),
					testAccCheckTaskQueueEquals(
						"google_task_queue.fizzbuzz", &queue),
				),
			},
			{
				ResourceName:      "google_task_queue.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTaskQueue_update(t *testing.T) {
	t.Parallel()

	var queue cloudtasks.Queue
	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())
	location := "us-central1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_basic(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaskQueueExists(
						"google_task_queue.fizzbuzz", location, queueName, &queue),
					testAccCheckTaskQueueEquals(
						"google_task_queue.fizzbuzz", &queue),
				),
			},
			{
				Config: testAccTaskQueue_withParams(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaskQueueExists(
						"google_task_queue.fizzbuzz", location, queueName, &queue),
					testAccCheckTaskQueueEquals(
						"google_task_queue.fizzbuzz", &queue),
				),
			},
		},
	})
}

func TestAccTaskQueue_forceDestroy(t *testing.T) {
	t.Parallel()

	var queue cloudtasks.Queue
	queueName := fmt.Sprintf("tf-test-task-queue-%d", acctest.RandInt())
	location := "us-central1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTaskQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskQueue_basic(queueName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaskQueueExists(
						"google_task_queue.fizzbuzz", location, queueName, &queue),
				),
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

  rate_limits {
    max_concurrent_dispatches = 1000
    max_dispatches_per_second = 500
  }

  retry {
    max_attempts  = 100
    max_backoff   = "3600s"
    min_backoff   = "0.100s"
    max_doublings = 16
  }
}`, name)
}

func testAccCheckTaskQueueExists(n string, location, queueName string, queue *cloudtasks.Queue) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project_ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		expectedName := fmt.Sprintf("projects/%s/locations/%s/queues/%s", config.Project, location, queueName)

		found, err := config.clientCloudTasks.Projects.Locations.Queues.Get(expectedName).Context(context.Background()).Do()
		if err != nil {
			return err
		}

		if found.Name != expectedName {
			return fmt.Errorf("expected name %s, got %s", expectedName, found.Name)
		}

		*queue = *found
		return nil
	}
}

func testAccCheckTaskQueueEquals(n string, queue *cloudtasks.Queue) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		name := rs.Primary.Attributes["name"]
		location := rs.Primary.Attributes["location"]
		expectedName := fmt.Sprintf("projects/%s/locations/%s/queues/%s", config.Project, location, name)

		maxBurstSize, _ := strconv.ParseInt(rs.Primary.Attributes["rate_limits.0.max_burst_size"], 10, 64)
		maxConcurrentDispatches, _ := strconv.ParseInt(rs.Primary.Attributes["rate_limits.0.max_concurrent_dispatches"], 10, 64)
		maxDispatchesPerSecond, _ := strconv.ParseFloat(rs.Primary.Attributes["rate_limits.0.max_dispatches_per_second"], 64)

		maxBackoff := rs.Primary.Attributes["retry.0.max_backoff"]
		maxAttempts, _ := strconv.ParseInt(rs.Primary.Attributes["retry.0.max_attempts"], 10, 64)
		minBackoff := rs.Primary.Attributes["retry.0.min_backoff"]
		maxDoublings, _ := strconv.ParseInt(rs.Primary.Attributes["retry.0.max_doublings"], 10, 64)

		if expectedName != queue.Name {
			return fmt.Errorf("Error name mismatch, (%s, %s)", expectedName, queue.Name)
		}

		if maxBurstSize != queue.RateLimits.MaxBurstSize {
			return fmt.Errorf("Error max_burst_size mismatch, (%v, %v)", maxBurstSize, queue.RateLimits.MaxBurstSize)
		}
		if maxConcurrentDispatches != queue.RateLimits.MaxConcurrentDispatches {
			return fmt.Errorf("Error max_concurrent_dispatches mismatch, (%v, %v)", maxConcurrentDispatches, queue.RateLimits.MaxConcurrentDispatches)
		}
		if maxDispatchesPerSecond != queue.RateLimits.MaxDispatchesPerSecond {
			return fmt.Errorf("Error max_dispatches_per_second mismatch, (%v, %v)", maxDispatchesPerSecond, queue.RateLimits.MaxDispatchesPerSecond)
		}

		if maxBackoff != queue.RetryConfig.MaxBackoff {
			return fmt.Errorf("Error max_backoff mismatch, (%s, %s)", maxBackoff, queue.RetryConfig.MaxBackoff)
		}
		if maxAttempts != queue.RetryConfig.MaxAttempts {
			return fmt.Errorf("Error max_attempts mismatch, (%v, %v)", maxAttempts, queue.RetryConfig.MaxAttempts)
		}
		if minBackoff != queue.RetryConfig.MinBackoff {
			return fmt.Errorf("Error min_backoff mismatch, (%s, %s)", minBackoff, queue.RetryConfig.MinBackoff)
		}
		if maxDoublings != queue.RetryConfig.MaxDoublings {
			return fmt.Errorf("Error max_doublings mismatch, (%v, %v)", maxDoublings, queue.RetryConfig.MaxDoublings)
		}

		return nil
	}
}

func testAccTaskQueueDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_task_queue" {
			continue
		}

		attributes := rs.Primary.Attributes

		expectedName := fmt.Sprintf("projects/%s/locations/%s/queues/%s", config.Project, attributes["location"], attributes["name"])

		_, err := config.clientCloudTasks.Projects.Locations.Queues.Get(expectedName).Context(context.Background()).Do()
		if err == nil {
			return fmt.Errorf("Task Queue still exists")
		}
	}

	return nil
}
