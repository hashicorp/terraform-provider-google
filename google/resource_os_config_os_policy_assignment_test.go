package google

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	osconfig "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/osconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccOsConfigOsPolicyAssignment_basicOsPolicyAssignment(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"zone":          getTestZoneFromEnv(),
		"random_suffix": randString(t, 10),
		"org_id":        getTestOrgFromEnv(t),
		"billing_act":   getTestBillingAccountFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOsConfigOsPolicyAssignmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOsConfigOsPolicyAssignment_PercentOsPolicyAssignment(context),
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

func testAccOsConfigOsPolicyAssignment_PercentOsPolicyAssignment(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_act}"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "osconfig" {
  project = google_project.project.project_id
  service = "osconfig.googleapis.com"
  depends_on = [google_project_service.compute]
}

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

  location = "%{zone}"
  name     = "tf-test-assignment%{random_suffix}"

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

    min_wait_duration = "3s"
  }

  description = "A test os policy assignment"
  project     = google_project.project.project_id
  depends_on = [google_project_service.compute, google_project_service.osconfig]
}


`, context)
}

func testAccCheckOsConfigOsPolicyAssignmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_os_config_os_policy_assignment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &osconfig.OSPolicyAssignment{
				Location:           dcl.String(rs.Primary.Attributes["location"]),
				Name:               dcl.String(rs.Primary.Attributes["name"]),
				Description:        dcl.String(rs.Primary.Attributes["description"]),
				Project:            dcl.StringOrNil(rs.Primary.Attributes["project"]),
				SkipAwaitRollout:   dcl.Bool(rs.Primary.Attributes["skip_await_rollout"] == "true"),
				Baseline:           dcl.Bool(rs.Primary.Attributes["baseline"] == "true"),
				Deleted:            dcl.Bool(rs.Primary.Attributes["deleted"] == "true"),
				Etag:               dcl.StringOrNil(rs.Primary.Attributes["etag"]),
				Reconciling:        dcl.Bool(rs.Primary.Attributes["reconciling"] == "true"),
				RevisionCreateTime: dcl.StringOrNil(rs.Primary.Attributes["revision_create_time"]),
				RevisionId:         dcl.StringOrNil(rs.Primary.Attributes["revision_id"]),
				RolloutState:       osconfig.OSPolicyAssignmentRolloutStateEnumRef(rs.Primary.Attributes["rollout_state"]),
				Uid:                dcl.StringOrNil(rs.Primary.Attributes["uid"]),
			}

			client := NewDCLOsConfigClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetOSPolicyAssignment(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_os_config_os_policy_assignment still exists %v", obj)
			}
		}
		return nil
	}
}
