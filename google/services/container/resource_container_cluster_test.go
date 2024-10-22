// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/container"
)

func TestAccContainerCluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_cluster.primary", "services_ipv4_cidr"),
					resource.TestCheckResourceAttrSet("google_container_cluster.primary", "self_link"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "networking_mode", "VPC_NATIVE"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportStateId:           fmt.Sprintf("us-central1-a/%s", clusterName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportStateId:           fmt.Sprintf("%s/us-central1-a/%s", envvar.GetTestProjectFromEnv(), clusterName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_resourceManagerTags(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()

	randomSuffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randomSuffix)

	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_resourceManagerTags(pid, clusterName, networkName, subnetworkName, randomSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_cluster.primary", "self_link"),
					resource.TestCheckResourceAttrSet("google_container_cluster.primary", "node_config.0.resource_manager_tags.%"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportStateId:           fmt.Sprintf("us-central1-a/%s", clusterName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_networkingModeRoutes(t *testing.T) {
	t.Parallel()

	firstClusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	secondClusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_networkingModeRoutes(firstClusterName, secondClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "networking_mode", "ROUTES"),
					resource.TestCheckResourceAttr("google_container_cluster.secondary", "networking_mode", "ROUTES")),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_container_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_misc(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_misc(clusterName, networkName, subnetworkName),
				// Explicitly check removing the default node pool since we won't
				// catch it by just importing.
				Check: resource.TestCheckResourceAttr(
					"google_container_cluster.primary", "node_pool.#", "0"),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection", "resource_labels", "terraform_labels"},
			},
			{
				Config: testAccContainerCluster_misc_update(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection", "resource_labels", "terraform_labels"},
			},
		},
	})
}

func TestAccContainerCluster_withAddons(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAddons(pid, clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_cluster.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: clean up this list in `4.0.0`, remove both `workload_identity_config` fields (same for below)
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateAddons(pid, clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			// Issue with cloudrun_config addon: https://github.com/hashicorp/terraform-provider-google/issues/11943
			// {
			// 	Config: testAccContainerCluster_withInternalLoadBalancer(pid, clusterName, networkName, subnetworkName),
			// },
			// {
			// 	ResourceName:            "google_container_cluster.primary",
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			// },
		},
	})
}

func TestAccContainerCluster_withDeletionProtection(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDeletionProtection(clusterName, networkName, subnetworkName, "false"),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withDeletionProtection(clusterName, networkName, subnetworkName, "true"),
			},
			{
				Config:      testAccContainerCluster_withDeletionProtection(clusterName, networkName, subnetworkName, "true"),
				Destroy:     true,
				ExpectError: regexp.MustCompile("Cannot destroy cluster because deletion_protection is set to true. Set it to false to proceed with cluster deletion."),
			},
			{
				Config: testAccContainerCluster_withDeletionProtection(clusterName, networkName, subnetworkName, "false"),
			},
		},
	})
}

func TestAccContainerCluster_withNotificationConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	newTopic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNotificationConfig(clusterName, topic, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNotificationConfig(clusterName, newTopic, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_disableNotificationConfig(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNotificationConfig(clusterName, newTopic, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withFilteredNotificationConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	newTopic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withFilteredNotificationConfig(clusterName, topic, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.filtered_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withFilteredNotificationConfigUpdate(clusterName, newTopic, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.filtered_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_disableFilteredNotificationConfig(clusterName, newTopic, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.filtered_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withConfidentialNodes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withConfidentialNodes(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.confidential_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_disableConfidentialNodes(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.confidential_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withConfidentialNodes(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.confidential_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withILBSubsetting(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_disableILBSubSetting(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.confidential_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withILBSubSetting(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.confidential_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_disableILBSubSetting(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.confidential_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withMultiNetworking(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_enableMultiNetworking(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAdditiveVPC(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAdditiveVPC(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withMasterAuthConfig_NoCert(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMasterAuthNoCert(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_auth_no_cert", "master_auth.0.client_certificate", ""),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_master_auth_no_cert",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAuthenticatorGroupsConfig(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	orgDomain := envvar.GetTestOrgDomainFromEnv(t)
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.primary",
						"authenticator_groups_config.0.enabled"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAuthenticatorGroupsConfigUpdate(clusterName, orgDomain, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary",
						"authenticator_groups_config.0.security_group", fmt.Sprintf("gke-security-groups@%s", orgDomain)),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAuthenticatorGroupsConfigUpdate2(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.primary",
						"authenticator_groups_config.0.enabled"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestUnitContainerCluster_Rfc3339TimeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"same time, format changed to have leading zero": {
			Old:                "2:00",
			New:                "02:00",
			ExpectDiffSuppress: true,
		},
		"same time, format changed not to have leading zero": {
			Old:                "02:00",
			New:                "2:00",
			ExpectDiffSuppress: true,
		},
		"different time, both without leading zero": {
			Old:                "2:00",
			New:                "3:00",
			ExpectDiffSuppress: false,
		},
		"different time, old with leading zero, new without": {
			Old:                "02:00",
			New:                "3:00",
			ExpectDiffSuppress: false,
		},
		"different time, new with leading zero, oldwithout": {
			Old:                "2:00",
			New:                "03:00",
			ExpectDiffSuppress: false,
		},
		"different time, both with leading zero": {
			Old:                "02:00",
			New:                "03:00",
			ExpectDiffSuppress: false,
		},
	}
	for tn, tc := range cases {
		if container.Rfc3339TimeDiffSuppress("time", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, '%s' => '%s' expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func testAccContainerCluster_enableMultiNetworking(clusterName string) string {
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
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }

  secondary_ip_range {
    range_name    = "another-pod"
    ip_cidr_range = "10.1.32.0/22"
  }

  lifecycle {
    ignore_changes = [
      # The auto nodepool creates a secondary range which diffs this resource.
      secondary_ip_range,
    ]
  }
}

resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
  release_channel {
	channel = "RAPID"
  }
  enable_multi_networking = true
  datapath_provider = "ADVANCED_DATAPATH"
  deletion_protection = false
}
`, clusterName, clusterName)
}

func testAccContainerCluster_withAdditiveVPC(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  dns_config {
    cluster_dns = "CLOUD_DNS"
    additive_vpc_scope_dns_domain = "test.com"
    cluster_dns_scope = "CLUSTER_SCOPE"
  }
  deletion_protection = false
}
`, clusterName)
}

func TestAccContainerCluster_withNetworkPolicyEnabled(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNetworkPolicyEnabled(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.#", "1"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_removeNetworkPolicy(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNetworkPolicyDisabled(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"network_policy.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNetworkPolicyConfigDisabled(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_network_policy_enabled",
						"addons_config.0.network_policy_config.0.disabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_network_policy_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config:             testAccContainerCluster_withNetworkPolicyConfigDisabled(clusterName, networkName, subnetworkName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccContainerCluster_withReleaseChannelEnabled(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "STABLE", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "UNSPECIFIED", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withReleaseChannelEnabledDefaultVersion(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withReleaseChannelEnabledDefaultVersion(clusterName, "REGULAR", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "REGULAR", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabledDefaultVersion(clusterName, "EXTENDED", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "EXTENDED", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withReleaseChannelEnabled(clusterName, "UNSPECIFIED", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_release_channel",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withInvalidReleaseChannel(t *testing.T) {
	// This is essentially a unit test, no interactions
	acctest.SkipIfVcr(t)
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withReleaseChannelEnabled(clusterName, "CANARY", networkName, subnetworkName),
				ExpectError: regexp.MustCompile(`expected release_channel\.0\.channel to be one of \["?UNSPECIFIED"? "?RAPID"? "?REGULAR"? "?STABLE"? "?EXTENDED"?\], got CANARY`),
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
				ResourceName:            "google_container_cluster.with_master_authorized_networks",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{"10.0.0.0/8", "8.8.8.8/32"}, ""),
			},
			{
				ResourceName:            "google_container_cluster.with_master_authorized_networks",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMasterAuthorizedNetworksConfig(clusterName, []string{}, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_master_authorized_networks",
						"master_authorized_networks_config.0.cidr_blocks.#", "0"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_master_authorized_networks",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_removeMasterAuthorizedNetworksConfig(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_master_authorized_networks",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withGcpPublicCidrsAccessEnabledToggle(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withoutGcpPublicCidrsAccessEnabled(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_gcp_public_cidrs_access_enabled",
						"master_authorized_networks_config.#", "0"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_gcp_public_cidrs_access_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withGcpPublicCidrsAccessEnabled(clusterName, "false", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_gcp_public_cidrs_access_enabled",
						"master_authorized_networks_config.0.gcp_public_cidrs_access_enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_gcp_public_cidrs_access_enabled",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withGcpPublicCidrsAccessEnabled(clusterName, "true", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_gcp_public_cidrs_access_enabled",
						"master_authorized_networks_config.0.gcp_public_cidrs_access_enabled", "true"),
				),
			},
		},
	})
}

func testAccContainerCluster_withGcpPublicCidrsAccessEnabled(clusterName string, flag, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, flag, networkName, subnetworkName)
}

func testAccContainerCluster_withoutGcpPublicCidrsAccessEnabled(clusterName, networkName, subnetworkName string) string {

	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_gcp_public_cidrs_access_enabled" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func TestAccContainerCluster_regional(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-regional-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_regional(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.regional",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_regionalWithNodePool(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-regional-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_regionalWithNodePool(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.regional",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_regionalWithNodeLocations(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_regionalNodeLocations(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_regionalUpdateNodeLocations(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ResourceName:            "google_container_cluster.with_private_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withPrivateClusterConfig(containerNetName, clusterName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_private_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ResourceName:            "google_container_cluster.with_private_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_private_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(clusterName, networkName, subnetworkName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_private_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withIntraNodeVisibility(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withIntraNodeVisibility(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_intranode_visibility", "enable_intranode_visibility", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_intranode_visibility",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateIntraNodeVisibility(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_intranode_visibility", "enable_intranode_visibility", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_intranode_visibility",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withVersion(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_version",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_updateVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLowerVersion(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_version",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateVersion(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_version",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfig(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_config.0.taint", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodeConfigUpdate(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_config.0.taint", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigGcfsConfig(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigGcfsConfig(clusterName, networkName, subnetworkName, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_node_config_gcfs_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodeConfigGcfsConfig(clusterName, networkName, subnetworkName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_node_config_gcfs_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigKubeletConfigSettingsUpdates(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigKubeletConfigSettingsBaseline(clusterName, networkName, subnetworkName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_node_config_kubelet_config_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodeConfigKubeletConfigSettingsUpdates(clusterName, "none", "100ms", "TRUE", networkName, subnetworkName, 2048, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_node_config_kubelet_config_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodeConfigKubeletConfigSettingsUpdates(clusterName, "static", "", "FALSE", networkName, subnetworkName, 1024, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_node_config_kubelet_config_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withInsecureKubeletReadonlyPortEnabledInNodePool(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePoolName := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledInNodePool(clusterName, nodePoolName, networkName, subnetworkName, "TRUE"),
			},
			{
				ResourceName:            "google_container_cluster.with_insecure_kubelet_readonly_port_enabled_in_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

// This is for `node_pool_defaults.node_config_defaults` - the default settings
// for newly created nodepools
func TestAccContainerCluster_withInsecureKubeletReadonlyPortEnabledDefaultsUpdates(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			// Test API default (no value set in config) first
			{
				Config: testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledDefaultsUpdateBaseline(clusterName, networkName, subnetworkName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_insecure_kubelet_readonly_port_enabled_node_pool_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledDefaultsUpdate(clusterName, networkName, subnetworkName, "TRUE"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_insecure_kubelet_readonly_port_enabled_node_pool_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledDefaultsUpdate(clusterName, networkName, subnetworkName, "FALSE"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_insecure_kubelet_readonly_port_enabled_node_pool_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledDefaultsUpdate(clusterName, networkName, subnetworkName, "TRUE"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectNoDelete(),
					},
				},
			},
			{
				ResourceName:            "google_container_cluster.with_insecure_kubelet_readonly_port_enabled_node_pool_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withLoggingVariantInNodeConfig(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLoggingVariantInNodeConfig(clusterName, "MAX_THROUGHPUT", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_logging_variant_in_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withLoggingVariantInNodePool(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePoolName := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLoggingVariantInNodePool(clusterName, nodePoolName, "MAX_THROUGHPUT", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_logging_variant_in_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withLoggingVariantUpdates(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, "DEFAULT", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_logging_variant_node_pool_default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, "MAX_THROUGHPUT", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_logging_variant_node_pool_default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, "DEFAULT", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_logging_variant_node_pool_default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAdvancedMachineFeaturesInNodePool(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePoolName := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAdvancedMachineFeaturesInNodePool(clusterName, nodePoolName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_advanced_machine_features_in_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolDefaults(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.primary",
						"node_pool_defaults.0.node_config_defaults.0.gcfs_config.0.enabled"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportStateId:           fmt.Sprintf("us-central1-a/%s", clusterName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodePoolDefaults(clusterName, "true", networkName, subnetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool_defaults",
						"node_pool_defaults.0.node_config_defaults.0.gcfs_config.#", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool_defaults",
						"node_pool_defaults.0.node_config_defaults.0.gcfs_config.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool_defaults",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodePoolDefaults(clusterName, "false", networkName, subnetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool_defaults",
						"node_pool_defaults.0.node_config_defaults.0.gcfs_config.#", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool_defaults",
						"node_pool_defaults.0.node_config_defaults.0.gcfs_config.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool_defaults",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigScopeAlias(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigScopeAlias(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_config_scope_alias",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigShieldedInstanceConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigShieldedInstanceConfig(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigReservationAffinity(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigReservationAffinity(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.#", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_config",
						"node_config.0.reservation_affinity.0.consume_reservation_type", "ANY_RESERVATION"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodeConfigReservationAffinitySpecific(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-reservation-%s", acctest.RandString(t, 10))
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodeConfigReservationAffinitySpecific(reservationName, clusterName, networkName, subnetworkName),
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
				ResourceName:            "google_container_cluster.with_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withWorkloadMetadataConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withWorkloadMetadataConfig(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_workload_metadata_config",
						"node_config.0.workload_metadata_config.0.mode", "GCE_METADATA"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_metadata_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withBootDiskKmsKey(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withBootDiskKmsKey(clusterName, kms.CryptoKey.Name, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_boot_disk_kms_key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
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
				ResourceName:            "google_container_cluster.with_net_ref_by_url",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_container_cluster.with_net_ref_by_name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_backend(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_backendRef(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolBasic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolBasic(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolUpdateVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolLowerVersion(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodePoolUpdateVersion(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolResize(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNodeLocations(clusterName, npName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "2"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodePoolResize(clusterName, npName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.node_count", "3"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolAutoscaling(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolAutoscaling(clusterName, npName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "3"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodePoolUpdateAutoscaling(clusterName, npName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count", "5"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodePoolBasic(clusterName, npName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.min_node_count"),
					resource.TestCheckNoResourceAttr("google_container_cluster.with_node_pool", "node_pool.0.autoscaling.0.max_node_count"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolCIA(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerRegionalCluster_withNodePoolCIA(clusterName, npName, networkName, subnetworkName),
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerRegionalClusterUpdate_withNodePoolCIA(clusterName, npName, networkName, subnetworkName),
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerRegionalCluster_withNodePoolBasic(clusterName, npName, networkName, subnetworkName),
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolNamePrefix(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	npNamePrefix := "tf-test-np-"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNamePrefix(clusterName, npNamePrefix, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool_name_prefix",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_pool.0.name_prefix", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withNodePoolMultiple(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	npNamePrefix := "tf-test-np-"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolMultiple(clusterName, npNamePrefix, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool_multiple",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withNodePoolNodeConfig(cluster, np, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_node_pool_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withMaintenanceWindow(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_window"
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMaintenanceWindow(clusterName, "03:00", networkName, subnetworkName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMaintenanceWindow(clusterName, "", networkName, subnetworkName),
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
				ImportStateVerifyIgnore: []string{"maintenance_policy.#", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withRecurringMaintenanceWindow(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_recurring_maintenance_window"
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withRecurringMaintenanceWindow(cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.daily_maintenance_window.0.start_time"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withRecurringMaintenanceWindow(cluster, "", "", networkName, subnetworkName),
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
				ImportStateVerifyIgnore: []string{"maintenance_policy.#", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withMaintenanceExclusionWindow(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_window"
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusion_RecurringMaintenanceWindow(cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", networkName, subnetworkName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withExclusion_DailyMaintenanceWindow(cluster, "2020-01-01T00:00:00Z", "2020-01-02T00:00:00Z", networkName, subnetworkName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withMaintenanceExclusionOptions(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_options"
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_MINOR_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_MINOR_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_OR_NODE_UPGRADES"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_deleteMaintenanceExclusionOptions(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_options"
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_OR_NODE_UPGRADES"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_NoExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope"),
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_updateMaintenanceExclusionOptions(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_options"
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

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
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope"),
					resource.TestCheckNoResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_MINOR_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_MINOR_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_OR_NODE_UPGRADES"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateExclusionOptions_RecurringMaintenanceWindow(
					cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", "NO_UPGRADES", "NO_MINOR_UPGRADES", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.0.exclusion_options.0.scope", "NO_UPGRADES"),
					resource.TestCheckResourceAttr(resourceName,
						"maintenance_policy.0.maintenance_exclusion.1.exclusion_options.0.scope", "NO_MINOR_UPGRADES"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_deleteExclusionWindow(t *testing.T) {
	t.Parallel()
	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	resourceName := "google_container_cluster.with_maintenance_exclusion_window"
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExclusion_DailyMaintenanceWindow(cluster, "2020-01-01T00:00:00Z", "2020-01-02T00:00:00Z", networkName, subnetworkName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withExclusion_RecurringMaintenanceWindow(cluster, "2019-01-01T00:00:00Z", "2019-01-02T00:00:00Z", "2019-05-01T00:00:00Z", "2019-05-02T00:00:00Z", networkName, subnetworkName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withExclusion_NoMaintenanceWindow(cluster, "2020-01-01T00:00:00Z", "2020-01-02T00:00:00Z", networkName, subnetworkName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ResourceName:            "google_container_cluster.with_ip_allocation_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ResourceName:            "google_container_cluster.with_ip_allocation_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ResourceName:            "google_container_cluster.with_ip_allocation_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioning(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioning(clusterName, networkName, subnetworkName, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning",
						"cluster_autoscaling.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autoprovisioning(clusterName, networkName, subnetworkName, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning",
						"cluster_autoscaling.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaults(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	includeMinCpuPlatform := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaults(clusterName, networkName, subnetworkName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning",
						"cluster_autoscaling.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config:             testAccContainerCluster_autoprovisioningDefaults(clusterName, networkName, subnetworkName, true),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsMinCpuPlatform(clusterName, networkName, subnetworkName, includeMinCpuPlatform),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsMinCpuPlatform(clusterName, networkName, subnetworkName, !includeMinCpuPlatform),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(clusterName, networkName, subnetworkName, 2, 1, "SURGE"),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_upgrade_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config:      testAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(clusterName, networkName, subnetworkName, 2, 1, "BLUE_GREEN"),
				ExpectError: regexp.MustCompile(`Surge upgrade settings max_surge/max_unavailable can only be used when strategy is set to SURGE`),
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsUpgradeSettingsWithBlueGreenStrategy(clusterName, networkName, subnetworkName, "3.500s", "BLUE_GREEN"),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_upgrade_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningNetworkTags(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioning(clusterName, networkName, subnetworkName, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning",
						"node_pool_auto_config.0.network_tags.0.tags.0", "test-network-tag"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withShieldedNodes(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withShieldedNodes(clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_shielded_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withShieldedNodes(clusterName, networkName, subnetworkName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_shielded_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autopilot", "networking_mode", "VPC_NATIVE"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
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

func TestAccContainerCluster_withAutopilotNetworkTags(t *testing.T) {
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
				Config: testAccContainerCluster_withAutopilot(pid, containerNetName, clusterName, "us-central1", true, true, ""),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAutopilotKubeletConfig(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randomSuffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutopilotKubeletConfigBaseline(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot_kubelet_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAutopilotKubeletConfigUpdates(clusterName, "FALSE"),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot_kubelet_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAutopilotKubeletConfigUpdates(clusterName, "TRUE"),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot_kubelet_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAutopilot_withNodePoolDefaults(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randomSuffix)
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutopilot_withNodePoolDefaults(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAutopilotResourceManagerTags(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()

	randomSuffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randomSuffix)
	clusterNetName := fmt.Sprintf("tf-test-container-net-%s", randomSuffix)
	clusterSubnetName := fmt.Sprintf("tf-test-container-subnet-%s", randomSuffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutopilotResourceManagerTags(pid, clusterName, clusterNetName, clusterSubnetName, randomSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_cluster.with_autopilot", "self_link"),
					resource.TestCheckResourceAttrSet("google_container_cluster.with_autopilot", "node_pool_auto_config.0.resource_manager_tags.%"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAutopilotResourceManagerTagsUpdate1(pid, clusterName, clusterNetName, clusterSubnetName, randomSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_container_cluster.with_autopilot", "node_pool_auto_config.0.resource_manager_tags.%"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAutopilotResourceManagerTagsUpdate2(pid, clusterName, clusterNetName, clusterSubnetName, randomSuffix),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withWorkloadIdentityConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withWorkloadIdentityConfigEnabled(pid, clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_identity_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateWorkloadIdentityConfig(pid, clusterName, networkName, subnetworkName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_identity_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateWorkloadIdentityConfig(pid, clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_identity_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withWorkloadIdentityConfigAutopilot(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withWorkloadIdentityConfigEnabledAutopilot(pid, clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_workload_identity_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withIdentityServiceConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withIdentityServiceConfigEnabled(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withIdentityServiceConfigUpdated(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withSecretManagerConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withSecretManagerConfigEnabled(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withSecretManagerConfigUpdated(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withLoggingConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withLoggingConfigEnabled(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withLoggingConfigDisabled(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withLoggingConfigUpdated(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigAdvancedDatapathObservabilityConfigDisabled(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withMonitoringConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigEnabled(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigDisabled(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigUpdated(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigPrometheusUpdated(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			// Back to basic settings to test setting Prometheus on its own
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigPrometheusOnly(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withMonitoringConfigPrometheusOnly2(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withSoleTenantGroup(t *testing.T) {
	t.Parallel()

	resourceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withSoleTenantGroup(resourceName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAutoscalingProfile(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutoscalingProfile(clusterName, "BALANCED", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.autoscaling_with_profile",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAutoscalingProfile(clusterName, "OPTIMIZE_UTILIZATION", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.autoscaling_with_profile",
				ImportStateIdPrefix:     "us-central1-a/",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withInvalidAutoscalingProfile(t *testing.T) {
	// This is essentially a unit test, no interactions
	acctest.SkipIfVcr(t)
	t.Parallel()
	clusterName := fmt.Sprintf("cluster-test-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withAutoscalingProfile(clusterName, "AS_CHEAP_AS_POSSIBLE", networkName, subnetworkName),
				ExpectError: regexp.MustCompile(`expected cluster_autoscaling\.0\.autoscaling_profile to be one of \["?BALANCED"? "?OPTIMIZE_UTILIZATION"?\], got AS_CHEAP_AS_POSSIBLE`),
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsDiskSizeGb(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	includeDiskSizeGb := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskSizeGb(clusterName, networkName, subnetworkName, includeDiskSizeGb),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskSizeGb(clusterName, networkName, subnetworkName, !includeDiskSizeGb),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsDiskType(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	includeDiskType := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskType(clusterName, networkName, subnetworkName, includeDiskType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsDiskType(clusterName, networkName, subnetworkName, !includeDiskType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsImageType(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	includeImageType := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsImageType(clusterName, networkName, subnetworkName, includeImageType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsImageType(clusterName, networkName, subnetworkName, !includeImageType),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsBootDiskKmsKey(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsBootDiskKmsKey(clusterName, kms.CryptoKey.Name, networkName, subnetworkName),
			},
			{
				ResourceName:      "google_container_cluster.nap_boot_disk_kms_key",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"min_master_version",
					"deletion_protection",
					"node_pool", // cluster_autoscaling (node auto-provisioning) creates new node pools automatically
				},
			},
		},
	})
}

func TestAccContainerCluster_nodeAutoprovisioningDefaultsShieldedInstance(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsShieldedInstance(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.nap_shielded_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_autoprovisioningDefaultsManagement(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsManagement(clusterName, networkName, subnetworkName, false, false),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_management",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningDefaultsManagement(clusterName, networkName, subnetworkName, true, true),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_management",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_autoprovisioningLocations(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autoprovisioningLocations(clusterName, networkName, subnetworkName, []string{"us-central1-a", "us-central1-f"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning_locations",
						"cluster_autoscaling.0.enabled", "true"),

					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning_locations",
						"cluster_autoscaling.0.auto_provisioning_locations.0", "us-central1-a"),

					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning_locations",
						"cluster_autoscaling.0.auto_provisioning_locations.1", "us-central1-f"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autoprovisioningLocations(clusterName, networkName, subnetworkName, []string{"us-central1-b", "us-central1-c"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning_locations",
						"cluster_autoscaling.0.enabled", "true"),

					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning_locations",
						"cluster_autoscaling.0.auto_provisioning_locations.0", "us-central1-b"),

					resource.TestCheckResourceAttr("google_container_cluster.with_autoprovisioning_locations",
						"cluster_autoscaling.0.auto_provisioning_locations.1", "us-central1-c"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autoprovisioning_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

// This resource originally cleaned up the dangling cluster directly, but now
// taints it, having Terraform clean it up during the next apply. This test
// name is now inexact, but is being preserved to maintain the test history.
func TestAccContainerCluster_errorCleanDanglingCluster(t *testing.T) {
	acctest.SkipIfVcr(t) // skipped because the timeout step doesn't record operation GET interactions
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", suffix)
	clusterNameError := fmt.Sprintf("tf-test-cluster-err-%s", suffix)
	clusterNameErrorWithTimeout := fmt.Sprintf("tf-test-cluster-timeout-%s", suffix)
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))

	initConfig := testAccContainerCluster_withInitialCIDR(containerNetName, clusterName)
	overlapConfig := testAccContainerCluster_withCIDROverlap(initConfig, clusterNameError)
	overlapConfigWithTimeout := testAccContainerCluster_withCIDROverlapWithTimeout(initConfig, clusterNameErrorWithTimeout, "1s")

	checkTaintApplied := func(st *terraform.State) error {
		// Return an error if there is no tainted (i.e. marked for deletion) cluster.
		ms := st.RootModule()
		errCluster, ok := ms.Resources["google_container_cluster.cidr_error_overlap"]
		if !ok {
			var resourceNames []string
			for rn := range ms.Resources {
				resourceNames = append(resourceNames, rn)
			}
			return fmt.Errorf("could not find google_container_cluster.cidr_error_overlap in resources: %v", resourceNames)
		}
		if !errCluster.Primary.Tainted {
			return fmt.Errorf("cluster with ID %s should be tainted, but is not", errCluster.Primary.ID)
		}
		return nil
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: initConfig,
			},
			{
				ResourceName:            "google_container_cluster.cidr_error_preempt",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				// First attempt to create the overlapping cluster with no timeout, this should fail and taint the resource.
				Config:      overlapConfig,
				ExpectError: regexp.MustCompile("Error waiting for creating GKE cluster"),
			},
			{
				// Check that the tainted resource is in the config.
				Config:             overlapConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check:              checkTaintApplied,
			},
			{
				// Next attempt to create the overlapping cluster with a 1s timeout. This will fail with a different error.
				Config:      overlapConfigWithTimeout,
				ExpectError: regexp.MustCompile("timeout while waiting for state to become 'DONE'"),
			},
			{
				// Check that the tainted resource is in the config.
				Config:             overlapConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check:              checkTaintApplied,
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
				ExpectError: regexp.MustCompile(`(Location "wonderland" does not exist)|(Permission denied on 'locations\/wonderland' \(or it may not exist\))`),
			},
		},
	})
}

func TestAccContainerCluster_withExternalIpsConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withExternalIpsConfig(pid, clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_external_ips_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withExternalIpsConfig(pid, clusterName, networkName, subnetworkName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_external_ips_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withMeshCertificatesConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withMeshCertificatesConfigEnabled(pid, clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_mesh_certificates_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateMeshCertificatesConfig(pid, clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_mesh_certificates_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateMeshCertificatesConfig(pid, clusterName, networkName, subnetworkName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_mesh_certificates_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withCostManagementConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	pid := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_updateCostManagementConfig(pid, clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.with_cost_management_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_updateCostManagementConfig(pid, clusterName, networkName, subnetworkName, false),
			},
			{
				ResourceName:            "google_container_cluster.with_cost_management_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withDatabaseEncryption(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

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
				Config: testAccContainerCluster_withDatabaseEncryption(clusterName, kmsData, networkName, subnetworkName),
				Check:  resource.TestCheckResourceAttrSet("data.google_kms_key_ring_iam_policy.test_key_ring_iam_policy", "policy_data"),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withAdvancedDatapath(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDatapathProvider(clusterName, "ADVANCED_DATAPATH", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_enableCiliumPolicies(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withDatapathProvider(clusterName, "ADVANCED_DATAPATH", networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "enable_cilium_clusterwide_network_policy", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_enableCiliumPolicies(clusterName, networkName, subnetworkName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "enable_cilium_clusterwide_network_policy", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_enableCiliumPolicies(clusterName, networkName, subnetworkName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "enable_cilium_clusterwide_network_policy", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_enableCiliumPolicies_withAutopilot(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randomSuffix)
	clusterNetName := fmt.Sprintf("tf-test-container-net-%s", randomSuffix)
	clusterSubnetName := fmt.Sprintf("tf-test-container-subnet-%s", randomSuffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_enableCiliumPolicies_withAutopilot(clusterName, clusterNetName, clusterSubnetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autopilot", "enable_cilium_clusterwide_network_policy", "false"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_enableCiliumPolicies_withAutopilotUpdate(clusterName, clusterNetName, clusterSubnetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autopilot", "enable_cilium_clusterwide_network_policy", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withResourceUsageExportConfig(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", suffix)
	datesetId := fmt.Sprintf("tf_test_cluster_resource_usage_%s", suffix)
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withResourceUsageExportConfig(clusterName, datesetId, "true", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_resource_usage_export_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withResourceUsageExportConfig(clusterName, datesetId, "false", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_resource_usage_export_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withResourceUsageExportConfigNoConfig(clusterName, datesetId, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_resource_usage_export_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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
				ResourceName:            "google_container_cluster.with_private_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withEnableKubernetesAlpha(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-np-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withEnableKubernetesAlpha(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withEnableKubernetesBetaAPIs(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withEnableKubernetesBetaAPIs(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withEnableKubernetesBetaAPIsOnExistingCluster(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withoutEnableKubernetesBetaAPIs(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withEnableKubernetesBetaAPIs(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withIncompatibleMasterVersionNodeVersion(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withIncompatibleMasterVersionNodeVersion(clusterName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Resource argument node_version`),
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
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withDNSConfig(clusterName, "CLOUD_DNS", domainName, "VPC_SCOPE", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withGatewayApiConfig(t *testing.T) {
	t.Parallel()
	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccContainerCluster_withGatewayApiConfig(clusterName, "CANARY", networkName, subnetworkName),
				ExpectError: regexp.MustCompile(`expected gateway_api_config\.0\.channel to be one of [^,]+, got CANARY`),
			},
			{
				Config: testAccContainerCluster_withGatewayApiConfig(clusterName, "CHANNEL_DISABLED", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withGatewayApiConfig(clusterName, "CHANNEL_STANDARD", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withSecurityPostureConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_SetSecurityPostureToStandard(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_security_posture_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_SetSecurityPostureToEnterprise(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_security_posture_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_SetWorkloadVulnerabilityToStandard(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_security_posture_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_SetWorkloadVulnerabilityToEnterprise(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_security_posture_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_DisableALL(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_security_posture_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_withFleetConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withFleetConfig(clusterName, projectID, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config:      testAccContainerCluster_withFleetConfig(clusterName, "random-project", networkName, subnetworkName),
				ExpectError: regexp.MustCompile(`changing existing fleet host project is not supported`),
			},
			{
				Config: testAccContainerCluster_DisableFleet(clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_withFleetConfig(name, projectID, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  fleet {
	project = "%s"
  }

  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, projectID, networkName, subnetworkName)
}

func testAccContainerCluster_DisableFleet(resource_name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, resource_name, networkName, subnetworkName)
}

func testAccContainerCluster_withIncompatibleMasterVersionNodeVersion(name string) string {
	return fmt.Sprintf(`
	resource "google_container_cluster" "gke_cluster" {
		name = "%s"
		location = "us-central1"

		min_master_version = "1.10.9-gke.5"
		node_version = "1.10.6-gke.11"
		initial_node_count = 1

	}
	`, name)
}

func testAccContainerCluster_SetSecurityPostureToStandard(resource_name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	mode = "BASIC"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, resource_name, networkName, subnetworkName)
}

func testAccContainerCluster_SetSecurityPostureToEnterprise(resource_name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	mode = "ENTERPRISE"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, resource_name, networkName, subnetworkName)
}

func testAccContainerCluster_SetWorkloadVulnerabilityToStandard(resource_name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	vulnerability_mode = "VULNERABILITY_BASIC"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, resource_name, networkName, subnetworkName)
}

func testAccContainerCluster_SetWorkloadVulnerabilityToEnterprise(resource_name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	vulnerability_mode = "VULNERABILITY_ENTERPRISE"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, resource_name, networkName, subnetworkName)
}

func testAccContainerCluster_DisableALL(resource_name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_security_posture_config" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  security_posture_config {
	mode = "DISABLED"
	vulnerability_mode = "VULNERABILITY_DISABLED"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, resource_name, networkName, subnetworkName)
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
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_autopilot_net_admin(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_autopilot_net_admin(clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autopilot_net_admin(clusterName, networkName, subnetworkName, false),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_autopilot_net_admin(clusterName, networkName, subnetworkName, true),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_additional_pod_ranges_config_on_create(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_additional_pod_ranges_config(clusterName, 1),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccContainerCluster_additional_pod_ranges_config_on_update(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_additional_pod_ranges_config(clusterName, 0),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_additional_pod_ranges_config(clusterName, 2),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_additional_pod_ranges_config(clusterName, 0),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_additional_pod_ranges_config(clusterName, 1),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_additional_pod_ranges_config(clusterName, 0),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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

func testAccContainerCluster_basic(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_networkingModeRoutes(firstName, secondName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  networking_mode    = "ROUTES"
  deletion_protection = false
}

resource "google_container_cluster" "secondary" {
	name               = "%s"
	location           = "us-central1-a"
	initial_node_count = 1
	cluster_ipv4_cidr  = "10.96.0.0/14"
	deletion_protection = false
  }
`, firstName, secondName)
}

func testAccContainerCluster_misc(name, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_misc_update(name, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withAddons(projectID, clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  min_master_version = "latest"
  release_channel {
    channel = "RAPID"
  }

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
    stateful_ha_config {
      enabled = false
    }
    ray_operator_config {
      enabled = false
    }
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_updateAddons(projectID, clusterName, networkName, subnetworkName string) string {
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
    stateful_ha_config {
      enabled = true
    }
    ray_operator_config {
      enabled = true
      ray_cluster_logging_config {
        enabled = true
      }
      ray_cluster_monitoring_config {
        enabled = true
      }
    }
	}
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, networkName, subnetworkName)
}

// Issue with cloudrun_config addon: https://github.com/hashicorp/terraform-provider-google/issues/11943/
// func testAccContainerCluster_withInternalLoadBalancer(projectID string, clusterName, networkName, subnetworkName string) string {
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
//   deletion_protection = false
//   network    = "%s"
//   subnetwork    = "%s"
// }
// `, projectID, clusterName, networkName, subnetworkName)
// }

func testAccContainerCluster_withNotificationConfig(clusterName, topic, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, topic, topic, clusterName, topic, networkName, subnetworkName)
}

func testAccContainerCluster_disableNotificationConfig(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withFilteredNotificationConfig(clusterName, topic, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, topic, topic, clusterName, topic, networkName, subnetworkName)
}

func testAccContainerCluster_withFilteredNotificationConfigUpdate(clusterName, topic, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, topic, topic, clusterName, topic, networkName, subnetworkName)
}

func testAccContainerCluster_disableFilteredNotificationConfig(clusterName, topic, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, topic, topic, clusterName, topic, networkName, subnetworkName)
}

func testAccContainerCluster_withConfidentialNodes(clusterName, npName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, npName, networkName, subnetworkName)
}

func testAccContainerCluster_disableConfidentialNodes(clusterName, npName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, npName, networkName, subnetworkName)
}

func testAccContainerCluster_withILBSubSetting(clusterName, npName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, npName, networkName, subnetworkName)
}

func testAccContainerCluster_disableILBSubSetting(clusterName, npName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, npName, networkName, subnetworkName)
}

func testAccContainerCluster_withNetworkPolicyEnabled(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withDeletionProtection(clusterName, networkName, subnetworkName, deletionProtection string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  deletion_protection = %s
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, deletionProtection, networkName, subnetworkName)
}

func testAccContainerCluster_withReleaseChannelEnabled(clusterName, channel, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_release_channel" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  release_channel {
    channel = "%s"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, channel, networkName, subnetworkName)
}

func testAccContainerCluster_withReleaseChannelEnabledDefaultVersion(clusterName, channel, networkName, subnetworkName string) string {
	return fmt.Sprintf(`

data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_release_channel" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  min_master_version = data.google_container_engine_versions.central1a.release_channel_default_version["%s"]
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, channel, networkName, subnetworkName)
}

func testAccContainerCluster_removeNetworkPolicy(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
  name                     = "%s"
  location                 = "us-central1-a"
  initial_node_count       = 1
  remove_default_node_pool = true
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNetworkPolicyDisabled(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_network_policy_enabled" {
  name                     = "%s"
  location                 = "us-central1-a"
  initial_node_count       = 1
  remove_default_node_pool = true

  network_policy {
    enabled = false
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNetworkPolicyConfigDisabled(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withAuthenticatorGroupsConfigUpdate(name, orgDomain, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  authenticator_groups_config {
    security_group = "gke-security-groups@%s"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, orgDomain, networkName, subnetworkName)
}

func testAccContainerCluster_withAuthenticatorGroupsConfigUpdate2(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  authenticator_groups_config {
    security_group = ""
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
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
  deletion_protection = false
}
`, clusterName, cidrBlocks)
}

func testAccContainerCluster_removeMasterAuthorizedNetworksConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_master_authorized_networks" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
}
`, clusterName)
}

func testAccContainerCluster_regional(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "regional" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
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
  deletion_protection = false
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
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
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
  deletion_protection = false
}
`, containerNetName, clusterName)
}

func TestAccContainerCluster_withCidrBlockWithoutPrivateEndpointSubnetwork(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	containerNetName := fmt.Sprintf("tf-test-container-net-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withCidrBlockWithoutPrivateEndpointSubnetwork(containerNetName, clusterName, "us-central1-a"),
			},
			{
				ResourceName:            "google_container_cluster.with_private_flexible_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_withCidrBlockWithoutPrivateEndpointSubnetwork(containerNetName, clusterName, location string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
}

resource "google_container_cluster" "with_private_flexible_cluster" {
  name               = "%s"
  location           = "%s"
  min_master_version = "1.29"
  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_nodes    = true
	master_ipv4_cidr_block  = "10.42.0.0/28"
  }
  deletion_protection = false
}
`, containerNetName, clusterName, location)
}

func TestAccContainerCluster_withEnablePrivateEndpointToggle(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withEnablePrivateEndpoint(clusterName, "true", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_enable_private_endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withEnablePrivateEndpoint(clusterName, "false", networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_enable_private_endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"min_master_version", "deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_withEnablePrivateEndpoint(clusterName, flag, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, flag, networkName, subnetworkName)
}

func testAccContainerCluster_regionalWithNodePool(cluster, nodePool, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "regional" {
  name     = "%s"
  location = "us-central1"

  node_pool {
    name = "%s"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, nodePool, networkName, subnetworkName)
}

func testAccContainerCluster_regionalNodeLocations(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_locations" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  node_locations = [
    "us-central1-f",
    "us-central1-c",
  ]
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_regionalUpdateNodeLocations(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_locations" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  node_locations = [
    "us-central1-f",
    "us-central1-b",
  ]
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withIntraNodeVisibility(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_intranode_visibility" {
  name                        = "%s"
  location                    = "us-central1-a"
  initial_node_count          = 1
  enable_intranode_visibility = true
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_updateIntraNodeVisibility(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_intranode_visibility" {
  name                        = "%s"
  location                    = "us-central1-a"
  initial_node_count          = 1
  enable_intranode_visibility = false
  private_ipv6_google_access  = "PRIVATE_IPV6_GOOGLE_ACCESS_BIDIRECTIONAL"
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withVersion(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withLowerVersion(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.release_channel_default_version["STABLE"]
  node_version       = data.google_container_engine_versions.central1a.release_channel_default_version["STABLE"]
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withMasterAuthNoCert(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_updateVersion(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_version" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.release_channel_latest_version["STABLE"]
  node_version       = data.google_container_engine_versions.central1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfig(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigGcfsConfig(clusterName, networkName, subnetworkName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config_gcfs_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    gcfs_config {
      enabled = %t
    }
  }

  deletion_protection = false
  network             = "%s"
  subnetwork          = "%s"
}
`, clusterName, enabled, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigKubeletConfigSettingsBaseline(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config_kubelet_config_settings" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    kubelet_config {
      pod_pids_limit = 1024
    }
  }
  deletion_protection = false
  network             = "%s"
  subnetwork          = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigKubeletConfigSettingsUpdates(clusterName, cpuManagerPolicy, cpuCfsQuotaPeriod, insecureKubeletReadonlyPortEnabled, networkName, subnetworkName string, podPidsLimit int, cpuCfsQuota bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_config_kubelet_config_settings" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    kubelet_config {
      cpu_manager_policy                     = "%s"
      cpu_cfs_quota                          = %v
      cpu_cfs_quota_period                   = "%s"
      insecure_kubelet_readonly_port_enabled = "%s"
      pod_pids_limit                         = %v
    }
  }
  deletion_protection = false
  network             = "%s"
  subnetwork          = "%s"
}
`, clusterName, cpuManagerPolicy, cpuCfsQuota, cpuCfsQuotaPeriod, insecureKubeletReadonlyPortEnabled, podPidsLimit, networkName, subnetworkName)
}

func testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledInNodePool(clusterName, nodePoolName, networkName, subnetworkName, insecureKubeletReadonlyPortEnabled string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_insecure_kubelet_readonly_port_enabled_in_node_pool" {
  name               = "%s"
  location           = "us-central1-f"

  node_pool {
    name               = "%s"
    initial_node_count = 1
    node_config {
      kubelet_config {
        cpu_manager_policy                     = "static"
        insecure_kubelet_readonly_port_enabled = "%s"
      }
    }
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, nodePoolName, insecureKubeletReadonlyPortEnabled, networkName, subnetworkName)
}

func testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledDefaultsUpdateBaseline(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_insecure_kubelet_readonly_port_enabled_node_pool_update" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withInsecureKubeletReadonlyPortEnabledDefaultsUpdate(clusterName, networkName, subnetworkName, insecureKubeletReadonlyPortEnabled string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_insecure_kubelet_readonly_port_enabled_node_pool_update" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_pool_defaults {
    node_config_defaults {
      insecure_kubelet_readonly_port_enabled = "%s"
    }
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, insecureKubeletReadonlyPortEnabled, networkName, subnetworkName)
}

func testAccContainerCluster_withLoggingVariantInNodeConfig(clusterName, loggingVariant, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_logging_variant_in_node_config" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_config {
    logging_variant = "%s"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, loggingVariant, networkName, subnetworkName)
}

func testAccContainerCluster_withLoggingVariantInNodePool(clusterName, nodePoolName, loggingVariant, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, nodePoolName, loggingVariant, networkName, subnetworkName)
}

func testAccContainerCluster_withLoggingVariantNodePoolDefault(clusterName, loggingVariant, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, loggingVariant, networkName, subnetworkName)
}

func testAccContainerCluster_withAdvancedMachineFeaturesInNodePool(clusterName, nodePoolName, networkName, subnetworkName string, nvEnabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_advanced_machine_features_in_node_pool" {
  name               = "%s"
  location           = "us-central1-f"

  node_pool {
    name               = "%s"
    initial_node_count = 1
    node_config {
      machine_type = "c2-standard-4"
      advanced_machine_features {
        threads_per_core = 1
        enable_nested_virtualization = "%t"
	    }
    }
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, nodePoolName, nvEnabled, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolDefaults(clusterName, enabled, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_defaults" {
  name               = "%s"
  location           = "us-central1-f"
  initial_node_count = 1

  node_pool_defaults {
    node_config_defaults {
      gcfs_config {
        enabled = "%s"
      }
    }
  }
  deletion_protection = false
  network    = "%s"
  subnetwork = "%s"
}
`, clusterName, enabled, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigUpdate(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigScopeAlias(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigShieldedInstanceConfig(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigReservationAffinity(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigReservationAffinitySpecific(reservation, clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
  depends_on = [google_project_service.container]
}
`, reservation, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withWorkloadMetadataConfig(clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withBootDiskKmsKey(clusterName, kmsKeyName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, kmsKeyName, networkName, subnetworkName)
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
  deletion_protection = false
}

resource "google_container_cluster" "with_net_ref_by_name" {
  name               = "%s-name"
  location           = "us-central1-a"
  initial_node_count = 1

  network = google_compute_network.container_network.name
  deletion_protection = false
}
`, network, cluster, cluster)
}

func testAccContainerCluster_autoprovisioningDefaultsManagement(clusterName, networkName, subnetworkName string, autoUpgrade, autoRepair bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, autoUpgrade, autoRepair, networkName, subnetworkName)
}

func testAccContainerCluster_autoprovisioningLocations(clusterName, networkName, subnetworkName string, locations []string) string {
	var autoprovisionLocationsStr string
	for i := 0; i < len(locations); i++ {
		autoprovisionLocationsStr += fmt.Sprintf("\"%s\",", locations[i])
	}
	var apl string
	if len(autoprovisionLocationsStr) > 0 {
		apl = fmt.Sprintf(`
			auto_provisioning_locations = [%s]
		`, autoprovisionLocationsStr)
	}

	return fmt.Sprintf(`
resource "google_container_cluster" "with_autoprovisioning_locations" {
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

    %s
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, apl, networkName, subnetworkName)
}

func testAccContainerCluster_backendRef(cluster, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, cluster, cluster, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolBasic(cluster, nodePool, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool" {
  name     = "%s"
  location = "us-central1-a"
  deletion_protection = false

  node_pool {
    name               = "%s"
    initial_node_count = 2
  }

  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, nodePool, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolLowerVersion(cluster, nodePool, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, nodePool, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolUpdateVersion(cluster, nodePool, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, nodePool, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolNodeLocations(cluster, nodePool, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, nodePool, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolResize(cluster, nodePool, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, nodePool, networkName, subnetworkName)
}

func testAccContainerCluster_withAutoscalingProfile(cluster, autoscalingProfile, networkName, subnetworkName string) string {
	config := fmt.Sprintf(`
resource "google_container_cluster" "autoscaling_with_profile" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  cluster_autoscaling {
    enabled             = false
    autoscaling_profile = "%s"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, autoscalingProfile, networkName, subnetworkName)
	return config
}

func testAccContainerCluster_autoprovisioning(cluster, networkName, subnetworkName string, autoprovisioning, withNetworkTag bool) string {
	config := fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
`, cluster, networkName, subnetworkName)
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

func testAccContainerCluster_autoprovisioningDefaults(cluster, networkName, subnetworkName string, monitoringWrite bool) string {
	config := fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autoprovisioning" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.latest_master_version
  initial_node_count = 1
  deletion_protection = false

  network    = "%s"
  subnetwork    = "%s"

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
		cluster, networkName, subnetworkName)

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

func testAccContainerCluster_autoprovisioningDefaultsMinCpuPlatform(cluster, networkName, subnetworkName string, includeMinCpuPlatform bool) string {
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
  network    = "%s"
  subnetwork    = "%s"

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
  deletion_protection = false
}
`, cluster, networkName, subnetworkName, minCpuPlatformCfg)
}

func testAccContainerCluster_autoprovisioningDefaultsUpgradeSettings(clusterName, networkName, subnetworkName string, maxSurge, maxUnavailable int, strategy string) string {
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
      deletion_protection = false
      network    = "%s"
      subnetwork    = "%s"
    }
  `, clusterName, maxSurge, maxUnavailable, strategy, blueGreenSettings, networkName, subnetworkName)
}

func testAccContainerCluster_autoprovisioningDefaultsUpgradeSettingsWithBlueGreenStrategy(clusterName, networkName, subnetworkName string, duration, strategy string) string {
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
        deletion_protection = false
        network    = "%s"
        subnetwork    = "%s"
      }
    `, clusterName, strategy, duration, duration, networkName, subnetworkName)
}

func testAccContainerCluster_autoprovisioningDefaultsDiskSizeGb(cluster, networkName, subnetworkName string, includeDiskSizeGb bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, DiskSizeGbCfg, networkName, subnetworkName)
}

func testAccContainerCluster_autoprovisioningDefaultsDiskType(cluster, networkName, subnetworkName string, includeDiskType bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, DiskTypeCfg, networkName, subnetworkName)
}

func testAccContainerCluster_autoprovisioningDefaultsImageType(cluster, networkName, subnetworkName string, includeImageType bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, imageTypeCfg, networkName, subnetworkName)
}

func testAccContainerCluster_autoprovisioningDefaultsBootDiskKmsKey(clusterName, kmsKeyName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, kmsKeyName, networkName, subnetworkName)
}

func testAccContainerCluster_autoprovisioningDefaultsShieldedInstance(cluster, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolAutoscaling(cluster, np, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, np, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolUpdateAutoscaling(cluster, np, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, np, networkName, subnetworkName)
}

func testAccContainerRegionalCluster_withNodePoolCIA(cluster, np, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, np, networkName, subnetworkName)
}

func testAccContainerRegionalClusterUpdate_withNodePoolCIA(cluster, np, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, np, networkName, subnetworkName)
}

func testAccContainerRegionalCluster_withNodePoolBasic(cluster, nodePool, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, nodePool, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolNamePrefix(cluster, npPrefix, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_node_pool_name_prefix" {
  name     = "%s"
  location = "us-central1-a"

  node_pool {
    name_prefix = "%s"
    node_count  = 2
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, npPrefix, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolMultiple(cluster, npPrefix, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, npPrefix, npPrefix, networkName, subnetworkName)
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
  deletion_protection = false
}
`, cluster, npPrefix, npPrefix)
}

func testAccContainerCluster_withNodePoolNodeConfig(cluster, np, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, np, networkName, subnetworkName)
}

func testAccContainerCluster_withMaintenanceWindow(clusterName, startTime, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, maintenancePolicy, networkName, subnetworkName)
}

func testAccContainerCluster_withRecurringMaintenanceWindow(clusterName, startTime, endTime, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, maintenancePolicy, networkName, subnetworkName)

}

func testAccContainerCluster_withExclusion_RecurringMaintenanceWindow(clusterName string, w1startTime, w1endTime, w2startTime, w2endTime, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, w1startTime, w1endTime, w1startTime, w1endTime, w2startTime, w2endTime, networkName, subnetworkName)
}

func testAccContainerCluster_withExclusionOptions_RecurringMaintenanceWindow(cclusterName, w1startTime, w1endTime, w2startTime, w2endTime, scope1, scope2, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cclusterName, w1startTime, w1endTime, w1startTime, w1endTime, scope1, w2startTime, w2endTime, scope2, networkName, subnetworkName)
}

func testAccContainerCluster_NoExclusionOptions_RecurringMaintenanceWindow(cclusterName, w1startTime, w1endTime, w2startTime, w2endTime, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cclusterName, w1startTime, w1endTime, w1startTime, w1endTime, w2startTime, w2endTime, networkName, subnetworkName)
}

func testAccContainerCluster_updateExclusionOptions_RecurringMaintenanceWindow(cclusterName, w1startTime, w1endTime, w2startTime, w2endTime, scope1, scope2, networkName, subnetworkName string) string {

	return fmt.Sprintf(`
resource "google_container_cluster" "with_maintenance_exclusion_options" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false

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
  network    = "%s"
  subnetwork    = "%s"
}
`, cclusterName, w1startTime, w1endTime, w1startTime, w1endTime, scope1, w2startTime, w2endTime, scope2, networkName, subnetworkName)
}

func testAccContainerCluster_withExclusion_NoMaintenanceWindow(clusterName string, w1startTime, w1endTime, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, w1startTime, w1endTime, networkName, subnetworkName)
}

func testAccContainerCluster_withExclusion_DailyMaintenanceWindow(clusterName, w1startTime, w1endTime, networkName, subnetworkName string) string {

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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, w1startTime, w1endTime, networkName, subnetworkName)
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
  deletion_protection = false
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
  deletion_protection = false
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
  deletion_protection = false
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

    initial_node_count = 1
    datapath_provider = "ADVANCED_DATAPATH"
    enable_l4_ilb_subsetting = true

    ip_allocation_policy {
        cluster_ipv4_cidr_block  = "10.0.0.0/16"
        services_ipv4_cidr_block = "10.1.0.0/16"
        stack_type = "IPV4_IPV6"
    }
	deletion_protection = false
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

    initial_node_count = 1
    enable_l4_ilb_subsetting = true

    ip_allocation_policy {
        cluster_ipv4_cidr_block  = "10.0.0.0/16"
        services_ipv4_cidr_block = "10.1.0.0/16"
        stack_type = "IPV4"
    }
	deletion_protection = false
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
	deletion_protection = false
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withResourceUsageExportConfig(clusterName, datasetId, enableMetering, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, datasetId, clusterName, enableMetering, networkName, subnetworkName)
}

func testAccContainerCluster_withResourceUsageExportConfigNoConfig(clusterName, datasetId, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, datasetId, clusterName, networkName, subnetworkName)
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
  deletion_protection = false
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
  deletion_protection = false
}
`, containerNetName, clusterName, masterGlobalAccessEnabled)
}

func testAccContainerCluster_withPrivateClusterConfigGlobalAccessEnabledOnly(clusterName, networkName, subnetworkName string, masterGlobalAccessEnabled bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, masterGlobalAccessEnabled, networkName, subnetworkName)
}

func testAccContainerCluster_withShieldedNodes(clusterName, networkName, subnetworkName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_shielded_nodes" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  enable_shielded_nodes = %v
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, enabled, networkName, subnetworkName)
}

func testAccContainerCluster_withWorkloadIdentityConfigEnabled(projectID, clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withWorkloadIdentityConfigEnabledAutopilot(projectID string, clusterName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_container_cluster" "with_workload_identity_config" {
  name               = "%s"
  location           = "us-central1"
  initial_node_count = 1

  workload_identity_config {
    workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
  }
  enable_autopilot = true
  deletion_protection = false
}
`, projectID, clusterName)
}

func testAccContainerCluster_updateWorkloadIdentityConfig(projectID, clusterName, networkName, subnetworkName string, enable bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, workloadIdentityConfig, networkName, subnetworkName)
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
  deletion_protection = false
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
  deletion_protection = false
}
`, initConfig, secondCluster)
}

func testAccContainerCluster_withCIDROverlapWithTimeout(initConfig, secondCluster, createTimeout string) string {
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
  deletion_protection = false
  timeouts {
    create = "%s"
  }
}
`, initConfig, secondCluster, createTimeout)
}

func testAccContainerCluster_withInvalidLocation(location string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_resource_labels" {
  name               = "invalid-gke-cluster"
  location           = "%s"
  initial_node_count = 1
  deletion_protection = false
}
`, location)
}

func testAccContainerCluster_withExternalIpsConfig(projectID, clusterName, networkName, subnetworkName string, enabled bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, enabled, networkName, subnetworkName)
}

func testAccContainerCluster_withMeshCertificatesConfigEnabled(projectID, clusterName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_updateMeshCertificatesConfig(projectID, clusterName, networkName, subnetworkName string, enabled bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, enabled, networkName, subnetworkName)
}

func testAccContainerCluster_updateCostManagementConfig(projectID, clusterName, networkName, subnetworkName string, enabled bool) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, projectID, clusterName, enabled, networkName, subnetworkName)
}

func testAccContainerCluster_withDatabaseEncryption(clusterName string, kmsData acctest.BootstrappedKMS, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
}

resource "google_kms_key_ring_iam_member" "test_key_ring_iam_policy" {
  key_ring_id = "%[1]s"
  role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
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
  deletion_protection = false
  network    = "%[4]s"
  subnetwork = "%[5]s"
}
`, kmsData.KeyRing.Name, kmsData.CryptoKey.Name, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withDatapathProvider(clusterName, datapathProvider, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, datapathProvider, networkName, subnetworkName)
}

func testAccContainerCluster_enableCiliumPolicies(clusterName, networkName, subnetworkName string, enableCilium bool) string {
	ciliumPolicies := ""
	if enableCilium {
		ciliumPolicies = "enable_cilium_clusterwide_network_policy = true"
	} else {
		ciliumPolicies = "enable_cilium_clusterwide_network_policy = false"
	}

	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  ip_allocation_policy {
  }

  datapath_provider = "ADVANCED_DATAPATH"
  %s

  release_channel {
    channel = "RAPID"
  }

  network    = "%s"
  subnetwork    = "%s"

  deletion_protection = false
}
`, clusterName, ciliumPolicies, networkName, subnetworkName)
}

func testAccContainerCluster_enableCiliumPolicies_withAutopilot(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%[2]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = "%[3]s"
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

resource "google_container_cluster" "with_autopilot" {
  name = "%[1]s"
  location = "us-central1"
  enable_autopilot = true

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

  datapath_provider = "ADVANCED_DATAPATH"

  deletion_protection = false
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_enableCiliumPolicies_withAutopilotUpdate(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "container_network" {
  name                    = "%[2]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = "%[3]s"
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

resource "google_container_cluster" "with_autopilot" {
  name = "%[1]s"
  location = "us-central1"
  enable_autopilot = true

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

  datapath_provider = "ADVANCED_DATAPATH"
  enable_cilium_clusterwide_network_policy = true

  deletion_protection = false
}
`, clusterName, networkName, subnetworkName)
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
  deletion_protection = false
}
`, containerNetName, clusterName)
}

func testAccContainerCluster_withEnableKubernetesAlpha(cluster, np, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, np, networkName, subnetworkName)
}

func testAccContainerCluster_withoutEnableKubernetesBetaAPIs(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "central1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.central1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withEnableKubernetesBetaAPIs(cluster, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1
  deletion_protection = false

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
  network    = "%s"
  subnetwork    = "%s"
}
`, cluster, networkName, subnetworkName)
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
  deletion_protection = false
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

resource "google_project_iam_member" "project" {
	project = "%[2]s"
	role    = "roles/container.nodeServiceAccount"
	member = "serviceAccount:%[1]s@%[2]s.iam.gserviceaccount.com"
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
	deletion_protection = false
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

func testAccContainerCluster_withDNSConfig(clusterName, clusterDns, clusterDnsDomain, clusterDnsScope, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
	name               = "%s"
	location           = "us-central1-a"
  initial_node_count = 1
  dns_config {
    cluster_dns 	   = "%s"
    cluster_dns_domain = "%s"
    cluster_dns_scope  = "%s"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, clusterDns, clusterDnsDomain, clusterDnsScope, networkName, subnetworkName)
}

func testAccContainerCluster_withGatewayApiConfig(clusterName, gatewayApiChannel, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, gatewayApiChannel, networkName, subnetworkName)
}

func testAccContainerCluster_withIdentityServiceConfigEnabled(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  identity_service_config {
    enabled = true
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withIdentityServiceConfigUpdated(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  identity_service_config {
    enabled = false
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withSecretManagerConfigEnabled(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  secret_manager_config {
    enabled = true
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withSecretManagerConfigUpdated(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  secret_manager_config {
    enabled = false
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withLoggingConfigEnabled(name, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withLoggingConfigDisabled(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  logging_config {
    enable_components = []
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withLoggingConfigUpdated(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  logging_config {
    enable_components = [ "SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER", "SCHEDULER", "KCP_CONNECTION", "KCP_SSHD"]
  }
  monitoring_config {
    enable_components = [ "SYSTEM_COMPONENTS" ]
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withMonitoringConfigEnabled(name, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withMonitoringConfigDisabled(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
      enable_components = []
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withMonitoringConfigUpdated(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  monitoring_config {
         enable_components = [ "SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER" ]
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withMonitoringConfigPrometheusUpdated(name, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withMonitoringConfigPrometheusOnly(name, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_withMonitoringConfigPrometheusOnly2(name, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
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
      enable_relay   = true
    }
  }
  deletion_protection = false
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
      enable_relay   = false
    }
  }
  deletion_protection = false
}
`, name, name)
}

func testAccContainerCluster_withSoleTenantGroup(name, networkName, subnetworkName string) string {
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

  initial_size	= 1
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, name, name, networkName, subnetworkName)
}

func testAccContainerCluster_autopilot_minimal(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name             = "%s"
  location         = "us-central1"
  enable_autopilot = true
  deletion_protection = false
}`, name)
}

func testAccContainerCluster_autopilot_net_admin(name, networkName, subnetworkName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name             = "%s"
  location         = "us-central1"
  enable_autopilot = true
  allow_net_admin  = %t
  min_master_version = 1.27
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, enabled, networkName, subnetworkName)
}

func TestAccContainerCluster_customPlacementPolicy(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	policy := fmt.Sprintf("tf-test-policy-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_customPlacementPolicy(cluster, np, policy, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.cluster", "node_pool.0.placement_policy.0.type", "COMPACT"),
					resource.TestCheckResourceAttr("google_container_cluster.cluster", "node_pool.0.placement_policy.0.policy_name", policy),
					resource.TestCheckResourceAttr("google_container_cluster.cluster", "node_pool.0.node_config.0.machine_type", "c2-standard-4"),
				),
			},
			{
				ResourceName:            "google_container_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_customPlacementPolicy(cluster, np, policyName, networkName, subnetworkName string) string {
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
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, policyName, cluster, np, networkName, subnetworkName)
}

func testAccContainerCluster_additional_pod_ranges_config(name string, nameCount int) string {
	var podRangeNamesStr string
	names := []string{"\"gke-autopilot-pods-add\",", "\"gke-autopilot-pods-add-2\""}
	for i := 0; i < nameCount; i++ {
		podRangeNamesStr += names[i]
	}
	var aprc string
	if len(podRangeNamesStr) > 0 {
		aprc = fmt.Sprintf(`
			additional_pod_ranges_config {
				pod_range_names = [%s]
			}
		`, podRangeNamesStr)
	}

	return fmt.Sprintf(`
	resource "google_compute_network" "main" {
		name                    = "%s"
		auto_create_subnetworks = false
	}
	resource "google_compute_subnetwork" "main" {
		ip_cidr_range = "10.10.0.0/16"
		name          = "%s"
		network       = google_compute_network.main.self_link
		region        = "us-central1"

		secondary_ip_range {
			range_name    = "gke-autopilot-services"
			ip_cidr_range = "10.11.0.0/20"
		}

		secondary_ip_range {
			range_name    = "gke-autopilot-pods"
			ip_cidr_range = "10.12.0.0/16"
		}

		secondary_ip_range {
			range_name    = "gke-autopilot-pods-add"
			ip_cidr_range = "10.100.0.0/16"
		}
		secondary_ip_range {
			range_name    = "gke-autopilot-pods-add-2"
			ip_cidr_range = "100.0.0.0/16"
		}
	}
	resource "google_container_cluster" "primary" {
		name     = "%s"
		location = "us-central1"

		enable_autopilot = true

		release_channel {
			channel = "REGULAR"
		}

		network    = google_compute_network.main.name
		subnetwork = google_compute_subnetwork.main.name

		private_cluster_config {
			enable_private_endpoint = false
			enable_private_nodes    = true
			master_ipv4_cidr_block  = "172.16.0.0/28"
		}

		# supresses permadiff
		dns_config {
			cluster_dns = "CLOUD_DNS"
			cluster_dns_domain = "cluster.local"
			cluster_dns_scope = "CLUSTER_SCOPE"
		}

		ip_allocation_policy {
			cluster_secondary_range_name  = "gke-autopilot-pods"
			services_secondary_range_name = "gke-autopilot-services"
			%s
		}
		deletion_protection = false
	}
	`, name, name, name, aprc)
}

func TestAccContainerCluster_withConfidentialBootDisk(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-node-pool-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withConfidentialBootDisk(clusterName, npName, kms.CryptoKey.Name, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_confidential_boot_disk",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_withConfidentialBootDisk(clusterName, npName, kmsKeyName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_confidential_boot_disk" {
  name               = "%s"
  location           = "us-central1-a"
  confidential_nodes {
  	enabled = true
  }
  release_channel {
    channel = "RAPID"
  }
  node_pool {
    name = "%s"
    initial_node_count = 1
    node_config {
      oauth_scopes = [
        "https://www.googleapis.com/auth/cloud-platform",
      ]
      image_type = "COS_CONTAINERD"
      boot_disk_kms_key = "%s"
      machine_type = "n2d-standard-2"
      enable_confidential_storage = true
      disk_type = "hyperdisk-balanced"
    }
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, npName, kmsKeyName, networkName, subnetworkName)
}

func TestAccContainerCluster_withConfidentialBootDiskNodeConfig(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withConfidentialBootDiskNodeConfig(clusterName, kms.CryptoKey.Name, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.with_confidential_boot_disk_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_withConfidentialBootDiskNodeConfig(clusterName, kmsKeyName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_confidential_boot_disk_node_config" {
  name               = "%s"
  location           = "us-central1-a"
  confidential_nodes {
  	enabled = true
  }
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
    machine_type = "n2d-standard-2"
    enable_confidential_storage = true
    disk_type = "hyperdisk-balanced"
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, kmsKeyName, networkName, subnetworkName)
}

func TestAccContainerCluster_withoutConfidentialBootDisk(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	npName := fmt.Sprintf("tf-test-cluster-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withoutConfidentialBootDisk(clusterName, npName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.without_confidential_boot_disk",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}
func testAccContainerCluster_withoutConfidentialBootDisk(clusterName, npName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "without_confidential_boot_disk" {
  name               = "%s"
  location           = "us-central1-a"
  release_channel {
    channel = "RAPID"
  }
  node_pool {
    name = "%s"
    initial_node_count = 1
    node_config {
      oauth_scopes = [
       "https://www.googleapis.com/auth/cloud-platform",
      ]
      image_type = "COS_CONTAINERD"
      machine_type = "n2-standard-2"
      enable_confidential_storage = false
      disk_type = "pd-balanced"
    }
  }
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, clusterName, npName, networkName, subnetworkName)
}

func testAccContainerCluster_withAutopilotKubeletConfigBaseline(name string) string {
	return fmt.Sprintf(`
  resource "google_container_cluster" "with_autopilot_kubelet_config" {
    name                = "%s"
    location            = "us-central1"
    initial_node_count  = 1
    enable_autopilot    = true
    deletion_protection = false
  }
`, name)
}

func testAccContainerCluster_withAutopilotKubeletConfigUpdates(name, insecureKubeletReadonlyPortEnabled string) string {
	return fmt.Sprintf(`
  resource "google_container_cluster" "with_autopilot_kubelet_config" {
    name               = "%s"
    location           = "us-central1"
    initial_node_count = 1

    node_pool_auto_config {
      node_kubelet_config {
        insecure_kubelet_readonly_port_enabled = "%s"
      }
    }

    enable_autopilot    = true
    deletion_protection = false
  }
`, name, insecureKubeletReadonlyPortEnabled)
}

func testAccContainerCluster_withAutopilot_withNodePoolDefaults(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name                = "%s"
  location            = "us-central1"
  enable_autopilot    = true

  node_pool_defaults {
    node_config_defaults {
    }
  }

  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
 }
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_resourceManagerTags(projectID, clusterName, networkName, subnetworkName, randomSuffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%[1]s"
}

resource "google_project_iam_member" "tagHoldAdmin" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagHoldAdmin"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "tagUser1" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "google_project_iam_member" "tagUser2" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:${data.google_project.project.number}@cloudservices.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"

  depends_on = [
    google_project_iam_member.tagHoldAdmin,
    google_project_iam_member.tagUser1,
    google_project_iam_member.tagUser2,
  ]
}

resource "google_tags_tag_key" "key" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz-%[2]s"
  description = "For foo/bar resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }
}

resource "google_tags_tag_value" "value" {
  parent = "tagKeys/${google_tags_tag_key.key.name}"
  short_name = "foo-%[2]s"
  description = "For foo resources"
}

data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "primary" {
  name               = "%[3]s"
  location           = "us-central1-a"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["STABLE"]
  initial_node_count = 1

  node_config {
    machine_type    = "n1-standard-1"  // can't be e2 because of local-ssd
    disk_size_gb    = 15

    resource_manager_tags = {
      "tagKeys/${google_tags_tag_key.key.name}" = "tagValues/${google_tags_tag_value.value.name}"
    }
  }

  deletion_protection = false
  network    = "%[4]s"
  subnetwork    = "%[5]s"

  depends_on = [time_sleep.wait_120_seconds]
}
`, projectID, randomSuffix, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withAutopilotResourceManagerTags(projectID, clusterName, networkName, subnetworkName, randomSuffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%[1]s"
}

resource "google_project_iam_member" "tagHoldAdmin" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagHoldAdmin"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "tagUser1" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "google_project_iam_member" "tagUser2" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:${data.google_project.project.number}@cloudservices.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"

  depends_on = [
    google_project_iam_member.tagHoldAdmin,
    google_project_iam_member.tagUser1,
    google_project_iam_member.tagUser2,
  ]
}

resource "google_tags_tag_key" "key1" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz1-%[2]s"
  description = "For foo/bar1 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [google_compute_network.container_network]
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "foo1-%[2]s"
  description = "For foo1 resources"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz2-%[2]s"
  description = "For foo/bar2 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [
    google_compute_network.container_network,
    google_tags_tag_key.key1
  ]
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "foo2-%[2]s"
  description = "For foo2 resources"
}

resource "google_compute_network" "container_network" {
  name                    = "%[4]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = "%[5]s"
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

data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autopilot" {
  name = "%[3]s"
  location = "us-central1"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["REGULAR"]
  enable_autopilot = true

  deletion_protection = false
  network       = google_compute_network.container_network.name
  subnetwork    = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }

  node_pool_auto_config {
    resource_manager_tags = {
      "tagKeys/${google_tags_tag_key.key1.name}" = "tagValues/${google_tags_tag_value.value1.name}"
	}
  }

  addons_config {
    horizontal_pod_autoscaling {
      disabled = false
    }
  }
  vertical_pod_autoscaling {
    enabled = true
  }

  depends_on = [time_sleep.wait_120_seconds]
}
`, projectID, randomSuffix, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withAutopilotResourceManagerTagsUpdate1(projectID, clusterName, networkName, subnetworkName, randomSuffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%[1]s"
}

resource "google_project_iam_member" "tagHoldAdmin" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagHoldAdmin"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "tagUser1" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "google_project_iam_member" "tagUser2" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:${data.google_project.project.number}@cloudservices.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"

  depends_on = [
    google_project_iam_member.tagHoldAdmin,
    google_project_iam_member.tagUser1,
    google_project_iam_member.tagUser2,
  ]
}

resource "google_tags_tag_key" "key1" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz1-%[2]s"
  description = "For foo/bar1 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [google_compute_network.container_network]
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "foo1-%[2]s"
  description = "For foo1 resources"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz2-%[2]s"
  description = "For foo/bar2 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [
    google_compute_network.container_network,
    google_tags_tag_key.key1
  ]
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "foo2-%[2]s"
  description = "For foo2 resources"
}

resource "google_compute_network" "container_network" {
  name                    = "%[4]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = "%[5]s"
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

data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autopilot" {
  name = "%[3]s"
  location = "us-central1"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["REGULAR"]
  enable_autopilot = true

  deletion_protection = false
  network       = google_compute_network.container_network.name
  subnetwork    = google_compute_subnetwork.container_subnetwork.name
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }

  node_pool_auto_config {
    resource_manager_tags = {
      "tagKeys/${google_tags_tag_key.key1.name}" = "tagValues/${google_tags_tag_value.value1.name}"
      "tagKeys/${google_tags_tag_key.key2.name}" = "tagValues/${google_tags_tag_value.value2.name}"
	}
  }

  addons_config {
    horizontal_pod_autoscaling {
      disabled = false
    }
  }
  vertical_pod_autoscaling {
    enabled = true
  }

  depends_on = [time_sleep.wait_120_seconds]
}
`, projectID, randomSuffix, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withAutopilotResourceManagerTagsUpdate2(projectID, clusterName, networkName, subnetworkName, randomSuffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%[1]s"
}

resource "google_project_iam_member" "tagHoldAdmin" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagHoldAdmin"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "tagUser1" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:service-${data.google_project.project.number}@container-engine-robot.iam.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "google_project_iam_member" "tagUser2" {
  project = "%[1]s"
  role    = "roles/resourcemanager.tagUser"
  member = "serviceAccount:${data.google_project.project.number}@cloudservices.gserviceaccount.com"

  depends_on = [google_project_iam_member.tagHoldAdmin]
}

resource "time_sleep" "wait_120_seconds" {
  create_duration = "120s"

  depends_on = [
    google_project_iam_member.tagHoldAdmin,
    google_project_iam_member.tagUser1,
    google_project_iam_member.tagUser2,
  ]
}

resource "google_tags_tag_key" "key1" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz1-%[2]s"
  description = "For foo/bar1 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [google_compute_network.container_network]
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "foo1-%[2]s"
  description = "For foo1 resources"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%[1]s"
  short_name = "foobarbaz2-%[2]s"
  description = "For foo/bar2 resources"
  purpose = "GCE_FIREWALL"
  purpose_data = {
    network = "%[1]s/%[4]s"
  }

  depends_on = [
    google_compute_network.container_network,
    google_tags_tag_key.key1
  ]
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "foo2-%[2]s"
  description = "For foo2 resources"
}

resource "google_compute_network" "container_network" {
  name                    = "%[4]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = "%[5]s"
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

data "google_container_engine_versions" "uscentral1a" {
  location = "us-central1-a"
}

resource "google_container_cluster" "with_autopilot" {
  name = "%[3]s"
  location = "us-central1"
  min_master_version = data.google_container_engine_versions.uscentral1a.release_channel_latest_version["REGULAR"]
  enable_autopilot = true

  deletion_protection = false
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

  depends_on = [time_sleep.wait_120_seconds]
}
`, projectID, randomSuffix, clusterName, networkName, subnetworkName)
}

func TestAccContainerCluster_privateRegistry(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	nodePoolName := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	secretID := fmt.Sprintf("tf-test-secret-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_privateRegistryEnabled(secretID, clusterName, networkName, subnetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_cluster.primary",
						"node_pool_defaults.0.node_config_defaults.0.containerd_config.0.private_registry_access_config.0.enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"google_container_cluster.primary",
						"node_pool_defaults.0.node_config_defaults.0.containerd_config.0.private_registry_access_config.0.certificate_authority_domain_config.#",
						"2",
					),
					// First CA config
					resource.TestCheckResourceAttr(
						"google_container_cluster.primary",
						"node_pool_defaults.0.node_config_defaults.0.containerd_config.0.private_registry_access_config.0.certificate_authority_domain_config.0.fqdns.0",
						"my.custom.domain",
					),
					// Second CA config
					resource.TestCheckResourceAttr(
						"google_container_cluster.primary",
						"node_pool_defaults.0.node_config_defaults.0.containerd_config.0.private_registry_access_config.0.certificate_authority_domain_config.1.fqdns.0",
						"10.1.2.32",
					),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_privateRegistryDisabled(clusterName, networkName, subnetworkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_container_cluster.primary",
						"node_pool_defaults.0.node_config_defaults.0.containerd_config.0.private_registry_access_config.0.enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"google_container_cluster.primary",
						"node_pool_defaults.0.node_config_defaults.0.containerd_config.0.private_registry_access_config.0.certificate_authority_domain_config.#",
						"0",
					),
				),
			},
			{
				Config: testAccContainerCluster_withNodePoolPrivateRegistry(secretID, clusterName, nodePoolName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withNodeConfigPrivateRegistry(secretID, clusterName, networkName, subnetworkName),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_privateRegistryEnabled(secretID, clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_project" "test_project" {
	}

resource "google_secret_manager_secret" "secret-basic" {
	secret_id     = "%s"
	replication {
		user_managed {
		replicas {
			location = "us-central1"
		}
		}
	}
}

resource "google_secret_manager_secret_version" "secret-version-basic" {
	secret = google_secret_manager_secret.secret-basic.id
	secret_data = "dummypassword"
  }

resource "google_secret_manager_secret_iam_member" "secret_iam" {
	secret_id  = google_secret_manager_secret.secret-basic.id
	role       = "roles/secretmanager.admin"
	member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
	depends_on = [google_secret_manager_secret_version.secret-version-basic]
  }

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
  node_pool_defaults {
    node_config_defaults {
      containerd_config {
        private_registry_access_config {
          enabled = true
          certificate_authority_domain_config {
            fqdns = [ "my.custom.domain" ]
            gcp_secret_manager_certificate_config {
              secret_uri = google_secret_manager_secret_version.secret-version-basic.name
            }
          }
		  certificate_authority_domain_config {
            fqdns = [ "10.1.2.32" ]
            gcp_secret_manager_certificate_config {
              secret_uri = google_secret_manager_secret_version.secret-version-basic.name
            }
          }
        }
      }
    }
  }
}
`, secretID, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_privateRegistryDisabled(clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"

  node_pool_defaults {
    node_config_defaults {
      containerd_config {
        private_registry_access_config {
          enabled = false
        }
      }
    }
  }
}
`, clusterName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodePoolPrivateRegistry(secretID, clusterName, nodePoolName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_project" "test_project" {
	}

resource "google_secret_manager_secret" "secret-basic" {
	secret_id     = "%s"
	replication {
		user_managed {
		replicas {
			location = "us-central1"
		}
		}
	}
}
resource "google_secret_manager_secret_version" "secret-version-basic" {
	secret = google_secret_manager_secret.secret-basic.id
	secret_data = "dummypassword"
  }

resource "google_secret_manager_secret_iam_member" "secret_iam" {
	secret_id  = google_secret_manager_secret.secret-basic.id
	role       = "roles/secretmanager.admin"
	member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
	depends_on = [google_secret_manager_secret_version.secret-version-basic]
  }
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"

  node_pool {
	name               = "%s"
	initial_node_count = 1
    node_config {
		oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
	machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
    containerd_config {
    	private_registry_access_config {
			enabled = true
			certificate_authority_domain_config {
			  fqdns = [ "my.custom.domain", "10.0.0.127:8888" ]
			  gcp_secret_manager_certificate_config {
				secret_uri = google_secret_manager_secret_version.secret-version-basic.name
			}
		}
    }
    }
}
}
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, secretID, clusterName, nodePoolName, networkName, subnetworkName)
}

func testAccContainerCluster_withNodeConfigPrivateRegistry(secretID, clusterName, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
data "google_project" "test_project" {
	}

resource "google_secret_manager_secret" "secret-basic" {
	secret_id     = "%s"
	replication {
		user_managed {
		replicas {
			location = "us-central1"
		}
		}
	}
}
resource "google_secret_manager_secret_version" "secret-version-basic" {
	secret = google_secret_manager_secret.secret-basic.id
	secret_data = "dummypassword"
  }

resource "google_secret_manager_secret_iam_member" "secret_iam" {
	secret_id  = google_secret_manager_secret.secret-basic.id
	role       = "roles/secretmanager.admin"
	member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
	depends_on = [google_secret_manager_secret_version.secret-version-basic]
  }
resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1

  node_config {
	  oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
    machine_type = "n1-standard-8"
    image_type = "COS_CONTAINERD"
    containerd_config {
    	private_registry_access_config {
			enabled = true
			certificate_authority_domain_config {
			  fqdns = [ "my.custom.domain", "10.0.0.127:8888" ]
			  gcp_secret_manager_certificate_config {
				secret_uri = google_secret_manager_secret_version.secret-version-basic.name
			}
		}
    }
    }
}
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, secretID, clusterName, networkName, subnetworkName)
}

func TestAccContainerCluster_withProviderDefaultLabels(t *testing.T) {
	// The test failed if VCR testing is enabled, because the cached provider config is used.
	// With the cached provider config, any changes in the provider default labels will not be applied.
	acctest.SkipIfVcr(t)
	t.Parallel()

	clusterName := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withProviderDefaultLabels(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "resource_labels.%", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "resource_labels.created-by", "terraform"),

					// goog-terraform-provisioned: true is added
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.created-by", "terraform"),

					resource.TestCheckResourceAttr("google_container_cluster.primary", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection", "resource_labels", "terraform_labels"},
			},
			{
				Config: testAccContainerCluster_resourceLabelsOverridesProviderDefaultLabels(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "resource_labels.%", "2"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "resource_labels.created-by", "terraform"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.default_key1", "value1"),

					// goog-terraform-provisioned: true is added
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.default_key1", "value1"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.created-by", "terraform"),

					resource.TestCheckResourceAttr("google_container_cluster.primary", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection", "resource_labels", "terraform_labels"},
			},
			{
				Config: testAccContainerCluster_moveResourceLabelToProviderDefaultLabels(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "resource_labels.%", "0"),

					// goog-terraform-provisioned: true is added
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.%", "3"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.default_key1", "default_value1"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.created-by", "terraform"),

					resource.TestCheckResourceAttr("google_container_cluster.primary", "effective_labels.%", "3"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection", "resource_labels", "terraform_labels"},
			},
			{
				Config: testAccContainerCluster_basic(clusterName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.primary", "resource_labels.%", "0"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "terraform_labels.%", "1"),
					resource.TestCheckResourceAttr("google_container_cluster.primary", "effective_labels.%", "1"),
				),
			},
			{
				ResourceName:            "google_container_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_default_node_pool", "deletion_protection", "resource_labels", "terraform_labels"},
			},
		},
	})
}

func testAccContainerCluster_withProviderDefaultLabels(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
  resource_labels = {
    created-by = "terraform"
  }
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_resourceLabelsOverridesProviderDefaultLabels(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
  resource_labels = {
    created-by = "terraform"
	default_key1 = "value1"
  }
}
`, name, networkName, subnetworkName)
}

func testAccContainerCluster_moveResourceLabelToProviderDefaultLabels(name, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
	created-by   = "terraform"
  }
}

resource "google_container_cluster" "primary" {
  name               = "%s"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
}
`, name, networkName, subnetworkName)
}

func TestAccContainerCluster_storagePoolsWithNodePool(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	np := fmt.Sprintf("tf-test-nodepool-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	location := envvar.GetTestZoneFromEnv()

	storagePoolResourceName := acctest.BootstrapComputeStoragePool(t, "basic-1", "hyperdisk-balanced")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_storagePoolsWithNodePool(cluster, location, networkName, subnetworkName, np, storagePoolResourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.storage_pools_with_node_pool", "node_pool.0.node_config.0.storage_pools.0", storagePoolResourceName),
				),
			},
			{
				ResourceName:            "google_container_cluster.storage_pools_with_node_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_storagePoolsWithNodePool(cluster, location, networkName, subnetworkName, np, storagePoolResourceName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "storage_pools_with_node_pool" {
  name               = "%s"
  location           = "%s"
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
  node_pool {
    name = "%s"
    initial_node_count = 1
    node_config {
      machine_type = "c3-standard-4"
      image_type = "COS_CONTAINERD"
      storage_pools = ["%s"]
	  disk_type = "hyperdisk-balanced"
    }
  }
}
`, cluster, location, networkName, subnetworkName, np, storagePoolResourceName)
}

func TestAccContainerCluster_storagePoolsWithNodeConfig(t *testing.T) {
	t.Parallel()

	cluster := fmt.Sprintf("tf-test-cluster-%s", acctest.RandString(t, 10))
	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)
	location := envvar.GetTestZoneFromEnv()

	storagePoolResourceName := acctest.BootstrapComputeStoragePool(t, "basic-1", "hyperdisk-balanced")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_storagePoolsWithNodeConfig(cluster, location, networkName, subnetworkName, storagePoolResourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.storage_pools_with_node_config", "node_config.0.storage_pools.0", storagePoolResourceName),
				),
			},
			{
				ResourceName:            "google_container_cluster.storage_pools_with_node_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_storagePoolsWithNodeConfig(cluster, location, networkName, subnetworkName, storagePoolResourceName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "storage_pools_with_node_config" {
  name               = "%s"
  location           = "%s"
  initial_node_count = 1
  deletion_protection = false
  network    = "%s"
  subnetwork    = "%s"
  node_config {
    machine_type = "c3-standard-4"
    image_type = "COS_CONTAINERD"
    storage_pools = ["%s"]
	disk_type = "hyperdisk-balanced"
  }
}
`, cluster, location, networkName, subnetworkName, storagePoolResourceName)
}

func TestAccContainerCluster_withAutopilotGcpFilestoreCsiDriver(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	clusterName := fmt.Sprintf("tf-test-cluster-%s", randomSuffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckContainerClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerCluster_withAutopilotGcpFilestoreCsiDriverDefault(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_container_cluster.with_autopilot_gcp_filestore", "addons_config.0.gcp_filestore_csi_driver_config.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot_gcp_filestore",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccContainerCluster_withAutopilotGcpFilestoreCsiDriverUpdated(clusterName),
			},
			{
				ResourceName:            "google_container_cluster.with_autopilot_gcp_filestore",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccContainerCluster_withAutopilotGcpFilestoreCsiDriverDefault(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_autopilot_gcp_filestore" {
  name                = "%s"
  location            = "us-central1"
  enable_autopilot    = true
  deletion_protection = false
}
`, name)
}

func testAccContainerCluster_withAutopilotGcpFilestoreCsiDriverUpdated(name string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "with_autopilot_gcp_filestore" {
  name                = "%s"
  location            = "us-central1"
  enable_autopilot    = true
  deletion_protection = false

  addons_config {
    gcp_filestore_csi_driver_config {
      enabled = false
    }
  }
}
`, name)
}
