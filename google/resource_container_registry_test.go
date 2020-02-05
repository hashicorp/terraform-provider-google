package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccContainerRegistry_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerRegistry_basic(),
			},
		},
	})
}

func TestAccContainerRegistry_iam(t *testing.T) {
	t.Parallel()
	account := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerRegistry_iam(account),
			},
		},
	})
}

func testAccContainerRegistry_basic() string {
	return `
resource "google_container_registry" "foobar" {
  location = "EU"
}
`
}

func testAccContainerRegistry_iam(account string) string {
	return fmt.Sprintf(`
resource "google_container_registry" "foobar" {
  location = "EU"
}

resource "google_service_account" "test-account-1" {
  account_id   = "acct-%s-1"
  display_name = "Container Registry Iam Testing Account"
}

resource "google_storage_bucket_iam_member" "viewer" {
  bucket = google_container_registry.foobar.id
  role = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.test-account-1.email}"
}
`, account)
}
