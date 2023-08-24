// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datafusion_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataFusionInstance_update(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataFusionInstance_basic(instanceName),
			},
			{
				ResourceName:      "google_data_fusion_instance.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataFusionInstance_updated(instanceName),
			},
			{
				ResourceName:      "google_data_fusion_instance.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataFusionInstance_basic(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name   = "%s"
  region = "us-central1"
  type   = "BASIC"
  # See supported versions here https://cloud.google.com/data-fusion/docs/support/version-support-policy
  version = "6.7.0"
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
  accelerators {
    accelerator_type = "CDC"
    state = "DISABLED"
  }
}
`, instanceName)
}

func testAccDataFusionInstance_updated(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name                          = "%s"
  region                        = "us-central1"
  type                          = "DEVELOPER"
  enable_stackdriver_monitoring = true
  enable_stackdriver_logging    = true

  labels = {
    label1 = "value1"
    label2 = "value2"
  }
  version = "6.8.0"

  accelerators {
    accelerator_type = "CCAI_INSIGHTS"
    state = "ENABLED"
  }
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
}
`, instanceName)
}

func TestAccDataFusionInstanceEnterprise_update(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataFusionInstanceEnterprise_basic(instanceName),
			},
			{
				ResourceName:      "google_data_fusion_instance.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataFusionInstanceEnterprise_updated(instanceName),
			},
			{
				ResourceName:      "google_data_fusion_instance.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataFusionInstanceEnterprise_basic(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name   = "%s"
  region = "us-central1"
  type   = "ENTERPRISE"
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
}
`, instanceName)
}

func testAccDataFusionInstanceEnterprise_updated(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name                          = "%s"
  region                        = "us-central1"
  type                          = "ENTERPRISE"
  enable_stackdriver_monitoring = true
  enable_stackdriver_logging    = true
  enable_rbac                   = true

  labels = {
    label1 = "value1"
    label2 = "value2"
  }
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
}
`, instanceName)
}

func TestAccDataFusionInstanceVersion_dataFusionInstanceUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"version":       "6.7.2",
	}

	contextUpdate := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"version":       "6.8.0",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataFusionInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataFusionInstanceVersion_dataFusionInstanceUpdate(context),
			},
			{
				ResourceName:            "google_data_fusion_instance.basic_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccDataFusionInstanceVersion_dataFusionInstanceUpdate(contextUpdate),
			},
			{
				ResourceName:            "google_data_fusion_instance.basic_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccDataFusionInstanceVersion_dataFusionInstanceUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_fusion_instance" "basic_instance" {
  name   = "tf-test-my-instance%{random_suffix}"
  region = "us-central1"
  type   = "BASIC"
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
    prober_test_run = "true"
  }
  version = "%{version}"
}
`, context)
}
