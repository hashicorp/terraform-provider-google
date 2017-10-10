package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDataSourceDnsManagedZone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsManagedZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceDnsManagedZone_basic,
				Check:  testAccDataSourceDnsManagedZoneCheck("data.google_dns_managed_zone.qa", "google_dns_managed_zone.foo"),
			},
		},
	})
}

func testAccDataSourceDnsManagedZoneCheck(dsName, rsName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[rsName]
		if !ok {
			return fmt.Errorf("can't find resource called %s in state", rsName)
		}

		rs, ok := s.RootModule().Resources[dsName]
		if !ok {
			return fmt.Errorf("can't find data source called %s in state", dsName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		attrsToTest := []string{
			"id",
			"name",
			"description",
			"dns_name",
			"name_servers",
		}

		for _, attrToTest := range attrsToTest {
			if dsAttr[attrToTest] != rsAttr[attrToTest] {
				return fmt.Errorf("%s is %s; want %s", attrToTest, dsAttr[attrToTest], rsAttr[attrToTest])
			}
		}

		return nil
	}
}

var testAccDataSourceDnsManagedZone_basic = fmt.Sprintf(`
resource "google_dns_managed_zone" "foo" {
	name		= "qa-zone-%s"
	dns_name	= "qa.test.com."
	description	= "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
	name	= "${google_dns_managed_zone.foo.name}"
}
`, acctest.RandString(10))
