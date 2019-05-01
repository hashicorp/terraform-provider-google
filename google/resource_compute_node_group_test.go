package google

import (
	"fmt"
	"testing"

	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeNodeGroup_updateNodeTemplate(t *testing.T) {
	t.Parallel()

	groupName := acctest.RandomWithPrefix("group-")
	tmplPrefix := acctest.RandomWithPrefix("tmpl-")

	var timeCreated time.Time
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNodeGroup_updateNodeTemplate(groupName, tmplPrefix, "tmpl1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNodeGroupCreationTimeBefore(&timeCreated),
				),
			},
			{
				ResourceName:      "google_compute_node_group.nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeNodeGroup_updateNodeTemplate(groupName, tmplPrefix, "tmpl2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNodeGroupCreationTimeBefore(&timeCreated),
				),
			},
			{
				ResourceName:      "google_compute_node_group.nodes",
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccComputeNodeGroup_updateNodeTemplate(groupName, tmplPrefix, tmplToUse string) string {
	return fmt.Sprintf(`
data "google_compute_node_types" "central1a" {
  zone = "us-central1-a"
}

resource "google_compute_node_template" "tmpl1" {
  name = "%s-first"
  region = "us-central1"
  node_type = "${data.google_compute_node_types.central1a.names[0]}"
}

resource "google_compute_node_template" "tmpl2" {
  name = "%s-second"
  region = "us-central1"
  node_type = "${data.google_compute_node_types.central1a.names[0]}"
}

resource "google_compute_node_group" "nodes" {
  name = "%s"
  zone = "us-central1-a"
  description = "example google_compute_node_group for Terraform Google Provider"

  size = 0
  node_template = "${google_compute_node_template.%s.self_link}"
}
`, tmplPrefix, tmplPrefix, groupName, tmplToUse)
}
