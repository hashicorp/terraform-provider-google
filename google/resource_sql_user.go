package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/sqladmin/v1beta4"
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

		SchemaVersion: 1,
		MigrateState:  resourceSqlUserMigrateState,

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
	op, err := config.clientSqlAdmin.Users.Insert(project, instance,
		user).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to insert "+
			"user %s into instance %s: %s", name, instance, err)
	}

	// This will include a double-slash (//) for postgres instances,
	// for which user.Host is an empty string.  That's okay.
	d.SetId(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))

	err = sqladminOperationWait(config, op, project, "Insert User")

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
		// The second part of this conditional is irrelevant for postgres instances because
		// host and currentUser.Host will always both be empty.
		if currentUser.Name == name && currentUser.Host == host {
			user = currentUser
			break
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
		host := d.Get("host").(string)
		password := d.Get("password").(string)

		user := &sqladmin.User{
			Name:     name,
			Instance: instance,
			Password: password,
			Host:     host,
		}

		mutexKV.Lock(instanceMutexKey(project, instance))
		defer mutexKV.Unlock(instanceMutexKey(project, instance))
		op, err := config.clientSqlAdmin.Users.Update(project, instance, name,
			user).Do()

		if err != nil {
			return fmt.Errorf("Error, failed to update"+
				"user %s into user %s: %s", name, instance, err)
		}

		err = sqladminOperationWait(config, op, project, "Insert User")

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
	instance := d.Get("instance").(string)
	host := d.Get("host").(string)

	mutexKV.Lock(instanceMutexKey(project, instance))
	defer mutexKV.Unlock(instanceMutexKey(project, instance))
	op, err := config.clientSqlAdmin.Users.Delete(project, instance, host, name).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to delete"+
			"user %s in instance %s: %s", name,
			instance, err)
	}

	err = sqladminOperationWait(config, op, project, "Delete User")

	if err != nil {
		return fmt.Errorf("Error, failure waiting for deletion of %s "+
			"in %s: %s", name, instance, err)
	}

	return nil
}

func resourceSqlUserImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) == 2 {
		d.Set("instance", parts[0])
		d.Set("name", parts[1])
	} else if len(parts) == 3 {
		d.Set("instance", parts[0])
		d.Set("host", parts[1])
		d.Set("name", parts[2])
	} else {
		return nil, fmt.Errorf("Invalid specifier. Expecting {instance}/{name} for postgres instance and {instance}/{host}/{name} for MySQL instance")
	}

	return []*schema.ResourceData{d}, nil
}
