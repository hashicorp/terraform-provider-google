package google

import (
	"fmt"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHealthcareFhirStoreIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedFhirStoreId string
		Config              *Config
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
			Config:              &Config{Project: "test-project"},
		},
		"id is in location/datasetName/fhirStoreName format without project in config": {
			ImportId:      "us-central1/test-dataset/test-store-name",
			ExpectedError: true,
			Config:        &Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		fhirStoreId, err := parseHealthcareFhirStoreId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if fhirStoreId.terraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, fhirStoreId.terraformId())
		}

		if fhirStoreId.fhirStoreId() != tc.ExpectedFhirStoreId {
			t.Fatalf("bad: %s, expected FhirStore ID to be `%s` but is `%s`", tn, tc.ExpectedFhirStoreId, fhirStoreId.fhirStoreId())
		}
	}
}

func TestAccHealthcareFhirStore_basic(t *testing.T) {
	t.Parallel()

	datasetName := fmt.Sprintf("tf-test-dataset-%s", randString(t, 10))
	fhirStoreName := fmt.Sprintf("tf-test-fhir-store-%s", randString(t, 10))
	pubsubTopic := fmt.Sprintf("tf-test-topic-%s", randString(t, 10))
	resourceName := "google_healthcare_fhir_store.default"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHealthcareFhirStoreDestroyProducer(t),
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

			config := googleProviderConfig(t)

			gcpResourceUri, err := replaceVarsForTest(config, rs, "{{dataset}}/fhirStores/{{name}}")
			if err != nil {
				return err
			}

			response, err := config.clientHealthcare.Projects.Locations.Datasets.FhirStores.Get(gcpResourceUri).Do()
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
