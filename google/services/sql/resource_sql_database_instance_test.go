// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/sql"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// Fields that should be ignored in import tests because they aren't returned
// from GCP (and thus can't be imported)
var ignoredReplicaConfigurationFields = []string{
	"deletion_protection",
	"root_password",
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
	"replica_configuration.0.failover_target",
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

	databaseName := "tf-test-sql-instance-" + acctest.RandString(t, 10)
	failoverName := "tf-test-sql-instance-failover-" + acctest.RandString(t, 10)
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

	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	testId := "sql-instance-clone-2"
	addressName := acctest.BootstrapSharedTestGlobalAddress(t, testId)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, testId)

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
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, false, false),
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
				ExpectError: regexp.MustCompile("Error, failed to create instance tf-test-\\d+-2: googleapi: Error 400: Invalid request: Invalid flag for instance role: Backups cannot be enabled for read replica instance"),
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

func TestAccSqlDatabaseInstance_maintenance_update_track_week5(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenance_week5, masterID),
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
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, "sql-instance-1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, false, false),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, true, false),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, true, true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName, true, false),
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

func TestAccSqlDatabaseInstance_withPSCEnabled_withoutAllowedConsumerProjects(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := "psctestproject" + acctest.RandString(t, 10)
	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withoutAllowedConsumerProjects(instanceName, projectId, orgId, billingAccount),
				Check:  resource.ComposeTestCheckFunc(verifyPscOperation("google_sql_database_instance.instance", true, true, nil)),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPSCEnabled_withEmptyAllowedConsumerProjects(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := "psctestproject" + acctest.RandString(t, 10)
	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withEmptyAllowedConsumerProjects(instanceName, projectId, orgId, billingAccount),
				Check:  resource.ComposeTestCheckFunc(verifyPscOperation("google_sql_database_instance.instance", true, true, []string{})),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPSCEnabled_withAllowedConsumerProjects(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := "psctestproject" + acctest.RandString(t, 10)
	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withAllowedConsumerProjects(instanceName, projectId, orgId, billingAccount),
				Check:  resource.ComposeTestCheckFunc(verifyPscOperation("google_sql_database_instance.instance", true, true, []string{envvar.GetTestProjectFromEnv()})),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPSCEnabled_thenAddAllowedConsumerProjects_thenRemoveAllowedConsumerProject(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := "psctestproject" + acctest.RandString(t, 10)
	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withoutAllowedConsumerProjects(instanceName, projectId, orgId, billingAccount),
				Check:  resource.ComposeTestCheckFunc(verifyPscOperation("google_sql_database_instance.instance", true, true, nil)),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withAllowedConsumerProjects(instanceName, projectId, orgId, billingAccount),
				Check:  resource.ComposeTestCheckFunc(verifyPscOperation("google_sql_database_instance.instance", true, true, []string{envvar.GetTestProjectFromEnv()})),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withoutAllowedConsumerProjects(instanceName, projectId, orgId, billingAccount),
				Check:  resource.ComposeTestCheckFunc(verifyPscOperation("google_sql_database_instance.instance", true, true, []string{})),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPSCEnabled_withoutPscAutoConnections(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withoutPscAutoConnections(instanceName),
				Check:  resource.ComposeTestCheckFunc(verifyPscAutoConnectionsOperation("google_sql_database_instance.instance", true, true, false, "", "")),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPSCEnabled_withPscAutoConnections(t *testing.T) {
	t.Parallel()

	testId := "test-psc-auto-con" + acctest.RandString(t, 10)
	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := envvar.GetTestProjectFromEnv()
	networkName := acctest.BootstrapSharedTestNetwork(t, testId)
	network_short_link_name := fmt.Sprintf("projects/%s/global/networks/%s", projectId, networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withPscAutoConnections(instanceName, projectId, networkName),
				Check:  resource.ComposeTestCheckFunc(verifyPscAutoConnectionsOperation("google_sql_database_instance.instance", true, true, true, network_short_link_name, projectId)),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPSCEnabled_thenAddPscAutoConnections_thenRemovePscAutoConnections(t *testing.T) {
	t.Parallel()

	testId := "test-psc-auto-con" + acctest.RandString(t, 10)
	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := envvar.GetTestProjectFromEnv()
	networkName := acctest.BootstrapSharedTestNetwork(t, testId)
	network_short_link_name := fmt.Sprintf("projects/%s/global/networks/%s", projectId, networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withoutPscAutoConnections(instanceName),
				Check:  resource.ComposeTestCheckFunc(verifyPscAutoConnectionsOperation("google_sql_database_instance.instance", true, true, false, "", "")),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withPscAutoConnections(instanceName, projectId, networkName),
				Check:  resource.ComposeTestCheckFunc(verifyPscAutoConnectionsOperation("google_sql_database_instance.instance", true, true, true, network_short_link_name, projectId)),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/", projectId),
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccSqlDatabaseInstance_withPSCEnabled_withoutPscAutoConnections(instanceName),
				Check:  resource.ComposeTestCheckFunc(verifyPscAutoConnectionsOperation("google_sql_database_instance.instance", true, true, false, "", "")),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPSCEnabled_withIpV4Enabled(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	projectId := "psctestproject" + acctest.RandString(t, 10)
	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSqlDatabaseInstance_withPSCEnabled_withIpV4Enable(instanceName, projectId, orgId, billingAccount),
				ExpectError: regexp.MustCompile("PSC connectivity cannot be enabled together with only public IP"),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(t *testing.T) {

	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	testId := "sql-instance-allocated-1"
	addressName := acctest.BootstrapSharedTestGlobalAddress(t, testId)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, testId)

	updateTestId := "sql-instance-allocated-update-1"
	addressName_update := acctest.BootstrapSharedTestGlobalAddress(t, updateTestId)
	networkName_update := acctest.BootstrapSharedServiceNetworkingConnection(t, updateTestId)

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

	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)

	testId := "sql-instance-replica-1"
	addressName := acctest.BootstrapSharedTestGlobalAddress(t, testId)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, testId)

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

	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	testId := "sql-instance-clone-1"
	addressName := acctest.BootstrapSharedTestGlobalAddress(t, testId)
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, testId)

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

func TestAccSqlDatabaseInstance_EnableGoogleMlIntegration(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_EnableGoogleMlIntegration(masterID, true, "POSTGRES_14", "db-custom-2-13312"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			// Test that updates to other settings work after google-ml-integration is enabled
			{
				Config: testGoogleSqlDatabaseInstance_EnableGoogleMlIntegration(masterID, true, "POSTGRES_14", "db-custom-2-10240"),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_EnableGoogleMlIntegration(masterID, false, "POSTGRES_14", "db-custom-2-10240"),
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

func TestAccSqlDatabaseInstance_EnableGoogleDataplexIntegration(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_EnableDataplexIntegration(masterID, true),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_EnableDataplexIntegration(masterID, false),
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
	networkName := acctest.BootstrapSharedServiceNetworkingConnection(t, "sql-instance-ad-1")
	rootPassword := acctest.RandString(t, 15)
	adDomainName := acctest.BootstrapSharedTestADDomain(t, "test-domain", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_ActiveDirectoryConfig(databaseName, networkName, rootPassword, adDomainName),
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

func TestAccSQLDatabaseInstance_DefaultEdition(t *testing.T) {
	t.Parallel()
	databaseName := "tf-test-" + acctest.RandString(t, 10)
	databaseVersion := "POSTGRES_16"
	enterprisePlusTier := "db-perf-optimized-N-2"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_DefaultEdition(databaseName, databaseVersion, enterprisePlusTier),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE_PLUS"),
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

func TestAccSqlDatabaseInstance_Edition(t *testing.T) {
	t.Parallel()
	enterprisePlusName := "tf-test-enterprise-plus" + acctest.RandString(t, 10)
	enterprisePlusTier := "db-perf-optimized-N-2"
	enterpriseName := "tf-test-enterprise-" + acctest.RandString(t, 10)
	enterpriseTier := "db-custom-2-13312"
	noEditionName := "tf-test-enterprise-noedition-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig_noEdition(noEditionName, enterpriseTier),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE"),
				),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			// Delete and recreate (ForceNew) triggered by passing in a new `name` value
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(enterprisePlusName, enterprisePlusTier, "ENTERPRISE_PLUS"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE_PLUS"),
				),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			// Delete and recreate (ForceNew) triggered by passing in a new `name` value
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(enterpriseName, enterpriseTier, "ENTERPRISE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE"),
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.data_cache_config.0.data_cache_enabled", "true"),
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

func TestAccSQLDatabaseInstance_sqlPostgresDataCacheConfig(t *testing.T) {
	t.Parallel()
	enterprisePlusInstanceName := "tf-test-enterprise-plus" + acctest.RandString(t, 10)
	enterprisePlusTier := "db-perf-optimized-N-2"
	enterpriseInstanceName := "tf-test-enterprise-" + acctest.RandString(t, 10)
	enterpriseTier := "db-custom-2-13312"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_sqlPostgresDataCacheConfig(enterprisePlusInstanceName, enterprisePlusTier, "ENTERPRISE_PLUS"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.data_cache_config.0.data_cache_enabled", "true"),
				),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_sqlPostgresDataCacheConfig(enterpriseInstanceName, enterpriseTier, "ENTERPRISE"),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf("Error, failed to create instance %s: googleapi: Error 400: Invalid request: Only ENTERPRISE PLUS edition supports data cache", enterpriseInstanceName)),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_Mysql_Edition_Upgrade(t *testing.T) {
	t.Parallel()
	enterpriseTier := "db-custom-2-13312"
	editionUpgrade := "tf-test-enterprise-upgrade-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_sqlMysql(editionUpgrade, enterpriseTier),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE"),
				),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_sqlMysqlDataCacheConfig(editionUpgrade),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE_PLUS"),
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.data_cache_config.0.data_cache_enabled", "true"),
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

func TestAccSqlDatabaseInstance_Postgres_Edition_Upgrade(t *testing.T) {
	t.Parallel()
	enterpriseTier := "db-custom-2-13312"
	enterprisePlusTier := "db-perf-optimized-N-2"
	editionUpgrade := "tf-test-enterprise-upgrade-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(editionUpgrade, enterpriseTier, "ENTERPRISE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE"),
				),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(editionUpgrade, enterprisePlusTier, "ENTERPRISE_PLUS"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE_PLUS"),
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

func TestAccSqlDatabaseInstance_Edition_Downgrade(t *testing.T) {
	t.Parallel()
	enterprisePlusTier := "db-perf-optimized-N-2"
	enterpriseTier := "db-custom-2-13312"
	editionDowngrade := "tf-test-enterprise-downgrade-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(editionDowngrade, enterprisePlusTier, "ENTERPRISE_PLUS"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE_PLUS"),
				),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_EditionConfig(editionDowngrade, enterpriseTier, "ENTERPRISE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_database_instance.instance", "settings.0.edition", "ENTERPRISE"),
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

func TestAccSqlDatabaseInstance_updateDifferentFlagOrder(t *testing.T) {
	t.Parallel()

	instance := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_flags(instance),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config:             testGoogleSqlDatabaseInstance_flags_update(instance),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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

func TestAccSqlDatabaseInstance_SqlServerTimezoneUpdate(t *testing.T) {
	t.Parallel()

	instanceName := "tf-test-" + acctest.RandString(t, 10)
	rootPassword := acctest.RandString(t, 15)
	timezone := "Eastern Standard Time"
	timezoneUpdate := "Pacific Standard Time"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_SqlServerTimezone(instanceName, rootPassword, timezone),
			},
			{
				ResourceName:            "google_sql_database_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "root_password"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_SqlServerTimezone(instanceName, rootPassword, timezoneUpdate),
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

	databaseName := "tf-test-sql-instance-" + acctest.RandString(t, 10)
	failoverName := "tf-test-sql-instance-failover-" + acctest.RandString(t, 10)
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
	databaseName := "tf-test-sql-instance-" + acctest.RandString(t, 10)
	failoverName := "tf-test-sql-instance-failover-" + acctest.RandString(t, 10)
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

	databaseName := "tf-test-sql-instance-" + acctest.RandString(t, 10)
	failoverName := "tf-test-sql-instance-failover-" + acctest.RandString(t, 10)
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

	databaseName := "tf-test-sql-instance-" + acctest.RandString(t, 10)
	failoverName := "tf-test-sql-instance-failover-" + acctest.RandString(t, 10)
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

	databaseName := "tf-test-sql-instance-" + acctest.RandString(t, 10)
	failoverName := "tf-test-sql-instance-failover-" + acctest.RandString(t, 10)
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

// Switchover between primary and cascadable replica sunny case
func TestAccSqlDatabaseInstance_SwitchoverSuccess(t *testing.T) {
	t.Parallel()
	primaryName := "tf-test-sql-instance-" + acctest.RandString(t, 10)
	replicaName := "tf-test-sql-instance-replica-" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_SqlServerwithCascadableReplica(primaryName, replicaName),
			},
			{
				ResourceName:            "google_sql_database_instance.original-primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
			{
				ResourceName:            "google_sql_database_instance.original-replica",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
			{
				// Split into two configs because current TestStep implementation checks diff before refreshing.
				Config: googleSqlDatabaseInstance_switchoverOnReplica(primaryName, replicaName),
			},
			{
				Config: googleSqlDatabaseInstance_updatePrimaryAfterSwitchover(primaryName, replicaName),
			},
			{
				RefreshState: true,
				Check:        resource.ComposeTestCheckFunc(resource.TestCheckTypeSetElemAttr("google_sql_database_instance.original-replica", "replica_names.*", primaryName), checkSwitchoverOriginalReplicaConfigurations("google_sql_database_instance.original-replica"), checkSwitchoverOriginalPrimaryConfigurations("google_sql_database_instance.original-primary", replicaName)),
			},
			{
				ResourceName:            "google_sql_database_instance.original-primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoredReplicaConfigurationFields,
			},
			{
				ResourceName:      "google_sql_database_instance.original-replica",
				ImportState:       true,
				ImportStateVerify: true,
				// original-replica is no longer a replica, but replica_configuration is O + C and cannot be unset
				ImportStateVerifyIgnore: []string{"replica_configuration", "deletion_protection", "root_password"},
			},
			{
				// Delete replica first so PostTestDestroy doesn't fail when deleting instances which have replicas. We've already validated switchover behavior, the remaining steps are cleanup
				Config: googleSqlDatabaseInstance_deleteReplicasAfterSwitchover(primaryName, replicaName),
				// We delete replica, but haven't updated the master's replica_names
				ExpectNonEmptyPlan: true,
			},
			{
				// Remove replica from primary's resource
				Config: googleSqlDatabaseInstance_removeReplicaFromPrimaryAfterSwitchover(replicaName),
			},
		},
	})
}

func TestAccSqlDatabaseInstance_updateSslOptionsForPostgreSQL(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	databaseVersion := "POSTGRES_14"
	resourceName := "google_sql_database_instance.instance"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),

		// We don't do ImportStateVerify for the ssl_mode because of the implementation. The ssl_mode is expected to be discarded if the local state doesn't have it.
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_setSslOptionsForPostgreSQL(databaseName, databaseVersion, "ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "settings.0.ip_configuration.0.ssl_mode", "ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "settings.0.ip_configuration.0.ssl_mode"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_setSslOptionsForPostgreSQL(databaseName, databaseVersion, "ENCRYPTED_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "settings.0.ip_configuration.0.ssl_mode", "ENCRYPTED_ONLY"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "settings.0.ip_configuration.0.ssl_mode"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_setSslOptionsForPostgreSQL(databaseName, databaseVersion, "TRUSTED_CLIENT_CERTIFICATE_REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "settings.0.ip_configuration.0.ssl_mode", "TRUSTED_CLIENT_CERTIFICATE_REQUIRED"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "settings.0.ip_configuration.0.ssl_mode"},
			},
			{
				Config: testGoogleSqlDatabaseInstance_setSslOptionsForPostgreSQL(databaseName, databaseVersion, "ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "settings.0.ip_configuration.0.ssl_mode", "ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "settings.0.ip_configuration.0.ssl_mode"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_useInternalCaByDefault(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	resourceName := "google_sql_database_instance.instance"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testGoogleSqlDatabaseInstance_basic3, databaseName),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(resourceName, "settings.0.ip_configuration.0.server_ca_mode", "GOOGLE_MANAGED_INTERNAL_CA")),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabaseInstance_useCasBasedServerCa(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(t, 10)
	resourceName := "google_sql_database_instance.instance"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),

		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_setCasServerCa(databaseName, "GOOGLE_MANAGED_CAS_CA"),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(resourceName, "settings.0.ip_configuration.0.server_ca_mode", "GOOGLE_MANAGED_CAS_CA")),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func testGoogleSqlDatabaseInstance_setCasServerCa(databaseName, serverCaMode string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "POSTGRES_15"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = "true"
      server_ca_mode  = "%s"
    }
  }
}
`, databaseName, serverCaMode)
}

func testGoogleSqlDatabaseInstance_setSslOptionsForPostgreSQL(databaseName string, databaseVersion string, sslMode string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "%s"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
      ipv4_enabled = true
      ssl_mode = "%s"
    }
  }
}`, databaseName, databaseVersion, sslMode)
}

func testAccSqlDatabaseInstance_sqlMysqlInstancePvpExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "mysql_pvp_instance_name" {
  name             = "tf-test-mysql-pvp-instance-name%{random_suffix}"
  region           = "asia-northeast1"
  database_version = "MYSQL_8_0"
  root_password = "abcABC123!"
  settings {
    tier              = "db-g1-small"
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
    tier = "db-g1-small"
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
    tier = "db-g1-small"
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
    tier = "db-g1-small"
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
      ssl_mode = "ENCRYPTED_ONLY"
    }
  }
}
`

func testGoogleSqlDatabaseInstance_ActiveDirectoryConfig(databaseName, networkName, rootPassword, adDomainName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance-with-ad" {
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
}`, networkName, databaseName, rootPassword, adDomainName)
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

func testGoogleSqlDatabaseInstance_DefaultEdition(databaseName, databaseVersion, tier string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-east1"
  database_version    = "%s"
  deletion_protection = false
  settings {
    tier = "%s"
  }
}`, databaseName, databaseVersion, tier)
}

func testGoogleSqlDatabaseInstance_EditionConfig_noEdition(databaseName, tier string) string {
	return fmt.Sprintf(`

resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-east1"
  database_version    = "POSTGRES_14"
  deletion_protection = false
  settings {
    tier = "%s"
  }
}`, databaseName, tier)
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
	backup_configuration {
	  transaction_log_retention_days = 7
    }
  }
}`, databaseName, tier, edition)
}

func testGoogleSqlDatabaseInstance_sqlMysql(databaseName, tier string) string {
	return fmt.Sprintf(`

resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-east1"
  database_version    = "MYSQL_8_0_31"
  deletion_protection = false
  settings {
    tier = "%s"
  }
}`, databaseName, tier)
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

func testGoogleSqlDatabaseInstance_sqlPostgresDataCacheConfig(instanceName, tier, edition string) string {
	return fmt.Sprintf(`

resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-east1"
  database_version    = "POSTGRES_14"
  deletion_protection = false
  settings {
    tier = "%s"
    edition = "%s"
    data_cache_config {
        data_cache_enabled = true
    }
  }
}`, instanceName, tier, edition)
}

func testGoogleSqlDatabaseInstance_SqlServerTimezone(instance, rootPassword, timezone string) string {
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
`, instance, rootPassword, timezone)
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

// Create SQL server primary with cascadable replica
func testGoogleSqlDatabaseInstanceConfig_SqlServerwithCascadableReplica(primaryName string, replicaName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "original-primary" {
  name                = "%s"
  region              = "us-east1"
  database_version    = "SQLSERVER_2019_ENTERPRISE"
  deletion_protection = false

  root_password = "sqlserver1"
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}

resource "google_sql_database_instance" "original-replica" {
  name                 = "%s"
  region               = "us-west2"
  database_version     = "SQLSERVER_2019_ENTERPRISE"
  master_instance_name = google_sql_database_instance.original-primary.name
  deletion_protection  = false
  root_password = "sqlserver1"
  replica_configuration {
    cascadable_replica = true
  }

  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}
`, primaryName, replicaName)
}

func googleSqlDatabaseInstance_switchoverOnReplica(primaryName string, replicaName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "original-primary" {
  name                = "%s"
  region              = "us-east1"
  database_version    = "SQLSERVER_2019_ENTERPRISE"
  deletion_protection = false

  root_password = "sqlserver1"
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}

resource "google_sql_database_instance" "original-replica" {
  name                 = "%s"
  region               = "us-west2"
  database_version     = "SQLSERVER_2019_ENTERPRISE"
  deletion_protection  = false
  root_password = "sqlserver1"
  instance_type = "CLOUD_SQL_INSTANCE"
  replica_names = [google_sql_database_instance.original-primary.name]
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}
`, primaryName, replicaName)
}

func googleSqlDatabaseInstance_updatePrimaryAfterSwitchover(primaryName string, replicaName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "original-primary" {
  name                = "%s"
  region              = "us-east1"
  database_version    = "SQLSERVER_2019_ENTERPRISE"
  deletion_protection = false
  root_password = "sqlserver1"
  instance_type = "READ_REPLICA_INSTANCE"
  master_instance_name = "%s"
  replica_configuration {
	cascadable_replica = true
  }
  replica_names = []
  settings {
	tier              = "db-perf-optimized-N-2"
	edition           = "ENTERPRISE_PLUS"
  }
}

  resource "google_sql_database_instance" "original-replica" {
  name                 = "%s"
  region               = "us-west2"
  database_version     = "SQLSERVER_2019_ENTERPRISE"
  deletion_protection  = false
  root_password = "sqlserver1"
  instance_type = "CLOUD_SQL_INSTANCE"
  replica_names = [google_sql_database_instance.original-primary.name]
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}
`, primaryName, replicaName, replicaName)
}

// After a switchover, the original-primary is now the replica and must be removed first.
func googleSqlDatabaseInstance_deleteReplicasAfterSwitchover(primaryName, replicaName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "original-replica" {
  name                 = "%s"
  region               = "us-west2"
  database_version     = "SQLSERVER_2019_ENTERPRISE"
  deletion_protection  = false
  root_password = "sqlserver1"
  instance_type = "CLOUD_SQL_INSTANCE"
  replica_names = ["%s"]
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}

`, replicaName, primaryName)
}

// Update original-replica replica_names after deleting original-primary
func googleSqlDatabaseInstance_removeReplicaFromPrimaryAfterSwitchover(replicaName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "original-replica" {
  name                 = "%s"
  region               = "us-west2"
  database_version     = "SQLSERVER_2019_ENTERPRISE"
  deletion_protection  = false
  root_password = "sqlserver1"
  instance_type = "CLOUD_SQL_INSTANCE"
  replica_names = []
  settings {
    tier              = "db-perf-optimized-N-2"
    edition           = "ENTERPRISE_PLUS"
  }
}
`, replicaName)
}

func testAccSqlDatabaseInstance_basicInstanceForPsc(instanceName string, projectId string, orgId string, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "testproject" {
  name                = "%s"
  project_id          = "%s"
  org_id              = "%s"
  billing_account     = "%s"
  deletion_policy = "DELETE"
}

resource "google_sql_database_instance" "instance" {
  project             = google_project.testproject.project_id
  name                = "%s"
  region              = "us-south1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	}
	availability_type = "REGIONAL"
  }
}
`, projectId, projectId, orgId, billingAccount, instanceName)
}

func testAccSqlDatabaseInstance_withPSCEnabled_withIpV4Enable(instanceName string, projectId string, orgId string, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "testproject" {
  name                = "%s"
  project_id          = "%s"
  org_id              = "%s"
  billing_account     = "%s"
  deletion_policy = "DELETE"
}

resource "google_sql_database_instance" "instance" {
  project             = google_project.testproject.project_id
  name                = "%s"
  region              = "us-south1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
		psc_config {
			psc_enabled = true
		}
		ipv4_enabled = true
    }
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	}
	availability_type = "REGIONAL"
  }
}
`, projectId, projectId, orgId, billingAccount, instanceName)
}

func testAccSqlDatabaseInstance_withPSCEnabled_withoutAllowedConsumerProjects(instanceName string, projectId string, orgId string, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "testproject" {
  name                = "%s"
  project_id          = "%s"
  org_id              = "%s"
  billing_account     = "%s"
  deletion_policy = "DELETE"
}

resource "google_sql_database_instance" "instance" {
  project             = google_project.testproject.project_id
  name                = "%s"
  region              = "us-south1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
		psc_config {
			psc_enabled = true
		}
		ipv4_enabled = false
    }
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	  }
	availability_type = "REGIONAL"
  }
}
`, projectId, projectId, orgId, billingAccount, instanceName)
}

func testAccSqlDatabaseInstance_withPSCEnabled_withEmptyAllowedConsumerProjects(instanceName string, projectId string, orgId string, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "testproject" {
  name                = "%s"
  project_id          = "%s"
  org_id              = "%s"
  billing_account     = "%s"
  deletion_policy = "DELETE"
}

resource "google_sql_database_instance" "instance" {
  project             = google_project.testproject.project_id
  name                = "%s"
  region              = "us-south1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
		psc_config {
			psc_enabled = true
			allowed_consumer_projects = []
		}
		ipv4_enabled = false
    }
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	}
	availability_type = "REGIONAL"
  }
}
`, projectId, projectId, orgId, billingAccount, instanceName)
}

func testAccSqlDatabaseInstance_withPSCEnabled_withAllowedConsumerProjects(instanceName string, projectId string, orgId string, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "testproject" {
  name                = "%s"
  project_id          = "%s"
  org_id              = "%s"
  billing_account     = "%s"
  deletion_policy = "DELETE"
}

resource "google_sql_database_instance" "instance" {
  project             = google_project.testproject.project_id
  name                = "%s"
  region              = "us-south1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
		psc_config {
			psc_enabled = true
			allowed_consumer_projects = ["%s"]
		}
		ipv4_enabled = false
    }
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	}
	availability_type = "REGIONAL"
  }
}
`, projectId, projectId, orgId, billingAccount, instanceName, projectId)
}

func verifyPscOperation(resourceName string, isPscConfigExpected bool, expectedPscEnabled bool, expectedAllowedConsumerProjects []string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Can't find %s in state", resourceName)
		}

		resourceAttributes := resource.Primary.Attributes
		_, ok = resourceAttributes["settings.0.ip_configuration.#"]
		if !ok {
			return fmt.Errorf("settings.0.ip_configuration.# block is not present in state for %s", resourceName)
		}

		if isPscConfigExpected {
			_, ok := resourceAttributes["settings.0.ip_configuration.0.psc_config.#"]
			if !ok {
				return fmt.Errorf("settings.0.ip_configuration.0.psc_config property is not present or set in state of %s", resourceName)
			}

			pscEnabledStr, ok := resourceAttributes["settings.0.ip_configuration.0.psc_config.0.psc_enabled"]
			pscEnabled, err := strconv.ParseBool(pscEnabledStr)
			if err != nil || pscEnabled != expectedPscEnabled {
				return fmt.Errorf("settings.0.ip_configuration.0.psc_config.0.psc_enabled property value is not set as expected in state of %s, expected %v, actual %v", resourceName, expectedPscEnabled, pscEnabled)
			}

			allowedConsumerProjectsStr, ok := resourceAttributes["settings.0.ip_configuration.0.psc_config.0.allowed_consumer_projects.#"]
			allowedConsumerProjects, err := strconv.Atoi(allowedConsumerProjectsStr)
			if !ok || allowedConsumerProjects != len(expectedAllowedConsumerProjects) {
				return fmt.Errorf("settings.0.ip_configuration.0.psc_config.0.allowed_consumer_projects property is not present or set as expected in state of %s", resourceName)
			}
		}

		return nil
	}
}

func verifyPscAutoConnectionsOperation(resourceName string, isPscConfigExpected bool, expectedPscEnabled bool, isPscAutoConnectionConfigExpected bool, expectedConsumerNetwork string, expectedConsumerProject string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Can't find %s in state", resourceName)
		}

		resourceAttributes := resource.Primary.Attributes
		_, ok = resourceAttributes["settings.0.ip_configuration.#"]
		if !ok {
			return fmt.Errorf("settings.0.ip_configuration.# block is not present in state for %s", resourceName)
		}

		if isPscConfigExpected {
			_, ok := resourceAttributes["settings.0.ip_configuration.0.psc_config.#"]
			if !ok {
				return fmt.Errorf("settings.0.ip_configuration.0.psc_config property is not present or set in state of %s", resourceName)
			}

			pscEnabledStr, ok := resourceAttributes["settings.0.ip_configuration.0.psc_config.0.psc_enabled"]
			pscEnabled, err := strconv.ParseBool(pscEnabledStr)
			if err != nil || pscEnabled != expectedPscEnabled {
				return fmt.Errorf("settings.0.ip_configuration.0.psc_config.0.psc_enabled property value is not set as expected in state of %s, expected %v, actual %v", resourceName, expectedPscEnabled, pscEnabled)
			}

			_, ok = resourceAttributes["settings.0.ip_configuration.0.psc_config.0.psc_auto_connections.#"]
			if !ok {
				return fmt.Errorf("settings.0.ip_configuration.0.psc_config.0.psc_auto_connections property is not present or set in state of %s", resourceName)
			}

			if isPscAutoConnectionConfigExpected {
				consumerNetwork, ok := resourceAttributes["settings.0.ip_configuration.0.psc_config.0.psc_auto_connections.0.consumer_network"]
				if !ok || consumerNetwork != expectedConsumerNetwork {
					return fmt.Errorf("settings.0.ip_configuration.0.psc_config.0.psc_auto_connections.0.consumer_network property is not present or set as expected in state of %s", resourceName)
				}

				consumerProject, ok := resourceAttributes["settings.0.ip_configuration.0.psc_config.0.psc_auto_connections.0.consumer_service_project_id"]
				if !ok || consumerProject != expectedConsumerProject {
					return fmt.Errorf("settings.0.ip_configuration.0.psc_config.0.psc_auto_connections.0.consumer_service_project_id property is not present or set as expected in state of %s", resourceName)
				}
			}
		}

		return nil
	}
}

func testAccSqlDatabaseInstance_withPSCEnabled_withoutPscAutoConnections(instanceName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-west2"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
		psc_config {
			psc_enabled = true
		}
		ipv4_enabled = false
    }
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	}
	availability_type = "REGIONAL"
  }
}
`, instanceName)
}

func testAccSqlDatabaseInstance_withPSCEnabled_withPscAutoConnections(instanceName string, projectId string, networkName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "testnetwork" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-west2"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
		psc_config {
			psc_enabled = true
			psc_auto_connections {
				consumer_network = "projects/%s/global/networks/%s"
				consumer_service_project_id = "%s"
			}
		}
		ipv4_enabled = false
    }
	backup_configuration {
		enabled = true
		binary_log_enabled = true
	}
	availability_type = "REGIONAL"
  }
}
`, networkName, instanceName, projectId, networkName, projectId)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withoutAllocatedIpRange(databaseName, networkName string, specifyPrivatePathOption bool, enablePrivatePath bool) string {
	privatePathOption := ""
	if specifyPrivatePathOption {
		privatePathOption = fmt.Sprintf("enable_private_path_for_google_cloud_services = %t", enablePrivatePath)
	}

	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      %s
    }
  }
}
`, networkName, databaseName, privatePathOption)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRange(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      allocated_ip_range = "%s"
    }
  }
}
`, networkName, databaseName, addressRangeName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeReplica(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
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
  name                = "%s-replica1"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
    ip_configuration {
      ipv4_enabled       = "false"
      private_network    = data.google_compute_network.servicenet.self_link
      allocated_ip_range = "%s"
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
`, networkName, databaseName, databaseName, addressRangeName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
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
    allocated_ip_range   = "%s"
  }

}
`, networkName, databaseName, databaseName, addressRangeName)
}

func testAccSqlDatabaseInstance_withPrivateNetwork_withAllocatedIpRangeClone_withSettings(databaseName, networkName, addressRangeName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name                    = "%s"
}

resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
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
    allocated_ip_range   = "%s"
  }

  settings {
    tier = "db-g1-small"
    backup_configuration {
      enabled = false
    }
  }
}
`, networkName, databaseName, databaseName, addressRangeName)
}

var testGoogleSqlDatabaseInstance_settings = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier                   = "db-g1-small"
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
    tier                   = "db-g1-small"
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
    tier                   = "db-g1-small"
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
    tier = "db-g1-small"
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
    tier                        = "db-g1-small"
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
    tier = "db-g1-small"
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
    tier = "db-g1-small"

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
    tier = "db-g1-small"
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
    tier = "db-g1-small"

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
    tier                  = "db-g1-small"
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
    tier = "db-g1-small"

    maintenance_window {
      day          = 7
      hour         = 3
      update_track = "canary"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_maintenance_week5 = `
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
      update_track = "week5"
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
    tier                   = "db-g1-small"

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
    tier                   = "db-g1-small"

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
    tier                   = "db-g1-small"
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
    tier = "db-g1-small"
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
    tier = "db-g1-small"
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
    tier = "db-g1-small"

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

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com"
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
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
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

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com"
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

  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

resource "google_kms_key_ring" "keyring-rep" {

  name     = "%{key_name}-rep"
  location = "us-east1"
}

resource "google_kms_crypto_key" "key-rep" {

  name     = "%{key_name}-rep"
  key_ring = google_kms_key_ring.keyring-rep.id
}

resource "google_kms_crypto_key_iam_member" "crypto_key_rep" {
  crypto_key_id = google_kms_crypto_key.key-rep.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloud-sql.iam.gserviceaccount.com"
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

  depends_on = [
    google_sql_database_instance.master,
    google_kms_crypto_key_iam_member.crypto_key_rep
  ]
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

func testGoogleSqlDatabaseInstance_EnableGoogleMlIntegration(masterID int, enableGoogleMlIntegration bool, dbVersion string, tier string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "%s"
  deletion_protection = false
  root_password		  = "rand-pwd-%d"
  settings {
    tier = "%s"
	enable_google_ml_integration = %t
  }
}
`, masterID, dbVersion, masterID, tier, enableGoogleMlIntegration)
}

func testGoogleSqlDatabaseInstance_EnableDataplexIntegration(masterID int, enableDataplexIntegration bool) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  root_password		  = "rand-pwd-%d"
  settings {
    tier = "db-custom-2-10240"
	enable_dataplex_integration = %t
  }
}
`, masterID, masterID, enableDataplexIntegration)
}

func testGoogleSqlDatabaseInstance_BackupRetention(masterID int) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "tf-test-%d"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-g1-small"
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
	tier = "db-g1-small"
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
	tier = "db-g1-small"
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
	tier = "db-g1-small"
	backup_configuration {
		enabled            = false
	}
  }

  clone {
    source_instance_name = data.google_sql_backup_run.backup.instance
    point_in_time = data.google_sql_backup_run.backup.start_time
	preferred_zone = "us-central1-b"
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

// Check that original-replica is now the primary
func checkSwitchoverOriginalReplicaConfigurations(replicaResourceName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		replicaResource, ok := s.RootModule().Resources[replicaResourceName]
		if !ok {
			return fmt.Errorf("Can't find %s in state", replicaResourceName)
		}
		replicaResourceAttributes := replicaResource.Primary.Attributes

		replicaInstanceType, ok := replicaResourceAttributes["instance_type"]
		if !ok {
			return fmt.Errorf("Instance type is not present in state for %s", replicaResourceName)
		}
		if replicaInstanceType != "CLOUD_SQL_INSTANCE" {
			return fmt.Errorf("Error in switchover. Original replica instance_type is %s, it should be CLOUD_SQL_INSTANCE.", replicaInstanceType)
		}

		replicaMasterInstanceName, ok := replicaResourceAttributes["master_instance_name"]
		if ok && replicaMasterInstanceName != "" {
			return fmt.Errorf("Error in switchover. master_instance_name should not be set on new primary")
		}
		return nil
	}
}

// Check that original-primary is now a replica
func checkSwitchoverOriginalPrimaryConfigurations(primaryResourceName string, replicaName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		primaryResource, ok := s.RootModule().Resources[primaryResourceName]
		if !ok {
			return fmt.Errorf("Can't find %s in state", primaryResourceName)
		}
		primaryResourceAttributes := primaryResource.Primary.Attributes
		primaryInstanceType, ok := primaryResourceAttributes["instance_type"]
		if !ok {
			return fmt.Errorf("Instance type is not present in state for %s", primaryResourceName)
		}
		if primaryInstanceType != "READ_REPLICA_INSTANCE" {
			return fmt.Errorf("Error in switchover. Original primary instance_type is %s, it should be READ_REPLICA_INSTANCE.", primaryInstanceType)
		}

		primaryMasterInstanceName, ok := primaryResourceAttributes["master_instance_name"]
		if !ok {
			return fmt.Errorf("Master instance name is not present in state for %s", primaryResourceName)
		}
		if primaryMasterInstanceName != replicaName {
			return fmt.Errorf("Error in switchover. master_instance_name should be %s", replicaName)
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

func testGoogleSqlDatabaseInstance_flags(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"

    database_flags {
      name  = "character_set_server"
      value = "utf8mb4"
    }
    database_flags {
      name  = "auto_increment_increment"
      value = "2"
    }
  }
}`, instance)
}

func testGoogleSqlDatabaseInstance_flags_update(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-g1-small"

    database_flags {
      name  = "auto_increment_increment"
      value = "2"
    }
    database_flags {
      name  = "character_set_server"
      value = "utf8mb4"
    }
  }
}`, instance)
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
    tier              = "db-g1-small"

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
    tier              = "db-g1-small"
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
    tier              = "db-g1-small"

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
    tier              = "db-g1-small"
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
    tier              = "db-g1-small"
    activation_policy = "%s"
  }
}
`, instance, databaseVersion, deletionProtection, activationPolicy)
}
