package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplateUpdate(t *testing.T) {
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
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplateStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplateUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_deidentify_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplateStart(context map[string]interface{}) string {
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
		}
	}
}
`, context)
}

func testAccDataLossPreventionDeidentifyTemplate_dlpDeidentifyTemplateUpdate(context map[string]interface{}) string {
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
		}
	}
}
`, context)
}
