package google

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccContainerCluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_cluster.primary", "services_ipv4_cidr"),
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

func TestAccContainerCluster_misc(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAddons(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_updateAddons(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuthConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuth(clusterName),
			},
			{
				ResourceName:      "google_container_cluster.with_master_auth",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_updateMasterAuth(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.username", "mr.yoda.adoy.mr"),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.password", "adoy.rm.123456789.mr.yoda"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_master_auth",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_disableMasterAuth(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.username", ""),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.password", ""),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_master_auth",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerCluster_updateMasterAuth(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.username", "mr.yoda.adoy.mr"),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.password", "adoy.rm.123456789.mr.yoda"),
				),
			},
			{
				ResourceName:      "google_container_cluster.with_master_auth",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuthConfig_NoCert(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuthNoCert(),
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
	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAuthenticatorGroupsConfig(clusterName),
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

func TestAccContainerCluster_withMasterAuthorizedNetworksConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-regional-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-regional-%s", acctest.RandString(10))
	npName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateClusterConfig(clusterName),
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(clusterName),
				ExpectError: regexp.MustCompile("master_ipv4_cidr_block must be set if enable_private_nodes == true"),
			},
		},
	})
}

func TestAccContainerCluster_withVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigScopeAlias(),
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

func TestAccContainerCluster_network(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_networkRef(),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_backendRef(),
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

	clusterName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))
	npName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))
	npName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))
	npName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))
	npName := fmt.Sprintf("tf-cluster-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
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
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNamePrefix(),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolMultiple(),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withNodePoolConflictingNameFields(),
				ExpectError: regexp.MustCompile("Cannot specify both name and name_prefix for a node_pool"),
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolNodeConfig(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNodeConfig(),
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
	clusterName := acctest.RandString(10)
	resourceName := "google_container_cluster.with_maintenance_window"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

func TestAccContainerCluster_withIPAllocationPolicy_existingSecondaryRanges(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy_existingSecondaryRanges(cluster),
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

	cluster := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy_specificIPRanges(cluster),
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

	cluster := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy_specificSizes(cluster),
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

func TestAccContainerCluster_errorCleanDanglingCluster(t *testing.T) {
	t.Parallel()

	prefix := acctest.RandString(10)
	clusterName := fmt.Sprintf("cluster-test-%s", prefix)
	clusterNameError := fmt.Sprintf("cluster-test-err-%s", prefix)

	initConfig := testAccContainerCluster_withInitialCIDR(clusterName)
	overlapConfig := testAccContainerCluster_withCIDROverlap(initConfig, clusterNameError)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withInvalidLocation("wonderland"),
				ExpectError: regexp.MustCompile(`Permission denied on 'locations/wonderland' \(or it may not exist\).`),
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuthorizedNetworksDisabled(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksDisabled(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccContainerCluster_masterAuthorizedNetworksDisabled("google_container_cluster.with_private_cluster"),
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

func testAccContainerCluster_masterAuthorizedNetworksDisabled(resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		config := testAccProvider.Meta().(*Config)
		attributes := rs.Primary.Attributes

		cluster, err := config.clientContainer.Projects.Zones.Clusters.Get(
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

func testAccCheckContainerClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_cluster" {
			continue
		}

		attributes := rs.Primary.Attributes
		_, err := config.clientContainer.Projects.Zones.Clusters.Get(
			config.Project, attributes["location"], attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
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

  enable_kubernetes_alpha = true
  enable_legacy_abac      = true

  logging_service    = "logging.googleapis.com"
  monitoring_service = "monitoring.googleapis.com"

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

  enable_kubernetes_alpha = true # Not updatable
  enable_legacy_abac      = false

  logging_service    = "none"
  monitoring_service = "none"

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

func testAccContainerCluster_withAddons(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

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
  }
}
`, clusterName)
}

func testAccContainerCluster_updateAddons(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

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
  }
}
`, clusterName)
}

func testAccContainerCluster_withMasterAuth(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3

  master_auth {
    username = "mr.yoda"
    password = "adoy.rm.123456789"
  }
}
`, clusterName)
}

func testAccContainerCluster_updateMasterAuth(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3

  master_auth {
    username = "mr.yoda.adoy.mr"
    password = "adoy.rm.123456789.mr.yoda"
  }
}
`, clusterName)
}

func testAccContainerCluster_disableMasterAuth(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 3

  master_auth {
    username = ""
    password = ""
  }
}
`, clusterName)
}

func testAccContainerCluster_withMasterAuthNoCert() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth_no_cert" {
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 3

  master_auth {
    username = "mr.yoda"
    password = "adoy.rm.123456789"
    client_certificate_config {
      issue_client_certificate = false
    }
  }
}
`, acctest.RandString(10))
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

func testAccContainerCluster_withAuthenticatorGroupsConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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
    security_group = "gke-security-groups@mydomain.tld"
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, clusterName, clusterName)
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

func testAccContainerCluster_withVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "cluster-test-%s"
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
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.valid_master_versions[2]
  initial_node_count = 1
}
`, clusterName)
}

func testAccContainerCluster_updateVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.valid_master_versions[1]
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
    machine_type    = "n1-standard-1"
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
    image_type = "COS"
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
    machine_type    = "n1-standard-1"
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

func testAccContainerCluster_withNodeConfigScopeAlias() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config_scope_alias" {
  name               = "cluster-test-%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type = "g1-small"
    disk_size_gb = 15
    oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
  }
}
`, acctest.RandString(10))
}

func testAccContainerCluster_withNodeConfigShieldedInstanceConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    machine_type    = "n1-standard-1"
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

func testAccContainerCluster_networkRef() string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
  auto_create_subnetworks = true
}

resource "google_container_cluster" "with_net_ref_by_url" {
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1

  network = google_compute_network.container_network.self_link
}

resource "google_container_cluster" "with_net_ref_by_name" {
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1

  network = google_compute_network.container_network.name
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerCluster_backendRef() string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "my-backend-service" {
  name      = "terraform-test-%s"
  port_name = "http"
  protocol  = "HTTP"

  backend {
    group = element(google_container_cluster.primary.instance_group_urls, 1)
  }

  health_checks = [google_compute_http_health_check.default.self_link]
}

resource "google_compute_http_health_check" "default" {
  name               = "terraform-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_container_cluster" "primary" {
  name               = "terraform-test-%s"
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
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
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

  min_master_version = data.google_container_engine_versions.central1a.valid_master_versions[1]

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

  min_master_version = data.google_container_engine_versions.central1a.valid_master_versions[1]

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

func testAccContainerCluster_withNodePoolNamePrefix() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_name_prefix" {
  name     = "tf-cluster-nodepool-test-%s"
  location = "us-central1-a"

  node_pool {
    name_prefix = "tf-np-test"
    node_count  = 2
  }
}
`, acctest.RandString(10))
}

func testAccContainerCluster_withNodePoolMultiple() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_multiple" {
  name     = "tf-cluster-nodepool-test-%s"
  location = "us-central1-a"

  node_pool {
    name       = "tf-cluster-nodepool-test-%s"
    node_count = 2
  }

  node_pool {
    name       = "tf-cluster-nodepool-test-%s"
    node_count = 3
  }
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerCluster_withNodePoolConflictingNameFields() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_multiple" {
  name     = "tf-cluster-nodepool-test-%s"
  location = "us-central1-a"

  node_pool {
    # ERROR: name and name_prefix cannot be both specified
    name        = "tf-cluster-nodepool-test-%s"
    name_prefix = "tf-cluster-nodepool-test-"
    node_count  = 1
  }
}
`, acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerCluster_withNodePoolNodeConfig() string {
	testId := acctest.RandString(10)
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_node_config" {
  name     = "tf-cluster-nodepool-test-%s"
  location = "us-central1-a"
  node_pool {
    name       = "tf-cluster-nodepool-test-%s"
    node_count = 2
    node_config {
      machine_type    = "n1-standard-1"
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
`, testId, testId)
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
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1
  %s
}
`, clusterName, maintenancePolicy)
}

func testAccContainerCluster_withIPAllocationPolicy_existingSecondaryRanges(cluster string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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

  initial_node_count = 1
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }
}
`, cluster, cluster)
}

func testAccContainerCluster_withIPAllocationPolicy_specificIPRanges(cluster string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
`, cluster, cluster)
}

func testAccContainerCluster_withIPAllocationPolicy_specificSizes(cluster string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "/16"
    services_ipv4_cidr_block = "/22"
  }
}
`, cluster, cluster)
}

func testAccContainerCluster_withPrivateClusterConfigMissingCidrBlock(clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
  }
  master_authorized_networks_config {
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, clusterName, clusterName)
}

func testAccContainerCluster_withPrivateClusterConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
  }
  master_authorized_networks_config {
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`, clusterName, clusterName)
}

func testAccContainerCluster_withInitialCIDR(clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  initial_node_count = 1

  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
`, clusterName, clusterName)
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

func testAccContainerCluster_withMasterAuthorizedNetworksDisabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
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
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1

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
`, clusterName, clusterName)
}
