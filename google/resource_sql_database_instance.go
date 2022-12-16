package google

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// Match fully-qualified or relative URLs
const privateNetworkLinkRegex = "^(?:http(?:s)?://.+/)?projects/(" + ProjectRegex + ")/global/networks/((?:[a-z](?:[-a-z0-9]*[a-z0-9])?))$"

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
		"settings.0.backup_configuration.0.point_in_time_recovery_enabled",
		"settings.0.backup_configuration.0.backup_retention_settings",
		"settings.0.backup_configuration.0.transaction_log_retention_days",
	}

	ipConfigurationKeys = []string{
		"settings.0.ip_configuration.0.authorized_networks",
		"settings.0.ip_configuration.0.ipv4_enabled",
		"settings.0.ip_configuration.0.require_ssl",
		"settings.0.ip_configuration.0.private_network",
		"settings.0.ip_configuration.0.allocated_ip_range",
	}

	maintenanceWindowKeys = []string{
		"settings.0.maintenance_window.0.day",
		"settings.0.maintenance_window.0.hour",
		"settings.0.maintenance_window.0.update_track",
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

	insightsConfigKeys = []string{
		"settings.0.insights_config.0.query_insights_enabled",
		"settings.0.insights_config.0.query_string_length",
		"settings.0.insights_config.0.record_application_tags",
		"settings.0.insights_config.0.record_client_address",
		"settings.0.insights_config.0.query_plans_per_minute",
	}

	sqlServerAuditConfigurationKeys = []string{
		"settings.0.sql_server_audit_config.0.bucket",
		"settings.0.sql_server_audit_config.0.retention_interval",
		"settings.0.sql_server_audit_config.0.upload_interval",
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
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("settings.0.disk_size", isDiskShrinkage),
			privateNetworkCustomizeDiff,
			pitrPostgresOnlyCustomizeDiff,
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region the instance will sit in. Note, Cloud SQL is not available in all regions. A valid region must be provided to use this resource. If a region is not provided in the resource definition, the provider region will be used instead, but this will be an apply-time error for instances if the provider region is not supported with Cloud SQL. If you choose not to provide the region argument for this resource, make sure you understand this.`,
			},
			"deletion_protection": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: `Used to block Terraform from deleting a SQL Instance. Defaults to true.`,
			},
			"settings": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"settings", "clone"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Used to make sure changes to the settings block are atomic.`,
						},
						"tier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The machine type to use. See tiers for more details and supported versions. Postgres supports only shared-core machine types, and custom machine types such as db-custom-2-13312. See the Custom Machine Type Documentation to learn about specifying custom machine types.`,
						},
						"activation_policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "ALWAYS",
							Description: `This specifies when the instance should be active. Can be either ALWAYS, NEVER or ON_DEMAND.`,
						},
						"active_directory_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Domain name of the Active Directory for SQL Server (e.g., mydomain.com).`,
									},
								},
							},
						},
						"deny_maintenance_period": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_date": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `End date before which maintenance will not take place. The date is in format yyyy-mm-dd i.e., 2020-11-01, or mm-dd, i.e., 11-01`,
									},
									"start_date": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Start date after which maintenance will not take place. The date is in format yyyy-mm-dd i.e., 2020-11-01, or mm-dd, i.e., 11-01`,
									},
									"time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Time in UTC when the "deny maintenance period" starts on start_date and ends on end_date. The time is in format: HH:mm:SS, i.e., 00:00:00`,
									},
								},
							},
						},
						"sql_server_audit_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: sqlServerAuditConfigurationKeys,
										Description:  `The name of the destination bucket (e.g., gs://mybucket).`,
									},
									"retention_interval": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: sqlServerAuditConfigurationKeys,
										Description:  `How long to keep generated audit files. A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s"..`,
									},
									"upload_interval": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: sqlServerAuditConfigurationKeys,
										Description:  `How often to upload generated audit files. A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".`,
									},
								},
							},
						},
						"time_zone": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: `The time_zone to be used by the database engine (supported only for SQL Server), in SQL Server timezone format.`,
						},
						"availability_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "ZONAL",
							ValidateFunc: validation.StringInSlice([]string{"REGIONAL", "ZONAL"}, false),
							Description: `The availability type of the Cloud SQL instance, high availability
(REGIONAL) or single zone (ZONAL). For all instances, ensure that
settings.backup_configuration.enabled is set to true.
For MySQL instances, ensure that settings.backup_configuration.binary_log_enabled is set to true.
For Postgres instances, ensure that settings.backup_configuration.point_in_time_recovery_enabled
is set to true. Defaults to ZONAL.`,
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
										Description:  `True if binary logging is enabled. If settings.backup_configuration.enabled is false, this must be as well. Can only be used with MySQL.`,
									},
									"enabled": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
										Description:  `True if backup configuration is enabled.`,
									},
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
										// start_time is randomly assigned if not set
										Computed:     true,
										AtLeastOneOf: backupConfigurationKeys,
										Description:  `HH:MM format time indicating when backup configuration starts.`,
									},
									"location": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
										Description:  `Location of the backup configuration.`,
									},
									"point_in_time_recovery_enabled": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
										Description:  `True if Point-in-time recovery is enabled.`,
									},
									"transaction_log_retention_days": {
										Type:         schema.TypeInt,
										Computed:     true,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
										Description:  `The number of days of transaction logs we retain for point in time restore, from 1-7.`,
									},
									"backup_retention_settings": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: backupConfigurationKeys,
										Computed:     true,
										MaxItems:     1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retained_backups": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `Number of backups to retain.`,
												},
												"retention_unit": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "COUNT",
													Description: `The unit that 'retainedBackups' represents. Defaults to COUNT`,
												},
											},
										},
									},
								},
							},
						},
						"collation": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The name of server instance collation.`,
						},
						"database_flags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Value of the flag.`,
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Name of the flag.`,
									},
								},
							},
						},
						"disk_autoresize": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Enables auto-resizing of the storage size. Defaults to true.`,
						},
						"disk_autoresize_limit": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: `The maximum size, in GB, to which storage capacity can be automatically increased. The default value is 0, which specifies that there is no limit.`,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
							// Default is likely 10gb, but it is undocumented and may change.
							Computed:    true,
							Description: `The size of data disk, in GB. Size of a running instance cannot be reduced but can be increased. The minimum value is 10GB.`,
						},
						"disk_type": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "PD_SSD",
							DiffSuppressFunc: caseDiffDashSuppress,
							Description:      `The type of data disk: PD_SSD or PD_HDD. Defaults to PD_SSD.`,
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
										Type:         schema.TypeBool,
										Optional:     true,
										Default:      true,
										AtLeastOneOf: ipConfigurationKeys,
										Description:  `Whether this Cloud SQL instance should be assigned a public IPV4 address. At least ipv4_enabled must be enabled or a private_network must be configured.`,
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
										Description:      `The VPC network from which the Cloud SQL instance is accessible for private IP. For example, projects/myProject/global/networks/default. Specifying a network enables private IP. At least ipv4_enabled must be enabled or a private_network must be configured. This setting can be updated, but it cannot be removed after it is set.`,
									},
									"allocated_ip_range": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: ipConfigurationKeys,
										Description:  `The name of the allocated ip range for the private ip CloudSQL instance. For example: "google-managed-services-default". If set, the instance ip will be created in the allocated range. The range name must comply with RFC 1035. Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])?.`,
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
										Description:  `A Google App Engine application whose zone to remain in. Must be in the same region as this instance.`,
									},
									"zone": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: []string{"settings.0.location_preference.0.follow_gae_application", "settings.0.location_preference.0.zone"},
										Description:  `The preferred compute engine zone.`,
									},
									"secondary_zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The preferred Compute Engine zone for the secondary/failover`,
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
										Description:  `Day of week (1-7), starting on Monday`,
									},
									"hour": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 23),
										AtLeastOneOf: maintenanceWindowKeys,
										Description:  `Hour of day (0-23), ignored if day not set`,
									},
									"update_track": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: maintenanceWindowKeys,
										Description:  `Receive updates earlier (canary) or later (stable)`,
									},
								},
							},
							Description: `Declares a one-hour maintenance window when an Instance can automatically restart to apply updates. The maintenance window is specified in UTC time.`,
						},
						"pricing_plan": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "PER_USE",
							Description: `Pricing plan for this instance, can only be PER_USE.`,
						},
						"user_labels": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `A set of key/value user label pairs to assign to the instance.`,
						},
						"insights_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"query_insights_enabled": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: insightsConfigKeys,
										Description:  `True if Query Insights feature is enabled.`,
									},
									"query_string_length": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1024,
										ValidateFunc: validation.IntBetween(256, 4500),
										AtLeastOneOf: insightsConfigKeys,
										Description:  `Maximum query length stored in bytes. Between 256 and 4500. Default to 1024.`,
									},
									"record_application_tags": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: insightsConfigKeys,
										Description:  `True if Query Insights will record application tags from query when enabled.`,
									},
									"record_client_address": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: insightsConfigKeys,
										Description:  `True if Query Insights will record client address when enabled.`,
									},
									"query_plans_per_minute": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IntBetween(0, 20),
										AtLeastOneOf: insightsConfigKeys,
										Description:  `Number of query execution plans captured by Insights per minute for all queries combined. Between 0 and 20. Default to 5.`,
									},
								},
							},
							Description: `Configuration of Query Insights.`,
						},
						"password_validation_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_length": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 2147483647),
										Description:  `Minimum number of characters allowed.`,
									},
									"complexity": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"COMPLEXITY_DEFAULT", "COMPLEXITY_UNSPECIFIED"}, false),
										Description:  `Password complexity.`,
									},
									"reuse_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 2147483647),
										Description:  `Number of previous passwords that cannot be reused.`,
									},
									"disallow_username_substring": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Disallow username as a part of the password.`,
									},
									"password_change_interval": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Minimum interval after which the password can be changed. This flag is only supported for PostgresSQL.`,
									},
									"enable_password_policy": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether the password policy is enabled or not.`,
									},
								},
							},
						},
						"connector_enforcement": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"NOT_REQUIRED", "REQUIRED"}, false),
							Description:  `Specifies if connections must use Cloud SQL connectors.`,
						},
						"deletion_protection_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Configuration to protect against accidental instance deletion.`,
						},
					},
				},
				Description: `The settings to use for the database. The configuration is detailed below.`,
			},

			"connection_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection name of the instance to be used in connection strings. For example, when connecting with Cloud SQL Proxy.`,
			},
			"maintenance_version": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				Description:      `Maintenance version.`,
				DiffSuppressFunc: maintenanceVersionDiffSuppress,
			},
			"available_maintenance_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: `Available Maintenance versions.`,
			},
			"database_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The MySQL, PostgreSQL or SQL Server (beta) version to use. Supported values include MYSQL_5_6, MYSQL_5_7, MYSQL_8_0, POSTGRES_9_6, POSTGRES_10, POSTGRES_11, POSTGRES_12, POSTGRES_13, POSTGRES_14, SQLSERVER_2017_STANDARD, SQLSERVER_2017_ENTERPRISE, SQLSERVER_2017_EXPRESS, SQLSERVER_2017_WEB. Database Version Policies includes an up-to-date reference of supported versions.`,
			},

			"encryption_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"root_password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Initial root password. Required for MS SQL Server.`,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The first IPv4 address of any type assigned. This is to support accessing the first address in the list in a terraform output when the resource is configured with a count.`,
			},

			"public_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `IPv4 address assigned. This is a workaround for an issue fixed in Terraform 0.12 but also provides a convenient way to access an IP of a specific type without performing filtering in a Terraform config.`,
			},

			"private_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `IPv4 address assigned. This is a workaround for an issue fixed in Terraform 0.12 but also provides a convenient way to access an IP of a specific type without performing filtering in a Terraform config.`,
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The name of the instance. If the name is left blank, Terraform will randomly generate one when the instance is first created. This is done because after a name is used, it cannot be reused for up to one week.`,
			},

			"master_instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The name of the instance that will act as the master in the replication setup. Note, this requires the master to have binary_log_enabled set, as well as existing backups.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
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
							Description:  `PEM representation of the trusted CA's x509 certificate.`,
						},
						"client_certificate": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `PEM representation of the replica's x509 certificate.`,
						},
						"client_key": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `PEM representation of the replica's private key. The corresponding public key in encoded in the client_certificate.`,
						},
						"connect_retry_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `The number of seconds between connect retries. MySQL's default is 60 seconds.`,
						},
						"dump_file_path": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `Path to a SQL file in Google Cloud Storage from which replica instances are created. Format is gs://bucket/filename.`,
						},
						"failover_target": {
							Type:         schema.TypeBool,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `Specifies if the replica is the failover target. If the field is set to true the replica will be designated as a failover replica. If the master instance fails, the replica instance will be promoted as the new master instance.`,
						},
						"master_heartbeat_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `Time in ms between replication heartbeats.`,
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Sensitive:    true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `Password for the replication connection.`,
						},
						"ssl_cipher": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `Permissible ciphers for use in SSL encryption.`,
						},
						"username": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `Username for replication connection.`,
						},
						"verify_server_certificate": {
							Type:         schema.TypeBool,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: replicaConfigurationKeys,
							Description:  `True if the master's common name value is checked during the SSL handshake.`,
						},
					},
				},
				Description: `The configuration for replication.`,
			},
			"server_ca_cert": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CA Certificate used to connect to the SQL Instance via SSL.`,
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CN valid for the CA Cert.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Creation time of the CA Cert.`,
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Expiration time of the CA Cert.`,
						},
						"sha1_fingerprint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `SHA Fingerprint of the CA Cert.`,
						},
					},
				},
			},
			"service_account_email_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service account email address assigned to the instance.`,
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},
			"restore_backup_context": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_run_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The ID of the backup run to restore from.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the instance that the backup was taken from.`,
						},
						"project": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The full project ID of the source instance.`,
						},
					},
				},
			},
			"clone": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     false,
				AtLeastOneOf: []string{"settings", "clone"},
				Description:  `Configuration for creating a new instance as a clone of another instance.`,
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_instance_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the instance from which the point in time should be restored.`,
						},
						"point_in_time": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: timestampDiffSuppress(time.RFC3339Nano),
							Description:      `The timestamp of the point in time that should be restored.`,
						},
						"allocated_ip_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the allocated ip range for the private ip CloudSQL instance. For example: "google-managed-services-default". If set, the cloned instance ip will be created in the allocated range. The range name must comply with [RFC 1035](https://tools.ietf.org/html/rfc1035). Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])?.`,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

// Makes private_network ForceNew if it is changing from set to nil. The API returns an error
// if this change is attempted in-place.
func privateNetworkCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	old, new := d.GetChange("settings.0.ip_configuration.0.private_network")

	if old != "" && new == "" {
		if err := d.ForceNew("settings.0.ip_configuration.0.private_network"); err != nil {
			return err
		}
	}

	return nil
}

// Point in time recovery for MySQL database instances needs binary_log_enabled set to true and
// not point_in_time_recovery_enabled, which is confusing to users. This checks for
// point_in_time_recovery_enabled being set to a non-PostgreSQL database instance and suggests
// binary_log_enabled.
func pitrPostgresOnlyCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
	pitr := diff.Get("settings.0.backup_configuration.0.point_in_time_recovery_enabled").(bool)
	dbVersion := diff.Get("database_version").(string)
	if pitr && !strings.Contains(dbVersion, "POSTGRES") {
		return fmt.Errorf("point_in_time_recovery_enabled is only available for Postgres. You may want to consider using binary_log_enabled instead.")
	}
	return nil
}

func resourceSqlDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

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

	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	// SQL Instances that fail to create are expensive- see https://github.com/hashicorp/terraform-provider-google/issues/7154
	// We can fail fast to stop instance names from getting reserved.
	network := d.Get("settings.0.ip_configuration.0.private_network").(string)
	if network != "" {
		err = sqlDatabaseInstanceServiceNetworkPrecheck(d, config, userAgent, network)
		if err != nil {
			return err
		}
	}

	instance := &sqladmin.DatabaseInstance{
		Name:                 name,
		Region:               region,
		DatabaseVersion:      d.Get("database_version").(string),
		MasterInstanceName:   d.Get("master_instance_name").(string),
		ReplicaConfiguration: expandReplicaConfiguration(d.Get("replica_configuration").([]interface{})),
	}

	cloneContext, cloneSource := expandCloneContext(d.Get("clone").([]interface{}))

	s, ok := d.GetOk("settings")
	desiredSettings := expandSqlDatabaseInstanceSettings(s.([]interface{}))
	if ok {
		instance.Settings = desiredSettings
	}

	if _, ok := d.GetOk("maintenance_version"); ok {
		instance.MaintenanceVersion = d.Get("maintenance_version").(string)
	}

	instance.RootPassword = d.Get("root_password").(string)

	// Modifying a replica during Create can cause problems if the master is
	// modified at the same time. Lock the master until we're done in order
	// to prevent that.
	if !sqlDatabaseIsMaster(d) {
		mutexKV.Lock(instanceMutexKey(project, instance.MasterInstanceName))
		defer mutexKV.Unlock(instanceMutexKey(project, instance.MasterInstanceName))
	}

	if k, ok := d.GetOk("encryption_key_name"); ok {
		instance.DiskEncryptionConfiguration = &sqladmin.DiskEncryptionConfiguration{
			KmsKeyName: k.(string),
		}
	}

	var patchData *sqladmin.DatabaseInstance

	// BinaryLogging can be enabled on replica instances but only after creation.
	if instance.MasterInstanceName != "" && instance.Settings != nil && instance.Settings.BackupConfiguration != nil && instance.Settings.BackupConfiguration.BinaryLogEnabled {
		settingsCopy := expandSqlDatabaseInstanceSettings(s.([]interface{}))
		bc := settingsCopy.BackupConfiguration
		patchData = &sqladmin.DatabaseInstance{Settings: &sqladmin.Settings{BackupConfiguration: bc}}

		instance.Settings.BackupConfiguration.BinaryLogEnabled = false
	}

	var op *sqladmin.Operation
	err = retryTimeDuration(func() (operr error) {
		if cloneContext != nil {
			cloneContext.DestinationInstanceName = name
			clodeReq := sqladmin.InstancesCloneRequest{CloneContext: cloneContext}
			op, operr = config.NewSqlAdminClient(userAgent).Instances.Clone(project, cloneSource, &clodeReq).Do()
		} else {
			op, operr = config.NewSqlAdminClient(userAgent).Instances.Insert(project, instance).Do()
		}
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

	err = sqlAdminOperationWaitTime(config, op, project, "Create Instance", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		d.SetId("")
		return err
	}

	// If a default root user was created with a wildcard ('%') hostname, delete it. Note it
	// appears to only be created for certain types of databases, like MySQL.
	// Users in a replica instance are inherited from the master instance and should be left alone.
	// This deletion is done immediately after the instance is created, in order to minimize the
	// risk of it being left on the instance, which would present a security concern.
	if sqlDatabaseIsMaster(d) {
		var users *sqladmin.UsersListResponse
		err = retryTimeDuration(func() error {
			users, err = config.NewSqlAdminClient(userAgent).Users.List(project, instance.Name).Do()
			return err
		}, d.Timeout(schema.TimeoutRead), isSqlOperationInProgressError)
		if err != nil {
			return fmt.Errorf("Error, attempting to list users associated with instance %s: %s", instance.Name, err)
		}
		for _, u := range users.Items {
			if u.Name == "root" && u.Host == "%" {
				err = retry(func() error {
					op, err = config.NewSqlAdminClient(userAgent).Users.Delete(project, instance.Name).Host(u.Host).Name(u.Name).Do()
					if err == nil {
						err = sqlAdminOperationWaitTime(config, op, project, "Delete default root User", userAgent, d.Timeout(schema.TimeoutCreate))
					}
					return err
				})
				if err != nil {
					return fmt.Errorf("Error, failed to delete default 'root'@'*' user, but the database was created successfully: %s", err)
				}
			}
		}
	}

	// patch any fields that need to be sent postcreation
	if patchData != nil {
		err = retryTimeDuration(func() (rerr error) {
			op, rerr = config.NewSqlAdminClient(userAgent).Instances.Patch(project, instance.Name, patchData).Do()
			return rerr
		}, d.Timeout(schema.TimeoutUpdate), isSqlOperationInProgressError)
		if err != nil {
			return fmt.Errorf("Error, failed to update instance settings for %s: %s", instance.Name, err)
		}
		err = sqlAdminOperationWaitTime(config, op, project, "Patch Instance", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	err = resourceSqlDatabaseInstanceRead(d, meta)
	if err != nil {
		return err
	}

	// Refresh settings from read as they may have defaulted from the API
	s = d.Get("settings")
	// If we've created an instance as a clone, we need to update it to set any user defined settings
	if len(s.([]interface{})) != 0 && cloneContext != nil && desiredSettings != nil {
		instanceUpdate := &sqladmin.DatabaseInstance{
			Settings: desiredSettings,
		}
		_settings := s.([]interface{})[0].(map[string]interface{})
		instanceUpdate.Settings.SettingsVersion = int64(_settings["version"].(int))
		var op *sqladmin.Operation
		err = retryTimeDuration(func() (rerr error) {
			op, rerr = config.NewSqlAdminClient(userAgent).Instances.Update(project, name, instanceUpdate).Do()
			return rerr
		}, d.Timeout(schema.TimeoutUpdate), isSqlOperationInProgressError)
		if err != nil {
			return fmt.Errorf("Error, failed to update instance settings for %s: %s", instance.Name, err)
		}

		err = sqlAdminOperationWaitTime(config, op, project, "Update Instance", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		// Refresh the state of the instance after updating the settings
		err = resourceSqlDatabaseInstanceRead(d, meta)
		if err != nil {
			return err
		}
	}

	// Perform a backup restore if the backup context exists
	if r, ok := d.GetOk("restore_backup_context"); ok {
		err = sqlDatabaseInstanceRestoreFromBackup(d, config, userAgent, project, name, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func expandSqlDatabaseInstanceSettings(configured []interface{}) *sqladmin.Settings {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_settings := configured[0].(map[string]interface{})
	settings := &sqladmin.Settings{
		// Version is unset in Create but is set during update
		SettingsVersion:           int64(_settings["version"].(int)),
		Tier:                      _settings["tier"].(string),
		ForceSendFields:           []string{"StorageAutoResize"},
		ActivationPolicy:          _settings["activation_policy"].(string),
		ActiveDirectoryConfig:     expandActiveDirectoryConfig(_settings["active_directory_config"].([]interface{})),
		DenyMaintenancePeriods:    expandDenyMaintenancePeriod(_settings["deny_maintenance_period"].([]interface{})),
		SqlServerAuditConfig:      expandSqlServerAuditConfig(_settings["sql_server_audit_config"].([]interface{})),
		TimeZone:                  _settings["time_zone"].(string),
		AvailabilityType:          _settings["availability_type"].(string),
		ConnectorEnforcement:      _settings["connector_enforcement"].(string),
		Collation:                 _settings["collation"].(string),
		DataDiskSizeGb:            int64(_settings["disk_size"].(int)),
		DataDiskType:              _settings["disk_type"].(string),
		PricingPlan:               _settings["pricing_plan"].(string),
		DeletionProtectionEnabled: _settings["deletion_protection_enabled"].(bool),
		UserLabels:                convertStringMap(_settings["user_labels"].(map[string]interface{})),
		BackupConfiguration:       expandBackupConfiguration(_settings["backup_configuration"].([]interface{})),
		DatabaseFlags:             expandDatabaseFlags(_settings["database_flags"].([]interface{})),
		IpConfiguration:           expandIpConfiguration(_settings["ip_configuration"].([]interface{})),
		LocationPreference:        expandLocationPreference(_settings["location_preference"].([]interface{})),
		MaintenanceWindow:         expandMaintenanceWindow(_settings["maintenance_window"].([]interface{})),
		InsightsConfig:            expandInsightsConfig(_settings["insights_config"].([]interface{})),
		PasswordValidationPolicy:  expandPasswordValidationPolicy(_settings["password_validation_policy"].([]interface{})),
	}

	resize := _settings["disk_autoresize"].(bool)
	settings.StorageAutoResize = &resize
	settings.StorageAutoResizeLimit = int64(_settings["disk_autoresize_limit"].(int))

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

func expandCloneContext(configured []interface{}) (*sqladmin.CloneContext, string) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, ""
	}

	_cloneConfiguration := configured[0].(map[string]interface{})

	return &sqladmin.CloneContext{
		PointInTime:      _cloneConfiguration["point_in_time"].(string),
		AllocatedIpRange: _cloneConfiguration["allocated_ip_range"].(string),
	}, _cloneConfiguration["source_instance_name"].(string)
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
		SecondaryZone:        _locationPreference["secondary_zone"].(string),
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
		AllocatedIpRange:   _ipConfiguration["allocated_ip_range"].(string),
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

func expandDatabaseFlags(configured []interface{}) []*sqladmin.DatabaseFlags {
	databaseFlags := make([]*sqladmin.DatabaseFlags, 0, len(configured))
	for _, _flag := range configured {
		if _flag == nil {
			continue
		}
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
		BinaryLogEnabled:            _backupConfiguration["binary_log_enabled"].(bool),
		BackupRetentionSettings:     expandBackupRetentionSettings(_backupConfiguration["backup_retention_settings"]),
		Enabled:                     _backupConfiguration["enabled"].(bool),
		StartTime:                   _backupConfiguration["start_time"].(string),
		Location:                    _backupConfiguration["location"].(string),
		TransactionLogRetentionDays: int64(_backupConfiguration["transaction_log_retention_days"].(int)),
		PointInTimeRecoveryEnabled:  _backupConfiguration["point_in_time_recovery_enabled"].(bool),
		ForceSendFields:             []string{"BinaryLogEnabled", "Enabled", "PointInTimeRecoveryEnabled"},
	}
}

func expandBackupRetentionSettings(configured interface{}) *sqladmin.BackupRetentionSettings {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &sqladmin.BackupRetentionSettings{
		RetainedBackups: int64(config["retained_backups"].(int)),
		RetentionUnit:   config["retention_unit"].(string),
	}
}

func expandActiveDirectoryConfig(configured interface{}) *sqladmin.SqlActiveDirectoryConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &sqladmin.SqlActiveDirectoryConfig{
		Domain: config["domain"].(string),
	}
}

func expandDenyMaintenancePeriod(configured []interface{}) []*sqladmin.DenyMaintenancePeriod {
	denyMaintenancePeriod := make([]*sqladmin.DenyMaintenancePeriod, 0, len(configured))

	for _, _flag := range configured {
		if _flag == nil {
			continue
		}
		_entry := _flag.(map[string]interface{})

		denyMaintenancePeriod = append(denyMaintenancePeriod, &sqladmin.DenyMaintenancePeriod{
			EndDate:   _entry["end_date"].(string),
			StartDate: _entry["start_date"].(string),
			Time:      _entry["time"].(string),
		})
	}
	return denyMaintenancePeriod

}

func expandSqlServerAuditConfig(configured interface{}) *sqladmin.SqlServerAuditConfig {
	l := configured.([]interface{})
	if len(l) == 0 {
		return nil
	}

	config := l[0].(map[string]interface{})
	return &sqladmin.SqlServerAuditConfig{
		Bucket:            config["bucket"].(string),
		RetentionInterval: config["retention_interval"].(string),
		UploadInterval:    config["upload_interval"].(string),
	}
}

func expandInsightsConfig(configured []interface{}) *sqladmin.InsightsConfig {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_insightsConfig := configured[0].(map[string]interface{})
	return &sqladmin.InsightsConfig{
		QueryInsightsEnabled:  _insightsConfig["query_insights_enabled"].(bool),
		QueryStringLength:     int64(_insightsConfig["query_string_length"].(int)),
		RecordApplicationTags: _insightsConfig["record_application_tags"].(bool),
		RecordClientAddress:   _insightsConfig["record_client_address"].(bool),
		QueryPlansPerMinute:   int64(_insightsConfig["query_plans_per_minute"].(int)),
	}
}

func expandPasswordValidationPolicy(configured []interface{}) *sqladmin.PasswordValidationPolicy {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_passwordValidationPolicy := configured[0].(map[string]interface{})
	return &sqladmin.PasswordValidationPolicy{
		MinLength:                 int64(_passwordValidationPolicy["min_length"].(int)),
		Complexity:                _passwordValidationPolicy["complexity"].(string),
		ReuseInterval:             int64(_passwordValidationPolicy["reuse_interval"].(int)),
		DisallowUsernameSubstring: _passwordValidationPolicy["disallow_username_substring"].(bool),
		PasswordChangeInterval:    _passwordValidationPolicy["password_change_interval"].(string),
		EnablePasswordPolicy:      _passwordValidationPolicy["enable_password_policy"].(bool),
	}
}

func resourceSqlDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	var instance *sqladmin.DatabaseInstance
	err = retryTimeDuration(func() (rerr error) {
		instance, rerr = config.NewSqlAdminClient(userAgent).Instances.Get(project, d.Get("name").(string)).Do()
		return rerr
	}, d.Timeout(schema.TimeoutRead), isSqlOperationInProgressError)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL Database Instance %q", d.Get("name").(string)))
	}

	if err := d.Set("name", instance.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("region", instance.Region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("database_version", instance.DatabaseVersion); err != nil {
		return fmt.Errorf("Error setting database_version: %s", err)
	}
	if err := d.Set("connection_name", instance.ConnectionName); err != nil {
		return fmt.Errorf("Error setting connection_name: %s", err)
	}
	if err := d.Set("maintenance_version", instance.MaintenanceVersion); err != nil {
		return fmt.Errorf("Error setting maintenance_version: %s", err)
	}
	if err := d.Set("available_maintenance_versions", instance.AvailableMaintenanceVersions); err != nil {
		return fmt.Errorf("Error setting available_maintenance_version: %s", err)
	}
	if err := d.Set("service_account_email_address", instance.ServiceAccountEmailAddress); err != nil {
		return fmt.Errorf("Error setting service_account_email_address: %s", err)
	}

	if err := d.Set("settings", flattenSettings(instance.Settings)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance Settings")
	}

	if instance.DiskEncryptionConfiguration != nil {
		if err := d.Set("encryption_key_name", instance.DiskEncryptionConfiguration.KmsKeyName); err != nil {
			return fmt.Errorf("Error setting encryption_key_name: %s", err)
		}
	}

	if err := d.Set("replica_configuration", flattenReplicaConfiguration(instance.ReplicaConfiguration, d)); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance Replica Configuration")
	}
	ipAddresses := flattenIpAddresses(instance.IpAddresses)
	if err := d.Set("ip_address", ipAddresses); err != nil {
		log.Printf("[WARN] Failed to set SQL Database Instance IP Addresses")
	}

	if len(ipAddresses) > 0 {
		if err := d.Set("first_ip_address", ipAddresses[0]["ip_address"]); err != nil {
			return fmt.Errorf("Error setting first_ip_address: %s", err)
		}
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

	if err := d.Set("public_ip_address", publicIpAddress); err != nil {
		return fmt.Errorf("Error setting public_ip_address: %s", err)
	}
	if err := d.Set("private_ip_address", privateIpAddress); err != nil {
		return fmt.Errorf("Error setting private_ip_address: %s", err)
	}

	if err := d.Set("server_ca_cert", flattenServerCaCerts([]*sqladmin.SslCert{instance.ServerCaCert})); err != nil {
		log.Printf("[WARN] Failed to set SQL Database CA Certificate")
	}

	if err := d.Set("master_instance_name", strings.TrimPrefix(instance.MasterInstanceName, project+":")); err != nil {
		return fmt.Errorf("Error setting master_instance_name: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("self_link", instance.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	d.SetId(instance.Name)

	return nil
}

func resourceSqlDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	var maintenance_version string
	if v, ok := d.GetOk("maintenance_version"); ok {
		maintenance_version = v.(string)
	}

	desiredSetting := d.Get("settings")
	var op *sqladmin.Operation
	var instance *sqladmin.DatabaseInstance

	// Check if the database version is being updated, because patching database version is an atomic operation and can not be
	// performed with other fields, we first patch database version before updating the rest of the fields.
	if v, ok := d.GetOk("database_version"); ok {
		instance = &sqladmin.DatabaseInstance{DatabaseVersion: v.(string)}
		err = retryTimeDuration(func() (rerr error) {
			op, rerr = config.NewSqlAdminClient(userAgent).Instances.Patch(project, d.Get("name").(string), instance).Do()
			return rerr
		}, d.Timeout(schema.TimeoutUpdate), isSqlOperationInProgressError)
		if err != nil {
			return fmt.Errorf("Error, failed to patch instance settings for %s: %s", instance.Name, err)
		}
		err = sqlAdminOperationWaitTime(config, op, project, "Patch Instance", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		err = resourceSqlDatabaseInstanceRead(d, meta)
		if err != nil {
			return err
		}
	}

	// Check if the maintenance version is being updated, because patching maintenance version is an atomic operation and can not be
	// performed with other fields, we first patch maintenance version before updating the rest of the fields.
	if d.HasChange("maintenance_version") {
		instance = &sqladmin.DatabaseInstance{MaintenanceVersion: maintenance_version}
		err = retryTimeDuration(func() (rerr error) {
			op, rerr = config.NewSqlAdminClient(userAgent).Instances.Patch(project, d.Get("name").(string), instance).Do()
			return rerr
		}, d.Timeout(schema.TimeoutUpdate), isSqlOperationInProgressError)
		if err != nil {
			return fmt.Errorf("Error, failed to patch instance settings for %s: %s", instance.Name, err)
		}
		err = sqlAdminOperationWaitTime(config, op, project, "Patch Instance", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		err = resourceSqlDatabaseInstanceRead(d, meta)
		if err != nil {
			return err
		}
	}

	s := d.Get("settings")
	instance = &sqladmin.DatabaseInstance{
		Settings: expandSqlDatabaseInstanceSettings(desiredSetting.([]interface{})),
	}
	_settings := s.([]interface{})[0].(map[string]interface{})
	// Instance.Patch operation on completion updates the settings proto version by +8. As terraform does not know this it tries
	// to make an update call with the proto version before patch and fails. To resolve this issue we update the setting version
	// before making the update call.
	instance.Settings.SettingsVersion = int64(_settings["version"].(int))
	// Collation cannot be included in the update request
	instance.Settings.Collation = ""

	// Lock on the master_instance_name just in case updating any replica
	// settings causes operations on the master.
	if v, ok := d.GetOk("master_instance_name"); ok {
		mutexKV.Lock(instanceMutexKey(project, v.(string)))
		defer mutexKV.Unlock(instanceMutexKey(project, v.(string)))
	}

	err = retryTimeDuration(func() (rerr error) {
		op, rerr = config.NewSqlAdminClient(userAgent).Instances.Update(project, d.Get("name").(string), instance).Do()
		return rerr
	}, d.Timeout(schema.TimeoutUpdate), isSqlOperationInProgressError)
	if err != nil {
		return fmt.Errorf("Error, failed to update instance settings for %s: %s", instance.Name, err)
	}

	err = sqlAdminOperationWaitTime(config, op, project, "Update Instance", userAgent, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	// Perform a backup restore if the backup context exists and has changed
	if r, ok := d.GetOk("restore_backup_context"); ok {
		if d.HasChange("restore_backup_context") {
			err = sqlDatabaseInstanceRestoreFromBackup(d, config, userAgent, project, d.Get("name").(string), r)
			if err != nil {
				return err
			}
		}
	}

	return resourceSqlDatabaseInstanceRead(d, meta)
}

func maintenanceVersionDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	// Ignore the database version part and only compare the last part of the maintenance version which represents the release date of the version.
	if len(old) > 14 && len(new) > 14 && old[len(old)-14:] >= new[len(new)-14:] {
		log.Printf("[DEBUG] Maintenance version in configuration [%s] is older than current maintenance version [%s] on instance. Suppressing diff", new, old)
		return true
	} else {
		return false
	}
}

func resourceSqlDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Check if deletion protection is enabled.

	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("Error, failed to delete instance because deletion_protection is set to true. Set it to false to proceed with instance deletion")
	}

	// Lock on the master_instance_name just in case deleting a replica causes
	// operations on the master.
	if v, ok := d.GetOk("master_instance_name"); ok {
		mutexKV.Lock(instanceMutexKey(project, v.(string)))
		defer mutexKV.Unlock(instanceMutexKey(project, v.(string)))
	}

	var op *sqladmin.Operation
	err = retryTimeDuration(func() (rerr error) {
		op, rerr = config.NewSqlAdminClient(userAgent).Instances.Delete(project, d.Get("name").(string)).Do()
		if rerr != nil {
			return rerr
		}
		err = sqlAdminOperationWaitTime(config, op, project, "Delete Instance", userAgent, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return err
		}
		return nil
	}, d.Timeout(schema.TimeoutDelete), isSqlOperationInProgressError, isSqlInternalError)
	if err != nil {
		return fmt.Errorf("Error, failed to delete instance %s: %s", d.Get("name").(string), err)
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

	if err := d.Set("deletion_protection", true); err != nil {
		return nil, fmt.Errorf("Error setting deletion_protection: %s", err)
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
		"availability_type":           settings.AvailabilityType,
		"collation":                   settings.Collation,
		"connector_enforcement":       settings.ConnectorEnforcement,
		"disk_type":                   settings.DataDiskType,
		"disk_size":                   settings.DataDiskSizeGb,
		"pricing_plan":                settings.PricingPlan,
		"user_labels":                 settings.UserLabels,
		"password_validation_policy":  settings.PasswordValidationPolicy,
		"time_zone":                   settings.TimeZone,
		"deletion_protection_enabled": settings.DeletionProtectionEnabled,
	}

	if settings.ActiveDirectoryConfig != nil {
		data["active_directory_config"] = flattenActiveDirectoryConfig(settings.ActiveDirectoryConfig)
	}

	if settings.DenyMaintenancePeriods != nil {
		data["deny_maintenance_period"] = flattenDenyMaintenancePeriod(settings.DenyMaintenancePeriods)
	}

	if settings.SqlServerAuditConfig != nil {
		data["sql_server_audit_config"] = flattenSqlServerAuditConfig(settings.SqlServerAuditConfig)
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

	if settings.InsightsConfig != nil {
		data["insights_config"] = flattenInsightsConfig(settings.InsightsConfig)
	}

	data["disk_autoresize"] = settings.StorageAutoResize
	data["disk_autoresize_limit"] = settings.StorageAutoResizeLimit

	if settings.UserLabels != nil {
		data["user_labels"] = settings.UserLabels
	}

	if settings.PasswordValidationPolicy != nil {
		data["password_validation_policy"] = flattenPasswordValidationPolicy(settings.PasswordValidationPolicy)
	}

	return []map[string]interface{}{data}
}

func flattenBackupConfiguration(backupConfiguration *sqladmin.BackupConfiguration) []map[string]interface{} {
	data := map[string]interface{}{
		"binary_log_enabled":             backupConfiguration.BinaryLogEnabled,
		"enabled":                        backupConfiguration.Enabled,
		"start_time":                     backupConfiguration.StartTime,
		"location":                       backupConfiguration.Location,
		"point_in_time_recovery_enabled": backupConfiguration.PointInTimeRecoveryEnabled,
		"backup_retention_settings":      flattenBackupRetentionSettings(backupConfiguration.BackupRetentionSettings),
		"transaction_log_retention_days": backupConfiguration.TransactionLogRetentionDays,
	}

	return []map[string]interface{}{data}
}

func flattenBackupRetentionSettings(b *sqladmin.BackupRetentionSettings) []map[string]interface{} {
	if b == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"retained_backups": b.RetainedBackups,
			"retention_unit":   b.RetentionUnit,
		},
	}
}

func flattenActiveDirectoryConfig(sqlActiveDirectoryConfig *sqladmin.SqlActiveDirectoryConfig) []map[string]interface{} {
	if sqlActiveDirectoryConfig == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"domain": sqlActiveDirectoryConfig.Domain,
		},
	}
}

func flattenDenyMaintenancePeriod(denyMaintenancePeriod []*sqladmin.DenyMaintenancePeriod) []map[string]interface{} {
	flags := make([]map[string]interface{}, 0, len(denyMaintenancePeriod))

	for _, flag := range denyMaintenancePeriod {
		data := map[string]interface{}{
			"end_date":   flag.EndDate,
			"start_date": flag.StartDate,
			"time":       flag.Time,
		}

		flags = append(flags, data)
	}

	return flags
}

func flattenSqlServerAuditConfig(sqlServerAuditConfig *sqladmin.SqlServerAuditConfig) []map[string]interface{} {
	if sqlServerAuditConfig == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"bucket":             sqlServerAuditConfig.Bucket,
			"retention_interval": sqlServerAuditConfig.RetentionInterval,
			"upload_interval":    sqlServerAuditConfig.UploadInterval,
		},
	}
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
		"ipv4_enabled":       ipConfiguration.Ipv4Enabled,
		"private_network":    ipConfiguration.PrivateNetwork,
		"allocated_ip_range": ipConfiguration.AllocatedIpRange,
		"require_ssl":        ipConfiguration.RequireSsl,
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
		"secondary_zone":         locationPreference.SecondaryZone,
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

func flattenServerCaCerts(caCerts []*sqladmin.SslCert) []map[string]interface{} {
	var certs []map[string]interface{}

	for _, caCert := range caCerts {
		if caCert != nil {
			data := map[string]interface{}{
				"cert":             caCert.Cert,
				"common_name":      caCert.CommonName,
				"create_time":      caCert.CreateTime,
				"expiration_time":  caCert.ExpirationTime,
				"sha1_fingerprint": caCert.Sha1Fingerprint,
			}

			certs = append(certs, data)
		}
	}

	return certs
}

func flattenInsightsConfig(insightsConfig *sqladmin.InsightsConfig) interface{} {
	data := map[string]interface{}{
		"query_insights_enabled":  insightsConfig.QueryInsightsEnabled,
		"query_string_length":     insightsConfig.QueryStringLength,
		"record_application_tags": insightsConfig.RecordApplicationTags,
		"record_client_address":   insightsConfig.RecordClientAddress,
		"query_plans_per_minute":  insightsConfig.QueryPlansPerMinute,
	}

	return []map[string]interface{}{data}
}

func flattenPasswordValidationPolicy(passwordValidationPolicy *sqladmin.PasswordValidationPolicy) interface{} {
	data := map[string]interface{}{
		"min_length":                  passwordValidationPolicy.MinLength,
		"complexity":                  passwordValidationPolicy.Complexity,
		"reuse_interval":              passwordValidationPolicy.ReuseInterval,
		"disallow_username_substring": passwordValidationPolicy.DisallowUsernameSubstring,
		"password_change_interval":    passwordValidationPolicy.PasswordChangeInterval,
		"enable_password_policy":      passwordValidationPolicy.EnablePasswordPolicy,
	}
	return []map[string]interface{}{data}
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

func sqlDatabaseInstanceServiceNetworkPrecheck(d *schema.ResourceData, config *Config, userAgent, network string) error {
	log.Printf("[DEBUG] checking network %q for at least one service networking connection", network)
	// This call requires projects.get permissions, which may not have been granted to the Terraform actor,
	// particularly in shared VPC setups. Most will! But it's not strictly required.
	serviceNetworkingNetworkName, err := retrieveServiceNetworkingNetworkName(d, config, network, userAgent)
	if err != nil {
		var gerr *googleapi.Error
		if errors.As(err, &gerr) {
			log.Printf("[DEBUG] retrieved googleapi error while creating sn name for %q. precheck skipped. code %v and message: %s", network, gerr.Code, gerr.Body)
			return nil
		}

		return err
	}

	response, err := config.NewServiceNetworkingClient(userAgent).Services.Connections.List("services/servicenetworking.googleapis.com").Network(serviceNetworkingNetworkName).Do()
	if err != nil {
		// It is possible that the actor creating the SQL Instance might not have permissions to call servicenetworking.services.connections.list
		log.Printf("[WARNING] Failed to list Service Networking of the project. Skipped Service Networking precheck.")
		return nil
	}

	if len(response.Connections) < 1 {
		return fmt.Errorf("Error, failed to create instance because the network doesn't have at least 1 private services connection. Please see https://cloud.google.com/sql/docs/mysql/private-ip#network_requirements for how to create this connection.")
	}

	return nil
}

func expandRestoreBackupContext(configured []interface{}) *sqladmin.RestoreBackupContext {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	_rc := configured[0].(map[string]interface{})
	return &sqladmin.RestoreBackupContext{
		BackupRunId: int64(_rc["backup_run_id"].(int)),
		InstanceId:  _rc["instance_id"].(string),
		Project:     _rc["project"].(string),
	}
}

func sqlDatabaseInstanceRestoreFromBackup(d *schema.ResourceData, config *Config, userAgent, project, instanceId string, r interface{}) error {
	log.Printf("[DEBUG] Initiating SQL database instance backup restore")
	restoreContext := r.([]interface{})

	backupRequest := &sqladmin.InstancesRestoreBackupRequest{
		RestoreBackupContext: expandRestoreBackupContext(restoreContext),
	}

	var op *sqladmin.Operation
	err := retryTimeDuration(func() (operr error) {
		op, operr = config.NewSqlAdminClient(userAgent).Instances.RestoreBackup(project, instanceId, backupRequest).Do()
		return operr
	}, d.Timeout(schema.TimeoutUpdate), isSqlOperationInProgressError)
	if err != nil {
		return fmt.Errorf("Error, failed to restore instance from backup %s: %s", instanceId, err)
	}

	err = sqlAdminOperationWaitTime(config, op, project, "Restore Backup", userAgent, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func caseDiffDashSuppress(_, old, new string, _ *schema.ResourceData) bool {
	postReplaceNew := strings.Replace(new, "-", "_", -1)
	return strings.ToUpper(postReplaceNew) == strings.ToUpper(old)
}
