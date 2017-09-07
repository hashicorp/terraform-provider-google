package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sqladmin/v1beta4"
)

func resourceSqlDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlDatabaseInstanceCreate,
		Read:   resourceSqlDatabaseInstanceRead,
		Update: resourceSqlDatabaseInstanceUpdate,
		Delete: resourceSqlDatabaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"settings": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tier": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"activation_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							// Defaults differ between first and second gen instances
							Computed: true,
						},
						"authorized_gae_applications": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_configuration": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"binary_log_enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"start_time": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										// start_time is randomly assigned if not set
										Computed: true,
									},
								},
							},
						},
						"crash_safe_replication": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"database_flags": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"disk_autoresize": &schema.Schema{
							Type:             schema.TypeBool,
							Optional:         true,
							DiffSuppressFunc: suppressFirstGen,
							// Set computed instead of default because this property is for second-gen only.
							Computed: true,
						},
						"disk_size": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							// Defaults differ between first and second gen instances
							Computed: true,
						},
						"disk_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							// Set computed instead of default because this property is for second-gen only.
							Computed: true,
						},
						"ip_configuration": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorized_networks": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"expiration_time": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"name": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"ipv4_enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										// Defaults differ between first and second gen instances
										Computed: true,
									},
									"require_ssl": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"location_preference": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"follow_gae_application": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"zone": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"maintenance_window": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": &schema.Schema{
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 7),
									},
									"hour": &schema.Schema{
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},
									"update_track": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"pricing_plan": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "PER_USE",
						},
						"replication_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "SYNCHRONOUS",
						},
					},
				},
			},

			"database_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MYSQL_5_6",
				ForceNew: true,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_to_retire": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"master_instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"replica_configuration": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// Returned from API on all replicas
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ca_certificate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"client_certificate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"client_key": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"connect_retry_interval": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"dump_file_path": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"failover_target": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"master_heartbeat_period": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ssl_cipher": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"verify_server_certificate": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Suppress diff with any disk_autoresize value on 1st Generation Instances
func suppressFirstGen(k, old, new string, d *schema.ResourceData) bool {
	settingsList := d.Get("settings").([]interface{})

	settings := settingsList[0].(map[string]interface{})
	tier := settings["tier"].(string)
	matched, err := regexp.MatchString("db*", tier)
	if err != nil {
		log.Printf("[ERR] error with regex in diff supression for disk_autoresize: %s", err)
	}
	if !matched {
		log.Printf("[DEBUG] suppressing diff on disk_autoresize due to 1st gen instance type")
		return true
	}
	return false
}

func resourceSqlDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	databaseVersion := d.Get("database_version").(string)

	_settingsList := d.Get("settings").([]interface{})

	_settings := _settingsList[0].(map[string]interface{})
	settings := &sqladmin.Settings{
		Tier:            _settings["tier"].(string),
		ForceSendFields: []string{"StorageAutoResize"},
	}

	if v, ok := _settings["activation_policy"]; ok {
		settings.ActivationPolicy = v.(string)
	}

	if v, ok := _settings["authorized_gae_applications"]; ok {
		settings.AuthorizedGaeApplications = make([]string, 0)
		for _, app := range v.([]interface{}) {
			settings.AuthorizedGaeApplications = append(settings.AuthorizedGaeApplications,
				app.(string))
		}
	}

	if v, ok := _settings["backup_configuration"]; ok {
		_backupConfigurationList := v.([]interface{})

		if len(_backupConfigurationList) == 1 && _backupConfigurationList[0] != nil {
			settings.BackupConfiguration = &sqladmin.BackupConfiguration{}
			_backupConfiguration := _backupConfigurationList[0].(map[string]interface{})

			if vp, okp := _backupConfiguration["binary_log_enabled"]; okp {
				settings.BackupConfiguration.BinaryLogEnabled = vp.(bool)
			}

			if vp, okp := _backupConfiguration["enabled"]; okp {
				settings.BackupConfiguration.Enabled = vp.(bool)
			}

			if vp, okp := _backupConfiguration["start_time"]; okp {
				settings.BackupConfiguration.StartTime = vp.(string)
			}
		}
	}

	if v, ok := _settings["crash_safe_replication"]; ok {
		settings.CrashSafeReplicationEnabled = v.(bool)
	}

	settings.StorageAutoResize = _settings["disk_autoresize"].(bool)

	if v, ok := _settings["disk_size"]; ok && v.(int) > 0 {
		settings.DataDiskSizeGb = int64(v.(int))
	}

	if v, ok := _settings["disk_type"]; ok && len(v.(string)) > 0 {
		settings.DataDiskType = v.(string)
	}

	if v, ok := _settings["database_flags"]; ok {
		settings.DatabaseFlags = make([]*sqladmin.DatabaseFlags, 0)
		_databaseFlagsList := v.([]interface{})
		for _, _flag := range _databaseFlagsList {
			_entry := _flag.(map[string]interface{})
			flag := &sqladmin.DatabaseFlags{}
			if vp, okp := _entry["name"]; okp {
				flag.Name = vp.(string)
			}

			if vp, okp := _entry["value"]; okp {
				flag.Value = vp.(string)
			}

			settings.DatabaseFlags = append(settings.DatabaseFlags, flag)
		}
	}

	if v, ok := _settings["ip_configuration"]; ok {
		_ipConfigurationList := v.([]interface{})

		if len(_ipConfigurationList) == 1 && _ipConfigurationList[0] != nil {
			settings.IpConfiguration = &sqladmin.IpConfiguration{}
			_ipConfiguration := _ipConfigurationList[0].(map[string]interface{})

			if vp, okp := _ipConfiguration["ipv4_enabled"]; okp {
				settings.IpConfiguration.Ipv4Enabled = vp.(bool)
			}

			if vp, okp := _ipConfiguration["require_ssl"]; okp {
				settings.IpConfiguration.RequireSsl = vp.(bool)
			}

			if vp, okp := _ipConfiguration["authorized_networks"]; okp {
				settings.IpConfiguration.AuthorizedNetworks = make([]*sqladmin.AclEntry, 0)
				_authorizedNetworksList := vp.([]interface{})
				for _, _acl := range _authorizedNetworksList {
					_entry := _acl.(map[string]interface{})
					entry := &sqladmin.AclEntry{}

					if vpp, okpp := _entry["expiration_time"]; okpp {
						entry.ExpirationTime = vpp.(string)
					}

					if vpp, okpp := _entry["name"]; okpp {
						entry.Name = vpp.(string)
					}

					if vpp, okpp := _entry["value"]; okpp {
						entry.Value = vpp.(string)
					}

					settings.IpConfiguration.AuthorizedNetworks = append(
						settings.IpConfiguration.AuthorizedNetworks, entry)
				}
			}
		}
	}

	if v, ok := _settings["location_preference"]; ok {
		_locationPreferenceList := v.([]interface{})

		if len(_locationPreferenceList) == 1 && _locationPreferenceList[0] != nil {
			settings.LocationPreference = &sqladmin.LocationPreference{}
			_locationPreference := _locationPreferenceList[0].(map[string]interface{})

			if vp, okp := _locationPreference["follow_gae_application"]; okp {
				settings.LocationPreference.FollowGaeApplication = vp.(string)
			}

			if vp, okp := _locationPreference["zone"]; okp {
				settings.LocationPreference.Zone = vp.(string)
			}
		}
	}

	if v, ok := _settings["maintenance_window"]; ok && len(v.([]interface{})) > 0 {
		settings.MaintenanceWindow = &sqladmin.MaintenanceWindow{}
		_maintenanceWindow := v.([]interface{})[0].(map[string]interface{})

		if vp, okp := _maintenanceWindow["day"]; okp {
			settings.MaintenanceWindow.Day = int64(vp.(int))
		}

		if vp, okp := _maintenanceWindow["hour"]; okp {
			settings.MaintenanceWindow.Hour = int64(vp.(int))
		}

		if vp, ok := _maintenanceWindow["update_track"]; ok {
			if len(vp.(string)) > 0 {
				settings.MaintenanceWindow.UpdateTrack = vp.(string)
			}
		}
	}

	if v, ok := _settings["pricing_plan"]; ok {
		settings.PricingPlan = v.(string)
	}

	if v, ok := _settings["replication_type"]; ok {
		settings.ReplicationType = v.(string)
	}

	instance := &sqladmin.DatabaseInstance{
		Region:          region,
		Settings:        settings,
		DatabaseVersion: databaseVersion,
	}

	if v, ok := d.GetOk("name"); ok {
		instance.Name = v.(string)
	} else {
		instance.Name = resource.UniqueId()
		d.Set("name", instance.Name)
	}

	if v, ok := d.GetOk("replica_configuration"); ok {
		_replicaConfigurationList := v.([]interface{})

		if len(_replicaConfigurationList) == 1 && _replicaConfigurationList[0] != nil {
			replicaConfiguration := &sqladmin.ReplicaConfiguration{}
			mySqlReplicaConfiguration := &sqladmin.MySqlReplicaConfiguration{}
			_replicaConfiguration := _replicaConfigurationList[0].(map[string]interface{})

			if vp, okp := _replicaConfiguration["failover_target"]; okp {
				replicaConfiguration.FailoverTarget = vp.(bool)
			}

			if vp, okp := _replicaConfiguration["ca_certificate"]; okp {
				mySqlReplicaConfiguration.CaCertificate = vp.(string)
			}

			if vp, okp := _replicaConfiguration["client_certificate"]; okp {
				mySqlReplicaConfiguration.ClientCertificate = vp.(string)
			}

			if vp, okp := _replicaConfiguration["client_key"]; okp {
				mySqlReplicaConfiguration.ClientKey = vp.(string)
			}

			if vp, okp := _replicaConfiguration["connect_retry_interval"]; okp {
				mySqlReplicaConfiguration.ConnectRetryInterval = int64(vp.(int))
			}

			if vp, okp := _replicaConfiguration["dump_file_path"]; okp {
				mySqlReplicaConfiguration.DumpFilePath = vp.(string)
			}

			if vp, okp := _replicaConfiguration["master_heartbeat_period"]; okp {
				mySqlReplicaConfiguration.MasterHeartbeatPeriod = int64(vp.(int))
			}

			if vp, okp := _replicaConfiguration["password"]; okp {
				mySqlReplicaConfiguration.Password = vp.(string)
			}

			if vp, okp := _replicaConfiguration["ssl_cipher"]; okp {
				mySqlReplicaConfiguration.SslCipher = vp.(string)
			}

			if vp, okp := _replicaConfiguration["username"]; okp {
				mySqlReplicaConfiguration.Username = vp.(string)
			}

			if vp, okp := _replicaConfiguration["verify_server_certificate"]; okp {
				mySqlReplicaConfiguration.VerifyServerCertificate = vp.(bool)
			}

			replicaConfiguration.MysqlReplicaConfiguration = mySqlReplicaConfiguration
			instance.ReplicaConfiguration = replicaConfiguration
		}
	}

	if v, ok := d.GetOk("master_instance_name"); ok {
		instance.MasterInstanceName = v.(string)
	}

	op, err := config.clientSqlAdmin.Instances.Insert(project, instance).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 {
			return fmt.Errorf("Error, the name %s is unavailable because it was used recently", instance.Name)
		} else {
			return fmt.Errorf("Error, failed to create instance %s: %s", instance.Name, err)
		}
	}

	d.SetId(instance.Name)

	err = sqladminOperationWait(config, op, "Create Instance")
	if err != nil {
		d.SetId("")
		return err
	}

	err = resourceSqlDatabaseInstanceRead(d, meta)
	if err != nil {
		return err
	}

	// If a root user exists with a wildcard ('%') hostname, delete it.
	users, err := config.clientSqlAdmin.Users.List(project, instance.Name).Do()
	if err != nil {
		return fmt.Errorf("Error, attempting to list users associated with instance %s: %s", instance.Name, err)
	}
	for _, u := range users.Items {
		if u.Name == "root" && u.Host == "%" {
			op, err = config.clientSqlAdmin.Users.Delete(project, instance.Name, u.Host, u.Name).Do()
			if err != nil {
				return fmt.Errorf("Error, failed to delete default 'root'@'*' user, but the database was created successfully: %s", err)
			}
			err = sqladminOperationWait(config, op, "Delete default root User")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceSqlDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance, err := config.clientSqlAdmin.Instances.Get(project,
		d.Id()).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL Database Instance %q", d.Get("name").(string)))
	}

	d.Set("name", instance.Name)
	d.Set("region", instance.Region)
	d.Set("database_version", instance.DatabaseVersion)

	if err := d.Set("settings", flattenSettings(instance.Settings)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance Settings")
	}
	if err := d.Set("replica_configuration", flattenReplicaConfiguration(instance.ReplicaConfiguration)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance Replica Configuration")
	}
	if err := d.Set("ip_address", flattenIpAddresses(instance.IpAddresses)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance IP Addresses")
	}

	d.Set("master_instance_name", strings.TrimPrefix(instance.MasterInstanceName, project+":"))

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

	d.Partial(true)

	instance, err := config.clientSqlAdmin.Instances.Get(project,
		d.Get("name").(string)).Do()

	if err != nil {
		return fmt.Errorf("Error retrieving instance %s: %s",
			d.Get("name").(string), err)
	}

	if d.HasChange("settings") {
		_oListCast, _settingsListCast := d.GetChange("settings")
		_oList := _oListCast.([]interface{})
		_o := _oList[0].(map[string]interface{})
		_settingsList := _settingsListCast.([]interface{})

		_settings := _settingsList[0].(map[string]interface{})
		settings := &sqladmin.Settings{
			Tier:              _settings["tier"].(string),
			SettingsVersion:   instance.Settings.SettingsVersion,
			StorageAutoResize: _settings["disk_autoresize"].(bool),
			ForceSendFields:   []string{"StorageAutoResize"},
		}

		if v, ok := _settings["activation_policy"]; ok {
			settings.ActivationPolicy = v.(string)
		}

		if v, ok := _settings["authorized_gae_applications"]; ok {
			settings.AuthorizedGaeApplications = make([]string, 0)
			for _, app := range v.([]interface{}) {
				settings.AuthorizedGaeApplications = append(settings.AuthorizedGaeApplications,
					app.(string))
			}
		}

		if v, ok := _settings["backup_configuration"]; ok {
			_backupConfigurationList := v.([]interface{})

			settings.BackupConfiguration = &sqladmin.BackupConfiguration{}
			if len(_backupConfigurationList) == 1 && _backupConfigurationList[0] != nil {
				_backupConfiguration := _backupConfigurationList[0].(map[string]interface{})

				if vp, okp := _backupConfiguration["binary_log_enabled"]; okp {
					settings.BackupConfiguration.BinaryLogEnabled = vp.(bool)
				}

				if vp, okp := _backupConfiguration["enabled"]; okp {
					settings.BackupConfiguration.Enabled = vp.(bool)
				}

				if vp, okp := _backupConfiguration["start_time"]; okp {
					settings.BackupConfiguration.StartTime = vp.(string)
				}
			}
		}

		if v, ok := _settings["crash_safe_replication"]; ok {
			settings.CrashSafeReplicationEnabled = v.(bool)
		}

		if v, ok := _settings["disk_size"]; ok {
			if v.(int) > 0 && int64(v.(int)) > instance.Settings.DataDiskSizeGb {
				settings.DataDiskSizeGb = int64(v.(int))
			}
		}

		if v, ok := _settings["disk_type"]; ok && len(v.(string)) > 0 {
			settings.DataDiskType = v.(string)
		}

		_oldDatabaseFlags := make([]interface{}, 0)
		if ov, ook := _o["database_flags"]; ook {
			_oldDatabaseFlags = ov.([]interface{})
		}

		if v, ok := _settings["database_flags"]; ok || len(_oldDatabaseFlags) > 0 {
			oldDatabaseFlags := settings.DatabaseFlags
			settings.DatabaseFlags = make([]*sqladmin.DatabaseFlags, 0)
			_databaseFlagsList := make([]interface{}, 0)
			if v != nil {
				_databaseFlagsList = v.([]interface{})
			}

			_odbf_map := make(map[string]interface{})
			for _, _dbf := range _oldDatabaseFlags {
				_entry := _dbf.(map[string]interface{})
				_odbf_map[_entry["name"].(string)] = true
			}

			// First read the flags from the server, and reinsert those that
			// were not previously defined
			for _, entry := range oldDatabaseFlags {
				_, ok_old := _odbf_map[entry.Name]
				if !ok_old {
					settings.DatabaseFlags = append(
						settings.DatabaseFlags, entry)
				}
			}
			// finally, insert only those that were previously defined
			// and are still defined.
			for _, _flag := range _databaseFlagsList {
				_entry := _flag.(map[string]interface{})
				flag := &sqladmin.DatabaseFlags{}
				if vp, okp := _entry["name"]; okp {
					flag.Name = vp.(string)
				}

				if vp, okp := _entry["value"]; okp {
					flag.Value = vp.(string)
				}

				settings.DatabaseFlags = append(settings.DatabaseFlags, flag)
			}
		}

		if v, ok := _settings["ip_configuration"]; ok {
			_ipConfigurationList := v.([]interface{})

			settings.IpConfiguration = &sqladmin.IpConfiguration{}
			if len(_ipConfigurationList) == 1 && _ipConfigurationList[0] != nil {
				_ipConfiguration := _ipConfigurationList[0].(map[string]interface{})

				if vp, okp := _ipConfiguration["ipv4_enabled"]; okp {
					settings.IpConfiguration.Ipv4Enabled = vp.(bool)
				}

				if vp, okp := _ipConfiguration["require_ssl"]; okp {
					settings.IpConfiguration.RequireSsl = vp.(bool)
				}

				_oldAuthorizedNetworkList := make([]interface{}, 0)
				if ov, ook := _o["ip_configuration"]; ook {
					_oldIpConfList := ov.([]interface{})
					if len(_oldIpConfList) > 0 {
						_oldIpConf := _oldIpConfList[0].(map[string]interface{})
						if ovp, ookp := _oldIpConf["authorized_networks"]; ookp {
							_oldAuthorizedNetworkList = ovp.([]interface{})
						}
					}
				}

				if vp, okp := _ipConfiguration["authorized_networks"]; okp || len(_oldAuthorizedNetworkList) > 0 {
					oldAuthorizedNetworks := instance.Settings.IpConfiguration.AuthorizedNetworks
					settings.IpConfiguration.AuthorizedNetworks = make([]*sqladmin.AclEntry, 0)

					_authorizedNetworksList := make([]interface{}, 0)
					if vp != nil {
						_authorizedNetworksList = vp.([]interface{})
					}
					_oipc_map := make(map[string]interface{})
					for _, _ipc := range _oldAuthorizedNetworkList {
						_entry := _ipc.(map[string]interface{})
						_oipc_map[_entry["value"].(string)] = true
					}
					// Next read the network tuples from the server, and reinsert those that
					// were not previously defined
					for _, entry := range oldAuthorizedNetworks {
						_, ok_old := _oipc_map[entry.Value]
						if !ok_old {
							settings.IpConfiguration.AuthorizedNetworks = append(
								settings.IpConfiguration.AuthorizedNetworks, entry)
						}
					}
					// finally, update old entries and insert new ones
					// and are still defined.
					for _, _ipc := range _authorizedNetworksList {
						_entry := _ipc.(map[string]interface{})
						entry := &sqladmin.AclEntry{}

						if vpp, okpp := _entry["expiration_time"]; okpp {
							entry.ExpirationTime = vpp.(string)
						}

						if vpp, okpp := _entry["name"]; okpp {
							entry.Name = vpp.(string)
						}

						if vpp, okpp := _entry["value"]; okpp {
							entry.Value = vpp.(string)
						}

						settings.IpConfiguration.AuthorizedNetworks = append(
							settings.IpConfiguration.AuthorizedNetworks, entry)
					}
				}
			}
		}

		if v, ok := _settings["location_preference"]; ok {
			_locationPreferenceList := v.([]interface{})

			settings.LocationPreference = &sqladmin.LocationPreference{}
			if len(_locationPreferenceList) == 1 && _locationPreferenceList[0] != nil {
				_locationPreference := _locationPreferenceList[0].(map[string]interface{})

				if vp, okp := _locationPreference["follow_gae_application"]; okp {
					settings.LocationPreference.FollowGaeApplication = vp.(string)
				}

				if vp, okp := _locationPreference["zone"]; okp {
					settings.LocationPreference.Zone = vp.(string)
				}
			}
		}

		if v, ok := _settings["maintenance_window"]; ok && len(v.([]interface{})) > 0 {
			_maintenanceWindowList := v.([]interface{})

			settings.MaintenanceWindow = &sqladmin.MaintenanceWindow{}
			if len(_maintenanceWindowList) == 1 && _maintenanceWindowList[0] != nil {
				_maintenanceWindow := _maintenanceWindowList[0].(map[string]interface{})

				if vp, okp := _maintenanceWindow["day"]; okp {
					settings.MaintenanceWindow.Day = int64(vp.(int))
				}

				if vp, okp := _maintenanceWindow["hour"]; okp {
					settings.MaintenanceWindow.Hour = int64(vp.(int))
				}

				if vp, ok := _maintenanceWindow["update_track"]; ok {
					if len(vp.(string)) > 0 {
						settings.MaintenanceWindow.UpdateTrack = vp.(string)
					}
				}
			}
		}

		if v, ok := _settings["pricing_plan"]; ok {
			settings.PricingPlan = v.(string)
		}

		if v, ok := _settings["replication_type"]; ok {
			settings.ReplicationType = v.(string)
		}

		instance.Settings = settings
	}

	d.Partial(false)

	op, err := config.clientSqlAdmin.Instances.Update(project, instance.Name, instance).Do()
	if err != nil {
		return fmt.Errorf("Error, failed to update instance %s: %s", instance.Name, err)
	}

	err = sqladminOperationWait(config, op, "Create Instance")
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

	op, err := config.clientSqlAdmin.Instances.Delete(project, d.Get("name").(string)).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to delete instance %s: %s", d.Get("name").(string), err)
	}

	err = sqladminOperationWait(config, op, "Delete Instance")
	if err != nil {
		return err
	}

	return nil
}

func flattenSettings(settings *sqladmin.Settings) []map[string]interface{} {
	data := map[string]interface{}{
		"version":                     settings.SettingsVersion,
		"tier":                        settings.Tier,
		"activation_policy":           settings.ActivationPolicy,
		"authorized_gae_applications": settings.AuthorizedGaeApplications,
		"crash_safe_replication":      settings.CrashSafeReplicationEnabled,
		"disk_autoresize":             settings.StorageAutoResize,
		"disk_type":                   settings.DataDiskType,
		"disk_size":                   settings.DataDiskSizeGb,
		"pricing_plan":                settings.PricingPlan,
		"replication_type":            settings.ReplicationType,
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

	return []map[string]interface{}{data}
}

func flattenBackupConfiguration(backupConfiguration *sqladmin.BackupConfiguration) []map[string]interface{} {
	data := map[string]interface{}{
		"binary_log_enabled": backupConfiguration.BinaryLogEnabled,
		"enabled":            backupConfiguration.Enabled,
		"start_time":         backupConfiguration.StartTime,
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
		"ipv4_enabled": ipConfiguration.Ipv4Enabled,
		"require_ssl":  ipConfiguration.RequireSsl,
	}

	if ipConfiguration.AuthorizedNetworks != nil {
		data["authorized_networks"] = flattenAuthorizedNetworks(ipConfiguration.AuthorizedNetworks)
	}

	return []map[string]interface{}{data}
}

func flattenAuthorizedNetworks(entries []*sqladmin.AclEntry) interface{} {
	networks := make([]map[string]interface{}, 0, len(entries))

	for _, entry := range entries {
		data := map[string]interface{}{
			"expiration_time": entry.ExpirationTime,
			"name":            entry.Name,
			"value":           entry.Value,
		}

		networks = append(networks, data)
	}

	return networks
}

func flattenLocationPreference(locationPreference *sqladmin.LocationPreference) interface{} {
	data := map[string]interface{}{
		"follow_gae_application": locationPreference.FollowGaeApplication,
		"zone": locationPreference.Zone,
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

func flattenReplicaConfiguration(replicaConfiguration *sqladmin.ReplicaConfiguration) []map[string]interface{} {
	rc := []map[string]interface{}{}

	if replicaConfiguration != nil {
		data := map[string]interface{}{
			"failover_target": replicaConfiguration.FailoverTarget,

			// Don't attempt to assign anything from replicaConfiguration.MysqlReplicaConfiguration,
			// since those fields are set on create and then not stored. See description at
			// https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/instances
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
			"time_to_retire": ip.TimeToRetire,
		}

		ips = append(ips, data)
	}

	return ips
}

func instanceMutexKey(project, instance_name string) string {
	return fmt.Sprintf("google-sql-database-instance-%s-%s", project, instance_name)
}
