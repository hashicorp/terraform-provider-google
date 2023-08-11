// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"google.golang.org/api/compute/v1"
)

func TestMinCpuPlatformDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"state: empty, conf: AUTOMATIC": {
			Old:                "",
			New:                "AUTOMATIC",
			ExpectDiffSuppress: true,
		},
		"state: empty, conf: automatic": {
			Old:                "",
			New:                "automatic",
			ExpectDiffSuppress: true,
		},
		"state: empty, conf: AuToMaTiC": {
			Old:                "",
			New:                "AuToMaTiC",
			ExpectDiffSuppress: true,
		},
		"state: empty, conf: Intel Haswell": {
			Old:                "",
			New:                "Intel Haswell",
			ExpectDiffSuppress: false,
		},
		// This case should never happen due to the field being
		// Optional + Computed; however, including for completeness.
		"state: Intel Haswell, conf: empty": {
			Old:                "Intel Haswell",
			New:                "",
			ExpectDiffSuppress: false,
		},
		// These cases should never happen given current API behavior; testing
		// in case API behavior changes in the future.
		"state: AUTOMATIC, conf: Intel Haswell": {
			Old:                "AUTOMATIC",
			New:                "Intel Haswell",
			ExpectDiffSuppress: false,
		},
		"state: Intel Haswell, conf: AUTOMATIC": {
			Old:                "Intel Haswell",
			New:                "AUTOMATIC",
			ExpectDiffSuppress: false,
		},
		"state: AUTOMATIC, conf: empty": {
			Old:                "AUTOMATIC",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"state: automatic, conf: empty": {
			Old:                "automatic",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"state: AuToMaTiC, conf: empty": {
			Old:                "AuToMaTiC",
			New:                "",
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		if tpgcompute.ComputeInstanceMinCpuPlatformEmptyOrAutomaticDiffSuppress("min_cpu_platform", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func computeInstanceImportStep(zone, instanceName string, additionalImportIgnores []string) resource.TestStep {
	// metadata is only read into state if set in the config
	// importing doesn't know whether metadata.startup_script vs metadata_startup_script is set in the config,
	// it always takes metadata.startup-script
	ignores := []string{"metadata.%", "metadata.startup-script", "metadata_startup_script", "boot_disk.0.initialize_params.0.resource_manager_tags.%", "params.0.resource_manager_tags.%"}

	return resource.TestStep{
		ResourceName:            "google_compute_instance.foobar",
		ImportState:             true,
		ImportStateId:           fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), zone, instanceName),
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: append(ignores, additionalImportIgnores...),
	}
}

func TestAccComputeInstance_basic1(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasInstanceId(&instance, "google_compute_instance.foobar"),
					testAccCheckComputeInstanceTag(&instance, "foo"),
					testAccCheckComputeInstanceLabel(&instance, "my_key", "my_value"),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeInstanceMetadata(&instance, "baz", "qux"),
					testAccCheckComputeInstanceDisk(&instance, instanceName, true, true),
					resource.TestCheckResourceAttr("google_compute_instance.foobar", "current_status", "RUNNING"),

					// by default, DeletionProtection is implicitly false. This should be false on any
					// instance resource without an explicit deletion_protection = true declaration.
					// Other tests check explicit true/false configs: TestAccComputeInstance_deletionProtectionExplicit[True | False]
					testAccCheckComputeInstanceHasConfiguredDeletionProtection(&instance, false),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"metadata.baz", "metadata.foo", "desired_status", "current_status"}),
		},
	})
}

func TestAccComputeInstance_basic2(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceTag(&instance, "foo"),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeInstanceDisk(&instance, instanceName, true, true),
				),
			},
		},
	})
}

func TestAccComputeInstance_basic3(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic3(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceTag(&instance, "foo"),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeInstanceDisk(&instance, instanceName, true, true),
				),
			},
		},
	})
}

func TestAccComputeInstance_basic4(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic4(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceTag(&instance, "foo"),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeInstanceDisk(&instance, instanceName, true, true),
				),
			},
		},
	})
}

func TestAccComputeInstance_basic5(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic5(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceTag(&instance, "foo"),
					testAccCheckComputeInstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeInstanceDisk(&instance, instanceName, true, true),
				),
			},
		},
	})
}

func TestAccComputeInstance_resourceManagerTags(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
		"instance_name": instanceName,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_resourceManagerTags(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance)),
			},
		},
	})
}

func TestAccComputeInstance_IP(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var ipName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_ip(ipName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceAccessConfigHasNatIP(&instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_IPv6(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var ipName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var ptrName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_ipv6(ipName, instanceName, ptrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceIpv6AccessConfigHasExternalIPv6(&instance),
				),
			},
			{
				ResourceName:      "google_compute_instance.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstance_ipv6ExternalReservation(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_ipv6ExternalReservation(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-west2-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_PTRRecord(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var ptrName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var ipName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_PTRRecord(ptrName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceAccessConfigHasPTR(&instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"metadata.baz", "metadata.foo"}),
			{
				Config: testAccComputeInstance_ip(ipName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceAccessConfigHasNatIP(&instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"metadata.baz", "metadata.foo"}),
		},
	})
}

func TestAccComputeInstance_networkTier(t *testing.T) {
	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_networkTier(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceAccessConfigHasNatIP(&instance),
					testAccCheckComputeInstanceHasAssignedNatIP,
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_diskEncryption(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	bootEncryptionKey := "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	bootEncryptionKeyHash := "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E="
	diskNameToEncryptionKey := map[string]*compute.CustomerEncryptionKey{
		fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10)): {
			RawKey: "Ym9vdDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "awJ7p57H+uVZ9axhJjl1D3lfC2MgA/wnt/z88Ltfvss=",
		},
		fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10)): {
			RawKey: "c2Vjb25kNzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "7TpIwUdtCOJpq2m+3nt8GFgppu6a2Xsj1t0Gexk13Yc=",
		},
		fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10)): {
			RawKey: "dGhpcmQ2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "b3pvaS7BjDbCKeLPPTx7yXBuQtxyMobCHN1QJR43xeM=",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disks_encryption(bootEncryptionKey, diskNameToEncryptionKey, instanceName, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDiskEncryptionKey("google_compute_instance.foobar", &instance, bootEncryptionKeyHash, diskNameToEncryptionKey),
				),
			},
		},
	})
}

func TestAccComputeInstance_diskEncryptionRestart(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	bootEncryptionKey := "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	bootEncryptionKeyHash := "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E="
	diskNameToEncryptionKey := map[string]*compute.CustomerEncryptionKey{
		fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10)): {
			RawKey: "Ym9vdDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "awJ7p57H+uVZ9axhJjl1D3lfC2MgA/wnt/z88Ltfvss=",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disks_encryption_restart(bootEncryptionKey, diskNameToEncryptionKey, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDiskEncryptionKey("google_compute_instance.foobar", &instance, bootEncryptionKeyHash, diskNameToEncryptionKey),
				),
			},
			{
				Config: testAccComputeInstance_disks_encryption_restartUpdate(bootEncryptionKey, diskNameToEncryptionKey, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDiskEncryptionKey("google_compute_instance.foobar", &instance, bootEncryptionKeyHash, diskNameToEncryptionKey),
				),
			},
		},
	})
}

func TestAccComputeInstance_kmsDiskEncryption(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKey(t)

	bootKmsKeyName := kms.CryptoKey.Name
	diskNameToEncryptionKey := map[string]*compute.CustomerEncryptionKey{
		fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10)): {
			KmsKeyName: kms.CryptoKey.Name,
		},
		fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10)): {
			KmsKeyName: kms.CryptoKey.Name,
		},
		fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10)): {
			KmsKeyName: kms.CryptoKey.Name,
		},
	}

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disks_kms(bootKmsKeyName, diskNameToEncryptionKey, instanceName, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDiskKmsEncryptionKey("google_compute_instance.foobar", &instance, bootKmsKeyName, diskNameToEncryptionKey),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_resourcePolicyUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var scheduleName1 = fmt.Sprintf("tf-tests-%s", acctest.RandString(t, 10))
	var scheduleName2 = fmt.Sprintf("tf-tests-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_instanceSchedule(instanceName, scheduleName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeResourcePolicy(&instance, "", 0),
				),
			},
			// check adding
			{
				Config: testAccComputeInstance_addResourcePolicy(instanceName, scheduleName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeResourcePolicy(&instance, scheduleName1, 1),
				),
			},
			// check updating
			{
				Config: testAccComputeInstance_updateResourcePolicy(instanceName, scheduleName1, scheduleName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeResourcePolicy(&instance, scheduleName2, 1),
				),
			},
			// check removing
			{
				Config: testAccComputeInstance_removeResourcePolicy(instanceName, scheduleName1, scheduleName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeResourcePolicy(&instance, "", 0),
				),
			},
		},
	})
}

func TestAccComputeInstance_attachedDisk(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_attachedDisk(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDisk(&instance, diskName, false, false),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_attachedDisk_sourceUrl(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_attachedDisk_sourceUrl(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDisk(&instance, diskName, false, false),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_attachedDisk_modeRo(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_attachedDisk_modeRo(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDisk(&instance, diskName, false, false),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_attachedDiskUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10))
	var diskName2 = fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_attachedDisk(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDisk(&instance, diskName, false, false),
				),
			},
			// check attaching
			{
				Config: testAccComputeInstance_addAttachedDisk(diskName, diskName2, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDisk(&instance, diskName, false, false),
					testAccCheckComputeInstanceDisk(&instance, diskName2, false, false),
				),
			},
			// check detaching
			{
				Config: testAccComputeInstance_detachDisk(diskName, diskName2, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDisk(&instance, diskName, false, false),
				),
			},
			// check updating
			{
				Config: testAccComputeInstance_updateAttachedDiskEncryptionKey(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDisk(&instance, diskName, false, false),
				),
			},
		},
	})
}

func TestAccComputeInstance_bootDisk_source(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_bootDisk_source(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceBootDisk(&instance, diskName),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_bootDisk_sourceUrl(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_bootDisk_sourceUrl(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceBootDisk(&instance, diskName),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_bootDisk_type(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskType = "pd-ssd"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_bootDisk_type(instanceName, diskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceBootDiskType(t, instanceName, diskType),
				),
			},
		},
	})
}

func TestAccComputeInstance_bootDisk_mode(t *testing.T) {
	t.Parallel()

	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskMode = "READ_WRITE"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_bootDisk_mode(instanceName, diskMode),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_with375GbScratchDisk(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_with375GbScratchDisk(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScratchDisk(&instance, []string{"NVME", "SCSI"}),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_with18TbScratchDisk(t *testing.T) {
	// Skip this test until the quota for the GitHub presubmit GCP project is increased
	// to handle the size of the resource this test spins up.
	t.Skip()
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_with18TbScratchDisk(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScratchDisk(&instance, []string{"NVME", "NVME", "NVME", "NVME", "NVME", "NVME"}),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_forceNewAndChangeMetadata(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			{
				Config: testAccComputeInstance_forceNewAndChangeMetadata(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceMetadata(
						&instance, "qux", "true"),
				),
			},
		},
	})
}

func TestAccComputeInstance_update(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			{
				Config: testAccComputeInstance_update(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceMetadata(
						&instance, "bar", "baz"),
					testAccCheckComputeInstanceLabel(&instance, "only_me", "nothing_else"),
					testAccCheckComputeInstanceTag(&instance, "baz"),
					testAccCheckComputeInstanceAccessConfig(&instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_stopInstanceToUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			// Set fields that require stopping the instance
			{
				Config: testAccComputeInstance_stopInstanceToUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			// Check that updating them works
			{
				Config: testAccComputeInstance_stopInstanceToUpdate2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			// Check that removing them works
			{
				Config: testAccComputeInstance_stopInstanceToUpdate3(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_serviceAccount(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_serviceAccount(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceServiceAccount(&instance,
						"https://www.googleapis.com/auth/compute.readonly"),
					testAccCheckComputeInstanceServiceAccount(&instance,
						"https://www.googleapis.com/auth/devstorage.read_only"),
					testAccCheckComputeInstanceServiceAccount(&instance,
						"https://www.googleapis.com/auth/userinfo.email"),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_serviceAccount_updated(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_serviceAccount_update0(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScopes(&instance, 0),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_serviceAccount_update01(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScopes(&instance, 0),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_serviceAccount_update02(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScopes(&instance, 0),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_serviceAccount_update3(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScopes(&instance, 3),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_serviceAccount_updated0to1to0scopes(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_serviceAccount_update01(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScopes(&instance, 0),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_serviceAccount_update4(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScopes(&instance, 1),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_serviceAccount_update01(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceScopes(&instance, 0),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_scheduling(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_scheduling(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
			{
				Config: testAccComputeInstance_schedulingUpdated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_advancedMachineFeatures(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_advancedMachineFeatures(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_advancedMachineFeaturesUpdated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_soleTenantNodeAffinities(t *testing.T) {
	t.Parallel()

	var instanceName = fmt.Sprintf("tf-test-soletenant-%s", acctest.RandString(t, 10))
	var templateName = fmt.Sprintf("tf-test-nodetmpl-%s", acctest.RandString(t, 10))
	var groupName = fmt.Sprintf("tf-test-nodegroup-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_withoutNodeAffinities(instanceName, templateName, groupName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_soleTenantNodeAffinities(instanceName, templateName, groupName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_soleTenantNodeAffinitiesUpdated(instanceName, templateName, groupName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_soleTenantNodeAffinitiesReduced(instanceName, templateName, groupName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_reservationAffinities(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-resaffinity-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_reservationAffinity_nonSpecificReservationConfig(instanceName, "NO_RESERVATION"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasReservationAffinity(&instance, "NO_RESERVATION"),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
			{
				Config: testAccComputeInstance_reservationAffinity_nonSpecificReservationConfig(instanceName, "ANY_RESERVATION"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasReservationAffinity(&instance, "ANY_RESERVATION"),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
			{
				Config: testAccComputeInstance_reservationAffinity_specificReservationConfig(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasReservationAffinity(&instance, "SPECIFIC_RESERVATION", instanceName),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_subnet_auto(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_subnet_auto(acctest.RandString(t, 10), instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasSubnet(&instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_subnet_custom(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_subnet_custom(acctest.RandString(t, 10), instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasSubnet(&instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_subnet_xpn(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	projectName := fmt.Sprintf("tf-test-xpn-%d", time.Now().Unix())

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_subnet_xpn(org, billingId, projectName, instanceName, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExistsInProject(
						t, "google_compute_instance.foobar", fmt.Sprintf("%s-service", projectName),
						&instance),
					testAccCheckComputeInstanceHasSubnet(&instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_networkIPAuto(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_networkIPAuto(acctest.RandString(t, 10), instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasAnyNetworkIP(&instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_network_ip_custom(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var ipAddress = "10.0.200.200"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_network_ip_custom(acctest.RandString(t, 10), instanceName, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasNetworkIP(&instance, ipAddress),
				),
			},
		},
	})
}

func TestAccComputeInstance_private_image_family(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10))
	var familyName = fmt.Sprintf("tf-testf-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_private_image_family(diskName, familyName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_networkPerformanceConfig(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", acctest.RandString(t, 10))
	var imageName = fmt.Sprintf("tf-testf-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_networkPerformanceConfig(imageName, diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasNetworkPerformanceConfig(&instance, "DEFAULT"),
				),
			},
		},
	})
}

func TestAccComputeInstance_forceChangeMachineTypeManually(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceUpdateMachineType(t, "google_compute_instance.foobar"),
				),
				ExpectNonEmptyPlan: true,
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"metadata.baz", "metadata.foo", "desired_status", "current_status"}),
		},
	})
}

func TestAccComputeInstance_multiNic(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	networkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetworkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_multiNic(instanceName, networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMultiNic(&instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_nictype_update(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_nictype(instanceName, instanceName, "GVNIC"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			{
				Config: testAccComputeInstance_nictype(instanceName, instanceName, "VIRTIO_NET"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
		},
	})
}

func TestAccComputeInstance_guestAccelerator(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_guestAccelerator(instanceName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasGuestAccelerator(&instance, "nvidia-tesla-k80", 1),
				),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{"metadata.baz", "metadata.foo"}),
		},
	})

}

func TestAccComputeInstance_guestAcceleratorSkip(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_guestAccelerator(instanceName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceLacksGuestAccelerator(&instance),
				),
			},
		},
	})

}

func TestAccComputeInstance_minCpuPlatform(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_minCpuPlatform(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMinCpuPlatform(&instance, "Intel Haswell"),
				),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_minCpuPlatform_remove(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMinCpuPlatform(&instance, ""),
				),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_deletionProtectionExplicitFalse(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic_deletionProtectionFalse(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasConfiguredDeletionProtection(&instance, false),
				),
			},
		},
	})
}

func TestAccComputeInstance_deletionProtectionExplicitTrueAndUpdateFalse(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic_deletionProtectionTrue(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasConfiguredDeletionProtection(&instance, true),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"metadata.foo"}),
			// Update deletion_protection to false, otherwise the test harness can't delete the instance
			{
				Config: testAccComputeInstance_basic_deletionProtectionFalse(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasConfiguredDeletionProtection(&instance, false),
				),
			},
		},
	})
}

func TestAccComputeInstance_primaryAliasIpRange(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_primaryAliasIpRange(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasAliasIpRange(&instance, "", "/24"),
				),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_secondaryAliasIpRange(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	networkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_secondaryAliasIpRange(networkName, subnetName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasAliasIpRange(&instance, "inst-test-secondary", "172.16.0.0/24"),
				),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{}),
			{
				Config: testAccComputeInstance_secondaryAliasIpRangeUpdate(networkName, subnetName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasAliasIpRange(&instance, "", "10.0.1.0/24"),
				),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_hostname(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_hostname(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_compute_instance.foobar", "hostname"),
					testAccCheckComputeInstanceLacksShieldedVmConfig(&instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_shieldedVmConfig(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_shieldedVmConfig(instanceName, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasShieldedVmConfig(&instance, true, true, true),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_shieldedVmConfig(instanceName, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasShieldedVmConfig(&instance, true, true, false),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstanceConfidentialInstanceConfigMain(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceConfidentialInstanceConfig(instanceName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasConfidentialInstanceConfig(&instance, true),
				),
			},
		},
	})
}

func TestAccComputeInstance_enableDisplay(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_enableDisplay(instanceName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_enableDisplayUpdated(instanceName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_desiredStatusOnCreation(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				ExpectError: regexp.MustCompile("When creating an instance, desired_status can only accept RUNNING value"),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "RUNNING", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
		},
	})
}

func TestAccComputeInstance_desiredStatusUpdateBasic(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "RUNNING", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "RUNNING", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
		},
	})
}

func TestAccComputeInstance_desiredStatusTerminatedUpdateFields(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_desiredStatusTerminatedUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceMetadata(
						&instance, "bar", "baz"),
					testAccCheckComputeInstanceLabel(&instance, "only_me", "nothing_else"),
					testAccCheckComputeInstanceTag(&instance, "baz"),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateRunning_desiredStatusRunning_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "RUNNING", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
		},
	})
}

const errorAllowStoppingMsg = "please set allow_stopping_for_update"

func TestAccComputeInstance_updateRunning_desiredStatusNotSet_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config:      testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "", false),
				ExpectError: regexp.MustCompile(errorAllowStoppingMsg),
			},
		},
	})
}

func TestAccComputeInstance_updateRunning_desiredStatusRunning_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config:      testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "RUNNING", false),
				ExpectError: regexp.MustCompile(errorAllowStoppingMsg),
			},
		},
	})
}

func TestAccComputeInstance_updateRunning_desiredStatusTerminated_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "TERMINATED", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateRunning_desiredStatusTerminated_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateTerminated_desiredStatusNotSet_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateTerminated_desiredStatusTerminated_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "TERMINATED", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateTerminated_desiredStatusNotSet_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateTerminated_desiredStatusTerminated_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateTerminated_desiredStatusRunning_allowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "RUNNING", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
		},
	})
}

func TestAccComputeInstance_updateTerminated_desiredStatusRunning_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic2(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-medium", "TERMINATED", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasStatus(&instance, "TERMINATED"),
				),
			},
			{
				Config: testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(instanceName, "e2-standard-2", "RUNNING", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMachineType(&instance, "e2-standard-2"),
					testAccCheckComputeInstanceHasStatus(&instance, "RUNNING"),
				),
			},
		},
	})
}

func TestAccComputeInstance_resourcePolicyCollocate(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_resourcePolicyCollocate(instanceName, acctest.RandString(t, 10)),
			},
			computeInstanceImportStep("us-east4-b", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_subnetworkUpdate(t *testing.T) {
	t.Parallel()
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	suffix := fmt.Sprintf("%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_subnetworkUpdate(suffix, instanceName),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_subnetworkUpdateTwo(suffix, instanceName),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{"allow_stopping_for_update"}),
			{
				Config: testAccComputeInstance_subnetworkUpdate(suffix, instanceName),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_queueCount(t *testing.T) {
	t.Parallel()
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_queueCountSet(instanceName),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_spotVM(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_spotVM(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_spotVM_update(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_scheduling(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
			{
				Config: testAccComputeInstance_spotVM(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_localSsdRecoveryTimeout(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var expectedLocalSsdRecoveryTimeout = compute.Duration{}
	// Define in testAccComputeInstance_localSsdRecoveryTimeout
	expectedLocalSsdRecoveryTimeout.Nanos = 0
	expectedLocalSsdRecoveryTimeout.Seconds = 3600

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_localSsdRecoveryTimeout(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceLocalSsdRecoveryTimeout(&instance, expectedLocalSsdRecoveryTimeout),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_localSsdRecoveryTimeout_update(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	// Define in testAccComputeInstance_localSsdRecoveryTimeout
	var expectedLocalSsdRecoveryTimeout = compute.Duration{}
	expectedLocalSsdRecoveryTimeout.Nanos = 0
	expectedLocalSsdRecoveryTimeout.Seconds = 3600
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_scheduling(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
			{
				Config: testAccComputeInstance_localSsdRecoveryTimeout(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceLocalSsdRecoveryTimeout(&instance, expectedLocalSsdRecoveryTimeout),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_metadataStartupScript_update(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_metadataStartupScript(instanceName, "e2-medium", "abc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
			{
				Config: testAccComputeInstance_metadataStartupScript(instanceName, "e2-standard-4", "xyz"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(
						t, "google_compute_instance.foobar", &instance),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceUpdateMachineType(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		op, err := config.NewComputeClient(config.UserAgent).Instances.Stop(config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return fmt.Errorf("Could not stop instance: %s", err)
		}
		err = tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "Waiting on stop", config.UserAgent, 20*time.Minute)
		if err != nil {
			return fmt.Errorf("Could not stop instance: %s", err)
		}

		machineType := compute.InstancesSetMachineTypeRequest{
			MachineType: "zones/us-central1-a/machineTypes/f1-micro",
		}

		op, err = config.NewComputeClient(config.UserAgent).Instances.SetMachineType(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"], &machineType).Do()
		if err != nil {
			return fmt.Errorf("Could not change machine type: %s", err)
		}
		err = tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "Waiting machine type change", config.UserAgent, 20*time.Minute)
		if err != nil {
			return fmt.Errorf("Could not change machine type: %s", err)
		}
		return nil
	}
}

func testAccCheckComputeInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_instance" {
				continue
			}

			_, err := config.NewComputeClient(config.UserAgent).Instances.Get(
				config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("Instance still exists")
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceExists(t *testing.T, n string, instance interface{}) resource.TestCheckFunc {
	if instance == nil {
		panic("Attempted to check existence of Instance that was nil.")
	}

	return testAccCheckComputeInstanceExistsInProject(t, n, envvar.GetTestProjectFromEnv(), instance.(*compute.Instance))
}

func testAccCheckComputeInstanceExistsInProject(t *testing.T, n, p string, instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewComputeClient(config.UserAgent).Instances.Get(
			p, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckComputeInstanceMetadata(
	instance *compute.Instance,
	k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return fmt.Errorf("no metadata")
		}

		for _, item := range instance.Metadata.Items {
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

func testAccCheckComputeInstanceAccessConfig(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instance.NetworkInterfaces {
			if len(i.AccessConfigs) == 0 {
				return fmt.Errorf("no access_config")
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceAccessConfigHasNatIP(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instance.NetworkInterfaces {
			for _, c := range i.AccessConfigs {
				if c.NatIP == "" {
					return fmt.Errorf("no NAT IP")
				}
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceIpv6AccessConfigHasExternalIPv6(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instance.NetworkInterfaces {
			for _, c := range i.Ipv6AccessConfigs {
				if c.ExternalIpv6 == "" {
					return fmt.Errorf("no External IPv6")
				}
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceAccessConfigHasPTR(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instance.NetworkInterfaces {
			for _, c := range i.AccessConfigs {
				if c.PublicPtrDomainName == "" {
					return fmt.Errorf("no PTR Record")
				}
			}
		}

		return nil
	}
}

func testAccCheckComputeResourcePolicy(instance *compute.Instance, scheduleName string, resourcePolicyCountWant int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourcePoliciesCountHave := len(instance.ResourcePolicies)
		if resourcePoliciesCountHave != resourcePolicyCountWant {
			return fmt.Errorf("number of resource polices does not match: have: %d; want: %d", resourcePoliciesCountHave, resourcePolicyCountWant)
		}

		if resourcePoliciesCountHave == 1 && !strings.Contains(instance.ResourcePolicies[0], scheduleName) {
			return fmt.Errorf("got the wrong schedule: have: %s; want: %s", instance.ResourcePolicies[0], scheduleName)
		}

		return nil
	}
}

func testAccCheckComputeInstanceLocalSsdRecoveryTimeout(instance *compute.Instance, instanceLocalSsdRecoveryTiemoutWant compute.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance == nil {
			return fmt.Errorf("instance is nil")
		}
		if instance.Scheduling == nil {
			return fmt.Errorf("no scheduling")
		}

		if !reflect.DeepEqual(*instance.Scheduling.LocalSsdRecoveryTimeout, instanceLocalSsdRecoveryTiemoutWant) {
			return fmt.Errorf("got the wrong instance local ssd recovery timeout action: have: %#v; want: %#v", instance.Scheduling.LocalSsdRecoveryTimeout, instanceLocalSsdRecoveryTiemoutWant)
		}

		return nil
	}
}

func testAccCheckComputeInstanceTerminationAction(instance *compute.Instance, instanceTerminationActionWant string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance == nil {
			return fmt.Errorf("instance is nil")
		}
		if instance.Scheduling == nil {
			return fmt.Errorf("no scheduling")
		}

		if instance.Scheduling.InstanceTerminationAction != instanceTerminationActionWant {
			return fmt.Errorf("got the wrong instance termniation action: have: %s; want: %s", instance.Scheduling.InstanceTerminationAction, instanceTerminationActionWant)
		}

		return nil
	}
}

func testAccCheckComputeInstanceDisk(instance *compute.Instance, source string, delete bool, boot bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Disks == nil {
			return fmt.Errorf("no disks")
		}

		for _, disk := range instance.Disks {
			if strings.HasSuffix(disk.Source, "/"+source) && disk.AutoDelete == delete && disk.Boot == boot {
				return nil
			}
		}

		return fmt.Errorf("Disk not found: %s", source)
	}
}

func testAccCheckComputeInstanceHasInstanceId(instance *compute.Instance, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		remote := fmt.Sprintf("%d", instance.Id)
		local := rs.Primary.Attributes["instance_id"]

		if remote != local {
			return fmt.Errorf("Instance id stored does not match: remote has %#v but local has %#v", remote,
				local)
		}

		return nil
	}
}

func testAccCheckComputeInstanceBootDisk(instance *compute.Instance, source string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Disks == nil {
			return fmt.Errorf("no disks")
		}

		for _, disk := range instance.Disks {
			if disk.Boot == true {
				if strings.HasSuffix(disk.Source, source) {
					return nil
				}
			}
		}

		return fmt.Errorf("Boot disk not found with source %q", source)
	}
}

func testAccCheckComputeInstanceBootDiskType(t *testing.T, instanceName string, diskType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		// boot disk is named the same as the Instance
		disk, err := config.NewComputeClient(config.UserAgent).Disks.Get(config.Project, "us-central1-a", instanceName).Do()
		if err != nil {
			return err
		}
		if strings.Contains(disk.Type, diskType) {
			return nil
		}

		return fmt.Errorf("Boot disk not found with type %q", diskType)
	}
}

func testAccCheckComputeInstanceScratchDisk(instance *compute.Instance, interfaces []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Disks == nil {
			return fmt.Errorf("no disks")
		}

		i := 0
		for _, disk := range instance.Disks {
			if disk.Type == "SCRATCH" {
				if i >= len(interfaces) {
					return fmt.Errorf("Expected %d scratch disks, found more", len(interfaces))
				}
				if disk.Interface != interfaces[i] {
					return fmt.Errorf("Mismatched interface on scratch disk #%d, expected: %q, found: %q",
						i, interfaces[i], disk.Interface)
				}
				i++
			}
		}

		if i != len(interfaces) {
			return fmt.Errorf("Expected %d scratch disks, found %d", len(interfaces), i)
		}

		return nil
	}
}

func testAccCheckComputeInstanceDiskEncryptionKey(n string, instance *compute.Instance, bootDiskEncryptionKey string, diskNameToEncryptionKey map[string]*compute.CustomerEncryptionKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		for i, disk := range instance.Disks {
			if disk.Boot {
				attr := rs.Primary.Attributes["boot_disk.0.disk_encryption_key_sha256"]
				if attr != bootDiskEncryptionKey {
					return fmt.Errorf("Boot disk has wrong encryption key in state.\nExpected: %s\nActual: %s", bootDiskEncryptionKey, attr)
				}
				if disk.DiskEncryptionKey == nil && attr != "" {
					return fmt.Errorf("Disk %d has mismatched encryption key.\nTF State: %+v\nGCP State: <empty>", i, attr)
				}
				if disk.DiskEncryptionKey != nil && attr != disk.DiskEncryptionKey.Sha256 {
					return fmt.Errorf("Disk %d has mismatched encryption key.\nTF State: %+v\nGCP State: %+v",
						i, attr, disk.DiskEncryptionKey.Sha256)
				}
			} else {
				if disk.DiskEncryptionKey != nil {
					expectedKey := diskNameToEncryptionKey[tpgresource.GetResourceNameFromSelfLink(disk.Source)].Sha256
					if disk.DiskEncryptionKey.Sha256 != expectedKey {
						return fmt.Errorf("Disk %d has unexpected encryption key in GCP.\nExpected: %s\nActual: %s", i, expectedKey, disk.DiskEncryptionKey.Sha256)
					}
				}
			}
		}

		numAttachedDisks, err := strconv.Atoi(rs.Primary.Attributes["attached_disk.#"])
		if err != nil {
			return fmt.Errorf("Error converting value of attached_disk.#")
		}
		for i := 0; i < numAttachedDisks; i++ {
			diskName := tpgresource.GetResourceNameFromSelfLink(rs.Primary.Attributes[fmt.Sprintf("attached_disk.%d.source", i)])
			encryptionKey := rs.Primary.Attributes[fmt.Sprintf("attached_disk.%d.disk_encryption_key_sha256", i)]
			if key, ok := diskNameToEncryptionKey[diskName]; ok {
				expectedEncryptionKey := key.Sha256
				if encryptionKey != expectedEncryptionKey {
					return fmt.Errorf("Attached disk %d has unexpected encryption key in state.\nExpected: %s\nActual: %s", i, expectedEncryptionKey, encryptionKey)
				}
			}
		}
		return nil
	}
}

func testAccCheckComputeInstanceDiskKmsEncryptionKey(n string, instance *compute.Instance, bootDiskEncryptionKey string, diskNameToEncryptionKey map[string]*compute.CustomerEncryptionKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		for i, disk := range instance.Disks {
			if disk.Boot {
				attr := rs.Primary.Attributes["boot_disk.0.kms_key_self_link"]
				if attr != bootDiskEncryptionKey {
					return fmt.Errorf("Boot disk has wrong encryption key in state.\nExpected: %s\nActual: %s", bootDiskEncryptionKey, attr)
				}
				if disk.DiskEncryptionKey == nil && attr != "" {
					return fmt.Errorf("Disk %d has mismatched encryption key.\nTF State: %+v\nGCP State: <empty>", i, attr)
				}
			} else {
				if disk.DiskEncryptionKey != nil {
					expectedKey := diskNameToEncryptionKey[tpgresource.GetResourceNameFromSelfLink(disk.Source)].KmsKeyName
					// The response for crypto keys often includes the version of the key which needs to be removed
					// format: projects/<project>/locations/<region>/keyRings/<keyring>/cryptoKeys/<key>/cryptoKeyVersions/1
					actualKey := strings.Split(disk.DiskEncryptionKey.KmsKeyName, "/cryptoKeyVersions")[0]
					if actualKey != expectedKey {
						return fmt.Errorf("Disk %d has unexpected encryption key in GCP.\nExpected: %s\nActual: %s", i, expectedKey, actualKey)
					}
				}
			}
		}

		numAttachedDisks, err := strconv.Atoi(rs.Primary.Attributes["attached_disk.#"])
		if err != nil {
			return fmt.Errorf("Error converting value of attached_disk.#")
		}
		for i := 0; i < numAttachedDisks; i++ {
			diskName := tpgresource.GetResourceNameFromSelfLink(rs.Primary.Attributes[fmt.Sprintf("attached_disk.%d.source", i)])
			kmsKeyName := rs.Primary.Attributes[fmt.Sprintf("attached_disk.%d.kms_key_self_link", i)]
			if key, ok := diskNameToEncryptionKey[diskName]; ok {
				expectedEncryptionKey := key.KmsKeyName
				if kmsKeyName != expectedEncryptionKey {
					return fmt.Errorf("Attached disk %d has unexpected encryption key in state.\nExpected: %s\nActual: %s", i, expectedEncryptionKey, kmsKeyName)
				}
			}
		}
		return nil
	}
}

func testAccCheckComputeInstanceTag(instance *compute.Instance, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Tags == nil {
			return fmt.Errorf("no tags")
		}

		for _, k := range instance.Tags.Items {
			if k == n {
				return nil
			}
		}

		return fmt.Errorf("tag not found: %s", n)
	}
}

func testAccCheckComputeInstanceLabel(instance *compute.Instance, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Labels == nil {
			return fmt.Errorf("no labels found on instance %s", instance.Name)
		}

		v, ok := instance.Labels[key]
		if !ok {
			return fmt.Errorf("No label found with key %s on instance %s", key, instance.Name)
		}
		if v != value {
			return fmt.Errorf("Expected value '%s' but found value '%s' for label '%s' on instance %s", value, v, key, instance.Name)
		}

		return nil
	}
}

func testAccCheckComputeInstanceServiceAccount(instance *compute.Instance, scope string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if count := len(instance.ServiceAccounts); count != 1 {
			return fmt.Errorf("Wrong number of ServiceAccounts: expected 1, got %d", count)
		}

		for _, val := range instance.ServiceAccounts[0].Scopes {
			if val == scope {
				return nil
			}
		}

		return fmt.Errorf("Scope not found: %s", scope)
	}
}

func testAccCheckComputeInstanceScopes(instance *compute.Instance, scopeCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if count := len(instance.ServiceAccounts); count == 0 {
			if scopeCount == 0 {
				return nil
			} else {
				return fmt.Errorf("Scope count expected: %s, but got %s", fmt.Sprint(scopeCount), fmt.Sprint(count))
			}
		} else {
			if count := len(instance.ServiceAccounts); count != 1 {
				return fmt.Errorf("Wrong number of ServiceAccounts: expected 1, got %d", count)
			}

			if scount := len(instance.ServiceAccounts[0].Scopes); scount == scopeCount {
				return nil
			} else {
				return fmt.Errorf("Scope count expected: %s, but got %s", fmt.Sprint(scopeCount), fmt.Sprint(scount))
			}
		}
	}
}

func testAccCheckComputeInstanceHasSubnet(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instance.NetworkInterfaces {
			if i.Subnetwork == "" {
				return fmt.Errorf("no subnet")
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasAnyNetworkIP(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instance.NetworkInterfaces {
			if i.NetworkIP == "" {
				return fmt.Errorf("no network_ip")
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasNetworkIP(instance *compute.Instance, networkIP string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instance.NetworkInterfaces {
			if i.NetworkIP != networkIP {
				return fmt.Errorf("Wrong network_ip found: expected %v, got %v", networkIP, i.NetworkIP)
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasNetworkPerformanceConfig(instance *compute.Instance, bandwidthTier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.NetworkPerformanceConfig == nil {
			return fmt.Errorf("Expected instance to have network performance config, but it was nil")
		}
		if instance.NetworkPerformanceConfig.TotalEgressBandwidthTier != bandwidthTier {
			return fmt.Errorf("Incorrect network_performance_config.total_egress_bandwidth_tier found: expected %v, got %v", bandwidthTier, instance.NetworkPerformanceConfig.TotalEgressBandwidthTier)
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasMultiNic(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(instance.NetworkInterfaces) < 2 {
			return fmt.Errorf("only saw %d nics", len(instance.NetworkInterfaces))
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasGuestAccelerator(instance *compute.Instance, acceleratorType string, acceleratorCount int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(instance.GuestAccelerators) != 1 {
			return fmt.Errorf("Expected only one guest accelerator")
		}

		if !strings.HasSuffix(instance.GuestAccelerators[0].AcceleratorType, acceleratorType) {
			return fmt.Errorf("Wrong accelerator type: expected %v, got %v", acceleratorType, instance.GuestAccelerators[0].AcceleratorType)
		}

		if instance.GuestAccelerators[0].AcceleratorCount != acceleratorCount {
			return fmt.Errorf("Wrong accelerator acceleratorCount: expected %d, got %d", acceleratorCount, instance.GuestAccelerators[0].AcceleratorCount)
		}

		return nil
	}
}

func testAccCheckComputeInstanceLacksGuestAccelerator(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(instance.GuestAccelerators) > 0 {
			return fmt.Errorf("Expected no guest accelerators")
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasMinCpuPlatform(instance *compute.Instance, minCpuPlatform string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.MinCpuPlatform != minCpuPlatform {
			return fmt.Errorf("Wrong minimum CPU platform: expected %s, got %s", minCpuPlatform, instance.MinCpuPlatform)
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasMachineType(instance *compute.Instance, machineType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		instanceMachineType := tpgresource.GetResourceNameFromSelfLink(instance.MachineType)
		if instanceMachineType != machineType {
			return fmt.Errorf("Wrong machine type: expected %s, got %s", machineType, instanceMachineType)
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasAliasIpRange(instance *compute.Instance, subnetworkRangeName, iPCidrRange string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, networkInterface := range instance.NetworkInterfaces {
			for _, aliasIpRange := range networkInterface.AliasIpRanges {
				if aliasIpRange.SubnetworkRangeName == subnetworkRangeName && (aliasIpRange.IpCidrRange == iPCidrRange || tpgresource.IpCidrRangeDiffSuppress("ip_cidr_range", aliasIpRange.IpCidrRange, iPCidrRange, nil)) {
					return nil
				}
			}
		}

		return fmt.Errorf("Alias ip range with name %s and cidr %s not present", subnetworkRangeName, iPCidrRange)
	}
}

func testAccCheckComputeInstanceHasAssignedNatIP(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_instance" {
			continue
		}
		ip := rs.Primary.Attributes["network_interface.0.access_config.0.nat_ip"]
		if ip == "" {
			return fmt.Errorf("No assigned NatIP for instance %s", rs.Primary.Attributes["name"])
		}
	}
	return nil
}

func testAccCheckComputeInstanceHasConfiguredDeletionProtection(instance *compute.Instance, configuredDeletionProtection bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.DeletionProtection != configuredDeletionProtection {
			return fmt.Errorf("Wrong deletion protection flag: expected %t, got %t", configuredDeletionProtection, instance.DeletionProtection)
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasReservationAffinity(instance *compute.Instance, reservationType string, specificReservationNames ...string) resource.TestCheckFunc {
	if len(specificReservationNames) > 1 {
		panic("too many specificReservationNames provided in test")
	}

	return func(*terraform.State) error {
		if instance.ReservationAffinity == nil {
			return fmt.Errorf("expected instance to have reservation affinity, but it was nil")
		}

		if instance.ReservationAffinity.ConsumeReservationType != reservationType {
			return fmt.Errorf("Wrong reservationAffinity consumeReservationType: expected %s, got, %s", reservationType, instance.ReservationAffinity.ConsumeReservationType)
		}

		if len(specificReservationNames) > 0 {
			const reservationNameKey = "compute.googleapis.com/reservation-name"
			if instance.ReservationAffinity.Key != reservationNameKey {
				return fmt.Errorf("Wrong reservationAffinity key: expected %s, got, %s", reservationNameKey, instance.ReservationAffinity.Key)
			}
			if len(instance.ReservationAffinity.Values) != 1 || instance.ReservationAffinity.Values[0] != specificReservationNames[0] {
				return fmt.Errorf("Wrong reservationAffinity values: expected %s, got, %s", specificReservationNames, instance.ReservationAffinity.Values)
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasShieldedVmConfig(instance *compute.Instance, enableSecureBoot bool, enableVtpm bool, enableIntegrityMonitoring bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		if instance.ShieldedInstanceConfig.EnableSecureBoot != enableSecureBoot {
			return fmt.Errorf("Wrong shieldedVmConfig enableSecureBoot: expected %t, got, %t", enableSecureBoot, instance.ShieldedInstanceConfig.EnableSecureBoot)
		}

		if instance.ShieldedInstanceConfig.EnableVtpm != enableVtpm {
			return fmt.Errorf("Wrong shieldedVmConfig enableVtpm: expected %t, got, %t", enableVtpm, instance.ShieldedInstanceConfig.EnableVtpm)
		}

		if instance.ShieldedInstanceConfig.EnableIntegrityMonitoring != enableIntegrityMonitoring {
			return fmt.Errorf("Wrong shieldedVmConfig enableIntegrityMonitoring: expected %t, got, %t", enableIntegrityMonitoring, instance.ShieldedInstanceConfig.EnableIntegrityMonitoring)
		}
		return nil
	}
}

func testAccCheckComputeInstanceHasConfidentialInstanceConfig(instance *compute.Instance, EnableConfidentialCompute bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		if instance.ConfidentialInstanceConfig.EnableConfidentialCompute != EnableConfidentialCompute {
			return fmt.Errorf("Wrong ConfidentialInstanceConfig EnableConfidentialCompute: expected %t, got, %t", EnableConfidentialCompute, instance.ConfidentialInstanceConfig.EnableConfidentialCompute)
		}

		return nil
	}
}

func testAccCheckComputeInstanceLacksShieldedVmConfig(instance *compute.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.ShieldedInstanceConfig != nil {
			return fmt.Errorf("Expected no shielded vm config")
		}

		return nil
	}
}

func testAccCheckComputeInstanceHasStatus(instance *compute.Instance, status string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Status != status {
			return fmt.Errorf("Instance has not status %s, status: %s", status, instance.Status)
		}
		return nil
	}
}

func testAccComputeInstance_basic(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]
  desired_status  = "RUNNING"

  //deletion_protection = false is implicit in this config due to default value

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo            = "bar"
    baz            = "qux"
    startup-script = "echo Hello"
  }

  labels = {
    my_key       = "my_value"
    my_other_key = "my_other_value"
  }
}
`, instance)
}

func testAccComputeInstance_basic2(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, instance)
}

func testAccComputeInstance_basic3(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, instance)
}

func testAccComputeInstance_basic4(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, instance)
}

func testAccComputeInstance_basic5(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, instance)
}

func testAccComputeInstance_resourceManagerTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
  parent = "projects/%{project}"
  short_name = "foobarbaz%{random_suffix}"
  description = "For foo/bar resources."
}

resource "google_tags_tag_value" "value" {
  parent = "tagKeys/${google_tags_tag_key.key.name}"
  short_name = "foo%{random_suffix}"
  description = "For foo resources."
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%{instance_name}"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["tag-key", "tag-value"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
      resource_manager_tags = {
        "tagKeys/${google_tags_tag_key.key.name}" = "tagValues/${google_tags_tag_value.value.name}"
      }
    }
  }

  params {
    resource_manager_tags = {
      "tagKeys/${google_tags_tag_key.key.name}" = "tagValues/${google_tags_tag_value.value.name}"
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, context)
}

func testAccComputeInstance_basic_deletionProtectionFalse(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name                = "%s"
  machine_type        = "e2-medium"
  zone                = "us-central1-a"
  can_ip_forward      = false
  tags                = ["foo", "bar"]
  deletion_protection = false

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}
`, instance)
}

func testAccComputeInstance_basic_deletionProtectionTrue(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name                = "%s"
  machine_type        = "e2-medium"
  zone                = "us-central1-a"
  can_ip_forward      = false
  tags                = ["foo", "bar"]
  deletion_protection = true

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}
`, instance)
}

// Update zone to ForceNew, and change metadata k/v entirely
// Generates diff mismatch
func testAccComputeInstance_forceNewAndChangeMetadata(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-b"
  tags         = ["baz"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }

  metadata = {
    qux = "true"
  }
}
`, instance)
}

// Update metadata, tags, and network_interface
func testAccComputeInstance_update(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = true
  tags           = ["baz"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }

  metadata = {
    bar            = "baz"
    startup-script = "echo Hello"
  }

  labels = {
    only_me = "nothing_else"
  }
}
`, instance)
}

func testAccComputeInstance_ip(ip, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_address" "foo" {
  name = "%s"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"
  tags         = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
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
`, ip, instance)
}

func testAccComputeInstance_ipv6(ip, instance, record string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_subnetwork" "subnetwork-ipv6" {
  name          = "%s-subnetwork"

  ip_cidr_range = "10.0.0.0/22"
  region        = "us-west2"

  stack_type       = "IPV4_IPV6"
  ipv6_access_type = "EXTERNAL"

  network       = google_compute_network.custom-test.id
}

resource "google_compute_network" "custom-test" {
  name                    = "%s-network"
  auto_create_subnetworks = false
}

resource "google_compute_address" "foo" {
  name = "%s"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-west2-a"
  tags         = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnetwork-ipv6.name
    stack_type = "IPV4_IPV6"
    ipv6_access_config {
      network_tier = "PREMIUM"
      public_ptr_domain_name = "%s.gcp.tfacc.hashicorptest.com."
    }
  }

  metadata = {
    foo = "bar"
  }
}
`, instance, instance, ip, instance, record)
}

func testAccComputeInstance_ipv6ExternalReservation(instance string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "ipv6-address" {
  region             = "us-west2"
  name               = "%s-address"
  address_type       = "EXTERNAL"
  ip_version         = "IPV6"
  network_tier       = "PREMIUM"
  ipv6_endpoint_type = "VM"
  subnetwork         = google_compute_subnetwork.subnetwork-ipv6.name
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_subnetwork" "subnetwork-ipv6" {
  name          = "%s-subnetwork"

  ip_cidr_range = "10.0.0.0/22"
  region        = "us-west2"

  stack_type       = "IPV4_IPV6"
  ipv6_access_type = "EXTERNAL"

  network       = google_compute_network.custom-test.id
}

resource "google_compute_network" "custom-test" {
  name                    = "%s-network"
  auto_create_subnetworks = false
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-west2-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork                    = google_compute_subnetwork.subnetwork-ipv6.name
    stack_type                    = "IPV4_IPV6"
    ipv6_access_config {
      external_ipv6               = google_compute_address.ipv6-address.address
      external_ipv6_prefix_length = 96
      name                        = "external-ipv6-access-config"
      network_tier                = "PREMIUM"
    }
  }
}
`, instance, instance, instance, instance)
}

func testAccComputeInstance_PTRRecord(record, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"
  tags         = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
    access_config {
      public_ptr_domain_name = "%s.gcp.tfacc.hashicorptest.com."
    }
  }

  metadata = {
    foo = "bar"
  }
}
`, instance, record)
}

func testAccComputeInstance_networkTier(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
    access_config {
      network_tier = "STANDARD"
    }
  }
}
`, instance)
}

func testAccComputeInstance_disks_encryption(bootEncryptionKey string, diskNameToEncryptionKey map[string]*compute.CustomerEncryptionKey, instance, suffix string) string {
	diskNames := []string{}
	for k := range diskNameToEncryptionKey {
		diskNames = append(diskNames, k)
	}
	sort.Strings(diskNames)
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    raw_key = "%s"
  }
}

resource "google_compute_disk" "foobar2" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    raw_key = "%s"
  }
}

resource "google_compute_disk" "foobar3" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    raw_key = "%s"
  }
}

resource "google_compute_disk" "foobar4" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
    disk_encryption_key_raw = "%s"
  }

  attached_disk {
    source                  = google_compute_disk.foobar.self_link
    disk_encryption_key_raw = "%s"
  }

  attached_disk {
    source                  = google_compute_disk.foobar2.self_link
    disk_encryption_key_raw = "%s"
  }

  attached_disk {
    source = google_compute_disk.foobar4.self_link
  }

  attached_disk {
    source                  = google_compute_disk.foobar3.self_link
    disk_encryption_key_raw = "%s"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  allow_stopping_for_update = true
}
`, diskNames[0], diskNameToEncryptionKey[diskNames[0]].RawKey,
		diskNames[1], diskNameToEncryptionKey[diskNames[1]].RawKey,
		diskNames[2], diskNameToEncryptionKey[diskNames[2]].RawKey,
		"tf-testd-"+suffix,
		instance, bootEncryptionKey,
		diskNameToEncryptionKey[diskNames[0]].RawKey, diskNameToEncryptionKey[diskNames[1]].RawKey, diskNameToEncryptionKey[diskNames[2]].RawKey)
}

func testAccComputeInstance_disks_encryption_restart(bootEncryptionKey string, diskNameToEncryptionKey map[string]*compute.CustomerEncryptionKey, instance string) string {
	diskNames := []string{}
	for k := range diskNameToEncryptionKey {
		diskNames = append(diskNames, k)
	}
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    raw_key = "%s"
  }
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
    disk_encryption_key_raw = "%s"
  }

  attached_disk {
    source                  = google_compute_disk.foobar.self_link
    disk_encryption_key_raw = "%s"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  allow_stopping_for_update = true
}
`, diskNames[0], diskNameToEncryptionKey[diskNames[0]].RawKey,
		instance, bootEncryptionKey,
		diskNameToEncryptionKey[diskNames[0]].RawKey)
}

func testAccComputeInstance_disks_encryption_restartUpdate(bootEncryptionKey string, diskNameToEncryptionKey map[string]*compute.CustomerEncryptionKey, instance string) string {
	diskNames := []string{}
	for k := range diskNameToEncryptionKey {
		diskNames = append(diskNames, k)
	}
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    raw_key = "%s"
  }
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-standard-2"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
    disk_encryption_key_raw = "%s"
  }

  attached_disk {
    source                  = google_compute_disk.foobar.self_link
    disk_encryption_key_raw = "%s"
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  allow_stopping_for_update = true
}
`, diskNames[0], diskNameToEncryptionKey[diskNames[0]].RawKey,
		instance, bootEncryptionKey,
		diskNameToEncryptionKey[diskNames[0]].RawKey)
}

func testAccComputeInstance_disks_kms(bootEncryptionKey string, diskNameToEncryptionKey map[string]*compute.CustomerEncryptionKey, instance, suffix string) string {
	diskNames := []string{}
	for k := range diskNameToEncryptionKey {
		diskNames = append(diskNames, k)
	}
	sort.Strings(diskNames)
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
  }
}

resource "google_compute_disk" "foobar2" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
  }
}

resource "google_compute_disk" "foobar3" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
  }
}

resource "google_compute_disk" "foobar4" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
    kms_key_self_link = "%s"
  }

  attached_disk {
    source            = google_compute_disk.foobar.self_link
    kms_key_self_link = "%s"
  }

  attached_disk {
    source            = google_compute_disk.foobar2.self_link
    kms_key_self_link = "%s"
  }

  attached_disk {
    source = google_compute_disk.foobar4.self_link
  }

  attached_disk {
    source = google_compute_disk.foobar3.self_link
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, diskNames[0], diskNameToEncryptionKey[diskNames[0]].KmsKeyName,
		diskNames[1], diskNameToEncryptionKey[diskNames[1]].KmsKeyName,
		diskNames[2], diskNameToEncryptionKey[diskNames[2]].KmsKeyName,
		"tf-testd-"+suffix,
		instance, bootEncryptionKey,
		diskNameToEncryptionKey[diskNames[0]].KmsKeyName, diskNameToEncryptionKey[diskNames[1]].KmsKeyName)
}

func testAccComputeInstance_instanceSchedule(instance, schedule string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_resource_policy" "instance_schedule" {
  name        = "%s"
  region      = "us-central1"
  instance_schedule_policy {
    vm_start_schedule {
      schedule = "1 1 1 1 1"
    }
    vm_stop_schedule {
      schedule = "2 2 2 2 2"
    }
    time_zone = "UTC"
  }
}
`, instance, schedule)
}

func testAccComputeInstance_addResourcePolicy(instance, schedule string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  resource_policies = [google_compute_resource_policy.instance_schedule.self_link]
}

resource "google_compute_resource_policy" "instance_schedule" {
  name        = "%s"
  region      = "us-central1"
  instance_schedule_policy {
    vm_start_schedule {
      schedule = "1 1 1 1 1"
    }
    vm_stop_schedule {
      schedule = "2 2 2 2 2"
    }
    time_zone = "UTC"
  }
}
`, instance, schedule)
}

func testAccComputeInstance_updateResourcePolicy(instance, schedule1, schedule2 string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  resource_policies = [google_compute_resource_policy.instance_schedule2.self_link]
}

resource "google_compute_resource_policy" "instance_schedule" {
  name        = "%s"
  region      = "us-central1"
  instance_schedule_policy {
    vm_start_schedule {
      schedule = "1 1 1 1 1"
    }
    vm_stop_schedule {
      schedule = "2 2 2 2 2"
    }
    time_zone = "UTC"
  }
}

resource "google_compute_resource_policy" "instance_schedule2" {
  name        = "%s"
  region      = "us-central1"
  instance_schedule_policy {
    vm_start_schedule {
      schedule = "2 2 2 2 2"
    }
    vm_stop_schedule {
      schedule = "3 3 3 3 3"
    }
    time_zone = "UTC"
  }
}
`, instance, schedule1, schedule2)
}

func testAccComputeInstance_removeResourcePolicy(instance, schedule1, schedule2 string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  resource_policies = null
}

resource "google_compute_resource_policy" "instance_schedule" {
  name        = "%s"
  region      = "us-central1"
  instance_schedule_policy {
    vm_start_schedule {
      schedule = "1 1 1 1 1"
    }
    vm_stop_schedule {
      schedule = "2 2 2 2 2"
    }
    time_zone = "UTC"
  }
}

resource "google_compute_resource_policy" "instance_schedule2" {
  name        = "%s"
  region      = "us-central1"
  instance_schedule_policy {
    vm_start_schedule {
      schedule = "2 2 2 2 2"
    }
    vm_stop_schedule {
      schedule = "3 3 3 3 3"
    }
    time_zone = "UTC"
  }
}
`, instance, schedule1, schedule2)
}

func testAccComputeInstance_attachedDisk(disk, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  attached_disk {
    source = google_compute_disk.foobar.name
  }

  network_interface {
    network = "default"
  }
}
`, disk, instance)
}

func testAccComputeInstance_attachedDisk_sourceUrl(disk, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  attached_disk {
    source = google_compute_disk.foobar.self_link
  }

  network_interface {
    network = "default"
  }
}
`, disk, instance)
}

func testAccComputeInstance_attachedDisk_modeRo(disk, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  attached_disk {
    source = google_compute_disk.foobar.self_link
    mode   = "READ_ONLY"
  }

  network_interface {
    network = "default"
  }
}
`, disk, instance)
}

func testAccComputeInstance_addAttachedDisk(disk, disk2, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_disk" "foobar2" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  attached_disk {
    source = google_compute_disk.foobar.name
  }

  attached_disk {
    source = google_compute_disk.foobar2.self_link
  }

  network_interface {
    network = "default"
  }
}
`, disk, disk2, instance)
}

func testAccComputeInstance_detachDisk(disk, disk2, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_disk" "foobar2" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  attached_disk {
    source = google_compute_disk.foobar.name
  }

  network_interface {
    network = "default"
  }
}
`, disk, disk2, instance)
}

func testAccComputeInstance_updateAttachedDiskEncryptionKey(disk, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"
  disk_encryption_key {
    raw_key = "c2Vjb25kNzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"
  }
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  attached_disk {
    source                  = google_compute_disk.foobar.name
    disk_encryption_key_raw = "c2Vjb25kNzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"
  }

  network_interface {
    network = "default"
  }
}
`, disk, instance)
}

func testAccComputeInstance_bootDisk_source(disk, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  zone  = "us-central1-a"
  image = data.google_compute_image.my_image.self_link
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    source = google_compute_disk.foobar.name
  }

  network_interface {
    network = "default"
  }
}
`, disk, instance)
}

func testAccComputeInstance_bootDisk_sourceUrl(disk, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  zone  = "us-central1-a"
  image = data.google_compute_image.my_image.self_link
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    source = google_compute_disk.foobar.self_link
  }

  network_interface {
    network = "default"
  }
}
`, disk, instance)
}

func testAccComputeInstance_bootDisk_type(instance string, diskType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
      type  = "%s"
    }
  }

  network_interface {
    network = "default"
  }
}
`, instance, diskType)
}

func testAccComputeInstance_bootDisk_mode(instance string, diskMode string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
      type  = "pd-ssd"
    }

    mode = "%s"
  }

  network_interface {
    network = "default"
  }
}
`, instance, diskMode)
}

func testAccComputeInstance_with375GbScratchDisk(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"   // can't be e2 because of local-ssd
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  scratch_disk {
    interface = "NVME"
  }

  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"
  }
}
`, instance)
}

func testAccComputeInstance_with18TbScratchDisk(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n2-standard-64"   // must be a large n2 to be paired with 18Tb local-ssd
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  scratch_disk {
    interface = "NVME"
    size      = 3000
  }

  scratch_disk {
    interface = "NVME"
    size      = 3000
  }

  scratch_disk {
    interface = "NVME"
    size      = 3000
  }

  scratch_disk {
    interface = "NVME"
    size      = 3000
  }

  scratch_disk {
    interface = "NVME"
    size      = 3000
  }

  scratch_disk {
    interface = "NVME"
    size      = 3000
  }

  network_interface {
    network = "default"
  }
}`, instance)
}

func testAccComputeInstance_serviceAccount(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = [
      "userinfo-email",
      "compute-ro",
      "storage-ro",
    ]
  }
}
`, instance)
}

func testAccComputeInstance_serviceAccount_update0(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_serviceAccount_update01(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = []
  }
  allow_stopping_for_update = true
}

data "google_compute_default_service_account" "default" {
}
`, instance)
}

func testAccComputeInstance_serviceAccount_update02(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  service_account {
    email = data.google_compute_default_service_account.default.email
    scopes = []
  }
  allow_stopping_for_update = true
}

data "google_compute_default_service_account" "default" {
}
`, instance)
}

func testAccComputeInstance_serviceAccount_update3(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = [
      "userinfo-email",
      "compute-ro",
      "storage-ro",
    ]
  }

  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_serviceAccount_update4(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}
resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"
  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }
  network_interface {
    network = "default"
  }
  service_account {
    scopes = [
      "userinfo-email",
    ]
  }
  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_scheduling(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    automatic_restart = false
  }
}
`, instance)
}

func testAccComputeInstance_schedulingUpdated(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    automatic_restart = false
    preemptible       = true
  }
}
`, instance)
}

func testAccComputeInstance_advancedMachineFeatures(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-10"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-2" // Nested Virt isn't supported on E2 and N2Ds https://cloud.google.com/compute/docs/instances/nested-virtualization/overview#restrictions and https://cloud.google.com/compute/docs/instances/disabling-smt#limitations
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  allow_stopping_for_update = true

}
`, instance)
}

func testAccComputeInstance_advancedMachineFeaturesUpdated(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-10"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-2" // Nested Virt isn't supported on E2 and N2Ds https://cloud.google.com/compute/docs/instances/nested-virtualization/overview#restrictions and https://cloud.google.com/compute/docs/instances/disabling-smt#limitations
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
  advanced_machine_features {
	threads_per_core = 1
	enable_nested_virtualization = true
	visible_core_count = 1
  }
  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_subnet_auto(suffix, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "inst-test-network" {
  name = "tf-test-network-%s"

  auto_create_subnetworks = true
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = google_compute_network.inst-test-network.name
    access_config {
    }
  }
}
`, suffix, instance)
}

func testAccComputeInstance_subnet_custom(suffix, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "inst-test-network" {
  name = "tf-test-network-%s"

  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "inst-test-subnetwork-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.inst-test-network.self_link
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.inst-test-subnetwork.self_link
    access_config {
    }
  }
}
`, suffix, suffix, instance)
}

func testAccComputeInstance_subnet_xpn(org, billingId, projectName, instance, suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

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

resource "google_compute_network" "inst-test-network" {
  name    = "tf-test-network-%s"
  project = google_compute_shared_vpc_host_project.host_project.project

  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "tf-test-subnetwork-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.inst-test-network.self_link
  project       = google_compute_shared_vpc_host_project.host_project.project
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"
  project      = google_compute_shared_vpc_service_project.service_project.service_project

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork         = google_compute_subnetwork.inst-test-subnetwork.name
    subnetwork_project = google_compute_subnetwork.inst-test-subnetwork.project
    access_config {
    }
  }
}
`, projectName, org, billingId, projectName, org, billingId, suffix, suffix, instance)
}

func testAccComputeInstance_networkIPAuto(suffix, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "inst-test-network" {
  name = "tf-test-network-%s"
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "tf-test-subnetwork-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.inst-test-network.self_link
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.inst-test-subnetwork.name
    access_config {
    }
  }
}
`, suffix, suffix, instance)
}

func testAccComputeInstance_network_ip_custom(suffix, instance, ipAddress string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "inst-test-network" {
  name = "tf-test-network-%s"
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "tf-test-subnetwork-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.inst-test-network.self_link
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.inst-test-subnetwork.name
    network_ip = "%s"
    access_config {
    }
  }
}
`, suffix, suffix, instance, ipAddress)
}

func testAccComputeInstance_private_image_family(disk, family, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  zone  = "us-central1-a"
  image = data.google_compute_image.my_image.self_link
}

resource "google_compute_image" "foobar" {
  name        = "%s-1"
  source_disk = google_compute_disk.foobar.self_link
  family      = "%s"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = google_compute_image.foobar.family
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }
}
`, disk, family, family, instance)
}

func testAccComputeInstance_networkPerformanceConfig(disk string, image string, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  zone  = "us-central1-a"
  image = data.google_compute_image.my_image.self_link
}

resource "google_compute_image" "foobar" {
  name              = "%s"
  source_disk       = google_compute_disk.foobar.self_link
  guest_os_features {
    type = "GVNIC"
  }
  guest_os_features {
    type = "VIRTIO_SCSI_MULTIQUEUE"
  }
	guest_os_features {
    type = "UEFI_COMPATIBLE"
   }
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n2-standard-2"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = google_compute_image.foobar.self_link
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral IP
    }
  }

  network_performance_config {
    total_egress_bandwidth_tier = "DEFAULT"
  }
}
`, disk, image, instance)
}

func testAccComputeInstance_multiNic(instance, network, subnetwork string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.inst-test-subnetwork.name
    access_config {
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_network" "inst-test-network" {
  name = "%s"
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.inst-test-network.self_link
}
`, instance, network, subnetwork)
}

func testAccComputeInstance_nictype(image, instance, nictype string) string {
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

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  //deletion_protection = false is implicit in this config due to default value

  boot_disk {
    initialize_params {
	  image = google_compute_image.example.id
    }
  }

  network_interface {
	network = "default"
	nic_type = "%s"
  }

  metadata = {
    foo            = "bar"
    baz            = "qux"
    startup-script = "echo Hello"
  }

  labels = {
    my_key       = "my_value"
    my_other_key = "my_other_value"
  }
}
`, image, instance, nictype)
}

func testAccComputeInstance_guestAccelerator(instance string, count uint8) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"   // can't be e2 because of guest_accelerator
  zone         = "us-east1-d"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
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
`, instance, count)
}

func testAccComputeInstance_minCpuPlatform(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"   // can't be e2 because of min_cpu_platform
  zone         = "us-east1-d"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  min_cpu_platform = "Intel Haswell"
  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_minCpuPlatform_remove(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-micro"
  zone         = "us-east1-d"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  min_cpu_platform = "AuToMaTiC"
  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_primaryAliasIpRange(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-east1-d"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"

    alias_ip_range {
      ip_cidr_range = "/24"
    }
  }
}
`, instance)
}

func testAccComputeInstance_secondaryAliasIpRange(network, subnet, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "inst-test-network" {
  name = "%s"
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.inst-test-network.self_link
  secondary_ip_range {
    range_name    = "inst-test-secondary"
    ip_cidr_range = "172.16.0.0/20"
  }
  secondary_ip_range {
    range_name    = "inst-test-tertiary"
    ip_cidr_range = "10.1.0.0/16"
  }
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-east1-d"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.inst-test-subnetwork.self_link

    alias_ip_range {
      subnetwork_range_name = google_compute_subnetwork.inst-test-subnetwork.secondary_ip_range[0].range_name
      ip_cidr_range         = "172.16.0.0/24"
    }

    alias_ip_range {
      subnetwork_range_name = google_compute_subnetwork.inst-test-subnetwork.secondary_ip_range[1].range_name
      ip_cidr_range         = "10.1.0.0/20"
    }
  }
}
`, network, subnet, instance)
}

func testAccComputeInstance_secondaryAliasIpRangeUpdate(network, subnet, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_network" "inst-test-network" {
  name = "%s"
}

resource "google_compute_subnetwork" "inst-test-subnetwork" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.inst-test-network.self_link
  secondary_ip_range {
    range_name    = "inst-test-secondary"
    ip_cidr_range = "172.16.0.0/20"
  }
  secondary_ip_range {
    range_name    = "inst-test-tertiary"
    ip_cidr_range = "10.1.0.0/16"
  }
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-east1-d"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.inst-test-subnetwork.self_link
    alias_ip_range {
      ip_cidr_range = "10.0.1.0/24"
    }
  }
}
`, network, subnet, instance)
}

func testAccComputeInstance_hostname(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  hostname = "%s.test"
}
`, instance, instance)
}

// Set fields that require stopping the instance: machine_type, min_cpu_platform, and service_account
func testAccComputeInstance_stopInstanceToUpdate(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"   // can't be e2 because of min_cpu_platform
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  min_cpu_platform = "Intel Broadwell"
  service_account {
    scopes = [
      "userinfo-email",
      "compute-ro",
      "storage-ro",
    ]
  }

  allow_stopping_for_update = true
}
`, instance)
}

// Update fields that require stopping the instance: machine_type, min_cpu_platform, and service_account
func testAccComputeInstance_stopInstanceToUpdate2(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-2"   // can't be e2 because of min_cpu_platform
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  min_cpu_platform = "Intel Skylake"
  service_account {
    scopes = [
      "userinfo-email",
      "compute-ro",
    ]
  }

  allow_stopping_for_update = true
}
`, instance)
}

// Remove fields that require stopping the instance: min_cpu_platform and service_account (machine_type is Required)
func testAccComputeInstance_stopInstanceToUpdate3(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-2"   // can't be e2 because of min_cpu_platform
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_withoutNodeAffinities(instance, nodeTemplate, nodeGroup string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-8"   // can't be e2 because of sole tenancy
  zone         = "us-central1-a"
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_node_template" "nodetmpl" {
  name   = "%s"
  region = "us-central1"

  node_affinity_labels = {
    tfacc = "test"
  }

  node_type = "n1-node-96-624"

  cpu_overcommit_type = "ENABLED"
}

resource "google_compute_node_group" "nodes" {
  name = "%s"
  zone = "us-central1-a"

  size          = 1
  node_template = google_compute_node_template.nodetmpl.self_link
}
`, instance, nodeTemplate, nodeGroup)
}

func testAccComputeInstance_soleTenantNodeAffinities(instance, nodeTemplate, nodeGroup string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-8"   // can't be e2 because of sole tenancy
  zone         = "us-central1-a"
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    node_affinities {
      key      = "tfacc"
      operator = "IN"
      values   = ["test"]
    }

    node_affinities {
      key      = "tfacc"
      operator = "NOT_IN"
      values   = ["not_here"]
    }

    node_affinities {
      key      = "compute.googleapis.com/node-group-name"
      operator = "IN"
      values   = [google_compute_node_group.nodes.name]
    }

    min_node_cpus = 4
  }
}

resource "google_compute_node_template" "nodetmpl" {
  name   = "%s"
  region = "us-central1"

  node_affinity_labels = {
    tfacc = "test"
  }

  node_type = "n1-node-96-624"

  cpu_overcommit_type = "ENABLED"
}

resource "google_compute_node_group" "nodes" {
  name = "%s"
  zone = "us-central1-a"

  size          = 1
  node_template = google_compute_node_template.nodetmpl.self_link
}
`, instance, nodeTemplate, nodeGroup)
}

func testAccComputeInstance_soleTenantNodeAffinitiesUpdated(instance, nodeTemplate, nodeGroup string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-8"   // can't be e2 because of sole tenancy
  zone         = "us-central1-a"
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    node_affinities {
      key      = "tfacc"
      operator = "IN"
      values   = ["test", "updatedlabel"]
    }

    node_affinities {
      key      = "tfacc"
      operator = "NOT_IN"
      values   = ["not_here"]
    }

    node_affinities {
      key      = "compute.googleapis.com/node-group-name"
      operator = "IN"
      values   = [google_compute_node_group.nodes.name]
    }

    min_node_cpus = 6
  }
}

resource "google_compute_node_template" "nodetmpl" {
  name   = "%s"
  region = "us-central1"

  node_affinity_labels = {
    tfacc = "test"
  }

  node_type = "n1-node-96-624"

  cpu_overcommit_type = "ENABLED"
}

resource "google_compute_node_group" "nodes" {
  name = "%s"
  zone = "us-central1-a"

  size          = 1
  node_template = google_compute_node_template.nodetmpl.self_link
}
`, instance, nodeTemplate, nodeGroup)
}

func testAccComputeInstance_soleTenantNodeAffinitiesReduced(instance, nodeTemplate, nodeGroup string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-8"   // can't be e2 because of sole tenancy
  zone         = "us-central1-a"
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    node_affinities {
      key      = "tfacc"
      operator = "IN"
      values   = ["test", "updatedlabel"]
    }

    node_affinities {
      key      = "compute.googleapis.com/node-group-name"
      operator = "IN"
      values   = [google_compute_node_group.nodes.name]
    }

    min_node_cpus = 6
  }
}

resource "google_compute_node_template" "nodetmpl" {
  name   = "%s"
  region = "us-central1"

  node_affinity_labels = {
    tfacc = "test"
  }

  node_type = "n1-node-96-624"

  cpu_overcommit_type = "ENABLED"
}

resource "google_compute_node_group" "nodes" {
  name = "%s"
  zone = "us-central1-a"

  size          = 1
  node_template = google_compute_node_template.nodetmpl.self_link
}
`, instance, nodeTemplate, nodeGroup)
}

func testAccComputeInstance_reservationAffinity_nonSpecificReservationConfig(instanceName, reservationType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  reservation_affinity {
    type = "%s"
  }
}`, instanceName, reservationType)
}

func testAccComputeInstance_reservationAffinity_specificReservationConfig(instanceName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      machine_type = "n1-standard-1"
    }
  }
  specific_reservation_required = true
}

resource "google_compute_instance" "foobar" {
  name         = "%[1]s"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  reservation_affinity {
    type = "SPECIFIC_RESERVATION"

	specific_reservation {
		key    = "compute.googleapis.com/reservation-name"
		values = ["%[1]s"]
	}
  }
}`, instanceName)
}

func testAccComputeInstance_shieldedVmConfig(instance string, enableSecureBoot bool, enableVtpm bool, enableIntegrityMonitoring bool) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "centos-7"
  project = "centos-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  shielded_instance_config {
    enable_secure_boot          = %t
    enable_vtpm                 = %t
    enable_integrity_monitoring = %t
  }

  allow_stopping_for_update = true
}
`, instance, enableSecureBoot, enableVtpm, enableIntegrityMonitoring)
}

func testAccComputeInstanceConfidentialInstanceConfig(instance string, enableConfidentialCompute bool) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family    = "ubuntu-2004-lts"
  project   = "ubuntu-os-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n2d-standard-2"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
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
`, instance, enableConfidentialCompute)
}

func testAccComputeInstance_enableDisplay(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "centos-7"
  project = "centos-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  enable_display = true

  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_enableDisplayUpdated(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "centos-7"
  project = "centos-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  enable_display = false

  allow_stopping_for_update = true
}
`, instance)
}

func testAccComputeInstance_machineType_desiredStatus_allowStoppingForUpdate(
	instance, machineType, desiredStatus string,
	allowStoppingForUpdate bool,
) string {
	desiredStatusConfigSection := ""
	if desiredStatus != "" {
		desiredStatusConfigSection = fmt.Sprintf(
			"desired_status = \"%s\"",
			desiredStatus,
		)
	}

	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-11"
	project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
	name           = "%s"
	machine_type   = "%s"
	zone           = "us-central1-a"
	can_ip_forward = false
	tags           = ["foo", "bar"]

	boot_disk {
		initialize_params{
			image = "${data.google_compute_image.my_image.self_link}"
		}
	}

	network_interface {
		network = "default"
	}

	%s

	metadata = {
		foo = "bar"
	}

	allow_stopping_for_update = %t
}
`, instance, machineType, desiredStatusConfigSection, allowStoppingForUpdate)
}

func testAccComputeInstance_desiredStatusTerminatedUpdate(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-11"
	project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
	name           = "%s"
	machine_type   = "e2-medium"
	zone           = "us-central1-a"
	can_ip_forward = false
	tags           = ["baz"]

	boot_disk {
		initialize_params{
			image = "${data.google_compute_image.my_image.self_link}"
		}
	}

	network_interface {
		network = "default"
	}

	desired_status = "TERMINATED"

	metadata = {
		bar = "baz"
	}

	labels = {
		only_me = "nothing_else"
	}
}
`, instance)
}

func testAccComputeInstance_resourcePolicyCollocate(instance, suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "c2-standard-4"
  zone           = "us-east4-b"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  //deletion_protection = false is implicit in this config due to default value

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    # Instances with resource policies do not support live migration.
    on_host_maintenance = "TERMINATE"
    automatic_restart = false
  }

  resource_policies = [google_compute_resource_policy.foo.self_link]
}

resource "google_compute_instance" "second" {
  name           = "%s-2"
  machine_type   = "c2-standard-4"
  zone           = "us-east4-b"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  //deletion_protection = false is implicit in this config due to default value

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    # Instances with resource policies do not support live migration.
    on_host_maintenance = "TERMINATE"
    automatic_restart = false
  }

  resource_policies = [google_compute_resource_policy.foo.self_link]
}

resource "google_compute_resource_policy" "foo" {
  name   = "tf-test-policy-%s"
  region = "us-east4"
  group_placement_policy {
    vm_count = 2
    collocation = "COLLOCATED"
  }
}

`, instance, instance, suffix)
}

func testAccComputeInstance_subnetworkUpdate(suffix, instance string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-11"
		project = "debian-cloud"
	}

	resource "google_compute_network" "inst-test-network" {
		name = "tf-test-network-%s"
		auto_create_subnetworks = false
	}

	resource "google_compute_network" "inst-test-network2" {
		name = "tf-test-network2-%s"
		auto_create_subnetworks = false
	}

	resource "google_compute_subnetwork" "inst-test-subnetwork" {
		name          = "tf-test-compute-subnet-%s"
		ip_cidr_range = "10.0.0.0/16"
		region        = "us-east1"
		network       = google_compute_network.inst-test-network.id
		secondary_ip_range {
			range_name    = "inst-test-secondary"
			ip_cidr_range = "172.16.0.0/20"
		}
		secondary_ip_range {
			range_name    = "inst-test-tertiary"
			ip_cidr_range = "10.1.0.0/16"
		}
	}

	resource "google_compute_subnetwork" "inst-test-subnetwork2" {
		name          = "tf-test-compute-subnet2-%s"
		ip_cidr_range = "10.3.0.0/16"
		region        = "us-east1"
		network       = google_compute_network.inst-test-network2.id
		secondary_ip_range {
			range_name    = "inst-test-secondary2"
			ip_cidr_range = "173.16.0.0/20"
		}
		secondary_ip_range {
			range_name    = "inst-test-tertiary2"
			ip_cidr_range = "10.4.0.0/16"
		}
	}

	resource "google_compute_instance" "foobar" {
		name         = "%s"
		machine_type = "e2-medium"
		zone         = "us-east1-d"
		allow_stopping_for_update = true

		boot_disk {
			initialize_params {
				image = data.google_compute_image.my_image.id
			}
		}

		network_interface {
			subnetwork = google_compute_subnetwork.inst-test-subnetwork.id
			access_config {
				network_tier = "STANDARD"
			}
			alias_ip_range {
				subnetwork_range_name = google_compute_subnetwork.inst-test-subnetwork.secondary_ip_range[0].range_name
				ip_cidr_range         = "172.16.0.0/24"
			}

			alias_ip_range {
				subnetwork_range_name = google_compute_subnetwork.inst-test-subnetwork.secondary_ip_range[1].range_name
				ip_cidr_range         = "10.1.0.0/20"
			}
		}
	}
`, suffix, suffix, suffix, suffix, instance)
}

func testAccComputeInstance_subnetworkUpdateTwo(suffix, instance string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-11"
		project = "debian-cloud"
	}

	resource "google_compute_network" "inst-test-network" {
		name = "tf-test-network-%s"
		auto_create_subnetworks = false
	}

	resource "google_compute_network" "inst-test-network2" {
		name = "tf-test-network2-%s"
		auto_create_subnetworks = false
	}

	resource "google_compute_subnetwork" "inst-test-subnetwork" {
		name          = "tf-test-compute-subnet-%s"
		ip_cidr_range = "10.0.0.0/16"
		region        = "us-east1"
		network       = google_compute_network.inst-test-network.id
		secondary_ip_range {
			range_name    = "inst-test-secondary"
			ip_cidr_range = "172.16.0.0/20"
		}
		secondary_ip_range {
			range_name    = "inst-test-tertiary"
			ip_cidr_range = "10.1.0.0/16"
		}
	}

	resource "google_compute_subnetwork" "inst-test-subnetwork2" {
		name          = "tf-test-compute-subnet2-%s"
		ip_cidr_range = "10.3.0.0/16"
		region        = "us-east1"
		network       = google_compute_network.inst-test-network2.id
		secondary_ip_range {
			range_name    = "inst-test-secondary2"
			ip_cidr_range = "173.16.0.0/20"
		}
		secondary_ip_range {
			range_name    = "inst-test-tertiary2"
			ip_cidr_range = "10.4.0.0/16"
		}
	}

	resource "google_compute_instance" "foobar" {
		name         = "%s"
		machine_type = "e2-medium"
		zone         = "us-east1-d"
		allow_stopping_for_update = true

		boot_disk {
			initialize_params {
				image = data.google_compute_image.my_image.id
			}
		}

		network_interface {
			subnetwork = google_compute_subnetwork.inst-test-subnetwork2.id
			network_ip = "10.3.0.3"
			access_config {
				network_tier = "STANDARD"
			}
			alias_ip_range {
				subnetwork_range_name = google_compute_subnetwork.inst-test-subnetwork2.secondary_ip_range[0].range_name
				ip_cidr_range         = "173.16.0.0/24"
			}
		}
	}
`, suffix, suffix, suffix, suffix, instance)
}

func testAccComputeInstance_queueCountSet(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-11"
	project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
	name         = "%s"
	machine_type = "e2-medium"
	zone         = "us-east1-d"
	allow_stopping_for_update = true

	boot_disk {
		initialize_params {
			image = data.google_compute_image.my_image.id
		}
	}

  network_interface {
    network = "default"
    queue_count = 2
  }
}
`, instance)
}

func testAccComputeInstance_spotVM(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family    = "ubuntu-2004-lts"
  project   = "ubuntu-os-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    provisioning_model = "SPOT"
    automatic_restart = false
    preemptible = true
		instance_termination_action = "STOP"
  }
}
`, instance)
}

func testAccComputeInstance_spotVM_maxRunDuration(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family    = "ubuntu-2004-lts"
  project   = "ubuntu-os-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    provisioning_model = "SPOT"
    automatic_restart = false
    preemptible = true
    instance_termination_action = "DELETE"
  }
}
`, instance)
}

func testAccComputeInstance_localSsdRecoveryTimeout(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family    = "ubuntu-2004-lts"
  project   = "ubuntu-os-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  scheduling {
    local_ssd_recovery_timeout {
        nanos = 0
        seconds = 3600
    }
  }

}
`, instance)
}

func testAccComputeInstance_metadataStartupScript(instance, machineType, metadata string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "%s"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "%s"
  }
  metadata_startup_script = "echo hi > /test.txt"
  allow_stopping_for_update = true
}
`, instance, machineType, metadata)
}
