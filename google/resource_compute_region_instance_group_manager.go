package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"fmt"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

var RegionInstanceGroupManagerBaseApiVersion = v1
var RegionInstanceGroupManagerVersionedFeatures = []Feature{Feature{Version: v0beta, Item: "auto_healing_policies"}}

func resourceComputeRegionInstanceGroupManager() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionInstanceGroupManagerCreate,
		Read:   resourceComputeRegionInstanceGroupManagerRead,
		Update: resourceComputeRegionInstanceGroupManagerUpdate,
		Delete: resourceComputeRegionInstanceGroupManagerDelete,
		Exists: resourceComputeRegionInstanceGroupManagerExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"base_instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_template": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"named_port": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_pools": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: selfLinkRelativePathHash,
			},

			"target_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"auto_healing_policies": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},

						"initial_delay_sec": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 3600),
						},
					},
				},
			},
		},
	}
}

func resourceComputeRegionInstanceGroupManagerCreate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	manager := &computeBeta.InstanceGroupManager{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		BaseInstanceName:    d.Get("base_instance_name").(string),
		InstanceTemplate:    d.Get("instance_template").(string),
		TargetSize:          int64(d.Get("target_size").(int)),
		NamedPorts:          getNamedPortsBeta(d.Get("named_port").([]interface{})),
		TargetPools:         convertStringSet(d.Get("target_pools").(*schema.Set)),
		AutoHealingPolicies: expandAutoHealingPolicies(d.Get("auto_healing_policies").([]interface{})),
		// Force send TargetSize to allow size of 0.
		ForceSendFields: []string{"TargetSize"},
	}

	var op interface{}
	switch computeApiVersion {
	case v1:
		managerV1 := &compute.InstanceGroupManager{}
		err = Convert(manager, managerV1)
		if err != nil {
			return err
		}
		managerV1.ForceSendFields = manager.ForceSendFields
		op, err = config.clientCompute.RegionInstanceGroupManagers.Insert(project, d.Get("region").(string), managerV1).Do()
	case v0beta:
		op, err = config.clientComputeBeta.RegionInstanceGroupManagers.Insert(project, d.Get("region").(string), manager).Do()
	}

	if err != nil {
		return fmt.Errorf("Error creating RegionInstanceGroupManager: %s", err)
	}

	d.SetId(manager.Name)

	// Wait for the operation to complete
	err = computeSharedOperationWait(config, op, project, "Creating InstanceGroupManager")
	if err != nil {
		return err
	}

	return resourceComputeRegionInstanceGroupManagerRead(d, meta)
}

func resourceComputeRegionInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	manager := &computeBeta.InstanceGroupManager{}
	switch computeApiVersion {
	case v1:
		v1Manager := &compute.InstanceGroupManager{}
		v1Manager, err = config.clientCompute.RegionInstanceGroupManagers.Get(project, region, d.Id()).Do()

		err = Convert(v1Manager, manager)
		if err != nil {
			return err
		}
	case v0beta:
		manager, err = config.clientComputeBeta.RegionInstanceGroupManagers.Get(project, region, d.Id()).Do()
	}

	if err != nil {
		handleNotFoundError(err, d, fmt.Sprintf("Region Instance Manager %q", d.Get("name").(string)))
	}

	d.Set("base_instance_name", manager.BaseInstanceName)
	d.Set("instance_template", manager.InstanceTemplate)
	d.Set("name", manager.Name)
	d.Set("region", GetResourceNameFromSelfLink(manager.Region))
	d.Set("description", manager.Description)
	d.Set("project", project)
	d.Set("target_size", manager.TargetSize)
	d.Set("target_pools", manager.TargetPools)
	d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts))
	d.Set("fingerprint", manager.Fingerprint)
	d.Set("instance_group", manager.InstanceGroup)
	d.Set("auto_healing_policies", flattenAutoHealingPolicies(manager.AutoHealingPolicies))
	d.Set("self_link", ConvertSelfLinkToV1(manager.SelfLink))

	return nil
}

func resourceComputeRegionInstanceGroupManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersionUpdate(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures, []Feature{})
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)

	d.Partial(true)

	if d.HasChange("target_pools") {
		targetPools := convertStringSet(d.Get("target_pools").(*schema.Set))

		// Build the parameter
		setTargetPools := &computeBeta.RegionInstanceGroupManagersSetTargetPoolsRequest{
			Fingerprint: d.Get("fingerprint").(string),
			TargetPools: targetPools,
		}

		var op interface{}
		switch computeApiVersion {
		case v1:
			setTargetPoolsV1 := &compute.RegionInstanceGroupManagersSetTargetPoolsRequest{}
			err = Convert(setTargetPools, setTargetPoolsV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.RegionInstanceGroupManagers.SetTargetPools(
				project, region, d.Id(), setTargetPoolsV1).Do()
		case v0beta:
			setTargetPoolsV0beta := &computeBeta.RegionInstanceGroupManagersSetTargetPoolsRequest{}
			err = Convert(setTargetPools, setTargetPoolsV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.RegionInstanceGroupManagers.SetTargetPools(
				project, region, d.Id(), setTargetPoolsV0beta).Do()
		}

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config, op, project, "Updating RegionInstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_pools")
	}

	if d.HasChange("instance_template") {
		// Build the parameter
		setInstanceTemplate := &computeBeta.RegionInstanceGroupManagersSetTemplateRequest{
			InstanceTemplate: d.Get("instance_template").(string),
		}

		var op interface{}
		switch computeApiVersion {
		case v1:
			setInstanceTemplateV1 := &compute.RegionInstanceGroupManagersSetTemplateRequest{}
			err = Convert(setInstanceTemplate, setInstanceTemplateV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.RegionInstanceGroupManagers.SetInstanceTemplate(
				project, region, d.Id(), setInstanceTemplateV1).Do()
		case v0beta:
			setInstanceTemplateV0beta := &computeBeta.RegionInstanceGroupManagersSetTemplateRequest{}
			err = Convert(setInstanceTemplate, setInstanceTemplateV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.RegionInstanceGroupManagers.SetInstanceTemplate(
				project, region, d.Id(), setInstanceTemplateV0beta).Do()
		}

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config, op, project, "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("instance_template")
	}

	if d.HasChange("named_port") {
		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").([]interface{}))
		setNamedPorts := &computeBeta.RegionInstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		var op interface{}
		switch computeApiVersion {
		case v1:
			setNamedPortsV1 := &compute.RegionInstanceGroupsSetNamedPortsRequest{}
			err = Convert(setNamedPorts, setNamedPortsV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.RegionInstanceGroups.SetNamedPorts(
				project, region, d.Id(), setNamedPortsV1).Do()
		case v0beta:
			setNamedPortsV0beta := &computeBeta.RegionInstanceGroupsSetNamedPortsRequest{}
			err = Convert(setNamedPorts, setNamedPortsV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.RegionInstanceGroups.SetNamedPorts(
				project, region, d.Id(), setNamedPortsV0beta).Do()
		}

		if err != nil {
			return fmt.Errorf("Error updating RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete:
		err = computeSharedOperationWait(config, op, project, "Updating RegionInstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("named_port")
	}

	if d.HasChange("target_size") {
		targetSize := int64(d.Get("target_size").(int))
		var op interface{}
		switch computeApiVersion {
		case v1:
			op, err = config.clientCompute.RegionInstanceGroupManagers.Resize(
				project, region, d.Id(), targetSize).Do()
		case v0beta:
			op, err = config.clientComputeBeta.RegionInstanceGroupManagers.Resize(
				project, region, d.Id(), targetSize).Do()
		}

		if err != nil {
			return fmt.Errorf("Error resizing RegionInstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config, op, project, "Resizing RegionInstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_size")
	}

	if d.HasChange("auto_healing_policies") {
		setAutoHealingPoliciesRequest := &computeBeta.RegionInstanceGroupManagersSetAutoHealingRequest{}
		if v, ok := d.GetOk("auto_healing_policies"); ok {
			setAutoHealingPoliciesRequest.AutoHealingPolicies = expandAutoHealingPolicies(v.([]interface{}))
		}

		op, err := config.clientComputeBeta.RegionInstanceGroupManagers.SetAutoHealingPolicies(
			project, region, d.Id(), setAutoHealingPoliciesRequest).Do()

		if err != nil {
			return fmt.Errorf("Error updating AutoHealingPolicies: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWait(config, op, project, "Updating AutoHealingPolicies")
		if err != nil {
			return err
		}

		d.SetPartial("auto_healing_policies")
	}

	d.Partial(false)

	return resourceComputeRegionInstanceGroupManagerRead(d, meta)
}

func resourceComputeRegionInstanceGroupManagerDelete(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)

	var op interface{}
	switch computeApiVersion {
	case v1:
		op, err = config.clientCompute.RegionInstanceGroupManagers.Delete(project, region, d.Id()).Do()
	case v0beta:
		op, err = config.clientComputeBeta.RegionInstanceGroupManagers.Delete(project, region, d.Id()).Do()
	}

	if err != nil {
		return fmt.Errorf("Error deleting region instance group manager: %s", err)
	}

	// Wait for the operation to complete
	err = computeSharedOperationWait(config, op, project, "Deleting RegionInstanceGroupManager")

	d.SetId("")
	return nil
}

func resourceComputeRegionInstanceGroupManagerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	computeApiVersion := getComputeApiVersion(d, RegionInstanceGroupManagerBaseApiVersion, RegionInstanceGroupManagerVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return false, err
	}

	region := d.Get("region").(string)

	switch computeApiVersion {
	case v1:
		_, err = config.clientCompute.RegionInstanceGroupManagers.Get(project, region, d.Id()).Do()
	case v0beta:
		_, err = config.clientComputeBeta.RegionInstanceGroupManagers.Get(project, region, d.Id()).Do()
	}

	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return false, nil
		}
		// There was some other error in reading the resource but we can't say for sure if it doesn't exist.
		return true, err
	}
	return true, nil

}
