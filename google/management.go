package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	containerBeta "google.golang.org/api/container/v1beta1"
)

var schemaManagement = &containerBeta.NodeManagement{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_repair": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"auto_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	},
}

func flattenManagement(n *containerBeta.NodeManagement) []map[string]interface{} {
	r := make(map[string]interface{})
	r["auto_repair"] = n.AutoRepair
	r["auto_upgrade"] = n.AutoUpgrade

	return []map[string]interface{}{r}
}

func expandManagement(configured interface{}) *containerBeta.NodeManagement {
	l, ok := configured.([]interface{})
	if !ok || l == nil || len(l) == 0 || l[0] == nil {
		return &containerBeta.NodeManagement{}
	}
	config := l[0].(map[string]interface{})
	result := &containerBeta.NodeManagement{
		ForceSendFields: []string{"AutoRepair", "AutoUpgrade"},
	}
	if v, ok := config["auto_repair"]; ok {
		result.AutoRepair = v.(bool)
	}
	if v, ok := config["auto_upgrade"]; ok {
		result.AutoUpgrade = v.(bool)
	}
	return result
}
