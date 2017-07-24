package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/sqladmin/v1beta4"
)

func resourceSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlDatabaseCreate,
		Read:   resourceSqlDatabaseRead,
		Update: resourceSqlDatabaseUpdate,
		Delete: resourceSqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"charset": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
			},

			"collation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8_general_ci",
			},
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
	op, err := config.clientSqlAdmin.Databases.Insert(project, instance_name,
		db).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to insert "+
			"database %s into instance %s: %s", database_name,
			instance_name, err)
	}

	err = sqladminOperationWait(config, op, "Insert Database")

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

	db, err := config.clientSqlAdmin.Databases.Get(project, instance_name,
		database_name).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL Database %q in instance %q", database_name, instance_name))
	}

	d.Set("instance", db.Instance)
	d.Set("name", db.Name)
	d.Set("self_link", db.SelfLink)
	d.SetId(instance_name + ":" + database_name)
	d.Set("charset", db.Charset)
	d.Set("collation", db.Collation)

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
	op, err := config.clientSqlAdmin.Databases.Update(project, instance_name, database_name,
		db).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to update "+
			"database %s in instance %s: %s", database_name,
			instance_name, err)
	}

	err = sqladminOperationWait(config, op, "Update Database")

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
	op, err := config.clientSqlAdmin.Databases.Delete(project, instance_name,
		database_name).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to delete"+
			"database %s in instance %s: %s", database_name,
			instance_name, err)
	}

	err = sqladminOperationWait(config, op, "Delete Database")

	if err != nil {
		return fmt.Errorf("Error, failure waiting for deletion of %s "+
			"in %s: %s", database_name, instance_name, err)
	}

	return nil
}
