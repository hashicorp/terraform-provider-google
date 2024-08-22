// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

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
		"matching image ubuntu amd64 canonical lts self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-2404-noble-amd64-v20240423",
			New:                "ubuntu-2404-lts-amd64",
			ExpectDiffSuppress: true,
		},
		"matching image ubuntu minimal amd64 canonical lts self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2404-noble-amd64-v20240423",
			New:                "ubuntu-minimal-2404-lts-amd64",
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
		"different image ubuntu amd64 canonical lts self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-2404-noble-amd64-v20240423",
			New:                "ubuntu-2404-lts",
			ExpectDiffSuppress: false,
		},
		"different image ubuntu minimal amd64 canonical lts self_link": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2404-noble-amd64-v20240423",
			New:                "ubuntu-minimal-2404-lts",
			ExpectDiffSuppress: false,
		},
		"different image ubuntu amd64 canonical lts family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-2404-noble-v20240423",
			New:                "ubuntu-2404-lts-amd64",
			ExpectDiffSuppress: false,
		},
		"different image ubuntu minimal amd64 canonical lts family": {
			Old:                "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-minimal-2404-noble-v20240423",
			New:                "ubuntu-minimal-2404-lts-amd64",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			if tpgcompute.DiskImageDiffSuppress("image", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
				t.Fatalf("%q => %q expect DiffSuppress to return %t", tc.Old, tc.New, tc.ExpectDiffSuppress)
			}
		})
	}
}

// Test that all the naming pattern for public images are supported.
func TestAccComputeDisk_imageDiffSuppressPublicVendorsFamilyNames(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}

	config := getInitializedConfig(t)

	for _, publicImageProject := range tpgcompute.ImageMap {
		token := ""
		// Hard limit on number of pages to prevent infinite loops
		// caused by the API always returning a pagination token
		page := 0
		maxPages := 10
		for paginate := true; paginate && page < maxPages; {
			resp, err := config.NewComputeClient(config.UserAgent).Images.List(publicImageProject).Filter("deprecated.replacement ne .*images.*").PageToken(token).Do()
			if err != nil {
				t.Fatalf("Can't list public images for project %q", publicImageProject)
			}

			for _, image := range resp.Items {
				if !tpgcompute.DiskImageDiffSuppress("image", image.SelfLink, "family/"+image.Family, nil) {
					t.Errorf("should suppress diff for image %q and family %q", image.SelfLink, image.Family)
				}
			}
			token := resp.NextPageToken
			paginate = token != ""
			page++
		}
	}
}

func TestAccComputeDisk_update(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskType := "pd-ssd"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_basic(diskName, diskType),
			},
			{
				ResourceName:            "google_compute_disk.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeDisk_updated(diskName, diskType),
			},
			{
				ResourceName:            "google_compute_disk.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccComputeDisk_fromTypeUrl(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskType := fmt.Sprintf("projects/%s/zones/us-central1-a/diskTypes/pd-ssd", envvar.GetTestProjectFromEnv())

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_basic(diskName, diskType),
			},
			{
				ResourceName:            "google_compute_disk.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccComputeDisk_pdHyperDiskProvisionedIopsLifeCycle(t *testing.T) {
	t.Parallel()

	context_1 := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"provisioned_iops": 10000,
		"lifecycle_bool":   true,
	}
	context_2 := map[string]interface{}{
		"random_suffix":    context_1["random_suffix"],
		"provisioned_iops": 11000,
		"lifecycle_bool":   true,
	}
	context_3 := map[string]interface{}{
		"random_suffix":    context_1["random_suffix"],
		"provisioned_iops": 11000,
		"lifecycle_bool":   false,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_pdHyperDiskProvisionedIopsLifeCycle(context_1),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDisk_pdHyperDiskProvisionedIopsLifeCycle(context_2),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDisk_pdHyperDiskProvisionedIopsLifeCycle(context_3),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_pdHyperDiskProvisionedThroughputLifeCycle(t *testing.T) {
	t.Parallel()

	context_1 := map[string]interface{}{
		"random_suffix":          acctest.RandString(t, 10),
		"provisioned_throughput": 180,
		"lifecycle_bool":         true,
	}
	context_2 := map[string]interface{}{
		"random_suffix":          context_1["random_suffix"],
		"provisioned_throughput": 20,
		"lifecycle_bool":         true,
	}
	context_3 := map[string]interface{}{
		"random_suffix":          context_1["random_suffix"],
		"provisioned_throughput": 20,
		"lifecycle_bool":         false,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_pdHyperDiskProvisionedThroughputLifeCycle(context_1),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDisk_pdHyperDiskProvisionedThroughputLifeCycle(context_2),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDisk_pdHyperDiskProvisionedThroughputLifeCycle(context_3),
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

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	firstDiskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	projectName := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
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

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_encryption(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", envvar.GetTestProjectFromEnv(), &disk),
					testAccCheckEncryptionKey(
						t, "google_compute_disk.foobar", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_encryptionKMS(t *testing.T) {
	t.Parallel()

	kms := acctest.BootstrapKMSKey(t)
	pid := envvar.GetTestProjectFromEnv()
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	importID := fmt.Sprintf("%s/%s/%s", pid, "us-central1-a", diskName)
	var disk compute.Disk

	if acctest.BootstrapPSARole(t, "service-", "compute-system", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_encryptionKMS(diskName, kms.CryptoKey.Name),
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

func TestAccComputeDisk_pdHyperDiskEnableConfidentialCompute(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"kms": acctest.BootstrapKMSKeyWithPurposeInLocationAndName(
			t,
			"ENCRYPT_DECRYPT",
			"us-central1",
			"tf-bootstrap-hyperdisk-key1").CryptoKey.Name, // regional KMS key
		"disk_size":            64,
		"confidential_compute": true,
	}

	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_pdHyperDiskEnableConfidentialCompute(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", envvar.GetTestProjectFromEnv(), &disk),
					testAccCheckEncryptionKey(
						t, "google_compute_disk.foobar", &disk),
				),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_deleteDetach(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
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
	acctest.SkipIfVcr(t)
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	diskName2 := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	mgrName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
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

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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

		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewComputeClient(config.UserAgent).Disks.Get(
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
	pid := envvar.GetTestProjectFromEnv()
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_diskClone(diskName, "self_link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.disk-clone", pid, &disk),
				),
			},
			{
				ResourceName:            "google_compute_disk.disk-clone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccComputeDisk_featuresUpdated(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_features(diskName),
			},
			{
				ResourceName:            "google_compute_disk.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccComputeDisk_featuresUpdated(diskName),
			},
			{
				ResourceName:            "google_compute_disk.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccComputeDisk_basic(diskName string, diskType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 50
  type  = "%s"
  zone  = "us-central1-a"
  labels = {
    my-label = "my-label-value"
  }
}
`, diskName, diskType)
}

func testAccComputeDisk_updated(diskName string, diskType string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 100
  type  = "%s"
  zone  = "us-central1-a"
  labels = {
    my-label    = "my-updated-label-value"
    a-new-label = "a-new-label-value"
  }
}
`, diskName, diskType)
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

func testAccComputeDisk_encryptionKMS(diskName, kmsKey string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"

  disk_encryption_key {
    kms_key_self_link = "%s"
  }
}
`, diskName, kmsKey)
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

func testAccComputeDisk_pdHyperDiskEnableConfidentialCompute(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_compute_disk" "foobar" {
		name                        = "tf-test-ecc-%{random_suffix}"
		size                        = %{disk_size}
		type                        = "hyperdisk-balanced"
		zone                        = "us-central1-a"
		enable_confidential_compute = %{confidential_compute}

		disk_encryption_key {
			kms_key_self_link       = "%{kms}"
		}

	}
`, context)
}

func testAccComputeDisk_pdHyperDiskProvisionedIopsLifeCycle(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_compute_disk" "foobar" {
		name                    = "tf-test-hyperdisk-%{random_suffix}"
		type                    = "hyperdisk-extreme"
		provisioned_iops        = %{provisioned_iops}
		size                    = 64
		lifecycle {
		  prevent_destroy       = %{lifecycle_bool}
		}
	  }
`, context)
}

func testAccComputeDisk_pdHyperDiskProvisionedThroughputLifeCycle(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_compute_disk" "foobar" {
		name                   = "tf-test-hyperdisk-%{random_suffix}"
		type                   = "hyperdisk-throughput"
		zone                   = "us-east4-c"
		provisioned_throughput = %{provisioned_throughput}
		size                   = 2048
		lifecycle {
		  prevent_destroy      = %{lifecycle_bool}
		}
	  }
`, context)
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

func TestAccComputeDisk_encryptionWithRSAEncryptedKey(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_encryptionWithRSAEncryptedKey(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar-1", envvar.GetTestProjectFromEnv(), &disk),
					testAccCheckEncryptionKey(
						t, "google_compute_disk.foobar-1", &disk),
				),
			},
		},
	})
}

func testAccComputeDisk_encryptionWithRSAEncryptedKey(diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar-1" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
  disk_encryption_key {
	rsa_encrypted_key = "fB6BS8tJGhGVDZDjGt1pwUo2wyNbkzNxgH1avfOtiwB9X6oPG94gWgenygitnsYJyKjdOJ7DyXLmxwQOSmnCYCUBWdKCSssyLV5907HL2mb5TfqmgHk5JcArI/t6QADZWiuGtR+XVXqiLa5B9usxFT2BTmbHvSKfkpJ7McCNc/3U0PQR8euFRZ9i75o/w+pLHFMJ05IX3JB0zHbXMV173PjObiV3ItSJm2j3mp5XKabRGSA5rmfMnHIAMz6stGhcuom6+bMri2u/axmPsdxmC6MeWkCkCmPjaKsVz1+uQUNCJkAnzesluhoD+R6VjFDm4WI7yYabu4MOOAOTaQXdEg=="
  }
}
`, diskName)
}

func testAccComputeDisk_features(diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
  name  = "%s"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
  labels = {
    my-label = "my-label-value"
  }

  guest_os_features {
    type = "SECURE_BOOT"
  }
}
`, diskName)
}

func testAccComputeDisk_featuresUpdated(diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
  name  = "%s"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
  labels = {
    my-label = "my-label-value"
  }

  guest_os_features {
    type = "SECURE_BOOT"
  }

  guest_os_features {
    type = "MULTI_IP_SUBNET"
  }
}
`, diskName)
}

func TestAccComputeDisk_attributionLabelOnCreation(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_attributionLabel(diskName, "true", "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.user-label", "foo"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.user-label", "foo"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "effective_labels.%", "2"),
				),
			},
			{
				Config: testAccComputeDisk_attributionLabelUpdated(diskName, "true", "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.user-label", "bar"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.user-label", "bar"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "effective_labels.%", "2"),
				),
			},
		},
	})
}

func TestAccComputeDisk_attributionLabelOnCreationSkip(t *testing.T) {
	// VCR tests cache provider configuration between steps, this test changes provider configuration and fails under VCR.
	acctest.SkipIfVcr(t)
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_attributionLabel(diskName, "false", "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.user-label", "foo"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.user-label", "foo"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "effective_labels.%", "1"),
				),
			},
			{
				Config: testAccComputeDisk_attributionLabelUpdated(diskName, "true", "CREATION_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.user-label", "bar"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.user-label", "bar"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "effective_labels.%", "1"),
				),
			},
		},
	})
}

func TestAccComputeDisk_attributionLabelProactive(t *testing.T) {
	// VCR tests cache provider configuration between steps, this test changes provider configuration and fails under VCR.
	acctest.SkipIfVcr(t)
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_attributionLabel(diskName, "false", "PROACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.user-label", "foo"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.user-label", "foo"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "effective_labels.%", "1"),
				),
			},
			{
				Config: testAccComputeDisk_attributionLabelUpdated(diskName, "true", "PROACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "labels.user-label", "bar"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.%", "2"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "terraform_labels.user-label", "bar"),

					resource.TestCheckResourceAttr("google_compute_disk.foobar", "effective_labels.%", "2"),
				),
			},
		},
	})
}

func testAccComputeDisk_attributionLabel(diskName, add, strategy string) string {
	return fmt.Sprintf(`
provider "google" {
	add_terraform_attribution_label               = %s
	terraform_attribution_label_addition_strategy = %q
}

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
    user-label = "foo"
  }
}
`, add, strategy, diskName)
}

func testAccComputeDisk_attributionLabelUpdated(diskName, add, strategy string) string {
	return fmt.Sprintf(`
provider "google" {
	add_terraform_attribution_label               = %s
	terraform_attribution_label_addition_strategy = %q
}

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
    user-label = "bar"
  }
}
`, add, strategy, diskName)
}

func TestAccComputeDisk_storagePoolSpecified(t *testing.T) {
	// Currently failing
	acctest.SkipIfVcr(t)
	t.Parallel()

	storagePoolName := fmt.Sprintf("tf-test-storage-pool-%s", acctest.RandString(t, 10))
	storagePoolUrl := fmt.Sprintf("/projects/%s/zones/%s/storagePools/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), storagePoolName)
	diskName := fmt.Sprintf("tf-test-disk-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: setupTestingStoragePool(t, storagePoolName),
				Config:    testAccComputeDisk_storagePoolSpecified(diskName, storagePoolUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "storage_pool", storagePoolName),
				),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	cleanupTestingStoragePool(t, storagePoolName)
}

func setupTestingStoragePool(t *testing.T, storagePoolName string) func() {
	return func() {
		config := acctest.GoogleProviderConfig(t)
		headers := make(http.Header)
		project := envvar.GetTestProjectFromEnv()
		zone := envvar.GetTestZoneFromEnv()
		url := fmt.Sprintf("%sprojects/%s/zones/%s/storagePools", config.ComputeBasePath, project, zone)
		storagePoolTypeUrl := fmt.Sprintf("/projects/%s/zones/%s/storagePoolTypes/hyperdisk-throughput", project, zone)
		defaultTimeout := 20 * time.Minute
		obj := make(map[string]interface{})
		obj["name"] = storagePoolName
		obj["poolProvisionedCapacityGb"] = 10240
		obj["poolProvisionedThroughput"] = 180
		obj["storagePoolType"] = storagePoolTypeUrl
		obj["capacityProvisioningType"] = "ADVANCED"

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      obj,
			Timeout:   defaultTimeout,
			Headers:   headers,
		})
		if err != nil {
			t.Errorf("Error creating StoragePool: %s", err)
		}

		err = tpgcompute.ComputeOperationWaitTime(config, res, project, "Creating StoragePool", config.UserAgent, defaultTimeout)
		if err != nil {
			t.Errorf("Error waiting to create StoragePool: %s", err)
		}
	}
}

func cleanupTestingStoragePool(t *testing.T, storagePoolName string) {
	config := acctest.GoogleProviderConfig(t)
	headers := make(http.Header)
	project := envvar.GetTestProjectFromEnv()
	zone := envvar.GetTestZoneFromEnv()
	url := fmt.Sprintf("%sprojects/%s/zones/%s/storagePools/%s", config.ComputeBasePath, project, zone, storagePoolName)
	defaultTimeout := 20 * time.Minute
	var obj map[string]interface{}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   project,
		RawURL:    url,
		UserAgent: config.UserAgent,
		Body:      obj,
		Timeout:   defaultTimeout,
		Headers:   headers,
	})
	if err != nil {
		t.Errorf("Error deleting StoragePool: %s", err)
	}

	err = tpgcompute.ComputeOperationWaitTime(config, res, project, "Deleting StoragePool", config.UserAgent, defaultTimeout)
	if err != nil {
		t.Errorf("Error waiting to delete StoragePool: %s", err)
	}
}

func testAccComputeDisk_storagePoolSpecified(diskName, storagePoolUrl string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
  name = "%s"
  type = "hyperdisk-throughput"
  size = 2048
  provisioned_throughput = 140
  storage_pool = "%s"
}
`, diskName, storagePoolUrl)
}

func TestAccComputeDisk_accessModeSpecified(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-disk-accessmode-%s", acctest.RandString(t, 10))
	accessModeForCreate := "READ_WRITE_SINGLE"
	accessModeForUpdate := "READ_ONLY_MANY"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create disk with Access Mode
			{
				Config: testAccComputeDisk_accessModeSpecified(diskName, accessModeForCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "access_mode", accessModeForCreate),
				),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update Access Mode
			{
				Config: testAccComputeDisk_accessModeSpecified(diskName, accessModeForUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "access_mode", accessModeForUpdate),
				),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeDisk_accessModeSpecified(diskName, accessMode string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
  name = "%s"
  type = "hyperdisk-ml"
  zone  = "us-central1-a"
  access_mode = "%s"
}
`, diskName, accessMode)
}
