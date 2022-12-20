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
                crypto_key_name = google_kms_crypto_key.my_key.id
              }
            }
            surrogate_info_type {
              name = "abc"
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
                crypto_key_name = google_kms_crypto_key.my_key.id
              }
            }
            radix = 10
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPE"
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
            }
          }
        }
      }
    }
  }
}

resource "google_kms_crypto_key" "my_key" {
  name     = "tf-test-example-k%{random_suffix}"
  key_ring = google_kms_key_ring.key_ring.id
}

resource "google_kms_key_ring" "key_ring" {
  name     = "tf-test-example-keyr%{random_suffix}"
  location = "global"
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
                crypto_key_name = google_kms_crypto_key.my_key.id
              }
            }
            surrogate_info_type {
              name = "abcd"
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
                crypto_key_name = google_kms_crypto_key.my_key.id
              }
            }
            radix = 10
            surrogate_info_type {
              name = "CUSTOM_INFO_TYPEF"
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
            }
          }
        }
      }
    }
  }
}

resource "google_kms_crypto_key" "my_key" {
  name     = "tf-test-example-k%{random_suffix}"
  key_ring = google_kms_key_ring.key_ring.id
}

resource "google_kms_key_ring" "key_ring" {
  name     = "tf-test-example-keyr%{random_suffix}"
  location = "global"
}
`, context)
}

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplate_recordTransformationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":  getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
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
    }
  }
}
`, context)
}
