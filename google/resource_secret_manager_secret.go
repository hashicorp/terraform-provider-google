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
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecretManagerSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretManagerSecretCreate,
		Read:   resourceSecretManagerSecretRead,
		Update: resourceSecretManagerSecretUpdate,
		Delete: resourceSecretManagerSecretDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSecretManagerSecretImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"replication": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Description: `The replication policy of the secret data attached to the Secret. It cannot be changed
after the Secret has been created.`,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"automatic": {
							Type:         schema.TypeBool,
							Optional:     true,
							Description:  `The Secret will automatically be replicated without any restrictions.`,
							ExactlyOneOf: []string{"replication.0.automatic", "replication.0.user_managed"},
						},
						"user_managed": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The Secret will automatically be replicated without any restrictions.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"replicas": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `The list of Replicas for this Secret. Cannot be empty.`,
										MinItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"location": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The canonical IDs of the location to replicate data. For example: "us-east1".`,
												},
											},
										},
									},
								},
							},
							ExactlyOneOf: []string{"replication.0.automatic", "replication.0.user_managed"},
						},
					},
				},
			},
			"secret_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `This must be unique within the project.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `The labels assigned to this Secret.

Label keys must be between 1 and 63 characters long, have a UTF-8 encoding of maximum 128 bytes,
and must conform to the following PCRE regular expression: [\p{Ll}\p{Lo}][\p{Ll}\p{Lo}\p{N}_-]{0,62}

Label values must be between 0 and 63 characters long, have a UTF-8 encoding of maximum 128 bytes,
and must conform to the following PCRE regular expression: [\p{Ll}\p{Lo}\p{N}_-]{0,63}

No more than 64 labels can be assigned to a given resource.

An object containing a list of "key": value pairs. Example:
{ "name": "wrench", "mass": "1.3kg", "count": "3" }.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time at which the Secret was created.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name of the Secret. Format:
'projects/{{project}}/secrets/{{secret_id}}'`,
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

func resourceSecretManagerSecretCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	labelsProp, err := expandSecretManagerSecretLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	replicationProp, err := expandSecretManagerSecretReplication(d.Get("replication"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("replication"); !isEmptyValue(reflect.ValueOf(replicationProp)) && (ok || !reflect.DeepEqual(v, replicationProp)) {
		obj["replication"] = replicationProp
	}

	url, err := replaceVars(d, config, "{{SecretManagerBasePath}}projects/{{project}}/secrets?secretId={{secret_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Secret: %#v", obj)
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

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Secret: %s", err)
	}
	if err := d.Set("name", flattenSecretManagerSecretName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/secrets/{{secret_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Secret %q: %#v", d.Id(), res)

	return resourceSecretManagerSecretRead(d, meta)
}

func resourceSecretManagerSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{SecretManagerBasePath}}projects/{{project}}/secrets/{{secret_id}}")
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

	res, err := sendRequest(config, "GET", billingProject, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SecretManagerSecret %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Secret: %s", err)
	}

	if err := d.Set("name", flattenSecretManagerSecretName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Secret: %s", err)
	}
	if err := d.Set("create_time", flattenSecretManagerSecretCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Secret: %s", err)
	}
	if err := d.Set("labels", flattenSecretManagerSecretLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Secret: %s", err)
	}
	if err := d.Set("replication", flattenSecretManagerSecretReplication(res["replication"], d, config)); err != nil {
		return fmt.Errorf("Error reading Secret: %s", err)
	}

	return nil
}

func resourceSecretManagerSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	obj := make(map[string]interface{})
	labelsProp, err := expandSecretManagerSecretLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	url, err := replaceVars(d, config, "{{SecretManagerBasePath}}projects/{{project}}/secrets/{{secret_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Secret %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("labels") {
		updateMask = append(updateMask, "labels")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating Secret %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Secret %q: %#v", d.Id(), res)
	}

	return resourceSecretManagerSecretRead(d, meta)
}

func resourceSecretManagerSecretDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{SecretManagerBasePath}}projects/{{project}}/secrets/{{secret_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Secret %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Secret")
	}

	log.Printf("[DEBUG] Finished deleting Secret %q: %#v", d.Id(), res)
	return nil
}

func resourceSecretManagerSecretImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/secrets/(?P<secret_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<secret_id>[^/]+)",
		"(?P<secret_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/secrets/{{secret_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenSecretManagerSecretName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenSecretManagerSecretCreateTime(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenSecretManagerSecretLabels(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenSecretManagerSecretReplication(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["automatic"] =
		flattenSecretManagerSecretReplicationAutomatic(original["automatic"], d, config)
	transformed["user_managed"] =
		flattenSecretManagerSecretReplicationUserManaged(original["userManaged"], d, config)
	return []interface{}{transformed}
}
func flattenSecretManagerSecretReplicationAutomatic(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v != nil
}

func flattenSecretManagerSecretReplicationUserManaged(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["replicas"] =
		flattenSecretManagerSecretReplicationUserManagedReplicas(original["replicas"], d, config)
	return []interface{}{transformed}
}
func flattenSecretManagerSecretReplicationUserManagedReplicas(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"location": flattenSecretManagerSecretReplicationUserManagedReplicasLocation(original["location"], d, config),
		})
	}
	return transformed
}
func flattenSecretManagerSecretReplicationUserManagedReplicasLocation(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandSecretManagerSecretLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandSecretManagerSecretReplication(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAutomatic, err := expandSecretManagerSecretReplicationAutomatic(original["automatic"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAutomatic); val.IsValid() && !isEmptyValue(val) {
		transformed["automatic"] = transformedAutomatic
	}

	transformedUserManaged, err := expandSecretManagerSecretReplicationUserManaged(original["user_managed"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUserManaged); val.IsValid() && !isEmptyValue(val) {
		transformed["userManaged"] = transformedUserManaged
	}

	return transformed, nil
}

func expandSecretManagerSecretReplicationAutomatic(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	if v == nil || !v.(bool) {
		return nil, nil
	}

	return struct{}{}, nil
}

func expandSecretManagerSecretReplicationUserManaged(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedReplicas, err := expandSecretManagerSecretReplicationUserManagedReplicas(original["replicas"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedReplicas); val.IsValid() && !isEmptyValue(val) {
		transformed["replicas"] = transformedReplicas
	}

	return transformed, nil
}

func expandSecretManagerSecretReplicationUserManagedReplicas(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedLocation, err := expandSecretManagerSecretReplicationUserManagedReplicasLocation(original["location"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedLocation); val.IsValid() && !isEmptyValue(val) {
			transformed["location"] = transformedLocation
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandSecretManagerSecretReplicationUserManagedReplicasLocation(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
