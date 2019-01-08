package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func resourceSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlDatabaseCreate,
		Read:   resourceSqlDatabaseRead,
		Update: resourceSqlDatabaseUpdate,
		Delete: resourceSqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSqlDatabaseImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"charset": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceSqlDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	database_name := d.Get("name").(string)
	instance_name := d.Get("instance").(string)
	d.SetId(instance_name + ":" + database_name)

	db := &sqladmin.Database{
		Name:      database_name,
		Instance:  instance_name,
		Charset:   d.Get("charset").(string),
		Collation: d.Get("collation").(string),
	}

	mutexKV.Lock(instanceMutexKey(project, instance_name))
	defer mutexKV.Unlock(instanceMutexKey(project, instance_name))

	var op *sqladmin.Operation
	err = retryTime(func() error {
		op, err = config.clientSqlAdmin.Databases.Insert(project, instance_name, db).Do()
		return err
	}, 5 /* minutes */)

	if err != nil {
		return fmt.Errorf("Error, failed to insert "+
			"database %s into instance %s: %s", database_name,
			instance_name, err)
	}

	err = sqladminOperationWaitTime(config, op, project, "Insert Database", int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if err != nil {
		return fmt.Errorf("Error, failure waiting for insertion of %s "+
			"into %s: %s", database_name, instance_name, err)
	}

	return resourceSqlDatabaseRead(d, meta)
}

func resourceSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	s := strings.Split(d.Id(), ":")

	if len(s) != 2 {
		return fmt.Errorf("Error, failure importing database %s. "+
			"ID format is instance:name", d.Id())
	}

	instance_name := s[0]
	database_name := s[1]

	var db *sqladmin.Database
	err = retryTime(func() error {
		db, err = config.clientSqlAdmin.Databases.Get(project, instance_name, database_name).Do()
		return err
	}, 5 /* minutes */)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL Database %q in instance %q", database_name, instance_name))
	}

	d.Set("instance", db.Instance)
	d.Set("name", db.Name)
	d.Set("self_link", db.SelfLink)
	d.SetId(instance_name + ":" + database_name)
	d.Set("charset", db.Charset)
	d.Set("collation", db.Collation)
	d.Set("project", project)

	return nil
}

func resourceSqlDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	database_name := d.Get("name").(string)
	instance_name := d.Get("instance").(string)

	db := &sqladmin.Database{
		Name:      database_name,
		Instance:  instance_name,
		Charset:   d.Get("charset").(string),
		Collation: d.Get("collation").(string),
	}

	mutexKV.Lock(instanceMutexKey(project, instance_name))
	defer mutexKV.Unlock(instanceMutexKey(project, instance_name))

	var op *sqladmin.Operation
	err = retryTime(func() error {
		op, err = config.clientSqlAdmin.Databases.Update(project, instance_name, database_name, db).Do()
		return err
	}, 5 /* minutes */)

	if err != nil {
		return fmt.Errorf("Error, failed to update "+
			"database %s in instance %s: %s", database_name,
			instance_name, err)
	}

	err = sqladminOperationWaitTime(config, op, project, "Update Database", int(d.Timeout(schema.TimeoutUpdate).Minutes()))

	if err != nil {
		return fmt.Errorf("Error, failure waiting for update of %s "+
			"into %s: %s", database_name, instance_name, err)
	}

	return resourceSqlDatabaseRead(d, meta)
}

func resourceSqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	database_name := d.Get("name").(string)
	instance_name := d.Get("instance").(string)

	mutexKV.Lock(instanceMutexKey(project, instance_name))
	defer mutexKV.Unlock(instanceMutexKey(project, instance_name))

	var op *sqladmin.Operation
	err = retryTimeDuration(func() error {
		op, err = config.clientSqlAdmin.Databases.Delete(project, instance_name, database_name).Do()
		return err
	}, d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return fmt.Errorf("Error, failed to delete"+
			"database %s in instance %s: %s", database_name,
			instance_name, err)
	}

	err = sqladminOperationWaitTime(config, op, project, "Delete Database", int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return fmt.Errorf("Error, failure waiting for deletion of %s "+
			"in %s: %s", database_name, instance_name, err)
	}

	return nil
}

func resourceSqlDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	parseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<name>[^/]+)",
		"instances/(?P<instance>[^/]+)/databases/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<instance>[^/]+)/(?P<name>[^/]+)",
		"(?P<instance>[^/]+)/(?P<name>[^/]+)",
		"(?P<instance>[^/]+):(?P<name>[^/]+)",
	}, d, config)

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{instance}}:{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
