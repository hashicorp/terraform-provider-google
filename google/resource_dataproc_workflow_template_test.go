package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataprocWorkflowTemplate_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
		"project":       GetTestProjectFromEnv(),
		"version":       "2.0.35-debian10",
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: funcAccTestDataprocWorkflowTemplateCheckDestroy(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocWorkflowTemplate_basic(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_dataproc_workflow_template.template",
			},
		},
	})
}

func TestAccDataprocWorkflowTemplate_withShieldedVMs(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
		"project":       GetTestProjectFromEnv(),
		"version":       "2.0.35-debian10",
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: funcAccTestDataprocWorkflowTemplateCheckDestroy(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocWorkflowTemplate_withShieldedVMs(context),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "google_dataproc_workflow_template.shielded_vms_template",
			},
		},
	})
}

func testAccDataprocWorkflowTemplate_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_dataproc_workflow_template" "template" {
  name = "template%{random_suffix}"
  location = "us-central1"
  placement {
    managed_cluster {
      cluster_name = "my-cluster"
      config {
        gce_cluster_config {
          zone = "us-central1-a"
          tags = ["foo", "bar"]
        }
        master_config {
          num_instances = 1
          machine_type = "n1-standard-1"
          disk_config {
            boot_disk_type = "pd-ssd"
            boot_disk_size_gb = 15
          }
        }
        worker_config {
          num_instances = 3
          machine_type = "n1-standard-2"
          disk_config {
            boot_disk_size_gb = 10
            num_local_ssds = 2
          }
        }

        secondary_worker_config {
          num_instances = 2
        }
        software_config {
          image_version = "%{version}"
        }
      }
    }
  }
  jobs {
    step_id = "someJob"
    spark_job {
      main_class = "SomeClass"
    }
  }
  jobs {
    step_id = "otherJob"
    prerequisite_step_ids = ["someJob"]
    presto_job {
      query_file_uri = "someuri"
    }
  }
}
`, context)
}

func testAccDataprocWorkflowTemplate_withShieldedVMs(context map[string]interface{}) string {
	return Nprintf(`
resource "google_dataproc_workflow_template" "shielded_vms_template" {
  name = "template%{random_suffix}"
  location = "us-central1"
  placement {
    managed_cluster {
      cluster_name = "my-shielded-cluster"
      config {
        gce_cluster_config {
          zone = "us-central1-a"
          tags = ["foo", "bar"]
          shielded_instance_config {
            enable_secure_boot = true
            enable_vtpm = true
            enable_integrity_monitoring = true
          }
        }
        master_config {
          num_instances = 1
          machine_type = "n1-standard-1"
          disk_config {
            boot_disk_type = "pd-ssd"
            boot_disk_size_gb = 15
          }
        }
        worker_config {
          num_instances = 3
          machine_type = "n1-standard-2"
          disk_config {
            boot_disk_size_gb = 10
            num_local_ssds = 2
          }
        }

        secondary_worker_config {
          num_instances = 2
        }
        software_config {
          image_version = "%{version}"
        }
      }
    }
  }
  jobs {
    step_id = "someJob"
    spark_job {
      main_class = "SomeClass"
    }
  }
  jobs {
    step_id = "otherJob"
    prerequisite_step_ids = ["someJob"]
    presto_job {
      query_file_uri = "someuri"
    }
  }
}
`, context)
}

func funcAccTestDataprocWorkflowTemplateCheckDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataproc_workflow_template" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{DataprocBasePath}}projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("DataprocWorkflowTemplate still exists at %s", url)
			}
		}

		return nil
	}
}
