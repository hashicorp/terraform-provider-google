// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeNodeGroup_update(t *testing.T) {
	t.Parallel()

	groupName := fmt.Sprintf("group--%d", acctest.RandInt(t))
	tmplPrefix := fmt.Sprintf("tmpl--%d", acctest.RandInt(t))

	var timeCreated time.Time
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNodeGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNodeGroup_update(groupName, tmplPrefix, "tmpl1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNodeGroupCreationTimeBefore(&timeCreated),
				),
			},
			{
				ResourceName:            "google_compute_node_group.nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_size"},
			},
			{
				Config: testAccComputeNodeGroup_update2(groupName, tmplPrefix, "tmpl2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNodeGroupCreationTimeBefore(&timeCreated),
				),
			},
			{
				ResourceName:            "google_compute_node_group.nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_size"},
			},
		},
	})
}

func TestAccComputeNodeGroup_fail(t *testing.T) {
	t.Parallel()

	groupName := fmt.Sprintf("group--%d", acctest.RandInt(t))
	tmplPrefix := fmt.Sprintf("tmpl--%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNodeGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeNodeGroup_fail(groupName, tmplPrefix, "tmpl1"),
				ExpectError: regexp.MustCompile("An initial_size or autoscaling_policy must be configured on node group creation."),
			},
		},
	})
}

func testAccCheckComputeNodeGroupCreationTimeBefore(prevTimeCreated *time.Time) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_node_group" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			timestampRaw, ok := rs.Primary.Attributes["creation_timestamp"]
			if !ok {
				return fmt.Errorf("expected creation_timestamp to be set in node group's state")
			}
			creationTimestamp, err := time.Parse(time.RFC3339Nano, timestampRaw)
			if err != nil {
				return fmt.Errorf("unexpected error while parsing creation_timestamp: %v", err)
			}

			if prevTimeCreated.IsZero() {
				*prevTimeCreated = creationTimestamp
				return nil
			}

			if creationTimestamp.After(prevTimeCreated.Add(time.Millisecond * 100)) {
				return fmt.Errorf(
					"Creation timestamp %q was after expected previous time of creation %q",
					timestampRaw, prevTimeCreated.Format(time.RFC3339Nano))
			}
		}

		return nil
	}
}

func testAccComputeNodeGroup_update(groupName, tmplPrefix, tmplToUse string) string {
	return fmt.Sprintf(`
resource "google_compute_node_template" "tmpl1" {
  name      = "%s-first"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_template" "tmpl2" {
  name      = "%s-second"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_group" "nodes" {
  name        = "%s"
  zone        = "us-central1-a"
  description = "example google_compute_node_group for Terraform Google Provider"

  initial_size = 1
  node_template = google_compute_node_template.%s.self_link
}

`, tmplPrefix, tmplPrefix, groupName, tmplToUse)
}

func testAccComputeNodeGroup_update2(groupName, tmplPrefix, tmplToUse string) string {
	return fmt.Sprintf(`
resource "google_compute_node_template" "tmpl1" {
  name      = "%s-first"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_template" "tmpl2" {
  name      = "%s-second"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_group" "nodes" {
  name        = "%s"
  zone        = "us-central1-a"
  description = "example google_compute_node_group for Terraform Google Provider"

  autoscaling_policy {
    mode      = "ONLY_SCALE_OUT"
    min_nodes = 1
    max_nodes = 10
  }
  node_template = google_compute_node_template.%s.self_link
}

`, tmplPrefix, tmplPrefix, groupName, tmplToUse)
}

func testAccComputeNodeGroup_fail(groupName, tmplPrefix, tmplToUse string) string {
	return fmt.Sprintf(`
resource "google_compute_node_template" "tmpl1" {
  name      = "%s-first"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_group" "nodes" {
  name        = "%s"
  zone        = "us-central1-a"
  description = "example google_compute_node_group for Terraform Google Provider"

  node_template = google_compute_node_template.%s.self_link
}

`, tmplPrefix, groupName, tmplToUse)
}
