// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/resourcemanager/Lien.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package resourcemanager

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceResourceManagerLien() *schema.Resource {
	return &schema.Resource{
		Create: resourceResourceManagerLienCreate,
		Read:   resourceResourceManagerLienRead,
		Delete: resourceResourceManagerLienDelete,

		Importer: &schema.ResourceImporter{
			State: resourceResourceManagerLienImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"origin": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `A stable, user-visible/meaningful string identifying the origin
of the Lien, intended to be inspected programmatically. Maximum length of
200 characters.`,
			},
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `A reference to the resource this Lien is attached to.
The server will validate the parent against those for which Liens are supported.
Since a variety of objects can have Liens against them, you must provide the type
prefix (e.g. "projects/my-project-name").`,
			},
			"reason": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Concise user-visible strings indicating why an action cannot be performed
on a resource. Maximum length of 200 characters.`,
			},
			"restrictions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Description: `The types of operations which should be blocked as a result of this Lien.
Each value should correspond to an IAM permission. The server will validate
the permissions against those for which Liens are supported.  An empty
list is meaningless and will be rejected.
e.g. ['resourcemanager.projects.delete']`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time of creation`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `A system-generated unique identifier for this Lien.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceResourceManagerLienCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	reasonProp, err := expandNestedResourceManagerLienReason(d.Get("reason"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("reason"); !tpgresource.IsEmptyValue(reflect.ValueOf(reasonProp)) && (ok || !reflect.DeepEqual(v, reasonProp)) {
		obj["reason"] = reasonProp
	}
	originProp, err := expandNestedResourceManagerLienOrigin(d.Get("origin"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("origin"); !tpgresource.IsEmptyValue(reflect.ValueOf(originProp)) && (ok || !reflect.DeepEqual(v, originProp)) {
		obj["origin"] = originProp
	}
	parentProp, err := expandNestedResourceManagerLienParent(d.Get("parent"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("parent"); !tpgresource.IsEmptyValue(reflect.ValueOf(parentProp)) && (ok || !reflect.DeepEqual(v, parentProp)) {
		obj["parent"] = parentProp
	}
	restrictionsProp, err := expandNestedResourceManagerLienRestrictions(d.Get("restrictions"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("restrictions"); !tpgresource.IsEmptyValue(reflect.ValueOf(restrictionsProp)) && (ok || !reflect.DeepEqual(v, restrictionsProp)) {
		obj["restrictions"] = restrictionsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ResourceManagerBasePath}}liens")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Lien: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating Lien: %s", err)
	}
	// Set computed resource properties from create API response so that they're available on the subsequent Read
	// call.
	res, err = resourceResourceManagerLienDecoder(d, meta, res)
	if err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	if res == nil {
		return fmt.Errorf("decoding response, could not find object")
	}
	if err := d.Set("name", flattenNestedResourceManagerLienName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// This resource is unusual - instead of returning an Operation from
	// Create, it returns the created object itself.  We don't parse
	// any of the values there, preferring to centralize that logic in
	// Read().  In this resource, Read is also unusual - it requires
	// us to know the server-side generated name of the object we're
	// trying to fetch, and the only way to know that is to capture
	// it here.  The following two lines do that.
	d.SetId(flattenNestedResourceManagerLienName(res["name"], d, config).(string))
	if err := d.Set("name", flattenNestedResourceManagerLienName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Lien %q: %#v", d.Id(), res)

	return resourceResourceManagerLienRead(d, meta)
}

func resourceResourceManagerLienRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ResourceManagerBasePath}}liens?parent={{parent}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ResourceManagerLien %q", d.Id()))
	}

	res, err = flattenNestedResourceManagerLien(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Object isn't there any more - remove it from the state.
		log.Printf("[DEBUG] Removing ResourceManagerLien because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	res, err = resourceResourceManagerLienDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing ResourceManagerLien because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("name", flattenNestedResourceManagerLienName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Lien: %s", err)
	}
	if err := d.Set("reason", flattenNestedResourceManagerLienReason(res["reason"], d, config)); err != nil {
		return fmt.Errorf("Error reading Lien: %s", err)
	}
	if err := d.Set("origin", flattenNestedResourceManagerLienOrigin(res["origin"], d, config)); err != nil {
		return fmt.Errorf("Error reading Lien: %s", err)
	}
	if err := d.Set("create_time", flattenNestedResourceManagerLienCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Lien: %s", err)
	}
	if err := d.Set("parent", flattenNestedResourceManagerLienParent(res["parent"], d, config)); err != nil {
		return fmt.Errorf("Error reading Lien: %s", err)
	}
	if err := d.Set("restrictions", flattenNestedResourceManagerLienRestrictions(res["restrictions"], d, config)); err != nil {
		return fmt.Errorf("Error reading Lien: %s", err)
	}

	return nil
}

func resourceResourceManagerLienDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{ResourceManagerBasePath}}liens?parent={{parent}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	// log the old URL to make the ineffassign linter happy
	// in theory, we should find a way to disable the default URL and not construct
	// both, but that's a problem for another day. Today, we cheat.
	log.Printf("[DEBUG] replacing URL %q with a custom delete URL", url)
	url, err = tpgresource.ReplaceVars(d, config, "{{ResourceManagerBasePath}}liens/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Lien %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Lien")
	}

	log.Printf("[DEBUG] Finished deleting Lien %q: %#v", d.Id(), res)
	return nil
}

func resourceResourceManagerLienImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^(?P<parent>[^/]+)/(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	parent, err := tpgresource.ReplaceVars(d, config, "projects/{{parent}}")
	if err != nil {
		return nil, err
	}
	if err := d.Set("parent", parent); err != nil {
		return nil, fmt.Errorf("Error setting parent: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenNestedResourceManagerLienName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func flattenNestedResourceManagerLienReason(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedResourceManagerLienOrigin(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedResourceManagerLienCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedResourceManagerLienParent(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedResourceManagerLienRestrictions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNestedResourceManagerLienReason(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedResourceManagerLienOrigin(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedResourceManagerLienParent(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedResourceManagerLienRestrictions(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func flattenNestedResourceManagerLien(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	var v interface{}
	var ok bool

	v, ok = res["liens"]
	if !ok || v == nil {
		return nil, nil
	}

	switch v.(type) {
	case []interface{}:
		break
	case map[string]interface{}:
		// Construct list out of single nested resource
		v = []interface{}{v}
	default:
		return nil, fmt.Errorf("expected list or map for value liens. Actual value: %v", v)
	}

	_, item, err := resourceResourceManagerLienFindNestedObjectInList(d, meta, v.([]interface{}))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func resourceResourceManagerLienFindNestedObjectInList(d *schema.ResourceData, meta interface{}, items []interface{}) (index int, item map[string]interface{}, err error) {
	expectedName := d.Get("name")
	expectedFlattenedName := flattenNestedResourceManagerLienName(expectedName, d, meta.(*transport_tpg.Config))

	// Search list for this resource.
	for idx, itemRaw := range items {
		if itemRaw == nil {
			continue
		}
		item := itemRaw.(map[string]interface{})

		// Decode list item before comparing.
		item, err := resourceResourceManagerLienDecoder(d, meta, item)
		if err != nil {
			return -1, nil, err
		}

		itemName := flattenNestedResourceManagerLienName(item["name"], d, meta.(*transport_tpg.Config))
		// IsEmptyValue check so that if one is nil and the other is "", that's considered a match
		if !(tpgresource.IsEmptyValue(reflect.ValueOf(itemName)) && tpgresource.IsEmptyValue(reflect.ValueOf(expectedFlattenedName))) && !reflect.DeepEqual(itemName, expectedFlattenedName) {
			log.Printf("[DEBUG] Skipping item with name= %#v, looking for %#v)", itemName, expectedFlattenedName)
			continue
		}
		log.Printf("[DEBUG] Found item for resource %q: %#v)", d.Id(), item)
		return idx, item, nil
	}
	return -1, nil, nil
}
func resourceResourceManagerLienDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	// The problem we're trying to solve here is that this property is a Project,
	// and there are a lot of ways to specify a Project, including the ID vs
	// Number, which is something that we can't address in a diffsuppress.
	// Since we can't enforce a particular method of entering the project,
	// we're just going to have to use whatever the user entered, whether
	// it's project/projectName, project/12345, projectName, or 12345.
	// The normal behavior of this method would be 'return res' - and that's
	// what we'll fall back to if any of our conditions aren't met.  Those
	// conditions are:
	// 1) if the new or old values contain '/', the prefix of that is 'projects'.
	// 2) if either is non-numeric, a project with that ID exists.
	// 3) the project IDs represented by both the new and old values are the same.
	config := meta.(*transport_tpg.Config)

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	new := res["parent"].(string)
	old := d.Get("parent").(string)
	if strings.HasPrefix(new, "projects/") {
		new = strings.Split(new, "/")[1]
	}
	if strings.HasPrefix(old, "projects/") {
		old = strings.Split(old, "/")[1]
	}
	log.Printf("[DEBUG] Trying to figure out whether to use %s or %s", old, new)
	// If there's still a '/' in there, the value must not be a project ID.
	if strings.Contains(old, "/") || strings.Contains(new, "/") {
		return res, nil
	}
	// If 'old' isn't entirely numeric, let's assume it's a project ID.
	// If it's a project ID
	var oldProjId int64
	var newProjId int64
	if oldVal, err := tpgresource.StringToFixed64(old); err == nil {
		log.Printf("[DEBUG] The old value was a real number: %d", oldVal)
		oldProjId = oldVal
	} else {
		pOld, err := config.NewResourceManagerClient(userAgent).Projects.Get(old).Do()
		if err != nil {
			return res, nil
		}
		oldProjId = pOld.ProjectNumber
	}
	if newVal, err := tpgresource.StringToFixed64(new); err == nil {
		log.Printf("[DEBUG] The new value was a real number: %d", newVal)
		newProjId = newVal
	} else {
		pNew, err := config.NewResourceManagerClient(userAgent).Projects.Get(new).Do()
		if err != nil {
			return res, nil
		}
		newProjId = pNew.ProjectNumber
	}
	if newProjId == oldProjId {
		res["parent"] = d.Get("parent")
	}
	return res, nil
}
