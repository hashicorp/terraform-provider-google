package google

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccCheckGoogleComposerEnvironmentConfig = `
data "google_composer_environment" "composer_env" {
	name = "data_google_composer_environment_test"
}
`

func TestAccDataSourceComposerEnvironment_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleComposerEnvironmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComposerEnvironmentMeta("data.google_composer_environment.composer_env"),
				),
			},
		},
	})
}

func testAccCheckGoogleComposerEnvironmentMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find environment data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("environment data source ID not set.")
		}

		configCountStr, ok := rs.Primary.Attributes["config.#"]
		if !ok {
			return errors.New("can't find 'config' attribute")
		}

		configCount, err := strconv.Atoi(configCountStr)
		if err != nil {
			return errors.New("failed to read number of valid config entries")
		}
		if configCount < 1 {
			return fmt.Errorf("expected at least 1 valid config entry, received %d, this is most likely a bug",
				configCount)
		}

		for i := 0; i < configCount; i++ {
			idx := "config." + strconv.Itoa(i)

			if v, ok := rs.Primary.Attributes[idx+".airflow_uri"]; !ok || v == "" {
				return fmt.Errorf("config %v is missing airflow_uri", i)
			}
			if v, ok := rs.Primary.Attributes[idx+".dag_gcs_prefix"]; !ok || v == "" {
				return fmt.Errorf("config %v is missing dag_gcs_prefix", i)
			}
			if v, ok := rs.Primary.Attributes[idx+".gke_cluster"]; !ok || v == "" {
				return fmt.Errorf("config %v is missing gke_cluster", i)
			}
		}

		return nil
	}
}
