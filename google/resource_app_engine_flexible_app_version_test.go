package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAppEngineFlexibleAppVersion_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_name": fmt.Sprintf("tf-test-ae-service-%s", acctest.RandString(10)),
		"project":      fmt.Sprintf("tf-test-project-%s", acctest.RandString(10)),
		"org_id":       getTestOrgFromEnv(t),
		"billing_acct": getTestBillingAccountFromEnv(t),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppEngineFlexibleAppVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineFlexibleAppVersion_python(context),
			},
			{
				ResourceName:            "google_app_engine_flexible_app_version.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_variables", "deployment", "entrypoint", "service", "delete_service_on_destroy"},
			},
			{
				Config: testAccAppEngineFlexibleAppVersion_pythonUpdate(context),
			},
			{
				ResourceName:            "google_app_engine_flexible_app_version.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_variables", "deployment", "entrypoint", "service", "delete_service_on_destroy"},
			},
		},
	})
}

func testAccAppEngineFlexibleAppVersion_python(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "my_project" {
  name            = "tf-test-project"
  project_id      = "%{project}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_acct}"
}

resource "google_app_engine_flexible_app_version" "foo" {
  project    = google_project.my_project.name
  version_id = "v1"
  service    = "%{service_name}"
  runtime    = "python"

  runtime_api_version = "1"

  resources {
    cpu       = 1
    memory_gb = 0.5
    disk_gb   = 10
  }

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

    files {
      name = "app.yaml"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.yaml.name}"
    }
  }

  liveness_check {
    path = "alive"
  }

  readiness_check {
    path = "ready"
  }

  env_variables = {
    port = "8000"
  }

  network {
    name       = "default"
    subnetwork = "default"
  }

  instance_class = "B1"

  manual_scaling {
    instances = 1
  }

  delete_service_on_destroy = true
}

resource "google_storage_bucket" "bucket" {
  project = google_project.my_project.project_id
  name    = "%{service_name}-bucket"
}

resource "google_storage_bucket_object" "yaml" {
  name   = "app.yaml"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world-flask/app.yaml"
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

func testAccAppEngineFlexibleAppVersion_pythonUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "my_project" {
  name            = "tf-test-project"
  project_id      = "%{project}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_acct}"
}

resource "google_app_engine_flexible_app_version" "foo" {
  project    = google_project.my_project.name
  version_id = "v1"
  service    = "%{service_name}"
  runtime    = "python"

  runtime_api_version = "1"

  resources {
    cpu       = 1
    memory_gb = 1
    disk_gb   = 10
  }

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

    files {
      name = "app.yaml"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.yaml.name}"
    }
  }

  liveness_check {
    path = ""
  }

  readiness_check {
    path = ""
  }

  env_variables = {
    port = "8000"
  }

  network {
    name       = "default"
    subnetwork = "default"
  }

  instance_class = "B2"

  manual_scaling {
    instances = 2
  }

  delete_service_on_destroy = true
}

resource "google_storage_bucket" "bucket" {
  project = google_project.my_project.project_id
  name    = "%{service_name}-bucket"
}

resource "google_storage_bucket_object" "yaml" {
  name   = "app.yaml"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world-flask/app.yaml"
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
