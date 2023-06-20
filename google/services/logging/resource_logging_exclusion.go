// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/logging/v2"
)

var LoggingExclusionBaseSchema = map[string]*schema.Schema{
	"filter": {
		Type:        schema.TypeString,
		Required:    true,
		Description: `The filter to apply when excluding logs. Only log entries that match the filter are excluded.`,
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The name of the logging exclusion.`,
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: `A human-readable description.`,
	},
	"disabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: `Whether this exclusion rule should be disabled or not. This defaults to false.`,
	},
}

func ResourceLoggingExclusion(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceLoggingExclusionUpdaterFunc, resourceIdParser tpgiamresource.ResourceIdParserFunc) *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingExclusionCreate(newUpdaterFunc),
		Read:   resourceLoggingExclusionRead(newUpdaterFunc),
		Update: resourceLoggingExclusionUpdate(newUpdaterFunc),
		Delete: resourceLoggingExclusionDelete(newUpdaterFunc),

		Importer: &schema.ResourceImporter{
			State: resourceLoggingExclusionImportState(resourceIdParser),
		},

		Schema:        tpgresource.MergeSchemas(LoggingExclusionBaseSchema, parentSpecificSchema),
		UseJSONNumber: true,
	}
}

func resourceLoggingExclusionCreate(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		id, exclusion := expandResourceLoggingExclusion(d, updater.GetResourceType(), updater.GetResourceId())

		// Logging exclusions don't seem to be able to be mutated in parallel, see
		// https://github.com/hashicorp/terraform-provider-google/issues/4796
		transport_tpg.MutexStore.Lock(id.parent())
		defer transport_tpg.MutexStore.Unlock(id.parent())

		err = updater.CreateLoggingExclusion(id.parent(), exclusion)
		if err != nil {
			return err
		}

		d.SetId(id.canonicalId())

		return resourceLoggingExclusionRead(newUpdaterFunc)(d, meta)
	}
}

func resourceLoggingExclusionRead(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		exclusion, err := updater.ReadLoggingExclusion(d.Id())

		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Logging Exclusion %s", d.Get("name").(string)))
		}

		if err := flattenResourceLoggingExclusion(d, exclusion); err != nil {
			return err
		}

		if updater.GetResourceType() == "projects" {
			if err := d.Set("project", updater.GetResourceId()); err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
		}

		return nil
	}
}

func resourceLoggingExclusionUpdate(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		id, _ := expandResourceLoggingExclusion(d, updater.GetResourceType(), updater.GetResourceId())
		exclusion, updateMask := expandResourceLoggingExclusionForUpdate(d)

		// Logging exclusions don't seem to be able to be mutated in parallel, see
		// https://github.com/hashicorp/terraform-provider-google/issues/4796
		transport_tpg.MutexStore.Lock(id.parent())
		defer transport_tpg.MutexStore.Unlock(id.parent())

		err = updater.UpdateLoggingExclusion(d.Id(), exclusion, updateMask)
		if err != nil {
			return err
		}

		return resourceLoggingExclusionRead(newUpdaterFunc)(d, meta)
	}
}

func resourceLoggingExclusionDelete(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		id, _ := expandResourceLoggingExclusion(d, updater.GetResourceType(), updater.GetResourceId())
		// Logging exclusions don't seem to be able to be mutated in parallel, see
		// https://github.com/hashicorp/terraform-provider-google/issues/4796
		transport_tpg.MutexStore.Lock(id.parent())
		defer transport_tpg.MutexStore.Unlock(id.parent())

		err = updater.DeleteLoggingExclusion(d.Id())
		if err != nil {
			return err
		}

		d.SetId("")
		return nil
	}
}

func resourceLoggingExclusionImportState(resourceIdParser tpgiamresource.ResourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		config := meta.(*transport_tpg.Config)
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}
		return []*schema.ResourceData{d}, nil
	}
}

func expandResourceLoggingExclusion(d *schema.ResourceData, resourceType, ResourceId string) (LoggingExclusionId, *logging.LogExclusion) {
	id := LoggingExclusionId{
		resourceType: resourceType,
		ResourceId:   ResourceId,
		name:         d.Get("name").(string),
	}

	exclusion := logging.LogExclusion{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Filter:      d.Get("filter").(string),
		Disabled:    d.Get("disabled").(bool),
	}
	return id, &exclusion
}

func flattenResourceLoggingExclusion(d *schema.ResourceData, exclusion *logging.LogExclusion) error {
	if err := d.Set("name", exclusion.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", exclusion.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("filter", exclusion.Filter); err != nil {
		return fmt.Errorf("Error setting filter: %s", err)
	}
	if err := d.Set("disabled", exclusion.Disabled); err != nil {
		return fmt.Errorf("Error setting disabled: %s", err)
	}

	return nil
}

func expandResourceLoggingExclusionForUpdate(d *schema.ResourceData) (*logging.LogExclusion, string) {
	// Can update description/filter/disabled right now.
	exclusion := logging.LogExclusion{}

	var updateMaskArr []string

	if d.HasChange("description") {
		exclusion.Description = d.Get("description").(string)
		exclusion.ForceSendFields = append(exclusion.ForceSendFields, "Description")
		updateMaskArr = append(updateMaskArr, "description")
	}

	if d.HasChange("filter") {
		exclusion.Filter = d.Get("filter").(string)
		exclusion.ForceSendFields = append(exclusion.ForceSendFields, "Filter")
		updateMaskArr = append(updateMaskArr, "filter")
	}

	if d.HasChange("disabled") {
		exclusion.Disabled = d.Get("disabled").(bool)
		exclusion.ForceSendFields = append(exclusion.ForceSendFields, "Disabled")
		updateMaskArr = append(updateMaskArr, "disabled")
	}

	updateMask := strings.Join(updateMaskArr, ",")
	return &exclusion, updateMask
}

// The ResourceLoggingExclusionUpdater interface is implemented for each GCP
// resource supporting log exclusions.
//
// Implementations should keep track of the resource identifier.
type ResourceLoggingExclusionUpdater interface {
	CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error
	ReadLoggingExclusion(id string) (*logging.LogExclusion, error)
	UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error
	DeleteLoggingExclusion(id string) error

	GetResourceType() string

	// Returns the unique resource identifier.
	GetResourceId() string

	// Textual description of this resource to be used in error message.
	// The description should include the unique resource identifier.
	DescribeResource() string
}

type newResourceLoggingExclusionUpdaterFunc func(d *schema.ResourceData, config *transport_tpg.Config) (ResourceLoggingExclusionUpdater, error)

// loggingExclusionResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingExclusionResourceTypes = []string{
	"billingAccounts",
	"folders",
	"organizations",
	"projects",
}

// LoggingExclusionId represents the parts that make up the canonical id used within terraform for a logging resource.
type LoggingExclusionId struct {
	resourceType string
	ResourceId   string
	name         string
}

// loggingExclusionIdRegex matches valid logging exclusion canonical ids
var loggingExclusionIdRegex = regexp.MustCompile("(.+)/(.+)/exclusions/(.+)")

// canonicalId returns the LoggingExclusionId as the canonical id used within terraform.
func (l LoggingExclusionId) canonicalId() string {
	return fmt.Sprintf("%s/%s/exclusions/%s", l.resourceType, l.ResourceId, l.name)
}

// parent returns the "parent-level" resource that the exclusion is in (e.g. `folders/foo` for id `folders/foo/exclusions/bar`)
func (l LoggingExclusionId) parent() string {
	return fmt.Sprintf("%s/%s", l.resourceType, l.ResourceId)
}

// ParseLoggingExclusionId parses a canonical id into a LoggingExclusionId, or returns an error on failure.
func ParseLoggingExclusionId(id string) (*LoggingExclusionId, error) {
	parts := loggingExclusionIdRegex.FindStringSubmatch(id)
	if parts == nil {
		return nil, fmt.Errorf("unable to parse logging exclusion id %#v", id)
	}
	// If our resourceType is not a valid logging exclusion resource type, complain loudly
	validLoggingExclusionResourceType := false
	for _, v := range loggingExclusionResourceTypes {
		if v == parts[1] {
			validLoggingExclusionResourceType = true
			break
		}
	}

	if !validLoggingExclusionResourceType {
		return nil, fmt.Errorf("Logging resource type %s is not valid. Valid resource types: %#v", parts[1],
			loggingExclusionResourceTypes)
	}
	return &LoggingExclusionId{
		resourceType: parts[1],
		ResourceId:   parts[2],
		name:         parts[3],
	}, nil
}
