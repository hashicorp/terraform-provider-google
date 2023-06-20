// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func DataSourceSqlDatabaseInstances() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceSqlDatabaseInstancesRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project that contains the instances.`,
			},
			"database_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `To filter out the database instances which are of the specified database version.`,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `To filter out the database instances which are located in this specified region.`,
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `To filter out the database instances which are located in this specified zone.`,
			},
			"tier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `To filter out the database instances based on the machine type.`,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `To filter out the database instances based on the current state of the database instance, valid values include : "SQL_INSTANCE_STATE_UNSPECIFIED", "RUNNABLE", "SUSPENDED", "PENDING_DELETE", "PENDING_CREATE", "MAINTENANCE" and "FAILED".`,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceSqlDatabaseInstance().Schema),
				},
			},
		},
	}
}

func dataSourceSqlDatabaseInstancesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	filter := ""

	if v, ok := d.GetOk("database_version"); ok {
		filter += fmt.Sprintf("databaseVersion:%s", v.(string))
	}
	if v, ok := d.GetOk("region"); ok {
		if filter != "" {
			filter += " AND "
		}
		filter += fmt.Sprintf("region:%s", v.(string))
	}
	if v, ok := d.GetOk("zone"); ok {
		if filter != "" {
			filter += " AND "
		}
		filter += fmt.Sprintf("gceZone:%s", v.(string))
	}
	if v, ok := d.GetOk("tier"); ok {
		if filter != "" {
			filter += " AND "
		}
		filter += fmt.Sprintf("settings.tier:%s", v.(string))
	}
	if v, ok := d.GetOk("state"); ok {
		if filter != "" {
			filter += " AND "
		}
		filter += fmt.Sprintf("state:%s", v.(string))
	}
	pageToken := ""
	databaseInstances := make([]map[string]interface{}, 0)
	for {
		var instances *sqladmin.InstancesListResponse
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (rerr error) {
				instances, rerr = config.NewSqlAdminClient(userAgent).Instances.List(project).Filter(filter).PageToken(pageToken).Do()
				return rerr
			},
			Timeout:              d.Timeout(schema.TimeoutRead),
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsSqlOperationInProgressError},
		})
		if err != nil {
			return err
		}

		pageInstances := flattenDatasourceGoogleDatabaseInstancesList(instances.Items, project)
		databaseInstances = append(databaseInstances, pageInstances...)

		pageToken = instances.NextPageToken
		if pageToken == "" {
			break
		}
	}

	if err := d.Set("instances", databaseInstances); err != nil {
		return fmt.Errorf("Error retrieving instances: %s", err)
	}

	d.SetId(fmt.Sprintf("database_instances_ds/%s/%s/%s/%s/%s/%s", project, d.Get("database_version").(string), d.Get("region").(string), d.Get("zone").(string), d.Get("tier").(string), d.Get("state").(string)))

	return nil
}

func flattenDatasourceGoogleDatabaseInstancesList(fetchedInstances []*sqladmin.DatabaseInstance, project string) []map[string]interface{} {
	if fetchedInstances == nil {
		return make([]map[string]interface{}, 0)
	}

	instances := make([]map[string]interface{}, 0, len(fetchedInstances))
	for _, rawInstance := range fetchedInstances {
		instance := make(map[string]interface{})
		instance["name"] = rawInstance.Name
		instance["region"] = rawInstance.Region
		instance["database_version"] = rawInstance.DatabaseVersion
		instance["connection_name"] = rawInstance.ConnectionName
		instance["maintenance_version"] = rawInstance.MaintenanceVersion
		instance["available_maintenance_versions"] = rawInstance.AvailableMaintenanceVersions
		instance["instance_type"] = rawInstance.InstanceType
		instance["service_account_email_address"] = rawInstance.ServiceAccountEmailAddress
		instance["settings"] = flattenSettings(rawInstance.Settings)

		if rawInstance.DiskEncryptionConfiguration != nil {
			instance["encryption_key_name"] = rawInstance.DiskEncryptionConfiguration.KmsKeyName
		}

		instance["replica_configuration"] = flattenReplicaConfigurationforDataSource(rawInstance.ReplicaConfiguration)

		ipAddresses := flattenIpAddresses(rawInstance.IpAddresses)
		instance["ip_address"] = ipAddresses

		if len(ipAddresses) > 0 {
			instance["first_ip_address"] = ipAddresses[0]["ip_address"]
		}

		publicIpAddress := ""
		privateIpAddress := ""
		for _, ip := range rawInstance.IpAddresses {
			if publicIpAddress == "" && ip.Type == "PRIMARY" {
				publicIpAddress = ip.IpAddress
			}

			if privateIpAddress == "" && ip.Type == "PRIVATE" {
				privateIpAddress = ip.IpAddress
			}
		}
		instance["public_ip_address"] = publicIpAddress
		instance["private_ip_address"] = privateIpAddress
		instance["server_ca_cert"] = flattenServerCaCerts([]*sqladmin.SslCert{rawInstance.ServerCaCert})
		instance["master_instance_name"] = strings.TrimPrefix(rawInstance.MasterInstanceName, project+":")
		instance["project"] = project
		instance["self_link"] = rawInstance.SelfLink

		instances = append(instances, instance)
	}

	return instances
}
func flattenReplicaConfigurationforDataSource(replicaConfiguration *sqladmin.ReplicaConfiguration) []map[string]interface{} {
	rc := []map[string]interface{}{}

	if replicaConfiguration != nil {
		data := map[string]interface{}{
			"failover_target": replicaConfiguration.FailoverTarget,
			// Don't attempt to assign anything from replicaConfiguration.MysqlReplicaConfiguration,
			// since those fields are set on create and then not stored. Hence, those fields are not shown.
		}
		rc = append(rc, data)
	}

	return rc
}
