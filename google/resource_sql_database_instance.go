package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

const privateNetworkLinkRegex = "projects/(" + ProjectRegex + ")/global/networks/((?:[a-z](?:[-a-z0-9]*[a-z0-9])?))$"

var sqlDatabaseAuthorizedNetWorkSchemaElem *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"expiration_time": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

var (
	backupConfigurationKeys = []string{
		"settings.0.backup_configuration.0.binary_log_enabled",
		"settings.0.backup_configuration.0.enabled",
		"settings.0.backup_configuration.0.start_time",
		"settings.0.backup_configuration.0.location",
	}

	ipConfigurationKeys = []string{
		"settings.0.ip_configuration.0.authorized_networks",
		"settings.0.ip_configuration.0.ipv4_enabled",
		"settings.0.ip_configuration.0.require_ssl",
		"settings.0.ip_configuration.0.private_network",
	}

	maintenanceWindowKeys = []string{
		"settings.0.maintenance_window.0.day",
		"settings.0.maintenance_window.0.hour",
		"settings.0.maintenance_window.0.update_track",
	}

	serverCertsKeys = []string{
		"server_ca_cert.0.cert",
		"server_ca_cert.0.common_name",
		"server_ca_cert.0.create_time",
		"server_ca_cert.0.expiration_time",
		"server_ca_cert.0.sha1_fingerprint",
	}

	replicaConfigurationKeys = []string{
		"replica_configuration.0.ca_certificate",
		"replica_configuration.0.client_certificate",
		"replica_configuration.0.client_key",
		"replica_configuration.0.connect_retry_interval",
		"replica_configuration.0.dump_file_path",
		"replica_configuration.0.failover_target",
		"replica_configuration.0.master_heartbeat_period",
		"replica_configuration.0.password",
		"replica_configuration.0.ssl_cipher",
		"replica_configuration.0.username",
		"replica_configuration.0.verify_server_certificate",
	}
)

func resourceSqlDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlDatabaseInstanceCreate,
		Read:   resourceSqlDatabaseInstanceRead,
		Update: resourceSqlDatabaseInstanceUpdate,
		Delete: resourceSqlDatabaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSqlDatabaseInstanceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("settings.0.disk_size", isDiskShrinkage)),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tier": {
							Type:     schema.TypeString,
							Required: true,
						},
						"activation_policy": {
							Type:     schema.TypeString,
							Optional: true,
							// Defaults differ between first and second gen instances
							Computed: true,
						},
						"authorized_gae_applications": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"availability_type": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressFirstGen,
							// Set computed instead of default because this property is for second-gen
							// only. The default when not provided is ZONAL, which means no explicit HA
							// configuration.
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"REGIONAL", "ZONAL"}, false),
						},
						"backup_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"binary_log_enabled": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
									},
									"enabled": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
									},
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
										// start_time is randomly assigned if not set
										Computed:     true,
										AtLeastOneOf: backupConfigurationKeys,
									},
									"location": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
									},
								},
							},
						},
						"crash_safe_replication": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"database_flags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"disk_autoresize": {
							Type:             schema.TypeBool,
							Optional:         true,
							Default:          true,
							DiffSuppressFunc: suppressFirstGen,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
							// Defaults differ between first and second gen instances
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							// Set computed instead of default because this property is for second-gen only.
							Computed: true,
						},
						"ip_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorized_networks": {
										Type:         schema.TypeSet,
										Optional:     true,
										Set:          schema.HashResource(sqlDatabaseAuthorizedNetWorkSchemaElem),
										Elem:         sqlDatabaseAuthorizedNetWorkSchemaElem,
										AtLeastOneOf: ipConfigurationKeys,
									},
									"ipv4_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										// Defaults differ between first and second gen instances
										Computed:     true,
										AtLeastOneOf: ipConfigurationKeys,
									},
									"require_ssl": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: ipConfigurationKeys,
									},
									"private_network": {
										Type:             schema.TypeString,
										Optional:         true,
										ValidateFunc:     orEmpty(validateRegexp(privateNetworkLinkRegex)),
										DiffSuppressFunc: compareSelfLinkRelativePaths,
										AtLeastOneOf:     ipConfigurationKeys,
									},
								},
							},
						},
						"location_preference": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"follow_gae_application": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: []string{"settings.0.location_preference.0.follow_gae_application", "settings.0.location_preference.0.zone"},
									},
									"zone": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: []string{"settings.0.location_preference.0.follow_gae_application", "settings.0.location_preference.0.zone"},
									},
								},
							},
						},
						"maintenance_window": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 7),
										AtLeastOneOf: maintenanceWindowKeys,
									},
									"hour": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 23),
										AtLeastOneOf: maintenanceWindowKeys,
									},
									"update_track": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: maintenanceWindowKeys,
									},
								},
							},
						},
						"pricing_plan": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "PER_USE",
						},
						"replication_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "SYNCHRONOUS",
						},
						"user_labels": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"connection_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"database_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MYSQL_5_6",
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_to_retire": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"first_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"master_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"replica_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// Returned from API on all replicas
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ca_certificate": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"client_certificate": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"client_key": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"connect_retry_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"dump_file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"failover_target": {
							Type:         schema.TypeBool,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"master_heartbeat_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Sensitive:    true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"ssl_cipher": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"username": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
						"verify_server_certificate": {
							Type:         schema.TypeBool,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
						},
					},
				},
			},
			"server_ca_cert": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert": {
							Type:         schema.TypeString,
							Computed:     true,
							AtLeastOneOf: serverCertsKeys,
						},
						"common_name": {
							Type:         schema.TypeString,
							Computed:     true,
							AtLeastOneOf: serverCertsKeys,
						},
						"create_time": {
							Type:         schema.TypeString,
							Computed:     true,
							AtLeastOneOf: serverCertsKeys,
						},
						"expiration_time": {
							Type:         schema.TypeString,
							Computed:     true,
							AtLeastOneOf: serverCertsKeys,
						},
						"sha1_fingerprint": {
							Type:         schema.TypeString,
							Computed:     true,
							AtLeastOneOf: serverCertsKeys,
						},
					},
				},
			},
			"service_account_email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Suppress diff with any attribute value that is not supported on 1st Generation
// Instances
func suppressFirstGen(k, old, new string, d *schema.ResourceData) bool {
	if isFirstGen(d) {
		log.Printf("[DEBUG] suppressing diff on %s due to 1st gen instance type", k)
		return true
	}

	return false
}

// Detects whether a database is 1st Generation by inspecting the tier name
func isFirstGen(d *schema.ResourceData) bool {
	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})
	tier := settings["tier"].(string)

	// 1st Generation databases have tiers like 'D0', as opposed to 2nd Generation which are
	// prefixed with 'db'
	return !regexp.MustCompile("db*").Match([]byte(tier))
}

func resourceSqlDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		name = resource.UniqueId()
	}

	d.Set("name", name)

	instance := &sqladmin.DatabaseInstance{
		Name:                 name,
		Region:               region,
		Settings:             expandSqlDatabaseInstanceSettings(d.Get("settings").([]interface{}), !isFirstGen(d)),
		DatabaseVersion:      d.Get("database_version").(string),
		MasterInstanceName:   d.Get("master_instance_name").(string),
		ReplicaConfiguration: expandReplicaConfiguration(d.Get("replica_configuration").([]interface{})),
	}

	// Modifying a replica during Create can cause problems if the master is
	// modified at the same time. Lock the master until we're done in order
	// to prevent that.
	if !sqlDatabaseIsMaster(d) {
		mutexKV.Lock(instanceMutexKey(project, instance.MasterInstanceName))
		defer mutexKV.Unlock(instanceMutexKey(project, instance.MasterInstanceName))
	}

	var op *sqladmin.Operation
	err = retryTimeDuration(func() (operr error) {
		op, operr = config.clientSqlAdmin.Instances.Insert(project, instance).Do()
		return operr
	}, d.Timeout(schema.TimeoutCreate), isSqlOperationInProgressError)
	if err != nil {
		return fmt.Errorf("Error, failed to create instance %s: %s", instance.Name, err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/instances/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = sqlAdminOperationWaitTime(config, op, project, "Create Instance", int(d.Timeout(schema.TimeoutCreate).Minutes()))
	if err != nil {
		d.SetId("")
		return err
	}

	err = resourceSqlDatabaseInstanceRead(d, meta)
	if err != nil {
		return err
	}

	// If a default root user was created with a wildcard ('%') hostname, delete it.
	// Users in a replica instance are inherited from the master instance and should be left alone.
	if sqlDatabaseIsMaster(d) {
		var users *sqladmin.UsersListResponse
		err = retryTimeDuration(func() error {
			users, err = config.clientSqlAdmin.Users.List(project, instance.Name).Do()
			return err
		}, d.Timeout(schema.TimeoutRead), isSqlOperationInProgressError)
		if err != nil {
			return fmt.Errorf("Error, attempting to list users associated with instance %s: %s", instance.Name, err)
		}
		for _, u := range users.Items {
			if u.Name == "root" && u.Host == "%" {
				err = retry(func() error {
					op, err = config.clientSqlAdmin.Users.Delete(project, instance.Name, u.Host, u.Name).Do()
					if err == nil {
						err = sqlAdminOperationWaitTime(config, op, project, "Delete default root User", int(d.Timeout(schema.TimeoutCreate).Minutes()))
					}
					return err
				})
				if err != nil {
					return fmt.Errorf("Error, failed to delete default 'root'@'*' user, but the database was created successfully: %s", err)
				}
			}
		}
	}

	return nil
}

func expandSqlDatabaseInstanceSettings(configured []interface{}, secondGen bool) *sqladmin.Settings {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_settings := configured[0].(map[string]interface{})
	settings := &sqladmin.Settings{
		// Version is unset in Create but is set during update
		SettingsVersion:             int64(_settings["version"].(int)),
		Tier:                        _settings["tier"].(string),
		ForceSendFields:             []string{"StorageAutoResize"},
		ActivationPolicy:            _settings["activation_policy"].(string),
		AvailabilityType:            _settings["availability_type"].(string),
		CrashSafeReplicationEnabled: _settings["crash_safe_replication"].(bool),
		DataDiskSizeGb:              int64(_settings["disk_size"].(int)),
		DataDiskType:                _settings["disk_type"].(string),
		PricingPlan:                 _settings["pricing_plan"].(string),
		ReplicationType:             _settings["replication_type"].(string),
		UserLabels:                  convertStringMap(_settings["user_labels"].(map[string]interface{})),
		BackupConfiguration:         expandBackupConfiguration(_settings["backup_configuration"].([]interface{})),
		DatabaseFlags:               expandDatabaseFlags(_settings["database_flags"].([]interface{})),
		AuthorizedGaeApplications:   expandAuthorizedGaeApplications(_settings["authorized_gae_applications"].([]interface{})),
		IpConfiguration:             expandIpConfiguration(_settings["ip_configuration"].([]interface{})),
		LocationPreference:          expandLocationPreference(_settings["location_preference"].([]interface{})),
		MaintenanceWindow:           expandMaintenanceWindow(_settings["maintenance_window"].([]interface{})),
	}

	// 1st Generation instances don't support the disk_autoresize parameter
	// and it defaults to true - so we shouldn't set it if this is first gen
	if secondGen {
		settings.StorageAutoResize = googleapi.Bool(_settings["disk_autoresize"].(bool))
	}

	return settings
}

func expandReplicaConfiguration(configured []interface{}) *sqladmin.ReplicaConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_replicaConfiguration := configured[0].(map[string]interface{})
	return &sqladmin.ReplicaConfiguration{
		FailoverTarget: _replicaConfiguration["failover_target"].(bool),

		// MysqlReplicaConfiguration has been flattened in the TF schema, so
		// we'll keep it flat here instead of another expand method.
		MysqlReplicaConfiguration: &sqladmin.MySqlReplicaConfiguration{
			CaCertificate:           _replicaConfiguration["ca_certificate"].(string),
			ClientCertificate:       _replicaConfiguration["client_certificate"].(string),
			ClientKey:               _replicaConfiguration["client_key"].(string),
			ConnectRetryInterval:    int64(_replicaConfiguration["connect_retry_interval"].(int)),
			DumpFilePath:            _replicaConfiguration["dump_file_path"].(string),
			MasterHeartbeatPeriod:   int64(_replicaConfiguration["master_heartbeat_period"].(int)),
			Password:                _replicaConfiguration["password"].(string),
			SslCipher:               _replicaConfiguration["ssl_cipher"].(string),
			Username:                _replicaConfiguration["username"].(string),
			VerifyServerCertificate: _replicaConfiguration["verify_server_certificate"].(bool),
		},
	}
}

func expandMaintenanceWindow(configured []interface{}) *sqladmin.MaintenanceWindow {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	window := configured[0].(map[string]interface{})
	return &sqladmin.MaintenanceWindow{
		Day:             int64(window["day"].(int)),
		Hour:            int64(window["hour"].(int)),
		UpdateTrack:     window["update_track"].(string),
		ForceSendFields: []string{"Hour"},
	}
}

func expandLocationPreference(configured []interface{}) *sqladmin.LocationPreference {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_locationPreference := configured[0].(map[string]interface{})
	return &sqladmin.LocationPreference{
		FollowGaeApplication: _locationPreference["follow_gae_application"].(string),
		Zone:                 _locationPreference["zone"].(string),
	}
}

func expandIpConfiguration(configured []interface{}) *sqladmin.IpConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_ipConfiguration := configured[0].(map[string]interface{})

	return &sqladmin.IpConfiguration{
		Ipv4Enabled:        _ipConfiguration["ipv4_enabled"].(bool),
		RequireSsl:         _ipConfiguration["require_ssl"].(bool),
		PrivateNetwork:     _ipConfiguration["private_network"].(string),
		AuthorizedNetworks: expandAuthorizedNetworks(_ipConfiguration["authorized_networks"].(*schema.Set).List()),
		ForceSendFields:    []string{"Ipv4Enabled", "RequireSsl"},
	}
}
func expandAuthorizedNetworks(configured []interface{}) []*sqladmin.AclEntry {
	an := make([]*sqladmin.AclEntry, 0, len(configured))
	for _, _acl := range configured {
		_entry := _acl.(map[string]interface{})
		an = append(an, &sqladmin.AclEntry{
			ExpirationTime: _entry["expiration_time"].(string),
			Name:           _entry["name"].(string),
			Value:          _entry["value"].(string),
		})
	}

	return an
}

func expandAuthorizedGaeApplications(configured []interface{}) []string {
	aga := make([]string, 0, len(configured))
	for _, app := range configured {
		aga = append(aga, app.(string))
	}
	return aga
}

func expandDatabaseFlags(configured []interface{}) []*sqladmin.DatabaseFlags {
	databaseFlags := make([]*sqladmin.DatabaseFlags, 0, len(configured))
	for _, _flag := range configured {
		_entry := _flag.(map[string]interface{})

		databaseFlags = append(databaseFlags, &sqladmin.DatabaseFlags{
			Name:  _entry["name"].(string),
			Value: _entry["value"].(string),
		})
	}
	return databaseFlags
}

func expandBackupConfiguration(configured []interface{}) *sqladmin.BackupConfiguration {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_backupConfiguration := configured[0].(map[string]interface{})
	return &sqladmin.BackupConfiguration{
		BinaryLogEnabled: _backupConfiguration["binary_log_enabled"].(bool),
		Enabled:          _backupConfiguration["enabled"].(bool),
		StartTime:        _backupConfiguration["start_time"].(string),
		Location:         _backupConfiguration["location"].(string),
	}
}

func resourceSqlDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	var instance *sqladmin.DatabaseInstance
	err = retryTimeDuration(func() (rerr error) {
		instance, rerr = config.clientSqlAdmin.Instances.Get(project, d.Get("name").(string)).Do()
		return rerr
	}, d.Timeout(schema.TimeoutRead), isSqlOperationInProgressError)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL Database Instance %q", d.Get("name").(string)))
	}

	d.Set("name", instance.Name)
	d.Set("region", instance.Region)
	d.Set("database_version", instance.DatabaseVersion)
	d.Set("connection_name", instance.ConnectionName)
	d.Set("service_account_email_address", instance.ServiceAccountEmailAddress)

	if err := d.Set("settings", flattenSettings(instance.Settings)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance Settings")
	}

	if err := d.Set("replica_configuration", flattenReplicaConfiguration(instance.ReplicaConfiguration, d)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance Replica Configuration")
	}
	ipAddresses := flattenIpAddresses(instance.IpAddresses)
	if err := d.Set("ip_address", ipAddresses); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance IP Addresses")
	}

	if len(ipAddresses) > 0 {
		d.Set("first_ip_address", ipAddresses[0]["ip_address"])
	}

	publicIpAddress := ""
	privateIpAddress := ""
	for _, ip := range instance.IpAddresses {
		if publicIpAddress == "" && ip.Type == "PRIMARY" {
			publicIpAddress = ip.IpAddress
		}

		if privateIpAddress == "" && ip.Type == "PRIVATE" {
			privateIpAddress = ip.IpAddress
		}
	}

	d.Set("public_ip_address", publicIpAddress)
	d.Set("private_ip_address", privateIpAddress)

	if err := d.Set("server_ca_cert", flattenServerCaCert(instance.ServerCaCert)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database CA Certificate")
	}

	d.Set("master_instance_name", strings.TrimPrefix(instance.MasterInstanceName, project+":"))
	d.Set("project", project)
	d.Set("self_link", instance.SelfLink)
	d.SetId(instance.Name)

	return nil
}

func resourceSqlDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Update only updates the settings, so they are all we need to set.
	instance := &sqladmin.DatabaseInstance{
		Settings: expandSqlDatabaseInstanceSettings(d.Get("settings").([]interface{}), !isFirstGen(d)),
	}

	// Lock on the master_instance_name just in case updating any replica
	// settings causes operations on the master.
	if v, ok := d.GetOk("master_instance_name"); ok {
		mutexKV.Lock(instanceMutexKey(project, v.(string)))
		defer mutexKV.Unlock(instanceMutexKey(project, v.(string)))
	}

	var op *sqladmin.Operation
	err = retryTimeDuration(func() (rerr error) {
		op, rerr = config.clientSqlAdmin.Instances.Update(project, d.Get("name").(string), instance).Do()
		return rerr
	}, d.Timeout(schema.TimeoutUpdate), isSqlOperationInProgressError)
	if err != nil {
		return fmt.Errorf("Error, failed to update instance settings for %s: %s", instance.Name, err)
	}

	err = sqlAdminOperationWaitTime(config, op, project, "Update Instance", int(d.Timeout(schema.TimeoutUpdate).Minutes()))
	if err != nil {
		return err
	}

	return resourceSqlDatabaseInstanceRead(d, meta)
}

func resourceSqlDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Lock on the master_instance_name just in case deleting a replica causes
	// operations on the master.
	if v, ok := d.GetOk("master_instance_name"); ok {
		mutexKV.Lock(instanceMutexKey(project, v.(string)))
		defer mutexKV.Unlock(instanceMutexKey(project, v.(string)))
	}

	var op *sqladmin.Operation
	err = retryTimeDuration(func() (rerr error) {
		op, rerr = config.clientSqlAdmin.Instances.Delete(project, d.Get("name").(string)).Do()
		return rerr
	}, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error, failed to delete instance %s: %s", d.Get("name").(string), err)
	}

	err = sqlAdminOperationWaitTime(config, op, project, "Delete Instance", int(d.Timeout(schema.TimeoutDelete).Minutes()))
	if err != nil {
		return err
	}

	return nil
}

func resourceSqlDatabaseInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/instances/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenSettings(settings *sqladmin.Settings) []map[string]interface{} {
	data := map[string]interface{}{
		"version":                     settings.SettingsVersion,
		"tier":                        settings.Tier,
		"activation_policy":           settings.ActivationPolicy,
		"authorized_gae_applications": settings.AuthorizedGaeApplications,
		"availability_type":           settings.AvailabilityType,
		"crash_safe_replication":      settings.CrashSafeReplicationEnabled,
		"disk_type":                   settings.DataDiskType,
		"disk_size":                   settings.DataDiskSizeGb,
		"pricing_plan":                settings.PricingPlan,
		"replication_type":            settings.ReplicationType,
		"user_labels":                 settings.UserLabels,
	}

	if settings.BackupConfiguration != nil {
		data["backup_configuration"] = flattenBackupConfiguration(settings.BackupConfiguration)
	}

	if settings.DatabaseFlags != nil {
		data["database_flags"] = flattenDatabaseFlags(settings.DatabaseFlags)
	}

	if settings.IpConfiguration != nil {
		data["ip_configuration"] = flattenIpConfiguration(settings.IpConfiguration)
	}

	if settings.LocationPreference != nil {
		data["location_preference"] = flattenLocationPreference(settings.LocationPreference)
	}

	if settings.MaintenanceWindow != nil {
		data["maintenance_window"] = flattenMaintenanceWindow(settings.MaintenanceWindow)
	}

	if settings.StorageAutoResize != nil {
		data["disk_autoresize"] = *settings.StorageAutoResize
	}

	if settings.UserLabels != nil {
		data["user_labels"] = settings.UserLabels
	}

	return []map[string]interface{}{data}
}

func flattenBackupConfiguration(backupConfiguration *sqladmin.BackupConfiguration) []map[string]interface{} {
	data := map[string]interface{}{
		"binary_log_enabled": backupConfiguration.BinaryLogEnabled,
		"enabled":            backupConfiguration.Enabled,
		"start_time":         backupConfiguration.StartTime,
		"location":           backupConfiguration.Location,
	}

	return []map[string]interface{}{data}
}

func flattenDatabaseFlags(databaseFlags []*sqladmin.DatabaseFlags) []map[string]interface{} {
	flags := make([]map[string]interface{}, 0, len(databaseFlags))

	for _, flag := range databaseFlags {
		data := map[string]interface{}{
			"name":  flag.Name,
			"value": flag.Value,
		}

		flags = append(flags, data)
	}

	return flags
}

func flattenIpConfiguration(ipConfiguration *sqladmin.IpConfiguration) interface{} {
	data := map[string]interface{}{
		"ipv4_enabled":    ipConfiguration.Ipv4Enabled,
		"private_network": ipConfiguration.PrivateNetwork,
		"require_ssl":     ipConfiguration.RequireSsl,
	}

	if ipConfiguration.AuthorizedNetworks != nil {
		data["authorized_networks"] = flattenAuthorizedNetworks(ipConfiguration.AuthorizedNetworks)
	}

	return []map[string]interface{}{data}
}

func flattenAuthorizedNetworks(entries []*sqladmin.AclEntry) interface{} {
	networks := schema.NewSet(schema.HashResource(sqlDatabaseAuthorizedNetWorkSchemaElem), []interface{}{})

	for _, entry := range entries {
		data := map[string]interface{}{
			"expiration_time": entry.ExpirationTime,
			"name":            entry.Name,
			"value":           entry.Value,
		}

		networks.Add(data)
	}

	return networks
}

func flattenLocationPreference(locationPreference *sqladmin.LocationPreference) interface{} {
	data := map[string]interface{}{
		"follow_gae_application": locationPreference.FollowGaeApplication,
		"zone":                   locationPreference.Zone,
	}

	return []map[string]interface{}{data}
}

func flattenMaintenanceWindow(maintenanceWindow *sqladmin.MaintenanceWindow) interface{} {
	data := map[string]interface{}{
		"day":          maintenanceWindow.Day,
		"hour":         maintenanceWindow.Hour,
		"update_track": maintenanceWindow.UpdateTrack,
	}

	return []map[string]interface{}{data}
}

func flattenReplicaConfiguration(replicaConfiguration *sqladmin.ReplicaConfiguration, d *schema.ResourceData) []map[string]interface{} {
	rc := []map[string]interface{}{}

	if replicaConfiguration != nil {
		data := map[string]interface{}{
			"failover_target": replicaConfiguration.FailoverTarget,

			// Don't attempt to assign anything from replicaConfiguration.MysqlReplicaConfiguration,
			// since those fields are set on create and then not stored. See description at
			// https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/instances.
			// Instead, set them to the values they previously had so we don't set them all to zero.
			"ca_certificate":            d.Get("replica_configuration.0.ca_certificate"),
			"client_certificate":        d.Get("replica_configuration.0.client_certificate"),
			"client_key":                d.Get("replica_configuration.0.client_key"),
			"connect_retry_interval":    d.Get("replica_configuration.0.connect_retry_interval"),
			"dump_file_path":            d.Get("replica_configuration.0.dump_file_path"),
			"master_heartbeat_period":   d.Get("replica_configuration.0.master_heartbeat_period"),
			"password":                  d.Get("replica_configuration.0.password"),
			"ssl_cipher":                d.Get("replica_configuration.0.ssl_cipher"),
			"username":                  d.Get("replica_configuration.0.username"),
			"verify_server_certificate": d.Get("replica_configuration.0.verify_server_certificate"),
		}
		rc = append(rc, data)
	}

	return rc
}

func flattenIpAddresses(ipAddresses []*sqladmin.IpMapping) []map[string]interface{} {
	var ips []map[string]interface{}

	for _, ip := range ipAddresses {
		data := map[string]interface{}{
			"ip_address":     ip.IpAddress,
			"type":           ip.Type,
			"time_to_retire": ip.TimeToRetire,
		}

		ips = append(ips, data)
	}

	return ips
}

func flattenServerCaCert(caCert *sqladmin.SslCert) []map[string]interface{} {
	var cert []map[string]interface{}

	if caCert != nil {
		data := map[string]interface{}{
			"cert":             caCert.Cert,
			"common_name":      caCert.CommonName,
			"create_time":      caCert.CreateTime,
			"expiration_time":  caCert.ExpirationTime,
			"sha1_fingerprint": caCert.Sha1Fingerprint,
		}

		cert = append(cert, data)
	}

	return cert
}

func instanceMutexKey(project, instance_name string) string {
	return fmt.Sprintf("google-sql-database-instance-%s-%s", project, instance_name)
}

// sqlDatabaseIsMaster returns true if the provided schema.ResourceData represents a
// master SQL Instance, and false if it is a replica.
func sqlDatabaseIsMaster(d *schema.ResourceData) bool {
	_, ok := d.GetOk("master_instance_name")
	return !ok
}
