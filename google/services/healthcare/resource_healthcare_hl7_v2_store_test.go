// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package healthcare_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/healthcare"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccHealthcareHl7V2StoreIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId             string
		ExpectedError        bool
		ExpectedTerraformId  string
		ExpectedHl7V2StoreId string
		Config               *transport_tpg.Config
	}{
		"id is in project/location/datasetName/hl7V2StoreName format": {
			ImportId:             "test-project/us-central1/test-dataset/test-store-name",
			ExpectedError:        false,
			ExpectedTerraformId:  "test-project/us-central1/test-dataset/test-store-name",
			ExpectedHl7V2StoreId: "projects/test-project/locations/us-central1/datasets/test-dataset/hl7V2Stores/test-store-name",
		},
		"id is in domain:project/location/datasetName/hl7V2StoreName format": {
			ImportId:             "example.com:test-project/us-central1/test-dataset/test-store-name",
			ExpectedError:        false,
			ExpectedTerraformId:  "example.com:test-project/us-central1/test-dataset/test-store-name",
			ExpectedHl7V2StoreId: "projects/example.com:test-project/locations/us-central1/datasets/test-dataset/hl7V2Stores/test-store-name",
		},
		"id is in location/datasetName/hl7V2StoreName format": {
			ImportId:             "us-central1/test-dataset/test-store-name",
			ExpectedError:        false,
			ExpectedTerraformId:  "test-project/us-central1/test-dataset/test-store-name",
			ExpectedHl7V2StoreId: "projects/test-project/locations/us-central1/datasets/test-dataset/hl7V2Stores/test-store-name",
			Config:               &transport_tpg.Config{Project: "test-project"},
		},
		"id is in location/datasetName/hl7V2StoreName format without project in config": {
			ImportId:      "us-central1/test-dataset/test-store-name",
			ExpectedError: true,
			Config:        &transport_tpg.Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		hl7V2StoreId, err := healthcare.ParseHealthcareHl7V2StoreId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if hl7V2StoreId.TerraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, hl7V2StoreId.TerraformId())
		}

		if hl7V2StoreId.Hl7V2StoreId() != tc.ExpectedHl7V2StoreId {
			t.Fatalf("bad: %s, expected Hl7V2Store ID to be `%s` but is `%s`", tn, tc.ExpectedHl7V2StoreId, hl7V2StoreId.Hl7V2StoreId())
		}
	}
}

func TestAccHealthcareHl7V2Store_basic(t *testing.T) {
	t.Parallel()

	datasetName := fmt.Sprintf("tf-test-dataset-%s", acctest.RandString(t, 10))
	hl7_v2StoreName := fmt.Sprintf("tf-test-hl7_v2-store-%s", acctest.RandString(t, 10))
	pubsubTopic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	resourceName := "google_healthcare_hl7_v2_store.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcareHl7V2StoreDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleHealthcareHl7V2Store_basic(hl7_v2StoreName, datasetName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleHealthcareHl7V2Store_update(hl7_v2StoreName, datasetName, pubsubTopic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleHealthcareHl7V2StoreUpdate(t, pubsubTopic),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleHealthcareHl7V2Store_basic(hl7_v2StoreName, datasetName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleHealthcareHl7V2Store_basic(hl7_v2StoreName, datasetName string) string {
	return fmt.Sprintf(`
resource "google_healthcare_hl7_v2_store" "default" {
  name     = "%s"
  dataset  = google_healthcare_dataset.dataset.id
}

resource "google_healthcare_dataset" "dataset" {
  name     = "%s"
  location = "us-central1"
}
`, hl7_v2StoreName, datasetName)
}

func testGoogleHealthcareHl7V2Store_update(hl7_v2StoreName, datasetName, pubsubTopic string) string {
	return fmt.Sprintf(`
resource "google_healthcare_hl7_v2_store" "default" {
  name     = "%s"
  dataset  = google_healthcare_dataset.dataset.id

  parser_config {
    allow_null_header  = true
    segment_terminator = "Jw=="
  }

  notification_configs {
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
`, hl7_v2StoreName, datasetName, pubsubTopic)
}

func testAccCheckGoogleHealthcareHl7V2StoreUpdate(t *testing.T, pubsubTopic string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var foundResource = false
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_healthcare_hl7_v2_store" {
				continue
			}
			foundResource = true

			config := acctest.GoogleProviderConfig(t)

			gcpResourceUri, err := tpgresource.ReplaceVarsForTest(config, rs, "{{dataset}}/hl7V2Stores/{{name}}")
			if err != nil {
				return err
			}

			response, err := config.NewHealthcareClient(config.UserAgent).Projects.Locations.Datasets.Hl7V2Stores.Get(gcpResourceUri).Do()
			if err != nil {
				return fmt.Errorf("Unexpected failure while verifying 'updated' dataset: %s", err)
			}

			if response.ParserConfig == nil {
				return fmt.Errorf("hl7_v2_store had no parser config: %s", gcpResourceUri)
			}

			if !response.ParserConfig.AllowNullHeader {
				return fmt.Errorf("hl7_v2_store allowNullHeader not changed to true: %s", gcpResourceUri)
			}

			if response.ParserConfig.SegmentTerminator != "Jw==" {
				return fmt.Errorf("hl7_v2_store segmentTerminator was not changed to 'JW==' as was expected: %s", gcpResourceUri)
			}

			if len(response.Labels) == 0 || response.Labels["label1"] != "labelvalue1" {
				return fmt.Errorf("hl7_v2_store labels not updated: %s", gcpResourceUri)
			}

			notifications := response.NotificationConfigs
			if len(notifications) > 0 {
				topicName := path.Base(notifications[0].PubsubTopic)
				if topicName != pubsubTopic {
					return fmt.Errorf("hl7_v2_store 'NotificationConfig' not updated ('%s' != '%s'): %s", topicName, pubsubTopic, gcpResourceUri)
				}
			}
		}

		if !foundResource {
			return fmt.Errorf("google_healthcare_hl7_v2_store resource was missing")
		}
		return nil
	}
}
