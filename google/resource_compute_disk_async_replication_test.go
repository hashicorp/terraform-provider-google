package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeDiskAsyncReplication(t *testing.T) {
	t.Parallel()

	region := GetTestRegionFromEnv()
	if !stringInSlice([]string{"europe-west2", "europe-west1", "us-central1", "us-east1", "us-west1", "us-east4", "asia-east1", "australia-southeast1"}, region) {
		return
	}
	secondaryRegion := region
	switch region {
	case "europe-west2":
		secondaryRegion = "europe-west1"
	case "europe-west1":
		secondaryRegion = "europe-west2"
	case "us-central1":
		secondaryRegion = "us-east1"
	case "us-east1", "us-west1", "us-east4":
		secondaryRegion = "us-central1"
	case "asia-east1":
		secondaryRegion = "asia-southeast1"
	case "asia-southeast1":
		secondaryRegion = "asia-east1"
	case "australia-southeast1":
		secondaryRegion = "australia-southeast2"
	case "australia-southeast2":
		secondaryRegion = "australia-southeast1"
	}

	primaryDisk := fmt.Sprintf("tf-test-disk-primary-%s", RandString(t, 10))
	secondaryDisk := fmt.Sprintf("tf-test-disk-secondary-%s", RandString(t, 10))
	primaryRegionalDisk := fmt.Sprintf("tf-test-disk-rprimary-%s", RandString(t, 10))
	secondaryRegionalDisk := fmt.Sprintf("tf-test-disk-rsecondary-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderBetaFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDiskAsyncReplication_basicZonal(region, secondaryRegion, primaryDisk, secondaryDisk),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk_async_replication.replication", "secondary_disk.0.state", "ACTIVE"),
				),
			},
			{
				ResourceName:      "google_compute_disk_async_replication.replication",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDiskAsyncReplication_basicRegional(region, secondaryRegion, primaryRegionalDisk, secondaryRegionalDisk),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_disk_async_replication.replication", "secondary_disk.0.state", "ACTIVE"),
				),
			},
			{
				ResourceName:      "google_compute_disk_async_replication.replication",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeDiskAsyncReplication_basicZonal(region, secondaryRegion, primaryDisk, secondaryDisk string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "primary" {
	provider = google-beta

	zone = "%s-a"
	name = "%s"
	type = "pd-ssd"
	
	physical_block_size_bytes = 4096
}
	
resource "google_compute_disk" "secondary" {
	provider = google-beta
	
	name = "%s"
	type = "pd-ssd"
	zone = "%s-b"
	
	async_primary_disk {
	    disk = google_compute_disk.primary.id
	}
	
	physical_block_size_bytes = 4096
}
	
resource "google_compute_disk_async_replication" "replication" {
	provider = google-beta
	
	primary_disk = google_compute_disk.primary.id

	secondary_disk {
		disk  = google_compute_disk.secondary.id
	}
}	  
`, region, primaryDisk, secondaryDisk, secondaryRegion)
}

func testAccComputeDiskAsyncReplication_basicRegional(region, secondaryRegion, primaryDisk, secondaryDisk string) string {
	return fmt.Sprintf(`
resource "google_compute_region_disk" "primary" {
	provider = google-beta
	
	region = "%s"
	name   = "%s"
	type   = "pd-ssd"
	
	physical_block_size_bytes = 4096

	replica_zones = [
		"%s-a",
		"%s-b"
	]
}
	
resource "google_compute_region_disk" "secondary" {
    provider = google-beta
	
	region = "%s"
	name   = "%s"
	type   = "pd-ssd"
	
	async_primary_disk {
	    disk = google_compute_region_disk.primary.id
	}
	
	physical_block_size_bytes = 4096

	replica_zones = [
		"%s-b",
		"%s-c"
	]
}
	
resource "google_compute_disk_async_replication" "replication" {
	provider = google-beta
	
	primary_disk = google_compute_region_disk.primary.id

	secondary_disk {
		disk  = google_compute_region_disk.secondary.id
	}
}	  
`, region, primaryDisk, region, region, secondaryRegion, secondaryDisk, secondaryRegion, secondaryRegion)
}
