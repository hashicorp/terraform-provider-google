package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataLossPreventionInspectTemplate_dlpInspectTemplateUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataLossPreventionInspectTemplateDestroyProducer(t),
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
	return Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
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
	return Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Updated"
	display_name = "Different"

	inspect_config {
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
