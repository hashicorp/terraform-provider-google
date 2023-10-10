// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package clouddeploy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccClouddeployTarget_withProviderDefaultLabels(t *testing.T) {
	// The test failed if VCR testing is enabled, because the cached provider config is used.
	// Any changes in the provider default labels will not be applied.
	acctest.SkipIfVcr(t)
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
				Config: testAccClouddeployTarget_withProviderDefaultLabels(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_second_label", "example-label-2"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "default_value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_clouddeploy_target.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccClouddeployTarget_resourceLabelsOverridesProviderDefaultLabels(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "3"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_second_label", "example-label-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_clouddeploy_target.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccClouddeployTarget_moveResourceLabelToProviderDefaultLabels(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_clouddeploy_target.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccClouddeployTarget_resourceLabelsOverridesProviderDefaultLabels(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "3"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_second_label", "example-label-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_clouddeploy_target.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccClouddeployTarget_withoutLabels(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_clouddeploy_target.primary", "labels.%"),
					resource.TestCheckNoResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%"),
					resource.TestCheckNoResourceAttr("google_clouddeploy_target.primary", "effective_labels.%"),
				),
			},
			{
				ResourceName:            "google_clouddeploy_target.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
		},
	})
}

func testAccClouddeployTarget_withProviderDefaultLabels(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
}

resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  deploy_parameters = {
    deployParameterKey = "deployParameterValue"
  }

  description = "basic description"

  gke {
    cluster = "projects/%{project_name}/locations/%{region}/clusters/example-cluster-name"
  }

  project          = "%{project_name}"
  require_approval = false

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  labels = {
    my_first_label = "example-label-1"
    my_second_label = "example-label-2"
  }
}
`, context)
}

func testAccClouddeployTarget_resourceLabelsOverridesProviderDefaultLabels(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
}

resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  deploy_parameters = {
    deployParameterKey = "deployParameterValue"
  }

  description = "basic description"

  gke {
    cluster = "projects/%{project_name}/locations/%{region}/clusters/example-cluster-name"
  }

  project          = "%{project_name}"
  require_approval = false

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  labels = {
    my_first_label = "example-label-1"
    my_second_label = "example-label-2"
    default_key1 = "value1"
  }
}
`, context)
}

func testAccClouddeployTarget_moveResourceLabelToProviderDefaultLabels(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
    my_second_label = "example-label-2"
  }
}

resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  deploy_parameters = {
    deployParameterKey = "deployParameterValue"
  }

  description = "basic description"

  gke {
    cluster = "projects/%{project_name}/locations/%{region}/clusters/example-cluster-name"
  }

  project          = "%{project_name}"
  require_approval = false

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  labels = {
    my_first_label = "example-label-1"
    default_key1 = "value1"
  }
}
`, context)
}

func testAccClouddeployTarget_withoutLabels(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_target" "primary" {
  location = "%{region}"
  name     = "tf-test-target%{random_suffix}"

  deploy_parameters = {
    deployParameterKey = "deployParameterValue"
  }

  description = "basic description"

  gke {
    cluster = "projects/%{project_name}/locations/%{region}/clusters/example-cluster-name"
  }

  project          = "%{project_name}"
  require_approval = false

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }
}
`, context)
}
