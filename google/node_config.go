package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/container/v1"
)

// Matches gke-default scope from https://cloud.google.com/sdk/gcloud/reference/container/clusters/create
var defaultOauthScopes = []string{
	"https://www.googleapis.com/auth/devstorage.read_only",
	"https://www.googleapis.com/auth/logging.write",
	"https://www.googleapis.com/auth/monitoring",
	"https://www.googleapis.com/auth/service.management.readonly",
	"https://www.googleapis.com/auth/servicecontrol",
	"https://www.googleapis.com/auth/trace.append",
}

var schemaNodeConfig = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	ForceNew: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},

			"guest_accelerator": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"type": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: linkDiffSuppress,
						},
					},
				},
			},

			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     schema.TypeString,
			},

			"local_ssd_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"machine_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     schema.TypeString,
			},

			"min_cpu_platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"oauth_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(v interface{}) string {
						return canonicalizeServiceScope(v.(string))
					},
				},
				Set: stringScopeHashcode,
			},

			"preemptible": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"service_account": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

func expandNodeConfig(v interface{}) *container.NodeConfig {
	nodeConfigs := v.([]interface{})
	nc := &container.NodeConfig{
		// Defaults can't be set on a list/set in the schema, so set the default on create here.
		OauthScopes: defaultOauthScopes,
	}
	if len(nodeConfigs) == 0 {
		return nc
	}

	nodeConfig := nodeConfigs[0].(map[string]interface{})

	if v, ok := nodeConfig["machine_type"]; ok {
		nc.MachineType = v.(string)
	}

	if v, ok := nodeConfig["guest_accelerator"]; ok {
		accels := v.([]interface{})
		guestAccelerators := make([]*container.AcceleratorConfig, 0, len(accels))
		for _, raw := range accels {
			data := raw.(map[string]interface{})
			if data["count"].(int) == 0 {
				continue
			}
			guestAccelerators = append(guestAccelerators, &container.AcceleratorConfig{
				AcceleratorCount: int64(data["count"].(int)),
				AcceleratorType:  data["type"].(string),
			})
		}
		nc.Accelerators = guestAccelerators
	}

	if v, ok := nodeConfig["disk_size_gb"]; ok {
		nc.DiskSizeGb = int64(v.(int))
	}

	if v, ok := nodeConfig["local_ssd_count"]; ok {
		nc.LocalSsdCount = int64(v.(int))
	}

	if scopes, ok := nodeConfig["oauth_scopes"]; ok {
		scopesSet := scopes.(*schema.Set)
		scopes := make([]string, scopesSet.Len())
		for i, scope := range scopesSet.List() {
			scopes[i] = canonicalizeServiceScope(scope.(string))
		}

		nc.OauthScopes = scopes
	}

	if v, ok := nodeConfig["service_account"]; ok {
		nc.ServiceAccount = v.(string)
	}

	if v, ok := nodeConfig["metadata"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.Metadata = m
	}

	if v, ok := nodeConfig["image_type"]; ok {
		nc.ImageType = v.(string)
	}

	if v, ok := nodeConfig["labels"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.Labels = m
	}

	if v, ok := nodeConfig["tags"]; ok {
		tagsList := v.([]interface{})
		tags := []string{}
		for _, v := range tagsList {
			tags = append(tags, v.(string))
		}
		nc.Tags = tags
	}
	// Preemptible Is Optional+Default, so it always has a value
	nc.Preemptible = nodeConfig["preemptible"].(bool)

	if v, ok := nodeConfig["min_cpu_platform"]; ok {
		nc.MinCpuPlatform = v.(string)
	}

	return nc
}

func flattenNodeConfig(c *container.NodeConfig) []map[string]interface{} {
	config := make([]map[string]interface{}, 0, 1)

	if c == nil {
		return config
	}

	config = append(config, map[string]interface{}{
		"machine_type":      c.MachineType,
		"disk_size_gb":      c.DiskSizeGb,
		"guest_accelerator": c.Accelerators,
		"local_ssd_count":   c.LocalSsdCount,
		"service_account":   c.ServiceAccount,
		"metadata":          c.Metadata,
		"image_type":        c.ImageType,
		"labels":            c.Labels,
		"tags":              c.Tags,
		"preemptible":       c.Preemptible,
		"min_cpu_platform":  c.MinCpuPlatform,
	})

	if len(c.OauthScopes) > 0 {
		config[0]["oauth_scopes"] = schema.NewSet(stringScopeHashcode, convertStringArrToInterface(c.OauthScopes))
	}

	return config
}
