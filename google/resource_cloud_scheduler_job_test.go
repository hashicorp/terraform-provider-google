package google

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestCloudScheduler_FlattenHttpHeaders(t *testing.T) {

	cases := []struct {
		Input  map[string]interface{}
		Output map[string]interface{}
	}{
		// simple, no headers included
		{
			Input: map[string]interface{}{
				"My-Header": "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the User-Agent header value Google-Cloud-Scheduler
		// Tests Removing User-Agent header
		{
			Input: map[string]interface{}{
				"User-Agent": "Google-Cloud-Scheduler",
				"My-Header":  "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the User-Agent header
		// Tests removing value AppEngine-Google; (+http://code.google.com/appengine)
		{
			Input: map[string]interface{}{
				"User-Agent": "My-User-Agent AppEngine-Google; (+http://code.google.com/appengine)",
				"My-Header":  "my-header-value",
			},
			Output: map[string]interface{}{
				"User-Agent": "My-User-Agent",
				"My-Header":  "my-header-value",
			},
		},

		// include the Content-Type header value application/octet-stream.
		// Tests Removing Content-Type header
		{
			Input: map[string]interface{}{
				"Content-Type": "application/octet-stream",
				"My-Header":    "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the Content-Length header
		// Tests Removing Content-Length header
		{
			Input: map[string]interface{}{
				"Content-Length": 7,
				"My-Header":      "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the X-Google- header
		// Tests Removing X-Google- header
		{
			Input: map[string]interface{}{
				"X-Google-My-Header": "x-google-my-header-value",
				"My-Header":          "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},
	}

	for _, c := range cases {
		d := &schema.ResourceData{}
		output := flattenCloudSchedulerJobAppEngineHttpTargetHeaders(c.Input, d, &transport_tpg.Config{})
		if !reflect.DeepEqual(output, c.Output) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", output, c.Output)
		}
	}
}

func TestAccCloudSchedulerJob_schedulerPausedExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudSchedulerJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSchedulerJob_schedulerPaused(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_scheduler_job.job", "paused", "true"),
					resource.TestCheckResourceAttr("google_cloud_scheduler_job.job", "state", "PAUSED"),
				),
			},
			{
				Config: testAccCloudSchedulerJob_schedulerUnPaused(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_scheduler_job.job", "paused", "false"),
					resource.TestCheckResourceAttr("google_cloud_scheduler_job.job", "state", "ENABLED"),
				),
			},
		},
	})
}

func testAccCloudSchedulerJob_schedulerPaused(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_scheduler_job" "job" {
  paused           = true
  name             = "tf-test-test-job%{random_suffix}"
  description      = "test http job with updated fields"
  schedule         = "*/8 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "https://example.com/ping"
    body        = base64encode("{\"foo\":\"bar\"}")
  }
}
`, context)
}

func testAccCloudSchedulerJob_schedulerUnPaused(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_scheduler_job" "job" {
  paused           = false # Has been flipped 
  name             = "tf-test-test-job%{random_suffix}"
  description      = "test http job with updated fields"
  schedule         = "*/8 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "https://example.com/ping"
    body        = base64encode("{\"foo\":\"bar\"}")
  }
}
`, context)
}
