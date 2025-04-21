// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package osconfigv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccOSConfigV2PolicyOrchestratorForFolder_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"zone":          envvar.GetTestZoneFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOSConfigV2PolicyOrchestratorForFolder_basic(context),
			},
			{
				ResourceName:            "google_os_config_v2_policy_orchestrator_for_folder.policy_orchestrator_for_folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder_id", "labels", "policy_orchestrator_id", "terraform_labels"},
			},
			{
				Config: testAccOSConfigV2PolicyOrchestratorForFolder_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_os_config_v2_policy_orchestrator_for_folder.policy_orchestrator_for_folder", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_os_config_v2_policy_orchestrator_for_folder.policy_orchestrator_for_folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder_id", "labels", "policy_orchestrator_id", "terraform_labels"},
			},
		},
	})
}

func testAccOSConfigV2PolicyOrchestratorForFolder_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "my_folder" {
    display_name        = "tf-test-po-folder%{random_suffix}"
    parent              = "organizations/%{org_id}"
    deletion_protection = false
}

resource "google_folder_service_identity" "osconfig_sa" {
  folder  = google_folder.my_folder.folder_id
  service = "osconfig.googleapis.com"
}

resource "google_folder_service_identity" "ripple_sa" {
  folder  = google_folder.my_folder.folder_id
  service = "progressiverollout.googleapis.com"
}

resource "time_sleep" "wait_30_sec" {
    depends_on = [
        google_folder_service_identity.osconfig_sa,
        google_folder_service_identity.ripple_sa,
    ]
    create_duration = "30s"
}

resource "google_folder_iam_member" "iam_osconfig_service_agent" {
    depends_on = [time_sleep.wait_30_sec]
    folder = google_folder.my_folder.folder_id
    role   = "roles/osconfig.serviceAgent"
    member = google_folder_service_identity.osconfig_sa.member
}

resource "google_folder_iam_member" "iam_osconfig_rollout_service_agent" {
    depends_on = [google_folder_iam_member.iam_osconfig_service_agent]
    folder     = google_folder.my_folder.folder_id
    role       = "roles/osconfig.rolloutServiceAgent"
    member     = "serviceAccount:service-folder-${google_folder.my_folder.folder_id}@gcp-sa-osconfig-rollout.iam.gserviceaccount.com"
}

resource "google_folder_iam_member" "iam_progressiverollout_service_agent" {
    depends_on = [google_folder_iam_member.iam_osconfig_rollout_service_agent]
    folder = google_folder.my_folder.folder_id
    role   = "roles/progressiverollout.serviceAgent"
    member = google_folder_service_identity.ripple_sa.member
}

resource "time_sleep" "wait_3_min" {
    depends_on = [google_folder_iam_member.iam_progressiverollout_service_agent]
    create_duration = "180s"
}


resource "google_os_config_v2_policy_orchestrator_for_folder" "policy_orchestrator_for_folder" {
    depends_on = [time_sleep.wait_3_min]

    policy_orchestrator_id = "tf-test-po-folder%{random_suffix}"
    folder_id = google_folder.my_folder.folder_id
    
    state = "ACTIVE"
    action = "UPSERT"
    
    orchestrated_resource {
        id = "tf-test-test-orchestrated-resource-folder%{random_suffix}"
        os_policy_assignment_v1_payload {
            description = "ospa for create"
            name = "ospa-1"
            os_policies {
                id = "tf-test-test-os-policy-folder%{random_suffix}"
                description = "policy for create"
                allow_no_resource_group_match = true
                mode = "VALIDATION"
                resource_groups {
                    inventory_filters {
                        os_short_name = "windows-10"
                        os_version    = "10.0.19044"
                    }
                    resources {
                        id = "resource-tf"
                        file {
                            content = "file-content-tf"
                            path = "file-path-tf-1"
                            state = "PRESENT"
                        }
                    }
                }
            }
            instance_filter {
                all = false
                inventories {
                    os_short_name = "windows-10"
                    os_version = "10.0.19044"
                }
                exclusion_labels {
                    labels = {
                        label1 = "test-exclusion-label-1"
                    }
                }
                inclusion_labels {
                    labels = {
                        label1 = "test-inclusion-label-1"
                    }
                }
            }
            rollout {
                disruption_budget {
                    percent = 100
                }
                min_wait_duration = "60s"
            }
        }
    }
    labels = {
        state = "active"
    }
    orchestration_scope {
        selectors {
            location_selector {
                included_locations = ["%{zone}"]
            }
        }
        selectors {
            resource_hierarchy_selector {
                included_folders  = ["folders/${google_folder.my_folder.folder_id}"]
            }
        }
    }
}
`, context)
}

func testAccOSConfigV2PolicyOrchestratorForFolder_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "my_folder" {
    display_name        = "tf-test-po-folder%{random_suffix}"
    parent              = "organizations/%{org_id}"
    deletion_protection = false
}

resource "google_os_config_v2_policy_orchestrator_for_folder" "policy_orchestrator_for_folder" {
    policy_orchestrator_id = "tf-test-po-folder%{random_suffix}"
    folder_id = google_folder.my_folder.folder_id

    state = "STOPPED"
    action = "DELETE"
    description = "Updated description"

    orchestrated_resource {
        id = "tf-test-updated-orchestrated-resource-folder%{random_suffix}"
        os_policy_assignment_v1_payload {
            description = "ospa for update"
            name = "ospa-2"
            os_policies {
                id = "tf-test-test-os-policy-folder%{random_suffix}"
                description = "policy for update"
                allow_no_resource_group_match = false
                mode = "ENFORCEMENT"
                resource_groups {
                    inventory_filters {
                        os_short_name = "debian"
                        os_version    = "11"
                    }
                    resources {
                        id = "resource-tf"
                        exec {
                            enforce {
                                args             = ["--arg1", "--arg2"]
                                interpreter      = "SHELL"
                                output_file_path = "/tmp/enforce_output.txt"
                                file {
                                    allow_insecure = false
                                    gcs {
                                        bucket     = "my-bucket"
                                        generation = 1
                                        object     = "scripts/enforce.sh"
                                    }
                                }
                            }
                            validate {
                                args             = ["--validate"]
                                interpreter      = "POWERSHELL"
                                output_file_path = "C:\\validate_out.txt"
                                file {
                                      allow_insecure = false
                                      gcs {
                                          bucket     = "my-bucket"
                                          generation = 2
                                          object     = "scripts/validate.ps1"
                                      }
                                }
                            }
                        }
                    }
                }
            }
            instance_filter {
                all = false
                inventories {
                    os_short_name = "debian"
                    os_version = "11"
                }
                exclusion_labels {
                    labels = {
                        label1 = "test-exclusion-label-2"
                    }
                }
                inclusion_labels {
                    labels = {
                        label1 = "test-inclusion-label-2"
                    }
                }
            }
            rollout {
                disruption_budget {
                    fixed = 1
                }
                min_wait_duration = "120s"
            }
        }
    }
    labels = {}
    orchestration_scope {
        selectors {
            location_selector {
                included_locations = ["us-central2-b"]
            }
        }
        selectors {
            resource_hierarchy_selector {
                included_folders  = ["folders/${google_folder.my_folder.folder_id}"]
            }
        }
    }
}
`, context)
}
