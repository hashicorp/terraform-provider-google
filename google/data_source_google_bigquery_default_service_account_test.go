package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleBigqueryDefaultServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_bigquery_default_service_account.bq_account"

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleBigqueryDefaultServiceAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "member"),
				),
			},
		},
	})
}

const testAccCheckGoogleBigqueryDefaultServiceAccount_basic = `
data "google_bigquery_default_service_account" "bq_account" {
}
`
