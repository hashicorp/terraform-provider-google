package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/runtimeconfig/v1beta1"
)

func TestAccRuntimeconfigConfig_basic(t *testing.T) {
	t.Parallel()

	var runtimeConfig runtimeconfig.RuntimeConfig
	configName := fmt.Sprintf("runtimeconfig-test-%s", acctest.RandString(10))
	description := "my test description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuntimeconfigConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuntimeconfigConfig_basicDescription(configName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeConfigExists(
						"google_runtimeconfig_config.foobar", &runtimeConfig),
					testAccCheckRuntimeConfigDescription(&runtimeConfig, description),
				),
			},
			{
				ResourceName:      "google_runtimeconfig_config.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRuntimeconfig_update(t *testing.T) {
	t.Parallel()

	var runtimeConfig runtimeconfig.RuntimeConfig
	configName := fmt.Sprintf("runtimeconfig-test-%s", acctest.RandString(10))
	firstDescription := "my test description"
	secondDescription := "my updated test description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuntimeconfigConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuntimeconfigConfig_basicDescription(configName, firstDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeConfigExists(
						"google_runtimeconfig_config.foobar", &runtimeConfig),
					testAccCheckRuntimeConfigDescription(&runtimeConfig, firstDescription),
				),
			}, {
				Config: testAccRuntimeconfigConfig_basicDescription(configName, secondDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeConfigExists(
						"google_runtimeconfig_config.foobar", &runtimeConfig),
					testAccCheckRuntimeConfigDescription(&runtimeConfig, secondDescription),
				),
			},
		},
	})
}

func TestAccRuntimeconfig_updateEmptyDescription(t *testing.T) {
	t.Parallel()

	var runtimeConfig runtimeconfig.RuntimeConfig
	configName := fmt.Sprintf("runtimeconfig-test-%s", acctest.RandString(10))
	description := "my test description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuntimeconfigConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuntimeconfigConfig_basicDescription(configName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeConfigExists(
						"google_runtimeconfig_config.foobar", &runtimeConfig),
					testAccCheckRuntimeConfigDescription(&runtimeConfig, description),
				),
			}, {
				Config: testAccRuntimeconfigConfig_emptyDescription(configName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeConfigExists(
						"google_runtimeconfig_config.foobar", &runtimeConfig),
					testAccCheckRuntimeConfigDescription(&runtimeConfig, ""),
				),
			},
		},
	})
}

func testAccCheckRuntimeConfigDescription(runtimeConfig *runtimeconfig.RuntimeConfig, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if runtimeConfig.Description != description {
			return fmt.Errorf("On runtime config '%s', expected description '%s', but found '%s'",
				runtimeConfig.Name, description, runtimeConfig.Description)
		}
		return nil
	}
}

func testAccCheckRuntimeConfigExists(resourceName string, runtimeConfig *runtimeconfig.RuntimeConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientRuntimeconfig.Projects.Configs.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		*runtimeConfig = *found

		return nil
	}
}

func testAccCheckRuntimeconfigConfigDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_runtimeconfig_config" {
			continue
		}

		_, err := config.clientRuntimeconfig.Projects.Configs.Get(rs.Primary.ID).Do()

		if err == nil {
			return fmt.Errorf("Runtimeconfig still exists")
		}
	}

	return nil
}

func testAccRuntimeconfigConfig_basicDescription(name, description string) string {
	return fmt.Sprintf(`
resource "google_runtimeconfig_config" "foobar" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}

func testAccRuntimeconfigConfig_emptyDescription(name string) string {
	return fmt.Sprintf(`
resource "google_runtimeconfig_config" "foobar" {
  name = "%s"
}
`, name)
}
