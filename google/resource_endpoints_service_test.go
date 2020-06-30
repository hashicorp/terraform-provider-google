package google

import (
	"reflect"
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccEndpointsService_basic(t *testing.T) {
	t.Parallel()
	serviceId := "tf-test" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckEndpointServiceDestroyProducer(t),
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointsService_basic(serviceId, getTestProjectFromEnv(), "1"),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
			},
			{
				Config: testAccEndpointsService_basic(serviceId, getTestProjectFromEnv(), "2"),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
			},
			{
				Config: testAccEndpointsService_basic(serviceId, getTestProjectFromEnv(), "3"),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
			},
		},
	})
}

func TestAccEndpointsService_grpc(t *testing.T) {
	t.Parallel()
	serviceId := "tf-test" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointsService_grpc(serviceId, getTestProjectFromEnv()),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
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

func testAccEndpointsService_basic(serviceId, project, rev string) string {
	return fmt.Sprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name   = "%[1]s.endpoints.%[2]s.cloud.goog"
  project        = "%[2]s"
  openapi_config = <<EOF
swagger: "2.0"
info:
  description: "A simple Google Cloud Endpoints API example."
  title: "Endpoints Example, rev. %[3]s"
  version: "1.0.0"
host: "%[1]s.endpoints.%[2]s.cloud.goog"
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

resource "random_id" "foo" {
  keepers = {
    config_id = google_endpoints_service.endpoints_service.config_id
  }
  byte_length = 8
}
`, serviceId, project, rev)
}

func testAccEndpointsService_grpc(serviceId, project string) string {
	return fmt.Sprintf(`
resource "google_endpoints_service" "endpoints_service" {
  service_name = "%[1]s.endpoints.%[2]s.cloud.goog"
  project      = "%[2]s"
  grpc_config  = <<EOF
type: google.api.Service
config_version: 3
name: %[1]s.endpoints.%[2]s.cloud.goog
usage:
  rules:
  - selector: endpoints.examples.bookstore.Bookstore.ListShelves
    allow_unregistered_calls: true
EOF

  protoc_output_base64 = filebase64("test-fixtures/test_api_descriptor.pb")
}
`, serviceId, project)
}

func testAccCheckEndpointExistsByName(t *testing.T, serviceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		service, err := config.clientServiceMan.Services.GetConfig(
			fmt.Sprintf("%s.endpoints.%s.cloud.goog", serviceId, config.Project)).Do()
		if err != nil {
			return err
		}
		if service != nil {
			return nil
		} else {
			return fmt.Errorf("Service %s.endpoints.%s.cloud.goog does not seem to exist.", serviceId, config.Project)
		}
	}
}

func testAccCheckEndpointServiceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for name, rs := range s.RootModule().Resources {
			if strings.HasPrefix(name, "data.") {
				continue
			}
			if rs.Type != "google_endpoints_service" {
				continue
			}

			serviceName := rs.Primary.Attributes["service_name"]
			service, err := config.clientServiceMan.Services.GetConfig(serviceName).Do()
			if err != nil {
				// ServiceManagement returns 403 if service doesn't exist.
				if !isGoogleApiErrorWithCode(err, 403) {
					return err
				}
			}
			if service != nil {
				return fmt.Errorf("expected service %q to have been destroyed, got %+v", service.Name, service)
			}
		}
		return nil
	}
}
