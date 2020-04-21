package google

import (
	"fmt"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHealthcareDicomStoreIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId             string
		ExpectedError        bool
		ExpectedTerraformId  string
		ExpectedDicomStoreId string
		Config               *Config
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
			Config:               &Config{Project: "test-project"},
		},
		"id is in location/datasetName/dicomStoreName format without project in config": {
			ImportId:      "us-central1/test-dataset/test-store-name",
			ExpectedError: true,
			Config:        &Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		dicomStoreId, err := parseHealthcareDicomStoreId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if dicomStoreId.terraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, dicomStoreId.terraformId())
		}

		if dicomStoreId.dicomStoreId() != tc.ExpectedDicomStoreId {
			t.Fatalf("bad: %s, expected DicomStore ID to be `%s` but is `%s`", tn, tc.ExpectedDicomStoreId, dicomStoreId.dicomStoreId())
		}
	}
}

func TestAccHealthcareDicomStore_basic(t *testing.T) {
	t.Parallel()

	datasetName := fmt.Sprintf("tf-test-dataset-%s", randString(t, 10))
	dicomStoreName := fmt.Sprintf("tf-test-dicom-store-%s", randString(t, 10))
	pubsubTopic := fmt.Sprintf("tf-test-topic-%s", randString(t, 10))
	resourceName := "google_healthcare_dicom_store.default"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHealthcareDicomStoreDestroyProducer(t),
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
			// TODO(b/148536607): Uncomment once b/148536607 is fixed.
			// {
			// 	Config: testGoogleHealthcareDicomStore_basic(dicomStoreName, datasetName),
			// },
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
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

			config := googleProviderConfig(t)

			gcpResourceUri, err := replaceVarsForTest(config, rs, "{{dataset}}/dicomStores/{{name}}")
			if err != nil {
				return err
			}

			response, err := config.clientHealthcare.Projects.Locations.Datasets.DicomStores.Get(gcpResourceUri).Do()
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
