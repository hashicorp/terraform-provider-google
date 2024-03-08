// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigQueryRoutine_bigQueryRoutine_Update(t *testing.T) {
	t.Parallel()

	dataset := fmt.Sprintf("tfmanualdataset%s", acctest.RandString(t, 10))
	routine := fmt.Sprintf("tfmanualroutine%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryRoutineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryRoutine_bigQueryRoutine(dataset, routine),
			},
			{
				ResourceName:      "google_bigquery_routine.sproc",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryRoutine_bigQueryRoutine_Update(dataset, routine),
			},
			{
				ResourceName:      "google_bigquery_routine.sproc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigQueryRoutine_bigQueryRoutine(dataset, routine string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_routine" "sproc" {
  dataset_id = google_bigquery_dataset.test.dataset_id
  routine_id     = "%s"
  routine_type = "SCALAR_FUNCTION"
  language = "SQL"
  definition_body = "1"
}
`, dataset, routine)
}

func testAccBigQueryRoutine_bigQueryRoutine_Update(dataset, routine string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_routine" "sproc" {
  dataset_id = google_bigquery_dataset.test.dataset_id
  routine_id     = "%s"
  routine_type = "SCALAR_FUNCTION"
  language = "JAVASCRIPT"
  definition_body = "CREATE FUNCTION multiplyInputs return x*y;"
  arguments {
    name = "x"
    data_type = "{\"typeKind\" :  \"FLOAT64\"}"
  }
  arguments {
    name = "y"
    data_type = "{\"typeKind\" :  \"FLOAT64\"}"
  }

  return_type = "{\"typeKind\" :  \"FLOAT64\"}"
}
`, dataset, routine)
}

func TestAccBigQueryRoutine_bigQueryRoutineSparkJar_Update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryRoutineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryRoutine_bigQueryRoutineSparkJar(context),
			},
			{
				ResourceName:      "google_bigquery_routine.spark_jar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryRoutine_bigQueryRoutineSparkJar_Update(context),
			},
			{
				ResourceName:      "google_bigquery_routine.spark_jar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigQueryRoutine_bigQueryRoutineSparkJar(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "tf_test_dataset_id%{random_suffix}"
}

resource "google_bigquery_connection" "test" {
  connection_id = "tf_test_connection_id%{random_suffix}"
  location      = "US"
  spark { }
}

resource "google_bigquery_routine" "spark_jar" {
  dataset_id      = google_bigquery_dataset.test.dataset_id
  routine_id      = "tf_test_routine_id%{random_suffix}"
  routine_type    = "PROCEDURE"
  language        = "SCALA"
  definition_body = ""
  spark_options {
    connection      = google_bigquery_connection.test.name
    runtime_version = "2.0"
    main_class      = "com.google.test.jar.MainClass"
    jar_uris        = [ "gs://test-bucket/testjar_spark_spark3.jar" ]
  }
}
`, context)
}

func testAccBigQueryRoutine_bigQueryRoutineSparkJar_Update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "tf_test_dataset_id%{random_suffix}"
}

resource "google_bigquery_connection" "test_updated" {
  connection_id = "tf_test_connection_updated_id%{random_suffix}"
  location      = "US"
  spark { }
}

resource "google_bigquery_routine" "spark_jar" {
  dataset_id      = google_bigquery_dataset.test.dataset_id
  routine_id      = "tf_test_routine_id%{random_suffix}"
  routine_type    = "PROCEDURE"
  language        = "SCALA"
  definition_body = ""
  spark_options {
    connection      = google_bigquery_connection.test_updated.name
    runtime_version = "2.1"
    container_image = "gcr.io/my-project-id/my-spark-image:latest"
    main_class      = "com.google.test.jar.MainClassUpdated"
    jar_uris        = [ "gs://test-bucket/uberjar_spark_spark3_updated.jar" ]
    properties      = {
      "spark.dataproc.scaling.version" : "2",
      "spark.reducer.fetchMigratedShuffle.enabled" : "true",
    }
  }
}
`, context)
}

func TestAccBigQueryRoutine_bigQueryRoutineRemoteFunction(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"zip_path":      "./test-fixtures/function-source.zip",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryRoutineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryRoutine_bigQueryRoutineRemoteFunction(context),
			},
			{
				ResourceName:      "google_bigquery_routine.remote_function_routine",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryRoutine_bigQueryRoutineRemoteFunction_Update(context),
			},
			{
				ResourceName:      "google_bigquery_routine.remote_function_routine",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigQueryRoutine_bigQueryRoutineRemoteFunction(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "%{random_suffix}-gcf-source"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.default.name
  source = "%{zip_path}"
}

resource "google_cloudfunctions2_function" "default" {
  name        = "function-v2-0"
  location    = "us-central1"
  description = "a new function"

  build_config {
    runtime     = "nodejs18"
    entry_point = "helloHttp"
    source {
      storage_source {
        bucket = google_storage_bucket.default.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count = 1
    available_memory   = "256M"
    timeout_seconds    = 60
  }
}

resource "google_bigquery_connection" "test" {
  connection_id = "tf_test_connection_id%{random_suffix}"
  location      = "US"
  cloud_resource { }
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "tf_test_dataset_id%{random_suffix}"
}

resource "google_bigquery_routine" "remote_function_routine" {
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"
  routine_id = "tf_test_routine_id%{random_suffix}"
  routine_type = "SCALAR_FUNCTION"
  definition_body = ""

  return_type = "{\"typeKind\" :  \"STRING\"}"

  remote_function_options {
    endpoint = google_cloudfunctions2_function.default.service_config[0].uri
    connection = "${google_bigquery_connection.test.name}"
    max_batching_rows = "10"
    user_defined_context = {
      "z": "1.5",
    }
  }
}
`, context)
}

func testAccBigQueryRoutine_bigQueryRoutineRemoteFunction_Update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "%{random_suffix}-gcf-source"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.default.name
  source = "%{zip_path}"
}

resource "google_cloudfunctions2_function" "default2" {
  name        = "function-v2-1"
  location    = "us-central1"
  description = "a new new function"

  build_config {
    runtime     = "nodejs18"
    entry_point = "helloHttp"
    source {
      storage_source {
        bucket = google_storage_bucket.default.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count = 1
    available_memory   = "256M"
    timeout_seconds    = 60
  }
}

resource "google_bigquery_connection" "test2" {
  connection_id = "tf_test_connection2_id%{random_suffix}"
  location      = "US"
  cloud_resource { }
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "tf_test_dataset_id%{random_suffix}"
}

resource "google_bigquery_routine" "remote_function_routine" {
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"
  routine_id = "tf_test_routine_id%{random_suffix}"
  routine_type = "SCALAR_FUNCTION"
  definition_body = ""

  return_type = "{\"typeKind\" :  \"STRING\"}"

  remote_function_options {
    endpoint = google_cloudfunctions2_function.default2.service_config[0].uri
    connection = "${google_bigquery_connection.test2.name}"
    max_batching_rows = "5"
    user_defined_context = {
      "z": "1.2",
      "w": "test",
    }
  }
}
`, context)
}
