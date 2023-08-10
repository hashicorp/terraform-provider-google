// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccContainerCluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_cluster.primary", "services_ipv4_cidr"),
					resource.TestCheckResourceAttrSet("google_container_cluster.primary", "self_link"),
				),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportStateId:     fmt.Sprintf("us-central1-a/%s", clusterName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportStateId:     fmt.Sprintf("%s/us-central1-a/%s", envvar.GetTestProjectFromEnv(), clusterName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_networkingModeRoutes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_networkingModeRoutes(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_misc(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_misc(clusterName),
				// Explicitly check removing the default node pool since we won't
				// catch it by just importing.
				Check: resource.TestCheckResourceAttr(
					"google_container_cluster.primary", "node_pool.#", "0"),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_misc_update(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
		},
	})
}

func TestAccContainerCluster_withAddons(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAddons(pid, clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: clean up this list in `4.0.0`, remove both `workload_identity_config` fields (same for below)
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_updateAddons(pid, clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			// Issue with cloudrun_config addon: https://github.com/hashicorp/terraform-provider-google/issues/11943
			// {
			// 	Config: testAccContainerCluster_withInternalLoadBalancer(pid, clusterName),
			// },
			// {
			// 	ResourceName:            "google_container_cluster.primary",
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"min_master_version"},
			// },
		},
	})
}

func TestAccContainerCluster_withNotificationConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	newTopic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNotificationConfig(clusterName, topic),
			},
			{
				ResourceName:      "google_container_cluster.notification_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withNotificationConfig(clusterName, newTopic),
			},
			{
				ResourceName:      "google_container_cluster.notification_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_disableNotificationConfig(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.notification_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withNotificationConfig(clusterName, newTopic),
			},
			{
				ResourceName:      "google_container_cluster.notification_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withFilteredNotificationConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	newTopic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withFilteredNotificationConfig(clusterName, topic),
			},
			{
				ResourceName:      "google_container_cluster.filtered_notification_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withFilteredNotificationConfigUpdate(clusterName, newTopic),
			},
			{
				ResourceName:      "google_container_cluster.filtered_notification_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_disableFilteredNotificationConfig(clusterName, newTopic),
			},
			{
				ResourceName:      "google_container_cluster.filtered_notification_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withConfidentialNodes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withConfidentialNodes(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.confidential_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_disableConfidentialNodes(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.confidential_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withConfidentialNodes(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.confidential_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withILBSubsetting(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_disableILBSubSetting(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.confidential_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withILBSubSetting(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.confidential_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_disableILBSubSetting(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.confidential_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuthConfig_NoCert(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuthNoCert(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth_no_cert", "master_auth.0.client_certificate", ""),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_master_auth_no_cert",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withAuthenticatorGroupsConfig(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	orgDomain := envvar.GetTestOrgDomainFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.primary",
						"authenticator_groups_config.0.enabled"),
				),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withAuthenticatorGroupsConfigUpdate(clusterName, orgDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary",
						"authenticator_groups_config.0.security_group", fmt.Sprintf("gke-security-groups@%s", orgDomain)),
				),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withAuthenticatorGroupsConfigUpdate2(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.primary",
						"authenticator_groups_config.0.enabled"),
				),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNetworkPolicyEnabled(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNetworkPolicyEnabled(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.#", "1"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_removeNetworkPolicy(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_withNetworkPolicyDisabled(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_withNetworkPolicyConfigDisabled(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"addons_config.0.network_policy_config.0.disabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config:             testAccContainerCluster_withNetworkPolicyConfigDisabled(clusterName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccContainerCluster_withReleaseChannelEnabled(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "STABLE"),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "UNSPECIFIED"),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withReleaseChannelEnabledDefaultVersion(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withReleaseChannelEnabledDefaultVersion(clusterName, "REGULAR"),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "REGULAR"),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "UNSPECIFIED"),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withInvalidReleaseChannel(t *testing.T) {
	// This is essentially a unit test, no interactions
	acctest.SkipIfVcr(t)
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withReleaseChannelEnabled(clusterName, "CANARY"),
				ExpectError: regexp.MustCompile(`expected release_channel\.0\.channel to be one of \[UNSPECIFIED RAPID REGULAR STABLE\], got CANARY`),
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuthorizedNetworksConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{}, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.#", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks.#", "0"),
				),
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{"8.8.8.8/32"}, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks.#", "1"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_master_authorized_networks",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{"10.0.0.0/8", "8.8.8.8/32"}, ""),
			},
			{
				ResourceName:      "google_container_cluster.with_master_authorized_networks",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{}, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks.#", "0"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_master_authorized_networks",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_removeMasterAuthorizedNetworksConfig(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_master_authorized_networks",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withGcpPublicCidrsAccessEnabledToggle(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withoutGcpPublicCidrsAccessEnabled(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_gcp_public_cidrs_access_enabled",
						"master_authorized_networks_config.#", "0"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_gcp_public_cidrs_access_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withGcpPublicCidrsAccessEnabled(clusterName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_gcp_public_cidrs_access_enabled",
						"master_authorized_networks_config.0.gcp_public_cidrs_access_enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_gcp_public_cidrs_access_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withGcpPublicCidrsAccessEnabled(clusterName, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_gcp_public_cidrs_access_enabled",
						"master_authorized_networks_config.0.gcp_public_cidrs_access_enabled", "true"),
				),
			},
		},
	})
}

func testAccContainerCluster_withGcpPublicCidrsAccessEnabled(clusterName string, flag string) string {

	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_gcp_public_cidrs_access_enabled" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1

  master_authorized_networks_config {
    gcp_public_cidrs_access_enabled = %s
  }
}
`, clusterName, flag)
}

func testAccContainerCluster_withoutGcpPublicCidrsAccessEnabled(clusterName string) string {

	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_gcp_public_cidrs_access_enabled" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1
}
`, clusterName)
}

func TestAccContainerCluster_regional(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-regional-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_regional(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.regional",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_regionalWithNodePool(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-regional-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_regionalWithNodePool(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.regional",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_regionalWithNodeLocations(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_regionalNodeLocations(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_node_locations",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_regionalUpdateNodeLocations(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_node_locations",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withPrivateClusterConfigBasic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateClusterConfig(containerNetName, clusterName, false),
			},
			{
				ResourceName:      "google_container_cluster.with_private_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withPrivateClusterConfig(containerNetName, clusterName, true),
			},
			{
				ResourceName:      "google_container_cluster.with_private_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(containerNetName, clusterName, "us-central1-a", false),
				ExpectError: regexp.MustCompile("master_ipv4_cidr_block must be set if enable_private_nodes is true"),
			},
		},
	})
}

func TestAccContainerCluster_withPrivateClusterConfigMissingCidrBlock_withAutopilot(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(containerNetName, clusterName, "us-central1", true),
			},
			{
				ResourceName:      "google_container_cluster.with_private_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(clusterName, true),
			},
			{
				ResourceName:      "google_container_cluster.with_private_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(clusterName, false),
			},
			{
				ResourceName:      "google_container_cluster.with_private_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withIntraNodeVisibility(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIntraNodeVisibility(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_intranode_visibility", "enable_intranode_visibility", "true"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_intranode_visibility",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_updateIntraNodeVisibility(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_intranode_visibility", "enable_intranode_visibility", "false"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_intranode_visibility",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withVersion(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_version",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_updateVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLowerVersion(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_version",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_updateVersion(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_version",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfig(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withNodeConfigUpdate(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withLoggingVariantInNodeConfig(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLoggingVariantInNodeConfig(clusterName, "MAX_THROUGHPUT"),
			},
			{
				ResourceName:      "google_container_cluster.with_logging_variant_in_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withLoggingVariantInNodePool(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePoolName := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLoggingVariantInNodePool(clusterName, nodePoolName, "MAX_THROUGHPUT"),
			},
			{
				ResourceName:      "google_container_cluster.with_logging_variant_in_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withLoggingVariantUpdates(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, "DEFAULT"),
			},
			{
				ResourceName:      "google_container_cluster.with_logging_variant_node_pool_default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, "MAX_THROUGHPUT"),
			},
			{
				ResourceName:      "google_container_cluster.with_logging_variant_node_pool_default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, "DEFAULT"),
			},
			{
				ResourceName:      "google_container_cluster.with_logging_variant_node_pool_default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigScopeAlias(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigScopeAlias(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_node_config_scope_alias",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigShieldedInstanceConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigShieldedInstanceConfig(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigReservationAffinity(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigReservationAffinity(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.#", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.0.consume_reservation_type", "ANY_RESERVATION"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigReservationAffinitySpecific(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-reservation-%s", acctest.RandString(t, 10))
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigReservationAffinitySpecific(reservationName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.#", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.0.consume_reservation_type", "SPECIFIC_RESERVATION"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.0.key", "compute.googleapis.com/reservation-name"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.0.values.#", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.0.values.0", reservationName),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withWorkloadMetadataConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withWorkloadMetadataConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_workload_metadata_config",
						"node_config.0.workload_metadata_config.0.mode", "GCE_METADATA"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_metadata_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withBootDiskKmsKey(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withBootDiskKmsKey(clusterName, kms.CryptoKey.Name),
			},
			{
				ResourceName:            "google_container_cluster.with_boot_disk_kms_key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_network(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_networkRef(clusterName, network),
			},
			{
				ResourceName:      "google_container_cluster.with_net_ref_by_url",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_container_cluster.with_net_ref_by_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_backend(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_backendRef(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolBasic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolBasic(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolUpdateVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolLowerVersion(clusterName, npName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withNodePoolUpdateVersion(clusterName, npName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolResize(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNodeLocations(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "2"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withNodePoolResize(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "3"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolAutoscaling(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolAutoscaling(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "3"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withNodePoolUpdateAutoscaling(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "5"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withNodePoolBasic(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolCIA(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerRegionalCluster_withNodePoolCIA(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.total_min_node_count", "3"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.total_max_node_count", "21"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.location_policy", "BALANCED"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerRegionalClusterUpdate_withNodePoolCIA(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "0"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.total_min_node_count", "4"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.total_max_node_count", "32"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.location_policy", "ANY"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerRegionalCluster_withNodePoolBasic(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.total_min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.total_max_node_count"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolNamePrefix(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npNamePrefix := "tf-test-np-"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNamePrefix(clusterName, npNamePrefix),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool_name_prefix",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_pool.0.name_prefix"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolMultiple(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npNamePrefix := "tf-test-np-"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolMultiple(clusterName, npNamePrefix),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool_multiple",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolConflictingNameFields(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npPrefix := "tf-test-np"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withNodePoolConflictingNameFields(clusterName, npPrefix),
				ExpectError: regexp.MustCompile("Cannot specify both name and name_prefix for a node_pool"),
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolNodeConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNodeConfig(cluster, np),
			},
			{
				ResourceName:      "google_container_cluster.with_node_pool_node_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withMaintenanceWindow(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_window"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMaintenanceWindow(clusterName, "03:00"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withMaintenanceWindow(clusterName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.daily_maintenance_window.0.start_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// maintenance_policy.# = 0 is equivalent to no maintenance policy at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"maintenance_policy.#"},
			},
		},
	})
}

func TestAccContainerCluster_withRecurringMaintenanceWindow(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_recurring_maintenance_window"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withRecurringMaintenanceWindow(cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.daily_maintenance_window.0.start_time"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withRecurringMaintenanceWindow(cluster, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.daily_maintenance_window.0.start_time"),
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.recurring_window.0.start_time"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
				// maintenance_policy.# = 0 is equivalent to no maintenance policy at all,
				// but will still cause an import diff
				ImportStateVerifyIgnore: []string{"maintenance_policy.#"},
			},
		},
	})
}

func TestAccContainerCluster_withMaintenanceExclusionWindow(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_window"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusion_RecurringMaintenanceWindow(cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z"),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withExclusion_DailyMaintenanceWindow(cluster, "2020-01-01T00:00:00Z", "2020-01-02T00:00:00Z"),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withMaintenanceExclusionOptions(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_options"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_MINOR_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_MINOR_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_OR_NODE_UPGRADES"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_deleteMaintenanceExclusionOptions(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_options"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_OR_NODE_UPGRADES"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_NoExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope"),
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_updateMaintenanceExclusionOptions(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_options"

	// step1: create a new cluster and initialize the maintenceExclusion without exclusion scopes,
	// step2: add exclusion scopes to the maintenancePolicy,
	// step3: update the maintenceExclusion with new scopes
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_NoExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope"),
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_MINOR_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_MINOR_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_OR_NODE_UPGRADES"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_UPGRADES", "NO_MINOR_UPGRADES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_UPGRADES"),
				),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_deleteExclusionWindow(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_window"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusion_DailyMaintenanceWindow(cluster, "2020-01-01T00:00:00Z", "2020-01-02T00:00:00Z"),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withExclusion_RecurringMaintenanceWindow(cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z"),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withExclusion_NoMaintenanceWindow(cluster, "2020-01-01T00:00:00Z", "2020-01-02T00:00:00Z"),
			},
			{
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withIPAllocationPolicy_existingSecondaryRanges(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy_existingSecondaryRanges(containerNetName, clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_ip_allocation_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withIPAllocationPolicy_specificIPRanges(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy_specificIPRanges(containerNetName, clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_ip_allocation_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withIPAllocationPolicy_specificSizes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy_specificSizes(containerNetName, clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_ip_allocation_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_stackType_withDualStack(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_stack_type"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_stackType_withDualStack(containerNetName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip_allocation_policy.0.stack_type", "IPV4_IPV6"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_stackType_withSingleStack(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_stack_type"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_stackType_withSingleStack(containerNetName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip_allocation_policy.0.stack_type", "IPV4"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_with_PodCIDROverprovisionDisabled(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_pco_disabled"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_with_PodCIDROverprovisionDisabled(containerNetName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ip_allocation_policy.0.pod_cidr_overprovision_config.0.disabled", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioning(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioning(clusterName, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning",
						"cluster_autoscaling.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autoprovisioning(clusterName, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning",
						"cluster_autoscaling.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaults(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	includeMinCpuPlatform := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaults(clusterName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning",
						"cluster_autoscaling.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config:             testAccContainerCluster_autoprovisioningDefaults(clusterName, true),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsMinCpuPlatform(clusterName, includeMinCpuPlatform),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsMinCpuPlatform(clusterName, !includeMinCpuPlatform),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(clusterName, 2, 1, "SURGE"),
			},
			{
				ResourceName:      "google_container_cluster.with_autoprovisioning_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(clusterName, 2, 1, "BLUE_GREEN"),
				ExpectError: regexp.MustCompile(`Surge upgrade settings max_surge/max_unavailable can only be used when strategy is set to SURGE`),
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsUpgradeSettingsWithBlueGreenStrategy(clusterName, "3.500s", "BLUE_GREEN"),
			},
			{
				ResourceName:      "google_container_cluster.with_autoprovisioning_upgrade_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withShieldedNodes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withShieldedNodes(clusterName, true),
			},
			{
				ResourceName:      "google_container_cluster.with_shielded_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withShieldedNodes(clusterName, false),
			},
			{
				ResourceName:      "google_container_cluster.with_shielded_nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withAutopilot(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutopilot(pid, containerNetName, clusterName, "us-central1", true, false, ""),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerClusterCustomServiceAccount_withAutopilot(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	serviceAccountName := fmt.Sprintf("tf-test-sa-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutopilot(pid, containerNetName, clusterName, "us-central1", true, false, serviceAccountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autopilot",
						"cluster_autoscaling.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_container_cluster.with_autopilot",
						"cluster_autoscaling.0.auto_provisioning_defaults.0.service_account",
						fmt.Sprintf("%s@%s.iam.gserviceaccount.com", serviceAccountName, pid)),
					resource.TestCheckResourceAttr("google_container_cluster.with_autopilot",
						"cluster_autoscaling.0.auto_provisioning_defaults.0.oauth_scopes.0", "https://www.googleapis.com/auth/cloud-platform"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_errorAutopilotLocation(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withAutopilot(pid, containerNetName, clusterName, "us-central1-a", true, false, ""),
				ExpectError: regexp.MustCompile(`Autopilot clusters must be regional clusters.`),
			},
		},
	})
}

func TestAccContainerCluster_withWorkloadIdentityConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withWorkloadIdentityConfigEnabled(pid, clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_identity_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_updateWorkloadIdentityConfig(pid, clusterName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_identity_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_updateWorkloadIdentityConfig(pid, clusterName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_identity_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
		},
	})
}

func TestAccContainerCluster_withLoggingConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withLoggingConfigEnabled(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withLoggingConfigDisabled(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withLoggingConfigUpdated(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_basic(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withMonitoringConfigAdvancedDatapathObservabilityConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMonitoringConfigAdvancedDatapathObservabilityConfigEnabled(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigAdvancedDatapathObservabilityConfigDisabled(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withMonitoringConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigEnabled(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigDisabled(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigUpdated(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigPrometheusUpdated(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			// Back to basic settings to test setting Prometheus on its own
			{
				Config: testAccContainerCluster_basic(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigPrometheusOnly(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigPrometheusOnly2(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_basic(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withSoleTenantGroup(t *testing.T) {
	t.Parallel()

	resourceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withSoleTenantGroup(resourceName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsDiskSizeGb(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	includeDiskSizeGb := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskSizeGb(clusterName, includeDiskSizeGb),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskSizeGb(clusterName, !includeDiskSizeGb),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsDiskType(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	includeDiskType := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskType(clusterName, includeDiskType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskType(clusterName, !includeDiskType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsImageType(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	includeImageType := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsImageType(clusterName, includeImageType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsImageType(clusterName, !includeImageType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsBootDiskKmsKey(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsBootDiskKmsKey(clusterName, kms.CryptoKey.Name),
			},
			{
				ResourceName:      "google_container_cluster.nap_boot_disk_kms_key",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"min_master_version",
					"node_pool", // cluster_autoscaling (node auto-provisioning) creates new node pools automatically
				},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsShieldedInstance(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsShieldedInstance(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.nap_shielded_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_autoprovisioningDefaultsManagement(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsManagement(clusterName, false, false),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_management",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsManagement(clusterName, true, true),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_management",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_errorCleanDanglingCluster(t *testing.T) {
	t.Parallel()

	prefix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", prefix)
	clusterNameError := fmt.Sprintf("tf-test-cluster-err-%s", prefix)
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))

	initConfig := testAccContainerCluster_withInitialCIDR(containerNetName, clusterName)
	overlapConfig := testAccContainerCluster_withCIDROverlap(initConfig, clusterNameError)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: initConfig,
			},
			{
				ResourceName:      "google_container_cluster.cidr_error_preempt",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      overlapConfig,
				ExpectError: regexp.MustCompile("Error waiting for creating GKE cluster"),
			},
			// If dangling cluster wasn't deleted, this plan will return an error
			{
				Config:             overlapConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccContainerCluster_errorNoClusterCreated(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withInvalidLocation("wonderland"),
				ExpectError: regexp.MustCompile(`Permission denied on 'locations/wonderland' \(or it may not exist\).`),
			},
		},
	})
}

func TestAccContainerCluster_withExternalIpsConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExternalIpsConfig(pid, clusterName, true),
			},
			{
				ResourceName:      "google_container_cluster.with_external_ips_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withExternalIpsConfig(pid, clusterName, false),
			},
			{
				ResourceName:      "google_container_cluster.with_external_ips_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withMeshCertificatesConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMeshCertificatesConfigEnabled(pid, clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_mesh_certificates_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_updateMeshCertificatesConfig(pid, clusterName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_mesh_certificates_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
			{
				Config: testAccContainerCluster_updateMeshCertificatesConfig(pid, clusterName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_mesh_certificates_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
			},
		},
	})
}

func TestAccContainerCluster_withCostManagementConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_updateCostManagementConfig(pid, clusterName, true),
			},
			{
				ResourceName:      "google_container_cluster.with_cost_management_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_updateCostManagementConfig(pid, clusterName, false),
			},
			{
				ResourceName:      "google_container_cluster.with_cost_management_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withDatabaseEncryption(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	// Use the bootstrapped KMS key so we can avoid creating keys needlessly
	// as they will pile up in the project because they can not be completely
	// deleted.  Also, we need to create the key in the same location as the
	// cluster as GKE does not support the "global" location for KMS keys.
	// See https://cloud.google.com/kubernetes-engine/docs/how-to/encrypting-secrets#creating_a_key
	kmsData := acctest.BootstrapKMSKeyInLocation(t, "us-central1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDatabaseEncryption(clusterName, kmsData),
				Check:  resource.TestCheckResourceAttrSet("data.google_kms_key_ring_iam_policy.test_key_ring_iam_policy", "policy_data"),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_basic(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withAdvancedDatapath(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDatapathProvider(clusterName, "ADVANCED_DATAPATH"),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withResourceUsageExportConfig(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", suffix)
	datesetId := fmt.Sprintf("tf_test_cluster_resource_usage_%s", suffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withResourceUsageExportConfig(clusterName, datesetId, "true"),
			},
			{
				ResourceName:      "google_container_cluster.with_resource_usage_export_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withResourceUsageExportConfig(clusterName, datesetId, "false"),
			},
			{
				ResourceName:      "google_container_cluster.with_resource_usage_export_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_withResourceUsageExportConfigNoConfig(clusterName, datesetId),
			},
			{
				ResourceName:      "google_container_cluster.with_resource_usage_export_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuthorizedNetworksDisabled(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksDisabled(containerNetName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccContainerCluster_masterAuthorizedNetworksDisabled(t, "google_container_cluster.with_private_cluster"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_private_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withEnableKubernetesAlpha(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withEnableKubernetesAlpha(clusterName, npName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withEnableKubernetesBetaAPIs(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withEnableKubernetesBetaAPIs(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withEnableKubernetesBetaAPIsOnExistingCluster(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withoutEnableKubernetesBetaAPIs(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withEnableKubernetesBetaAPIs(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withIPv4Error(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withIPv4Error(clusterName),
				ExpectError: regexp.MustCompile("master_ipv4_cidr_block can only be set if"),
			},
		},
	})
}

func TestAccContainerCluster_withDNSConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	domainName := fmt.Sprintf("tf-test-domain-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDNSConfig(clusterName, "CLOUD_DNS", domainName, "VPC_SCOPE"),
			},
			{
				ResourceName:      "google_container_cluster.with_dns_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withGatewayApiConfig(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withGatewayApiConfig(clusterName, "CANARY"),
				ExpectError: regexp.MustCompile(`expected gateway_api_config\.0\.channel to be one of \[CHANNEL_DISABLED CHANNEL_EXPERIMENTAL CHANNEL_STANDARD\], got CANARY`),
			},
			{
				Config: testAccContainerCluster_withGatewayApiConfig(clusterName, "CHANNEL_DISABLED"),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withGatewayApiConfig(clusterName, "CHANNEL_STANDARD"),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_withSecurityPostureConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_SetSecurityPostureToStandard(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_security_posture_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_SetWorkloadVulnerabilityToStandard(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_security_posture_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_DisableALL(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_security_posture_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerCluster_SetSecurityPostureToStandard(resource_name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	mode = "BASIC"
  }
}
`, resource_name)
}

func testAccContainerCluster_SetWorkloadVulnerabilityToStandard(resource_name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	vulnerability_mode = "VULNERABILITY_BASIC"
  }
}
`, resource_name)
}

func testAccContainerCluster_DisableALL(resource_name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	mode = "DISABLED"
	vulnerability_mode = "VULNERABILITY_DISABLED"
  }
}
`, resource_name)
}

func TestAccContainerCluster_autopilot_minimal(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autopilot_minimal(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_autopilot_net_admin(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autopilot_net_admin(clusterName, true),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autopilot_net_admin(clusterName, false),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_autopilot_net_admin(clusterName, true),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func testAccContainerCluster_masterAuthorizedNetworksDisabled(t *testing.T, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		config := acctest.GoogleProviderConfig(t)
		attributes := rs.Primary.Attributes

		cluster, err := config.NewContainerClient(config.UserAgent).Projects.Zones.Clusters.Get(
			config.Project, attributes["location"], attributes["name"]).Do()
		if err != nil {
			return err
		}

		if cluster.MasterAuthorizedNetworksConfig.Enabled {
			return fmt.Errorf("Cluster's master authorized networks config is enabled, but expected to be disabled.")
		}

		return nil
	}
}

func testAccCheckContainerClusterDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_container_cluster" {
				continue
			}

			attributes := rs.Primary.Attributes
			_, err := config.NewContainerClient(config.UserAgent).Projects.Locations.Clusters.Get(
				fmt.Sprintf("projects/%s/locations/%s/clusters/%s", config.Project, attributes["location"], attributes["name"])).Do()
			if err == nil {
				return fmt.Errorf("Cluster still exists")
			}
		}

		return nil
	}
}

func testAccContainerCluster_basic(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
}
`, name)
}

func testAccContainerCluster_networkingModeRoutes(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  networking_mode    = "ROUTES"
}
`, name)
}

func testAccContainerCluster_misc(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  remove_default_node_pool = true

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]

  enable_legacy_abac      = true

  resource_labels = {
    created-by = "terraform"
  }

  vertical_pod_autoscaling {
    enabled = true
  }

  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }
}
`, name)
}

func testAccContainerCluster_misc_update(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  remove_default_node_pool = true # Not worth updating

  node_locations = [
    "us-central1-f",
    "us-central1-c",
  ]

  enable_legacy_abac      = false

  resource_labels = {
    created-by = "terraform-update"
    new-label  = "update"
  }

  vertical_pod_autoscaling {
    enabled = true
  }

  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }
}
`, name)
}

func testAccContainerCluster_withAddons(projectID string, clusterName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  min_master_version = "latest"

  workload_identity_config {
    workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
  }

  addons_config {
    http_load_balancing {
      disabled = true
    }
    horizontal_pod_autoscaling {
      disabled = true
    }
    network_policy_config {
      disabled = true
    }
    gcp_filestore_csi_driver_config {
      enabled = false
    }
    cloudrun_config {
      disabled = true
    }
    dns_cache_config {
      enabled = false
    }
    gce_persistent_disk_csi_driver_config {
      enabled = false
    }
	gke_backup_agent_config {
	  enabled = false
	}
	config_connector_config {
	  enabled = false
	}
    gcs_fuse_csi_driver_config {
      enabled = false
    }
  }
}
`, projectID, clusterName)
}

func testAccContainerCluster_updateAddons(projectID string, clusterName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  min_master_version = "latest"

  workload_identity_config {
    workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
  }

  addons_config {
    http_load_balancing {
      disabled = false
    }
    horizontal_pod_autoscaling {
      disabled = false
    }
    network_policy_config {
      disabled = false
    }
    gcp_filestore_csi_driver_config {
      enabled = true
    }
    cloudrun_config {
	  # https://github.com/hashicorp/terraform-provider-google/issues/11943
      # disabled = false
      disabled = true
    }
    dns_cache_config {
      enabled = true
    }
    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
	gke_backup_agent_config {
	  enabled = true
	}
	config_connector_config {
	  enabled = true
	}
    gcs_fuse_csi_driver_config {
      enabled = true
    }
  }
}
`, projectID, clusterName)
}

// Issue with cloudrun_config addon: https://github.com/hashicorp/terraform-provider-google/issues/11943/
// func testAccContainerCluster_withInternalLoadBalancer(projectID string, clusterName string) string {
// 	return fmt.Sprintf(`
// data "google_project" "project" {
//   project_id = "%s"
// }

// resource "google_container_cluster" "primary" {
//   name               = "%s"
//   location           = "us-central1-a"
//   initial_node_count = 1

//   min_master_version = "latest"

//   workload_identity_config {
//     workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
//   }

//   addons_config {
//     http_load_balancing {
//       disabled = false
//     }
//     horizontal_pod_autoscaling {
//       disabled = false
//     }
//     network_policy_config {
//       disabled = false
//     }
//     cloudrun_config {
// 	  disabled = false
// 	  load_balancer_type = "LOAD_BALANCER_TYPE_INTERNAL"
//     }
//   }
// }
// `, projectID, clusterName)
// }

func testAccContainerCluster_withNotificationConfig(clusterName string, topic string) string {
	return fmt.Sprintf(`

resource "google_pubsub_topic" "%s" {
  name = "%s"
}

resource "google_container_cluster" "notification_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  notification_config {
	pubsub {
	  enabled = true
	  topic   = google_pubsub_topic.%s.id
	}
  }
}
`, topic, topic, clusterName, topic)
}

func testAccContainerCluster_disableNotificationConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "notification_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  notification_config {
	pubsub {
	  enabled = false
	}
  }
}
`, clusterName)
}

func testAccContainerCluster_withFilteredNotificationConfig(clusterName string, topic string) string {

	return fmt.Sprintf(`

resource "google_pubsub_topic" "%s" {
  name = "%s"
}

resource "google_container_cluster" "filtered_notification_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  notification_config {
	pubsub {
	  enabled = true
	  topic   = google_pubsub_topic.%s.id
	  filter {
		event_type = ["UPGRADE_EVENT", "SECURITY_BULLETIN_EVENT"]
	  }
	}
  }
}
`, topic, topic, clusterName, topic)
}

func testAccContainerCluster_withFilteredNotificationConfigUpdate(clusterName string, topic string) string {

	return fmt.Sprintf(`

resource "google_pubsub_topic" "%s" {
  name = "%s"
}

resource "google_container_cluster" "filtered_notification_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  notification_config {
	pubsub {
	  enabled = true
	  topic   = google_pubsub_topic.%s.id
	  filter {
		event_type = ["UPGRADE_AVAILABLE_EVENT"]
	  }
	}
  }
}
`, topic, topic, clusterName, topic)
}

func testAccContainerCluster_disableFilteredNotificationConfig(clusterName string, topic string) string {

	return fmt.Sprintf(`

resource "google_pubsub_topic" "%s" {
  name = "%s"
}

resource "google_container_cluster" "filtered_notification_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  notification_config {
	pubsub {
	  enabled = true
	  topic   = google_pubsub_topic.%s.id
	}
  }
}
`, topic, topic, clusterName, topic)
}

func testAccContainerCluster_withConfidentialNodes(clusterName string, npName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "confidential_nodes" {
  name               = "%s"
  location           = "us-central1-a"
  release_channel {
    channel = "RAPID"
  }

  node_pool {
    name = "%s"
    initial_node_count = 1
    node_config {
      machine_type = "n2d-standard-2" // can't be e2 because Confidential Nodes require AMD CPUs
    }
  }

  confidential_nodes {
    enabled = true
  }
}
`, clusterName, npName)
}

func testAccContainerCluster_disableConfidentialNodes(clusterName string, npName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "confidential_nodes" {
  name               = "%s"
  location           = "us-central1-a"
  release_channel {
    channel = "RAPID"
  }

  node_pool {
    name = "%s"
    initial_node_count = 1
    node_config {
      machine_type = "n2d-standard-2"
    }
  }

  confidential_nodes {
    enabled = false
  }
}
`, clusterName, npName)
}

func testAccContainerCluster_withILBSubSetting(clusterName string, npName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "confidential_nodes" {
  name               = "%s"
  location           = "us-central1-a"
  release_channel {
    channel = "RAPID"
  }

  node_pool {
    name = "%s"
    initial_node_count = 1
    node_config {
      machine_type = "e2-medium"
    }
  }

  enable_l4_ilb_subsetting = true
}
`, clusterName, npName)
}

func testAccContainerCluster_disableILBSubSetting(clusterName string, npName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "confidential_nodes" {
  name               = "%s"
  location           = "us-central1-a"
  release_channel {
    channel = "RAPID"
  }

  node_pool {
    name = "%s"
    initial_node_count = 1
    node_config {
      machine_type = "e2-medium"
    }
  }

  enable_l4_ilb_subsetting = false
}
`, clusterName, npName)
}

func testAccContainerCluster_withNetworkPolicyEnabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
  name                     = "%s"
  location                 = "us-central1-a"
  initial_node_count       = 1
  remove_default_node_pool = true

  network_policy {
    enabled  = true
    provider = "CALICO"
  }

  addons_config {
    network_policy_config {
      disabled = false
    }
  }
}
`, clusterName)
}

func testAccContainerCluster_withReleaseChannelEnabled(clusterName string, channel string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_release_channel" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  release_channel {
    channel = "%s"
  }
}
`, clusterName, channel)
}

func testAccContainerCluster_withReleaseChannelEnabledDefaultVersion(clusterName string, channel string) string {
	return fmt.Sprintf(`

data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_release_channel" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.release_channel_default_version["%s"]
}
`, clusterName, channel)
}

func testAccContainerCluster_removeNetworkPolicy(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
  name                     = "%s"
  location                 = "us-central1-a"
  initial_node_count       = 1
  remove_default_node_pool = true
}
`, clusterName)
}

func testAccContainerCluster_withNetworkPolicyDisabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
  name                     = "%s"
  location                 = "us-central1-a"
  initial_node_count       = 1
  remove_default_node_pool = true

  network_policy {
    enabled = false
  }
}
`, clusterName)
}

func testAccContainerCluster_withNetworkPolicyConfigDisabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
  name                     = "%s"
  location                 = "us-central1-a"
  initial_node_count       = 1
  remove_default_node_pool = true

  network_policy {
    enabled = false
  }

  addons_config {
    network_policy_config {
      disabled = true
    }
  }
}
`, clusterName)
}

func testAccContainerCluster_withAuthenticatorGroupsConfigUpdate(name string, orgDomain string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
	name               = "%s"
	location           = "us-central1-a"
	initial_node_count = 1

	authenticator_groups_config {
		security_group = "gke-security-groups@%s"
	}
}
`, name, orgDomain)
}

func testAccContainerCluster_withAuthenticatorGroupsConfigUpdate2(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
	name               = "%s"
	location           = "us-central1-a"
	initial_node_count = 1

	authenticator_groups_config {
		security_group = ""
	}
}
`, name)
}

func testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName string, cidrs []string, emptyValue string) string {

	cidrBlocks := emptyValue
	if len(cidrs) > 0 {
		var buf bytes.Buffer
		for _, c := range cidrs {
			buf.WriteString(fmt.Sprintf(`
			cidr_blocks {
				cidr_block = "%s"
				display_name = "disp-%s"
			}`, c, c))
		}
		cidrBlocks = buf.String()
	}

	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_authorized_networks" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  master_authorized_networks_config {
    %s
  }
}
`, clusterName, cidrBlocks)
}

func testAccContainerCluster_removeMasterAuthorizedNetworksConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_authorized_networks" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
}
`, clusterName)
}

func testAccContainerCluster_regional(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "regional" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1
}
`, clusterName)
}

func TestAccContainerCluster_withPrivateEndpointSubnetwork(t *testing.T) {
	t.Parallel()

	r := acctest.RandString(t, 10)

	subnet1Name := fmt.Sprintf("tf-test-container-subnetwork1-%s", r)
	subnet1Cidr := "10.0.36.0/24"

	subnet2Name := fmt.Sprintf("tf-test-container-subnetwork2-%s", r)
	subnet2Cidr := "10.9.26.0/24"

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", r)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateEndpointSubnetwork(containerNetName, clusterName, subnet1Name, subnet1Cidr, subnet2Name, subnet2Cidr),
			},
			{
				ResourceName:            "google_container_cluster.with_private_endpoint_subnetwork",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func testAccContainerCluster_withPrivateEndpointSubnetwork(containerNetName, clusterName, s1Name, s1Cidr, s2Name, s2Cidr string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork1" {
  name                     = "%s"
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "%s"
  region                   = "us-central1"
  private_ip_google_access = true
}

resource "google_compute_subnetwork" "container_subnetwork2" {
  name                     = "%s"
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "%s"
  region                   = "us-central1"
  private_ip_google_access = true
}

resource "google_container_cluster" "with_private_endpoint_subnetwork" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork1.name

  private_cluster_config {
    private_endpoint_subnetwork = google_compute_subnetwork.container_subnetwork2.name
  }
}
`, containerNetName, s1Name, s1Cidr, s2Name, s2Cidr, clusterName)
}

func TestAccContainerCluster_withPrivateClusterConfigPrivateEndpointSubnetwork(t *testing.T) {
	t.Parallel()

	r := acctest.RandString(t, 10)

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", r)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateClusterConfigPrivateEndpointSubnetwork(containerNetName, clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_private_endpoint_subnetwork",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func testAccContainerCluster_withPrivateClusterConfigPrivateEndpointSubnetwork(containerNetName, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "with_private_endpoint_subnetwork" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  networking_mode    = "VPC_NATIVE"

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_nodes        = true
    enable_private_endpoint     = true
    private_endpoint_subnetwork = google_compute_subnetwork.container_subnetwork.name
  }
  master_authorized_networks_config {
    gcp_public_cidrs_access_enabled = false
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, containerNetName, clusterName)
}

func TestAccContainerCluster_withEnablePrivateEndpointToggle(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withEnablePrivateEndpoint(clusterName, "true"),
			},
			{
				ResourceName:            "google_container_cluster.with_enable_private_endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_withEnablePrivateEndpoint(clusterName, "false"),
			},
			{
				ResourceName:            "google_container_cluster.with_enable_private_endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
		},
	})
}

func TestAccContainerCluster_failedCreation(t *testing.T) {
	// Test that in a scenario where the cluster fails to create, a subsequent apply will delete the resource.
	// Skip this test for now as we don't have a good way to force cluster creation to fail. https://github.com/hashicorp/terraform-provider-google/issues/13711
	t.Skip()
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	project := acctest.BootstrapProject(t, "tf-fail-cluster-", envvar.GetTestBillingAccountFromEnv(t), []string{"container.googleapis.com"})
	acctest.RemoveContainerServiceAgentRoleFromContainerEngineRobot(t, project)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_failedCreation(clusterName, project.ProjectId),
				ExpectError: regexp.MustCompile("timeout while waiting for state to become 'DONE'"),
			},
			{
				Config:      testAccContainerCluster_failedCreation_update(clusterName, project.ProjectId),
				ExpectError: regexp.MustCompile("Failed to create cluster"),
				Check:       testAccCheckContainerClusterDestroyProducer(t),
			},
		},
	})
}

func testAccContainerCluster_withEnablePrivateEndpoint(clusterName string, flag string) string {

	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_enable_private_endpoint" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1

  master_authorized_networks_config {
    gcp_public_cidrs_access_enabled = false
  }

  private_cluster_config {
    enable_private_endpoint = %s
  }
}
`, clusterName, flag)
}

func testAccContainerCluster_regionalWithNodePool(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "regional" {
  name     = "%s"
  location = "us-central1"

  node_pool {
    name = "%s"
  }
}
`, cluster, nodePool)
}

func testAccContainerCluster_regionalNodeLocations(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_locations" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  node_locations = [
    "us-central1-f",
    "us-central1-c",
  ]
}
`, clusterName)
}

func testAccContainerCluster_regionalUpdateNodeLocations(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_locations" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  node_locations = [
    "us-central1-f",
    "us-central1-b",
  ]
}
`, clusterName)
}

func testAccContainerCluster_withIntraNodeVisibility(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_intranode_visibility" {
  name                        = "%s"
  location                    = "us-central1-a"
  initial_node_count          = 1
  enable_intranode_visibility = true
}
`, clusterName)
}

func testAccContainerCluster_updateIntraNodeVisibility(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_intranode_visibility" {
  name                        = "%s"
  location                    = "us-central1-a"
  initial_node_count          = 1
  enable_intranode_visibility = false
  private_ipv6_google_access  = "PRIVATE_IPV6_GOOGLE_ACCESS_BIDIRECTIONAL"
}
`, clusterName)
}

func testAccContainerCluster_withVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1
}
`, clusterName)
}

func testAccContainerCluster_withLowerVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.valid_master_versions[3]
  initial_node_count = 1
}
`, clusterName)
}

func testAccContainerCluster_withMasterAuthNoCert(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth_no_cert" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3
  master_auth {
    client_certificate_config {
      issue_client_certificate = false
    }
  }
}
`, clusterName)
}

func testAccContainerCluster_updateVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  node_version       = data.google_container_engine_versions.central1a.valid_node_versions[1]
  initial_node_count = 1
}
`, clusterName)
}

func testAccContainerCluster_withNodeConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15
    disk_type       = "pd-ssd"
    local_ssd_count = 1
    oauth_scopes = [
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
    ]
    service_account = "default"
    metadata = {
      foo                      = "bar"
      disable-legacy-endpoints = "true"
    }
    labels = {
      foo = "bar"
    }
    tags             = ["foo", "bar"]
    preemptible      = true
    min_cpu_platform = "Intel Broadwell"

    taint {
      key    = "taint_key"
      value  = "taint_value"
      effect = "PREFER_NO_SCHEDULE"
    }

    taint {
      key    = "taint_key2"
      value  = "taint_value2"
      effect = "NO_EXECUTE"
    }

    // Updatable fields
    image_type = "COS_CONTAINERD"
  }
}
`, clusterName)
}

func testAccContainerCluster_withLoggingVariantInNodeConfig(clusterName, loggingVariant string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging_variant_in_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    logging_variant = "%s"
  }
}
`, clusterName, loggingVariant)
}

func testAccContainerCluster_withLoggingVariantInNodePool(clusterName, nodePoolName, loggingVariant string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging_variant_in_node_pool" {
  name               = "%s"
  location           = "us-central1-f"

  node_pool {
    name               = "%s"
    initial_node_count = 1
    node_config {
      logging_variant = "%s"
    }
  }
}
`, clusterName, nodePoolName, loggingVariant)
}

func testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, loggingVariant string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging_variant_node_pool_default" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_pool_defaults {
    node_config_defaults {
      logging_variant = "%s"
    }
  }
}
`, clusterName, loggingVariant)
}

func testAccContainerCluster_withNodeConfigUpdate(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15
    disk_type       = "pd-ssd"
    local_ssd_count = 1
    oauth_scopes = [
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
    ]
    service_account = "default"
    metadata = {
      foo                      = "bar"
      disable-legacy-endpoints = "true"
    }
    labels = {
      foo = "bar"
    }
    tags             = ["foo", "bar"]
    preemptible      = true
    min_cpu_platform = "Intel Broadwell"

    taint {
      key    = "taint_key"
      value  = "taint_value"
      effect = "PREFER_NO_SCHEDULE"
    }

    taint {
      key    = "taint_key2"
      value  = "taint_value2"
      effect = "NO_EXECUTE"
    }

    // Updatable fields
    image_type = "UBUNTU_CONTAINERD"
  }
}
`, clusterName)
}

func testAccContainerCluster_withNodeConfigScopeAlias(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config_scope_alias" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type = "e2-medium"
    disk_size_gb = 15
    oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
  }
}
`, clusterName)
}

func testAccContainerCluster_withNodeConfigShieldedInstanceConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type    = "e2-medium"
    disk_size_gb    = 15
    disk_type       = "pd-ssd"
    oauth_scopes = [
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
    ]
    service_account = "default"
    metadata = {
      foo                      = "bar"
      disable-legacy-endpoints = "true"
    }
    labels = {
      foo = "bar"
    }
    tags             = ["foo", "bar"]
    preemptible      = true

    // Updatable fields
    image_type = "COS_CONTAINERD"

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }
  }
}
`, clusterName)
}

func testAccContainerCluster_withNodeConfigReservationAffinity(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type    = "e2-medium"
    disk_size_gb    = 15
    disk_type       = "pd-ssd"
    oauth_scopes = [
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
    ]
    service_account = "default"
    metadata = {
      foo                      = "bar"
      disable-legacy-endpoints = "true"
    }
    labels = {
      foo = "bar"
    }
    tags             = ["foo", "bar"]
    preemptible      = true

    // Updatable fields
    image_type = "COS_CONTAINERD"

    reservation_affinity {
      consume_reservation_type = "ANY_RESERVATION"
    }
  }
}
`, clusterName)
}

func testAccContainerCluster_withNodeConfigReservationAffinitySpecific(reservation, clusterName string) string {
	return fmt.Sprintf(`

resource "google_project_service" "compute" {
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "container" {
  service = "container.googleapis.com"
  disable_on_destroy = false
  depends_on = [google_project_service.compute]
}


resource "google_compute_reservation" "gce_reservation" {
  name = "%s"
  zone = "us-central1-f"

  specific_reservation {
    count = 1
    instance_properties {
      machine_type     = "n1-standard-1"
    }
  }

  specific_reservation_required = true
  depends_on = [google_project_service.compute]
}

resource "google_container_cluster" "with_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type    = "n1-standard-1"
    disk_size_gb    = 15
    disk_type       = "pd-ssd"
    oauth_scopes = [
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
    ]
    service_account = "default"
    metadata = {
      foo                      = "bar"
      disable-legacy-endpoints = "true"
    }
    labels = {
      foo = "bar"
    }
    tags             = ["foo", "bar"]

    // Updatable fields
    image_type = "COS_CONTAINERD"

    reservation_affinity {
      consume_reservation_type = "SPECIFIC_RESERVATION"
      key = "compute.googleapis.com/reservation-name"
      values = [
        google_compute_reservation.gce_reservation.name
      ]
    }
  }
  depends_on = [google_project_service.container]
}
`, reservation, clusterName)
}

func testAccContainerCluster_withWorkloadMetadataConfig(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_workload_metadata_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    workload_metadata_config {
      mode = "GCE_METADATA"
    }
  }
}
`, clusterName)
}

func testAccContainerCluster_withBootDiskKmsKey(clusterName, kmsKeyName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_boot_disk_kms_key" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  release_channel {
    channel = "RAPID"
  }
  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]

    image_type = "COS_CONTAINERD"

    boot_disk_kms_key = "%s"
  }
}
`, clusterName, kmsKeyName)
}

func testAccContainerCluster_networkRef(cluster, network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = true
}

resource "google_container_cluster" "with_net_ref_by_url" {
  name               = "%s-url"
  location           = "us-central1-a"
  initial_node_count = 1

  network = google_compute_network.container_network.self_link
}

resource "google_container_cluster" "with_net_ref_by_name" {
  name               = "%s-name"
  location           = "us-central1-a"
  initial_node_count = 1

  network = google_compute_network.container_network.name
}
`, network, cluster, cluster)
}

func testAccContainerCluster_autoprovisioningDefaultsManagement(clusterName string, autoUpgrade, autoRepair bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_autoprovisioning_management" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  cluster_autoscaling {
    enabled = true

	resource_limits {
	  resource_type = "cpu"
	  maximum       = 2
	}

	resource_limits {
	  resource_type = "memory"
	  maximum       = 2048
	}

    auto_provisioning_defaults {
      management {
        auto_upgrade    = %t
        auto_repair     = %t
      }
    }
  }
}
`, clusterName, autoUpgrade, autoRepair)
}

func testAccContainerCluster_backendRef(cluster string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "my-backend-service" {
  name      = "%s-backend"
  port_name = "http"
  protocol  = "HTTP"

  backend {
    group = element(google_container_cluster.primary.node_pool[0].managed_instance_group_urls, 1)
  }

  health_checks = [google_compute_http_health_check.default.self_link]
}

resource "google_compute_http_health_check" "default" {
  name               = "%s-hc"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
`, cluster, cluster, cluster)
}

func testAccContainerCluster_withNodePoolBasic(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"

  node_pool {
    name               = "%s"
    initial_node_count = 2
  }
}
`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolLowerVersion(cluster, nodePool string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"

  min_master_version = data.google_container_engine_versions.central1a.latest_master_version

  node_pool {
    name               = "%s"
    initial_node_count = 2
    version            = data.google_container_engine_versions.central1a.valid_node_versions[2]
  }
}
`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolUpdateVersion(cluster, nodePool string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"

  min_master_version = data.google_container_engine_versions.central1a.latest_master_version

  node_pool {
    name               = "%s"
    initial_node_count = 2
    version            = data.google_container_engine_versions.central1a.valid_node_versions[1]
  }
}
`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolNodeLocations(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]

  node_pool {
    name       = "%s"
    node_count = 2
  }
}
`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolResize(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"

  node_locations = [
    "us-central1-b",
    "us-central1-c",
  ]

  node_pool {
    name       = "%s"
    node_count = 3
  }
}
`, cluster, nodePool)
}

func testAccContainerCluster_autoprovisioning(cluster string, autoprovisioning, withNetworkTag bool) string {
	config := fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1
`, cluster)
	if autoprovisioning {
		config += `
  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }
  }`
	} else {
		config += `
  cluster_autoscaling {
    enabled = false
  }`
	}
	if withNetworkTag {
		config += `
  node_pool_auto_config {
    network_tags {
      tags = ["test-network-tag"]
    }
  }`
	}
	config += `
}`
	return config
}

func testAccContainerCluster_autoprovisioningDefaults(cluster string, monitoringWrite bool) string {
	config := fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1

  logging_service    = "none"
  monitoring_service = "none"

  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }

    auto_provisioning_defaults {
      oauth_scopes = [
        "https://www.googleapis.com/auth/pubsub",
        "https://www.googleapis.com/auth/devstorage.read_only",`,
		cluster)

	if monitoringWrite {
		config += `
        "https://www.googleapis.com/auth/monitoring.write",
`
	}
	config += `
      ]
    }
  }
}`
	return config
}

func testAccContainerCluster_autoprovisioningDefaultsMinCpuPlatform(cluster string, includeMinCpuPlatform bool) string {
	minCpuPlatformCfg := ""
	if includeMinCpuPlatform {
		minCpuPlatformCfg = `min_cpu_platform = "Intel Haswell"`
	}

	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  min_master_version = data.google_container_engine_versions.central1a.latest_master_version

  cluster_autoscaling {
    enabled = true

    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }

    auto_provisioning_defaults {
      %s
    }
  }
}`, cluster, minCpuPlatformCfg)
}

func testAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(clusterName string, maxSurge, maxUnavailable int, strategy string) string {
	blueGreenSettings := ""
	if strategy == "BLUE_GREEN" {
		blueGreenSettings = `
      blue_green_settings {
        node_pool_soak_duration = "3.500s"
        standard_rollout_policy {
        batch_percentage    = 0.5
        batch_soak_duration = "3.500s"
        }
      }
    `
	}

	return fmt.Sprintf(`
    resource "google_container_cluster" "with_autoprovisioning_upgrade_settings" {
      name               = "%s"
      location           = "us-central1-f"
      initial_node_count = 1

      cluster_autoscaling {
        enabled = true

        resource_limits {
          resource_type = "cpu"
          maximum       = 2
        }

        resource_limits {
          resource_type = "memory"
          maximum       = 2048
        }

        auto_provisioning_defaults {
          upgrade_settings {
            max_surge       = %d
            max_unavailable = %d
            strategy        = "%s"
            %s
          }
        }
      }
    }
  `, clusterName, maxSurge, maxUnavailable, strategy, blueGreenSettings)
}

func testAccContainerCluster_autoprovisioningDefaultsUpgradeSettingsWithBlueGreenStrategy(clusterName string, duration, strategy string) string {
	return fmt.Sprintf(`
      resource "google_container_cluster" "with_autoprovisioning_upgrade_settings" {
        name               = "%s"
        location           = "us-central1-f"
        initial_node_count = 1

        cluster_autoscaling {
          enabled = true

          resource_limits {
            resource_type = "cpu"
            maximum       = 2
          }

          resource_limits {
            resource_type = "memory"
            maximum       = 2048
          }

          auto_provisioning_defaults {
            upgrade_settings {
              strategy        = "%s"
              blue_green_settings {
                node_pool_soak_duration = "%s"
                standard_rollout_policy {
                  batch_percentage    = 0.5
                  batch_soak_duration = "%s"
                }
              }
            }
          }
        }
      }
    `, clusterName, strategy, duration, duration)
}

func testAccContainerCluster_autoprovisioningDefaultsDiskSizeGb(cluster string, includeDiskSizeGb bool) string {
	DiskSizeGbCfg := ""
	if includeDiskSizeGb {
		DiskSizeGbCfg = `disk_size = 120`
	}

	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}
resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }
    auto_provisioning_defaults {
      %s
    }
  }
}`, cluster, DiskSizeGbCfg)
}

func testAccContainerCluster_autoprovisioningDefaultsDiskType(cluster string, includeDiskType bool) string {
	DiskTypeCfg := ""
	if includeDiskType {
		DiskTypeCfg = `disk_type = "pd-balanced"`
	}

	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}
resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }
    auto_provisioning_defaults {
      %s
    }
  }
}`, cluster, DiskTypeCfg)
}

func testAccContainerCluster_autoprovisioningDefaultsImageType(cluster string, includeImageType bool) string {
	imageTypeCfg := ""
	if includeImageType {
		imageTypeCfg = `image_type = "COS_CONTAINERD"`
	}

	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}
resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }
    auto_provisioning_defaults {
      %s
    }
  }
}`, cluster, imageTypeCfg)
}

func testAccContainerCluster_autoprovisioningDefaultsBootDiskKmsKey(clusterName, kmsKeyName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "nap_boot_disk_kms_key" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  release_channel {
    channel = "RAPID"
  }
  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }
    auto_provisioning_defaults {
	  boot_disk_kms_key = "%s"
    }
  }
}
`, clusterName, kmsKeyName)
}

func testAccContainerCluster_autoprovisioningDefaultsShieldedInstance(cluster string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}
resource "google_container_cluster" "nap_shielded_instance" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  cluster_autoscaling {
    enabled = true
    resource_limits {
      resource_type = "cpu"
      maximum       = 2
    }
    resource_limits {
      resource_type = "memory"
      maximum       = 2048
    }
    auto_provisioning_defaults {
	  shielded_instance_config {
	    enable_integrity_monitoring = true
	    enable_secure_boot          = true
	  }
    }
  }
}`, cluster)
}

func testAccContainerCluster_withNodePoolAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"

  node_pool {
    name               = "%s"
    initial_node_count = 2
    autoscaling {
      min_node_count = 1
      max_node_count = 3
    }
  }
}
`, cluster, np)
}

func testAccContainerCluster_withNodePoolUpdateAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"

  node_pool {
    name               = "%s"
    initial_node_count = 2
    autoscaling {
      min_node_count = 1
      max_node_count = 5
    }
  }
}
`, cluster, np)
}

func testAccContainerRegionalCluster_withNodePoolCIA(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]

  node_pool {
    name               = "%s"
    initial_node_count = 2
    autoscaling {
      total_min_node_count = 3
      total_max_node_count = 21
      location_policy = "BALANCED"
    }
  }
}
`, cluster, np)
}

func testAccContainerRegionalClusterUpdate_withNodePoolCIA(cluster, np string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]

  node_pool {
    name               = "%s"
    initial_node_count = 2
    autoscaling {
      total_min_node_count = 4
      total_max_node_count = 32
      location_policy = "ANY"
    }
  }
}
`, cluster, np)
}

func testAccContainerRegionalCluster_withNodePoolBasic(cluster, nodePool string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]

  node_pool {
    name               = "%s"
    initial_node_count = 2
  }
}
`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolNamePrefix(cluster, npPrefix string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_name_prefix" {
  name     = "%s"
  location = "us-central1-a"

  node_pool {
    name_prefix = "%s"
    node_count  = 2
  }
}
`, cluster, npPrefix)
}

func testAccContainerCluster_withNodePoolMultiple(cluster, npPrefix string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_multiple" {
  name     = "%s"
  location = "us-central1-a"

  node_pool {
    name       = "%s-one"
    node_count = 2
  }

  node_pool {
    name       = "%s-two"
    node_count = 3
  }
}
`, cluster, npPrefix, npPrefix)
}

func testAccContainerCluster_withNodePoolConflictingNameFields(cluster, npPrefix string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_multiple" {
  name     = "%s"
  location = "us-central1-a"

  node_pool {
    # ERROR: name and name_prefix cannot be both specified
    name        = "%s-notok"
    name_prefix = "%s"
    node_count  = 1
  }
}
`, cluster, npPrefix, npPrefix)
}

func testAccContainerCluster_withNodePoolNodeConfig(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_node_config" {
  name     = "%s"
  location = "us-central1-a"
  node_pool {
    name       = "%s"
    node_count = 2
    node_config {
      machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
      disk_size_gb    = 15
      local_ssd_count = 1
      oauth_scopes = [
        "https://www.googleapis.com/auth/compute",
        "https://www.googleapis.com/auth/devstorage.read_only",
        "https://www.googleapis.com/auth/logging.write",
        "https://www.googleapis.com/auth/monitoring",
      ]
      service_account = "default"
      metadata = {
        foo                      = "bar"
        disable-legacy-endpoints = "true"
      }
      image_type = "COS_CONTAINERD"
      labels = {
        foo = "bar"
      }
      tags = ["foo", "bar"]
    }
  }
}
`, cluster, np)
}

func testAccContainerCluster_withMaintenanceWindow(clusterName string, startTime string) string {
	maintenancePolicy := ""
	if len(startTime) > 0 {
		maintenancePolicy = fmt.Sprintf(`
	maintenance_policy {
		daily_maintenance_window {
			start_time = "%s"
		}
	}`, startTime)
	}

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_window" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  %s
}
`, clusterName, maintenancePolicy)
}

func testAccContainerCluster_withRecurringMaintenanceWindow(clusterName string, startTime, endTime string) string {
	maintenancePolicy := ""
	if len(startTime) > 0 {
		maintenancePolicy = fmt.Sprintf(`
	maintenance_policy {
		recurring_window {
			start_time = "%s"
			end_time = "%s"
			recurrence = "FREQ=DAILY"
		}
	}`, startTime, endTime)
	}

	return fmt.Sprintf(`
resource "google_container_cluster" "with_recurring_maintenance_window" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  %s
}
`, clusterName, maintenancePolicy)

}

func testAccContainerCluster_withExclusion_RecurringMaintenanceWindow(clusterName string, w1startTime, w1endTime, w2startTime, w2endTime string) string {

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_exclusion_window" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  maintenance_policy {
	recurring_window {
		start_time = "%s"
		end_time = "%s"
		recurrence = "FREQ=DAILY"
	}
	maintenance_exclusion {
		exclusion_name = "batch job"
		start_time = "%s"
		end_time = "%s"
	}
	maintenance_exclusion {
		exclusion_name = "holiday data load"
		start_time = "%s"
		end_time = "%s"
	}
 }
}
`, clusterName, w1startTime, w1endTime, w1startTime, w1endTime, w2startTime, w2endTime)
}

func testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(cclusterName string, w1startTime, w1endTime, w2startTime, w2endTime string, scope1, scope2 string) string {

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_exclusion_options" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  maintenance_policy {
	recurring_window {
		start_time = "%s"
		end_time = "%s"
		recurrence = "FREQ=DAILY"
	}
	maintenance_exclusion {
		exclusion_name = "batch job"
		start_time = "%s"
		end_time = "%s"
		exclusion_options {
			scope = "%s"
    	}
	}
	maintenance_exclusion {
		exclusion_name = "holiday data load"
		start_time = "%s"
		end_time = "%s"
		exclusion_options {
			scope = "%s"
    	}
	}
 }
}
`, cclusterName, w1startTime, w1endTime, w1startTime, w1endTime, scope1, w2startTime, w2endTime, scope2)
}

func testAccContainerCluster_NoExclusionOptions_RecurringMaintenanceWindow(cclusterName string, w1startTime, w1endTime, w2startTime, w2endTime string) string {

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_exclusion_options" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  maintenance_policy {
	recurring_window {
		start_time = "%s"
		end_time = "%s"
		recurrence = "FREQ=DAILY"
	}
	maintenance_exclusion {
		exclusion_name = "batch job"
		start_time = "%s"
		end_time = "%s"
	}
	maintenance_exclusion {
		exclusion_name = "holiday data load"
		start_time = "%s"
		end_time = "%s"
	}
 }
}
`, cclusterName, w1startTime, w1endTime, w1startTime, w1endTime, w2startTime, w2endTime)
}

func testAccContainerCluster_updateExclusionOptions_RecurringMaintenanceWindow(cclusterName string, w1startTime, w1endTime, w2startTime, w2endTime string, scope1, scope2 string) string {

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_exclusion_options" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  maintenance_policy {
	recurring_window {
		start_time = "%s"
		end_time = "%s"
		recurrence = "FREQ=DAILY"
	}
	maintenance_exclusion {
		exclusion_name = "batch job"
		start_time = "%s"
		end_time = "%s"
		exclusion_options {
			scope = "%s"
    	}
	}
	maintenance_exclusion {
		exclusion_name = "holiday data load"
		start_time = "%s"
		end_time = "%s"
		exclusion_options {
			scope = "%s"
    	}
	}
 }
}
`, cclusterName, w1startTime, w1endTime, w1startTime, w1endTime, scope1, w2startTime, w2endTime, scope2)
}

func testAccContainerCluster_withExclusion_NoMaintenanceWindow(clusterName string, w1startTime, w1endTime string) string {

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_exclusion_window" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  maintenance_policy {
	recurring_window {
		start_time = "%s"
		end_time = "%s"
		recurrence = "FREQ=DAILY"
	}
 }
}
`, clusterName, w1startTime, w1endTime)
}

func testAccContainerCluster_withExclusion_DailyMaintenanceWindow(clusterName string, w1startTime, w1endTime string) string {

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_exclusion_window" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  maintenance_policy {
	daily_maintenance_window {
		start_time = "03:00"
	}
	maintenance_exclusion {
		exclusion_name = "batch job"
		start_time = "%s"
		end_time = "%s"
	}
 }
}
`, clusterName, w1startTime, w1endTime)
}

func testAccContainerCluster_withIPAllocationPolicy_existingSecondaryRanges(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name    = google_compute_network.container_network.name
  network = google_compute_network.container_network.name
  region  = "us-central1"

  ip_cidr_range = "10.0.0.0/24"

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = "10.1.0.0/16"
  }
  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = "10.2.0.0/20"
  }
}

resource "google_container_cluster" "with_ip_allocation_policy" {
  name     = "%s"
  location = "us-central1-a"

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  networking_mode = "VPC_NATIVE"
  initial_node_count = 1
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withIPAllocationPolicy_specificIPRanges(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name    = google_compute_network.container_network.name
  network = google_compute_network.container_network.name
  region  = "us-central1"

  ip_cidr_range = "10.2.0.0/16"
}

resource "google_container_cluster" "with_ip_allocation_policy" {
  name       = "%s"
  location   = "us-central1-a"
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withIPAllocationPolicy_specificSizes(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name    = google_compute_network.container_network.name
  network = google_compute_network.container_network.name
  region  = "us-central1"

  ip_cidr_range = "10.2.0.0/16"
}

resource "google_container_cluster" "with_ip_allocation_policy" {
  name       = "%s"
  location   = "us-central1-a"
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "/16"
    services_ipv4_cidr_block = "/22"
  }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_stackType_withDualStack(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
    name                    = "%s"
    auto_create_subnetworks = false
}

    resource "google_compute_subnetwork" "container_subnetwork" {
    name    = google_compute_network.container_network.name
    network = google_compute_network.container_network.name
    region  = "us-central1"

    ip_cidr_range = "10.2.0.0/16"
    stack_type = "IPV4_IPV6"
    ipv6_access_type = "EXTERNAL"
}

resource "google_container_cluster" "with_stack_type" {
    name       = "%s"
    location   = "us-central1-a"
    network    = google_compute_network.container_network.name
    subnetwork = google_compute_subnetwork.container_subnetwork.name

    min_master_version = "1.25"
    initial_node_count = 1
    datapath_provider = "ADVANCED_DATAPATH"
    enable_l4_ilb_subsetting = true

    ip_allocation_policy {
        cluster_ipv4_cidr_block  = "10.0.0.0/16"
        services_ipv4_cidr_block = "10.1.0.0/16"
        stack_type = "IPV4_IPV6"
    }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_stackType_withSingleStack(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
    name                    = "%s"
    auto_create_subnetworks = false
}

    resource "google_compute_subnetwork" "container_subnetwork" {
    name    = google_compute_network.container_network.name
    network = google_compute_network.container_network.name
    region  = "us-central1"

    ip_cidr_range = "10.2.0.0/16"
}

resource "google_container_cluster" "with_stack_type" {
    name       = "%s"
    location   = "us-central1-a"
    network    = google_compute_network.container_network.name
    subnetwork = google_compute_subnetwork.container_subnetwork.name

    min_master_version = "1.25"
    initial_node_count = 1
    enable_l4_ilb_subsetting = true

    ip_allocation_policy {
        cluster_ipv4_cidr_block  = "10.0.0.0/16"
        services_ipv4_cidr_block = "10.1.0.0/16"
        stack_type = "IPV4"
    }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_with_PodCIDROverprovisionDisabled(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
    name                    = "%s"
    auto_create_subnetworks = false
}

    resource "google_compute_subnetwork" "container_subnetwork" {
    name    = google_compute_network.container_network.name
    network = google_compute_network.container_network.name
    region  = "us-central1"

    ip_cidr_range = "10.0.0.0/16"
}

resource "google_container_cluster" "with_pco_disabled" {
    name       = "%s"
    location   = "us-central1-a"
    network    = google_compute_network.container_network.name
    subnetwork = google_compute_subnetwork.container_subnetwork.name

    min_master_version = "1.27"
    initial_node_count = 1
    datapath_provider = "ADVANCED_DATAPATH"

    ip_allocation_policy {
        cluster_ipv4_cidr_block  = "10.1.0.0/16"
        services_ipv4_cidr_block = "10.2.0.0/16"
	pod_cidr_overprovision_config {
		disabled = true
	}
    }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withResourceUsageExportConfig(clusterName, datasetId, enableMetering string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                 = "user-project-override"
  user_project_override = true
}
resource "google_bigquery_dataset" "default" {
  dataset_id                 = "%s"
  description                = "gke resource usage dataset tests"
  delete_contents_on_destroy = true
}

resource "google_container_cluster" "with_resource_usage_export_config" {
  provider           = google.user-project-override
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  resource_usage_export_config {
    enable_network_egress_metering = true
    enable_resource_consumption_metering = %s
    bigquery_destination {
      dataset_id = google_bigquery_dataset.default.dataset_id
    }
  }
}
`, datasetId, clusterName, enableMetering)
}

func testAccContainerCluster_withResourceUsageExportConfigNoConfig(clusterName, datasetId string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "default" {
  dataset_id                 = "%s"
  description                = "gke resource usage dataset tests"
  delete_contents_on_destroy = true
}

resource "google_container_cluster" "with_resource_usage_export_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
}
`, datasetId, clusterName)
}

func testAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(containerNetName, clusterName, location string, autopilotEnabled bool) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "with_private_cluster" {
  name               = "%s"
  location           = "%s"
  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
  }

  enable_autopilot = %t

  master_authorized_networks_config {}

  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, containerNetName, clusterName, location, autopilotEnabled)
}

func testAccContainerCluster_withPrivateClusterConfig(containerNetName string, clusterName string, masterGlobalAccessEnabled bool) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "with_private_cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  default_snat_status {
    disabled = true
  }
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
    master_global_access_config {
      enabled = %t
	}
  }
  master_authorized_networks_config {
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, containerNetName, clusterName, masterGlobalAccessEnabled)
}

func testAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(clusterName string, masterGlobalAccessEnabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_private_cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  private_cluster_config {
    enable_private_endpoint = false
    master_global_access_config {
      enabled = %t
	}
  }
}
`, clusterName, masterGlobalAccessEnabled)
}

func testAccContainerCluster_withShieldedNodes(clusterName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_shielded_nodes" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  enable_shielded_nodes = %v
}
`, clusterName, enabled)
}

func testAccContainerCluster_withWorkloadIdentityConfigEnabled(projectID string, clusterName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_container_cluster" "with_workload_identity_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  workload_identity_config {
    workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
  }
  remove_default_node_pool = true

}
`, projectID, clusterName)
}

func testAccContainerCluster_updateWorkloadIdentityConfig(projectID string, clusterName string, enable bool) string {
	workloadIdentityConfig := ""
	if enable {
		workloadIdentityConfig = `
			workload_identity_config {
		  workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
		}`
	} else {
		workloadIdentityConfig = `
			workload_identity_config {
			workload_pool = ""
		}`
	}
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_container_cluster" "with_workload_identity_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  remove_default_node_pool = true
  %s
}
`, projectID, clusterName, workloadIdentityConfig)
}

func testAccContainerCluster_withInitialCIDR(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name          = google_compute_network.container_network.name
  network       = google_compute_network.container_network.name
  ip_cidr_range = "10.128.0.0/9"
}

resource "google_container_cluster" "cidr_error_preempt" {
  name     = "%s"
  location = "us-central1-a"

  networking_mode = "VPC_NATIVE"
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  initial_node_count = 1

  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withCIDROverlap(initConfig, secondCluster string) string {
	return fmt.Sprintf(`
  %s

resource "google_container_cluster" "cidr_error_overlap" {
  name     = "%s"
  location = "us-central1-a"

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
`, initConfig, secondCluster)
}

func testAccContainerCluster_withInvalidLocation(location string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_resource_labels" {
  name               = "invalid-gke-cluster"
  location           = "%s"
  initial_node_count = 1
}
`, location)
}

func testAccContainerCluster_withExternalIpsConfig(projectID string, clusterName string, enabled bool) string {
	return fmt.Sprintf(`
	data "google_project" "project" {
  		project_id = "%s"
	}

	resource "google_container_cluster" "with_external_ips_config" {
		name               = "%s"
		location           = "us-central1-a"
		initial_node_count = 1
		service_external_ips_config {
			enabled = %v
		}
	}`, projectID, clusterName, enabled)
}

func testAccContainerCluster_withMeshCertificatesConfigEnabled(projectID string, clusterName string) string {
	return fmt.Sprintf(`
	data "google_project" "project" {
		project_id = "%s"
	}

	resource "google_container_cluster" "with_mesh_certificates_config" {
	name               = "%s"
	location           = "us-central1-a"
	initial_node_count = 1
	remove_default_node_pool = true
	workload_identity_config {
		workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
	}
	mesh_certificates {
		enable_certificates = true
	}
	}
`, projectID, clusterName)
}

func testAccContainerCluster_updateMeshCertificatesConfig(projectID string, clusterName string, enabled bool) string {
	return fmt.Sprintf(`
	data "google_project" "project" {
  		project_id = "%s"
	}

	resource "google_container_cluster" "with_mesh_certificates_config" {
		name               = "%s"
		location           = "us-central1-a"
		initial_node_count = 1
		remove_default_node_pool = true
		workload_identity_config {
			workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
			}
			mesh_certificates {
			enable_certificates = %v
			}
	}`, projectID, clusterName, enabled)
}

func testAccContainerCluster_updateCostManagementConfig(projectID string, clusterName string, enabled bool) string {
	return fmt.Sprintf(`
	data "google_project" "project" {
  		project_id = "%s"
	}

	resource "google_container_cluster" "with_cost_management_config" {
		name               = "%s"
		location           = "us-central1-a"
		initial_node_count = 1
		cost_management_config {
			enabled = %v
		}
	}`, projectID, clusterName, enabled)
}

func testAccContainerCluster_withDatabaseEncryption(clusterName string, kmsData acctest.BootstrappedKMS) string {
	return fmt.Sprintf(`
data "google_project" "project" {
}

data "google_iam_policy" "test_kms_binding" {
  binding {
    role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

    members = [
      "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com",
    ]
  }
}

resource "google_kms_key_ring_iam_policy" "test_key_ring_iam_policy" {
  key_ring_id = "%[1]s"
  policy_data = data.google_iam_policy.test_kms_binding.policy_data
}

data "google_kms_key_ring_iam_policy" "test_key_ring_iam_policy" {
  key_ring_id = "%[1]s"
}

resource "google_container_cluster" "primary" {
  name               = "%[3]s"
  location           = "us-central1-a"
  initial_node_count = 1

  database_encryption {
    state    = "ENCRYPTED"
    key_name = "%[2]s"
  }
}
`, kmsData.KeyRing.Name, kmsData.CryptoKey.Name, clusterName)
}

func testAccContainerCluster_withDatapathProvider(clusterName, datapathProvider string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  ip_allocation_policy {
  }

  datapath_provider = "%s"

  release_channel {
    channel = "RAPID"
  }
}
`, clusterName, datapathProvider)
}

func testAccContainerCluster_withMasterAuthorizedNetworksDisabled(containerNetName string, clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "with_private_cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_endpoint = false
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withEnableKubernetesAlpha(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  enable_kubernetes_alpha = true

  node_pool {
    name = "%s"
	initial_node_count = 1
	management {
		auto_repair = false
		auto_upgrade = false
	}
  }
}
`, cluster, np)
}

func testAccContainerCluster_withoutEnableKubernetesBetaAPIs(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1
}
`, clusterName)
}

func testAccContainerCluster_withEnableKubernetesBetaAPIs(cluster string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1

  # This feature has been available since GKE 1.27, and currently the only
  # supported Beta API is authentication.k8s.io/v1beta1/selfsubjectreviews.
  # However, in the future, more Beta APIs will be supported, such as the
  # resource.k8s.io group. At the same time, some existing Beta APIs will be
  # deprecated as the feature will be GAed, and the Beta API will be eventually
  # removed. In the case of the SelfSubjectReview API, it is planned to be GAed
  # in Kubernetes as of 1.28. And, the Beta API of SelfSubjectReview will be removed
  # after at least 3 minor version bumps, so it will be removed as of Kubernetes 1.31
  # or later.
  # https://pr.k8s.io/117713
  # https://kubernetes.io/docs/reference/using-api/deprecation-guide/
  #
  # The new Beta APIs will be available since GKE 1.28
  # - admissionregistration.k8s.io/v1beta1/validatingadmissionpolicies
  # - admissionregistration.k8s.io/v1beta1/validatingadmissionpolicybindings
  # https://pr.k8s.io/118644
  #
  # Removing the Beta API from Kubernetes will break the test.
  # TODO: Replace the Beta API with one available on the version of GKE
  # if the test is broken.
  enable_k8s_beta_apis {
    enabled_apis = ["authentication.k8s.io/v1beta1/selfsubjectreviews"]
  }
}
`, cluster)
}

func testAccContainerCluster_withIPv4Error(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
	initial_node_count = 1
	private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = false
    master_ipv4_cidr_block  = "10.42.0.0/28"
  }
}
`, name)
}

func testAccContainerCluster_withAutopilot(projectID string, containerNetName string, clusterName string, location string, enabled bool, withNetworkTag bool, serviceAccount string) string {
	config := ""
	clusterAutoscaling := ""
	if serviceAccount != "" {
		config += fmt.Sprintf(`
resource "google_service_account" "service_account" {
	account_id   = "%[1]s"
	project      = "%[2]s"
	display_name = "Service Account"
}

resource "google_project_iam_binding" "project" {
	project = "%[2]s"
	role    = "roles/container.nodeServiceAccount"
	members = [
		"serviceAccount:%[1]s@%[2]s.iam.gserviceaccount.com",
	]
}`, serviceAccount, projectID)

		clusterAutoscaling = fmt.Sprintf(`
	cluster_autoscaling {
		auto_provisioning_defaults {
			service_account = "%s@%s.iam.gserviceaccount.com"
			oauth_scopes = ["https://www.googleapis.com/auth/cloud-platform"]
		}
	}`, serviceAccount, projectID)
	}

	config += fmt.Sprintf(`

resource "google_compute_network" "container_network" {
	name                    = "%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
	name                     = google_compute_network.container_network.name
	network                  = google_compute_network.container_network.name
	ip_cidr_range            = "10.0.36.0/24"
	region                   = "us-central1"
	private_ip_google_access = true

	secondary_ip_range {
	  range_name    = "pod"
	  ip_cidr_range = "10.0.0.0/19"
	}

	secondary_ip_range {
	  range_name    = "svc"
	  ip_cidr_range = "10.0.32.0/22"
	}
}

data "google_container_engine_versions" "central1a" {
	location = "us-central1-a"
}

resource "google_container_cluster" "with_autopilot" {
	name               = "%s"
	location           = "%s"
	enable_autopilot   = %v
	min_master_version = "latest"
	release_channel {
		channel = "RAPID"
	}
	network       = google_compute_network.container_network.name
	subnetwork    = google_compute_subnetwork.container_subnetwork.name
	ip_allocation_policy {
		cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
		services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
	}
	addons_config {
		horizontal_pod_autoscaling {
			disabled = false
		}
	}
	%s
	vertical_pod_autoscaling {
		enabled = true
	}`, containerNetName, clusterName, location, enabled, clusterAutoscaling)
	if withNetworkTag {
		config += `
	node_pool_auto_config {
		network_tags {
			tags = ["test-network-tag"]
		}
	}`
	}
	config += `
}`
	return config
}

func testAccContainerCluster_withDNSConfig(clusterName string, clusterDns string, clusterDnsDomain string, clusterDnsScope string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_dns_config" {
	name               = "%s"
	location           = "us-central1-f"
	initial_node_count = 1
	dns_config {
		cluster_dns 	   = "%s"
		cluster_dns_domain = "%s"
		cluster_dns_scope  = "%s"
	}
}
`, clusterName, clusterDns, clusterDnsDomain, clusterDnsScope)
}

func testAccContainerCluster_withGatewayApiConfig(clusterName string, gatewayApiChannel string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
	name               = "%s"
	location           = "us-central1-f"
	initial_node_count = 1
	min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
	gateway_api_config {
		channel = "%s"
	}
}
`, clusterName, gatewayApiChannel)
}

func testAccContainerCluster_withLoggingConfigEnabled(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  logging_config {
	  enable_components = [ "SYSTEM_COMPONENTS" ]
  }
  monitoring_config {
      enable_components = [ "SYSTEM_COMPONENTS" ]
  }
}
`, name)
}

func testAccContainerCluster_withLoggingConfigDisabled(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  logging_config {
	  enable_components = []
  }
}
`, name)
}

func testAccContainerCluster_withLoggingConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  logging_config {
	  enable_components = [ "SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER", "SCHEDULER"]
  }
  monitoring_config {
	  enable_components = [ "SYSTEM_COMPONENTS" ]
  }
}
`, name)
}

func testAccContainerCluster_withMonitoringConfigEnabled(name string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
      enable_components = [ "SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER", "SCHEDULER" ]
  }
}
`, name)
}

func testAccContainerCluster_withMonitoringConfigDisabled(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
      enable_components = []
  }
}
`, name)
}

func testAccContainerCluster_withMonitoringConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
         enable_components = [ "SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER" ]
  }
}
`, name)
}

func testAccContainerCluster_withMonitoringConfigPrometheusUpdated(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
         enable_components = [ "SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER", "SCHEDULER" ]
         managed_prometheus {
                 enabled = true
         }
  }
}
`, name)
}

func testAccContainerCluster_withMonitoringConfigPrometheusOnly(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
	     enable_components = []
         managed_prometheus {
                enabled = true
         }
  }
}
`, name)
}

func testAccContainerCluster_withMonitoringConfigPrometheusOnly2(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
         managed_prometheus {
                enabled = true
         }
  }
}
`, name)
}

func testAccContainerCluster_withMonitoringConfigAdvancedDatapathObservabilityConfigEnabled(name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s-nw"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "192.168.1.0/24"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "192.168.64.0/22"
  }
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  datapath_provider = "ADVANCED_DATAPATH"

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }

  monitoring_config {
    enable_components = []
    advanced_datapath_observability_config {
      enable_metrics = true
      relay_mode     = "INTERNAL_VPC_LB"
    }
  }
}
`, name, name)
}

func testAccContainerCluster_withMonitoringConfigAdvancedDatapathObservabilityConfigDisabled(name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s-nw"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "192.168.1.0/24"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "192.168.64.0/22"
  }
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  datapath_provider  = "ADVANCED_DATAPATH"

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }

  monitoring_config {
    enable_components = []
    advanced_datapath_observability_config {
      enable_metrics = false
      relay_mode     = "DISABLED"
    }
  }
}
`, name, name)
}

func testAccContainerCluster_withSoleTenantGroup(name string) string {
	return fmt.Sprintf(`
resource "google_compute_node_template" "soletenant-tmpl" {
  name      = "%s"
  region    = "us-central1"
  node_type = "n1-node-96-624"
}

resource "google_compute_node_group" "group" {
  name        = "%s"
  zone        = "us-central1-f"
  description = "example google_compute_node_group for Terraform Google Provider"

  size          = 1
  node_template = google_compute_node_template.soletenant-tmpl.id
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1
  node_config {
    machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15
    disk_type       = "pd-ssd"
    node_group = google_compute_node_group.group.name
  }
}
`, name, name, name)
}

func testAccContainerCluster_failedCreation(cluster, project string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  project            = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  workload_identity_config {
    workload_pool = "%s.svc.id.goog"
  }

  timeouts {
    create = "40s"
  }
}`, cluster, project, project)
}

func testAccContainerCluster_failedCreation_update(cluster, project string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  project            = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  workload_identity_config {
    workload_pool = "%s.svc.id.goog"
  }
}`, cluster, project, project)
}

func testAccContainerCluster_autopilot_minimal(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name             = "%s"
  location         = "us-central1"
  enable_autopilot = true
}`, name)
}

func testAccContainerCluster_autopilot_net_admin(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name             = "%s"
  location         = "us-central1"
  enable_autopilot = true
  allow_net_admin  = %t
  min_master_version = 1.27
}`, name, enabled)
}

func TestAccContainerCluster_customPlacementPolicy(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	policy := fmt.Sprintf("tf-test-policy-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_customPlacementPolicy(cluster, np, policy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.cluster", "node_pool.0.placement_policy.0.type", "COMPACT"),
					resource.TestCheckResourceAttr("google_container_cluster.cluster", "node_pool.0.placement_policy.0.policy_name", policy),
					resource.TestCheckResourceAttr("google_container_cluster.cluster", "node_pool.0.node_config.0.machine_type", "c2-standard-4"),
				),
			},
			{
				ResourceName:      "google_container_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerCluster_customPlacementPolicy(cluster, np, policyName string) string {
	return fmt.Sprintf(`

resource "google_compute_resource_policy" "policy" {
  name = "%s"
  region = "us-central1"
  group_placement_policy {
    collocation = "COLLOCATED"
  }
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  
  node_pool {
    name               = "%s"
    initial_node_count = 2

    node_config {
      machine_type = "c2-standard-4"
    }

    placement_policy {
      type = "COMPACT"
      policy_name = google_compute_resource_policy.policy.name
    }
  }
}`, policyName, cluster, np)
}
