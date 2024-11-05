// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	gkehub "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/gkehub"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccGKEHubFeatureMembership_gkehubFeatureAcmUpdate(t *testing.T) {
	// Multiple fine-grained resources cause VCR to fail
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmUpdateStart(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test2%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmMembershipUpdate(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test2%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member_2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmAddHierarchyController(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipNotPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test2%s", context["random_suffix"])),
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test3%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member_3",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmRemoveFields(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipNotPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test2%s", context["random_suffix"])),
					testAccCheckGkeHubFeatureMembershipNotPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("basic1%s", context["random_suffix"])),
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test3%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member_3",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmUpdateStart(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member_1" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      source_format = "hierarchy"
      git {
        sync_repo   = "https://github.com/GoogleCloudPlatform/magic-modules"
        secret_type = "none"
      }
    }
  }
}

resource "google_gke_hub_feature_membership" "feature_member_2" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_second.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      source_format = "hierarchy"
      git {
        sync_repo   = "https://github.com/terraform-providers/terraform-provider-google"
        secret_type = "none"
      }
    }
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmMembershipUpdate(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "changed"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member_1" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      source_format = "hierarchy"
      enabled       = true
      git {
        sync_repo   = "https://github.com/GoogleCloudPlatform/magic-modules"
        secret_type = "none"
      }
    }
    management = "MANAGEMENT_AUTOMATIC"
  }
}

resource "google_gke_hub_feature_membership" "feature_member_2" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_second.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      source_format = "hierarchy"
      git {
        sync_repo   = "https://github.com/terraform-providers/terraform-provider-google-beta"
        secret_type = "none"
      }
    }
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmAddHierarchyController(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "changed"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member_2" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_second.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      source_format = "unstructured"
      git {
        sync_repo   = "https://github.com/terraform-providers/terraform-provider-google-beta"
        secret_type = "none"
      }
    }
    hierarchy_controller {
      enable_hierarchical_resource_quota = true
      enable_pod_tree_labels = false
      enabled = true
    }
  }
}

resource "google_gke_hub_feature_membership" "feature_member_3" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_third.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      source_format = "hierarchy"
      git {
        sync_repo   = "https://github.com/hashicorp/terraform"
        secret_type = "none"
      }
    }
    hierarchy_controller {
      enable_hierarchical_resource_quota = false
      enable_pod_tree_labels = true
      enabled = false
    }
  }
}

resource "google_gke_hub_feature_membership" "feature_member_4" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_fourth.membership_id
  configmanagement {
    version = "1.18.2"
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmRemoveFields(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "changed"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member_3" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_third.membership_id
  configmanagement {
    version = "1.18.2"
  }
}
`, context)
}

func TestAccGKEHubFeatureMembership_gkehubFeatureAcmAllFields(t *testing.T) {
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmFewFields(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmAllFields(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmFewFields(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureWithPreventDriftField(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmAllFields(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  project = google_project.project.project_id
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      git {
        sync_repo      = "https://github.com/hashicorp/terraform"
        https_proxy    = "https://example.com"
        policy_dir     = "google/"
        secret_type    = "none"
        sync_branch    = "some-branch"
        sync_rev       = "v3.60.0"
        sync_wait_secs = "30"
      }
    }
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_gkehubFeatureWithPreventDriftField(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  project = google_project.project.project_id
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      git {
        sync_repo      = "https://github.com/hashicorp/terraform"
        https_proxy    = "https://example.com"
        policy_dir     = "google/"
        secret_type    = "none"
        sync_branch    = "some-branch"
        sync_rev       = "v3.60.0"
        sync_wait_secs = "30"
      }
      prevent_drift = true
    }
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmFewFields(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  project = google_project.project.project_id
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_service_account" "feature_sa" {
  project = google_project.project.project_id
  account_id = "feature-sa"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      git {
        sync_repo   = "https://github.com/hashicorp/terraform"
        secret_type = "none"
      }
    }
  }
}
`, context)
}

func TestAccGKEHubFeatureMembership_gkehubFeatureAcmOci(t *testing.T) {
	// Multiple fine-grained resources cause VCR to fail
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmOciStart(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmOciUpdate(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_gkehubFeatureAcmOciRemoveFields(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "configmanagement", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmOciStart(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup_ACMOCI(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_service_account" "feature_sa" {
  project = google_project.project.project_id
  account_id = "feature-sa"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_acmoci.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      source_format = "unstructured"
      oci {
        sync_repo = "us-central1-docker.pkg.dev/sample-project/config-repo/config-sync-gke:latest"
        policy_dir = "config-connector"
        sync_wait_secs = "20"
        secret_type = "gcpserviceaccount"
        gcp_service_account_email = google_service_account.feature_sa.email
      }
      prevent_drift = true
    }
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmOciUpdate(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup_ACMOCI(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_service_account" "feature_sa" {
  project = google_project.project.project_id
  account_id = "feature-sa"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_acmoci.membership_id
  configmanagement {
    version = "1.18.2"
    config_sync {
      enabled = true
      source_format = "hierarchy"
      oci {
        sync_repo = "us-central1-docker.pkg.dev/sample-project/config-repo/config-sync-gke:latest"
        policy_dir = "config-sync"
        sync_wait_secs = "15"
        secret_type = "gcenode"
        gcp_service_account_email = google_service_account.feature_sa.email
      }
      prevent_drift = true
    }
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_gkehubFeatureAcmOciRemoveFields(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup_ACMOCI(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_service_account" "feature_sa" {
  project = google_project.project.project_id
  account_id = "feature-sa"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership_acmoci.membership_id
  configmanagement {
    version = "1.18.2"
  }
}
`, context)
}

func TestAccGKEHubFeatureMembership_gkehubFeatureMesh(t *testing.T) {
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeatureMembership_meshStart(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "servicemesh", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_meshUpdateManagement(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "servicemesh", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_meshUpdateControlPlane(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "servicemesh", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeatureMembership_meshStart(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  project = google_project.project.project_id
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "servicemesh"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_service_account" "feature_sa" {
  project = google_project.project.project_id
  account_id = "feature-sa"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  mesh {
    management = "MANAGEMENT_AUTOMATIC"
    control_plane = "AUTOMATIC"
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_meshUpdateManagement(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  project = google_project.project.project_id
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "servicemesh"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_service_account" "feature_sa" {
  project = google_project.project.project_id
  account_id = "feature-sa"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  mesh {
    management = "MANAGEMENT_MANUAL"
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_meshUpdateControlPlane(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  project = google_project.project.project_id
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "servicemesh"
  location = "global"

  labels = {
    foo = "bar"
  }
  depends_on = [time_sleep.wait_120s]
}

resource "google_service_account" "feature_sa" {
  project = google_project.project.project_id
  account_id = "feature-sa"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  mesh {
    control_plane = "MANUAL"
  }
}
`, context)
}

func TestAccGKEHubFeatureMembership_gkehubFeaturePolicyController(t *testing.T) {
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeatureMembership_policycontrollerStart(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "policycontroller", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_policycontrollerUpdateDefaultFields(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "policycontroller", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeatureMembership_policycontrollerUpdateMaps(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGkeHubFeatureMembershipPresent(t, fmt.Sprintf("tf-test-gkehub%s", context["random_suffix"]), "global", "policycontroller", fmt.Sprintf("tf-test1%s", context["random_suffix"])),
				),
			},
			{
				ResourceName:      "google_gke_hub_feature_membership.feature_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeatureMembership_policycontrollerStart(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "policycontroller"
  location = "global"
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  policycontroller {
    policy_controller_hub_config {
      install_spec = "INSTALL_SPEC_ENABLED"
      exemptable_namespaces = ["foo"]
      audit_interval_seconds = 30
      referential_rules_enabled = true
    }
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_policycontrollerUpdateDefaultFields(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "policycontroller"
  location = "global"
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  policycontroller {
    policy_controller_hub_config {
      install_spec = "INSTALL_SPEC_SUSPENDED"
      constraint_violation_limit = 50
      referential_rules_enabled = true
      log_denies_enabled = true
      mutation_enabled = true
      monitoring {
        backends = [
          "PROMETHEUS"
        ]
      }
      deployment_configs {
        component_name = "admission"
        replica_count = 3
        pod_affinity = "ANTI_AFFINITY"
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
        pod_tolerations {
          key = "key1"
          operator = "Equal"
          value = "value1"
          effect = "NoSchedule"
        }
      }
      deployment_configs {
        component_name = "mutation"
        replica_count = 3
        pod_affinity = "ANTI_AFFINITY"
      }
      policy_content {
        template_library {
          installation = "ALL"
        }
        bundles {
          bundle_name = "pci-dss-v3.2.1"
          exempted_namespaces = ["sample-namespace"]
        }
        bundles {
          bundle_name = "nist-sp-800-190"
        }
      }
    }
    version = "1.17.0"
  }
}
`, context)
}

func testAccGKEHubFeatureMembership_policycontrollerUpdateMaps(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetup(context) + gkeHubClusterMembershipSetup(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  project = google_project.project.project_id
  name = "policycontroller"
  location = "global"
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_feature_membership" "feature_member" {
  project = google_project.project.project_id
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  policycontroller {
    policy_controller_hub_config {
      install_spec = "INSTALL_SPEC_SUSPENDED"
      constraint_violation_limit = 50
      referential_rules_enabled = true
      log_denies_enabled = true
      mutation_enabled = true
      monitoring {
        backends = [
          "PROMETHEUS"
        ]
      }
      deployment_configs {
        component_name = "admission"
        pod_affinity = "NO_AFFINITY"
      }
      deployment_configs {
        component_name = "audit"
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
      }
    }
    version = "1.17.0"
  }
}
`, context)
}

func gkeHubClusterMembershipSetup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_container_cluster" "secondary" {
  name               = "tf-test-cl2%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_container_cluster" "tertiary" {
  name               = "tf-test-cl3%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}


resource "google_container_cluster" "quarternary" {
  name               = "tf-test-cl4%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
}

resource "google_gke_hub_membership" "membership_second" {
  project = google_project.project.project_id
  membership_id = "tf-test2%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.secondary.id}"
    }
  }
}

resource "google_gke_hub_membership" "membership_third" {
  project = google_project.project.project_id
  membership_id = "tf-test3%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.tertiary.id}"
    }
  }
}

resource "google_gke_hub_membership" "membership_fourth" {
  project = google_project.project.project_id
  membership_id = "tf-test4%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.quarternary.id}"
    }
  }
}
`, context)
}

func gkeHubClusterMembershipSetup_ACMOCI(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_network" "testnetwork" {
    project                 = google_project.project.project_id
    name                    = "testnetwork"
    auto_create_subnetworks = true
    depends_on = [time_sleep.wait_120s]
}

resource "google_container_cluster" "container_acmoci" {
  name               = "tf-test-cl%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  network = google_compute_network.testnetwork.self_link
  project = google_project.project.project_id
  deletion_protection = false
  depends_on = [time_sleep.wait_120s]
}

resource "google_gke_hub_membership" "membership_acmoci" {
  project = google_project.project.project_id
  membership_id = "tf-test1%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.container_acmoci.id}"
    }
  }
}
`, context)
}

func testAccCheckGkeHubFeatureMembershipPresent(t *testing.T, project, location, feature, membership string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		obj := &gkehub.FeatureMembership{
			Feature:    dcl.StringOrNil(feature),
			Location:   dcl.StringOrNil(location),
			Membership: dcl.StringOrNil(membership),
			Project:    dcl.String(project),
		}

		_, err := transport_tpg.NewDCLGkeHubClient(config, "", "", 0).GetFeatureMembership(context.Background(), obj)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckGkeHubFeatureMembershipNotPresent(t *testing.T, project, location, feature, membership string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		obj := &gkehub.FeatureMembership{
			Feature:    dcl.StringOrNil(feature),
			Location:   dcl.StringOrNil(location),
			Membership: dcl.StringOrNil(membership),
			Project:    dcl.String(project),
		}

		_, err := transport_tpg.NewDCLGkeHubClient(config, "", "", 0).GetFeatureMembership(context.Background(), obj)
		if err == nil {
			return fmt.Errorf("Did not expect to find GKE Feature Membership for projects/%s/locations/%s/features/%s/membershipId/%s", project, location, feature, membership)
		}
		if dcl.IsNotFound(err) {
			return nil
		}
		return err
	}
}

// Copy this function from the package gkehub2_test to here
func gkeHubFeatureProjectSetup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "tf-test-gkehub%{random_suffix}"
  project_id      = "tf-test-gkehub%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "anthos" {
  project = google_project.project.project_id
  service = "anthos.googleapis.com"
  disable_on_destroy = false
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

resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}

// It needs waiting until the API services are really activated.
resource "time_sleep" "wait_120s" {
  create_duration = "120s"
  depends_on = [
    google_project_service.anthos,
    google_project_service.mesh,
    google_project_service.mci,
    google_project_service.acm,
    google_project_service.poco,
    google_project_service.mcsd,
    google_project_service.compute,
    google_project_service.container,
    google_project_service.gkehub,
  ]
}
`, context)
}

// Copy this function from the package gkehub2_test to here
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
