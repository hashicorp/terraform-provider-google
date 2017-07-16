package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/stretchr/testify/assert"
	"strconv"
)

const base10 = 10

func TestExtractLastResourceFromUri_withUrl(t *testing.T) {
	r := extractLastResourceFromUri("http://something.com/one/two/three")
	assert.Equal(t, "three", r)
}

func TestExtractLastResourceFromUri_WithStaticValue(t *testing.T) {
	r := extractLastResourceFromUri("three")
	assert.Equal(t, "three", r)
}

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
						"google_dataproc_cluster.basic"),
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
						"google_dataproc_cluster.with_timeout"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withMasterConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withMasterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_master_config"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withBucketRef(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withBucket,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_bucket"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withInitAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withInitAction(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_init_action"),
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

func TestAccDataprocCluster_withServiceAcc(t *testing.T) {

	saEmail := "TODO-compute@developer.gserviceaccount.com"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withServiceAcc(saEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_service_account"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withImageVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withImageVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocCluster(
						"google_dataproc_cluster.with_image_version"),
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

			{"bucket", cluster.Config.ConfigBucket},
			{"image_version", cluster.Config.SoftwareConfig.ImageVersion},
			{"zone", extractLastResourceFromUri(cluster.Config.GceClusterConfig.ZoneUri)},

			{"subnetwork", extractLastResourceFromUri(cluster.Config.GceClusterConfig.SubnetworkUri)},
			{"service_account", cluster.Config.GceClusterConfig.ServiceAccount},
			{"service_account_scopes", cluster.Config.GceClusterConfig.ServiceAccountScopes},
			{"metadata", cluster.Config.GceClusterConfig.Metadata},
			{"labels", cluster.Labels},
			{"tags", cluster.Config.GceClusterConfig.Tags},

			{"master_config.0.num_masters", strconv.FormatInt(cluster.Config.MasterConfig.NumInstances, base10)},
			{"master_config.0.boot_disk_size_gb", strconv.FormatInt(cluster.Config.MasterConfig.DiskConfig.BootDiskSizeGb, base10)},
			{"master_config.0.num_local_ssds", strconv.FormatInt(cluster.Config.MasterConfig.DiskConfig.NumLocalSsds, base10)},
			{"master_config.0.machine_type", extractLastResourceFromUri(cluster.Config.MasterConfig.MachineTypeUri)},

			{"worker_config.0.num_workers", strconv.FormatInt(cluster.Config.WorkerConfig.NumInstances, base10)},
			{"worker_config.0.boot_disk_size_gb", strconv.FormatInt(cluster.Config.WorkerConfig.DiskConfig.BootDiskSizeGb, base10)},
			{"worker_config.0.num_local_ssds", strconv.FormatInt(cluster.Config.WorkerConfig.DiskConfig.NumLocalSsds, base10)},
			{"worker_config.0.machine_type", extractLastResourceFromUri(cluster.Config.WorkerConfig.MachineTypeUri)},
		}

		extracted := false
		if len(cluster.Config.InitializationActions) > 0 {
			actions := []string{}
			for _, v := range cluster.Config.InitializationActions {
				actions = append(actions, v.ExecutableFile)

				if !extracted {
					tsec := v.ExecutionTimeout[:len(v.ExecutionTimeout)-2]
					tsecI, err := strconv.Atoi(tsec)
					if err != nil {
						return err
					}
					clusterTests = append(clusterTests, clusterTestField{"initialization_action_timeout_sec", tsecI})
					extracted = true
				}
			}
			clusterTests = append(clusterTests, clusterTestField{"initialization_actions", actions})
		} else {
			clusterTests = append(clusterTests, clusterTestField{"initialization_actions", []string{}})
		}

		if cluster.Config.SecondaryWorkerConfig != nil {
			clusterTests = append(clusterTests,
				clusterTestField{"worker_config.0.preemptible_num_workers", strconv.FormatInt(cluster.Config.SecondaryWorkerConfig.NumInstances, base10)},
				clusterTestField{"worker_config.0.preemptible_boot_disk_size_gb", strconv.FormatInt(cluster.Config.SecondaryWorkerConfig.DiskConfig.BootDiskSizeGb, base10)})
		}

		for _, attrs := range clusterTests {
			if c := checkMatch(attributes, attrs.tf_attr, attrs.gcp_attr); c != "" {
				return fmt.Errorf(c)
			}
		}

		// A few attributes need to be done separately in order to normalise them.
		// Network
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
resource "google_dataproc_cluster" "basic" {
	name = "cluster-test-%s"
	zone = "us-central1-a"
}`, acctest.RandString(10))

var testAccDataprocCluster_withTimeout = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_timeout" {
	name = "cluster-test-%s"
	zone = "us-central1-a"

	timeouts {
		create = "30m"
		delete = "30m"
		update = "30m"
	}
}`, acctest.RandString(10))

var testAccDataprocCluster_withMasterConfig = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_master_config" {
	name = "cluster-test-%s"
	zone = "us-central1-f"

	master_config {
		num_masters       = 1
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
		num_local_ssds    = 1
	}
}`, acctest.RandString(10))

func testAccDataprocCluster_withInitAction(rnd string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "init_bucket" {
    name          = "tf-dataproc-acctest-%s"
    force_destroy = "true"
}

resource "google_storage_bucket_object" "init_script" {
  name    = "tf-acctest-init-script-%s.sh"
  bucket  = "${google_storage_bucket.init_bucket.name}"
  content = <<EOL
#!/bin/bash
ROLE=$$(/usr/share/google/get_metadata_value attributes/dataproc-role)
if [[ "$${ROLE}" == 'Master' ]]; then
  echo "on the master" >> /tmp/msg.txt
else
  echo "on the worker" >> /tmp/msg.txt
fi
EOL

}

resource "google_dataproc_cluster" "with_init_action" {
	name   = "cluster-test-%s"
	zone   = "us-central1-f"
	initialization_action_timeout_sec = 500
	initialization_actions = [
	   "${google_storage_bucket.init_bucket.url}/${google_storage_bucket_object.init_script.name}"
	]

	worker_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}
}`, rnd, rnd, rnd)
}

var testAccDataprocCluster_withBucket = fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
    name          = "tf-dataproc-bucket-%s"
    force_destroy = "true"
}

resource "google_dataproc_cluster" "with_bucket" {
	name   = "cluster-test-%s"
	zone   = "us-central1-f"
	bucket = "${google_storage_bucket.bucket.name}"

	worker_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}
}`, acctest.RandString(10), acctest.RandString(10))

var testAccDataprocCluster_withWorkerConfig = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_worker_config" {
	name = "cluster-test-%s"
	zone = "us-central1-f"

	worker_config {
		num_workers       = 2
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
		num_local_ssds    = 1

		preemptible_num_workers = 1
		preemptible_boot_disk_size_gb = 10
	}
}`, acctest.RandString(10))

var testAccDataprocCluster_withImageVersion = fmt.Sprintf(`
resource "google_dataproc_cluster" "with_image_version" {
	name = "cluster-test-%s"
	zone = "us-central1-f"
	image_version = "1.0.44"
}`, acctest.RandString(10))

func testAccDataprocCluster_withServiceAcc(saEmail string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_service_account" {
	name = "cluster-test-%s"
	zone = "us-central1-f"

    master_config {
        machine_type = "n1-standard-1"
        boot_disk_size_gb = 10
    }

	service_account = "%s"

	service_account_scopes = [
        #    The following scopes necessary for the cluster to function properly are
		#	always added, even if not explicitly specified:
		#		useraccounts-ro: https://www.googleapis.com/auth/cloud.useraccounts.readonly
		#		storage-rw:      https://www.googleapis.com/auth/devstorage.read_write
		#		logging-write:   https://www.googleapis.com/auth/logging.write
        #
		#	So user is expected to add these explicitly (in this order) otherwise terraform
		#   will think there is a change to resource
		"useraccounts-ro","storage-rw","logging-write",

	    # Additional ones specifically desired by user (Note for now must be in alpha order
	    # of fully qualified scope name)
	    "monitoring"

	]

	worker_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}
}`, acctest.RandString(10), saEmail)
}

var rndDPName = acctest.RandString(10)
var testAccDataprocCluster_networkRef = fmt.Sprintf(`
resource "google_compute_network" "dataproc_network" {
	name = "dataproc-net-%s"
	auto_create_subnetworks = true
}

resource "google_compute_firewall" "dataproc_network_firewall" {
	name = "dataproc-net-%s-allow-internal"
	description = "Firewall rules for dataproc Terraform acceptance testing"
	network = "${google_compute_network.dataproc_network.name}"

	allow {
	    protocol = "icmp"
	}

	allow {
		protocol = "tcp"
		ports    = ["0-65535"]
	}

	allow {
		protocol = "udp"
		ports    = ["0-65535"]
	}
}

resource "google_dataproc_cluster" "with_net_ref_by_name" {
	name = "cluster-test-%s-name"
	zone = "us-central1-a"
	depends_on = ["google_compute_firewall.dataproc_network_firewall"]

    # to minimise cost for tests, using smaller instances
    master_config {
        machine_type = "n1-standard-1"
        boot_disk_size_gb = 10
    }

    worker_config {
        machine_type = "n1-standard-1"
        boot_disk_size_gb = 10
    }

	network = "${google_compute_network.dataproc_network.name}"
}

resource "google_dataproc_cluster" "with_net_ref_by_url" {
	name = "cluster-test-%s-url"
	zone = "us-central1-a"
    depends_on = ["google_compute_firewall.dataproc_network_firewall"]

    # to minimise cost for tests, using smaller instances
    master_config {
        machine_type = "n1-standard-1"
        boot_disk_size_gb = 10
    }

    worker_config {
        machine_type = "n1-standard-1"
        boot_disk_size_gb = 10
    }

	network = "${google_compute_network.dataproc_network.self_link}"
}

`, rndDPName, rndDPName, rndDPName, rndDPName)
