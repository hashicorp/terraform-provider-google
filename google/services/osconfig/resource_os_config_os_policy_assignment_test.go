// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package osconfig_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccOSConfigOSPolicyAssignment_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOSConfigOSPolicyAssignmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOSConfigOSPolicyAssignment_basic(context),
			},
			{
				ResourceName:            "google_os_config_os_policy_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rollout.0.min_wait_duration"},
			},
			{
				Config: testAccOSConfigOSPolicyAssignment_update(context),
			},
			{
				ResourceName:            "google_os_config_os_policy_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rollout.0.min_wait_duration"},
			},
		},
	})
}

func testAccOSConfigOSPolicyAssignment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = false
    exclusion_labels {
      labels = {
        label-two = "value-two"
      }
    }
    inclusion_labels {
      labels = {
        label-one = "value-one"
      }
    }
    inventories {
      os_short_name = "centos"
      os_version    = "8.*"
    }
  }

  location = "us-central1-a"
  name     = "tf-test-policy-assignment%{random_suffix}"

  os_policies {
    id   = "policy"
    mode = "VALIDATION"

    resource_groups {
      resources {
        id = "apt-to-yum"

        repository {
          apt {
            archive_type = "DEB"
            components   = ["doc"]
            distribution = "debian"
            uri          = "https://atl.mirrors.clouvider.net/debian"
            gpg_key      = ".gnupg/pubring.kbx"
          }
        }
      }
      inventory_filters {
        os_short_name = "centos"
        os_version    = "8.*"
      }

      resources {
        id = "exec1"
        exec {
          validate {
            interpreter = "SHELL"
            args        = ["arg1"]
            file {
              local_path = "$HOME/script.sh"
            }
            output_file_path = "$HOME/out"
          }
          enforce {
            interpreter = "SHELL"
            args        = ["arg1"]
            file {
              allow_insecure = true
              remote {
                uri             = "https://www.example.com/script.sh"
                sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
              }
            }
            output_file_path = "$HOME/out"
          }
        }
      }
    }
    allow_no_resource_group_match = false
    description                   = "A test os policy"
  }

  rollout {
    disruption_budget {
      percent = 100
    }

    min_wait_duration = "3.2s"
  }

  description = "A test os policy assignment"
}
`, context)
}

func testAccOSConfigOSPolicyAssignment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = false
    inventories {
      os_short_name = "centos"
      os_version    = "9.*"
    }
  }

  location = "us-central1-a"
  name     = "tf-test-policy-assignment%{random_suffix}"

  os_policies {
    id   = "policy"
    mode = "ENFORCEMENT"

    resource_groups {
      resources {
        id = "apt-to-yum"

        repository {
          yum {
            id           = "new-yum"
            display_name = "new-yum"
            base_url     = "http://mirrors.rcs.alaska.edu/centos/"
            gpg_keys     = ["RPM-GPG-KEY-CentOS-Debug-7"]
          }
        }
      }
      inventory_filters {
        os_short_name = "centos"
        os_version    = "8.*"
      }

      resources {
        id = "new-exec1"
        exec {
          validate {
            interpreter = "POWERSHELL"
            args        = ["arg2"]
            file {
              local_path = "$HOME/script.bat"
            }
            output_file_path = "$HOME/out"
          }
          enforce {
            interpreter = "POWERSHELL"
            args        = ["arg2"]
            file {
              allow_insecure = false
              remote {
                uri             = "https://www.example.com/script.bat"
                sha256_checksum = "9f8e5818ccb47024d01000db713c0a333679b64678ff5fe2d9bea0a23014dd54"
              }
            }
            output_file_path = "$HOME/out"
          }
        }
      }
    }
    allow_no_resource_group_match = true
    description                   = "An updated test os policy"
  }

  rollout {
    disruption_budget {
      percent = 90
    }

    min_wait_duration = "3.1s"
  }

  description = "An updated test os policy assignment"
}
`, context)
}

func testAccCheckOSConfigOSPolicyAssignmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_os_config_os_policy_assignment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{OSConfigBasePath}}projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("OSConfigOSPolicyAssignment still exists at %s", url)
			}
		}

		return nil
	}
}
