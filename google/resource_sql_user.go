package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func resourceSqlUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlUserCreate,
		Read:   resourceSqlUserRead,
		Update: resourceSqlUserUpdate,
		Delete: resourceSqlUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSqlUserImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		MigrateState:  resourceSqlUserMigrateState,

		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The host the user can connect from. This is only supported for MySQL instances. Don't set this field for PostgreSQL instances. Can be an IP address. Changing this forces a new resource to be created.`,
			},

			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the Cloud SQL instance. Changing this forces a new resource to be created.`,
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the user. Changing this forces a new resource to be created.`,
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The password for the user. Can be updated. For Postgres instances this is a Required field.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},
	}
}

func resourceSqlUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	instance := d.Get("instance").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)

	user := &sqladmin.User{
		Name:     name,
		Instance: instance,
		Password: password,
		Host:     host,
	}

	mutexKV.Lock(instanceMutexKey(project, instance))
	defer mutexKV.Unlock(instanceMutexKey(project, instance))
	var op *sqladmin.Operation
	insertFunc := func() error {
		op, err = config.clientSqlAdmin.Users.Insert(project, instance,
			user).Do()
		return err
	}
	err = retryTimeDuration(insertFunc, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("Error, failed to insert "+
			"user %s into instance %s: %s", name, instance, err)
	}

	// This will include a double-slash (//) for postgres instances,
	// for which user.Host is an empty string.  That's okay.
	d.SetId(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))

	err = sqlAdminOperationWaitTime(config, op, project, "Insert User", d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("Error, failure waiting for insertion of %s "+
			"into %s: %s", name, instance, err)
	}

	return resourceSqlUserRead(d, meta)
}

func resourceSqlUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	name := d.Get("name").(string)
	host := d.Get("host").(string)

	var users *sqladmin.UsersListResponse
	err = nil
	err = retryTime(func() error {
		users, err = config.clientSqlAdmin.Users.List(project, instance).Do()
		return err
	}, 5)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL User %q in instance %q", name, instance))
	}

	var user *sqladmin.User
	for _, currentUser := range users.Items {
		if currentUser.Name == name {
			// Host can only be empty for postgres instances,
			// so don't compare the host if the API host is empty.
			if currentUser.Host == "" || currentUser.Host == host {
				user = currentUser
				break
			}
		}
	}

	if user == nil {
		log.Printf("[WARN] Removing SQL User %q because it's gone", d.Get("name").(string))
		d.SetId("")

		return nil
	}

	d.Set("host", user.Host)
	d.Set("instance", user.Instance)
	d.Set("name", user.Name)
	d.Set("project", project)
	d.SetId(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))
	return nil
}

func resourceSqlUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("password") {
		project, err := getProject(d, config)
		if err != nil {
			return err
		}

		name := d.Get("name").(string)
		instance := d.Get("instance").(string)
		password := d.Get("password").(string)
		host := d.Get("host").(string)

		user := &sqladmin.User{
			Name:     name,
			Instance: instance,
			Password: password,
		}

		mutexKV.Lock(instanceMutexKey(project, instance))
		defer mutexKV.Unlock(instanceMutexKey(project, instance))
		var op *sqladmin.Operation
		updateFunc := func() error {
			op, err = config.clientSqlAdmin.Users.Update(project, instance, user).Host(host).Name(name).Do()
			return err
		}
		err = retryTimeDuration(updateFunc, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return fmt.Errorf("Error, failed to update"+
				"user %s into user %s: %s", name, instance, err)
		}

		err = sqlAdminOperationWaitTime(config, op, project, "Insert User", d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return fmt.Errorf("Error, failure waiting for update of %s "+
				"in %s: %s", name, instance, err)
		}

		return resourceSqlUserRead(d, meta)
	}

	return nil
}

func resourceSqlUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	host := d.Get("host").(string)
	instance := d.Get("instance").(string)

	mutexKV.Lock(instanceMutexKey(project, instance))
	defer mutexKV.Unlock(instanceMutexKey(project, instance))

	var op *sqladmin.Operation
	err = retryTimeDuration(func() error {
		op, err = config.clientSqlAdmin.Users.Delete(project, instance).Host(host).Name(name).Do()
		if err != nil {
			return err
		}

		if err := sqlAdminOperationWaitTime(config, op, project, "Delete User", d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
		return nil
	}, d.Timeout(schema.TimeoutDelete), isSqlOperationInProgressError, isSqlInternalError)

	if err != nil {
		return fmt.Errorf("Error, failed to delete"+
			"user %s in instance %s: %s", name,
			instance, err)
	}

	return nil
}

func resourceSqlUserImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) == 3 {
		d.Set("project", parts[0])
		d.Set("instance", parts[1])
		d.Set("name", parts[2])
	} else if len(parts) == 4 {
		d.Set("project", parts[0])
		d.Set("instance", parts[1])
		d.Set("host", parts[2])
		d.Set("name", parts[3])
	} else {
		return nil, fmt.Errorf("Invalid specifier. Expecting {project}/{instance}/{name} for postgres instance and {project}/{instance}/{host}/{name} for MySQL instance")
	}

	return []*schema.ResourceData{d}, nil
}
