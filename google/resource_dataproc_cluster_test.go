package google

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/api/googleapi"

	"google.golang.org/api/dataproc/v1"
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

	rnd := randString(t, 10)
	vcrTest(t, resource.TestCase{
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

	rnd := randString(t, 10)
	vcrTest(t, resource.TestCase{
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
	rnd := randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_basic(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.basic", &cluster),

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

func TestAccDataprocVirtualCluster_basic(t *testing.T) {
	t.Parallel()

	var cluster dataproc.Cluster
	rnd := randString(t, 10)
	pid := getTestProjectFromEnv()
	version := "3.1-dataproc-7"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocVirtualCluster_basic(pid, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.virtual_cluster", &cluster),

					// Expect 1 dataproc on gke instances with computed values
					resource.TestCheckResourceAttr("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.#", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.#", "1"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.kubernetes_namespace"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.kubernetes_software_config.#", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.kubernetes_software_config.0.component_version.SPARK", version),

					resource.TestCheckResourceAttr("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.#", "1"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.gke_cluster_target"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.node_pool_target.#", "1"),
					resource.TestCheckResourceAttrSet("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.node_pool_target.0.node_pool"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.virtual_cluster", "virtual_cluster_config.0.kubernetes_cluster_config.0.gke_cluster_config.0.node_pool_target.0.roles.#", "1"),
					testAccCheckDataprocGkeClusterNodePoolsHaveRoles(&cluster, "DEFAULT"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withAccelerators(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster

	project := getTestProjectFromEnv()
	acceleratorType := "nvidia-tesla-k80"
	zone := "us-central1-c"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withAccelerators(rnd, acceleratorType, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.accelerated_cluster", &cluster),
					testAccCheckDataprocClusterAccelerator(&cluster, project, 1, 1),
				),
			},
		},
	})
}

func testAccCheckDataprocClusterAccelerator(cluster *dataproc.Cluster, project string, masterCount int, workerCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expectedUri := fmt.Sprintf("projects/%s/zones/.*/acceleratorTypes/nvidia-tesla-k80", project)
		r := regexp.MustCompile(expectedUri)

		master := cluster.Config.MasterConfig.Accelerators
		if len(master) != 1 {
			return fmt.Errorf("Saw %d master accelerator types instead of 1", len(master))
		}

		if int(master[0].AcceleratorCount) != masterCount {
			return fmt.Errorf("Saw %d master accelerators instead of %d", int(master[0].AcceleratorCount), masterCount)
		}

		matches := r.FindStringSubmatch(master[0].AcceleratorTypeUri)
		if len(matches) != 1 {
			return fmt.Errorf("Saw %s master accelerator type instead of %s", master[0].AcceleratorTypeUri, expectedUri)
		}

		worker := cluster.Config.WorkerConfig.Accelerators
		if len(worker) != 1 {
			return fmt.Errorf("Saw %d worker accelerator types instead of 1", len(worker))
		}

		if int(worker[0].AcceleratorCount) != workerCount {
			return fmt.Errorf("Saw %d worker accelerators instead of %d", int(worker[0].AcceleratorCount), workerCount)
		}

		matches = r.FindStringSubmatch(worker[0].AcceleratorTypeUri)
		if len(matches) != 1 {
			return fmt.Errorf("Saw %s worker accelerator type instead of %s", worker[0].AcceleratorTypeUri, expectedUri)
		}

		return nil
	}
}

func TestAccDataprocCluster_withInternalIpOnlyTrueAndShieldedConfig(t *testing.T) {
	t.Parallel()

	var cluster dataproc.Cluster
	rnd := randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withInternalIpOnlyTrueAndShieldedConfig(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.basic", &cluster),

					// Testing behavior for Dataproc to use only internal IP addresses
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.internal_ip_only", "true"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_integrity_monitoring", "true"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_secure_boot", "true"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.basic", "cluster_config.0.gce_cluster_config.0.shielded_instance_config.0.enable_vtpm", "true"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withMetadataAndTags(t *testing.T) {
	t.Parallel()

	var cluster dataproc.Cluster
	rnd := randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withMetadataAndTags(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.basic", &cluster),

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

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_singleNodeCluster(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.single_node_cluster", &cluster),
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

	rnd := randString(t, 10)
	var cluster dataproc.Cluster

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_updatable(rnd, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.updatable", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.master_config.0.num_instances", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.worker_config.0.num_instances", "2"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.preemptible_worker_config.0.num_instances", "1")),
			},
			{
				Config: testAccDataprocCluster_updatable(rnd, 2, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.updatable", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.master_config.0.num_instances", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.worker_config.0.num_instances", "2"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.updatable", "cluster_config.0.preemptible_worker_config.0.num_instances", "0")),
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

func TestAccDataprocCluster_nonPreemptibleSecondary(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_nonPreemptibleSecondary(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.non_preemptible_secondary", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.non_preemptible_secondary", "cluster_config.0.preemptible_worker_config.0.preemptibility", "NON_PREEMPTIBLE"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_spotSecondary(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_spotSecondary(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.spot_secondary", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.spot_secondary", "cluster_config.0.preemptible_worker_config.0.preemptibility", "SPOT"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withStagingBucket(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	clusterName := fmt.Sprintf("tf-test-dproc-%s", rnd)
	bucketName := fmt.Sprintf("%s-bucket", clusterName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withStagingBucketAndCluster(clusterName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_bucket", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_bucket", "cluster_config.0.staging_bucket", bucketName),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_bucket", "cluster_config.0.bucket", bucketName)),
			},
			{
				// Simulate destroy of cluster by removing it from definition,
				// but leaving the storage bucket (should not be auto deleted)
				Config: testAccDataprocCluster_withStagingBucketOnly(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocStagingBucketExists(t, bucketName),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withTempBucket(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	clusterName := fmt.Sprintf("tf-test-dproc-%s", rnd)
	bucketName := fmt.Sprintf("%s-temp-bucket", clusterName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withTempBucketAndCluster(clusterName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_bucket", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_bucket", "cluster_config.0.temp_bucket", bucketName)),
			},
			{
				// Simulate destroy of cluster by removing it from definition,
				// but leaving the temp bucket (should not be auto deleted)
				Config: testAccDataprocCluster_withTempBucketOnly(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocTempBucketExists(t, bucketName),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withInitAction(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	bucketName := fmt.Sprintf("tf-test-dproc-%s-init-bucket", rnd)
	objectName := "msg.txt"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withInitAction(rnd, bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_init_action", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_init_action", "cluster_config.0.initialization_action.#", "2"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_init_action", "cluster_config.0.initialization_action.0.timeout_sec", "500"),
					testAccCheckDataprocClusterInitActionSucceeded(t, bucketName, objectName),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withConfigOverrides(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withConfigOverrides(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_config_overrides", &cluster),
					validateDataprocCluster_withConfigOverrides("google_dataproc_cluster.with_config_overrides", &cluster),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withServiceAcc(t *testing.T) {
	t.Parallel()

	sa := "a" + randString(t, 10)
	saEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", sa, getTestProjectFromEnv())
	rnd := randString(t, 10)

	var cluster dataproc.Cluster

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withServiceAcc(sa, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(
						t, "google_dataproc_cluster.with_service_account", &cluster),
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

	rnd := randString(t, 10)
	version := "2.0.35-debian10"

	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withImageVersion(rnd, version),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_image_version", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_image_version", "cluster_config.0.software_config.0.image_version", version),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withOptionalComponents(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withOptionalComponents(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_opt_components", &cluster),
					testAccCheckDataprocClusterHasOptionalComponents(&cluster, "ZOOKEEPER", "DOCKER"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withLifecycleConfigIdleDeleteTtl(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withLifecycleConfigIdleDeleteTtl(rnd, "600s"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_lifecycle_config", &cluster),
				),
			},
			{
				Config: testAccDataprocCluster_withLifecycleConfigIdleDeleteTtl(rnd, "610s"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_lifecycle_config", &cluster),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withLifecycleConfigAutoDeletion(t *testing.T) {
	// Uses time.Now
	skipIfVcr(t)
	t.Parallel()

	rnd := randString(t, 10)
	now := time.Now()
	fmtString := "2006-01-02T15:04:05.072Z"

	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withLifecycleConfigAutoDeletionTime(rnd, now.Add(time.Hour*10).Format(fmtString)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_lifecycle_config", &cluster),
				),
			},
			{
				Config: testAccDataprocCluster_withLifecycleConfigAutoDeletionTime(rnd, now.Add(time.Hour*20).Format(fmtString)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_lifecycle_config", &cluster),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withLabels(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withLabels(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					// We only provide one, but GCP adds three and we added goog-dataproc-autozone internally, so expect 5.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key1", "value1"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withNetworkRefs(t *testing.T) {
	// Multiple fine-grained resources
	skipIfVcr(t)
	t.Parallel()

	var c1, c2 dataproc.Cluster
	rnd := randString(t, 10)
	netName := fmt.Sprintf(`dproc-cluster-test-%s-net`, rnd)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withNetworkRefs(rnd, netName),
				Check: resource.ComposeTestCheckFunc(
					// successful creation of the clusters is good enough to assess it worked
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_net_ref_by_url", &c1),
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_net_ref_by_name", &c2),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withEndpointConfig(t *testing.T) {
	t.Parallel()

	var cluster dataproc.Cluster
	rnd := randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withEndpointConfig(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_endpoint_config", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_endpoint_config", "cluster_config.0.endpoint_config.0.enable_http_port_access", "true"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_KMS(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	kms := BootstrapKMSKey(t)
	pid := getTestProjectFromEnv()

	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_KMS(pid, rnd, kms.CryptoKey.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.kms", &cluster),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withKerberos(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)
	kms := BootstrapKMSKey(t)

	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withKerberos(rnd, kms.CryptoKey.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.kerb", &cluster),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withAutoscalingPolicy(t *testing.T) {
	t.Parallel()

	rnd := randString(t, 10)

	var cluster dataproc.Cluster
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withAutoscalingPolicy(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.basic", &cluster),
					testAccCheckDataprocClusterAutoscaling(t, &cluster, true),
				),
			},
			{
				Config: testAccDataprocCluster_removeAutoscalingPolicy(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.basic", &cluster),
					testAccCheckDataprocClusterAutoscaling(t, &cluster, false),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withMetastoreConfig(t *testing.T) {
	t.Parallel()

	pid := getTestProjectFromEnv()
	msName_basic := fmt.Sprintf("projects/%s/locations/us-central1/services/metastore-srv", pid)
	msName_update := fmt.Sprintf("projects/%s/locations/us-central1/services/metastore-srv-update", pid)

	var cluster dataproc.Cluster
	rnd := randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withMetastoreConfig(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_metastore_config", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_metastore_config", "cluster_config.0.metastore_config.0.dataproc_metastore_service", msName_basic),
				),
			},
			{
				Config: testAccDataprocCluster_withMetastoreConfig_update(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_metastore_config", &cluster),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_metastore_config", "cluster_config.0.metastore_config.0.dataproc_metastore_service", msName_update),
				),
			},
		},
	})
}

func testAccCheckDataprocClusterDestroy(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

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
			_, err = config.NewDataprocClient(config.userAgent).Projects.Regions.Clusters.Get(
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

func testAccCheckDataprocClusterAutoscaling(t *testing.T, cluster *dataproc.Cluster, expectAutoscaling bool) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		if cluster.Config.AutoscalingConfig == nil && expectAutoscaling {
			return fmt.Errorf("Cluster does not contain AutoscalingConfig, expected it would")
		} else if cluster.Config.AutoscalingConfig != nil && !expectAutoscaling {
			return fmt.Errorf("Cluster contains AutoscalingConfig, expected it not to")
		}

		return nil
	}
}

func validateBucketExists(bucket string, config *Config) (bool, error) {
	_, err := config.NewStorageClient(config.userAgent).Buckets.Get(bucket).Do()

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

func testAccCheckDataprocStagingBucketExists(t *testing.T, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config := googleProviderConfig(t)

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

func testAccCheckDataprocTempBucketExists(t *testing.T, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		config := googleProviderConfig(t)

		exists, err := validateBucketExists(bucketName, config)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Temp Bucket %s does not exist", bucketName)
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

func testAccCheckDataprocClusterInitActionSucceeded(t *testing.T, bucket, object string) resource.TestCheckFunc {

	// The init script will have created an object in the specified bucket.
	// Ensure it exists
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		_, err := config.NewStorageClient(config.userAgent).Objects.Get(bucket, object).Do()
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
			{"cluster_config.0.master_config.0.disk_config.0.boot_disk_size_gb", "35", strconv.Itoa(int(cluster.Config.MasterConfig.DiskConfig.BootDiskSizeGb))},
			{"cluster_config.0.master_config.0.disk_config.0.num_local_ssds", "0", strconv.Itoa(int(cluster.Config.MasterConfig.DiskConfig.NumLocalSsds))},
			{"cluster_config.0.master_config.0.disk_config.0.boot_disk_type", "pd-ssd", cluster.Config.MasterConfig.DiskConfig.BootDiskType},
			{"cluster_config.0.master_config.0.machine_type", "n1-standard-2", GetResourceNameFromSelfLink(cluster.Config.MasterConfig.MachineTypeUri)},
			{"cluster_config.0.master_config.0.instance_names.#", "3", strconv.Itoa(len(cluster.Config.MasterConfig.InstanceNames))},
			{"cluster_config.0.master_config.0.min_cpu_platform", "Intel Skylake", cluster.Config.MasterConfig.MinCpuPlatform},

			{"cluster_config.0.worker_config.0.num_instances", "3", strconv.Itoa(int(cluster.Config.WorkerConfig.NumInstances))},
			{"cluster_config.0.worker_config.0.disk_config.0.boot_disk_size_gb", "35", strconv.Itoa(int(cluster.Config.WorkerConfig.DiskConfig.BootDiskSizeGb))},
			{"cluster_config.0.worker_config.0.disk_config.0.num_local_ssds", "1", strconv.Itoa(int(cluster.Config.WorkerConfig.DiskConfig.NumLocalSsds))},
			{"cluster_config.0.worker_config.0.disk_config.0.boot_disk_type", "pd-standard", cluster.Config.WorkerConfig.DiskConfig.BootDiskType},
			{"cluster_config.0.worker_config.0.machine_type", "n1-standard-2", GetResourceNameFromSelfLink(cluster.Config.WorkerConfig.MachineTypeUri)},
			{"cluster_config.0.worker_config.0.instance_names.#", "3", strconv.Itoa(len(cluster.Config.WorkerConfig.InstanceNames))},
			{"cluster_config.0.worker_config.0.min_cpu_platform", "Intel Broadwell", cluster.Config.WorkerConfig.MinCpuPlatform},

			{"cluster_config.0.preemptible_worker_config.0.num_instances", "1", strconv.Itoa(int(cluster.Config.SecondaryWorkerConfig.NumInstances))},
			{"cluster_config.0.preemptible_worker_config.0.disk_config.0.boot_disk_size_gb", "35", strconv.Itoa(int(cluster.Config.SecondaryWorkerConfig.DiskConfig.BootDiskSizeGb))},
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

func testAccCheckDataprocClusterExists(t *testing.T, n string, cluster *dataproc.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Terraform resource Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for Dataproc cluster")
		}

		config := googleProviderConfig(t)
		project, err := getTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		parts := strings.Split(rs.Primary.ID, "/")
		clusterId := parts[len(parts)-1]
		found, err := config.NewDataprocClient(config.userAgent).Projects.Regions.Clusters.Get(
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
  name   = "tf-test-dproc-%s"
  region = "global"
}
`, rnd)
}

func testAccCheckDataproc_missingZoneGlobalRegion2(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name   = "tf-test-dproc-%s"
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
  name   = "tf-test-dproc-%s"
  region = "us-central1"
}
`, rnd)
}

func testAccDataprocVirtualCluster_basic(projectID string, rnd string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_container_cluster" "primary" {
  name     = "tf-test-gke-%s"
  location = "us-central1-a"

  initial_node_count = 1

  workload_identity_config {
    workload_pool = "${data.google_project.project.project_id}.svc.id.goog"
  }
}

resource "google_project_iam_binding" "workloadidentity" {
  project = "%s"
  role    = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${data.google_project.project.project_id}.svc.id.goog[tf-test-dproc-%s/agent]",
    "serviceAccount:${data.google_project.project.project_id}.svc.id.goog[tf-test-dproc-%s/spark-driver]",
    "serviceAccount:${data.google_project.project.project_id}.svc.id.goog[tf-test-dproc-%s/spark-executor]",
  ]
}

resource "google_dataproc_cluster" "virtual_cluster" {
	depends_on = [
	  google_project_iam_binding.workloadidentity
	]
  
	name   	= "tf-test-dproc-%s"
	region  = "us-central1"
  
	virtual_cluster_config {
	  kubernetes_cluster_config {
		kubernetes_namespace = "tf-test-dproc-%s"
		kubernetes_software_config {
		  component_version = {
			"SPARK": "3.1-dataproc-7",
		  }
		}
		gke_cluster_config {
		  gke_cluster_target = google_container_cluster.primary.id
		  node_pool_target {
			node_pool = "tf-test-gke-np-%s"
			roles = [
			  "DEFAULT"
			]
		  }
		} 
	  }
	}
  }
`, projectID, rnd, projectID, rnd, rnd, rnd, rnd, rnd, rnd)
}

func testAccCheckDataprocGkeClusterNodePoolsHaveRoles(cluster *dataproc.Cluster, roles ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {

		for _, nodePool := range cluster.VirtualClusterConfig.KubernetesClusterConfig.GkeClusterConfig.NodePoolTarget {
			if reflect.DeepEqual(roles, nodePool.Roles) {
				return nil
			}
		}

		return fmt.Errorf("Cluster NodePools does not contain expected roles : %v", roles)
	}
}

func testAccDataprocCluster_withAccelerators(rnd, acceleratorType, zone string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "accelerated_cluster" {
  name   = "tf-test-dproc-%s"
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

func testAccDataprocCluster_withInternalIpOnlyTrueAndShieldedConfig(rnd string) string {
	return fmt.Sprintf(`
variable "subnetwork_cidr" {
  default = "10.0.0.0/16"
}

resource "google_compute_network" "dataproc_network" {
  name                    = "tf-test-dproc-net-%s"
  auto_create_subnetworks = false
}

#
# Create a subnet with Private IP Access enabled to test
# deploying a Dataproc cluster with Internal IP Only enabled.
#
resource "google_compute_subnetwork" "dataproc_subnetwork" {
  name                     = "tf-test-dproc-subnet-%s"
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
  name        = "tf-test-dproc-firewall-%s"
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
  name       = "tf-test-dproc-%s"
  region     = "us-central1"
  depends_on = [google_compute_firewall.dataproc_network_firewall]

  cluster_config {
    gce_cluster_config {
      subnetwork       = google_compute_subnetwork.dataproc_subnetwork.name
      internal_ip_only = true
      shielded_instance_config{
        enable_integrity_monitoring = true
        enable_secure_boot          = true
        enable_vtpm                 = true
      }
    }
  }
}
`, rnd, rnd, rnd, rnd)
}

func testAccDataprocCluster_withMetadataAndTags(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name   = "tf-test-dproc-%s"
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
  name   = "tf-test-dproc-%s"
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
  name     = "tf-test-dproc-%s"
  region   = "us-central1"

  cluster_config {
    master_config {
      num_instances = 3
      machine_type  = "n1-standard-2"  // can't be e2 because of min_cpu_platform
      disk_config {
        boot_disk_type    = "pd-ssd"
        boot_disk_size_gb = 35
      }
      min_cpu_platform = "Intel Skylake"
    }

    worker_config {
      num_instances = 3
      machine_type  = "n1-standard-2"  // can't be e2 because of min_cpu_platform
      disk_config {
        boot_disk_type    = "pd-standard"
        boot_disk_size_gb = 35
        num_local_ssds    = 1
      }

      min_cpu_platform = "Intel Broadwell"
    }

    preemptible_worker_config {
      num_instances = 1
      disk_config {
        boot_disk_type    = "pd-ssd"
        boot_disk_size_gb = 35
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
  location      = "US"
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
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    master_config {
      machine_type = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
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
  name   = "tf-test-dproc-%s"
  region = "us-central1"
  graceful_decommission_timeout = "0.2s"

  cluster_config {
    master_config {
      num_instances = "1"
      machine_type  = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
      }
    }

    worker_config {
      num_instances = "%d"
      machine_type  = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
      }
    }

    preemptible_worker_config {
      num_instances = "%d"
      disk_config {
        boot_disk_size_gb = 35
      }
    }
  }
}
`, rnd, w, p)
}

func testAccDataprocCluster_nonPreemptibleSecondary(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "non_preemptible_secondary" {
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    master_config {
	  num_instances = "1"
	  machine_type  = "e2-medium"
	  disk_config {
		boot_disk_size_gb = 35
	  }
	}
  
	worker_config {
	  num_instances = "2"
	  machine_type  = "e2-medium"
	  disk_config {
		boot_disk_size_gb = 35
	  }
	}
  
	preemptible_worker_config {
	  num_instances = "1"
	  preemptibility = "NON_PREEMPTIBLE"
	  disk_config {
		boot_disk_size_gb = 35
	  }
	}
  }
}
	`, rnd)
}

func testAccDataprocCluster_spotSecondary(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "spot_secondary" {
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    master_config {
      num_instances = "1"
      machine_type  = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
      }
    }

    worker_config {
      num_instances = "2"
      machine_type  = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
      }
    }

    preemptible_worker_config {
      num_instances = "1"
      preemptibility = "SPOT"
      disk_config {
        boot_disk_size_gb = 35
      }
    }
  }
}
	`, rnd)
}

func testAccDataprocCluster_withStagingBucketOnly(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = "true"
}
`, bucketName)
}

func testAccDataprocCluster_withTempBucketOnly(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
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
      machine_type = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
      }
    }
  }
}
`, testAccDataprocCluster_withStagingBucketOnly(bucketName), clusterName)
}

func testAccDataprocCluster_withTempBucketAndCluster(clusterName, bucketName string) string {
	return fmt.Sprintf(`
%s

resource "google_dataproc_cluster" "with_bucket" {
  name   = "%s"
  region = "us-central1"

  cluster_config {
    temp_bucket = google_storage_bucket.bucket.name

    # Keep the costs down with smallest config we can get away with
    software_config {
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    master_config {
      machine_type = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
      }
    }
  }
}
`, testAccDataprocCluster_withTempBucketOnly(bucketName), clusterName)
}

func testAccDataprocCluster_withLabels(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_labels" {
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  labels = {
    key1 = "value1"
  }
}
`, rnd)
}

func testAccDataprocCluster_withEndpointConfig(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_endpoint_config" {
	name                  = "tf-test-%s"
	region                = "us-central1"

	cluster_config {
		endpoint_config {
			enable_http_port_access = "true"
		}
	}
}
`, rnd)
}

func testAccDataprocCluster_withImageVersion(rnd, version string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_image_version" {
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    software_config {
      image_version = "%s"
    }
  }
}
`, rnd, version)
}

func testAccDataprocCluster_withOptionalComponents(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_opt_components" {
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    software_config {
      optional_components = ["DOCKER", "ZOOKEEPER"]
    }
  }
}
`, rnd)
}

func testAccDataprocCluster_withLifecycleConfigIdleDeleteTtl(rnd, tm string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_lifecycle_config" {
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    lifecycle_config {
      idle_delete_ttl = "%s"
    }
  }
}
`, rnd, tm)
}

func testAccDataprocCluster_withLifecycleConfigAutoDeletionTime(rnd, tm string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_lifecycle_config" {
 name   = "tf-test-dproc-%s"
 region = "us-central1"

 cluster_config {
   lifecycle_config {
     auto_delete_time = "%s"
   }
 }
}
`, rnd, tm)
}

func testAccDataprocCluster_withServiceAcc(sa string, rnd string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "service_account" {
  account_id = "%s"
}

resource "google_project_iam_member" "service_account" {
  project = data.google_project.project.project_id
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
      machine_type = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
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
  name          = "tf-test-dproc-%s"
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
  name       = "tf-test-dproc-net-%s"
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
      machine_type = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
      }
    }

    gce_cluster_config {
      network = google_compute_network.dataproc_network.name
    }
  }
}

resource "google_dataproc_cluster" "with_net_ref_by_url" {
  name       = "tf-test-dproc-url-%s"
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
      machine_type = "e2-medium"
      disk_config {
        boot_disk_size_gb = 35
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

  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    encryption_config {
      kms_key_name = "%s"
    }
  }
}
`, pid, rnd, kmsKey)
}

func testAccDataprocCluster_withKerberos(rnd, kmsKey string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-dproc-%s"
  location = "US"
}
resource "google_storage_bucket_object" "password" {
  name = "dataproc-password-%s"
  bucket = google_storage_bucket.bucket.name
  content = "hunter2"
}

resource "google_dataproc_cluster" "kerb" {
  name   = "tf-test-dproc-%s"
  region = "us-central1"

  cluster_config {
    security_config {
      kerberos_config {
        root_principal_password_uri = google_storage_bucket_object.password.self_link
        kms_key_uri = "%s"
      }
    }
  }
}
`, rnd, rnd, rnd, kmsKey)
}

func testAccDataprocCluster_withAutoscalingPolicy(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name     = "tf-test-dataproc-policy-%s"
  region   = "us-central1"

  cluster_config {
    autoscaling_config {
      policy_uri = google_dataproc_autoscaling_policy.asp.id
    }
  }
}

resource "google_dataproc_autoscaling_policy" "asp" {
  policy_id = "tf-test-dataproc-policy-%s"
  location  = "us-central1"

  worker_config {
    max_instances = 3
  }

  basic_algorithm {
    yarn_config {
      graceful_decommission_timeout = "30s"
      scale_up_factor   = 0.5
      scale_down_factor = 0.5
    }
  }
}
`, rnd, rnd)
}

func testAccDataprocCluster_removeAutoscalingPolicy(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
  name     = "tf-test-dataproc-policy-%s"
  region   = "us-central1"

  cluster_config {
    autoscaling_config {
      policy_uri = ""
    }
  }
}

resource "google_dataproc_autoscaling_policy" "asp" {
  policy_id = "tf-test-dataproc-policy-%s"
  location  = "us-central1"

  worker_config {
    max_instances = 3
  }

  basic_algorithm {
    yarn_config {
      graceful_decommission_timeout = "30s"
      scale_up_factor   = 0.5
      scale_down_factor = 0.5
    }
  }
}
`, rnd, rnd)
}

func testAccDataprocCluster_withMetastoreConfig(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_metastore_config" {
	name                  = "tf-test-%s"
	region                = "us-central1"

	cluster_config {
		metastore_config {
			dataproc_metastore_service = google_dataproc_metastore_service.ms.name
		}
	}
}

resource "google_dataproc_metastore_service" "ms" {
	service_id = "metastore-srv"
	location   = "us-central1"
	port       = 9080
	tier       = "DEVELOPER"

	maintenance_window {
		hour_of_day = 2
		day_of_week = "SUNDAY"
	}

	hive_metastore_config {
		version = "3.1.2"
	}
}
`, rnd)
}

func testAccDataprocCluster_withMetastoreConfig_update(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_metastore_config" {
	name                  = "tf-test-%s"
	region                = "us-central1"

	cluster_config {
		metastore_config {
			dataproc_metastore_service = google_dataproc_metastore_service.ms.name
		}
	}
}

resource "google_dataproc_metastore_service" "ms" {
	service_id = "metastore-srv-update"
	location   = "us-central1"
	port       = 9080
	tier       = "DEVELOPER"

	maintenance_window {
		hour_of_day = 2
		day_of_week = "SUNDAY"
	}

	hive_metastore_config {
		version = "3.1.2"
	}
}
`, rnd)
}
