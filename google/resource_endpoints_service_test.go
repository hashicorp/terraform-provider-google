package google

import (
	"reflect"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/servicemanagement/v1"
)

func TestAccEndpointsService_basic(t *testing.T) {
	t.Parallel()
	random_name := "t-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointsService_basic(random_name),
				Check:  testAccCheckEndpointExistsByName(random_name),
			},
		},
	})
}

func TestAccEndpointsService_grpc(t *testing.T) {
	t.Parallel()
	random_name := "t-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointsService_grpc(random_name),
				Check:  testAccCheckEndpointExistsByName(random_name),
			},
		},
	})
}

func TestEndpointsService_grpcMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		Attributes         map[string]string
		ExpectedAttributes map[string]string
		Meta               interface{}
	}{
		"update from protoc_output to protoc_output_base64": {
			StateVersion: 0,
			Attributes: map[string]string{
				"protoc_output": "123456789",
				"name":          "testcase",
			},
			ExpectedAttributes: map[string]string{
				"protoc_output_base64": "MTIzNDU2Nzg5",
				"protoc_output":        "",
				"name":                 "testcase",
			},
			Meta: &Config{Project: "gcp-project", Region: "us-central1"},
		},
		"update from non-protoc_output": {
			StateVersion: 0,
			Attributes: map[string]string{
				"openapi_config": "foo bar baz",
				"name":           "testcase-2",
			},
			ExpectedAttributes: map[string]string{
				"openapi_config": "foo bar baz",
				"name":           "testcase-2",
			},
			Meta: &Config{Project: "gcp-project", Region: "us-central1"},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.Attributes["name"],
			Attributes: tc.Attributes,
		}

		is, err := migrateEndpointsService(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if !reflect.DeepEqual(is.Attributes, tc.ExpectedAttributes) {
			t.Fatalf("Attributes should be `%s` but are `%s`", tc.ExpectedAttributes, is.Attributes)
		}
	}
}

func testAccEndpointsService_basic(random_name string) string {
	return fmt.Sprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name   = "%s.endpoints.%s.cloud.goog"
  project        = "%s"
  openapi_config = <<EOF
swagger: "2.0"
info:
  description: "A simple Google Cloud Endpoints API example."
  title: "Endpoints Example"
  version: "1.0.0"
host: "%s.endpoints.%s.cloud.goog"
basePath: "/"
consumes:
- "application/json"
produces:
- "application/json"
schemes:
- "https"
paths:
  "/echo":
    post:
      description: "Echo back a given message."
      operationId: "echo"
      produces:
      - "application/json"
      responses:
        200:
          description: "Echo"
          schema:
            $ref: "#/definitions/echoMessage"
      parameters:
      - description: "Message to echo"
        in: body
        name: message
        required: true
        schema:
          $ref: "#/definitions/echoMessage"
      security:
      - api_key: []
definitions:
  echoMessage:
    properties:
      message:
        type: "string"
EOF

}
`, random_name, getTestProjectFromEnv(), getTestProjectFromEnv(), random_name, getTestProjectFromEnv())
}

func testAccEndpointsService_grpc(random_name string) string {
	return fmt.Sprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name = "%s.endpoints.%s.cloud.goog"
  project      = "%s"
  grpc_config  = <<EOF
type: google.api.Service
config_version: 3
name: %s.endpoints.%s.cloud.goog
usage:
  rules:
  - selector: endpoints.examples.bookstore.Bookstore.ListShelves
    allow_unregistered_calls: true
EOF

  protoc_output_base64 = filebase64("test-fixtures/test_api_descriptor.pb")
}
`, random_name, getTestProjectFromEnv(), getTestProjectFromEnv(), random_name, getTestProjectFromEnv())
}

func testAccCheckEndpointExistsByName(random_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		servicesService := servicemanagement.NewServicesService(config.clientServiceMan)
		service, err := servicesService.GetConfig(fmt.Sprintf("%s.endpoints.%s.cloud.goog", random_name, config.Project)).Do()
		if err != nil {
			return err
		}
		if service != nil {
			return nil
		} else {
			return fmt.Errorf("Service %s.endpoints.%s.cloud.goog does not seem to exist.", random_name, config.Project)
		}
	}
}
