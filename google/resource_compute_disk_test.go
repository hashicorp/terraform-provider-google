package google

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	}

	for tn, tc := range cases {
		if diskImageDiffSuppress("image", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

// Test that all the naming pattern for public images are supported.
func TestAccComputeDisk_imageDiffSuppressPublicVendorsFamilyNames(t *testing.T) {
	t.Parallel()

	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use %s=1 to enable", resource.TestEnvVar))
	}

	config := getInitializedConfig(t)

	for _, publicImageProject := range imageMap {
		token := ""
		for paginate := true; paginate; {
			resp, err := config.clientCompute.Images.List(publicImageProject).Filter("deprecated.replacement ne .*images.*").PageToken(token).Do()
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

func TestAccComputeDisk_basic(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_basic(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					testAccCheckComputeDiskHasLabel(&disk, "my-label", "my-label-value"),
					testAccCheckComputeDiskHasLabelFingerprint(&disk, "google_compute_disk.foobar"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_timeout(t *testing.T) {
	t.Parallel()

	diskName := acctest.RandomWithPrefix("tf-test-disk")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccComputeDisk_timeout(diskName),
				ExpectError: regexp.MustCompile("timeout"),
			},
		},
	})
}

func TestAccComputeDisk_update(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_basic(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "size", "50"),
					testAccCheckComputeDiskHasLabel(&disk, "my-label", "my-label-value"),
					testAccCheckComputeDiskHasLabelFingerprint(&disk, "google_compute_disk.foobar"),
				),
			},
			{
				Config: testAccComputeDisk_updated(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "size", "100"),
					testAccCheckComputeDiskHasLabel(&disk, "my-label", "my-updated-label-value"),
					testAccCheckComputeDiskHasLabel(&disk, "a-new-label", "a-new-label-value"),
					testAccCheckComputeDiskHasLabelFingerprint(&disk, "google_compute_disk.foobar"),
				),
			},
		},
	})
}

func TestAccComputeDisk_fromSnapshot(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	firstDiskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	projectName := getTestProjectFromEnv()

	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_fromSnapshot(projectName, firstDiskName, snapshotName, diskName, "self_link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.seconddisk", &disk),
				),
			},
			resource.TestStep{
				Config: testAccComputeDisk_fromSnapshot(projectName, firstDiskName, snapshotName, diskName, "name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.seconddisk", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_encryption(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_encryption(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					testAccCheckEncryptionKey(
						"google_compute_disk.foobar", &disk),
				),
			},
			// Update from top-level attribute to nested.
			resource.TestStep{
				Config: testAccComputeDisk_encryptionMigrate(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					testAccCheckEncryptionKey(
						"google_compute_disk.foobar", &disk),
				),
			},
			// Update from nested attribute back to top-level.
			resource.TestStep{
				Config: testAccComputeDisk_encryption(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					testAccCheckEncryptionKey(
						"google_compute_disk.foobar", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_deleteDetach(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetach(instanceName, diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
				),
			},
			// this needs to be a second step so we refresh and see the instance
			// listed as attached to the disk; the instance is created after the
			// disk. and the disk's properties aren't refreshed unless there's
			// another step
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetach(instanceName, diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
					testAccCheckComputeDiskInstances(
						"google_compute_disk.foo", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_deleteDetachIGM(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	diskName2 := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	mgrName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetachIGM(diskName, mgrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
				),
			},
			// this needs to be a second step so we refresh and see the instance
			// listed as attached to the disk; the instance is created after the
			// disk. and the disk's properties aren't refreshed unless there's
			// another step
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetachIGM(diskName, mgrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
					testAccCheckComputeDiskInstances(
						"google_compute_disk.foo", &disk),
				),
			},
			// Change the disk name to recreate the instances
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetachIGM(diskName2, mgrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
				),
			},
			// Add the extra step like before
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetachIGM(diskName2, mgrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
					testAccCheckComputeDiskInstances(
						"google_compute_disk.foo", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_computeDiskUserRegex(t *testing.T) {

	shouldPass := []string{

		"https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1/instances/123",
		"https://www.googleapis.com/compute/v1/projects/123123/zones/us-central1/instances/123",
		"https://www.googleapis.com/compute/v1/projects/hashicorptest.net:project-123/zones/us-central1/instances/123",
		"https://www.googleapis.com/compute/v1/projects/123/zones/456/instances/789",
	}

	shouldFail := []string{
		"https://www.googleapis.com/compute/v1/projects/project#/zones/us-central1/instances/123",
		"https://www.googleapis.com/compute/v1/projects/project/zones/us-central#/instances/123",
		"https://www.googleapis.com/compute/v1/projects/project/zones/us-central1/instances/?",
		"https://www.googleapis.com/compute/v1/projects/foo.com:bar:baz/zones/us-central1/instances/?",
		"https://www.googleapis.com/compute/v1/projects/foo.com:/zones/us-central1/instances/?",
	}

	for _, element := range shouldPass {
		if !computeDiskUserRegex.MatchString(element) {
			t.Error("computeDiskUserRegex should match on '" + element + "' but doesn't")
		}
	}

	for _, element := range shouldFail {
		if computeDiskUserRegex.MatchString(element) {
			t.Error("computeDiskUserRegex shouldn't match on '" + element + "' but does")
		}
	}

}

func testAccCheckComputeDiskDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_disk" {
			continue
		}

		_, err := config.clientCompute.Disks.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Disk still exists")
		}
	}

	return nil
}

func testAccCheckComputeDiskExists(n string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		p := getTestProjectFromEnv()
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Disks.Get(
			p, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Disk not found")
		}

		*disk = *found

		return nil
	}
}

func testAccCheckComputeDiskHasLabel(disk *compute.Disk, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val, ok := disk.Labels[key]
		if !ok {
			return fmt.Errorf("Label with key %s not found", key)
		}

		if val != value {
			return fmt.Errorf("Label value did not match for key %s: expected %s but found %s", key, value, val)
		}
		return nil
	}
}

func testAccCheckComputeDiskHasLabelFingerprint(disk *compute.Disk, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		state := s.RootModule().Resources[resourceName]
		if state == nil {
			return fmt.Errorf("Unable to find resource named %s", resourceName)
		}

		labelFingerprint := state.Primary.Attributes["label_fingerprint"]
		if labelFingerprint != disk.LabelFingerprint {
			return fmt.Errorf("Label fingerprints do not match: api returned %s but state has %s",
				disk.LabelFingerprint, labelFingerprint)
		}

		return nil
	}
}

func testAccCheckEncryptionKey(n string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attr := rs.Primary.Attributes["disk_encryption_key_sha256"]
		if disk.DiskEncryptionKey == nil {
			return fmt.Errorf("Disk %s has mismatched encryption key.\nTF State: %+v\nGCP State: <empty>", n, attr)
		} else if attr != disk.DiskEncryptionKey.Sha256 {
			return fmt.Errorf("Disk %s has mismatched encryption key.\nTF State: %+v.\nGCP State: %+v",
				n, attr, disk.DiskEncryptionKey.Sha256)
		}
		return nil
	}
}

func testAccCheckComputeDiskInstances(n string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attr := rs.Primary.Attributes["users.#"]
		if strconv.Itoa(len(disk.Users)) != attr {
			return fmt.Errorf("Disk %s has mismatched users.\nTF State: %+v\nGCP State: %+v", n, rs.Primary.Attributes["users"], disk.Users)
		}

		for pos, user := range disk.Users {
			if rs.Primary.Attributes["users."+strconv.Itoa(pos)] != user {
				return fmt.Errorf("Disk %s has mismatched users.\nTF State: %+v.\nGCP State: %+v",
					n, rs.Primary.Attributes["users"], disk.Users)
			}
		}
		return nil
	}
}

func testAccComputeDisk_basic(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
	labels {
		my-label = "my-label-value"
	}
}`, diskName)
}

func testAccComputeDisk_timeout(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name  = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	type  = "pd-ssd"
	zone  = "us-central1-a"

	timeouts {
		create = "1s"
	}
}`, diskName)
}

func testAccComputeDisk_updated(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 100
	type = "pd-ssd"
	zone = "us-central1-a"
	labels {
		my-label = "my-updated-label-value"
		a-new-label = "a-new-label-value"
	}
}`, diskName)
}

func testAccComputeDisk_fromSnapshot(projectName, firstDiskName, snapshotName, diskName, ref_selector string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "d1-%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
	project = "%s"
}

resource "google_compute_snapshot" "snapdisk" {
	name = "%s"
	source_disk = "${google_compute_disk.foobar.name}"
	zone = "us-central1-a"
	project = "%s"
}

resource "google_compute_disk" "seconddisk" {
	name = "d2-%s"
	snapshot = "${google_compute_snapshot.snapdisk.%s}"
	type = "pd-ssd"
	zone = "us-central1-a"
	project = "%s"
}`, firstDiskName, projectName, snapshotName, projectName, diskName, ref_selector, projectName)
}

func testAccComputeDisk_encryption(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
	disk_encryption_key_raw = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
}`, diskName)
}

func testAccComputeDisk_encryptionMigrate(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
	disk_encryption_key {
		raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	}
}`, diskName)
}

func testAccComputeDisk_deleteDetach(instanceName, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foo" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
}

resource "google_compute_instance" "bar" {
	name = "%s"
	machine_type = "n1-standard-1"
	zone = "us-central1-a"

	boot_disk {
		initialize_params {
			image = "${data.google_compute_image.my_image.self_link}"
		}
	}

	attached_disk {
		source = "${google_compute_disk.foo.self_link}"
	}

	network_interface {
		network = "default"
	}
}`, diskName, instanceName)
}

func testAccComputeDisk_deleteDetachIGM(diskName, mgrName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foo" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
}

resource "google_compute_instance_template" "template" {
	machine_type = "g1-small"

	disk {
		boot        = true
		source      = "${google_compute_disk.foo.name}"
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
  base_instance_name = "disk-igm"
  instance_template  = "${google_compute_instance_template.template.self_link}"
  zone               = "us-central1-a"
  target_size        = 1
}`, diskName, mgrName)
}
