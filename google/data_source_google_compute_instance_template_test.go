package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccInstanceTemplateDatasource_name(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceTemplate_name(getTestProjectFromEnv(), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_instance_template.default",
						"google_compute_instance_template.default",
						map[string]struct{}{},
					),
				),
			},
		},
	})
}

func TestAccInstanceTemplateDatasource_filter(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceTemplate_filter(getTestProjectFromEnv(), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_instance_template.default",
						"google_compute_instance_template.c",
						map[string]struct{}{},
					),
				),
			},
		},
	})
}

func TestAccInstanceTemplateDatasource_filter_mostRecent(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceTemplate_filter_mostRecent(getTestProjectFromEnv(), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_instance_template.default",
						"google_compute_instance_template.c",
						map[string]struct{}{},
					),
				),
			},
		},
	})
}

func testAccInstanceTemplate_name(project, suffix string) string {
	return Nprintf(`
resource "google_compute_instance_template" "default" {
  name        = "tf-test-template-%{suffix}"
  description = "Example template."

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

data "google_compute_instance_template" "default" {
  project = "%{project}"
  name    = google_compute_instance_template.default.name
}
`, map[string]interface{}{"project": project, "suffix": suffix})
}

func testAccInstanceTemplate_filter(project, suffix string) string {
	return Nprintf(`
resource "google_compute_instance_template" "a" {
  name        = "tf-test-template-a-%{suffix}"
  description = "Example template."

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
resource "google_compute_instance_template" "b" {
  name        = "tf-test-template-b-%{suffix}"
  description = "Example template."

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
resource "google_compute_instance_template" "c" {
  name        = "tf-test-template-c-%{suffix}"
  description = "Example template."

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

data "google_compute_instance_template" "default" {
  // Hack to prevent depends_on bug triggering datasource recreate due to https://github.com/hashicorp/terraform/issues/11806
  // This bug is fixed in 0.13+.
  project = "%{project}${replace(google_compute_instance_template.a.id, "/.*/", "")}${replace(google_compute_instance_template.b.id, "/.*/", "")}${replace(google_compute_instance_template.c.id, "/.*/", "")}"
  filter  = "name = tf-test-template-c-%{suffix}"
}
`, map[string]interface{}{"project": project, "suffix": suffix})
}

func testAccInstanceTemplate_filter_mostRecent(project, suffix string) string {
	return Nprintf(`
resource "google_compute_instance_template" "a" {
  name        = "tf-test-template-%{suffix}-a"
  description = "tf-test-instance-template"

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
resource "google_compute_instance_template" "b" {
  name        = "tf-test-template-%{suffix}-b"
  description = "tf-test-instance-template"

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
    google_compute_instance_template.a,
    google_compute_instance_template.c,
  ]
}
resource "google_compute_instance_template" "c" {
  name        = "tf-test-template-%{suffix}-c"
  description = "tf-test-instance-template"

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
    google_compute_instance_template.a,
  ]
}

data "google_compute_instance_template" "default" {
  // Hack to prevent depends_on bug triggering datasource recreate due to https://github.com/hashicorp/terraform/issues/11806
  // This bug is fixed in 0.13+.
  project = "%{project}${replace(google_compute_instance_template.b.id, "/.*/", "")}"
  filter      = "(name != tf-test-template-%{suffix}-b) (description = tf-test-instance-template)"
  most_recent = true
}
`, map[string]interface{}{"project": project, "suffix": suffix})
}
