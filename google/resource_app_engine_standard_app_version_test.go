package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAppEngineStandardAppVersion_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
		"random_suffix":   randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppEngineStandardAppVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineStandardAppVersion_python(context),
			},
			{
				ResourceName:            "google_app_engine_standard_app_version.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_variables", "deployment", "entrypoint", "service", "noop_on_destroy"},
			},
			{
				Config: testAccAppEngineStandardAppVersion_pythonUpdate(context),
			},
			{
				ResourceName:            "google_app_engine_standard_app_version.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_variables", "deployment", "entrypoint", "service", "noop_on_destroy"},
			},
		},
	})
}

func testAccAppEngineStandardAppVersion_python(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "my_project" {
  name = "tf-test-appeng-std%{random_suffix}"
  project_id = "tf-test-appeng-std%{random_suffix}"
  org_id = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_app_engine_application" "app" {
  project     = google_project.my_project.project_id
  location_id = "us-central"
}

resource "google_project_service" "project" {
  project = google_project.my_project.project_id
  service = "appengine.googleapis.com"

  disable_dependent_services = false
}

resource "google_app_engine_standard_app_version" "foo" {
  project    = google_project_service.project.project
  version_id = "v1"
  service    = "default"
  runtime    = "python38"

  entrypoint {
    shell = "gunicorn -b :$PORT main:app"
  }

  deployment {
    files {
      name = "main.py"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.main.name}"
    }

    files {
      name = "requirements.txt"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.requirements.name}"
    }
  }

  env_variables = {
    port = "8000"
  }

  instance_class = "F2"

  automatic_scaling {
    max_concurrent_requests = 10
    min_idle_instances = 1
    max_idle_instances = 3
    min_pending_latency = "1s"
    max_pending_latency = "5s"
    standard_scheduler_settings {
      target_cpu_utilization = 0.5
      target_throughput_utilization = 0.75
      min_instances = 2
      max_instances = 10
    }
  }

  noop_on_destroy = true
}

resource "google_storage_bucket" "bucket" {
  project = google_project.my_project.project_id
  name = "tf-test-%{random_suffix}-standard-ae-bucket"
}

resource "google_storage_bucket_object" "requirements" {
  name   = "requirements.txt"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world-flask/requirements.txt"
}

resource "google_storage_bucket_object" "main" {
  name   = "main.py"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world-flask/main.py"
}`, context)
}

func testAccAppEngineStandardAppVersion_pythonUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "my_project" {
  name = "tf-test-appeng-std%{random_suffix}"
  project_id = "tf-test-appeng-std%{random_suffix}"
  org_id = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_app_engine_application" "app" {
  project     = google_project.my_project.project_id
  location_id = "us-central"
}

resource "google_project_service" "project" {
  project = google_project.my_project.project_id
  service = "appengine.googleapis.com"

  disable_dependent_services = false
}

resource "google_app_engine_standard_app_version" "foo" {
  project    = google_project_service.project.project
  version_id = "v1"
  service    = "default"
  runtime    = "python38"

  entrypoint {
    shell = "gunicorn -b :$PORT main:app"
  }

  deployment {
    files {
      name = "main.py"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.main.name}"
    }

    files {
      name = "requirements.txt"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.requirements.name}"
    }
  }

  env_variables = {
    port = "8000"
  }

  instance_class = "B2"

  basic_scaling {
    max_instances = 5
  }

  noop_on_destroy = true
}

resource "google_storage_bucket" "bucket" {
  project = google_project.my_project.project_id
  name = "tf-test-%{random_suffix}-standard-ae-bucket"
}

resource "google_storage_bucket_object" "requirements" {
  name   = "requirements.txt"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world-flask/requirements.txt"
}

resource "google_storage_bucket_object" "main" {
  name   = "main.py"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world-flask/main.py"
}`, context)
}
