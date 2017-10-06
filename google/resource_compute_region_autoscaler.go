package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
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

			"autoscaling_policy": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_replicas": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"max_replicas": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"cooldown_period": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  60,
						},

						"cpu_utilization": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": &schema.Schema{
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},

						"metric": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"target": &schema.Schema{
										Type:     schema.TypeFloat,
										Required: true,
									},

									"type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"load_balancing_utilization": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": &schema.Schema{
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

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

func buildRegionAutoscaler(d *schema.ResourceData) (*compute.Autoscaler, error) {
	// Build the parameter
	scaler := &compute.Autoscaler{
		Name:   d.Get("name").(string),
		Target: d.Get("target").(string),
	}

	// Optional fields
	if v, ok := d.GetOk("description"); ok {
		scaler.Description = v.(string)
	}

	aspCount := d.Get("autoscaling_policy.#").(int)
	if aspCount != 1 {
		return nil, fmt.Errorf("The autoscaler must have exactly one autoscaling_policy, found %d.", aspCount)
	}

	prefix := "autoscaling_policy.0."

	scaler.AutoscalingPolicy = &compute.AutoscalingPolicy{
		MaxNumReplicas:    int64(d.Get(prefix + "max_replicas").(int)),
		MinNumReplicas:    int64(d.Get(prefix + "min_replicas").(int)),
		CoolDownPeriodSec: int64(d.Get(prefix + "cooldown_period").(int)),
	}

	// Check that only one autoscaling policy is defined

	policyCounter := 0
	if _, ok := d.GetOk(prefix + "cpu_utilization"); ok {
		if d.Get(prefix+"cpu_utilization.0.target").(float64) != 0 {
			cpuUtilCount := d.Get(prefix + "cpu_utilization.#").(int)
			if cpuUtilCount != 1 {
				return nil, fmt.Errorf("The autoscaling_policy must have exactly one cpu_utilization, found %d.", cpuUtilCount)
			}
			policyCounter++
			scaler.AutoscalingPolicy.CpuUtilization = &compute.AutoscalingPolicyCpuUtilization{
				UtilizationTarget: d.Get(prefix + "cpu_utilization.0.target").(float64),
			}
		}
	}
	if _, ok := d.GetOk("autoscaling_policy.0.metric"); ok {
		if d.Get(prefix+"metric.0.name") != "" {
			policyCounter++
			metricCount := d.Get(prefix + "metric.#").(int)
			if metricCount != 1 {
				return nil, fmt.Errorf("The autoscaling_policy must have exactly one metric, found %d.", metricCount)
			}
			scaler.AutoscalingPolicy.CustomMetricUtilizations = []*compute.AutoscalingPolicyCustomMetricUtilization{
				{
					Metric:                d.Get(prefix + "metric.0.name").(string),
					UtilizationTarget:     d.Get(prefix + "metric.0.target").(float64),
					UtilizationTargetType: d.Get(prefix + "metric.0.type").(string),
				},
			}
		}

	}
	if _, ok := d.GetOk("autoscaling_policy.0.load_balancing_utilization"); ok {
		if d.Get(prefix+"load_balancing_utilization.0.target").(float64) != 0 {
			policyCounter++
			lbuCount := d.Get(prefix + "load_balancing_utilization.#").(int)
			if lbuCount != 1 {
				return nil, fmt.Errorf("The autoscaling_policy must have exactly one load_balancing_utilization, found %d.", lbuCount)
			}
			scaler.AutoscalingPolicy.LoadBalancingUtilization = &compute.AutoscalingPolicyLoadBalancingUtilization{
				UtilizationTarget: d.Get(prefix + "load_balancing_utilization.0.target").(float64),
			}
		}
	}

	if policyCounter != 1 {
		return nil, fmt.Errorf("One policy must be defined for an autoscaler.")
	}

	return scaler, nil
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

	err = computeOperationWait(config, op, project, "Creating Autoscaler")
	if err != nil {
		return err
	}

	return resourceComputeRegionAutoscalerRead(d, meta)
}

func flattenRegionAutoscalingPolicy(policy *compute.AutoscalingPolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	policyMap := make(map[string]interface{})
	policyMap["max_replicas"] = policy.MaxNumReplicas
	policyMap["min_replicas"] = policy.MinNumReplicas
	policyMap["cooldown_period"] = policy.CoolDownPeriodSec
	if policy.CpuUtilization != nil {
		cpuUtils := make([]map[string]interface{}, 0, 1)
		cpuUtil := make(map[string]interface{})
		cpuUtil["target"] = policy.CpuUtilization.UtilizationTarget
		cpuUtils = append(cpuUtils, cpuUtil)
		policyMap["cpu_utilization"] = cpuUtils
	}
	if policy.LoadBalancingUtilization != nil {
		loadBalancingUtils := make([]map[string]interface{}, 0, 1)
		loadBalancingUtil := make(map[string]interface{})
		loadBalancingUtil["target"] = policy.LoadBalancingUtilization.UtilizationTarget
		loadBalancingUtils = append(loadBalancingUtils, loadBalancingUtil)
		policyMap["load_balancing_utilization"] = loadBalancingUtils
	}
	if policy.CustomMetricUtilizations != nil {
		metricUtils := make([]map[string]interface{}, 0, len(policy.CustomMetricUtilizations))
		for _, customMetricUtilization := range policy.CustomMetricUtilizations {
			metricUtil := make(map[string]interface{})
			metricUtil["target"] = customMetricUtilization.UtilizationTarget
			metricUtils = append(metricUtils, metricUtil)
		}
		policyMap["metric"] = metricUtils
	}
	result = append(result, policyMap)
	return result
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

	regionUrl := strings.Split(scaler.Region, "/")
	d.Set("self_link", scaler.SelfLink)
	d.Set("name", scaler.Name)
	d.Set("target", scaler.Target)
	d.Set("region", regionUrl[len(regionUrl)-1])
	d.Set("description", scaler.Description)
	if scaler.AutoscalingPolicy != nil {
		d.Set("autoscaling_policy", flattenRegionAutoscalingPolicy(scaler.AutoscalingPolicy))
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

	err = computeOperationWait(config, op, project, "Updating Autoscaler")
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

	err = computeOperationWait(config, op, project, "Deleting Autoscaler")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
