package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLoggingView_basic(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	bucketId := fmt.Sprintf("projects/%s/locations/global/buckets/_Default", projectId)
	viewId := "tf-test-view-" + randString(t, 10)
	notTestProjectFilter := fmt.Sprintf("NOT source(projects/%s)", getTestProjectFromEnv())

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingView_basic(bucketId, viewId, "All logs", ""),
			},
			{
				ResourceName:      "google_logging_view.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingView_basic(bucketId, viewId, "All logs (except test project)", notTestProjectFilter),
			},
			{
				ResourceName:      "google_logging_view.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingView_basic(bucketId, viewId, "All logs", ""),
			},
			{
				ResourceName:      "google_logging_view.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingView_billingAccount(t *testing.T) {
	t.Parallel()

	billingAccountId := getTestBillingAccountFromEnv(t)
	bucketId := fmt.Sprintf("billingAccounts/%s/locations/global/buckets/_Default", billingAccountId)
	viewId := "tf-test-view-" + randString(t, 10)
	notTestProjectFilter := fmt.Sprintf("NOT source(projects/%s)", getTestProjectFromEnv())

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingView_basic(bucketId, viewId, "All logs", ""),
			},
			{
				ResourceName:      "google_logging_view.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingView_basic(bucketId, viewId, "All logs (except test project)", notTestProjectFilter),
			},
			{
				ResourceName:      "google_logging_view.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingView_organization(t *testing.T) {
	t.Parallel()

	organizationId := getTestOrgFromEnv(t)
	bucketId := fmt.Sprintf("organizations/%s/locations/global/buckets/_Default", organizationId)
	viewId := "tf-test-view-" + randString(t, 10)
	notTestProjectFilter := fmt.Sprintf("NOT source(projects/%s)", getTestProjectFromEnv())

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingView_basic(bucketId, viewId, "All logs", ""),
			},
			{
				ResourceName:      "google_logging_view.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingView_basic(bucketId, viewId, "All logs (except test project)", notTestProjectFilter),
			},
			{
				ResourceName:      "google_logging_view.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingView_folder(t *testing.T) {
	t.Parallel()

	organizationId := getTestOrgFromEnv(t)
	folderName := "tf-test-folder-" + randString(t, 10)
	viewId := "tf-test-view-" + randString(t, 10)
	notTestProjectFilter := fmt.Sprintf("NOT source(projects/%s)", getTestProjectFromEnv())

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingView_folder(organizationId, folderName, viewId, "All logs", ""),
			},
			{
				ResourceName:      "google_logging_view.folder",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingView_folder(organizationId, folderName, viewId, "All logs (except test project)", notTestProjectFilter),
			},
			{
				ResourceName:      "google_logging_view.folder",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingView_customBucket(t *testing.T) {
	t.Parallel()

	projectId := getTestProjectFromEnv()
	bucketName := "tf-test-bucket-" + randString(t, 10)
	viewId := "tf-test-view-" + randString(t, 10)
	notTestProjectFilter := fmt.Sprintf("NOT source(projects/%s)", getTestProjectFromEnv())

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingView_customBucket(projectId, bucketName, viewId, "All logs", ""),
			},
			{
				ResourceName:      "google_logging_view.custom_bucket_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingView_customBucket(projectId, bucketName, viewId, "All logs (except test project)", notTestProjectFilter),
			},
			{
				ResourceName:      "google_logging_view.custom_bucket_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLogViewDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_view" {
				continue
			}

			attributes := rs.Primary.Attributes

			_, err := config.NewLoggingClient(config.userAgent).Locations.Buckets.Views.Get(attributes["id"]).Do()
			if err == nil {
				return fmt.Errorf("Log View still exists")
			}
		}

		return nil
	}
}

func testAccLoggingView_basic(bucketId, viewId, description, filter string) string {
	return fmt.Sprintf(`
resource "google_logging_view" "basic" {
  bucket      = "%s"
  view_id     = "%s"
  description = "%s"
  filter      = "%s"
}
`, bucketId, viewId, description, filter)
}

func testAccLoggingView_folder(organizationId, folderName, viewId, description, filter string) string {
	return fmt.Sprintf(`
resource "google_folder" "test_folder" {
  parent = "organizations/%s"
  display_name = "%s"
}

resource "google_logging_view" "folder" {
  view_id     = "%s"
  bucket      = "${google_folder.test_folder.name}/locations/global/buckets/_Default"
  description = "%s"
  filter      = "%s"
}
`, organizationId, folderName, viewId, description, filter)
}

func testAccLoggingView_customBucket(projectId, bucketName, viewId, description, filter string) string {
	return fmt.Sprintf(`
resource "google_logging_project_bucket_config" "custom_bucket" {
  project        = "%s"
  bucket_id      = "%s"
  location       = "us-west1"
  retention_days = 30
  description    = "Log View test"
}

resource "google_logging_view" "custom_bucket_view" {
  view_id     = "%s"
  bucket      = google_logging_project_bucket_config.custom_bucket.id
  description = "%s"
  filter      = "%s"
}
`, projectId, bucketName, viewId, description, filter)
}
