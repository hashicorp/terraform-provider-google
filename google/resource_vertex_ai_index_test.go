package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVertexAIIndex_updated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckVertexAIIndexDestroyProducer_basic(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIIndex_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_index.index",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "region", "metadata.0.contents_delta_uri", "metadata.0.is_complete_overwrite"},
			},
			{
				Config: testAccVertexAIIndex_updated(context),
			},
			{
				ResourceName:            "google_vertex_ai_index.index",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "region", "metadata.0.contents_delta_uri", "metadata.0.is_complete_overwrite"},
			},
		},
	})
}

func testAccVertexAIIndex_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-%{random_suffix}"
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

resource "google_storage_bucket_object" "data-v2" {
  name   = "contents-v2/data.json"
  bucket = google_storage_bucket.bucket.name
  content = <<EOF
{"id": "1", "embedding": [0.5, 1.0], "restricts": [{"namespace": "class", "allow": ["cat", "pet"]},{"namespace": "category", "allow": ["feline"]}]}
{"id": "2", "embedding": [0.6, 1.0], "restricts": [{"namespace": "class", "allow": ["dog", "pet"]},{"namespace": "category", "allow": ["canine"]}]}
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
`, context)
}

func testAccVertexAIIndex_updated(context map[string]interface{}) string {
	return Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-%{random_suffix}"
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

resource "google_storage_bucket_object" "data-v2" {
  name   = "contents-v2/data.json"
  bucket = google_storage_bucket.bucket.name
  content = <<EOF
{"id": "1", "embedding": [0.5, 1.0], "restricts": [{"namespace": "class", "allow": ["cat", "pet"]},{"namespace": "category", "allow": ["feline"]}]}
{"id": "2", "embedding": [0.6, 1.0], "restricts": [{"namespace": "class", "allow": ["dog", "pet"]},{"namespace": "category", "allow": ["canine"]}]}
EOF
}


resource "google_vertex_ai_index" "index" {
  labels = {
    foo = "bar"
  }
  region   = "us-central1"
  display_name = "tf-test-test-index%{random_suffix}"
  description = "index for test (updated)"
  metadata {
    contents_delta_uri = "gs://${google_storage_bucket.bucket.name}/contents-v2"
    is_complete_overwrite = true
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
`, context)
}

func testAccCheckVertexAIIndexDestroyProducer_basic(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_index" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{VertexAIBasePath}}projects/{{project}}/locations/{{region}}/indexes/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("VertexAIIndex still exists at %s", url)
			}
		}

		return nil
	}
}
