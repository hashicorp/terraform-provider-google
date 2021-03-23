package google

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func init() {
	resource.AddTestSweepers("ComputeInstance", &resource.Sweeper{
		Name: "ComputeInstance",
		F:    testSweepComputeInstance,
	})
}

// At the time of writing, the CI only passes us-central1 as the region.
// Since we can read all instances across zones, we don't really use this param.
func testSweepComputeInstance(region string) error {
	resourceName := "ComputeInstance"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	found, err := config.NewComputeClient(config.userAgent).Instances.AggregatedList(config.Project).Do()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request: %s", err)
		return nil
	}

	// Keep count of items that aren't sweepable for logging.
	nonPrefixCount := 0
	for zone, itemList := range found.Items {
		for _, instance := range itemList.Instances {
			if !isSweepableTestResource(instance.Name) {
				nonPrefixCount++
				continue
			}

			// Don't wait on operations as we may have a lot to delete
			_, err := config.NewComputeClient(config.userAgent).Instances.Delete(config.Project, GetResourceNameFromSelfLink(zone), instance.Name).Do()
			if err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error deleting %s resource %s : %s", resourceName, instance.Name, err)
			} else {
				log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, instance.Name)
			}
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonPrefixCount)
	}

	return nil
}

func computeInstanceImportStep(zone, instanceName string, additionalImportIgnores []string) resource.TestStep {
	// metadata is only read into state if set in the config
	// since importing doesn't know whether metadata.startup_script vs metadata_startup_script is set in the config,
	// it guesses metadata_startup_script
	ignores := []string{"metadata.%", "metadata.startup-script", "metadata_startup_script"}

	return resource.TestStep{
		ResourceName:            "google_compute_instance.foobar",
		ImportState:             true,
		ImportStateId:           fmt.Sprintf("%s/%s/%s", getTestProjectFromEnv(), zone, instanceName),
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: append(ignores, additionalImportIgnores...),
	}
}

func TestAccComputeInstance_basic1(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
					// by default, DeletionProtection is implicitly false. This should be false on any
					// instance resource without an explicit deletion_protection = true declaration.
					// Other tests check explicit true/false configs: TestAccComputeInstance_deletionProtectionExplicit[True | False]
					testAccCheckComputeInstanceHasConfiguredDeletionProtection(&instance, false),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"metadata.baz", "metadata.foo"}),
		},
	})
}

func TestAccComputeInstance_basic2(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

func TestAccComputeInstance_IP(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var ipName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

func TestAccComputeInstance_PTRRecord(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var ptrName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var ipName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	bootEncryptionKey := "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	bootEncryptionKeyHash := "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E="
	diskNameToEncryptionKey := map[string]*compute.CustomerEncryptionKey{
		fmt.Sprintf("tf-testd-%s", randString(t, 10)): {
			RawKey: "Ym9vdDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "awJ7p57H+uVZ9axhJjl1D3lfC2MgA/wnt/z88Ltfvss=",
		},
		fmt.Sprintf("tf-testd-%s", randString(t, 10)): {
			RawKey: "c2Vjb25kNzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "7TpIwUdtCOJpq2m+3nt8GFgppu6a2Xsj1t0Gexk13Yc=",
		},
		fmt.Sprintf("tf-testd-%s", randString(t, 10)): {
			RawKey: "dGhpcmQ2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "b3pvaS7BjDbCKeLPPTx7yXBuQtxyMobCHN1QJR43xeM=",
		},
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disks_encryption(bootEncryptionKey, diskNameToEncryptionKey, instanceName, randString(t, 10)),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	bootEncryptionKey := "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	bootEncryptionKeyHash := "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E="
	diskNameToEncryptionKey := map[string]*compute.CustomerEncryptionKey{
		fmt.Sprintf("tf-testd-%s", randString(t, 10)): {
			RawKey: "Ym9vdDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
			Sha256: "awJ7p57H+uVZ9axhJjl1D3lfC2MgA/wnt/z88Ltfvss=",
		},
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	kms := BootstrapKMSKey(t)

	bootKmsKeyName := kms.CryptoKey.Name
	diskNameToEncryptionKey := map[string]*compute.CustomerEncryptionKey{
		fmt.Sprintf("tf-testd-%s", randString(t, 10)): {
			KmsKeyName: kms.CryptoKey.Name,
		},
		fmt.Sprintf("tf-testd-%s", randString(t, 10)): {
			KmsKeyName: kms.CryptoKey.Name,
		},
		fmt.Sprintf("tf-testd-%s", randString(t, 10)): {
			KmsKeyName: kms.CryptoKey.Name,
		},
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disks_kms(getTestProjectFromEnv(), bootKmsKeyName, diskNameToEncryptionKey, instanceName, randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceDiskKmsEncryptionKey("google_compute_instance.foobar", &instance, bootKmsKeyName, diskNameToEncryptionKey),
				),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_attachedDisk(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", randString(t, 10))
	var diskName2 = fmt.Sprintf("tf-testd-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskType = "pd-ssd"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskMode = "READ_WRITE"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_bootDisk_mode(instanceName, diskMode),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_scratchDisk(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_scratchDisk(instanceName),
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

func TestAccComputeInstance_forceNewAndChangeMetadata(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

func TestAccComputeInstance_scheduling(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

func TestAccComputeInstance_soleTenantNodeAffinities(t *testing.T) {
	t.Parallel()

	var instanceName = fmt.Sprintf("tf-test-soletenant-%s", randString(t, 10))
	var templateName = fmt.Sprintf("tf-test-nodetmpl-%s", randString(t, 10))
	var groupName = fmt.Sprintf("tf-test-nodegroup-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_soleTenantNodeAffinities(instanceName, templateName, groupName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
			{
				Config: testAccComputeInstance_soleTenantNodeAffinitiesUpdated(instanceName, templateName, groupName),
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_subnet_auto(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_subnet_auto(randString(t, 10), instanceName),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_subnet_custom(randString(t, 10), instanceName),
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
	skipIfVcr(t)
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	projectName := fmt.Sprintf("tf-test-xpn-%d", time.Now().Unix())

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_subnet_xpn(org, billingId, projectName, instanceName, randString(t, 10)),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_networkIPAuto(randString(t, 10), instanceName),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var ipAddress = "10.0.200.200"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_network_ip_custom(randString(t, 10), instanceName, ipAddress),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))
	var diskName = fmt.Sprintf("tf-testd-%s", randString(t, 10))
	var familyName = fmt.Sprintf("tf-testf-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

func TestAccComputeInstance_forceChangeMachineTypeManually(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceUpdateMachineType(t, "google_compute_instance.foobar"),
				),
				ExpectNonEmptyPlan: true,
			},
			computeInstanceImportStep("us-central1-a", instanceName, []string{"metadata.baz", "metadata.foo"}),
		},
	})
}

func TestAccComputeInstance_multiNic(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	networkName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	subnetworkName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_minCpuPlatform(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(t, "google_compute_instance.foobar", &instance),
					testAccCheckComputeInstanceHasMinCpuPlatform(&instance, "Intel Haswell"),
				),
			},
			computeInstanceImportStep("us-east1-d", instanceName, []string{}),
		},
	})
}

func TestAccComputeInstance_deletionProtectionExplicitFalse(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	networkName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	subnetName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

	var instance computeBeta.Instance
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

	var instance computeBeta.Instance
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

	var instance computeBeta.Instance
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

const errorAllowStoppingMsg = "Changing the machine_type, min_cpu_platform, service_account, enable_display, shielded_instance_config, or network_interface.\\[#d\\].\\(network/subnetwork/subnetwork_project\\) on a started instance requires stopping it. To acknowledge this, please set allow_stopping_for_update = true in your config. You can also stop it by setting desired_status = \"TERMINATED\", but the instance will not be restarted after the update."

func TestAccComputeInstance_updateRunning_desiredStatusNotSet_notAllowStoppingForUpdate(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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
	var instanceName = fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_resourcePolicyCollocate(instanceName, randString(t, 10)),
			},
			computeInstanceImportStep("us-east4-b", instanceName, []string{"allow_stopping_for_update"}),
		},
	})
}

func TestAccComputeInstance_subnetworkUpdate(t *testing.T) {
	t.Parallel()
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	suffix := fmt.Sprintf("%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
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

func TestComputeInstance_networkIPCustomizedDiff(t *testing.T) {
	t.Parallel()

	d := &ResourceDiffMock{
		Before: map[string]interface{}{
			"network_interface.#": 0,
		},
		After: map[string]interface{}{
			"network_interface.#": 1,
		},
	}

	err := forceNewIfNetworkIPNotUpdatableFunc(d)
	if err != nil {
		t.Error(err)
	}

	if d.IsForceNew {
		t.Errorf("Expected not force new if network_interface array size changes")
	}

	type NetworkInterface struct {
		Network           string
		Subnetwork        string
		SubnetworkProject string
		NetworkIP         string
	}
	NIBefore := NetworkInterface{
		Network:           "a",
		Subnetwork:        "a",
		SubnetworkProject: "a",
		NetworkIP:         "a",
	}

	cases := map[string]struct {
		ExpectedForceNew bool
		Before           NetworkInterface
		After            NetworkInterface
	}{
		"NetworkIP only change": {
			ExpectedForceNew: true,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and Network change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "b",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and Subnetwork change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "b",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and SubnetworkProject change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "b",
				NetworkIP:         "b",
			},
		},
		"All change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "b",
				Subnetwork:        "b",
				SubnetworkProject: "b",
				NetworkIP:         "b",
			},
		},
		"No change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "a",
			},
		},
	}

	for tn, tc := range cases {
		d := &ResourceDiffMock{
			Before: map[string]interface{}{
				"network_interface.#":                    1,
				"network_interface.0.network":            tc.Before.Network,
				"network_interface.0.subnetwork":         tc.Before.Subnetwork,
				"network_interface.0.subnetwork_project": tc.Before.SubnetworkProject,
				"network_interface.0.network_ip":         tc.Before.NetworkIP,
			},
			After: map[string]interface{}{
				"network_interface.#":                    1,
				"network_interface.0.network":            tc.After.Network,
				"network_interface.0.subnetwork":         tc.After.Subnetwork,
				"network_interface.0.subnetwork_project": tc.After.SubnetworkProject,
				"network_interface.0.network_ip":         tc.After.NetworkIP,
			},
		}
		err := forceNewIfNetworkIPNotUpdatableFunc(d)
		if err != nil {
			t.Error(err)
		}
		if tc.ExpectedForceNew != d.IsForceNew {
			t.Errorf("%v: expected d.IsForceNew to be %v, but was %v", tn, tc.ExpectedForceNew, d.IsForceNew)
		}
	}
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

		config := googleProviderConfig(t)

		op, err := config.NewComputeClient(config.userAgent).Instances.Stop(config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return fmt.Errorf("Could not stop instance: %s", err)
		}
		err = computeOperationWaitTime(config, op, config.Project, "Waiting on stop", config.userAgent, 20*time.Minute)
		if err != nil {
			return fmt.Errorf("Could not stop instance: %s", err)
		}

		machineType := compute.InstancesSetMachineTypeRequest{
			MachineType: "zones/us-central1-a/machineTypes/f1-micro",
		}

		op, err = config.NewComputeClient(config.userAgent).Instances.SetMachineType(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"], &machineType).Do()
		if err != nil {
			return fmt.Errorf("Could not change machine type: %s", err)
		}
		err = computeOperationWaitTime(config, op, config.Project, "Waiting machine type change", config.userAgent, 20*time.Minute)
		if err != nil {
			return fmt.Errorf("Could not change machine type: %s", err)
		}
		return nil
	}
}

func testAccCheckComputeInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_instance" {
				continue
			}

			_, err := config.NewComputeClient(config.userAgent).Instances.Get(
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

	switch instance.(type) {
	case *compute.Instance:
		return testAccCheckComputeInstanceExistsInProject(t, n, getTestProjectFromEnv(), instance.(*compute.Instance))
	case *computeBeta.Instance:
		return testAccCheckComputeBetaInstanceExistsInProject(t, n, getTestProjectFromEnv(), instance.(*computeBeta.Instance))
	default:
		panic("Attempted to check existence of an Instance of unknown type.")
	}
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

		config := googleProviderConfig(t)

		found, err := config.NewComputeClient(config.userAgent).Instances.Get(
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

func testAccCheckComputeBetaInstanceExistsInProject(t *testing.T, n, p string, instance *computeBeta.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)

		found, err := config.NewComputeBetaClient(config.userAgent).Instances.Get(
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
		config := googleProviderConfig(t)

		// boot disk is named the same as the Instance
		disk, err := config.NewComputeClient(config.userAgent).Disks.Get(config.Project, "us-central1-a", instanceName).Do()
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
					expectedKey := diskNameToEncryptionKey[GetResourceNameFromSelfLink(disk.Source)].Sha256
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
			diskName := GetResourceNameFromSelfLink(rs.Primary.Attributes[fmt.Sprintf("attached_disk.%d.source", i)])
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
					expectedKey := diskNameToEncryptionKey[GetResourceNameFromSelfLink(disk.Source)].KmsKeyName
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
			diskName := GetResourceNameFromSelfLink(rs.Primary.Attributes[fmt.Sprintf("attached_disk.%d.source", i)])
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
		instanceMachineType := GetResourceNameFromSelfLink(instance.MachineType)
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
				if aliasIpRange.SubnetworkRangeName == subnetworkRangeName && (aliasIpRange.IpCidrRange == iPCidrRange || ipCidrRangeDiffSuppress("ip_cidr_range", aliasIpRange.IpCidrRange, iPCidrRange, nil)) {
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

func testAccCheckComputeInstanceHasShieldedVmConfig(instance *computeBeta.Instance, enableSecureBoot bool, enableVtpm bool, enableIntegrityMonitoring bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		if instance.ShieldedVmConfig.EnableSecureBoot != enableSecureBoot {
			return fmt.Errorf("Wrong shieldedVmConfig enableSecureBoot: expected %t, got, %t", enableSecureBoot, instance.ShieldedVmConfig.EnableSecureBoot)
		}

		if instance.ShieldedVmConfig.EnableVtpm != enableVtpm {
			return fmt.Errorf("Wrong shieldedVmConfig enableVtpm: expected %t, got, %t", enableVtpm, instance.ShieldedVmConfig.EnableVtpm)
		}

		if instance.ShieldedVmConfig.EnableIntegrityMonitoring != enableIntegrityMonitoring {
			return fmt.Errorf("Wrong shieldedVmConfig enableIntegrityMonitoring: expected %t, got, %t", enableIntegrityMonitoring, instance.ShieldedVmConfig.EnableIntegrityMonitoring)
		}
		return nil
	}
}

func testAccCheckComputeInstanceHasConfidentialInstanceConfig(instance *computeBeta.Instance, EnableConfidentialCompute bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		if instance.ConfidentialInstanceConfig.EnableConfidentialCompute != EnableConfidentialCompute {
			return fmt.Errorf("Wrong ConfidentialInstanceConfig EnableConfidentialCompute: expected %t, got, %t", EnableConfidentialCompute, instance.ConfidentialInstanceConfig.EnableConfidentialCompute)
		}

		return nil
	}
}

func testAccCheckComputeInstanceLacksShieldedVmConfig(instance *computeBeta.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.ShieldedVmConfig != nil {
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
  family  = "debian-9"
  project = "debian-cloud"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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

func testAccComputeInstance_basic_deletionProtectionFalse(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
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
  family  = "debian-9"
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

func testAccComputeInstance_PTRRecord(record, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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

func testAccComputeInstance_disks_kms(pid string, bootEncryptionKey string, diskNameToEncryptionKey map[string]*compute.CustomerEncryptionKey, instance, suffix string) string {
	diskNames := []string{}
	for k := range diskNameToEncryptionKey {
		diskNames = append(diskNames, k)
	}
	sort.Strings(diskNames)
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_project_iam_member" "kms-project-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com"
}

resource "google_compute_disk" "foobar" {
  depends_on = [google_project_iam_member.kms-project-binding]

  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
  }
}

resource "google_compute_disk" "foobar2" {
  depends_on = [google_project_iam_member.kms-project-binding]

  name = "%s"
  size = 10
  type = "pd-ssd"
  zone = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
  }
}

resource "google_compute_disk" "foobar3" {
  depends_on = [google_project_iam_member.kms-project-binding]

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
  depends_on = [google_project_iam_member.kms-project-binding]

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
`, pid, diskNames[0], diskNameToEncryptionKey[diskNames[0]].KmsKeyName,
		diskNames[1], diskNameToEncryptionKey[diskNames[1]].KmsKeyName,
		diskNames[2], diskNameToEncryptionKey[diskNames[2]].KmsKeyName,
		"tf-testd-"+suffix,
		instance, bootEncryptionKey,
		diskNameToEncryptionKey[diskNames[0]].KmsKeyName, diskNameToEncryptionKey[diskNames[1]].KmsKeyName)
}

func testAccComputeInstance_attachedDisk(disk, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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

func testAccComputeInstance_scratchDisk(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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

func testAccComputeInstance_serviceAccount(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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

func testAccComputeInstance_scheduling(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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
  family  = "debian-9"
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

func testAccComputeInstance_subnet_auto(suffix, instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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

func testAccComputeInstance_multiNic(instance, network, subnetwork string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
}
`, instance)
}

func testAccComputeInstance_primaryAliasIpRange(instance string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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
  family  = "debian-9"
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

func testAccComputeInstance_soleTenantNodeAffinities(instance, nodeTemplate, nodeGroup string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-8"   // can't be e2 because of sole tenancy
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
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-8"   // can't be e2 because of sole tenancy
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
	family  = "debian-9"
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
	family  = "debian-9"
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
  family  = "debian-9"
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
		family  = "debian-9"
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
		family  = "debian-9"
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
