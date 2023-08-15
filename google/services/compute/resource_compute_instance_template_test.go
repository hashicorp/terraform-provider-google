// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"google.golang.org/api/compute/v1"
)

const DEFAULT_MIN_CPU_TEST_VALUE = "Intel Haswell"

func TestAccComputeInstanceTemplate_basic(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_basic(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateTag(&instanceTemplate, "foo"),
					testAccCheckComputeInstanceTemplateMetadata(&instanceTemplate, "foo", "bar"),
					testAccCheckComputeInstanceTemplateContainsLabel(&instanceTemplate, "my_label", "foobar"),
					testAccCheckComputeInstanceTemplateLacksShieldedVmConfig(&instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_imageShorthand(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_imageShorthand(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_preemptible(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_preemptible(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateAutomaticRestart(&instanceTemplate, false),
					testAccCheckComputeInstanceTemplatePreemptible(&instanceTemplate, true),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_IP(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_ip(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetwork(&instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_IPv6(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_ipv6(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_networkTier(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_networkTier(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_networkIP(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	networkIP := "10.128.0.2"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_networkIP(acctest.RandString(t, 10), networkIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetwork(&instanceTemplate),
					testAccCheckComputeInstanceTemplateNetworkIP(
						"google_compute_instance_template.foobar", networkIP, &instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_networkIPAddress(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	ipAddress := "10.128.0.2"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_networkIPAddress(acctest.RandString(t, 10), ipAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetwork(&instanceTemplate),
					testAccCheckComputeInstanceTemplateNetworkIPAddress(
						"google_compute_instance_template.foobar", ipAddress, &instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_disks(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_disks(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_disksInvalid(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeInstanceTemplate_disksInvalid(acctest.RandString(t, 10)),
				ExpectError: regexp.MustCompile("Cannot use `source`.*"),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_regionDisks(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_regionDisks(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_diskIops(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_diskIops(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_subnet_auto(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	network := "tf-test-network-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_subnet_auto(network, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetworkName(&instanceTemplate, network),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_subnet_custom(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_subnet_custom(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateSubnetwork(&instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_subnet_xpn(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	projectName := fmt.Sprintf("tf-testxpn-%d", time.Now().Unix())

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_subnet_xpn(org, billingId, projectName, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExistsInProject(
						t, "google_compute_instance_template.foobar", fmt.Sprintf("%s-service", projectName),
						&instanceTemplate),
					testAccCheckComputeInstanceTemplateSubnetwork(&instanceTemplate),
				),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_metadata_startup_script(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_startup_script(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateStartupScript(&instanceTemplate, "echo 'Hello'"),
				),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_primaryAliasIpRange(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_primaryAliasIpRange(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasAliasIpRange(&instanceTemplate, "", "/24"),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_secondaryAliasIpRange(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_secondaryAliasIpRange(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasAliasIpRange(&instanceTemplate, "inst-test-secondary", "/24"),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_guestAccelerator(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_guestAccelerator(acctest.RandString(t, 10), 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasGuestAccelerator(&instanceTemplate, "nvidia-tesla-k80", 1),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccComputeInstanceTemplate_guestAcceleratorSkip(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_guestAccelerator(acctest.RandString(t, 10), 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateLacksGuestAccelerator(&instanceTemplate),
				),
			},
		},
	})

}

func TestAccComputeInstanceTemplate_minCpuPlatform(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_minCpuPlatform(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasMinCpuPlatform(&instanceTemplate, DEFAULT_MIN_CPU_TEST_VALUE),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_EncryptKMS(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	kms := acctest.BootstrapKMSKey(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_encryptionKMS(acctest.RandString(t, 10), kms.CryptoKey.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_soleTenantNodeAffinities(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_soleTenantInstanceTemplate(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_instanceResourcePolicies(t *testing.T) {
	t.Parallel()

	var template compute.InstanceTemplate
	var policyName = "tf-test-policy-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_instanceResourcePolicyCollocated(acctest.RandString(t, 10), policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &template),
					testAccCheckComputeInstanceTemplateHasInstanceResourcePolicies(&template, policyName),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_reservationAffinities(t *testing.T) {
	t.Parallel()

	var template compute.InstanceTemplate
	var templateName = acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_reservationAffinityInstanceTemplate_nonSpecificReservation(templateName, "NO_RESERVATION"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &template),
					testAccCheckComputeInstanceTemplateHasReservationAffinity(&template, "NO_RESERVATION"),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeInstanceTemplate_reservationAffinityInstanceTemplate_nonSpecificReservation(templateName, "ANY_RESERVATION"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &template),
					testAccCheckComputeInstanceTemplateHasReservationAffinity(&template, "ANY_RESERVATION"),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeInstanceTemplate_reservationAffinityInstanceTemplate_specificReservation(templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &template),
					testAccCheckComputeInstanceTemplateHasReservationAffinity(&template, "SPECIFIC_RESERVATION", templateName),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_shieldedVmConfig1(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_shieldedVmConfig(acctest.RandString(t, 10), true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasShieldedVmConfig(&instanceTemplate, true, true, true),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_shieldedVmConfig2(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_shieldedVmConfig(acctest.RandString(t, 10), true, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasShieldedVmConfig(&instanceTemplate, true, true, false),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_ConfidentialInstanceConfigMain(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplateConfidentialInstanceConfig(acctest.RandString(t, 10), true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasConfidentialInstanceConfig(&instanceTemplate, true),
				),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_AdvancedMachineFeatures(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplateAdvancedMachineFeatures(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_invalidDiskType(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeInstanceTemplate_invalidDiskType(acctest.RandString(t, 10)),
				ExpectError: regexp.MustCompile("SCRATCH disks must have a disk_type of local-ssd"),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_withScratchDisk(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_with375GbScratchDisk(acctest.RandString(t, 10)),
			},
			{
				ResourceName:            "google_compute_instance_template.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func TestAccComputeInstanceTemplate_with18TbScratchDisk(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_with18TbScratchDisk(acctest.RandString(t, 10)),
			},
			{
				ResourceName:            "google_compute_instance_template.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func TestAccComputeInstanceTemplate_imageResourceTest(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()
	diskName := "tf-test-disk-" + acctest.RandString(t, 10)
	computeImage := "tf-test-image-" + acctest.RandString(t, 10)
	imageDesc1 := "Some description"
	imageDesc2 := "Some other description"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_imageResourceTest(diskName, computeImage, imageDesc1),
			},
			{
				ResourceName:            "google_compute_instance_template.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
			{
				Config: testAccComputeInstanceTemplate_imageResourceTest(diskName, computeImage, imageDesc2),
			},
			{
				ResourceName:            "google_compute_instance_template.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func TestAccComputeInstanceTemplate_diskResourcePolicies(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	policyName := "tf-test-policy-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_diskResourcePolicies(acctest.RandString(t, 10), policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasDiskResourcePolicy(&instanceTemplate, policyName),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_nictype_update(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	var instanceTemplateName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_nictype(instanceTemplateName, instanceTemplateName, "GVNIC"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
			{
				Config: testAccComputeInstanceTemplate_nictype(instanceTemplateName, instanceTemplateName, "VIRTIO_NET"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_queueCount(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	var instanceTemplateName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_queueCount(instanceTemplateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_managedEnvoy(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_managedEnvoy(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_spot(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_spot(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateAutomaticRestart(&instanceTemplate, false),
					testAccCheckComputeInstanceTemplatePreemptible(&instanceTemplate, true),
					testAccCheckComputeInstanceTemplateProvisioningModel(&instanceTemplate, "SPOT"),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_localSsdRecoveryTimeout(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	var expectedLocalSsdRecoveryTimeout = compute.Duration{}
	expectedLocalSsdRecoveryTimeout.Nanos = 0
	expectedLocalSsdRecoveryTimeout.Seconds = 3600

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_localSsdRecoveryTimeout(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateAutomaticRestart(&instanceTemplate, false),
					testAccCheckComputeInstanceTemplateLocalSsdRecoveryTimeout(&instanceTemplate, expectedLocalSsdRecoveryTimeout),
				),
			},
			{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_sourceSnapshotEncryptionKey(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	kmsKey := acctest.BootstrapKMSKeyInLocation(t, "us-central1")

	context := map[string]interface{}{
		"kms_ring_name": tpgresource.GetResourceNameFromSelfLink(kmsKey.KeyRing.Name),
		"kms_key_name":  tpgresource.GetResourceNameFromSelfLink(kmsKey.CryptoKey.Name),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_sourceSnapshotEncryptionKey(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.template", &instanceTemplate),
				),
			},
			{
				ResourceName:            "google_compute_instance_template.template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk.0.source_snapshot", "disk.0.source_snapshot_encryption_key"},
			},
		},
	})
}

func TestAccComputeInstanceTemplate_sourceImageEncryptionKey(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	kmsKey := acctest.BootstrapKMSKeyInLocation(t, "us-central1")

	context := map[string]interface{}{
		"kms_ring_name": tpgresource.GetResourceNameFromSelfLink(kmsKey.KeyRing.Name),
		"kms_key_name":  tpgresource.GetResourceNameFromSelfLink(kmsKey.CryptoKey.Name),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceTemplate_sourceImageEncryptionKey(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						t, "google_compute_instance_template.template", &instanceTemplate),
				),
			},
			{
				ResourceName:            "google_compute_instance_template.template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk.0.source_image_encryption_key"},
			},
		},
	})
}

func testAccCheckComputeInstanceTemplateDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_instance_template" {
				continue
			}

			splits := strings.Split(rs.Primary.ID, "/")
			_, err := config.NewComputeClient(config.UserAgent).InstanceTemplates.Get(
				config.Project, splits[len(splits)-1]).Do()
			if err == nil {
				return fmt.Errorf("Instance template still exists")
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateExists(t *testing.T, n string, instanceTemplate interface{}) resource.TestCheckFunc {
	if instanceTemplate == nil {
		panic("Attempted to check existence of Instance template that was nil.")
	}

	return testAccCheckComputeInstanceTemplateExistsInProject(t, n, envvar.GetTestProjectFromEnv(), instanceTemplate.(*compute.InstanceTemplate))
}

func testAccCheckComputeInstanceTemplateExistsInProject(t *testing.T, n, p string, instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		splits := strings.Split(rs.Primary.ID, "/")
		templateName := splits[len(splits)-1]
		found, err := config.NewComputeClient(config.UserAgent).InstanceTemplates.Get(
			p, templateName).Do()
		if err != nil {
			return err
		}

		if found.Name != templateName {
			return fmt.Errorf("Instance template not found")
		}
		if strings.Contains(rs.Primary.ID, "uniqueId") {
			return fmt.Errorf("unique ID is not supposed to be present in the Terraform resource ID")
		}
		selfLink := rs.Primary.Attributes["self_link"]
		if strings.Contains(selfLink, "uniqueId") {
			return fmt.Errorf("unique ID is not supposed to be present in selfLink")
		}

		actualSelfLinkUnique := rs.Primary.Attributes["self_link_unique"]
		foundId := strconv.FormatUint(found.Id, 10)
		expectedSelfLinkUnique := selfLink + "?uniqueId=" + foundId
		if actualSelfLinkUnique != expectedSelfLinkUnique {
			return fmt.Errorf("self_link_unique should be %v but it is: %v", expectedSelfLinkUnique, actualSelfLinkUnique)
		}

		*instanceTemplate = *found

		return nil
	}
}

func testAccCheckComputeInstanceTemplateMetadata(
	instanceTemplate *compute.InstanceTemplate,
	k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Metadata == nil {
			return fmt.Errorf("no metadata")
		}

		for _, item := range instanceTemplate.Properties.Metadata.Items {
			if k != item.Key {
				continue
			}

			if item.Value != nil && v == *item.Value {
				return nil
			}

			return fmt.Errorf("bad value for %s: %s", k, *item.Value)
		}

		return fmt.Errorf("metadata not found: %s", k)
	}
}

func testAccCheckComputeInstanceTemplateNetwork(instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instanceTemplate.Properties.NetworkInterfaces {
			for _, c := range i.AccessConfigs {
				if c.NatIP == "" {
					return fmt.Errorf("no NAT IP")
				}
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateNetworkName(instanceTemplate *compute.InstanceTemplate, network string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instanceTemplate.Properties.NetworkInterfaces {
			if !strings.Contains(i.Network, network) {
				return fmt.Errorf("Network doesn't match expected value, Expected: %s Actual: %s", network, i.Network[strings.LastIndex("/", i.Network)+1:])
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateSubnetwork(instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instanceTemplate.Properties.NetworkInterfaces {
			if i.Subnetwork == "" {
				return fmt.Errorf("no subnet")
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateTag(instanceTemplate *compute.InstanceTemplate, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Tags == nil {
			return fmt.Errorf("no tags")
		}

		for _, k := range instanceTemplate.Properties.Tags.Items {
			if k == n {
				return nil
			}
		}

		return fmt.Errorf("tag not found: %s", n)
	}
}

func testAccCheckComputeInstanceTemplatePreemptible(instanceTemplate *compute.InstanceTemplate, preemptible bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Scheduling.Preemptible != preemptible {
			return fmt.Errorf("Expected preemptible value %v, got %v", preemptible, instanceTemplate.Properties.Scheduling.Preemptible)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateProvisioningModel(instanceTemplate *compute.InstanceTemplate, provisioning_model string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Scheduling.ProvisioningModel != provisioning_model {
			return fmt.Errorf("Expected provisioning_model  %v, got %v", provisioning_model, instanceTemplate.Properties.Scheduling.ProvisioningModel)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateInstanceTerminationAction(instanceTemplate *compute.InstanceTemplate, instance_termination_action string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Scheduling.InstanceTerminationAction != instance_termination_action {
			return fmt.Errorf("Expected instance_termination_action  %v, got %v", instance_termination_action, instanceTemplate.Properties.Scheduling.InstanceTerminationAction)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateLocalSsdRecoveryTimeout(instanceTemplate *compute.InstanceTemplate, instance_local_ssd_recovery_timeout_want compute.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !reflect.DeepEqual(*instanceTemplate.Properties.Scheduling.LocalSsdRecoveryTimeout, instance_local_ssd_recovery_timeout_want) {
			return fmt.Errorf("gExpected LocalSsdRecoveryTimeout: %#v; got %#v", instance_local_ssd_recovery_timeout_want, instanceTemplate.Properties.Scheduling.LocalSsdRecoveryTimeout)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateAutomaticRestart(instanceTemplate *compute.InstanceTemplate, automaticRestart bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ar := instanceTemplate.Properties.Scheduling.AutomaticRestart
		if ar == nil {
			return fmt.Errorf("Expected to see a value for AutomaticRestart, but got nil")
		}
		if *ar != automaticRestart {
			return fmt.Errorf("Expected automatic restart value %v, got %v", automaticRestart, ar)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateStartupScript(instanceTemplate *compute.InstanceTemplate, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Metadata == nil && n == "" {
			return nil
		} else if instanceTemplate.Properties.Metadata == nil && n != "" {
			return fmt.Errorf("Expected metadata.startup-script to be '%s', metadata wasn't set at all", n)
		}
		for _, item := range instanceTemplate.Properties.Metadata.Items {
			if item.Key != "startup-script" {
				continue
			}
			if item.Value != nil && *item.Value == n {
				return nil
			} else if item.Value == nil && n == "" {
				return nil
			} else if item.Value == nil && n != "" {
				return fmt.Errorf("Expected metadata.startup-script to be '%s', wasn't set", n)
			} else if *item.Value != n {
				return fmt.Errorf("Expected metadata.startup-script to be '%s', got '%s'", n, *item.Value)
			}
		}
		return fmt.Errorf("This should never be reached.")
	}
}

func testAccCheckComputeInstanceTemplateNetworkIP(n, networkIP string, instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ip := instanceTemplate.Properties.NetworkInterfaces[0].NetworkIP
		err := resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", ip)(s)
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", networkIP)(s)
	}
}

func testAccCheckComputeInstanceTemplateNetworkIPAddress(n, ipAddress string, instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ip := instanceTemplate.Properties.NetworkInterfaces[0].NetworkIP
		err := resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", ip)(s)
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", ipAddress)(s)
	}
}

func testAccCheckComputeInstanceTemplateContainsLabel(instanceTemplate *compute.InstanceTemplate, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := instanceTemplate.Properties.Labels[key]
		if !ok {
			return fmt.Errorf("Expected label with key '%s' not found", key)
		}
		if v != value {
			return fmt.Errorf("Incorrect label value for key '%s': expected '%s' but found '%s'", key, value, v)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasAliasIpRange(instanceTemplate *compute.InstanceTemplate, subnetworkRangeName, iPCidrRange string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, networkInterface := range instanceTemplate.Properties.NetworkInterfaces {
			for _, aliasIpRange := range networkInterface.AliasIpRanges {
				if aliasIpRange.SubnetworkRangeName == subnetworkRangeName && (aliasIpRange.IpCidrRange == iPCidrRange || tpgresource.IpCidrRangeDiffSuppress("ip_cidr_range", aliasIpRange.IpCidrRange, iPCidrRange, nil)) {
					return nil
				}
			}
		}

		return fmt.Errorf("Alias ip range with name %s and cidr %s not present", subnetworkRangeName, iPCidrRange)
	}
}

func testAccCheckComputeInstanceTemplateHasGuestAccelerator(instanceTemplate *compute.InstanceTemplate, acceleratorType string, acceleratorCount int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(instanceTemplate.Properties.GuestAccelerators) != 1 {
			return fmt.Errorf("Expected only one guest accelerator")
		}

		if !strings.HasSuffix(instanceTemplate.Properties.GuestAccelerators[0].AcceleratorType, acceleratorType) {
			return fmt.Errorf("Wrong accelerator type: expected %v, got %v", acceleratorType, instanceTemplate.Properties.GuestAccelerators[0].AcceleratorType)
		}

		if instanceTemplate.Properties.GuestAccelerators[0].AcceleratorCount != acceleratorCount {
			return fmt.Errorf("Wrong accelerator acceleratorCount: expected %d, got %d", acceleratorCount, instanceTemplate.Properties.GuestAccelerators[0].AcceleratorCount)
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateLacksGuestAccelerator(instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(instanceTemplate.Properties.GuestAccelerators) > 0 {
			return fmt.Errorf("Expected no guest accelerators")
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasMinCpuPlatform(instanceTemplate *compute.InstanceTemplate, minCpuPlatform string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.MinCpuPlatform != minCpuPlatform {
			return fmt.Errorf("Wrong minimum CPU platform: expected %s, got %s", minCpuPlatform, instanceTemplate.Properties.MinCpuPlatform)
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasInstanceResourcePolicies(instanceTemplate *compute.InstanceTemplate, resourcePolicy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourcePolicyActual := instanceTemplate.Properties.ResourcePolicies[0]
		if resourcePolicyActual != resourcePolicy {
			return fmt.Errorf("Wrong instance resource policy: expected %s, got %s", resourcePolicy, resourcePolicyActual)
		}

		return nil
	}

}

func testAccCheckComputeInstanceTemplateHasReservationAffinity(instanceTemplate *compute.InstanceTemplate, consumeReservationType string, specificReservationNames ...string) resource.TestCheckFunc {
	if len(specificReservationNames) > 1 {
		panic("too many specificReservationNames in test")
	}

	return func(*terraform.State) error {
		if instanceTemplate.Properties.ReservationAffinity == nil {
			return fmt.Errorf("expected template to have reservation affinity, but it was nil")
		}

		if actualReservationType := instanceTemplate.Properties.ReservationAffinity.ConsumeReservationType; actualReservationType != consumeReservationType {
			return fmt.Errorf("Wrong reservationAffinity consumeReservationType: expected %s, got, %s", consumeReservationType, actualReservationType)
		}

		if len(specificReservationNames) > 0 {
			const reservationNameKey = "compute.googleapis.com/reservation-name"
			if actualKey := instanceTemplate.Properties.ReservationAffinity.Key; actualKey != reservationNameKey {
				return fmt.Errorf("Wrong reservationAffinity key: expected %s, got, %s", reservationNameKey, actualKey)
			}

			reservationAffinityValues := instanceTemplate.Properties.ReservationAffinity.Values
			if len(reservationAffinityValues) != 1 || reservationAffinityValues[0] != specificReservationNames[0] {
				return fmt.Errorf("Wrong reservationAffinity values: expected %s, got, %s", specificReservationNames, reservationAffinityValues)
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasShieldedVmConfig(instanceTemplate *compute.InstanceTemplate, enableSecureBoot bool, enableVtpm bool, enableIntegrityMonitoring bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		if instanceTemplate.Properties.ShieldedInstanceConfig.EnableSecureBoot != enableSecureBoot {
			return fmt.Errorf("Wrong shieldedVmConfig enableSecureBoot: expected %t, got, %t", enableSecureBoot, instanceTemplate.Properties.ShieldedInstanceConfig.EnableSecureBoot)
		}

		if instanceTemplate.Properties.ShieldedInstanceConfig.EnableVtpm != enableVtpm {
			return fmt.Errorf("Wrong shieldedVmConfig enableVtpm: expected %t, got, %t", enableVtpm, instanceTemplate.Properties.ShieldedInstanceConfig.EnableVtpm)
		}

		if instanceTemplate.Properties.ShieldedInstanceConfig.EnableIntegrityMonitoring != enableIntegrityMonitoring {
			return fmt.Errorf("Wrong shieldedVmConfig enableIntegrityMonitoring: expected %t, got, %t", enableIntegrityMonitoring, instanceTemplate.Properties.ShieldedInstanceConfig.EnableIntegrityMonitoring)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasConfidentialInstanceConfig(instanceTemplate *compute.InstanceTemplate, EnableConfidentialCompute bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		if instanceTemplate.Properties.ConfidentialInstanceConfig.EnableConfidentialCompute != EnableConfidentialCompute {
			return fmt.Errorf("Wrong ConfidentialInstanceConfig EnableConfidentialCompute: expected %t, got, %t", EnableConfidentialCompute, instanceTemplate.Properties.ConfidentialInstanceConfig.EnableConfidentialCompute)
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateLacksShieldedVmConfig(instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.ShieldedInstanceConfig != nil {
			return fmt.Errorf("Expected no shielded vm config")
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasDiskResourcePolicy(instanceTemplate *compute.InstanceTemplate, resourcePolicy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourcePolicyActual := instanceTemplate.Properties.Disks[0].InitializeParams.ResourcePolicies[0]
		if resourcePolicyActual != resourcePolicy {
			return fmt.Errorf("Wrong disk resource policy: expected %s, got %s", resourcePolicy, resourcePolicyActual)
		}

		return nil
	}
}

func testAccComputeInstanceTemplate_basic(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = false
    automatic_restart = true
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }

  labels = {
    my_label = "foobar"
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_imageShorthand(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_image" "foobar" {
  name        = "tf-test-%s"
  description = "description-test"
  family      = "family-test"
  raw_disk {
    source = "https://storage.googleapis.com/bosh-gce-raw-stemcells/bosh-stemcell-97.98-google-kvm-ubuntu-xenial-go_agent-raw-1557960142.tar.gz"
  }
  labels = {
    my-label    = "my-label-value"
    empty-label = ""
  }
  timeouts {
    create = "5m"
  }
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = google_compute_image.foobar.name
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = false
    automatic_restart = true
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }

  labels = {
    my_label = "foobar"
  }
}
`, suffix, suffix)
}

func testAccComputeInstanceTemplate_preemptible(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = true
    automatic_restart = false
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_ip(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "foo" {
  name = "tf-test-instance-template-%s"
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"
  tags         = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
  }

  network_interface {
    network = "default"
    access_config {
      nat_ip = google_compute_address.foo.address
    }
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, suffix)
}

func testAccComputeInstanceTemplate_ipv6(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "foo" {
  name = "tf-test-instance-template-%s"
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "foo" {
  name                    = "tf-test-network-%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork-ipv6" {
  name          = "tf-test-subnetwork-%s"

  ip_cidr_range = "10.0.0.0/22"
  region        = "us-west2"

  stack_type       = "IPV4_IPV6"
  ipv6_access_type = "EXTERNAL"

  network       = google_compute_network.foo.id
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"
  region       = "us-west2"
  tags         = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork-ipv6.name
    stack_type = "IPV4_IPV6"
    ipv6_access_config {
      network_tier = "PREMIUM"
    }
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, suffix, suffix, suffix)
}

func testAccComputeInstanceTemplate_networkTier(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
  }

  network_interface {
    network = "default"
    access_config {
      network_tier = "STANDARD"
    }
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_networkIP(suffix, networkIP string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"
  tags         = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
  }

  network_interface {
    network    = "default"
    network_ip = "%s"
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, networkIP)
}

func testAccComputeInstanceTemplate_networkIPAddress(suffix, ipAddress string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"
  tags         = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
  }

  network_interface {
    network    = "default"
    network_ip = "%s"
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, ipAddress)
}

func testAccComputeInstanceTemplate_disks(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "tf-test-instance-template-%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
    labels = {
      foo = "bar"
    }
  }

  disk {
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = false
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, suffix)
}

func testAccComputeInstanceTemplate_disksInvalid(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "tf-test-instance-template-%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
  }

  disk {
    source       = google_compute_disk.foobar.name
    disk_size_gb = 50
    auto_delete  = false
    boot         = false
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, suffix)
}

func testAccComputeInstanceTemplate_with375GbScratchDisk(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "centos-7"
	project = "centos-cloud"
}
resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  disk {
    source_image = data.google_compute_image.my_image.name
    auto_delete  = true
    boot         = true
  }
  disk {
    auto_delete  = true
    disk_size_gb = 375
    type         = "SCRATCH"
    disk_type    = "local-ssd"
  }
  network_interface {
    network = "default"
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_with18TbScratchDisk(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "centos-7"
	project = "centos-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "n2-standard-16"
  can_ip_forward = false
  disk {
    source_image = data.google_compute_image.my_image.name
    auto_delete  = true
    boot         = true
  }
  disk {
    auto_delete  = true
    disk_size_gb = 3000
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    interface    = "NVME"
  }
  disk {
    auto_delete  = true
    disk_size_gb = 3000
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    interface    = "NVME"
  }
  disk {
    auto_delete  = true
    disk_size_gb = 3000
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    interface    = "NVME"
  }
  disk {
    auto_delete  = true
    disk_size_gb = 3000
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    interface    = "NVME"
  }
  disk {
    auto_delete  = true
    disk_size_gb = 3000
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    interface    = "NVME"
  }
  disk {
    auto_delete  = true
    disk_size_gb = 3000
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    interface    = "NVME"
  }
  network_interface {
    network = "default"
  }
}`, suffix)
}

func testAccComputeInstanceTemplate_regionDisks(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_region_disk" "foobar" {
  name          = "tf-test-instance-template-%s"
  size          = 10
  type          = "pd-ssd"
  region        = "us-central1"
  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
  }

  disk {
    source      = google_compute_region_disk.foobar.name
    auto_delete = false
    boot        = false
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, suffix)
}

func testAccComputeInstanceTemplate_diskIops(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete      = true
    disk_size_gb     = 100
    boot             = true
    provisioned_iops = 10000
    labels = {
      foo = "bar"
    }
  }

  network_interface {
    network = "default"
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_subnet_auto(network, suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "auto-network" {
  name                    = "%s"
  auto_create_subnetworks = true
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  network_interface {
    network = google_compute_network.auto-network.name
  }

  metadata = {
    foo = "bar"
  }
}
`, network, suffix)
}

func testAccComputeInstanceTemplate_subnet_custom(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-network-%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "subnetwork-%s"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.network.self_link
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"
  region       = "us-central1"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork.name
  }

  metadata = {
    foo = "bar"
  }
}
`, suffix, suffix, suffix)
}

func testAccComputeInstanceTemplate_subnet_xpn(org, billingId, projectName, suffix string) string {
	return fmt.Sprintf(`
resource "google_project" "host_project" {
  name            = "Test Project XPN Host"
  project_id      = "%s-host"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "host_project" {
  project = google_project.host_project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_shared_vpc_host_project" "host_project" {
  project = google_project_service.host_project.project
}

resource "google_project" "service_project" {
  name            = "Test Project XPN Service"
  project_id      = "%s-service"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "service_project" {
  project = google_project.service_project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_shared_vpc_service_project" "service_project" {
  host_project    = google_compute_shared_vpc_host_project.host_project.project
  service_project = google_project_service.service_project.project
}

resource "google_compute_network" "network" {
  name                    = "tf-test-network-%s"
  auto_create_subnetworks = false
  project                 = google_compute_shared_vpc_host_project.host_project.project
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "subnetwork-%s"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.network.self_link
  project       = google_compute_shared_vpc_host_project.host_project.project
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"
  region       = "us-central1"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  network_interface {
    subnetwork         = google_compute_subnetwork.subnetwork.name
    subnetwork_project = google_compute_subnetwork.subnetwork.project
  }

  metadata = {
    foo = "bar"
  }
  project = google_compute_shared_vpc_service_project.service_project.service_project
}
`, projectName, org, billingId, projectName, org, billingId, suffix, suffix, suffix)
}

func testAccComputeInstanceTemplate_startup_script(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  metadata = {
    foo = "bar"
  }

  network_interface {
    network = "default"
  }

  metadata_startup_script = "echo 'Hello'"
}
`, suffix)
}

func testAccComputeInstanceTemplate_primaryAliasIpRange(i string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  metadata = {
    foo = "bar"
  }

  network_interface {
    network = "default"
    alias_ip_range {
      ip_cidr_range = "/24"
    }
  }
}
`, i)
}

func testAccComputeInstanceTemplate_secondaryAliasIpRange(i string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "inst-test-network" {
  name = "tf-test-network-%s"
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "inst-test-subnetwork-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.inst-test-network.self_link
  secondary_ip_range {
    range_name    = "inst-test-secondary"
    ip_cidr_range = "172.16.0.0/20"
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  metadata = {
    foo = "bar"
  }

  network_interface {
    subnetwork = google_compute_subnetwork.inst-test-subnetwork.self_link

    // Note that unlike compute instances, instance templates seem to be
    // only able to specify the netmask here. Trying a full CIDR string
    // results in:
    // Invalid value for field 'resource.properties.networkInterfaces[0].aliasIpRanges[0].ipCidrRange':
    // '172.16.0.0/24'. Alias IP CIDR range must be a valid netmask starting with '/' (e.g. '/24')
    alias_ip_range {
      subnetwork_range_name = google_compute_subnetwork.inst-test-subnetwork.secondary_ip_range[0].range_name
      ip_cidr_range         = "/24"
    }
  }
}
`, i, i, i)
}

func testAccComputeInstanceTemplate_guestAccelerator(i string, count uint8) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    # Instances with guest accelerators do not support live migration.
    on_host_maintenance = "TERMINATE"
  }

  guest_accelerator {
    count = %d
    type  = "nvidia-tesla-k80"
  }
}
`, i, count)
}

func testAccComputeInstanceTemplate_minCpuPlatform(i string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-medium"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    disk_size_gb = 10
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    # Instances with guest accelerators do not support live migration.
    on_host_maintenance = "TERMINATE"
  }

  min_cpu_platform = "%s"
}
`, i, DEFAULT_MIN_CPU_TEST_VALUE)
}

func testAccComputeInstanceTemplate_encryptionKMS(suffix, kmsLink string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false

  disk {
    source_image = data.google_compute_image.my_image.self_link
    disk_encryption_key {
      kms_key_self_link = "%s"
    }
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }

  labels = {
    my_label = "foobar"
  }
}
`, suffix, kmsLink)
}

func testAccComputeInstanceTemplate_soleTenantInstanceTemplate(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-standard-4"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = false
    automatic_restart = true
    node_affinities {
      key      = "tfacc"
      operator = "IN"
      values   = ["testinstancetemplate"]
    }

    min_node_cpus = 2
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_instanceResourcePolicyCollocated(suffix string, policyName string) string {
	return fmt.Sprintf(`
resource "google_compute_resource_policy" "foo" {
  name = "%s"
  region = "us-central1"
  group_placement_policy {
    vm_count  = 2
    collocation = "COLLOCATED"
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "e2-standard-4"

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = false
    automatic_restart = false
  }

  resource_policies = [google_compute_resource_policy.foo.self_link]

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}
`, policyName, suffix)
}

func testAccComputeInstanceTemplate_reservationAffinityInstanceTemplate_nonSpecificReservation(templateName, consumeReservationType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instancet-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  reservation_affinity {
    type = "%s"
  }
}
`, templateName, consumeReservationType)
}

func testAccComputeInstanceTemplate_reservationAffinityInstanceTemplate_specificReservation(templateName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instancet-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  reservation_affinity {
    type = "SPECIFIC_RESERVATION"

	specific_reservation {
		key = "compute.googleapis.com/reservation-name"
		values = ["%s"]
	}
  }
}
`, templateName, templateName)
}

func testAccComputeInstanceTemplate_shieldedVmConfig(suffix string, enableSecureBoot bool, enableVtpm bool, enableIntegrityMonitoring bool) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "centos-7"
  project = "centos-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  shielded_instance_config {
    enable_secure_boot          = %t
    enable_vtpm                 = %t
    enable_integrity_monitoring = %t
  }
}
`, suffix, enableSecureBoot, enableVtpm, enableIntegrityMonitoring)
}

func testAccComputeInstanceTemplateConfidentialInstanceConfig(suffix string, enableConfidentialCompute bool) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "ubuntu-2004-lts"
  project = "ubuntu-os-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "n2d-standard-2"

  disk {
    source_image = data.google_compute_image.my_image.self_link
	auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  confidential_instance_config {
    enable_confidential_compute       = %t
  }

  scheduling {
	  on_host_maintenance = "TERMINATE"
  }

}
`, suffix, enableConfidentialCompute)
}

func testAccComputeInstanceTemplateAdvancedMachineFeatures(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "ubuntu-2004-lts"
  project = "ubuntu-os-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name         = "tf-test-instance-template-%s"
  machine_type = "n2-standard-2" // Nested Virt isn't supported on E2 and N2Ds https://cloud.google.com/compute/docs/instances/nested-virtualization/overview#restrictions and https://cloud.google.com/compute/docs/instances/disabling-smt#limitations

  disk {
    source_image = data.google_compute_image.my_image.self_link
	auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  advanced_machine_features {
	threads_per_core = 1
	enable_nested_virtualization = true
	visible_core_count = 1
  }

  scheduling {
	  on_host_maintenance = "TERMINATE"
  }

}
`, suffix)
}

func testAccComputeInstanceTemplate_invalidDiskType(suffix string) string {
	return fmt.Sprintf(`
# Use this datasource insead of hardcoded values when https://github.com/hashicorp/terraform/issues/22679
# is resolved.
# data "google_compute_image" "my_image" {
# 	family  = "centos-7"
# 	project = "centos-cloud"
# }

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  disk {
    source_image = "https://www.googleapis.com/compute/v1/projects/centos-cloud/global/images/centos-7-v20210217"
    auto_delete  = true
    boot         = true
  }
  disk {
    auto_delete  = true
    disk_size_gb = 375
    type         = "SCRATCH"
    disk_type    = "local-ssd"
  }
  disk {
    source_image = "https://www.googleapis.com/compute/v1/projects/centos-cloud/global/images/centos-7-v20210217"
    auto_delete  = true
    type         = "SCRATCH"
  }
  network_interface {
    network = "default"
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_imageResourceTest(diskName string, imageName string, imageDescription string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-11"
	project = "debian-cloud"
}

resource "google_compute_disk" "my_disk" {
	name  = "%s"
	zone  = "us-central1-a"
	image = data.google_compute_image.my_image.self_link
}

resource "google_compute_image" "diskimage" {
	name = "%s"
	description = "%s"
	source_disk = google_compute_disk.my_disk.self_link
}

resource "google_compute_instance_template" "foobar" {
	name_prefix = "tf-test-instance-template-"
	machine_type         = "e2-medium"
	disk {
		source_image = google_compute_image.diskimage.self_link
	}
	network_interface {
		network = "default"
		access_config {}
	}
}
`, diskName, imageName, imageDescription)
}

func testAccComputeInstanceTemplate_diskResourcePolicies(suffix string, policyName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}
resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  disk {
    source_image = data.google_compute_image.my_image.self_link
    resource_policies = [google_compute_resource_policy.foo.id]
  }
  network_interface {
    network = "default"
  }
  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
  labels = {
    my_label = "foobar"
  }
}

resource "google_compute_resource_policy" "foo" {
  name   = "%s"
  region = "us-central1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time    = "04:00"
      }
    }
  }
}
`, suffix, policyName)
}

func testAccComputeInstanceTemplate_nictype(image, instance, nictype string) string {
	return fmt.Sprintf(`
resource "google_compute_image" "example" {
	name = "%s"
	raw_disk {
		source = "https://storage.googleapis.com/bosh-gce-raw-stemcells/bosh-stemcell-97.98-google-kvm-ubuntu-xenial-go_agent-raw-1557960142.tar.gz"
	}

	guest_os_features {
		type = "SECURE_BOOT"
	}

	guest_os_features {
		type = "MULTI_IP_SUBNET"
	}

	guest_os_features {
		type = "GVNIC"
	}
}

resource "google_compute_instance_template" "foobar" {
	name           = "tf-test-instance-template-%s"
	machine_type   = "e2-medium"
	can_ip_forward = false
	tags           = ["foo", "bar"]

	disk {
		source_image = google_compute_image.example.name
		auto_delete  = true
		boot         = true
	}

	network_interface {
		network = "default"
		nic_type = "%s"
	}

	scheduling {
		preemptible       = false
		automatic_restart = true
	}

	metadata = {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}

	labels = {
		my_label = "foobar"
	}
}
`, image, instance, nictype)
}

func testAccComputeInstanceTemplate_queueCount(instanceTemplateName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-11"
	project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
	name = "%s"
	machine_type         = "e2-medium"
	network_interface {
		network = "default"
		access_config {}
		queue_count = 2
	}
  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
}
`, instanceTemplateName)
}

func testAccComputeInstanceTemplate_managedEnvoy(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_default_service_account" "default" {
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = false
    automatic_restart = true
  }

  metadata = {
    gce-software-declaration = <<-EOF
    {
	  "softwareRecipes": [{
	    "name": "install-gce-service-proxy-agent",
	    "desired_state": "INSTALLED",
	    "installSteps": [{
	      "scriptRun": {
	        "script": "#! /bin/bash\nZONE=$(curl --silent http://metadata.google.internal/computeMetadata/v1/instance/zone -H Metadata-Flavor:Google | cut -d/ -f4 )\nexport SERVICE_PROXY_AGENT_DIRECTORY=$(mktemp -d)\nsudo gsutil cp   gs://gce-service-proxy-"$ZONE"/service-proxy-agent/releases/service-proxy-agent-0.2.tgz   "$SERVICE_PROXY_AGENT_DIRECTORY"   || sudo gsutil cp     gs://gce-service-proxy/service-proxy-agent/releases/service-proxy-agent-0.2.tgz     "$SERVICE_PROXY_AGENT_DIRECTORY"\nsudo tar -xzf "$SERVICE_PROXY_AGENT_DIRECTORY"/service-proxy-agent-0.2.tgz -C "$SERVICE_PROXY_AGENT_DIRECTORY"\n"$SERVICE_PROXY_AGENT_DIRECTORY"/service-proxy-agent/service-proxy-agent-bootstrap.sh"
	      }
	    }]
	  }]
	}
    EOF
    gce-service-proxy        = <<-EOF
    {
      "api-version": "0.2",
      "proxy-spec": {
        "proxy-port": 15001,
        "network": "my-network",
        "tracing": "ON",
        "access-log": "/var/log/envoy/access.log"
      }
      "service": {
        "serving-ports": [80, 81]
      },
     "labels": {
       "app_name": "bookserver_app",
       "app_version": "STABLE"
      }
    }
    EOF
    enable-guest-attributes = "true"
    enable-osconfig         = "true"

  }

  service_account {
  	email  = data.google_compute_default_service_account.default.email
    scopes = ["cloud-platform"]
  }

  labels = {
    gce-service-proxy = "on"
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_spot(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = true
    automatic_restart = false
    provisioning_model = "SPOT"
    instance_termination_action = "STOP"
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_spot_maxRunDuration(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    preemptible       = true
    automatic_restart = false
    provisioning_model = "SPOT"
    instance_termination_action = "DELETE"

  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_localSsdRecoveryTimeout(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "tf-test-instance-template-%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  scheduling {
    automatic_restart = false
    local_ssd_recovery_timeout {
		nanos = 0
		seconds = 3600
    }
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}
`, suffix)
}

func testAccComputeInstanceTemplate_sourceSnapshotEncryptionKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_kms_key_ring" "ring" {
  name     = "%{kms_ring_name}"
  location = "us-central1"
}

data "google_kms_crypto_key" "key" {
  name     = "%{kms_key_name}"
  key_ring = data.google_kms_key_ring.ring.id
}

resource "google_service_account" "test" {
  account_id   = "tf-test-sa-%{random_suffix}"
  display_name = "KMS Ops Account"
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = data.google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${google_service_account.test.email}"
}

data "google_compute_image" "debian" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "persistent" {
  name  = "tf-test-debian-disk-%{random_suffix}"
  image = data.google_compute_image.debian.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "snapshot" {
  name        = "tf-test-my-snapshot-%{random_suffix}"
  source_disk = google_compute_disk.persistent.id
  zone        = "us-central1-a"
  snapshot_encryption_key {
    kms_key_self_link       = data.google_kms_crypto_key.key.id
    kms_key_service_account = google_service_account.test.email
  }
}

resource "google_compute_instance_template" "template" {
  name           = "tf-test-instance-template-%{random_suffix}"
  machine_type   = "e2-medium"

  disk {
    source_snapshot = google_compute_snapshot.snapshot.self_link
    source_snapshot_encryption_key {
      kms_key_self_link       = data.google_kms_crypto_key.key.id
      kms_key_service_account = google_service_account.test.email
    }
    auto_delete = true
    boot        = true
  }

  network_interface {
    network = "default"
  }
}
`, context)
}

func testAccComputeInstanceTemplate_sourceImageEncryptionKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_kms_key_ring" "ring" {
  name     = "%{kms_ring_name}"
  location = "us-central1"
}

data "google_kms_crypto_key" "key" {
  name     = "%{kms_key_name}"
  key_ring = data.google_kms_key_ring.ring.id
}

resource "google_service_account" "test" {
  account_id   = "tf-test-sa-%{random_suffix}"
  display_name = "KMS Ops Account"
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = data.google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${google_service_account.test.email}"
}

data "google_compute_image" "debian" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_image" "image" {
  name         = "debian-image"
  source_image = data.google_compute_image.debian.self_link
  image_encryption_key {
    kms_key_self_link       = data.google_kms_crypto_key.key.id
    kms_key_service_account = google_service_account.test.email
  }
}


resource "google_compute_instance_template" "template" {
  name           = "tf-test-instance-template-%{random_suffix}"
  machine_type   = "e2-medium"

  disk {
    source_image = google_compute_image.image.self_link
    source_image_encryption_key {
      kms_key_self_link       = data.google_kms_crypto_key.key.id
      kms_key_service_account = google_service_account.test.email
    }
    auto_delete = true
    boot        = true
  }

  network_interface {
    network = "default"
  }
}
`, context)
}
