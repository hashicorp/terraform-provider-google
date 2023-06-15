// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package bigqueryreservation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	bigqueryreservation "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/bigqueryreservation"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceBigqueryReservationAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigqueryReservationAssignmentCreate,
		Read:   resourceBigqueryReservationAssignmentRead,
		Delete: resourceBigqueryReservationAssignmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceBigqueryReservationAssignmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"assignee": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The resource which will use the reservation. E.g. projects/myproject, folders/123, organizations/456.",
			},

			"job_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Types of job, which could be specified when using the reservation. Possible values: JOB_TYPE_UNSPECIFIED, PIPELINE, QUERY",
			},

			"reservation": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The reservation for the resource",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The resource name of the assignment.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assignment will remain in PENDING state if no active capacity commitment is present. It will become ACTIVE when some capacity commitment becomes active. Possible values: STATE_UNSPECIFIED, PENDING, ACTIVE",
			},
		},
	}
}

func resourceBigqueryReservationAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &bigqueryreservation.Assignment{
		Assignee:    dcl.String(d.Get("assignee").(string)),
		JobType:     bigqueryreservation.AssignmentJobTypeEnumRef(d.Get("job_type").(string)),
		Reservation: dcl.String(d.Get("reservation").(string)),
		Location:    dcl.StringOrNil(d.Get("location").(string)),
		Project:     dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLBigqueryReservationClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyAssignment(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Assignment: %s", err)
	}

	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	// ID has a server-generated value, set again after creation.

	id, err = res.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Assignment %q: %#v", d.Id(), res)

	return resourceBigqueryReservationAssignmentRead(d, meta)
}

func resourceBigqueryReservationAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &bigqueryreservation.Assignment{
		Assignee:    dcl.String(d.Get("assignee").(string)),
		JobType:     bigqueryreservation.AssignmentJobTypeEnumRef(d.Get("job_type").(string)),
		Reservation: dcl.String(d.Get("reservation").(string)),
		Location:    dcl.StringOrNil(d.Get("location").(string)),
		Project:     dcl.String(project),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLBigqueryReservationClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetAssignment(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("BigqueryReservationAssignment %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("assignee", res.Assignee); err != nil {
		return fmt.Errorf("error setting assignee in state: %s", err)
	}
	if err = d.Set("job_type", res.JobType); err != nil {
		return fmt.Errorf("error setting job_type in state: %s", err)
	}
	if err = d.Set("reservation", res.Reservation); err != nil {
		return fmt.Errorf("error setting reservation in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("state", res.State); err != nil {
		return fmt.Errorf("error setting state in state: %s", err)
	}

	return nil
}

func resourceBigqueryReservationAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &bigqueryreservation.Assignment{
		Assignee:    dcl.String(d.Get("assignee").(string)),
		JobType:     bigqueryreservation.AssignmentJobTypeEnumRef(d.Get("job_type").(string)),
		Reservation: dcl.String(d.Get("reservation").(string)),
		Location:    dcl.StringOrNil(d.Get("location").(string)),
		Project:     dcl.String(project),
		Name:        dcl.StringOrNil(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] Deleting Assignment %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLBigqueryReservationClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteAssignment(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Assignment: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Assignment %q", d.Id())
	return nil
}

func resourceBigqueryReservationAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/reservations/(?P<reservation>[^/]+)/assignments/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<reservation>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<reservation>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
