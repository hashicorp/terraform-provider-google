// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXPage_update(t *testing.T) {
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
				Config: testAccDialogflowCXPage_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_page.my_page",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXPage_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_page.my_page",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXPage_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_dialogflow_cx_agent" "agent_page" {
    display_name             = "tf-test-%{random_suffix}"
    location                 = "global"
    default_language_code    = "en"
    supported_language_codes = ["fr", "de", "es"]
    time_zone                = "America/New_York"
    description              = "Description 1."
    avatar_uri               = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
  }

  resource "google_dialogflow_cx_page" "my_page" {
    parent       = google_dialogflow_cx_agent.agent_page.start_flow
    display_name = "MyPage"
  }
`, context)
}

func testAccDialogflowCXPage_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_dialogflow_cx_agent" "agent_page" {
    display_name               = "tf-test-%{random_suffix}update"
    location                   = "global"
    default_language_code      = "en"
    supported_language_codes   = ["no"]
    time_zone                  = "Europe/London"
    description                = "Description 2!"
    avatar_uri                 = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo-2.png"
    enable_stackdriver_logging = true
    enable_spell_correction    = true
    speech_to_text_settings {
      enable_speech_adaptation = true
    }
  }

  resource "google_dialogflow_cx_page" "my_page" {
    parent       = google_dialogflow_cx_agent.agent_page.start_flow
    display_name = "MyPage"

    entry_fulfillment {
      messages {
        channel = "some-channel"
        text {
          text = ["Welcome to page"]
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

    event_handlers {
      event = "some-event"
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

    form {
      parameters {
        display_name = "param1"
        entity_type  = "projects/-/locations/-/agents/-/entityTypes/sys.date"
        default_value = jsonencode("2000-01-01")
        fill_behavior {
          initial_prompt_fulfillment {
            messages {
              channel = "some-channel"
              text {
                text = ["Please provide param1"]
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
          reprompt_event_handlers {
            event = "sys.no-match-1"
            trigger_fulfillment {
              return_partial_responses = true
              webhook = google_dialogflow_cx_webhook.my_webhook.id
              tag = "some-tag"

              messages {
                channel = "some-channel"
                text {
                  text = ["Please provide param1"]
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
          reprompt_event_handlers {
            event = "sys.no-match-2"
            target_flow = google_dialogflow_cx_agent.agent_page.start_flow
          }
          reprompt_event_handlers {
            event = "sys.no-match-3"
            target_page = google_dialogflow_cx_page.my_page2.id
          }
        }
        required = "true"
        redact   = "true"
        advanced_settings {
          dtmf_settings {
            enabled      = true
            max_digits   = 1
            finish_digit = "#"
          }
        }
      }
    }

    transition_routes {
      condition = "$page.params.status = 'FINAL'"
      trigger_fulfillment {
        messages {
          channel = "some-channel"
          text {
            text = ["information completed, navigating to page 2"]
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
      target_page = google_dialogflow_cx_page.my_page2.id
    }

    advanced_settings {
      dtmf_settings {
        enabled      = true
        max_digits   = 1
        finish_digit = "#"
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
      data_store_connections {
        data_store_type = "UNSTRUCTURED"
        data_store = "projects/${data.google_project.project.number}/locations/${google_dialogflow_cx_agent.agent_page.location}/collections/default_collection/dataStores/datastore-page-update"
        document_processing_mode = "DOCUMENTS"
      }
      target_page = google_dialogflow_cx_page.my_page2.id
    }
  }

  resource "google_dialogflow_cx_page" "my_page2" {
    parent       = google_dialogflow_cx_agent.agent_page.start_flow
    display_name = "MyPage2"
  }

  resource "google_discovery_engine_data_store" "my_datastore" {
    location          = "global"
    data_store_id     = "datastore-page-update"
    display_name      = "datastore-page-update"
    industry_vertical = "GENERIC"
    content_config    = "NO_CONTENT"
  }

  resource "google_dialogflow_cx_webhook" "my_webhook" {
    parent       = google_dialogflow_cx_agent.agent_page.id
    display_name = "MyWebhook"
    generic_web_service {
      uri = "https://example.com"
    }
  }

  data "google_project" "project" {
  }
`, context)
}
