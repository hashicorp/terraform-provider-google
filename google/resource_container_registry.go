package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceContainerRegistry() *schema.Resource {
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
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"bucket_self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceContainerRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
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
	url, err := replaceVars(d, config, fmt.Sprintf("%s?service=gcr.io&scope=repository:{{project}}/my-repo:push,pull", urlBase))
	if err != nil {
		return err
	}

	_, err = sendRequestWithTimeout(config, "GET", project, url, nil, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return err
	}
	return resourceContainerRegistryRead(d, meta)
}

func resourceContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	location := d.Get("location").(string)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	name := ""
	if location != "" {
		name = fmt.Sprintf("%s.artifacts.%s.appspot.com", strings.ToLower(location), project)
	} else {
		name = fmt.Sprintf("artifacts.%s.appspot.com", project)
	}

	res, err := config.clientStorage.Buckets.Get(name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Container Registry Storage Bucket %q", name))
	}
	log.Printf("[DEBUG] Read bucket %v at location %v\n\n", res.Name, res.SelfLink)

	// Update the ID according to the bucket ID
	d.Set("bucket_self_link", res.SelfLink)

	d.SetId(res.Id)
	return nil
}

func resourceContainerRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	// Don't delete the backing bucket as this is not a supported GCR action
	return nil
}
