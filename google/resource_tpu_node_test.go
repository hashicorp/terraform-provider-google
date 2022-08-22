package google

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTPUNode_tpuNodeBUpdateTensorFlowVersion(t *testing.T) {
	t.Parallel()

	nodeId := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTPUNodeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTpuNode_tpuNodeTensorFlow(nodeId, 0),
			},
			{
				ResourceName:            "google_tpu_node.tpu",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccTpuNode_tpuNodeTensorFlow(nodeId, 1),
			},
			{
				ResourceName:            "google_tpu_node.tpu",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

func testAccTpuNode_tpuNodeTensorFlow(nodeId string, versionIdx int) string {
	return fmt.Sprintf(`
data "google_tpu_tensorflow_versions" "available" {
}

resource "google_tpu_node" "tpu" {
  name = "%s"
  zone = "us-central1-b"

  accelerator_type   = "v3-8"
  tensorflow_version = data.google_tpu_tensorflow_versions.available.versions[%d]
}
`, nodeId, versionIdx)
}
