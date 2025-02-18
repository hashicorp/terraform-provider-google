// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/cloudtasks/Queue.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/examples/base_configs/iam_test_file.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package cloudtasks_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCloudTasksQueueIamBindingGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueueIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_cloud_tasks_queue_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/queues/%s roles/viewer", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-cloud-tasks-queue-test%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccCloudTasksQueueIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_cloud_tasks_queue_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/queues/%s roles/viewer", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-cloud-tasks-queue-test%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudTasksQueueIamMemberGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccCloudTasksQueueIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_cloud_tasks_queue_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/queues/%s roles/viewer user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-cloud-tasks-queue-test%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudTasksQueueIamPolicyGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTasksQueueIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_cloud_tasks_queue_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_cloud_tasks_queue_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/queues/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-cloud-tasks-queue-test%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudTasksQueueIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_cloud_tasks_queue_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/queues/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), fmt.Sprintf("tf-test-cloud-tasks-queue-test%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudTasksQueueIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "tf-test-cloud-tasks-queue-test%{random_suffix}"
  location = "us-central1"
}

resource "google_cloud_tasks_queue_iam_member" "foo" {
  project = google_cloud_tasks_queue.default.project
  location = google_cloud_tasks_queue.default.location
  name = google_cloud_tasks_queue.default.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccCloudTasksQueueIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "tf-test-cloud-tasks-queue-test%{random_suffix}"
  location = "us-central1"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_cloud_tasks_queue_iam_policy" "foo" {
  project = google_cloud_tasks_queue.default.project
  location = google_cloud_tasks_queue.default.location
  name = google_cloud_tasks_queue.default.name
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_cloud_tasks_queue_iam_policy" "foo" {
  project = google_cloud_tasks_queue.default.project
  location = google_cloud_tasks_queue.default.location
  name = google_cloud_tasks_queue.default.name
  depends_on = [
    google_cloud_tasks_queue_iam_policy.foo
  ]
}
`, context)
}

func testAccCloudTasksQueueIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "tf-test-cloud-tasks-queue-test%{random_suffix}"
  location = "us-central1"
}

data "google_iam_policy" "foo" {
}

resource "google_cloud_tasks_queue_iam_policy" "foo" {
  project = google_cloud_tasks_queue.default.project
  location = google_cloud_tasks_queue.default.location
  name = google_cloud_tasks_queue.default.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccCloudTasksQueueIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "tf-test-cloud-tasks-queue-test%{random_suffix}"
  location = "us-central1"
}

resource "google_cloud_tasks_queue_iam_binding" "foo" {
  project = google_cloud_tasks_queue.default.project
  location = google_cloud_tasks_queue.default.location
  name = google_cloud_tasks_queue.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccCloudTasksQueueIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_tasks_queue" "default" {
  name = "tf-test-cloud-tasks-queue-test%{random_suffix}"
  location = "us-central1"
}

resource "google_cloud_tasks_queue_iam_binding" "foo" {
  project = google_cloud_tasks_queue.default.project
  location = google_cloud_tasks_queue.default.location
  name = google_cloud_tasks_queue.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
