package google

import (
	"context"
	"fmt"
	"testing"

	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/storage/v1"
)

const testComposerEnvironmentPrefix = "tf-test-composer-env"
const testComposerNetworkPrefix = "tf-test-composer-net"

func init() {
	resource.AddTestSweepers("gcp_composer_environment", &resource.Sweeper{
		Name: "gcp_composer_environment",
		F:    testSweepComposerResources,
	})
}

func allComposerServiceAgents() []string {
	return []string{
		"cloudcomposer-accounts",
		"compute-system",
		"container-engine-robot",
		"gcp-sa-artifactregistry",
		"gcp-sa-pubsub",
	}
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
		{"preview matches", "composer-1.17.0-preview.0-airflow-2.0.1", "composer-1.17.0-preview.0-airflow-2.0.1", true},
		{"old latest", "composer-latest-airflow-1.10.0", "composer-1.4.1-airflow-1.10.0", true},
		{"new latest", "composer-1.4.1-airflow-1.10.0", "composer-latest-airflow-1.10.0", true},
		{"composer major alias equivalent", "composer-1.4.0-airflow-1.10.0", "composer-1-airflow-1.10", true},
		{"composer major alias different", "composer-1.4.0-airflow-2.1.4", "composer-2-airflow-2.2", false},
		{"composer different", "composer-1.4.0-airflow-1.10.0", "composer-1.4.1-airflow-1.10.0", false},
		{"airflow major alias equivalent", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1", true},
		{"airflow major alias different", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-2", false},
		{"airflow major.minor alias equivalent", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1.10", true},
		{"airflow major.minor alias different", "composer-1.4.0-airflow-2.1.4", "composer-1.4.0-airflow-2.2", false},
		{"airflow different", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1.9.0", false},
	}

	for _, tc := range cases {
		if actual := composerImageVersionDiffSuppress("", tc.old, tc.new, nil); actual != tc.expected {
			t.Errorf("'%s' failed, expected %v but got %v", tc.name, tc.expected, actual)
		}
	}
}

// Checks environment creation with minimum required information.
func TestAccComposerEnvironment_basic(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
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
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", GetTestProjectFromEnv(), "us-central1", envName),
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_basic(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks that all updatable fields can be updated in one apply
// (PATCH for Environments only is per-field)
func TestAccComposerEnvironment_update(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
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
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_update(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks private environment creation for composer 1 and 2.
func TestAccComposerEnvironmentComposer1_private(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer1_private(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", GetTestProjectFromEnv(), "us-central1", envName),
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironmentComposer1_private(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentComposer2_private(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer2_private(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", GetTestProjectFromEnv(), "us-central1", envName),
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironmentComposer2_private(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks environment creation with minimum required information.
func TestAccComposerEnvironment_privateWithWebServerControl(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_privateWithWebServerControl(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComposerEnvironment_privateWithWebServerControlUpdated(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", GetTestProjectFromEnv(), "us-central1", envName),
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_privateWithWebServerControlUpdated(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_withDatabaseConfig(t *testing.T) {
	t.Parallel()
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_databaseCfg(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironment_databaseCfgUpdated(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_databaseCfgUpdated(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_withWebServerConfig(t *testing.T) {
	t.Parallel()
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	grantServiceAgentsRole(t, "service-", []string{"gcp-sa-cloudbuild"}, "roles/cloudbuild.builds.builder")

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_webServerCfg(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironment_webServerCfgUpdated(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_webServerCfgUpdated(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_withEncryptionConfigComposer1(t *testing.T) {
	t.Parallel()

	kms := BootstrapKMSKeyInLocation(t, "us-central1")
	pid := GetTestProjectFromEnv()
	grantServiceAgentsRole(t, "service-", allComposerServiceAgents(), "roles/cloudkms.cryptoKeyEncrypterDecrypter")
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_encryptionCfg(pid, "1", "1", envName, kms.CryptoKey.Name, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_encryptionCfg(pid, "1", "1", envName, kms.CryptoKey.Name, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_withEncryptionConfigComposer2(t *testing.T) {
	t.Parallel()

	kms := BootstrapKMSKeyInLocation(t, "us-central1")
	pid := GetTestProjectFromEnv()
	grantServiceAgentsRole(t, "service-", allComposerServiceAgents(), "roles/cloudkms.cryptoKeyEncrypterDecrypter")
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_encryptionCfg(pid, "2", "2", envName, kms.CryptoKey.Name, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_encryptionCfg(pid, "2", "2", envName, kms.CryptoKey.Name, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_withMaintenanceWindow(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_maintenanceWindow(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_maintenanceWindow(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_maintenanceWindowUpdate(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_maintenanceWindow(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironment_maintenanceWindowUpdate(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_maintenanceWindowUpdate(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_ComposerV2(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_composerV2(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_UpdateComposerV2(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironment_updateComposerV2(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_updateComposerV2(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_composerV2PrivateServiceConnect(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2PrivateServiceConnect(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_composerV2PrivateServiceConnect(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_composerV1MasterAuthNetworks(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_MasterAuthNetworks("1", "1", envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_MasterAuthNetworks("1", "1", envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_composerV2MasterAuthNetworks(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_MasterAuthNetworks("2", "2", envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_MasterAuthNetworks("2", "2", envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_composerV1MasterAuthNetworksUpdate(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_MasterAuthNetworks("1", "1", envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironment_MasterAuthNetworksUpdate("1", "1", envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_MasterAuthNetworksUpdate("1", "1", envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_composerV2MasterAuthNetworksUpdate(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_MasterAuthNetworks("2", "2", envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironment_MasterAuthNetworksUpdate("2", "2", envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO(dzarmola): Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_MasterAuthNetworksUpdate("2", "2", envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks behavior of node config, including dependencies on Compute resources.
func TestAccComposerEnvironment_withNodeConfig(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
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
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_nodeCfg(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentAirflow2_withRecoveryConfig(t *testing.T) {
	t.Parallel()
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_airflow2RecoveryCfg(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComposerEnvironmentUpdate_airflow2RecoveryCfg(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironmentUpdate_airflow2RecoveryCfg(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_withSoftwareConfig(t *testing.T) {
	t.Parallel()
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
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
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_softwareCfg(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentAirflow2_withSoftwareConfig(t *testing.T) {
	t.Parallel()
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_airflow2SoftwareCfg(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComposerEnvironmentUpdate_airflow2SoftwareCfg(envName, network, subnetwork),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironmentUpdate_airflow2SoftwareCfg(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks behavior of config for creation for attributes that must
// be updated during create.
func TestAccComposerEnvironment_withUpdateOnCreate(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
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
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironment_updateOnlyFields(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_fixPyPiPackages(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComposerEnvironment_fixPyPiPackages(envName, network, subnetwork, serviceAccount),
				ExpectError: regexp.MustCompile("Failed to install pypi packages"),
			},
			{
				Config: testAccComposerEnvironment_fixPyPiPackagesUpdate(envName, network, subnetwork, serviceAccount),
			},
			{
				ResourceName:      "google_composer_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This bootstraps the IAM roles needed for the service agents.
func grantServiceAgentsRole(t *testing.T, prefix string, agentNames []string, role string) {
	if BootstrapAllPSARole(t, prefix, agentNames, role) {
		// Fail this test run because the policy needs time to reconcile.
		t.Fatal("Stopping test because permissions were added.")
	}
}

func testAccComposerEnvironmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := GoogleProviderConfig(t)

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

			_, err := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.Get(envName.resourceName()).Do()
			if err == nil {
				return fmt.Errorf("environment %s still exists", envName.resourceName())
			}
		}

		return nil
	}
}

func testAccComposerEnvironment_basic(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network        = google_compute_network.test.self_link
      subnetwork     = google_compute_subnetwork.test.self_link
      zone           = "us-central1-a"
      machine_type  = "n1-standard-1"
      ip_allocation_policy {
        use_ip_aliases          = true
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }
    software_config {
      image_version = "composer-1-airflow-2.3"
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

func testAccComposerEnvironmentComposer1_private(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"
      enable_ip_masq_agent = true
      ip_allocation_policy {
        use_ip_aliases          = true
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }
    software_config {
      image_version = "composer-1-airflow-2"
    }
    private_environment_config {
      enable_private_endpoint = true
      enable_privately_used_public_ips = true
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

func testAccComposerEnvironmentComposer2_private(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      enable_ip_masq_agent = true
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }
    software_config {
      image_version  = "composer-2-airflow-2"
    }
    private_environment_config {
      enable_private_endpoint = true
      enable_privately_used_public_ips = true
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

func testAccComposerEnvironment_privateWithWebServerControl(name, network, subnetwork string) string {
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
        cluster_ipv4_cidr_block = "10.56.0.0/14"
        services_ipv4_cidr_block = "10.122.0.0/20"
      }
    }
    private_environment_config {
      enable_private_endpoint = false
      web_server_ipv4_cidr_block = "172.30.240.0/24"
      cloud_sql_ipv4_cidr_block = "10.32.0.0/12"
      master_ipv4_cidr_block =  "172.17.50.0/28"
    }
    software_config {
      image_version = "composer-1-airflow-2"
    }
    web_server_network_access_control {
      allowed_ip_range {
        value = "192.168.0.1"
        description = "my range1"
      }
      allowed_ip_range {
        value = "0.0.0.0/0"
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
  name                     = "%s"
  ip_cidr_range            = "10.2.0.0/16"
  region                   = "us-central1"
  network                  = google_compute_network.test.self_link
  private_ip_google_access = true
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_privateWithWebServerControlUpdated(name, network, subnetwork string) string {
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
        cluster_ipv4_cidr_block = "10.56.0.0/14"
        services_ipv4_cidr_block = "10.122.0.0/20"
      }
    }
    private_environment_config {
      enable_private_endpoint = false
      web_server_ipv4_cidr_block = "172.30.240.0/24"
      cloud_sql_ipv4_cidr_block = "10.32.0.0/12"
      master_ipv4_cidr_block =  "172.17.50.0/28"
    }
    software_config {
      image_version = "composer-1-airflow-2"
    }
    web_server_network_access_control {
      allowed_ip_range {
        value = "192.168.0.1"
        description = "my range1"
      }
      allowed_ip_range {
        value = "0.0.0.0/0"
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
  name                     = "%s"
  ip_cidr_range            = "10.2.0.0/16"
  region                   = "us-central1"
  network                  = google_compute_network.test.self_link
  private_ip_google_access = true
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_databaseCfg(name, network, subnetwork string) string {
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
    database_config {
      machine_type  = "db-n1-standard-4"
    }
    software_config {
      image_version = "composer-1-airflow-2"
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

func testAccComposerEnvironment_databaseCfgUpdated(name, network, subnetwork string) string {
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
    database_config {
      machine_type  = "db-n1-standard-8"
    }
    software_config {
      image_version = "composer-1-airflow-2"
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

func testAccComposerEnvironment_webServerCfg(name, network, subnetwork string) string {
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
      image_version = "composer-1-airflow-2"
    }
    web_server_config {
      machine_type  = "composer-n1-webserver-4"
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

func testAccComposerEnvironment_webServerCfgUpdated(name, network, subnetwork string) string {
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
      image_version = "composer-1-airflow-2"
    }
    web_server_config {
      machine_type  = "composer-n1-webserver-8"
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

func testAccComposerEnvironment_encryptionCfg(pid, compVersion, airflowVersion, name, kmsKey, network, subnetwork string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = "%s"
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@gs-project-accounts.iam.gserviceaccount.com"
}
resource "google_composer_environment" "test" {
  depends_on = [google_kms_crypto_key_iam_member.iam]
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
    }

    software_config {
      image_version  = "composer-%s-airflow-%s"
    }

    encryption_config {
      kms_key_name  = "%s"
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
`,
		pid, kmsKey, name, compVersion, airflowVersion, kmsKey, network, subnetwork)
}

func testAccComposerEnvironment_maintenanceWindow(envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    maintenance_window {
      start_time = "2019-08-01T01:00:00Z"
      end_time = "2019-08-01T07:00:00Z"
      recurrence = "FREQ=WEEKLY;BYDAY=TU,WE"
    }
  }
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

`, envName, network, subnetwork)
}

func testAccComposerEnvironment_maintenanceWindowUpdate(envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    maintenance_window {
      start_time = "2019-08-01T01:00:00Z"
      end_time = "2019-08-01T07:00:00Z"
      recurrence = "FREQ=DAILY"
    }
  }
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

`, envName, network, subnetwork)
}

func testAccComposerEnvironment_composerV2WithDisabledTriggerer(envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }

    software_config {
      image_version = "composer-2-airflow-2"
    }

    workloads_config {
      scheduler {
        cpu          = 1.25
        memory_gb    = 2.5
        storage_gb   = 5.4
        count        = 2
      }
      web_server {
        cpu          = 1.75
        memory_gb    = 3.0
        storage_gb   = 4.4
      }
      worker {
        cpu          = 0.5
        memory_gb    = 2.0
        storage_gb   = 3.4
        min_count    = 2
        max_count    = 5
      }
    }
    environment_size = "ENVIRONMENT_SIZE_MEDIUM"
    private_environment_config {
      enable_private_endpoint                 = true
      cloud_composer_network_ipv4_cidr_block   = "10.3.192.0/24"
      master_ipv4_cidr_block                   = "172.16.194.0/23"
      cloud_sql_ipv4_cidr_block               = "10.3.224.0/20"
    }
  }

}

resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-east1"
   network       = google_compute_network.test.self_link
  private_ip_google_access = true
}

`, envName, network, subnetwork)
}

func testAccComposerEnvironment_composerV2(envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

    config {
      node_config {
        network          = google_compute_network.test.self_link
        subnetwork       = google_compute_subnetwork.test.self_link
        ip_allocation_policy {
          cluster_ipv4_cidr_block = "10.0.0.0/16"
        }
      }

      software_config {
        image_version = "composer-2-airflow-2"
      }

      workloads_config {
        scheduler {
          cpu         = 1.25
          memory_gb   = 2.5
          storage_gb   = 5.4
          count       = 2
        }
        web_server {
          cpu         = 1.75
          memory_gb   = 3.0
          storage_gb   = 4.4
        }
        worker {
          cpu         = 0.5
          memory_gb   = 2.0
          storage_gb   = 3.4
          min_count   = 2
          max_count   = 5
        }
				      }
      environment_size = "ENVIRONMENT_SIZE_MEDIUM"
      private_environment_config {
        enable_private_endpoint                 = true
        cloud_composer_network_ipv4_cidr_block   = "10.3.192.0/24"
        master_ipv4_cidr_block                   = "172.16.194.0/23"
        cloud_sql_ipv4_cidr_block               = "10.3.224.0/20"
      }
    }

}

resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-east1"
   network       = google_compute_network.test.self_link
  private_ip_google_access = true
}

`, envName, network, subnetwork)
}

func testAccComposerEnvironment_composerV2PrivateServiceConnect(envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

    config {
      node_config {
        network          = google_compute_network.test.self_link
        subnetwork       = google_compute_subnetwork.test.self_link
      }

      software_config {
        image_version = "composer-2-airflow-2"
      }
      private_environment_config {
        cloud_composer_connection_subnetwork    = google_compute_subnetwork.test.self_link
      }
    }

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
  private_ip_google_access = true
}

`, envName, network, subnetwork)
}

func testAccComposerEnvironment_MasterAuthNetworks(compVersion, airflowVersion, envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_config {
      network        = google_compute_network.test.self_link
      subnetwork     = google_compute_subnetwork.test.self_link
    }

    software_config {
      image_version = "composer-%s-airflow-%s"
    }

    master_authorized_networks_config {
      enabled  = true
      cidr_blocks {
        display_name  = "foo"
        cidr_block    = "8.8.8.8/32"
      }
      cidr_blocks {
        cidr_block    = "8.8.8.0/24"
      }
    }
  }
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

`, envName, compVersion, airflowVersion, network, subnetwork)
}

func testAccComposerEnvironment_MasterAuthNetworksUpdate(compVersion, airflowVersion, envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_config {
      network        = google_compute_network.test.self_link
      subnetwork     = google_compute_subnetwork.test.self_link
    }

    software_config {
      image_version = "composer-%s-airflow-%s"
    }

    master_authorized_networks_config {
      enabled  = true
      cidr_blocks {
        display_name  = "foo_update"
        cidr_block    = "9.9.9.8/30"
      }
    }
  }
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

`, envName, compVersion, airflowVersion, network, subnetwork)
}

func testAccComposerEnvironment_update(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_count = 4
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"
      machine_type  = "n1-standard-1"
      ip_allocation_policy {
        use_ip_aliases          = true
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }

    software_config {
      image_version = "composer-1-airflow-2"

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

func testAccComposerEnvironment_updateComposerV2(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

    config {
      node_config {
        network          = google_compute_network.test.self_link
        subnetwork       = google_compute_subnetwork.test.self_link
        ip_allocation_policy {
          cluster_ipv4_cidr_block = "10.0.0.0/16"
        }
      }

      software_config {
        image_version = "composer-2-airflow-2"
      }

      workloads_config {
        scheduler {
          cpu         = 2.25
          memory_gb   = 3.5
          storage_gb  = 6.4
          count       = 3
        }
        web_server {
          cpu         = 2.75
          memory_gb   = 4.0
          storage_gb  = 5.4
        }
        worker {
          cpu         = 1.5
          memory_gb   = 3.0
          storage_gb  = 4.4
          min_count   = 3
          max_count   = 6
        }
				      }
      environment_size = "ENVIRONMENT_SIZE_LARGE"
      private_environment_config {
        enable_private_endpoint                 = true
        cloud_composer_network_ipv4_cidr_block  = "10.3.192.0/24"
        master_ipv4_cidr_block                  = "172.16.194.0/23"
        cloud_sql_ipv4_cidr_block               = "10.3.224.0/20"
      }
    }

}

resource "google_compute_network" "test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.test.self_link
  private_ip_google_access = true
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_nodeCfg(environment, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link
      zone       = "us-central1-a"

      service_account = google_service_account.test.name
      ip_allocation_policy {
        use_ip_aliases          = true
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }
    software_config {
      image_version = "composer-1-airflow-2"
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
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}
`, environment, network, subnetwork, serviceAccount)
}

func testAccComposerEnvironment_airflow2RecoveryCfg(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }

    software_config {
      image_version = "composer-2-airflow-2"
    }

    recovery_config {
      scheduled_snapshots_config {
        enabled =                    true
        snapshot_location =          "gs://example-bucket/environment_snapshots"
        snapshot_creation_schedule = "0 4 * * *"
        time_zone =                  "UTC+01"
      }
    }
  }

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
  private_ip_google_access = true
}
`, name, network, subnetwork)
}

func testAccComposerEnvironmentUpdate_airflow2RecoveryCfg(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }

    software_config {
      image_version = "composer-2-airflow-2"
    }

    recovery_config {
		  scheduled_snapshots_config {
			  enabled =                    true
			  snapshot_location =          "gs://example-bucket/environment_snapshots2"
			  snapshot_creation_schedule = "1 2 * * *"
			  time_zone =                  "UTC+02"
      }
    }
  }

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
  private_ip_google_access = true
}
`, name, network, subnetwork)
}

func testAccComposerEnvironment_softwareCfg(name, network, subnetwork string) string {
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
      image_version  = "composer-1-airflow-1"
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
      image_version = "composer-1-airflow-2"
      pypi_packages = {
        numpy = ""
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

func testAccComposerEnvironment_airflow2SoftwareCfg(name, network, subnetwork string) string {
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
      image_version  = "composer-1-airflow-2"
      scheduler_count = 2
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

func testAccComposerEnvironmentUpdate_airflow2SoftwareCfg(name, network, subnetwork string) string {
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
      image_version  = "composer-1-airflow-2"
      scheduler_count = 3
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

func testAccComposerEnvironment_fixPyPiPackages(environment, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {

    software_config {
      image_version = "composer-2-airflow-2"

      pypi_packages = {
        "google-cloud-bigquery" = "==1"
      }
    }

    private_environment_config {
      enable_private_endpoint = true
      master_ipv4_cidr_block  = "10.10.0.0/28"
    }

    workloads_config {
      scheduler {
        cpu        = 0.5
        memory_gb  = 1.875
        storage_gb = 1
        count      = 1
      }
      web_server {
        cpu        = 0.5
        memory_gb  = 1.875
        storage_gb = 1
      }
      worker {
        cpu = 0.5
        memory_gb  = 1.875
        storage_gb = 1
        min_count  = 1
        max_count  = 3
      }
    }

    environment_size = "ENVIRONMENT_SIZE_SMALL"

    node_config {
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
      service_account = google_service_account.test.name
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

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

data "google_project" "project" {}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role    = "roles/composer.worker"
  member  = "serviceAccount:${google_service_account.test.email}"
}`, environment, network, subnetwork, serviceAccount)
}

func testAccComposerEnvironment_fixPyPiPackagesUpdate(environment, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {

    software_config {
      image_version = "composer-2-airflow-2"
    }

    private_environment_config {
      enable_private_endpoint = true
      master_ipv4_cidr_block  = "10.10.0.0/28"
    }

    workloads_config {
      scheduler {
        cpu        = 0.5
        memory_gb  = 1.875
        storage_gb = 1
        count      = 1
      }
      web_server {
        cpu        = 0.5
        memory_gb  = 1.875
        storage_gb = 1
      }
      worker {
        cpu = 0.5
        memory_gb  = 1.875
        storage_gb = 1
        min_count  = 1
        max_count  = 3
      }
    }

    environment_size = "ENVIRONMENT_SIZE_SMALL"

    node_config {
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
      service_account = google_service_account.test.name
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

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

data "google_project" "project" {}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role    = "roles/composer.worker"
  member  = "serviceAccount:${google_service_account.test.email}"
}
`, environment, network, subnetwork, serviceAccount)
}

/**
 * CLEAN UP HELPER FUNCTIONS
 * Because the environments are flaky and bucket deletion rates can be
 * rate-limited, for now just warn instead of returning actual errors.
 */
func testSweepComposerResources(region string) error {
	config, err := SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	// us-central is passed as the region for our sweepers, but there are also
	// many tests that use the us-east1 region
	regions := []string{"us-central1", "us-east1"}
	for _, r := range regions {
		// Environments need to be cleaned up because the service is flaky.
		if err := testSweepComposerEnvironments(config, r); err != nil {
			log.Printf("[WARNING] unable to clean up all environments: %s", err)
		}

		// Buckets need to be cleaned up because they just don't get deleted on purpose.
		if err := testSweepComposerEnvironmentBuckets(config, r); err != nil {
			log.Printf("[WARNING] unable to clean up all environment storage buckets: %s", err)
		}
	}

	return nil
}

func testSweepComposerEnvironments(config *Config, region string) error {
	found, err := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.List(
		fmt.Sprintf("projects/%s/locations/%s", config.Project, region)).Do()
	if err != nil {
		return fmt.Errorf("error listing storage buckets for composer environment: %s", err)
	}

	if len(found.Environments) == 0 {
		log.Printf("composer: no environments need to be cleaned up")
		return nil
	}

	log.Printf("composer: %d environments need to be cleaned up", len(found.Environments))

	var allErrors error
	for _, e := range found.Environments {
		createdAt, err := time.Parse(time.RFC3339Nano, e.CreateTime)
		if err != nil {
			return fmt.Errorf("composer: environment %q has invalid create time %q", e.Name, e.CreateTime)
		}
		// Skip environments that were created in same day
		// This sweeper should really only clean out very old environments.
		if time.Since(createdAt) < time.Hour*24 {
			log.Printf("composer: skipped environment %q, it was created today", e.Name)
			continue
		}

		switch e.State {
		case "CREATING":
			fallthrough
		case "UPDATING":
			log.Printf("composer: skipping pending Environment %q with state %q", e.Name, e.State)
		case "DELETING":
			log.Printf("composer: skipping pending Environment %q that is currently deleting", e.Name)
		case "RUNNING":
			fallthrough
		case "ERROR":
			fallthrough
		default:
			op, deleteErr := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.Delete(e.Name).Do()
			if deleteErr != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("composer: unable to delete environment %q: %s", e.Name, deleteErr))
				continue
			}
			waitErr := ComposerOperationWaitTime(config, op, config.Project, "Sweeping old test environments", config.UserAgent, 10*time.Minute)
			if waitErr != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("composer: unable to delete environment %q: %s", e.Name, waitErr))
			}
		}
	}
	return allErrors
}

func testSweepComposerEnvironmentBuckets(config *Config, region string) error {
	artifactsBName := fmt.Sprintf("artifacts.%s.appspot.com", config.Project)
	artifactBucket, err := config.NewStorageClient(config.UserAgent).Buckets.Get(artifactsBName).Do()
	if err != nil {
		if IsGoogleApiErrorWithCode(err, 404) {
			log.Printf("composer environment bucket %q not found, doesn't need to be cleaned up", artifactsBName)
		} else {
			return err
		}
	} else if err = testSweepComposerEnvironmentCleanUpBucket(config, artifactBucket); err != nil {
		return err
	}

	found, err := config.NewStorageClient(config.UserAgent).Buckets.List(config.Project).Prefix(region).Do()
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
	objList, err := config.NewStorageClient(config.UserAgent).Objects.List(bucket.Name).Do()
	if err != nil {
		allErrors = multierror.Append(allErrors,
			fmt.Errorf("Unable to list objects to delete for bucket %q: %s", bucket.Name, err))
	}

	for _, o := range objList.Items {
		if err := config.NewStorageClient(config.UserAgent).Objects.Delete(bucket.Name, o.Name).Do(); err != nil {
			allErrors = multierror.Append(allErrors,
				fmt.Errorf("Unable to delete object %q from bucket %q: %s", o.Name, bucket.Name, err))
		}
	}

	if err := config.NewStorageClient(config.UserAgent).Buckets.Delete(bucket.Name).Do(); err != nil {
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
func testAccCheckClearComposerEnvironmentFirewalls(t *testing.T, networkName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GoogleProviderConfig(t)
		config.Project = GetTestProjectFromEnv()
		network, err := config.NewComputeClient(config.UserAgent).Networks.Get(GetTestProjectFromEnv(), networkName).Do()
		if err != nil {
			return err
		}

		foundFirewalls, err := config.NewComputeClient(config.UserAgent).Firewalls.List(config.Project).Do()
		if err != nil {
			return fmt.Errorf("Unable to list firewalls for network %q: %s", network.Name, err)
		}

		var allErrors error
		for _, firewall := range foundFirewalls.Items {
			if !strings.HasPrefix(firewall.Name, testComposerNetworkPrefix) {
				continue
			}
			log.Printf("[DEBUG] Deleting firewall %q for test-resource network %q", firewall.Name, network.Name)
			op, err := config.NewComputeClient(config.UserAgent).Firewalls.Delete(config.Project, firewall.Name).Do()
			if err != nil {
				allErrors = multierror.Append(allErrors,
					fmt.Errorf("Unable to delete firewalls for network %q: %s", network.Name, err))
				continue
			}

			waitErr := ComputeOperationWaitTime(config, op, config.Project,
				"Sweeping test composer environment firewalls", config.UserAgent, 10)
			if waitErr != nil {
				allErrors = multierror.Append(allErrors,
					fmt.Errorf("Error while waiting to delete firewall %q: %s", firewall.Name, waitErr))
			}
		}
		return allErrors
	}
}
