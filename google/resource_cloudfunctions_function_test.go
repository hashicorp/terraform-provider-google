package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudfunctions/v1"
	"strings"
)

const (
	FUNCTION_TRIGGER_HTTP  = iota
	FUNCTION_TRIGGER_TOPIC
	FUNCTION_TRIGGER_BUCKET
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
					testAccCloudFunctionsFunctionDescription("test function", &function),
					testAccCloudFunctionsFunctionMemory(128, &function),
					testAccCloudFunctionsFunctionSource("gs://test-cloudfunctions-sk/index.zip", &function),
					testAccCloudFunctionsFunctionTrigger(FUNCTION_TRIGGER_HTTP, &function),
					testAccCloudFunctionsFunctionTimeout(360, &function),
					testAccCloudFunctionsFunctionEntryPoint("helloGET", &function),

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

func testAccCloudFunctionsFunctionSource(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if n != function.SourceArchiveUrl {
			return fmt.Errorf("Expected source to be %v, got %v", n, function.EntryPoint)
		}
		return nil
	}
}

func testAccCloudFunctionsFunctionEntryPoint(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if n != function.EntryPoint {
			return fmt.Errorf("Expected entry_point to be %v, got %v", n, function.EntryPoint)
		}
		return nil
	}
}

func testAccCloudFunctionsFunctionMemory(n int, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if int64(n) != function.AvailableMemoryMb {
			return fmt.Errorf("Expected memory to be %v, got %v", n, function.AvailableMemoryMb)
		}
		return nil
	}
}
func testAccCloudFunctionsFunctionTrigger(n int, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		switch n {
		case FUNCTION_TRIGGER_HTTP:
			if function.HttpsTrigger == nil {
				return fmt.Errorf("Expected trigger_http to be set")
			}
		case FUNCTION_TRIGGER_BUCKET:
			if function.EventTrigger == nil {
				return fmt.Errorf("Expected trigger_bucket to be set")
			}
			if strings.Index(function.EventTrigger.EventType, "cloud.storage") == -1 {
				return fmt.Errorf("Expected trigger_bucket to be set")
			}
		case FUNCTION_TRIGGER_TOPIC:
			if function.EventTrigger == nil {
				return fmt.Errorf("Expected trigger_bucket to be set")
			}
			if strings.Index(function.EventTrigger.EventType, "cloud.pubsub") == -1 {
				return fmt.Errorf("Expected trigger_topic to be set")
			}
		default:
			return fmt.Errorf("testAccCloudFunctionsFunctionTrigger expects only FUNCTION_TRIGGER_HTTP, " +
				"FUNCTION_TRIGGER_BUCKET or FUNCTION_TRIGGER_TOPIC")
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
