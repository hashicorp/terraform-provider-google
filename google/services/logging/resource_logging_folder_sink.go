// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceLoggingFolderSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingFolderSinkCreate,
		Read:   resourceLoggingFolderSinkRead,
		Delete: resourceLoggingFolderSinkDelete,
		Update: resourceLoggingFolderSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("folder"),
		},
		UseJSONNumber: true,
	}
	schm.Schema["folder"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The folder to be exported to the sink. Note that either [FOLDER_ID] or "folders/[FOLDER_ID]" is accepted.`,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "folders/", "", 1)
		},
	}
	schm.Schema["include_children"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: `Whether or not to include children folders in the sink export. If true, logs associated with child projects are also exported; otherwise only logs relating to the provided folder are included.`,
	}

	return schm
}

func resourceLoggingFolderSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	folder := resourcemanager.ParseFolderId(d.Get("folder"))
	id, sink := expandResourceLoggingSink(d, "folders", folder)
	sink.IncludeChildren = d.Get("include_children").(bool)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).Folders.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingFolderSinkRead(d, meta)
}

func resourceLoggingFolderSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, err := config.NewLoggingClient(userAgent).Folders.Sinks.Get(d.Id()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Folder Logging Sink %s", d.Get("name").(string)))
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	if err := d.Set("include_children", sink.IncludeChildren); err != nil {
		return fmt.Errorf("Error setting include_children: %s", err)
	}

	return nil
}

func resourceLoggingFolderSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)
	// It seems the API might actually accept an update for include_children; this is not in the list of updatable
	// properties though and might break in the future. Always include the value to prevent it changing.
	sink.IncludeChildren = d.Get("include_children").(bool)
	sink.ForceSendFields = append(sink.ForceSendFields, "IncludeChildren")

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).Folders.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	return resourceLoggingFolderSinkRead(d, meta)
}

func resourceLoggingFolderSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	_, err = config.NewLoggingClient(userAgent).Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}
