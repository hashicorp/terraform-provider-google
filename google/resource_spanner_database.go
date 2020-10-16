// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customizeDiff func for additional checks on google_spanner_database properties:
func resourceSpannerDBDdlCustomDiffFunc(diff TerraformResourceDiff) error {
	old, new := diff.GetChange("ddl")
	oldDdls := old.([]interface{})
	newDdls := new.([]interface{})
	var err error

	if len(newDdls) < len(oldDdls) {
		err = diff.ForceNew("ddl")
		if err != nil {
			return fmt.Errorf("ForceNew failed for ddl, old - %v and new - %v", oldDdls, newDdls)
		}
		return nil
	}

	for i := range oldDdls {
		if newDdls[i].(string) != oldDdls[i].(string) {
			err = diff.ForceNew("ddl")
			if err != nil {
				return fmt.Errorf("ForceNew failed for ddl, old - %v and new - %v", oldDdls, newDdls)
			}
			return nil
		}
	}
	return nil
}

func resourceSpannerDBDdlCustomDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return resourceSpannerDBDdlCustomDiffFunc(diff)
}

func resourceSpannerDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpannerDatabaseCreate,
		Read:   resourceSpannerDatabaseRead,
		Update: resourceSpannerDatabaseUpdate,
		Delete: resourceSpannerDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSpannerDatabaseImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		CustomizeDiff: resourceSpannerDBDdlCustomDiff,

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The instance to create the database on.`,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp(`^[a-z][a-z0-9_\-]*[a-z0-9]$`),
				Description: `A unique identifier for the database, which cannot be changed after
the instance is created. Values are of the form [a-z][-a-z0-9]*[a-z0-9].`,
			},
			"ddl": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `An optional list of DDL statements to run inside the newly created
database. Statements can create tables, indexes, etc. These statements
execute atomically with the creation of the database: if there is an
error in any statement, the database is not created.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `An explanation of the status of the database.`,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSpannerDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandSpannerDatabaseName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	extraStatementsProp, err := expandSpannerDatabaseDdl(d.Get("ddl"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ddl"); !isEmptyValue(reflect.ValueOf(extraStatementsProp)) && (ok || !reflect.DeepEqual(v, extraStatementsProp)) {
		obj["extraStatements"] = extraStatementsProp
	}
	instanceProp, err := expandSpannerDatabaseInstance(d.Get("instance"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("instance"); !isEmptyValue(reflect.ValueOf(instanceProp)) && (ok || !reflect.DeepEqual(v, instanceProp)) {
		obj["instance"] = instanceProp
	}

	obj, err = resourceSpannerDatabaseEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances/{{instance}}/databases")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Database: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Database: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{instance}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = spannerOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Creating Database", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Database: %s", err)
	}

	opRes, err = resourceSpannerDatabaseDecoder(d, meta, opRes)
	if err != nil {
		return fmt.Errorf("Error decoding response from operation: %s", err)
	}
	if opRes == nil {
		return fmt.Errorf("Error decoding response from operation, could not find object")
	}

	if err := d.Set("name", flattenSpannerDatabaseName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = replaceVars(d, config, "{{instance}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Database %q: %#v", d.Id(), res)

	return resourceSpannerDatabaseRead(d, meta)
}

func resourceSpannerDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances/{{instance}}/databases/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SpannerDatabase %q", d.Id()))
	}

	res, err = resourceSpannerDatabaseDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing SpannerDatabase because it no longer exists.")
		d.SetId("")
		return nil
	}

	// Explicitly set virtual fields to default values if unset
	if _, ok := d.GetOkExists("deletion_protection"); !ok {
		if err := d.Set("deletion_protection", true); err != nil {
			return fmt.Errorf("Error setting deletion_protection: %s", err)
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}

	if err := d.Set("name", flattenSpannerDatabaseName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}
	if err := d.Set("state", flattenSpannerDatabaseState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}
	if err := d.Set("instance", flattenSpannerDatabaseInstance(res["instance"], d, config)); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}

	return nil
}

func resourceSpannerDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	config.userAgent = userAgent

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	d.Partial(true)

	if d.HasChange("ddl") {
		obj := make(map[string]interface{})

		extraStatementsProp, err := expandSpannerDatabaseDdl(d.Get("ddl"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("ddl"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, extraStatementsProp)) {
			obj["extraStatements"] = extraStatementsProp
		}

		obj, err = resourceSpannerDatabaseUpdateEncoder(d, meta, obj)
		if err != nil {
			return err
		}

		url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances/{{instance}}/databases/{{name}}/ddl")
		if err != nil {
			return err
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating Database %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Database %q: %#v", d.Id(), res)
		}

		err = spannerOperationWaitTime(
			config, res, project, "Updating Database", userAgent,
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceSpannerDatabaseRead(d, meta)
}

func resourceSpannerDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	config.userAgent = userAgent

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances/{{instance}}/databases/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("cannot destroy instance without setting deletion_protection=false and running `terraform apply`")
	}
	log.Printf("[DEBUG] Deleting Database %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Database")
	}

	err = spannerOperationWaitTime(
		config, res, project, "Deleting Database", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Database %q: %#v", d.Id(), res)
	return nil
}

func resourceSpannerDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<name>[^/]+)",
		"instances/(?P<instance>[^/]+)/databases/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<instance>[^/]+)/(?P<name>[^/]+)",
		"(?P<instance>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{instance}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Explicitly set virtual fields to default values on import
	if err := d.Set("deletion_protection", true); err != nil {
		return nil, fmt.Errorf("Error setting deletion_protection: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenSpannerDatabaseName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return NameFromSelfLinkStateFunc(v)
}

func flattenSpannerDatabaseState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenSpannerDatabaseInstance(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func expandSpannerDatabaseName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerDatabaseDdl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerDatabaseInstance(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("instances", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for instance: %s", err)
	}
	return f.RelativeLink(), nil
}

func resourceSpannerDatabaseEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	obj["createStatement"] = fmt.Sprintf("CREATE DATABASE `%s`", obj["name"])
	delete(obj, "name")
	delete(obj, "instance")
	return obj, nil
}

func resourceSpannerDatabaseUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	old, new := d.GetChange("ddl")
	oldDdls := old.([]interface{})
	newDdls := new.([]interface{})
	updateDdls := []string{}

	//Only new ddl statments to be add to update call
	for i := len(oldDdls); i < len(newDdls); i++ {
		updateDdls = append(updateDdls, newDdls[i].(string))
	}

	obj["statements"] = strings.Join(updateDdls, ",")
	delete(obj, "name")
	delete(obj, "instance")
	delete(obj, "extraStatements")
	return obj, nil
}

func resourceSpannerDatabaseDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*Config)
	d.SetId(res["name"].(string))
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}
	res["project"] = d.Get("project").(string)
	res["instance"] = d.Get("instance").(string)
	res["name"] = d.Get("name").(string)
	id, err := replaceVars(d, config, "{{instance}}/{{name}}")
	if err != nil {
		return nil, err
	}
	d.SetId(id)
	return res, nil
}
