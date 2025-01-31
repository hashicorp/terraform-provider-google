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
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.
func testAccAccessContextManagerServicePerimeters_basicTest(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	projectNumber := envvar.GetTestProjectNumberFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerServicePerimetersDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerServicePerimeters_basic(org, "my policy", "level", "storage_perimeter", "bigtable_perimeter", "bigquery_omni_perimeter"),
			},
			{
				ResourceName:            "google_access_context_manager_service_perimeters.test-access",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_perimeters"},
			},
			{
				Config: testAccAccessContextManagerServicePerimeters_update(org, "my policy", "level", "storage_perimeter", "bigquery_perimeter", "bigtable_perimeter", "bigquery_omni_perimeter", projectNumber),
			},
			{
				ResourceName:            "google_access_context_manager_service_perimeters.test-access",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_perimeters"},
			},
			{
				Config: testAccAccessContextManagerServicePerimeters_empty(org, "my policy", "level"),
			},
			{
				ResourceName:            "google_access_context_manager_service_perimeters.test-access",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_perimeters"},
			},
		},
	})
}

func testAccCheckAccessContextManagerServicePerimetersDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_service_perimeters" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}{{parent}}/servicePerimeters")
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
				return fmt.Errorf("ServicePerimeters still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccAccessContextManagerServicePerimeters_basic(org, policyTitle, levelTitleName, perimeterTitleName1, perimeterTitleName2, perimeterTitleName3 string) string {
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

resource "google_access_context_manager_service_perimeters" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["storage.googleapis.com"]
    }
  }

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["bigtable.googleapis.com"]
    }
  }

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["bigquery.googleapis.com"]
      egress_policies {
        egress_to {
          external_resources = ["s3://bucket1"]
          operations {
            service_name = "bigquery.googleapis.com"
            method_selectors {
              method = "*"
            }
          }
        }
        egress_from {
          identity_type = "ANY_USER_ACCOUNT"
          sources {
            access_level = google_access_context_manager_access_level.test-access.name
          }
          source_restriction = "SOURCE_RESTRICTION_ENABLED"
        }
      }
    }
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName1, perimeterTitleName1, perimeterTitleName2, perimeterTitleName2, perimeterTitleName3, perimeterTitleName3)
}

func testAccAccessContextManagerServicePerimeters_update(org, policyTitle, levelTitleName, perimeterTitleName1, perimeterTitleName2, perimeterTitleName3, perimeterTitleName4, projectNumber string) string {
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

resource "google_access_context_manager_service_perimeters" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["storage.googleapis.com"]
      access_levels       = [google_access_context_manager_access_level.test-access.name]
    }
  }

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["bigquery.googleapis.com"]
      access_levels       = [google_access_context_manager_access_level.test-access.name]
    }
  }

  service_perimeters {
    name           = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}/servicePerimeters/%s"
    title          = "%s"
    perimeter_type = "PERIMETER_TYPE_REGULAR"
    status {
      restricted_services = ["bigtable.googleapis.com"]
    }
  }

  service_perimeters {
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
    		ingress_from {
    			identities = ["user:test@google.com"]
    		}
    		ingress_to {
    			resources = ["*"]
    		}
    	}
    
    	egress_policies {
    		egress_from {
    			identity_type = "ANY_USER_ACCOUNT"
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
    		egress_from {
    			identities = ["user:test@google.com"]
    		}
    		egress_to {
    			resources = ["*"]
    		}
    	}
      egress_policies {
    		egress_from {
    			sources {
            resource = "projects/%s"
          }
          source_restriction = "SOURCE_RESTRICTION_ENABLED"
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
        ingress_from {
          identities = ["user:test@google.com"]
        }
        ingress_to {
          resources = ["*"]
        }
      }

      egress_policies {
        egress_from {
          identity_type = "ANY_USER_ACCOUNT"
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
        egress_from {
          identities = ["user:test@google.com"]
        }
        egress_to {
          resources = ["*"]
        }
      }
      egress_policies {
    		egress_from {
    			sources {
            resource = "projects/%s"
          }
          source_restriction = "SOURCE_RESTRICTION_ENABLED"
    		}
    	}
    }
  }
}
`, org, policyTitle, levelTitleName, levelTitleName, perimeterTitleName1, perimeterTitleName1, perimeterTitleName2, perimeterTitleName2, perimeterTitleName3, perimeterTitleName3, perimeterTitleName4, perimeterTitleName4, projectNumber, projectNumber)
}

func testAccAccessContextManagerServicePerimeters_empty(org, policyTitle, levelTitleName string) string {
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

resource "google_access_context_manager_service_perimeters" "test-access" {
  parent         = "accessPolicies/${google_access_context_manager_access_policy.test-access.name}"
}
`, org, policyTitle, levelTitleName, levelTitleName)
}
