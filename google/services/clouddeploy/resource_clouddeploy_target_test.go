// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package clouddeploy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

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
  add_terraform_attribution_label = false
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
  add_terraform_attribution_label = false
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
  add_terraform_attribution_label = false
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
provider "google" {
  add_terraform_attribution_label = false
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
}
`, context)
}

func TestAccClouddeployTarget_withAttributionDisabled(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":         envvar.GetTestProjectFromEnv(),
		"region":               envvar.GetTestRegionFromEnv(),
		"random_suffix":        acctest.RandString(t, 10),
		"add_attribution":      "false",
		"attribution_strategy": "CREATION_ONLY",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployTargetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployTarget_createWithAttribution(context),
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
				Config: testAccClouddeployTarget_updateWithAttribution(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-updated-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_second_label", "example-label-updated-2"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-updated-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-updated-2"),
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
				Config: testAccClouddeployTarget_clearWithAttribution(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_clouddeploy_target.primary", "labels.%"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "default_value1"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "1"),
				),
			},
		},
	})
}

func TestAccClouddeployTarget_withCreationOnlyAttribution(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":         envvar.GetTestProjectFromEnv(),
		"region":               envvar.GetTestRegionFromEnv(),
		"random_suffix":        acctest.RandString(t, 10),
		"add_attribution":      "true",
		"attribution_strategy": "CREATION_ONLY",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployTargetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployTarget_createWithAttribution(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_second_label", "example-label-2"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "4"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.goog-terraform-provisioned", "true"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "4"),
				),
			},
			{
				ResourceName:            "google_clouddeploy_target.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccClouddeployTarget_updateWithAttribution(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-updated-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_second_label", "example-label-updated-2"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "4"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-updated-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-updated-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.goog-terraform-provisioned", "true"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "4"),
				),
			},
			{
				ResourceName:            "google_clouddeploy_target.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "annotations"},
			},
			{
				Config: testAccClouddeployTarget_clearWithAttribution(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_clouddeploy_target.primary", "labels.%"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.goog-terraform-provisioned", "true"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "2"),
				),
			},
		},
	})
}

func TestAccClouddeployTarget_withProactiveAttribution(t *testing.T) {
	// VCR tests cache provider configuration between steps, this test changes provider configuration and fails under VCR.
	acctest.SkipIfVcr(t)
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	disabledContext := map[string]interface{}{
		"project_name":         envvar.GetTestProjectFromEnv(),
		"region":               envvar.GetTestRegionFromEnv(),
		"random_suffix":        suffix,
		"add_attribution":      "false",
		"attribution_strategy": "PROACTIVE",
	}
	enabledContext := map[string]interface{}{
		"project_name":         envvar.GetTestProjectFromEnv(),
		"region":               envvar.GetTestRegionFromEnv(),
		"random_suffix":        suffix,
		"add_attribution":      "true",
		"attribution_strategy": "PROACTIVE",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployTargetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployTarget_createWithAttribution(disabledContext),
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
				Config: testAccClouddeployTarget_updateWithAttribution(enabledContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_first_label", "example-label-updated-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "labels.my_second_label", "example-label-updated-2"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.%", "4"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_first_label", "example-label-updated-1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.my_second_label", "example-label-updated-2"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "terraform_labels.goog-terraform-provisioned", "true"),

					resource.TestCheckResourceAttr("google_clouddeploy_target.primary", "effective_labels.%", "4"),
				),
			},
		},
	})
}

func testAccClouddeployTarget_createWithAttribution(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
  add_terraform_attribution_label               = %{add_attribution}
  terraform_attribution_label_addition_strategy = "%{attribution_strategy}"
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
    dns_endpoint = true
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

  associated_entities {
    entity_id = "test"
    anthos_clusters {
      membership = "projects/%{project_name}/locations/%{region}/memberships/membership-a"
    }

    gke_clusters {
      cluster     = "projects/%{project_name}/locations/%{region}/clusters/cluster-a"
      internal_ip = true
      proxy_url   = "http://10.0.0.1"
    }
  }
}
`, context)
}

func testAccClouddeployTarget_updateWithAttribution(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
  add_terraform_attribution_label               = %{add_attribution}
  terraform_attribution_label_addition_strategy = "%{attribution_strategy}"
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
    my_first_label = "example-label-updated-1"
    my_second_label = "example-label-updated-2"
  }
}
`, context)
}

func testAccClouddeployTarget_clearWithAttribution(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
  add_terraform_attribution_label               = %{add_attribution}
  terraform_attribution_label_addition_strategy = "%{attribution_strategy}"
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
}
`, context)
}
