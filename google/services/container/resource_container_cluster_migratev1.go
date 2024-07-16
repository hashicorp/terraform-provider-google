// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceContainerClusterUpgradeV1(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Applying container cluster migration to schema version V2.")

	rawState["deletion_protection"] = true
	return rawState, nil
}

func resourceContainerClusterResourceV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the cluster, unique within the project and location.`,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) > 40 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 40 characters", k))
					}
					if !regexp.MustCompile("^[a-z0-9-]+$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q can only contain lowercase letters, numbers and hyphens", k))
					}
					if !regexp.MustCompile("^[a-z]").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must start with a letter", k))
					}
					if !regexp.MustCompile("[a-z0-9]$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must end with a number or a letter", k))
					}
					return
				},
			},

			"operation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The location (region or zone) in which the cluster master will be created, as well as the default node location. If you specify a zone (such as us-central1-a), the cluster will be a zonal cluster with a single cluster master. If you specify a region (such as us-west1), the cluster will be a regional cluster with multiple masters spread across zones in the region, and with default node locations in those zones as well.`,
			},

			"node_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of zones in which the cluster's nodes are located. Nodes must be in the region of their regional cluster or in the same region as their cluster's zone for zonal clusters. If this is specified for a zonal cluster, omit the cluster's zone.`,
			},

			"addons_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The configuration for addons supported by GKE.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_load_balancing": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the HTTP (L7) load balancing controller addon, which makes it easy to set up HTTP load balancers for services in a cluster. It is enabled by default; set disabled = true to disable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"horizontal_pod_autoscaling": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the Horizontal Pod Autoscaling addon, which increases or decreases the number of replica pods a replication controller has based on the resource usage of the existing pods. It ensures that a Heapster pod is running in the cluster, which is also used by the Cloud Monitoring service. It is enabled by default; set disabled = true to disable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"network_policy_config": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							AtLeastOneOf:  addonsConfigKeys,
							MaxItems:      1,
							Description:   `Whether we should enable the network policy addon for the master. This must be enabled in order to enable network policy for the nodes. To enable this, you must also define a network_policy block, otherwise nothing will happen. It can only be disabled if the nodes already do not have network policies enabled. Defaults to disabled; set disabled = false to enable.`,
							ConflictsWith: []string{"enable_autopilot"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gcp_filestore_csi_driver_config": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							AtLeastOneOf:  addonsConfigKeys,
							MaxItems:      1,
							Description:   `The status of the Filestore CSI driver addon, which allows the usage of filestore instance as volumes. Defaults to disabled; set enabled = true to enable.`,
							ConflictsWith: []string{"enable_autopilot"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"cloudrun_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the CloudRun addon. It is disabled by default. Set disabled = false to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"load_balancer_type": {
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"LOAD_BALANCER_TYPE_INTERNAL"}, false),
										Optional:     true,
									},
								},
							},
						},
						"dns_cache_config": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							AtLeastOneOf:  addonsConfigKeys,
							MaxItems:      1,
							Description:   `The status of the NodeLocal DNSCache addon. It is disabled by default. Set enabled = true to enable.`,
							ConflictsWith: []string{"enable_autopilot"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gce_persistent_disk_csi_driver_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `Whether this cluster should enable the Google Compute Engine Persistent Disk Container Storage Interface (CSI) Driver. Set enabled = true to enable. The Compute Engine persistent disk CSI Driver is enabled by default on newly created clusters for the following versions: Linux clusters: GKE version 1.18.10-gke.2100 or later, or 1.19.3-gke.2100 or later.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gke_backup_agent_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the Backup for GKE Agent addon. It is disabled by default. Set enabled = true to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gcs_fuse_csi_driver_config": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							AtLeastOneOf:  addonsConfigKeys,
							MaxItems:      1,
							Description:   `The status of the GCS Fuse CSI driver addon, which allows the usage of gcs bucket as volumes. Defaults to disabled; set enabled = true to enable.`,
							ConflictsWith: []string{"enable_autopilot"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"config_connector_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The of the Config Connector addon.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"cluster_autoscaling": {
				Type:     schema.TypeList,
				MaxItems: 1,
				// This field is Optional + Computed because we automatically set the
				// enabled value to false if the block is not returned in API responses.
				Optional:    true,
				Computed:    true,
				Description: `Per-cluster configuration of Node Auto-Provisioning with Cluster Autoscaler to automatically adjust the size of the cluster and create/delete node pools based on the current needs of the cluster's workload. See the guide to using Node Auto-Provisioning for more details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"enable_autopilot"},
							Description:   `Whether node auto-provisioning is enabled. Resource limits for cpu and memory must be defined to enable node auto-provisioning.`,
						},
						"resource_limits": {
							Type:             schema.TypeList,
							Optional:         true,
							ConflictsWith:    []string{"enable_autopilot"},
							DiffSuppressFunc: suppressDiffForAutopilot,
							Description:      `Global constraints for machine resources in the cluster. Configuring the cpu and memory types is required if node auto-provisioning is enabled. These limits will apply to node pool autoscaling in addition to node auto-provisioning.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The type of the resource. For example, cpu and memory. See the guide to using Node Auto-Provisioning for a list of types.`,
									},
									"minimum": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Minimum amount of the resource in the cluster.`,
									},
									"maximum": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Maximum amount of the resource in the cluster.`,
									},
								},
							},
						},
						"auto_provisioning_defaults": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: `Contains defaults for a node pool created by NAP.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oauth_scopes": {
										Type:             schema.TypeList,
										Optional:         true,
										Computed:         true,
										Elem:             &schema.Schema{Type: schema.TypeString},
										DiffSuppressFunc: containerClusterAddedScopesSuppress,
										Description:      `Scopes that are used by NAP when creating node pools.`,
									},
									"service_account": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "default",
										Description: `The Google Cloud Platform Service Account to be used by the node VMs.`,
									},
									"disk_size": {
										Type:             schema.TypeInt,
										Optional:         true,
										Default:          100,
										Description:      `Size of the disk attached to each node, specified in GB. The smallest allowed disk size is 10GB.`,
										DiffSuppressFunc: suppressDiffForAutopilot,
										ValidateFunc:     validation.IntAtLeast(10),
									},
									"disk_type": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "pd-standard",
										Description:      `Type of the disk attached to each node.`,
										DiffSuppressFunc: suppressDiffForAutopilot,
										ValidateFunc:     validation.StringInSlice([]string{"pd-standard", "pd-ssd", "pd-balanced"}, false),
									},
									"image_type": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "COS_CONTAINERD",
										Description:      `The default image type used by NAP once a new node pool is being created.`,
										DiffSuppressFunc: suppressDiffForAutopilot,
										ValidateFunc:     validation.StringInSlice([]string{"COS_CONTAINERD", "COS", "UBUNTU_CONTAINERD", "UBUNTU"}, false),
									},
									"min_cpu_platform": {
										Type:             schema.TypeString,
										Optional:         true,
										DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("automatic"),
										Description:      `Minimum CPU platform to be used by this instance. The instance may be scheduled on the specified or newer CPU platform. Applicable values are the friendly names of CPU platforms, such as Intel Haswell.`,
									},
									"boot_disk_kms_key": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: `The Customer Managed Encryption Key used to encrypt the boot disk attached to each node in the node pool.`,
									},
									"shielded_instance_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Shielded Instance options.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_secure_boot": {
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
													Description: `Defines whether the instance has Secure Boot enabled.`,
													AtLeastOneOf: []string{
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_secure_boot",
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_integrity_monitoring",
													},
												},
												"enable_integrity_monitoring": {
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     true,
													Description: `Defines whether the instance has integrity monitoring enabled.`,
													AtLeastOneOf: []string{
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_secure_boot",
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_integrity_monitoring",
													},
												},
											},
										},
									},
									"management": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: `NodeManagement configuration for this NodePool.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_upgrade": {
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: `Specifies whether node auto-upgrade is enabled for the node pool. If enabled, node auto-upgrade helps keep the nodes in your node pool up to date with the latest release version of Kubernetes.`,
												},
												"auto_repair": {
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: `Specifies whether the node auto-repair is enabled for the node pool. If enabled, the nodes in this node pool will be monitored and, if they fail health checks too many times, an automatic repair action will be triggered.`,
												},
												"upgrade_options": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: `Specifies the Auto Upgrade knobs for the node pool.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"auto_upgrade_start_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: `This field is set when upgrades are about to commence with the approximate start time for the upgrades, in RFC3339 text format.`,
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: `This field is set when upgrades are about to commence with the description of the upgrade.`,
															},
														},
													},
												},
											},
										},
									},
									"upgrade_settings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Specifies the upgrade settings for NAP created node pools`,
										Computed:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_surge": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: `The maximum number of nodes that can be created beyond the current size of the node pool during the upgrade process.`,
												},
												"max_unavailable": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: `The maximum number of nodes that can be simultaneously unavailable during the upgrade process.`,
												},
												"strategy": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													Description:  `Update strategy of the node pool.`,
													ValidateFunc: validation.StringInSlice([]string{"NODE_POOL_UPDATE_STRATEGY_UNSPECIFIED", "BLUE_GREEN", "SURGE"}, false),
												},
												"blue_green_settings": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													MaxItems:    1,
													Description: `Settings for blue-green upgrade strategy.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"node_pool_soak_duration": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
																Description: `Time needed after draining entire blue pool. After this period, blue pool will be cleaned up.

																A duration in seconds with up to nine fractional digits, ending with 's'. Example: "3.5s".`,
															},
															"standard_rollout_policy": {
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																MaxItems:    1,
																Description: `Standard policy for the blue-green upgrade.`,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"batch_percentage": {
																			Type:         schema.TypeFloat,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.FloatBetween(0.0, 1.0),
																			ExactlyOneOf: []string{
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_percentage",
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_node_count",
																			},
																			Description: `Percentage of the bool pool nodes to drain in a batch. The range of this field should be (0.0, 1.0].`,
																		},
																		"batch_node_count": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			Computed: true,
																			ExactlyOneOf: []string{
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_percentage",
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_node_count",
																			},
																			Description: `Number of blue nodes to drain in a batch.`,
																		},
																		"batch_soak_duration": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Default:  "0s",
																			Description: `Soak time after each batch gets drained.

																			A duration in seconds with up to nine fractional digits, ending with 's'. Example: "3.5s".`,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"cluster_ipv4_cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  verify.OrEmpty(verify.ValidateRFC1918Network(8, 32)),
				ConflictsWith: []string{"ip_allocation_policy"},
				Description:   `The IP address range of the Kubernetes pods in this cluster in CIDR notation (e.g. 10.96.0.0/14). Leave blank to have one automatically chosen or specify a /14 block in 10.0.0.0/8. This field will only work for routes-based clusters, where ip_allocation_policy is not defined.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: ` Description of the cluster.`,
			},

			"binary_authorization": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: BinaryAuthorizationDiffSuppress,
				MaxItems:         1,
				Description:      "Configuration options for the Binary Authorization feature.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:          schema.TypeBool,
							Optional:      true,
							Deprecated:    "Deprecated in favor of evaluation_mode.",
							Description:   "Enable Binary Authorization for this cluster.",
							ConflictsWith: []string{"enable_autopilot", "binary_authorization.0.evaluation_mode"},
						},
						"evaluation_mode": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringInSlice([]string{"DISABLED", "PROJECT_SINGLETON_POLICY_ENFORCE"}, false),
							Description:   "Mode of operation for Binary Authorization policy evaluation.",
							ConflictsWith: []string{"binary_authorization.0.enabled"},
						},
					},
				},
			},

			"enable_kubernetes_alpha": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: `Whether to enable Kubernetes Alpha features for this cluster. Note that when this option is enabled, the cluster cannot be upgraded and will be automatically deleted after 30 days.`,
			},

			"enable_k8s_beta_apis": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Configuration for Kubernetes Beta APIs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled_apis": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Enabled Kubernetes Beta APIs.`,
						},
					},
				},
			},

			"enable_tpu": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to enable Cloud TPU resources in this cluster.`,
			},

			"enable_legacy_abac": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether the ABAC authorizer is enabled for this cluster. When enabled, identities in the system, including service accounts, nodes, and controllers, will have statically granted permissions beyond those provided by the RBAC configuration or IAM. Defaults to false.`,
			},

			"enable_shielded_nodes": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       true,
				Description:   `Enable Shielded Nodes features on all nodes in this cluster. Defaults to true.`,
				ConflictsWith: []string{"enable_autopilot"},
			},

			"enable_autopilot": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Enable Autopilot for this cluster.`,
				// ConflictsWith: many fields, see https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview#comparison. The conflict is only set one-way, on other fields w/ this field.
			},

			"allow_net_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Enable NET_ADMIN for this cluster.`,
			},

			"authenticator_groups_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration for the Google Groups for GKE feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the RBAC security group for use with Google security groups in Kubernetes RBAC. Group name must be in format gke-security-groups@yourdomain.com.`,
						},
					},
				},
			},

			"initial_node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `The number of nodes to create in this cluster's default node pool. In regional or multi-zonal clusters, this is the number of nodes per zone. Must be set if node_pool is not set. If you're using google_container_node_pool objects with no default node pool, you'll need to set this to a value of at least 1, alongside setting remove_default_node_pool to true.`,
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Logging configuration for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_components": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `GKE components exposing logs. Valid values include SYSTEM_COMPONENTS, APISERVER, CONTROLLER_MANAGER, SCHEDULER, and WORKLOADS.`,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER", "SCHEDULER", "WORKLOADS"}, false),
							},
						},
					},
				},
			},

			"logging_service": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"logging.googleapis.com", "logging.googleapis.com/kubernetes", "none"}, false),
				Description:  `The logging service that the cluster should write logs to. Available options include logging.googleapis.com(Legacy Stackdriver), logging.googleapis.com/kubernetes(Stackdriver Kubernetes Engine Logging), and none. Defaults to logging.googleapis.com/kubernetes.`,
			},

			"maintenance_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The maintenance policy to use for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_maintenance_window": {
							Type:     schema.TypeList,
							Optional: true,
							ExactlyOneOf: []string{
								"maintenance_policy.0.daily_maintenance_window",
								"maintenance_policy.0.recurring_window",
							},
							MaxItems:    1,
							Description: `Time window specified for daily maintenance operations. Specify start_time in RFC3339 format "HH:MM”, where HH : [00-23] and MM : [00-59] GMT.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateFunc:     verify.ValidateRFC3339Time,
										DiffSuppressFunc: Rfc3339TimeDiffSuppress,
									},
									"duration": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"recurring_window": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ExactlyOneOf: []string{
								"maintenance_policy.0.daily_maintenance_window",
								"maintenance_policy.0.recurring_window",
							},
							Description: `Time window for recurring maintenance operations.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"end_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"recurrence": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: rfc5545RecurrenceDiffSuppress,
									},
								},
							},
						},
						"maintenance_exclusion": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    20,
							Description: `Exceptions to maintenance window. Non-emergency maintenance should not occur in these windows.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclusion_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"end_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"exclusion_options": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: `Maintenance exclusion related options.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scope": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"NO_UPGRADES", "NO_MINOR_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES"}, false),
													Description:  `The scope of automatic upgrades to restrict in the exclusion window.`,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"security_posture_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `Defines the config needed to enable/disable features for the Security Posture API`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringInSlice([]string{"DISABLED", "BASIC", "ENTERPRISE", "MODE_UNSPECIFIED"}, false),
							Description:      `Sets the mode of the Kubernetes security posture API's off-cluster features. Available options include DISABLED, BASIC, and ENTERPRISE.`,
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("MODE_UNSPECIFIED"),
						},
						"vulnerability_mode": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringInSlice([]string{"VULNERABILITY_DISABLED", "VULNERABILITY_BASIC", "VULNERABILITY_ENTERPRISE", "VULNERABILITY_MODE_UNSPECIFIED"}, false),
							Description:      `Sets the mode of the Kubernetes security posture API's workload vulnerability scanning. Available options include VULNERABILITY_DISABLED, VULNERABILITY_BASIC and VULNERABILITY_ENTERPRISE.`,
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("VULNERABILITY_MODE_UNSPECIFIED"),
						},
					},
				},
			},
			"monitoring_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Monitoring configuration for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_components": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `GKE components exposing metrics. Valid values include SYSTEM_COMPONENTS, APISERVER, SCHEDULER, CONTROLLER_MANAGER, STORAGE, HPA, POD, DAEMONSET, DEPLOYMENT, STATEFULSET and DCGM.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"managed_prometheus": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `Configuration for Google Cloud Managed Services for Prometheus.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not the managed collection is enabled.`,
									},
								},
							},
						},
						"advanced_datapath_observability_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    2,
							Description: `Configuration of Advanced Datapath Observability features.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_metrics": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not the advanced datapath metrics are enabled.`,
									},
									"enable_relay": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Whether or not Relay is enabled.`,
									},
									"relay_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										Description:  `Mode used to make Relay available.`,
										ValidateFunc: validation.StringInSlice([]string{"DISABLED", "INTERNAL_VPC_LB", "EXTERNAL_LB"}, false),
									},
								},
							},
						},
					},
				},
			},

			"notification_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The notification config for sending cluster upgrade notifications`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pubsub": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: `Notification config for Cloud Pub/Sub`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not the notification config is enabled`,
									},
									"topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The pubsub topic to push upgrade notifications to. Must be in the same project as the cluster. Must be in the format: projects/{project}/topics/{topic}.`,
									},
									"filter": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: `Allows filtering to one or more specific event types. If event types are present, those and only those event types will be transmitted to the cluster. Other types will be skipped. If no filter is specified, or no event types are present, all event types will be sent`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"event_type": {
													Type:        schema.TypeList,
													Required:    true,
													Description: `Can be used to filter what notifications are sent. Valid values include include UPGRADE_AVAILABLE_EVENT, UPGRADE_EVENT and SECURITY_BULLETIN_EVENT`,
													Elem: &schema.Schema{
														Type:         schema.TypeString,
														ValidateFunc: validation.StringInSlice([]string{"UPGRADE_AVAILABLE_EVENT", "UPGRADE_EVENT", "SECURITY_BULLETIN_EVENT"}, false),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"confidential_nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `Configuration for the confidential nodes feature, which makes nodes run on confidential VMs. Warning: This configuration can't be changed (or added/removed) after cluster creation without deleting and recreating the entire cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: `Whether Confidential Nodes feature is enabled for all nodes in this cluster.`,
						},
					},
				},
			},

			"master_auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `The authentication information for accessing the Kubernetes master. Some values in this block are only returned by the API if your service account has permission to get credentials for your GKE cluster. If you see an unexpected diff unsetting your client cert, ensure you have the container.clusters.getCredentials permission.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_certificate_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							ForceNew:    true,
							Description: `Whether client certificate authorization is enabled for this cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issue_client_certificate": {
										Type:        schema.TypeBool,
										Required:    true,
										ForceNew:    true,
										Description: `Whether client certificate authorization is enabled for this cluster.`,
									},
								},
							},
						},

						"client_certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Base64 encoded public certificate used by clients to authenticate to the cluster endpoint.`,
						},

						"client_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `Base64 encoded private key used by clients to authenticate to the cluster endpoint.`,
						},

						"cluster_ca_certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Base64 encoded public certificate that is the root of trust for the cluster.`,
						},
					},
				},
			},

			"master_authorized_networks_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        masterAuthorizedNetworksConfig,
				Description: `The desired configuration options for master authorized networks. Omit the nested cidr_blocks attribute to disallow external access (except the cluster node IPs, which GKE automatically whitelists).`,
			},

			"min_master_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The minimum version of the master. GKE will auto-update the master to new versions, so this does not guarantee the current master version--use the read-only master_version field to obtain that. If unset, the cluster's version will be set by GKE to the version of the most recent official release (which is not necessarily the latest version).`,
			},

			"monitoring_service": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"monitoring.googleapis.com", "monitoring.googleapis.com/kubernetes", "none"}, false),
				Description:  `The monitoring service that the cluster should write metrics to. Automatically send metrics from pods in the cluster to the Google Cloud Monitoring API. VM metrics will be collected by Google Compute Engine regardless of this setting Available options include monitoring.googleapis.com(Legacy Stackdriver), monitoring.googleapis.com/kubernetes(Stackdriver Kubernetes Engine Monitoring), and none. Defaults to monitoring.googleapis.com/kubernetes.`,
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "default",
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name or self_link of the Google Compute Engine network to which the cluster is connected. For Shared VPC, set this to the self link of the shared network.`,
			},

			"network_policy": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				Description:      `Configuration options for the NetworkPolicy feature.`,
				ConflictsWith:    []string{"enable_autopilot"},
				DiffSuppressFunc: containerClusterNetworkPolicyDiffSuppress,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether network policy is enabled on the cluster.`,
						},
						"provider": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"PROVIDER_UNSPECIFIED", "CALICO"}, false),
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("PROVIDER_UNSPECIFIED"),
							Description:      `The selected network policy provider.`,
						},
					},
				},
			},

			"node_config": clusterSchemaNodeConfig(),

			"node_pool": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true, // TODO: Add ability to add/remove nodePools
				Elem: &schema.Resource{
					Schema: schemaNodePool,
				},
				Description:   `List of node pools associated with this cluster. See google_container_node_pool for schema. Warning: node pools defined inside a cluster can't be changed (or added/removed) after cluster creation without deleting and recreating the entire cluster. Unless you absolutely need the ability to say "these are the only node pools associated with this cluster", use the google_container_node_pool resource instead of this property.`,
				ConflictsWith: []string{"enable_autopilot"},
			},

			"node_pool_defaults": clusterSchemaNodePoolDefaults(),

			"node_pool_auto_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Node pool configs that apply to all auto-provisioned node pools in autopilot clusters and node auto-provisioning enabled clusters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Collection of Compute Engine network tags that can be applied to a node's underlying VM instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tags": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `List of network tags applied to auto-provisioned node pools.`,
									},
								},
							},
						},
					},
				},
			},

			"node_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The Kubernetes version on the nodes. Must either be unset or set to the same value as min_master_version on create. Defaults to the default version set by GKE which is not necessarily the latest version. This only affects nodes in the default node pool. While a fuzzy version can be specified, it's recommended that you specify explicit versions as Terraform will see spurious diffs when fuzzy versions are used. See the google_container_engine_versions data source's version_prefix field to approximate fuzzy versions in a Terraform-compatible way. To update nodes in other node pools, use the version attribute on the node pool.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name or self_link of the Google Compute Engine subnetwork in which the cluster's instances are launched.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Server-defined URL for the resource.`,
			},

			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of this cluster's Kubernetes master.`,
			},

			"master_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current version of the master in the cluster. This may be different than the min_master_version set in the config if the master has been updated by GKE.`,
			},

			"services_ipv4_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address range of the Kubernetes services in this cluster, in CIDR notation (e.g. 1.2.3.4/29). Service addresses are typically put in the last /16 from the container CIDR.`,
			},

			"ip_allocation_policy": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ForceNew:      true,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"cluster_ipv4_cidr"},
				Description:   `Configuration of cluster IP allocation for VPC-native clusters. Adding this block enables IP aliasing, making the cluster VPC-native instead of routes-based.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// GKE creates/deletes secondary ranges in VPC
						"cluster_ipv4_cidr_block": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    ipAllocationRangeFields,
							DiffSuppressFunc: tpgresource.CidrOrSizeDiffSuppress,
							Description:      `The IP address range for the cluster pod IPs. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.`,
						},

						"services_ipv4_cidr_block": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    ipAllocationRangeFields,
							DiffSuppressFunc: tpgresource.CidrOrSizeDiffSuppress,
							Description:      `The IP address range of the services IPs in this cluster. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.`,
						},

						// User manages secondary ranges manually
						"cluster_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
							Description:   `The name of the existing secondary range in the cluster's subnetwork to use for pod IP addresses. Alternatively, cluster_ipv4_cidr_block can be used to automatically create a GKE-managed one.`,
						},

						"services_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
							Description:   `The name of the existing secondary range in the cluster's subnetwork to use for service ClusterIPs. Alternatively, services_ipv4_cidr_block can be used to automatically create a GKE-managed one.`,
						},

						"stack_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "IPV4",
							ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV4_IPV6"}, false),
							Description:  `The IP Stack type of the cluster. Choose between IPV4 and IPV4_IPV6. Default type is IPV4 Only if not set`,
						},
						"pod_cidr_overprovision_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: `Configuration for cluster level pod cidr overprovision. Default is disabled=false.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"additional_pod_ranges_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: `AdditionalPodRangesConfig is the configuration for additional pod secondary ranges supporting the ClusterUpdate message.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_range_names": {
										Type:        schema.TypeSet,
										MinItems:    1,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Name for pod secondary ipv4 range which has the actual range defined ahead.`,
									},
								},
							},
						},
					},
				},
			},

			"networking_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC_NATIVE", "ROUTES"}, false),
				Description:  `Determines whether alias IPs or routes will be used for pod IPs in the cluster.`,
			},

			"remove_default_node_pool": {
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   `If true, deletes the default node pool upon cluster creation. If you're using google_container_node_pool resources with no default node pool, this should be set to true, alongside setting initial_node_count to at least 1.`,
				ConflictsWith: []string{"enable_autopilot"},
			},

			"private_cluster_config": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
				Description:      `Configuration for private clusters, clusters with private nodes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// enable_private_endpoint is orthogonal to private_endpoint_subnetwork.
						// User can create a private_cluster_config block without including
						// either one of those two fields. Both fields are optional.
						// At the same time, we use 'AtLeastOneOf' to prevent an empty block
						// like 'private_cluster_config{}'
						"enable_private_endpoint": {
							Type:             schema.TypeBool,
							Optional:         true,
							AtLeastOneOf:     privateClusterConfigKeys,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
							Description:      `When true, the cluster's private endpoint is used as the cluster endpoint and access through the public endpoint is disabled. When false, either endpoint can be used.`,
						},
						"enable_private_nodes": {
							Type:             schema.TypeBool,
							Optional:         true,
							ForceNew:         true,
							AtLeastOneOf:     privateClusterConfigKeys,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
							Description:      `Enables the private cluster feature, creating a private endpoint on the cluster. In a private cluster, nodes only have RFC 1918 private addresses and communicate with the master's private endpoint via private networking.`,
						},
						"master_ipv4_cidr_block": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: privateClusterConfigKeys,
							ValidateFunc: verify.OrEmpty(validation.IsCIDRNetwork(28, 28)),
							Description:  `The IP range in CIDR notation to use for the hosted master network. This range will be used for assigning private IP addresses to the cluster master(s) and the ILB VIP. This range must not overlap with any other ranges in use within the cluster's network, and it must be a /28 subnet. See Private Cluster Limitations for more details. This field only applies to private clusters, when enable_private_nodes is true.`,
						},
						"peering_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the peering between this cluster and the Google owned VPC.`,
						},
						"private_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The internal IP address of this cluster's master endpoint.`,
						},
						"private_endpoint_subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							AtLeastOneOf:     privateClusterConfigKeys,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `Subnetwork in cluster's network where master's endpoint will be provisioned.`,
						},
						"public_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The external IP address of this cluster's master endpoint.`,
						},
						"master_global_access_config": {
							Type:         schema.TypeList,
							MaxItems:     1,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: privateClusterConfigKeys,
							Description:  "Controls cluster master global access settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether the cluster master is accessible globally or not.`,
									},
								},
							},
						},
					},
				},
			},

			"resource_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The GCE resource labels (a map of key/value pairs) to be applied to the cluster.`,
			},

			"label_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fingerprint of the set of labels for this cluster.`,
			},

			"default_max_pods_per_node": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Description:   `The default maximum number of pods per node in this cluster. This doesn't work on "routes-based" clusters, clusters that don't have IP Aliasing enabled.`,
				ConflictsWith: []string{"enable_autopilot"},
			},

			"vertical_pod_autoscaling": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Vertical Pod Autoscaling automatically adjusts the resources of pods controlled by it.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enables vertical pod autoscaling.`,
						},
					},
				},
			},
			"workload_identity_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// Computed is unsafe to remove- this API may return `"workloadIdentityConfig": {},` or omit the key entirely
				// and both will be valid. Note that we don't handle the case where the API returns nothing & the user has defined
				// workload_identity_config today.
				Computed:      true,
				Description:   `Configuration for the use of Kubernetes Service Accounts in GCP IAM policies.`,
				ConflictsWith: []string{"enable_autopilot"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workload_pool": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The workload pool to attach all Kubernetes service accounts to.",
						},
					},
				},
			},

			"identity_service_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Configuration for Identity Service which allows customers to use external identity providers with the K8S API.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable the Identity Service component.",
						},
					},
				},
			},

			"service_external_ips_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `If set, and enabled=true, services with external ips field will not be blocked`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `When enabled, services with external ips specified will be allowed.`,
						},
					},
				},
			},

			"mesh_certificates": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `If set, and enable_certificates=true, the GKE Workload Identity Certificates controller and node agent will be deployed in the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_certificates": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `When enabled the GKE Workload Identity Certificates controller and node agent will be deployed in the cluster.`,
						},
					},
				},
			},

			"database_encryption": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Application-layer Secrets Encryption settings. The object format is {state = string, key_name = string}. Valid values of state are: "ENCRYPTED"; "DECRYPTED". key_name is the name of a CloudKMS key.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ENCRYPTED", "DECRYPTED"}, false),
							Description:  `ENCRYPTED or DECRYPTED.`,
						},
						"key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The key to use to encrypt/decrypt secrets.`,
						},
					},
				},
			},

			"release_channel": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration options for the Release channel feature, which provide more control over automatic upgrades of your GKE clusters. Note that removing this field from your config will not unenroll it. Instead, use the "UNSPECIFIED" channel.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"UNSPECIFIED", "RAPID", "REGULAR", "STABLE"}, false),
							Description: `The selected release channel. Accepted values are:
* UNSPECIFIED: Not set.
* RAPID: Weekly upgrade cadence; Early testers and developers who requires new features.
* REGULAR: Multiple per month upgrade cadence; Production users who need features not yet offered in the Stable channel.
* STABLE: Every few months upgrade cadence; Production users who need stability above all else, and for whom frequent upgrades are too risky.`,
						},
					},
				},
			},

			"tpu_ipv4_cidr_block": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: `The IP address range of the Cloud TPUs in this cluster, in CIDR notation (e.g. 1.2.3.4/29).`,
			},

			"default_snat_status": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Whether the cluster disables default in-node sNAT rules. In-node sNAT rules will be disabled when defaultSnatStatus is disabled.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `When disabled is set to false, default IP masquerade rules will be applied to the nodes to prevent sNAT on cluster internal traffic.`,
						},
					},
				},
			},

			"datapath_provider": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				Description:      `The desired datapath provider for this cluster. By default, uses the IPTables-based kube-proxy implementation.`,
				ValidateFunc:     validation.StringInSlice([]string{"DATAPATH_PROVIDER_UNSPECIFIED", "LEGACY_DATAPATH", "ADVANCED_DATAPATH"}, false),
				DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("DATAPATH_PROVIDER_UNSPECIFIED"),
			},

			"enable_intranode_visibility": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				Description:   `Whether Intra-node visibility is enabled for this cluster. This makes same node pod to pod traffic visible for VPC network.`,
				ConflictsWith: []string{"enable_autopilot"},
			},
			"enable_l4_ilb_subsetting": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether L4ILB Subsetting is enabled for this cluster.`,
				Default:     false,
			},
			"private_ipv6_google_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The desired state of IPv6 connectivity to Google Services. By default, no private IPv6 access to or from Google Services (all access will be via IPv4).`,
				Computed:    true,
			},

			"cost_management_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Cost management configuration for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether to enable GKE cost allocation. When you enable GKE cost allocation, the cluster name and namespace of your GKE workloads appear in the labels field of the billing export to BigQuery. Defaults to false.`,
						},
					},
				},
			},

			"resource_usage_export_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: `Configuration for the ResourceUsageExportConfig feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_network_egress_metering": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: `Whether to enable network egress metering for this cluster. If enabled, a daemonset will be created in the cluster to meter network egress traffic.`,
						},
						"enable_resource_consumption_metering": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Whether to enable resource consumption metering on this cluster. When enabled, a table will be created in the resource export BigQuery dataset to store resource consumption data. The resulting table can be joined with the resource usage table or with BigQuery billing export. Defaults to true.`,
						},
						"bigquery_destination": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: `Parameters for using BigQuery as the destination of resource usage export.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataset_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The ID of a BigQuery Dataset.`,
									},
								},
							},
						},
					},
				},
			},
			"dns_config": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: suppressDiffForAutopilot,
				Description:      `Configuration for Cloud DNS for Kubernetes Engine.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_dns": {
							Type:         schema.TypeString,
							Default:      "PROVIDER_UNSPECIFIED",
							ValidateFunc: validation.StringInSlice([]string{"PROVIDER_UNSPECIFIED", "PLATFORM_DEFAULT", "CLOUD_DNS"}, false),
							Description:  `Which in-cluster DNS provider should be used.`,
							Optional:     true,
						},
						"cluster_dns_scope": {
							Type:         schema.TypeString,
							Default:      "DNS_SCOPE_UNSPECIFIED",
							ValidateFunc: validation.StringInSlice([]string{"DNS_SCOPE_UNSPECIFIED", "CLUSTER_SCOPE", "VPC_SCOPE"}, false),
							Description:  `The scope of access to cluster DNS records.`,
							Optional:     true,
						},
						"cluster_dns_domain": {
							Type:        schema.TypeString,
							Description: `The suffix used for all cluster service records.`,
							Optional:    true,
						},
					},
				},
			},
			"gateway_api_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration for GKE Gateway API controller.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CHANNEL_DISABLED", "CHANNEL_EXPERIMENTAL", "CHANNEL_STANDARD"}, false),
							Description:  `The Gateway API release channel to use for Gateway API.`,
						},
					},
				},
			},
		},
	}
}
