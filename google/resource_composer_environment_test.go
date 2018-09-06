package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/composer/v1"
	"strings"
)

// Checks environment creation with minimum required information.
func TestAccComposerEnvironment_basic(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix("tf-test")
	var env composer.Environment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_basic(envName),
				Check:  testAccCheckComposerEnvironmentExists("google_composer_environment.test", &env),
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
		},
	})

	if env.Config == nil {
		t.Errorf("expected read value to have non-nil config")
	}
	if len(env.Config.AirflowUri) == 0 {
		t.Errorf("expected computed airflow URI value to be set")
	}
}

// Checks that all updatable fields can be updated in one apply
// (PATCH for Environments only is per-field) and that reverting
// config force-updates back to default.
func TestAccComposerEnvironment_update(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix("tf-test")
	var env composer.Environment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_basic(envName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComposerEnvironmentExists("google_composer_environment.test", &env),
				),
			},
			{
				Config: testAccComposerEnvironment_update(envName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComposerEnvironmentExists("google_composer_environment.test", &env),
				),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if env.Config == nil || env.Config.SoftwareConfig == nil {
		t.Fatalf("expected read value to have non-nil config")
	}

	if env.Config.NodeCount != 4 {
		t.Errorf("expected node count to be updated to 4, got %d", env.Config.NodeCount)
	}

	if len(env.Config.SoftwareConfig.PypiPackages) != 1 {
		t.Errorf(`expected PypiPackages to have one key-value { "numpy": "" }, got: %#v`, env.Config.SoftwareConfig.PypiPackages)
	} else if v, ok := env.Config.SoftwareConfig.PypiPackages["numpy"]; !ok || v != "" {
		t.Errorf(`expected PypiPackages to contain { "numpy": "" }, got: %#v`, env.Config.SoftwareConfig.PypiPackages)
	}

	if len(env.Config.SoftwareConfig.AirflowConfigOverrides) != 1 {
		t.Errorf(`expected AirflowConfigOverrides to have one key-value {"core-load_example": "True" }, got: %#v`, env.Config.SoftwareConfig.AirflowConfigOverrides)
	} else if v, ok := env.Config.SoftwareConfig.AirflowConfigOverrides["core-load_example"]; !ok || v != "True" {
		t.Errorf(`expected AirflowConfigOverrides to contain { "core-load_example": "True" }, got: %#v`, env.Config.SoftwareConfig.AirflowConfigOverrides)
	}

	if len(env.Config.SoftwareConfig.EnvVariables) != 1 {
		t.Errorf(`expected EnvVariables to have one key-value { "FOO": "bar" }, got: %#v`, env.Config.SoftwareConfig.EnvVariables)
	} else if v, ok := env.Config.SoftwareConfig.EnvVariables["FOO"]; !ok || v != "bar" {
		t.Errorf(`expected EnvVariables to contain { "FOO": "bar" }, got: %#v`, env.Config.SoftwareConfig.EnvVariables)
	}
}

// Checks behavior of node config, including dependencies on Compute resources.
func TestAccComposerEnvironment_withNodeConfig(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix("tf-test")
	var env composer.Environment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_nodeCfg(envName),
				Check:  testAccCheckComposerEnvironmentExists("google_composer_environment.test", &env),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Checks behavior of config for creation for attributes that must
// be updated during create.
func TestAccComposerEnvironment_withUpdateOnCreate(t *testing.T) {
	t.Parallel()

	envName := acctest.RandomWithPrefix("tf-test")
	var env composer.Environment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_updateOnlyFields(envName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComposerEnvironmentExists("google_composer_environment.test", &env),
				),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if env.Config == nil {
		t.Fatalf("expected read value to have non-nil config")
	}

	if env.Config.SoftwareConfig == nil {
		t.Fatalf("expected non-nil SoftwareConfig")
	}

	if len(env.Config.SoftwareConfig.PypiPackages) != 1 {
		t.Errorf(`expected PypiPackages to have one key-value { "scipy": "==1.1.0" }, got: %#v`, env.Config.SoftwareConfig.PypiPackages)
	} else if v, ok := env.Config.SoftwareConfig.PypiPackages["numpy"]; !ok || v != "==1.1.0" {
		t.Errorf(`expected PypiPackages to contain { "numpy": "" }, got: %#v`, env.Config.SoftwareConfig.PypiPackages)
	}
}

func testAccCheckComposerEnvironmentExists(n string, environment *composer.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		idTokens := strings.Split(rs.Primary.ID, "/")
		if len(idTokens) != 3 {
			return fmt.Errorf("Invalid ID %q, expected format {project}/{region}/{environment}", rs.Primary.ID)
		}
		envName := &composerEnvironmentName{
			Project:     idTokens[0],
			Region:      idTokens[1],
			Environment: idTokens[2],
		}

		nameFromId := envName.resourceName()
		config := testAccProvider.Meta().(*Config)

		found, err := config.clientComposer.Projects.Locations.Environments.Get(nameFromId).Do()
		if err != nil {
			return err
		}

		if found.Name != nameFromId {
			return fmt.Errorf("Environment not found")
		}

		*environment = *found
		return nil
	}
}

func testAccComposerEnvironmentDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_composer_environment" {
			continue
		}

		idTokens := strings.Split(rs.Primary.ID, "/")
		if len(idTokens) != 3 {
			return fmt.Errorf("Invalid ID %q, expected format {project}/{region}/{environment}", rs.Primary.ID)
		}
		envName := &composerEnvironmentName{
			Project:     idTokens[0],
			Region:      idTokens[1],
			Environment: idTokens[2],
		}

		_, err := config.clientComposer.Projects.Locations.Environments.Get(envName.resourceName()).Do()
		if err == nil {
			return fmt.Errorf("environment %s still exists", envName.resourceName())
		}
	}

	return nil
}

func testAccComposerEnvironment_basic(name string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name           = "%s"
  region         = "us-central1"
}
`, name)
}

func testAccComposerEnvironment_update(name string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
	name = "%s"
	region = "us-central1"

	config {
		node_count = 4

		software_config {
			airflow_config_overrides {
			  core-load_example = "True"
			}

			pypi_packages {
			  numpy = ""
			}

			env_variables {
			   FOO = "bar"
			}
		}
 	}

	labels {
   		foo = "bar"
		anotherlabel = "boo"
 	}
}
`, name)
}

func testAccComposerEnvironment_nodeCfg(name string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
	name = "%s"
	region = "us-central1"
	config {
		node_config {
			network = "${google_compute_network.test.self_link}"
			subnetwork =  "${google_compute_subnetwork.test.self_link}"

			service_account = "${google_service_account.test.name}"
		}
	}

	depends_on = ["google_project_iam_member.composer-worker"]
}

resource "google_compute_network" "test" {
	name 					= "composer-test-network"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
	name          = "composer-test-subnetwork"
	ip_cidr_range = "10.2.0.0/16"
	region        = "us-central1"
	network       = "${google_compute_network.test.self_link}"
}

resource "google_service_account" "test" {
  account_id   = "composer-env-account"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  role    = "roles/composer.worker"
  member  = "serviceAccount:${google_service_account.test.email}"
}
`, name)
}

func testAccComposerEnvironment_updateOnlyFields(name string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
	name = "%s"
	region = "us-central1"
	config {
		software_config {
			pypi_packages {
			  scipy = "==1.1.0"
			}
		}
	}
}
`, name)
}
