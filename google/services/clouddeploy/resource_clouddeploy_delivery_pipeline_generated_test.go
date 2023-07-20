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

package clouddeploy_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	clouddeploy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/clouddeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccClouddeployDeliveryPipeline_DeliveryPipeline(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployDeliveryPipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployDeliveryPipeline_DeliveryPipeline(context),
			},
			{
				ResourceName:      "google_clouddeploy_delivery_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClouddeployDeliveryPipeline_DeliveryPipelineUpdate0(context),
			},
			{
				ResourceName:      "google_clouddeploy_delivery_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClouddeployDeliveryPipeline_DeliveryPipeline(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = "%{region}"
  name     = "tf-test-pipeline%{random_suffix}"

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  description = "basic description"

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }

  project = "%{project_name}"

  serial_pipeline {
    stages {
      deploy_parameters {
        values = {
          deployParameterKey = "deployParameterValue"
        }

        match_target_labels = {}
      }

      profiles  = ["example-profile-one", "example-profile-two"]
      target_id = "example-target-one"
    }

    stages {
      profiles  = []
      target_id = "example-target-two"
    }
  }
}


`, context)
}

func testAccClouddeployDeliveryPipeline_DeliveryPipelineUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = "%{region}"
  name     = "tf-test-pipeline%{random_suffix}"

  annotations = {
    my_second_annotation = "updated-example-annotation-2"

    my_third_annotation = "example-annotation-3"
  }

  description = "updated description"

  labels = {
    my_second_label = "updated-example-label-2"

    my_third_label = "example-label-3"
  }

  project = "%{project_name}"

  serial_pipeline {
    stages {
      profiles  = ["new-example-profile"]
      target_id = "example-target-two"
    }

    stages {
      profiles  = ["example-profile-four", "example-profile-five"]
      target_id = "example-target-three"
    }
  }

  suspended = true
}


`, context)
}

func testAccCheckClouddeployDeliveryPipelineDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_clouddeploy_delivery_pipeline" {
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

			obj := &clouddeploy.DeliveryPipeline{
				Location:    dcl.String(rs.Primary.Attributes["location"]),
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				Description: dcl.String(rs.Primary.Attributes["description"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				Suspended:   dcl.Bool(rs.Primary.Attributes["suspended"] == "true"),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Etag:        dcl.StringOrNil(rs.Primary.Attributes["etag"]),
				Uid:         dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLClouddeployClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetDeliveryPipeline(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_clouddeploy_delivery_pipeline still exists %v", obj)
			}
		}
		return nil
	}
}
