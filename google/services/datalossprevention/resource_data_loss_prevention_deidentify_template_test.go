// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datalossprevention_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformationsStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformationsUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformationsStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "PHONE_NUMBER"
        }
        info_types {
          name = "CREDIT_CARD_NUMBER"
          sensitivity_score {
            score = "SENSITIVITY_MODERATE"
          }
        }

        primitive_transformation {
          replace_config {
            new_value {
              integer_value = 9
            }
          }
        }
      }

      transformations {
        info_types {
          name = "EMAIL_ADDRESS"
        }
        info_types {
          name = "LAST_NAME"
        }

        primitive_transformation {
          character_mask_config {
            masking_character = "X"
            number_to_mask = 4
            reverse_order = true
            characters_to_ignore {
              common_characters_to_ignore = "PUNCTUATION"
            }
          }
        }
      }

      transformations {
        info_types {
          name = "DATE_OF_BIRTH"
        }

        primitive_transformation {
          replace_config {
            new_value {
              date_value {
                year  = 2020
                month = 1
                day   = 1
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "CREDIT_CARD_NUMBER2322"
        }

        primitive_transformation {
          crypto_deterministic_config {
            context {
              name = "sometweak"
            }
            crypto_key {
              kms_wrapped {
                wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                crypto_key_name = "%{kms_key_name}"
              }
            }
            surrogate_info_type {
              name = "abc"
              version = "version-1"
            }
          }
        }
      }

      transformations {
        info_types {
          name = "CREDIT_CARD_NUMBER23224"
        }

        primitive_transformation {
          crypto_deterministic_config {
            context {
              name = "sometweak"
            }
            crypto_key {
              unwrapped {
                key     = "VVdWVWFGZHRXbkUwZERkM0lYb2xRdz09"
              }
            }
            surrogate_info_type {
              name = "abc"
              version = "version-1"
            }
          }
        }
      }

      transformations {
        info_types {
          name = "CUSTOM_INFO_TYPE"
        }

        primitive_transformation {
          crypto_deterministic_config {
            crypto_key {
              transient {
                name = "beep"
              }
            }
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
              version = "version-1"
              sensitivity_score {
                score = "SENSITIVITY_MODERATE"
              }
            }
          }
        }
      }
      transformations {
        info_types {
          name = "PHONE_NUMBER2"
        }
        primitive_transformation {
          crypto_replace_ffx_fpe_config {
            context {
              name = "someTweak"
            }
            crypto_key {
              kms_wrapped {
                wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                crypto_key_name = "%{kms_key_name}"
              }
            }
            radix = 10
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
              version = "version-1"
              sensitivity_score {
                score = "SENSITIVITY_LOW"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "SSN"
        }
        primitive_transformation {
          crypto_replace_ffx_fpe_config {
            common_alphabet = "UPPER_CASE_ALPHA_NUMERIC"
            context {
              name = "someTweak"
            }
            crypto_key {
              transient {
                name = "beep"
              }
            }
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
              version = "version-1"
            }
          }
        }
      }

      transformations {
        info_types {
          name = "SSN33"
        }
        primitive_transformation {
          crypto_replace_ffx_fpe_config {
            common_alphabet = "UPPER_CASE_ALPHA_NUMERIC"
            crypto_key {
              unwrapped {
                key = "VVdWVWFGZHRXbkUwZERkM0lYb2xRdz09"
              }
            }
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
              version = "version-1"
            }
          }
        }
      }

      transformations {
        info_types {
          name = "RDC_EXAMPLE"
        }
        primitive_transformation {
          replace_dictionary_config {
            word_list {
              words = [
                "foo",
                "bar",
                "baz",
              ]
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformationsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "CREDIT_CARD_NUMBER"
          sensitivity_score {
            score = "SENSITIVITY_HIGH"
          }
        }

        primitive_transformation {
          replace_config {
            new_value {
              integer_value = 9
            }
          }
        }
      }

      transformations {
        info_types {
          name = "EMAIL_ADDRESS"
          version = "0.1"
        }
        info_types {
          name = "LAST_NAME"
        }

        primitive_transformation {
          character_mask_config {
            number_to_mask = 3
            reverse_order = true
          }
        }
      }

      transformations {
        info_types {
          name = "DATE_OF_BIRTH"
        }

        primitive_transformation {
          replace_config {
            new_value {
              date_value {
                year  = 2020
                month = 1
                day   = 1
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "CREDIT_CARD_NUMBERR"
        }

        primitive_transformation {
          crypto_deterministic_config {
            context {
              name = "sometweakd"
            }
            crypto_key {
              kms_wrapped {
                wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                crypto_key_name = "%{kms_key_name}"
              }
            }
            surrogate_info_type {
              name = "abcd"
              version = "version-2"
            }
          }
        }
      }
      transformations {
        info_types {
          name = "CUSTOM_INFO_TYPE"
        }

        primitive_transformation {
          crypto_deterministic_config {
            crypto_key {
              transient {
                name = "beeper"
              }
            }
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPEf"
              version = "version-2"
              sensitivity_score {
                score="SENSITIVITY_LOW"
              }
            }
          }
        }
      }
      transformations {
        info_types {
          name = "PHONE_NUMBER2"
        }
        primitive_transformation {
          crypto_replace_ffx_fpe_config {
            context {
              name = "someTweaker"
            }
            crypto_key {
              kms_wrapped {
                wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                crypto_key_name = "%{kms_key_name}"
              }
            }
            radix = 10
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPEF"
              version = "version-2"
              sensitivity_score {
                score = "SENSITIVITY_MODERATE"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "SSN"
        }
        primitive_transformation {
          crypto_replace_ffx_fpe_config {
            common_alphabet = "UPPER_CASE_ALPHA_NUMERIC"
            context {
              name = "someTweak2"
            }
            crypto_key {
              transient {
                name = "beepf"
              }
            }
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
              version = "version-2"
            }
          }
        }
      }

      transformations {
        info_types {
          name = "RDC_EXAMPLE"
        }
        primitive_transformation {
          replace_dictionary_config {
            word_list {
              words = [
                # update list - deletion and addition
                "foo",
                "baz",
                "fizz",
                "buzz",
              ]
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_start(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_update(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_start(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field3"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "fizzbuzz"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "details.pii.email"
        }
        condition {
          expressions {
            conditions {
              conditions {
                field {
                  name = "details.pii.country_code"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "US"
                }
              }
              conditions {
                field {
                  name = "details.pii.date_of_birth"
                }
                operator = "GREATER_THAN_OR_EQUALS"
                value {
                  date_value {
                    year = 2001
                    month = 6
                    day = 29
                  }
                }
              }
            }
          }
        }
        primitive_transformation {
          replace_config {
            new_value {
              string_value = "born.after.shrek@example.com"
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-redacted-field"
        }
        primitive_transformation {
          redact_config {}
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-char-masked-field"
        }
        primitive_transformation {
          character_mask_config {
            masking_character = "x"
            number_to_mask = 8
            characters_to_ignore {
              characters_to_skip = "-"
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-replace-ffx-fpe-field"
        }
        primitive_transformation {
          crypto_replace_ffx_fpe_config {
            context {
              name = "someTweak"
            }
            crypto_key {
              kms_wrapped {
                wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                crypto_key_name = "%{kms_key_name}"
              }
            }
            radix = 10
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
              version = "version-1"
              sensitivity_score {
                score = "SENSITIVITY_LOW"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-fixed-size-bucketing-field"
        }
        primitive_transformation {
          fixed_size_bucketing_config {
            lower_bound {
              integer_value = 0
            }
            upper_bound {
              integer_value = 100
            } 
            bucket_size = 10
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-field"
        }
        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                string_value = "00:00:00"
              }
              max {
                string_value = "11:59:59"
              }
              replacement_value {
                string_value = "AM"
              }
            }
            buckets {
              min {
                string_value = "12:00:00"
              }
              max {
                string_value = "23:59:59"
              }
              replacement_value {
                string_value = "PM"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-time-part-field"
        }
        primitive_transformation {
          time_part_config {
            part_to_extract = "YEAR"
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-hash-field"
        }
        primitive_transformation {
          crypto_hash_config {
            crypto_key {
              transient {
                name = "beep" # Copy-pasting from existing test that uses this field
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-date-shift-field"
        }
        primitive_transformation {
          date_shift_config {
            upper_bound_days = 30
            lower_bound_days = -30
            context {
              name = "unconditionally-date-shift-field"
            }
            crypto_key {
              transient {
                name = "beep"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-deterministic-field"
        }
        primitive_transformation {
          crypto_deterministic_config {
            crypto_key {
              transient {
                name = "beep"
              }
            }
            surrogate_info_type {
              name = "CREDIT_CARD_NUMBER"
              version = "version-1"
              sensitivity_score {
                score = "SENSITIVITY_HIGH"
              }
            }
            context {
              name = "unconditionally-crypto-deterministic-field"
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-dictionary-field"
        }
        primitive_transformation {
          replace_dictionary_config {
            word_list {
              words = [
                "foo",
                "bar",
                "baz",
              ]
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field3"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR-updated"
                }
              }

              # update includes deleting condition affecting field2

              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "fizzbuzz-updated"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "details.pii.email"
        }
        condition {
          expressions {
            conditions {
              # update to remove condition checking the details.pii.country_code field
              # update to add a new condition
              conditions {
                field {
                  name = "details.pii.gender"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "M"
                }
              }
              conditions {
                field {
                  name = "details.pii.date_of_birth"
                }
                operator = "GREATER_THAN_OR_EQUALS"
                value {
                  # update date values
                  date_value {
                    year = 2004
                    month = 7
                    day = 2
                  }
                }
              }
            }
          }
        }
        primitive_transformation {
          # update values inside replace_config
          replace_config {
            new_value {
              string_value = "dude.born.after.shrek2@example.com"
            }
          }
        }
      }

      # update to remove field_transformations block using redact_config

      field_transformations {
        fields {
          name = "unconditionally-char-masked-field"
        }
        primitive_transformation {
          character_mask_config {
            masking_character = "x"
            number_to_mask = 8
            # update to delete old characters_to_ignore block and add new ones
            characters_to_ignore {
              common_characters_to_ignore = "PUNCTUATION"
            }
            characters_to_ignore {
              common_characters_to_ignore = "ALPHA_UPPER_CASE"
            }
            characters_to_ignore {
              common_characters_to_ignore = "ALPHA_LOWER_CASE"
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-replace-ffx-fpe-field"
        }
        primitive_transformation {
          crypto_replace_ffx_fpe_config {
            common_alphabet = "UPPER_CASE_ALPHA_NUMERIC"
            context {
              name = "someTweak2"
            }
            crypto_key {
              transient {
                name = "beep" # Copy-pasting from existing test that uses this field
              }
            }
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
              version = "version-2"
              sensitivity_score {
                score = "SENSITIVITY_MODERATE"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-fixed-size-bucketing-field"
        }
        primitive_transformation {
          # update values
          fixed_size_bucketing_config {
            lower_bound {
              integer_value = 0
            }
            upper_bound {
              integer_value = 200
            } 
            bucket_size = 20
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-field"
        }
        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                string_value = "00:00:00"
              }
              max {
                string_value = "11:59:59"
              }
              replacement_value {
                string_value = "AM"
              }
            }
            # Add new bucket
            buckets {
              min {
                string_value = "12:00:00"
              }
              max {
                string_value = "13:59:59"
              }
              replacement_value {
                string_value = "Lunchtime"
              }
            }
            buckets {
              min {
                string_value = "14:00:00"
              }
              max {
                string_value = "23:59:59"
              }
              replacement_value {
                string_value = "PM"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-time-part-field"
        }
        primitive_transformation {
          time_part_config {
            part_to_extract = "MONTH"
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-hash-field"
        }
        primitive_transformation {
          crypto_hash_config {
            crypto_key {
              transient {
                # update value
                name = "beepy-beep-updated"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-date-shift-field"
        }
        primitive_transformation {
          date_shift_config {
            # update values
            upper_bound_days = 60
            lower_bound_days = -60
            context {
              name = "unconditionally-date-shift-field"
            }
            crypto_key {
              transient {
                # update value
                name = "beepy-beep-updated"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-deterministic-field"
        }
        primitive_transformation {
          crypto_deterministic_config {
            crypto_key {
              transient {
                # update value
                name = "beepy-beep-updated"
              }
            }
            surrogate_info_type {
              # update info type
              name = "CREDIT_CARD_TRACK_NUMBER"
              version = "version-2"
              sensitivity_score {
                score = "SENSITIVITY_LOW"
              }
            }
            context {
              name = "unconditionally-crypto-deterministic-field"
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-dictionary-field"
        }
        primitive_transformation {
          replace_dictionary_config {
            word_list {
              words = [
                # update list - deletion and addition
                "foo",
                "baz",
                "fizz",
                "buzz",
              ]
            }
          }
        }
      }

    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_imageTransformationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_imageTransformationsBasic(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_imageTransformationsUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_imageTransformationsBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    image_transformations {
      transforms {
        redaction_color {
          red = 0.5
          blue = 1
          green = 0.2
        }
        selected_info_types {
          info_types {
            name = "COLOR_INFO"
            version = "latest"
            sensitivity_score {
              score = "SENSITIVITY_LOW"
            }
          }
        }
      }

      transforms {
        all_info_types {}
      }

      transforms {
        all_text {}
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_imageTransformationsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    image_transformations {
      transforms {
        redaction_color {
          red = 0.3
          blue = 0.5
          green = 0.9
        }
        selected_info_types {
          info_types {
            name = "COLOR_EXAMPLE"
            version = "0.1"
            sensitivity_score {
              score = "SENSITIVITY_MODERATE"
            }
          }
        }
      }
      # Update allInfoTypes and allText by removing the block
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformationsStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformationsUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformationsStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "PHONE_NUMBER"
        }
        info_types {
          name = "DATE_SHIFT_EXAMPLE"
        }

        primitive_transformation {
          date_shift_config {
            upper_bound_days = 30
            lower_bound_days = -30
            context {
              name = "DATE_SHIFT_EXAMPLE"
            }
            crypto_key {
              transient {
                name = "beep"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "EMAIL_ADDRESS"
        }
        info_types {
          name = "FIXED_BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          fixed_size_bucketing_config {
            lower_bound {
              integer_value = 0
            }
            upper_bound {
              integer_value = 100
            }
            bucket_size = 10
          }
        }
      }

      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                string_value = "00:00:00"
              }
              max {
                string_value = "11:59:59"
              }
              replacement_value {
                string_value = "AM"
              }
            }
            buckets {
              min {
                string_value = "12:00:00"
              }
              max {
                string_value = "23:59:59"
              }
              replacement_value {
                string_value = "PM"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "TIME_PART_EXAMPLE"
        }

        primitive_transformation {
          time_part_config {
            part_to_extract = "YEAR"
          }
        }
      }

      transformations {
        info_types {
          name = "CRYPTO_HASH_TRANSIENT_EXAMPLE"
        }

        primitive_transformation {
          crypto_hash_config {
            crypto_key {
              transient {
                name = "beep" # Copy-pasting from existing test that uses this field
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "CRYPTO_HASH_UNWRAPPED_EXAMPLE"
        }

        primitive_transformation {
          crypto_hash_config {
            crypto_key {
              unwrapped {
                key     = "VVdWVWFGZHRXbkUwZERkM0lYb2xRdz09"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "REDACT_EXAMPLE"
        }

        primitive_transformation {
          redact_config {}
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformationsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "DATE_SHIFT_EXAMPLE"
        }

        primitive_transformation {
          date_shift_config {
            # update values
            upper_bound_days = 60
            lower_bound_days = -60
            context {
              name = "DATE_SHIFT_EXAMPLE"
            }
            crypto_key {
              transient {
                # update value
                name = "beepy-beep-updated"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "EMAIL_ADDRESS"
        }
        info_types {
          name = "FIXED_BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          # update values
          fixed_size_bucketing_config {
            lower_bound {
              integer_value = 0
            }
            upper_bound {
              integer_value = 200
            }
            bucket_size = 20
          }
        }
      }

      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                string_value = "00:00:00"
              }
              max {
                string_value = "11:59:59"
              }
              replacement_value {
                string_value = "AM"
              }
            }
            # Add new bucket
            buckets {
              min {
                string_value = "12:00:00"
              }
              max {
                string_value = "13:59:59"
              }
              replacement_value {
                string_value = "Lunchtime"
              }
            }
            buckets {
              min {
                string_value = "14:00:00"
              }
              max {
                string_value = "23:59:59"
              }
              replacement_value {
                string_value = "PM"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "TIME_PART_EXAMPLE"
        }

        primitive_transformation {
          time_part_config {
            part_to_extract = "MONTH"
          }
        }
      }

      transformations {
        info_types {
          name = "CRYPTO_HASH_TRANSIENT_UPDATED_EXAMPLE"
        }

        primitive_transformation {
          crypto_hash_config {
            crypto_key {
              transient {
                # update value
                name = "beepy-beep-updated"
              }
            }
          }
        }
      }

      transformations {
        info_types {
          name = "CRYPTO_HASH_WRAPPED_EXAMPLE"
        }

        primitive_transformation {
          crypto_hash_config {
            crypto_key {
              kms_wrapped {
                wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                crypto_key_name = "%{kms_key_name}"
              }
            }
          }
        }
      }

      # update to remove transformations block using redact_config
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_integerValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_floatValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_timestampValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_timeValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_dateValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_dayOfWeekValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_integerValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                integer_value = 921
              }
              max {
                integer_value = 3010
              }
              replacement_value {
                integer_value = 1212
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_floatValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                float_value = 10.50
              }
              max {
                float_value = 310.75
              }
              replacement_value {
                float_value = 5.37
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_timestampValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                timestamp_value = "2014-10-02T15:01:23Z"
              }
              max {
                timestamp_value = "2015-06-29T18:46:39Z"
              }
              replacement_value {
                timestamp_value = "2014-12-24T09:19:50Z"
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_timeValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                time_value {
                  hours   = 09
                  minutes = 30
                  seconds = 45
                  nanos   = 123412
                }
              }
              max {
                time_value {
                  hours   = 15
                  minutes = 45
                  seconds = 00
                  nanos   = 523278
                }
              }
              replacement_value {
                time_value {
                  hours   = 23
                  minutes = 59
                  seconds = 59
                  nanos   = 999999
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_dateValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                date_value {
                  year = 1969
                  month = 11
                  day = 23
                }
              }
              max {
                date_value {
                  year = 2010
                  month = 12
                  day = 31
                }
              }
              replacement_value {
                date_value {
                  year = 2011
                  month = 05
                  day = 19
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_bucketingConfig_dayOfWeekValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          bucketing_config {
            buckets {
              min {
                day_of_week_value = "FRIDAY"
              }
              max {
                day_of_week_value = "SUNDAY"
              }
              replacement_value {
                day_of_week_value = "MONDAY"
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_fixedSizeBucketingConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_fixedSizeBucketingConfig_integerValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_fixedSizeBucketingConfig_floatValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_fixedSizeBucketingConfig_integerValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "FIXED_BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          fixed_size_bucketing_config {
            lower_bound {
              integer_value = 0
            }
            upper_bound {
              integer_value = 200
            }
            bucket_size = 20
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_fixedSizeBucketingConfig_floatValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "FIXED_BUCKETING_EXAMPLE"
        }

        primitive_transformation {
          fixed_size_bucketing_config {
            lower_bound {
              float_value = 0.5
            }
            upper_bound {
              float_value = 20.5
            }
            bucket_size = 20
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_dateShiftConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_dateShiftConfig_transient(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_dateShiftConfig_unwrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_dateShiftConfig_kmsWrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_dateShiftConfig_transient(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "DATE_SHIFT_EXAMPLE"
        }

        primitive_transformation {
          date_shift_config {
            upper_bound_days = 30
            lower_bound_days = -30
            context {
              name = "some-context-field"
            }
            crypto_key {
              transient {
                name = "someRandomTerraformKey"
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_dateShiftConfig_unwrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "DATE_SHIFT_EXAMPLE"
        }

        primitive_transformation {
          date_shift_config {
            upper_bound_days = 30
            lower_bound_days = -30
            context {
              name = "some-context-field"
            }
            crypto_key {
              unwrapped {
                key = "0836c61118ac590243bdadb25f0bb08e"
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformations_primitiveTransformations_dateShiftConfig_kmsWrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "config" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "DATE_SHIFT_EXAMPLE"
        }

        primitive_transformation {
          date_shift_config {
            upper_bound_days = 30
            lower_bound_days = -30
            context {
              name = "some-context-field"
            }
            crypto_key {
              kms_wrapped {
                wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                crypto_key_name = "%{kms_key_name}"
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization": envvar.GetTestOrgFromEnv(t),
		"kms_key_name": acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_start(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_update(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_start(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field3"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "fizzbuzz"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "details.pii.email"
        }
        condition {
          expressions {
            conditions {
              conditions {
                field {
                  name = "details.pii.country_code"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "US"
                }
              }
              conditions {
                field {
                  name = "details.pii.date_of_birth"
                }
                operator = "GREATER_THAN_OR_EQUALS"
                value {
                  date_value {
                    year = 2001
                    month = 6
                    day = 29
                  }
                }
              }
            }
          }
        }
        info_type_transformations {
          transformations {   
            info_types {
              name = "PHONE_NUMBER"
              version = "0.1"
            }
            info_types {
              name = "CREDIT_CARD_NUMBER"
              version = "1.2"
              sensitivity_score {
                score = "SENSITIVITY_HIGH"
              }
            } 
            primitive_transformation {
              replace_config {
                new_value {
                  integer_value = 9
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-redacted-field"
        }
        info_type_transformations {
          transformations {    
            primitive_transformation {
              redact_config {}
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-char-masked-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "EMAIL_ADDRESS"
              version = "latest"
            }
            info_types {
              name = "LAST_NAME"
            }
            primitive_transformation {
              character_mask_config {
                masking_character = "x"
                number_to_mask = 8
                characters_to_ignore {
                  characters_to_skip = "-"
                }
                reverse_order = true
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-replace-ffx-fpe-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "SSN"
            }
            primitive_transformation {
              crypto_replace_ffx_fpe_config {
                context {
                  name = "someTweak"
                }
                crypto_key {
                  kms_wrapped {
                    wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                    crypto_key_name = "%{kms_key_name}"
                  }
                }
                radix = 10
                surrogate_info_type {
                  name = "CUSTOM_INFO_TYPE"
                  version = "version-1"
                  sensitivity_score {
                    score = "SENSITIVITY_LOW"
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-fixed-size-bucketing-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "AGE"
            }
            primitive_transformation {
              fixed_size_bucketing_config {
                lower_bound {
                  integer_value = 0
                }
                upper_bound {
                  integer_value = 100
                } 
                bucket_size = 10
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-field"
        }
        info_type_transformations {
          transformations {  
            info_types {
              name = "CREATED_TIME"
            }  
            primitive_transformation {
              bucketing_config {
                buckets {
                  min {
                    string_value = "00:00:00"
                  }
                  max {
                    string_value = "11:59:59"
                  }
                  replacement_value {
                    string_value = "AM"
                  }
                }
                buckets {
                  min {
                    string_value = "12:00:00"
                  }
                  max {
                    string_value = "23:59:59"
                  }
                  replacement_value {
                    string_value = "PM"
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-time-part-field"
        }
        info_type_transformations {
          transformations { 
            info_types {
              name = "DATE_OF_BIRTH"
            }   
            primitive_transformation {
              time_part_config {
                part_to_extract = "YEAR"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-hash-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "CREDIT_CARD_SECRET"
            } 
            primitive_transformation {
              crypto_hash_config {
                crypto_key {
                  transient {
                    name = "beep" # Copy-pasting from existing test that uses this field
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-date-shift-field"
        }
        info_type_transformations {
          transformations {   
            info_types {
              name = "EXTRACT_DATE"
            }  
            primitive_transformation {
              date_shift_config {
                upper_bound_days = 30
                lower_bound_days = -30
                context {
                  name = "unconditionally-date-shift-field"
                }
                crypto_key {
                  transient {
                    name = "beep"
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-deterministic-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "CREDIT_CARD_SECRET1234"
            } 
            primitive_transformation {
              crypto_deterministic_config {
                crypto_key {
                  transient {
                    name = "beep"
                  }
                }
                surrogate_info_type {
                  name = "CREDIT_CARD_NUMBER"
                  version = "version-1"
                  sensitivity_score {
                    score = "SENSITIVITY_LOW"
                  }
                }
                context {
                  name = "unconditionally-crypto-deterministic-field"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-dictionary-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "RANDOM_FIELD"
            } 
            primitive_transformation {
              replace_dictionary_config {
                word_list {
                  words = [
                    "foo",
                    "bar",
                    "baz",
                  ]
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-with-info-type-config"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "FIRST_NAME"
            }
            primitive_transformation {
              replace_with_info_type_config {}
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field3"
                }
                operator = "EQUAL_TO"
                value {

                  # update the condition for field3

                  string_value = "FOO-BAR-updated"
                }
              }

              # update includes deleting condition affecting field2

              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {

                  # update the condition for field1

                  string_value = "fizzbuzz-updated"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "details.pii.email"
        }
        condition {
          expressions {
            conditions {
              conditions {

                # update to remove condition checking the details.pii.country_code field
                # update to add a new condition
                
                field {
                  name = "details.pii.gender"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "M"
                }
              }
              conditions {
                field {
                  name = "details.pii.date_of_birth"
                }
                operator = "GREATER_THAN_OR_EQUALS"
                value {

                  # update date values

                  date_value {
                    year = 2007
                    month = 11
                    day = 9
                  }
                }
              }
            }
          }
        }
        info_type_transformations {
          transformations {    
            
            # removing the info_types

            info_types {
              name = "CREDIT_CARD_NUMBER"
              version = "1.5"
              sensitivity_score {
                score = "SENSITIVITY_MODERATE"
              }
            } 
            primitive_transformation {

              # update values inside replace_config

              replace_config {
                new_value {
                  float_value = 652.23
                }
              }
            }
          }
        }
      }

      # update to remove field_transformations block using redact_config

      field_transformations {
        fields {
          name = "unconditionally-char-masked-field"
        }
        info_type_transformations {
          transformations {  
            info_types {
              name = "EMAIL_ADDRESS"
              version = "latest"
            }

            # adding the info_types

            info_types {
              name = "FIRST_NAME"
            }
            info_types {
              name = "LAST_NAME"
              version = "0.5"
            }  
            primitive_transformation {
              character_mask_config {
                masking_character = "x"
                number_to_mask = 8

                # update to delete old characters_to_ignore block and add new ones
                
                characters_to_ignore {
                  common_characters_to_ignore = "PUNCTUATION"
                }
                characters_to_ignore {
                  common_characters_to_ignore = "ALPHA_UPPER_CASE"
                }
                characters_to_ignore {
                  common_characters_to_ignore = "ALPHA_LOWER_CASE"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-replace-ffx-fpe-field"
        }
        info_type_transformations {
          transformations {  

            # updated the info_types

            info_types {
              name = "SSN33"
            }  
            primitive_transformation {
              crypto_replace_ffx_fpe_config {
                common_alphabet = "UPPER_CASE_ALPHA_NUMERIC"
                context {
                  name = "someTweak2"
                }
                crypto_key {
                  transient {
                    name = "beep"
                  }
                }
                surrogate_info_type {
                  name = "CUSTOM_INFO_TYPE"
                  version = "version-2"
                  sensitivity_score {
                    score = "SENSITIVITY_MODERATE"
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-fixed-size-bucketing-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "AGE"
            }
            primitive_transformation {

              # update values inside fixed_size_bucketing_config

              fixed_size_bucketing_config {
                lower_bound {
                  float_value = 23.5
                }
                upper_bound {
                  float_value = 71.75
                } 
                bucket_size = 20
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "CREATED_TIME"
            } 
            primitive_transformation {
              bucketing_config {
                buckets {
                  min {
                    string_value = "00:00:00"
                  }
                  max {
                    string_value = "11:59:59"
                  }
                  replacement_value {
                    string_value = "AM"
                  }
                }

                # Add new bucket

                buckets {
                  min {
                    string_value = "12:00:00"
                  }
                  max {
                    string_value = "13:59:59"
                  }
                  replacement_value {
                    string_value = "Lunchtime"
                  }
                }
                buckets {
                  min {
                    string_value = "14:00:00"
                  }
                  max {
                    string_value = "23:59:59"
                  }
                  replacement_value {
                    string_value = "PM"
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-time-part-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "DATE_OF_BIRTH"
            }
            primitive_transformation {
              time_part_config {
                part_to_extract = "MONTH"
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-hash-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "CREDIT_CARD_SECRET"
            } 
            primitive_transformation {
              crypto_hash_config {
                crypto_key {
                  transient {

                    # update the value

                    name = "beepy-beep-updated"
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-date-shift-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "EXTRACT_DATE"
            } 
            primitive_transformation {

              # update the value

              date_shift_config {
                upper_bound_days = 60
                lower_bound_days = -60
                context {
                  name = "unconditionally-date-shift-field"
                }
                crypto_key {
                  transient {
                    name = "beepy-beep-updated"
                  }
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-deterministic-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "CREDIT_CARD_SECRET1234"
            }
            primitive_transformation {
              crypto_deterministic_config {
                crypto_key {
                  transient {

                    # update the value

                    name = "beepy-beep-updated"
                  }
                }
                surrogate_info_type {

                  # update info type

                  name = "CREDIT_CARD_TRACK_NUMBER"
                  version = "version-2"
                  sensitivity_score {
                    score = "SENSITIVITY_MODERATE"
                  }
                }
                context {
                  name = "unconditionally-crypto-deterministic-field"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-dictionary-field"
        }
        info_type_transformations {
          transformations {    
            info_types {
              name = "RANDOM_FIELD"
            }
            primitive_transformation {
              replace_dictionary_config {
                word_list {
                  words = [

                  # update the list
                  
                    "foo",
                    "fizz",
                    "some",
                    "bar",
                  ]
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-with-info-type-config"
        }
        info_type_transformations {
          transformations {    
            info_types {
              
              # updated the value

              name = "LAST_NAME"
            }
            primitive_transformation {
              replace_with_info_type_config {}
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigString(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigBoolean(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigTimestamp(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigTimevalue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigDatevalue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigDayOfWeek(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigString(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-config"
        }

        info_type_transformations {
          transformations {   
            info_types {
              name = "PHONE_NUMBER"
            }
            info_types {
              name = "CREDIT_CARD_NUMBER"
            } 
            primitive_transformation {
              replace_config {
                new_value {
                  string_value = "someVal"
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigBoolean(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-config"
        }

        info_type_transformations {
          transformations {   
            info_types {
              name = "PHONE_NUMBER"
            }
            info_types {
              name = "CREDIT_CARD_NUMBER"
            } 
            primitive_transformation {
              replace_config {
                new_value {
                  boolean_value = true
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigTimestamp(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-config"
        }

        info_type_transformations {
          transformations {   
            info_types {
              name = "PHONE_NUMBER"
            }
            info_types {
              name = "CREDIT_CARD_NUMBER"
            } 
            primitive_transformation {
              replace_config {
                new_value {
                  timestamp_value = "2021-11-16T17:28:52Z"
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigTimevalue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-config"
        }

        info_type_transformations {
          transformations {   
            info_types {
              name = "PHONE_NUMBER"
            }
            info_types {
              name = "CREDIT_CARD_NUMBER"
            } 
            primitive_transformation {
              replace_config {
                new_value {
                  time_value {
                    hours = 22
                    minutes = 43
                    seconds = 54
                    nanos = 428947264
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigDatevalue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-config"
        }

        info_type_transformations {
          transformations {   
            info_types {
              name = "PHONE_NUMBER"
            }
            info_types {
              name = "CREDIT_CARD_NUMBER"
            } 
            primitive_transformation {
              replace_config {
                new_value {
                  date_value {
                    day = 24
                    month = 8
                    year = 2020
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_replaceConfigDayOfWeek(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-replace-config"
        }

        info_type_transformations {
          transformations {   
            info_types {
              name = "PHONE_NUMBER"
            }
            info_types {
              name = "CREDIT_CARD_NUMBER"
            } 
            primitive_transformation {
              replace_config {
                new_value {
                  day_of_week_value = "WEDNESDAY"
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoReplaceFfxFpeConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization": envvar.GetTestOrgFromEnv(t),
		"kms_key_name": acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoReplaceFfxFpeConfigTransient(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoReplaceFfxFpeConfigUnwrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoReplaceFfxFpeConfigKmswrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoReplaceFfxFpeConfigTransient(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-replace-ffx-fpe-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_replace_ffx_fpe_config {
                context {
                  name = "someTweak"
                }
                crypto_key {
                  transient {
                    name = "someRandomTerraformKey"
                  }
                }
                custom_alphabet = "ASE13RT76"
                surrogate_info_type {
                  name    = "CUSTOM_INFO_TYPE"
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoReplaceFfxFpeConfigUnwrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-replace-ffx-fpe-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_replace_ffx_fpe_config {
                context {
                  name = "someTweak2"
                }
                crypto_key {
                  unwrapped {
                    key = "0836c61118ac590243bdadb25f0bb08e"
                  }
                }
                common_alphabet = "HEXADECIMAL"
                surrogate_info_type {
                  name    = "CUSTOM_INFO_TYPE"
                  version = "version-1"
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoReplaceFfxFpeConfigKmswrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-replace-ffx-fpe-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_replace_ffx_fpe_config {
                context {
                  name = "someTweak3"
                }
                crypto_key {
                  kms_wrapped {
                    wrapped_key = "B64/WRAPPED/TOKENIZATION/KEY"
                    crypto_key_name = "%{kms_key_name}"
                  }
                }
                radix = 57
                surrogate_info_type {
                  name    = "CUSTOM_INFO_TYPE"
                  version = "version-2"
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigInteger(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigFloat(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigTimestamp(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigTimeValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigDateValue(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigDayOfTheWeek(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigInteger(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-config"
        }

        info_type_transformations {
          transformations {  
            info_types {
              name = "CREATED_TIME"
            }  
            primitive_transformation {
              bucketing_config {
                buckets {
                  min {
                    integer_value = 921
                  }
                  max {
                    integer_value = 3010
                  }
                  replacement_value {
                    integer_value = 1212
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigFloat(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-config"
        }

        info_type_transformations {
          transformations {  
            info_types {
              name = "CREATED_TIME"
            }  
            primitive_transformation {
              bucketing_config {
                buckets {
                  min {
                    float_value = 10.50
                  }
                  max {
                    float_value = 310.75
                  }
                  replacement_value {
                    float_value = 5.37
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigTimestamp(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-config"
        }

        info_type_transformations {
          transformations {  
            info_types {
              name = "CREATED_TIME"
            }  
            primitive_transformation {
              bucketing_config {
                buckets {
                  min {
                    timestamp_value = "2014-10-02T15:01:23Z"
                  }
                  max {
                    timestamp_value = "2015-06-29T18:46:39Z"
                  }
                  replacement_value {
                    timestamp_value = "2014-12-24T09:19:50Z"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigTimeValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-config"
        }

        info_type_transformations {
          transformations {  
            info_types {
              name = "CREATED_TIME"
            }  
            primitive_transformation {
              bucketing_config {
                buckets {
                  min {
                    time_value {
                      hours   = 09
                      minutes = 30
                      seconds = 45
                      nanos   = 123412
                    }
                  }
                  max {
                    time_value {
                      hours   = 15
                      minutes = 45
                      seconds = 00
                      nanos   = 523278
                    }
                  }
                  replacement_value {
                    time_value {
                      hours   = 23
                      minutes = 59
                      seconds = 59
                      nanos   = 999999
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigDateValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-config"
        }

        info_type_transformations {
          transformations {  
            info_types {
              name = "CREATED_TIME"
            }  
            primitive_transformation {
              bucketing_config{
                buckets {
                  min {
                    date_value {
                      year = 1969
                      month = 11
                      day = 23
                    }
                  }
                  max {
                    date_value {
                      year = 2010
                      month = 12
                      day = 31
                    }
                  }
                  replacement_value {
                    date_value {
                      year = 2011
                      month = 05
                      day = 19
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_bucketingConfigDayOfTheWeek(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "projects/%{project}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-bucketing-config"
        }

        info_type_transformations {
          transformations {  
            info_types {
              name = "CREATED_TIME"
            }  
            primitive_transformation {
              bucketing_config {
                buckets {
                  min {
                    day_of_week_value = "MONDAY"
                  }
                  max {
                    day_of_week_value = "THURSDAY"
                  }
                  replacement_value {
                    day_of_week_value = "FRIDAY"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoHashConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization": envvar.GetTestOrgFromEnv(t),
		"kms_key_name": acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoHashConfigTransient(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoHashConfigUnwrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoHashConfigKmswrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoHashConfigTransient(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-hash-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_hash_config {
                crypto_key {
                  transient {
                    name = "someRandomTerraformKey"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoHashConfigUnwrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-hash-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_hash_config {
                crypto_key {
                  unwrapped {
                    key = "0836c61118ac590243bdadb25f0bb08e"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoHashConfigKmswrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-hash-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_hash_config {
                crypto_key {
                  kms_wrapped {
                    wrapped_key = "B64/WRAPPED/TOKENIZATION/KEY"
                    crypto_key_name = "%{kms_key_name}"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_dateShiftConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization": envvar.GetTestOrgFromEnv(t),
		"kms_key_name": acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_dateShiftConfigTransient(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_dateShiftConfigUnwrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_dateShiftConfigKmswrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_dateShiftConfigTransient(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-date-shift-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              date_shift_config {
                upper_bound_days = 30
                lower_bound_days = -30
                context {
                  name = "some-context-field"
                }
                crypto_key {
                  transient {
                    name = "someRandomTerraformKey"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_dateShiftConfigUnwrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-date-shift-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              date_shift_config {
                upper_bound_days = 30
                lower_bound_days = -30
                context {
                  name = "some-context-field"
                }
                crypto_key {
                  unwrapped {
                    key = "0836c61118ac590243bdadb25f0bb08e"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_dateShiftConfigKmswrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-date-shift-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              date_shift_config {
                upper_bound_days = 30
                lower_bound_days = -30
                context {
                  name = "some-context-field"
                }
                crypto_key {
                  kms_wrapped {
                    wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                    crypto_key_name = "%{kms_key_name}"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoDeterministicConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization": envvar.GetTestOrgFromEnv(t),
		"kms_key_name": acctest.BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoDeterministicConfigTransient(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoDeterministicConfigUnwrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoDeterministicConfigKmswrapped(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoDeterministicConfigTransient(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-deterministic-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_deterministic_config {
                surrogate_info_type {
                  name = "SECRET_NUMBER"
                }
                crypto_key {
                  transient {
                    name = "someRandomTerraformKey"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoDeterministicConfigUnwrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-deterministic-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_deterministic_config {
                surrogate_info_type {
                  name = "SECRET_NUMBER"
                  version = "1.0"
                }
                context {
                  name = "some-context-field"
                }
                crypto_key {
                  unwrapped {
                    key = "0836c61118ac590243bdadb25f0bb08e"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformations_with_infoTypeTransformations_cryptoDeterministicConfigKmswrapped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    record_transformations {
      record_suppressions {
        condition {
          expressions {
            logical_operator = "AND"
            conditions {
              conditions {
                field {
                  name = "field1"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "FOO-BAR"
                }
              }
              conditions {
                field {
                  name = "field2"
                }
                operator = "EQUAL_TO"
                value {
                  string_value = "foobar"
                }
              }
            }
          }
        }
      }
      field_transformations {
        fields {
          name = "unconditionally-crypto-deterministic-config"
        }

        info_type_transformations {
          transformations {
            primitive_transformation {
              crypto_deterministic_config {
                surrogate_info_type {
                  name = "SECRET_NUMBER"
                  version = "2.0"
                }
                context {
                  name = "updated-context-field"
                }
                crypto_key {
                  kms_wrapped {
                    wrapped_key     = "B64/WRAPPED/TOKENIZATION/KEY"
                    crypto_key_name = "%{kms_key_name}"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`, context)
}
