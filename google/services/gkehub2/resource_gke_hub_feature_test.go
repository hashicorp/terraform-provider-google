// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccGKEHubFeature_gkehubFeatureFleetObservability(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_gkehubFeatureFleetObservability(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate1(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate2(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeature_gkehubFeatureFleetObservability(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "time_sleep" "wait_for_gkehub_enablement" {
  create_duration = "150s"
  depends_on = [google_project_service.gkehub]
}

resource "google_gke_hub_feature" "feature" {
  name = "fleetobservability"
  location = "global"
  project = google_project.project.project_id
  spec {
    fleetobservability {
      logging_config {
        default_config {
    mode = "MOVE"
        }
        fleet_scope_logs_config {
          mode = "COPY"
        }
      }
    }
  }
  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate1(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "time_sleep" "wait_for_gkehub_enablement" {
  create_duration = "150s"
  depends_on = [google_project_service.gkehub]
}

resource "google_gke_hub_feature" "feature" {
  name = "fleetobservability"
  location = "global"
  project = google_project.project.project_id
  spec {
    fleetobservability {
      logging_config {
        default_config {
    mode = "MOVE"
        }
      }
    }
  }
  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate2(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "time_sleep" "wait_for_gkehub_enablement" {
  create_duration = "150s"
  depends_on = [google_project_service.gkehub]
}

resource "google_gke_hub_feature" "feature" {
  name = "fleetobservability"
  location = "global"
  project = google_project.project.project_id
  spec {
    fleetobservability {
      logging_config {
        fleet_scope_logs_config {
          mode = "COPY"
        }
      }
    }
  }
  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func TestAccGKEHubFeature_gkehubFeatureMciUpdate(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_gkehubFeatureMciUpdateStart(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time"},
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureMciChangeMembership(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEHubFeature_gkehubFeatureMciUpdateStart(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`

resource "google_container_cluster" "primary" {
  name               = "tf-test%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_container_cluster" "secondary" {
  name               = "tf-test2%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "tf-test%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_membership" "membership_second" {
  membership_id = "tf-test2%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.secondary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_feature" "feature" {
  name = "multiclusteringress"
  location = "global"
  spec {
    multiclusteringress {
      config_membership = google_gke_hub_membership.membership.id
    }
  }
  project = google_project.project.project_id
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureMciChangeMembership(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_container_cluster" "secondary" {
  name               = "tf-test2%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "tf-test%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_membership" "membership_second" {
  membership_id = "tf-test2%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.secondary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_feature" "feature" {
  name = "multiclusteringress"
  location = "global"
  spec {
    multiclusteringress {
      config_membership = google_gke_hub_membership.membership_second.id
    }
  }
  labels = {
    foo = "bar"
  }
  project = google_project.project.project_id
}
`, context)
}

func TestAccGKEHubFeature_FleetDefaultMemberConfigServiceMesh(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_FleetDefaultMemberConfigServiceMesh(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccGKEHubFeature_FleetDefaultMemberConfigServiceMeshUpdate(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeature_FleetDefaultMemberConfigServiceMesh(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "servicemesh"
  location = "global"
  fleet_default_member_config {
    mesh {
      management = "MANAGEMENT_AUTOMATIC"
    }
  }
  depends_on = [google_project_service.anthos, google_project_service.gkehub, google_project_service.mesh]
  project = google_project.project.project_id
}
`, context)
}

func testAccGKEHubFeature_FleetDefaultMemberConfigServiceMeshUpdate(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "servicemesh"
  location = "global"
  fleet_default_member_config {
    mesh {
      management = "MANAGEMENT_MANUAL"
    }
  }
  depends_on = [google_project_service.anthos, google_project_service.gkehub, google_project_service.mesh]
  project = google_project.project.project_id
}
`, context)
}

func TestAccGKEHubFeature_FleetDefaultMemberConfigConfigManagement(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_FleetDefaultMemberConfigConfigManagement(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccGKEHubFeature_FleetDefaultMemberConfigConfigManagementUpdate(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeature_FleetDefaultMemberConfigConfigManagement(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "configmanagement"
  location = "global"
  fleet_default_member_config {
    configmanagement {
      version = "1.16.0"
      config_sync {
        source_format = "hierarchy"
        git {
          sync_repo = "https://github.com/GoogleCloudPlatform/magic-modules"
          sync_branch = "master"
          policy_dir = "."
          sync_rev = "HEAD"
          secret_type = "none"
          sync_wait_secs = "15"
        }
      }
    }
  }
  depends_on = [google_project_service.anthos, google_project_service.gkehub, google_project_service.acm]
  project = google_project.project.project_id
}
`, context)
}

func testAccGKEHubFeature_FleetDefaultMemberConfigConfigManagementUpdate(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "configmanagement"
  location = "global"
  fleet_default_member_config {
    configmanagement {
      version = "1.16.1"
      config_sync {
        source_format = "unstructured"
        oci {
          sync_repo = "us-central1-docker.pkg.dev/corp-gke-build-artifacts/acm/configs:latest"
          policy_dir = "/acm/nonprod-root/"
          secret_type = "gcpserviceaccount"
          sync_wait_secs = "15"
          gcp_service_account_email = "gke-cluster@gke-foo-nonprod.iam.gserviceaccount.com"
        }
      }
    }
  }
  depends_on = [google_project_service.anthos, google_project_service.gkehub, google_project_service.acm]
  project = google_project.project.project_id
}
`, context)
}

func TestAccGKEHubFeature_FleetDefaultMemberConfigPolicyController(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_FleetDefaultMemberConfigPolicyController(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "update_time"},
			},
			{
				Config: testAccGKEHubFeature_FleetDefaultMemberConfigPolicyControllerUpdate(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeature_FleetDefaultMemberConfigPolicyControllerUpdateSetEmpty(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeature_FleetDefaultMemberConfigPolicyController(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "policycontroller"
  location = "global"
  fleet_default_member_config {
    policycontroller {
      policy_controller_hub_config {
        install_spec = "INSTALL_SPEC_ENABLED"
        exemptable_namespaces = ["foo"]
        policy_content {
          bundles {
            bundle = "policy-essentials-v2022"
            exempted_namespaces = ["foo", "bar"]
          }
        }
        audit_interval_seconds = 30
        referential_rules_enabled = true
      }
    }
  }
  depends_on = [google_project_service.anthos, google_project_service.gkehub, google_project_service.poco]
  project = google_project.project.project_id
}
`, context)
}

func testAccGKEHubFeature_FleetDefaultMemberConfigPolicyControllerUpdate(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "policycontroller"
  location = "global"
  fleet_default_member_config {
    policycontroller {
      policy_controller_hub_config {
        install_spec = "INSTALL_SPEC_SUSPENDED"
        policy_content {
          bundles {
            bundle = "pci-dss-v3.2.1"
            exempted_namespaces = ["baz", "bar"]
          }
          bundles {
            bundle = "nist-sp-800-190"
            exempted_namespaces = []
          }
          template_library {
            installation = "ALL"
          }
        }
        constraint_violation_limit = 50
        referential_rules_enabled = true
        log_denies_enabled = true
        mutation_enabled = true
        deployment_configs {
          component = "admission"
          replica_count = 2
          pod_affinity = "ANTI_AFFINITY"
        }
        deployment_configs {
          component = "audit"
          container_resources {
            limits {
              memory = "1Gi"
              cpu = "1.5"
            }
            requests {
              memory = "500Mi"
              cpu = "150m"
            }
          }
          pod_toleration {
            key = "key1"
            operator = "Equal"
            value = "value1"
            effect = "NoSchedule"
          }
        }
        monitoring {
          backends = [
            "PROMETHEUS"
          ]
        }
      }
    }
  }
  depends_on = [google_project_service.anthos, google_project_service.gkehub, google_project_service.poco]
  project = google_project.project.project_id
}
`, context)
}

func testAccGKEHubFeature_FleetDefaultMemberConfigPolicyControllerUpdateSetEmpty(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "policycontroller"
  location = "global"
  fleet_default_member_config {
    policycontroller {
      policy_controller_hub_config {
        install_spec = "INSTALL_SPEC_ENABLED"
        policy_content {}
        constraint_violation_limit = 50
        referential_rules_enabled = true
        log_denies_enabled = true
        mutation_enabled = true
        deployment_configs {
          component = "admission"
        }
        monitoring {}
      }
    }
  }
  depends_on = [google_project_service.anthos, google_project_service.gkehub, google_project_service.poco]
  project = google_project.project.project_id
}
`, context)
}

func TestAccGKEHubFeature_gkehubFeatureMcsd(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_gkehubFeatureMcsd(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "labels", "terraform_labels"},
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureMcsdUpdate(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEHubFeature_gkehubFeatureMcsd(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = "projects/${google_project.project.project_id}"
  labels = {
    foo = "bar"
  }
  depends_on = [google_project_service.mcsd]
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureMcsdUpdate(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = google_project.project.project_id
  labels = {
    foo = "quux"
    baz = "qux"
  }
  depends_on = [google_project_service.mcsd]
}
`, context)
}

func gkeHubFeatureProjectSetupForGA(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "tf-test-gkehub%{random_suffix}"
  project_id      = "tf-test-gkehub%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "mesh" {
  project = google_project.project.project_id
  service = "meshconfig.googleapis.com"
}

resource "google_project_service" "mci" {
  project = google_project.project.project_id
  service = "multiclusteringress.googleapis.com"
}

resource "google_project_service" "acm" {
  project = google_project.project.project_id
  service = "anthosconfigmanagement.googleapis.com"
}

resource "google_project_service" "poco" {
  project = google_project.project.project_id
  service = "anthospolicycontroller.googleapis.com"
}

resource "google_project_service" "mcsd" {
  project = google_project.project.project_id
  service = "multiclusterservicediscovery.googleapis.com"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "container" {
  project = google_project.project.project_id
  service = "container.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "anthos" {
  project = google_project.project.project_id
  service = "anthos.googleapis.com"
}

resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}
`, context)
}

func testAccCheckGKEHubFeatureDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gke_hub_feature" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{GKEHub2BasePath}}projects/{{project}}/locations/{{location}}/features/{{name}}")
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
				return fmt.Errorf("GKEHubFeature still exists at %s", url)
			}
		}

		return nil
	}
}
