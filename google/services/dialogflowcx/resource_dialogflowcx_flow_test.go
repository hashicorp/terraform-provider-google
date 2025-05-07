// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				ResourceName:            "google_dialogflow_cx_flow.my_flow",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_settings.0.logging_settings"},
			},
			{
				Config: testAccDialogflowCXFlow_full(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_flow.my_flow",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_settings.0.logging_settings"},
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

  resource "google_storage_bucket" "bucket" {
    name                        = "tf-test-dialogflowcx-bucket%{random_suffix}"
    location                    = "US"
    uniform_bucket_level_access = true
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

    advanced_settings {
      audio_export_gcs_destination {
        uri = "${google_storage_bucket.bucket.url}/prefix-"
      }
      speech_settings {
        endpointer_sensitivity        = 30
        no_speech_timeout             = "3.500s"
        use_timeout_based_endpointing = true
        models = {
          name : "wrench"
          mass : "1.3kg"
          count : "3"
        }
      }
      dtmf_settings {
        enabled      = true
        max_digits   = 1
        finish_digit = "#"
      }
      logging_settings {
        enable_stackdriver_logging     = true
        enable_interaction_logging     = true
        enable_consent_based_redaction = true
      }
    }

    knowledge_connector_settings {
      enabled = true
      trigger_fulfillment {
        messages {
          channel = "some-channel"
          output_audio_text {
            text = "some output text"
          }
        }
      }
      target_flow = google_dialogflow_cx_agent.agent_entity.start_flow
    }
  }
`, context)
}

func TestAccDialogflowCXFlow_defaultStartFlow(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Note: this isn't actually a "create" test; it creates a resource in the TF state, but is actually importing the default object GCP has created, then updating it.
				Config: testAccDialogflowCXFlow_defaultStartFlow_create(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dialogflow_cx_flow.default_start_flow", "name", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttrPair(
						"google_dialogflow_cx_flow.default_start_flow", "id",
						"google_dialogflow_cx_agent.agent", "start_flow",
					),
				),
			},
			{
				ResourceName:      "google_dialogflow_cx_flow.default_start_flow",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// This is testing updating the default object without having to create it in the TF state first.
				Config: testAccDialogflowCXFlow_defaultStartFlow_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dialogflow_cx_flow.default_start_flow", "name", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttrPair(
						"google_dialogflow_cx_flow.default_start_flow", "id",
						"google_dialogflow_cx_agent.agent", "start_flow",
					),
				),
			},
			{
				ResourceName:      "google_dialogflow_cx_flow.default_start_flow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXFlow_defaultStartFlow_create(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
  display_name          = "tf-test-dialogflowcx-agent%{random_suffix}"
  location              = "global"
  default_language_code = "en"
  time_zone             = "America/New_York"
}

resource "google_dialogflow_cx_intent" "default_welcome_intent" {
  parent                    = google_dialogflow_cx_agent.agent.id
  is_default_welcome_intent = true
  display_name              = "Default Welcome Intent"
  priority                  = 1
  training_phrases {
    parts {
      text = "Hello"
    }
    repeat_count = 1
  }
}


resource "google_dialogflow_cx_flow" "default_start_flow" {
  parent                = google_dialogflow_cx_agent.agent.id
  is_default_start_flow = true
  display_name          = "Default Start Flow"
  description           = "A start flow created along with the agent"

  nlu_settings {
    classification_threshold = 0.3
    model_type               = "MODEL_TYPE_STANDARD"
  }

  transition_routes {
    intent = google_dialogflow_cx_intent.default_welcome_intent.id
    trigger_fulfillment {
      messages {
        text {
          text = ["Response to default welcome intent."]
        }
      }
    }
  }

  event_handlers {
    event = "custom-event"
    trigger_fulfillment {
      messages {
        text {
          text = ["Handle a custom event!"]
        }
      }
    }
  }

  event_handlers {
    event = "sys.no-match-default"
    trigger_fulfillment {
      messages {
        text {
          text = ["This is the flow no-match response."]
        }
      }
    }
  }

  event_handlers {
    event = "sys.no-input-default"
    trigger_fulfillment {
      messages {
        text {
          text = ["This is the flow no-input response."]
        }
      }
    }
  }
}
`, context)
}

func testAccDialogflowCXFlow_defaultStartFlow_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
  display_name          = "tf-test-dialogflowcx-agent%{random_suffix}"
  location              = "global"
  default_language_code = "en"
  time_zone             = "America/New_York"
}

resource "google_dialogflow_cx_intent" "default_welcome_intent" {
  parent                    = google_dialogflow_cx_agent.agent.id
  is_default_welcome_intent = true
  display_name              = "Default Welcome Intent"
  priority                  = 1
  training_phrases {
    parts {
      text = "Hello"
    }
    repeat_count = 1
  }
}

resource "google_dialogflow_cx_page" "my_page" {
  parent       = google_dialogflow_cx_agent.agent.start_flow
  display_name = "MyPage"
}

resource "google_dialogflow_cx_flow" "default_start_flow" {
  parent                = google_dialogflow_cx_agent.agent.id
  is_default_start_flow = true
  display_name          = "Default Start Flow"
  description           = "A start flow created along with the agent"

  nlu_settings {
    classification_threshold = 0.5
    model_type               = "MODEL_TYPE_STANDARD"
  }

  transition_routes {
    intent = google_dialogflow_cx_intent.default_welcome_intent.id
    trigger_fulfillment {
      messages {
        text {
          text = ["We can update the default welcome intent response!"]
        }
      }
    }
  }

  // delete the custom-event handler to show we can

  event_handlers {
    event = "sys.no-match-default"
    trigger_fulfillment {
      messages {
        text {
          text = ["We an also update the no-match response!"]
        }
      }
    }
  }

  event_handlers {
    event = "sys.no-input-default"
    trigger_fulfillment {
      messages {
        text {
          text = ["The no-input response has been updated too!"]
        }
      }
    }
  }

  knowledge_connector_settings {
    enabled = false
    trigger_fulfillment {
      messages {
        output_audio_text {
          text = "We can update the knowledge_connector_settings in this flow!"
        }
      }
    }
    target_page = google_dialogflow_cx_page.my_page.id
  }
}
`, context)
}
