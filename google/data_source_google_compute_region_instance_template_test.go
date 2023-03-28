package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRegionInstanceTemplateDatasource_name(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceTemplate_name(GetTestProjectFromEnv(), RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_region_instance_template.default",
						"google_compute_region_instance_template.default",
						map[string]struct{}{},
					),
				),
			},
		},
	})
}

func TestAccRegionInstanceTemplateDatasource_filter(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceTemplate_filter(GetTestProjectFromEnv(), RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_region_instance_template.default",
						"google_compute_region_instance_template.c",
						map[string]struct{}{},
					),
				),
			},
		},
	})
}

func TestAccRegionInstanceTemplateDatasource_filter_mostRecent(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceTemplate_filter_mostRecent(GetTestProjectFromEnv(), RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_region_instance_template.default",
						"google_compute_region_instance_template.c",
						map[string]struct{}{},
					),
				),
			},
		},
	})
}

func testAccRegionInstanceTemplate_name(project, suffix string) string {
	return Nprintf(`
resource "google_compute_region_instance_template" "default" {
  name        = "tf-test-template-%{suffix}"
  description = "Example template."
  region = "us-central1"

  machine_type = "e2-small"

  tags = ["foo", "bar"]

  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }
}

data "google_compute_region_instance_template" "default" {
  project = "%{project}"
  region = "us-central1"
  name    = google_compute_region_instance_template.default.name
}
`, map[string]interface{}{"project": project, "suffix": suffix})
}

func testAccRegionInstanceTemplate_filter(project, suffix string) string {
	return Nprintf(`
resource "google_compute_region_instance_template" "a" {
  name        = "tf-test-template-a-%{suffix}"
  description = "Example template."
  region = "us-central1"
  machine_type = "e2-small"
  tags = ["foo", "bar", "a"]
  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }
}

resource "google_compute_region_instance_template" "b" {
  name        = "tf-test-template-b-%{suffix}"
  description = "Example template."
  region = "us-central1"
  machine_type = "e2-small"
  tags = ["foo", "bar", "b"]
  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }
}

resource "google_compute_region_instance_template" "c" {
  name        = "tf-test-template-c-%{suffix}"
  description = "Example template."
  region = "us-central1"
  machine_type = "e2-small"
  tags = ["foo", "bar", "c"]
  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }
}

data "google_compute_region_instance_template" "default" {
  project = "%{project}"
  region = "us-central1"
  filter  = "name = tf-test-template-c-%{suffix}"
  depends_on = [
    google_compute_region_instance_template.a,
    google_compute_region_instance_template.b,
    google_compute_region_instance_template.c,
  ]
}
`, map[string]interface{}{"project": project, "suffix": suffix})
}

func testAccRegionInstanceTemplate_filter_mostRecent(project, suffix string) string {
	return Nprintf(`
resource "google_compute_region_instance_template" "a" {
  name        = "tf-test-template-%{suffix}-a"
  description = "tf-test-instance-template"
  region = "us-central1"

  machine_type = "e2-small"

  tags = ["foo", "bar", "a"]

  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }
}
resource "google_compute_region_instance_template" "b" {
  name        = "tf-test-template-%{suffix}-b"
  description = "tf-test-instance-template"
  region = "us-central1"

  machine_type = "e2-small"

  tags = ["foo", "bar", "b"]

  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  depends_on = [
    google_compute_region_instance_template.a,
    google_compute_region_instance_template.c,
  ]
}
resource "google_compute_region_instance_template" "c" {
  name        = "tf-test-template-%{suffix}-c"
  description = "tf-test-instance-template"
  region = "us-central1"

  machine_type = "e2-small"

  tags = ["foo", "bar", "c"]

  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  depends_on = [
    google_compute_region_instance_template.a,
  ]
}

data "google_compute_region_instance_template" "default" {
  region = "us-central1"
  filter      = "(name != tf-test-template-%{suffix}-b) (description = tf-test-instance-template)"
  most_recent = true
  depends_on = [
    google_compute_region_instance_template.a,
    google_compute_region_instance_template.b,
    google_compute_region_instance_template.c,
  ]
}
`, map[string]interface{}{"project": project, "suffix": suffix})
}
