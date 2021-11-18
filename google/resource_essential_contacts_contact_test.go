package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEssentialContactsContact_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssentialContactsContactDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEssentialContactsContact_v1(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_essential_contacts_contact.contact",
						"email", "foo_v1@bar.com"),
				),
			},
			{
				ResourceName:            "google_essential_contacts_contact.contact",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccEssentialContactsContact_v2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_essential_contacts_contact.contact",
						"email", "foo_v2@bar.com"),
				),
			},
			{
				ResourceName:            "google_essential_contacts_contact.contact",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccEssentialContactsContact_v1(context map[string]interface{}) string {
	return Nprintf(`
data "google_project" "project" {
}

resource "google_essential_contacts_contact" "contact" {
  parent = data.google_project.project.id
  email = "foo_v1@bar.com"
  language_tag = "en-GB"
  notification_category_subscriptions = ["ALL"]
}
`, context)
}

func testAccEssentialContactsContact_v2(context map[string]interface{}) string {
	return Nprintf(`
data "google_project" "project" {
}

resource "google_essential_contacts_contact" "contact" {
  parent = data.google_project.project.id
  email = "foo_v2@bar.com"
  language_tag = "en-GB"
  notification_category_subscriptions = ["ALL"]
}
`, context)
}
