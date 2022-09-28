package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSnapshotDatasource_name(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshot_name(getTestProjectFromEnv(), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_snapshot.default",
						"google_compute_snapshot.default",
						map[string]struct{}{"zone": {}},
					),
				),
			},
		},
	})
}

func TestAccSnapshotDatasource_filter(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshot_filter(getTestProjectFromEnv(), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_snapshot.default",
						"google_compute_snapshot.c",
						map[string]struct{}{"zone": {}},
					),
				),
			},
		},
	})
}

func TestAccSnapshotDatasource_filterMostRecent(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshot_filter_mostRecent(getTestProjectFromEnv(), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_snapshot.default",
						"google_compute_snapshot.c",
						map[string]struct{}{"zone": {}},
					),
				),
			},
		},
	})
}

func testAccSnapshot_name(project, suffix string) string {
	return Nprintf(`
	data "google_compute_image" "tf-test-image" {
		family  = "debian-11"
		project = "debian-cloud"
	}
	resource "google_compute_disk" "tf-test-disk" {
		name  = "debian-disk-%{suffix}"
		image = data.google_compute_image.tf-test-image.self_link
		size  = 10
		type  = "pd-ssd"
		zone  = "us-central1-a"
	  }

	resource "google_compute_snapshot" "default" {
		name = "tf-test-snapshot-%{suffix}"
		description = "Example snapshot."
		source_disk = google_compute_disk.tf-test-disk.id
		zone        = "us-central1-a"
		labels = {
			my_label = "value"
		}
		storage_locations = ["us-central1"]
	}
	data "google_compute_snapshot" "default" {
		project = "%{project}"
		name = google_compute_snapshot.default.name
	}

	`, map[string]interface{}{"project": project, "suffix": suffix})
}

func testAccSnapshot_filter(project, suffix string) string {
	return Nprintf(`
	data "google_compute_image" "tf-test-image" {
		family  = "debian-11"
		project = "debian-cloud"
	}
	resource "google_compute_disk" "tf-test-disk" {
		name  = "debian-disk-%{suffix}"
		image = data.google_compute_image.tf-test-image.self_link
		size  = 10
		type  = "pd-ssd"
		zone  = "us-central1-a"
	}
	resource "google_compute_snapshot" "a" {
		name = "tf-test-snapshot-a-%{suffix}"
		description = "Example snapshot."
		source_disk = google_compute_disk.tf-test-disk.id
		zone        = "us-central1-a"
		labels = {
			my_label = "a"
		}
		storage_locations = ["us-central1"]
	}
	resource "google_compute_snapshot" "b" {
		name = "tf-test-snapshot-b-%{suffix}"
		description = "Example snapshot."
		source_disk = google_compute_disk.tf-test-disk.id
		zone        = "us-central1-a"
		labels = {
			my_label = "b"
		}
		storage_locations = ["us-central1"]
	}
	resource "google_compute_snapshot" "c" {
		name = "tf-test-snapshot-c-%{suffix}"
		description = "Example snapshot."
		source_disk = google_compute_disk.tf-test-disk.id
		zone        = "us-central1-a"
		labels = {
			my_label = "c"
		}
		storage_locations = ["us-central1"]
	}
	data "google_compute_snapshot" "default" {
		project = "%{project}"
		filter  = "name = tf-test-snapshot-c-%{suffix}"
		depends_on = [google_compute_snapshot.c]
	}
`, map[string]interface{}{"project": project, "suffix": suffix})
}

func testAccSnapshot_filter_mostRecent(project, suffix string) string {
	return Nprintf(`
	data "google_compute_image" "tf-test-image" {
		family  = "debian-11"
		project = "debian-cloud"
	}
	resource "google_compute_disk" "tf-test-disk" {
		name  = "debian-disk-%{suffix}"
		image = data.google_compute_image.tf-test-image.self_link
		size  = 10
		type  = "pd-ssd"
		zone  = "us-central1-a"
	}
	resource "google_compute_snapshot" "a" {
		name = "tf-test-snapshot-a-%{suffix}"
		description = "Example snapshot."
		source_disk = google_compute_disk.tf-test-disk.id
		zone        = "us-central1-a"
		labels = {
			my_label = "a"
		}
		storage_locations = ["us-central1"]
	}
	resource "google_compute_snapshot" "b" {
		name = "tf-test-snapshot-b-%{suffix}"
		description = "Example snapshot."
		source_disk = google_compute_disk.tf-test-disk.id
		zone        = "us-central1-a"
		labels = {
			my_label = "b"
		}
		storage_locations = ["us-central1"]
	}
	resource "google_compute_snapshot" "c" {
		name = "tf-test-snapshot-c-%{suffix}"
		description = "Example snapshot."
		source_disk = google_compute_disk.tf-test-disk.id
		zone        = "us-central1-a"
		labels = {
			my_label = "c"
		}
		storage_locations = ["us-central1"]
	}
	data "google_compute_snapshot" "default" {
		project = "%{project}"
		most_recent = true
		filter  = "name = tf-test-snapshot-c-%{suffix}"
		depends_on = [google_compute_snapshot.c]
	}
`, map[string]interface{}{"project": project, "suffix": suffix})
}
