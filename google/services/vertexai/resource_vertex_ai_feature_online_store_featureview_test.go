// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_updated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIFeatureOnlineStoreFeatureviewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_online_store_featureview.featureview",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "feature_online_store", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_update(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_online_store_featureview.featureview",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "feature_online_store", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_vertex_ai_feature_online_store" "featureonlinestore" {
    name = "tf_test_featureonlinestore%{random_suffix}"
    labels = {
      foo = "bar"
    }
    region = "us-central1"
    bigtable {
      auto_scaling {
        min_node_count         = 1
        max_node_count         = 2
        cpu_utilization_target = 80
      }
    }
  }
  
  resource "google_bigquery_dataset" "tf-test-dataset" {
  
    dataset_id    = "tf_test_dataset1_featureview%{random_suffix}"
    friendly_name = "test"
    description   = "This is a test description"
    location      = "US"
  }
  
  resource "google_bigquery_table" "tf-test-table" {
    deletion_protection = false
  
    dataset_id = google_bigquery_dataset.tf-test-dataset.dataset_id
    table_id   = "tf_test_bq_table%{random_suffix}"
    schema     = <<EOF
      [
      {
        "name": "entity_id",
        "mode": "NULLABLE",
        "type": "STRING",
        "description": "Test default entity_id"
      },
        {
        "name": "test_entity_column",
        "mode": "NULLABLE",
        "type": "STRING",
        "description": "test secondary entity column"
      },
      {
        "name": "feature_timestamp",
        "mode": "NULLABLE",
        "type": "TIMESTAMP",
        "description": "Default timestamp value"
      }
    ]
    EOF
  }
  
  resource "google_vertex_ai_feature_online_store_featureview" "featureview" {
    name   = "tf_test_fv%{random_suffix}"
    region = "us-central1"
    labels = {
      foo = "bar"
    }
    feature_online_store = google_vertex_ai_feature_online_store.featureonlinestore.name
    sync_config {
      cron = "0 0 * * *"
    }
    big_query_source {
      uri               = "bq://${google_bigquery_table.tf-test-table.project}.${google_bigquery_table.tf-test-table.dataset_id}.${google_bigquery_table.tf-test-table.table_id}"
      entity_id_columns = ["test_entity_column"]
  
    }
  }
  
  data "google_project" "project" {
    provider = google
  }  
`, context)
}

func testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_vertex_ai_feature_online_store" "featureonlinestore" {
    name = "tf_test_featureonlinestore%{random_suffix}"
    labels = {
      foo = "bar"
    }
    region = "us-central1"
    bigtable {
      auto_scaling {
        min_node_count         = 1
        max_node_count         = 2
        cpu_utilization_target = 80
      }
    }
  }
  
  resource "google_bigquery_dataset" "tf-test-dataset" {
  
    dataset_id    = "tf_test_dataset1_featureview%{random_suffix}"
    friendly_name = "test"
    description   = "This is a test description"
    location      = "US"
  }
  
  resource "google_bigquery_table" "tf-test-table" {
    deletion_protection = false
  
    dataset_id = google_bigquery_dataset.tf-test-dataset.dataset_id
    table_id   = "tf_test_bq_table%{random_suffix}"
    schema     = <<EOF
  [
  {
  "name": "entity_id",
  "mode": "NULLABLE",
  "type": "STRING",
  "description": "Test default entity_id"
  },
  {
  "name": "test_entity_column",
  "mode": "NULLABLE",
  "type": "STRING",
  "description": "test secondary entity column"
  },
  {
  "name": "feature_timestamp",
  "mode": "NULLABLE",
  "type": "TIMESTAMP",
  "description": "Default timestamp value"
  }
  ]
  EOF
  }
  
  resource "google_vertex_ai_feature_online_store_featureview" "featureview" {
    name   = "tf_test_fv%{random_suffix}"
    region = "us-central1"
    labels = {
      foo1 = "bar1"
    }
    feature_online_store = google_vertex_ai_feature_online_store.featureonlinestore.name
    sync_config {
      cron = "0 4 * * *"
    }
    big_query_source {
      uri               = "bq://${google_bigquery_table.tf-test-table.project}.${google_bigquery_table.tf-test-table.dataset_id}.${google_bigquery_table.tf-test-table.table_id}"
      entity_id_columns = ["test_entity_column"]
  
    }
  }
  
  data "google_project" "project" {
    provider = google
  }
`, context)
}
