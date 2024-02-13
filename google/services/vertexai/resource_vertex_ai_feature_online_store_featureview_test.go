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

func TestAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_featureRegistry_updated(t *testing.T) {
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
				Config: testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_featureRegistry_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_online_store_featureview.featureregistry_featureview",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "feature_online_store", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_featureRegistry_update(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_online_store_featureview.featureregistry_featureview",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "feature_online_store", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_featureRegistry_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_vertex_ai_feature_online_store" "featureregistry_featureonlinestore" {
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
  
  resource "google_bigquery_dataset" "featureregistry-tf-test-dataset" {
  
    dataset_id    = "tf_test_dataset1_featureview%{random_suffix}"
    friendly_name = "test"
    description   = "This is a test description"
    location      = "US"
  }
  
  resource "google_bigquery_table" "sample_table" {
    deletion_protection = false
  
    dataset_id = google_bigquery_dataset.featureregistry-tf-test-dataset.dataset_id
    table_id   = "tf_test_bq_table%{random_suffix}"
    schema     = <<EOF
      [
        {
          "name": "feature_id",
          "type": "STRING",
          "mode": "NULLABLE"
      },
      {
        "name": "feature_id_updated",
        "type": "STRING",
        "mode": "NULLABLE"
    },
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

  resource "google_vertex_ai_feature_group" "sample_feature_group" {
    name = "tf_test_feature_group%{random_suffix}"
    description = "A sample feature group"
    region = "us-central1"
    labels = {
        label-one = "value-one"
    }
    big_query {
      big_query_source {
          # The source table must have a column named 'feature_timestamp' of type TIMESTAMP.
          input_uri = "bq://${google_bigquery_table.sample_table.project}.${google_bigquery_table.sample_table.dataset_id}.${google_bigquery_table.sample_table.table_id}"
      }
      entity_id_columns = ["feature_id"]
    }
  }
  
  
  
  resource "google_vertex_ai_feature_group_feature" "sample_feature" {
    name = "feature_id"
    region = "us-central1"
    feature_group = google_vertex_ai_feature_group.sample_feature_group.name
    description = "A sample feature"
    labels = {
        label-one = "value-one"
    }
  }  
  resource "google_vertex_ai_feature_group_feature" "updated_feature" {
    name = "feature_id_updated"
    region = "us-central1"
    feature_group = google_vertex_ai_feature_group.sample_feature_group.name
    version_column_name = "feature_id_updated"
    description = "Updated feature"
    labels = {
        label-one = "value-one"
    }
  }
  
  resource "google_vertex_ai_feature_online_store_featureview" "featureregistry_featureview" {
    name   = "tf_test_fv%{random_suffix}"
    region = "us-central1"
    labels = {
      foo = "bar"
    }
    feature_online_store = google_vertex_ai_feature_online_store.featureregistry_featureonlinestore.name
    sync_config {
      cron = "0 0 * * *"
    }
    feature_registry_source {
    
      feature_groups { 
          feature_group_id = google_vertex_ai_feature_group.sample_feature_group.name
          feature_ids      = [google_vertex_ai_feature_group_feature.sample_feature.name]
         }
    }
  }
  
  data "google_project" "project" {
    provider = google
  }  
`, context)
}

func testAccVertexAIFeatureOnlineStoreFeatureview_vertexAiFeatureonlinestoreFeatureview_featureRegistry_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_vertex_ai_feature_online_store" "featureregistry_featureonlinestore" {
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
  
  resource "google_bigquery_dataset" "featureregistry-tf-test-dataset" {
  
    dataset_id    = "tf_test_dataset1_featureview%{random_suffix}"
    friendly_name = "test"
    description   = "This is a test description"
    location      = "US"
  }
  
  resource "google_bigquery_table" "sample_table" {
    deletion_protection = false
  
    dataset_id = google_bigquery_dataset.featureregistry-tf-test-dataset.dataset_id
    table_id   = "tf_test_bq_table%{random_suffix}"
    schema     = <<EOF
      [
        {
          "name": "feature_id",
          "type": "STRING",
          "mode": "NULLABLE"
      },
      {
        "name": "feature_id_updated",
        "type": "STRING",
        "mode": "NULLABLE"
    },
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

  resource "google_vertex_ai_feature_group" "sample_feature_group" {
    name = "tf_test_feature_group%{random_suffix}"
    description = "A sample feature group"
    region = "us-central1"
    labels = {
        label-one = "value-one"
    }
    big_query {
      big_query_source {
          # The source table must have a column named 'feature_timestamp' of type TIMESTAMP.
          input_uri = "bq://${google_bigquery_table.sample_table.project}.${google_bigquery_table.sample_table.dataset_id}.${google_bigquery_table.sample_table.table_id}"
      }
      entity_id_columns = ["feature_id"]
    }
  }
  
  
  
  resource "google_vertex_ai_feature_group_feature" "sample_feature" {
    name = "feature_id"
    region = "us-central1"
    feature_group = google_vertex_ai_feature_group.sample_feature_group.name
    description = "A sample feature"
    labels = {
        label-one = "value-one"
    }
  }  
  resource "google_vertex_ai_feature_group_feature" "updated_feature" {
    name = "feature_id_updated"
    region = "us-central1"
    feature_group = google_vertex_ai_feature_group.sample_feature_group.name
    version_column_name = "feature_id_updated"
    description = "Updated feature"
    labels = {
        label-one = "value-one"
    }
  }
  
  resource "google_vertex_ai_feature_online_store_featureview" "featureregistry_featureview" {
    name   = "tf_test_fv%{random_suffix}"
    region = "us-central1"
    labels = {
      foo = "bar"
    }
    feature_online_store = google_vertex_ai_feature_online_store.featureregistry_featureonlinestore.name
    sync_config {
      cron = "0 0 * * *"
    }
    feature_registry_source {
    
      feature_groups { 
          feature_group_id = google_vertex_ai_feature_group.sample_feature_group.name
          feature_ids      = [google_vertex_ai_feature_group_feature.updated_feature.name]
         }
    }
  }
  
  data "google_project" "project" {
    provider = google
  }
`, context)
}
