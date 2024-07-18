// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package integrationconnectors_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIntegrationConnectorsEndpointAttachment_integrationConnectorsEndpointAttachmentExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsEndpointAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsEndpointAttachment_integrationConnectorsEndpointAttachmentExample_full(context),
			},
			{
				ResourceName:            "google_integration_connectors_endpoint_attachment.sampleendpointattachment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccIntegrationConnectorsEndpointAttachment_integrationConnectorsEndpointAttachmentExample_update(context),
			},
			{
				ResourceName:            "google_integration_connectors_endpoint_attachment.sampleendpointattachment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccIntegrationConnectorsEndpointAttachment_integrationConnectorsEndpointAttachmentExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_integration_connectors_endpoint_attachment" "sampleendpointattachment" {
  name     = "tf-test-test-endpoint-attachment%{random_suffix}"
  location = "us-central1"
  description = "tf created description"
  # Third party dependency, see https://github.com/GoogleCloudPlatform/magic-modules/pull/9616#discussion_r1429029155
  service_attachment = "projects/connectors-example/regions/us-central1/serviceAttachments/test"
  labels = {
    foo = "bar"
  }
  endpoint_global_access = false
}
`, context)
}

func testAccIntegrationConnectorsEndpointAttachment_integrationConnectorsEndpointAttachmentExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_integration_connectors_endpoint_attachment" "sampleendpointattachment" {
  name     = "tf-test-test-endpoint-attachment%{random_suffix}"
  location = "us-central1"
  description = "tf updated description"
  # Third party dependency, see https://github.com/GoogleCloudPlatform/magic-modules/pull/9616#discussion_r1429029155
  service_attachment = "projects/connectors-example/regions/us-central1/serviceAttachments/test"
  labels = {
    bar = "foo"
  }
  endpoint_global_access = true
}
`, context)
}
