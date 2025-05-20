// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
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

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@compute-system.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

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

func testAccCheckComputeDisk_removeBackupSnapshot(t *testing.T, parentDiskName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		snapshot, err := config.NewComputeClient(config.UserAgent).Snapshots.List(envvar.GetTestProjectFromEnv()).Filter(fmt.Sprintf("name eq %s.*", parentDiskName)).Do()
		if err != nil {
			return err
		}

		if len(snapshot.Items) == 0 {
			return fmt.Errorf("No snapshot found")
		}

		op, err := config.NewComputeClient(config.UserAgent).Snapshots.Delete(envvar.GetTestProjectFromEnv(), snapshot.Items[0].Name).Do()
		if err != nil {
			return err
		}
		return tpgcompute.ComputeOperationWaitTime(config, op, envvar.GetTestProjectFromEnv(), "Deleting Snapshot", config.UserAgent, 10*time.Minute)
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

func TestAccComputeDisk_architecture(t *testing.T) {
	t.Parallel()

	context_1 := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"architecture":  "X86_64",
	}
	context_2 := map[string]interface{}{
		"random_suffix": context_1["random_suffix"],
		"architecture":  "ARM64",
	}
	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_architecture(context_1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", envvar.GetTestProjectFromEnv(), &disk),
				),
			},
			{
				Config: testAccComputeDisk_architecture(context_2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", envvar.GetTestProjectFromEnv(), &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_sourceStorageObject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":         acctest.RandString(t, 10),
		"source_storage_object": "test-fixtures/empty-image.tar.gz",
	}

	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_sourceStorageObject(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", envvar.GetTestProjectFromEnv(), &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_resourceManagerTags(t *testing.T) {
	t.Parallel()
	pid := envvar.GetTestProjectFromEnv()
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_id":    pid,
	}

	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_resourceManagerTags(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", pid, &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_sourceInstantSnapshot(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	var disk compute.Disk

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_sourceInstantSnapshot(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.foobar", envvar.GetTestProjectFromEnv(), &disk),
				),
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

func TestAccComputeDisk_createSnapshotBeforeDestroy(t *testing.T) {
	acctest.SkipIfVcr(t) // Disk cleanup test check
	t.Parallel()

	var disk1 compute.Disk
	var disk2 compute.Disk
	var disk3 compute.Disk
	context := map[string]interface{}{
		"disk_name1":        fmt.Sprintf("tf-test-disk-%s", acctest.RandString(t, 10)),
		"disk_name2":        fmt.Sprintf("test-%s", acctest.RandString(t, 44)), //this is over the snapshot character creation limit of 48
		"disk_name3":        fmt.Sprintf("tf-test-disk-%s", acctest.RandString(t, 10)),
		"snapshot_prefix":   fmt.Sprintf("tf-test-snapshot-%s", acctest.RandString(t, 10)),
		"kms_key_self_link": acctest.BootstrapKMSKey(t).CryptoKey.Name,
		"raw_key":           "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=",
		"rsa_encrypted_key": "ieCx/NcW06PcT7Ep1X6LUTc/hLvUDYyzSZPPVCVPTVEohpeHASqC8uw5TzyO9U+Fka9JFHz0mBibXUInrC/jEk014kCK/NPjYgEMOyssZ4ZINPKxlUh2zn1bV+MCaTICrdmuSBTWlUUiFoDD6PYznLwh8ZNdaheCeZ8ewEXgFQ8V+sDroLaN3Xs3MDTXQEMMoNUXMCZEIpg9Vtp9x2oeQ5lAbtt7bYAAHf5l+gJWw3sUfs0/Glw5fpdjT8Uggrr+RMZezGrltJEF293rvTIjWOEB3z5OHyHwQkvdrPDFcTqsLfh+8Hr8g+mf+7zVPEC8nEbqpdl3GPv3A7AwpFp7MA==",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_createSnapshotBeforeDestroy_init(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.raw-encrypted-name", envvar.GetTestProjectFromEnv(), &disk1),
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.rsa-encrypted-prefix", envvar.GetTestProjectFromEnv(), &disk2),
					testAccCheckComputeDiskExists(
						t, "google_compute_disk.kms-encrypted-name", envvar.GetTestProjectFromEnv(), &disk3),
				),
			},
			{
				Config:  testAccComputeDisk_createSnapshotBeforeDestroy_init(context),
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDisk_removeBackupSnapshot(t, context["disk_name1"].(string)),
					testAccCheckComputeDisk_removeBackupSnapshot(t, context["snapshot_prefix"].(string)),
					testAccCheckComputeDisk_removeBackupSnapshot(t, context["disk_name3"].(string)),
				),
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
	t.Parallel()

	storagePoolNameLong := acctest.BootstrapComputeStoragePool(t, "basic-1", "hyperdisk-throughput")
	diskName := fmt.Sprintf("tf-test-disk-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_storagePoolSpecified(diskName, storagePoolNameLong),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeDisk_storagePoolSpecified_nameOnly(t *testing.T) {
	t.Parallel()

	acctest.BootstrapComputeStoragePool(t, "basic-2", "hyperdisk-throughput")
	diskName := fmt.Sprintf("tf-test-disk-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_storagePoolSpecified(diskName, "tf-bootstrap-storage-pool-hyperdisk-throughput-basic-2"),
			},
			{
				ResourceName:      "google_compute_disk.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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

func TestExpandStoragePoolUrl_withDataProjectAndZone(t *testing.T) {
	config := &transport_tpg.Config{
		ComputeBasePath: "https://www.googleapis.com/compute/v1/",
		Project:         "other-project",
		Zone:            "other-zone",
	}

	data := &tpgresource.ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project": "test-project",
			"zone":    "test-zone",
		},
	}

	name := "test-storage-pool"
	zoneUrl := "zones/test-zone/storagePools/" + name
	projectUrl := "projects/test-project/" + zoneUrl
	fullUrl := config.ComputeBasePath + projectUrl

	cases := []struct {
		name     string
		inputStr string
	}{
		{
			name:     "full url",
			inputStr: fullUrl,
		},
		{
			name:     "project/{project}/zones/{zone}/storagePools/{storagePool}",
			inputStr: projectUrl,
		},
		{
			name:     "/project/{project}/zones/{zone}/storagePools/{storagePool}",
			inputStr: "/" + projectUrl,
		},
		{
			name:     "zones/{zone}/storagePools/{storagePool}",
			inputStr: zoneUrl,
		},
		{
			name:     "/zones/{zone}/storagePools/{storagePool}",
			inputStr: "/" + zoneUrl,
		},
		{
			name:     "{storagePool}",
			inputStr: name,
		},
		{
			name:     "/{storagePool}",
			inputStr: "/" + name,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, _ := tpgcompute.ExpandStoragePoolUrl(tc.inputStr, data, config)
			if result != fullUrl {
				t.Fatalf("%s does not match with expected full url: %s", result, fullUrl)
			}
		})
	}
}

func TestExpandStoragePoolUrl_withConfigProjectAndZone(t *testing.T) {
	config := &transport_tpg.Config{
		ComputeBasePath: "https://www.googleapis.com/compute/v1/",
		Project:         "test-project",
		Zone:            "test-zone",
	}

	data := &tpgresource.ResourceDataMock{}

	name := "test-storage-pool"
	zoneUrl := "zones/test-zone/storagePools/" + name
	projectUrl := "projects/test-project/" + zoneUrl
	fullUrl := config.ComputeBasePath + projectUrl

	cases := []struct {
		name     string
		inputStr string
	}{
		{
			name:     "full url",
			inputStr: fullUrl,
		},
		{
			name:     "project/{project}/zones/{zone}/storagePools/{storagePool}",
			inputStr: projectUrl,
		},
		{
			name:     "/project/{project}/zones/{zone}/storagePools/{storagePool}",
			inputStr: "/" + projectUrl,
		},
		{
			name:     "zones/{zone}/storagePools/{storagePool}",
			inputStr: zoneUrl,
		},
		{
			name:     "/zones/{zone}/storagePools/{storagePool}",
			inputStr: "/" + zoneUrl,
		},
		{
			name:     "{storagePool}",
			inputStr: name,
		},
		{
			name:     "/{storagePool}",
			inputStr: "/" + name,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, _ := tpgcompute.ExpandStoragePoolUrl(tc.inputStr, data, config)
			if result != fullUrl {
				t.Fatalf("%s does not match with expected full url: %s", result, fullUrl)
			}
		})
	}
}

func TestExpandStoragePoolUrl_noProjectAndZoneFromConfigAndData(t *testing.T) {
	config := &transport_tpg.Config{
		ComputeBasePath: "https://www.googleapis.com/compute/v1/",
	}

	data := &tpgresource.ResourceDataMock{}

	name := "test-storage-pool"
	zoneUrl := "zones/test-zone/storagePools/" + name
	projectUrl := "projects/test-project/" + zoneUrl
	fullUrl := config.ComputeBasePath + projectUrl

	cases := []struct {
		name     string
		inputStr string
	}{
		{
			name:     "full url",
			inputStr: fullUrl,
		},
		{
			name:     "project/{project}/zones/{zone}/storagePools/{storagePool}",
			inputStr: projectUrl,
		},
		{
			name:     "/project/{project}/zones/{zone}/storagePools/{storagePool}",
			inputStr: "/" + projectUrl,
		},
		{
			name:     "zones/{zone}/storagePools/{storagePool}",
			inputStr: zoneUrl,
		},
		{
			name:     "/zones/{zone}/storagePools/{storagePool}",
			inputStr: "/" + zoneUrl,
		},
		{
			name:     "{storagePool}",
			inputStr: name,
		},
		{
			name:     "/{storagePool}",
			inputStr: "/" + name,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := tpgcompute.ExpandStoragePoolUrl(tc.inputStr, data, config)
			if err == nil {
				t.Fatal("Should return error when no project and zone available from config or resource data")
			}
		})
	}
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

func testAccComputeDisk_createSnapshotBeforeDestroy_init(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_disk" "raw-encrypted-name" {
  name = "%{disk_name1}"
  type = "pd-ssd"
  size = 10
  zone  = "us-central1-a"

  disk_encryption_key {
	raw_key = "%{raw_key}"
  }

  create_snapshot_before_destroy = true
}

resource "google_compute_disk" "rsa-encrypted-prefix" {
  name = "%{disk_name2}"
  type = "pd-ssd"
  size = 10
  zone  = "us-central1-a"

  disk_encryption_key {
	rsa_encrypted_key = "%{rsa_encrypted_key}"
  }

  create_snapshot_before_destroy = true
  create_snapshot_before_destroy_prefix = "%{snapshot_prefix}"
}

resource "google_compute_disk" "kms-encrypted-name" {
  name = "%{disk_name3}"
  type = "pd-ssd"
  size = 10
  zone  = "us-central1-a"

  disk_encryption_key {
	kms_key_self_link = "%{kms_key_self_link}"
  }

  create_snapshot_before_destroy = true
}`, context)
}

func testAccComputeDisk_architecture(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_disk" "foobar" {
	name = "tf-test-disk-%{random_suffix}"
	type = "pd-ssd"
	size = 10
	zone = "us-central1-a"
	architecture = "%{architecture}"
}
`, context)
}

func testAccComputeDisk_sourceInstantSnapshot(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_disk" "to-snapshot" {
	name = "tf-test-disk-1-%{random_suffix}"
	type = "pd-ssd"
	size = 10
	zone = "us-central1-a"
}

resource "google_compute_instant_snapshot" "test" {
	name = "tf-test-instant-snapshot-%{random_suffix}"
	zone = "us-central1-a"
	source_disk = google_compute_disk.to-snapshot.id
}

resource "google_compute_disk" "foobar" {
	name = "tf-test-disk-2-%{random_suffix}"
	type = "pd-ssd"
	size = 10
	zone = "us-central1-a"
	source_instant_snapshot = google_compute_instant_snapshot.test.id
}
`, context)
}

func testAccComputeDisk_sourceStorageObject(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
	name = "tf-test-bucket-%{random_suffix}"
	location = "US"
}

resource "google_storage_bucket_object" "object" {
	name = "tf-test-object-%{random_suffix}.tar.gz"
	bucket = google_storage_bucket.bucket.name
	source = "%{source_storage_object}"
}

resource "google_compute_disk" "foobar" {
	name = "tf-test-disk-%{random_suffix}"
	type = "pd-ssd"
	size = 10
	zone = "us-central1-a"
	source_storage_object = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"

	depends_on = [google_storage_bucket_object.object]
}
`, context)
}

func testAccComputeDisk_resourceManagerTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "tag_key" {
  parent = "projects/%{project_id}"
  short_name = "test-%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value" {
  parent = "tagKeys/${google_tags_tag_key.tag_key.name}"
  short_name = "name-%{random_suffix}"
}

resource "google_compute_disk" "foobar" {
  name = "tf-test-disk-%{random_suffix}"
  type = "pd-ssd"
  size = 10
  zone = "us-central1-a"
  params {
	resource_manager_tags = {
	  "${google_tags_tag_key.tag_key.id}" = "${google_tags_tag_value.tag_value.id}"
  	}
  }
}
`, context)
}
