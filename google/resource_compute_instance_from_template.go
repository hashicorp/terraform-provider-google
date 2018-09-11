package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	strcase "github.com/stoewer/go-strcase"
)

func resourceComputeInstanceFromTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceFromTemplateCreate,
		Read:   resourceComputeInstanceRead,
		Update: resourceComputeInstanceUpdate,
		Delete: resourceComputeInstanceDelete,

		// Import doesn't really make sense, because you could just import
		// as a google_compute_instance.

		Timeouts: resourceComputeInstance().Timeouts,

		Schema:        computeInstanceFromTemplateSchema(),
		CustomizeDiff: resourceComputeInstance().CustomizeDiff,
	}
}

func computeInstanceFromTemplateSchema() map[string]*schema.Schema {
	s := resourceComputeInstance().Schema

	for _, field := range []string{"boot_disk", "machine_type", "network_interface"} {
		// The user can set these fields as an override, but doesn't need to -
		// the template values will be used if they're unset.
		s[field].Required = false
		s[field].Optional = true
	}

	// Remove deprecated/removed fields that are never d.Set. We can't
	// programatically remove all of them, because some of them still have d.Set
	// calls.
	for _, field := range []string{"create_timeout", "disk", "network"} {
		delete(s, field)
	}

	recurseOnSchema(s, func(field *schema.Schema) {
		// We don't want to accidentally use default values to override the instance
		// template, so remove defaults.
		field.Default = nil

		// Make non-required fields computed since they'll be set by the template.
		// Leave deprecated and removed fields alone because we don't set them.
		if !field.Required && !(field.Deprecated != "" || field.Removed != "") {
			field.Computed = true
		}
	})

	s["source_instance_template"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}

	return s
}

func recurseOnSchema(s map[string]*schema.Schema, f func(*schema.Schema)) {
	for _, field := range s {
		f(field)
		if e := field.Elem; e != nil {
			if r, ok := e.(*schema.Resource); ok {
				recurseOnSchema(r.Schema, f)
			}
		}
	}
}

func resourceComputeInstanceFromTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the zone
	z, err := getZone(d, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Loading zone: %s", z)
	zone, err := config.clientCompute.Zones.Get(project, z).Do()
	if err != nil {
		return fmt.Errorf("Error loading zone '%s': %s", z, err)
	}

	instance, err := expandComputeInstance(project, zone, d, config)
	if err != nil {
		return err
	}

	// Force send all top-level fields in case they're overridden to zero values.
	// TODO: consider doing so for nested fields as well.
	for f, s := range computeInstanceFromTemplateSchema() {
		// It seems that GetOkExists always returns true for sets.
		// TODO: confirm this and file issue against Terraform core.
		// In the meantime, don't force send sets.
		if s.Type == schema.TypeSet {
			continue
		}

		if _, exists := d.GetOkExists(f); exists {
			// Assume for now that all fields are exact snake_case versions of the API fields.
			// This won't necessarily always be true, but it serves as a good approximation and
			// can be adjusted later as we discover issues.
			instance.ForceSendFields = append(instance.ForceSendFields, strcase.UpperCamelCase(f))
		}
	}

	tpl, err := ParseInstanceTemplateFieldValue(d.Get("source_instance_template").(string), d, config)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Requesting instance creation")
	op, err := config.clientComputeBeta.Instances.Insert(project, zone.Name, instance).SourceInstanceTemplate(tpl.RelativeLink()).Do()
	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	// Store the ID now
	d.SetId(instance.Name)

	// Wait for the operation to complete
	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), "instance to create")
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	return resourceComputeInstanceRead(d, meta)
}
