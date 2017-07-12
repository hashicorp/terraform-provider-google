package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataprocCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.primary"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withTimeout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.primary"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withWorkerConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withWorkerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_worker_config"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withWorkerConfigScopeAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withWorkerConfigScopeAlias,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_worker_config_scope_alias"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_network(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_networkRef,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_net_ref_by_url"),
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_net_ref_by_name"),
				),
			},
		},
	})
}

func testAccCheckDataprocClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_dataproc_cluster" {
			continue
		}

		attributes := rs.Primary.Attributes
		_, err := config.clientDataproc.Projects.Regions.Clusters.Get(
			config.Project, attributes["region"], attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckDataprocCluster(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)
		cluster, err := config.clientDataproc.Projects.Regions.Clusters.Get(
			config.Project, attributes["region"], attributes["name"]).Do()
		if err != nil {
			return err
		}

		if cluster.ClusterName != attributes["name"] {
			return fmt.Errorf("Cluster %s not found, found %s instead", attributes["name"], cluster.ClusterName)
		}

		type clusterTestField struct {
			tf_attr  string
			gcp_attr interface{}
		}

		clusterTests := []clusterTestField{
			{"zone", cluster.Config.GceClusterConfig.ZoneUri},
			{"bucket", cluster.Config.ConfigBucket},
			{"subnetwork", cluster.Config.GceClusterConfig.SubnetworkUri},

			// TODO finish
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
		gcp, err := getNetworkNameFromSelfLink(cluster.Config.GceClusterConfig.NetworkUri)
		if err != nil {
			return err
		}
		if tf != gcp {
			return fmt.Errorf(matchError("network", tf, gcp))
		}

		return nil
	}
}

var testAccDataprocCluster_basic = fmt.Sprintf(`
resource "google_dataproc_cluster" "primary" {
	name = "cluster-test-%s"
	zone = "us-central1-a"

	master_config {
        machine_type = "n1-standard-1"
        boot_disk_size_gb = 10
    }

    worker_config {
        machine_type = "n1-standard-1"
        boot_disk_size_gb = 10
    }
}`, acctest.RandString(10))

var testAccDataprocCluster_withTimeout = fmt.Sprintf(`
resource "google_dataproc_cluster" "primary" {
	name = "cluster-test-%s"
	zone = "us-central1-a"

	timeouts {
		create = "30m"
		delete = "30m"
		update = "30m"
	}
}`, acctest.RandString(10))

var testAccDataprocCluster_withMasterConfig = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_master_auth" {
	name = "cluster-test-%s"
	zone = "us-central1-a"

	master_config {
		num_masters = 1
	}
}`, acctest.RandString(10))

var testAccDataprocCluster_withVersion = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_version" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
	image_version = "1.1"
}`, acctest.RandString(10))

var testAccDataprocCluster_withWorkerConfig = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_worker_config" {
	name = "cluster-test-%s"
	zone = "us-central1-f"

	worker_config {
		num_workers = 3
		machine_type = "n1-standard-1"
		boot_disk_size_gb = 10
		num_local_ssds = 1
		service_account_scopes = [
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring"
		]
		service_account = "default"
	}
}`, acctest.RandString(10))

var testAccDataprocCluster_withWorkerConfigScopeAlias = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_worker_config_scope_alias" {
	name = "cluster-test-%s"
	zone = "us-central1-f"

	worker_config {
		num_workers = 3
		boot_disk_size_gb = "n1-standard-1"
		disk_size_gb = 10
		num_local_ssds = 1
		service_account_scopes = [ "compute-rw", "storage-ro", "logging-write", "monitoring" ]
		service_account = "default"
	}
}`, acctest.RandString(10))

var testAccDataprocCluster_networkRef = fmt.Sprintf(`
resource "google_compute_network" "dataproc_network" {
	name = "dataproc-net-%s"
	auto_create_subnetworks = true
}

resource "google_dataproc_cluster" "with_net_ref_by_url" {
	name = "cluster-test-%s"
	zone = "us-central1-a"

	network = "${google_compute_network.dataproc_network.self_link}"
}

resource "google_dataproc_cluster" "with_net_ref_by_name" {
	name = "cluster-test-%s"
	zone = "us-central1-a"

	network = "${google_compute_network.dataproc_network.name}"
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
