// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataproc

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func resourceDataprocClusterResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the cluster, unique within the project and zone.`,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) > 55 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 55 characters", k))
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

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the cluster will exist. If it is not provided, the provider project is used.`,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				ForceNew:    true,
				Description: `The region in which the cluster and associated nodes will be created in. Defaults to global.`,
			},

			"graceful_decommission_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0s",
				Description: `The timeout duration which allows graceful decomissioning when you change the number of worker nodes directly through a terraform apply`,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// GCP automatically adds labels
				DiffSuppressFunc: resourceDataprocLabelDiffSuppress,
				Computed:         true,
				Description:      `The list of labels (key/value pairs) to be applied to instances in the cluster. GCP generates some itself including goog-dataproc-cluster-name which is the name of the cluster.`,
			},

			"virtual_cluster_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The virtual cluster config is used when creating a Dataproc cluster that does not directly control the underlying compute resources, for example, when creating a Dataproc-on-GKE cluster. Dataproc may set default values, and values may change when clusters are updated. Exactly one of config or virtualClusterConfig must be specified.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"staging_bucket": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: virtualClusterConfigKeys,
							ForceNew:     true,
							Description:  `A Cloud Storage bucket used to stage job dependencies, config files, and job driver console output. If you do not specify a staging bucket, Cloud Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's staging bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket.`,
						},

						"auxiliary_services_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							MaxItems:     1,
							AtLeastOneOf: virtualClusterConfigKeys,
							Description:  `Auxiliary services configuration for a Cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"metastore_config": {
										Type:         schema.TypeList,
										Optional:     true,
										MaxItems:     1,
										AtLeastOneOf: auxiliaryServicesConfigKeys,
										Description:  `The Hive Metastore configuration for this workload.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"dataproc_metastore_service": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													AtLeastOneOf: auxiliaryServicesMetastoreConfigKeys,
													Description:  `The Hive Metastore configuration for this workload.`,
												},
											},
										},
									},

									"spark_history_server_config": {
										Type:         schema.TypeList,
										Optional:     true,
										MaxItems:     1,
										AtLeastOneOf: auxiliaryServicesConfigKeys,
										Description:  `The Spark History Server configuration for the workload.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"dataproc_cluster": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													AtLeastOneOf: auxiliaryServicesSparkHistoryServerConfigKeys,
													Description:  `Resource name of an existing Dataproc Cluster to act as a Spark History Server for the workload.`,
												},
											},
										},
									},
								},
							},
						},

						"kubernetes_cluster_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							MaxItems:     1,
							AtLeastOneOf: virtualClusterConfigKeys,
							Description:  `The configuration for running the Dataproc cluster on Kubernetes.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"kubernetes_namespace": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										AtLeastOneOf: kubernetesClusterConfigKeys,
										Description:  `A namespace within the Kubernetes cluster to deploy into. If this namespace does not exist, it is created. If it exists, Dataproc verifies that another Dataproc VirtualCluster is not installed into it. If not specified, the name of the Dataproc Cluster is used.`,
									},

									"kubernetes_software_config": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: `The software configuration for this Dataproc cluster running on Kubernetes.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"component_version": {
													Type:        schema.TypeMap,
													Required:    true,
													ForceNew:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: `The components that should be installed in this Dataproc cluster. The key must be a string from the KubernetesComponent enumeration. The value is the version of the software to be installed.`,
												},

												"properties": {
													Type:             schema.TypeMap,
													Optional:         true,
													ForceNew:         true,
													DiffSuppressFunc: resourceDataprocPropertyDiffSuppress,
													Elem:             &schema.Schema{Type: schema.TypeString},
													Computed:         true,
													Description:      `The properties to set on daemon config files. Property keys are specified in prefix:property format, for example spark:spark.kubernetes.container.image.`,
												},
											},
										},
									},

									"gke_cluster_config": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: `The configuration for running the Dataproc cluster on GKE.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												"gke_cluster_target": {
													Type:         schema.TypeString,
													ForceNew:     true,
													Optional:     true,
													AtLeastOneOf: gkeClusterConfigKeys,
													Description:  `A target GKE cluster to deploy to. It must be in the same project and region as the Dataproc cluster (the GKE cluster can be zonal or regional). Format: 'projects/{project}/locations/{location}/clusters/{cluster_id}'`,
												},

												"node_pool_target": {
													Type:         schema.TypeList,
													Optional:     true,
													AtLeastOneOf: gkeClusterConfigKeys,
													MinItems:     1,
													Description:  `GKE node pools where workloads will be scheduled. At least one node pool must be assigned the DEFAULT GkeNodePoolTarget.Role. If a GkeNodePoolTarget is not specified, Dataproc constructs a DEFAULT GkeNodePoolTarget.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{

															"node_pool": {
																Type:        schema.TypeString,
																ForceNew:    true,
																Required:    true,
																Description: `The target GKE node pool. Format: 'projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{nodePool}'`,
															},

															"roles": {
																Type:        schema.TypeSet,
																Elem:        &schema.Schema{Type: schema.TypeString},
																ForceNew:    true,
																Required:    true,
																Description: `The roles associated with the GKE node pool.`,
															},

															"node_pool_config": {
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																MaxItems:    1,
																Description: `Input only. The configuration for the GKE node pool.`,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{

																		"config": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Computed:    true,
																			MaxItems:    1,
																			Description: `The node pool configuration.`,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{

																					"machine_type": {
																						Type:        schema.TypeString,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The name of a Compute Engine machine type.`,
																					},

																					"local_ssd_count": {
																						Type:        schema.TypeInt,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The minimum number of nodes in the node pool. Must be >= 0 and <= maxNodeCount.`,
																					},

																					"preemptible": {
																						Type:        schema.TypeBool,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `Whether the nodes are created as preemptible VM instances. Preemptible nodes cannot be used in a node pool with the CONTROLLER role or in the DEFAULT node pool if the CONTROLLER role is not assigned (the DEFAULT node pool will assume the CONTROLLER role).`,
																					},

																					"min_cpu_platform": {
																						Type:        schema.TypeString,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `Minimum CPU platform to be used by this instance. The instance may be scheduled on the specified or a newer CPU platform. Specify the friendly names of CPU platforms, such as "Intel Haswell" or "Intel Sandy Bridge".`,
																					},

																					"spot": {
																						Type:        schema.TypeBool,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `Spot flag for enabling Spot VM, which is a rebrand of the existing preemptible flag.`,
																					},
																				},
																			},
																		},

																		"locations": {
																			Type:        schema.TypeSet,
																			Elem:        &schema.Schema{Type: schema.TypeString},
																			ForceNew:    true,
																			Required:    true,
																			Description: `The list of Compute Engine zones where node pool nodes associated with a Dataproc on GKE virtual cluster will be located.`,
																		},

																		"autoscaling": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Computed:    true,
																			MaxItems:    1,
																			Description: `The autoscaler configuration for this node pool. The autoscaler is enabled only when a valid configuration is present.`,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{

																					"min_node_count": {
																						Type:        schema.TypeInt,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The minimum number of nodes in the node pool. Must be >= 0 and <= maxNodeCount.`,
																					},

																					"max_node_count": {
																						Type:        schema.TypeInt,
																						ForceNew:    true,
																						Optional:    true,
																						Description: `The maximum number of nodes in the node pool. Must be >= minNodeCount, and must be > 0.`,
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
					},
				},
			},

			"cluster_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Allows you to configure various aspects of the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"staging_bucket": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							ForceNew:     true,
							Description:  `The Cloud Storage staging bucket used to stage files, such as Hadoop jars, between client machines and the cluster. Note: If you don't explicitly specify a staging_bucket then GCP will auto create / assign one for you. However, you are not guaranteed an auto generated bucket which is solely dedicated to your cluster; it may be shared with other clusters in the same region/zone also choosing to use the auto generation option.`,
						},
						// If the user does not specify a staging bucket, GCP will allocate one automatically.
						// The staging_bucket field provides a way for the user to supply their own
						// staging bucket. The bucket field is purely a computed field which details
						// the definitive bucket allocated and in use (either the user supplied one via
						// staging_bucket, or the GCP generated one)
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ` The name of the cloud storage bucket ultimately used to house the staging data for the cluster. If staging_bucket is specified, it will contain this value, otherwise it will be the auto generated name.`,
						},

						"temp_bucket": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: clusterConfigKeys,
							ForceNew:     true,
							Description:  `The Cloud Storage temp bucket used to store ephemeral cluster and jobs data, such as Spark and MapReduce history files. Note: If you don't explicitly specify a temp_bucket then GCP will auto create / assign one for you.`,
						},

						"gce_cluster_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `Common config settings for resources of Google Compute Engine cluster instances, applicable to all instances in the cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"zone": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Description:  `The GCP zone where your data is stored and used (i.e. where the master and the worker nodes will be created in). If region is set to 'global' (default) then zone is mandatory, otherwise GCP is able to make use of Auto Zone Placement to determine this automatically for you. Note: This setting additionally determines and restricts which computing resources are available for use with other configs such as cluster_config.master_config.machine_type and cluster_config.worker_config.machine_type.`,
									},

									"network": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										AtLeastOneOf:     gceClusterConfigKeys,
										ForceNew:         true,
										ConflictsWith:    []string{"cluster_config.0.gce_cluster_config.0.subnetwork"},
										DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
										Description:      `The name or self_link of the Google Compute Engine network to the cluster will be part of. Conflicts with subnetwork. If neither is specified, this defaults to the "default" network.`,
									},

									"subnetwork": {
										Type:             schema.TypeString,
										Optional:         true,
										AtLeastOneOf:     gceClusterConfigKeys,
										ForceNew:         true,
										ConflictsWith:    []string{"cluster_config.0.gce_cluster_config.0.network"},
										DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
										Description:      `The name or self_link of the Google Compute Engine subnetwork the cluster will be part of. Conflicts with network.`,
									},

									"tags": {
										Type:         schema.TypeSet,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Elem:         &schema.Schema{Type: schema.TypeString},
										Description:  `The list of instance tags applied to instances in the cluster. Tags are used to identify valid sources or targets for network firewalls.`,
									},

									"service_account": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Description:  `The service account to be used by the Node VMs. If not specified, the "default" service account is used.`,
									},

									"service_account_scopes": {
										Type:         schema.TypeSet,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Description:  `The set of Google API scopes to be made available on all of the node VMs under the service_account specified. These can be either FQDNs, or scope aliases.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											StateFunc: func(v interface{}) string {
												return tpgresource.CanonicalizeServiceScope(v.(string))
											},
										},
										Set: tpgresource.StringScopeHashcode,
									},

									"internal_ip_only": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										ForceNew:     true,
										Default:      false,
										Description:  `By default, clusters are not restricted to internal IP addresses, and will have ephemeral external IP addresses assigned to each instance. If set to true, all instances in the cluster will only have internal IP addresses. Note: Private Google Access (also known as privateIpGoogleAccess) must be enabled on the subnetwork that the cluster will be launched in.`,
									},

									"metadata": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ForceNew:     true,
										Description:  `A map of the Compute Engine metadata entries to add to all instances`,
									},

									"shielded_instance_config": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Computed:     true,
										MaxItems:     1,
										Description:  `Shielded Instance Config for clusters using Compute Engine Shielded VMs.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_secure_boot": {
													Type:         schema.TypeBool,
													Optional:     true,
													Default:      false,
													AtLeastOneOf: schieldedInstanceConfigKeys,
													ForceNew:     true,
													Description:  `Defines whether instances have Secure Boot enabled.`,
												},
												"enable_vtpm": {
													Type:         schema.TypeBool,
													Optional:     true,
													Default:      false,
													AtLeastOneOf: schieldedInstanceConfigKeys,
													ForceNew:     true,
													Description:  `Defines whether instances have the vTPM enabled.`,
												},
												"enable_integrity_monitoring": {
													Type:         schema.TypeBool,
													Optional:     true,
													Default:      false,
													AtLeastOneOf: schieldedInstanceConfigKeys,
													ForceNew:     true,
													Description:  `Defines whether instances have integrity monitoring enabled.`,
												},
											},
										},
									},

									"reservation_affinity": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Computed:     true,
										MaxItems:     1,
										Description:  `Reservation Affinity for consuming Zonal reservation.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"consume_reservation_type": {
													Type:         schema.TypeString,
													Optional:     true,
													AtLeastOneOf: reservationAffinityKeys,
													ForceNew:     true,
													ValidateFunc: validation.StringInSlice([]string{"NO_RESERVATION", "ANY_RESERVATION", "SPECIFIC_RESERVATION"}, false),
													Description:  `Type of reservation to consume.`,
												},
												"key": {
													Type:         schema.TypeString,
													Optional:     true,
													AtLeastOneOf: reservationAffinityKeys,
													ForceNew:     true,
													Description:  `Corresponds to the label key of reservation resource.`,
												},
												"values": {
													Type:         schema.TypeSet,
													Elem:         &schema.Schema{Type: schema.TypeString},
													Optional:     true,
													AtLeastOneOf: reservationAffinityKeys,
													ForceNew:     true,
													Description:  `Corresponds to the label values of reservation resource.`,
												},
											},
										},
									},

									"node_group_affinity": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: gceClusterConfigKeys,
										Computed:     true,
										MaxItems:     1,
										Description:  `Node Group Affinity for sole-tenant clusters.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"node_group_uri": {
													Type:             schema.TypeString,
													ForceNew:         true,
													Required:         true,
													Description:      `The URI of a sole-tenant that the cluster will be created on.`,
													DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
												},
											},
										},
									},
								},
							},
						},

						"master_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `The Compute Engine config settings for the cluster's master instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_instances": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Computed:    true,
										Description: `Specifies the number of master nodes to create. If not specified, GCP will default to a predetermined computed value.`,
										AtLeastOneOf: []string{
											"cluster_config.0.master_config.0.num_instances",
											"cluster_config.0.master_config.0.image_uri",
											"cluster_config.0.master_config.0.machine_type",
											"cluster_config.0.master_config.0.accelerators",
										},
									},

									"image_uri": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.master_config.0.num_instances",
											"cluster_config.0.master_config.0.image_uri",
											"cluster_config.0.master_config.0.machine_type",
											"cluster_config.0.master_config.0.accelerators",
										},
										ForceNew:    true,
										Description: `The URI for the image to use for this master`,
									},

									"machine_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.master_config.0.num_instances",
											"cluster_config.0.master_config.0.image_uri",
											"cluster_config.0.master_config.0.machine_type",
											"cluster_config.0.master_config.0.accelerators",
										},
										ForceNew:    true,
										Description: `The name of a Google Compute Engine machine type to create for the master`,
									},

									"min_cpu_platform": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.master_config.0.num_instances",
											"cluster_config.0.master_config.0.image_uri",
											"cluster_config.0.master_config.0.machine_type",
											"cluster_config.0.master_config.0.accelerators",
										},
										ForceNew:    true,
										Description: `The name of a minimum generation of CPU family for the master. If not specified, GCP will default to a predetermined computed value for each zone.`,
									},
									"disk_config": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.master_config.0.num_instances",
											"cluster_config.0.master_config.0.image_uri",
											"cluster_config.0.master_config.0.machine_type",
											"cluster_config.0.master_config.0.accelerators",
										},
										MaxItems:    1,
										Description: `Disk Config`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_local_ssds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: `The amount of local SSD disks that will be attached to each master cluster node. Defaults to 0.`,
													AtLeastOneOf: []string{
														"cluster_config.0.master_config.0.disk_config.0.num_local_ssds",
														"cluster_config.0.master_config.0.disk_config.0.boot_disk_size_gb",
														"cluster_config.0.master_config.0.disk_config.0.boot_disk_type",
													},
													ForceNew: true,
												},

												"boot_disk_size_gb": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: `Size of the primary disk attached to each node, specified in GB. The primary disk contains the boot volume and system libraries, and the smallest allowed disk size is 10GB. GCP will default to a predetermined computed value if not set (currently 500GB). Note: If SSDs are not attached, it also contains the HDFS data blocks and Hadoop working directories.`,
													AtLeastOneOf: []string{
														"cluster_config.0.master_config.0.disk_config.0.num_local_ssds",
														"cluster_config.0.master_config.0.disk_config.0.boot_disk_size_gb",
														"cluster_config.0.master_config.0.disk_config.0.boot_disk_type",
													},
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(10),
												},

												"boot_disk_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The disk type of the primary disk attached to each node. Such as "pd-ssd" or "pd-standard". Defaults to "pd-standard".`,
													AtLeastOneOf: []string{
														"cluster_config.0.master_config.0.disk_config.0.num_local_ssds",
														"cluster_config.0.master_config.0.disk_config.0.boot_disk_size_gb",
														"cluster_config.0.master_config.0.disk_config.0.boot_disk_type",
													},
													ForceNew: true,
													Default:  "pd-standard",
												},
											},
										},
									},
									"accelerators": {
										Type:     schema.TypeSet,
										Optional: true,
										AtLeastOneOf: []string{
											"cluster_config.0.master_config.0.num_instances",
											"cluster_config.0.master_config.0.image_uri",
											"cluster_config.0.master_config.0.machine_type",
											"cluster_config.0.master_config.0.accelerators",
										},
										ForceNew:    true,
										Elem:        acceleratorsSchema(),
										Description: `The Compute Engine accelerator (GPU) configuration for these instances. Can be specified multiple times.`,
									},

									"instance_names": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `List of master instance names which have been assigned to the cluster.`,
									},
								},
							},
						},
						"worker_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `The Compute Engine config settings for the cluster's worker instances.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_instances": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    false,
										Computed:    true,
										Description: `Specifies the number of worker nodes to create. If not specified, GCP will default to a predetermined computed value.`,
										AtLeastOneOf: []string{
											"cluster_config.0.worker_config.0.num_instances",
											"cluster_config.0.worker_config.0.image_uri",
											"cluster_config.0.worker_config.0.machine_type",
											"cluster_config.0.worker_config.0.accelerators",
											"cluster_config.0.worker_config.0.min_num_instances",
										},
									},

									"image_uri": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.worker_config.0.num_instances",
											"cluster_config.0.worker_config.0.image_uri",
											"cluster_config.0.worker_config.0.machine_type",
											"cluster_config.0.worker_config.0.accelerators",
											"cluster_config.0.worker_config.0.min_num_instances",
										},
										ForceNew:    true,
										Description: `The URI for the image to use for this master/worker`,
									},

									"machine_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.worker_config.0.num_instances",
											"cluster_config.0.worker_config.0.image_uri",
											"cluster_config.0.worker_config.0.machine_type",
											"cluster_config.0.worker_config.0.accelerators",
											"cluster_config.0.worker_config.0.min_num_instances",
										},
										ForceNew:    true,
										Description: `The name of a Google Compute Engine machine type to create for the master/worker`,
									},

									"min_cpu_platform": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.worker_config.0.num_instances",
											"cluster_config.0.worker_config.0.image_uri",
											"cluster_config.0.worker_config.0.machine_type",
											"cluster_config.0.worker_config.0.accelerators",
											"cluster_config.0.worker_config.0.min_num_instances",
										},
										ForceNew:    true,
										Description: `The name of a minimum generation of CPU family for the master/worker. If not specified, GCP will default to a predetermined computed value for each zone.`,
									},
									"disk_config": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"cluster_config.0.worker_config.0.num_instances",
											"cluster_config.0.worker_config.0.image_uri",
											"cluster_config.0.worker_config.0.machine_type",
											"cluster_config.0.worker_config.0.accelerators",
											"cluster_config.0.worker_config.0.min_num_instances",
											"cluster_config.0.worker_config.0.disk_config",
										},
										MaxItems:    1,
										Description: `Disk Config`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_local_ssds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: `The amount of local SSD disks that will be attached to each master cluster node. Defaults to 0.`,
													AtLeastOneOf: []string{
														"cluster_config.0.worker_config.0.disk_config.0.num_local_ssds",
														"cluster_config.0.worker_config.0.disk_config.0.boot_disk_size_gb",
														"cluster_config.0.worker_config.0.disk_config.0.boot_disk_type",
													},
													ForceNew: true,
												},

												"boot_disk_size_gb": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: `Size of the primary disk attached to each node, specified in GB. The primary disk contains the boot volume and system libraries, and the smallest allowed disk size is 10GB. GCP will default to a predetermined computed value if not set (currently 500GB). Note: If SSDs are not attached, it also contains the HDFS data blocks and Hadoop working directories.`,
													AtLeastOneOf: []string{
														"cluster_config.0.worker_config.0.disk_config.0.num_local_ssds",
														"cluster_config.0.worker_config.0.disk_config.0.boot_disk_size_gb",
														"cluster_config.0.worker_config.0.disk_config.0.boot_disk_type",
													},
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(10),
												},

												"boot_disk_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The disk type of the primary disk attached to each node. Such as "pd-ssd" or "pd-standard". Defaults to "pd-standard".`,
													AtLeastOneOf: []string{
														"cluster_config.0.worker_config.0.disk_config.0.num_local_ssds",
														"cluster_config.0.worker_config.0.disk_config.0.boot_disk_size_gb",
														"cluster_config.0.worker_config.0.disk_config.0.boot_disk_type",
													},
													ForceNew: true,
													Default:  "pd-standard",
												},
											},
										},
									},

									// Note: preemptible workers don't support accelerators
									"accelerators": {
										Type:     schema.TypeSet,
										Optional: true,
										AtLeastOneOf: []string{
											"cluster_config.0.worker_config.0.num_instances",
											"cluster_config.0.worker_config.0.image_uri",
											"cluster_config.0.worker_config.0.machine_type",
											"cluster_config.0.worker_config.0.accelerators",
											"cluster_config.0.worker_config.0.min_num_instances",
										},
										ForceNew:    true,
										Elem:        acceleratorsSchema(),
										Description: `The Compute Engine accelerator (GPU) configuration for these instances. Can be specified multiple times.`,
									},

									"instance_names": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `List of master/worker instance names which have been assigned to the cluster.`,
									},
									"min_num_instances": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: false,
										AtLeastOneOf: []string{
											"cluster_config.0.worker_config.0.num_instances",
											"cluster_config.0.worker_config.0.image_uri",
											"cluster_config.0.worker_config.0.machine_type",
											"cluster_config.0.worker_config.0.accelerators",
											"cluster_config.0.worker_config.0.min_num_instances",
										},
										Description: `The minimum number of primary worker instances to create.`,
									},
								},
							},
						},
						// preemptible_worker_config has a slightly different config
						"preemptible_worker_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `The Google Compute Engine config settings for the additional (aka preemptible) instances in a cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_instances": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the number of preemptible nodes to create. Defaults to 0.`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
											"cluster_config.0.preemptible_worker_config.0.preemptibility",
											"cluster_config.0.preemptible_worker_config.0.disk_config",
										},
									},

									// API does not honour this if set ...
									// It always uses whatever is specified for the worker_config
									// "machine_type": { ... }
									// "min_cpu_platform": { ... }
									"preemptibility": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the preemptibility of the secondary nodes. Defaults to PREEMPTIBLE.`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
											"cluster_config.0.preemptible_worker_config.0.preemptibility",
											"cluster_config.0.preemptible_worker_config.0.disk_config",
										},
										ForceNew:     true,
										ValidateFunc: validation.StringInSlice([]string{"PREEMPTIBILITY_UNSPECIFIED", "NON_PREEMPTIBLE", "PREEMPTIBLE", "SPOT"}, false),
										Default:      "PREEMPTIBLE",
									},

									"disk_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: `Disk Config`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
											"cluster_config.0.preemptible_worker_config.0.preemptibility",
											"cluster_config.0.preemptible_worker_config.0.disk_config",
										},
										MaxItems: 1,

										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_local_ssds": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													AtLeastOneOf: preemptibleWorkerDiskConfigKeys,
													ForceNew:     true,
													Description:  `The amount of local SSD disks that will be attached to each preemptible worker node. Defaults to 0.`,
												},

												"boot_disk_size_gb": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													AtLeastOneOf: preemptibleWorkerDiskConfigKeys,
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(10),
													Description:  `Size of the primary disk attached to each preemptible worker node, specified in GB. The smallest allowed disk size is 10GB. GCP will default to a predetermined computed value if not set (currently 500GB). Note: If SSDs are not attached, it also contains the HDFS data blocks and Hadoop working directories.`,
												},

												"boot_disk_type": {
													Type:         schema.TypeString,
													Optional:     true,
													AtLeastOneOf: preemptibleWorkerDiskConfigKeys,
													ForceNew:     true,
													Default:      "pd-standard",
													Description:  `The disk type of the primary disk attached to each preemptible worker node. Such as "pd-ssd" or "pd-standard". Defaults to "pd-standard".`,
												},
											},
										},
									},

									"instance_names": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `List of preemptible instance names which have been assigned to the cluster.`,
									},
									"instance_flexibility_policy": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: `Instance flexibility Policy allowing a mixture of VM shapes and provisioning models.`,
										AtLeastOneOf: []string{
											"cluster_config.0.preemptible_worker_config.0.num_instances",
											"cluster_config.0.preemptible_worker_config.0.preemptibility",
											"cluster_config.0.preemptible_worker_config.0.disk_config",
											"cluster_config.0.preemptible_worker_config.0.instance_flexibility_policy",
										},
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_selection_list": {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													ForceNew: true,
													AtLeastOneOf: []string{
														"cluster_config.0.preemptible_worker_config.0.instance_flexibility_policy.0.instance_selection_list",
													},
													Description: `List of instance selection options that the group will use when creating new VMs.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"machine_types": {
																Type:        schema.TypeList,
																Computed:    true,
																Optional:    true,
																ForceNew:    true,
																Elem:        &schema.Schema{Type: schema.TypeString},
																Description: `Full machine-type names, e.g. "n1-standard-16".`,
															},
															"rank": {
																Type:        schema.TypeInt,
																Computed:    true,
																Optional:    true,
																ForceNew:    true,
																Elem:        &schema.Schema{Type: schema.TypeInt},
																Description: `Preference of this instance selection. Lower number means higher preference. Dataproc will first try to create a VM based on the machine-type with priority rank and fallback to next rank based on availability. Machine types and instance selections with the same priority have the same preference.`,
															},
														},
													},
												},
												"instance_selection_results": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: `A list of instance selection results in the group.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"machine_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Elem:        &schema.Schema{Type: schema.TypeString},
																Description: `Full machine-type names, e.g. "n1-standard-16".`,
															},
															"vm_count": {
																Type:        schema.TypeInt,
																Computed:    true,
																Elem:        &schema.Schema{Type: schema.TypeInt},
																Description: `Number of VM provisioned with the machine_type.`,
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

						"security_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Security related configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kerberos_config": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Kerberos related configuration",
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cross_realm_trust_admin_server": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The admin server (IP or hostname) for the remote trusted realm in a cross realm trust relationship.`,
												},
												"cross_realm_trust_kdc": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The KDC (IP or hostname) for the remote trusted realm in a cross realm trust relationship.`,
												},
												"cross_realm_trust_realm": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The remote realm the Dataproc on-cluster KDC will trust, should the user enable cross realm trust.`,
												},
												"cross_realm_trust_shared_password_uri": {
													Type:     schema.TypeString,
													Optional: true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the shared password between the on-cluster
Kerberos realm and the remote trusted realm, in a cross realm trust relationship.`,
												},
												"enable_kerberos": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: `Flag to indicate whether to Kerberize the cluster.`,
												},
												"kdc_db_key_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the master key of the KDC database.`,
												},
												"key_password_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the password to the user provided key. For the self-signed certificate, this password is generated by Dataproc.`,
												},
												"keystore_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of the keystore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.`,
												},
												"keystore_password_uri": {
													Type:     schema.TypeString,
													Optional: true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing
the password to the user provided keystore. For the self-signed certificate, this password is generated
by Dataproc`,
												},
												"kms_key_uri": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The uri of the KMS key used to encrypt various sensitive files.`,
												},
												"realm": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The name of the on-cluster Kerberos realm. If not specified, the uppercased domain of hostnames will be the realm.`,
												},
												"root_principal_password_uri": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The cloud Storage URI of a KMS encrypted file containing the root principal password.`,
												},
												"tgt_lifetime_hours": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: `The lifetime of the ticket granting ticket, in hours.`,
												},
												"truststore_password_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of a KMS encrypted file containing the password to the user provided truststore. For the self-signed certificate, this password is generated by Dataproc.`,
												},
												"truststore_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The Cloud Storage URI of the truststore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.`,
												},
											},
										},
									},
								},
							},
						},

						"software_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							Computed:     true,
							MaxItems:     1,
							Description:  `The config settings for software inside the cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_version": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										AtLeastOneOf:     clusterSoftwareConfigKeys,
										ForceNew:         true,
										DiffSuppressFunc: dataprocImageVersionDiffSuppress,
										Description:      `The Cloud Dataproc image version to use for the cluster - this controls the sets of software versions installed onto the nodes when you create clusters. If not specified, defaults to the latest version.`,
									},
									"override_properties": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: clusterSoftwareConfigKeys,
										ForceNew:     true,
										Elem:         &schema.Schema{Type: schema.TypeString},
										Description:  `A list of override and additional properties (key/value pairs) used to modify various aspects of the common configuration files used when creating a cluster.`,
									},

									"properties": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: `A list of the properties used to set the daemon config files. This will include any values supplied by the user via cluster_config.software_config.override_properties`,
									},

									// We have two versions of the properties field here because by default
									// dataproc will set a number of default properties for you out of the
									// box. If you want to override one or more, if we only had one field,
									// you would need to add in all these values as well otherwise you would
									// get a diff. To make this easier, 'properties' simply contains the computed
									// values (including overrides) for all properties, whilst override_properties
									// is only for properties the user specifically wants to override. If nothing
									// is overridden, this will be empty.

									"optional_components": {
										Type:         schema.TypeSet,
										Optional:     true,
										AtLeastOneOf: clusterSoftwareConfigKeys,
										Description:  `The set of optional components to activate on the cluster.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"initialization_action": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							ForceNew:     true,
							Description:  `Commands to execute on each node after config is completed. You can specify multiple versions of these.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"script": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `The script to be executed during initialization of the cluster. The script must be a GCS file with a gs:// prefix.`,
									},

									"timeout_sec": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     300,
										ForceNew:    true,
										Description: `The maximum duration (in seconds) which script is allowed to take to execute its action. GCP will default to a predetermined computed value if not set (currently 300).`,
									},
								},
							},
						},
						"encryption_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							MaxItems:     1,
							Description:  `The Customer managed encryption keys settings for the cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The Cloud KMS key name to use for PD disk encryption for all instances in the cluster.`,
									},
								},
							},
						},
						"autoscaling_config": {
							Type:             schema.TypeList,
							Optional:         true,
							AtLeastOneOf:     clusterConfigKeys,
							MaxItems:         1,
							Description:      `The autoscaling policy config associated with the cluster.`,
							DiffSuppressFunc: tpgresource.EmptyOrUnsetBlockDiffSuppress,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_uri": {
										Type:             schema.TypeString,
										Required:         true,
										Description:      `The autoscaling policy used by the cluster.`,
										DiffSuppressFunc: tpgresource.LocationDiffSuppress,
									},
								},
							},
						},
						"metastore_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: clusterConfigKeys,
							MaxItems:     1,
							Description:  `Specifies a Metastore configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataproc_metastore_service": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `Resource name of an existing Dataproc Metastore service.`,
									},
								},
							},
						},
						"lifecycle_config": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							AtLeastOneOf: clusterConfigKeys,
							Description:  `The settings for auto deletion cluster schedule.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"idle_delete_ttl": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The duration to keep the cluster alive while idling (no jobs running). After this TTL, the cluster will be deleted. Valid range: [10m, 14d].`,
										AtLeastOneOf: []string{
											"cluster_config.0.lifecycle_config.0.idle_delete_ttl",
											"cluster_config.0.lifecycle_config.0.auto_delete_time",
										},
									},
									"idle_start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Time when the cluster became idle (most recent job finished) and became eligible for deletion due to idleness.`,
									},
									// the API also has the auto_delete_ttl option in its request, however,
									// the value is not returned in the response, rather the auto_delete_time
									// after calculating ttl with the update time is returned, thus, for now
									// we will only allow auto_delete_time to updated.
									"auto_delete_time": {
										Type:             schema.TypeString,
										Optional:         true,
										Description:      `The time when cluster will be auto-deleted. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".`,
										DiffSuppressFunc: tpgresource.TimestampDiffSuppress(time.RFC3339Nano),
										AtLeastOneOf: []string{
											"cluster_config.0.lifecycle_config.0.idle_delete_ttl",
											"cluster_config.0.lifecycle_config.0.auto_delete_time",
										},
									},
								},
							},
						},
						"endpoint_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							MaxItems:     1,
							Description:  `The config settings for port access on the cluster. Structure defined below.`,
							AtLeastOneOf: clusterConfigKeys,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_http_port_access": {
										Type:        schema.TypeBool,
										Required:    true,
										ForceNew:    true,
										Description: `The flag to enable http access to specific ports on the cluster from external sources (aka Component Gateway). Defaults to false.`,
									},
									"http_ports": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: `The map of port descriptions to URLs. Will only be populated if enable_http_port_access is true.`,
									},
								},
							},
						},

						"dataproc_metric_config": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Description:  `The config for Dataproc metrics.`,
							AtLeastOneOf: clusterConfigKeys,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metrics": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `Metrics sources to enable.`,
										Elem:        metricsSchema(),
									},
								},
							},
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func ResourceDataprocClusterStateUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return tpgresource.LabelsStateUpgrade(rawState, resourceDataprocGoogleLabelPrefix)
}
