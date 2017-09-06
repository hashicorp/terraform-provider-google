package google

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
				ValidateFunc: validateCredentials,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_PROJECT",
					"GCLOUD_PROJECT",
					"CLOUDSDK_CORE_PROJECT",
				}, nil),
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_REGION",
					"GCLOUD_REGION",
					"CLOUDSDK_COMPUTE_REGION",
				}, nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"google_dns_managed_zone":          dataSourceDnsManagedZone(),
			"google_client_config":             dataSourceGoogleClientConfig(),
			"google_compute_network":           dataSourceGoogleComputeNetwork(),
			"google_compute_subnetwork":        dataSourceGoogleComputeSubnetwork(),
			"google_compute_zones":             dataSourceGoogleComputeZones(),
			"google_compute_instance_group":    dataSourceGoogleComputeInstanceGroup(),
			"google_container_engine_versions": dataSourceGoogleContainerEngineVersions(),
			"google_iam_policy":                dataSourceGoogleIamPolicy(),
			"google_storage_object_signed_url": dataSourceGoogleSignedUrl(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"google_bigquery_dataset":               resourceBigQueryDataset(),
			"google_bigquery_table":                 resourceBigQueryTable(),
			"google_bigtable_instance":              resourceBigtableInstance(),
			"google_bigtable_table":                 resourceBigtableTable(),
			"google_compute_autoscaler":             resourceComputeAutoscaler(),
			"google_compute_address":                resourceComputeAddress(),
			"google_compute_backend_bucket":         resourceComputeBackendBucket(),
			"google_compute_backend_service":        resourceComputeBackendService(),
			"google_compute_disk":                   resourceComputeDisk(),
			"google_compute_snapshot":               resourceComputeSnapshot(),
			"google_compute_firewall":               resourceComputeFirewall(),
			"google_compute_forwarding_rule":        resourceComputeForwardingRule(),
			"google_compute_global_address":         resourceComputeGlobalAddress(),
			"google_compute_global_forwarding_rule": resourceComputeGlobalForwardingRule(),
			"google_compute_health_check":           resourceComputeHealthCheck(),
			"google_compute_http_health_check":      resourceComputeHttpHealthCheck(),
			"google_compute_https_health_check":     resourceComputeHttpsHealthCheck(),
			"google_compute_image":                  resourceComputeImage(),
			"google_compute_instance":               resourceComputeInstance(),
			"google_compute_instance_group":         resourceComputeInstanceGroup(),
			"google_compute_instance_group_manager": resourceComputeInstanceGroupManager(),
			"google_compute_instance_template":      resourceComputeInstanceTemplate(),
			"google_compute_network":                resourceComputeNetwork(),
			"google_compute_network_peering":        resourceComputeNetworkPeering(),
			"google_compute_project_metadata":       resourceComputeProjectMetadata(),
			"google_compute_project_metadata_item":  resourceComputeProjectMetadataItem(),
			"google_compute_region_backend_service": resourceComputeRegionBackendService(),
			"google_compute_route":                  resourceComputeRoute(),
			"google_compute_router":                 resourceComputeRouter(),
			"google_compute_router_interface":       resourceComputeRouterInterface(),
			"google_compute_router_peer":            resourceComputeRouterPeer(),
			"google_compute_ssl_certificate":        resourceComputeSslCertificate(),
			"google_compute_subnetwork":             resourceComputeSubnetwork(),
			"google_compute_target_http_proxy":      resourceComputeTargetHttpProxy(),
			"google_compute_target_https_proxy":     resourceComputeTargetHttpsProxy(),
			"google_compute_target_pool":            resourceComputeTargetPool(),
			"google_compute_url_map":                resourceComputeUrlMap(),
			"google_compute_vpn_gateway":            resourceComputeVpnGateway(),
			"google_compute_vpn_tunnel":             resourceComputeVpnTunnel(),
			"google_container_cluster":              resourceContainerCluster(),
			"google_container_node_pool":            resourceContainerNodePool(),
			"google_dns_managed_zone":               resourceDnsManagedZone(),
			"google_dns_record_set":                 resourceDnsRecordSet(),
			"google_sourcerepo_repository":          resourceSourceRepoRepository(),
			"google_spanner_instance":               resourceSpannerInstance(),
			"google_spanner_database":               resourceSpannerDatabase(),
			"google_sql_database":                   resourceSqlDatabase(),
			"google_sql_database_instance":          resourceSqlDatabaseInstance(),
			"google_sql_user":                       resourceSqlUser(),
			"google_project":                        resourceGoogleProject(),
			"google_project_iam_policy":             resourceGoogleProjectIamPolicy(),
			"google_project_iam_binding":            resourceGoogleProjectIamBinding(),
			"google_project_iam_member":             resourceGoogleProjectIamMember(),
			"google_project_services":               resourceGoogleProjectServices(),
			"google_pubsub_topic":                   resourcePubsubTopic(),
			"google_pubsub_subscription":            resourcePubsubSubscription(),
			"google_runtimeconfig_config":           resourceRuntimeconfigConfig(),
			"google_runtimeconfig_variable":         resourceRuntimeconfigVariable(),
			"google_service_account":                resourceGoogleServiceAccount(),
			"google_storage_bucket":                 resourceStorageBucket(),
			"google_storage_bucket_acl":             resourceStorageBucketAcl(),
			"google_storage_bucket_object":          resourceStorageBucketObject(),
			"google_storage_object_acl":             resourceStorageObjectAcl(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	credentials := d.Get("credentials").(string)
	config := Config{
		Credentials: credentials,
		Project:     d.Get("project").(string),
		Region:      d.Get("region").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	var account accountFile
	if err := json.Unmarshal([]byte(creds), &account); err != nil {
		errors = append(errors,
			fmt.Errorf("credentials are not valid JSON '%s': %s", creds, err))
	}

	return
}
