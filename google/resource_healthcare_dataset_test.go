package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHealthcareDatasetIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedDatasetId   string
		Config              *Config
	}{
		"id is in project/location/datasetName format": {
			ImportId:            "test-project/us-central1/test-dataset",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-dataset",
			ExpectedDatasetId:   "projects/test-project/locations/us-central1/datasets/test-dataset",
		},
		"id is in domain:project/location/datasetName format": {
			ImportId:            "example.com:test-project/us-central1/test-dataset",
			ExpectedError:       false,
			ExpectedTerraformId: "example.com:test-project/us-central1/test-dataset",
			ExpectedDatasetId:   "projects/example.com:test-project/locations/us-central1/datasets/test-dataset",
		},
		"id is in location/datasetName format": {
			ImportId:            "us-central1/test-dataset",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-dataset",
			ExpectedDatasetId:   "projects/test-project/locations/us-central1/datasets/test-dataset",
			Config:              &Config{Project: "test-project"},
		},
		"id is in location/datasetName format without project in config": {
			ImportId:      "us-central1/test-dataset",
			ExpectedError: true,
			Config:        &Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		datasetId, err := parseHealthcareDatasetId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if datasetId.terraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, datasetId.terraformId())
		}

		if datasetId.datasetId() != tc.ExpectedDatasetId {
			t.Fatalf("bad: %s, expected Dataset ID to be `%s` but is `%s`", tn, tc.ExpectedDatasetId, datasetId.datasetId())
		}
	}
}

func TestAccHealthcareDataset_basic(t *testing.T) {
	t.Parallel()

	location := "us-central1"
	datasetName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	timeZone := "America/New_York"
	resourceName := "google_healthcare_dataset.dataset"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHealthcareDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleHealthcareDataset_basic(datasetName, location),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleHealthcareDataset_update(datasetName, location, timeZone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleHealthcareDatasetUpdate(t, timeZone),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleHealthcareDatasetUpdate(t *testing.T, timeZone string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_healthcare_dataset" {
				continue
			}

			config := googleProviderConfig(t)

			gcpResourceUri, err := replaceVarsForTest(config, rs, "projects/{{project}}/locations/{{location}}/datasets/{{name}}")
			if err != nil {
				return err
			}

			response, err := config.clientHealthcare.Projects.Locations.Datasets.Get(gcpResourceUri).Do()
			if err != nil {
				return fmt.Errorf("Unexpected failure while verifying 'updated' dataset: %s", err)
			}

			if response.TimeZone != timeZone {
				return fmt.Errorf("Dataset timeZone was not set to '%s' as expected: %s", timeZone, gcpResourceUri)
			}
		}

		return nil
	}
}

func testGoogleHealthcareDataset_basic(datasetName, location string) string {
	return fmt.Sprintf(`
resource "google_healthcare_dataset" "dataset" {
  name     = "%s"
  location = "%s"
}
`, datasetName, location)
}

func testGoogleHealthcareDataset_update(datasetName, location, timeZone string) string {
	return fmt.Sprintf(`
resource "google_healthcare_dataset" "dataset" {
  name      = "%s"
  location  = "%s"
  time_zone = "%s"
}
`, datasetName, location, timeZone)
}
