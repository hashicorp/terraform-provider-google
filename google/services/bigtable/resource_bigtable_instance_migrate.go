// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigtableInstanceResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"num_nodes": {
							Type:     schema.TypeInt,
							Optional: true,
							// DEVELOPMENT instances could get returned with either zero or one node,
							// so mark as computed.
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"storage_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SSD",
							ValidateFunc: validation.StringInSlice([]string{"SSD", "HDD"}, false),
						},
					},
				},
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PRODUCTION",
				ValidateFunc: validation.StringInSlice([]string{"DEVELOPMENT", "PRODUCTION"}, false),
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func ResourceBigtableInstanceUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)

	rawState["deletion_protection"] = true

	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}
