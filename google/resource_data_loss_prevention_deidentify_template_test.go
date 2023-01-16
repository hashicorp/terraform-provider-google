package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
		"kms_key_name":  BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
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
	return Nprintf(`
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
    }
  }
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_infoTypeTransformationsUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_deidentify_template" "basic" {
  parent = "organizations/%{organization}"
  description = "Description"
  display_name = "Displayname"

  deidentify_config {
    info_type_transformations {
      transformations {
        info_types {
          name = "CREDIT_CARD_NUMBER"
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
    }
  }
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
		"kms_key_name":  BootstrapKMSKey(t).CryptoKey.Name, // global KMS key
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataLossPreventionDeidentifyTemplateDestroyProducer(t),
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
	return Nprintf(`
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
	return Nprintf(`
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
