package google

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSpannerInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cir := &spanner.CreateInstanceRequest{
		Instance: &spanner.Instance{},
	}

	if v, ok := d.GetOk("name"); ok {
		cir.InstanceId = v.(string)
	} else {
		cir.InstanceId = genSpannerInstanceName()
		d.Set("name", cir.InstanceId)
	}

	if v, ok := d.GetOk("labels"); ok {
		cir.Instance.Labels = convertStringMap(v.(map[string]interface{}))
	}

	id, err := buildSpannerInstanceId(d, config)
	if err != nil {
		return err
	}

	cir.Instance.Config = id.instanceConfigUri(d.Get("config").(string))
	cir.Instance.DisplayName = d.Get("display_name").(string)
	cir.Instance.NodeCount = int64(d.Get("num_nodes").(int))

	op, err := config.clientSpanner.Projects.Instances.Create(
		id.parentProjectUri(), cir).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusConflict {
			return fmt.Errorf("Error, the name %s is not unique within project %s", id.Instance, id.Project)
		}
		return fmt.Errorf("Error, failed to create instance %s: %s", id.terraformId(), err)
	}

	d.SetId(id.terraformId())

	// Wait until it's created
	timeoutMins := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := spannerInstanceOperationWait(config, op, "Creating Spanner instance", timeoutMins)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Spanner instance %s has been created", id.terraformId())
	return resourceSpannerInstanceRead(d, meta)
}

func resourceSpannerInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := buildSpannerInstanceId(d, config)
	if err != nil {
		return err
	}

	instance, err := config.clientSpanner.Projects.Instances.Get(
		id.instanceUri()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Spanner instance %s", id.terraformId()))
	}

	d.Set("config", extractInstanceConfigFromUri(instance.Config))
	d.Set("labels", instance.Labels)
	d.Set("display_name", instance.DisplayName)
	d.Set("num_nodes", instance.NodeCount)
	d.Set("state", instance.State)

	return nil
}

func resourceSpannerInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[INFO] About to update Spanner Instance %s ", d.Id())
	uir := &spanner.UpdateInstanceRequest{
		Instance: &spanner.Instance{},
	}

	id, err := buildSpannerInstanceId(d, config)
	if err != nil {
		return err
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
	if d.HasChange("labels") {
		fieldMask = append(fieldMask, "labels")
		uir.Instance.Labels = convertStringMap(d.Get("labels").(map[string]interface{}))
	}

	uir.FieldMask = strings.Join(fieldMask, ",")
	op, err := config.clientSpanner.Projects.Instances.Patch(
		id.instanceUri(), uir).Do()
	if err != nil {
		return err
	}

	// Wait until it's updated
	timeoutMins := int(d.Timeout(schema.TimeoutUpdate).Minutes())
	err = spannerInstanceOperationWait(config, op, "Update Spanner Instance", timeoutMins)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Spanner Instance %s has been updated ", id.terraformId())
	return resourceSpannerInstanceRead(d, meta)
}

func resourceSpannerInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := buildSpannerInstanceId(d, config)
	if err != nil {
		return err
	}

	_, err = config.clientSpanner.Projects.Instances.Delete(
		id.instanceUri()).Do()
	if err != nil {
		return fmt.Errorf("Error, failed to delete Spanner Instance %s in project %s: %s", id.Instance, id.Project, err)
	}

	d.SetId("")
	return nil
}

func resourceSpannerInstanceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	id, err := importSpannerInstanceId(d.Id())
	if err != nil {
		return nil, err
	}

	if id.Project != "" {
		d.Set("project", id.Project)
	} else {
		project, err := getProject(d, config)
		if err != nil {
			return nil, err
		}
		id.Project = project
	}

	d.Set("name", id.Instance)
	d.SetId(id.terraformId())

	return []*schema.ResourceData{d}, nil
}

func buildSpannerInstanceId(d *schema.ResourceData, config *Config) (*spannerInstanceId, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}
	return &spannerInstanceId{
		Project:  project,
		Instance: d.Get("name").(string),
	}, nil
}

func extractInstanceConfigFromUri(configUri string) string {
	return extractLastResourceFromUri(configUri)
}

func extractInstanceNameFromUri(nameUri string) string {
	return extractLastResourceFromUri(nameUri)
}

func extractLastResourceFromUri(uri string) string {
	rUris := strings.Split(uri, "/")
	return rUris[len(rUris)-1]
}

func genSpannerInstanceName() string {
	return resource.PrefixedUniqueId("tfgen-spanid-")[:30]
}

type spannerInstanceId struct {
	Project  string
	Instance string
}

func (s spannerInstanceId) terraformId() string {
	return fmt.Sprintf("%s/%s", s.Project, s.Instance)
}

func (s spannerInstanceId) parentProjectUri() string {
	return fmt.Sprintf("projects/%s", s.Project)
}

func (s spannerInstanceId) instanceUri() string {
	return fmt.Sprintf("%s/instances/%s", s.parentProjectUri(), s.Instance)
}

func (s spannerInstanceId) instanceConfigUri(c string) string {
	return fmt.Sprintf("%s/instanceConfigs/%s", s.parentProjectUri(), c)
}

func importSpannerInstanceId(id string) (*spannerInstanceId, error) {
	if !regexp.MustCompile("^[a-z0-9-]+$").Match([]byte(id)) &&
		!regexp.MustCompile("^[a-z0-9-]+/[a-z0-9-]+$").Match([]byte(id)) {
		return nil, fmt.Errorf("Invalid spanner instance specifier. " +
			"Expecting either {projectId}/{instanceId} OR " +
			"{instanceId} (where project is to be derived from that specified in provider)")
	}

	parts := strings.Split(id, "/")
	if len(parts) == 1 {
		log.Printf("[INFO] Spanner instance import format of {instanceId} specified: %s", id)
		return &spannerInstanceId{Instance: parts[0]}, nil
	}

	log.Printf("[INFO] Spanner instance import format of {projectId}/{instanceId} specified: %s", id)
	return extractSpannerInstanceId(id)
}

func extractSpannerInstanceId(id string) (*spannerInstanceId, error) {
	if !regexp.MustCompile("^[a-z0-9-]+/[a-z0-9-]+$").Match([]byte(id)) {
		return nil, fmt.Errorf("Invalid spanner id format, expecting {projectId}/{instanceId}")
	}
	parts := strings.Split(id, "/")
	return &spannerInstanceId{
		Project:  parts[0],
		Instance: parts[1],
	}, nil
}
