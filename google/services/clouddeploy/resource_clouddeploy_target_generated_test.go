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

func TestAccClouddeployTarget_Target(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployTargetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployTarget_Target(context),
			},
			{
				ResourceName:      "google_clouddeploy_target.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClouddeployTarget_TargetUpdate0(context),
			},
			{
				ResourceName:      "google_clouddeploy_target.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClouddeployTarget_TargetUpdate1(context),
			},
			{
				ResourceName:      "google_clouddeploy_target.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClouddeployTarget_TargetUpdate2(context),
			},
			{
				ResourceName:      "google_clouddeploy_target.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClouddeployTarget_TargetUpdate3(context),
			},
			{
				ResourceName:      "google_clouddeploy_target.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClouddeployTarget_Target(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  deploy_parameters = {
    deployParameterKey = "deployParameterValue"
  }

  description = "basic description"

  gke {
    cluster = "projects/%{project_name}/locations/%{region}/clusters/example-cluster-name"
  }

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }

  project          = "%{project_name}"
  require_approval = false
}


`, context)
}

func testAccClouddeployTarget_TargetUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  annotations = {
    my_second_annotation = "updated-example-annotation-2"

    my_third_annotation = "example-annotation-3"
  }

  deploy_parameters = {}
  description       = "updated description"

  gke {
    cluster     = "projects/%{project_name}/locations/%{region}/clusters/different-example-cluster-name"
    internal_ip = true
  }

  labels = {
    my_second_label = "updated-example-label-2"

    my_third_label = "example-label-3"
  }

  project          = "%{project_name}"
  require_approval = true
}


`, context)
}

func testAccClouddeployTarget_TargetUpdate1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  annotations = {
    my_second_annotation = "updated-example-annotation-2"

    my_third_annotation = "example-annotation-3"
  }

  deploy_parameters = {}
  description       = "updated description"

  execution_configs {
    usages           = ["RENDER", "DEPLOY"]
    artifact_storage = "gs://my-bucket/my-dir"
    service_account  = "pool-owner@%{project_name}.iam.gserviceaccount.com"
  }

  gke {
    cluster     = "projects/%{project_name}/locations/%{region}/clusters/different-example-cluster-name"
    internal_ip = true
  }

  labels = {
    my_second_label = "updated-example-label-2"

    my_third_label = "example-label-3"
  }

  project          = "%{project_name}"
  require_approval = true
}


`, context)
}

func testAccClouddeployTarget_TargetUpdate2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  annotations = {
    my_second_annotation = "updated-example-annotation-2"

    my_third_annotation = "example-annotation-3"
  }

  deploy_parameters = {}
  description       = "updated description"

  execution_configs {
    usages           = ["RENDER"]
    artifact_storage = "gs://my-bucket/my-dir"
    service_account  = "pool-owner@%{project_name}.iam.gserviceaccount.com"
  }

  execution_configs {
    usages           = ["DEPLOY"]
    artifact_storage = "gs://deploy-bucket/deploy-dir"
    service_account  = "deploy-pool-owner@%{project_name}.iam.gserviceaccount.com"
    worker_pool      = "projects/%{project_name}/locations/%{region}/workerPools/my-deploy-pool"
  }

  gke {
    cluster     = "projects/%{project_name}/locations/%{region}/clusters/different-example-cluster-name"
    internal_ip = true
  }

  labels = {
    my_second_label = "updated-example-label-2"

    my_third_label = "example-label-3"
  }

  project          = "%{project_name}"
  require_approval = true
}


`, context)
}

func testAccClouddeployTarget_TargetUpdate3(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  annotations = {
    my_second_annotation = "updated-example-annotation-2"

    my_third_annotation = "example-annotation-3"
  }

  deploy_parameters = {}
  description       = "updated description"

  execution_configs {
    usages           = ["RENDER"]
    artifact_storage = "gs://other-bucket/other-dir"
    service_account  = "other-owner@%{project_name}.iam.gserviceaccount.com"
  }

  execution_configs {
    usages           = ["DEPLOY"]
    artifact_storage = "gs://deploy-bucket/deploy-dir"
    service_account  = "deploy-pool-owner@%{project_name}.iam.gserviceaccount.com"
    worker_pool      = "projects/%{project_name}/locations/%{region}/workerPools/my-deploy-pool"
  }

  gke {
    cluster     = "projects/%{project_name}/locations/%{region}/clusters/different-example-cluster-name"
    internal_ip = true
  }

  labels = {
    my_second_label = "updated-example-label-2"

    my_third_label = "example-label-3"
  }

  project          = "%{project_name}"
  require_approval = true
}


`, context)
}

func testAccCheckClouddeployTargetDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_clouddeploy_target" {
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

			obj := &clouddeploy.Target{
				Location:        dcl.String(rs.Primary.Attributes["location"]),
				Name:            dcl.String(rs.Primary.Attributes["name"]),
				Description:     dcl.String(rs.Primary.Attributes["description"]),
				Project:         dcl.StringOrNil(rs.Primary.Attributes["project"]),
				RequireApproval: dcl.Bool(rs.Primary.Attributes["require_approval"] == "true"),
				CreateTime:      dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Etag:            dcl.StringOrNil(rs.Primary.Attributes["etag"]),
				TargetId:        dcl.StringOrNil(rs.Primary.Attributes["target_id"]),
				Uid:             dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:      dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLClouddeployClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetTarget(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_clouddeploy_target still exists %v", obj)
			}
		}
		return nil
	}
}
