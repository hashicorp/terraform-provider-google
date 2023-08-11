// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package healthcare_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/healthcare"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccHealthcareFhirStoreIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedFhirStoreId string
		Config              *transport_tpg.Config
	}{
		"id is in project/location/datasetName/fhirStoreName format": {
			ImportId:            "test-project/us-central1/test-dataset/test-store-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-dataset/test-store-name",
			ExpectedFhirStoreId: "projects/test-project/locations/us-central1/datasets/test-dataset/fhirStores/test-store-name",
		},
		"id is in domain:project/location/datasetName/fhirStoreName format": {
			ImportId:            "example.com:test-project/us-central1/test-dataset/test-store-name",
			ExpectedError:       false,
			ExpectedTerraformId: "example.com:test-project/us-central1/test-dataset/test-store-name",
			ExpectedFhirStoreId: "projects/example.com:test-project/locations/us-central1/datasets/test-dataset/fhirStores/test-store-name",
		},
		"id is in location/datasetName/fhirStoreName format": {
			ImportId:            "us-central1/test-dataset/test-store-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-dataset/test-store-name",
			ExpectedFhirStoreId: "projects/test-project/locations/us-central1/datasets/test-dataset/fhirStores/test-store-name",
			Config:              &transport_tpg.Config{Project: "test-project"},
		},
		"id is in location/datasetName/fhirStoreName format without project in config": {
			ImportId:      "us-central1/test-dataset/test-store-name",
			ExpectedError: true,
			Config:        &transport_tpg.Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		fhirStoreId, err := healthcare.ParseHealthcareFhirStoreId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if fhirStoreId.TerraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, fhirStoreId.TerraformId())
		}

		if fhirStoreId.FhirStoreId() != tc.ExpectedFhirStoreId {
			t.Fatalf("bad: %s, expected FhirStore ID to be `%s` but is `%s`", tn, tc.ExpectedFhirStoreId, fhirStoreId.FhirStoreId())
		}
	}
}

func TestAccHealthcareFhirStore_basic(t *testing.T) {
	t.Parallel()

	datasetName := fmt.Sprintf("tf-test-dataset-%s", acctest.RandString(t, 10))
	fhirStoreName := fmt.Sprintf("tf-test-fhir-store-%s", acctest.RandString(t, 10))
	pubsubTopic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	resourceName := "google_healthcare_fhir_store.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcareFhirStoreDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleHealthcareFhirStore_basic(fhirStoreName, datasetName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleHealthcareFhirStore_update(fhirStoreName, datasetName, pubsubTopic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleHealthcareFhirStoreUpdate(t, pubsubTopic),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleHealthcareFhirStore_basic(fhirStoreName, datasetName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleHealthcareFhirStore_basic(fhirStoreName, datasetName string) string {
	return fmt.Sprintf(`
resource "google_healthcare_fhir_store" "default" {
  name     = "%s"
  dataset  = google_healthcare_dataset.dataset.id

  enable_update_create          = false
  disable_referential_integrity = false
  disable_resource_versioning   = false
  enable_history_import         = false
  version                       = "R4"
}

resource "google_healthcare_dataset" "dataset" {
  name     = "%s"
  location = "us-central1"
}
`, fhirStoreName, datasetName)
}

func testGoogleHealthcareFhirStore_update(fhirStoreName, datasetName, pubsubTopic string) string {
	return fmt.Sprintf(`
resource "google_healthcare_fhir_store" "default" {
  name     = "%s"
  dataset  = google_healthcare_dataset.dataset.id

  enable_update_create = true
  version              = "R4"


  notification_config {
    pubsub_topic = google_pubsub_topic.topic.id
  }

  labels = {
    label1 = "labelvalue1"
  }
}

resource "google_healthcare_dataset" "dataset" {
  name     = "%s"
  location = "us-central1"
}

resource "google_pubsub_topic" "topic" {
  name = "%s"
}
`, fhirStoreName, datasetName, pubsubTopic)
}

func testAccCheckGoogleHealthcareFhirStoreUpdate(t *testing.T, pubsubTopic string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var foundResource = false
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_healthcare_fhir_store" {
				continue
			}
			foundResource = true

			config := acctest.GoogleProviderConfig(t)

			gcpResourceUri, err := tpgresource.ReplaceVarsForTest(config, rs, "{{dataset}}/fhirStores/{{name}}")
			if err != nil {
				return err
			}

			response, err := config.NewHealthcareClient(config.UserAgent).Projects.Locations.Datasets.FhirStores.Get(gcpResourceUri).Do()
			if err != nil {
				return fmt.Errorf("Unexpected failure while verifying 'updated' dataset: %s", err)
			}

			if !response.EnableUpdateCreate {
				return fmt.Errorf("fhirStore 'EnableUpdateCreate' not updated: %s", gcpResourceUri)
			}

			// because the GET for the FHIR store resource does not return the "enableHistoryImport" flag, this value
			// will always be false and cannot be relied upon

			//if !response.EnableHistoryImport {
			//	return fmt.Errorf("fhirStore 'EnableHistoryImport' not updated: %s", gcpResourceUri)
			//}

			if len(response.Labels) == 0 || response.Labels["label1"] != "labelvalue1" {
				return fmt.Errorf("fhirStore labels not updated: %s", gcpResourceUri)
			}

			topicName := path.Base(response.NotificationConfig.PubsubTopic)
			if topicName != pubsubTopic {
				return fmt.Errorf("fhirStore 'NotificationConfig' not updated ('%s' != '%s'): %s", topicName, pubsubTopic, gcpResourceUri)
			}
		}

		if !foundResource {
			return fmt.Errorf("google_healthcare_fhir_store resource was missing")
		}
		return nil
	}
}
