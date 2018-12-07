package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/http"
	"regexp"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/spanner/v1"
)

const (
	spannerDatabaseNameFormat = "^[a-z][a-z0-9_-]*[a-z0-9]$"
)

func resourceSpannerDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpannerDatabaseCreate,
		Read:   resourceSpannerDatabaseRead,
		Delete: resourceSpannerDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSpannerDatabaseImport("name"),
		},

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateResourceSpannerDatabaseName,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"ddl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSpannerDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := buildSpannerDatabaseId(d, config)
	if err != nil {
		return err
	}

	cdr := &spanner.CreateDatabaseRequest{}
	cdr.CreateStatement = fmt.Sprintf("CREATE DATABASE `%s`", id.Database)
	if v, ok := d.GetOk("ddl"); ok {
		cdr.ExtraStatements = convertStringArr(v.([]interface{}))
	}

	op, err := config.clientSpanner.Projects.Instances.Databases.Create(
		id.parentInstanceUri(), cdr).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusConflict {
			return fmt.Errorf("Error, A database with name %s already exists in this instance", id.Database)
		}
		return fmt.Errorf("Error, failed to create database %s: %s", id.Database, err)
	}

	d.SetId(id.terraformId())

	// Wait until it's created
	timeoutMins := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := spannerDatabaseOperationWait(config, op, "Creating Spanner database", timeoutMins)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Spanner database %s has been created", id.terraformId())
	return resourceSpannerDatabaseRead(d, meta)
}

func resourceSpannerDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := buildSpannerDatabaseId(d, config)
	if err != nil {
		return err
	}

	db, err := config.clientSpanner.Projects.Instances.Databases.Get(
		id.databaseUri()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Spanner database %q", id.databaseUri()))
	}

	d.Set("state", db.State)
	d.Set("project", id.Project)
	return nil
}

func resourceSpannerDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := buildSpannerDatabaseId(d, config)
	if err != nil {
		return err
	}

	_, err = config.clientSpanner.Projects.Instances.Databases.DropDatabase(
		id.databaseUri()).Do()
	if err != nil {
		return fmt.Errorf("Error, failed to delete Spanner Database %s: %s", id.databaseUri(), err)
	}

	d.SetId("")
	return nil
}

func resourceSpannerDatabaseImport(databaseField string) schema.StateFunc {
	return func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		config := meta.(*Config)
		err := parseImportId([]string{
			fmt.Sprintf("projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<%s>[^/]+)", databaseField),
			fmt.Sprintf("instances/(?P<instance>[^/]+)/databases/(?P<%s>[^/]+)", databaseField),
			fmt.Sprintf("(?P<project>[^/]+)/(?P<instance>[^/]+)/(?P<%s>[^/]+)", databaseField),
			fmt.Sprintf("(?P<instance>[^/]+)/(?P<%s>[^/]+)", databaseField),
		}, d, config)
		if err != nil {
			return nil, fmt.Errorf("Error constructing id: %s", err)
		}

		id, err := buildSpannerDatabaseId(d, config)
		if err != nil {
			return nil, fmt.Errorf("Error constructing id: %s", err)
		}

		d.SetId(id.terraformId())

		return []*schema.ResourceData{d}, nil
	}
}

func buildSpannerDatabaseId(d *schema.ResourceData, config *Config) (*spannerDatabaseId, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}
	database, ok := d.GetOk("name")
	if !ok {
		database = d.Get("database")
	}

	dbName := database.(string)
	instanceName := d.Get("instance").(string)

	return &spannerDatabaseId{
		Project:  project,
		Instance: instanceName,
		Database: dbName,
	}, nil
}

type spannerDatabaseId struct {
	Project  string
	Instance string
	Database string
}

func (s spannerDatabaseId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Instance, s.Database)
}

func (s spannerDatabaseId) parentProjectUri() string {
	return fmt.Sprintf("projects/%s", s.Project)
}

func (s spannerDatabaseId) parentInstanceUri() string {
	return fmt.Sprintf("%s/instances/%s", s.parentProjectUri(), s.Instance)
}

func (s spannerDatabaseId) databaseUri() string {
	return fmt.Sprintf("%s/databases/%s", s.parentInstanceUri(), s.Database)
}

func validateResourceSpannerDatabaseName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) < 2 && len(value) > 30 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 2 and 30 characters in length", k))
	}

	if !regexp.MustCompile(spannerDatabaseNameFormat).MatchString(value) {
		errors = append(errors, fmt.Errorf("database name %q must match regexp %q", value, spannerDatabaseNameFormat))
	}
	return
}
