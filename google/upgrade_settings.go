package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	containerBeta "google.golang.org/api/container/v1beta1"
)

var schemaUpgradeSettings = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_surge": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"max_unavailable": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	},
}

func flattenUpgradeSettings(u *containerBeta.UpgradeSettings) []map[string]interface{} {
	r := make(map[string]interface{})
	if u == nil {
		return nil
	}
	r["max_surge"] = u.MaxSurge
	r["max_unavailable"] = u.MaxUnavailable

	return []map[string]interface{}{r}
}

func expandUpgradeSettings(configured interface{}) *containerBeta.UpgradeSettings {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &containerBeta.UpgradeSettings{}
	}
	config := l[0].(map[string]interface{})
	result := &containerBeta.UpgradeSettings{
		ForceSendFields: []string{"MaxSurge", "MaxUnavailable"},
	}
	if v, ok := config["max_surge"]; ok {
		result.MaxSurge = int64(v.(int))
	}
	if v, ok := config["max_unavailable"]; ok {
		result.MaxUnavailable = int64(v.(int))
	}

	return result
}
