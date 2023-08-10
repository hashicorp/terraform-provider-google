// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceVertexAIIndex_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIIndexDestroyProducer_basic(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVertexAIIndex_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vertex_ai_index.foo",
						"google_vertex_ai_index.index",
						// The projects.locations.indexes.get doesn't return the following fields
						map[string]struct{}{
							"metadata.0.contents_delta_uri":    {},
							"metadata.0.is_complete_overwrite": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceVertexAIIndex_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-%{project}-vertex-ai-index-%{random_suffix}"  # Every bucket name must be globally unique
  location = "us-central1"
  uniform_bucket_level_access = true
}

# The sample data comes from the following link:
# https://cloud.google.com/vertex-ai/docs/matching-engine/filtering#specify-namespaces-tokens
resource "google_storage_bucket_object" "data" {
  name   = "contents/data.json"
  bucket = google_storage_bucket.bucket.name
  content = <<EOF
{"id": "42", "embedding": [0.5, 1.0], "restricts": [{"namespace": "class", "allow": ["cat", "pet"]},{"namespace": "category", "allow": ["feline"]}]}
{"id": "43", "embedding": [0.6, 1.0], "restricts": [{"namespace": "class", "allow": ["dog", "pet"]},{"namespace": "category", "allow": ["canine"]}]}
EOF
}

resource "google_vertex_ai_index" "index" {
  labels = {
    foo = "bar"
  }
  region   = "us-central1"
  display_name = "tf-test-test-index%{random_suffix}"
  description = "index for test"
  metadata {
    contents_delta_uri = "gs://${google_storage_bucket.bucket.name}/contents"
    config {
      dimensions = 2
      approximate_neighbors_count = 150
      distance_measure_type = "DOT_PRODUCT_DISTANCE"
      algorithm_config {
        tree_ah_config {
          leaf_node_embedding_count = 500
          leaf_nodes_to_search_percent = 7
        }
      }
    }
  }
  index_update_method = "BATCH_UPDATE"
}

data "google_vertex_ai_index" "foo" {
  name = google_vertex_ai_index.index.name
  region = google_vertex_ai_index.index.region
  project = google_vertex_ai_index.index.project
}

`, context)
}
