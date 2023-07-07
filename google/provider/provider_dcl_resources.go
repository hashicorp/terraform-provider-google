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

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/services/apikeys"
	"github.com/hashicorp/terraform-provider-google/google/services/assuredworkloads"
	"github.com/hashicorp/terraform-provider-google/google/services/bigqueryreservation"
	"github.com/hashicorp/terraform-provider-google/google/services/cloudbuild"
	"github.com/hashicorp/terraform-provider-google/google/services/cloudbuildv2"
	"github.com/hashicorp/terraform-provider-google/google/services/clouddeploy"
	"github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/services/containeraws"
	"github.com/hashicorp/terraform-provider-google/google/services/containerazure"
	"github.com/hashicorp/terraform-provider-google/google/services/dataplex"
	"github.com/hashicorp/terraform-provider-google/google/services/dataproc"
	"github.com/hashicorp/terraform-provider-google/google/services/eventarc"
	"github.com/hashicorp/terraform-provider-google/google/services/firebaserules"
	"github.com/hashicorp/terraform-provider-google/google/services/networkconnectivity"
	"github.com/hashicorp/terraform-provider-google/google/services/orgpolicy"
	"github.com/hashicorp/terraform-provider-google/google/services/privateca"
	"github.com/hashicorp/terraform-provider-google/google/services/recaptchaenterprise"
)

var dclResources = map[string]*schema.Resource{
	"google_apikeys_key":                                        apikeys.ResourceApikeysKey(),
	"google_assured_workloads_workload":                         assuredworkloads.ResourceAssuredWorkloadsWorkload(),
	"google_bigquery_reservation_assignment":                    bigqueryreservation.ResourceBigqueryReservationAssignment(),
	"google_cloudbuild_worker_pool":                             cloudbuild.ResourceCloudbuildWorkerPool(),
	"google_cloudbuildv2_connection":                            cloudbuildv2.ResourceCloudbuildv2Connection(),
	"google_cloudbuildv2_repository":                            cloudbuildv2.ResourceCloudbuildv2Repository(),
	"google_clouddeploy_delivery_pipeline":                      clouddeploy.ResourceClouddeployDeliveryPipeline(),
	"google_clouddeploy_target":                                 clouddeploy.ResourceClouddeployTarget(),
	"google_compute_firewall_policy":                            compute.ResourceComputeFirewallPolicy(),
	"google_compute_firewall_policy_association":                compute.ResourceComputeFirewallPolicyAssociation(),
	"google_compute_firewall_policy_rule":                       compute.ResourceComputeFirewallPolicyRule(),
	"google_compute_region_network_firewall_policy":             compute.ResourceComputeRegionNetworkFirewallPolicy(),
	"google_compute_network_firewall_policy":                    compute.ResourceComputeNetworkFirewallPolicy(),
	"google_compute_network_firewall_policy_association":        compute.ResourceComputeNetworkFirewallPolicyAssociation(),
	"google_compute_region_network_firewall_policy_association": compute.ResourceComputeRegionNetworkFirewallPolicyAssociation(),
	"google_compute_network_firewall_policy_rule":               compute.ResourceComputeNetworkFirewallPolicyRule(),
	"google_compute_region_network_firewall_policy_rule":        compute.ResourceComputeRegionNetworkFirewallPolicyRule(),
	"google_container_aws_cluster":                              containeraws.ResourceContainerAwsCluster(),
	"google_container_aws_node_pool":                            containeraws.ResourceContainerAwsNodePool(),
	"google_container_azure_client":                             containerazure.ResourceContainerAzureClient(),
	"google_container_azure_cluster":                            containerazure.ResourceContainerAzureCluster(),
	"google_container_azure_node_pool":                          containerazure.ResourceContainerAzureNodePool(),
	"google_dataplex_asset":                                     dataplex.ResourceDataplexAsset(),
	"google_dataplex_lake":                                      dataplex.ResourceDataplexLake(),
	"google_dataplex_zone":                                      dataplex.ResourceDataplexZone(),
	"google_dataproc_workflow_template":                         dataproc.ResourceDataprocWorkflowTemplate(),
	"google_eventarc_channel":                                   eventarc.ResourceEventarcChannel(),
	"google_eventarc_google_channel_config":                     eventarc.ResourceEventarcGoogleChannelConfig(),
	"google_eventarc_trigger":                                   eventarc.ResourceEventarcTrigger(),
	"google_firebaserules_release":                              firebaserules.ResourceFirebaserulesRelease(),
	"google_firebaserules_ruleset":                              firebaserules.ResourceFirebaserulesRuleset(),
	"google_network_connectivity_hub":                           networkconnectivity.ResourceNetworkConnectivityHub(),
	"google_network_connectivity_spoke":                         networkconnectivity.ResourceNetworkConnectivitySpoke(),
	"google_org_policy_policy":                                  orgpolicy.ResourceOrgPolicyPolicy(),
	"google_privateca_certificate_template":                     privateca.ResourcePrivatecaCertificateTemplate(),
	"google_recaptcha_enterprise_key":                           recaptchaenterprise.ResourceRecaptchaEnterpriseKey(),
}
