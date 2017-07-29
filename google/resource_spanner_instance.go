package google

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/spanner/v1"
)

func resourceSpannerInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpannerInstanceCreate,
		Read:   resourceSpannerInstanceRead,
		Update: resourceSpannerInstanceUpdate,
		Delete: resourceSpannerInstanceDelete,

		Schema: map[string]*schema.Schema{

			"config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"num_nodes": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     schema.TypeString,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSpannerInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	timeoutMins := int(d.Timeout(schema.TimeoutCreate).Minutes())

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	cir := &spanner.CreateInstanceRequest{
		Instance: &spanner.Instance{},
	}

	cir.Instance.Config = instanceConfigForApi(project, d.Get("config").(string))
	cir.Instance.DisplayName = d.Get("display_name").(string)
	cir.Instance.NodeCount = int64(d.Get("num_nodes").(int))

	if v, ok := d.GetOk("labels"); ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		cir.Instance.Labels = m
	}
	if v, ok := d.GetOk("name"); ok {
		cir.InstanceId = v.(string)
	} else {
		cir.InstanceId = resource.UniqueId()
		d.Set("name", cir.InstanceId)
	}

	op, err := config.clientSpanner.Projects.Instances.Create(projectNameForApi(project), cir).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusConflict {
			return fmt.Errorf("Error, the name %s is not unique and already used", cir.InstanceId)
		}
		return fmt.Errorf("Error, failed to create instance %s: %s", cir.InstanceId, err)
	}

	// Wait until it's created
	waitErr := spannerInstanceOperationWait(config, op, "Creating Spanner instance", timeoutMins)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Spanner instance %s has been created", cir.Instance.Name)
	d.SetId(cir.InstanceId)

	return resourceSpannerInstanceRead(d, meta)

}

func resourceSpannerInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	instanceName := d.Get("name").(string)

	instance, err := config.clientSpanner.Projects.Instances.Get(instanceNameForApi(project, instanceName)).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Spanner instance %q", instanceName))
	}

	d.Set("config", extractInstanceConfigFromApi(instance.Config))
	d.Set("labels", instance.Labels)
	d.Set("name", extractInstanceNameFromApi(instance.Name))
	d.Set("display_name", instance.DisplayName)
	d.Set("num_nodes", instance.NodeCount)
	d.Set("status", instance.State)

	return nil
}

func resourceSpannerInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] About to update Spanner Instance %s ", d.Id())
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	timeoutMins := int(d.Timeout(schema.TimeoutUpdate).Minutes())
	instanceName := d.Get("name").(string)

	uir := &spanner.UpdateInstanceRequest{
		Instance:  &spanner.Instance{},
		FieldMask: "",
	}

	if d.HasChange("num_nodes") {
		uir.FieldMask = "nodeCount"
		uir.Instance.NodeCount = int64(d.Get("num_nodes").(int))
	}
	if d.HasChange("display_name") {
		if uir.FieldMask != "" {
			uir.FieldMask = uir.FieldMask + ","
		}
		uir.FieldMask = uir.FieldMask + "displayName"
		uir.Instance.DisplayName = d.Get("display_name").(string)
	}

	op, err := config.clientSpanner.Projects.Instances.Patch(
		instanceNameForApi(project, instanceName), uir).Do()

	// Wait until it's updated
	err = spannerInstanceOperationWait(config, op, "Update Spanner Instance", timeoutMins)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Spanner Instance %s has been updated ", d.Id())
	return resourceSpannerInstanceRead(d, meta)
}

func resourceSpannerInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	instanceName := d.Get("name").(string)
	_, err = config.clientSpanner.Projects.Instances.Delete(
		instanceNameForApi(project, instanceName)).Do()
	if err != nil {
		return fmt.Errorf("Error, failed to delete Spanner Instance %s: %s", d.Get("name").(string), err)
	}

	d.SetId("")
	return nil
}

func extractInstanceConfigFromApi(nameUri string) string {
	rUris := strings.Split(nameUri, "/")
	return rUris[len(rUris)-1]
}

func extractInstanceNameFromApi(nameUri string) string {
	rUris := strings.Split(nameUri, "/")
	return rUris[len(rUris)-1]
}

func instanceNameForApi(p, i string) string {
	return projectNameForApi(p) + "/instances/" + i
}

func instanceConfigForApi(p, c string) string {
	return projectNameForApi(p) + "/instanceConfigs/" + c
}

func projectNameForApi(p string) string {
	return "projects/" + p
}
