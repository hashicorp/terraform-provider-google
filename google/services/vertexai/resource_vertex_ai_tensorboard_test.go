// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAITensorboard_Update(t *testing.T) {
	t.Parallel()

	random_suffix := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAITensorboardDestroyProducer(t),
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
