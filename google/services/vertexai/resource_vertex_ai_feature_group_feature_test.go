// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIFeatureGroupFeature_updated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIFeatureGroupFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIFeatureGroupFeature_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_group_feature.feature_group_feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"feature_group", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIFeatureGroupFeature_updated(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_group_feature.feature_group_feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"feature_group", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIFeatureGroupFeature_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_feature_group_feature" "feature_group_feature" {
  name = "tf_test_example_feature%{random_suffix}"
  region = "us-central1"
  feature_group = google_vertex_ai_feature_group.sample_feature_group.name
  description = "A sample feature"
  labels = {
      label-one = "value-one"
  }
  version_column_name = "tf_test_example_feature_v1_%{random_suffix}"
}


resource "google_vertex_ai_feature_group" "sample_feature_group" {
  name = "tf_test_example_feature_group%{random_suffix}"
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

resource "google_bigquery_dataset" "sample_dataset" {
  dataset_id                  = "tf_test_job_load%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_table" "sample_table" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.sample_dataset.dataset_id
  table_id   = "tf_test_job_load%{random_suffix}_table"

  schema = <<EOF
[
    {
        "name": "feature_id",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "tf_test_example_feature_v1_%{random_suffix}",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "tf_test_example_feature_v2_%{random_suffix}",
        "type": "STRING",
        "mode": "NULLABLE"
    },
	{
        "name": "feature_timestamp",
        "type": "TIMESTAMP",
        "mode": "NULLABLE"
    }
]
EOF
}
`, context)
}

func testAccVertexAIFeatureGroupFeature_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_feature_group_feature" "feature_group_feature" {
  name = "tf_test_example_feature%{random_suffix}"
  region = "us-central1"
  feature_group = google_vertex_ai_feature_group.sample_feature_group.name
  description = "A sample feature (updated)"
  labels = {
      label-one = "value-one"
	  label-two = "value-two"
  }
  version_column_name = "tf_test_example_feature_v2_%{random_suffix}"
}


resource "google_vertex_ai_feature_group" "sample_feature_group" {
  name = "tf_test_example_feature_group%{random_suffix}"
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

resource "google_bigquery_dataset" "sample_dataset" {
  dataset_id                  = "tf_test_job_load%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_table" "sample_table" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.sample_dataset.dataset_id
  table_id   = "tf_test_job_load%{random_suffix}_table"

  schema = <<EOF
[
    {
        "name": "feature_id",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "tf_test_example_feature_v1_%{random_suffix}",
        "type": "STRING",
        "mode": "NULLABLE"
    },
    {
        "name": "tf_test_example_feature_v2_%{random_suffix}",
        "type": "STRING",
        "mode": "NULLABLE"
    },
	{
        "name": "feature_timestamp",
        "type": "TIMESTAMP",
        "mode": "NULLABLE"
    }
]
EOF
}
`, context)
}
