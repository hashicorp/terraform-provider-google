package google

import (
	"bytes"
	"fmt"
	"testing"

	"strconv"

	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
			},
			{
				ResourceName:        "google_container_cluster.primary",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				ResourceName:        "google_container_cluster.primary",
				ImportStateIdPrefix: fmt.Sprintf("%s/us-central1-a/", getTestProjectFromEnv()),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withTimeout(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withTimeout(),
			},
			{
				ResourceName:        "google_container_cluster.primary",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.http_load_balancing.0.disabled", "true"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.kubernetes_dashboard.0.disabled", "true"),
				),
			},
			{
				ResourceName:        "google_container_cluster.primary",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateAddons(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.horizontal_pod_autoscaling.0.disabled", "true"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.http_load_balancing.0.disabled", "false"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.kubernetes_dashboard.0.disabled", "true"),
				),
			},
			{
				ResourceName:        "google_container_cluster.primary",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.with_master_auth",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateMasterAuth(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.username", "mr.yoda.adoy.mr"),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.password", "adoy.rm.123456789.mr.yoda"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_master_auth",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_disableMasterAuth(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.username", ""),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.password", ""),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_master_auth",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateMasterAuth(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.username", "mr.yoda.adoy.mr"),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth", "master_auth.0.password", "adoy.rm.123456789.mr.yoda"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_master_auth",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.with_master_auth_no_cert",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ImportStateIdPrefix:     "us-central1-a/",
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
				ImportStateIdPrefix:     "us-central1-a/",
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
				ImportStateIdPrefix:     "us-central1-a/",
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
				ImportStateIdPrefix:     "us-central1-a/",
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
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{}, "cidr_blocks = []"),
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
				ResourceName:        "google_container_cluster.with_master_authorized_networks",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: "us-central1-a/",
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{"10.0.0.0/8", "8.8.8.8/32"}, ""),
			},
			{
				ResourceName:        "google_container_cluster.with_master_authorized_networks",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: "us-central1-a/",
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{}, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_master_authorized_networks",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.regional",
				ImportStateIdPrefix: "us-central1/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.regional",
				ImportStateIdPrefix: "us-central1/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withAdditionalZones(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAdditionalZones(clusterName),
			},
			{
				ResourceName:        "google_container_cluster.with_additional_zones",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateAdditionalZones(clusterName),
			},
			{
				ResourceName:        "google_container_cluster.with_additional_zones",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_regionalWithAdditionalZones(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_regionalAdditionalZones(clusterName),
			},
			{
				ResourceName:        "google_container_cluster.with_additional_zones",
				ImportStateIdPrefix: "us-central1/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_regionalUpdateAdditionalZones(clusterName),
			},
			{
				ResourceName:        "google_container_cluster.with_additional_zones",
				ImportStateIdPrefix: "us-central1/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withKubernetesAlpha(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withKubernetesAlpha(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_kubernetes_alpha", "enable_kubernetes_alpha", "true"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_kubernetes_alpha",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withTpu(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withTpu(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_tpu", "enable_tpu", "true"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_tpu",
				ImportStateIdPrefix: "us-central1-b/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withPrivateCluster(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateCluster(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_private_cluster", "private_cluster", "true"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_private_cluster",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withLegacyAbac(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLegacyAbac(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_legacy_abac", "enable_legacy_abac", "true"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_legacy_abac",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateLegacyAbac(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_legacy_abac", "enable_legacy_abac", "false"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_legacy_abac",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

/*
	Since GKE disables legacy ABAC by default in Kubernetes version 1.8+, and the default Kubernetes
	version for GKE is also 1.8+, this test will ensure that legacy ABAC is disabled by default to be
	more consistent with default settings in the Cloud Console
*/
func TestAccContainerCluster_withDefaultLegacyAbac(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_defaultLegacyAbac(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.default_legacy_abac", "enable_legacy_abac", "false"),
				),
			},
			{
				ResourceName:        "google_container_cluster.default_legacy_abac",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ImportStateIdPrefix:     "us-central1-a/",
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
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
			},
			{
				Config: testAccContainerCluster_updateVersion(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_version",
				ImportStateIdPrefix:     "us-central1-a/",
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
				ResourceName:        "google_container_cluster.with_node_config",
				ImportStateIdPrefix: "us-central1-f/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withNodeConfigUpdate(clusterName),
			},
			{
				ResourceName:        "google_container_cluster.with_node_config",
				ImportStateIdPrefix: "us-central1-f/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.with_node_config_scope_alias",
				ImportStateIdPrefix: "us-central1-f/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigTaints(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigTaints(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config", "node_config.0.taint.#", "2"),
				),
			},
			// Don't include an import step because beta features can't yet be imported.
			// Once taints are in GA, consider merging this test with the _withNodeConfig test.
		},
	})
}

func TestAccContainerCluster_withWorkloadMetadataConfig(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withWorkloadMetadataConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_workload_metadata_config",
						"node_config.0.workload_metadata_config.0.node_metadata", "SECURE"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_metadata_config",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version"},
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
				ResourceName:        "google_container_cluster.with_net_ref_by_url",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				ResourceName:        "google_container_cluster.with_net_ref_by_name",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.primary",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withLogging(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLogging(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_logging", "logging_service", "logging.googleapis.com"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_logging",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateLogging(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_logging", "logging_service", "none"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_logging",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withMonitoring(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMonitoring(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_monitoring", "monitoring_service", "monitoring.googleapis.com"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_monitoring",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_updateMonitoring(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_monitoring", "monitoring_service", "none"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_monitoring",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.with_node_pool",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				Config: testAccContainerCluster_withNodePoolAdditionalZones(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "2"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_node_pool",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withNodePoolResize(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "3"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_node_pool",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
			resource.TestStep{
				Config: testAccContainerCluster_withNodePoolAutoscaling(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "3"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_node_pool",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			resource.TestStep{
				Config: testAccContainerCluster_withNodePoolUpdateAutoscaling(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "5"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_node_pool",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			resource.TestStep{
				Config: testAccContainerCluster_withNodePoolBasic(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_node_pool",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ImportStateIdPrefix:     "us-central1-a/",
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
				ResourceName:        "google_container_cluster.with_node_pool_multiple",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.with_node_pool_node_config",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withDefaultNodePoolRemoved(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDefaultNodePoolRemoved(),
				Check: resource.TestCheckResourceAttr(
					"google_container_cluster.with_default_node_pool_removed", "node_pool.#", "0"),
			},
			{
				ResourceName:            "google_container_cluster.with_default_node_pool_removed",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool"},
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
				ResourceName:        resourceName,
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withMaintenanceWindow(clusterName, ""),
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
				ResourceName:        "google_container_cluster.with_ip_allocation_policy",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.with_ip_allocation_policy",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
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
				ResourceName:        "google_container_cluster.with_ip_allocation_policy",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withIPAllocationPolicy_createSubnetwork(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy_createSubnetwork(cluster),
			},
			{
				ResourceName:        "google_container_cluster.with_ip_allocation_policy",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withPodSecurityPolicy(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPodSecurityPolicy(clusterName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_pod_security_policy",
						"pod_security_policy_config.0.enabled", "true"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_pod_security_policy",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withPodSecurityPolicy(clusterName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_pod_security_policy",
						"pod_security_policy_config.0.enabled", "false"),
				),
			},
			{
				ResourceName:        "google_container_cluster.with_pod_security_policy",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_sharedVpc(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	projectName := fmt.Sprintf("tf-xpntest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_sharedVpc(org, billingId, projectName, clusterName),
			},
			{
				ResourceName:        "google_container_cluster.shared_vpc_cluster",
				ImportStateIdPrefix: fmt.Sprintf("%s-service/us-central1-a/", projectName),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withResourceLabels(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withResourceLabels(clusterName),
			},
			{
				ResourceName:        "google_container_cluster.with_resource_labels",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withResourceLabelsUpdate(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withoutResourceLabels(clusterName),
			},
			{
				Config: testAccContainerCluster_withResourceLabels(clusterName),
			},
			{
				ResourceName:        "google_container_cluster.with_resource_labels",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccContainerCluster_withBinaryAuthorization(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withBinaryAuthorization(clusterName, true),
			},
			{
				ResourceName:        "google_container_cluster.with_binary_authorization",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
			{
				Config: testAccContainerCluster_withBinaryAuthorization(clusterName, false),
			},
			{
				ResourceName:        "google_container_cluster.with_binary_authorization",
				ImportStateIdPrefix: "us-central1-a/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func testAccCheckContainerClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_cluster" {
			continue
		}

		attributes := rs.Primary.Attributes
		_, err := config.clientContainer.Projects.Zones.Clusters.Get(
			config.Project, attributes["zone"], attributes["name"]).Do()
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
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3
}`, name)
}

func testAccContainerCluster_withTimeout() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 3

	timeouts {
		create = "30m"
		delete = "30m"
		update = "30m"
	}
}`, acctest.RandString(10))
}

func testAccContainerCluster_withAddons(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	addons_config {
		http_load_balancing { disabled = true }
		kubernetes_dashboard { disabled = true }
	}
}`, clusterName)
}

func testAccContainerCluster_updateAddons(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	addons_config {
		http_load_balancing { disabled = false }
		kubernetes_dashboard { disabled = true }
		horizontal_pod_autoscaling { disabled = true }
	}
}`, clusterName)
}

func testAccContainerCluster_withMasterAuth(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm.123456789"
	}
}`, clusterName)
}

func testAccContainerCluster_updateMasterAuth(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda.adoy.mr"
		password = "adoy.rm.123456789.mr.yoda"
	}
}`, clusterName)
}

func testAccContainerCluster_disableMasterAuth(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = ""
		password = ""
	}
}`, clusterName)
}

func testAccContainerCluster_updateMasterAuthNoCert() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm.123456789"
		client_certificate_config {
			issue_client_certificate = false
		}
	}
}`, acctest.RandString(10))
}

func testAccContainerCluster_withMasterAuthNoCert() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth_no_cert" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm.123456789"
		client_certificate_config {
			issue_client_certificate = false
		}
	}
}`, acctest.RandString(10))
}

func testAccContainerCluster_withNetworkPolicyEnabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
	remove_default_node_pool = true

	network_policy {
		enabled = true
		provider = "CALICO"
	}
	addons_config {
		network_policy_config {
			disabled = false
		}
	}
}`, clusterName)
}

func testAccContainerCluster_removeNetworkPolicy(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
	remove_default_node_pool = true
}`, clusterName)
}

func testAccContainerCluster_withNetworkPolicyDisabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
	remove_default_node_pool = true

	network_policy = {}
}`, clusterName)
}

func testAccContainerCluster_withNetworkPolicyConfigDisabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
	remove_default_node_pool = true

	network_policy = {}
	addons_config {
		network_policy_config {
			disabled = true
		}
	}
}`, clusterName)
}

func testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName string, cidrs []string, emptyValue string) string {

	cidrBlocks := emptyValue
	if len(cidrs) > 0 {
		var buf bytes.Buffer
		buf.WriteString("cidr_blocks = [")
		for _, c := range cidrs {
			buf.WriteString(fmt.Sprintf(`
			{
				cidr_block = "%s"
				display_name = "disp-%s"
			},`, c, c))
		}
		buf.WriteString("]")
		cidrBlocks = buf.String()
	}

	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_authorized_networks" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	master_authorized_networks_config {
		%s
	}
}`, clusterName, cidrBlocks)
}

func testAccContainerCluster_regional(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "regional" {
	name = "%s"
	region = "us-central1"
	initial_node_count = 1
}`, clusterName)
}

func testAccContainerCluster_regionalWithNodePool(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "regional" {
	name   = "%s"
	region = "us-central1"

	node_pool {
		name = "%s"
	}
}`, cluster, nodePool)
}

func testAccContainerCluster_withAdditionalZones(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_additional_zones" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	additional_zones = [
		"us-central1-b",
		"us-central1-c"
	]
}`, clusterName)
}

func testAccContainerCluster_updateAdditionalZones(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_additional_zones" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	additional_zones = [
		"us-central1-f",
		"us-central1-c",
	]
}`, clusterName)
}

func testAccContainerCluster_regionalAdditionalZones(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_additional_zones" {
	name = "%s"
	region = "us-central1"
	initial_node_count = 1

	additional_zones = [
		"us-central1-f",
		"us-central1-c",
	]
}`, clusterName)
}

func testAccContainerCluster_regionalUpdateAdditionalZones(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_additional_zones" {
	name = "%s"
	region = "us-central1"
	initial_node_count = 1

	additional_zones = [
		"us-central1-f",
		"us-central1-b",
	]
}`, clusterName)
}

func testAccContainerCluster_withKubernetesAlpha(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_kubernetes_alpha" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	enable_kubernetes_alpha = true
}`, clusterName)
}

func testAccContainerCluster_withTpu(clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
	name                     = "${google_compute_network.container_network.name}"
	network                  = "${google_compute_network.container_network.name}"
	ip_cidr_range            = "10.0.35.0/24"
	region                   = "us-central1"

	secondary_ip_range {
		range_name    = "pod"
		ip_cidr_range = "10.1.0.0/19"
	}

	secondary_ip_range {
		range_name    = "svc"
		ip_cidr_range = "10.2.0.0/22"
	}
}

resource "google_container_cluster" "with_tpu" {
	name = "cluster-test-%s"
	zone = "us-central1-b"
	initial_node_count = 1

	enable_tpu = true

	network = "${google_compute_network.container_network.name}"
	subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

	master_ipv4_cidr_block = "10.42.0.0/28"
	master_authorized_networks_config { cidr_blocks = [] }
	ip_allocation_policy {
		cluster_secondary_range_name  = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.0.range_name}"
		services_secondary_range_name = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.1.range_name}"
	}
}`, clusterName, clusterName)
}

func testAccContainerCluster_defaultLegacyAbac(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "default_legacy_abac" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1
}`, clusterName)
}

func testAccContainerCluster_withLegacyAbac(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_legacy_abac" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	enable_legacy_abac = true
}`, clusterName)
}

func testAccContainerCluster_updateLegacyAbac(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_legacy_abac" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	enable_legacy_abac = false
}`, clusterName)
}

func testAccContainerCluster_withVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
	zone = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	min_master_version = "${data.google_container_engine_versions.central1a.latest_master_version}"
	initial_node_count = 1
}`, clusterName)
}

func testAccContainerCluster_withLowerVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
	zone = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	min_master_version = "${data.google_container_engine_versions.central1a.valid_master_versions.2}"
	initial_node_count = 1
}`, clusterName)
}

func testAccContainerCluster_updateVersion(clusterName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
	zone = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	min_master_version = "${data.google_container_engine_versions.central1a.valid_master_versions.1}"
	node_version = "${data.google_container_engine_versions.central1a.valid_node_versions.1}"
	initial_node_count = 1
}`, clusterName)
}

func testAccContainerCluster_withNodeConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
	name = "%s"
	zone = "us-central1-f"
	initial_node_count = 1

	node_config {
		machine_type = "n1-standard-1"
		disk_size_gb = 15
		disk_type = "pd-ssd"
		local_ssd_count = 1
		oauth_scopes = [
			"https://www.googleapis.com/auth/monitoring",
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write"
		]
		service_account = "default"
		metadata {
			foo = "bar"
		}
		labels {
			foo = "bar"
		}
		tags = ["foo", "bar"]
		preemptible = true
		min_cpu_platform = "Intel Broadwell"

		// Updatable fields
		image_type = "COS"
	}
}`, clusterName)
}

func testAccContainerCluster_withNodeConfigUpdate(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
	name = "%s"
	zone = "us-central1-f"
	initial_node_count = 1

	node_config {
		machine_type = "n1-standard-1"
		disk_size_gb = 15
		disk_type = "pd-ssd"
		local_ssd_count = 1
		oauth_scopes = [
			"https://www.googleapis.com/auth/monitoring",
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write"
		]
		service_account = "default"
		metadata {
			foo = "bar"
		}
		labels {
			foo = "bar"
		}
		tags = ["foo", "bar"]
		preemptible = true
		min_cpu_platform = "Intel Broadwell"

		// Updatable fields
		image_type = "UBUNTU"
	}
}`, clusterName)
}

func testAccContainerCluster_withNodeConfigScopeAlias() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config_scope_alias" {
	name = "cluster-test-%s"
	zone = "us-central1-f"
	initial_node_count = 1

	node_config {
		machine_type = "g1-small"
		disk_size_gb = 15
		oauth_scopes = [ "compute-rw", "storage-ro", "logging-write", "monitoring" ]
	}
}`, acctest.RandString(10))
}

func testAccContainerCluster_withNodeConfigTaints() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
	name = "cluster-test-%s"
	zone = "us-central1-f"
	initial_node_count = 1

	node_config {
		taint {
			key = "taint_key"
			value = "taint_value"
			effect = "PREFER_NO_SCHEDULE"
		}
		taint {
			key = "taint_key2"
			value = "taint_value2"
			effect = "NO_EXECUTE"
		}
	}
}`, acctest.RandString(10))
}

func testAccContainerCluster_withWorkloadMetadataConfig() string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  zone = "us-central1-a"
}

resource "google_container_cluster" "with_workload_metadata_config" {
  name               = "cluster-test-%s"
  zone               = "us-central1-a"
  initial_node_count = 1
  min_master_version = "${data.google_container_engine_versions.central1a.latest_master_version}"
  node_version       = "${data.google_container_engine_versions.central1a.latest_node_version}"

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring"
    ]

    workload_metadata_config {
      node_metadata = "SECURE"
    }
  }
}
`, acctest.RandString(10))
}

func testAccContainerCluster_networkRef() string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = true
}

resource "google_container_cluster" "with_net_ref_by_url" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	network = "${google_compute_network.container_network.self_link}"
}

resource "google_container_cluster" "with_net_ref_by_name" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	network = "${google_compute_network.container_network.name}"
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerCluster_backendRef() string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "my-backend-service" {
  name      = "terraform-test-%s"
  port_name = "http"
  protocol  = "HTTP"

  backend {
    group = "${element(google_container_cluster.primary.instance_group_urls, 1)}"
  }

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_http_health_check" "default" {
  name               = "terraform-test-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_container_cluster" "primary" {
  name               = "terraform-test-%s"
  zone               = "us-central1-a"
  initial_node_count = 3

  additional_zones = [
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

func testAccContainerCluster_withLogging(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging" {
	name               = "cluster-test-%s"
	zone               = "us-central1-a"
	initial_node_count = 1

	logging_service = "logging.googleapis.com"
}`, clusterName)
}

func testAccContainerCluster_updateLogging(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging" {
	name               = "cluster-test-%s"
	zone               = "us-central1-a"
	initial_node_count = 1

	logging_service = "none"
}`, clusterName)
}

func testAccContainerCluster_withMonitoring(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_monitoring" {
	name               = "cluster-test-%s"
	zone               = "us-central1-a"
	initial_node_count = 1

	monitoring_service = "monitoring.googleapis.com"
}`, clusterName)
}

func testAccContainerCluster_updateMonitoring(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_monitoring" {
	name               = "cluster-test-%s"
	zone               = "us-central1-a"
	initial_node_count = 1

	monitoring_service = "none"
}`, clusterName)
}

func testAccContainerCluster_withNodePoolBasic(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
	name = "%s"
	zone = "us-central1-a"

	node_pool {
		name               = "%s"
		initial_node_count = 2
	}
}`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolAdditionalZones(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
	name = "%s"
	zone = "us-central1-a"

	additional_zones = [
		"us-central1-b",
		"us-central1-c"
	]

	node_pool {
		name       = "%s"
		node_count = 2
	}
}`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolResize(cluster, nodePool string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
	name = "%s"
	zone = "us-central1-a"

	additional_zones = [
		"us-central1-b",
		"us-central1-c"
	]

	node_pool {
		name       = "%s"
		node_count = 3
	}
}`, cluster, nodePool)
}

func testAccContainerCluster_withNodePoolAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
	name = "%s"
	zone = "us-central1-a"

	node_pool {
		name               = "%s"
		initial_node_count = 2
		autoscaling {
			min_node_count = 1
			max_node_count = 3
		}
	}
}`, cluster, np)
}

func testAccContainerCluster_withNodePoolUpdateAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
	name = "%s"
	zone = "us-central1-a"

	node_pool {
		name               = "%s"
		initial_node_count = 2
		autoscaling {
			min_node_count = 1
			max_node_count = 5
		}
	}
}`, cluster, np)
}

func testAccContainerCluster_withNodePoolNamePrefix() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_name_prefix" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"

	node_pool {
		name_prefix = "tf-np-test"
		node_count  = 2
	}
}`, acctest.RandString(10))
}

func testAccContainerCluster_withNodePoolMultiple() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_multiple" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"

	node_pool {
		name       = "tf-cluster-nodepool-test-%s"
		node_count = 2
	}

	node_pool {
		name       = "tf-cluster-nodepool-test-%s"
		node_count = 3
	}
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerCluster_withNodePoolConflictingNameFields() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_multiple" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"

	node_pool {
		# ERROR: name and name_prefix cannot be both specified
		name        = "tf-cluster-nodepool-test-%s"
		name_prefix = "tf-cluster-nodepool-test-"
		node_count  = 1
	}
}`, acctest.RandString(10), acctest.RandString(10))
}

func testAccContainerCluster_withNodePoolNodeConfig() string {
	testId := acctest.RandString(10)
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_node_config" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"
	node_pool {
		name = "tf-cluster-nodepool-test-%s"
		node_count = 2
		node_config {
			machine_type = "n1-standard-1"
			disk_size_gb = 15
			local_ssd_count = 1
			oauth_scopes = [
				"https://www.googleapis.com/auth/compute",
				"https://www.googleapis.com/auth/devstorage.read_only",
				"https://www.googleapis.com/auth/logging.write",
				"https://www.googleapis.com/auth/monitoring"
			]
			service_account = "default"
			metadata {
				foo = "bar"
			}
			image_type = "COS"
			labels {
				foo = "bar"
			}
			tags = ["foo", "bar"]
		}
	}

}
`, testId, testId)
}

func testAccContainerCluster_withDefaultNodePoolRemoved() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_default_node_pool_removed" {
	name               = "cluster-test-%s"
	zone               = "us-central1-a"
	initial_node_count = 1

	remove_default_node_pool = true
}
`, acctest.RandString(10))
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
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	%s
}`, clusterName, maintenancePolicy)
}

func testAccContainerCluster_withIPAllocationPolicy_existingSecondaryRanges(cluster string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
	name          = "${google_compute_network.container_network.name}"
	network       = "${google_compute_network.container_network.name}"
	ip_cidr_range = "10.0.0.0/24"
	region        = "us-central1"

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
	name = "%s"
	zone = "us-central1-a"

	network = "${google_compute_network.container_network.name}"
	subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

	initial_node_count = 1
	ip_allocation_policy {
		cluster_secondary_range_name = "pods"
		services_secondary_range_name = "services"
	}
}`, cluster, cluster)
}

func testAccContainerCluster_withIPAllocationPolicy_specificIPRanges(cluster string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_ip_allocation_policy" {
	name = "%s"
	zone = "us-central1-a"

	initial_node_count = 1
	ip_allocation_policy {
		cluster_ipv4_cidr_block = "10.90.0.0/19"
		services_ipv4_cidr_block = "10.40.0.0/19"
	}
}`, cluster)
}

func testAccContainerCluster_withIPAllocationPolicy_specificSizes(cluster string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
	name          = "${google_compute_network.container_network.name}"
	network       = "${google_compute_network.container_network.name}"
	ip_cidr_range = "10.0.0.0/24"
	region        = "us-central1"
}

resource "google_container_cluster" "with_ip_allocation_policy" {
	name = "%s"
	zone = "us-central1-a"

	network = "${google_compute_network.container_network.name}"
	subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

	initial_node_count = 1
	ip_allocation_policy {
		cluster_ipv4_cidr_block = "/16"
		services_ipv4_cidr_block = "/22"
	}
}`, cluster, cluster)
}

func testAccContainerCluster_withIPAllocationPolicy_createSubnetwork(cluster string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_ip_allocation_policy" {
	name = "%s"
	zone = "us-central1-a"

	initial_node_count = 1
	ip_allocation_policy {
		create_subnetwork = true
	}
}`, cluster)
}

func testAccContainerCluster_withPodSecurityPolicy(clusterName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_pod_security_policy" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	pod_security_policy_config {
		enabled = %v
	}
}`, clusterName, enabled)
}

func testAccContainerCluster_withPrivateCluster(clusterName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
	name                     = "${google_compute_network.container_network.name}"
	network                  = "${google_compute_network.container_network.name}"
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
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1

	network = "${google_compute_network.container_network.name}"
	subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

	private_cluster = true
	master_ipv4_cidr_block = "10.42.0.0/28"
	master_authorized_networks_config { cidr_blocks = [] }
	ip_allocation_policy {
		cluster_secondary_range_name  = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.0.range_name}"
		services_secondary_range_name = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.1.range_name}"
	}
}`, clusterName, clusterName)
}

func testAccContainerCluster_sharedVpc(org, billingId, projectName, name string) string {
	return fmt.Sprintf(`
resource "google_project" "host_project" {
	name            = "Test Project XPN Host"
	project_id      = "%s-host"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "host_project" {
	project = "${google_project.host_project.project_id}"
	service = "container.googleapis.com"
}

resource "google_compute_shared_vpc_host_project" "host_project" {
	project = "${google_project_service.host_project.project}"
}

resource "google_project" "service_project" {
	name            = "Test Project XPN Service"
	project_id      = "%s-service"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "service_project" {
	project = "${google_project.service_project.project_id}"
	service = "container.googleapis.com"
}

resource "google_compute_shared_vpc_service_project" "service_project" {
	host_project    = "${google_compute_shared_vpc_host_project.host_project.project}"
	service_project = "${google_project_service.service_project.project}"
}

resource "google_project_iam_member" "host_service_agent" {
	project = "${google_project_service.host_project.project}"
	role    = "roles/container.hostServiceAgentUser"
	member  = "serviceAccount:service-${google_project.service_project.number}@container-engine-robot.iam.gserviceaccount.com"

	depends_on = ["google_project_service.service_project"]
}

resource "google_compute_subnetwork_iam_member" "service_network_cloud_services" {
	project       = "${google_compute_shared_vpc_host_project.host_project.project}"
	subnetwork    = "${google_compute_subnetwork.shared_subnetwork.name}"
	role          = "roles/compute.networkUser"
	member        = "serviceAccount:${google_project.service_project.number}@cloudservices.gserviceaccount.com"
}

resource "google_compute_subnetwork_iam_member" "service_network_gke_user" {
	project       = "${google_compute_shared_vpc_host_project.host_project.project}"
	subnetwork    = "${google_compute_subnetwork.shared_subnetwork.name}"
	role          = "roles/compute.networkUser"
	member        = "serviceAccount:service-${google_project.service_project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_compute_network" "shared_network" {
	name    = "test-%s"
	project = "${google_compute_shared_vpc_host_project.host_project.project}"

	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "shared_subnetwork" {
	name          = "test-%s"
	ip_cidr_range = "10.0.0.0/16"
	region        = "us-central1"
	network       = "${google_compute_network.shared_network.self_link}"
	project       = "${google_compute_shared_vpc_host_project.host_project.project}"

	secondary_ip_range {
		range_name = "pods"
		ip_cidr_range = "10.1.0.0/16"
	}

	secondary_ip_range {
		range_name = "services"
		ip_cidr_range = "10.2.0.0/20"
	}
}

resource "google_container_cluster" "shared_vpc_cluster" {
	name               = "%s"
	zone               = "us-central1-a"
	initial_node_count = 1
	project            = "${google_compute_shared_vpc_service_project.service_project.service_project}"

	network    = "${google_compute_network.shared_network.self_link}"
	subnetwork = "${google_compute_subnetwork.shared_subnetwork.self_link}"

	ip_allocation_policy {
		cluster_secondary_range_name  = "${google_compute_subnetwork.shared_subnetwork.secondary_ip_range.0.range_name}"
		services_secondary_range_name = "${google_compute_subnetwork.shared_subnetwork.secondary_ip_range.1.range_name}"
	}

	depends_on = [
		"google_project_iam_member.host_service_agent",
		"google_compute_subnetwork_iam_member.service_network_cloud_services",
		"google_compute_subnetwork_iam_member.service_network_gke_user"
	]
}`, projectName, org, billingId, projectName, org, billingId, acctest.RandString(10), acctest.RandString(10), name)
}

func testAccContainerCluster_withoutResourceLabels(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_resource_labels" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
}
`, clusterName)
}

func testAccContainerCluster_withResourceLabels(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_resource_labels" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	resource_labels {
		created-by = "terraform"
	}
}
`, clusterName)
}

func testAccContainerCluster_withBinaryAuthorization(clusterName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_binary_authorization" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	enable_binary_authorization = %v
}
`, clusterName, enabled)
}
