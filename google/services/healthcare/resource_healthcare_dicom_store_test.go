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

func TestAccHealthcareDicomStoreIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId             string
		ExpectedError        bool
		ExpectedTerraformId  string
		ExpectedDicomStoreId string
		Config               *transport_tpg.Config
	}{
		"id is in project/location/datasetName/dicomStoreName format": {
			ImportId:             "test-project/us-central1/test-dataset/test-store-name",
			ExpectedError:        false,
			ExpectedTerraformId:  "test-project/us-central1/test-dataset/test-store-name",
			ExpectedDicomStoreId: "projects/test-project/locations/us-central1/datasets/test-dataset/dicomStores/test-store-name",
		},
		"id is in domain:project/location/datasetName/dicomStoreName format": {
			ImportId:             "example.com:test-project/us-central1/test-dataset/test-store-name",
			ExpectedError:        false,
			ExpectedTerraformId:  "example.com:test-project/us-central1/test-dataset/test-store-name",
			ExpectedDicomStoreId: "projects/example.com:test-project/locations/us-central1/datasets/test-dataset/dicomStores/test-store-name",
		},
		"id is in location/datasetName/dicomStoreName format": {
			ImportId:             "us-central1/test-dataset/test-store-name",
			ExpectedError:        false,
			ExpectedTerraformId:  "test-project/us-central1/test-dataset/test-store-name",
			ExpectedDicomStoreId: "projects/test-project/locations/us-central1/datasets/test-dataset/dicomStores/test-store-name",
			Config:               &transport_tpg.Config{Project: "test-project"},
		},
		"id is in location/datasetName/dicomStoreName format without project in config": {
			ImportId:      "us-central1/test-dataset/test-store-name",
			ExpectedError: true,
			Config:        &transport_tpg.Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		dicomStoreId, err := healthcare.ParseHealthcareDicomStoreId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if dicomStoreId.TerraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, dicomStoreId.TerraformId())
		}

		if dicomStoreId.DicomStoreId() != tc.ExpectedDicomStoreId {
			t.Fatalf("bad: %s, expected DicomStore ID to be `%s` but is `%s`", tn, tc.ExpectedDicomStoreId, dicomStoreId.DicomStoreId())
		}
	}
}

func TestAccHealthcareDicomStore_basic(t *testing.T) {
	t.Parallel()

	datasetName := fmt.Sprintf("tf-test-dataset-%s", acctest.RandString(t, 10))
	dicomStoreName := fmt.Sprintf("tf-test-dicom-store-%s", acctest.RandString(t, 10))
	pubsubTopic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	resourceName := "google_healthcare_dicom_store.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcareDicomStoreDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleHealthcareDicomStore_basic(dicomStoreName, datasetName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleHealthcareDicomStore_update(dicomStoreName, datasetName, pubsubTopic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleHealthcareDicomStoreUpdate(t, pubsubTopic),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleHealthcareDicomStore_basic(dicomStoreName, datasetName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleHealthcareDicomStore_basic(dicomStoreName, datasetName string) string {
	return fmt.Sprintf(`
resource "google_healthcare_dicom_store" "default" {
  name     = "%s"
  dataset  = google_healthcare_dataset.dataset.id
}

resource "google_healthcare_dataset" "dataset" {
  name     = "%s"
  location = "us-central1"
}
`, dicomStoreName, datasetName)
}

func testGoogleHealthcareDicomStore_update(dicomStoreName, datasetName, pubsubTopic string) string {
	return fmt.Sprintf(`
resource "google_healthcare_dicom_store" "default" {
  name     = "%s"
  dataset  = google_healthcare_dataset.dataset.id

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
`, dicomStoreName, datasetName, pubsubTopic)
}

func testAccCheckGoogleHealthcareDicomStoreUpdate(t *testing.T, pubsubTopic string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var foundResource = false
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_healthcare_dicom_store" {
				continue
			}
			foundResource = true

			config := acctest.GoogleProviderConfig(t)

			gcpResourceUri, err := tpgresource.ReplaceVarsForTest(config, rs, "{{dataset}}/dicomStores/{{name}}")
			if err != nil {
				return err
			}

			response, err := config.NewHealthcareClient(config.UserAgent).Projects.Locations.Datasets.DicomStores.Get(gcpResourceUri).Do()
			if err != nil {
				return fmt.Errorf("Unexpected failure while verifying 'updated' dataset: %s", err)
			}

			if len(response.Labels) == 0 || response.Labels["label1"] != "labelvalue1" {
				return fmt.Errorf("dicomStore labels not updated: %s", gcpResourceUri)
			}

			topicName := path.Base(response.NotificationConfig.PubsubTopic)
			if topicName != pubsubTopic {
				return fmt.Errorf("dicomStore 'NotificationConfig' not updated ('%s' != '%s'): %s", topicName, pubsubTopic, gcpResourceUri)
			}
		}

		if !foundResource {
			return fmt.Errorf("google_healthcare_dicom_store resource was missing")
		}
		return nil
	}
}
