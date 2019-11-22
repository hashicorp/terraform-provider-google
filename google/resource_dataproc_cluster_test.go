package google

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	dataproc "google.golang.org/api/dataproc/v1beta2"
	"google.golang.org/api/googleapi"
)

func TestDataprocExtractInitTimeout(t *testing.T) {
	t.Parallel()

	actual, err := extractInitTimeout("500s")
	expected := 500
	if err != nil {
		t.Fatalf("Expected %d, but got error %v", expected, err)
	}
	if actual != expected {
		t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func TestDataprocExtractInitTimeout_nonSeconds(t *testing.T) {
	t.Parallel()

	actual, err := extractInitTimeout("5m")
	expected := 300
	if err != nil {
		t.Fatalf("Expected %d, but got error %v", expected, err)
	}
	if actual != expected {
		t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func TestDataprocExtractInitTimeout_empty(t *testing.T) {
	t.Parallel()

	_, err := extractInitTimeout("")
	expected := "time: invalid duration"
	if err != nil && err.Error() != expected {
		return
	}
	t.Fatalf("Expected an error with message '%s', but got %v", expected, err.Error())
}

func TestDataprocParseImageVersion(t *testing.T) {
	t.Parallel()

	testCases := map[string]dataprocImageVersion{
		"1.2":             {"1", "2", "", ""},
		"1.2.3":           {"1", "2", "3", ""},
		"1.2.3rc":         {"1", "2", "3rc", ""},
		"1.2-debian9":     {"1", "2", "", "debian9"},
		"1.2.3-debian9":   {"1", "2", "3", "debian9"},
		"1.2.3rc-debian9": {"1", "2", "3rc", "debian9"},
	}

	for v, expected := range testCases {
		actual, err := parseDataprocImageVersion(v)
		if actual.major != expected.major {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if err != nil {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.minor != expected.minor {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.subminor != expected.subminor {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.osName != expected.osName {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
	}

	errorTestCases := []string{
		"",
		"1",
		"notaversion",
		"1-debian",
	}
	for _, v := range errorTestCases {
		if _, err := parseDataprocImageVersion(v); err == nil {
			t.Errorf("expected parsing invalid version %q to return error", v)
		}
	}
}

func TestDataprocDiffSuppress(t *testing.T) {
	t.Parallel()

	doSuppress := [][]string{
		{"1.3.10-debian9", "1.3"},
		{"1.3.10-debian9", "1.3-debian9"},
		{"1.3.10", "1.3"},
		{"1.3-debian9", "1.3"},
	}

	noSuppress := [][]string{
		{"1.3.10-debian9", "1.3.10-ubuntu"},
		{"1.3.10-debian9", "1.3.9-debian9"},
		{"1.3.10-debian9", "1.3-ubuntu"},
		{"1.3.10-debian9", "1.3.9"},
		{"1.3.10-debian9", "1.4"},
		{"1.3.10-debian9", "2.3"},
		{"1.3.10", "1.3.10-debian9"},
		{"1.3", "1.3.10"},
		{"1.3", "1.3.10-debian9"},
		{"1.3", "1.3-debian9"},
	}

	for _, tup := range doSuppress {
		if !dataprocImageVersionDiffSuppress("", tup[0], tup[1], nil) {
			t.Errorf("expected (old: %q, new: %q) to be suppressed", tup[0], tup[1])
		}
	}
	for _, tup := range noSuppress {
		if dataprocImageVersionDiffSuppress("", tup[0], tup[1], nil) {
			t.Errorf("expected (old: %q, new: %q) to not be suppressed", tup[0], tup[1])
		}
	}
}

func TestAccDataprocCluster_missingZoneGlobalRegion1(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataproc_missingZoneGlobalRegion1(rnd),
				ExpectError: regexp.MustCompile("zone is mandatory when region is set to 'global'"),
			},
		},
	})
}

func TestAccDataprocCluster_missingZoneGlobalRegion2(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataproc_missingZoneGlobalRegion2(rnd),
				ExpectError: regexp.MustCompile("zone is mandatory when region is set to 'global'"),
			},
		},
	})
}

func TestAccDataprocCluster_basic(t *testing.T) {
	t.Parallel()

	var cluster dataproc.Cluster
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_basic(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.basic", &cluster),

					// Default behaviour is for Dataproc to autogen or autodiscover a config bucket
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.bucket"),

					// Default behavior is for Dataproc to not use only internal IP addresses
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.internal_ip_only", "false"),

					// Expect 1 master instances with computed values
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.master_config.#", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.master_config.0.num_instances", "1"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.master_config.0.disk_config.0.boot_disk_size_gb"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.master_config.0.disk_config.0.num_local_ssds"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.master_config.0.disk_config.0.boot_disk_type"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.master_config.0.machine_type"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.master_config.0.instance_names.#", "1"),

					// Expect 2 worker instances with computed values
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.worker_config.#", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.worker_config.0.num_instances", "2"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.worker_config.0.disk_config.0.boot_disk_size_gb"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.worker_config.0.disk_config.0.num_local_ssds"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.worker_config.0.disk_config.0.boot_disk_type"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.basic", "cluster_config.0.worker_config.0.machine_type"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.worker_config.0.instance_names.#", "2"),

					// Expect 0 preemptible worker instances
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.preemptible_worker_config.#", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.preemptible_worker_config.0.num_instances", "0"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.preemptible_worker_config.0.instance_names.#", "0"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withAccelerators(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster

	project := getTestProjectFromEnv()
	zone := "us-central1-a"
	acceleratorType := "nvidia-tesla-k80"
	acceleratorLink := fmt.Sprintf("https://www.googleapis.com/compute/beta/projects/%s/zones/%s/acceleratorTypes/%s", project, zone, acceleratorType)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withAccelerators(rnd, zone, acceleratorType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.accelerated_cluster", &cluster),
					testAccCheckDataprocClusterAccelerator(&cluster, 1, acceleratorLink, 1, acceleratorLink),
				),
			},
		},
	})
}

func testAccCheckDataprocClusterAccelerator(cluster *dataproc.Cluster, masterCount int, masterAccelerator string, workerCount int, workerAccelerator string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		master := cluster.Config.MasterConfig.Accelerators
		if len(master) != 1 {
			return fmt.Errorf("Saw %d master accelerator types instead of 1", len(master))
		}

		if int(master[0].AcceleratorCount) != masterCount {
			return fmt.Errorf("Saw %d master accelerators instead of %d", int(master[0].AcceleratorCount), masterCount)
		}

		if master[0].AcceleratorTypeUri != masterAccelerator {
			return fmt.Errorf("Saw %s master accelerator type instead of %s", master[0].AcceleratorTypeUri, masterAccelerator)
		}

		worker := cluster.Config.WorkerConfig.Accelerators
		if len(worker) != 1 {
			return fmt.Errorf("Saw %d worker accelerator types instead of 1", len(worker))
		}

		if int(worker[0].AcceleratorCount) != workerCount {
			return fmt.Errorf("Saw %d worker accelerators instead of %d", int(worker[0].AcceleratorCount), workerCount)
		}

		if worker[0].AcceleratorTypeUri != workerAccelerator {
			return fmt.Errorf("Saw %s worker accelerator type instead of %s", worker[0].AcceleratorTypeUri, workerAccelerator)
		}

		return nil
	}
}

func TestAccDataprocCluster_withInternalIpOnlyTrue(t *testing.T) {
	t.Parallel()

	var cluster dataproc.Cluster
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withInternalIpOnlyTrue(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.basic", &cluster),

					// Testing behavior for Dataproc to use only internal IP addresses
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.internal_ip_only", "true"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withMetadataAndTags(t *testing.T) {
	t.Parallel()

	var cluster dataproc.Cluster
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withMetadataAndTags(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.basic", &cluster),

					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.metadata.foo", "bar"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.metadata.baz", "qux"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.tags.#", "4"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_singleNodeCluster(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_singleNodeCluster(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.single_node_cluster", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.single_node_cluster", "cluster_config.0.master_config.0.num_instances", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.single_node_cluster", "cluster_config.0.worker_config.0.num_instances", "0"),

					// We set the "dataproc:dataproc.allow.zero.workers" override property.
					// GCP should populate the 'properties' value with this value, as well as many others
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.single_node_cluster", "cluster_config.0.software_config.0.properties.%"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_updatable(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_updatable(rnd, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.updatable", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.master_config.0.num_instances", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.worker_config.0.num_instances", "2"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.preemptible_worker_config.0.num_instances", "1")),
			},
			{
				Config: testAccDataprocCluster_updatable(rnd, 3, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.master_config.0.num_instances", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.worker_config.0.num_instances", "3"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.preemptible_worker_config.0.num_instances", "2")),
			},
		},
	})
}

func TestAccDataprocCluster_withStagingBucket(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	clusterName := fmt.Sprintf("dproc-cluster-test-%s", rnd)
	bucketName := fmt.Sprintf("%s-bucket", clusterName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withStagingBucketAndCluster(clusterName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_bucket", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_bucket", "cluster_config.0.staging_bucket", bucketName),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_bucket", "cluster_config.0.bucket", bucketName)),
			},
			{
				// Simulate destroy of cluster by removing it from definition,
				// but leaving the storage bucket (should not be auto deleted)
				Config: testAccDataprocCluster_withStagingBucketOnly(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocStagingBucketExists(bucketName),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withInitAction(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	bucketName := fmt.Sprintf("dproc-cluster-test-%s-init-bucket", rnd)
	objectName := "msg.txt"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withInitAction(rnd, bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_init_action", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_init_action", "cluster_config.0.initialization_action.#", "2"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_init_action", "cluster_config.0.initialization_action.0.timeout_sec", "500"),
					testAccCheckDataprocClusterInitActionSucceeded(bucketName, objectName),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withConfigOverrides(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withConfigOverrides(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_config_overrides", &cluster),
					validateDataprocCluster_withConfigOverrides("google_dataproc_cluster.with_config_overrides", &cluster),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withServiceAcc(t *testing.T) {
	t.Parallel()

	sa := "a" + acctest.RandString(10)
	saEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", sa, getTestProjectFromEnv())
	rnd := acctest.RandString(10)

	var cluster dataproc.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withServiceAcc(sa, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(
						"google_dataproc_cluster.with_service_account", &cluster),
					testAccCheckDataprocClusterHasServiceScopes(t, &cluster,
						"https://www.googleapis.com/auth/cloud.useraccounts.readonly",
						"https://www.googleapis.com/auth/devstorage.read_write",
						"https://www.googleapis.com/auth/logging.write",
						"https://www.googleapis.com/auth/monitoring",
					),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_service_account", "cluster_config.0.gce_cluster_config.0.service_account", saEmail),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withImageVersion(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withImageVersion(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_image_version", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_image_version", "cluster_config.0.software_config.0.image_version", "1.3.7-deb9"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withOptionalComponents(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withOptionalComponents(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_opt_components", &cluster),
					testAccCheckDataprocClusterHasOptionalComponents(&cluster, "ANACONDA", "ZOOKEEPER"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withLabels(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withLabels(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_labels", &cluster),

					// We only provide one, but GCP adds three, so expect 4. This means unfortunately a
					// diff will exist unless the user adds these in. An alternative approach would
					// be to follow the same approach as properties, i.e. split in into labels
					// and override_labels
					//
					// The config is currently configured with ignore_changes = ["labels"] to handle this
					//
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "4"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key1", "value1"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withNetworkRefs(t *testing.T) {
	t.Parallel()

	var c1, c2 dataproc.Cluster
	rnd := acctest.RandString(10)
	netName := fmt.Sprintf(`dproc-cluster-test-%s-net`, rnd)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withNetworkRefs(rnd, netName),
				Check: resource.ComposeTestCheckFunc(
					// successful creation of the clusters is good enough to assess it worked
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_net_ref_by_url", &c1),
					testAccCheckDataprocClusterExists("google_dataproc_cluster.with_net_ref_by_name", &c2),
				),
			},
		},
	})
}

func TestAccDataprocCluster_KMS(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	kms := BootstrapKMSKey(t)
	pid := getTestProjectFromEnv()

	var cluster dataproc.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_KMS(pid, rnd, kms.CryptoKey.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists("google_dataproc_cluster.kms", &cluster),
				),
			},
		},
	})
}

func testAccCheckDataprocClusterDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataproc_cluster" {
				continue
			}

			if rs.Primary.ID == "" {
				return fmt.Errorf("Unable to verify delete of dataproc cluster, ID is empty")
			}

			attributes := rs.Primary.Attributes
			project, err := getTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			parts := strings.Split(rs.Primary.ID, "/")
			clusterId := parts[len(parts)-1]
			_, err = config.clientDataprocBeta.Projects.Regions.Clusters.Get(
				project, attributes["region"], clusterId).Do()

			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
					return nil
				} else if ok {
					return fmt.Errorf("Error validating cluster deleted. Code: %d. Message: %s", gerr.Code, gerr.Message)
				}
				return fmt.Errorf("Error validating cluster deleted. %s", err.Error())
			}
			return fmt.Errorf("Dataproc cluster still exists")
		}

		return nil
	}
}

func testAccCheckDataprocClusterHasServiceScopes(t *testing.T, cluster *dataproc.Cluster, scopes ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {

		if !reflect.DeepEqual(scopes, cluster.Config.GceClusterConfig.ServiceAccountScopes) {
			return fmt.Errorf("Cluster does not contain expected set of service account scopes : %v : instead %v",
				scopes, cluster.Config.GceClusterConfig.ServiceAccountScopes)
		}
		return nil
	}
}

func validateBucketExists(bucket string, config *Config) (bool, error) {
	_, err := config.clientStorage.Buckets.Get(bucket).Do()

	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			return false, nil
		} else if ok {
			return false, fmt.Errorf("Error validating bucket exists: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
		}
		return false, fmt.Errorf("Error validating bucket exists: %s", err.Error())
	}
	return true, nil
}

func testAccCheckDataprocStagingBucketExists(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config := testAccProvider.Meta().(*Config)

		exists, err := validateBucketExists(bucketName, config)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Staging Bucket %s does not exist", bucketName)
		}
		return nil
	}
}

func testAccCheckDataprocClusterHasOptionalComponents(cluster *dataproc.Cluster, components ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {

		if !reflect.DeepEqual(components, cluster.Config.SoftwareConfig.OptionalComponents) {
			return fmt.Errorf("Cluster does not contain expected optional components : %v : instead %v",
				components, cluster.Config.SoftwareConfig.OptionalComponents)
		}
		return nil
	}
}

func testAccCheckDataprocClusterInitActionSucceeded(bucket, object string) resource.TestCheckFunc {

	// The init script will have created an object in the specified bucket.
	// Ensure it exists
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientStorage.Objects.Get(bucket, object).Do()
		if err != nil {
			return fmt.Errorf("Unable to verify init action success: Error reading object %s in bucket %s: %v", object, bucket, err)
		}

		return nil
	}
}

func validateDataprocCluster_withConfigOverrides(n string, cluster *dataproc.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		type tfAndGCPTestField struct {
			tfAttr       string
			expectedVal  string
			actualGCPVal string
		}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Terraform resource Not found: %s", n)
		}

		if cluster.Config.MasterConfig == nil || cluster.Config.WorkerConfig == nil || cluster.Config.SecondaryWorkerConfig == nil {
			return fmt.Errorf("Master/Worker/SecondaryConfig values not set in GCP, expecting values")
		}

		clusterTests := []tfAndGCPTestField{
			{"cluster_config.0.master_config.0.num_instances", "3", strconv.Itoa(int(cluster.Config.MasterConfig.NumInstances))},
			{"cluster_config.0.master_config.0.disk_config.0.boot_disk_size_gb", "15", strconv.Itoa(int(cluster.Config.MasterConfig.DiskConfig.BootDiskSizeGb))},
			{"cluster_config.0.master_config.0.disk_config.0.num_local_ssds", "0", strconv.Itoa(int(cluster.Config.MasterConfig.DiskConfig.NumLocalSsds))},
			{"cluster_config.0.master_config.0.disk_config.0.boot_disk_type", "pd-ssd", cluster.Config.MasterConfig.DiskConfig.BootDiskType},
			{"cluster_config.0.master_config.0.machine_type", "n1-standard-1", GetResourceNameFromSelfLink(cluster.Config.MasterConfig.MachineTypeUri)},
			{"cluster_config.0.master_config.0.instance_names.#", "3", strconv.Itoa(len(cluster.Config.MasterConfig.InstanceNames))},
			{"cluster_config.0.master_config.0.min_cpu_platform", "Intel Skylake", cluster.Config.MasterConfig.MinCpuPlatform},

			{"cluster_config.0.worker_config.0.num_instances", "3", strconv.Itoa(int(cluster.Config.WorkerConfig.NumInstances))},
			{"cluster_config.0.worker_config.0.disk_config.0.boot_disk_size_gb", "16", strconv.Itoa(int(cluster.Config.WorkerConfig.DiskConfig.BootDiskSizeGb))},
			{"cluster_config.0.worker_config.0.disk_config.0.num_local_ssds", "1", strconv.Itoa(int(cluster.Config.WorkerConfig.DiskConfig.NumLocalSsds))},
			{"cluster_config.0.worker_config.0.disk_config.0.boot_disk_type", "pd-standard", cluster.Config.WorkerConfig.DiskConfig.BootDiskType},
			{"cluster_config.0.worker_config.0.machine_type", "n1-standard-1", GetResourceNameFromSelfLink(cluster.Config.WorkerConfig.MachineTypeUri)},
			{"cluster_config.0.worker_config.0.instance_names.#", "3", strconv.Itoa(len(cluster.Config.WorkerConfig.InstanceNames))},
			{"cluster_config.0.worker_config.0.min_cpu_platform", "Intel Broadwell", cluster.Config.WorkerConfig.MinCpuPlatform},

			{"cluster_config.0.preemptible_worker_config.0.num_instances", "1", strconv.Itoa(int(cluster.Config.SecondaryWorkerConfig.NumInstances))},
			{"cluster_config.0.preemptible_worker_config.0.disk_config.0.boot_disk_size_gb", "17", strconv.Itoa(int(cluster.Config.SecondaryWorkerConfig.DiskConfig.BootDiskSizeGb))},
			{"cluster_config.0.preemptible_worker_config.0.disk_config.0.num_local_ssds", "1", strconv.Itoa(int(cluster.Config.SecondaryWorkerConfig.DiskConfig.NumLocalSsds))},
			{"cluster_config.0.preemptible_worker_config.0.disk_config.0.boot_disk_type", "pd-ssd", cluster.Config.SecondaryWorkerConfig.DiskConfig.BootDiskType},
			{"cluster_config.0.preemptible_worker_config.0.instance_names.#", "1", strconv.Itoa(len(cluster.Config.SecondaryWorkerConfig.InstanceNames))},
		}

		for _, attrs := range clusterTests {
			tfVal := rs.Primary.Attributes[attrs.tfAttr]
			if tfVal != attrs.expectedVal {
				return fmt.Errorf("%s: Terraform Attribute value '%s' is not as expected '%s' ", attrs.tfAttr, tfVal, attrs.expectedVal)
			}
			if attrs.actualGCPVal != tfVal {
				return fmt.Errorf("%s: Terraform Attribute value '%s' is not aligned with that in GCP '%s' ", attrs.tfAttr, tfVal, attrs.actualGCPVal)
			}
		}

		return nil
	}
}

func testAccCheckDataprocClusterExists(n string, cluster *dataproc.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Terraform resource Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for Dataproc cluster")
		}

		config := testAccProvider.Meta().(*Config)
		project, err := getTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		parts := strings.Split(rs.Primary.ID, "/")
		clusterId := parts[len(parts)-1]
		found, err := config.clientDataprocBeta.Projects.Regions.Clusters.Get(
			project, rs.Primary.Attributes["region"], clusterId).Do()
		if err != nil {
			return err
		}

		if found.ClusterName != clusterId {
			return fmt.Errorf("Dataproc cluster %s not found, found %s instead", clusterId, cluster.ClusterName)
		}

		*cluster = *found

		return nil
	}
}

func testAccCheckDataproc_missingZoneGlobalRegion1(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name   = "dproc-cluster-test-%s"
  region = "global"
}
`, rnd)
}

func testAccCheckDataproc_missingZoneGlobalRegion2(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name   = "dproc-cluster-test-%s"
  region = "global"

  cluster_config {
    gce_cluster_config {
      network = "default"
    }
  }
}
`, rnd)
}

func testAccDataprocCluster_basic(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"
}
`, rnd)
}

func testAccDataprocCluster_withAccelerators(rnd, zone, acceleratorType string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "accelerated_cluster" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    gce_cluster_config {
      zone = "%s"
    }

    master_config {
      accelerators {
        accelerator_type  = "%s"
        accelerator_count = "1"
      }
    }

    worker_config {
      accelerators {
        accelerator_type  = "%s"
        accelerator_count = "1"
      }
    }
  }
}
`, rnd, zone, acceleratorType, acceleratorType)
}

func testAccDataprocCluster_withInternalIpOnlyTrue(rnd string) string {
	return fmt.Sprintf(`
variable "subnetwork_cidr" {
  default = "10.0.0.0/16"
}

resource "google_compute_network" "dataproc_network" {
  name                    = "dataproc-internalip-network-%s"
  auto_create_subnetworks = false
}

#
# Create a subnet with Private IP Access enabled to test
# deploying a Dataproc cluster with Internal IP Only enabled.
#
resource "google_compute_subnetwork" "dataproc_subnetwork" {
  name                     = "dataproc-internalip-subnetwork-%s"
  ip_cidr_range            = var.subnetwork_cidr
  network                  = google_compute_network.dataproc_network.self_link
  region                   = "us-central1"
  private_ip_google_access = true
}

#
# The default network within GCP already comes pre configured with
# certain firewall rules open to allow internal communication. As we
# are creating a new one here for this test, we need to additionally
# open up similar rules to allow the nodes to talk to each other
# internally as part of their configuration or this will just hang.
#
resource "google_compute_firewall" "dataproc_network_firewall" {
  name        = "dproc-cluster-test-allow-internal"
  description = "Firewall rules for dataproc Terraform acceptance testing"
  network     = google_compute_network.dataproc_network.name

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

  source_ranges = [var.subnetwork_cidr]
}

resource "google_dataproc_cluster" "basic" {
  name       = "dproc-cluster-test-%s"
  region     = "us-central1"
  depends_on = [google_compute_firewall.dataproc_network_firewall]

  cluster_config {
    gce_cluster_config {
      subnetwork       = google_compute_subnetwork.dataproc_subnetwork.name
      internal_ip_only = true
    }
  }
}
`, rnd, rnd, rnd)
}

func testAccDataprocCluster_withMetadataAndTags(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    gce_cluster_config {
      metadata = {
        foo = "bar"
        baz = "qux"
      }
      tags = ["my-tag", "your-tag", "our-tag", "their-tag"]
    }
  }
}
`, rnd)
}

func testAccDataprocCluster_singleNodeCluster(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "single_node_cluster" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }
  }
}
`, rnd)
}

func testAccDataprocCluster_withConfigOverrides(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_config_overrides" {
  name     = "dproc-cluster-test-%s"
  region   = "us-central1"

  cluster_config {
    master_config {
      num_instances = 3
      machine_type  = "n1-standard-1"
      disk_config {
        boot_disk_type    = "pd-ssd"
        boot_disk_size_gb = 15
      }
      min_cpu_platform = "Intel Skylake"
    }

    worker_config {
      num_instances = 3
      machine_type  = "n1-standard-1"
      disk_config {
        boot_disk_type    = "pd-standard"
        boot_disk_size_gb = 16
        num_local_ssds    = 1
      }

      min_cpu_platform = "Intel Broadwell"
    }

    preemptible_worker_config {
      num_instances = 1
      disk_config {
        boot_disk_type    = "pd-ssd"
        boot_disk_size_gb = 17
        num_local_ssds    = 1
      }
    }
  }
}
`, rnd)
}

func testAccDataprocCluster_withInitAction(rnd, bucket, objName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "init_bucket" {
  name          = "%s"
  force_destroy = "true"
}

resource "google_storage_bucket_object" "init_script" {
  name    = "dproc-cluster-test-%s-init-script.sh"
  bucket  = google_storage_bucket.init_bucket.name
  content = <<EOL
#!/bin/bash
echo "init action success" >> /tmp/%s
gsutil cp /tmp/%s ${google_storage_bucket.init_bucket.url}
EOL

}

resource "google_dataproc_cluster" "with_init_action" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    master_config {
      machine_type = "n1-standard-1"
      disk_config {
        boot_disk_size_gb = 15
      }
    }

    initialization_action {
      script      = "${google_storage_bucket.init_bucket.url}/${google_storage_bucket_object.init_script.name}"
      timeout_sec = 500
    }
    initialization_action {
      script = "${google_storage_bucket.init_bucket.url}/${google_storage_bucket_object.init_script.name}"
    }
  }
}
`, bucket, rnd, objName, objName, rnd)
}

func testAccDataprocCluster_updatable(rnd string, w, p int) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "updatable" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    master_config {
      num_instances = "1"
      machine_type  = "n1-standard-1"
      disk_config {
        boot_disk_size_gb = 15
      }
    }

    worker_config {
      num_instances = "%d"
      machine_type  = "n1-standard-1"
      disk_config {
        boot_disk_size_gb = 15
      }
    }

    preemptible_worker_config {
      num_instances = "%d"
      disk_config {
        boot_disk_size_gb = 15
      }
    }
  }
}
`, rnd, w, p)
}

func testAccDataprocCluster_withStagingBucketOnly(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  force_destroy = "true"
}
`, bucketName)
}

func testAccDataprocCluster_withStagingBucketAndCluster(clusterName, bucketName string) string {
	return fmt.Sprintf(`
%s

resource "google_dataproc_cluster" "with_bucket" {
  name   = "%s"
  region = "us-central1"

  cluster_config {
    staging_bucket = google_storage_bucket.bucket.name

    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    master_config {
      machine_type = "n1-standard-1"
      disk_config {
        boot_disk_size_gb = 15
      }
    }
  }
}
`, testAccDataprocCluster_withStagingBucketOnly(bucketName), clusterName)
}

func testAccDataprocCluster_withLabels(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_labels" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  labels = {
    key1 = "value1"
  }

  # This is because GCP automatically adds its own labels as well.
  # In this case we just want to test our newly added label is there
  lifecycle {
    ignore_changes = [labels]
  }
}
`, rnd)
}

func testAccDataprocCluster_withImageVersion(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_image_version" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    software_config {
      image_version = "1.3.7-deb9"
    }
  }
}
`, rnd)
}

func testAccDataprocCluster_withOptionalComponents(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_opt_components" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    software_config {
      optional_components = ["ANACONDA", "ZOOKEEPER"]
    }
  }
}
`, rnd)
}

func testAccDataprocCluster_withServiceAcc(sa string, rnd string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
  account_id = "%s"
}

resource "google_project_iam_member" "service_account" {
  role   = "roles/dataproc.worker"
  member = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_dataproc_cluster" "with_service_account" {
  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    master_config {
      machine_type = "n1-standard-1"
      disk_config {
        boot_disk_size_gb = 15
      }
    }

    gce_cluster_config {
      service_account = google_service_account.service_account.email
      service_account_scopes = [
		#	User supplied scopes
        "https://www.googleapis.com/auth/monitoring",
		#	The following scopes necessary for the cluster to function properly are
		#	always added, even if not explicitly specified:
		#		useraccounts-ro: https://www.googleapis.com/auth/cloud.useraccounts.readonly
		#		storage-rw:      https://www.googleapis.com/auth/devstorage.read_write
		#		logging-write:   https://www.googleapis.com/auth/logging.write
        "useraccounts-ro",
        "storage-rw",
        "logging-write",
      ]
    }
  }

  depends_on = [google_project_iam_member.service_account]
}
`, sa, rnd)
}

func testAccDataprocCluster_withNetworkRefs(rnd, netName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "dataproc_network" {
  name                    = "%s"
  auto_create_subnetworks = true
}

#
# The default network within GCP already comes pre configured with
# certain firewall rules open to allow internal communication. As we
# are creating a new one here for this test, we need to additionally
# open up similar rules to allow the nodes to talk to each other
# internally as part of their configuration or this will just hang.
#
resource "google_compute_firewall" "dataproc_network_firewall" {
  name          = "dproc-cluster-test-%s-allow-internal"
  description   = "Firewall rules for dataproc Terraform acceptance testing"
  network       = google_compute_network.dataproc_network.name
  source_ranges = ["192.168.0.0/16"]

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
  name       = "dproc-cluster-test-%s-name"
  region     = "us-central1"
  depends_on = [google_compute_firewall.dataproc_network_firewall]

  cluster_config {
    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    master_config {
      machine_type = "n1-standard-1"
      disk_config {
        boot_disk_size_gb = 15
      }
    }

    gce_cluster_config {
      network = google_compute_network.dataproc_network.name
    }
  }
}

resource "google_dataproc_cluster" "with_net_ref_by_url" {
  name       = "dproc-cluster-test-%s-url"
  region     = "us-central1"
  depends_on = [google_compute_firewall.dataproc_network_firewall]

  cluster_config {
    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    master_config {
      machine_type = "n1-standard-1"
      disk_config {
        boot_disk_size_gb = 15
      }
    }

    gce_cluster_config {
      network = google_compute_network.dataproc_network.self_link
    }
  }
}
`, netName, rnd, rnd, rnd)
}

func testAccDataprocCluster_KMS(pid, rnd, kmsKey string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_project_iam_member" "kms-project-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com"
}

resource "google_dataproc_cluster" "kms" {
  depends_on = [google_project_iam_member.kms-project-binding]

  name   = "dproc-cluster-test-%s"
  region = "us-central1"

  cluster_config {
    encryption_config {
      kms_key_name = "%s"
    }
  }
}
`, pid, rnd, kmsKey)
}
