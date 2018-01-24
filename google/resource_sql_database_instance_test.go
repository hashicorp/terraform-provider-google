package google

/**
 * Note! You must run these tests once at a time. Google Cloud SQL does
 * not allow you to reuse a database for a short time after you reserved it,
 * and for this reason the tests will fail if the same config is used serveral
 * times in short succession.
 */

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/sqladmin/v1beta4"
)

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

	err = config.loadAndValidate()
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

			err = sqladminOperationWait(config, op, config.Project, "Stop Replica")
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

			err = sqladminOperationWait(config, op, config.Project, "Delete Instance")
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

func TestAccGoogleSqlDatabaseInstance_basic(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_basic2(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleSqlDatabaseInstance_basic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_basic3(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(
						&instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_dontDeleteDefaultUserOnReplica(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	databaseName := "sql-instance-test-" + acctest.RandString(10)
	failoverName := "sql-instance-test-failover-" + acctest.RandString(10)
	// 1. Create an instance.
	// 2. Add a root@'%' user.
	// 3. Create a replica and assert it succeeds (it'll fail if we try to delete the root user thinking it's a
	//    default user)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleSqlDatabaseInstanceConfig_withoutReplica(databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			}, resource.TestStep{
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
					err = sqladminOperationWait(config, op, config.Project, "Waiting for user to insert")
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

func TestAccGoogleSqlDatabaseInstance_settings_basic(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_slave(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	masterID := acctest.RandInt()
	slaveID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_slave, masterID, slaveID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance_master", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance_master", &instance),
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance_slave", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance_slave", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_highAvailability(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	instanceID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_highAvailability, instanceID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
					// Check that we've set our high availability type correctly, and it's been
					// accepted by the API
					func(s *terraform.State) error {
						if instance.Settings.AvailabilityType != "REGIONAL" {
							return fmt.Errorf("Database %s was not configured with Regional HA", instance.Name)
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_diskspecs(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	masterID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_diskspecs, masterID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_maintenance(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	masterID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenance, masterID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_settings_upgrade(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_settingsDowngrade(t *testing.T) {
	t.Parallel()

	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

// GH-4222
func TestAccGoogleSqlDatabaseInstance_authNets(t *testing.T) {
	t.Parallel(
	// var instance sqladmin.DatabaseInstance
	)

	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step2, databaseID),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
		},
	})
}

// Tests that a SQL instance can be referenced from more than one other resource without
// throwing an error during provisioning, see #9018.
func TestAccGoogleSqlDatabaseInstance_multipleOperations(t *testing.T) {
	t.Parallel()

	databaseID, instanceID, userID := acctest.RandString(8), acctest.RandString(8), acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_multipleOperations, databaseID, instanceID, userID),
			},
		},
	})
}

func testAccCheckGoogleSqlDatabaseInstanceEquals(n string,
	instance *sqladmin.DatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		attributes := rs.Primary.Attributes

		server := instance.Name
		local := attributes["name"]
		if server != local {
			return fmt.Errorf("Error name mismatch, (%s, %s)", server, local)
		}

		server = instance.Settings.Tier
		local = attributes["settings.0.tier"]
		if server != local {
			return fmt.Errorf("Error settings.tier mismatch, (%s, %s)", server, local)
		}

		server = strings.TrimPrefix(instance.MasterInstanceName, instance.Project+":")
		local = attributes["master_instance_name"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error master_instance_name mismatch, (%s, %s)", server, local)
		}

		ip_len, err := strconv.Atoi(attributes["ip_address.#"])
		if err != nil {
			return fmt.Errorf("Error parsing ip_addresses.# : %s", err.Error())
		}
		if ip_len != len(instance.IpAddresses) {
			return fmt.Errorf("Error ip_addresses.# mismatch, server has %d but local has %d", len(instance.IpAddresses), ip_len)
		}
		// For now, assume the order matches
		for idx, ip := range instance.IpAddresses {
			server = attributes["ip_address."+strconv.Itoa(idx)+".ip_address"]
			local = ip.IpAddress
			if server != local {
				return fmt.Errorf("Error ip_addresses.%d.ip_address mismatch, server has %s but local has %s", idx, server, local)
			}

			server = attributes["ip_address."+strconv.Itoa(idx)+".time_to_retire"]
			local = ip.TimeToRetire
			if server != local {
				return fmt.Errorf("Error ip_addresses.%d.time_to_retire mismatch, server has %s but local has %s", idx, server, local)
			}
		}

		server = instance.Settings.ActivationPolicy
		local = attributes["settings.0.activation_policy"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.activation_policy mismatch, (%s, %s)", server, local)
		}

		server = instance.Settings.AvailabilityType
		local = attributes["settings.0.availability_type"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.availability_type mismatch, (%s, %s)", server, local)
		}

		if instance.Settings.BackupConfiguration != nil {
			server = strconv.FormatBool(instance.Settings.BackupConfiguration.BinaryLogEnabled)
			local = attributes["settings.0.backup_configuration.0.binary_log_enabled"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.backup_configuration.binary_log_enabled mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatBool(instance.Settings.BackupConfiguration.Enabled)
			local = attributes["settings.0.backup_configuration.0.enabled"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.backup_configuration.enabled mismatch, (%s, %s)", server, local)
			}

			server = instance.Settings.BackupConfiguration.StartTime
			local = attributes["settings.0.backup_configuration.0.start_time"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.backup_configuration.start_time mismatch, (%s, %s)", server, local)
			}
		}

		server = strconv.FormatBool(instance.Settings.CrashSafeReplicationEnabled)
		local = attributes["settings.0.crash_safe_replication"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.crash_safe_replication mismatch, (%s, %s)", server, local)
		}

		// First generation CloudSQL instances will not have any value for StorageAutoResize.
		// We need to check if this value has been omitted before we potentially deference a
		// nil pointer.
		if instance.Settings.StorageAutoResize != nil {
			server = strconv.FormatBool(*instance.Settings.StorageAutoResize)
			local = attributes["settings.0.disk_autoresize"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.disk_autoresize mismatch, (%s, %s)", server, local)
			}
		}

		server = strconv.FormatInt(instance.Settings.DataDiskSizeGb, 10)
		local = attributes["settings.0.disk_size"]
		if server != local && len(server) > 0 && len(local) > 0 && local != "0" {
			return fmt.Errorf("Error settings.disk_size mismatch, (%s, %s)", server, local)
		}

		server = instance.Settings.DataDiskType
		local = attributes["settings.0.disk_type"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.disk_type mismatch, (%s, %s)", server, local)
		}

		if instance.Settings.IpConfiguration != nil {
			server = strconv.FormatBool(instance.Settings.IpConfiguration.Ipv4Enabled)
			local = attributes["settings.0.ip_configuration.0.ipv4_enabled"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.ip_configuration.ipv4_enabled mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatBool(instance.Settings.IpConfiguration.RequireSsl)
			local = attributes["settings.0.ip_configuration.0.require_ssl"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.ip_configuration.require_ssl mismatch, (%s, %s)", server, local)
			}
		}

		if instance.Settings.LocationPreference != nil {
			server = instance.Settings.LocationPreference.FollowGaeApplication
			local = attributes["settings.0.location_preference.0.follow_gae_application"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.location_preference.follow_gae_application mismatch, (%s, %s)", server, local)
			}

			server = instance.Settings.LocationPreference.Zone
			local = attributes["settings.0.location_preference.0.zone"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.location_preference.zone mismatch, (%s, %s)", server, local)
			}
		}

		if instance.Settings.MaintenanceWindow != nil {
			server = strconv.FormatInt(instance.Settings.MaintenanceWindow.Day, 10)
			local = attributes["settings.0.maintenance_window.0.day"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.maintenance_window.day mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatInt(instance.Settings.MaintenanceWindow.Hour, 10)
			local = attributes["settings.0.maintenance_window.0.hour"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.maintenance_window.hour mismatch, (%s, %s)", server, local)
			}

			server = instance.Settings.MaintenanceWindow.UpdateTrack
			local = attributes["settings.0.maintenance_window.0.update_track"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.maintenance_window.update_track mismatch, (%s, %s)", server, local)
			}
		}

		server = instance.Settings.PricingPlan
		local = attributes["settings.0.pricing_plan"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.pricing_plan mismatch, (%s, %s)", server, local)
		}

		if instance.ReplicaConfiguration != nil {
			server = strconv.FormatBool(instance.ReplicaConfiguration.FailoverTarget)
			local = attributes["replica_configuration.0.failover_target"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.failover_target mismatch, (%s, %s)", server, local)
			}
		}

		server = instance.ConnectionName
		local = attributes["connection_name"]
		if server != local {
			return fmt.Errorf("Error connection_name mismatch. (%s, %s)", server, local)
		}

		return nil
	}
}

func testAccCheckGoogleSqlDatabaseInstanceExists(n string,
	instance *sqladmin.DatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		found, err := config.clientSqlAdmin.Instances.Get(config.Project,
			rs.Primary.Attributes["name"]).Do()

		*instance = *found

		if err != nil {
			return fmt.Errorf("Not found: %s", n)
		}

		return nil
	}
}

func testAccGoogleSqlDatabaseInstanceDestroy(s *terraform.State) error {
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

func testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(
	instance *sqladmin.DatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		users, err := config.clientSqlAdmin.Users.List(config.Project, instance.Name).Do()

		if err != nil {
			return fmt.Errorf("Could not list database users for %q: %s", instance.Name, err)
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
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false
	}
}
`

var testGoogleSqlDatabaseInstance_basic2 = `
resource "google_sql_database_instance" "instance" {
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false
	}
}
`
var testGoogleSqlDatabaseInstance_basic3 = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central1"
	settings {
		tier = "db-f1-micro"
	}
}
`

func testGoogleSqlDatabaseInstanceConfig_withoutReplica(instanceName string) string {
	return fmt.Sprintf(`resource "google_sql_database_instance" "instance" {
  name               = "%s"
  region             = "us-central1"
  database_version   = "MYSQL_5_7"

  settings {
    tier             = "db-n1-standard-1"

    backup_configuration {
        binary_log_enabled = "true"
        enabled            = "true"
        start_time         = "18:00"
    }
  }
}`, instanceName)
}

func testGoogleSqlDatabaseInstanceConfig_withReplica(instanceName, failoverName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name               = "%s"
  region             = "us-central1"
  database_version   = "MYSQL_5_7"

  settings {
    tier             = "db-n1-standard-1"

    backup_configuration {
        binary_log_enabled = "true"
        enabled            = "true"
        start_time         = "18:00"
    }
  }
}

resource "google_sql_database_instance" "instance-failover" {
  name               = "%s"
  region             = "us-central1"
  database_version   = "MYSQL_5_7"
  master_instance_name = "${google_sql_database_instance.instance.name}"

  replica_configuration {
    failover_target        = "true"
  }

  settings {
    tier             = "db-n1-standard-1"
  }
}
`, instanceName, failoverName)
}

var testGoogleSqlDatabaseInstance_settings = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false
		replication_type = "ASYNCHRONOUS"
		location_preference {
			zone = "us-central1-f"
		}

		ip_configuration {
			ipv4_enabled = "true"
			authorized_networks {
				value = "108.12.12.12"
				name = "misc"
				expiration_time = "2050-11-15T16:19:00.094Z"
			}
		}

		backup_configuration {
			enabled = "true"
			start_time = "19:19"
		}

		activation_policy = "ON_DEMAND"
	}
}
`

// Note - this test is not feasible to run unless we generate
// backups first.
var testGoogleSqlDatabaseInstance_replica = `
resource "google_sql_database_instance" "instance_master" {
	name = "tf-lw-%d"
	database_version = "MYSQL_5_6"
	region = "us-east1"

	settings {
		tier = "D0"
		crash_safe_replication = true

		backup_configuration {
			enabled = true
			start_time = "00:00"
			binary_log_enabled = true
		}
	}
}

resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	database_version = "MYSQL_5_6"
	region = "us-central"

	settings {
		tier = "D0"
	}

	master_instance_name = "${google_sql_database_instance.instance_master.name}"

	replica_configuration {
		ca_certificate = "${file("~/tmp/fake.pem")}"
		client_certificate = "${file("~/tmp/fake.pem")}"
		client_key = "${file("~/tmp/fake.pem")}"
		connect_retry_interval = 100
		master_heartbeat_period = 10000
		password = "password"
		username = "username"
		ssl_cipher = "ALL"
		verify_server_certificate = false
	}
}
`

var testGoogleSqlDatabaseInstance_slave = `
resource "google_sql_database_instance" "instance_master" {
	name = "tf-lw-%d"
	region = "us-central1"

	settings {
		tier = "db-f1-micro"

		backup_configuration {
			enabled = true
			binary_log_enabled = true
		}
	}
}

resource "google_sql_database_instance" "instance_slave" {
	name = "tf-lw-%d"
	region = "us-central1"

	master_instance_name = "${google_sql_database_instance.instance_master.name}"

	settings {
		tier = "db-f1-micro"
	}
}
`

var testGoogleSqlDatabaseInstance_highAvailability = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central1"
	database_version = "POSTGRES_9_6"

	settings {
		tier = "db-f1-micro"

		availability_type = "REGIONAL"

		backup_configuration {
			enabled = true
			binary_log_enabled = true
		}
	}
}
`

var testGoogleSqlDatabaseInstance_diskspecs = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central1"

	settings {
		tier = "db-f1-micro"
		disk_autoresize = true
		disk_size = 15
		disk_type = "PD_HDD"
	}
}
`

var testGoogleSqlDatabaseInstance_maintenance = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central1"

	settings {
		tier = "db-f1-micro"

		maintenance_window {
		  day  = 7
		  hour = 3
			update_track = "canary"
	  }
	}
}
`

var testGoogleSqlDatabaseInstance_authNets_step1 = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false

		ip_configuration {
			ipv4_enabled = "true"
			authorized_networks {
				value = "108.12.12.12"
				name = "misc"
				expiration_time = "2050-11-15T16:19:00.094Z"
			}
		}
	}
}
`

var testGoogleSqlDatabaseInstance_authNets_step2 = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false

		ip_configuration {
			ipv4_enabled = "true"
		}
	}
}
`

var testGoogleSqlDatabaseInstance_multipleOperations = `
resource "google_sql_database_instance" "instance" {
	name = "tf-test-%s"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false
	}
}

resource "google_sql_database" "database" {
	name = "tf-test-%s"
	instance = "${google_sql_database_instance.instance.name}"
}

resource "google_sql_user" "user" {
	name = "tf-test-%s"
	instance = "${google_sql_database_instance.instance.name}"
	host = "google.com"
	password = "hunter2"
}
`
