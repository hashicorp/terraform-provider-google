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

package assuredworkloads

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceAssuredWorkloadsWorkload() *schema.Resource {
	return &schema.Resource{
		Create: resourceAssuredWorkloadsWorkloadCreate,
		Read:   resourceAssuredWorkloadsWorkloadRead,
		Update: resourceAssuredWorkloadsWorkloadUpdate,
		Delete: resourceAssuredWorkloadsWorkloadDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAssuredWorkloadsWorkloadImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
		),

		Schema: map[string]*schema.Schema{
			"compliance_regime": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Immutable. Compliance Regime associated with this workload. Possible values: COMPLIANCE_REGIME_UNSPECIFIED, IL4, CJIS, FEDRAMP_HIGH, FEDRAMP_MODERATE, US_REGIONAL_ACCESS, HIPAA, HITRUST, EU_REGIONS_AND_SUPPORT, CA_REGIONS_AND_SUPPORT, ITAR, AU_REGIONS_AND_US_SUPPORT, ASSURED_WORKLOADS_FOR_PARTNERS, ISR_REGIONS, ISR_REGIONS_AND_SUPPORT, CA_PROTECTED_B, IL5, IL2, JP_REGIONS_AND_SUPPORT, KSA_REGIONS_AND_SUPPORT_WITH_SOVEREIGNTY_CONTROLS, REGIONAL_CONTROLS, HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS, HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS_WITH_US_SUPPORT",
			},

			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The user-assigned display name of the Workload. When present it must be between 4 to 30 characters. Allowed characters are: lowercase and uppercase letters, numbers, hyphen, and spaces. Example: My Workload",
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"organization": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The organization for the resource",
			},

			"billing_account": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. Input only. The billing account used for the resources which are direct children of workload. This billing account is initially associated with the resources created as part of Workload creation. After the initial creation of these resources, the customer can change the assigned billing account. The resource name has the form `billingAccounts/{billing_account_id}`. For example, `billingAccounts/012345-567890-ABCDEF`.",
			},

			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.",
			},

			"enable_sovereign_controls": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Indicates the sovereignty status of the given workload. Currently meant to be used by Europe/Canada customers.",
			},

			"kms_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "**DEPRECATED** Input only. Settings used to create a CMEK crypto key. When set, a project with a KMS CMEK key is provisioned. This field is deprecated as of Feb 28, 2022. In order to create a Keyring, callers should specify, ENCRYPTION_KEYS_PROJECT or KEYRING in ResourceSettings.resource_type field.",
				MaxItems:    1,
				Elem:        AssuredWorkloadsWorkloadKmsSettingsSchema(),
			},

			"partner": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Partner regime associated with this workload. Possible values: PARTNER_UNSPECIFIED, LOCAL_CONTROLS_BY_S3NS, SOVEREIGN_CONTROLS_BY_T_SYSTEMS, SOVEREIGN_CONTROLS_BY_SIA_MINSAIT, SOVEREIGN_CONTROLS_BY_PSN, SOVEREIGN_CONTROLS_BY_CNTXT, SOVEREIGN_CONTROLS_BY_CNTXT_NO_EKM",
			},

			"partner_permissions": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Permissions granted to the AW Partner SA account for the customer workload",
				MaxItems:    1,
				Elem:        AssuredWorkloadsWorkloadPartnerPermissionsSchema(),
			},

			"partner_services_billing_account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Input only. Billing account necessary for purchasing services from Sovereign Partners. This field is required for creating SIA/PSN/CNTXT partner workloads. The caller should have 'billing.resourceAssociations.create' IAM permission on this billing-account. The format of this string is billingAccounts/AAAAAA-BBBBBB-CCCCCC.",
			},

			"provisioned_resources_parent": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Input only. The parent resource for the resources managed by this Assured Workload. May be either empty or a folder resource which is a child of the Workload parent. If not specified all resources are created under the parent organization. Format: folders/{folder_id}",
			},

			"resource_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Input only. Resource properties that are used to customize workload resources. These properties (such as custom project id) will be used to create workload resources if possible. This field is optional.",
				Elem:        AssuredWorkloadsWorkloadResourceSettingsSchema(),
			},

			"violation_notifications_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Indicates whether the e-mail notification for a violation is enabled for a workload. This value will be by default True, and if not present will be considered as true. This should only be updated via updateWorkload call. Any Changes to this field during the createWorkload call will not be honored. This will always be true while creating the workload.",
			},

			"workload_options": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Used to specify certain options for a workload during workload creation - currently only supporting KAT Optionality for Regional Controls workloads.",
				MaxItems:    1,
				Elem:        AssuredWorkloadsWorkloadWorkloadOptionsSchema(),
			},

			"compliance_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Count of active Violations in the Workload.",
				Elem:        AssuredWorkloadsWorkloadComplianceStatusSchema(),
			},

			"compliant_but_disallowed_services": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Urls for services which are compliant for this Assured Workload, but which are currently disallowed by the ResourceUsageRestriction org policy. Invoke workloads.restrictAllowedResources endpoint to allow your project developers to use these services in their environment.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Immutable. The Workload creation timestamp.",
			},

			"ekm_provisioning_response": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Optional. Represents the Ekm Provisioning State of the given workload.",
				Elem:        AssuredWorkloadsWorkloadEkmProvisioningResponseSchema(),
			},

			"kaj_enrollment_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Represents the KAJ enrollment state of the given workload. Possible values: KAJ_ENROLLMENT_STATE_UNSPECIFIED, KAJ_ENROLLMENT_STATE_PENDING, KAJ_ENROLLMENT_STATE_COMPLETE",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Labels applied to the workload.\n\n**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.\nPlease refer to the field `effective_labels` for all of the labels present on the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The resource name of the workload.",
			},

			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The resources associated with this workload. These resources will be created when creating the workload. If any of the projects already exist, the workload creation will fail. Always read only.",
				Elem:        AssuredWorkloadsWorkloadResourcesSchema(),
			},

			"saa_enrollment_response": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Represents the SAA enrollment response of the given workload. SAA enrollment response is queried during workloads.get call. In failure cases, user friendly error message is shown in SAA details page.",
				Elem:        AssuredWorkloadsWorkloadSaaEnrollmentResponseSchema(),
			},

			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The combination of labels configured directly on the resource and default labels configured on the provider.",
			},
		},
	}
}

func AssuredWorkloadsWorkloadKmsSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"next_rotation_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Input only. Immutable. The time at which the Key Management Service will automatically create a new version of the crypto key and mark it as the primary.",
			},

			"rotation_period": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Input only. Immutable. will be advanced by this period when the Key Management Service automatically rotates a key. Must be at least 24 hours and at most 876,000 hours.",
			},
		},
	}
}

func AssuredWorkloadsWorkloadPartnerPermissionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"assured_workloads_monitoring": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Allow partner to view violation alerts.",
			},

			"data_logs_viewer": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Allow the partner to view inspectability logs and monitoring violations.",
			},

			"service_access_approver": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Allow partner to view access approval logs.",
			},
		},
	}
}

func AssuredWorkloadsWorkloadResourceSettingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "User-assigned resource display name. If not empty it will be used to create a resource with the specified name.",
			},

			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Resource identifier. For a project this represents projectId. If the project is already taken, the workload creation will fail. For KeyRing, this represents the keyring_id. For a folder, don't set this value as folder_id is assigned by Google.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates the type of resource. This field should be specified to correspond the id to the right project type (CONSUMER_PROJECT or ENCRYPTION_KEYS_PROJECT) Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER",
			},
		},
	}
}

func AssuredWorkloadsWorkloadWorkloadOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kaj_enrollment_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates type of KAJ enrollment for the workload. Currently, only specifiying KEY_ACCESS_TRANSPARENCY_OFF is implemented to not enroll in KAT-level KAJ enrollment for Regional Controls workloads. Possible values: KAJ_ENROLLMENT_TYPE_UNSPECIFIED, FULL_KAJ, EKM_ONLY, KEY_ACCESS_TRANSPARENCY_OFF",
			},
		},
	}
}

func AssuredWorkloadsWorkloadComplianceStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"acknowledged_violation_count": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Number of current orgPolicy violations which are acknowledged.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},

			"active_violation_count": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Number of current orgPolicy violations which are not acknowledged.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func AssuredWorkloadsWorkloadEkmProvisioningResponseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ekm_provisioning_error_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates Ekm provisioning error if any. Possible values: EKM_PROVISIONING_ERROR_DOMAIN_UNSPECIFIED, UNSPECIFIED_ERROR, GOOGLE_SERVER_ERROR, EXTERNAL_USER_ERROR, EXTERNAL_PARTNER_ERROR, TIMEOUT_ERROR",
			},

			"ekm_provisioning_error_mapping": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Detailed error message if Ekm provisioning fails Possible values: EKM_PROVISIONING_ERROR_MAPPING_UNSPECIFIED, INVALID_SERVICE_ACCOUNT, MISSING_METRICS_SCOPE_ADMIN_PERMISSION, MISSING_EKM_CONNECTION_ADMIN_PERMISSION",
			},

			"ekm_provisioning_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates Ekm enrollment Provisioning of a given workload. Possible values: EKM_PROVISIONING_STATE_UNSPECIFIED, EKM_PROVISIONING_STATE_PENDING, EKM_PROVISIONING_STATE_FAILED, EKM_PROVISIONING_STATE_COMPLETED",
			},
		},
	}
}

func AssuredWorkloadsWorkloadResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Resource identifier. For a project this represents project_number.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of resource. Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER",
			},
		},
	}
}

func AssuredWorkloadsWorkloadSaaEnrollmentResponseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"setup_errors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Indicates SAA enrollment setup error if any.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"setup_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates SAA enrollment status of a given workload. Possible values: SETUP_STATE_UNSPECIFIED, STATUS_PENDING, STATUS_COMPLETE",
			},
		},
	}
}

func resourceAssuredWorkloadsWorkloadCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		ComplianceRegime:              assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                   dcl.String(d.Get("display_name").(string)),
		Location:                      dcl.String(d.Get("location").(string)),
		Organization:                  dcl.String(d.Get("organization").(string)),
		BillingAccount:                dcl.String(d.Get("billing_account").(string)),
		Labels:                        tpgresource.CheckStringMap(d.Get("effective_labels")),
		EnableSovereignControls:       dcl.Bool(d.Get("enable_sovereign_controls").(bool)),
		KmsSettings:                   expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Partner:                       assuredworkloads.WorkloadPartnerEnumRef(d.Get("partner").(string)),
		PartnerPermissions:            expandAssuredWorkloadsWorkloadPartnerPermissions(d.Get("partner_permissions")),
		PartnerServicesBillingAccount: dcl.String(d.Get("partner_services_billing_account").(string)),
		ProvisionedResourcesParent:    dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:              expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
		ViolationNotificationsEnabled: dcl.Bool(d.Get("violation_notifications_enabled").(bool)),
		WorkloadOptions:               expandAssuredWorkloadsWorkloadWorkloadOptions(d.Get("workload_options")),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyWorkload(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Workload: %s", err)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	// ID has a server-generated value, set again after creation.

	id, err = res.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Workload %q: %#v", d.Id(), res)

	return resourceAssuredWorkloadsWorkloadRead(d, meta)
}

func resourceAssuredWorkloadsWorkloadRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		ComplianceRegime:              assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                   dcl.String(d.Get("display_name").(string)),
		Location:                      dcl.String(d.Get("location").(string)),
		Organization:                  dcl.String(d.Get("organization").(string)),
		BillingAccount:                dcl.String(d.Get("billing_account").(string)),
		Labels:                        tpgresource.CheckStringMap(d.Get("effective_labels")),
		EnableSovereignControls:       dcl.Bool(d.Get("enable_sovereign_controls").(bool)),
		KmsSettings:                   expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Partner:                       assuredworkloads.WorkloadPartnerEnumRef(d.Get("partner").(string)),
		PartnerPermissions:            expandAssuredWorkloadsWorkloadPartnerPermissions(d.Get("partner_permissions")),
		PartnerServicesBillingAccount: dcl.String(d.Get("partner_services_billing_account").(string)),
		ProvisionedResourcesParent:    dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:              expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
		ViolationNotificationsEnabled: dcl.Bool(d.Get("violation_notifications_enabled").(bool)),
		WorkloadOptions:               expandAssuredWorkloadsWorkloadWorkloadOptions(d.Get("workload_options")),
		Name:                          dcl.StringOrNil(d.Get("name").(string)),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetWorkload(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("AssuredWorkloadsWorkload %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("compliance_regime", res.ComplianceRegime); err != nil {
		return fmt.Errorf("error setting compliance_regime in state: %s", err)
	}
	if err = d.Set("display_name", res.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("organization", res.Organization); err != nil {
		return fmt.Errorf("error setting organization in state: %s", err)
	}
	if err = d.Set("billing_account", res.BillingAccount); err != nil {
		return fmt.Errorf("error setting billing_account in state: %s", err)
	}
	if err = d.Set("effective_labels", res.Labels); err != nil {
		return fmt.Errorf("error setting effective_labels in state: %s", err)
	}
	if err = d.Set("enable_sovereign_controls", res.EnableSovereignControls); err != nil {
		return fmt.Errorf("error setting enable_sovereign_controls in state: %s", err)
	}
	if err = d.Set("kms_settings", flattenAssuredWorkloadsWorkloadKmsSettings(res.KmsSettings)); err != nil {
		return fmt.Errorf("error setting kms_settings in state: %s", err)
	}
	if err = d.Set("partner", res.Partner); err != nil {
		return fmt.Errorf("error setting partner in state: %s", err)
	}
	if err = d.Set("partner_permissions", flattenAssuredWorkloadsWorkloadPartnerPermissions(res.PartnerPermissions)); err != nil {
		return fmt.Errorf("error setting partner_permissions in state: %s", err)
	}
	if err = d.Set("partner_services_billing_account", res.PartnerServicesBillingAccount); err != nil {
		return fmt.Errorf("error setting partner_services_billing_account in state: %s", err)
	}
	if err = d.Set("provisioned_resources_parent", res.ProvisionedResourcesParent); err != nil {
		return fmt.Errorf("error setting provisioned_resources_parent in state: %s", err)
	}
	if err = d.Set("resource_settings", flattenAssuredWorkloadsWorkloadResourceSettingsArray(res.ResourceSettings)); err != nil {
		return fmt.Errorf("error setting resource_settings in state: %s", err)
	}
	if err = d.Set("violation_notifications_enabled", res.ViolationNotificationsEnabled); err != nil {
		return fmt.Errorf("error setting violation_notifications_enabled in state: %s", err)
	}
	if err = d.Set("workload_options", flattenAssuredWorkloadsWorkloadWorkloadOptions(res.WorkloadOptions)); err != nil {
		return fmt.Errorf("error setting workload_options in state: %s", err)
	}
	if err = d.Set("compliance_status", flattenAssuredWorkloadsWorkloadComplianceStatus(res.ComplianceStatus)); err != nil {
		return fmt.Errorf("error setting compliance_status in state: %s", err)
	}
	if err = d.Set("compliant_but_disallowed_services", res.CompliantButDisallowedServices); err != nil {
		return fmt.Errorf("error setting compliant_but_disallowed_services in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("ekm_provisioning_response", flattenAssuredWorkloadsWorkloadEkmProvisioningResponse(res.EkmProvisioningResponse)); err != nil {
		return fmt.Errorf("error setting ekm_provisioning_response in state: %s", err)
	}
	if err = d.Set("kaj_enrollment_state", res.KajEnrollmentState); err != nil {
		return fmt.Errorf("error setting kaj_enrollment_state in state: %s", err)
	}
	if err = d.Set("labels", flattenAssuredWorkloadsWorkloadLabels(res.Labels, d)); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("resources", flattenAssuredWorkloadsWorkloadResourcesArray(res.Resources)); err != nil {
		return fmt.Errorf("error setting resources in state: %s", err)
	}
	if err = d.Set("saa_enrollment_response", flattenAssuredWorkloadsWorkloadSaaEnrollmentResponse(res.SaaEnrollmentResponse)); err != nil {
		return fmt.Errorf("error setting saa_enrollment_response in state: %s", err)
	}
	if err = d.Set("terraform_labels", flattenAssuredWorkloadsWorkloadTerraformLabels(res.Labels, d)); err != nil {
		return fmt.Errorf("error setting terraform_labels in state: %s", err)
	}

	return nil
}
func resourceAssuredWorkloadsWorkloadUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		ComplianceRegime:              assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                   dcl.String(d.Get("display_name").(string)),
		Location:                      dcl.String(d.Get("location").(string)),
		Organization:                  dcl.String(d.Get("organization").(string)),
		BillingAccount:                dcl.String(d.Get("billing_account").(string)),
		Labels:                        tpgresource.CheckStringMap(d.Get("effective_labels")),
		EnableSovereignControls:       dcl.Bool(d.Get("enable_sovereign_controls").(bool)),
		KmsSettings:                   expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Partner:                       assuredworkloads.WorkloadPartnerEnumRef(d.Get("partner").(string)),
		PartnerPermissions:            expandAssuredWorkloadsWorkloadPartnerPermissions(d.Get("partner_permissions")),
		PartnerServicesBillingAccount: dcl.String(d.Get("partner_services_billing_account").(string)),
		ProvisionedResourcesParent:    dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:              expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
		ViolationNotificationsEnabled: dcl.Bool(d.Get("violation_notifications_enabled").(bool)),
		WorkloadOptions:               expandAssuredWorkloadsWorkloadWorkloadOptions(d.Get("workload_options")),
		Name:                          dcl.StringOrNil(d.Get("name").(string)),
	}
	// Construct state hint from old values
	old := &assuredworkloads.Workload{
		ComplianceRegime:              assuredworkloads.WorkloadComplianceRegimeEnumRef(tpgdclresource.OldValue(d.GetChange("compliance_regime")).(string)),
		DisplayName:                   dcl.String(tpgdclresource.OldValue(d.GetChange("display_name")).(string)),
		Location:                      dcl.String(tpgdclresource.OldValue(d.GetChange("location")).(string)),
		Organization:                  dcl.String(tpgdclresource.OldValue(d.GetChange("organization")).(string)),
		BillingAccount:                dcl.String(tpgdclresource.OldValue(d.GetChange("billing_account")).(string)),
		Labels:                        tpgresource.CheckStringMap(tpgdclresource.OldValue(d.GetChange("effective_labels"))),
		EnableSovereignControls:       dcl.Bool(tpgdclresource.OldValue(d.GetChange("enable_sovereign_controls")).(bool)),
		KmsSettings:                   expandAssuredWorkloadsWorkloadKmsSettings(tpgdclresource.OldValue(d.GetChange("kms_settings"))),
		Partner:                       assuredworkloads.WorkloadPartnerEnumRef(tpgdclresource.OldValue(d.GetChange("partner")).(string)),
		PartnerPermissions:            expandAssuredWorkloadsWorkloadPartnerPermissions(tpgdclresource.OldValue(d.GetChange("partner_permissions"))),
		PartnerServicesBillingAccount: dcl.String(tpgdclresource.OldValue(d.GetChange("partner_services_billing_account")).(string)),
		ProvisionedResourcesParent:    dcl.String(tpgdclresource.OldValue(d.GetChange("provisioned_resources_parent")).(string)),
		ResourceSettings:              expandAssuredWorkloadsWorkloadResourceSettingsArray(tpgdclresource.OldValue(d.GetChange("resource_settings"))),
		ViolationNotificationsEnabled: dcl.Bool(tpgdclresource.OldValue(d.GetChange("violation_notifications_enabled")).(bool)),
		WorkloadOptions:               expandAssuredWorkloadsWorkloadWorkloadOptions(tpgdclresource.OldValue(d.GetChange("workload_options"))),
		Name:                          dcl.StringOrNil(tpgdclresource.OldValue(d.GetChange("name")).(string)),
	}
	directive := tpgdclresource.UpdateDirective
	directive = append(directive, dcl.WithStateHint(old))
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyWorkload(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Workload: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Workload %q: %#v", d.Id(), res)

	return resourceAssuredWorkloadsWorkloadRead(d, meta)
}

func resourceAssuredWorkloadsWorkloadDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	obj := &assuredworkloads.Workload{
		ComplianceRegime:              assuredworkloads.WorkloadComplianceRegimeEnumRef(d.Get("compliance_regime").(string)),
		DisplayName:                   dcl.String(d.Get("display_name").(string)),
		Location:                      dcl.String(d.Get("location").(string)),
		Organization:                  dcl.String(d.Get("organization").(string)),
		BillingAccount:                dcl.String(d.Get("billing_account").(string)),
		Labels:                        tpgresource.CheckStringMap(d.Get("effective_labels")),
		EnableSovereignControls:       dcl.Bool(d.Get("enable_sovereign_controls").(bool)),
		KmsSettings:                   expandAssuredWorkloadsWorkloadKmsSettings(d.Get("kms_settings")),
		Partner:                       assuredworkloads.WorkloadPartnerEnumRef(d.Get("partner").(string)),
		PartnerPermissions:            expandAssuredWorkloadsWorkloadPartnerPermissions(d.Get("partner_permissions")),
		PartnerServicesBillingAccount: dcl.String(d.Get("partner_services_billing_account").(string)),
		ProvisionedResourcesParent:    dcl.String(d.Get("provisioned_resources_parent").(string)),
		ResourceSettings:              expandAssuredWorkloadsWorkloadResourceSettingsArray(d.Get("resource_settings")),
		ViolationNotificationsEnabled: dcl.Bool(d.Get("violation_notifications_enabled").(bool)),
		WorkloadOptions:               expandAssuredWorkloadsWorkloadWorkloadOptions(d.Get("workload_options")),
		Name:                          dcl.StringOrNil(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Deleting Workload %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLAssuredWorkloadsClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteWorkload(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Workload: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Workload %q", d.Id())
	return nil
}

func resourceAssuredWorkloadsWorkloadImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"organizations/(?P<organization>[^/]+)/locations/(?P<location>[^/]+)/workloads/(?P<name>[^/]+)",
		"(?P<organization>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "organizations/{{organization}}/locations/{{location}}/workloads/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandAssuredWorkloadsWorkloadKmsSettings(o interface{}) *assuredworkloads.WorkloadKmsSettings {
	if o == nil {
		return assuredworkloads.EmptyWorkloadKmsSettings
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return assuredworkloads.EmptyWorkloadKmsSettings
	}
	obj := objArr[0].(map[string]interface{})
	return &assuredworkloads.WorkloadKmsSettings{
		NextRotationTime: dcl.String(obj["next_rotation_time"].(string)),
		RotationPeriod:   dcl.String(obj["rotation_period"].(string)),
	}
}

func flattenAssuredWorkloadsWorkloadKmsSettings(obj *assuredworkloads.WorkloadKmsSettings) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"next_rotation_time": obj.NextRotationTime,
		"rotation_period":    obj.RotationPeriod,
	}

	return []interface{}{transformed}

}

func expandAssuredWorkloadsWorkloadPartnerPermissions(o interface{}) *assuredworkloads.WorkloadPartnerPermissions {
	if o == nil {
		return assuredworkloads.EmptyWorkloadPartnerPermissions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return assuredworkloads.EmptyWorkloadPartnerPermissions
	}
	obj := objArr[0].(map[string]interface{})
	return &assuredworkloads.WorkloadPartnerPermissions{
		AssuredWorkloadsMonitoring: dcl.Bool(obj["assured_workloads_monitoring"].(bool)),
		DataLogsViewer:             dcl.Bool(obj["data_logs_viewer"].(bool)),
		ServiceAccessApprover:      dcl.Bool(obj["service_access_approver"].(bool)),
	}
}

func flattenAssuredWorkloadsWorkloadPartnerPermissions(obj *assuredworkloads.WorkloadPartnerPermissions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"assured_workloads_monitoring": obj.AssuredWorkloadsMonitoring,
		"data_logs_viewer":             obj.DataLogsViewer,
		"service_access_approver":      obj.ServiceAccessApprover,
	}

	return []interface{}{transformed}

}
func expandAssuredWorkloadsWorkloadResourceSettingsArray(o interface{}) []assuredworkloads.WorkloadResourceSettings {
	if o == nil {
		return make([]assuredworkloads.WorkloadResourceSettings, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]assuredworkloads.WorkloadResourceSettings, 0)
	}

	items := make([]assuredworkloads.WorkloadResourceSettings, 0, len(objs))
	for _, item := range objs {
		i := expandAssuredWorkloadsWorkloadResourceSettings(item)
		items = append(items, *i)
	}

	return items
}

func expandAssuredWorkloadsWorkloadResourceSettings(o interface{}) *assuredworkloads.WorkloadResourceSettings {
	if o == nil {
		return assuredworkloads.EmptyWorkloadResourceSettings
	}

	obj := o.(map[string]interface{})
	return &assuredworkloads.WorkloadResourceSettings{
		DisplayName:  dcl.String(obj["display_name"].(string)),
		ResourceId:   dcl.String(obj["resource_id"].(string)),
		ResourceType: assuredworkloads.WorkloadResourceSettingsResourceTypeEnumRef(obj["resource_type"].(string)),
	}
}

func flattenAssuredWorkloadsWorkloadResourceSettingsArray(objs []assuredworkloads.WorkloadResourceSettings) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenAssuredWorkloadsWorkloadResourceSettings(&item)
		items = append(items, i)
	}

	return items
}

func flattenAssuredWorkloadsWorkloadResourceSettings(obj *assuredworkloads.WorkloadResourceSettings) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"display_name":  obj.DisplayName,
		"resource_id":   obj.ResourceId,
		"resource_type": obj.ResourceType,
	}

	return transformed

}

func expandAssuredWorkloadsWorkloadWorkloadOptions(o interface{}) *assuredworkloads.WorkloadWorkloadOptions {
	if o == nil {
		return assuredworkloads.EmptyWorkloadWorkloadOptions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return assuredworkloads.EmptyWorkloadWorkloadOptions
	}
	obj := objArr[0].(map[string]interface{})
	return &assuredworkloads.WorkloadWorkloadOptions{
		KajEnrollmentType: assuredworkloads.WorkloadWorkloadOptionsKajEnrollmentTypeEnumRef(obj["kaj_enrollment_type"].(string)),
	}
}

func flattenAssuredWorkloadsWorkloadWorkloadOptions(obj *assuredworkloads.WorkloadWorkloadOptions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"kaj_enrollment_type": obj.KajEnrollmentType,
	}

	return []interface{}{transformed}

}

func flattenAssuredWorkloadsWorkloadComplianceStatus(obj *assuredworkloads.WorkloadComplianceStatus) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"acknowledged_violation_count": obj.AcknowledgedViolationCount,
		"active_violation_count":       obj.ActiveViolationCount,
	}

	return []interface{}{transformed}

}

func flattenAssuredWorkloadsWorkloadEkmProvisioningResponse(obj *assuredworkloads.WorkloadEkmProvisioningResponse) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"ekm_provisioning_error_domain":  obj.EkmProvisioningErrorDomain,
		"ekm_provisioning_error_mapping": obj.EkmProvisioningErrorMapping,
		"ekm_provisioning_state":         obj.EkmProvisioningState,
	}

	return []interface{}{transformed}

}

func flattenAssuredWorkloadsWorkloadResourcesArray(objs []assuredworkloads.WorkloadResources) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenAssuredWorkloadsWorkloadResources(&item)
		items = append(items, i)
	}

	return items
}

func flattenAssuredWorkloadsWorkloadResources(obj *assuredworkloads.WorkloadResources) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"resource_id":   obj.ResourceId,
		"resource_type": obj.ResourceType,
	}

	return transformed

}

func flattenAssuredWorkloadsWorkloadSaaEnrollmentResponse(obj *assuredworkloads.WorkloadSaaEnrollmentResponse) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"setup_errors": flattenAssuredWorkloadsWorkloadSaaEnrollmentResponseSetupErrorsArray(obj.SetupErrors),
		"setup_status": obj.SetupStatus,
	}

	return []interface{}{transformed}

}

func flattenAssuredWorkloadsWorkloadLabels(v map[string]string, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	if l, ok := d.Get("labels").(map[string]interface{}); ok {
		for k, _ := range l {
			transformed[k] = v[k]
		}
	}

	return transformed
}

func flattenAssuredWorkloadsWorkloadTerraformLabels(v map[string]string, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	if l, ok := d.Get("terraform_labels").(map[string]interface{}); ok {
		for k, _ := range l {
			transformed[k] = v[k]
		}
	}

	return transformed
}

func flattenAssuredWorkloadsWorkloadSaaEnrollmentResponseSetupErrorsArray(obj []assuredworkloads.WorkloadSaaEnrollmentResponseSetupErrorsEnum) interface{} {
	if obj == nil {
		return nil
	}
	items := []string{}
	for _, item := range obj {
		items = append(items, string(item))
	}
	return items
}
