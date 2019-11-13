package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceComputeImage(t *testing.T) {
	t.Parallel()

	family := acctest.RandomWithPrefix("tf-test")
	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePublicImageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCheckPublicImage(),
				),
			},
			{
				Config: testAccDataSourceCustomImageConfig(family, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_image.from_name",
						"name", name),
					resource.TestCheckResourceAttr("data.google_compute_image.from_name",
						"family", family),
					resource.TestCheckResourceAttrSet("data.google_compute_image.from_name",
						"self_link"),
				),
			},
		},
	})
}

func testAccDataSourceCheckPublicImage() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		data_source_name := "data.google_compute_image.debian"
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		ds_attr := ds.Primary.Attributes
		attrs_to_test := map[string]string{
			"family": "debian-9",
		}

		for attr, expect_value := range attrs_to_test {
			if ds_attr[attr] != expect_value {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr,
					ds_attr[attr],
					expect_value,
				)
			}
		}

		selfLink := "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-9-stretch-v20171129"

		if !compareSelfLinkOrResourceName("", ds_attr["self_link"], selfLink, nil) && ds_attr["self_link"] != selfLink {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], selfLink)
		}

		return nil
	}
}

var testAccDataSourcePublicImageConfig = `
data "google_compute_image" "debian" {
  project = "debian-cloud"
  name    = "debian-9-stretch-v20171129"
}
`

func testAccDataSourceCustomImageConfig(family, name string) string {
	return fmt.Sprintf(`
resource "google_compute_image" "image" {
  family      = "%s"
  name        = "%s"
  source_disk = google_compute_disk.disk.self_link
}

resource "google_compute_disk" "disk" {
  name = "%s-disk"
  zone = "us-central1-b"
}

data "google_compute_image" "from_name" {
  project = google_compute_image.image.project
  name    = google_compute_image.image.name
}

data "google_compute_image" "from_family" {
  project = google_compute_image.image.project
  family  = google_compute_image.image.family
}
`, family, name, name)
}
