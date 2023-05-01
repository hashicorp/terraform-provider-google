package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSecurityCenterSource_basic(t *testing.T) {
	t.Parallel()

	orgId := acctest.GetTestOrgFromEnv(t)
	suffix := RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterSource_sccSourceBasicExample(orgId, suffix, "My description"),
			},
			{
				ResourceName:      "google_scc_source.custom_source",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterSource_sccSourceBasicExample(orgId, suffix, ""),
			},
			{
				ResourceName:      "google_scc_source.custom_source",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterSource_sccSourceBasicExample(orgId, suffix, description string) string {
	return fmt.Sprintf(`
resource "google_scc_source" "custom_source" {
  display_name = "TFSrc %s"
  organization = "%s"
  description  = "%s"
}
`, suffix, orgId, description)
}
