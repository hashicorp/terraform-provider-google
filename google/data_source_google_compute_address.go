package google

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var (
	computeAddressIdTemplate = "projects/%s/regions/%s/addresses/%s"
	computeAddressLinkRegex  = regexp.MustCompile("projects/(.+)/regions/(.+)/addresses/(.+)$")
)

func dataSourceGoogleComputeAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeAddressRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleComputeAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	address, err := config.clientCompute.Addresses.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Address Not Found : %s", name))
	}

	d.Set("address", address.Address)
	d.Set("status", address.Status)
	d.Set("self_link", address.SelfLink)
	d.Set("project", project)
	d.Set("region", region)

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/addresses/%s", project, region, name))
	return nil
}

type computeAddressId struct {
	Project string
	Region  string
	Name    string
}

func (s computeAddressId) canonicalId() string {
	return fmt.Sprintf(computeAddressIdTemplate, s.Project, s.Region, s.Name)
}

func parseComputeAddressId(id string, config *Config) (*computeAddressId, error) {
	var parts []string
	if computeAddressLinkRegex.MatchString(id) {
		parts = computeAddressLinkRegex.FindStringSubmatch(id)

		return &computeAddressId{
			Project: parts[1],
			Region:  parts[2],
			Name:    parts[3],
		}, nil
	} else {
		parts = strings.Split(id, "/")
	}

	if len(parts) == 3 {
		return &computeAddressId{
			Project: parts[0],
			Region:  parts[1],
			Name:    parts[2],
		}, nil
	} else if len(parts) == 2 {
		// Project is optional.
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{region}/{name}` id format.")
		}

		return &computeAddressId{
			Project: config.Project,
			Region:  parts[0],
			Name:    parts[1],
		}, nil
	} else if len(parts) == 1 {
		// Project and region is optional
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{name}` id format.")
		}
		if config.Region == "" {
			return nil, fmt.Errorf("The default region for the provider must be set when using the `{name}` id format.")
		}

		return &computeAddressId{
			Project: config.Project,
			Region:  config.Region,
			Name:    parts[0],
		}, nil
	}

	return nil, fmt.Errorf("Invalid compute address id. Expecting resource link, `{project}/{region}/{name}`, `{region}/{name}` or `{name}` format.")
}
