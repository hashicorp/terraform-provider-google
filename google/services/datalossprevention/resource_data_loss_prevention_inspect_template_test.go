// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datalossprevention_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataLossPreventionInspectTemplate_dlpInspectTemplateUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionInspectTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplateBasic(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplateUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplateBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
		info_types {
			name    = "PERSON_NAME"
			version = "latest"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}

		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					regex {
						pattern = ".+@example.com"
					}
					matching_type = "MATCHING_TYPE_FULL_MATCH"
				}
			}
		}
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}

		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				hotword_rule {
					hotword_regex {
						pattern = "patient"
					}
					proximity {
						window_before = 50
					}
					likelihood_adjustment {
						fixed_likelihood = "VERY_LIKELY"
					}
				}
			}
		}

		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name = "PERSON_NAME"
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplateUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Updated"
	display_name = "Different"

	inspect_config {
		info_types {
			name    = "PERSON_NAME"
			version = "stable"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}

		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}

		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				hotword_rule {
					hotword_regex {
						pattern = "not-a-patient"
					}
					proximity {
						window_before = 50
					}
					likelihood_adjustment {
						fixed_likelihood = "UNLIKELY"
					}
				}
			}
		}

		limits {
			max_findings_per_item    = 1
			max_findings_per_request = 5
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "PERSON_NAME"
				}
			}
			max_findings_per_info_type {
				max_findings = "20"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}

func TestAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withInfoTypesVersion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionInspectTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withInfoTypesVersionBasic(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withInfoTypesVersionUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withInfoTypesVersionBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		custom_info_types {
			info_type {
				name    = "MY_CUSTOM_TYPE"
				version = "0.1"
			}
	  
			likelihood = "UNLIKELY"
	  
			regex {
				pattern = "test*"
			}
		}
		info_types {
			name = "EMAIL_ADDRESS"
		}
		info_types {
			name    = "PERSON_NAME"
			version = "latest"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}

		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					matching_type = "MATCHING_TYPE_FULL_MATCH"
					exclude_info_types {
						info_types {
							name    = "EMAIL_ADDRESS"
							version = "0.1"
						}
						info_types {
							name    = "FIRST_NAME"
							version = "0.3"
						}
					}
				}
			}
		}
		rule_set {
			info_types {
				name    = "EMAIL_ADDRESS"
				version = "0.1"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name    = "PHONE_NUMBER"
				version = "0.4"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}

		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				hotword_rule {
					hotword_regex {
						pattern = "patient"
					}
					proximity {
						window_before = 50
					}
					likelihood_adjustment {
						fixed_likelihood = "VERY_LIKELY"
					}
				}
			}
		}

		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name    = "PERSON_NAME"
					version = "1.0"
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name    = "LAST_NAME"
					version = "0.5"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withInfoTypesVersionUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		custom_info_types {
			info_type {
				name    = "MY_CUSTOM_TYPE"
				version = "0.3"
			}
	  
			likelihood = "UNLIKELY"
	  
			regex {
				pattern = "test*"
			}
		}
		info_types {
			name = "EMAIL_ADDRESS"
		}
		info_types {
			name = "PERSON_NAME"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}

		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					matching_type = "MATCHING_TYPE_FULL_MATCH"
					exclude_info_types {
						info_types {
							name    = "EMAIL_ADDRESS"
							version = "0.5"
						}
						info_types {
							name = "FIRST_NAME"
						}
					}
				}
			}
		}
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name    = "PHONE_NUMBER"
				version = "0.5"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name    = "FIRST_NAME"
				version = "0.1"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}

		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				hotword_rule {
					hotword_regex {
						pattern = "patient"
					}
					proximity {
						window_before = 50
					}
					likelihood_adjustment {
						fixed_likelihood = "VERY_LIKELY"
					}
				}
			}
		}

		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name    = "PERSON_NAME"
					version = "1.4"
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}

func TestAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withExcludeByHotword(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionInspectTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withExcludeByHotwordBasic(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withExcludeByHotwordUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withExcludeByHotwordBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		custom_info_types {
			info_type {
				name = "MY_CUSTOM_TYPE"
			}
	  
			likelihood = "UNLIKELY"
	  
			regex {
				pattern = "test*"
			}
		}
		info_types {
			name = "EMAIL_ADDRESS"
		}
		info_types {
			name    = "PERSON_NAME"
			version = "latest"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}

		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					matching_type = "MATCHING_TYPE_FULL_MATCH"
					exclude_info_types {
						info_types {
							name = "EMAIL_ADDRESS"
						}
						info_types {
							name = "FIRST_NAME"
						}
					}
				}
			}
		}
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}
		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				exclusion_rule {
					matching_type = "MATCHING_TYPE_FULL_MATCH"
					exclude_by_hotword {
						hotword_regex {
							pattern = "test*"
						}
						proximity {
							window_before = 12
							window_after  = 14
						}
					}
				}
			}
		}

		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				hotword_rule {
					hotword_regex {
						pattern = "patient"
					}
					proximity {
						window_before = 50
					}
					likelihood_adjustment {
						fixed_likelihood = "VERY_LIKELY"
					}
				}
			}
		}

		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name = "PERSON_NAME"
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withExcludeByHotwordUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		custom_info_types {
			info_type {
				name = "MY_CUSTOM_TYPE"
			}
	  
			likelihood = "UNLIKELY"
	  
			regex {
				pattern = "test*"
			}
		}
		info_types {
			name = "EMAIL_ADDRESS"
		}
		info_types {
			name    = "PERSON_NAME"
			version = "latest"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}

		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					matching_type = "MATCHING_TYPE_FULL_MATCH"
					exclude_info_types {
						info_types {
							name = "EMAIL_ADDRESS"
						}
						info_types {
							name = "FIRST_NAME"
						}
					}
				}
			}
		}
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}
		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				exclusion_rule {
					matching_type = "MATCHING_TYPE_FULL_MATCH"
					exclude_by_hotword {
						hotword_regex {
							pattern       = "updatetest*"
							group_indexes = [2]
						}
						proximity {
							window_before = 2
							window_after  = 4
						}
					}
				}
			}
		}

		rule_set {
			info_types {
				name = "PERSON_NAME"
			}
			rules {
				hotword_rule {
					hotword_regex {
						pattern = "patient"
					}
					proximity {
						window_before = 50
					}
					likelihood_adjustment {
						fixed_likelihood = "VERY_LIKELY"
					}
				}
			}
		}

		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name = "PERSON_NAME"
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}

func TestAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withSensitivityScore(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionInspectTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withSensitivityScoreBasic(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withSensitivityScoreUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withSensitivityScoreUpdate2(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_inspect_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withSensitivityScoreBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent       = "projects/%{project}"
	description  = "Description"
	display_name = "Display"
	
	inspect_config {
		custom_info_types {
			info_type {
				name = "MY_CUSTOM_TYPE"
				sensitivity_score {
					score = "SENSITIVITY_MODERATE"
				}
			}
			sensitivity_score {
				score = "SENSITIVITY_HIGH"
			}
			surrogate_type {}
		}
		info_types {
			name = "EMAIL_ADDRESS"
			sensitivity_score {
				score = "SENSITIVITY_MODERATE"
			}
		}
		info_types {
			name    = "PERSON_NAME"
			version = "latest"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}
	
		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					exclude_info_types {
						info_types {
							name = "LAST_NAME"
							sensitivity_score {
								score = "SENSITIVITY_LOW"
							}
						}
					}
					matching_type = "MATCHING_TYPE_FULL_MATCH"
				}
			}
		}
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
				sensitivity_score {
					score = "SENSITIVITY_LOW"
				}
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}
	
		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name = "PERSON_NAME"
					sensitivity_score {
						score = "SENSITIVITY_HIGH"
					}
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withSensitivityScoreUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent       = "projects/%{project}"
	description  = "Description"
	display_name = "Display"
	
	inspect_config {
		custom_info_types {
			info_type {
				name = "MY_CUSTOM_TYPE"
			}
			sensitivity_score {
				score = "SENSITIVITY_LOW"
			}
			surrogate_type {}
		}
		info_types {
			name = "EMAIL_ADDRESS"
			sensitivity_score {
				score = "SENSITIVITY_LOW"
			}
		}
		info_types {
			name    = "PERSON_NAME"
			version = "latest"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}
	
		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					exclude_info_types {
						info_types {
							name = "LAST_NAME"
						}
					}
					matching_type = "MATCHING_TYPE_FULL_MATCH"
				}
			}
		}
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}
	
		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name = "PERSON_NAME"
					sensitivity_score {
						score = "SENSITIVITY_MODERATE"
					}
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionInspectTemplate_dlpInspectTemplate_withSensitivityScoreUpdate2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent       = "projects/%{project}"
	description  = "Description"
	display_name = "Display"
	
	inspect_config {
		custom_info_types {
			info_type {
				name = "MY_CUSTOM_TYPE"
				sensitivity_score {
					score = "SENSITIVITY_LOW"
				}
			}
			surrogate_type {}
		}
		info_types {
			name = "EMAIL_ADDRESS"
			sensitivity_score {
				score = "SENSITIVITY_HIGH"
			}
		}
		info_types {
			name    = "PERSON_NAME"
			version = "latest"
		}
		info_types {
			name = "LAST_NAME"
		}
		info_types {
			name = "DOMAIN_NAME"
		}
		info_types {
			name = "PHONE_NUMBER"
		}
		info_types {
			name = "FIRST_NAME"
		}
	
		min_likelihood = "UNLIKELY"
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			rules {
				exclusion_rule {
					exclude_info_types {
						info_types {
							name = "LAST_NAME"
							sensitivity_score {
								score = "SENSITIVITY_HIGH"
							}
						}
					}
					matching_type = "MATCHING_TYPE_FULL_MATCH"
				}
			}
		}
		rule_set {
			info_types {
				name = "EMAIL_ADDRESS"
			}
			info_types {
				name = "DOMAIN_NAME"
			}
			info_types {
				name = "PHONE_NUMBER"
			}
			info_types {
				name = "PERSON_NAME"
			}
			info_types {
				name = "FIRST_NAME"
			}
			rules {
				exclusion_rule {
					dictionary {
						word_list {
							words = ["TEST"]
						}
					}
					matching_type = "MATCHING_TYPE_PARTIAL_MATCH"
				}
			}
		}
	
		limits {
			max_findings_per_item    = 10
			max_findings_per_request = 50
			max_findings_per_info_type {
				max_findings = "75"
				info_type {
					name = "PERSON_NAME"
					sensitivity_score {
						score = "SENSITIVITY_LOW"
					}
				}
			}
			max_findings_per_info_type {
				max_findings = "80"
				info_type {
					name = "LAST_NAME"
				}
			}
		}
	}
}
`, context)
}
