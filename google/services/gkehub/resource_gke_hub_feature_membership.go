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

			"mesh": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Manage Mesh Features",
				MaxItems:    1,
				Elem:        GkeHubFeatureMembershipMeshSchema(),
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
				Optional:    true,
				Description: "Binauthz configuration for the cluster.",
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

func resourceGkeHubFeatureMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &gkehub.FeatureMembership{
		Feature:          dcl.String(d.Get("feature").(string)),
		Location:         dcl.String(d.Get("location").(string)),
		Membership:       dcl.String(d.Get("membership").(string)),
		Configmanagement: expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		Mesh:             expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Project:          dcl.String(project),
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
		Feature:          dcl.String(d.Get("feature").(string)),
		Location:         dcl.String(d.Get("location").(string)),
		Membership:       dcl.String(d.Get("membership").(string)),
		Configmanagement: expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		Mesh:             expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Project:          dcl.String(project),
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
	if err = d.Set("mesh", flattenGkeHubFeatureMembershipMesh(res.Mesh)); err != nil {
		return fmt.Errorf("error setting mesh in state: %s", err)
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
		Feature:          dcl.String(d.Get("feature").(string)),
		Location:         dcl.String(d.Get("location").(string)),
		Membership:       dcl.String(d.Get("membership").(string)),
		Configmanagement: expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		Mesh:             expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Project:          dcl.String(project),
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
		Feature:          dcl.String(d.Get("feature").(string)),
		Location:         dcl.String(d.Get("location").(string)),
		Membership:       dcl.String(d.Get("membership").(string)),
		Configmanagement: expandGkeHubFeatureMembershipConfigmanagement(d.Get("configmanagement")),
		Mesh:             expandGkeHubFeatureMembershipMesh(d.Get("mesh")),
		Project:          dcl.String(project),
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
		return gkehub.EmptyFeatureMembershipConfigmanagementBinauthz
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return gkehub.EmptyFeatureMembershipConfigmanagementBinauthz
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
		Git:          expandGkeHubFeatureMembershipConfigmanagementConfigSyncGit(obj["git"]),
		Oci:          expandGkeHubFeatureMembershipConfigmanagementConfigSyncOci(obj["oci"]),
		PreventDrift: dcl.Bool(obj["prevent_drift"].(bool)),
		SourceFormat: dcl.String(obj["source_format"].(string)),
	}
}

func flattenGkeHubFeatureMembershipConfigmanagementConfigSync(obj *gkehub.FeatureMembershipConfigmanagementConfigSync) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"git":           flattenGkeHubFeatureMembershipConfigmanagementConfigSyncGit(obj.Git),
		"oci":           flattenGkeHubFeatureMembershipConfigmanagementConfigSyncOci(obj.Oci),
		"prevent_drift": obj.PreventDrift,
		"source_format": obj.SourceFormat,
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
