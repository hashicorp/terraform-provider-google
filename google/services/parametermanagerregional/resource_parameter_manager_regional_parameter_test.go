// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanagerregional_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccParameterManagerRegionalRegionalParameter_import(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameter_import(context),
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-import",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
		},
	})
}

func testAccParameterManagerRegionalRegionalParameter_import(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-import" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "YAML"

  labels = {
    key1 = "val1"
    key2 = "val2"
    key3 = "val3"
    key4 = "val4"
    key5 = "val5"
  }
}
`, context)
}

func TestAccParameterManagerRegionalRegionalParameter_labelsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameter_withoutLabels(context),
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
			{
				Config: testAccParameterManagerRegionalRegionalParameter_labelsUpdate(context),
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
			{
				Config: testAccParameterManagerRegionalRegionalParameter_labelsUpdateOther(context),
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
			{
				Config: testAccParameterManagerRegionalRegionalParameter_withoutLabels(context),
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
		},
	})
}

func testAccParameterManagerRegionalRegionalParameter_withoutLabels(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-with-labels" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"
}
`, context)
}

func testAccParameterManagerRegionalRegionalParameter_labelsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-with-labels" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"

  labels = {
    key1 = "val1"
    key2 = "val2"
    key3 = "val3"
    key4 = "val4"
    key5 = "val5"
  }
}
`, context)
}

func testAccParameterManagerRegionalRegionalParameter_labelsUpdateOther(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_regional_parameter" "regional-parameter-with-labels" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"

  labels = {
    key1 = "val1"
    key2 = "updateval2"
    updatekey3 = "val3"
    updatekey4 = "updateval4"
    key6 = "val6"
  }
}
`, context)
}

func TestAccParameterManagerRegionalRegionalParameter_kmskeyUpdate(t *testing.T) {
	t.Parallel()

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-pm.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	context := map[string]interface{}{
		"kms_key":       acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-parameter-manager-managed-central-key1").CryptoKey.Name,
		"kms_key_other": acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-parameter-manager-managed-central-key2").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerRegionalRegionalParameterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerRegionalRegionalParameter_withoutKmsKey(context),
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-kms-key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
			{
				Config: testAccParameterManagerRegionalRegionalParameter_kmsKeyUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_parameter_manager_regional_parameter.regional-parameter-with-kms-key", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-kms-key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
			{
				Config: testAccParameterManagerRegionalRegionalParameter_kmsKeyUpdateOther(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_parameter_manager_regional_parameter.regional-parameter-with-kms-key", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-kms-key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
			{
				Config: testAccParameterManagerRegionalRegionalParameter_withoutKmsKey(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_parameter_manager_regional_parameter.regional-parameter-with-kms-key", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_parameter_manager_regional_parameter.regional-parameter-with-kms-key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "parameter_id", "terraform_labels"},
			},
		},
	})
}

func testAccParameterManagerRegionalRegionalParameter_withoutKmsKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_parameter_manager_regional_parameter" "regional-parameter-with-kms-key" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"
}
`, context)
}

func testAccParameterManagerRegionalRegionalParameter_kmsKeyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_parameter_manager_regional_parameter" "regional-parameter-with-kms-key" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"

  kms_key = "%{kms_key}"
}
`, context)
}

func testAccParameterManagerRegionalRegionalParameter_kmsKeyUpdateOther(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_parameter_manager_regional_parameter" "regional-parameter-with-kms-key" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  location = "us-central1"
  format = "JSON"

  kms_key = "%{kms_key_other}"
}
`, context)
}
