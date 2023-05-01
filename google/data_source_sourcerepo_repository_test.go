package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleSourceRepoRepository_basic(t *testing.T) {
	t.Parallel()

	name := "tf-repository-" + RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSourceRepoRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleSourceRepoRepositoryConfig(name),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_sourcerepo_repository.bar", "google_sourcerepo_repository.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSourceRepoRepositoryConfig(name string) string {
	return fmt.Sprintf(`
resource "google_sourcerepo_repository" "foo" {
  name               = "%s"
}

data "google_sourcerepo_repository" "bar" {
  name = google_sourcerepo_repository.foo.name
  depends_on = [
    google_sourcerepo_repository.foo,
  ]
}
`, name)
}
