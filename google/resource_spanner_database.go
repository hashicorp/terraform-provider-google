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
		},
	}
}

func resourceSpannerDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	timeoutMins := int(d.Timeout(schema.TimeoutCreate).Minutes())

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := d.Get("instance").(string)
	dbName := d.Get("name").(string)

	cdr := &spanner.CreateDatabaseRequest{}
	cdr.CreateStatement = fmt.Sprintf("CREATE DATABASE `%s`", dbName)

	if v, ok := d.GetOk("ddl"); ok {
		ddlList := v.([]interface{})
		ddls := []string{}
		for _, v := range ddlList {
			ddls = append(ddls, v.(string))
		}
		cdr.ExtraStatements = ddls
	}

	op, err := config.clientSpanner.Projects.Instances.Databases.Create(
		instanceNameForApi(project, instanceName), cdr).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusConflict {
			return fmt.Errorf("Error, A database with name %s already exists", dbName)
		}
		return fmt.Errorf("Error, failed to create database %s: %s", dbName, err)
	}

	// Wait until it's created
	waitErr := spannerDatabaseOperationWait(config, op, "Creating Spanner database", timeoutMins)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Spanner database %s has been created", dbName)
	d.SetId(dbName)

	return resourceSpannerDatabaseRead(d, meta)

}

func resourceSpannerDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	dbName := d.Get("name").(string)
	instanceName := d.Get("instance").(string)
	_, err = config.clientSpanner.Projects.Instances.Databases.Get(
		databaseNameForApi(project, instanceName, dbName)).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Spanner database %q", dbName))
	}

	return nil
}

func resourceSpannerDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	dbName := d.Get("name").(string)
	instanceName := d.Get("instance").(string)
	_, err = config.clientSpanner.Projects.Instances.Databases.DropDatabase(
		databaseNameForApi(project, instanceName, dbName)).Do()
	if err != nil {
		return fmt.Errorf("Error, failed to delete Spanner Database %s: %s", dbName, err)
	}

	d.SetId("")
	return nil
}

func databaseNameForApi(p, i, d string) string {
	return instanceNameForApi(p, i) + "/databases/" + d
}

func resourceSpannerDatabaseImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importId, err := extractSpannerDatabaseImportIds(d.Id())
	if err != nil {
		return nil, err
	}

	if importId.Project != "" {
		d.Set("project", importId.Project)
	}
	d.Set("instance", importId.Instance)
	d.Set("name", importId.Database)
	d.SetId(importId.Database)

	return []*schema.ResourceData{d}, nil
}

type spannerDbImportId struct {
	Project  string
	Instance string
	Database string
}

func extractSpannerDatabaseImportIds(id string) (*spannerDbImportId, error) {
	parts := strings.Split(id, "/")
	if id == "" || strings.HasPrefix(id, "/") || strings.HasSuffix(id, "/") ||
		(len(parts) != 2 && len(parts) != 3) {
		return nil, fmt.Errorf("Invalid spanner database specifier. " +
			"Expecting either {projectId}/{instanceId}/{dbId} OR " +
			"{instanceId}/{dbId} (where project will be derived from the provider)")
	}

	sid := &spannerDbImportId{}

	if len(parts) == 2 {
		log.Printf("[INFO] Spanner instance import format of {instanceId}/{dbId} specified: %s", id)
		sid.Instance = parts[0]
		sid.Database = parts[1]
	}
	if len(parts) == 3 {
		log.Printf("[INFO] Spanner instance import format of {projectId}/{instanceId}/{dbId} specified: %s", id)
		sid.Project = parts[0]
		sid.Instance = parts[1]
		sid.Database = parts[2]
	}

	return sid, nil
}
