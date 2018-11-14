package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleSelfLink_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleSelfLink_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_self_link.instance", "relative_uri", "projects/my-gcp-project/regions/us-central1/instances/my-instance"),
					resource.TestCheckResourceAttr("data.google_self_link.instance", "name", "my-instance"),
				),
			},
		},
	})
}

var testAccGoogleSelfLink_basic = `
data "google_self_link" "instance" {
  self_link = "https://www.googleapis.com/compute/v1/projects/my-gcp-project/regions/us-central1/instances/my-instance"
}
`
