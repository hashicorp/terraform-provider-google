// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/composer"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"testing"

	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testComposerEnvironmentPrefix = "tf-test-composer-env"
const testComposerNetworkPrefix = "tf-test-composer-net"
const testComposerBucketPrefix = "tf-test-composer-bucket"
const testComposerNetworkAttachmentPrefix = "tf-test-composer-nta"

func bootstrapComposerServiceAgents(t *testing.T) {
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@cloudcomposer-accounts.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
		{
			Member: "serviceAccount:service-{project_number}@compute-system.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
		{
			Member: "serviceAccount:service-{project_number}@container-engine-robot.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-artifactregistry.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-pubsub.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})
}

// Checks environment creation with minimum required information.
func TestAccComposerEnvironment_basic(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", envvar.GetTestProjectFromEnv(), "us-central1", envName),
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

// Checks private environment creation for composer 1 and 2.
func TestAccComposerEnvironmentComposer1_private(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", envvar.GetTestProjectFromEnv(), "us-central1", envName),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", envvar.GetTestProjectFromEnv(), "us-central1", envName),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/environments/%s", envvar.GetTestProjectFromEnv(), "us-central1", envName),
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
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

func TestAccComposerEnvironment_withEncryptionConfigComposer1(t *testing.T) {
	t.Parallel()

	kms := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-composer1-key1")
	pid := envvar.GetTestProjectFromEnv()
	bootstrapComposerServiceAgents(t)
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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
	acctest.SkipIfVcr(t)
	t.Parallel()

	kms := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-composer2-key1")
	pid := envvar.GetTestProjectFromEnv()
	bootstrapComposerServiceAgents(t)
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2(envName, network, subnetwork, serviceAccount),
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
				Config:             testAccComposerEnvironment_composerV2(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_UpdateComposerV2ImageVersion(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerOldVersion(envName, network, subnetwork, serviceAccount),
			},
			{
				Config: testAccComposerEnvironment_composerNewVersion(envName, network, subnetwork, serviceAccount),
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
				Config:             testAccComposerEnvironment_composerNewVersion(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_UpdateComposerV2ResilienceMode(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2HighResilience(envName, network, subnetwork, serviceAccount),
			},
			{
				Config: testAccComposerEnvironment_updateComposerV2StandardResilience(envName, network, subnetwork, serviceAccount),
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
				Config:             testAccComposerEnvironment_updateComposerV2StandardResilience(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_ComposerV2HighResilience(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2HighResilience(envName, network, subnetwork, serviceAccount),
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
				Config:             testAccComposerEnvironment_composerV2HighResilience(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_UpdateComposerV2WithTriggerer(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2(envName, network, subnetwork, serviceAccount),
			},
			{
				Config: testAccComposerEnvironment_composerV2WithDisabledTriggerer(envName, network, subnetwork, serviceAccount),
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
				Config:             testAccComposerEnvironment_composerV2WithDisabledTriggerer(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_UpdateComposerV2(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_composerV2(envName, network, subnetwork, serviceAccount),
			},
			{
				Config: testAccComposerEnvironment_updateComposerV2(envName, network, subnetwork, serviceAccount),
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
				Config:             testAccComposerEnvironment_updateComposerV2(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_composerV2PrivateServiceConnect(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

func TestAccComposer2Environment_withNodeConfig(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposer2Environment_nodeCfg(envName, network, subnetwork, serviceAccount),
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
				Config:             testAccComposer2Environment_nodeCfg(envName, network, subnetwork, serviceAccount),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentAirflow2_withRecoveryConfig(t *testing.T) {
	t.Parallel()
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
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

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	serviceAccount := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComposerEnvironment_fixPyPiPackages(envName, network, subnetwork, serviceAccount),
				ExpectError: regexp.MustCompile("Failed to install Python packages"),
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

func testAccComposerEnvironmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_composer_environment" {
				continue
			}

			idTokens := strings.Split(rs.Primary.ID, "/")
			if len(idTokens) != 6 {
				return fmt.Errorf("Invalid ID %q, expected format projects/{project}/regions/{region}/environments/{environment}", rs.Primary.ID)
			}
			envName := &composer.ComposerEnvironmentName{
				Project:     idTokens[1],
				Region:      idTokens[3],
				Environment: idTokens[5],
			}

			_, err := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.Get(envName.ResourceName()).Do()
			if err == nil {
				return fmt.Errorf("environment %s still exists", envName.ResourceName())
			}
		}

		return nil
	}
}

// Checks environment creation with custom bucket
func TestAccComposerEnvironment_customBucket(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("%s-%d", testComposerBucketPrefix, acctest.RandInt(t))
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_customBucket(bucketName, envName, network, subnetwork),
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
				Config:             testAccComposerEnvironment_customBucket(bucketName, envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironment_customBucketWithUrl(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("%s-%d", testComposerBucketPrefix, acctest.RandInt(t))
	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironment_customBucketWithUrl(bucketName, envName, network, subnetwork),
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
				Config:             testAccComposerEnvironment_customBucketWithUrl(bucketName, envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks Composer 3 environment creation with new fields.
func TestAccComposerEnvironmentComposer3_basic(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_basic(envName, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_basic(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks Composer 3 specific updatable fields.
func TestAccComposerEnvironmentComposer3_update(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_basic(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironmentComposer3_update(envName, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_update(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentComposer3_withNetworkSubnetworkAndAttachment_expectError(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	networkAttachment := fmt.Sprintf("%s-%d", testComposerNetworkAttachmentPrefix, acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComposerEnvironmentComposer3_withNetworkSubnetworkAndAttachment_expectError(envName, networkAttachment, network, subnetwork),
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Config:             testAccComposerEnvironmentComposer3_basic(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentComposer3_databaseRetention(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_databaseRetention(envName, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_databaseRetention(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentComposer3_withNetworkAttachment(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	networkAttachment := fmt.Sprintf("%s-%d", testComposerNetworkAttachmentPrefix, acctest.RandInt(t))
	fullFormNetworkAttachmentName := fmt.Sprintf("projects/%s/regions/%s/networkAttachments/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), networkAttachment)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_withNetworkAttachment(envName, networkAttachment, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_withNetworkAttachment(envName, fullFormNetworkAttachmentName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccComposerEnvironmentComposer3_updateWithNetworkAttachment(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	networkAttachment := fmt.Sprintf("%s-%d", testComposerNetworkAttachmentPrefix, acctest.RandInt(t))
	fullFormNetworkAttachmentName := fmt.Sprintf("projects/%s/regions/%s/networkAttachments/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), networkAttachment)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_withNetworkAndSubnetwork(envName, networkAttachment, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironmentComposer3_withNetworkAttachment(envName, networkAttachment, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_withNetworkAttachment(envName, fullFormNetworkAttachmentName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccComposerEnvironmentComposer3_updateWithNetworkAndSubnetwork(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	networkAttachment := fmt.Sprintf("%s-%d", testComposerNetworkAttachmentPrefix, acctest.RandInt(t))
	fullFormNetworkAttachmentName := fmt.Sprintf("projects/%s/regions/%s/networkAttachments/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), networkAttachment)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_withNetworkAttachment(envName, networkAttachment, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironmentComposer3_withNetworkAndSubnetwork(envName, networkAttachment, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_withNetworkAttachment(envName, fullFormNetworkAttachmentName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Checks Composer 3 specific updatable fields.
func TestAccComposerEnvironmentComposer3_updateToEmpty(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_basic(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironmentComposer3_empty(envName, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_empty(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

// Checks Composer 3 specific updatable fields.
func TestAccComposerEnvironmentComposer3_updateFromEmpty(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer3_empty(envName, network, subnetwork),
			},
			{
				Config: testAccComposerEnvironmentComposer3_update(envName, network, subnetwork),
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
				Config:             testAccComposerEnvironmentComposer3_update(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentComposer3_upgrade_expectError(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	network := fmt.Sprintf("%s-%d", testComposerNetworkPrefix, acctest.RandInt(t))
	subnetwork := network + "-1"
	errorRegExp, _ := regexp.Compile(".*upgrade to composer 3 is not yet supported.*")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerEnvironmentComposer2_empty(envName, network, subnetwork),
			},
			{
				Config:      testAccComposerEnvironmentComposer3_empty(envName, network, subnetwork),
				ExpectError: errorRegExp,
			},
			// This is a terrible clean-up step in order to get destroy to succeed,
			// due to dangling firewall rules left by the Composer Environment blocking network deletion.
			// TODO: Remove this check if firewall rules bug gets fixed by Composer.
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Config:             testAccComposerEnvironmentComposer2_empty(envName, network, subnetwork),
				Check:              testAccCheckClearComposerEnvironmentFirewalls(t, network),
			},
		},
	})
}

func TestAccComposerEnvironmentComposer2_usesUnsupportedField_expectError(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	errorRegExp, _ := regexp.Compile(".*error in configuration, .* should only be used in Composer 3.*")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComposerEnvironmentComposer2_usesUnsupportedField(envName),
				ExpectError: errorRegExp,
			},
		},
	})
}

func TestAccComposerEnvironmentComposer3_usesUnsupportedField_expectError(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	errorRegExp, _ := regexp.Compile(".*error in configuration, .* should not be used in Composer 3.*")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComposerEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComposerEnvironmentComposer3_usesUnsupportedField(envName),
				ExpectError: errorRegExp,
			},
		},
	})
}

func testAccComposerEnvironment_customBucket(bucketName, envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "test" {
  name   = "%s"
  location = "us-central1"
  force_destroy = true
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network        = google_compute_network.test.self_link
      subnetwork     = google_compute_subnetwork.test.self_link
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }
    software_config {
      image_version = "composer-2-airflow-2"
    }
  }
  storage_config {
    bucket = google_storage_bucket.test.name
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
`, bucketName, envName, network, subnetwork)
}

func testAccComposerEnvironment_customBucketWithUrl(bucketName, envName, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "test" {
  name   = "%s"
  location = "us-central1"
  force_destroy = true
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network        = google_compute_network.test.self_link
      subnetwork     = google_compute_subnetwork.test.self_link
    }
    software_config {
      image_version = "composer-2-airflow-2"
    }
  }
  storage_config {
		bucket = google_storage_bucket.test.url
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
`, bucketName, envName, network, subnetwork)
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
        cluster_ipv4_cidr_block = "10.56.0.0/14"
      }
    }
    software_config {
      image_version  = "composer-2-airflow-2"
    }
    private_environment_config {
      connection_type = "VPC_PEERING"
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

func testAccComposerEnvironment_composerV2WithDisabledTriggerer(envName, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      service_account  = google_service_account.test.name
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
  depends_on = [google_project_iam_member.composer-worker]
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

`, serviceAccount, envName, network, subnetwork)
}

func testAccComposerEnvironment_composerV2(envName, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      service_account  = google_service_account.test.name
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }

    software_config {
      image_version = "composer-2-airflow-2"
      cloud_data_lineage_integration {
        enabled = true
      }
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
      triggerer {
        cpu         = 0.5
        memory_gb   = 2.0
        count         = 1
      }
    }
    database_config {
      zone = "us-east1-c"
    }
    environment_size = "ENVIRONMENT_SIZE_MEDIUM"
    data_retention_config {
      task_logs_retention_config {
        storage_mode = "CLOUD_LOGGING_ONLY"
      }
    }
    private_environment_config {
      enable_private_endpoint                 = true
      cloud_composer_network_ipv4_cidr_block   = "10.3.192.0/24"
      master_ipv4_cidr_block                   = "172.16.194.0/23"
      cloud_sql_ipv4_cidr_block               = "10.3.224.0/20"
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
  region        = "us-east1"
   network       = google_compute_network.test.self_link
  private_ip_google_access = true
}

`, serviceAccount, envName, network, subnetwork)
}

func testAccComposerEnvironment_composerOldVersion(envName, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      service_account  = google_service_account.test.name
    }

    software_config {
      image_version = "composer-2.10.0-airflow-2.10.2"
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
  region        = "us-east1"
   network       = google_compute_network.test.self_link
  private_ip_google_access = true
}

`, serviceAccount, envName, network, subnetwork)
}

func testAccComposerEnvironment_composerNewVersion(envName, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      service_account  = google_service_account.test.name
    }

    software_config {
      image_version = "composer-2.10.1-airflow-2.10.2"
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
  region        = "us-east1"
   network       = google_compute_network.test.self_link
  private_ip_google_access = true
}

`, serviceAccount, envName, network, subnetwork)
}

func testAccComposerEnvironment_composerV2HighResilience(envName, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      service_account  = google_service_account.test.name
    }

    software_config {
      image_version = "composer-2-airflow-2"
    }

    workloads_config {
      scheduler {
        cpu         = 1.25
        memory_gb   = 2.5
        storage_gb  = 5.4
        count       = 2
      }
      web_server {
        cpu         = 1.75
        memory_gb   = 3.0
        storage_gb  = 4.4
      }
      worker {
        cpu         = 0.5
        memory_gb   = 2.0
        storage_gb  = 3.4
        min_count   = 2
        max_count   = 5
      }
    }
    environment_size = "ENVIRONMENT_SIZE_MEDIUM"
    resilience_mode = "HIGH_RESILIENCE"
    private_environment_config {
      enable_private_endpoint                  = true
      cloud_composer_network_ipv4_cidr_block   = "10.3.192.0/24"
      master_ipv4_cidr_block                   = "172.16.194.0/23"
      cloud_sql_ipv4_cidr_block                = "10.3.224.0/20"
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
  region        = "us-east1"
   network       = google_compute_network.test.self_link
  private_ip_google_access = true
}

`, serviceAccount, envName, network, subnetwork)
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

func testAccComposerEnvironment_updateComposerV2StandardResilience(envName, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      service_account  = google_service_account.test.name
    }

    software_config {
      image_version = "composer-2-airflow-2"
    }

    workloads_config {
      scheduler {
        cpu         = 1.25
        memory_gb   = 2.5
        storage_gb  = 5.4
        count       = 2
      }
      web_server {
        cpu         = 1.75
        memory_gb   = 3.0
        storage_gb  = 4.4
      }
      worker {
        cpu         = 0.5
        memory_gb   = 2.0
        storage_gb  = 3.4
        min_count   = 2
        max_count   = 5
      }
    }
    environment_size = "ENVIRONMENT_SIZE_MEDIUM"
    resilience_mode = "STANDARD_RESILIENCE"
    private_environment_config {
      enable_private_endpoint                  = true
      cloud_composer_network_ipv4_cidr_block   = "10.3.192.0/24"
      master_ipv4_cidr_block                   = "172.16.194.0/23"
      cloud_sql_ipv4_cidr_block                = "10.3.224.0/20"
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
  region        = "us-east1"
   network       = google_compute_network.test.self_link
  private_ip_google_access = true
}

`, serviceAccount, envName, network, subnetwork)
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

func testAccComposerEnvironment_updateComposerV2(name, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%s"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-east1"

  config {
    node_config {
      network          = google_compute_network.test.self_link
      subnetwork       = google_compute_subnetwork.test.self_link
      service_account  = google_service_account.test.name
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
    }

    software_config {
      image_version = "composer-2-airflow-2"
      cloud_data_lineage_integration {
      enabled = false
      }
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
      triggerer {
        cpu         = 0.75
        memory_gb   = 2
        count       = 1
      }
    }
    environment_size = "ENVIRONMENT_SIZE_LARGE"
    data_retention_config {
      task_logs_retention_config {
        storage_mode = "CLOUD_LOGGING_AND_CLOUD_STORAGE"
      }
    }
    private_environment_config {
      enable_private_endpoint                 = true
      cloud_composer_network_ipv4_cidr_block  = "10.3.192.0/24"
      master_ipv4_cidr_block                  = "172.16.194.0/23"
      cloud_sql_ipv4_cidr_block               = "10.3.224.0/20"
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
  region        = "us-east1"
  network       = google_compute_network.test.self_link
  private_ip_google_access = true
}
`, serviceAccount, name, network, subnetwork)
}

func testAccComposer2Environment_nodeCfg(environment, network, subnetwork, serviceAccount string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.self_link
      subnetwork = google_compute_subnetwork.test.self_link

      service_account = google_service_account.test.name
      ip_allocation_policy {
        cluster_ipv4_cidr_block = "10.0.0.0/16"
      }
	  tags = toset(["t1", "t2"])
    }
    software_config {
      image_version = "composer-2-airflow-2"
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
  depends_on = [google_project_iam_member.composer-worker]
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
  depends_on = [google_project_iam_member.composer-worker]
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

func testAccComposerEnvironmentComposer2_empty(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    software_config {
      image_version = "composer-2-airflow-2"
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

func testAccComposerEnvironmentComposer3_empty(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
    node_config {
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
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

func testAccComposerEnvironmentComposer2_usesUnsupportedField(name string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    software_config {
      image_version = "composer-2-airflow-2"
	  web_server_plugins_mode = "ENABLED"
    }
  }
}
`, name)
}

func testAccComposerEnvironmentComposer3_usesUnsupportedField(name string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
	  node_config {
	    enable_ip_masq_agent = true
	  }
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
`, name)
}

func testAccComposerEnvironmentComposer3_basic(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      composer_internal_ipv4_cidr_block = "100.64.128.0/20"
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
    }
    software_config {
      image_version = "composer-3-airflow-2"
    }
    workloads_config {
      dag_processor {
        cpu          = 1
        memory_gb    = 2.5
        storage_gb   = 2
        count        = 1
      }
    }
    enable_private_environment = true
    enable_private_builds_only = true
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

func testAccComposerEnvironmentComposer3_update(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test_1.id
      subnetwork = google_compute_subnetwork.test_1.id
      composer_internal_ipv4_cidr_block = "100.64.128.0/20"
    }
    software_config {
      web_server_plugins_mode = "DISABLED"
      image_version = "composer-3-airflow-2"
    }
    workloads_config {
      dag_processor {
        cpu          = 2
        memory_gb    = 2
        storage_gb   = 1
        count        = 2
      }
    }
    enable_private_environment = false
    enable_private_builds_only = false
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

resource "google_compute_network" "test_1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test_1" {
  name          = "%s"
  ip_cidr_range = "10.3.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test_1.self_link
}
`, name, network, subnetwork, network+"-update", subnetwork+"update")
}

func testAccComposerEnvironmentComposer3_withNetworkAttachment(name, networkAttachment, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      composer_network_attachment = google_compute_network_attachment.test.id
    }
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}

resource "google_compute_network_attachment" "test" {
  name = "%s"
  region = "us-central1"
  subnetworks = [ google_compute_subnetwork.test-att.id ]
  connection_preference = "ACCEPT_MANUAL"
  // Composer 3 is modifying producer_accept_lists outside terraform, ignoring this change for now
  lifecycle {
    ignore_changes = [producer_accept_lists]
  }
}

resource "google_compute_network" "test-att" {
  name                    = "%s-att"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test-att" {
  name          = "%s-att"
  ip_cidr_range = "10.3.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test-att.self_link
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
`, name, networkAttachment, network, subnetwork, network, subnetwork)
}

func testAccComposerEnvironmentComposer3_withNetworkAndSubnetwork(name, networkAttachment, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
    }
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}

resource "google_compute_network_attachment" "test" {
  name = "%s"
  region = "us-central1"
  subnetworks = [ google_compute_subnetwork.test-att.id ]
  connection_preference = "ACCEPT_MANUAL"
  // Composer 3 is modifying producer_accept_lists outside terraform, ignoring this change for now
  lifecycle {
    ignore_changes = [producer_accept_lists]
  }
}

resource "google_compute_network" "test-att" {
  name                    = "%s-att"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test-att" {
  name          = "%s-att"
  ip_cidr_range = "10.3.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test-att.self_link
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
`, name, networkAttachment, network, subnetwork, network, subnetwork)
}

func testAccComposerEnvironmentComposer3_withNetworkSubnetworkAndAttachment_expectError(name, networkAttachment, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    node_config {
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
			composer_network_attachment = google_compute_network_attachment.test.id
    }
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}

resource "google_compute_network_attachment" "test" {
  name = "%s"
  region = "us-central1"
  subnetworks = [ google_compute_subnetwork.test.id ]
  connection_preference = "ACCEPT_MANUAL"
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
`, name, networkAttachment, network, subnetwork)
}

func testAccComposerEnvironmentComposer3_databaseRetention(name, network, subnetwork string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  region = "us-central1"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
    node_config {
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
    }
	data_retention_config {
      airflow_metadata_retention_config {
        retention_mode = "RETENTION_MODE_ENABLED"
        retention_days = 61
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

// WARNING: This is not actually a check and is a terrible clean-up step because Composer Environments
// have a bug that hasn't been fixed. Composer will add firewalls to non-default networks for environments
// but will not remove them when the Environment is deleted.
//
// Destroy test step for config with a network will fail unless we clean up the firewalls before.
func testAccCheckClearComposerEnvironmentFirewalls(t *testing.T, networkName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		config.Project = envvar.GetTestProjectFromEnv()
		network, err := config.NewComputeClient(config.UserAgent).Networks.Get(envvar.GetTestProjectFromEnv(), networkName).Do()
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

			waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project,
				"Sweeping test composer environment firewalls", config.UserAgent, 10)
			if waitErr != nil {
				allErrors = multierror.Append(allErrors,
					fmt.Errorf("Error while waiting to delete firewall %q: %s", firewall.Name, waitErr))
			}
		}
		return allErrors
	}
}
