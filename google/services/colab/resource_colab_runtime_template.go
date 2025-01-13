// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package colab

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceColabRuntimeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceColabRuntimeTemplateCreate,
		Read:   resourceColabRuntimeTemplateRead,
		Update: resourceColabRuntimeTemplateUpdate,
		Delete: resourceColabRuntimeTemplateDelete,

		Importer: &schema.ResourceImporter{
			State: resourceColabRuntimeTemplateImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Required. The display name of the Runtime Template.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The location for the resource: https://cloud.google.com/colab/docs/locations`,
			},
			"data_persistent_disk_spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The configuration for the data disk of the runtime.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size_gb": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: `The disk size of the runtime in GB. If specified, the diskType must also be specified. The minimum size is 10GB and the maximum is 65536GB.`,
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: `The type of the persistent disk.`,
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The description of the Runtime Template.`,
			},
			"encryption_spec": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `Customer-managed encryption key spec for the notebook runtime.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The Cloud KMS encryption key (customer-managed encryption key) used to protect the runtime.`,
						},
					},
				},
			},
			"euc_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `EUC configuration of the NotebookRuntimeTemplate.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"euc_disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: `Disable end user credential access for the runtime.`,
						},
					},
				},
			},
			"idle_shutdown_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `Notebook Idle Shutdown configuration for the runtime.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idle_timeout": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: `The duration after which the runtime is automatically shut down. An input of 0s disables the idle shutdown feature, and a valid range is [10m, 24h].`,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Labels to identify and group the runtime template.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"machine_spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `'The machine configuration of the runtime.'`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerator_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: `The number of accelerators used by the runtime.`,
						},
						"accelerator_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The type of hardware accelerator used by the runtime. If specified, acceleratorCount must also be specified.`,
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: `The Compute Engine machine type selected for the runtime.`,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The resource name of the Runtime Template`,
			},
			"network_spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The network configuration for the runtime.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_internet_access": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: `Enable public internet access for the runtime.`,
						},
						"network": {
							Type:             schema.TypeString,
							Computed:         true,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
							Description:      `The name of the VPC that this runtime is in.`,
						},
						"subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
							Description:      `The name of the subnetwork that this runtime is in.`,
						},
					},
				},
			},
			"network_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `Applies the given Compute Engine tags to the runtime.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"shielded_vm_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `Runtime Shielded VM spec.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_secure_boot": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: `Enables secure boot for the runtime.`,
						},
					},
				},
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				ForceNew:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceColabRuntimeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandColabRuntimeTemplateName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	displayNameProp, err := expandColabRuntimeTemplateDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandColabRuntimeTemplateDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	machineSpecProp, err := expandColabRuntimeTemplateMachineSpec(d.Get("machine_spec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("machine_spec"); !tpgresource.IsEmptyValue(reflect.ValueOf(machineSpecProp)) && (ok || !reflect.DeepEqual(v, machineSpecProp)) {
		obj["machineSpec"] = machineSpecProp
	}
	dataPersistentDiskSpecProp, err := expandColabRuntimeTemplateDataPersistentDiskSpec(d.Get("data_persistent_disk_spec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("data_persistent_disk_spec"); !tpgresource.IsEmptyValue(reflect.ValueOf(dataPersistentDiskSpecProp)) && (ok || !reflect.DeepEqual(v, dataPersistentDiskSpecProp)) {
		obj["dataPersistentDiskSpec"] = dataPersistentDiskSpecProp
	}
	networkSpecProp, err := expandColabRuntimeTemplateNetworkSpec(d.Get("network_spec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("network_spec"); !tpgresource.IsEmptyValue(reflect.ValueOf(networkSpecProp)) && (ok || !reflect.DeepEqual(v, networkSpecProp)) {
		obj["networkSpec"] = networkSpecProp
	}
	idleShutdownConfigProp, err := expandColabRuntimeTemplateIdleShutdownConfig(d.Get("idle_shutdown_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("idle_shutdown_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(idleShutdownConfigProp)) && (ok || !reflect.DeepEqual(v, idleShutdownConfigProp)) {
		obj["idleShutdownConfig"] = idleShutdownConfigProp
	}
	eucConfigProp, err := expandColabRuntimeTemplateEucConfig(d.Get("euc_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("euc_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(eucConfigProp)) && (ok || !reflect.DeepEqual(v, eucConfigProp)) {
		obj["eucConfig"] = eucConfigProp
	}
	shieldedVmConfigProp, err := expandColabRuntimeTemplateShieldedVmConfig(d.Get("shielded_vm_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("shielded_vm_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(shieldedVmConfigProp)) && (ok || !reflect.DeepEqual(v, shieldedVmConfigProp)) {
		obj["shieldedVmConfig"] = shieldedVmConfigProp
	}
	networkTagsProp, err := expandColabRuntimeTemplateNetworkTags(d.Get("network_tags"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("network_tags"); !tpgresource.IsEmptyValue(reflect.ValueOf(networkTagsProp)) && (ok || !reflect.DeepEqual(v, networkTagsProp)) {
		obj["networkTags"] = networkTagsProp
	}
	encryptionSpecProp, err := expandColabRuntimeTemplateEncryptionSpec(d.Get("encryption_spec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("encryption_spec"); !tpgresource.IsEmptyValue(reflect.ValueOf(encryptionSpecProp)) && (ok || !reflect.DeepEqual(v, encryptionSpecProp)) {
		obj["encryptionSpec"] = encryptionSpecProp
	}
	labelsProp, err := expandColabRuntimeTemplateEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ColabBasePath}}projects/{{project}}/locations/{{location}}/notebookRuntimeTemplates?notebook_runtime_template_id={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new RuntimeTemplate: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RuntimeTemplate: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating RuntimeTemplate: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/notebookRuntimeTemplates/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = ColabOperationWaitTime(
		config, res, project, "Creating RuntimeTemplate", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create RuntimeTemplate: %s", err)
	}

	// The operation for this resource contains the generated name that we need
	// in order to perform a READ. We need to access the object inside of it as
	// a map[string]interface, so let's do that.

	resp := res["response"].(map[string]interface{})
	name := tpgresource.GetResourceNameFromSelfLink(resp["name"].(string))
	log.Printf("[DEBUG] Setting resource name, id to %s", name)
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	log.Printf("[DEBUG] Finished creating RuntimeTemplate %q: %#v", d.Id(), res)

	return resourceColabRuntimeTemplateRead(d, meta)
}

func resourceColabRuntimeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ColabBasePath}}projects/{{project}}/locations/{{location}}/notebookRuntimeTemplates/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RuntimeTemplate: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ColabRuntimeTemplate %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}

	if err := d.Set("name", flattenColabRuntimeTemplateName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("display_name", flattenColabRuntimeTemplateDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("description", flattenColabRuntimeTemplateDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("machine_spec", flattenColabRuntimeTemplateMachineSpec(res["machineSpec"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("data_persistent_disk_spec", flattenColabRuntimeTemplateDataPersistentDiskSpec(res["dataPersistentDiskSpec"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("network_spec", flattenColabRuntimeTemplateNetworkSpec(res["networkSpec"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("labels", flattenColabRuntimeTemplateLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("idle_shutdown_config", flattenColabRuntimeTemplateIdleShutdownConfig(res["idleShutdownConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("euc_config", flattenColabRuntimeTemplateEucConfig(res["eucConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("shielded_vm_config", flattenColabRuntimeTemplateShieldedVmConfig(res["shieldedVmConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("network_tags", flattenColabRuntimeTemplateNetworkTags(res["networkTags"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("encryption_spec", flattenColabRuntimeTemplateEncryptionSpec(res["encryptionSpec"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("terraform_labels", flattenColabRuntimeTemplateTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}
	if err := d.Set("effective_labels", flattenColabRuntimeTemplateEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuntimeTemplate: %s", err)
	}

	return nil
}

func resourceColabRuntimeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	// Only the root field "labels" and "terraform_labels" are mutable
	return resourceColabRuntimeTemplateRead(d, meta)
}

func resourceColabRuntimeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RuntimeTemplate: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{ColabBasePath}}projects/{{project}}/locations/{{location}}/notebookRuntimeTemplates/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting RuntimeTemplate %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "RuntimeTemplate")
	}

	err = ColabOperationWaitTime(
		config, res, project, "Deleting RuntimeTemplate", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting RuntimeTemplate %q: %#v", d.Id(), res)
	return nil
}

func resourceColabRuntimeTemplateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/notebookRuntimeTemplates/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/notebookRuntimeTemplates/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenColabRuntimeTemplateName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenColabRuntimeTemplateDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateMachineSpec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["machine_type"] =
		flattenColabRuntimeTemplateMachineSpecMachineType(original["machineType"], d, config)
	transformed["accelerator_type"] =
		flattenColabRuntimeTemplateMachineSpecAcceleratorType(original["acceleratorType"], d, config)
	transformed["accelerator_count"] =
		flattenColabRuntimeTemplateMachineSpecAcceleratorCount(original["acceleratorCount"], d, config)
	return []interface{}{transformed}
}
func flattenColabRuntimeTemplateMachineSpecMachineType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateMachineSpecAcceleratorType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateMachineSpecAcceleratorCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenColabRuntimeTemplateDataPersistentDiskSpec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["disk_type"] =
		flattenColabRuntimeTemplateDataPersistentDiskSpecDiskType(original["diskType"], d, config)
	transformed["disk_size_gb"] =
		flattenColabRuntimeTemplateDataPersistentDiskSpecDiskSizeGb(original["diskSizeGb"], d, config)
	return []interface{}{transformed}
}
func flattenColabRuntimeTemplateDataPersistentDiskSpecDiskType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateDataPersistentDiskSpecDiskSizeGb(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateNetworkSpec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["enable_internet_access"] =
		flattenColabRuntimeTemplateNetworkSpecEnableInternetAccess(original["enableInternetAccess"], d, config)
	transformed["network"] =
		flattenColabRuntimeTemplateNetworkSpecNetwork(original["network"], d, config)
	transformed["subnetwork"] =
		flattenColabRuntimeTemplateNetworkSpecSubnetwork(original["subnetwork"], d, config)
	return []interface{}{transformed}
}
func flattenColabRuntimeTemplateNetworkSpecEnableInternetAccess(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateNetworkSpecNetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateNetworkSpecSubnetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenColabRuntimeTemplateIdleShutdownConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["idle_timeout"] =
		flattenColabRuntimeTemplateIdleShutdownConfigIdleTimeout(original["idleTimeout"], d, config)
	return []interface{}{transformed}
}
func flattenColabRuntimeTemplateIdleShutdownConfigIdleTimeout(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateEucConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["euc_disabled"] =
		flattenColabRuntimeTemplateEucConfigEucDisabled(original["eucDisabled"], d, config)
	return []interface{}{transformed}
}
func flattenColabRuntimeTemplateEucConfigEucDisabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateShieldedVmConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["enable_secure_boot"] =
		flattenColabRuntimeTemplateShieldedVmConfigEnableSecureBoot(original["enableSecureBoot"], d, config)
	return []interface{}{transformed}
}
func flattenColabRuntimeTemplateShieldedVmConfigEnableSecureBoot(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateNetworkTags(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateEncryptionSpec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["kms_key_name"] =
		flattenColabRuntimeTemplateEncryptionSpecKmsKeyName(original["kmsKeyName"], d, config)
	return []interface{}{transformed}
}
func flattenColabRuntimeTemplateEncryptionSpecKmsKeyName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenColabRuntimeTemplateTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("terraform_labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenColabRuntimeTemplateEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandColabRuntimeTemplateName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateMachineSpec(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMachineType, err := expandColabRuntimeTemplateMachineSpecMachineType(original["machine_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMachineType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["machineType"] = transformedMachineType
	}

	transformedAcceleratorType, err := expandColabRuntimeTemplateMachineSpecAcceleratorType(original["accelerator_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAcceleratorType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["acceleratorType"] = transformedAcceleratorType
	}

	transformedAcceleratorCount, err := expandColabRuntimeTemplateMachineSpecAcceleratorCount(original["accelerator_count"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAcceleratorCount); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["acceleratorCount"] = transformedAcceleratorCount
	}

	return transformed, nil
}

func expandColabRuntimeTemplateMachineSpecMachineType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateMachineSpecAcceleratorType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateMachineSpecAcceleratorCount(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateDataPersistentDiskSpec(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDiskType, err := expandColabRuntimeTemplateDataPersistentDiskSpecDiskType(original["disk_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskType"] = transformedDiskType
	}

	transformedDiskSizeGb, err := expandColabRuntimeTemplateDataPersistentDiskSpecDiskSizeGb(original["disk_size_gb"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDiskSizeGb); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["diskSizeGb"] = transformedDiskSizeGb
	}

	return transformed, nil
}

func expandColabRuntimeTemplateDataPersistentDiskSpecDiskType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateDataPersistentDiskSpecDiskSizeGb(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateNetworkSpec(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnableInternetAccess, err := expandColabRuntimeTemplateNetworkSpecEnableInternetAccess(original["enable_internet_access"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnableInternetAccess); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enableInternetAccess"] = transformedEnableInternetAccess
	}

	transformedNetwork, err := expandColabRuntimeTemplateNetworkSpecNetwork(original["network"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNetwork); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["network"] = transformedNetwork
	}

	transformedSubnetwork, err := expandColabRuntimeTemplateNetworkSpecSubnetwork(original["subnetwork"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSubnetwork); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["subnetwork"] = transformedSubnetwork
	}

	return transformed, nil
}

func expandColabRuntimeTemplateNetworkSpecEnableInternetAccess(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateNetworkSpecNetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateNetworkSpecSubnetwork(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateIdleShutdownConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedIdleTimeout, err := expandColabRuntimeTemplateIdleShutdownConfigIdleTimeout(original["idle_timeout"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIdleTimeout); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["idleTimeout"] = transformedIdleTimeout
	}

	return transformed, nil
}

func expandColabRuntimeTemplateIdleShutdownConfigIdleTimeout(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateEucConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEucDisabled, err := expandColabRuntimeTemplateEucConfigEucDisabled(original["euc_disabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEucDisabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["eucDisabled"] = transformedEucDisabled
	}

	return transformed, nil
}

func expandColabRuntimeTemplateEucConfigEucDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateShieldedVmConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnableSecureBoot, err := expandColabRuntimeTemplateShieldedVmConfigEnableSecureBoot(original["enable_secure_boot"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnableSecureBoot); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enableSecureBoot"] = transformedEnableSecureBoot
	}

	return transformed, nil
}

func expandColabRuntimeTemplateShieldedVmConfigEnableSecureBoot(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateNetworkTags(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateEncryptionSpec(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedKmsKeyName, err := expandColabRuntimeTemplateEncryptionSpecKmsKeyName(original["kms_key_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKmsKeyName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["kmsKeyName"] = transformedKmsKeyName
	}

	return transformed, nil
}

func expandColabRuntimeTemplateEncryptionSpecKmsKeyName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandColabRuntimeTemplateEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
