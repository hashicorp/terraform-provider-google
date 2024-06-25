// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudscheduler_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/cloudscheduler"
)

func TestAccCloudSchedulerJob_schedulerPausedExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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

func TestUnitCloudSchedulerJob_LastSlashDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"slash to no slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app/",
			New:                "https://hello-rehvs75zla-uc.a.run.app",
			ExpectDiffSuppress: true,
		},
		"no slash to slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app",
			New:                "https://hello-rehvs75zla-uc.a.run.app/",
			ExpectDiffSuppress: true,
		},
		"slash to slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app/",
			New:                "https://hello-rehvs75zla-uc.a.run.app/",
			ExpectDiffSuppress: true,
		},
		"no slash to no slash": {
			Old:                "https://hello-rehvs75zla-uc.a.run.app",
			New:                "https://hello-rehvs75zla-uc.a.run.app",
			ExpectDiffSuppress: true,
		},
		"different domains": {
			Old:                "https://x.a.run.app/",
			New:                "https://y.a.run.app",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if cloudscheduler.LastSlashDiffSuppress("uri", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func testAccCloudSchedulerJob_schedulerPaused(context map[string]interface{}) string {
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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
