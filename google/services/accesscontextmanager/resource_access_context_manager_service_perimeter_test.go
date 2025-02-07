// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/accesscontextmanager"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.
func testAccAccessContextManagerServicePerimeter_basicTest(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerServicePerimeterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeter_basic(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAccessContextManagerServicePerimeter_updateTest(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	projectNumber := envvar.GetTestProjectNumberFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerServicePerimeterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeter_basic(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeter_update(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeter_updateAllowed(org, "my policy", "level", "perimeter", projectNumber),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeter_updateDryrun(org, "my policy", "level", "perimeter"),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerServicePerimeter_updateAllowed(org, "my policy", "level", "perimeter", projectNumber),
			},
			{
				ResourceName:      "google_access_context_manager_service_perimeter.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerServicePerimeterDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_service_perimeter" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ServicePerimeter still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerServicePerimeter_basic(org, policyTitle, levelTitleName, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_level" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
  title       = "%s"
  description = "hello"
  basic {
    combining_function = "AND"
    conditions {
      ip_subnetworks = ["192.0.4.0/24"]
    }
  }
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  perimeter_type = "PERIMETER_TYPE_REGULAR"
  status {
    restricted_services = ["storage.googleapis.com"]
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName)
}

func testAccAccessContextManagerServicePerimeter_update(org, policyTitle, levelTitleName, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_level" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
  title       = "%s"
  description = "hello"
  basic {
    combining_function = "AND"
    conditions {
      ip_subnetworks = ["192.0.4.0/24"]
    }
  }
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  perimeter_type = "PERIMETER_TYPE_REGULAR"
  status {
    restricted_services = ["bigquery.googleapis.com"]
    access_levels       = [google_access_context_manager_access_level.test-access.name]
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName)
}

func testAccAccessContextManagerServicePerimeter_updateAllowed(org, policyTitle, levelTitleName, perimeterTitleName, projectNumber string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_level" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
  title       = "%s"
  description = "hello"
  basic {
    combining_function = "AND"
    conditions {
      ip_subnetworks = ["192.0.4.0/24"]
    }
  }
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  perimeter_type = "PERIMETER_TYPE_REGULAR"
  use_explicit_dry_run_spec = true
  spec {
    restricted_services = ["bigquery.googleapis.com", "storage.googleapis.com"]
		access_levels       = [google_access_context_manager_access_level.test-access.name]

		vpc_accessible_services {
			enable_restriction = true
			allowed_services   = ["bigquery.googleapis.com", "storage.googleapis.com"]
		}

		ingress_policies {
			title = "ingress policy 1"

			ingress_from {
				sources {
					access_level = google_access_context_manager_access_level.test-access.name
				}
				identity_type = "ANY_IDENTITY"
			}

			ingress_to {
				resources = [ "*" ]
				operations {
					service_name = "bigquery.googleapis.com"

					method_selectors {
						method = "BigQueryStorage.ReadRows"
					}

					method_selectors {
						method = "TableService.ListTables"
					}

					method_selectors {
						permission = "bigquery.jobs.get"
					}
				}

				operations {
					service_name = "storage.googleapis.com"

					method_selectors {
						method = "google.storage.objects.create"
					}
				}
			}
		}
		ingress_policies {
			title = "ingress policy 2"
			ingress_from {
				identities = ["user:test@google.com"]
			}
			ingress_to {
				resources = ["*"]
			}
		}

		egress_policies {
			title = "egress policy 1"
			egress_from {
				identity_type = "ANY_USER_ACCOUNT"
				sources {
					access_level = google_access_context_manager_access_level.test-access.name
				}
					
				sources {
					resource = "projects/%s"
				}
					
				source_restriction = "SOURCE_RESTRICTION_ENABLED"
			}
			egress_to {
				operations {
					service_name = "bigquery.googleapis.com"
					method_selectors {
						permission = "externalResource.read"
					}
				}
				external_resources = ["s3://bucket1"]
			}
		}
		egress_policies {
			title = "egress policy 2"
			egress_from {
				identities = ["user:test@google.com"]
			}
			egress_to {
				resources = ["*"]
			}
		}
  }
  status {
    restricted_services = ["bigquery.googleapis.com", "storage.googleapis.com"]
		access_levels       = [google_access_context_manager_access_level.test-access.name]

		vpc_accessible_services {
			enable_restriction = true
			allowed_services   = ["bigquery.googleapis.com", "storage.googleapis.com"]
		}

		ingress_policies {
			title = "ingress policy 1"

			ingress_from {
				sources {
					access_level = google_access_context_manager_access_level.test-access.name
				}
				identity_type = "ANY_IDENTITY"
			}

			ingress_to {
				resources = [ "*" ]
				operations {
					service_name = "bigquery.googleapis.com"

					method_selectors {
						method = "BigQueryStorage.ReadRows"
					}

					method_selectors {
						method = "TableService.ListTables"
					}

					method_selectors {
						permission = "bigquery.jobs.get"
					}
				}

				operations {
					service_name = "storage.googleapis.com"

					method_selectors {
						method = "google.storage.objects.create"
					}
				}
			}
		}
		ingress_policies {
			title = "ingress policy 2"
			ingress_from {
				identities = ["user:test@google.com"]
			}
			ingress_to {
				resources = ["*"]
			}
		}

		egress_policies {
			title = "egress policy 1"
			egress_from {
				identity_type = "ANY_USER_ACCOUNT"
				sources {
					access_level = google_access_context_manager_access_level.test-access.name
				}

				sources {
					resource = "projects/%s"
				}
					
				source_restriction = "SOURCE_RESTRICTION_ENABLED"
			}
			egress_to {
				operations {
					service_name = "bigquery.googleapis.com"
					method_selectors {
						permission = "externalResource.read"
					}
				}
				external_resources = ["s3://bucket1"]
			}
		}
		egress_policies {
			title = "egress policy 2"
			egress_from {
				identities = ["user:test@google.com"]
			}
			egress_to {
				resources = ["*"]
			}
		}
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName, projectNumber, projectNumber)
}

func testAccAccessContextManagerServicePerimeter_updateDryrun(org, policyTitle, levelTitleName, perimeterTitleName string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}

resource "google_access_context_manager_access_level" "test-access" {
  parent      = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name        = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/accessLevels/%s"
  title       = "%s"
  description = "hello"
  basic {
    combining_function = "AND"
    conditions {
      ip_subnetworks = ["192.0.4.0/24"]
    }
  }
}

resource "google_access_context_manager_service_perimeter" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
  name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
  title          = "%s"
  perimeter_type = "PERIMETER_TYPE_REGULAR"
  status {
    restricted_services = ["bigquery.googleapis.com"]
  }

  spec {
    restricted_services = ["storage.googleapis.com"]
	access_levels       = [google_access_context_manager_access_level.test-access.name]
  }

  use_explicit_dry_run_spec = true
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName, perimeterTitleName)
}

type IdentityTypeDiffSuppressFuncDiffSuppressTestCase struct {
	Name     string
	AreEqual bool
	Before   string
	After    string
}

var identityTypeDiffSuppressTestCases = []IdentityTypeDiffSuppressFuncDiffSuppressTestCase{
	{
		AreEqual: false,
		Before:   "A",
		After:    "B",
	},
	{
		AreEqual: true,
		Before:   "A",
		After:    "A",
	},
	{
		AreEqual: false,
		Before:   "",
		After:    "A",
	},
	{
		AreEqual: false,
		Before:   "A",
		After:    "",
	},
	{
		AreEqual: true,
		Before:   "",
		After:    "IDENTITY_TYPE_UNSPECIFIED",
	},
	{
		AreEqual: false,
		Before:   "IDENTITY_TYPE_UNSPECIFIED",
		After:    "",
	},
}

func TestUnitAccessContextManagerServicePerimeter_identityTypeDiff(t *testing.T) {
	for _, tc := range identityTypeDiffSuppressTestCases {
		tc.Test(t)
	}
}

func (tc *IdentityTypeDiffSuppressFuncDiffSuppressTestCase) Test(t *testing.T) {
	actual := accesscontextmanager.AccessContextManagerServicePerimeterIdentityTypeDiffSuppressFunc("", tc.Before, tc.After, nil)
	if actual != tc.AreEqual {
		t.Errorf(
			"Unexpected difference found. Before: \"%s\", after: \"%s\", actual: %t, expected: %t",
			tc.Before, tc.After, actual, tc.AreEqual)
	}
}
