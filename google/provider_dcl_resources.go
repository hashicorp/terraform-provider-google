// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"google_apikeys_key":                                        ResourceApikeysKey(),
	"google_assured_workloads_workload":                         ResourceAssuredWorkloadsWorkload(),
	"google_bigquery_reservation_assignment":                    ResourceBigqueryReservationAssignment(),
	"google_cloudbuild_worker_pool":                             ResourceCloudbuildWorkerPool(),
	"google_clouddeploy_delivery_pipeline":                      ResourceClouddeployDeliveryPipeline(),
	"google_clouddeploy_target":                                 ResourceClouddeployTarget(),
	"google_compute_firewall_policy":                            ResourceComputeFirewallPolicy(),
	"google_compute_firewall_policy_association":                ResourceComputeFirewallPolicyAssociation(),
	"google_compute_firewall_policy_rule":                       ResourceComputeFirewallPolicyRule(),
	"google_compute_region_network_firewall_policy":             ResourceComputeRegionNetworkFirewallPolicy(),
	"google_compute_network_firewall_policy":                    ResourceComputeNetworkFirewallPolicy(),
	"google_compute_network_firewall_policy_association":        ResourceComputeNetworkFirewallPolicyAssociation(),
	"google_compute_region_network_firewall_policy_association": ResourceComputeRegionNetworkFirewallPolicyAssociation(),
	"google_compute_network_firewall_policy_rule":               ResourceComputeNetworkFirewallPolicyRule(),
	"google_compute_region_network_firewall_policy_rule":        ResourceComputeRegionNetworkFirewallPolicyRule(),
	"google_container_aws_cluster":                              ResourceContainerAwsCluster(),
	"google_container_aws_node_pool":                            ResourceContainerAwsNodePool(),
	"google_container_azure_client":                             ResourceContainerAzureClient(),
	"google_container_azure_cluster":                            ResourceContainerAzureCluster(),
	"google_container_azure_node_pool":                          ResourceContainerAzureNodePool(),
	"google_dataplex_asset":                                     ResourceDataplexAsset(),
	"google_dataplex_lake":                                      ResourceDataplexLake(),
	"google_dataplex_zone":                                      ResourceDataplexZone(),
	"google_dataproc_workflow_template":                         ResourceDataprocWorkflowTemplate(),
	"google_eventarc_channel":                                   ResourceEventarcChannel(),
	"google_eventarc_google_channel_config":                     ResourceEventarcGoogleChannelConfig(),
	"google_eventarc_trigger":                                   ResourceEventarcTrigger(),
	"google_firebaserules_release":                              ResourceFirebaserulesRelease(),
	"google_firebaserules_ruleset":                              ResourceFirebaserulesRuleset(),
	"google_monitoring_monitored_project":                       ResourceMonitoringMonitoredProject(),
	"google_network_connectivity_hub":                           ResourceNetworkConnectivityHub(),
	"google_network_connectivity_spoke":                         ResourceNetworkConnectivitySpoke(),
	"google_org_policy_policy":                                  ResourceOrgPolicyPolicy(),
	"google_os_config_os_policy_assignment":                     ResourceOsConfigOsPolicyAssignment(),
	"google_privateca_certificate_template":                     ResourcePrivatecaCertificateTemplate(),
	"google_recaptcha_enterprise_key":                           ResourceRecaptchaEnterpriseKey(),
}
