package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudfunctions/v1"
)

func TestAccCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

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
						"google_cloudfunctions_function.function", &function),
					testAccCloudFunctionsFunctionName(functionName, &function),
					testAccCloudFunctionsFunctionTimeout(360, &function),
					testAccCloudFunctionsFunctionDescription("test function", &function),
				),
			},
			{
				ResourceName:      "google_cloudfunctions_function.function",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudFunctionsFunctionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudfunctions_function" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		_, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(
			createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
		if err == nil {
			return fmt.Errorf("CloudFunctions still exists")
		}

	}

	return nil
}

func testAccCloudFunctionsFunctionExists(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
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
		found, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(
			createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
		if err != nil {
			return fmt.Errorf("CloudFunctions Function not present")
		}

		*function = *found

		return nil
	}
}

func testAccCloudFunctionsFunctionName(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expected, err := getCloudFunctionName(function.Name)
		if err != nil {
			return err
		}
		if n != expected {
			return fmt.Errorf("Expected function name %s, got %s", n, expected)
		}

		return nil
	}
}

func testAccCloudFunctionsFunctionTimeout(n int, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expected, err := readTimeout(function.Timeout)
		if err != nil {
			return err
		}
		if n != expected {
			return fmt.Errorf("Expected timeout to be %v, got %v", n, expected)
		}

		return nil
	}
}

func testAccCloudFunctionsFunctionDescription(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if n != function.Description {
			return fmt.Errorf("Expected description to be %v, got %v", n, function.Description)
		}

		return nil
	}
}

func testAccCloudFunctionsFunction(functionName string) string {
	return fmt.Sprintf(`
resource "google_cloudfunctions_function" "function" {
  name          = "%s"
  description   = "test function"
  memory		= 128
  source        = "gs://test-cloudfunctions-sk/index.zip"
  trigger_http  = true
  timeout		= 360
  entry_point   = "helloGET"
}
`, functionName)
}
