// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIFeatureOnlineStore_updated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIFeatureOnlineStoreDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIFeatureOnlineStore_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_online_store.feature_online_store",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "region", "force_destroy", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIFeatureOnlineStore_updated(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_online_store.feature_online_store",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "region", "force_destroy", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIFeatureOnlineStore_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource google_vertex_ai_feature_online_store "feature_online_store" {
    name = "tf_test_feature_online_store%{random_suffix}"
    region = "us-central1"
    labels = {
        label-one = "value-one"
    }

    bigtable {
        auto_scaling {
            min_node_count = 1
            max_node_count = 2
            cpu_utilization_target = 60
        }
    }
  force_destroy = true
}
`, context)
}

func testAccVertexAIFeatureOnlineStore_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource google_vertex_ai_feature_online_store "feature_online_store" {
    name = "tf_test_feature_online_store%{random_suffix}"
    region = "us-central1"
    labels = {
        label-one = "value-one"
		label-two = "value-two"
    }

    bigtable {
        auto_scaling {
            min_node_count = 2
            max_node_count = 3
        }
    }
  force_destroy = true
}
`, context)
}
