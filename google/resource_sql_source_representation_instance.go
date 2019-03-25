package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func resourceSqlSourceRepresentationInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlSourceRepresentationInstanceCreate,
		Read:   resourceSqlSourceRepresentationInstanceRead,
		Delete: resourceSqlSourceRepresentationInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSqlSourceRepresentationInstanceImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"database_version": {
				Type:     schema.TypeString,
				Default:  "MYSQL_5_7",
				Optional: true,
				ForceNew: true,
			},

			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port": {
				Type:         schema.TypeInt,
				Default:      3306,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSqlSourceRepresentationInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		name = resource.UniqueId()
	}
	d.Set("name", name)

	host := d.Get("host").(string)
	port := d.Get("port").(int)
	hostPort := fmt.Sprintf("%s:%d", host, port)

	instance := &sqladmin.DatabaseInstance{
		Name:            name,
		Region:          region,
		DatabaseVersion: d.Get("database_version").(string),
		OnPremisesConfiguration: &sqladmin.OnPremisesConfiguration{
			HostPort: hostPort,
		},
	}

	backoff := time.Second
	var op *sqladmin.Operation
	for {
		op, err = config.clientSqlAdmin.Instances.Insert(project, instance).Do()
		if err == nil {
			break
		}

		// When deleting and recreating a source representation instance, the (re)create operation fails with an invalidState
		// for a very short window after the deletion succeeds. This is easily fixed by retrying when encountering that
		// error.
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 && strings.Contains(gerr.Body, "invalidState") {
			log.Printf("[DEBUG]: Got an invalid state error, retrying after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return errwrap.Wrapf(fmt.Sprintf("Error, constantly failing to create source representation instance %s. Too many invalid state errors. Latest error {{err}}", instance.Name), err)
			}
		} else {
			return errwrap.Wrapf(fmt.Sprintf("Failed to create source representation instance %s: {{err}}", instance.Name), err)
		}
	}

	d.SetId(instance.Name)

	err = sqladminOperationWaitTime(config, op, project, "Create Source Representation Instance", int(d.Timeout(schema.TimeoutCreate).Minutes()))
	if err != nil {
		d.SetId("")
		return err
	}

	return resourceSqlSourceRepresentationInstanceRead(d, meta)
}

func resourceSqlSourceRepresentationInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance, err := config.clientSqlAdmin.Instances.Get(project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL Source Representation Instance %q", d.Get("name").(string)))
	}

	d.Set("name", instance.Name)
	d.Set("region", instance.Region)
	d.Set("database_version", instance.DatabaseVersion)
	hostPort := strings.Split(instance.OnPremisesConfiguration.HostPort, ":")
	d.Set("host", hostPort[0])

	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return err
	}
	d.Set("port", port)

	d.Set("project", project)
	d.Set("self_link", instance.SelfLink)
	d.SetId(instance.Name)

	return nil
}

func resourceSqlSourceRepresentationInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	var op *sqladmin.Operation
	err = retryTimeDuration(func() error {
		op, err = config.clientSqlAdmin.Instances.Delete(project, d.Get("name").(string)).Do()
		return err
	}, d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return fmt.Errorf("Error, failed to delete source representation instance %s: %s", d.Get("name").(string), err)
	}

	err = sqladminOperationWaitTime(config, op, project, "Delete Source Representation Instance", int(d.Timeout(schema.TimeoutDelete).Minutes()))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSqlSourceRepresentationInstanceImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
