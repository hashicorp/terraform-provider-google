package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputeRegionAutoscaler() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionAutoscalerCreate,
		Read:   resourceComputeRegionAutoscalerRead,
		Update: resourceComputeRegionAutoscalerUpdate,
		Delete: resourceComputeRegionAutoscalerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"target": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"autoscaling_policy": autoscalingPolicy,

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeRegionAutoscalerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the region
	log.Printf("[DEBUG] Loading region: %s", d.Get("region").(string))
	region, err := config.clientCompute.Regions.Get(
		project, d.Get("region").(string)).Do()
	if err != nil {
		return fmt.Errorf(
			"Error loading region '%s': %s", d.Get("region").(string), err)
	}

	scaler, err := buildAutoscaler(d)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.RegionAutoscalers.Insert(
		project, region.Name, scaler).Do()
	if err != nil {
		return fmt.Errorf("Error creating Autoscaler: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(scaler.Name)

	err = computeOperationWait(config.clientCompute, op, project, "Creating Autoscaler")
	if err != nil {
		return err
	}

	return resourceComputeRegionAutoscalerRead(d, meta)
}

func resourceComputeRegionAutoscalerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	scaler, err := config.clientCompute.RegionAutoscalers.Get(
		project, region, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Autoscaler %q", d.Id()))
	}

	if scaler == nil {
		log.Printf("[WARN] Removing Autoscaler %q because it's gone", d.Get("name").(string))
		d.SetId("")
		return nil
	}

	d.Set("self_link", scaler.SelfLink)
	d.Set("name", scaler.Name)
	d.Set("target", scaler.Target)
	d.Set("region", GetResourceNameFromSelfLink(scaler.Region))
	d.Set("description", scaler.Description)
	if scaler.AutoscalingPolicy != nil {
		d.Set("autoscaling_policy", flattenAutoscalingPolicy(scaler.AutoscalingPolicy))
	}

	return nil
}

func resourceComputeRegionAutoscalerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)

	scaler, err := buildAutoscaler(d)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.RegionAutoscalers.Update(
		project, region, scaler).Do()
	if err != nil {
		return fmt.Errorf("Error updating Autoscaler: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(scaler.Name)

	err = computeOperationWait(config.clientCompute, op, project, "Updating Autoscaler")
	if err != nil {
		return err
	}

	return resourceComputeRegionAutoscalerRead(d, meta)
}

func resourceComputeRegionAutoscalerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	op, err := config.clientCompute.RegionAutoscalers.Delete(
		project, region, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting autoscaler: %s", err)
	}

	err = computeOperationWait(config.clientCompute, op, project, "Deleting Autoscaler")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
