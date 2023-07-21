// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/sql"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// Fields that should be ignored in import tests because they aren't returned
// from GCP (and thus can't be imported)
var ignoredReplicaConfigurationFields = []string{
	"replica_configuration.0.ca_certificate",
	"replica_configuration.0.client_certificate",
	"replica_configuration.0.client_key",
	"replica_configuration.0.connect_retry_interval",
	"replica_configuration.0.dump_file_path",
	"replica_configuration.0.master_heartbeat_period",
	"replica_configuration.0.password",
	"replica_configuration.0.ssl_cipher",
	"replica_configuration.0.username",
	"replica_configuration.0.verify_server_certificate",
	"deletion_protection",
}

func TestAccSqlDatabaseInstance_basicInferredName(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_basic2,
				Check: resource.ComposeTestCheckFunc(
					checkInstanceTypeIsPresent("google_sql_database_instance.instance"),
				),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicSecondGen(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
				Check: testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(t, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicMSSQL(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	rootPassword := acctest.RandString(t, 15)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_mssql, databaseName, rootPassword),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_update_mssql, databaseName, rootPassword),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_dontDeleteDefaultUserOnReplica(t *testing.T) {
	t.Parallel()

	databaseName := "sql-instance-test-" + acctest.RandString(t, 10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(t, 10)
	// 1. Create an instance.
	// 2. Add a root@'%' user.
	// 3. Create a replica and assert it succeeds (it'll fail if we try to delete the root user thinking it's a
	//    default user)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withoutReplica(databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				PreConfig: func() {
					// Add a root user
					config := acctest.GoogleProviderConfig(t)
					user := sqladmin.User{
						Name:     "root",
						Host:     "%",
						Password: acctest.RandString(t, 26),
					}
					op, err := config.NewSqlAdminClient(config.UserAgent).Users.Insert(config.Project, databaseName, &user).Do()
					if err != nil {
						t.Errorf("Error while inserting root@%% user: %s", err)
						return
					}
					err = sql.SqlAdminOperationWaitTime(config, op, config.Project, "Waiting for user to insert", config.UserAgent, 10*time.Minute)
					if err != nil {
						t.Errorf("Error while waiting for user insert operation to complete: %s", err.Error())
					}
					// User was created, now create replica
				},
				Config: testGoogleSqlDatabaseInstanceConfig_withReplica(databaseName, failoverName),
			},
		},
	})
}

// This test requires an arbitrary error to occur during the SQL instance creation. Currently, we
// are relying on an error that occurs when settings are used on a MySQL clone. Note that this is
// somewhat brittle, and the test could begin failing if that error no longer behaves the same way.
// If this test begins failing, we will want to be sure that the root user is removed as early as
// possible, and we should attempt to find any other scenarios where the root user could otherwise
// be left on the instance.
func TestAccSqlDatabaseInstance_deleteDefaultUserBeforeSubsequentApiCalls(t *testing.T) {
	// Service Networking
	acctest.SkipIfVcr(t)
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	addressName := "tf-test-" + acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedTestNetwork(t, "sql-instance-private-clone-2")

	// 1. Create an instance.
	// 2. Add a root@'%' user.
	// 3. Create a clone with settings and assert it fails after the instance creation API call.
	// 4. Check root user was deleted.
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressName, false, false),
			},
			{
				PreConfig: func() {
					// Add a root user
					config := acctest.GoogleProviderConfig(t)
					user := sqladmin.User{
						Name:     "root",
						Host:     "%",
						Password: acctest.RandString(t, 26),
					}
					op, err := config.NewSqlAdminClient(config.UserAgent).Users.Insert(config.Project, databaseName, &user).Do()
					if err != nil {
						t.Errorf("Error while inserting root@%% user: %s", err)
						return
					}
					err = sql.SqlAdminOperationWaitTime(config, op, config.Project, "Waiting for user to insert", config.UserAgent, 10*time.Minute)
					if err != nil {
						t.Errorf("Error while waiting for user insert operation to complete: %s", err.Error())
					}
				},
				Config:      testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone_withSettings(databaseName, networkName, addressName),
				ExpectError: regexp.MustCompile("Error, failed to update instance settings"),
			},
			{
				// This PreConfig does a check on the previous step. It is needed because the
				// previous step expects an error, so it cannot perform a check of its own.
				PreConfig: func() {
					var s *terraform.State
					err := testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(t, databaseName+"-clone1")(s)
					if err != nil {
						t.Errorf("Failed to verify that root user was removed: %s", err)
					}
				},
				Config:             testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone_withSettings(databaseName, networkName, addressName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_basic(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_secondary(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_secondary, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_deletionProtection(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtection, databaseName, "true"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtection, databaseName, "true"),
				Destroy:     true,
				ExpectError: regexp.MustCompile("Error, failed to delete instance because deletion_protection is set to true. Set it to false to proceed with instance deletion"),
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtection, databaseName, "false"),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_maintenanceVersion(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenanceVersionWithOldVersion, databaseName),
				ExpectError: regexp.MustCompile(
					`.*Maintenance version \(MYSQL_5_7_37.R20210508.01_03\) must not be set.*`),
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenanceVersionWithOldVersion, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_deletionProtectionEnabled(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtectionEnabled, databaseName, "true"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtectionEnabled, databaseName, "true"),
				Destroy:     true,
				ExpectError: regexp.MustCompile(fmt.Sprintf("Error, failed to delete instance %s: googleapi: Error 400: The instance is protected. Please disable the deletion protection and try again. To disable deletion protection, update the instance settings with deletionProtectionEnabled set to false.", databaseName)),
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_deletionProtectionEnabled, databaseName, "false"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_checkServiceNetworking(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings_checkServiceNetworking, databaseName, databaseName),
				ExpectError: regexp.MustCompile("Error, failed to create instance because the network doesn't have at least 1 private services connection. Please see https://cloud.google.com/sql/docs/mysql/private-ip#network_requirements for how to create this connection."),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_replica(t *testing.T) {
	t.Parallel()

	databaseID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_replica, databaseID, databaseID, databaseID, "true"),
				ExpectError: regexp.MustCompile("Error, failed to create instance tf-test-\\d+-2: googleapi: Error 400: Invalid request: Invalid flag for instance role: Backups cannot be enabled for read replica instance.., invalid"),
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_replica, databaseID, databaseID, databaseID, "false"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance_master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.replica1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
			{
				ResourceName:            "google_sql_database_instance.replica2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_slave(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)
	slaveID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_slave, masterID, slaveID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance_master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance_slave",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_highAvailability(t *testing.T) {
	t.Parallel()

	instanceID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_highAvailability, instanceID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_diskspecs(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_diskspecs, masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_maintenance(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenance, masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_upgrade(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settingsDowngrade(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

// GH-4222
func TestAccSqlDatabaseInstance_authNets(t *testing.T) {
	t.Parallel()

	databaseID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step2, databaseID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

// Tests that a SQL instance can be referenced from more than one other resource without
// throwing an error during provisioning, see #9018.
func TestAccSqlDatabaseInstance_multipleOperations(t *testing.T) {
	t.Parallel()

	databaseID, instanceID, userID := acctest.RandString(t, 8), acctest.RandString(t, 8), acctest.RandString(t, 8)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_multipleOperations, databaseID, instanceID, userID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basic_with_user_labels(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_with_user_labels, databaseName),
				Check: testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(t, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_with_user_labels_update, databaseName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	addressName := "tf-test-" + acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedTestNetwork(t, "sql-instance-private")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressName, false, false),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressName, true, false),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressName, true, true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressName, true, false),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(t *testing.T) {
	// Service Networking
	acctest.SkipIfVcr(t)
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	addressName := "tf-test-" + acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedTestNetwork(t, "sql-instance-private-allocated-ip-range")
	addressName_update := "tf-test-" + acctest.RandString(t, 10) + "update"
	networkName_update := acctest.BootstrapSharedTestNetwork(t, "sql-instance-private-allocated-ip-range-update")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(databaseName, networkName, addressName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(databaseName, networkName_update, addressName_update),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeReplica(t *testing.T) {
	// Service Networking
	acctest.SkipIfVcr(t)
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	addressName := "tf-test-" + acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedTestNetwork(t, "sql-instance-private-replica")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeReplica(databaseName, networkName, addressName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.replica1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone(t *testing.T) {
	// Service Networking
	acctest.SkipIfVcr(t)
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	addressName := "tf-test-" + acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedTestNetwork(t, "sql-instance-private-clone")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone(databaseName, networkName, addressName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.clone1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "clone"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_createFromBackup(t *testing.T) {
	// Sqladmin client
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"original_db_name": acctest.BootstrapSharedSQLInstanceBackupRun(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_restoreFromBackup(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "restore_backup_context"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_backupUpdate(t *testing.T) {
	// Sqladmin client
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"original_db_name": acctest.BootstrapSharedSQLInstanceBackupRun(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_beforeBackup(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_restoreFromBackup(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "restore_backup_context"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicClone(t *testing.T) {
	// Sqladmin client
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"original_db_name": acctest.BootstrapSharedSQLInstanceBackupRun(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_basicClone(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "clone"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_cloneWithSettings(t *testing.T) {
	// Sqladmin client
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"original_db_name": acctest.BootstrapSharedSQLInstanceBackupRun(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_cloneWithSettings(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "clone"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_cloneWithDatabaseNames(t *testing.T) {
	// Sqladmin client
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"original_db_name": acctest.BootstrapSharedSQLInstanceBackupRun(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_cloneWithDatabaseNames(context),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "clone"},
			},
		},
	})
}

func testAccSqlDatabaseInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			config := acctest.GoogleProviderConfig(t)
			if rs.Type != "google_sql_database_instance" {
				continue
			}

			_, err := config.NewSqlAdminClient(config.UserAgent).Instances.Get(config.Project,
				rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("Database Instance still exists")
			}
		}

		return nil
	}
}

func testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(t *testing.T, instance string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		users, err := config.NewSqlAdminClient(config.UserAgent).Users.List(config.Project, instance).Do()

		if err != nil {
			return fmt.Errorf("Could not list database users for %q: %s", instance, err)
		}

		for _, u := range users.Items {
			if u.Name == "root" && u.Host == "%" {
				return fmt.Errorf("%v@%v user still exists", u.Name, u.Host)
			}
		}

		return nil
	}
}

func TestAccSqlDatabaseInstance_BackupRetention(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_BackupRetention(masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_PointInTimeRecoveryEnabled(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID, true, "POSTGRES_9_6"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID, false, "POSTGRES_9_6"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_PointInTimeRecoveryEnabledForSqlServer(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID, true, "SQLSERVER_2017_STANDARD"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID, false, "SQLSERVER_2017_STANDARD"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_insights(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_insights, masterID),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_encryptionKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
		"instance_name": "tf-test-sql-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: acctest.Nprintf(
					testGoogleSqlDatabaseInstance_encryptionKey, context),
			},
			{
				ResourceName:            "google_sql_database_instance.replica",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_encryptionKey_replicaInDifferentRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
		"instance_name": "tf-test-sql-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: acctest.Nprintf(
					testGoogleSqlDatabaseInstance_encryptionKey_replicaInDifferentRegion, context),
			},
			{
				ResourceName:            "google_sql_database_instance.replica",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_ActiveDirectory(t *testing.T) {
	// skip the test until Active Directory setup issue gets resolved
	// see https://github.com/hashicorp/terraform-provider-google/issues/13517
	t.Skip()

	t.Parallel()
	databaseName := "tf-test-" + acctest.RandString(t, 10)
	networkName := acctest.BootstrapSharedTestNetwork(t, "sql-instance-private-test-ad")
	addressName := "tf-test-" + acctest.RandString(t, 10)
	rootPassword := acctest.RandString(t, 15)
	adDomainName := acctest.BootstrapSharedTestADDomain(t, "test-domain", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_ActiveDirectoryConfig(databaseName, networkName, addressName, rootPassword, adDomainName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance-with-ad",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSQLDatabaseInstance_DenyMaintenancePeriod(t *testing.T) {
	t.Parallel()
	databaseName := "tf-test-" + acctest.RandString(t, 10)
	endDate := "2022-12-5"
	startDate := "2022-10-5"
	time := "00:00:00"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_DenyMaintenancePeriodConfig(databaseName, endDate, startDate, time),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_Edition(t *testing.T) {
	t.Parallel()
	enterprisePlusName := "tf-test-enterprise-plus" + acctest.RandString(t, 10)
	enterprisePlusTier := "db-perf-optimized-N-2"
	enterpriseName := "tf-test-enterprise-" + acctest.RandString(t, 10)
	enterpriseTier := "db-custom-2-13312"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(enterprisePlusName, enterprisePlusTier, "ENTERPRISE_PLUS"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(enterpriseName, enterpriseTier, "ENTERPRISE"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSQLDatabaseInstance_sqlMysqlDataCacheConfig(t *testing.T) {
	t.Parallel()
	instanceName := "tf-test-enterprise-plus" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_sqlMysqlDataCacheConfig(instanceName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_SqlServerAuditConfig(t *testing.T) {
	// Service Networking
	acctest.SkipIfVcr(t)
	t.Parallel()
	databaseName := "tf-test-" + acctest.RandString(t, 10)
	rootPassword := acctest.RandString(t, 15)
	bucketName := fmt.Sprintf("%s-%d", "tf-test-bucket", acctest.RandInt(t))
	uploadInterval := "900s"
	retentionInterval := "86400s"
	bucketNameUpdate := fmt.Sprintf("%s-%d", "tf-test-bucket", acctest.RandInt(t)) + "update"
	uploadIntervalUpdate := "1200s"
	retentionIntervalUpdate := "172800s"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_SqlServerAuditConfig(databaseName, rootPassword, bucketName, uploadInterval, retentionInterval),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_SqlServerAuditConfig(databaseName, rootPassword, bucketNameUpdate, uploadIntervalUpdate, retentionIntervalUpdate),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_SqlServerAuditOptionalBucket(t *testing.T) {
	t.Parallel()
	databaseName := "tf-test-" + acctest.RandString(t, 10)
	rootPassword := acctest.RandString(t, 15)
	uploadInterval := "900s"
	retentionInterval := "86400s"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_SqlServerAuditOptionalBucket(databaseName, rootPassword, uploadInterval, retentionInterval),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_Smt(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	rootPassword := acctest.RandString(t, 15)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_Smt(databaseName, rootPassword, 1),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_Smt(databaseName, rootPassword, 2),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_NullSmt(databaseName, rootPassword),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_Timezone(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	rootPassword := acctest.RandString(t, 15)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_Timezone(databaseName, rootPassword, "Pacific Standard Time"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_sqlMysqlInstancePvpExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_sqlMysqlInstancePvpExample(context),
			},
			{
				ResourceName:            "google_sql_database_instance.mysql_pvp_instance_name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_updateReadReplicaWithBinaryLogEnabled(t *testing.T) {
	t.Parallel()

	instance := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_readReplica(instance),
			},
			{
				ResourceName:            "google_sql_database_instance.replica",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_updateReadReplica(instance),
			},
			{
				ResourceName:            "google_sql_database_instance.replica",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_rootPasswordShouldBeUpdatable(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	rootPwd := "rootPassword-1-" + acctest.RandString(t, 10)
	newRootPwd := "rootPassword-2-" + acctest.RandString(t, 10)
	databaseVersion := "SQLSERVER_2017_STANDARD"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_updateRootPassword(databaseName, databaseVersion, rootPwd),
			},
			{
				ResourceName:            "google_sql_database_instance.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_updateRootPassword(databaseName, databaseVersion, newRootPwd),
			},
			{
				ResourceName:            "google_sql_database_instance.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_updateRootPassword(databaseName, databaseVersion, ""),
				ExpectError: regexp.MustCompile(
					`Error, root password cannot be empty for SQL Server instance.`),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_activationPolicy(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_activationPolicy(instanceName, "MYSQL_5_7", "ALWAYS", true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_activationPolicy(instanceName, "MYSQL_5_7", "NEVER", true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_activationPolicy(instanceName, "MYSQL_8_0_18", "ALWAYS", true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_activationPolicy(instanceName, "MYSQL_8_0_26", "NEVER", true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_activationPolicy(instanceName, "MYSQL_8_0_26", "ALWAYS", false),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_ReplicaPromoteSuccessful(t *testing.T) {
	t.Parallel()

	databaseName := "sql-instance-test-" + acctest.RandString(t, 10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withReplica(databaseName, failoverName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: googleSqlDatabaseInstance_replicaPromote(databaseName, failoverName),
				Check:  resource.ComposeTestCheckFunc(checkPromoteReplicaConfigurations("google_sql_database_instance.instance-failover")),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_ReplicaPromoteFailedWithMasterInstanceNamePresent(t *testing.T) {
	t.Parallel()
	databaseName := "sql-instance-test-" + acctest.RandString(t, 10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withReplica(databaseName, failoverName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config:      googleSqlDatabaseInstance_replicaPromoteWithMasterInstanceName(databaseName, failoverName),
				ExpectError: regexp.MustCompile("Replica promote configuration check failed. Please remove master_instance_name and try again."),
				Check:       resource.ComposeTestCheckFunc(checkPromoteReplicaSkipConfigurations("google_sql_database_instance.instance-failover")),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_ReplicaPromoteFailedWithReplicaConfigurationPresent(t *testing.T) {
	t.Parallel()

	databaseName := "sql-instance-test-" + acctest.RandString(t, 10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withReplica(databaseName, failoverName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config:      googleSqlDatabaseInstance_replicaPromoteWithReplicaConfiguration(databaseName, failoverName),
				ExpectError: regexp.MustCompile("Replica promote configuration check failed. Please remove replica_configuration and try again."),
				Check:       resource.ComposeTestCheckFunc(checkPromoteReplicaSkipConfigurations("google_sql_database_instance.instance-failover")),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_ReplicaPromoteFailedWithMasterInstanceNameAndReplicaConfigurationPresent(t *testing.T) {
	t.Parallel()

	databaseName := "sql-instance-test-" + acctest.RandString(t, 10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withReplica(databaseName, failoverName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config:      googleSqlDatabaseInstance_replicaPromoteWithMasterInstanceNameAndReplicaConfiguration(databaseName, failoverName),
				ExpectError: regexp.MustCompile("Replica promote configuration check failed. Please remove master_instance_name and try again."),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_ReplicaPromoteSkippedWithNoMasterInstanceNameAndNoReplicaConfigurationPresent(t *testing.T) {
	t.Parallel()

	databaseName := "sql-instance-test-" + acctest.RandString(t, 10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withReplica(databaseName, failoverName),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: googleSqlDatabaseInstance_replicaPromoteSkippedWithNoMasterInstanceNameAndNoReplicaConfiguration(databaseName, failoverName),
				Check:  resource.ComposeTestCheckFunc(checkPromoteReplicaSkipConfigurations("google_sql_database_instance.instance-failover")),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            "google_sql_database_instance.instance-failover",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testAccSqlDatabaseInstance_sqlMysqlInstancePvpExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "mysql_pvp_instance_name" {
  name             = "tf-test-mysql-pvp-instance-name%{random_suffix}"
  region           = "asia-northeast1"
  database_version = "MYSQL_8_0"
  root_password = "abcABC123!"
  settings {
    tier              = "db-f1-micro"
    password_validation_policy {
      min_length  = 6
      complexity  =  "COMPLEXITY_DEFAULT"
      reuse_interval = 2
      disallow_username_substring = true
      enable_password_policy = true
    }
  }
  deletion_protection =  "%{deletion_protection}"
}
`, context)
}

var testGoogleSqlDatabaseInstance_basic2 = `
resource "google_sql_database_instance" "instance" {
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_basic3 = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_basic3_update = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_basic_mssql = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  database_version    = "SQLSERVER_2019_STANDARD"
  root_password       = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    collation = "Polish_CI_AS"
  }
}
`

var testGoogleSqlDatabaseInstance_update_mssql = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  database_version    = "SQLSERVER_2019_STANDARD"
  root_password       = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    collation = "Polish_CI_AS"
    ip_configuration {
      ipv4_enabled = true
      require_ssl = true
    }
  }
}
`

func testGoogleSqlDatabaseInstance_ActiveDirectoryConfig(databaseName, networkName, addressRangeName, rootPassword, adDomainName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance-with-ad" {
  depends_on = [google_service_networking_connection.foobar]
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-2-7680"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }

    active_directory_config {
      domain = "%s"
    }
  }
}`, networkName, addressRangeName, databaseName, rootPassword, adDomainName)
}

func testGoogleSqlDatabaseInstance_DenyMaintenancePeriodConfig(databaseName, endDate, startDate, time string) string {
	return fmt.Sprintf(`

resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-custom-4-26624"
    deny_maintenance_period {
      end_date     	= "%s"
      start_date	= "%s"
      time 		= "%s"
    }
  }
}`, databaseName, endDate, startDate, time)
}

func testGoogleSqlDatabaseInstance_EditionConfig(databaseName, tier, edition string) string {
	return fmt.Sprintf(`

resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-east1"
  database_version    = "POSTGRES_14"
  deletion_protection = false
  settings {
    tier = "%s"
    edition = "%s"
  }
}`, databaseName, tier, edition)
}

func testGoogleSqlDatabaseInstance_sqlMysqlDataCacheConfig(instanceName string) string {
	return fmt.Sprintf(`

resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-east1"
  database_version    = "MYSQL_8_0_31"
  deletion_protection = false
  settings {
    tier = "db-perf-optimized-N-2"
    edition = "ENTERPRISE_PLUS"
    data_cache_config {
        data_cache_enabled = true
    }
  }
}`, instanceName)
}

func testGoogleSqlDatabaseInstance_SqlServerAuditConfig(databaseName, rootPassword, bucketName, uploadInterval, retentionInterval string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "gs-bucket" {
  name                      	= "%s"
  location                  	= "US"
  uniform_bucket_level_access = true
}

resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    ip_configuration {
      ipv4_enabled       = "true"
    }
    sql_server_audit_config {
      bucket = "gs://%s"
      retention_interval = "%s"
      upload_interval = "%s"
    }
  }
}
`, bucketName, databaseName, rootPassword, bucketName, retentionInterval, uploadInterval)
}

func testGoogleSqlDatabaseInstance_SqlServerAuditOptionalBucket(databaseName, rootPassword, uploadInterval, retentionInterval string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    sql_server_audit_config {
      retention_interval = "%s"
      upload_interval = "%s"
    }
  }
}
`, databaseName, rootPassword, retentionInterval, uploadInterval)
}

func testGoogleSqlDatabaseInstance_NullSmt(databaseName, rootPassword string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-8-53248"
    ip_configuration {
      ipv4_enabled       = "true"
    }
    advanced_machine_features {
    }
  }
}
`, databaseName, rootPassword)
}

func testGoogleSqlDatabaseInstance_Smt(databaseName, rootPassword string, threadsPerCore int) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-8-53248"
    ip_configuration {
      ipv4_enabled       = "true"
    }
    advanced_machine_features {
      threads_per_core = "%d"
    }
  }
}
`, databaseName, rootPassword, threadsPerCore)
}

func testGoogleSqlDatabaseInstance_Timezone(databaseName, rootPassword, timezone string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "SQLSERVER_2017_STANDARD"
  root_password    = "%s"
  deletion_protection = false
  settings {
    tier = "db-custom-1-3840"
    ip_configuration {
      ipv4_enabled       = "true"
    }
    time_zone = "%s"
  }
}
`, databaseName, rootPassword, timezone)
}

func testGoogleSqlDatabaseInstanceConfig_withoutReplica(instanceName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}
`, instanceName)
}

func testGoogleSqlDatabaseInstanceConfig_withReplica(instanceName, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name                 = "%s"
  region               = "us-central1"
  database_version     = "MYSQL_5_7"
  master_instance_name = google_sql_database_instance.instance.name
  deletion_protection  = false

  replica_configuration {
    failover_target = "true"
  }

  settings {
    tier = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

func googleSqlDatabaseInstance_replicaPromote(instanceName, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name                 = "%s"
  region               = "us-central1"
  database_version     = "MYSQL_5_7"
  deletion_protection  = false

  instance_type = "CLOUD_SQL_INSTANCE"
  settings {
    tier = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

func googleSqlDatabaseInstance_replicaPromoteWithMasterInstanceName(instanceName, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name                 = "%s"
  region               = "us-central1"
  master_instance_name = google_sql_database_instance.instance.name
  database_version     = "MYSQL_5_7"
  deletion_protection  = false
  instance_type = "CLOUD_SQL_INSTANCE"
  settings {
    tier = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

func googleSqlDatabaseInstance_replicaPromoteWithMasterInstanceNameAndReplicaConfiguration(instanceName string, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name                 = "%s"
  region               = "us-central1"
  master_instance_name = google_sql_database_instance.instance.name
  database_version     = "MYSQL_5_7"
  deletion_protection  = false

  replica_configuration {
    failover_target = "true"
  }

  instance_type = "CLOUD_SQL_INSTANCE"
  settings {
    tier = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

func googleSqlDatabaseInstance_replicaPromoteSkippedWithNoMasterInstanceNameAndNoReplicaConfiguration(instanceName string, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name                 = "%s"
  region               = "us-central1"
  database_version     = "MYSQL_5_7"
  deletion_protection  = false
  settings {
    tier = "db-n1-standard-1"
  }
  depends_on = [google_sql_database_instance.instance]
}
`, instanceName, failoverName)
}

func googleSqlDatabaseInstance_replicaPromoteWithReplicaConfiguration(instanceName string, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      binary_log_enabled = "true"
      enabled            = "true"
      start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name                 = "%s"
  region               = "us-central1"
  database_version     = "MYSQL_5_7"
  deletion_protection  = false

  replica_configuration {
    failover_target = "true"
  }

  instance_type = "CLOUD_SQL_INSTANCE"
  settings {
    tier = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, addressRangeName string, specifyPrivatePathOption bool, enablePrivatePath bool) string {
	privatePathOption := ""
	if specifyPrivatePathOption {
		privatePathOption = fmt.Sprintf("enable_private_path_for_google_cloud_services = %t", enablePrivatePath)
	}

	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      %s
    }
  }
}
`, networkName, addressRangeName, databaseName, privatePathOption)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      allocated_ip_range = google_compute_global_address.foobar.name
    }
  }
}
`, networkName, addressRangeName, databaseName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeReplica(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}
resource "google_sql_database_instance" "replica1" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s-replica1"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      allocated_ip_range = google_compute_global_address.foobar.name
    }
  }

  master_instance_name = google_sql_database_instance.instance.name

  replica_configuration {
    connect_retry_interval    = 100
    master_heartbeat_period   = 10000
    password                  = "password"
    username                  = "username"
    ssl_cipher                = "ALL"
    verify_server_certificate = false
  }
}
`, networkName, addressRangeName, databaseName, databaseName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "clone1" {
  name                = "%s-clone1"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  clone {
    source_instance_name = google_sql_database_instance.instance.name
    allocated_ip_range   = google_compute_global_address.foobar.name
  }

}
`, networkName, addressRangeName, databaseName, databaseName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone_withSettings(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}

resource "google_sql_database_instance" "instance" {
  depends_on = [google_service_networking_connection.foobar]
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
    }
    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "clone1" {
  name                = "%s-clone1"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  clone {
    source_instance_name = google_sql_database_instance.instance.name
    allocated_ip_range   = google_compute_global_address.foobar.name
  }

  settings {
    tier = "db-f1-micro"
    backup_configuration {
      enabled = false
    }
  }
}
`, networkName, addressRangeName, databaseName, databaseName)
}

var testGoogleSqlDatabaseInstance_settings = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"
    location_preference {
      zone = "us-central1-f"
    }

    ip_configuration {
      ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }

    backup_configuration {
      enabled    = "true"
      start_time = "19:19"
    }

    activation_policy = "ALWAYS"
  }
}
`

var testGoogleSqlDatabaseInstance_settings_secondary = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"
    availability_type      = "REGIONAL"
    location_preference {
      zone           = "us-central1-f"
      secondary_zone = "us-central1-a"
    }

    ip_configuration {
      ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }

    backup_configuration {
      enabled    = "true"
      start_time = "19:19"
      binary_log_enabled = true
    }

    activation_policy = "ALWAYS"
    connector_enforcement = "REQUIRED"
  }
}
`

var testGoogleSqlDatabaseInstance_settings_deletionProtection = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = %s
  settings {
    tier                   = "db-f1-micro"
    location_preference {
      zone = "us-central1-f"
	}

    ip_configuration {
	  ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }

    backup_configuration {
      enabled    = "true"
      start_time = "19:19"
    }

    activation_policy = "ALWAYS"
  }
}
`

var testGoogleSqlDatabaseInstance_maintenanceVersionWithOldVersion = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  maintenance_version = "MYSQL_5_7_37.R20210508.01_03"
  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_settings_deletionProtectionEnabled = `
resource "google_sql_database_instance" "instance" {
  name                        = "%s"
  region                      = "us-central1"
  database_version            = "MYSQL_5_7"
  deletion_protection         = false
  settings {
	deletion_protection_enabled = %s
    tier                        = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_settings_checkServiceNetworking = `
resource "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = "false"
      private_network = google_compute_network.servicenet.self_link
    }
  }
}
`

var testGoogleSqlDatabaseInstance_replica = `
resource "google_sql_database_instance" "instance_master" {
  name                = "tf-test-%d"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "replica1" {
  name                = "tf-test-%d-1"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"
    backup_configuration {
      enabled = false
      binary_log_enabled = true
    }
  }

  master_instance_name = google_sql_database_instance.instance_master.name

  replica_configuration {
    connect_retry_interval    = 100
    master_heartbeat_period   = 10000
    password                  = "password"
    username                  = "username"
    ssl_cipher                = "ALL"
    verify_server_certificate = false
  }
}

resource "google_sql_database_instance" "replica2" {
  name                = "tf-test-%d-2"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false

  settings {
    tier = "db-n1-standard-1"
    backup_configuration {
      enabled = %s
    }
  }

  master_instance_name = google_sql_database_instance.instance_master.name

  replica_configuration {
    connect_retry_interval    = 100
    master_heartbeat_period   = 10000
    password                  = "password"
    username                  = "username"
    ssl_cipher                = "ALL"
    verify_server_certificate = false
  }
}
`

var testGoogleSqlDatabaseInstance_slave = `
resource "google_sql_database_instance" "instance_master" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "instance_slave" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  master_instance_name = google_sql_database_instance.instance_master.name

  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_highAvailability = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    availability_type = "REGIONAL"

    backup_configuration {
      enabled  = true
      location = "us"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_diskspecs = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier                  = "db-f1-micro"
    disk_autoresize       = true
    disk_autoresize_limit = 50
    disk_size             = 15
    disk_type             = "PD_HDD"
  }
}
`

var testGoogleSqlDatabaseInstance_maintenance = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    maintenance_window {
      day          = 7
      hour         = 3
      update_track = "canary"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_authNets_step1 = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"

    ip_configuration {
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2037-11-15T16:19:00.094Z"
      }
    }
  }
}
`

var testGoogleSqlDatabaseInstance_authNets_step2 = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"

    ip_configuration {
      ipv4_enabled = "true"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_multipleOperations = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-f1-micro"
  }
}

resource "google_sql_database" "database" {
  name     = "tf-test-%s"
  instance = google_sql_database_instance.instance.name
}

resource "google_sql_user" "user" {
  name     = "tf-test-%s"
  instance = google_sql_database_instance.instance.name
  host     = "google.com"
  password = "hunter2"
}
`

var testGoogleSqlDatabaseInstance_basic_with_user_labels = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    user_labels = {
      track    = "production"
      location = "western-division"
    }
  }
}
`
var testGoogleSqlDatabaseInstance_basic_with_user_labels_update = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    user_labels = {
      track = "production"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_insights = `
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"

    insights_config {
      query_insights_enabled  = true
      query_string_length     = 256
      record_application_tags = true
      record_client_address   = true
      query_plans_per_minute  = 10
    }
  }
}
`
var testGoogleSqlDatabaseInstance_encryptionKey = `
data "google_project" "project" {
  project_id = "%{project_id}"
}
resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
  "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com",
  ]
}

resource "google_sql_database_instance" "master" {
  name                = "%{instance_name}-master"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false
  encryption_key_name = google_kms_crypto_key.key.id

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "replica" {
  name                 = "%{instance_name}-replica"
  database_version     = "MYSQL_5_7"
  region               = "us-central1"
  master_instance_name = google_sql_database_instance.master.name
  deletion_protection  = false

  settings {
    tier = "db-n1-standard-1"
  }

  depends_on = [google_sql_database_instance.master]
}
`

var testGoogleSqlDatabaseInstance_encryptionKey_replicaInDifferentRegion = `

data "google_project" "project" {
  project_id = "%{project_id}"
}

resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "key" {

  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com",
  ]
}

resource "google_sql_database_instance" "master" {
  name                = "%{instance_name}-master"
  database_version    = "MYSQL_5_7"
  region              = "us-central1"
  deletion_protection = false
  encryption_key_name = google_kms_crypto_key.key.id

  settings {
    tier = "db-n1-standard-1"

    backup_configuration {
      enabled            = true
      start_time         = "00:00"
      binary_log_enabled = true
    }
  }
}

resource "google_kms_key_ring" "keyring-rep" {

  name     = "%{key_name}-rep"
  location = "us-east1"
}

resource "google_kms_crypto_key" "key-rep" {

  name     = "%{key_name}-rep"
  key_ring = google_kms_key_ring.keyring-rep.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_rep" {
  crypto_key_id = google_kms_crypto_key.key-rep.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com",
  ]
}

resource "google_sql_database_instance" "replica" {
  name                 = "%{instance_name}-replica"
  database_version     = "MYSQL_5_7"
  region               = "us-east1"
  master_instance_name = google_sql_database_instance.master.name
  encryption_key_name = google_kms_crypto_key.key-rep.id
  deletion_protection  = false

  settings {
    tier = "db-n1-standard-1"
  }

  depends_on = [google_sql_database_instance.master]
}
`

func testGoogleSqlDatabaseInstance_PointInTimeRecoveryEnabled(masterID int, pointInTimeRecoveryEnabled bool, dbVersion string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "%s"
  deletion_protection = false
  root_password		  = "rand-pwd-%d"
  settings {
    tier = "db-custom-2-13312"
    backup_configuration {
      enabled                        = true
      start_time                     = "00:00"
      point_in_time_recovery_enabled = %t
    }
  }
}
`, masterID, dbVersion, masterID, pointInTimeRecoveryEnabled)
}

func testGoogleSqlDatabaseInstance_BackupRetention(masterID int) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    backup_configuration {
      enabled                        = true
      start_time                     = "00:00"
      binary_log_enabled             = true
	  transaction_log_retention_days = 2
	  backup_retention_settings {
	    retained_backups = 4
	  }
    }
  }
}
`, masterID)
}

func testAccSqlDatabaseInstance_beforeBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = "false"
	}
  }

  deletion_protection = false
}
`, context)
}

func testAccSqlDatabaseInstance_restoreFromBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = "false"
	}
  }

  restore_backup_context {
    backup_run_id = data.google_sql_backup_run.backup.backup_id
    instance_id = data.google_sql_backup_run.backup.instance
  }

  // Ignore changes, since the most recent backup may change during the test
  lifecycle{
	ignore_changes = [restore_backup_context[0].backup_run_id]
  }

  deletion_protection = false
}

data "google_sql_backup_run" "backup" {
	instance = "%{original_db_name}"
	most_recent = true
}
`, context)
}

func testAccSqlDatabaseInstance_basicClone(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  clone {
    source_instance_name = data.google_sql_backup_run.backup.instance
    point_in_time = data.google_sql_backup_run.backup.start_time
  }

  deletion_protection = false

  // Ignore changes, since the most recent backup may change during the test
  lifecycle{
	ignore_changes = [clone[0].point_in_time]
  }
}

data "google_sql_backup_run" "backup" {
	instance = "%{original_db_name}"
	most_recent = true
}
`, context)
}

func testAccSqlDatabaseInstance_cloneWithSettings(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = false
	}
  }

  clone {
    source_instance_name = data.google_sql_backup_run.backup.instance
    point_in_time = data.google_sql_backup_run.backup.start_time
  }

  deletion_protection = false

  // Ignore changes, since the most recent backup may change during the test
  lifecycle{
	ignore_changes = [clone[0].point_in_time]
  }
}

data "google_sql_backup_run" "backup" {
	instance = "%{original_db_name}"
	most_recent = true
}
`, context)
}

func testAccSqlDatabaseInstance_cloneWithDatabaseNames(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  clone {
    source_instance_name = data.google_sql_backup_run.backup.instance
    point_in_time = data.google_sql_backup_run.backup.start_time
    database_names = ["userdb1"]
  }

  deletion_protection = false

  // Ignore changes, since the most recent backup may change during the test
  lifecycle{
	ignore_changes = [clone[0].point_in_time]
  }
}

data "google_sql_backup_run" "backup" {
	instance = "%{original_db_name}"
	most_recent = true
}
`, context)
}

func checkPromoteReplicaSkipConfigurations(resourceName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Can't find %s in state", resourceName)
		}

		resourceAttributes := resource.Primary.Attributes
		instanceType, ok := resourceAttributes["instance_type"]
		if !ok {
			return fmt.Errorf("Instance type is not present in state for %s", resourceName)
		}
		if instanceType != "READ_REPLICA_INSTANCE" {
			return fmt.Errorf("instance_type is %s, it should be READ_REPLICA_INSTANCE.", instanceType)
		}

		masterInstanceName, ok := resourceAttributes["master_instance_name"]
		if !ok && masterInstanceName != "" {
			return fmt.Errorf("master_instance_name should be present in %s state.", resourceName)
		}

		return nil
	}
}

func checkPromoteReplicaConfigurations(resourceName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Can't find %s in state", resourceName)
		}

		resourceAttributes := resource.Primary.Attributes
		instanceType, ok := resourceAttributes["instance_type"]
		if !ok {
			return fmt.Errorf("Instance type is not present in state for %s", resourceName)
		}
		if instanceType != "CLOUD_SQL_INSTANCE" {
			return fmt.Errorf("Error in replica promotion. instance_type is %s, it should be CLOUD_SQL_INSTANCE.", instanceType)
		}

		masterInstanceName, ok := resourceAttributes["master_instance_name"]
		if ok && masterInstanceName != "" {
			return fmt.Errorf("Error in replica promotion. master_instance_name should not be present in %s state.", resourceName)
		}

		replicaConfiguration, ok := resourceAttributes["replica_configuration"]
		if ok && replicaConfiguration != "" {
			return fmt.Errorf("Error in replica promotion. replica_configuration should not be present in %s state.", resourceName)
		}

		return nil
	}
}

func checkInstanceTypeIsPresent(resourceName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}
		rsAttr := rs.Primary.Attributes
		_, ok = rsAttr["instance_type"]
		if !ok {
			return fmt.Errorf("Instance type is not computed for %s", resourceName)
		}
		return nil
	}
}

func testGoogleSqlDatabaseInstance_readReplica(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "master" {
  region           = "asia-northeast1"
  name             = "%s-master"
  database_version = "MYSQL_5_7"
  deletion_protection  = false
  settings {
    availability_type = "ZONAL"
    disk_autoresize   = true
    disk_size         = 10
    disk_type         = "PD_SSD"
    tier              = "db-f1-micro"

    activation_policy = "ALWAYS"
    pricing_plan      = "PER_USE"

    backup_configuration {
      binary_log_enabled = true
      enabled            = true
      location           = "asia"
      start_time         = "18:00"
    }

    database_flags {
      name  = "character_set_server"
      value = "utf8mb4"
    }
  }
}


resource "google_sql_database_instance" "replica" {
  depends_on           = [google_sql_database_instance.master]
  name                 = "%s-replica"
  master_instance_name = google_sql_database_instance.master.name
  region               = "asia-northeast1"
  database_version     = "MYSQL_5_7"
  deletion_protection  = false
  replica_configuration {
    failover_target = false
  }
  settings {
    tier              = "db-f1-micro"
    availability_type = "ZONAL"
    pricing_plan      = "PER_USE"
    disk_autoresize   = true
    disk_size         = 10

    backup_configuration {
      binary_log_enabled = true
    }

    database_flags {
      name  = "slave_parallel_workers"
      value = "3"
    }

    database_flags {
      name  = "slave_parallel_type"
      value = "LOGICAL_CLOCK"
    }

    database_flags {
      name  = "slave_pending_jobs_size_max"
      value = "536870912" # 512MB
    }
  }
}`, instance, instance)
}

func testGoogleSqlDatabaseInstance_updateReadReplica(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "master" {
  region           = "asia-northeast1"
  name             = "%s-master"
  database_version = "MYSQL_5_7"
  deletion_protection  = false
  settings {
    availability_type = "ZONAL"
    disk_autoresize   = true
    disk_size         = 10
    disk_type         = "PD_SSD"
    tier              = "db-f1-micro"

    activation_policy = "ALWAYS"
    pricing_plan      = "PER_USE"

    backup_configuration {
      binary_log_enabled = true
      enabled            = true
      location           = "asia"
      start_time         = "18:00"
    }

    database_flags {
      name  = "character_set_server"
      value = "utf8mb4"
    }
  }
}


resource "google_sql_database_instance" "replica" {
  depends_on           = [google_sql_database_instance.master]
  name                 = "%s-replica"
  master_instance_name = google_sql_database_instance.master.name
  region               = "asia-northeast1"
  database_version     = "MYSQL_5_7"
  deletion_protection  = false
  replica_configuration {
    failover_target = false
  }
  settings {
    tier              = "db-f1-micro"
    availability_type = "ZONAL"
    pricing_plan      = "PER_USE"
    disk_autoresize   = true
    disk_size         = 10

    backup_configuration {
      binary_log_enabled = true
    }

    database_flags {
      name  = "slave_parallel_workers"
      value = "2"
    }

    database_flags {
      name  = "slave_parallel_type"
      value = "LOGICAL_CLOCK"
    }

    database_flags {
      name  = "slave_pending_jobs_size_max"
      value = "536870912" # 512MB
    }
  }
}`, instance, instance)
}

func testGoogleSqlDatabaseInstance_updateRootPassword(instance, databaseVersion, rootPassword string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "main" {
    name             = "%s"
	database_version = "%s"
	region           = "us-central1"
	deletion_protection = false
	root_password = "%s"
	settings {
		tier = "db-custom-2-13312"
	}
}`, instance, databaseVersion, rootPassword)
}

func testGoogleSqlDatabaseInstance_activationPolicy(instance, databaseVersion, activationPolicy string, deletionProtection bool) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "%s"
  deletion_protection = %t
  settings {
    tier              = "db-f1-micro"
    activation_policy = "%s"
  }
}
`, instance, databaseVersion, deletionProtection, activationPolicy)
}
