// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package containeranalysis

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerRegistryCreate,
		Read:   resourceContainerRegistryRead,
		Delete: resourceContainerRegistryDelete,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(s interface{}) string {
					return strings.ToUpper(s.(string))
				},
				Description: `The location of the registry. One of ASIA, EU, US or not specified. See the official documentation for more information on registry locations.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"bucket_self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceContainerRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Project: %s", project)

	location := d.Get("location").(string)
	log.Printf("[DEBUG] location: %s", location)
	urlBase := "https://gcr.io/v2/token"
	if location != "" {
		urlBase = fmt.Sprintf("https://%s.gcr.io/v2/token", strings.ToLower(location))
	}

	// Performing a token handshake with the GCR API causes the backing bucket to create if it hasn't already.
	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("%s?service=gcr.io&scope=repository:{{project}}/my-repo:push,pull", urlBase))
	if err != nil {
		return err
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})

	if err != nil {
		return err
	}
	return resourceContainerRegistryRead(d, meta)
}

func resourceContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	location := d.Get("location").(string)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	name := ""
	if location != "" {
		name = fmt.Sprintf("%s.artifacts.%s.appspot.com", strings.ToLower(location), project)
	} else {
		name = fmt.Sprintf("artifacts.%s.appspot.com", project)
	}

	res, err := config.NewStorageClient(userAgent).Buckets.Get(name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Container Registry Storage Bucket %q", name))
	}
	log.Printf("[DEBUG] Read bucket %v at location %v\n\n", res.Name, res.SelfLink)

	// Update the ID according to the bucket ID
	if err := d.Set("bucket_self_link", res.SelfLink); err != nil {
		return fmt.Errorf("Error setting bucket_self_link: %s", err)
	}

	d.SetId(res.Id)
	return nil
}

func resourceContainerRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	// Don't delete the backing bucket as this is not a supported GCR action
	return nil
}
