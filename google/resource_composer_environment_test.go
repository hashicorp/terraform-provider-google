package google

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/storage/v1"
)

const testComposerEnvironmentPrefix = "tf-cc-testenv"
const testComposerNetworkPrefix = "tf-cc-testnet"

func init() {
	resource.AddTestSweepers("gcp_composer_environment", &resource.Sweeper{
		Name: "gcp_composer_environment",
		F:    testSweepComposerResources,
	})
}

func TestComposerImageVersionDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{"matches", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1.10.0", true},
		{"old latest", "composer-latest-airflow-1.10.0", "composer-1.4.1-airflow-1.10.0", true},
		{"new latest", "composer-1.4.1-airflow-1.10.0", "composer-latest-airflow-1.10.0", true},
		{"airflow equivalent", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1.10", true},
		{"airflow different", "composer-1.4.0-airflow-1.10.0", "composer-1.4-airflow-1.9.0", false},
		{"airflow different composer latest", "composer-1.4.0-airflow-1.10.0", "composer-latest-airflow-1.9.0", false},
	}

	for _, tc := range cases {
		if actual := composerImageVersionDiffSuppress("", tc.old, tc.new, nil); actual != tc.expected {
			t.Fatalf("'%s' failed, expected %v but got %v", tc.name, tc.expected, actual)
		}
	}
}

// Checks environment creation with minimum required information.
func TestAccComposerEnvironment_basic(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix(testComposerEnvironmentPrefix)
	network := acctest.RandomWithPrefix(testComposerNetworkPrefix)
	subnetwork := network + "-1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_basic(envName, network, subnetwork),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_composer_environment.test", "config.0.airflow_uri"),
					resource.TestCheckResourceAttrSet("google_composer_environment.test", "config.0.gke_cluster"),
					resource.TestCheckResourceAttrSet("google_composer_environment.test", "config.0.node_count"),
					resource.TestCheckResourceAttrSet("google_composer_environment.test", "config.0.node_config.0.zone"),
					resource.TestCheckResourceAttrSet("google_composer_environment.test", "config.0.node_config.0.machine_type")),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", getTestProjectFromEnv(), "us-central1", envName),
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(emilyye): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_basic(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(network),
			},
		},
	})
}

// Checks that all updatable fields can be updated in one apply
// (PATCH for Environments only is per-field)
func TestAccComposerEnvironment_update(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix(testComposerEnvironmentPrefix)
	network := acctest.RandomWithPrefix(testComposerNetworkPrefix)
	subnetwork := network + "-1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_basic(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironment_update(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(emilyye): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_update(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(network),
			},
		},
	})
}

// Checks environment creation with minimum required information.
func TestAccComposerEnvironment_private(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix(testComposerEnvironmentPrefix)
	network := acctest.RandomWithPrefix(testComposerNetworkPrefix)
	subnetwork := network + "-1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_private(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", getTestProjectFromEnv(), "us-central1", envName),
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(emilyye): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_private(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(network),
			},
		},
	})
}

// Checks behavior of node config, including dependencies on Compute resources.
func TestAccComposerEnvironment_withNodeConfig(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix(testComposerEnvironmentPrefix)
	network := acctest.RandomWithPrefix(testComposerNetworkPrefix)
	subnetwork := network + "-1"
	serviceAccount := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_nodeCfg(envName, network, subnetwork, serviceAccount),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(emilyye): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_nodeCfg(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(network),
			},
		},
	})
}

func TestAccComposerEnvironment_withSoftwareConfig(t *testing.T) {
	t.Parallel()
	envName := acctest.RandomWithPrefix(testComposerEnvironmentPrefix)
	network := acctest.RandomWithPrefix(testComposerNetworkPrefix)
	subnetwork := network + "-1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_softwareCfg(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(emilyye): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_softwareCfg(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(network),
			},
		},
	})
}

// Checks behavior of config for creation for attributes that must
// be updated during create.
func TestAccComposerEnvironment_withUpdateOnCreate(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix(testComposerEnvironmentPrefix)
	network := acctest.RandomWithPrefix(testComposerNetworkPrefix)
	subnetwork := network + "-1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_updateOnlyFields(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(emilyye): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_updateOnlyFields(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(network),
			},
		},
	})
}

func testAccComposerEnvironmentDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_composer_environment" {
			continue
		}

		idTokens := strings.Split(rs.Primary.ID, "/")
		if len(idTokens) != 6 {
			return fmt.Errorf("Invalid ID %q, expected format projects/{project}/regions/{region}/environments/{environment}", rs.Primary.ID)
		}
		envName := &composerEnvironmentName{
			Project:     idTokens[1],
			Region:      idTokens[3],
			Environment: idTokens[5],
		}

		_, err := config.clientComposer.Projects.Locations.Environments.Get(envName.resourceName()).Do()
		if err == nil {
			return fmt.Errorf("environment %s still exists", envName.resourceName())
		}
	}

	return nil
}

func testAccComposerEnvironment_basic(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"
    }
  }
}

// use a separate network to avoid conflicts with other tests running in parallel
// that use the default network/subnet
resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test.self_link
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_private(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"
      ip_allocation_policy {
        use_ip_aliases          = true
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }
    private_environment_config {
      enable_private_endpoint = true
    }
  }
}

// use a separate network to avoid conflicts with other tests running in parallel
// that use the default network/subnet
resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name                     = "%s"
  ip_cidr_range            = "10.2.0.0/16"
  region                   = "us-central1"
  network                  = google_compute_network.test.self_link
  private_ip_google_access = true
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_update(name, network, subnetwork string) string {
	return fmt.Sprintf(`
data "google_composer_image_versions" "all" {
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_count = 4
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"
    }

    software_config {
      image_version = data.google_composer_image_versions.all.image_versions[0].image_version_id

      airflow_config_overrides = {
        core-load_example = "True"
      }

      pypi_packages = {
        numpy = ""
      }

      env_variables = {
        FOO = "bar"
      }
    }
  }

  labels = {
    foo          = "bar"
    anotherlabel = "boo"
  }
}

// use a separate network to avoid conflicts with other tests running in parallel
// that use the default network/subnet
resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test.self_link
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_nodeCfg(environment, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"

      service_account = google_service_account.test.name
    }
  }

  depends_on = [google_project_iam_member.composer-worker]
}

resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test.self_link
}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}
`, environment, network, subnetwork, serviceAccount)
}

func testAccComposerEnvironment_softwareCfg(name, network, subnetwork string) string {
	return fmt.Sprintf(`
data "google_composer_image_versions" "all" {
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"
    }
    software_config {
      image_version  = data.google_composer_image_versions.all.image_versions[0].image_version_id
      python_version = "3"
    }
  }
}

// use a separate network to avoid conflicts with other tests running in parallel
// that use the default network/subnet
resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test.self_link
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_updateOnlyFields(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"
    }
    software_config {
      pypi_packages = {
        scipy = "==1.1.0"
      }
    }
  }
}

// use a separate network to avoid conflicts with other tests running in parallel
// that use the default network/subnet
resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test.self_link
}
`, name, network, subnetwork)
}

/**
 * CLEAN UP HELPER FUNCTIONS
 * Because the environments are flaky and bucket deletion rates can be
 * rate-limited, for now just warn instead of returning actual errors.
 */
func testSweepComposerResources(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate()
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	// Environments need to be cleaned up because the service is flaky.
	if err := testSweepComposerEnvironments(config); err != nil {
		log.Printf("[WARNING] unable to clean up all environments: %s", err)
	}

	// Buckets need to be cleaned up because they just don't get deleted on purpose.
	if err := testSweepComposerEnvironmentBuckets(config); err != nil {
		log.Printf("[WARNING] unable to clean up all environment storage buckets: %s", err)
	}

	return nil
}

func testSweepComposerEnvironments(config *Config) error {
	found, err := config.clientComposer.Projects.Locations.Environments.List(
		fmt.Sprintf("projects/%s/locations/%s", config.Project, config.Region)).Do()
	if err != nil {
		return fmt.Errorf("error listing storage buckets for composer environment: %s", err)
	}

	if len(found.Environments) == 0 {
		log.Printf("No environment need to be cleaned up")
		return nil
	}

	var allErrors error
	for _, e := range found.Environments {
		createdAt, err := time.Parse(time.RFC3339Nano, e.CreateTime)
		if err != nil {
			return fmt.Errorf("[ERROR] Environment %q has invalid create time %q", e.Name, e.CreateTime)
		}
		// Skip environments that were created in same day
		// This sweeper should really only clean out very old environments.
		if time.Since(createdAt) < time.Hour*24 {
			continue
		}

		switch e.State {
		case "CREATING":
		case "UPDATING":
			log.Printf("Skipping pending Environment %q with state %q", e.Name, e.State)
		case "DELETING":
			log.Printf("Skipping pending Environment %q that is currently deleting", e.Name)
		case "RUNNING":
		case "ERROR":
		default:
			op, deleteErr := config.clientComposer.Projects.Locations.Environments.Delete(e.Name).Do()
			if deleteErr != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("Unable to delete environment %q: %s", e.Name, deleteErr))
				continue
			}
			waitErr := composerOperationWaitTime(config.clientComposer, op, config.Project, "Sweeping old test environments", 10)
			if waitErr != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("Unable to delete environment %q: %s", e.Name, waitErr))
			}
		}
	}
	return allErrors
}

func testSweepComposerEnvironmentBuckets(config *Config) error {
	artifactsBName := fmt.Sprintf("artifacts.%s.appspot.com", config.Project)
	artifactBucket, err := config.clientStorage.Buckets.Get(artifactsBName).Do()
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			log.Printf("composer environment bucket %q not found, doesn't need to be clean up", artifactsBName)
		} else {
			return err
		}
	} else if err = testSweepComposerEnvironmentCleanUpBucket(config, artifactBucket); err != nil {
		return err
	}

	found, err := config.clientStorage.Buckets.List(config.Project).Prefix(config.Region).Do()
	if err != nil {
		return fmt.Errorf("error listing storage buckets created when testing composer environment: %s", err)
	}
	if len(found.Items) == 0 {
		log.Printf("No environment-specific buckets need to be cleaned up")
		return nil
	}

	for _, bucket := range found.Items {
		if _, ok := bucket.Labels["goog-composer-environment"]; !ok {
			continue
		}
		if err := testSweepComposerEnvironmentCleanUpBucket(config, bucket); err != nil {
			return err
		}
	}
	return nil
}

func testSweepComposerEnvironmentCleanUpBucket(config *Config, bucket *storage.Bucket) error {
	var allErrors error
	objList, err := config.clientStorage.Objects.List(bucket.Name).Do()
	if err != nil {
		allErrors = multierror.Append(allErrors,
			fmt.Errorf("Unable to list objects to delete for bucket %q: %s", bucket.Name, err))
	}

	for _, o := range objList.Items {
		if err := config.clientStorage.Objects.Delete(bucket.Name, o.Name).Do(); err != nil {
			allErrors = multierror.Append(allErrors,
				fmt.Errorf("Unable to delete object %q from bucket %q: %s", o.Name, bucket.Name, err))
		}
	}

	if err := config.clientStorage.Buckets.Delete(bucket.Name).Do(); err != nil {
		allErrors = multierror.Append(allErrors, fmt.Errorf("Unable to delete bucket %q: %s", bucket.Name, err))
	}

	if allErrors != nil {
		return fmt.Errorf("Unable to clean up bucket %q: %v", bucket.Name, allErrors)
	}

	log.Printf("Cleaned up bucket %q for composer environment tests", bucket.Name)
	return nil
}

// WARNING: This is not actually a check and is a terrible clean-up step because Composer Environments
// have a bug that hasn't been fixed. Composer will add firewalls to non-default networks for environments
// but will not remove them when the Environment is deleted.
//
// Destroy test step for config with a network will fail unless we clean up the firewalls before.
func testAccCheckClearComposerEnvironmentFirewalls(networkName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		config.Project = getTestProjectFromEnv()
		network, err := config.clientCompute.Networks.Get(getTestProjectFromEnv(), networkName).Do()
		if err != nil {
			return err
		}

		foundFirewalls, err := config.clientCompute.Firewalls.List(config.Project).Do()
		if err != nil {
			return fmt.Errorf("Unable to list firewalls for network %q: %s", network.Name, err)
		}

		var allErrors error
		for _, firewall := range foundFirewalls.Items {
			if !strings.HasPrefix(firewall.Name, testComposerNetworkPrefix) {
				continue
			}
			log.Printf("[DEBUG] Deleting firewall %q for test-resource network %q", firewall.Name, network.Name)
			op, err := config.clientCompute.Firewalls.Delete(config.Project, firewall.Name).Do()
			if err != nil {
				allErrors = multierror.Append(allErrors,
					fmt.Errorf("Unable to delete firewalls for network %q: %s", network.Name, err))
				continue
			}

			waitErr := computeOperationWaitTime(config, op, config.Project,
				"Sweeping test composer environment firewalls", 10)
			if waitErr != nil {
				allErrors = multierror.Append(allErrors,
					fmt.Errorf("Error while waiting to delete firewall %q: %s", firewall.Name, waitErr))
			}
		}
		return allErrors
	}
}
