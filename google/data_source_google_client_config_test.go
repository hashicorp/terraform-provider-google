package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleClientConfig_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_client_config.current"

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "zone"),
					resource.TestCheckResourceAttrSet(resourceName, "access_token"),
				),
			},
		},
	})
}

const testAccCheckGoogleClientConfig_basic = `
data "google_client_config" "current" { }
`
