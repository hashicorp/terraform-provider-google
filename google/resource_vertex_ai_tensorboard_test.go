package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVertexAITensorboard_Update(t *testing.T) {
	t.Parallel()

	random_suffix := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVertexAITensorboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAITensorboard_Update(random_suffix, random_suffix, random_suffix, random_suffix),
			},
			{
				ResourceName:            "google_vertex_ai_tensorboard.tensorboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
			{
				Config: testAccVertexAITensorboard_Update(random_suffix+"new", random_suffix, random_suffix, random_suffix),
			},
			{
				ResourceName:            "google_vertex_ai_tensorboard.tensorboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
			{
				Config: testAccVertexAITensorboard_Update(random_suffix+"new", random_suffix+"new", random_suffix, random_suffix),
			},
			{
				ResourceName:            "google_vertex_ai_tensorboard.tensorboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
			{
				Config: testAccVertexAITensorboard_Update(random_suffix+"new", random_suffix+"new", random_suffix+"new", random_suffix),
			},
			{
				ResourceName:            "google_vertex_ai_tensorboard.tensorboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
			{
				Config: testAccVertexAITensorboard_Update(random_suffix+"new", random_suffix+"new", random_suffix+"new", random_suffix+"new"),
			},
			{
				ResourceName:            "google_vertex_ai_tensorboard.tensorboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
			{
				Config: testAccVertexAITensorboard_Update(random_suffix, random_suffix, random_suffix, random_suffix),
			},
			{
				ResourceName:            "google_vertex_ai_tensorboard.tensorboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "project"},
			},
		},
	})
}

func testAccVertexAITensorboard_Update(displayName, description, labelKey, labelVal string) string {
	return fmt.Sprintf(`
resource "google_vertex_ai_tensorboard" "tensorboard" {
  display_name = "%s"
  description  = "%s"
  labels       = {
    "%s" : "%s",
  }
  region       = "us-central1"
}
`, displayName, description, labelKey, labelVal)
}
