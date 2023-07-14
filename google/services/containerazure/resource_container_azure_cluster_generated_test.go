// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package containerazure_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	containerazure "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccContainerAzureCluster_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"azure_app":           "00000000-0000-0000-0000-17aad2f0f61f",
		"azure_config_secret": "07d4b1f1a7cb4b1b91f070c30ae761a1",
		"azure_sub":           "00000000-0000-0000-0000-17aad2f0f61f",
		"azure_tenant":        "00000000-0000-0000-0000-17aad2f0f61f",
		"byo_prefix":          "mmv2",
		"project_name":        envvar.GetTestProjectFromEnv(),
		"project_number":      envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAzureClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAzureCluster_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_container_azure_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fleet.0.project"},
			},
			{
				Config: testAccContainerAzureCluster_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_container_azure_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fleet.0.project"},
			},
		},
	})
}

func testAccContainerAzureCluster_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_container_azure_versions" "versions" {
  project = "%{project_name}"
  location = "us-west1"
}

resource "google_container_azure_cluster" "primary" {
  authorization {
    admin_users {
      username = "mmv2@google.com"
    }
  }

  azure_region = "westus2"
  client       = "projects/%{project_number}/locations/us-west1/azureClients/${google_container_azure_client.basic.name}"

  control_plane {
    ssh_config {
      authorized_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC8yaayO6lnb2v+SedxUMa2c8vtIEzCzBjM3EJJsv8Vm9zUDWR7dXWKoNGARUb2mNGXASvI6mFIDXTIlkQ0poDEPpMaXR0g2cb5xT8jAAJq7fqXL3+0rcJhY/uigQ+MrT6s+ub0BFVbsmGHNrMQttXX9gtmwkeAEvj3mra9e5pkNf90qlKnZz6U0SVArxVsLx07vHPHDIYrl0OPG4zUREF52igbBPiNrHJFDQJT/4YlDMJmo/QT/A1D6n9ocemvZSzhRx15/Arjowhr+VVKSbaxzPtEfY0oIg2SrqJnnr/l3Du5qIefwh5VmCZe4xopPUaDDoOIEFriZ88sB+3zz8ib8sk8zJJQCgeP78tQvXCgS+4e5W3TUg9mxjB6KjXTyHIVhDZqhqde0OI3Fy1UuVzRUwnBaLjBnAwP5EoFQGRmDYk/rEYe7HTmovLeEBUDQocBQKT4Ripm/xJkkWY7B07K/tfo56dGUCkvyIVXKBInCh+dLK7gZapnd4UWkY0xBYcwo1geMLRq58iFTLA2j/JmpmHXp7m0l7jJii7d44uD3tTIFYThn7NlOnvhLim/YcBK07GMGIN7XwrrKZKmxXaspw6KBWVhzuw1UPxctxshYEaMLfFg/bwOw8HvMPr9VtrElpSB7oiOh91PDIPdPBgHCi7N2QgQ5l/ZDBHieSpNrQ== thomasrodgers"
    }

    subnet_id = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-byo/providers/Microsoft.Network/virtualNetworks/%{byo_prefix}-dev-vnet/subnets/default"
    version   = "${data.google_container_azure_versions.versions.valid_versions[0]}"
  }

  fleet {
    project = "%{project_number}"
  }

  location = "us-west1"
  name     = "tf-test-name%{random_suffix}"

  networking {
    pod_address_cidr_blocks     = ["10.200.0.0/16"]
    service_address_cidr_blocks = ["10.32.0.0/24"]
    virtual_network_id          = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-byo/providers/Microsoft.Network/virtualNetworks/%{byo_prefix}-dev-vnet"
  }

  resource_group_id = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-cluster"
  project           = "%{project_name}"
}

resource "google_container_azure_client" "basic" {
  application_id = "%{azure_app}"
  location       = "us-west1"
  name           = "tf-test-client-name%{random_suffix}"
  tenant_id      = "%{azure_tenant}"
  project        = "%{project_name}"
}


`, context)
}

func testAccContainerAzureCluster_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_container_azure_versions" "versions" {
  project = "%{project_name}"
  location = "us-west1"
}

resource "google_container_azure_cluster" "primary" {
  authorization {
    admin_users {
      username = "mmv2@google.com"
    }
  }

  azure_region = "westus2"
  client       = "projects/%{project_number}/locations/us-west1/azureClients/${google_container_azure_client.basic.name}"

  control_plane {
    ssh_config {
      authorized_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC8yaayO6lnb2v+SedxUMa2c8vtIEzCzBjM3EJJsv8Vm9zUDWR7dXWKoNGARUb2mNGXASvI6mFIDXTIlkQ0poDEPpMaXR0g2cb5xT8jAAJq7fqXL3+0rcJhY/uigQ+MrT6s+ub0BFVbsmGHNrMQttXX9gtmwkeAEvj3mra9e5pkNf90qlKnZz6U0SVArxVsLx07vHPHDIYrl0OPG4zUREF52igbBPiNrHJFDQJT/4YlDMJmo/QT/A1D6n9ocemvZSzhRx15/Arjowhr+VVKSbaxzPtEfY0oIg2SrqJnnr/l3Du5qIefwh5VmCZe4xopPUaDDoOIEFriZ88sB+3zz8ib8sk8zJJQCgeP78tQvXCgS+4e5W3TUg9mxjB6KjXTyHIVhDZqhqde0OI3Fy1UuVzRUwnBaLjBnAwP5EoFQGRmDYk/rEYe7HTmovLeEBUDQocBQKT4Ripm/xJkkWY7B07K/tfo56dGUCkvyIVXKBInCh+dLK7gZapnd4UWkY0xBYcwo1geMLRq58iFTLA2j/JmpmHXp7m0l7jJii7d44uD3tTIFYThn7NlOnvhLim/YcBK07GMGIN7XwrrKZKmxXaspw6KBWVhzuw1UPxctxshYEaMLfFg/bwOw8HvMPr9VtrElpSB7oiOh91PDIPdPBgHCi7N2QgQ5l/ZDBHieSpNrQ== thomasrodgers"
    }

    subnet_id = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-byo/providers/Microsoft.Network/virtualNetworks/%{byo_prefix}-dev-vnet/subnets/default"
    version   = "${data.google_container_azure_versions.versions.valid_versions[0]}"

    database_encryption {
      key_id = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-cluster/providers/Microsoft.KeyVault/vaults/%{byo_prefix}-dev-keyvault/keys/%{byo_prefix}-dev-key"
    }

    main_volume {
      size_gib = 8
    }

    proxy_config {
      resource_group_id = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-cluster"
      secret_id         = "https://%{byo_prefix}-dev-keyvault.vault.azure.net/secrets/%{byo_prefix}-dev-secret/%{azure_config_secret}"
    }

    replica_placements {
      azure_availability_zone = "1"
      subnet_id               = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-byo/providers/Microsoft.Network/virtualNetworks/%{byo_prefix}-dev-vnet/subnets/default"
    }

    root_volume {
      size_gib = 32
    }

    tags = {
      owner = "mmv2"
    }

    vm_size = "Standard_DS2_v2"
  }

  fleet {
    project = "%{project_number}"
  }

  location = "us-west1"
  name     = "tf-test-name%{random_suffix}"

  networking {
    pod_address_cidr_blocks     = ["10.200.0.0/16"]
    service_address_cidr_blocks = ["10.32.0.0/24"]
    virtual_network_id          = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-byo/providers/Microsoft.Network/virtualNetworks/%{byo_prefix}-dev-vnet"
  }

  resource_group_id = "/subscriptions/%{azure_sub}/resourceGroups/%{byo_prefix}-dev-cluster"

  annotations = {
    annotation-one = "value-one"
  }

  description = "An updated sample azure cluster"
  project     = "%{project_name}"
}

resource "google_container_azure_client" "basic" {
  application_id = "%{azure_app}"
  location       = "us-west1"
  name           = "tf-test-client-name%{random_suffix}"
  tenant_id      = "%{azure_tenant}"
  project        = "%{project_name}"
}



`, context)
}

func testAccCheckContainerAzureClusterDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_container_azure_cluster" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &containerazure.Cluster{
				AzureRegion:     dcl.String(rs.Primary.Attributes["azure_region"]),
				Location:        dcl.String(rs.Primary.Attributes["location"]),
				Name:            dcl.String(rs.Primary.Attributes["name"]),
				ResourceGroupId: dcl.String(rs.Primary.Attributes["resource_group_id"]),
				Client:          dcl.String(rs.Primary.Attributes["client"]),
				Description:     dcl.String(rs.Primary.Attributes["description"]),
				Project:         dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:      dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Endpoint:        dcl.StringOrNil(rs.Primary.Attributes["endpoint"]),
				Etag:            dcl.StringOrNil(rs.Primary.Attributes["etag"]),
				Reconciling:     dcl.Bool(rs.Primary.Attributes["reconciling"] == "true"),
				State:           containerazure.ClusterStateEnumRef(rs.Primary.Attributes["state"]),
				Uid:             dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:      dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLContainerAzureClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetCluster(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_container_azure_cluster still exists %v", obj)
			}
		}
		return nil
	}
}
