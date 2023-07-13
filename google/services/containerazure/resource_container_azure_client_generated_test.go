// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package containerazure_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	containerazure "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccContainerAzureClient_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"azure_app":     "00000000-0000-0000-0000-17aad2f0f61f",
		"azure_tenant":  "00000000-0000-0000-0000-17aad2f0f61f",
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAzureClientDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAzureClient_BasicHandWritten(context),
			},
			{
				ResourceName:      "google_container_azure_client.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerAzureClient_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_azure_client" "primary" {
  application_id = "%{azure_app}"
  location       = "us-west1"
  name           = "tf-test-client-name%{random_suffix}"
  tenant_id      = "%{azure_tenant}"
  project        = "%{project_name}"
}

`, context)
}

func testAccCheckContainerAzureClientDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_container_azure_client" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &containerazure.AzureClient{
				ApplicationId: dcl.String(rs.Primary.Attributes["application_id"]),
				Location:      dcl.String(rs.Primary.Attributes["location"]),
				Name:          dcl.String(rs.Primary.Attributes["name"]),
				TenantId:      dcl.String(rs.Primary.Attributes["tenant_id"]),
				Project:       dcl.StringOrNil(rs.Primary.Attributes["project"]),
				Certificate:   dcl.StringOrNil(rs.Primary.Attributes["certificate"]),
				CreateTime:    dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Uid:           dcl.StringOrNil(rs.Primary.Attributes["uid"]),
			}

			client := transport_tpg.NewDCLContainerAzureClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetClient(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_container_azure_client still exists %v", obj)
			}
		}
		return nil
	}
}
