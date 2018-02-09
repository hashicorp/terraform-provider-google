package google

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.primary"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.primary"),
				),
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
					testAccCheckContainerCluster(
						"google_container_cluster.primary"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.http_load_balancing.0.disabled", "true"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.kubernetes_dashboard.0.disabled", "true"),
				),
			},
			{
				Config: testAccContainerCluster_updateAddons(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.primary"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.horizontal_pod_autoscaling.0.disabled", "true"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.http_load_balancing.0.disabled", "false"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "addons_config.0.kubernetes_dashboard.0.disabled", "true"),
				),
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuth(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuth(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_master_auth"),
				),
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
					testAccCheckContainerCluster(
						"google_container_cluster.with_network_policy_enabled"),
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.#", "1"),
				),
			},
			{
				Config: testAccContainerCluster_removeNetworkPolicy(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_network_policy_enabled"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy"),
				),
			},
			{
				Config: testAccContainerCluster_withNetworkPolicyDisabled(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_network_policy_enabled"),
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.0.enabled", "false"),
				),
			},
			{
				Config:             testAccContainerCluster_withNetworkPolicyDisabled(clusterName),
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
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{"0.0.0.0/0"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster("google_container_cluster.with_master_authorized_networks"),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks.#", "1"),
				),
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster("google_container_cluster.with_master_authorized_networks"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks"),
				),
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{"8.8.8.8/32"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster("google_container_cluster.with_master_authorized_networks"),
					resource.TestCheckResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks.#", "1"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_additional_zones"),
				),
			},
			{
				Config: testAccContainerCluster_updateAdditionalZones(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_additional_zones"),
				),
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
					testAccCheckContainerCluster(
						"google_container_cluster.with_kubernetes_alpha"),
					resource.TestCheckResourceAttr("google_container_cluster.with_kubernetes_alpha", "enable_kubernetes_alpha", "true"),
				),
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
					testAccCheckContainerCluster(
						"google_container_cluster.with_legacy_abac"),
					resource.TestCheckResourceAttr("google_container_cluster.with_legacy_abac", "enable_legacy_abac", "true"),
				),
			},
			{
				Config: testAccContainerCluster_updateLegacyAbac(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_legacy_abac"),
					resource.TestCheckResourceAttr("google_container_cluster.with_legacy_abac", "enable_legacy_abac", "false"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_version"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_version"),
				),
			},
			{
				Config: testAccContainerCluster_updateVersion(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_version"),
				),
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfig(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_config"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_config_scope_alias"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_net_ref_by_url"),
					testAccCheckContainerCluster(
						"google_container_cluster.with_net_ref_by_name"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.primary"),
				),
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
					testAccCheckContainerCluster(
						"google_container_cluster.with_logging"),
					resource.TestCheckResourceAttr("google_container_cluster.with_logging", "logging_service", "logging.googleapis.com"),
				),
			},
			{
				Config: testAccContainerCluster_updateLogging(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_logging"),
					resource.TestCheckResourceAttr("google_container_cluster.with_logging", "logging_service", "none"),
				),
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
					testAccCheckContainerCluster(
						"google_container_cluster.with_monitoring"),
					resource.TestCheckResourceAttr("google_container_cluster.with_monitoring", "monitoring_service", "monitoring.googleapis.com"),
				),
			},
			{
				Config: testAccContainerCluster_updateMonitoring(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_monitoring"),
					resource.TestCheckResourceAttr("google_container_cluster.with_monitoring", "monitoring_service", "none"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_pool"),
				),
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
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_pool"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "2"),
				),
			},
			{
				Config: testAccContainerCluster_withNodePoolResize(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_pool"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "3"),
				),
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
					testAccCheckContainerCluster("google_container_cluster.with_node_pool"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "3"),
				),
			},
			resource.TestStep{
				Config: testAccContainerCluster_withNodePoolUpdateAutoscaling(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster("google_container_cluster.with_node_pool"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "5"),
				),
			},
			resource.TestStep{
				Config: testAccContainerCluster_withNodePoolBasic(clusterName, npName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster("google_container_cluster.with_node_pool"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_pool_name_prefix"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_pool_multiple"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_node_pool_node_config"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(resourceName),
				),
			},
			{
				Config: testAccContainerCluster_withMaintenanceWindow(clusterName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.daily_maintenance_window.0.start_time"),
					testAccCheckContainerCluster(resourceName),
				),
			},
		},
	})
}

func TestAccContainerCluster_withIPAllocationPolicy(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("cluster-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIPAllocationPolicy(
					cluster,
					map[string]string{
						"pods":     "10.1.0.0/16",
						"services": "10.2.0.0/20",
					},
					map[string]string{
						"cluster_secondary_range_name":  "pods",
						"services_secondary_range_name": "services",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerCluster(
						"google_container_cluster.with_ip_allocation_policy"),
					resource.TestCheckResourceAttr("google_container_cluster.with_ip_allocation_policy",
						"ip_allocation_policy.0.cluster_secondary_range_name", "pods"),
					resource.TestCheckResourceAttr("google_container_cluster.with_ip_allocation_policy",
						"ip_allocation_policy.0.services_secondary_range_name", "services"),
				),
			},
			{
				Config: testAccContainerCluster_withIPAllocationPolicy(
					cluster,
					map[string]string{
						"pods":     "10.1.0.0/16",
						"services": "10.2.0.0/20",
					},
					map[string]string{},
				),
				ExpectError: regexp.MustCompile("clusters using IP aliases must specify secondary ranges"),
			},
			{
				Config: testAccContainerCluster_withIPAllocationPolicy(
					cluster,
					map[string]string{
						"pods": "10.1.0.0/16",
					},
					map[string]string{
						"cluster_secondary_range_name":  "pods",
						"services_secondary_range_name": "services",
					},
				),
				ExpectError: regexp.MustCompile("services secondary range \"pods\" not found in subnet"),
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

var setFields = []string{
	"additional_zones",
	"node_config.0.oauth_scopes",
	"node_pool.[0-9]*.node_config.0.oauth_scopes",
}

func testAccCheckContainerCluster(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)
		cluster, err := config.clientContainer.Projects.Zones.Clusters.Get(
			config.Project, attributes["zone"], attributes["name"]).Do()
		if err != nil {
			return err
		}

		if cluster.Name != attributes["name"] {
			return fmt.Errorf("Cluster %s not found, found %s instead", attributes["name"], cluster.Name)
		}

		type clusterTestField struct {
			tf_attr  string
			gcp_attr interface{}
		}

		var igUrls []string
		if igUrls, err = getInstanceGroupUrlsFromManagerUrls(config, cluster.InstanceGroupUrls); err != nil {
			return err
		}
		clusterTests := []clusterTestField{
			{"initial_node_count", strconv.FormatInt(cluster.InitialNodeCount, 10)},
			{"master_auth.0.client_certificate", cluster.MasterAuth.ClientCertificate},
			{"master_auth.0.client_key", cluster.MasterAuth.ClientKey},
			{"master_auth.0.cluster_ca_certificate", cluster.MasterAuth.ClusterCaCertificate},
			{"master_auth.0.password", cluster.MasterAuth.Password},
			{"master_auth.0.username", cluster.MasterAuth.Username},
			{"zone", cluster.Zone},
			{"cluster_ipv4_cidr", cluster.ClusterIpv4Cidr},
			{"description", cluster.Description},
			{"enable_kubernetes_alpha", strconv.FormatBool(cluster.EnableKubernetesAlpha)},
			{"enable_legacy_abac", strconv.FormatBool(cluster.LegacyAbac.Enabled)},
			{"endpoint", cluster.Endpoint},
			{"instance_group_urls", igUrls},
			{"logging_service", cluster.LoggingService},
			{"monitoring_service", cluster.MonitoringService},
			{"subnetwork", cluster.Subnetwork},
			{"node_config.0.machine_type", cluster.NodeConfig.MachineType},
			{"node_config.0.disk_size_gb", strconv.FormatInt(cluster.NodeConfig.DiskSizeGb, 10)},
			{"node_config.0.local_ssd_count", strconv.FormatInt(cluster.NodeConfig.LocalSsdCount, 10)},
			{"node_config.0.oauth_scopes", cluster.NodeConfig.OauthScopes},
			{"node_config.0.service_account", cluster.NodeConfig.ServiceAccount},
			{"node_config.0.metadata", cluster.NodeConfig.Metadata},
			{"node_config.0.image_type", cluster.NodeConfig.ImageType},
			{"node_config.0.labels", cluster.NodeConfig.Labels},
			{"node_config.0.tags", cluster.NodeConfig.Tags},
			{"node_config.0.preemptible", cluster.NodeConfig.Preemptible},
			{"node_config.0.min_cpu_platform", cluster.NodeConfig.MinCpuPlatform},
			{"node_version", cluster.CurrentNodeVersion},
		}

		if cluster.NetworkPolicy != nil {
			clusterTests = append(clusterTests,
				clusterTestField{"network_policy.0.enabled", cluster.NetworkPolicy.Enabled},
				clusterTestField{"network_policy.0.provider", cluster.NetworkPolicy.Provider},
			)
		}
		// Remove Zone from additional_zones since that's what the resource writes in state
		additionalZones := []string{}
		for _, location := range cluster.Locations {
			if location != cluster.Zone {
				additionalZones = append(additionalZones, location)
			}
		}
		clusterTests = append(clusterTests, clusterTestField{"additional_zones", additionalZones})

		// AddonsConfig is neither Required or Computed, so the API may return nil for it.
		httpLoadBalancingDisabled := false
		if cluster.AddonsConfig != nil && cluster.AddonsConfig.HttpLoadBalancing != nil {
			httpLoadBalancingDisabled = cluster.AddonsConfig.HttpLoadBalancing.Disabled
		}
		horizontalPodAutoscalingDisabled := false
		if cluster.AddonsConfig != nil && cluster.AddonsConfig.HorizontalPodAutoscaling != nil {
			horizontalPodAutoscalingDisabled = cluster.AddonsConfig.HorizontalPodAutoscaling.Disabled
		}
		kubernetesDashboardDisabled := false
		if cluster.AddonsConfig != nil && cluster.AddonsConfig.KubernetesDashboard != nil {
			kubernetesDashboardDisabled = cluster.AddonsConfig.KubernetesDashboard.Disabled
		}
		clusterTests = append(clusterTests, clusterTestField{"addons_config.0.http_load_balancing.0.disabled", httpLoadBalancingDisabled})
		clusterTests = append(clusterTests, clusterTestField{"addons_config.0.horizontal_pod_autoscaling.0.disabled", horizontalPodAutoscalingDisabled})
		clusterTests = append(clusterTests, clusterTestField{"addons_config.0.kubernetes_dashboard.0.disabled", kubernetesDashboardDisabled})

		if cluster.MaintenancePolicy != nil {
			clusterTests = append(clusterTests, clusterTestField{"maintenance_policy.0.daily_maintenance_window.0.start_time", cluster.MaintenancePolicy.Window.DailyMaintenanceWindow.StartTime})
			clusterTests = append(clusterTests, clusterTestField{"maintenance_policy.0.daily_maintenance_window.0.duration", cluster.MaintenancePolicy.Window.DailyMaintenanceWindow.Duration})
		}

		if cluster.IpAllocationPolicy != nil && cluster.IpAllocationPolicy.UseIpAliases {
			clusterTests = append(clusterTests, clusterTestField{"ip_allocation_policy.0.cluster_secondary_range_name", cluster.IpAllocationPolicy.ClusterSecondaryRangeName})
			clusterTests = append(clusterTests, clusterTestField{"ip_allocation_policy.0.services_secondary_range_name", cluster.IpAllocationPolicy.ServicesSecondaryRangeName})
		}

		for i, np := range cluster.NodePools {
			prefix := fmt.Sprintf("node_pool.%d.", i)
			clusterTests = append(clusterTests, clusterTestField{prefix + "name", np.Name})
			if np.Config != nil {
				clusterTests = append(clusterTests,
					clusterTestField{prefix + "node_config.0.machine_type", np.Config.MachineType},
					clusterTestField{prefix + "node_config.0.disk_size_gb", strconv.FormatInt(np.Config.DiskSizeGb, 10)},
					clusterTestField{prefix + "node_config.0.local_ssd_count", strconv.FormatInt(np.Config.LocalSsdCount, 10)},
					clusterTestField{prefix + "node_config.0.oauth_scopes", np.Config.OauthScopes},
					clusterTestField{prefix + "node_config.0.service_account", np.Config.ServiceAccount},
					clusterTestField{prefix + "node_config.0.metadata", np.Config.Metadata},
					clusterTestField{prefix + "node_config.0.image_type", np.Config.ImageType},
					clusterTestField{prefix + "node_config.0.labels", np.Config.Labels},
					clusterTestField{prefix + "node_config.0.tags", np.Config.Tags})

			}
			tfAS := attributes[prefix+"autoscaling.#"] == "1"
			if gcpAS := np.Autoscaling != nil && np.Autoscaling.Enabled == true; tfAS != gcpAS {
				return fmt.Errorf("Mismatched autoscaling status. TF State: %t. GCP State: %t", tfAS, gcpAS)
			}
			if tfAS {
				if tf := attributes[prefix+"autoscaling.0.min_node_count"]; strconv.FormatInt(np.Autoscaling.MinNodeCount, 10) != tf {
					return fmt.Errorf("Mismatched Autoscaling.MinNodeCount. TF State: %s. GCP State: %d",
						tf, np.Autoscaling.MinNodeCount)
				}

				if tf := attributes[prefix+"autoscaling.0.max_node_count"]; strconv.FormatInt(np.Autoscaling.MaxNodeCount, 10) != tf {
					return fmt.Errorf("Mismatched Autoscaling.MaxNodeCount. TF State: %s. GCP State: %d",
						tf, np.Autoscaling.MaxNodeCount)
				}
			}
		}

		for _, attrs := range clusterTests {
			if c := checkMatch(attributes, attrs.tf_attr, attrs.gcp_attr); c != "" {
				return fmt.Errorf(c)
			}
		}

		// Network has to be done separately in order to normalize the two values
		tf, err := getNetworkNameFromSelfLink(attributes["network"])
		if err != nil {
			return err
		}
		gcp, err := getNetworkNameFromSelfLink(cluster.Network)
		if err != nil {
			return err
		}
		if tf != gcp {
			return fmt.Errorf(matchError("network", tf, gcp))
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
		for _, setField := range setFields {
			if match, _ := regexp.MatchString(setField, attr); match {
				return checkSetMatch(attributes, attr, gcpList)
			}
		}
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

func checkSetMatch(attributes map[string]string, attr string, gcpList []string) string {
	num, err := strconv.Atoi(attributes[attr+".#"])
	if err != nil {
		return fmt.Sprintf("Error in number conversion for attribute %s: %s", attr, err)
	}
	if num != len(gcpList) {
		return fmt.Sprintf("Cluster has mismatched %s size.\nTF Size: %d\nGCP Size: %d", attr, num, len(gcpList))
	}

	// We don't know the exact keys of the elements, so go through the whole list looking for matching ones
	tfAttr := []string{}
	for k, v := range attributes {
		if strings.HasPrefix(k, attr) && !strings.HasSuffix(k, "#") {
			tfAttr = append(tfAttr, v)
		}
	}
	sort.Strings(tfAttr)
	sort.Strings(gcpList)
	if reflect.DeepEqual(tfAttr, gcpList) {
		return ""
	}
	return matchError(attr, tfAttr, gcpList)
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

func testAccContainerCluster_withMasterAuth() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_auth" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm.123456789"
	}
}`, acctest.RandString(10))
}

func testAccContainerCluster_withNetworkPolicyEnabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	network_policy {
		enabled = true
		provider = "CALICO"
	}	
}`, clusterName)
}

func testAccContainerCluster_removeNetworkPolicy(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1
}`, clusterName)
}

func testAccContainerCluster_withNetworkPolicyDisabled(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 1

	network_policy = {}
}`, clusterName)
}

func testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName string, cidrs []string) string {

	cidrBlocks := ""
	if len(cidrs) > 0 {
		var buf bytes.Buffer
		for _, c := range cidrs {
			buf.WriteString(fmt.Sprintf(`
			cidr_blocks {
				cidr_block = "%s"
			}`, c))
		}
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
		"us-central1-b",
		"us-central1-c",
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

func testAccContainerCluster_withNodeConfig() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config" {
	name = "cluster-test-%s"
	zone = "us-central1-f"
	initial_node_count = 1

	node_config {
		machine_type = "n1-standard-1"
		disk_size_gb = 15
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
		image_type = "COS"
		labels {
			foo = "bar"
		}
		tags = ["foo", "bar"]
		preemptible = true
		min_cpu_platform = "Intel Broadwell"
	}
}`, acctest.RandString(10))
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

func testAccContainerCluster_withIPAllocationPolicy(cluster string, ranges, policy map[string]string) string {

	var secondaryRanges bytes.Buffer
	for rangeName, cidr := range ranges {
		secondaryRanges.WriteString(fmt.Sprintf(`
	secondary_ip_range {
	    range_name    = "%s"
	    ip_cidr_range = "%s"
	}`, rangeName, cidr))
	}

	var ipAllocationPolicy bytes.Buffer
	for key, value := range policy {
		ipAllocationPolicy.WriteString(fmt.Sprintf(`
		%s = "%s"`, key, value))
	}

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

	%s
}

resource "google_container_cluster" "with_ip_allocation_policy" {
	name = "%s"
	zone = "us-central1-a"

	network = "${google_compute_network.container_network.name}"
	subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

	initial_node_count = 1
	ip_allocation_policy {
	    %s
	}
}`, acctest.RandString(10), secondaryRanges.String(), cluster, ipAllocationPolicy.String())
}
