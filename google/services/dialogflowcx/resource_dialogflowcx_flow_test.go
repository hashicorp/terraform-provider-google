// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXFlow_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXFlow_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_flow.my_flow",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXFlow_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_flow.my_flow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXFlow_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_dialogflow_cx_agent" "agent_entity" {
    display_name             = "tf-test-%{random_suffix}"
    location                 = "global"
    default_language_code    = "en"
    supported_language_codes = ["fr", "de", "es"]
    time_zone                = "America/New_York"
    description              = "Description 1."
    avatar_uri               = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
  }

  resource "google_dialogflow_cx_flow" "my_flow" {
    parent       = google_dialogflow_cx_agent.agent_entity.id
    display_name = "MyFlow"

    nlu_settings {
      classification_threshold = 0.3
      model_type               = "MODEL_TYPE_STANDARD"
    }
  }
`, context)
}

func testAccDialogflowCXFlow_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_dialogflow_cx_agent" "agent_entity" {
    display_name               = "tf-test-dialogflowcx-agent%{random_suffix}update"
    location                   = "global"
    default_language_code      = "en"
    supported_language_codes   = ["fr", "de", "es"]
    time_zone                  = "America/New_York"
    description                = "Example description."
    avatar_uri                 = "https://cloud.google.com/_static/images/cloud/icons/favicons/onecloud/super_cloud.png"
    enable_stackdriver_logging = true
    enable_spell_correction    = true
    speech_to_text_settings {
      enable_speech_adaptation = true
    }
  }

  resource "google_dialogflow_cx_flow" "my_flow" {
    parent       = google_dialogflow_cx_agent.agent_entity.id
    display_name = "MyFlow"
    description  = "Test Flow"
  
    nlu_settings {
      classification_threshold = 0.3
      model_type               = "MODEL_TYPE_STANDARD"
    }
  
    event_handlers {
      event = "custom-event"
      trigger_fulfillment {
        return_partial_responses = false
        messages {
          text {
            text = ["I didn't get that. Can you say it again?"]
          }
        }
      }
    }

    event_handlers {
      event = "sys.no-match-default"
      trigger_fulfillment {
        return_partial_responses = false
        messages {
          text {
            text = ["Sorry, could you say that again?"]
          }
        }
      }
    }

    event_handlers {
      event = "sys.no-input-default"
      trigger_fulfillment {
        return_partial_responses = false
        messages {
          text {
            text = ["One more time?"]
          }
        }
      }
    }

    event_handlers {
      event = "another-event"
      trigger_fulfillment {
        return_partial_responses = true
        messages {
          channel = "some-channel"
          text {
            text = ["Some text"]
          }
        }
        messages {
          payload = <<EOF
            {"some-key": "some-value", "other-key": ["other-value"]}
          EOF
        }
        messages {
          conversation_success {
            metadata = <<EOF
              {"some-metadata-key": "some-value", "other-metadata-key": 1234}
            EOF
          }
        }
        messages {
          output_audio_text {
            text = "some output text"
          }
        }
        messages {
          output_audio_text {
            ssml = <<EOF
              <speak>Some example <say-as interpret-as="characters">SSML XML</say-as></speak>
            EOF
          }
        }
        messages {
          live_agent_handoff {
            metadata = <<EOF
              {"some-metadata-key": "some-value", "other-metadata-key": 1234}
            EOF
          }
        }
        messages {
          play_audio {
            audio_uri = "http://example.com/some-audio-file.mp3"
          }
        }
        messages {
          telephony_transfer_call {
            phone_number = "1-234-567-8901"
          }
        }

        set_parameter_actions {
          parameter = "some-param"
          value     = "123.45"
        }
        set_parameter_actions {
          parameter = "another-param"
          value     = jsonencode("abc")
        }
        set_parameter_actions {
          parameter = "other-param"
          value     = jsonencode(["foo"])
        }

        conditional_cases {
          cases = jsonencode([
            {
              condition = "$sys.func.RAND() < 0.5",
              caseContent = [
                {
                  message = { text = { text = ["First case"] } }
                },
                {
                  additionalCases = {
                    cases = [
                      {
                        condition = "$sys.func.RAND() < 0.2"
                        caseContent = [
                          {
                            message = { text = { text = ["Nested case"] } }
                          }
                        ]
                      }
                    ]
                  }
                }
              ]
            },
            {
              caseContent = [
                {
                  message = { text = { text = ["Final case"] } }
                }
              ]
            },
          ])
        }
      }
    }

    transition_routes {
      condition = "true"
      trigger_fulfillment {
        return_partial_responses = true
        messages {
          channel = "some-channel"
          text {
            text = ["Some text"]
          }
        }
        messages {
          payload = <<EOF
            {"some-key": "some-value", "other-key": ["other-value"]}
          EOF
        }
        messages {
          conversation_success {
            metadata = <<EOF
              {"some-metadata-key": "some-value", "other-metadata-key": 1234}
            EOF
          }
        }
        messages {
          output_audio_text {
            text = "some output text"
          }
        }
        messages {
          output_audio_text {
            ssml = <<EOF
              <speak>Some example <say-as interpret-as="characters">SSML XML</say-as></speak>
            EOF
          }
        }
        messages {
          live_agent_handoff {
            metadata = <<EOF
              {"some-metadata-key": "some-value", "other-metadata-key": 1234}
            EOF
          }
        }
        messages {
          play_audio {
            audio_uri = "http://example.com/some-audio-file.mp3"
          }
        }
        messages {
          telephony_transfer_call {
            phone_number = "1-234-567-8901"
          }
        }

        set_parameter_actions {
          parameter = "some-param"
          value     = "123.45"
        }
        set_parameter_actions {
          parameter = "another-param"
          value     = jsonencode("abc")
        }
        set_parameter_actions {
          parameter = "other-param"
          value     = jsonencode(["foo"])
        }

        conditional_cases {
          cases = jsonencode([
            {
              condition = "$sys.func.RAND() < 0.5",
              caseContent = [
                {
                  message = { text = { text = ["First case"] } }
                },
                {
                  additionalCases = {
                    cases = [
                      {
                        condition = "$sys.func.RAND() < 0.2"
                        caseContent = [
                          {
                            message = { text = { text = ["Nested case"] } }
                          }
                        ]
                      }
                    ]
                  }
                }
              ]
            },
            {
              caseContent = [
                {
                  message = { text = { text = ["Final case"] } }
                }
              ]
            },
          ])
        }
      }
      target_flow = google_dialogflow_cx_agent.agent_entity.start_flow
    }
  }
`, context)
}
