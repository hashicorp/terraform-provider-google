package google

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("gcp_container_cluster", &resource.Sweeper{
		Name: "gcp_container_cluster",
		F:    testSweepContainerClusters,
	})
}

func testSweepContainerClusters(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Fatalf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	// List clusters for all zones by using "-" as the zone name
	found, err := config.NewContainerClient(config.userAgent).Projects.Zones.Clusters.List(config.Project, "-").Do()
	if err != nil {
		log.Printf("error listing container clusters: %s", err)
		return nil
	}

	if len(found.Clusters) == 0 {
		log.Printf("No container clusters found.")
		return nil
	}

	for _, cluster := range found.Clusters {
		if isSweepableTestResource(cluster.Name) {
			log.Printf("Sweeping Container Cluster: %s", cluster.Name)
			clusterURL := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", config.Project, cluster.Location, cluster.Name)
			_, err := config.NewContainerClient(config.userAgent).Projects.Locations.Clusters.Delete(clusterURL).Do()

			if err != nil {
				log.Printf("Error, failed to delete cluster %s: %s", cluster.Name, err)
				return nil
			}
		}
	}

	return nil
}

func TestAccContainerCluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
				ImportStateId:     fmt.Sprintf("%s/us-central1-a/%s", getTestProjectFromEnv(), clusterName),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	pid := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
			{
				Config: testAccContainerCluster_withInternalLoadBalancer(pid, clusterName),
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

func TestAccContainerCluster_withConfidentialNodes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_withMasterAuthConfig_NoCert(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))
	orgDomain := getTestOrgDomainFromEnv(t)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAuthenticatorGroupsConfig(containerNetName, clusterName, orgDomain),
			},
			{
				ResourceName:      "google_container_cluster.with_authenticator_groups",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withNetworkPolicyEnabled(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
					resource.TestCheckNoResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy"),
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
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
	skipIfVcr(t)
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
					resource.TestCheckNoResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks"),
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

func TestAccContainerCluster_regional(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-regional-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-regional-%s", randString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_withPrivateClusterConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateClusterConfig(containerNetName, clusterName),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(containerNetName, clusterName),
				ExpectError: regexp.MustCompile("master_ipv4_cidr_block must be set if enable_private_nodes is true"),
			},
		},
	})
}

func TestAccContainerCluster_withIntraNodeVisibility(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_withNodeConfigScopeAlias(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_withWorkloadMetadataConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	kms := BootstrapKMSKeyInLocation(t, "us-central1")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withBootDiskKmsKey(getTestProjectFromEnv(), clusterName, kms.CryptoKey.Name),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	network := fmt.Sprintf("tf-test-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroyProducer(t),
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

func TestAccContainerCluster_withNodePoolNamePrefix(t *testing.T) {
	// Randomness
	skipIfVcr(t)
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	npNamePrefix := "tf-test-np-"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	npNamePrefix := "tf-test-np-"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	npPrefix := "tf-test-np"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	np := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_window"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	resourceName := "google_container_cluster.with_recurring_maintenance_window"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_window"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_deleteExclusionWindow(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_window"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_nodeAutoprovisioning(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioning(clusterName, true),
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
				Config: testAccContainerCluster_autoprovisioning(clusterName, false),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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
		},
	})
}

func TestAccContainerCluster_withShieldedNodes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutopilot(containerNetName, clusterName, "us-central1", true),
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

	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withAutopilot(containerNetName, clusterName, "us-central1-a", true),
				ExpectError: regexp.MustCompile(`Autopilot clusters must be regional clusters.`),
			},
		},
	})
}

func TestAccContainerCluster_withWorkloadIdentityConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	pid := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_withSoleTenantGroup(t *testing.T) {
	t.Parallel()

	resourceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_nodeAutoprovisioningDefaultsImageType(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	includeImageType := true

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_errorCleanDanglingCluster(t *testing.T) {
	t.Parallel()

	prefix := randString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", prefix)
	clusterNameError := fmt.Sprintf("tf-test-cluster-err-%s", prefix)
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))

	initConfig := testAccContainerCluster_withInitialCIDR(containerNetName, clusterName)
	overlapConfig := testAccContainerCluster_withCIDROverlap(initConfig, clusterNameError)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withInvalidLocation("wonderland"),
				ExpectError: regexp.MustCompile(`Permission denied on 'locations/wonderland' \(or it may not exist\).`),
			},
		},
	})
}

func TestAccContainerCluster_withDatabaseEncryption(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	// Use the bootstrapped KMS key so we can avoid creating keys needlessly
	// as they will pile up in the project because they can not be completely
	// deleted.  Also, we need to create the key in the same location as the
	// cluster as GKE does not support the "global" location for KMS keys.
	// See https://cloud.google.com/kubernetes-engine/docs/how-to/encrypting-secrets#creating_a_key
	kmsData := BootstrapKMSKeyInLocation(t, "us-central1")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDatabaseEncryption(clusterName, kmsData),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	suffix := randString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", suffix)
	datesetId := fmt.Sprintf("tf_test_cluster_resource_usage_%s", suffix)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	npName := fmt.Sprintf("tf-test-np-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func TestAccContainerCluster_withIPv4Error(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

	clusterName := fmt.Sprintf("tf-test-cluster-%s", randString(t, 10))
	domainName := fmt.Sprintf("tf-test-domain-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
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

func testAccContainerCluster_masterAuthorizedNetworksDisabled(t *testing.T, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		config := googleProviderConfig(t)
		attributes := rs.Primary.Attributes

		cluster, err := config.NewContainerClient(config.userAgent).Projects.Zones.Clusters.Get(
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
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_container_cluster" {
				continue
			}

			attributes := rs.Primary.Attributes
			_, err := config.NewContainerClient(config.userAgent).Projects.Zones.Clusters.Get(
				config.Project, attributes["location"], attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("Cluster still exists")
			}
		}

		return nil
	}
}

func getResourceAttributes(n string, s *terraform.State) (map[string]string, error) {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", n)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	return rs.Primary.Attributes, nil
}

func checkMatch(attributes map[string]string, attr string, gcp interface{}) string {
	if gcpList, ok := gcp.([]string); ok {
		return checkListMatch(attributes, attr, gcpList)
	}
	if gcpMap, ok := gcp.(map[string]string); ok {
		return checkMapMatch(attributes, attr, gcpMap)
	}
	if gcpBool, ok := gcp.(bool); ok {
		return checkBoolMatch(attributes, attr, gcpBool)
	}

	tf := attributes[attr]
	if tf != gcp {
		return matchError(attr, tf, gcp)
	}
	return ""
}

func checkListMatch(attributes map[string]string, attr string, gcpList []string) string {
	// A bunch of the TestAccDataprocJob_* tests fail without this. It's likely an inaccuracy that happens when shimming the terraform-json
	// representation of state back to the old framework's representation of state. So, in the past we would get x.# = 0 whereas now we get x.# = ''.
	// It's likely not intentional, however, shouldn't be a big problem - but if we notice it is the sdk team can address it.
	if attributes[attr+".#"] == "" {
		attributes[attr+".#"] = "0"
	}

	num, err := strconv.Atoi(attributes[attr+".#"])
	if err != nil {
		return fmt.Sprintf("Error in number conversion for attribute %s: %s", attr, err)
	}
	if num != len(gcpList) {
		return fmt.Sprintf("Cluster has mismatched %s size.\nTF Size: %d\nGCP Size: %d", attr, num, len(gcpList))
	}

	for i, gcp := range gcpList {
		if tf := attributes[fmt.Sprintf("%s.%d", attr, i)]; tf != gcp {
			return matchError(fmt.Sprintf("%s[%d]", attr, i), tf, gcp)
		}
	}

	return ""
}

func checkMapMatch(attributes map[string]string, attr string, gcpMap map[string]string) string {
	// A bunch of the TestAccDataprocJob_* tests fail without this. It's likely an inaccuracy that happens when shimming the terraform-json
	// representation of state back to the old framework's representation of state. So, in the past we would get x.# = 0 whereas now we get x.# = ''.
	// It's likely not intentional, however, shouldn't be a big problem - but if we notice it is the sdk team can address it.
	if attributes[attr+".%"] == "" {
		attributes[attr+".%"] = "0"
	}

	num, err := strconv.Atoi(attributes[attr+".%"])
	if err != nil {
		return fmt.Sprintf("Error in number conversion for attribute %s: %s", attr, err)
	}
	if num != len(gcpMap) {
		return fmt.Sprintf("Cluster has mismatched %s size.\nTF Size: %d\nGCP Size: %d", attr, num, len(gcpMap))
	}

	for k, gcp := range gcpMap {
		if tf := attributes[fmt.Sprintf("%s.%s", attr, k)]; tf != gcp {
			return matchError(fmt.Sprintf("%s[%s]", attr, k), tf, gcp)
		}
	}

	return ""
}

func checkBoolMatch(attributes map[string]string, attr string, gcpBool bool) string {
	// Handle the case where an unset value defaults to false
	var tf bool
	var err error
	if attributes[attr] == "" {
		tf = false
	} else {
		tf, err = strconv.ParseBool(attributes[attr])
		if err != nil {
			return fmt.Sprintf("Error converting attribute %s to boolean: value is %s", attr, attributes[attr])
		}
	}

	if tf != gcpBool {
		return matchError(attr, tf, gcpBool)
	}

	return ""
}

func matchError(attr, tf interface{}, gcp interface{}) string {
	return fmt.Sprintf("Cluster has mismatched %s.\nTF State: %+v\nGCP State: %+v", attr, tf, gcp)
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
      disabled = false
    }
  }
}
`, projectID, clusterName)
}

func testAccContainerCluster_withInternalLoadBalancer(projectID string, clusterName string) string {
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
    cloudrun_config {
	  disabled = false
	  load_balancer_type = "LOAD_BALANCER_TYPE_INTERNAL"
    }
  }
}
`, projectID, clusterName)
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

func testAccContainerCluster_withAuthenticatorGroupsConfig(containerNetName string, clusterName string, orgDomain string) string {
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

resource "google_container_cluster" "with_authenticator_groups" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  network            = google_compute_network.container_network.name
  subnetwork         = google_compute_subnetwork.container_subnetwork.name

  authenticator_groups_config {
    security_group = "gke-security-groups-test@%s"
  }

  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, containerNetName, clusterName, orgDomain)
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
  min_master_version = data.google_container_engine_versions.central1a.valid_master_versions[2]
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
    image_type = "cos"
  }
}
`, clusterName)
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
    image_type = "UBUNTU"
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
    image_type = "COS"

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }
  }
}
`, clusterName)
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

func testAccContainerCluster_withBootDiskKmsKey(project, clusterName, kmsKeyName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_project_iam_member" "kms-project-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com"
}

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
`, project, clusterName, kmsKeyName)
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

func testAccContainerCluster_autoprovisioning(cluster string, autoprovisioning bool) string {
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
      image_type = "COS"
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

func testAccContainerCluster_withResourceUsageExportConfig(clusterName, datasetId, enableMetering string) string {
	return fmt.Sprintf(`
provider "google" {
  user_project_override = true
}
resource "google_bigquery_dataset" "default" {
  dataset_id                 = "%s"
  description                = "gke resource usage dataset tests"
  delete_contents_on_destroy = true
}

resource "google_container_cluster" "with_resource_usage_export_config" {
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

func testAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(containerNetName string, clusterName string) string {
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
    enable_private_endpoint = true
    enable_private_nodes    = true
  }

  master_authorized_networks_config {}

  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withPrivateClusterConfig(containerNetName string, clusterName string) string {
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
      enabled = true
	}
  }
  master_authorized_networks_config {
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, containerNetName, clusterName)
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

func testAccContainerCluster_withDatabaseEncryption(clusterName string, kmsData bootstrappedKMS) string {
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

func testAccContainerCluster_withAutopilot(containerNetName string, clusterName string, location string, enabled bool) string {
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
	vertical_pod_autoscaling {
		enabled = true
	}
}
`, containerNetName, clusterName, location, enabled)
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

func testAccContainerCluster_withLoggingConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  logging_config {
	  enable_components = [ "SYSTEM_COMPONENTS", "WORKLOADS" ]
  }
  monitoring_config {
	  enable_components = [ "SYSTEM_COMPONENTS" ]
  }
}
`, name)
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
