package google

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccContainerNodePool_basic(t *testing.T) {
	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_basic(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
			},
		},
	})
}

func TestAccContainerNodePool_withNodeConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_withNodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np_with_node_config"),
				),
			},
		},
	})
}

func TestAccContainerNodePool_withNodeConfigScopeAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_withNodeConfigScopeAlias,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np_with_node_config_scope_alias"),
				),
			},
		},
	})
}

func TestAccContainerNodePool_autoscaling(t *testing.T) {
	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerNodePool_autoscaling(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
			},
			resource.TestStep{
				Config: testAccContainerNodePool_updateAutoscaling(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
			},
			resource.TestStep{
				Config: testAccContainerNodePool_basic(cluster, np),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerNodePoolMatches("google_container_node_pool.np"),
				),
			},
		},
	})
}

func testAccCheckContainerNodePoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_node_pool" {
			continue
		}

		attributes := rs.Primary.Attributes
		_, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
			config.Project, attributes["zone"], attributes["cluster"], attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("NodePool still exists")
		}
	}

	return nil
}

func testAccCheckContainerNodePoolMatches(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		attributes := rs.Primary.Attributes
		nodepool, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
			config.Project, attributes["zone"], attributes["cluster"], attributes["name"]).Do()
		if err != nil {
			return err
		}

		if nodepool.Name != attributes["name"] {
			return fmt.Errorf("NodePool not found")
		}

		type nodepoolTestField struct {
			tfAttr  string
			gcpAttr interface{}
		}

		nodepoolTests := []nodepoolTestField{
			{"initial_node_count", strconv.FormatInt(nodepool.InitialNodeCount, 10)},
			{"node_config.0.machine_type", nodepool.Config.MachineType},
			{"node_config.0.disk_size_gb", strconv.FormatInt(nodepool.Config.DiskSizeGb, 10)},
			{"node_config.0.local_ssd_count", strconv.FormatInt(nodepool.Config.LocalSsdCount, 10)},
			{"node_config.0.oauth_scopes", nodepool.Config.OauthScopes},
			{"node_config.0.service_account", nodepool.Config.ServiceAccount},
			{"node_config.0.metadata", nodepool.Config.Metadata},
			{"node_config.0.image_type", nodepool.Config.ImageType},
			{"node_config.0.labels", nodepool.Config.Labels},
			{"node_config.0.tags", nodepool.Config.Tags},
			{"node_config.0.preemptible", nodepool.Config.Preemptible},
		}

		for _, attrs := range nodepoolTests {
			if c := nodepoolCheckMatch(attributes, attrs.tfAttr, attrs.gcpAttr); c != "" {
				return fmt.Errorf(c)
			}
		}

		tfAS := attributes["autoscaling.#"] == "1"
		if gcpAS := nodepool.Autoscaling != nil && nodepool.Autoscaling.Enabled == true; tfAS != gcpAS {
			return fmt.Errorf("Mismatched autoscaling status. TF State: %t. GCP State: %t", tfAS, gcpAS)
		}
		if tfAS {
			if tf := attributes["autoscaling.0.min_node_count"]; strconv.FormatInt(nodepool.Autoscaling.MinNodeCount, 10) != tf {
				return fmt.Errorf("Mismatched Autoscaling.MinNodeCount. TF State: %s. GCP State: %d",
					tf, nodepool.Autoscaling.MinNodeCount)
			}

			if tf := attributes["autoscaling.0.max_node_count"]; strconv.FormatInt(nodepool.Autoscaling.MaxNodeCount, 10) != tf {
				return fmt.Errorf("Mismatched Autoscaling.MaxNodeCount. TF State: %s. GCP State: %d",
					tf, nodepool.Autoscaling.MaxNodeCount)
			}

		}

		return nil
	}
}

func testAccContainerNodePool_basic(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
}`, cluster, np)
}

func testAccContainerNodePool_autoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
	autoscaling {
		min_node_count = 1
		max_node_count = 3
	}
}`, cluster, np)
}

func testAccContainerNodePool_updateAutoscaling(cluster, np string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "%s"
	zone = "us-central1-a"
	initial_node_count = 3

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}

resource "google_container_node_pool" "np" {
	name = "%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 2
	autoscaling {
		min_node_count = 1
		max_node_count = 5
	}
}`, cluster, np)
}

func nodepoolCheckMatch(attributes map[string]string, attr string, gcp interface{}) string {
	if gcpList, ok := gcp.([]string); ok {
		return nodepoolCheckListMatch(attributes, attr, gcpList)
	}
	if gcpMap, ok := gcp.(map[string]string); ok {
		return nodepoolCheckMapMatch(attributes, attr, gcpMap)
	}
	if gcpBool, ok := gcp.(bool); ok {
		return checkBoolMatch(attributes, attr, gcpBool)
	}
	tf := attributes[attr]
	if tf != gcp {
		return nodepoolMatchError(attr, tf, gcp)
	}
	return ""
}

func nodepoolCheckSetMatch(attributes map[string]string, attr string, gcpList []string) string {
	num, err := strconv.Atoi(attributes[attr+".#"])
	if err != nil {
		return fmt.Sprintf("Error in number conversion for attribute %s: %s", attr, err)
	}
	if num != len(gcpList) {
		return fmt.Sprintf("NodePool has mismatched %s size.\nTF Size: %d\nGCP Size: %d", attr, num, len(gcpList))
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
	return nodepoolMatchError(attr, tfAttr, gcpList)
}

func nodepoolCheckListMatch(attributes map[string]string, attr string, gcpList []string) string {
	num, err := strconv.Atoi(attributes[attr+".#"])
	if err != nil {
		return fmt.Sprintf("Error in number conversion for attribute %s: %s", attr, err)
	}
	if num != len(gcpList) {
		return fmt.Sprintf("NodePool has mismatched %s size.\nTF Size: %d\nGCP Size: %d", attr, num, len(gcpList))
	}

	for i, gcp := range gcpList {
		if tf := attributes[fmt.Sprintf("%s.%d", attr, i)]; tf != gcp {
			return nodepoolMatchError(fmt.Sprintf("%s[%d]", attr, i), tf, gcp)
		}
	}

	return ""
}

func nodepoolCheckMapMatch(attributes map[string]string, attr string, gcpMap map[string]string) string {
	num, err := strconv.Atoi(attributes[attr+".%"])
	if err != nil {
		return fmt.Sprintf("Error in number conversion for attribute %s: %s", attr, err)
	}
	if num != len(gcpMap) {
		return fmt.Sprintf("NodePool has mismatched %s size.\nTF Size: %d\nGCP Size: %d", attr, num, len(gcpMap))
	}

	for k, gcp := range gcpMap {
		if tf := attributes[fmt.Sprintf("%s.%s", attr, k)]; tf != gcp {
			return nodepoolMatchError(fmt.Sprintf("%s[%s]", attr, k), tf, gcp)
		}
	}

	return ""
}

func nodepoolMatchError(attr, tf interface{}, gcp interface{}) string {
	return fmt.Sprintf("NodePool has mismatched %s.\nTF State: %+v\nGCP State: %+v", attr, tf, gcp)
}

var testAccContainerNodePool_withNodeConfig = fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1
	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}
resource "google_container_node_pool" "np_with_node_config" {
	name = "tf-nodepool-test-%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1
	node_config {
		machine_type = "g1-small"
		disk_size_gb = 10
		oauth_scopes = [
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring"
		]
		preemptible = true
	}
}`, acctest.RandString(10), acctest.RandString(10))

var testAccContainerNodePool_withNodeConfigScopeAlias = fmt.Sprintf(`
resource "google_container_cluster" "cluster" {
	name = "tf-cluster-nodepool-test-%s"
	zone = "us-central1-a"
	initial_node_count = 1
	master_auth {
		username = "mr.yoda"
		password = "adoy.rm"
	}
}
resource "google_container_node_pool" "np_with_node_config_scope_alias" {
	name = "tf-nodepool-test-%s"
	zone = "us-central1-a"
	cluster = "${google_container_cluster.cluster.name}"
	initial_node_count = 1
	node_config {
		machine_type = "g1-small"
		disk_size_gb = 10
		oauth_scopes = ["compute-rw", "storage-ro", "logging-write", "monitoring"]
	}
}`, acctest.RandString(10), acctest.RandString(10))
