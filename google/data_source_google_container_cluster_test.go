package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccContainerClusterDatasource_zonal(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_zonal(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleContainerClusterCheck("data.google_container_cluster.kubes", "google_container_cluster.kubes"),
				),
			},
		},
	})
}

func TestAccContainerClusterDatasource_regional(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_regional(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleContainerClusterCheck("data.google_container_cluster.kubes", "google_container_cluster.kubes"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleContainerClusterCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		clusterAttrToCheck := []string{
			"name",
			"zone",
			"additional_zones",
			"addons_config",
			"cluster_ipv4_cidr",
			"description",
			"enable_kubernetes_alpha",
			"enable_tpu",
			"enable_legacy_abac",
			"endpoint",
			"enable_legacy_abac",
			"instance_group_urls",
			"ip_allocation_policy",
			"logging_service",
			"maintenance_policy",
			"master_auth",
			"master_auth.0.password",
			"master_auth.0.username",
			"master_auth.0.client_certificate_config.0.issue_client_certificate",
			"master_auth.0.client_certificate",
			"master_auth.0.client_key",
			"master_auth.0.cluster_ca_certificate",
			"master_authorized_networks_config",
			"master_version",
			"min_master_version",
			"monitoring_service",
			"network",
			"network_policy",
			"node_version",
			"subnetwork",
		}

		for _, attr := range clusterAttrToCheck {
			if dsAttr[attr] != rsAttr[attr] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr,
					dsAttr[attr],
					rsAttr[attr],
				)
			}
		}

		return nil
	}
}

func testAccContainerClusterDatasource_zonal() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
	name               = "cluster-test-%s"
	zone               = "us-central1-a"
	initial_node_count = 1
	
	master_auth {
		username = "mr.yoda"
		password = "adoy.rm.123456789"
	}
}
	
data "google_container_cluster" "kubes" {
	name = "${google_container_cluster.kubes.name}"
	zone = "${google_container_cluster.kubes.zone}"
}
`, acctest.RandString(10))
}

func testAccContainerClusterDatasource_regional() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
	name               = "cluster-test-%s"
	region             = "us-central1"
	initial_node_count = 1
}
	
data "google_container_cluster" "kubes" {
	name   = "${google_container_cluster.kubes.name}"
	region = "${google_container_cluster.kubes.region}"
}
`, acctest.RandString(10))
}
