// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicemanagement_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccEndpointsService_basic(t *testing.T) {
	// Uses random provider
	acctest.SkipIfVcr(t)
	t.Parallel()
	serviceId := "tf-test" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		CheckDestroy:             testAccCheckEndpointServiceDestroyProducer(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointsService_basic(serviceId, envvar.GetTestProjectFromEnv(), "1"),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
			},
			{
				Config: testAccEndpointsService_basic(serviceId, envvar.GetTestProjectFromEnv(), "2"),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
			},
			{
				Config: testAccEndpointsService_basic(serviceId, envvar.GetTestProjectFromEnv(), "3"),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
			},
		},
	})
}

func TestAccEndpointsService_grpc(t *testing.T) {
	t.Parallel()
	serviceId := "tf-test" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEndpointServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointsService_grpc(serviceId, envvar.GetTestProjectFromEnv()),
				Check:  testAccCheckEndpointExistsByName(t, serviceId),
			},
		},
	})
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
		config := acctest.GoogleProviderConfig(t)
		service, err := config.NewServiceManClient(config.UserAgent).Services.GetConfig(
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
		config := acctest.GoogleProviderConfig(t)

		for name, rs := range s.RootModule().Resources {
			if strings.HasPrefix(name, "data.") {
				continue
			}
			if rs.Type != "google_endpoints_service" {
				continue
			}

			serviceName := rs.Primary.Attributes["service_name"]
			service, err := config.NewServiceManClient(config.UserAgent).Services.GetConfig(serviceName).Do()
			if err != nil {
				// ServiceManagement returns 403 if service doesn't exist.
				if !transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
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
