package google

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleContainerCluster() *schema.Resource {
	return &schema.Resource{
		Read: datasourceContainerClusterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"additional_zones": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"addons_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_load_balancing": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"horizontal_pod_autoscaling": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"kubernetes_dashboard": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"cluster_ipv4_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_kubernetes_alpha": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_legacy_abac": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"initial_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"logging_service": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"maintenance_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_maintenance_window": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"duration": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"master_auth": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"client_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"client_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"cluster_ca_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"master_authorized_networks_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_blocks": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"min_master_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"monitoring_service": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"node_config": schemaNodeConfig,

			"node_pool": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaNodePool,
				},
			},

			"node_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnetwork": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_group_urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"master_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_allocation_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_secondary_range_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"services_secondary_range_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	clusterName := d.Get("name").(string)

	d.SetId(clusterName)

	return resourceContainerClusterRead(d, meta)
}
