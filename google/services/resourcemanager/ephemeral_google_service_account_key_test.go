package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestEphemeralServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "basic", serviceAccount)
	keyName := fmt.Sprintf("projects/-/serviceAccounts/%s/keys/123", targetServiceAccountEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_basic(keyName),
			},
		},
	})
}

func TestEphemeralServiceAccountKey_withPublicKeyType(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "pktype", serviceAccount)
	keyName := fmt.Sprintf("projects/-/serviceAccounts/%s/keys/123", targetServiceAccountEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_withPublicKeyType(keyName),
			},
		},
	})
}

func TestEphemeralServiceAccountKey_withProject(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "project", serviceAccount)
	keyName := fmt.Sprintf("projects/-/serviceAccounts/%s/keys/123", targetServiceAccountEmail)
	project := envvar.GetTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_withProject(keyName, project),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_basic(keyName string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_key" "key" {
  name = "%s"
}
`, keyName)
}

func testAccEphemeralServiceAccountKey_withPublicKeyType(keyName string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_key" "key" {
  name            = "%s"
  public_key_type = "TYPE_RAW"
}
`, keyName)
}

func testAccEphemeralServiceAccountKey_withProject(keyName, project string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_key" "key" {
  name    = "%s"
  project = "%s"
}
`, keyName, project)
}
