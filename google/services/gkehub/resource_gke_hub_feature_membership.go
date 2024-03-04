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

package gkehub

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	gkehub "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/gkehub"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceGkeHubFeatureMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGkeHubFeatureMembershipCreate,
		Read:   resourceGkeHubFeatureMembershipRead,
		Update: resourceGkeHubFeatureMembershipUpdate,
		Delete: resourceGkeHubFeatureMembershipDelete,

		Importer: &schema.ResourceImporter{
			State: resourceGkeHubFeatureMembershipImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"feature": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The name of the feature",
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location of the feature",
			},

			"membership": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The name of the membership",
			},

			"configmanagement": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Config Management-specific spec.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementSchema(),
			},

			"membership_location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of the membership",
			},

			"mesh": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Manage Mesh Features",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipMeshSchema(),
			},

			"policycontroller": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Policy Controller-specific spec.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerSchema(),
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project of the feature",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"binauthz": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "**DEPRECATED** Binauthz configuration for the cluster. This field will be ignored and should not be set.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementBinauthzSchema(),
			},

			"config_sync": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Config Sync configuration for the cluster.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementConfigSyncSchema(),
			},

			"hierarchy_controller": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Hierarchy Controller configuration for the cluster.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementHierarchyControllerSchema(),
			},

			"policy_controller": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Policy Controller configuration for the cluster.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementPolicyControllerSchema(),
			},

			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Version of ACM to install. Defaults to the latest version.",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementBinauthzSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether binauthz is enabled in this cluster.",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementConfigSyncSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"git": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementConfigSyncGitSchema(),
			},

			"metrics_gcp_service_account_email": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The Email of the Google Cloud Service Account (GSA) used for exporting Config Sync metrics to Cloud Monitoring. The GSA should have the Monitoring Metric Writer(roles/monitoring.metricWriter) IAM role. The Kubernetes ServiceAccount `default` in the namespace `config-management-monitoring` should be bound to the GSA.",
			},

			"oci": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementConfigSyncOciSchema(),
			},

			"prevent_drift": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Set to true to enable the Config Sync admission webhook to prevent drifts. If set to `false`, disables the Config Sync admission webhook and does not prevent drifts.",
			},

			"source_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the Config Sync Repo is in \"hierarchical\" or \"unstructured\" mode.",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementConfigSyncGitSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gcp_service_account_email": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The GCP Service Account Email used for auth when secretType is gcpServiceAccount.",
			},

			"https_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL for the HTTPS proxy to be used when communicating with the Git repo.",
			},

			"policy_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path within the Git repository that represents the top level of the repo to sync. Default: the root directory of the repository.",
			},

			"secret_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of secret configured for access to the Git repo. Must be one of ssh, cookiefile, gcenode, token, gcpserviceaccount or none. The validation of this is case-sensitive.",
			},

			"sync_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The branch of the repository to sync from. Default: master.",
			},

			"sync_repo": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the Git repository to use as the source of truth.",
			},

			"sync_rev": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Git revision (tag or hash) to check out. Default HEAD.",
			},

			"sync_wait_secs": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Period in seconds between consecutive syncs. Default: 15.",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementConfigSyncOciSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gcp_service_account_email": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The GCP Service Account Email used for auth when secret_type is gcpserviceaccount. ",
			},

			"policy_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The absolute path of the directory that contains the local resources. Default: the root directory of the image.",
			},

			"secret_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of secret configured for access to the OCI Image. Must be one of gcenode, gcpserviceaccount or none. The validation of this is case-sensitive.",
			},

			"sync_repo": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OCI image repository URL for the package to sync from. e.g. LOCATION-docker.pkg.dev/PROJECT_ID/REPOSITORY_NAME/PACKAGE_NAME.",
			},

			"sync_wait_secs": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Period in seconds(int64 format) between consecutive syncs. Default: 15.",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementHierarchyControllerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_hierarchical_resource_quota": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether hierarchical resource quota is enabled in this cluster.",
			},

			"enable_pod_tree_labels": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether pod tree labels are enabled in this cluster.",
			},

			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether Hierarchy Controller is enabled in this cluster.",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementPolicyControllerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"audit_interval_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the interval for Policy Controller Audit Scans (in seconds). When set to 0, this disables audit functionality altogether.",
			},

			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables the installation of Policy Controller. If false, the rest of PolicyController fields take no effect.",
			},

			"exemptable_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The set of namespaces that are excluded from Policy Controller checks. Namespaces do not need to currently exist on the cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"log_denies_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Logs all denies and dry run failures.",
			},

			"monitoring": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Specifies the backends Policy Controller should export metrics to. For example, to specify metrics should be exported to Cloud Monitoring and Prometheus, specify backends: [\"cloudmonitoring\", \"prometheus\"]. Default: [\"cloudmonitoring\", \"prometheus\"]",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoringSchema(),
			},

			"mutation_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable mutation in policy controller. If true, mutation CRDs, webhook and controller deployment will be deployed to the cluster.",
			},

			"referential_rules_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables the ability to use Constraint Templates that reference to objects other than the object currently being evaluated.",
			},

			"template_library_installed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Installs the default template library along with Policy Controller.",
			},
		},
	}
}

func GkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoringSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backends": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: " Specifies the list of backends Policy Controller will export to. Specifying an empty value `[]` disables metrics export.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func GkeHubFeatureMembershipMeshSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"control_plane": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "**DEPRECATED** Whether to automatically manage Service Mesh control planes. Possible values: CONTROL_PLANE_MANAGEMENT_UNSPECIFIED, AUTOMATIC, MANUAL",
				Deprecated:  "Deprecated in favor of the `management` field",
			},

			"management": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether to automatically manage Service Mesh. Possible values: MANAGEMENT_UNSPECIFIED, MANAGEMENT_AUTOMATIC, MANAGEMENT_MANUAL",
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_controller_hub_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Policy Controller configuration for the cluster.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigSchema(),
			},

			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Version of Policy Controller to install. Defaults to the latest version.",
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"audit_interval_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the interval for Policy Controller Audit Scans (in seconds). When set to 0, this disables audit functionality altogether.",
			},

			"constraint_violation_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of audit violations to be stored in a constraint. If not set, the internal default of 20 will be used.",
			},

			"deployment_configs": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Description: "Map of deployment configs to deployments (\"admission\", \"audit\", \"mutation\").",
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsSchema(),
				Set:         schema.HashResource(GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsSchema()),
			},

			"exemptable_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The set of namespaces that are excluded from Policy Controller checks. Namespaces do not need to currently exist on the cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"install_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Configures the mode of the Policy Controller installation. Possible values: INSTALL_SPEC_UNSPECIFIED, INSTALL_SPEC_NOT_INSTALLED, INSTALL_SPEC_ENABLED, INSTALL_SPEC_SUSPENDED, INSTALL_SPEC_DETACHED",
			},

			"log_denies_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Logs all denies and dry run failures.",
			},

			"monitoring": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Specifies the backends Policy Controller should export metrics to. For example, to specify metrics should be exported to Cloud Monitoring and Prometheus, specify backends: [\"cloudmonitoring\", \"prometheus\"]. Default: [\"cloudmonitoring\", \"prometheus\"]",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringSchema(),
			},

			"mutation_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables the ability to mutate resources using Policy Controller.",
			},

			"policy_content": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Specifies the desired policy content on the cluster.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentSchema(),
			},

			"referential_rules_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables the ability to use Constraint Templates that reference to objects other than the object currently being evaluated.",
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"component_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the key in the map for which this object is mapped to in the API",
			},

			"container_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Container resource requirements.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesSchema(),
			},

			"pod_affinity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pod affinity configuration. Possible values: AFFINITY_UNSPECIFIED, NO_AFFINITY, ANTI_AFFINITY",
			},

			"pod_tolerations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Pod tolerations of node taints.",
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerationsSchema(),
			},

			"replica_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pod replica count.",
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"limits": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Limits describes the maximum amount of compute resources allowed for use by the running container.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimitsSchema(),
			},

			"requests": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Requests describes the amount of compute resources reserved for the container by the kube-scheduler.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequestsSchema(),
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimitsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cpu": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CPU requirement expressed in Kubernetes resource units.",
			},

			"memory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Memory requirement expressed in Kubernetes resource units.",
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequestsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cpu": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CPU requirement expressed in Kubernetes resource units.",
			},

			"memory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Memory requirement expressed in Kubernetes resource units.",
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"effect": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Matches a taint effect.",
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Matches a taint key (not necessarily unique).",
			},

			"operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Matches a taint operator.",
			},

			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Matches a taint value.",
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backends": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: " Specifies the list of backends Policy Controller will export to. Specifying an empty value `[]` disables metrics export.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bundles": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "map of bundle name to BundleInstallSpec. The bundle name maps to the `bundleName` key in the `policycontroller.gke.io/constraintData` annotation on a constraint.",
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundlesSchema(),
				Set:         schema.HashResource(GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundlesSchema()),
			},

			"template_library": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Configures the installation of the Template Library.",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrarySchema(),
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundlesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bundle_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the key in the map for which this object is mapped to in the API",
			},

			"exempted_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The set of namespaces to be exempted from the bundle.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func GkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrarySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"installation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Configures the manner in which the template library is installed on the cluster. Possible values: INSTALLATION_UNSPECIFIED, NOT_INSTALLED, ALL",
			},
		},
	}
}

func resourceGkeHubFeatureMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &gkehub.FeatureMembership{
		Feature:            dcl.String(d.Get("feature").(string)),
		Location:           dcl.String(d.Get("location").(string)),
		Membership:         dcl.String(d.Get("membership").(string)),
		Configmanagement:   expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		MembershipLocation: dcl.String(d.Get("membership_location").(string)),
		Mesh:               expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Policycontroller:   expandGkeHubFeatureMembershipPolicycontroller(d.Get("policycontroller")),
		Project:            dcl.String(project),
	}
	lockName, err := tpgresource.ReplaceVarsForId(d, config, "{{project}}/{{location}}/{{feature}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}")
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLGkeHubClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyFeatureMembership(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating FeatureMembership: %s", err)
	}

	log.Printf("[DEBUG] Finished creating FeatureMembership %q: %#v", d.Id(), res)

	return resourceGkeHubFeatureMembershipRead(d, meta)
}

func resourceGkeHubFeatureMembershipRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &gkehub.FeatureMembership{
		Feature:            dcl.String(d.Get("feature").(string)),
		Location:           dcl.String(d.Get("location").(string)),
		Membership:         dcl.String(d.Get("membership").(string)),
		Configmanagement:   expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		MembershipLocation: dcl.String(d.Get("membership_location").(string)),
		Mesh:               expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Policycontroller:   expandGkeHubFeatureMembershipPolicycontroller(d.Get("policycontroller")),
		Project:            dcl.String(project),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLGkeHubClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetFeatureMembership(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("GkeHubFeatureMembership %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("feature", res.Feature); err != nil {
		return fmt.Errorf("error setting feature in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("membership", res.Membership); err != nil {
		return fmt.Errorf("error setting membership in state: %s", err)
	}
	if err = d.Set("configmanagement", flattenGkeHubFeatureMembershipConfigmanagement(res.Configmanagement)); err != nil {
		return fmt.Errorf("error setting configmanagement in state: %s", err)
	}
	if err = d.Set("membership_location", res.MembershipLocation); err != nil {
		return fmt.Errorf("error setting membership_location in state: %s", err)
	}
	if err = d.Set("mesh", flattenGkeHubFeatureMembershipMesh(res.Mesh)); err != nil {
		return fmt.Errorf("error setting mesh in state: %s", err)
	}
	if err = d.Set("policycontroller", flattenGkeHubFeatureMembershipPolicycontroller(res.Policycontroller)); err != nil {
		return fmt.Errorf("error setting policycontroller in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}

	return nil
}
func resourceGkeHubFeatureMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &gkehub.FeatureMembership{
		Feature:            dcl.String(d.Get("feature").(string)),
		Location:           dcl.String(d.Get("location").(string)),
		Membership:         dcl.String(d.Get("membership").(string)),
		Configmanagement:   expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		MembershipLocation: dcl.String(d.Get("membership_location").(string)),
		Mesh:               expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Policycontroller:   expandGkeHubFeatureMembershipPolicycontroller(d.Get("policycontroller")),
		Project:            dcl.String(project),
	}
	lockName, err := tpgresource.ReplaceVarsForId(d, config, "{{project}}/{{location}}/{{feature}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLGkeHubClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyFeatureMembership(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating FeatureMembership: %s", err)
	}

	log.Printf("[DEBUG] Finished creating FeatureMembership %q: %#v", d.Id(), res)

	return resourceGkeHubFeatureMembershipRead(d, meta)
}

func resourceGkeHubFeatureMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &gkehub.FeatureMembership{
		Feature:            dcl.String(d.Get("feature").(string)),
		Location:           dcl.String(d.Get("location").(string)),
		Membership:         dcl.String(d.Get("membership").(string)),
		Configmanagement:   expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		MembershipLocation: dcl.String(d.Get("membership_location").(string)),
		Mesh:               expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Policycontroller:   expandGkeHubFeatureMembershipPolicycontroller(d.Get("policycontroller")),
		Project:            dcl.String(project),
	}
	lockName, err := tpgresource.ReplaceVarsForId(d, config, "{{project}}/{{location}}/{{feature}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	log.Printf("[DEBUG] Deleting FeatureMembership %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLGkeHubClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteFeatureMembership(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting FeatureMembership: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting FeatureMembership %q", d.Id())
	return nil
}

func resourceGkeHubFeatureMembershipImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/features/(?P<feature>[^/]+)/membershipId/(?P<membership>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<feature>[^/]+)/(?P<membership>[^/]+)",
		"(?P<location>[^/]+)/(?P<feature>[^/]+)/(?P<membership>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandGkeHubFeatureMembershipConfigmanagement(o interface{}) *gkehub.FeatureMembershipConfigmanagement {
	if o == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagement
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagement
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagement{
		Binauthz:            expandGkeHubFeatureMembershipConfigmanagementBinauthz(obj["binauthz"]),
		ConfigSync:          expandGkeHubFeatureMembershipConfigmanagementConfigSync(obj["config_sync"]),
		HierarchyController: expandGkeHubFeatureMembershipConfigmanagementHierarchyController(obj["hierarchy_controller"]),
		PolicyController:    expandGkeHubFeatureMembershipConfigmanagementPolicyController(obj["policy_controller"]),
		Version:             dcl.StringOrNil(obj["version"].(string)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagement(obj *gkehub.FeatureMembershipConfigmanagement) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"binauthz":             flattenGkeHubFeatureMembershipConfigmanagementBinauthz(obj.Binauthz),
		"config_sync":          flattenGkeHubFeatureMembershipConfigmanagementConfigSync(obj.ConfigSync),
		"hierarchy_controller": flattenGkeHubFeatureMembershipConfigmanagementHierarchyController(obj.HierarchyController),
		"policy_controller":    flattenGkeHubFeatureMembershipConfigmanagementPolicyController(obj.PolicyController),
		"version":              obj.Version,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipConfigmanagementBinauthz(o interface{}) *gkehub.FeatureMembershipConfigmanagementBinauthz {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagementBinauthz{
		Enabled: dcl.Bool(obj["enabled"].(bool)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementBinauthz(obj *gkehub.FeatureMembershipConfigmanagementBinauthz) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled": obj.Enabled,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipConfigmanagementConfigSync(o interface{}) *gkehub.FeatureMembershipConfigmanagementConfigSync {
	if o == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementConfigSync
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementConfigSync
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagementConfigSync{
		Git:                           expandGkeHubFeatureMembershipConfigmanagementConfigSyncGit(obj["git"]),
		MetricsGcpServiceAccountEmail: dcl.String(obj["metrics_gcp_service_account_email"].(string)),
		Oci:                           expandGkeHubFeatureMembershipConfigmanagementConfigSyncOci(obj["oci"]),
		PreventDrift:                  dcl.Bool(obj["prevent_drift"].(bool)),
		SourceFormat:                  dcl.String(obj["source_format"].(string)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementConfigSync(obj *gkehub.FeatureMembershipConfigmanagementConfigSync) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"git":                               flattenGkeHubFeatureMembershipConfigmanagementConfigSyncGit(obj.Git),
		"metrics_gcp_service_account_email": obj.MetricsGcpServiceAccountEmail,
		"oci":                               flattenGkeHubFeatureMembershipConfigmanagementConfigSyncOci(obj.Oci),
		"prevent_drift":                     obj.PreventDrift,
		"source_format":                     obj.SourceFormat,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipConfigmanagementConfigSyncGit(o interface{}) *gkehub.FeatureMembershipConfigmanagementConfigSyncGit {
	if o == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementConfigSyncGit
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementConfigSyncGit
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagementConfigSyncGit{
		GcpServiceAccountEmail: dcl.String(obj["gcp_service_account_email"].(string)),
		HttpsProxy:             dcl.String(obj["https_proxy"].(string)),
		PolicyDir:              dcl.String(obj["policy_dir"].(string)),
		SecretType:             dcl.String(obj["secret_type"].(string)),
		SyncBranch:             dcl.String(obj["sync_branch"].(string)),
		SyncRepo:               dcl.String(obj["sync_repo"].(string)),
		SyncRev:                dcl.String(obj["sync_rev"].(string)),
		SyncWaitSecs:           dcl.String(obj["sync_wait_secs"].(string)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementConfigSyncGit(obj *gkehub.FeatureMembershipConfigmanagementConfigSyncGit) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"gcp_service_account_email": obj.GcpServiceAccountEmail,
		"https_proxy":               obj.HttpsProxy,
		"policy_dir":                obj.PolicyDir,
		"secret_type":               obj.SecretType,
		"sync_branch":               obj.SyncBranch,
		"sync_repo":                 obj.SyncRepo,
		"sync_rev":                  obj.SyncRev,
		"sync_wait_secs":            obj.SyncWaitSecs,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipConfigmanagementConfigSyncOci(o interface{}) *gkehub.FeatureMembershipConfigmanagementConfigSyncOci {
	if o == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementConfigSyncOci
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementConfigSyncOci
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagementConfigSyncOci{
		GcpServiceAccountEmail: dcl.String(obj["gcp_service_account_email"].(string)),
		PolicyDir:              dcl.String(obj["policy_dir"].(string)),
		SecretType:             dcl.String(obj["secret_type"].(string)),
		SyncRepo:               dcl.String(obj["sync_repo"].(string)),
		SyncWaitSecs:           dcl.String(obj["sync_wait_secs"].(string)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementConfigSyncOci(obj *gkehub.FeatureMembershipConfigmanagementConfigSyncOci) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"gcp_service_account_email": obj.GcpServiceAccountEmail,
		"policy_dir":                obj.PolicyDir,
		"secret_type":               obj.SecretType,
		"sync_repo":                 obj.SyncRepo,
		"sync_wait_secs":            obj.SyncWaitSecs,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipConfigmanagementHierarchyController(o interface{}) *gkehub.FeatureMembershipConfigmanagementHierarchyController {
	if o == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementHierarchyController
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementHierarchyController
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagementHierarchyController{
		EnableHierarchicalResourceQuota: dcl.Bool(obj["enable_hierarchical_resource_quota"].(bool)),
		EnablePodTreeLabels:             dcl.Bool(obj["enable_pod_tree_labels"].(bool)),
		Enabled:                         dcl.Bool(obj["enabled"].(bool)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementHierarchyController(obj *gkehub.FeatureMembershipConfigmanagementHierarchyController) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"enable_hierarchical_resource_quota": obj.EnableHierarchicalResourceQuota,
		"enable_pod_tree_labels":             obj.EnablePodTreeLabels,
		"enabled":                            obj.Enabled,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipConfigmanagementPolicyController(o interface{}) *gkehub.FeatureMembershipConfigmanagementPolicyController {
	if o == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementPolicyController
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementPolicyController
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagementPolicyController{
		AuditIntervalSeconds:     dcl.String(obj["audit_interval_seconds"].(string)),
		Enabled:                  dcl.Bool(obj["enabled"].(bool)),
		ExemptableNamespaces:     tpgdclresource.ExpandStringArray(obj["exemptable_namespaces"]),
		LogDeniesEnabled:         dcl.Bool(obj["log_denies_enabled"].(bool)),
		Monitoring:               expandGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoring(obj["monitoring"]),
		MutationEnabled:          dcl.Bool(obj["mutation_enabled"].(bool)),
		ReferentialRulesEnabled:  dcl.Bool(obj["referential_rules_enabled"].(bool)),
		TemplateLibraryInstalled: dcl.Bool(obj["template_library_installed"].(bool)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementPolicyController(obj *gkehub.FeatureMembershipConfigmanagementPolicyController) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"audit_interval_seconds":     obj.AuditIntervalSeconds,
		"enabled":                    obj.Enabled,
		"exemptable_namespaces":      obj.ExemptableNamespaces,
		"log_denies_enabled":         obj.LogDeniesEnabled,
		"monitoring":                 flattenGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoring(obj.Monitoring),
		"mutation_enabled":           obj.MutationEnabled,
		"referential_rules_enabled":  obj.ReferentialRulesEnabled,
		"template_library_installed": obj.TemplateLibraryInstalled,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoring(o interface{}) *gkehub.FeatureMembershipConfigmanagementPolicyControllerMonitoring {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipConfigmanagementPolicyControllerMonitoring{
		Backends: expandGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsArray(obj["backends"]),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoring(obj *gkehub.FeatureMembershipConfigmanagementPolicyControllerMonitoring) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"backends": flattenGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsArray(obj.Backends),
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipMesh(o interface{}) *gkehub.FeatureMembershipMesh {
	if o == nil {
		return gkehub.EmptyFeatureMembershipMesh
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipMesh
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipMesh{
		ControlPlane: gkehub.FeatureMembershipMeshControlPlaneEnumRef(obj["control_plane"].(string)),
		Management:   gkehub.FeatureMembershipMeshManagementEnumRef(obj["management"].(string)),
	}
}

func flattenGkeHubFeatureMembershipMesh(obj *gkehub.FeatureMembershipMesh) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"control_plane": obj.ControlPlane,
		"management":    obj.Management,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipPolicycontroller(o interface{}) *gkehub.FeatureMembershipPolicycontroller {
	if o == nil {
		return gkehub.EmptyFeatureMembershipPolicycontroller
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipPolicycontroller
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontroller{
		PolicyControllerHubConfig: expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfig(obj["policy_controller_hub_config"]),
		Version:                   dcl.StringOrNil(obj["version"].(string)),
	}
}

func flattenGkeHubFeatureMembershipPolicycontroller(obj *gkehub.FeatureMembershipPolicycontroller) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"policy_controller_hub_config": flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfig(obj.PolicyControllerHubConfig),
		"version":                      obj.Version,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfig(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfig {
	if o == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfig{
		AuditIntervalSeconds:     dcl.Int64(int64(obj["audit_interval_seconds"].(int))),
		ConstraintViolationLimit: dcl.Int64(int64(obj["constraint_violation_limit"].(int))),
		DeploymentConfigs:        expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsMap(obj["deployment_configs"]),
		ExemptableNamespaces:     tpgdclresource.ExpandStringArray(obj["exemptable_namespaces"]),
		InstallSpec:              gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigInstallSpecEnumRef(obj["install_spec"].(string)),
		LogDeniesEnabled:         dcl.Bool(obj["log_denies_enabled"].(bool)),
		Monitoring:               expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring(obj["monitoring"]),
		MutationEnabled:          dcl.Bool(obj["mutation_enabled"].(bool)),
		PolicyContent:            expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent(obj["policy_content"]),
		ReferentialRulesEnabled:  dcl.Bool(obj["referential_rules_enabled"].(bool)),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfig(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"audit_interval_seconds":     obj.AuditIntervalSeconds,
		"constraint_violation_limit": obj.ConstraintViolationLimit,
		"deployment_configs":         flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsMap(obj.DeploymentConfigs),
		"exemptable_namespaces":      obj.ExemptableNamespaces,
		"install_spec":               obj.InstallSpec,
		"log_denies_enabled":         obj.LogDeniesEnabled,
		"monitoring":                 flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring(obj.Monitoring),
		"mutation_enabled":           obj.MutationEnabled,
		"policy_content":             flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent(obj.PolicyContent),
		"referential_rules_enabled":  obj.ReferentialRulesEnabled,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsMap(o interface{}) map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs {
	if o == nil {
		return nil
	}

	o = o.(*schema.Set).List()

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return nil
	}

	items := make(map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs)
	for _, item := range objs {
		i := expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs(item)
		if item != nil {
			items[item.(map[string]interface{})["component_name"].(string)] = *i
		}
	}

	return items
}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs {
	if o == nil {
		return nil
	}

	obj := o.(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs{
		ContainerResources: expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources(obj["container_resources"]),
		PodAffinity:        gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodAffinityEnumRef(obj["pod_affinity"].(string)),
		PodTolerations:     expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerationsArray(obj["pod_tolerations"]),
		ReplicaCount:       dcl.Int64(int64(obj["replica_count"].(int))),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsMap(objs map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for name, item := range objs {
		i := flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs(&item, name)
		items = append(items, i)
	}

	return items
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs, name string) interface{} {
	if obj == nil {
		return nil
	}
	transformed := map[string]interface{}{
		"container_resources": flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources(obj.ContainerResources),
		"pod_affinity":        obj.PodAffinity,
		"pod_tolerations":     flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerationsArray(obj.PodTolerations),
		"replica_count":       obj.ReplicaCount,
	}

	transformed["component_name"] = name

	return transformed

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources {
	if o == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources{
		Limits:   expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits(obj["limits"]),
		Requests: expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests(obj["requests"]),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"limits":   flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits(obj.Limits),
		"requests": flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests(obj.Requests),
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits {
	if o == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits{
		Cpu:    dcl.String(obj["cpu"].(string)),
		Memory: dcl.String(obj["memory"].(string)),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cpu":    obj.Cpu,
		"memory": obj.Memory,
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests {
	if o == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests{
		Cpu:    dcl.String(obj["cpu"].(string)),
		Memory: dcl.String(obj["memory"].(string)),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cpu":    obj.Cpu,
		"memory": obj.Memory,
	}

	return []interface{}{transformed}

}
func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerationsArray(o interface{}) []gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations {
	if o == nil {
		return make([]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations, 0)
	}

	items := make([]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations, 0, len(objs))
	for _, item := range objs {
		i := expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations(item)
		items = append(items, *i)
	}

	return items
}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations {
	if o == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations
	}

	obj := o.(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations{
		Effect:   dcl.String(obj["effect"].(string)),
		Key:      dcl.String(obj["key"].(string)),
		Operator: dcl.String(obj["operator"].(string)),
		Value:    dcl.String(obj["value"].(string)),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerationsArray(objs []gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations(&item)
		items = append(items, i)
	}

	return items
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"effect":   obj.Effect,
		"key":      obj.Key,
		"operator": obj.Operator,
		"value":    obj.Value,
	}

	return transformed

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring{
		Backends: expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsArray(obj["backends"]),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"backends": flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsArray(obj.Backends),
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent{
		Bundles:         expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundlesMap(obj["bundles"]),
		TemplateLibrary: expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary(obj["template_library"]),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"bundles":          flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundlesMap(obj.Bundles),
		"template_library": flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary(obj.TemplateLibrary),
	}

	return []interface{}{transformed}

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundlesMap(o interface{}) map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles {
	if o == nil {
		return make(map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles)
	}

	o = o.(*schema.Set).List()

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make(map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles)
	}

	items := make(map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles)
	for _, item := range objs {
		i := expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles(item)
		if item != nil {
			items[item.(map[string]interface{})["bundle_name"].(string)] = *i
		}
	}

	return items
}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles {
	if o == nil {
		return gkehub.EmptyFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles
	}

	obj := o.(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles{
		ExemptedNamespaces: tpgdclresource.ExpandStringArray(obj["exempted_namespaces"]),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundlesMap(objs map[string]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for name, item := range objs {
		i := flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles(&item, name)
		items = append(items, i)
	}

	return items
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles, name string) interface{} {
	if obj == nil {
		return nil
	}
	transformed := map[string]interface{}{
		"exempted_namespaces": obj.ExemptedNamespaces,
	}

	transformed["bundle_name"] = name

	return transformed

}

func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary(o interface{}) *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary{
		Installation: gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibraryInstallationEnumRef(obj["installation"].(string)),
	}
}

func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary(obj *gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"installation": obj.Installation,
	}

	return []interface{}{transformed}

}

func flattenGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsArray(obj []gkehub.FeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsEnum) interface{} {
	if obj == nil {
		return nil
	}
	items := []string{}
	for _, item := range obj {
		items = append(items, string(item))
	}
	return items
}
func expandGkeHubFeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsArray(o interface{}) []gkehub.FeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsEnum {
	objs := o.([]interface{})
	items := make([]gkehub.FeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsEnum, 0, len(objs))
	for _, item := range objs {
		i := gkehub.FeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsEnumRef(item.(string))
		items = append(items, *i)
	}
	return items
}
func flattenGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsArray(obj []gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsEnum) interface{} {
	if obj == nil {
		return nil
	}
	items := []string{}
	for _, item := range obj {
		items = append(items, string(item))
	}
	return items
}
func expandGkeHubFeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsArray(o interface{}) []gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsEnum {
	objs := o.([]interface{})
	items := make([]gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsEnum, 0, len(objs))
	for _, item := range objs {
		i := gkehub.FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsEnumRef(item.(string))
		items = append(items, *i)
	}
	return items
}
