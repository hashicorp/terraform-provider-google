// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package clouddeploy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccClouddeployCustomTargetType_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployCustomTargetTypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployCustomTargetType_basic(context),
			},
			{
				ResourceName:            "google_clouddeploy_custom_target_type.custom-target-type",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccClouddeployCustomTargetType_update(context),
			},
			{
				ResourceName:            "google_clouddeploy_custom_target_type.custom-target-type",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "annotations", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccClouddeployCustomTargetType_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_custom_target_type" "custom-target-type" {
    location = "us-central1"
    name = "tf-test-my-custom-target-type%{random_suffix}"
    description = "My custom target type"
    custom_actions {
      render_action = "renderAction"
      deploy_action = "deployAction"
    }
}
`, context)
}

func testAccClouddeployCustomTargetType_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_custom_target_type" "custom-target-type" {
    location = "us-central1"
    name = "tf-test-my-custom-target-type%{random_suffix}"
    description = "My custom target type"
    custom_actions {
      render_action = "renderAction"
      deploy_action = "deployAction"
      include_skaffold_modules {
        configs = ["my-config"]
        google_cloud_storage {
          source = "gs://example-bucket/dir/configs/*"
          path = "skaffold.yaml"
        }
      }
      include_skaffold_modules {
        configs = ["my-config2"]
        git {
          repo = "http://github.com/example/example-repo.git"
          path = "configs/skaffold.yaml"
          ref = "main"
        }
      }
      include_skaffold_modules {
        configs = ["my-config3"]
        google_cloud_build_repo {
          repository = "projects/example/locations/us-central1/connections/git/repositories/example-repo"
          path = "configs/skaffold.yaml"
          ref = "main"
        }
      }
    }
}
`, context)
}
