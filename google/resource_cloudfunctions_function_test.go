package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction(functionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function"),
				),
			},
		},
	})
}

func testAccCheckCloudFunctionsFunctionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudfunctions_function" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		_, err := config.clientCloudFunctions.Projects.Locations.Functions.Delete(
			createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
		if err == nil {
			return fmt.Errorf("CloudFunction still exists")
		}

	}

	return nil
}

func testAccCloudFunctionsFunctionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		getOpt, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(
			createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
		if err != nil {
			return fmt.Errorf("CloudFunctions Function not present")
		}
		nameFromGet, err := getCloudFunctionName(getOpt.Name)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(nameFromGet, rs.Primary.Attributes["name"]) {
			return fmt.Errorf("CloudFunctions Function name does not match expected value")
		}

		return nil
	}
}

func testAccCloudFunctionsFunction(functionName string) string {
	return fmt.Sprintf(`
resource "google_cloudfunctions_function" "function" {
  name          = "%s"
  source        = "gs://test-cloudfunctions-sk/index.zip"
  trigger_http  = ""
}
`, functionName)
}
