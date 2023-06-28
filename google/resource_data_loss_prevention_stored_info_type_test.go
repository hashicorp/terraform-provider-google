// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionStoredInfoTypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStart(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Displayname"

	regex {
		pattern = "patient"
		group_indexes = [2]
	}
}
`, context)
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Updated Description"
	display_name = "display_name"

	dictionary {
		word_list {
			words = ["word", "word2"]
		}
	}
}
`, context)
}

func TestAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeGroupIndexUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": envvar.GetTestProjectFromEnv(),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionStoredInfoTypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeWithoutGroupIndex(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeGroupIndexUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeWithoutGroupIndex(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeWithoutGroupIndex(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Displayname"

	regex {
		pattern = "patient"
	}
}
`, context)
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeGroupIndexUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Displayname"

	regex {
		pattern = "patient"
		group_indexes = [3]
	}
}
`, context)
}

func TestAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStoredInfoTypeId(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionStoredInfoTypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStoredInfoTypeId(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStoredInfoTypeIdUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStoredInfoTypeId(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Displayname"
	stored_info_type_id = "tf-test-%{random_suffix}"

	regex {
		pattern = "patient"
		group_indexes = [2]
	}
}
`, context)
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStoredInfoTypeIdUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Displayname"
	stored_info_type_id = "tf-test-update-%{random_suffix}"

	regex {
		pattern = "patient"
		group_indexes = [2]
	}
}
`, context)
}
