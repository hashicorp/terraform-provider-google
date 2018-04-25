package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceDnsManagedZone_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsManagedZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceDnsManagedZone_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDnsManagedZoneCheck("data.google_dns_managed_zone.foo", "google_dns_managed_zone.foo"),
					testAccDataSourceDnsManagedZoneCheck("data.google_dns_managed_zone.bar", "google_dns_managed_zone.bar"),
				),
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

func testAccDataSourceDnsManagedZone_basic() string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foo" {
	name		= "foo-zone-%s"
	dns_name	= "foo.test.com."
	description	= "Foo DNS zone"
}

data "google_dns_managed_zone" "foo" {
	name	= "${google_dns_managed_zone.foo.name}"
}

resource "google_dns_managed_zone" "bar" {
	name		= "bar-zone-%s"
	dns_name	= "bar.test.com."
	description	= "Bar DNS zone"
}

data "google_dns_managed_zone" "bar" {
	dns_name = "${google_dns_managed_zone.bar.dns_name}"
}
`, acctest.RandString(10), acctest.RandString(10))
}
