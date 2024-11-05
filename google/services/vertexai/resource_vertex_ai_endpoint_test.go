// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAIEndpoint_vertexAiEndpointNetwork(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"endpoint_name": fmt.Sprint(acctest.RandInt(t) % 9999999999),
		"kms_key_name":  acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "vertex-ai-endpoint-update-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpoint_vertexAiEndpointNetwork(context),
			},
			{
				ResourceName:            "google_vertex_ai_endpoint.endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "location", "region", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIEndpoint_vertexAiEndpointNetworkUpdate(context),
			},
			{
				ResourceName:            "google_vertex_ai_endpoint.endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "location", "region", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIEndpoint_vertexAiEndpointNetwork(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "%{endpoint_name}"
  display_name = "sample-endpoint"
  description  = "A sample vertex endpoint"
  location     = "us-central1"
  region       = "us-central1"
  labels       = {
    label-one = "value-one"
  }
  network      = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.vertex_network.name}"
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
  predict_request_response_logging_config {
    bigquery_destination {
      output_uri = "bq://${data.google_project.project.project_id}.${google_bigquery_dataset.bq_dataset.dataset_id}.request_response_logging"
    }
    enabled       = true
    sampling_rate = 0.1
  }

  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

data "google_compute_network" "vertex_network" {
  name       = "%{network_name}"
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "bq_dataset" {
  dataset_id                 = "some_dataset%{endpoint_name}"
  friendly_name              = "logging dataset"
  description                = "This is a dataset that requests are logged to"
  location                   = "US"
  delete_contents_on_destroy = true
}

data "google_project" "project" {}
`, context)
}

func testAccVertexAIEndpoint_vertexAiEndpointNetworkUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint" "endpoint" {
  name         = "%{endpoint_name}"
  display_name = "new-sample-endpoint"
  description  = "An updated sample vertex endpoint"
  location     = "us-central1"
  region       = "us-central1"
  labels       = {
    label-two = "value-two"
  }
  network      = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.vertex_network.name}"
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }

  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

data "google_compute_network" "vertex_network" {
  name       = "%{network_name}"
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com"
}

data "google_project" "project" {}
`, context)
}
