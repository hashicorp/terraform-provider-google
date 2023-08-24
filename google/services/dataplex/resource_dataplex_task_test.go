// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexTaskDataplexTask_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexTaskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexTask_dataplexTaskPrimary(context),
			},
			{
				ResourceName:            "google_dataplex_task.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "lake", "task_id"},
			},
			{
				Config: testAccDataplexTask_dataplexTaskPrimaryUpdate(context),
			},
			{
				ResourceName:            "google_dataplex_task.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "lake", "task_id"},
			},
		},
	})
}

func testAccDataplexTask_dataplexTaskPrimary(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {

}

resource "google_dataplex_lake" "example" {
  name         = "tf-test-lake%{random_suffix}"
  location     = "us-central1"
  project = "%{project_name}"
}


resource "google_dataplex_task" "example" {

    task_id      = "tf-test-task%{random_suffix}"
    location     = "us-central1"
    lake         = google_dataplex_lake.example.name
    trigger_spec  {
        type = "ON_DEMAND"
    }
    
    execution_spec {
        service_account = "${data.google_project.project.number}-compute@developer.gserviceaccount.com"
    }
    
    spark {
        python_script_file = "gs://dataproc-examples/pyspark/hello-world/hello-world.py"
    }
    
    project = "%{project_name}"
    
}
`, context)
}

func testAccDataplexTask_dataplexTaskPrimaryUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {

}

resource "google_dataplex_lake" "example" {
  name         = "tf-test-lake%{random_suffix}"
  location     = "us-central1"
  project = "%{project_name}"
}


resource "google_dataplex_task" "example" {

    task_id      = "tf-test-task%{random_suffix}"
    location     = "us-central1"
    lake         = google_dataplex_lake.example.name
    trigger_spec  {
        type = "ON_DEMAND"
    }
    
    execution_spec {
        service_account = "${data.google_project.project.number}-compute@developer.gserviceaccount.com"
    }
    
    spark {
        python_script_file = "gs://dataplex-clouddq-api-integration-test/clouddq_pyspark_driver.py"
    }
    
    project = "%{project_name}"
    
}
`, context)
}
