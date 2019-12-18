package google

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
}

func init() {
	resource.AddTestSweepers("gcp_sql_db_instance", &resource.Sweeper{
		Name: "gcp_sql_db_instance",
		F:    testSweepDatabases,
	})
}

func testSweepDatabases(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate()
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	found, err := config.clientSqlAdmin.Instances.List(config.Project).Do()
	if err != nil {
		log.Fatalf("error listing databases: %s", err)
	}

	if len(found.Items) == 0 {
		log.Printf("No databases found")
		return nil
	}

	running := map[string]struct{}{}

	for _, d := range found.Items {
		var testDbInstance bool
		for _, testName := range []string{"tf-lw-", "sqldatabasetest"} {
			// only destroy instances we know to fit our test naming pattern
			if strings.HasPrefix(d.Name, testName) {
				testDbInstance = true
			}
		}

		if !testDbInstance {
			continue
		}
		if d.State != "RUNNABLE" {
			continue
		}
		running[d.Name] = struct{}{}
	}

	for _, d := range found.Items {
		// don't delete replicas, we'll take care of that
		// when deleting the database they replicate
		if d.ReplicaConfiguration != nil {
			continue
		}
		log.Printf("Destroying SQL Instance (%s)", d.Name)

		// replicas need to be stopped and destroyed before destroying a master
		// instance. The ordering slice tracks replica databases for a given master
		// and we call destroy on them before destroying the master
		var ordering []string
		for _, replicaName := range d.ReplicaNames {
			// don't try to stop replicas that aren't running
			if _, ok := running[replicaName]; !ok {
				ordering = append(ordering, replicaName)
				continue
			}

			// need to stop replication before being able to destroy a database
			op, err := config.clientSqlAdmin.Instances.StopReplica(config.Project, replicaName).Do()

			if err != nil {
				return fmt.Errorf("error, failed to stop replica instance (%s) for instance (%s): %s", replicaName, d.Name, err)
			}

			err = sqlAdminOperationWait(config, op, config.Project, "Stop Replica")
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					log.Printf("Replication operation not found")
				} else {
					return err
				}
			}

			ordering = append(ordering, replicaName)
		}

		// ordering has a list of replicas (or none), now add the primary to the end
		ordering = append(ordering, d.Name)

		for _, db := range ordering {
			// destroy instances, replicas first
			op, err := config.clientSqlAdmin.Instances.Delete(config.Project, db).Do()

			if err != nil {
				if strings.Contains(err.Error(), "409") {
					// the GCP api can return a 409 error after the delete operation
					// reaches a successful end
					log.Printf("Operation not found, got 409 response")
					continue
				}

				return fmt.Errorf("Error, failed to delete instance %s: %s", db, err)
			}

			err = sqlAdminOperationWait(config, op, config.Project, "Delete Instance")
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					log.Printf("SQL instance not found")
					continue
				}
				return err
			}
		}
	}

	return nil
}

func TestAccSqlDatabaseInstance_basicFirstGen(t *testing.T) {
	t.Parallel()

	instanceID := acctest.RandInt()
	instanceName := fmt.Sprintf("tf-lw-%d", instanceID)
	resourceName := "google_sql_database_instance.instance"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testGoogleSqlDatabaseInstance_basic, instanceID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     fmt.Sprintf("projects/%s/instances/%s", getTestProjectFromEnv(), instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     fmt.Sprintf("%s/%s", getTestProjectFromEnv(), instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicInferredName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstance_basic2,
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basicSecondGen(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseName),
				Check: testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(databaseName),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_dontDeleteDefaultUserOnReplica(t *testing.T) {
	t.Parallel()

	databaseName := "sql-instance-test-" + acctest.RandString(10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(10)
	// 1. Create an instance.
	// 2. Add a root@'%' user.
	// 3. Create a replica and assert it succeeds (it'll fail if we try to delete the root user thinking it's a
	//    default user)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlDatabaseInstanceConfig_withoutReplica(databaseName),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					// Add a root user
					config := testAccProvider.Meta().(*Config)
					user := sqladmin.User{
						Name:     "root",
						Host:     "%",
						Password: acctest.RandString(26),
					}
					op, err := config.clientSqlAdmin.Users.Insert(config.Project, databaseName, &user).Do()
					if err != nil {
						t.Errorf("Error while inserting root@%% user: %s", err)
						return
					}
					err = sqlAdminOperationWait(config, op, config.Project, "Waiting for user to insert")
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

func TestAccSqlDatabaseInstance_settings_basic(t *testing.T) {
	t.Parallel()

	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_replica(t *testing.T) {
	t.Parallel()

	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_replica, databaseID, databaseID, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance_master",
				ImportState:       true,
				ImportStateVerify: true,
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

	masterID := acctest.RandInt()
	slaveID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_slave, masterID, slaveID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance_master",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_sql_database_instance.instance_slave",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_highAvailability(t *testing.T) {
	t.Parallel()

	instanceID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_highAvailability, instanceID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_diskspecs(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_diskspecs, masterID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_maintenance(t *testing.T) {
	t.Parallel()

	masterID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenance, masterID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settings_upgrade(t *testing.T) {
	t.Parallel()

	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_settingsDowngrade(t *testing.T) {
	t.Parallel()

	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// GH-4222
func TestAccSqlDatabaseInstance_authNets(t *testing.T) {
	t.Parallel()

	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step2, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Tests that a SQL instance can be referenced from more than one other resource without
// throwing an error during provisioning, see #9018.
func TestAccSqlDatabaseInstance_multipleOperations(t *testing.T) {
	t.Parallel()

	databaseID, instanceID, userID := acctest.RandString(8), acctest.RandString(8), acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_multipleOperations, databaseID, instanceID, userID),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlDatabaseInstance_basic_with_user_labels(t *testing.T) {
	t.Parallel()

	databaseName := "tf-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_with_user_labels, databaseName),
				Check: testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(databaseName),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic_with_user_labels_update, databaseName),
			},
			{
				ResourceName:      "google_sql_database_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSqlDatabaseInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		config := testAccProvider.Meta().(*Config)
		if rs.Type != "google_sql_database_instance" {
			continue
		}

		_, err := config.clientSqlAdmin.Instances.Get(config.Project,
			rs.Primary.Attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("Database Instance still exists")
		}
	}

	return nil
}

func testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(instance string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		users, err := config.clientSqlAdmin.Users.List(config.Project, instance).Do()

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

var testGoogleSqlDatabaseInstance_basic = `
resource "google_sql_database_instance" "instance" {
  name   = "tf-lw-%d"
  region = "us-central"
  settings {
    tier                   = "D0"
    crash_safe_replication = false
  }
}
`

var testGoogleSqlDatabaseInstance_basic2 = `
resource "google_sql_database_instance" "instance" {
  region = "us-central"
  settings {
    tier                   = "D0"
    crash_safe_replication = false
  }
}
`

var testGoogleSqlDatabaseInstance_basic3 = `
resource "google_sql_database_instance" "instance" {
  name   = "%s"
  region = "us-central1"
  settings {
    tier = "db-f1-micro"
  }
}
`

func testGoogleSqlDatabaseInstanceConfig_withoutReplica(instanceName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "MYSQL_5_7"

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
  name             = "%s"
  region           = "us-central1"
  database_version = "MYSQL_5_7"

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

  replica_configuration {
    failover_target = "true"
  }

  settings {
    tier = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

var testGoogleSqlDatabaseInstance_settings = `
resource "google_sql_database_instance" "instance" {
  name   = "tf-lw-%d"
  region = "us-central"
  settings {
    tier                   = "D0"
    crash_safe_replication = false
    replication_type       = "ASYNCHRONOUS"
    location_preference {
      zone = "us-central1-f"
    }

    ip_configuration {
      ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2050-11-15T16:19:00.094Z"
      }
    }

    backup_configuration {
      enabled    = "true"
      start_time = "19:19"
    }

    activation_policy = "ON_DEMAND"
  }
}
`

var testGoogleSqlDatabaseInstance_replica = `
resource "google_sql_database_instance" "instance_master" {
  name             = "tf-lw-%d"
  database_version = "MYSQL_5_6"
  region           = "us-central1"

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
  name             = "tf-lw-%d-1"
  database_version = "MYSQL_5_6"
  region           = "us-central1"

  settings {
    tier = "db-n1-standard-1"
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
  name             = "tf-lw-%d-2"
  database_version = "MYSQL_5_6"
  region           = "us-central1"

  settings {
    tier = "db-n1-standard-1"
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
  name   = "tf-lw-%d"
  region = "us-central1"

  settings {
    tier = "db-f1-micro"

    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }
  }
}

resource "google_sql_database_instance" "instance_slave" {
  name   = "tf-lw-%d"
  region = "us-central1"

  master_instance_name = google_sql_database_instance.instance_master.name

  settings {
    tier = "db-f1-micro"
  }
}
`

var testGoogleSqlDatabaseInstance_highAvailability = `
resource "google_sql_database_instance" "instance" {
  name             = "tf-lw-%d"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"

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
  name   = "tf-lw-%d"
  region = "us-central1"

  settings {
    tier            = "db-f1-micro"
    disk_autoresize = true
    disk_size       = 15
    disk_type       = "PD_HDD"
  }
}
`

var testGoogleSqlDatabaseInstance_maintenance = `
resource "google_sql_database_instance" "instance" {
  name   = "tf-lw-%d"
  region = "us-central1"

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
  name   = "tf-lw-%d"
  region = "us-central"
  settings {
    tier                   = "D0"
    crash_safe_replication = false

    ip_configuration {
      ipv4_enabled = "true"
      authorized_networks {
        value           = "108.12.12.12"
        name            = "misc"
        expiration_time = "2050-11-15T16:19:00.094Z"
      }
    }
  }
}
`

var testGoogleSqlDatabaseInstance_authNets_step2 = `
resource "google_sql_database_instance" "instance" {
  name   = "tf-lw-%d"
  region = "us-central"
  settings {
    tier                   = "D0"
    crash_safe_replication = false

    ip_configuration {
      ipv4_enabled = "true"
    }
  }
}
`

var testGoogleSqlDatabaseInstance_multipleOperations = `
resource "google_sql_database_instance" "instance" {
  name   = "tf-test-%s"
  region = "us-central"
  settings {
    tier                   = "D0"
    crash_safe_replication = false
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
  name   = "%s"
  region = "us-central1"
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
  name   = "%s"
  region = "us-central1"
  settings {
    tier = "db-f1-micro"
    user_labels = {
      track = "production"
    }
  }
}
`
