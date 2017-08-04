package google

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

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
		Importer: &schema.ResourceImporter{
			State: resourceSpannerInstanceImportState,
		},

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
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) < 6 && len(value) > 30 {
						errors = append(errors, fmt.Errorf(
							"%q must be between 6 and 30 characters in length", k))
					}
					if !regexp.MustCompile("^[a-z0-9-]+$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q can only contain lowercase letters, numbers and hyphens", k))
					}
					if !regexp.MustCompile("^[a-z]").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must start with a letter", k))
					}
					if !regexp.MustCompile("[a-z0-9]$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must end with a number or a letter", k))
					}
					return
				},
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) < 4 && len(value) > 30 {
						errors = append(errors, fmt.Errorf(
							"%q must be between 4 and 30 characters in length", k))
					}
					return
				},
			},

			"num_nodes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSpannerInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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

	if v, ok := d.GetOk("name"); ok {
		cir.InstanceId = v.(string)
	} else {
		cir.InstanceId = genSpannerInstanceId()
		d.Set("name", cir.InstanceId)
	}
	if v, ok := d.GetOk("labels"); ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		cir.Instance.Labels = m
	}

	op, err := config.clientSpanner.Projects.Instances.Create(
		projectNameForApi(project), cir).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusConflict {
			return fmt.Errorf("Error, the name %s is not unique and already used", cir.InstanceId)
		}
		return fmt.Errorf("Error, failed to create instance %s: %s", cir.InstanceId, err)
	}

	d.SetId(cir.InstanceId)

	// Wait until it's created
	timeoutMins := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := spannerInstanceOperationWait(config, op, "Creating Spanner instance", timeoutMins)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Spanner instance %s has been created", cir.Instance.Name)

	return resourceSpannerInstanceRead(d, meta)

}

func resourceSpannerInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := d.Get("name").(string)
	instance, err := config.clientSpanner.Projects.Instances.Get(
		instanceNameForApi(project, instanceName)).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Spanner instance %q", instanceName))
	}

	d.Set("config", extractInstanceConfigFromApi(instance.Config))
	d.Set("labels", instance.Labels)
	d.Set("display_name", instance.DisplayName)
	d.Set("num_nodes", instance.NodeCount)

	return nil
}

func resourceSpannerInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] About to update Spanner Instance %s ", d.Id())
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	uir := &spanner.UpdateInstanceRequest{
		Instance: &spanner.Instance{},
	}

	fieldMask := []string{}
	if d.HasChange("num_nodes") {
		fieldMask = append(fieldMask, "nodeCount")
		uir.Instance.NodeCount = int64(d.Get("num_nodes").(int))
	}
	if d.HasChange("display_name") {
		fieldMask = append(fieldMask, "displayName")
		uir.Instance.DisplayName = d.Get("display_name").(string)
	}

	instanceName := d.Get("name").(string)
	uir.FieldMask = strings.Join(fieldMask, ",")
	op, err := config.clientSpanner.Projects.Instances.Patch(
		instanceNameForApi(project, instanceName), uir).Do()

	// Wait until it's updated
	timeoutMins := int(d.Timeout(schema.TimeoutUpdate).Minutes())
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

func resourceSpannerInstanceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	sid, err := extractSpannerInstanceImportIds(d.Id())
	if err != nil {
		return nil, err
	}

	if sid.Project != "" {
		d.Set("project", sid.Project)
	}
	d.Set("name", sid.Instance)
	d.SetId(sid.Instance)

	return []*schema.ResourceData{d}, nil
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

func genSpannerInstanceId() string {
	return fmt.Sprintf("tfgen-spanid-%010d", rand.Int63n(999999))
}

type spannerInstanceImportId struct {
	Project  string
	Instance string
}

func extractSpannerInstanceImportIds(id string) (*spannerInstanceImportId, error) {
	parts := strings.Split(id, "/")
	if id == "" || strings.HasPrefix(id, "/") || strings.HasSuffix(id, "/") ||
		(len(parts) != 1 && len(parts) != 2) {
		return nil, fmt.Errorf("Invalid spanner database specifier. " +
			"Expecting either {projectId}/{instanceId} OR " +
			"{instanceId} (where project is to be derived from that specified in provider)")
	}

	sid := &spannerInstanceImportId{}

	if len(parts) == 1 {
		log.Printf("[INFO] Spanner instance import format of {instanceId} specified: %s", id)
		sid.Instance = parts[0]
	}
	if len(parts) == 2 {
		log.Printf("[INFO] Spanner instance import format of {projectId}/{instanceId} specified: %s", id)
		sid.Project = parts[0]
		sid.Instance = parts[1]
	}
	return sid, nil

}
