// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexDedicatedResourcesExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexDedicatedResourcesBasicExample_full(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint_deployed_index.dedicated_resources_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"index_endpoint"},
			},
			{
				Config: testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexDedicatedResourcesBasicExample_update(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint_deployed_index.dedicated_resources_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"index_endpoint"},
			},
		},
	})
}

func testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexDedicatedResourcesBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_index_endpoint_deployed_index" "dedicated_resources_update" {
  depends_on = [ google_vertex_ai_index_endpoint.vertex_endpoint ]
  index_endpoint = google_vertex_ai_index_endpoint.vertex_endpoint.id
  index = google_vertex_ai_index.index.id // this is the index that will be deployed onto an endpoint
  deployed_index_id = "tf_test_deployed_index_id%{random_suffix}"
  display_name = "tf-test-vertex-deployed-index%{random_suffix}"
  dedicated_resources {
    machine_spec {
      machine_type      = "n1-standard-32"
    }
    max_replica_count = 2
    min_replica_count = 1
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "tf-test-bucket-name%{random_suffix}"
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
      shard_size = "SHARD_SIZE_MEDIUM"
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


resource "google_vertex_ai_index_endpoint" "vertex_endpoint" {
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  public_endpoint_enabled = true
}

data "google_project" "project" {}
`, context)
}

func testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexDedicatedResourcesBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_index_endpoint_deployed_index" "dedicated_resources_update" {
  depends_on = [ google_vertex_ai_index_endpoint.vertex_endpoint ]
  index_endpoint = google_vertex_ai_index_endpoint.vertex_endpoint.id
  index = google_vertex_ai_index.index.id // this is the index that will be deployed onto an endpoint
  deployed_index_id = "tf_test_deployed_index_id%{random_suffix}"
  display_name = "tf-test-vertex-deployed-index%{random_suffix}"
  dedicated_resources {
    machine_spec {
      machine_type      = "n1-standard-32"
    }
    max_replica_count = 3
    min_replica_count = 2
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "tf-test-bucket-name%{random_suffix}"
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
      shard_size = "SHARD_SIZE_MEDIUM"
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


resource "google_vertex_ai_index_endpoint" "vertex_endpoint" {
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  public_endpoint_enabled = true
}

data "google_project" "project" {}
`, context)
}

func TestAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointAutomaticResourcesExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexAutomaticResourcesBasicExample_full(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint_deployed_index.automatic_resources_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"index_endpoint"},
			},
			{
				Config: testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexAutomaticResourcesBasicExample_update(context),
			},
			{
				ResourceName:            "google_vertex_ai_index_endpoint_deployed_index.automatic_resources_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"index_endpoint"},
			},
		},
	})
}

func testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexAutomaticResourcesBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_index_endpoint_deployed_index" "automatic_resources_update" {
  depends_on = [ google_vertex_ai_index_endpoint.vertex_endpoint ]
  index_endpoint = google_vertex_ai_index_endpoint.vertex_endpoint.id
  index = google_vertex_ai_index.index.id // this is the index that will be deployed onto an endpoint
  deployed_index_id = "tf_test_deployed_index_id%{random_suffix}"
  display_name = "tf-test-vertex-deployed-index%{random_suffix}"
  automatic_resources {
	min_replica_count = 3
	max_replica_count = 4
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "tf-test-bucket-name%{random_suffix}"
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
      shard_size = "SHARD_SIZE_MEDIUM"
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


resource "google_vertex_ai_index_endpoint" "vertex_endpoint" {
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  public_endpoint_enabled = true
}

data "google_project" "project" {}
`, context)
}

func testAccVertexAIIndexEndpointDeployedIndex_vertexAiIndexEndpointDeployedIndexAutomaticResourcesBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_index_endpoint_deployed_index" "automatic_resources_update" {
  depends_on = [ google_vertex_ai_index_endpoint.vertex_endpoint ]
  index_endpoint = google_vertex_ai_index_endpoint.vertex_endpoint.id
  index = google_vertex_ai_index.index.id // this is the index that will be deployed onto an endpoint
  deployed_index_id = "tf_test_deployed_index_id%{random_suffix}"
  display_name = "tf-test-vertex-deployed-index%{random_suffix}"
  automatic_resources {
	min_replica_count = 5
	max_replica_count = 6
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "tf-test-bucket-name%{random_suffix}"
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
      shard_size = "SHARD_SIZE_MEDIUM"
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


resource "google_vertex_ai_index_endpoint" "vertex_endpoint" {
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  public_endpoint_enabled = true
}

data "google_project" "project" {}
`, context)
}
