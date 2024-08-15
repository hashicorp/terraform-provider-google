// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquerydatapolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDatapolicyDataPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyBasicExample(context),
			},
			{
				ResourceName:            "google_bigquery_datapolicy_data_policy.data_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyUpdate(context),
			},
		},
	})
}

func TestAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyRoutineUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDatapolicyDataPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyRoutineExample(context),
			},
			{
				ResourceName:            "google_bigquery_datapolicy_data_policy.data_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyRoutineUpdate(context),
			},
		},
	})
}

func testAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_datapolicy_data_policy" "data_policy" {
    location         = "us-central1"
    data_policy_id   = "tf_test_data_policy%{random_suffix}"
    policy_tag       = google_data_catalog_policy_tag.policy_tag_updated.name
    data_policy_type = "COLUMN_LEVEL_SECURITY_POLICY"
  }

  resource "google_data_catalog_policy_tag" "policy_tag" {
    taxonomy     = google_data_catalog_taxonomy.taxonomy.id
    display_name = "Low security"
    description  = "A policy tag normally associated with low security items"
  }

  resource "google_data_catalog_policy_tag" "policy_tag_updated" {
    taxonomy     = google_data_catalog_taxonomy.taxonomy.id
    display_name = "Low security updated"
    description  = "A policy tag normally associated with low security items"
  }  

  resource "google_bigquery_datapolicy_data_policy" "policy_tag_with_data_masking_policy" {
    location         = "us-central1"
    data_policy_id   = "masking_policy_test"
    policy_tag       = google_data_catalog_policy_tag.policy_tag_updated.name
    data_policy_type = "DATA_MASKING_POLICY"
    data_masking_policy {
        predefined_expression = "SHA256"
    }
  }

  resource "google_data_catalog_taxonomy" "taxonomy" {
    region                 = "us-central1"
    display_name           = "taxonomy%{random_suffix}"
    description            = "A collection of policy tags"
    activated_policy_types = ["FINE_GRAINED_ACCESS_CONTROL"]
  }
`, context)
}

func testAccBigqueryDatapolicyDataPolicy_bigqueryDatapolicyDataPolicyRoutineUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_datapolicy_data_policy" "data_policy" {
  location         = "us-central1"
  data_policy_id   = "tf_test_data_policy%{random_suffix}"
  policy_tag       = google_data_catalog_policy_tag.policy_tag.name
  data_policy_type = "DATA_MASKING_POLICY"  
  data_masking_policy {
    routine = google_bigquery_routine.custom_masking_routine_updated.id
  }
}

resource "google_data_catalog_policy_tag" "policy_tag" {
  taxonomy     = google_data_catalog_taxonomy.taxonomy.id
  display_name = "Low security"
  description  = "A policy tag normally associated with low security items"
}
  
resource "google_data_catalog_taxonomy" "taxonomy" {
  region                 = "us-central1"
  display_name           = "taxonomy%{random_suffix}"
  description            = "A collection of policy tags"
  activated_policy_types = ["FINE_GRAINED_ACCESS_CONTROL"]
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "tf_test_dataset_id%{random_suffix}"
  location   = "us-central1"
}

resource "google_bigquery_routine" "custom_masking_routine" {
	dataset_id           = google_bigquery_dataset.test.dataset_id
	routine_id           = "custom_masking_routine"
	routine_type         = "SCALAR_FUNCTION"
	language             = "SQL"
	data_governance_type = "DATA_MASKING"
	definition_body      = "SAFE.REGEXP_REPLACE(ssn, '[0-9]', 'X')"
	return_type          = "{\"typeKind\" :  \"STRING\"}"

	arguments {
	  name = "ssn"
	  data_type = "{\"typeKind\" :  \"STRING\"}"
	} 
}

resource "google_bigquery_routine" "custom_masking_routine_updated" {
	dataset_id           = google_bigquery_dataset.test.dataset_id
	routine_id           = "custom_masking_routine_update"
	routine_type         = "SCALAR_FUNCTION"
	language             = "SQL"
	data_governance_type = "DATA_MASKING"
	definition_body      = "SAFE.REGEXP_REPLACE(ssn, '[0-9]', 'X')"
	return_type          = "{\"typeKind\" :  \"STRING\"}"

	arguments {
	  name = "ssn"
	  data_type = "{\"typeKind\" :  \"STRING\"}"
	} 
}
`, context)
}
