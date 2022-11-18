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

package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var dclResources = map[string]*schema.Resource{
	"google_apikeys_key":                                        resourceApikeysKey(),
	"google_assured_workloads_workload":                         resourceAssuredWorkloadsWorkload(),
	"google_bigquery_reservation_assignment":                    resourceBigqueryReservationAssignment(),
	"google_cloudbuild_worker_pool":                             resourceCloudbuildWorkerPool(),
	"google_clouddeploy_delivery_pipeline":                      resourceClouddeployDeliveryPipeline(),
	"google_clouddeploy_target":                                 resourceClouddeployTarget(),
	"google_compute_firewall_policy":                            resourceComputeFirewallPolicy(),
	"google_compute_firewall_policy_association":                resourceComputeFirewallPolicyAssociation(),
	"google_compute_firewall_policy_rule":                       resourceComputeFirewallPolicyRule(),
	"google_compute_region_network_firewall_policy":             resourceComputeRegionNetworkFirewallPolicy(),
	"google_compute_network_firewall_policy":                    resourceComputeNetworkFirewallPolicy(),
	"google_compute_network_firewall_policy_association":        resourceComputeNetworkFirewallPolicyAssociation(),
	"google_compute_region_network_firewall_policy_association": resourceComputeRegionNetworkFirewallPolicyAssociation(),
	"google_compute_network_firewall_policy_rule":               resourceComputeNetworkFirewallPolicyRule(),
	"google_compute_region_network_firewall_policy_rule":        resourceComputeRegionNetworkFirewallPolicyRule(),
	"google_container_aws_cluster":                              resourceContainerAwsCluster(),
	"google_container_aws_node_pool":                            resourceContainerAwsNodePool(),
	"google_container_azure_client":                             resourceContainerAzureClient(),
	"google_container_azure_cluster":                            resourceContainerAzureCluster(),
	"google_container_azure_node_pool":                          resourceContainerAzureNodePool(),
	"google_dataplex_asset":                                     resourceDataplexAsset(),
	"google_dataplex_lake":                                      resourceDataplexLake(),
	"google_dataplex_zone":                                      resourceDataplexZone(),
	"google_dataproc_workflow_template":                         resourceDataprocWorkflowTemplate(),
	"google_eventarc_channel":                                   resourceEventarcChannel(),
	"google_eventarc_google_channel_config":                     resourceEventarcGoogleChannelConfig(),
	"google_eventarc_trigger":                                   resourceEventarcTrigger(),
	"google_firebaserules_release":                              resourceFirebaserulesRelease(),
	"google_firebaserules_ruleset":                              resourceFirebaserulesRuleset(),
	"google_logging_log_view":                                   resourceLoggingLogView(),
	"google_monitoring_monitored_project":                       resourceMonitoringMonitoredProject(),
	"google_network_connectivity_hub":                           resourceNetworkConnectivityHub(),
	"google_network_connectivity_spoke":                         resourceNetworkConnectivitySpoke(),
	"google_org_policy_policy":                                  resourceOrgPolicyPolicy(),
	"google_os_config_os_policy_assignment":                     resourceOsConfigOsPolicyAssignment(),
	"google_privateca_certificate_template":                     resourcePrivatecaCertificateTemplate(),
	"google_recaptcha_enterprise_key":                           resourceRecaptchaEnterpriseKey(),
}
