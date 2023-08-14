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

package containeraws_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	containeraws "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containeraws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccContainerAwsNodePool_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"aws_acct_id":    "111111111111",
		"aws_db_key":     "00000000-0000-0000-0000-17aad2f0f61f",
		"aws_region":     "us-west-2",
		"aws_sg":         "sg-0b3f63cb91b247628",
		"aws_subnet":     "subnet-0b3f63cb91b247628",
		"aws_vol_key":    "00000000-0000-0000-0000-17aad2f0f61f",
		"aws_vpc":        "vpc-0b3f63cb91b247628",
		"byo_prefix":     "mmv2",
		"project_name":   envvar.GetTestProjectFromEnv(),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"service_acct":   envvar.GetTestServiceAccountFromEnv(t),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAwsNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAwsNodePool_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fleet.0.project"},
			},
			{
				Config: testAccContainerAwsNodePool_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fleet.0.project"},
			},
		},
	})
}
func TestAccContainerAwsNodePool_BasicEnumHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"aws_acct_id":    "111111111111",
		"aws_db_key":     "00000000-0000-0000-0000-17aad2f0f61f",
		"aws_region":     "us-west-2",
		"aws_sg":         "sg-0b3f63cb91b247628",
		"aws_subnet":     "subnet-0b3f63cb91b247628",
		"aws_vol_key":    "00000000-0000-0000-0000-17aad2f0f61f",
		"aws_vpc":        "vpc-0b3f63cb91b247628",
		"byo_prefix":     "mmv2",
		"project_name":   envvar.GetTestProjectFromEnv(),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"service_acct":   envvar.GetTestServiceAccountFromEnv(t),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAwsNodePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAwsNodePool_BasicEnumHandWritten(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fleet.0.project"},
			},
			{
				Config: testAccContainerAwsNodePool_BasicEnumHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_container_aws_node_pool.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fleet.0.project"},
			},
		},
	})
}

func testAccContainerAwsNodePool_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_container_aws_versions" "versions" {
  project = "%{project_name}"
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  authorization {
    admin_users {
      username = "%{service_acct}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-oneplatform"
      role_session_name = "%{byo_prefix}-1p-dev-session"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "t3.medium"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      owner = "%{service_acct}"
    }
  }

  fleet {
    project = "%{project_number}"
  }

  location = "us-west1"
  name     = "tf-test-name%{random_suffix}"

  networking {
    pod_address_cidr_blocks     = ["10.2.0.0/16"]
    service_address_cidr_blocks = ["10.1.0.0/16"]
    vpc_id                      = "%{aws_vpc}"
  }

  annotations = {
    label-one = "value-one"
  }

  description = "A sample aws cluster"
  project     = "%{project_name}"
}


resource "google_container_aws_node_pool" "primary" {
  autoscaling {
    max_node_count = 5
    min_node_count = 1
  }

  cluster = google_container_aws_cluster.primary.name

  config {
    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool"
    instance_type        = "t3.medium"

    labels = {
      label-one = "value-one"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}"]

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      tag-one = "value-one"
    }

    taints {
      effect = "PREFER_NO_SCHEDULE"
      key    = "taint-key"
      value  = "taint-value"
    }
  }

  location = "us-west1"

  max_pods_constraint {
    max_pods_per_node = 110
  }

  name      = "tf-test-node-pool-name%{random_suffix}"
  subnet_id = "%{aws_subnet}"
  version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"

  annotations = {
    label-one = "value-one"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccContainerAwsNodePool_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_container_aws_versions" "versions" {
  project = "%{project_name}"
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  authorization {
    admin_users {
      username = "%{service_acct}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-oneplatform"
      role_session_name = "%{byo_prefix}-1p-dev-session"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "t3.medium"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      owner = "%{service_acct}"
    }
  }

  fleet {
    project = "%{project_number}"
  }

  location = "us-west1"
  name     = "tf-test-name%{random_suffix}"

  networking {
    pod_address_cidr_blocks     = ["10.2.0.0/16"]
    service_address_cidr_blocks = ["10.1.0.0/16"]
    vpc_id                      = "%{aws_vpc}"
  }

  annotations = {
    label-one = "value-one"
  }

  description = "A sample aws cluster"
  project     = "%{project_name}"
}

resource "google_container_aws_node_pool" "primary" {
  autoscaling {
    max_node_count = 5
    min_node_count = 1
  }

  cluster = google_container_aws_cluster.primary.name

  config {
    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool"
    instance_type        = "t3.large"

    labels = {
      label-one = "value-one"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}"]

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      tag-one = "value-one"
    }

    taints {
      effect = "PREFER_NO_SCHEDULE"
      key    = "taint-key"
      value  = "taint-value"
    }
  }

  location = "us-west1"

  max_pods_constraint {
    max_pods_per_node = 110
  }

  name      = "tf-test-node-pool-name%{random_suffix}"
  subnet_id = "%{aws_subnet}"
  version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"

  annotations = {
    label-two = "value-two"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccContainerAwsNodePool_BasicEnumHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_container_aws_versions" "versions" {
  project = "%{project_name}"
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  authorization {
    admin_users {
      username = "%{service_acct}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-oneplatform"
      role_session_name = "%{byo_prefix}-1p-dev-session"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "t3.medium"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      owner = "%{service_acct}"
    }
  }

  fleet {
    project = "%{project_number}"
  }

  location = "us-west1"
  name     = "tf-test-name%{random_suffix}"

  networking {
    pod_address_cidr_blocks     = ["10.2.0.0/16"]
    service_address_cidr_blocks = ["10.1.0.0/16"]
    vpc_id                      = "%{aws_vpc}"
  }

  annotations = {
    label-one = "value-one"
  }

  description = "A sample aws cluster"
  project     = "%{project_name}"
}


resource "google_container_aws_node_pool" "primary" {
  autoscaling {
    max_node_count = 5
    min_node_count = 1
  }

  cluster = google_container_aws_cluster.primary.name

  config {
    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool"
    instance_type        = "t3.medium"

    labels = {
      label-one = "value-one"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "gp3"
    }

    security_group_ids = ["%{aws_sg}"]

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      tag-one = "value-one"
    }

    taints {
      effect = "prefer_no_schedule"
      key    = "taint-key"
      value  = "taint-value"
    }
  }

  location = "us-west1"

  max_pods_constraint {
    max_pods_per_node = 110
  }

  name      = "tf-test-node-pool-name%{random_suffix}"
  subnet_id = "%{aws_subnet}"
  version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"

  annotations = {
    label-one = "value-one"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccContainerAwsNodePool_BasicEnumHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_container_aws_versions" "versions" {
  project = "%{project_name}"
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  authorization {
    admin_users {
      username = "%{service_acct}"
    }
  }

  aws_region = "%{aws_region}"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::%{aws_acct_id}:role/%{byo_prefix}-1p-dev-oneplatform"
      role_session_name = "%{byo_prefix}-1p-dev-session"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-controlplane"
    subnet_ids           = ["%{aws_subnet}"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "t3.medium"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["%{aws_sg}"]

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      owner = "%{service_acct}"
    }
  }

  fleet {
    project = "%{project_number}"
  }

  location = "us-west1"
  name     = "tf-test-name%{random_suffix}"

  networking {
    pod_address_cidr_blocks     = ["10.2.0.0/16"]
    service_address_cidr_blocks = ["10.1.0.0/16"]
    vpc_id                      = "%{aws_vpc}"
  }

  annotations = {
    label-one = "value-one"
  }

  description = "A sample aws cluster"
  project     = "%{project_name}"
}

resource "google_container_aws_node_pool" "primary" {
  autoscaling {
    max_node_count = 5
    min_node_count = 1
  }

  cluster = google_container_aws_cluster.primary.name

  config {
    config_encryption {
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_db_key}"
    }

    iam_instance_profile = "%{byo_prefix}-1p-dev-nodepool"
    instance_type        = "t3.large"

    labels = {
      label-one = "value-one"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:%{aws_region}:%{aws_acct_id}:key/%{aws_vol_key}"
      size_gib    = 10
      volume_type = "gp3"
    }

    security_group_ids = ["%{aws_sg}"]

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    ssh_config {
      ec2_key_pair = "%{byo_prefix}-1p-dev-ssh"
    }

    tags = {
      tag-one = "value-one"
    }

    taints {
      effect = "prefer_no_schedule"
      key    = "taint-key"
      value  = "taint-value"
    }
  }

  location = "us-west1"

  max_pods_constraint {
    max_pods_per_node = 110
  }

  name      = "tf-test-node-pool-name%{random_suffix}"
  subnet_id = "%{aws_subnet}"
  version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"

  annotations = {
    label-two = "value-two"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccCheckContainerAwsNodePoolDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_container_aws_node_pool" {
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

			obj := &containeraws.NodePool{
				Cluster:     dcl.String(rs.Primary.Attributes["cluster"]),
				Location:    dcl.String(rs.Primary.Attributes["location"]),
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				SubnetId:    dcl.String(rs.Primary.Attributes["subnet_id"]),
				Version:     dcl.String(rs.Primary.Attributes["version"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Etag:        dcl.StringOrNil(rs.Primary.Attributes["etag"]),
				Reconciling: dcl.Bool(rs.Primary.Attributes["reconciling"] == "true"),
				State:       containeraws.NodePoolStateEnumRef(rs.Primary.Attributes["state"]),
				Uid:         dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLContainerAwsClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetNodePool(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_container_aws_node_pool still exists %v", obj)
			}
		}
		return nil
	}
}
