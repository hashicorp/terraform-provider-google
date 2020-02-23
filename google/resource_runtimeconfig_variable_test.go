package google

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/runtimeconfig/v1beta1"
)

func TestAccRuntimeconfigVariable_basic(t *testing.T) {
	t.Parallel()

	var variable runtimeconfig.Variable

	varName := fmt.Sprintf("variable-test-%s", acctest.RandString(10))
	varText := "this is my test value"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuntimeconfigVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuntimeconfigVariable_basicText(varName, varText),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeconfigVariableExists(
						"google_runtimeconfig_variable.foobar", &variable),
					testAccCheckRuntimeconfigVariableText(&variable, varText),
					testAccCheckRuntimeconfigVariableUpdateTime("google_runtimeconfig_variable.foobar"),
				),
			},
			{
				ResourceName:      "google_runtimeconfig_variable.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRuntimeconfigVariable_basicUpdate(t *testing.T) {
	t.Parallel()

	var variable runtimeconfig.Variable

	configName := fmt.Sprintf("some-name-%s", acctest.RandString(10))
	varName := fmt.Sprintf("variable-test-%s", acctest.RandString(10))
	varText := "this is my test value"
	varText2 := "this is my updated value"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuntimeconfigVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuntimeconfigVariable_basicTextUpdate(configName, varName, varText),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeconfigVariableExists(
						"google_runtimeconfig_variable.foobar", &variable),
					testAccCheckRuntimeconfigVariableText(&variable, varText),
				),
			}, {
				Config: testAccRuntimeconfigVariable_basicTextUpdate(configName, varName, varText2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeconfigVariableExists(
						"google_runtimeconfig_variable.foobar", &variable),
					testAccCheckRuntimeconfigVariableText(&variable, varText2),
				),
			},
		},
	})
}

func TestAccRuntimeconfigVariable_basicValue(t *testing.T) {
	t.Parallel()

	var variable runtimeconfig.Variable

	varName := fmt.Sprintf("variable-test-%s", acctest.RandString(10))
	varValue := "Zm9vYmFyCg=="

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuntimeconfigVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRuntimeconfigVariable_basicValue(varName, varValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuntimeconfigVariableExists(
						"google_runtimeconfig_variable.foobar", &variable),
					testAccCheckRuntimeconfigVariableValue(&variable, varValue),
					testAccCheckRuntimeconfigVariableUpdateTime("google_runtimeconfig_variable.foobar"),
				),
			},
			{
				ResourceName:      "google_runtimeconfig_variable.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRuntimeconfigVariable_errorsOnBothValueAndText(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccRuntimeconfigVariable_invalidBothTextValue(),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

func TestAccRuntimeconfigVariable_errorsOnMissingValueAndText(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccRuntimeconfigVariable_invalidMissingTextValue(),
				ExpectError: regexp.MustCompile("You must specify one of value or text"),
			},
		},
	})
}

func testAccCheckRuntimeconfigVariableExists(resourceName string, variable *runtimeconfig.Variable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientRuntimeconfig.Projects.Configs.Variables.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		*variable = *found

		return nil
	}
}

func testAccCheckRuntimeconfigVariableUpdateTime(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		updateTime := rs.Primary.Attributes["update_time"]
		if updateTime == "" {
			return fmt.Errorf("No update time set for resource %s", resourceName)
		}

		// Make sure it's a valid rfc 3339 date
		_, err := time.Parse(time.RFC3339, updateTime)
		if err != nil {
			return fmt.Errorf("Error while parsing update time for resource %s: %s", resourceName, err.Error())
		}

		return nil
	}
}

func testAccCheckRuntimeconfigVariableText(variable *runtimeconfig.Variable, text string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if variable.Text != text {
			return fmt.Errorf("Variable %s had incorrect text: expected '%s' but found '%s'", variable.Name,
				text, variable.Text)
		}

		return nil
	}
}

func testAccCheckRuntimeconfigVariableValue(variable *runtimeconfig.Variable, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if variable.Value != value {
			return fmt.Errorf("Variable %s had incorrect value: expected '%s' but found '%s'", variable.Name,
				value, variable.Value)
		}

		return nil
	}
}

func testAccCheckRuntimeconfigVariableDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_runtimeconfig_variable" {
			continue
		}

		_, err := config.clientRuntimeconfig.Projects.Configs.Variables.Get(rs.Primary.ID).Do()

		if err == nil {
			return fmt.Errorf("Runtimeconfig variable still exists")
		}
	}

	return nil
}

func testAccRuntimeconfigVariable_basicText(name, text string) string {
	return fmt.Sprintf(`
resource "google_runtimeconfig_config" "foobar" {
  name = "some-config-%s"
}

resource "google_runtimeconfig_variable" "foobar" {
  parent = google_runtimeconfig_config.foobar.name
  name   = "%s"
  text   = "%s"
}
`, acctest.RandString(10), name, text)
}

func testAccRuntimeconfigVariable_basicTextUpdate(configName, name, text string) string {
	return fmt.Sprintf(`
resource "google_runtimeconfig_config" "foobar" {
  name = "%s"
}

resource "google_runtimeconfig_variable" "foobar" {
  parent = google_runtimeconfig_config.foobar.name
  name   = "%s"
  text   = "%s"
}
`, configName, name, text)
}

func testAccRuntimeconfigVariable_basicValue(name, value string) string {
	return fmt.Sprintf(`
resource "google_runtimeconfig_config" "foobar" {
  name = "some-config-%s"
}

resource "google_runtimeconfig_variable" "foobar" {
  parent = google_runtimeconfig_config.foobar.name
  name   = "%s"
  value  = "%s"
}
`, acctest.RandString(10), name, value)
}

func testAccRuntimeconfigVariable_invalidBothTextValue() string {
	return fmt.Sprintf(`
resource "google_runtimeconfig_config" "foobar" {
  name = "some-config-%s"
}

resource "google_runtimeconfig_variable" "foobar" {
  parent = google_runtimeconfig_config.foobar.name
  name   = "%s"
  text   = "here's my value"
  value  = "Zm9vYmFyCg=="
}
`, acctest.RandString(10), acctest.RandString(10))
}

func testAccRuntimeconfigVariable_invalidMissingTextValue() string {
	return fmt.Sprintf(`
resource "google_runtimeconfig_config" "foobar" {
  name = "some-config-%s"
}

resource "google_runtimeconfig_variable" "foobar" {
  parent = google_runtimeconfig_config.foobar.name
  name   = "my-variable-namespace/%s"
}
`, acctest.RandString(10), acctest.RandString(10))
}
