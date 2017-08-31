package google

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/spanner/v1"
)

func resourceSpannerDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpannerDatabaseCreate,
		Read:   resourceSpannerDatabaseRead,
		Delete: resourceSpannerDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSpannerDatabaseImportState,
		},

		Schema: map[string]*schema.Schema{

			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) < 2 && len(value) > 30 {
						errors = append(errors, fmt.Errorf(
							"%q must be between 2 and 30 characters in length", k))
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

			"project": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceSpannerDatabaseImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	id, err := importSpannerDatabaseId(d.Id())
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

	d.Set("instance", id.Instance)
	d.Set("name", id.Database)
	d.SetId(id.terraformId())

	return []*schema.ResourceData{d}, nil
}

func buildSpannerDatabaseId(d *schema.ResourceData, config *Config) (*spannerDatabaseId, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}
	dbName := d.Get("name").(string)
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

func importSpannerDatabaseId(id string) (*spannerDatabaseId, error) {
	if !regexp.MustCompile("^[a-z0-9-]+/[a-z0-9-]+$").Match([]byte(id)) &&
		!regexp.MustCompile("^[a-z0-9-]+/[a-z0-9-]+/[a-z0-9-]+$").Match([]byte(id)) {
		return nil, fmt.Errorf("Invalid spanner database specifier. " +
			"Expecting either {projectId}/{instanceId}/{dbId} OR " +
			"{instanceId}/{dbId} (where project will be derived from the provider)")
	}

	parts := strings.Split(id, "/")
	if len(parts) == 2 {
		log.Printf("[INFO] Spanner database import format of {instanceId}/{dbId} specified: %s", id)
		return &spannerDatabaseId{Instance: parts[0], Database: parts[1]}, nil
	}

	log.Printf("[INFO] Spanner database import format of {projectId}/{instanceId}/{dbId} specified: %s", id)
	return extractSpannerDatabaseId(id)
}

func extractSpannerDatabaseId(id string) (*spannerDatabaseId, error) {
	if !regexp.MustCompile("^[a-z0-9-]+/[a-z0-9-]+/[a-z0-9-]+$").Match([]byte(id)) {
		return nil, fmt.Errorf("Invalid spanner id format, expecting {projectId}/{instanceId}/{databaseId}")
	}
	parts := strings.Split(id, "/")
	return &spannerDatabaseId{
		Project:  parts[0],
		Instance: parts[1],
		Database: parts[2],
	}, nil
}
