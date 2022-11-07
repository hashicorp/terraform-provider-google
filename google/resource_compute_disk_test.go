package google

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/api/compute/v1"
)

func TestDiskImageDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		// Full & partial links
		"matching self_link with different api version": {
			Old:                "https://www.googleapis.com/compute/beta/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			ExpectDiffSuppress: true,
		},
		"matching image partial self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			ExpectDiffSuppress: true,
		},
		"matching image partial no project self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "global/images/debian-8-jessie-v20171213",
			ExpectDiffSuppress: true,
		},
		"different image self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-7-jessie-v20171213",
			ExpectDiffSuppress: false,
		},
		"different image partial self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "projects/debian-cloud/global/images/debian-7-jessie-v20171213",
			ExpectDiffSuppress: false,
		},
		"different image partial no project self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "global/images/debian-7-jessie-v20171213",
			ExpectDiffSuppress: false,
		},
		// Image name
		"matching image name": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "debian-8-jessie-v20171213",
			ExpectDiffSuppress: true,
		},
		"different image name": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "debian-7-jessie-v20171213",
			ExpectDiffSuppress: false,
		},
		// Image short hand
		"matching image short hand": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "debian-cloud/debian-8-jessie-v20171213",
			ExpectDiffSuppress: true,
		},
		"matching image short hand but different project": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "different-cloud/debian-8-jessie-v20171213",
			ExpectDiffSuppress: false,
		},
		"different image short hand": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "debian-cloud/debian-7-jessie-v20171213",
			ExpectDiffSuppress: false,
		},
		// Image Family
		"matching image family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "family/debian-8",
			ExpectDiffSuppress: true,
		},
		"matching image family self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/family/debian-8",
			ExpectDiffSuppress: true,
		},
		"matching unconventional image family self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1404-trusty-v20180122",
			New:                "https://www.googleapis.com/compute/v1/projects/projects/ubuntu-os-cloud/global/images/family/ubuntu-1404-lts",
			ExpectDiffSuppress: true,
		},
		"matching image family partial self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "projects/debian-cloud/global/images/family/debian-8",
			ExpectDiffSuppress: true,
		},
		"matching unconventional image family partial self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1404-trusty-v20180122",
			New:                "projects/ubuntu-os-cloud/global/images/family/ubuntu-1404-lts",
			ExpectDiffSuppress: true,
		},
		"matching image family partial no project self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "global/images/family/debian-8",
			ExpectDiffSuppress: true,
		},
		"matching image family short hand": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "debian-cloud/debian-8",
			ExpectDiffSuppress: true,
		},
		"matching image family short hand with project short name": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "debian/debian-8",
			ExpectDiffSuppress: true,
		},
		"matching unconventional image family short hand": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1404-trusty-v20180122",
			New:                "ubuntu-os-cloud/ubuntu-1404-lts",
			ExpectDiffSuppress: true,
		},
		"matching unconventional image family - minimal": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-1804-bionic-v20180705",
			New:                "ubuntu-minimal-1804-lts",
			ExpectDiffSuppress: true,
		},
		"matching unconventional image family - cos": {
			Old:                "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-85-13310-1209-17",
			New:                "cos-85-lts",
			ExpectDiffSuppress: true,
		},
		"different image family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "family/debian-7",
			ExpectDiffSuppress: false,
		},
		"different image family self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/family/debian-7",
			ExpectDiffSuppress: false,
		},
		"different image family partial self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "projects/debian-cloud/global/images/family/debian-7",
			ExpectDiffSuppress: false,
		},
		"different image family partial no project self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "global/images/family/debian-7",
			ExpectDiffSuppress: false,
		},
		"matching image family but different project in self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "https://www.googleapis.com/compute/v1/projects/other-cloud/global/images/family/debian-8",
			ExpectDiffSuppress: false,
		},
		"different image family but different project in partial self link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "projects/other-cloud/global/images/family/debian-8",
			ExpectDiffSuppress: false,
		},
		"different image family short hand": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "debian-cloud/debian-7",
			ExpectDiffSuppress: false,
		},
		"matching image family shorthand but different project": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20171213",
			New:                "different-cloud/debian-8",
			ExpectDiffSuppress: false,
		},
		// arm images
		"matching image opensuse arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/opensuse-cloud/global/images/opensuse-leap-15-4-v20220713-arm64",
			New:                "opensuse-leap-arm64",
			ExpectDiffSuppress: true,
		},
		"matching image sles arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/suse-cloud/global/images/sles-15-sp4-v20220713-arm64",
			New:                "sles-15-arm64",
			ExpectDiffSuppress: true,
		},
		"matching image ubuntu arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1804-bionic-arm64-v20220712",
			New:                "ubuntu-1804-lts-arm64",
			ExpectDiffSuppress: true,
		},
		"matching image ubuntu-minimal arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2004-focal-arm64-v20220713",
			New:                "ubuntu-minimal-2004-lts-arm64",
			ExpectDiffSuppress: true,
		},
		"matching image debian arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-11-bullseye-arm64-v20220719",
			New:                "debian-11-arm64",
			ExpectDiffSuppress: true,
		},
		"different architecture image opensuse arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/opensuse-cloud/global/images/opensuse-leap-15-4-v20220713-arm64",
			New:                "opensuse-leap",
			ExpectDiffSuppress: false,
		},
		"different architecture image sles arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/suse-cloud/global/images/sles-15-sp4-v20220713-arm64",
			New:                "sles-15",
			ExpectDiffSuppress: false,
		},
		"different architecture image ubuntu arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1804-bionic-arm64-v20220712",
			New:                "ubuntu-1804-lts",
			ExpectDiffSuppress: false,
		},
		"different architecture image ubuntu-minimal arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2004-focal-arm64-v20220713",
			New:                "ubuntu-minimal-2004-lts",
			ExpectDiffSuppress: false,
		},
		"different architecture image debian arm64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-11-bullseye-arm64-v20220719",
			New:                "debian-11",
			ExpectDiffSuppress: false,
		},
		"different architecture image opensuse arm64 family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/opensuse-cloud/global/images/opensuse-leap-15-2-v20200702",
			New:                "opensuse-leap-arm64",
			ExpectDiffSuppress: false,
		},
		"different architecture image sles arm64 family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/suse-cloud/global/images/sles-15-sp4-v20220722-x86-64",
			New:                "sles-15-arm64",
			ExpectDiffSuppress: false,
		},
		"different architecture image ubuntu arm64 family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1804-bionic-v20220712",
			New:                "ubuntu-1804-lts-arm64",
			ExpectDiffSuppress: false,
		},
		"different architecture image ubuntu-minimal arm64 family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2004-focal-v20220713",
			New:                "ubuntu-minimal-2004-lts-arm64",
			ExpectDiffSuppress: false,
		},
		"different architecture image debian arm64 family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-11-bullseye-v20220719",
			New:                "debian-11-arm64",
			ExpectDiffSuppress: false,
		},
		// amd images
		"matching image ubuntu amd64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-2210-kinetic-amd64-v20221022",
			New:                "ubuntu-2210-amd64",
			ExpectDiffSuppress: true,
		},
		"matching image ubuntu-minimal amd64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2210-kinetic-amd64-v20221022",
			New:                "ubuntu-minimal-2210-amd64",
			ExpectDiffSuppress: true,
		},
		"different architecture image ubuntu amd64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-2210-kinetic-amd64-v20221022",
			New:                "ubuntu-2210",
			ExpectDiffSuppress: false,
		},
		"different architecture image ubuntu-minimal amd64 self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2210-kinetic-amd64-v20221022",
			New:                "ubuntu-minimal-2210",
			ExpectDiffSuppress: false,
		},
		"different architecture image ubuntu amd64 family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-2210-kinetic-v20221022",
			New:                "ubuntu-2210-amd64",
			ExpectDiffSuppress: false,
		},
		"different architecture image ubuntu-minimal amd64 family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2210-kinetic-v20221022",
			New:                "ubuntu-minimal-2210-amd64",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			if diskImageDiffSuppress("image", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
				t.Fatalf("%q => %q expect DiffSuppress to return %t", tc.Old, tc.New, tc.ExpectDiffSuppress)
			}
		})
	}
}

// Test that all the naming pattern for public images are supported.
func TestAccComputeDisk_imageDiffSuppressPublicVendorsFamilyNames(t *testing.T) {
	t.Parallel()

	if os.Getenv(TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", TestEnvVar)
	}

	config := getInitializedConfig(t)

	for _, publicImageProject := range imageMap {
		token := ""
		for paginate := true; paginate; {
			resp, err := config.NewComputeClient(config.userAgent).Images.List(publicImageProject).Filter("deprecated.replacement ne .*images.*").PageToken(token).Do()
			if err != nil {
				t.Fatalf("Can't list public images for project %q", publicImageProject)
			}

			for _, image := range resp.Items {
				if !diskImageDiffSuppress("image", image.SelfLink, "family/"+image.Family, nil) {
					t.Errorf("should suppress diff for image %q and family %q", image.SelfLink, image.Family)
				}
			}
			token := resp.NextPageToken
			paginate = token != ""
		}
	}
}

func TestAccComputeDisk_timeout(t *testing.T) {
	// Vcr speeds up test, so it doesn't time out
	skipIfVcr(t)
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-disk-%d", randInt(t))
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeDisk_timeout(diskName),
				ExpectError: regexp.MustCompile("timeout"),
			},
		},
	})
}

func TestAccComputeDisk_update(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_basic(diskName),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDisk_updated(diskName),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_fromSnapshot(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	firstDiskName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	snapshotName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	projectName := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_fromSnapshot(projectName, firstDiskName, snapshotName, diskName, "self_link"),
			},
			{
				ResourceName:      "google_compute_disk.seconddisk",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDisk_fromSnapshot(projectName, firstDiskName, snapshotName, diskName, "name"),
			},
			{
				ResourceName:      "google_compute_disk.seconddisk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_encryption(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	var disk compute.Disk

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_encryption(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", getTestProjectFromEnv(), &disk),
					testAccCheckEncryptionKey(
						t, "google_compute_disk.foobar", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_encryptionKMS(t *testing.T) {
	t.Parallel()

	kms := BootstrapKMSKey(t)
	pid := getTestProjectFromEnv()
	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	importID := fmt.Sprintf("%s/%s/%s", pid, "us-central1-a", diskName)
	var disk compute.Disk

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_encryptionKMS(pid, diskName, kms.CryptoKey.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", pid, &disk),
					testAccCheckEncryptionKey(
						t, "google_compute_disk.foobar", &disk),
				),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportStateId:     importID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_deleteDetach(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_deleteDetach(instanceName, diskName),
			},
			{
				ResourceName:      "google_compute_disk.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// this needs to be a second step so we refresh and see the instance
			// listed as attached to the disk; the instance is created after the
			// disk. and the disk's properties aren't refreshed unless there's
			// another step
			{
				Config: testAccComputeDisk_deleteDetach(instanceName, diskName),
			},
			{
				ResourceName:      "google_compute_disk.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_deleteDetachIGM(t *testing.T) {
	// Randomness in instance template
	skipIfVcr(t)
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	diskName2 := fmt.Sprintf("tf-test-%s", randString(t, 10))
	mgrName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_deleteDetachIGM(diskName, mgrName),
			},
			{
				ResourceName:      "google_compute_disk.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// this needs to be a second step so we refresh and see the instance
			// listed as attached to the disk; the instance is created after the
			// disk. and the disk's properties aren't refreshed unless there's
			// another step
			{
				Config: testAccComputeDisk_deleteDetachIGM(diskName, mgrName),
			},
			{
				ResourceName:      "google_compute_disk.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change the disk name to recreate the instances
			{
				Config: testAccComputeDisk_deleteDetachIGM(diskName2, mgrName),
			},
			{
				ResourceName:      "google_compute_disk.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Add the extra step like before
			{
				Config: testAccComputeDisk_deleteDetachIGM(diskName2, mgrName),
			},
			{
				ResourceName:      "google_compute_disk.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_pdExtremeImplicitProvisionedIops(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_pdExtremeImplicitProvisionedIops(diskName),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeDiskExists(t *testing.T, n, p string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)

		found, err := config.NewComputeClient(config.userAgent).Disks.Get(
			p, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Disk not found")
		}

		*disk = *found

		return nil
	}
}

func testAccCheckEncryptionKey(t *testing.T, n string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attr := rs.Primary.Attributes["disk_encryption_key.0.sha256"]
		if disk.DiskEncryptionKey == nil {
			return fmt.Errorf("Disk %s has mismatched encryption key.\nTF State: %+v\nGCP State: <empty>", n, attr)
		} else if attr != disk.DiskEncryptionKey.Sha256 {
			return fmt.Errorf("Disk %s has mismatched encryption key.\nTF State: %+v.\nGCP State: %+v",
				n, attr, disk.DiskEncryptionKey.Sha256)
		}
		return nil
	}
}

func TestAccComputeDisk_cloneDisk(t *testing.T) {
	t.Parallel()
	pid := getTestProjectFromEnv()
	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	var disk compute.Disk

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_diskClone(diskName, "self_link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.disk-clone", pid, &disk),
				),
			},
			{
				ResourceName:      "google_compute_disk.disk-clone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeDisk_basic(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
  labels = {
    my-label = "my-label-value"
  }
}
`, diskName)
}

func testAccComputeDisk_timeout(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  type  = "pd-ssd"
  zone  = "us-central1-a"

  timeouts {
    create = ".5s"
  }
}
`, diskName)
}

func testAccComputeDisk_updated(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 100
  type  = "pd-ssd"
  zone  = "us-central1-a"
  labels = {
    my-label    = "my-updated-label-value"
    a-new-label = "a-new-label-value"
  }
}
`, diskName)
}

func testAccComputeDisk_fromSnapshot(projectName, firstDiskName, snapshotName, diskName, ref_selector string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name    = "%s-d1"
  image   = data.google_compute_image.my_image.self_link
  size    = 50
  type    = "pd-ssd"
  zone    = "us-central1-a"
  project = "%s"
}

resource "google_compute_snapshot" "snapdisk" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  project     = "%s"
}

resource "google_compute_disk" "seconddisk" {
  name     = "%s-d2"
  snapshot = google_compute_snapshot.snapdisk.%s
  type     = "pd-ssd"
  zone     = "us-central1-a"
  project  = "%s"
}
`, firstDiskName, projectName, snapshotName, projectName, diskName, ref_selector, projectName)
}

func testAccComputeDisk_encryption(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
  disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
`, diskName)
}

func testAccComputeDisk_encryptionKMS(pid, diskName, kmsKey string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_project_iam_member" "kms-project-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@compute-system.iam.gserviceaccount.com"
}

resource "google_compute_disk" "foobar" {
  depends_on = [google_project_iam_member.kms-project-binding]

  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
  }
}
`, pid, diskName, kmsKey)
}

func testAccComputeDisk_deleteDetach(instanceName, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foo" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance" "bar" {
  name         = "%s"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  attached_disk {
    source = google_compute_disk.foo.self_link
  }

  network_interface {
    network = "default"
  }
}
`, diskName, instanceName)
}

func testAccComputeDisk_deleteDetachIGM(diskName, mgrName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foo" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_instance_template" "template" {
  machine_type = "g1-small"

  disk {
    boot        = true
    source      = google_compute_disk.foo.name
    auto_delete = false
  }

  network_interface {
    network = "default"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_instance_group_manager" "manager" {
  name               = "%s"
  base_instance_name = "tf-test-disk-igm"
  version {
    instance_template = google_compute_instance_template.template.self_link
    name              = "primary"
  }
  update_policy {
    minimal_action        = "RESTART"
    type                  = "PROACTIVE"
    max_unavailable_fixed = 1
  }
  zone        = "us-central1-a"
  target_size = 1

  // block on instances being ready so that when they get deleted, we don't try
  // to continue interacting with them in other resources
  wait_for_instances = true
}
`, diskName, mgrName)
}

func testAccComputeDisk_pdExtremeImplicitProvisionedIops(diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
  name  = "%s"
  type = "pd-extreme"
  size = 1
}
`, diskName)
}

func testAccComputeDisk_diskClone(diskName, refSelector string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-11"
		project = "debian-cloud"
	}

	resource "google_compute_disk" "foobar" {
		name  = "%s"
		image = data.google_compute_image.my_image.self_link
		size  = 50
		type  = "pd-ssd"
		zone  = "us-central1-a"
		labels = {
			my-label = "my-label-value"
		}
	}

	resource "google_compute_disk" "disk-clone" {
		name  = "%s"
		source_disk = google_compute_disk.foobar.%s
		type  = "pd-ssd"
		zone  = "us-central1-a"
		labels = {
			my-label = "my-label-value"
		}
	}
`, diskName, diskName+"-clone", refSelector)
}
